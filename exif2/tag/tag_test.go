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
	{11, TypeFloat, TypeFloatSize, "FLOAT", nil},
	{12, TypeDouble, TypeDoubleSize, "DOUBLE", nil},
	{0xf0, TypeASCIINoNul, TypeASCIINoNulSize, "_ASCII_NO_NUL", nil},
	{0xf1, TypeIfd, TypeIfdSize, "IFD", nil},
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

func TestTagID(t *testing.T) {
	id := ID(0x0010)
	if id.String() != "0x0010" {
		t.Errorf("Incorrect ID String wanted %v got %v", "0x0010", id.String())
	}
}

func TestType(t *testing.T) {
	if !TypeASCII.Is(TypeASCII) {
		t.Errorf("Incorrect Type should be equal")
	}
}
