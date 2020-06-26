package tag

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var tagTypeTests = []struct {
	rawTagType uint16
	tagType    Type
	tagSize    uint32
	tagString  string
}{
	{1, TypeByte, TypeByteSize, "BYTE"},
	{2, TypeASCII, TypeASCIISize, "ASCII"},
	{3, TypeShort, TypeShortSize, "SHORT"},
	{4, TypeLong, TypeLongSize, "LONG"},
	{5, TypeRational, TypeRationalSize, "RATIONAL"},
	{7, TypeUndefined, 0, "UNDEFINED"},
	{9, TypeSignedLong, TypeSignedLongSize, "SLONG"},
	{10, TypeSignedRational, TypeSignedRationalSize, "SRATIONAL"},
	{0xf0, TypeASCIINoNul, TypeASCIINoNulSize, "_ASCII_NO_NUL"},
}

func TestTypeFromRaw(t *testing.T) {
	for _, tag := range tagTypeTests {
		t.Run(tag.tagType.String(), func(t *testing.T) {
			ty := TypeFromRaw(tag.rawTagType)
			if ty != tag.tagType {
				t.Errorf("Incorrect Tag Type wanted %s got %s", tag.tagType, ty)
			}
		})
	}
}

func TestTagTypeSizeAndString(t *testing.T) {
	for _, tag := range tagTypeTests {
		t.Run(tag.tagType.String(), func(t *testing.T) {
			var s uint32
			if tag.tagType != TypeUndefined {
				s = tag.tagType.Size()
			} else {
				s = 0
			}
			if s != tag.tagSize {
				t.Errorf("Incorrect Tag Type %s wanted %d got %d", tag.tagType, tag.tagSize, s)
			}
			str := tag.tagType.String()
			if str != tag.tagString {
				t.Errorf("Incorrect Tag Type %s string wanted %s got %s", tag.tagType, tag.tagString, str)
			}
		})
	}
}

func TestTag(t *testing.T) {
	rawValueOffset := RawValueOffset([4]byte{0, 0, 0, 2})
	valueOffset := binary.BigEndian.Uint32(rawValueOffset[:])
	tag := NewTag(ID(0x0000), TypeASCII, 1, valueOffset, rawValueOffset)

	if tag.TagID != ID(0x0000) {
		t.Errorf("Incorrect Tag ID wanted 0x%04x got 0x%04x", ID(0x0000), tag.TagID)
	}
	if tag.TagType != TypeASCII {
		t.Errorf("Incorrect Tag Type wanted %s got %s", TypeASCII, tag.TagType)
	}
	if tag.UnitCount != 1 {
		t.Errorf("Incorrect Tag UnitCount wanted %d got %d", 1, tag.UnitCount)
	}
	if tag.Offset() != 0x0002 {
		t.Errorf("Incorrect Tag Offset wanted 0x%04x got 0x%04x", 0x0002, tag.Offset())
	}
	if !bytes.Equal(tag.rawValueOffset[:], rawValueOffset[:]) {
		t.Errorf("Incorrect Tag Type wanted %b got %b", []byte{0, 0, 0, 2}, tag.rawValueOffset[:])
	}

}
