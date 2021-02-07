package exif

import (
	"encoding/binary"
	"errors"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// Errors
var (
	// ErrDataLength is an error for data length
	ErrDataLength = errors.New("error the data is not long enough")

	// ErrIfdBufferLength
	ErrIfdBufferLength = errors.New("ifd buffer length insufficient")
)

type ifdTagEnumerator struct {
	exifReader *reader
	byteOrder  binary.ByteOrder
	ifdOffset  uint32
	offset     uint32
	rawBuffer  [4]byte // Length of uint32
}

// scan moves through an ifd at the specified offset and enumerates over the IfdTags
func scan(er *reader, e *ExifData, ifd ifds.IFD, offset uint32) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	for ifdIndex := 0; ; ifdIndex++ {
		enumerator := getTagEnumerator(offset, er)
		//fmt.Printf("Parsing IFD [%s] (%d) at offset (0x%04x).\n", ifd, ifdIndex, offset)
		nextIfdOffset, err := enumerator.ParseIfd(e, ifd, ifdIndex, true)
		if err != nil {
			return err
		}
		if nextIfdOffset == 0 {
			break
		}

		offset = nextIfdOffset
	}
	return
}

// scanSubIfds moves through the subIfds at the specified offsetes and enumerates over their IfdTags
func scanSubIfds(er *reader, e *ExifData, t tag.Tag) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// Fetch SubIfd Values from []Uint32 (LongType)
	offsets, err := t.Uint32Values(er)
	if err != nil {
		return err
	}

	for ifdIndex := 0; ifdIndex < len(offsets); ifdIndex++ {
		enumerator := getTagEnumerator(offsets[ifdIndex], er)
		//fmt.Printf("Parsing IFD [%s] (%d) at offset (0x%04x).\n", ifd, ifdIndex, offset)
		if _, err := enumerator.ParseIfd(e, ifds.SubIFD, ifdIndex, false); err != nil {
			// Log Error
			continue
		}
	}

	return
}

// ifdTagEnumerator implements the io.Reader interface using
// an underlying exifReader, ifdOffset and offset
func (ite *ifdTagEnumerator) Read(p []byte) (n int, err error) {

	// Read from underlying exifReader io.ReaderAt interface
	n, err = ite.exifReader.ReadAt(p, int64(ite.offset+ite.ifdOffset))

	// Update reader offset
	ite.offset += uint32(n)

	return
}

func getTagEnumerator(offset uint32, er *reader) *ifdTagEnumerator {
	return &ifdTagEnumerator{
		exifReader: er,
		byteOrder:  er.ByteOrder(),
		ifdOffset:  offset,
	}
}

// rawValueOffset safely copies the ifdTagEnumerator's raw Buffer
func (ite *ifdTagEnumerator) rawValueOffset() (rawValueOffset tag.RawValueOffset, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	if n := copy(rawValueOffset[:], ite.rawBuffer[:]); n < 4 {
		panic(ErrIfdBufferLength)
	}
	return
}

// parseUndefinedIfds
// Makernotes and AdobeDNGData
func (ite *ifdTagEnumerator) parseUndefinedIfds(e *ExifData, ifd ifds.IFD) bool {
	if ifd == ifds.MknoteIFD {
		switch e.make {
		case "Canon":
			// Canon Makernotes do not have a Makernote Header
			// offset 0
			// ByteOrder is the same as RootIfd
			return true
		case "NIKON CORPORATION", "Nikon":
			// Nikon v3 maker note is a self-contained Ifd
			// (offsets are relative to the start of the maker note)
			byteOrder, err := mknote.NikonMkNoteHeader(ite)
			if err != nil {
				return false
			}
			ite.byteOrder = byteOrder
			return true
		}
		return false
	}

	// TODO: Adobe DNG data
	return true
}

