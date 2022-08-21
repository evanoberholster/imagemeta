package exif

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// reader errors
var (
	ErrReadNegativeOffset = errors.New("error read at negative offset")
)

const rawBufferSize = 64

// reader is an EXIF Reader that uses an underlying ReaderAt and rawBuffer.
type reader struct {
	// Underlying reader and offset
	u io.ReaderAt

	// Exif byteOrder
	byteOrder binary.ByteOrder

	// Offsets for multiple Ifds
	ifdExifOffset [8]uint32

	// rawBuffer for parsing Tags
	rawBuffer [rawBufferSize]byte

	exifOffset uint32
	exifLength uint32
}

// newReader returns a new Reader. It reads from reader according to byteOrder from exifOffset
func newReader(r io.ReaderAt, header meta.ExifHeader) *reader {
	return &reader{
		u:          r,
		byteOrder:  header.ByteOrder,
		exifLength: header.ExifLength,
		exifOffset: header.TiffHeaderOffset,
	}
}

// scanIFD scans through an ifd at the specified offset and enumerates over the IfdTags
func (r *reader) scanIFD(e *Data, ifd ifds.Ifd) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	var nextIfdOffset uint32

	for ifd.Index = 0; ; ifd.Index++ {
		r.ifdExifOffset[ifd.Type] = uint32(r.exifOffset)
		ifd.Offset += r.exifOffset

		if nextIfdOffset, err = r.parseIfd(e, ifd, true); err != nil {
			return err
		}
		if nextIfdOffset == 0 {
			break
		}
		ifd.Offset = nextIfdOffset
	}
	return
}

// scanSubIFD scans through the subIfd at the specified offset and enumerates over their IfdTags
func (r *reader) scanSubIFD(e *Data, t tag.Tag) (err error) {
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

	for ifdIndex := uint8(0); ifdIndex < uint8(len(offsets)); ifdIndex++ {
		ifdOffset := offsets[ifdIndex]
		ifd := ifds.NewIFD(ifds.SubIFD, ifdIndex, ifdOffset)
		if _, err = r.parseIfd(e, ifd, false); err != nil {
			return errors.WithMessage(err, "ScanSubIfds: ParseIfd Error")
		}
	}
	return
}

// ParseIfd - enumerates over the ifd using the enumerator.ifdReader
func (r *reader) parseIfd(e *Data, ifd ifds.Ifd, doDescend bool) (nextIfdOffset uint32, err error) {
	byteOrder := r.byteOrder

	// Parse MakerNoteIfds
	if ifd.IsType(ifds.MknoteIFD) {
		ifd, byteOrder = r.parseMknoteIFD(e, ifd)
		if byteOrder == nil {
			return 0, nil
		}
	}

	offset := ifd.Offset

	var tagCount uint16
	var t tag.Tag

	// Determine tagCount
	if tagCount, offset, err = r.ReadUint16(byteOrder, offset); err != nil {
		return 0, errors.Wrapf(err, "Tag Count: %d for %s", tagCount, ifd.String())
	}
	if tagCount > 255 {
		return 0, errors.Errorf("Tagcount too high. Tag Count: %d for %s", tagCount, ifd.String())
	}

	// Log Ifd Info
	if isInfo() {
		logIfdInfo(ifd, tagCount, offset)
	}

	for i := 0; i < int(tagCount); i++ {
		if t, offset, err = r.ReadTag(ifd, byteOrder, offset); err != nil {
			if err == tag.ErrTagTypeNotValid {
				//if errors.Is(err, tag.ErrTagTypeNotValid) {
				// Log TagNotValid Error
				//ifdEnumerateLogger.Warningf(nil, "Tag in IFD [%s] at position (%d) has invalid type and will be skipped.", fqIfdPath, i)
				continue
			}
			return nextIfdOffset, err
		}

		// Log Tag Info
		if isInfo() {
			logTagInfo(ifd, t, offset)
		}

		// Tag is an Ifd then descend
		if t.IsIfd() {
			// Descend into Child IFD
			childIfd := ifd.ChildIfd(t)
			if childIfd.IsType(ifds.SubIFD) {
				if err := r.scanSubIFD(e, t); err != nil {
					return offset, err
				}
			} else {
				if err := r.scanIFD(e, childIfd); err != nil {
					return offset, err
				}
			}
		} else {
			e.addTag(ifd, t) // Add Tag to Map
		}
	}

	// NextIfdOffset
	if nextIfdOffset, _, err = r.ReadUint32(byteOrder, offset); err != nil {
		return nextIfdOffset, err
	}

	// Adjust for incorrect Makernotes NextIfd Offsets set nextIfdOffset to 0x0000.
	if ifd.IsType(ifds.MknoteIFD) {
		nextIfdOffset = 0x0000
	}

	return
}

func (r *reader) embeddedTagValue(valueOffset uint32) []byte {
	r.byteOrder.PutUint32(r.rawBuffer[:4], valueOffset)
	return r.rawBuffer[:4]
}

// ReadValue returns the Tag's Value as a byte slice.
func (r *reader) ReadValue(t tag.Tag) (buf []byte, err error) {
	if t.IsEmbedded() {
		return r.embeddedTagValue(t.ValueOffset), nil // return tag Value if Embedded
	}

	byteLength := int(t.Size())           // Tag Value Size
	valueOffset := t.ValueOffset          // Tag Value Offset
	valueOffset += r.ifdExifOffset[t.Ifd] // Exif Offset for the given Tag's Ifd

	return r.ReadBufferAt(byteLength, int(valueOffset))
}

