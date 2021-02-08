package exif

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/metadata"
)

// reader errors
var (
	ErrReadNegativeOffset = errors.New("error read at negative offset")
)

// reader -
type reader struct {
	reader io.ReaderAt

	// Exif Header
	byteOrder  binary.ByteOrder
	exifOffset int64
	exifLength uint32

	// reader interface offset
	offset int64
}

// newExifReader returns a new ExifReader. It reads from reader according to byteOrder from exifOffset
func newExifReader(r io.ReaderAt, byteOrder binary.ByteOrder, exifOffset uint32) *reader {
	return &reader{
		reader:     r,
		byteOrder:  byteOrder,
		exifOffset: int64(exifOffset),
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

// ByteOrder returns the ExifReader's byteOrder
func (er *reader) ByteOrder() binary.ByteOrder {
	return er.byteOrder
}

// SetHeader sets the ByteOrder, exifOffset and exifLength of an ExifReader
// from a TiffHeader and sets the ExifReader read offset to 0
func (er *reader) SetHeader(header metadata.TiffHeader) error {
	if !header.IsValid() {
		return ErrInvalidHeader
	}
	er.byteOrder = header.ByteOrder
	er.exifOffset = int64(header.TiffHeaderOffset)
	er.exifLength = header.ExifLength
	er.offset = 0
	return nil
}
