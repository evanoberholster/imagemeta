package exif

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif/tag"
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

func (er *reader) TagValue(t tag.Tag) (buf []byte, err error) {
	// check if Value is Embedded
	if t.IsEmbedded() {
		er.ByteOrder().PutUint32(er.rawBuffer[:4], t.ValueOffset)
		return er.rawBuffer[:4], nil
	}

	byteLength := t.Size()
	if byteLength <= len(er.rawBuffer) {
		buf = er.rawBuffer[:byteLength]
	} else {
		buf = make([]byte, byteLength)
	}
	exifOffset := er.ifdExifOffset[t.Ifd]
	_, err = er.reader.ReadAt(buf[:byteLength], int64(exifOffset+t.ValueOffset))
	return buf[:byteLength], err
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

// ByteOrder returns the ExifReader's byteOrder
func (er *reader) ByteOrder() binary.ByteOrder {
	return er.byteOrder
}

// SetHeader sets the ByteOrder, exifOffset and exifLength of an ExifReader
// from a TiffHeader and sets the ExifReader read offset to 0
func (er *reader) SetHeader(header Header) error {
	if !header.IsValid() {
		return ErrInvalidHeader
	}
	er.byteOrder = header.ByteOrder
	er.exifOffset = int64(header.TiffHeaderOffset)
	er.exifLength = header.ExifLength
	er.offset = 0
	return nil
}
