// Package meta contains meta types for image metadata
package meta

import (
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

// Common Errors
var (
	// ErrInvalidHeader is an error for an Invalid ExifHeader
	ErrInvalidHeader = errors.New("error ExifHeader is not valid")

	// ErrNoExif is an error for when no exif is found
	ErrNoExif = errors.New("error no Exif")

	// ErrBufLength
	ErrBufLength = errors.New("error buffer length insufficient")
)

// Reader that is compatible with imagemeta
type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// DecodeFn is a function for decoding Metadata.
//
// For Exif Metadata use ExifHeader and ExifFn.
//
// For Xmp Metadata use XmpHeader and XmpFn.
type DecodeFn func(io.Reader, *Metadata) error

// Metadata is common metadata among image parsers
type Metadata struct {
	// Exif Decode Function with ExifHeader
	ExifFn     DecodeFn
	ExifHeader ExifHeader

	// Xmp Decoding Function with XmpHeader
	XmpFn     DecodeFn
	XmpHeader XmpHeader

	// Dimenions of primary Image
	Dim Dimensions

	// Image Type
	It imagetype.ImageType
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
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	ByteOrder        utils.ByteOrder
	FirstIfd         ifds.IfdType
	ImageType        imagetype.ImageType
}

// IsValid returns true if the ExifHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h ExifHeader) IsValid() bool {
	return h.ByteOrder != utils.UnknownEndian && h.FirstIfdOffset > 0 && h.FirstIfd != ifds.NullIFD
}

func (h ExifHeader) String() string {
	return fmt.Sprintf("ByteOrder: %s, Ifd: %s, Offset: 0x%.4x TiffOffset: 0x%.4x Length: %d Imagetype: %s", h.ByteOrder, h.FirstIfd, h.FirstIfdOffset, h.TiffHeaderOffset, h.ExifLength, h.ImageType)
}

// NewExifHeader returns a new ExifHeader.
func NewExifHeader(byteOrder utils.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32, imageType imagetype.ImageType) ExifHeader {
	return ExifHeader{
		ByteOrder:        byteOrder,
		FirstIfd:         ifds.IFD0,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
		ImageType:        imageType,
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (h ExifHeader) MarshalZerologObject(e *zerolog.Event) {
	e.Str("FirstIfd", h.FirstIfd.String()).Uint32("FirstIfdOffset", h.FirstIfdOffset).Uint32("TiffHeaderOffset", h.TiffHeaderOffset).Uint32("ExifLength", h.ExifLength).Str("Endian", h.ByteOrder.String()).Str("ImageType", h.ImageType.String())
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

type PreviewHeader struct {
	Size   uint32
	Width  uint16
	Height uint16
}
