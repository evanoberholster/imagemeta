package ifd

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// Type identifies an EXIF/TIFF image file directory kind.
type Type uint8

const (
	Unknown Type = iota
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

func (t Type) String() string {
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

func (t Type) IsSubIFD() bool {
	return t >= SubIFD0 && t <= SubIFD7
}

func (t Type) IsRootIFD() bool {
	return t == IFD0 || t == IFD1 || t == IFD2
}

func (t Type) NextRootIFD() (Type, bool) {
	switch t {
	case IFD0:
		return IFD1, true
	case IFD1:
		return IFD2, true
	default:
		return Unknown, false
	}
}

func (t Type) IsValid() bool {
	return t != Unknown
}

// Directory describes a TIFF/EXIF IFD to read.
type Directory struct {
	Offset     uint32
	BaseOffset uint32
	ByteOrder  utils.ByteOrder
	Type       Type
	Index      int8
}

// New returns a Directory instance.
func New(byteOrder utils.ByteOrder, directoryType Type, index int8, ifdOffset uint32, baseOffset uint32) Directory {
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