// ParseIfd - enumerates over the ifd using the enumerator.ifdReader
func (ite *ifdTagEnumerator) ParseIfd(e *ExifData, ifd ifds.IFD, ifdIndex int, doDescend bool) (nextIfdOffset uint32, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// Parse undefined Ifds
	if !ite.parseUndefinedIfds(e, ifd) {
		return 0, nil
	}

	// Determine tagCount
	tagCount, err := ite.uint16()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Parsing \"%s\" with %d tags at offset [0x%04x]\n", ifd.String(), tagCount, ite.ifdOffset)

	// Log info
	// Remove log for now until we have a better solution
	//log.Info().
	//	Str("ifd", ifd.String()).
	//	Uint32("offset", ite.ifdOffset).
	//	Uint8("ifdIndex", uint8(ifdIndex)).
	//	Uint16("tagcount", tagCount).
	//	Msg("Parsing IFD")

	if tagCount > 256 {
		panic(errors.New("error Tagcount too high"))
	}

	// Add Ifd to Exif
	e.AddIfd(ifd)

	for i := 0; i < int(tagCount); i++ {
		t, err := ite.ParseTag()
		if err != nil {
			if err == tag.ErrTagTypeNotValid {
				//if errors.Is(err, tag.ErrTagTypeNotValid) {
				// Log TagNotValid Error
				//ifdEnumerateLogger.Warningf(nil, "Tag in IFD [%s] at position (%d) has invalid type and will be skipped.", fqIfdPath, i)
				continue
			}
			if err != nil {
				panic(err)
			}
		}

		// Descend to Child IFD
		childIFD := ifd.IsChildIfd(t)
		switch childIFD {
		case ifds.NullIFD:
			e.AddTag(ifd, ifdIndex, t)
		case ifds.SubIFD:
			if err := scanSubIfds(ite.exifReader, e, t); err != nil {
				panic(err)
			}
		default:
			if err := scan(ite.exifReader, e, childIFD, t.Offset()); err != nil {
				panic(err)
			}
		}

		//fmt.Printf("%s %s \t %v\n", ifd, ifd.TagName(t.TagID), t.TagType)
	}

	// NextIfdOffset
	if nextIfdOffset, err = ite.uint32(); err != nil {
		panic(err)
	}

	// Adjust for incorrect Makernotes NextIfd Offsets
	// set nextIfdOffset to 0x0000.
	if ifd == ifds.MknoteIFD {
		nextIfdOffset = 0x0000
	}

	return
}

// ParseTag parses the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an ifdTagEnumerator
func (ite *ifdTagEnumerator) ParseTag() (t tag.Tag, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// TagID
	tagIDRaw, err := ite.uint16()
	if err != nil {
		panic(err)
	}

	// TagType
	tagTypeRaw, err := ite.uint16()
	if err != nil {
		panic(err)
	}

	// UnitCount
	unitCount, err := ite.uint32()
	if err != nil {
		panic(err)
	}

	// ValueOffset
	valueOffset, err := ite.uint32()
	if err != nil {
		panic(err)
	}

	// RawBytes for ValueOffset
	rawValueOffset, err := ite.rawValueOffset()
	if err != nil {
		panic(err)
	}

	// Creates a newTag. If the TypeFromRaw is unsupported, it panics.
	t = tag.NewTag(tag.ID(tagIDRaw), tag.TypeFromRaw(tagTypeRaw), unitCount, valueOffset, rawValueOffset)

	return
}

// uint16 reads a uint16 and advances both our current and our current
// accumulator (which allows us to know how far to seek to the beginning of the
// next IFD when it's time to jump).
func (ite *ifdTagEnumerator) uint16() (value uint16, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	if n, err := ite.Read(ite.rawBuffer[:2]); err != nil || n != 2 { // Uint16 = 2bytes
		panic(err)
	}

	value = ite.byteOrder.Uint16(ite.rawBuffer[:2])

	return value, nil
}

// uint32 reads a uint32 and advances both our current and our current
// accumulator (which allows us to know how far to seek to the beginning of the
// next IFD when it's time to jump).
func (ite *ifdTagEnumerator) uint32() (value uint32, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	if n, err := ite.Read(ite.rawBuffer[:]); err != nil || n != 4 { // Uint32 = 4bytes
		panic(err)
	}

	value = ite.byteOrder.Uint32(ite.rawBuffer[:])

	return value, nil
}
