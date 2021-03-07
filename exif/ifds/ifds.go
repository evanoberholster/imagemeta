// Package ifds provides types and functions for decoding tiff Ifds
package ifds

import (
	"fmt"

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
	DNGAdobeDataIFD // TODO: Need to implement this
)

func (ifd IFD) String() string {
	switch ifd {
	case RootIFD:
		return "Ifd"
	case SubIFD:
		return "Ifd/SubIfd"
	case ExifIFD:
		return "Ifd/Exif"
	case GPSIFD:
		return "Ifd/GPS"
	case IopIFD:
		return "Ifd/Iop"
	case MknoteIFD:
		return "Ifd/Exif/Makernote"
	case DNGAdobeDataIFD:
		return "Ifd/DNGAdobeData"
	}
	return "UnknownIfd"
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
		name = fmt.Sprintf("0x%04x", id)
	}
	return
}
