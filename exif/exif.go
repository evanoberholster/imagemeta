// Package exif provides functions for parsing and extracting Exif Information.
package exif

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
)

// Errors
var (
	// Alias to meta Errors
	ErrInvalidHeader = meta.ErrInvalidHeader
	ErrNoExif        = meta.ErrNoExif
)

// ScanExif identifies the imageType based on magic bytes and
// searches for exif headers, then it parses the io.ReaderAt for exif
// information and returns it.
// Sets exif imagetype from magicbytes, if not found sets imagetype
// to imagetypeUnknown.
//
// If no exif information is found ScanExif will return ErrNoExif.
func ScanExif(r meta.Reader) (e *Data, err error) {
	br := bufio.NewReaderSize(r, 64)

	// Identify Image Type
	it, err := imagetype.ScanBuf(br)
	if err != nil {
		return
	}

	// Search Image for Metadata Header using ImageType
	header, err := tiff.Scan(br)
	if err != nil {
		return
	}

	// Update Imagetype in ExifHeader
	header.ImageType = it

	// Set FirstIfd to RootIfd
	header.FirstIfd = ifds.RootIFD

	return ParseExif(r, header)
}

// ParseExif parses Exif metadata from an io.ReaderAt and a tiff.Header and
// returns exif and an error.
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt, header meta.ExifHeader) (e *Data, err error) {
	e, err = e.ParseExif(r, header)
	return e, err
}

func (e *Data) ParseExif(r io.ReaderAt, header meta.ExifHeader) (*Data, error) {
	if !header.IsValid() {
		return e, ErrInvalidHeader
	}

	if e == nil {
		e = newData(newExifReader(r, header.ByteOrder, header.TiffHeaderOffset, header.ExifLength), header.ImageType)
	}

	e.er.exifOffset = int64(header.TiffHeaderOffset)
	e.er.exifLength = header.ExifLength
	e.er.offset = 0

	if header.FirstIfd == ifds.NullIFD {
		header.FirstIfd = ifds.RootIFD
	}
	// Scan the FirstIfd with the FirstIfdOffset from the ExifReader
	err := scan(e.er, e, header.FirstIfd, header.FirstIfdOffset)
	return e, err
}
