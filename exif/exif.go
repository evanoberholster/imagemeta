// Package exif provides functions for parsing and extracting Exif Information.
package exif

import (
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
)

// Errors
var (
	// Alias to meta Errors
	ErrInvalidHeader = meta.ErrInvalidHeader
	ErrNoExif        = meta.ErrNoExif
	ErrEmptyTag      = errors.New("error empty tag")
)

// ScanExif identifies the imageType based on magic bytes and
// searches for exif headers, then it parses the io.ReaderAt for exif
// information and returns it.
// Sets exif imagetype from magicbytes, if not found sets imagetype
// to imagetypeUnknown.
//
// If no exif information is found ScanExif will return ErrNoExif.
func ScanExif(r meta.Reader) (e *Data, err error) {
	// Search Image for Metadata Header using ImageType
	header, err := tiff.ScanTiffHeader(r, imagetype.ImageTiff)
	if !header.IsValid() || err != nil {
		return nil, ErrNoExif
	}
	// Set FirstIfd to RootIfd
	header.FirstIfd = ifds.IFD0

	return ParseExif(r, header)
}

// ParseExif parses Exif metadata from an io.ReaderAt and a TiffHeader
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt, header meta.ExifHeader) (*Data, error) {
	var err error
	if !header.IsValid() {
		return nil, ErrInvalidHeader
	}

	if header.FirstIfd == ifds.NullIFD {
		header.FirstIfd = ifds.IFD0
	}

	reader := newReader(r, header)

	e := newData(reader, header.ImageType)

	// Scan the FirstIfd with the FirstIfdOffset from the ExifReader
	err = reader.scanIFD(e, ifds.NewIFD(header.FirstIfd, 0, header.FirstIfdOffset))
	//fmt.Println(reader.bufferedLen())
	return e, err
}

// newData creates a new initialized Exif object
func newData(r *reader, it imagetype.ImageType) *Data {
	return &Data{
		reader:    r,
		imageType: it,
		tagMap:    make(ifds.TagMap, 50),
	}
}

// Data struct contains parsed Exif information
type Data struct {
	reader      *reader
	tagMap      ifds.TagMap
	make        string
	model       string
	width       uint16
	height      uint16
	exifVersion uint16
	imageType   imagetype.ImageType
}

// GetTag returns a tag from Exif and returns an error if tag doesn't exist
func (e *Data) GetTag(ifd ifds.IfdType, ifdIndex uint8, tagID tag.ID) (tag.Tag, error) {
	if t, ok := e.tagMap[ifds.NewKey(ifd, ifdIndex, tagID)]; ok {
		return t, nil
	}
	return tag.Tag{}, ErrEmptyTag
}

// RangeTags returns a chan tag.Tag for the
// ranging over tags in exif.Data
func (e *Data) RangeTags() chan tag.Tag {
	c := make(chan tag.Tag)
	go func() {
		for _, t := range e.tagMap {
			c <- t
		}
		close(c)
	}()
	return c
}

// GetTagValue returns the tag's value as an interface.
//
// For performance reasons its preferable to use the Parse* functions.
func (e *Data) GetTagValue(t tag.Tag) (value interface{}) {
	asciiLimit := 64 // Limit ascii values to length

	switch t.Type() {
	case tag.TypeASCII, tag.TypeASCIINoNul, tag.TypeByte:
		str, _ := e.ParseASCIIValue(t)
		if len(str) > asciiLimit {
			value = str[:asciiLimit]
		} else {
			value = str
		}
	case tag.TypeShort:
		if t.UnitCount > 1 {
			value, _ = e.ParseUint16Values(t)
		} else {
			value, _ = e.ParseUint16Value(t)
		}
	case tag.TypeLong:
		if t.UnitCount > 1 {
			value, _ = e.ParseUint32Values(t)
		} else {
			value, _ = e.ParseUint32Value(t)
		}
	case tag.TypeRational:
		value, _ = e.ParseRationalValues(t)
	case tag.TypeSignedRational:
		value, _ = e.ParseSRationalValues(t)
	}
	return
}
