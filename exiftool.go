// Package exiftool provides functions for scanning for Exif Information and extracting it
package exiftool

import (
	"bufio"
	"io"

	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/imagetype"
	"github.com/evanoberholster/exiftool/meta"
	"github.com/evanoberholster/exiftool/meta/tiffmeta"
)

// Errors
var (
	// Alias to tiffmeta ErrInvalidHeader
	ErrInvalidHeader = tiffmeta.ErrInvalidHeader
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
	er := newExifReader(r, nil, 0)
	br := bufio.NewReader(er)

	// Identify Image Type
	t, err := imagetype.ScanBuf(br)
	if err != nil {
		return
	}

	// Search Image for Metadata Header using
	// Imagetype information
	m, err := meta.ScanBuf(br, t)
	if err != nil {
		if err != ErrNoExif {
			return
		}
	}

	// NewExif with an ExifReader attached
	e = newExifData(er, t)
	e.SetMetadata(m)

	if err == nil {
		header := m.TiffHeader()
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
	m, err := meta.ScanBuf(br, imagetype.ImageUnknown)
	if err != nil {
		if err != ErrNoExif {
			return
		}
	}

	// NewExif with an ExifReader attached
	e = newExifData(er, imagetype.ImageUnknown)
	e.SetMetadata(m)

	if err == nil {
		header := m.TiffHeader()
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
