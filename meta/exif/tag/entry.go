package tag

import (
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// Entry is a decoded EXIF tag header.
//
// Layout is intentionally compact and friendly for fixed-size parser queues.
type Entry struct {
	ValueOffset uint32
	// valueHigh carries the upper 4 bytes of a 5-8 byte value packed inline in a
	// BigTIFF entry's 8-byte value field. Unused (0) for classic TIFF.
	valueHigh uint32
	UnitCount uint32
	ID        ID
	Type      Type
	IfdType   IfdType
	IfdIndex  int8
	// embedded forces the inline-value path regardless of the type-size
	// heuristic; set for BigTIFF values carried in the 8-byte value field.
	embedded  bool
	ByteOrder utils.ByteOrder
}

// NewEntry returns a new tag Entry.
func NewEntry(id ID, typ Type, unitCount, valueOffset uint32, directoryType IfdType, ifdIndex int8, byteOrder utils.ByteOrder) Entry {
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

// NewInlineEntry returns an entry whose value is packed inline in up to 8 bytes
// (BigTIFF), with lo and hi holding the value's low and high 4 bytes. Such an
// entry always reports IsEmbedded, so value parsers read the packed bytes
// rather than treating ValueOffset as a file offset.
func NewInlineEntry(id ID, typ Type, unitCount, lo, hi uint32, directoryType IfdType, ifdIndex int8, byteOrder utils.ByteOrder) Entry {
	return Entry{
		ValueOffset: lo,
		valueHigh:   hi,
		UnitCount:   unitCount,
		ID:          id,
		Type:        typ,
		IfdType:     directoryType,
		IfdIndex:    ifdIndex,
		embedded:    true,
		ByteOrder:   byteOrder,
	}
}

func (t Entry) Name() string {
	return NameFor(t.IfdType, t.ID)
}

func (t Entry) Size() uint32 {
	return uint32(t.Type.Size()) * t.UnitCount
}

// IsEmbedded checks inline-value eligibility for common TIFF/EXIF types.
func (t Entry) IsEmbedded() bool {
	if t.embedded {
		return true
	}
	switch t.Type {
	case TypeByte, TypeASCII, TypeUndefined, TypeASCIINoNul:
		return t.UnitCount <= 4
	case TypeShort, TypeSignedShort:
		return t.UnitCount <= 2
	case TypeLong, TypeSignedLong, TypeFloat:
		return t.UnitCount <= 1
	default:
		return false
	}
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

// EmbeddedValue writes the packed value bytes into dst: the low 4 bytes from
// ValueOffset and, when dst has room, the high 4 bytes from valueHigh (BigTIFF
// 8-byte inline values). Callers pass an 8-byte dst to read the full value.
func (t Entry) EmbeddedValue(dst []byte) {
	t.ByteOrder.PutUint32(dst, t.ValueOffset)
	if len(dst) >= 8 {
		t.ByteOrder.PutUint32(dst[4:8], t.valueHigh)
	}
}

// EmbeddedShort returns the first embedded SHORT value from ValueOffset.
func (t Entry) EmbeddedShort() uint16 {
	if t.ByteOrder == utils.BigEndian {
		return high16(t.ValueOffset)
	}
	return low16(t.ValueOffset)
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
		dst[0] = high16(t.ValueOffset)
		if n > 1 {
			dst[1] = low16(t.ValueOffset)
		}
		return n
	}

	dst[0] = low16(t.ValueOffset)
	if n > 1 {
		dst[1] = high16(t.ValueOffset)
	}
	return n
}

func low16(v uint32) uint16 {
	//nolint:gosec // G115: deliberate lower 16-bit extraction.
	return uint16(v & 0xFFFF)
}

func high16(v uint32) uint16 {
	//nolint:gosec // G115: deliberate upper 16-bit extraction.
	return uint16((v >> 16) & 0xFFFF)
}

// EmbeddedLong returns the embedded LONG/IFD value from ValueOffset.
func (t Entry) EmbeddedLong() uint32 {
	return t.ValueOffset
}

// ChildDirectory resolves known child-IFD pointers for this tag.
func (t Entry) ChildDirectory() Directory {
	switch t.IfdType {
	case IFD0:
		switch t.ID {
		case TagExifIFDPointer:
			return NewDirectory(t.ByteOrder, ExifIFD, t.IfdIndex, t.ValueOffset, 0)
		case TagGPSIFDPointer:
			return NewDirectory(t.ByteOrder, GPSIFD, t.IfdIndex, t.ValueOffset, 0)
		case TagNextIFD:
			return NewDirectory(t.ByteOrder, IFD1, t.IfdIndex+1, t.ValueOffset, 0)
		}
	case IFD1:
		if t.ID == TagNextIFD {
			return NewDirectory(t.ByteOrder, IFD2, t.IfdIndex+1, t.ValueOffset, 0)
		}
	case ExifIFD:
		if t.ID == TagMakerNote {
			return NewDirectory(t.ByteOrder, MakerNoteIFD, t.IfdIndex, t.ValueOffset, 0)
		}
		if t.ID == TagInteropIFDPointer {
			// Parse InteropIFD tags through ExifIFD semantics.
			return NewDirectory(t.ByteOrder, ExifIFD, t.IfdIndex, t.ValueOffset, 0)
		}
	case SubIFD0, SubIFD1, SubIFD2, SubIFD3, SubIFD4, SubIFD5, SubIFD6, SubIFD7:
		return NewDirectory(t.ByteOrder, t.IfdType, t.IfdIndex, t.ValueOffset, 0)
	}
	return NewDirectory(t.ByteOrder, Unknown, t.IfdIndex, t.ValueOffset, 0)
}

// MarshalLogObject implements structured object marshaling.
func (t Entry) MarshalLogObject(e *metalog.Event) {
	e.Stringer("id", t.ID).
		Str("name", t.Name()).
		Stringer("type", t.Type).
		Str("ifd", t.IfdType.String()).
		Uint32("units", t.UnitCount).
		Str("offset", HexUint32LowerMinWidth(t.ValueOffset, 4))
}
