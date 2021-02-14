package exif

import (
	"errors"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/imagetype"
)

// API Errors
var (
	ErrEmptyTag = errors.New("error empty tag")
)

// ExifData struct contains the parsed Exif information
type ExifData struct {
	exifReader *reader
	rootIfd    []ifds.TagMap
	subIfd     []ifds.TagMap
	exifIfd    ifds.TagMap
	gpsIfd     ifds.TagMap
	mkNote     ifds.TagMap
	make       string
	model      string
	XMP        []byte
	width      uint16
	height     uint16
	imageType  imagetype.ImageType
}

// newExifData creates a new initialized Exif object
func newExifData(er *reader, it imagetype.ImageType) *ExifData {
	return &ExifData{
		exifReader: er,
		imageType:  it,
	}
}

// SetMetadata sets the imagetype metadata in exif
func (e *ExifData) SetMetadata(m imagemeta.Metadata) {
	// Set Exif Width, Height from Metadata Image Size
	e.width, e.height = m.Size()

	// Set Exif XMP form Metadata XML
	e.XMP = []byte(m.XMP())

}

// AddIfd adds an Ifd to a TagMap
func (e *ExifData) AddIfd(ifd ifds.IFD) {
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
func (e *ExifData) AddTag(ifd ifds.IFD, ifdIndex int, t tag.Tag) {
	if ifd == ifds.RootIFD {
		// Add Make and Model to Exif struct for future decoding of Makernotes
		switch t.TagID {
		case ifds.Make:
			e.make, _ = t.ASCIIValue(e.exifReader)
		case ifds.Model:
			e.model, _ = t.ASCIIValue(e.exifReader)
		}
	}
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

func (e *ExifData) getTagMap(ifd ifds.IFD, ifdIndex int) ifds.TagMap {
	switch ifd {
	case ifds.RootIFD:
		if len(e.rootIfd) > ifdIndex {
			return e.rootIfd[ifdIndex]
		}
	case ifds.SubIFD:
		if len(e.subIfd) > ifdIndex {
			return e.subIfd[ifdIndex]
		}
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
func (e *ExifData) GetTag(ifd ifds.IFD, ifdIndex int, tagID tag.ID) (t tag.Tag, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	if tm := e.getTagMap(ifd, ifdIndex); tm != nil {
		var ok bool
		if t, ok = tm[tagID]; ok {
			return
		}
	}
	err = ErrEmptyTag
	return
}
