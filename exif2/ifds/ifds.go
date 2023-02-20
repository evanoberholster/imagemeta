// Package ifds provides types and functions for decoding tiff Ifds
package ifds

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/apple"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/canon"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/nikon"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/sony"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

// IfdType is the Type of Information Directory
type IfdType uint8

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
	MkNoteAppleIFD
	MkNoteSonyIFD

	// SubIfds
	SubIfd0
	SubIfd1
	SubIfd2
	SubIfd3
	SubIfd4
	SubIfd5
	SubIfd6
	SubIfd7

	// IFD Stringer String
	_IFDStringerString = "UnknownIfdIfdIfd/SubIfdIfd/ExifIfd/GPSIfd/IopIfd/Exif/MakernoteIfd/DNGAdobeDataIfd/Exif/MakernoteIfd/Exif/MakernoteIfd/Exif/MakernoteIfd/Exif/MakernoteIfd/SubIfd0Ifd/SubIfd1Ifd/SubIfd2Ifd/SubIfd3Ifd/SubIfd4Ifd/SubIfd5Ifd/SubIfd6Ifd/SubIfd7"
)

var (
	// IFD Stringer Index
	_IFDStringerIndex = [...]uint8{0, 10, 13, 23, 31, 38, 45, 63, 79, 97, 115, 133, 151, 162, 173, 184, 195, 206, 217, 228, 239}
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
	case MkNoteCanonIFD:
		return canon.TagCanonString(id)
	case MkNoteNikonIFD:
		return nikon.TagNikonString(id)
	case MkNoteAppleIFD:
		return apple.TagAppleString(id)
	case MkNoteSonyIFD:
		return sony.TagSonyString(id)
	case SubIfd0, SubIfd1, SubIfd2, SubIfd3, SubIfd4, SubIfd5, SubIfd6, SubIfd7:
		return TagSubIfdString(id, ifdType)
	default:
		return id.String()
	}
}

// Ifd is a Tiff Information directory. Contains Offset, Type, and Index.
type Ifd struct {
	Offset     uint32
	BaseOffset uint32
	ByteOrder  utils.ByteOrder
	Type       IfdType
	Index      int8
}

// NewIFD returns a new IFD from IfdType, index, and offset.
func NewIFD(byteOrder utils.ByteOrder, ifdType IfdType, index int8, ifdOffset uint32, baseOffset uint32) Ifd {
	return Ifd{
		ByteOrder:  byteOrder,
		Type:       ifdType,
		Offset:     ifdOffset,
		BaseOffset: baseOffset,
		Index:      index,
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

// MarshalZerologObject is a zerolog interface for logging
func (ifd Ifd) MarshalZerologObject(e *zerolog.Event) {
	e.Str("IfdType", ifd.Type.String()).Int8("idx", ifd.Index).Str("offset", fmt.Sprintf("0x%04x", ifd.Offset))
}
