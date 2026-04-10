package tag

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// IfdType identifies an EXIF/TIFF image file directory kind.
type IfdType uint8

const (
	Unknown IfdType = iota
	IFD0
	IFD1
	IFD2
	ExifIFD
	GPSIFD
	MakerNoteIFD
	SubIFD0
	SubIFD1
	SubIFD2
	SubIFD3
	SubIFD4
	SubIFD5
	SubIFD6
	SubIFD7
)

func (t IfdType) String() string {
	switch t {
	case IFD0:
		return "IFD0"
	case IFD1:
		return "IFD1"
	case IFD2:
		return "IFD2"
	case ExifIFD:
		return "ExifIFD"
	case GPSIFD:
		return "GPSIFD"
	case MakerNoteIFD:
		return "MakerNoteIFD"
	case SubIFD0:
		return "SubIFD0"
	case SubIFD1:
		return "SubIFD1"
	case SubIFD2:
		return "SubIFD2"
	case SubIFD3:
		return "SubIFD3"
	case SubIFD4:
		return "SubIFD4"
	case SubIFD5:
		return "SubIFD5"
	case SubIFD6:
		return "SubIFD6"
	case SubIFD7:
		return "SubIFD7"
	default:
		return "UnknownIFD"
	}
}

func (t IfdType) IsSubIFD() bool {
	return t >= SubIFD0 && t <= SubIFD7
}

func (t IfdType) IsRootIFD() bool {
	return t == IFD0 || t == IFD1 || t == IFD2
}

func (t IfdType) NextRootIFD() (IfdType, bool) {
	switch t {
	case IFD0:
		return IFD1, true
	case IFD1:
		return IFD2, true
	default:
		return Unknown, false
	}
}

func (t IfdType) IsValid() bool {
	return t != Unknown
}

// Directory describes a TIFF/EXIF IFD to read.
type Directory struct {
	Offset     uint32
	BaseOffset uint32
	ByteOrder  utils.ByteOrder
	Type       IfdType
	Index      int8
}

// NewDirectory returns a Directory instance.
func NewDirectory(byteOrder utils.ByteOrder, directoryType IfdType, index int8, ifdOffset uint32, baseOffset uint32) Directory {
	return Directory{
		Offset:     ifdOffset,
		BaseOffset: baseOffset,
		ByteOrder:  byteOrder,
		Type:       directoryType,
		Index:      index,
	}
}

func (d Directory) String() string {
	return fmt.Sprintf("IFD[%s](%d)@0x%04x", d.Type, d.Index, d.Offset)
}
