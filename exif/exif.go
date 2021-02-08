// Package exif provides functions for parsing and extracting Exif Information.
package exif

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/metadata"
)

// Errors
var (
	// Alias to meta Errors
	ErrInvalidHeader = metadata.ErrInvalidHeader
	ErrNoExif        = metadata.ErrNoExif
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
	m, err := metadata.ScanBuf(br, it)
	if err != nil {
		if err != ErrNoExif {
			return
		}
	}

	// ExifData with an ExifReader attached
	e = newExifData(er, it)
	e.SetMetadata(m)

	if err == nil {
		header := m.Header()
		// Set TiffHeader sets the ExifReader and checks
		// the header validity.
		// Returns ErrInvalidHeader if header is not valid.
		if err = er.SetHeader(header); err != nil {
			return
		}

		// Scan the RootIFD with the FirstIfdOffset from the ExifReader
		err = scan(er, e, ifds.RootIFD, header.FirstIfdOffset)
	}
	return
}

// ParseExif parses a tiff header from the io.ReaderAt and
// returns exif and an error.
// Sets exif imagetype as imageTypeUnknown
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt) (e *ExifData, err error) {
	er := newExifReader(r, nil, 0)
	br := bufio.NewReader(er)

	// Search Image for Metadata Header using
	// Imagetype information
	m, err := metadata.ScanBuf(br, imagetype.ImageUnknown)
	if err != nil {
		if err != ErrNoExif {
			return
		}
	}

	// ExifData with an ExifReader attached
	e = newExifData(er, imagetype.ImageUnknown)
	e.SetMetadata(m)

	if err == nil {
		header := m.Header()
		// Set TiffHeader sets the ExifReader and checks
		// the header validity.
		// Returns ErrInvalidHeader if header is not valid.
		if err = er.SetHeader(header); err != nil {
			return
		}

		// Scan the RootIFD with the FirstIfdOffset from the ExifReader
		err = scan(er, e, ifds.RootIFD, header.FirstIfdOffset)
	}
	return
}
