package exif

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/tiff"
)

// DecodeFn is a function for decoding Exif Metadata
type DecodeFn func(r io.Reader, header Header) (err error)

// Header is the same as a tiff.Header
type Header tiff.Header

// IsValid returns true if the TiffHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h Header) IsValid() bool {
	return h.ByteOrder != nil || h.FirstIfdOffset > 0
}

// NewHeader returns a new exif.Header.
func NewHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32) Header {
	return Header{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}
