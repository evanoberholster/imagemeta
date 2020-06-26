package tiffmeta

import (
	"encoding/binary"
	"errors"
)

// Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif Header
	ErrInvalidHeader = errors.New("Error Tiff Header is not valid")
)

// Header represents a Tiff Header.
// The first 8 bytes in a Tiff Directory.
//
// Header contains the byte Order, first Ifd Offset,
// tiff Header offset, Exif Length (0 if unknown) and
// Image type for the parsing of the Exif information.
type Header struct {
	ByteOrder        binary.ByteOrder
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	//Imagetype
}

// NewHeader returns a new Tiff Header.
func NewHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32) Header {
	return Header{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}

// IsValid returns true if the Header ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h Header) IsValid() bool {
	return h.ByteOrder != nil || h.FirstIfdOffset > 0
}
