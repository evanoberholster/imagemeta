// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/cr3"
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

// ImageMeta interface for Image Metadata
type ImageMeta interface {
	Dimensions() meta.Dimensions
	ImageType() imagetype.ImageType
	PreviewImage() io.Reader
	Exif() (exif.Exif, error)
	Xmp() (xmp.XMP, error)
}

// Parse meta.Reader for Image Metadata returns ImageMeta corresponding
// to identified image type.
func Parse(r meta.Reader) (ImageMeta, error) {
	t, err := imagetype.ReadAt(r)
	if err != nil {
		return nil, err
	}
	switch t {

	case imagetype.ImageJPEG:
		return jpeg.ScanJPEG(r, nil, nil)
	case imagetype.ImageCR3:
		return cr3.Parse(r)
	case imagetype.ImageTiff, imagetype.ImageCR2, imagetype.ImageARW, imagetype.ImageHEIF, imagetype.ImageNEF, imagetype.ImagePanaRAW:
		return tiff.Parse(r, t)
	}
	return nil, nil
}

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
	case imagetype.ImageWebP:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageNEF:
		return m.parseTiff(br)
	case imagetype.ImageCR2:
		return m.parseTiff(br)
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

// parseHeic uses the 'heic' package to identify the metadata and the
// 'exif' and 'xmp' packages parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseHeic(br *bufio.Reader) (err error) {
	//if _, err = m.r.Seek(0, 0); err != nil {
	//	return
	//}
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

// parseTiff uses the 'tiff' package to identify the metadata and
// the 'exif' and 'xmp' packages to parse the metadata.
//
// Will use the custom decode functions: XmpDecodeFn and
// ExifDecodeFn if they are not nil.
func (m *Metadata) parseTiff(br *bufio.Reader) (err error) {
	// package tiff -> exif
	m.ExifHeader, err = tiff.ScanTiffHeader(br, imagetype.ImageTiff)
	if err != nil {
		return
	}
	m.ExifHeader.ImageType = m.It
	if m.ExifFn != nil {
		return m.ExifFn(m.r, m.Metadata)
	}
	return err
}
