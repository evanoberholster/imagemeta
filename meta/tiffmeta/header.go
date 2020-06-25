package tiffmeta

import "encoding/binary"

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
	ExifLength       uint16
	//Imagetype
}

// NewHeader returns a new Tiff Header.
func NewHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint16) Header {
	return Header{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
	}
}
