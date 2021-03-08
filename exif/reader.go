package exif

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/pkg/errors"
)

// reader errors
var (
	ErrReadNegativeOffset = errors.New("error read at negative offset")
)

const rawBufferSize = 24

// reader -
type reader struct {
	reader io.ReaderAt
	// current reader offset
	offset int64

	// Exif Header
	byteOrder binary.ByteOrder

	// Offsets for multiple Ifds
	ifdExifOffset [7]uint32

	// rawBuffer for parsing Tags
	rawBuffer [rawBufferSize]byte

	exifOffset int64
	exifLength uint32
}

// newExifReader returns a new ExifReader. It reads from reader according to byteOrder from exifOffset
func newExifReader(r io.ReaderAt, byteOrder binary.ByteOrder, exifOffset uint32, exifLength uint32) *reader {
	er, ok := r.(*reader)
	if ok {
		return er
	}
	return &reader{
		reader:     r,
		byteOrder:  byteOrder,
		exifOffset: int64(exifOffset),
		exifLength: exifLength,
	}
}

// Read reads from ExifReader and moves the placement marker
func (er *reader) Read(p []byte) (n int, err error) {
	// Buffer is empty
	if len(p) == 0 {
		return 0, nil
	}
	n, err = er.reader.ReadAt(p, er.exifOffset+er.offset)
	er.offset += int64(n)

	return n, err
}

// ReadAt reads from ExifReader at the given offset
func (er reader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ErrReadNegativeOffset
	}
	n, err = er.reader.ReadAt(p, er.exifOffset+off)
	return
}

// ReadBufferAt reads from ExifReader at the give offset and returns
// the reader's underlying buffer. Only valid until next read.
func (er reader) ReadBufferAt(n int, off int64) ([]byte, error) {
	n, err := er.ReadAt(er.rawBuffer[:n], off)
	return er.rawBuffer[:n], err
}

// TagValue returns the Tag's Value as a byte slice.
// It allocates a new []byte when the value is larger than the exifReader's underlying rawBuffer
// and it is not a an embedded tag.
func (er *reader) TagValue(t tag.Tag) (buf []byte, err error) {
	if t.IsEmbedded() { // check if Value is Embedded
		return er.embeddedTagValue(t.ValueOffset), nil
	}

	byteLength := t.Size()
	if byteLength <= rawBufferSize {
		buf = er.rawBuffer[:byteLength]
	} else {
		buf = make([]byte, byteLength)
	}
	exifOffset := er.ifdExifOffset[t.Ifd] // Offset for the given Tag's Ifd
	n, err := er.reader.ReadAt(buf[:byteLength], int64(exifOffset+t.ValueOffset))
	if n < t.Size() {
		err = errors.Wrap(err, tag.ErrNotEnoughData.Error()) // FixMe Please
	}
	return buf[:byteLength], err
}

func (er *reader) embeddedTagValue(valueOffset uint32) []byte {
	er.byteOrder.PutUint32(er.rawBuffer[:4], valueOffset)
	return er.rawBuffer[:4]
}

// ByteOrder returns the ExifReader's byteOrder
func (er *reader) ByteOrder() binary.ByteOrder {
	return er.byteOrder
}
