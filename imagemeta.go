// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/heic"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/evanoberholster/imagemeta/xmp"
)

// Errors
var (
	ErrNoExif               = meta.ErrNoExif
	ErrNoExifDecodeFn       = errors.New("error no Exif Decode Func set")
	ErrNoXmpDecodeFn        = errors.New("error no Xmp Decode Func set")
	ErrImageTypeNotFound    = imagetype.ErrImageTypeNotFound
	ErrMetadataNotSupported = errors.New("error metadata reading not supported for this imagetype")
)

// Reader that is compatible with imagemeta
type Reader interface {
	io.ReaderAt
	io.ReadSeeker
}

// Metadata from an Image. The ExifDecodeFn and XmpDecodeFn
// are responsible for decoding their respective data.
type Metadata struct {
	r            Reader
	ExifDecodeFn exif.DecodeFn
	XmpDecodeFn  xmp.DecodeFn
	images       uint16
	exifHeader   exif.Header
	xmpHeader    xmp.Header
	size         meta.Dimensions
	t            imagetype.ImageType
	//Thumbnail Offsets
}

// Dimensions returns the primary image's width and height dimensions
func (m Metadata) Dimensions() meta.Dimensions {
	return m.size
}

// NewMetadata creates a new Metadata
func NewMetadata(r Reader, xmpDecodeFn xmp.DecodeFn, exifDecodeFn exif.DecodeFn) (meta *Metadata, err error) {
	meta = &Metadata{
		r:            r,
		XmpDecodeFn:  xmpDecodeFn,
		ExifDecodeFn: exifDecodeFn}
	// Create New bufio.Reader w/ 6KB because of XMP processing
	br := bufio.NewReaderSize(r, 160)
	// Pool
	// Identify image Type
	meta.t, err = imagetype.ScanBuf(br) // no discard
	if err != nil {
		return
	}
	// Parse ImageMetadata
	err = meta.parse(br)
	return
}

func (m *Metadata) parse(br *bufio.Reader) (err error) {
	switch m.t {
	case imagetype.ImageXMP:
		return m.parseXmp(br)
	case imagetype.ImageJPEG:
		return m.parseJpeg(br)
	case imagetype.ImageWebP:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageNEF:
		return m.parseTiff(br, m.t)
	case imagetype.ImageCR2:
		return m.parseTiff(br, m.t)
	case imagetype.ImageCR3:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageHEIF:
		return m.parseHeic(br)
	case imagetype.ImageAVIF:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImagePNG:
		err = ErrMetadataNotSupported
		return
	default:
		// process as Tiff
		// Bruteforce search for Exif header
		return m.parseTiff(br, m.t)
	}
}

// parseJpeg uses the 'jpeg' package to identify the metadata and the
// 'exif' and 'xmp' packages parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseJpeg(br *bufio.Reader) (err error) {
	jpegMeta, err := jpeg.ScanJPEG(br, m.XmpDecodeFn, m.ExifDecodeFn)
	if err != nil {
		return
	}
	m.exifHeader = jpegMeta.ExifHeader
	m.xmpHeader = jpegMeta.XmpHeader
	width, height := jpegMeta.Size()
	m.size = meta.NewDimensions(uint32(width), uint32(height))
	m.images = 1
	return nil
}

// parseXmp uses the 'xmp' package to identify and parse the metadata.
//
// Will use the custom decode function: XmpDecodeFn if it is not nil.
func (m *Metadata) parseXmp(br *bufio.Reader) (err error) {
	if m.XmpDecodeFn != nil {
		return m.XmpDecodeFn(br, xmp.NewHeader(0, 0))
	}
	return ErrNoXmpDecodeFn
}

// parseHeic uses the 'heic' package to identify the metadata and the
// 'exif' and 'xmp' packages parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseHeic(br *bufio.Reader) (err error) {
	if _, err = m.r.Seek(0, 0); err != nil {
		return
	}
	hm, err := heic.NewMetadata(m.r)
	if err != nil {
		return err
	}
	hm.ExifDecodeFn = m.ExifDecodeFn
	m.exifHeader = hm.ExifHeader
	m.xmpHeader = hm.XmpHeader
	m.size = hm.Dimensions()
	m.images = hm.Images()
	if m.ExifDecodeFn != nil {
		err = hm.DecodeExif(m.r)
	}
	if m.XmpDecodeFn != nil {
		return m.XmpDecodeFn(m.r, hm.XmpHeader)
	}
	// Add Support for XMP
	//hm.DecodeXmp(m.r)
	return err
}

// parseTiff uses the 'tiff' package to identify the metadata and
// the 'exif' and 'xmp' packages to parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseTiff(br *bufio.Reader, imageType imagetype.ImageType) (err error) {
	// package tiff -> exif
	header, err := tiff.Scan(br)
	if err != nil {
		return
	}
	// Update Header's ImageType
	header.ImageType = imageType
	return m.ExifDecodeFn(m.r, exif.Header(header))
}
