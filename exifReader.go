package exiftool

import (
	"encoding/binary"
	"errors"
	"io"
)

// ExifReader -
type ExifReader struct {
	reader io.ReaderAt

	// previously ExifHeader
	byteOrder  binary.ByteOrder
	exifOffset int64

	// reader interface offset
	offset int64
}

// Read reads from ExifReader and moves the placement marker
func (er *ExifReader) Read(p []byte) (n int, err error) {
	// Buffer is empty
	if len(p) == 0 {
		return 0, nil
	}
	n, err = er.reader.ReadAt(p, er.exifOffset+er.offset)
	er.offset += int64(n)

	return n, err
}

// ReadAt reads from ExifReader at the given offset
func (er *ExifReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("ExifReader.ReadAt: negative offset")
	}
	n, err = er.reader.ReadAt(p, er.exifOffset+off)
	return
}

// NewExifReader returns a new ExifReader. It reads from reader according to byteOrder from exifOffset
func NewExifReader(reader io.ReaderAt, byteOrder binary.ByteOrder, exifOffset uint32) *ExifReader {
	return &ExifReader{
		reader:     reader,
		byteOrder:  byteOrder,
		exifOffset: int64(exifOffset),
	}
}

// ByteOrder returns the ExifReader's byteOrder
func (er ExifReader) ByteOrder() binary.ByteOrder {
	return er.byteOrder
}
