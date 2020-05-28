package exiftool

import (
	"errors"
	"io"

	"github.com/evanoberholster/exiftool/exif"
	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/tag"
)

// Errors
var (
	// ErrInvalidHeader is an error for an Invalid Exif Header
	ErrInvalidHeader = errors.New("error exif header not valid")

	// ErrDataLength is an error for data length
	ErrDataLength = errors.New("Error the data is not long enough")

	// ErrIfdBufferLength
	ErrIfdBufferLength = errors.New("Ifd buffer length insufficient")
)

// IsValid returns true for a Valid ExifHeader
func (eh ExifHeader) isValid() bool {
	return eh.byteOrder != nil || eh.firstIfdOffset > 0x0000
}

// ParseExif parses an io.ReaderAt for exif informationan and returns it
func (eh ExifHeader) ParseExif(r io.ReaderAt) (e *exif.Exif, err error) {
	if !eh.isValid() {
		err = ErrInvalidHeader
		return
	}
	er := NewExifReader(r, eh.byteOrder, eh.tiffHeaderOffset)

	e = exif.NewExif(er)
	if err = er.scan(e, ifds.RootIFD, eh.firstIfdOffset); err != nil {
		return
	}
	return
}

// scan moves through an ifd at the specified offset and enumerates over the IfdTags
func (er *ExifReader) scan(e *exif.Exif, ifd ifds.IFD, offset uint32) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	for ifdIndex := 0; ; ifdIndex++ {
		enumerator := getTagEnumerator(offset, er)
		//fmt.Printf("Parsing IFD [%s] (%d) at offset (0x%04x).\n", ifd, ifdIndex, offset)
		nextIfdOffset, err := enumerator.ParseIfd(e, ifd, ifdIndex, true)
		if err != nil {
			return err
		}
		if nextIfdOffset == 0 {
			break
		}

		offset = nextIfdOffset
	}
	return
}

// scanSubIfds moves through the subIfds at the specified offsetes and enumerates over their IfdTags
func (er *ExifReader) scanSubIfds(e *exif.Exif, t tag.Tag) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// Fetch SubIfd Values from []Uint32 (LongType)
	offsets, err := t.Uint32Values(er)
	if err != nil {
		return err
	}

	for ifdIndex := 0; ifdIndex < len(offsets); ifdIndex++ {
		enumerator := getTagEnumerator(offsets[ifdIndex], er)
		//fmt.Printf("Parsing IFD [%s] (%d) at offset (0x%04x).\n", ifd, ifdIndex, offset)
		if _, err := enumerator.ParseIfd(e, ifds.SubIFD, ifdIndex, false); err != nil {
			// Log Error
			continue
		}
	}

	return
}
