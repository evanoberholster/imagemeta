package tag

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

// Entry is a decoded EXIF tag header.
//
// Layout is intentionally compact and friendly for fixed-size parser queues.
type Entry struct {
	ValueOffset uint32
	UnitCount   uint32
	ID          ID
	Type        Type
	IfdType     ifd.Type
	IfdIndex    int8
	ByteOrder   utils.ByteOrder
}

var embeddedMaxUnitsByType = [256]uint8{
	TypeByte:        4,
	TypeASCII:       4,
	TypeUndefined:   4,
	TypeShort:       2,
	TypeSignedShort: 2,
	TypeLong:        1,
	TypeSignedLong:  1,
	TypeFloat:       1,
	TypeASCIINoNul:  4,
}

// NewEntry returns a new tag Entry.
func NewEntry(id ID, typ Type, unitCount, valueOffset uint32, directoryType ifd.Type, ifdIndex int8, byteOrder utils.ByteOrder) Entry {
	return Entry{
		ValueOffset: valueOffset,
		UnitCount:   unitCount,
		ID:          id,
		Type:        typ,
		IfdType:     directoryType,
		IfdIndex:    ifdIndex,
		ByteOrder:   byteOrder,
	}
}

func (t Entry) Name() string {
	return NameFor(t.IfdType, t.ID)
}

func (t Entry) Size() uint32 {
	return uint32(t.Type.Size()) * t.UnitCount
}

func (t Entry) IsEmbedded() bool {
	maxUnits := embeddedMaxUnitsByType[uint8(t.Type)]
	return maxUnits != 0 && t.UnitCount <= uint32(maxUnits)
}

func (t Entry) IsType(tt Type) bool {
	return t.Type == tt
}

func (t Entry) IsIfd() bool {
	return t.Type == TypeIfd
}

func (t Entry) IsValid() bool {
	return t.Type.IsValid()
}

// EmbeddedValue writes the packed value bytes (up to 4 bytes) into dst.
func (t Entry) EmbeddedValue(dst []byte) {
	t.ByteOrder.PutUint32(dst, t.ValueOffset)
}

// EmbeddedShort returns the first embedded SHORT value from ValueOffset.
func (t Entry) EmbeddedShort() uint16 {
	if t.ByteOrder == utils.BigEndian {
		return uint16(t.ValueOffset >> 16)
	}
	return uint16(t.ValueOffset)
}

// EmbeddedShorts decodes embedded SHORT-like values into dst and returns count.
func (t Entry) EmbeddedShorts(dst []uint16) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}

	n := int(t.UnitCount)
	if n > 2 {
		n = 2
	}
	if n > len(dst) {
		n = len(dst)
	}
	if n == 0 {
		return 0
	}

	if t.ByteOrder == utils.BigEndian {
		dst[0] = uint16(t.ValueOffset >> 16)
		if n > 1 {
			dst[1] = uint16(t.ValueOffset)
		}
		return n
	}

	dst[0] = uint16(t.ValueOffset)
	if n > 1 {
		dst[1] = uint16(t.ValueOffset >> 16)
	}
	return n
}

// EmbeddedLong returns the embedded LONG/IFD value from ValueOffset.
func (t Entry) EmbeddedLong() uint32 {
	return t.ValueOffset
}

// ChildDirectory resolves known child-IFD pointers for this tag.
func (t Entry) ChildDirectory() ifd.Directory {
	switch t.IfdType {
	case ifd.IFD0:
		switch t.ID {
		case TagExifIFDPointer:
			return ifd.New(t.ByteOrder, ifd.ExifIFD, t.IfdIndex, t.ValueOffset, 0)
		case TagGPSIFDPointer:
			return ifd.New(t.ByteOrder, ifd.GPSIFD, t.IfdIndex, t.ValueOffset, 0)
		case TagNextIFD:
			return ifd.New(t.ByteOrder, ifd.IFD1, t.IfdIndex+1, t.ValueOffset, 0)
		}
	case ifd.IFD1:
		if t.ID == TagNextIFD {
			return ifd.New(t.ByteOrder, ifd.IFD2, t.IfdIndex+1, t.ValueOffset, 0)
		}
	case ifd.ExifIFD:
		if t.ID == TagMakerNote {
			return ifd.New(t.ByteOrder, ifd.MakerNoteIFD, t.IfdIndex, t.ValueOffset, 0)
		}
	case ifd.SubIFD0, ifd.SubIFD1, ifd.SubIFD2, ifd.SubIFD3, ifd.SubIFD4, ifd.SubIFD5, ifd.SubIFD6, ifd.SubIFD7:
		return ifd.New(t.ByteOrder, t.IfdType, t.IfdIndex, t.ValueOffset, 0)
	}
	return ifd.New(t.ByteOrder, ifd.Unknown, t.IfdIndex, t.ValueOffset, 0)
}

// MarshalZerologObject implements zerolog object marshaling.
func (t Entry) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("id", t.ID).
		Str("name", t.Name()).
		Stringer("type", t.Type).
		Str("ifd", t.IfdType.String()).
		Uint32("units", t.UnitCount).
		Str("offset", fmt.Sprintf("0x%04x", t.ValueOffset))
}
