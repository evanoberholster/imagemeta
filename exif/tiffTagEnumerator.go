package exif

import (
	"encoding/binary"
	"errors"
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
