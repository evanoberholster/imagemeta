package ifds

import (
	"fmt"

	"github.com/evanoberholster/exiftool/ifds/exififd"
	"github.com/evanoberholster/exiftool/ifds/gpsifd"
	"github.com/evanoberholster/exiftool/tag"
)

// IFD -
type IFD uint8

// TagMap is a map of Tags
type TagMap map[tag.ID]tag.Tag

// IfdMap is a map of ifds to an array of tagMaps
//type IfdMap map[IFD][]tag.TagMap

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
	return "Unknown"
}

// IsChildIfd returns the IFD if it is a Child of the current ifd
// if it is not, it returns NullIFD
func (ifd IFD) IsChildIfd(t tag.Tag) IFD {

	// RootIfd Children
	if ifd == RootIFD {
		switch t.TagID {
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
		switch t.TagID {
		case exififd.MakerNote:
			return MknoteIFD
		}
	}

	return NullIFD
}

// TagName returns the tagName for the given IFD and tag.ID
// if tag name is not known returns uint32 representation
func (ifd IFD) TagName(id tag.ID) (name string) {
	switch ifd {
	case RootIFD, SubIFD:
		name = RootIfdTagIDMap[id]
	case ExifIFD:
		name = exififd.TagIDMap[id]
	case GPSIFD:
		name = gpsifd.TagIDMap[id]
	case MknoteIFD:
	}
	if name == "" {
		name = fmt.Sprintf("0x%04x", id)
	}
	return
}
