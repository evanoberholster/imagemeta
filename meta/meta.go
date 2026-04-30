// Package meta contains meta types for image metadata
package meta

import (
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
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

// ExifReader is invoked with Exif payload bytes and its parsed header.
type ExifReader func(r io.Reader, h ExifHeader) error

// XMPReader is invoked with XMP payload bytes and XPacket metadata.
type XMPReader func(r io.Reader, h XPacketHeader) error

// PreviewImageReader is invoked with preview image bytes and metadata header.
type PreviewImageReader func(r io.Reader, h PreviewHeader) error

// Metadata is common metadata among image parsers
type Metadata struct {
	// Exif Decode Function with ExifHeader
	ExifFn     DecodeFn
	ExifHeader ExifHeader

	// Xmp Decoding Function with XmpHeader
	XmpFn     DecodeFn
	XmpHeader XmpHeader

	// Dimensions of primary Image
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
	FirstIfd         tag.IfdType
	ImageType        imagetype.ImageType
}

// IsValid returns true if the ExifHeader ByteOrder is not nil and
// the FirstIfdOffset is greater than 0
func (h ExifHeader) IsValid() bool {
	return h.ByteOrder != utils.UnknownEndian && h.FirstIfdOffset > 0 && h.FirstIfd != tag.Unknown
}

func (h ExifHeader) String() string {
	return fmt.Sprintf("ByteOrder: %s, Ifd: %s, Offset: 0x%.4x TiffOffset: 0x%.4x Length: %d Imagetype: %s", h.ByteOrder, h.FirstIfd, h.FirstIfdOffset, h.TiffHeaderOffset, h.ExifLength, h.ImageType)
}

// NewExifHeader returns a new ExifHeader.
func NewExifHeader(byteOrder utils.ByteOrder, firstIfdOffset, tiffHeaderOffset uint32, exifLength uint32, imageType imagetype.ImageType) ExifHeader {
	return ExifHeader{
		ByteOrder:        byteOrder,
		FirstIfd:         tag.IFD0,
		FirstIfdOffset:   firstIfdOffset,
		TiffHeaderOffset: tiffHeaderOffset,
		ExifLength:       exifLength,
		ImageType:        imageType,
	}
}

// MarshalLogObject is a structured logging interface
func (h ExifHeader) MarshalLogObject(e *metalog.Event) {
	e.
		Str("firstIfd", h.FirstIfd.String()).
		Uint32("firstIfdOffset", h.FirstIfdOffset).
		Uint32("tiffHeaderOffset", h.TiffHeaderOffset).
		Uint32("exifLength", h.ExifLength).
		Str("byteOrder", h.ByteOrder.String()).
		Str("imageType", h.ImageType.String())
}

// XmpHeader is an XMP header of an image file.
// Contains Offset and Length of XMP metadata.
type XmpHeader struct {
	Offset, Length uint32
}

// XPacketHeader describes discovered XPacket payload metadata.
type XPacketHeader struct {
	Offset       uint64
	Length       int
	HasXPacketPI bool
	HasXMPMeta   bool
}

// NewXMPHeader returns a new xmp.Header with an offset
// and length of where to read XMP metadata.
func NewXMPHeader(offset, length uint32) XmpHeader {
	return XmpHeader{offset, length}
}

// MarshalLogObject is a structured logging interface.
func (h XmpHeader) MarshalLogObject(e *metalog.Event) {
	e.Uint32("offset", h.Offset).Uint32("length", h.Length)
}

// MarshalLogObject is a structured logging interface.
func (h XPacketHeader) MarshalLogObject(e *metalog.Event) {
	e.
		Uint64("offset", h.Offset).
		Int("length", h.Length).
		Bool("hasXPacketPI", h.HasXPacketPI).
		Bool("hasXMPMeta", h.HasXMPMeta)
}

type PreviewHeader struct {
	Size      uint32
	Width     uint16
	Height    uint16
	ImageType imagetype.ImageType
	Source    PreviewSource
}

// MarshalLogObject is a structured logging interface.
func (h PreviewHeader) MarshalLogObject(e *metalog.Event) {
	e.
		Uint32("size", h.Size).
		Uint16("width", h.Width).
		Uint16("height", h.Height).
		Str("imageType", h.ImageType.String()).
		Str("source", h.Source.String())
}

// PreviewSource identifies where preview image bytes were extracted from.
type PreviewSource uint8

const (
	PreviewSourceUnknown PreviewSource = iota
	PreviewSourceTHMB
	PreviewSourcePRVW
)

func (s PreviewSource) String() string {
	switch s {
	case PreviewSourceTHMB:
		return "THMB"
	case PreviewSourcePRVW:
		return "PRVW"
	default:
		return "unknown"
	}
}
