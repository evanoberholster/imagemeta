// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"
	"errors"

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

// Metadata from an Image. The ExifDecodeFn and XmpDecodeFn
// are responsible for decoding their respective data.
type Metadata struct {
	r meta.Reader
	*meta.Metadata
	images uint16
	//Thumbnail Offsets
}

// NewMetadata creates a new Metadata
func NewMetadata(r meta.Reader, xmpFn meta.DecodeFn, exifFn meta.DecodeFn) (m *Metadata, err error) {
	m = &Metadata{r: r}
	m.Metadata = &meta.Metadata{
		XmpFn:  xmpFn,
		ExifFn: exifFn,
	}
	// Create New bufio.Reader w/ 6KB because of XMP processing
	br := bufio.NewReaderSize(r, 6*1024)
	// Pool
	// Identify image Type
	if m.It, err = imagetype.ScanBuf(br); err != nil {
		return
	}
	// Parse ImageMetadata
	err = m.parse(br)
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
	case imagetype.ImagePNG, imagetype.ImageBMP, imagetype.ImageGIF:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageCRW:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageUnknown:
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
	_, err = jpeg.ScanJPEG(br, m.Metadata)
	m.images = 1
	return err
}

// parseXmp uses the 'xmp' package to identify and parse the metadata.
//
// Will use the custom decode function: XmpDecodeFn if it is not nil.
func (m *Metadata) parseXmp(br *bufio.Reader) (err error) {
	if m.XmpFn != nil {
		m.Metadata.XmpHeader = meta.NewXMPHeader(0, 0)
		return m.XmpFn(br, m.Metadata)
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
	hm, err := heic.NewMetadata(br, m.Metadata)
	if err != nil {
		return err
	}
	m.images = hm.Images()
	if err = hm.ReadExif(m.r); err != nil {
		return
	}
	if err = hm.ReadXmp(m.r); err != nil {
		return
	}
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
	m.ExifHeader, err = tiff.Scan(br)
	if err != nil {
		return
	}
	m.ExifHeader.ImageType = m.It
	if m.ExifFn != nil {
		return m.ExifFn(m.r, m.Metadata)
	}
	return err
}
