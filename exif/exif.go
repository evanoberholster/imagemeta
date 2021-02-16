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
func ScanExif(r io.ReaderAt) (e *ExifData, err error) {
	var it imagetype.ImageType
	er := newExifReader(r, nil, 0)
	br := bufio.NewReader(er)

	// Identify Image Type
	if it, err = imagetype.ScanBuf(br); err != nil {
		return
	}

	// Search Image for Metadata Header using
	// Imagetype information

	header, err := tiff.ScanTiff(br)
	if err != nil {
		if err != ErrNoExif {
			return
		}
	}

	// ExifData with an ExifReader attached
	e = newExifData(er, it)

	if err == nil {
		// Set TiffHeader sets the ExifReader and checks
		// the header validity.
		// Returns ErrInvalidHeader if header is not valid.
		if err = er.SetHeader(Header(header)); err != nil {
			return
		}

		// Scan the RootIFD with the FirstIfdOffset from the ExifReader
		err = scan(er, e, ifds.RootIFD, header.FirstIfdOffset)
	}
	return
}

// ParseExif parses Exif metadata from an io.ReaderAt and a tiff.Header and
// returns exif and an error.
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt, it imagetype.ImageType, header Header) (e *ExifData, err error) {
	er := newExifReader(r, nil, 0)

	// ExifData with an ExifReader attached
	e = newExifData(er, it)

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
