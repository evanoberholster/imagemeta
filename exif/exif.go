// Package exif provides functions for parsing and extracting Exif Information.
package exif

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sort"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
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

// ParseExif parses Exif metadata from an io.ReaderAt and a TiffHeader and
// returns exif and an error.
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
func ParseExif(r io.ReaderAt, header meta.ExifHeader) (e *Data, err error) {
	e, err = e.ParseExif(r, header)
	return e, err
}

// ParseExif parses Exif metadata from an io.ReaderAt and a TiffHeader
//
// If the header is invalid ParseExif will return ErrInvalidHeader.
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

// ParseExifWithMetadata parses Exif Metdata from a meta.Reader and meta.Metadata
func (e *Data) ParseExifWithMetadata(r meta.Reader, m *meta.Metadata) (*Data, error) {
	var err error
	e, err = e.ParseExif(r, m.ExifHeader)
	if m.Dim == 0 {
		if e.width != 0 && e.height != 0 {
			m.Dim = e.Dimensions()
		}
	}
	if m.It == imagetype.ImageTiff {
		m.It = e.imageType
	}
	return e, err
}

// Data struct contains the parsed Exif information
type Data struct {
	er        *reader
	tagMap    ifds.TagMap
	make      string
	model     string
	width     uint16
	height    uint16
	imageType imagetype.ImageType
}

// newData creates a new initialized Exif object
func newData(er *reader, it imagetype.ImageType) *Data {
	return &Data{
		er:        er,
		imageType: it,
		tagMap:    make(ifds.TagMap, 20),
	}
}

// AddTag adds a Tag to a tag.TagMap
func (e *Data) addTag(ifd ifds.IFD, ifdIndex uint8, t tag.Tag) {
	switch ifd {
	case ifds.RootIFD:
		switch t.ID {
		case ifds.Make: // Add Make and Model to Exif struct for future decoding of Makernotes
			e.make, _ = e.ParseASCIIValue(t)
		case ifds.Model:
			e.model, _ = e.ParseASCIIValue(t)
		case ifds.ImageWidth:
			if ifdIndex == 0 {
				e.width, _ = e.ParseUint16Value(t)
			}
		case ifds.ImageLength:
			if ifdIndex == 0 {
				e.height, _ = e.ParseUint16Value(t)
			}
		}
	case ifds.ExifIFD:
		switch t.ID {
		case exififd.PixelXDimension:
			e.width, _ = e.ParseUint16Value(t)
		case exififd.PixelYDimension:
			e.height, _ = e.ParseUint16Value(t)
		}
	case ifds.SubIFD, ifds.MknoteIFD, ifds.GPSIFD:

	default:
		// trace UnknownIFD
		return
	}
	e.tagMap[ifds.NewKey(ifd, ifdIndex, t.ID)] = t
}

// GetTag returns a tag from Exif and returns an error if tag doesn't exist
func (e *Data) GetTag(ifd ifds.IFD, ifdIndex uint8, tagID tag.ID) (tag.Tag, error) {
	if t, ok := e.tagMap[ifds.NewKey(ifd, ifdIndex, tagID)]; ok {
		return t, nil
	}
	return tag.Tag{}, ErrEmptyTag
}

// MarshalJSON implements the JSONMarshaler interface that is used by encoding/json
// This is mostly used for testing and debuging.
func (e *Data) MarshalJSON() ([]byte, error) {
	je := jsonExif{It: e.imageType, Make: e.CameraMake(), Model: e.CameraModel(), Width: e.width, Height: e.height}
	for k, t := range e.tagMap {
		ifd, ifdIndex, _ := k.Val()
		value := e.GetTagValue(t)
		je.addTag(ifd, ifdIndex, t, value)
	}

	return json.Marshal(je)
}

// GetTagValue returns the tag's value as an interface.
//
// For performance reasons its preferable to use the Parse* functions.
func (e *Data) GetTagValue(t tag.Tag) (value interface{}) {
	switch t.Type() {
	case tag.TypeASCII, tag.TypeASCIINoNul, tag.TypeByte:
		str, _ := e.ParseASCIIValue(t)
		if len(str) > 64 {
			value = str[:256]
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

// jsonExif for testing purposes.

type jsonExif struct {
	Ifds   map[string]map[uint8]jsonIfds `json:"Ifds"`
	It     imagetype.ImageType           `json:"ImageType"`
	Make   string
	Model  string
	Width  uint16
	Height uint16
}

type jsonIfds struct {
	Tags []jsonTags `json:"Tags"`
}

func (je *jsonExif) addTag(ifd ifds.IFD, ifdIndex uint8, t tag.Tag, v interface{}) {
	if je.Ifds == nil {
		je.Ifds = make(map[string]map[uint8]jsonIfds)
	}
	ji, ok := je.Ifds[ifd.String()]
	if !ok {
		je.Ifds[ifd.String()] = make(map[uint8]jsonIfds)
		ji = je.Ifds[ifd.String()]
	}
	jm, ok := ji[ifdIndex]
	if !ok {
		ji[ifdIndex] = jsonIfds{make([]jsonTags, 0)}
		jm = ji[ifdIndex]
	}
	jm.insertSorted(jsonTags{Name: ifd.TagName(t.ID), Type: t.Type(), ID: t.ID, Count: t.UnitCount, Value: v})
	je.Ifds[ifd.String()][ifdIndex] = jm
}

func (ji *jsonIfds) insertSorted(e jsonTags) {
	i := sort.Search(len(ji.Tags), func(i int) bool { return ji.Tags[i].ID > e.ID })
	ji.Tags = append(ji.Tags, jsonTags{})
	copy(ji.Tags[i+1:], ji.Tags[i:])
	ji.Tags[i] = e
}

type jsonTags struct {
	ID    tag.ID
	Name  string
	Count uint16
	Type  tag.Type
	Value interface{}
}

func (jt jsonTags) MarshalJSON() ([]byte, error) {
	st := struct {
		ID    string
		Name  string
		Count uint16
		Type  string
		Value interface{} `json:"Val"`
	}{
		ID:    jt.ID.String(),
		Name:  jt.Name,
		Count: jt.Count,
		Type:  jt.Type.String(),
		Value: jt.Value,
	}
	return json.Marshal(st)
}
