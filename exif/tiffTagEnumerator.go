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
}

// scan moves through an ifd at the specified offset and enumerates over the IfdTags
func scan(er *reader, e *Data, ifd ifds.IFD, offset uint32) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	var ifdIndex uint8
	for ifdIndex = 0; ; ifdIndex++ {
		er.ifdExifOffset[ifd] = uint32(er.exifOffset)

		enumerator := newTagEnumerator(offset, er)
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
func scanSubIfds(er *reader, e *Data, t tag.Tag) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// Fetch SubIfd Values from []Uint32 (LongType)
	offsets, err := e.ParseUint32Values(t)
	if err != nil {
		return err
	}
	var ifdIndex uint8
	for ifdIndex = 0; ifdIndex < uint8(len(offsets)); ifdIndex++ {
		enumerator := newTagEnumerator(offsets[ifdIndex], er)
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

	ite.offset += uint32(n) // Update reader offset
	return
}

// ReadBuffer reads the buffer from the underlying exifReader and advances our ifdTagEnumerator offset
// (which allows us to know how far to seek to the beginning of the next IFD when it's time to jump).
// Returns the temporary buffer or an error. Buffer is only valid until the next read.
func (ite *ifdTagEnumerator) ReadBuffer(n int) (buf []byte, err error) {
	if n > len(ite.exifReader.rawBuffer) {
		return nil, ErrDataLength
	}
	// Read from underlying exifReader io.ReaderAt interface
	n, err = ite.exifReader.ReadAt(ite.exifReader.rawBuffer[:n], int64(ite.offset+ite.ifdOffset))

	ite.offset += uint32(n) // Update reader offset

	return ite.exifReader.rawBuffer[:n], err
}

func newTagEnumerator(offset uint32, er *reader) *ifdTagEnumerator {
	return &ifdTagEnumerator{
		exifReader: er,
		byteOrder:  er.ByteOrder(),
		ifdOffset:  offset,
	}
}

// parseUndefinedIfds
// Makernotes and AdobeDNGData
func (ite *ifdTagEnumerator) parseUndefinedIfds(e *Data, ifd ifds.IFD) bool {
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
func (ite *ifdTagEnumerator) ParseIfd(e *Data, ifd ifds.IFD, ifdIndex uint8, doDescend bool) (nextIfdOffset uint32, err error) {
	// Parse undefined Ifds
	if !ite.parseUndefinedIfds(e, ifd) {
		return 0, nil
	}

	// Determine tagCount
	tagCount, err := ite.ReadUint16()
	if err != nil {
		return 0, err
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

	for i := 0; i < int(tagCount); i++ {
		t, err := ite.ReadTag(ifd)
		if err != nil {
			if err == tag.ErrTagTypeNotValid {
				//if errors.Is(err, tag.ErrTagTypeNotValid) {
				// Log TagNotValid Error
				//ifdEnumerateLogger.Warningf(nil, "Tag in IFD [%s] at position (%d) has invalid type and will be skipped.", fqIfdPath, i)
				continue
			}
			return nextIfdOffset, err
		}

		// Descend into Child IFD
		childIFD := ifd.IsChildIfd(t)
		switch childIFD {
		case ifds.NullIFD:
			e.AddTag(ifd, ifdIndex, t)
		case ifds.SubIFD:
			if err := scanSubIfds(ite.exifReader, e, t); err != nil {
				return nextIfdOffset, err
			}
		default:
			if err := scan(ite.exifReader, e, childIFD, t.ValueOffset); err != nil {
				return nextIfdOffset, err
			}
		}

		//fmt.Printf("%s %s \t %v\n", ifd, ifd.TagName(t.TagID), t.TagType)
	}

	// NextIfdOffset
	if nextIfdOffset, err = ite.ReadUint32(); err != nil {
		return nextIfdOffset, err
	}

	// Adjust for incorrect Makernotes NextIfd Offsets
	// set nextIfdOffset to 0x0000.
	if ifd == ifds.MknoteIFD {
		nextIfdOffset = 0x0000
	}

	return
}

// ReadTag reads the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an ifdTagEnumerator
func (ite *ifdTagEnumerator) ReadTag(ifd ifds.IFD) (t tag.Tag, err error) {
	// Read 12 bytes of Tag
	buf, err := ite.ReadBuffer(12)
	if err != nil {
		return
	}
	tagID := tag.ID(ite.byteOrder.Uint16(buf[:2])) // TagID

	tagTypeRaw := ite.byteOrder.Uint16(buf[2:4]) // TagType

	unitCount := ite.byteOrder.Uint32(buf[4:8]) // UnitCount

	valueOffset := ite.byteOrder.Uint32(buf[8:12]) // ValueOffset

	tagType, err := tag.NewTagType(tagTypeRaw)
	if err != nil {
		return t, err
	}
	// Creates a newTag. If the TypeFromRaw is unsupported, it returns tag.ErrTagTypeNotValid.
	return tag.NewTag(tagID, tagType, unitCount, valueOffset, uint8(ifd)), err
}

// ReadUint16 reads a uint16 from an ifdTagEnumerator.
func (ite *ifdTagEnumerator) ReadUint16() (uint16, error) {
	buf, err := ite.ReadBuffer(2)
	if err != nil {
		return 0, err
	}
	return ite.byteOrder.Uint16(buf[:2]), nil
}

// ReadUint32 reads a uint32 from an ifdTagEnumerator.
func (ite *ifdTagEnumerator) ReadUint32() (uint32, error) {
	buf, err := ite.ReadBuffer(4)
	if err != nil {
		return 0, err
	}
	return ite.byteOrder.Uint32(buf[:4]), nil
}