// Read Lengths
const (
	tagByteLength    = 12
	uint16ByteLength = 2
	uint32ByteLength = 4
)

// ReadTag reads the tagID uint16, tagType uint16, unitCount uint32 and valueOffset uint32
// from an ifdTagEnumerator. Returns Tag and error. If the tagType is unsupported, returns tag.ErrTagTypeNotValid.
func (r *reader) ReadTag(ifd ifds.Ifd, byteOrder binary.ByteOrder, offset uint32) (tag.Tag, uint32, error) {
	buf, err := r.ReadBufferAt(tagByteLength, int(offset))
	if err != nil {
		return tag.Tag{}, offset, err
	}
	tagID := tag.ID(byteOrder.Uint16(buf[:2]))      // TagID
	tagType := tag.Type(byteOrder.Uint16(buf[2:4])) // TagType
	unitCount := byteOrder.Uint32(buf[4:8])         // UnitCount
	valueOffset := byteOrder.Uint32(buf[8:12])      // ValueOffset

	tagType = tagIsIfd(ifd, tagID, tagType)

	t, err := tag.NewTag(tagID, tagType, unitCount, valueOffset, uint8(ifd.Type)) // NewTag
	return t, offset + tagByteLength, err
}

// ReadUint16 reads a uint16 from an ifdTagEnumerator.
func (r *reader) ReadUint16(byteOrder binary.ByteOrder, offset uint32) (val uint16, off uint32, err error) {
	buf, err := r.ReadBufferAt(uint16ByteLength, int(offset))
	return byteOrder.Uint16(buf), offset + uint16ByteLength, err
}

// ReadUint32 reads a uint32 from an ifdTagEnumerator.
func (r *reader) ReadUint32(byteOrder binary.ByteOrder, offset uint32) (val uint32, off uint32, err error) {
	buf, err := r.ReadBufferAt(uint32ByteLength, int(offset))
	return byteOrder.Uint32(buf), offset + uint32ByteLength, err
}

// ReadBufferAt reads n at offset from the underlying reader.
func (r *reader) ReadBufferAt(n int, offset int) (buf []byte, err error) {
	if n <= rawBufferSize {
		buf = r.rawBuffer[:n]
	} else {
		buf = make([]byte, n)
	}

	nn, err := r.u.ReadAt(buf[:n], int64(offset))
	if nn < n {
		return nil, errors.Wrapf(err, "ReadBufferAt error wanted %d bytes got %d bytes", n, nn)
	}
	return buf[:n], nil
}

func tagIsIfd(ifd ifds.Ifd, tagID tag.ID, tagType tag.Type) tag.Type {
	if tagType.Is(tag.TypeLong) {
		// RootIfd Children
		if ifd.IsType(ifds.IFD0) {
			switch tagID {
			case ifds.ExifTag:
				return tag.TypeIfd
			case ifds.GPSTag:
				return tag.TypeIfd
			case ifds.SubIFDs:
				return tag.TypeIfd
			}
		}
	}
	if tagType.Is(tag.TypeUndefined) {
		// ExifIfd Children
		if ifd.IsType(ifds.ExifIFD) {
			switch tagID {
			case exififd.MakerNote:
				return tag.TypeIfd
			}
		}
	}
	return tagType
}

// AddTag adds a Tag to a tag.TagMap
func (e *Data) addTag(ifd ifds.Ifd, t tag.Tag) {
	if !ifd.IsValid() {
		return // don't add invalid IFD tags
	}

	// add tag to tagMap
	e.tagMap[ifds.Key{Type: ifd.Type, Index: ifd.Index, TagID: t.ID}] = t

	// Special Ifd0 Tags
	if ifd.IsType(ifds.IFD0) {
		switch t.ID {
		case ifds.Make: // Add Make and Model to Exif struct for future decoding of Makernotes
			e.make, _ = e.ParseASCIIValue(t)
		case ifds.Model:
			e.model, _ = e.ParseASCIIValue(t)
		case ifds.ImageWidth:
			if ifd.Index == 0 {
				e.width, _ = e.ParseUint16Value(t)
			}
		case ifds.ImageLength:
			if ifd.Index == 0 {
				e.height, _ = e.ParseUint16Value(t)
			}
		}
	}
	// Special ExifIfd Tags
	if ifd.IsType(ifds.ExifIFD) {
		switch t.ID {
		case exififd.PixelXDimension:
			e.width, _ = e.ParseUint16Value(t)
		case exififd.PixelYDimension:
			e.height, _ = e.ParseUint16Value(t)
		case exififd.ExifVersion:
			e.exifVersion = e.parseExifVersion(t)
		}
	}
}

// ReadBufferAt reads n at offset from the underlying reader.
//func (r *reader2) ReadBufferAt2(n int, offset int) (buf []byte, err error) {
//	if offset+n > r.bufferedLen() {
//		if err = r.grow(int(n + offset - r.bufferedLen())); err != nil {
//			return nil, ErrDataLength
//		}
//	}
//	return r.buf[offset : offset+n], nil
//}
//
//// grow the underlying buffer
//func (r *reader2) grow(n int) error {
//	if n < minBufferLength {
//		n = minBufferLength
//	}
//	buf := makeSlice(cap(r.buf) + int(n))
//	copy(buf, r.buf[:r.uOffset])
//	r.buf = buf
//	//n, err := r.u.Read(r.buf[r.uOffset:])
//	r.uOffset += uint64(n)
//	//return err
//	return nil
//}
//
