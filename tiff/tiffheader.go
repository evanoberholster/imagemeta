package tiff

import (
	"encoding/binary"
	"errors"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
)

// Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif TiffHeader
	ErrInvalidHeader = errors.New("error TiffHeader is not valid")
)

// Header is the first 8 bytes of a Tiff Directory.
//
// A Header contains the byte Order, first Ifd Offset,
// tiff Header offset, Exif Length (0 if unknown) and
// Image type for the parsing of the Exif information from
// a Tiff Directory.
type Header struct {
	ByteOrder        binary.ByteOrder
	FirstIfd         ifds.IFD
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	ImageType        imagetype.ImageType
}

// NewHeader returns a new TiffHeader.
func NewHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32, imageType imagetype.ImageType) Header {
	return Header{
		ByteOrder:        byteOrder,
		FirstIfd:         ifds.RootIFD,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}

// IsValid returns true if the TiffHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (th Header) IsValid() bool {
	return th.ByteOrder != nil || th.FirstIfdOffset > 0
}
