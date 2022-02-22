// Package ifds provides types and functions for decoding tiff Ifds
package ifds

import (
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// IFD is an Information Directory
type IFD uint8

// Key is a TagMap Key
type Key uint32

// NewKey returns a new TagMap Key
func NewKey(ifd IFD, ifdIndex uint8, tagID tag.ID) Key {
	var key uint32
	key |= (uint32(ifd) << 24)
	key |= (uint32(ifdIndex) << 16)
	key |= (uint32(tagID))
	return Key(key)
}

// Val returns the TagMap's Key as an ifd, ifdIndex and a tagID
func (k Key) Val() (ifd IFD, ifdIndex uint8, tagID tag.ID) {
	return IFD(k >> 24), uint8(k << 8 >> 24), tag.ID(k << 16 >> 16)
}

// TagMap is a map of Tags
type TagMap map[Key]tag.Tag

// List of IFDs
const (
	NullIFD IFD = iota
	RootIFD
	SubIFD
	ExifIFD
	GPSIFD
	IopIFD
	MknoteIFD
	DNGAdobeDataIFD

	// IFD Stringer String
	_IFDStringerString = "UnknownIfdIfdIfd/SubIfdIfd/ExifIfd/GPSIfd/IopIfd/Exif/MakernoteIfd/DNGAdobeData"
)

var (
	// IFD Stringer Index
	_IFDStringerIndex = [...]uint8{0, 10, 13, 23, 31, 38, 45, 63, 79}
)

func (ifd IFD) String() string {
	if int(ifd) < len(_IFDStringerIndex)-1 {
		return _IFDStringerString[_IFDStringerIndex[ifd]:_IFDStringerIndex[ifd+1]]
	}
	return IFD(0).String()
}

// IsChildIfd returns the IFD if it is a Child of the current ifd
// if it is not, it returns NullIFD
func (ifd IFD) IsChildIfd(t tag.Tag) IFD {

	// RootIfd Children
	if ifd == RootIFD {
		switch t.ID {
		case ExifTag:
			return ExifIFD
		case GPSTag:
			return GPSIFD
		case SubIFDs:
			return SubIFD
		}
	}

	// ExifIfd Children
	if ifd == ExifIFD {
		switch t.ID {
		case exififd.MakerNote:
			return MknoteIFD
		}
	}

	return NullIFD
}

// TagName returns the tagName for the given IFD and tag.ID
// if tag name is not known returns uint32 representation
func (ifd IFD) TagName(id tag.ID) (name string) {
	var ok bool
	switch ifd {
	case RootIFD, SubIFD:
		name, ok = RootIfdTagIDMap[id]
	case ExifIFD:
		name, ok = exififd.TagIDMap[id]
	case GPSIFD:
		name, ok = gpsifd.TagIDMap[id]
	case MknoteIFD:
	}
	if !ok {
		name = id.String()
	}
	return
}
