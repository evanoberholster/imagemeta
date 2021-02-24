package exif

import (
	"errors"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
)

// API Errors
var (
	ErrEmptyTag = errors.New("error empty tag")
)

// Data struct contains the parsed Exif information
type Data struct {
	er        *reader
	ifdMap    ifds.TagMap
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
		ifdMap:    make(ifds.TagMap, 20),
	}
}

// AddTag adds a Tag to a tag.TagMap
func (e *Data) AddTag(ifd ifds.IFD, ifdIndex uint8, t tag.Tag) {
	if ifd == ifds.RootIFD {
		// Add Make and Model to Exif struct for future decoding of Makernotes
		switch t.TagID {
		case ifds.Make:
			e.make, _ = e.ParseASCIIValue(t)
		case ifds.Model:
			e.model, _ = e.ParseASCIIValue(t)
		}
	}
	switch ifd {
	case ifds.RootIFD, ifds.SubIFD, ifds.ExifIFD, ifds.GPSIFD, ifds.MknoteIFD:
		e.ifdMap[ifds.NewKey(ifd, ifdIndex, t.TagID)] = t
	}
}

// GetTag returns a tag from Exif and returns an error if tag doesn't exist
func (e *Data) GetTag(ifd ifds.IFD, ifdIndex uint8, tagID tag.ID) (tag.Tag, error) {
	if t, ok := e.ifdMap[ifds.NewKey(ifd, ifdIndex, tagID)]; ok {
		return t, nil
	}
	return tag.Tag{}, ErrEmptyTag
}

// RangeTags returns a chan tag.Tag for the
// ranging over tags in exif.Data
func (e *Data) RangeTags() chan tag.Tag {
	c := make(chan tag.Tag)
	go func() {
		for _, t := range e.ifdMap {
			c <- t
		}
		close(c)
	}()
	return c
}
