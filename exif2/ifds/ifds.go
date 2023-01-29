// Package ifds provides types and functions for decoding tiff Ifds
package ifds

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// IfdType is the Type of Information Directory
type IfdType uint8

// Key is a TagMap Key
type Key struct {
	Type  IfdType
	Index uint8
	TagID tag.ID
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

// IsValid returns true if IFD is valid
func (ifdType IfdType) IsValid() bool {
	return ifdType != NullIFD && int(ifdType) < len(_IFDStringerIndex)-1
}

// String is a stringer interface for ifdType
func (ifdType IfdType) String() string {
	if int(ifdType) < len(_IFDStringerIndex)-1 {
		return _IFDStringerString[_IFDStringerIndex[ifdType]:_IFDStringerIndex[ifdType+1]]
	}
	return NullIFD.String()
}

// TagName returns the tagName for the given IFD and tag.ID
// if tag name is not known returns uint32 representation
func (ifdType IfdType) TagName(id tag.ID) string {
	switch ifdType {
	case IFD0, SubIFD:
		return TagString(id)
	case ExifIFD:
		return exififd.TagString(id)
	case GPSIFD:
		return gpsifd.TagString(id)
	case MknoteIFD:
		//case MkNoteCanonIFD:
		return mknote.TagCanonString(id)
	}

	return id.String()
}

// Ifd is a Tiff Information directory. Contains Offset, Type, and Index.
type Ifd struct {
	Offset    tag.Offset
	ByteOrder utils.ByteOrder
	Type      IfdType
	Index     int8
}

// NewIFD returns a new IFD from IfdType, index, and offset.
func NewIFD(byteOrder utils.ByteOrder, ifdType IfdType, index int8, offset tag.Offset) Ifd {
	return Ifd{
		ByteOrder: byteOrder,
		Type:      ifdType,
		Offset:    offset,
		Index:     index,
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

// String is a stringer interface for ifd
func (ifd Ifd) String() string {
	return fmt.Sprintf("IFD [%s] (%d) at offset (0x%04x)", ifd.Type, ifd.Index, ifd.Offset)
}

// ChildIfd returns the Ifd if it is a Child of the current Tag
// if it is not, it returns NullIFD
func ChildIfd(t tag.Tag) Ifd {
	// RootIfd Children
	switch IfdType(t.Ifd) {
	case IFD0: // IFD0 Children
		switch t.ID {
		case ExifTag:
			return NewIFD(t.ByteOrder, ExifIFD, t.IfdIndex, t.ValueOffset)
		case GPSTag:
			return NewIFD(t.ByteOrder, GPSIFD, t.IfdIndex, t.ValueOffset)
		case SubIFDs:
			return NewIFD(t.ByteOrder, SubIFD, t.IfdIndex, t.ValueOffset)
		}
	case ExifIFD: // ExifIfd Children
		switch t.ID {
		case exififd.MakerNote:
			return NewIFD(t.ByteOrder, MknoteIFD, t.IfdIndex, t.ValueOffset)
		}
	}

	return NewIFD(t.ByteOrder, NullIFD, t.IfdIndex, t.ValueOffset)
}
