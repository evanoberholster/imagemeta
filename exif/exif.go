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
func ScanExif(r io.ReaderAt) (e *Data, err error) {
	er := newExifReader(r, nil, 0)

	br := bufio.NewReaderSize(er, 64)

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

	return ParseExif(er, Header(header))
}

// ParseExif parses Exif metadata from an io.ReaderAt and a tiff.Header and
// returns exif and an error.
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt, header Header) (e *Data, err error) {
	er := newExifReader(r, nil, 0)

	// ExifData with an ExifReader attached
	e = newData(er, header.ImageType)

	// Set TiffHeader sets the ExifReader and checks
	// the header's validity.
	// Returns ErrInvalidHeader if header is not valid.
	if err = er.SetHeader(header); err != nil {
		return
	}

	// Scan the RootIFD with the FirstIfdOffset from the ExifReader
	err = scan(er, e, ifds.RootIFD, header.FirstIfdOffset)
	return
}
