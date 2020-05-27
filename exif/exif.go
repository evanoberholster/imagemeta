package exif

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/tag"
)

// API Errors
var (
	ErrEmptyTag = errors.New("Error empty tag")
)

type Reader interface {
	io.ReaderAt
	ByteOrder() binary.ByteOrder
}

// Exif struct contains the parsed Exif information
type Exif struct {
	exifReader Reader
	rootIfd    []ifds.TagMap
	subIfd     []ifds.TagMap
	exifIfd    ifds.TagMap
	gpsIfd     ifds.TagMap
	mkNote     ifds.TagMap
	Make       string
	Model      string
}

// NewExif creates a new initalized Exif object
func NewExif(exifReader Reader) *Exif {
	return &Exif{
		exifReader: exifReader,
		//exifIfd:    make(ifds.TagMap),
		//gpsIfd:     make(ifds.TagMap),
		//mkNote:     make(ifds.TagMap),
	}
}

// AddIfd adds an Ifd to a TagMap
func (e *Exif) AddIfd(ifd ifds.IFD) {
	switch ifd {
	case ifds.RootIFD:
		e.rootIfd = append(e.rootIfd, make(ifds.TagMap))
	case ifds.SubIFD:
		e.subIfd = append(e.subIfd, make(ifds.TagMap))
	case ifds.ExifIFD:
		if e.exifIfd == nil {
			e.exifIfd = make(ifds.TagMap)
		}
	case ifds.GPSIFD:
		if e.gpsIfd == nil {
			e.gpsIfd = make(ifds.TagMap)
		}
	case ifds.MknoteIFD:
		if e.mkNote == nil {
			e.mkNote = make(ifds.TagMap)
		}
	}
}

// AddTag adds a Tag to an Ifd -> IfdIndex -> tag.TagMap
func (e *Exif) AddTag(ifd ifds.IFD, ifdIndex int, t tag.Tag) {
	switch ifd {
	case ifds.RootIFD:
		e.rootIfd[ifdIndex][t.TagID] = t
	case ifds.SubIFD:
		e.subIfd[ifdIndex][t.TagID] = t
	case ifds.ExifIFD:
		e.exifIfd[t.TagID] = t
	case ifds.GPSIFD:
		e.gpsIfd[t.TagID] = t
	case ifds.MknoteIFD:
		e.mkNote[t.TagID] = t
	}
}

func (e *Exif) getTagMap(ifd ifds.IFD, ifdIndex int) ifds.TagMap {
	switch ifd {
	case ifds.RootIFD:
		return e.rootIfd[ifdIndex]
	case ifds.SubIFD:
		return e.subIfd[ifdIndex]
	case ifds.ExifIFD:
		return e.exifIfd
	case ifds.GPSIFD:
		return e.gpsIfd
	case ifds.MknoteIFD:
		return e.mkNote
	}
	return nil
}

// GetTag returns a tag from Exif and returns an error if tag doesn't exist
func (e *Exif) GetTag(ifd ifds.IFD, ifdIndex int, tagID tag.ID) (t tag.Tag, err error) {
	if tm := e.getTagMap(ifd, ifdIndex); tm != nil {
		var ok bool
		if t, ok = tm[tagID]; ok {
			return
		}
	}
	err = ErrEmptyTag
	return
}
