// Package meta contains meta types for image metadata
package meta

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
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
	io.ReaderAt
	io.ReadSeeker
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
	ByteOrder        binary.ByteOrder
	FirstIfd         ifds.IfdType
	FirstIfdOffset   uint32
	TiffHeaderOffset uint32
	ExifLength       uint32
	ImageType        imagetype.ImageType
}

// IsValid returns true if the ExifHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h ExifHeader) IsValid() bool {
	return h.ByteOrder != nil && h.FirstIfdOffset > 0 && h.FirstIfd != ifds.NullIFD
}

// NewExifHeader returns a new ExifHeader.
func NewExifHeader(byteOrder binary.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32, imageType imagetype.ImageType) ExifHeader {
	return ExifHeader{
		ByteOrder:        byteOrder,
		FirstIfd:         ifds.IFD0,
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

// BinaryOrder returns the binary.ByteOrder for a Tiff Header based
// on 4 bytes from the buf.
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func BinaryOrder(buf []byte) binary.ByteOrder {
	if isTiffBigEndian(buf[:4]) {
		return binary.BigEndian
	}
	if isTiffLittleEndian(buf[:4]) {
		return binary.LittleEndian
	}
	return nil
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func isTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func isTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}
