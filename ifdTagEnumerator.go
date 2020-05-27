package exiftool

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/evanoberholster/exiftool/exif"
	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/ifds/mknote"
	"github.com/evanoberholster/exiftool/tag"
)

// Errors
var (
	ErrIfdBufferLength = fmt.Errorf("Ifd buffer length insufficient")
)

type ifdTagEnumerator struct {
	exifReader *ExifReader
	byteOrder  binary.ByteOrder
	ifdOffset  uint32
	offset     uint32
	rawBuffer  [4]byte // Length of uint32
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

func getTagEnumerator(offset uint32, er *ExifReader) *ifdTagEnumerator {
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
func (ite *ifdTagEnumerator) parseUndefinedIfds(e *exif.Exif, ifd ifds.IFD) bool {
	if ifd == ifds.MknoteIFD {
		switch e.Make {
		case "Canon":
			// Canon Makernotes do not have a Makernote Header
			// offset 0
			// ByteOrder is the same as RootIfd
			return true
		case "NIKON CORPORATION", "Nikon":
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
func (ite *ifdTagEnumerator) ParseIfd(e *exif.Exif, ifd ifds.IFD, ifdIndex int, doDescend bool) (nextIfdOffset uint32, err error) {
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

	if tagCount > 256 {
		panic(errors.New("error Tagcount too high"))
	}

	// Add Ifd to Exif
	e.AddIfd(ifd)

	for i := 0; i < int(tagCount); i++ {
		t, err := ite.ParseTag()
		if err != nil {
			if err == ErrTagTypeNotValid {
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
			if err := ite.exifReader.scanSubIfds(e, t); err != nil {
				panic(err)
			}
		default:
			if err := ite.exifReader.scan(e, childIFD, t.Offset()); err != nil {
				panic(err)
			}
		}

		// Add Make and Model to Exif struct for future decoding of Makernotes
		switch t.TagID {
		case ifds.Make:
			e.Make, _ = t.ASCIIValue(ite.exifReader)
		case ifds.Model:
			e.Model, _ = t.ASCIIValue(ite.exifReader)
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
