// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/cr3"
	"github.com/evanoberholster/imagemeta/heic"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
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
	meta.Metadata
	r      Reader
	images uint16
	//Thumbnail Offsets
}

// NewMetadata creates a new Metadata
func NewMetadata(r Reader, xmpDecodeFn meta.XmpDecodeFn, exifDecodeFn meta.ExifDecodeFn) (meta *Metadata, err error) {
	meta = &Metadata{r: r}
	meta.XmpDecodeFn = xmpDecodeFn
	meta.ExifDecodeFn = exifDecodeFn
	// Create New bufio.Reader w/ 6KB because of XMP processing
	br := bufio.NewReaderSize(r, 160)
	// Pool
	// Identify image Type
	meta.It, err = imagetype.ScanBuf(br) // no discard
	if err != nil {
		return
	}
	// Parse ImageMetadata
	err = meta.parse(br)
	return
}

func (m *Metadata) parse(br *bufio.Reader) (err error) {
	switch m.It {
	case imagetype.ImageXMP:
		return m.parseXmp(br)
	case imagetype.ImageJPEG:
		return m.parseJpeg(br)
	case imagetype.ImageWebP:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageNEF:
		return m.parseTiff(br)
	case imagetype.ImageCR2:
		return m.parseTiff(br)
	case imagetype.ImageCR3:
		return m.parseCR3(br)
	case imagetype.ImageHEIF:
		return m.parseHeic(br)
	case imagetype.ImageAVIF:
		return m.parseHeic(br)
		//err = ErrMetadataNotSupported
		//return
	case imagetype.ImagePNG:
		err = ErrMetadataNotSupported
		return
	default:
		// process as Tiff
		// Bruteforce search for Exif header
		return m.parseTiff(br)
	}
}

// parseJpeg uses the 'jpeg' package to identify the metadata and the
// 'exif' and 'xmp' packages parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseJpeg(br *bufio.Reader) (err error) {
	jpegMeta, err := jpeg.ScanJPEG(br, m.Metadata)
	if err != nil {
		return
	}
	m.Metadata = jpegMeta.Metadata
	m.images = 1
	return nil
}

// parseXmp uses the 'xmp' package to identify and parse the metadata.
//
// Will use the custom decode function: XmpDecodeFn if it is not nil.
func (m *Metadata) parseXmp(br *bufio.Reader) (err error) {
	if m.XmpDecodeFn != nil {
		return m.XmpDecodeFn(br, meta.NewXMPHeader(0, 0))
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
	hm, err := heic.NewMetadata(m.r, m.Metadata)
	if err != nil {
		return err
	}
	//m.Metadata = hm.Metadata
	m.images = hm.Images()
	if m.ExifDecodeFn != nil {
		err = hm.DecodeExif(m.r)
	}
	if hm.XmpDecodeFn != nil {
		return hm.XmpDecodeFn(m.r, hm.XmpHeader)
	}
	m.Metadata = hm.Metadata
	// Add Support for XMP
	//hm.DecodeXmp(m.r)
	return err
}

func (m *Metadata) parseCR3(br *bufio.Reader) (err error) {
	cr3, err := cr3.NewMetadata(br, m.Metadata)
	if err != nil {
		return err
	}

	if err = cr3.DecodeExif(br); err != nil {
		return
	}
	if err = cr3.DecodeXMP(br); err != nil {
		return
	}
	m.Metadata = cr3.Metadata
	return nil
}

// parseTiff uses the 'tiff' package to identify the metadata and
// the 'exif' and 'xmp' packages to parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseTiff(br *bufio.Reader) (err error) {
	// package tiff -> exif
	header, err := tiff.Scan(br)
	if err != nil {
		return
	}
	m.Metadata.ExifHeader = header
	// Update Header's ImageType
	header.ImageType = m.It
	return m.ExifDecodeFn(m.r, header)
}
