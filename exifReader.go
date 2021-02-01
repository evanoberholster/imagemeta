package imagemeta

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// ExifReader -
type ExifReader struct {
	reader io.ReaderAt

	// Exif Header
	byteOrder  binary.ByteOrder
	exifOffset int64
	exifLength uint32

	// reader interface offset
	offset int64
}

// newExifReader returns a new ExifReader. It reads from reader according to byteOrder from exifOffset
func newExifReader(reader io.ReaderAt, byteOrder binary.ByteOrder, exifOffset uint32) *ExifReader {
	return &ExifReader{
		reader:     reader,
		byteOrder:  byteOrder,
		exifOffset: int64(exifOffset),
	}
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
func (er ExifReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("ExifReader.ReadAt: negative offset")
	}
	n, err = er.reader.ReadAt(p, er.exifOffset+off)
	return
}

// ByteOrder returns the ExifReader's byteOrder
func (er *ExifReader) ByteOrder() binary.ByteOrder {
	return er.byteOrder
}

// SetHeader sets the ByteOrder, exifOffset and exifLength of an ExifReader
// from a TiffHeader and sets the ExifReader read offset to 0
func (er *ExifReader) SetHeader(header meta.TiffHeader) error {
	if !header.IsValid() {
		return ErrInvalidHeader
	}
	er.byteOrder = header.ByteOrder
	er.exifOffset = int64(header.TiffHeaderOffset)
	er.exifLength = header.ExifLength
	er.offset = 0
	return nil
}
