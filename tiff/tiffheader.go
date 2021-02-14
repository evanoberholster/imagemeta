package tiff

import (
	"encoding/binary"
	"errors"
)

// Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif TiffHeader
	ErrInvalidHeader = errors.New("error TiffHeader is not valid")
)

// TiffHeader is the first 8 bytes of a Tiff Directory.
//
// A TiffHeader contains the byte Order, first Ifd Offset,
// tiff Header offset, Exif Length (0 if unknown) and
// Image type for the parsing of the Exif information from
// a Tiff Directory.
type TiffHeader struct {
	ByteOrder        binary.ByteOrder
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
}

// NewTiffHeader returns a new TiffHeader.
func NewTiffHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32) TiffHeader {
	return TiffHeader{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}

// IsValid returns true if the TiffHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (th TiffHeader) IsValid() bool {
	return th.ByteOrder != nil || th.FirstIfdOffset > 0
}
