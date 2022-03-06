// Package ifds provides types and functions for decoding tiff Ifds
package ifds

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// IFD is an Information Directory
type IfdType uint8

// Key is a TagMap Key
type Key uint32

// NewKey returns a new TagMap Key
func NewKey(ifdType IfdType, ifdIndex uint8, tagID tag.ID) Key {
	var key uint32
	key |= (uint32(ifdType) << 24)
	key |= (uint32(ifdIndex) << 16)
	key |= (uint32(tagID))
	return Key(key)
}

// Val returns the TagMap's Key as an ifd, ifdIndex and a tagID
func (k Key) Val() (ifdType IfdType, ifdIndex uint8, tagID tag.ID) {
	return IfdType(k >> 24), uint8(k << 8 >> 24), tag.ID(k << 16 >> 16)
}

// TagMap is a map of Tags
type TagMap map[Key]tag.Tag

// List of IFDs
const (
	NullIFD IfdType = iota
	IFD0
	SubIFD
	ExifIFD
	GPSIFD
	IopIFD
	MknoteIFD
	DNGAdobeDataIFD
	MkNoteCanonIFD
	MkNoteNikonIFD

	// IFD Stringer String
	_IFDStringerString = "UnknownIfdIfdIfd/SubIfdIfd/ExifIfd/GPSIfd/IopIfd/Exif/MakernoteIfd/DNGAdobeDataIfd/Exif/MakernoteIfd/Exif/Makernote"
)

var (
	// IFD Stringer Index
	_IFDStringerIndex = [...]uint8{0, 10, 13, 23, 31, 38, 45, 63, 79, 97, 115}
)

// Valid returns true if IFD is valid
func (ifdType IfdType) IsValid() bool {
	return ifdType != NullIFD && int(ifdType) < len(_IFDStringerIndex)-1
}

func (ifdType IfdType) String() string {
	if int(ifdType) < len(_IFDStringerIndex)-1 {
		return _IFDStringerString[_IFDStringerIndex[ifdType]:_IFDStringerIndex[ifdType+1]]
	}
	return NullIFD.String()
}

// TagName returns the tagName for the given IFD and tag.ID
// if tag name is not known returns uint32 representation
func (ifdType IfdType) TagName(id tag.ID) (name string) {
	var ok bool
	switch ifdType {
	case IFD0, SubIFD:
		name, ok = RootIfdTagIDMap[id]
	case ExifIFD:
		name, ok = exififd.TagIDMap[id]
	case GPSIFD:
		name, ok = gpsifd.TagIDMap[id]
	case MknoteIFD:
		//case MkNoteCanonIFD:
		name, ok = mknote.TagCanonIDMap[id]
	}
	if !ok {
		name = id.String()
	}
	return
}

// Ifd is a Tiff Information directory. Contains Offset, Type, and Index.
type Ifd struct {
	Offset uint32
	Type   IfdType
	Index  uint8
}

// NewIFD returns a new IFD from IfdType, index, and offset.
func NewIFD(ifdType IfdType, index uint8, offset uint32) Ifd {
	return Ifd{
		Type:   ifdType,
		Offset: offset,
		Index:  index,
	}
}

// TagName returns the Tagname for the given tag.ID
func (ifd Ifd) TagName(id tag.ID) (name string) {
	return ifd.Type.TagName(id)
}

// IsType returns true if ifdType equals IfdType
func (ifd Ifd) IsType(t IfdType) bool {
	return ifd.Type == t
}

// IsValid returns IfdType.IsValid
func (ifd Ifd) IsValid() bool {
	return ifd.Type.IsValid()
}

func (ifd Ifd) String() string {
	return fmt.Sprintf("IFD [%s] (%d) at offset (0x%04x)", ifd.Type, ifd.Index, ifd.Offset)
}

// ChildIfd returns the Ifd if it is a Child of the current Tag
// if it is not, it returns NullIFD
func (ifd Ifd) ChildIfd(t tag.Tag) Ifd {
	// RootIfd Children
	if ifd.IsType(IFD0) {
		switch t.ID {
		case ExifTag:
			return NewIFD(ExifIFD, 0, t.ValueOffset)
		case GPSTag:
			return NewIFD(GPSIFD, 0, t.ValueOffset)
		case SubIFDs:
			return NewIFD(SubIFD, 0, t.ValueOffset)
		}
	}

	// ExifIfd Children
	if ifd.IsType(ExifIFD) {
		switch t.ID {
		case exififd.MakerNote:
			return NewIFD(MknoteIFD, 0, t.ValueOffset)
		}
	}

	return NewIFD(NullIFD, 0, t.ValueOffset)
}
