// Package meta contains meta types for image metadata
package meta

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
)

// Reader that is compatible with imagemeta
type Reader interface {
	io.ReaderAt
	io.ReadSeeker
}

// ExifDecodeFn is a function for decoding Exif Metadata
type ExifDecodeFn func(io.Reader, ExifHeader) error

// XmpDecodeFn is a function for decoding Xmp Metadata
type XmpDecodeFn func(io.Reader, XmpHeader) error

// Metadata is common metadata among image parsers
type Metadata struct {
	ExifDecodeFn ExifDecodeFn
	ExifHeader   ExifHeader
	XmpDecodeFn  XmpDecodeFn
	XmpHeader    XmpHeader
	Dim          Dimensions
	It           imagetype.ImageType
}

// Dimensions returns the Dimensions of the Primary Image.
func (m Metadata) Dimensions() Dimensions {
	return m.Dim
}

// ImageType returns the ImageType of the Primary Image.
func (m Metadata) ImageType() imagetype.ImageType {
	return m.It
}

// ExifHeader is the first 8 bytes of a Tiff Directory.
//
// A Header contains the byte Order, first Ifd Offset,
// Tiff Header offset, Exif Length (0 if unknown) and
// Image type for the parsing of the Exif information from
// a Tiff Directory.
type ExifHeader struct {
	ByteOrder        binary.ByteOrder
	FirstIfd         ifds.IFD
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	ImageType        imagetype.ImageType
}

// IsValid returns true if the ExifHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h ExifHeader) IsValid() bool {
	return h.ByteOrder != nil || h.FirstIfdOffset > 0
}

// NewExifHeader returns a new ExifHeader.
func NewExifHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32, imageType imagetype.ImageType) ExifHeader {
	return ExifHeader{
		ByteOrder:        byteOrder,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
		ImageType:        imageType,
	}
}

// XmpHeader is an XMP header of an image file.
// Contains Offset and Length of XMP metadata.
type XmpHeader struct {
	Offset, Length uint32
}

// NewXMPHeader returns a new xmp.Header with an offset
// and length of where to read XMP metadata.
func NewXMPHeader(offset, length uint32) XmpHeader {
	return XmpHeader{offset, length}
}
