package tag

import (
	"testing"
)

var tagTypeTests = []struct {
	rawTagType uint16
	tagType    Type
	tagSize    uint8
	tagString  string
	err        error
}{
	{1, TypeByte, TypeByteSize, "BYTE", nil},
	{2, TypeASCII, TypeASCIISize, "ASCII", nil},
	{3, TypeShort, TypeShortSize, "SHORT", nil},
	{4, TypeLong, TypeLongSize, "LONG", nil},
	{5, TypeRational, TypeRationalSize, "RATIONAL", nil},
	{7, TypeUndefined, 0, "UNDEFINED", nil},
	{9, TypeSignedLong, TypeSignedLongSize, "SLONG", nil},
	{10, TypeSignedRational, TypeSignedRationalSize, "SRATIONAL", nil},
	{0xf0, TypeASCIINoNul, TypeASCIINoNulSize, "_ASCII_NO_NUL", nil},
	{0, TypeUnknown, 0, "Unknown", ErrTagTypeNotValid},
	{100, 100, 0, "Unknown", ErrTagTypeNotValid},
}

func TestNewTagType(t *testing.T) {
	for _, tag := range tagTypeTests {
		t.Run(tag.tagType.String(), func(t *testing.T) {
			ty := Type(tag.rawTagType)
			if ty != tag.tagType {
				if ty.IsValid() {
					t.Errorf("Incorrect Tag Type wanted %s got %s", tag.tagType, ty)
				}
			}
			if !ty.IsValid() && tag.err != ErrTagTypeNotValid {
				t.Errorf("Incorrect err %s", tag.err)
			}
		})
	}
}

func TestTagType(t *testing.T) {
	for _, tag := range tagTypeTests {
		t.Run(tag.tagType.String(), func(t *testing.T) {
			var s uint8
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
	tag, err := NewTag(ID(0x0010), TypeASCII, 16, 0x0002, 0)
	if err != nil {
		t.Error(err)
	}

	if tag.ID != ID(0x0010) {
		t.Errorf("Incorrect Tag ID wanted 0x%04x got 0x%04x", ID(0x0010), tag.ID)
	}
	if tag.Type() != TypeASCII {
		t.Errorf("Incorrect Tag Type wanted %s got %s", TypeASCII, tag.Type())
	}
	if tag.UnitCount != 16 {
		t.Errorf("Incorrect Tag UnitCount wanted %d got %d", 16, tag.UnitCount)
	}
	if tag.ValueOffset != 0x0002 {
		t.Errorf("Incorrect Tag Offset wanted 0x%04x got 0x%04x", 0x0002, tag.ValueOffset)
	}
	if tag.IsEmbedded() {
		t.Errorf("ValueIsEmbedded is true when equal or less than 4 bytes")
	}
	if tag.Size() != 16 {
		t.Errorf("Incorrect Tag Size wanted %d got %d", 16, tag.Size())
	}
	if tag.String() != "0x0010\t | ASCII " {
		t.Errorf("Incorrect Tag String wanted %v got %v", "0x0010\t | ASCII ", tag.String())
	}
	if tag.ID.String() != "0x0010" {
		t.Errorf("Incorrect ID String wanted %v got %v", "0x0010", tag.ID.String())
	}

	tag, err = NewTag(ID(0x0010), 100, 16, 0x0002, 0)
	if err != ErrTagTypeNotValid {
		t.Errorf("Incorrect error wanted %s, got %s", ErrTagTypeNotValid, err)
	}
}
