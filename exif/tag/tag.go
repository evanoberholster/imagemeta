// Package tag provides types and functions for decoding Exif Tags
package tag

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrEmptyTag      = errors.New("error empty tag")
	ErrTagNotValid   = errors.New("error tag not valid")
	ErrNotEnoughData = errors.New("error not enough data to parse tag")
)

// ID is the uint16 representation of an IFD tag
type ID uint16

func (id ID) String() string {
	return fmt.Sprintf("0x%04x", uint16(id))
}

// Offset for reading data
type Offset uint32

// Rational is a rational value
type Rational struct {
	Numerator   uint32
	Denominator uint32
}

// SRational is a signed rational value
type SRational struct {
	Numerator   int32
	Denominator int32
}

// Tag is an Exif Tag
type Tag struct {
	ValueOffset uint32
	UnitCount   uint16 // 4 bytes
	ID          ID     // 2 bytes
	t           Type   // 1 byte
	Ifd         uint8  // 1 byte
}

// NewTag returns a new Tag from tagID, tagType, unitCount, valueOffset and rawValueOffset
func NewTag(tagID ID, tagType Type, unitCount uint32, valueOffset uint32, ifd uint8) Tag {
	return Tag{
		ID:          tagID,
		t:           tagType,
		UnitCount:   uint16(unitCount),
		ValueOffset: valueOffset,
		Ifd:         ifd,
	}
}

func (t Tag) String() string {
	return fmt.Sprintf("0x%04x\t | %s ", uint16(t.ID), t.t)
}

// IsEmbedded checks if the Tag's value is embedded in the Tag.ValueOffset
func (t Tag) IsEmbedded() bool {
	return t.Size() <= 4
}

// Size returns the size of the Tag's value
func (t Tag) Size() int {
	return int(t.t.Size() * uint32(t.UnitCount))
}

// Type returns the type of Tag
func (t Tag) Type() Type {
	return t.t
}

// Errors
var (
	ErrTagTypeNotValid = errors.New("Tag type not valid")
)

// Type is the type of Tag
type Type uint8

// TagTypes defined
// Copied from dsoprea/go-exif
const (
	TypeUnknown Type = 0
	// TypeByte describes an encoded list of bytes.
	TypeByte Type = 1

	// TypeASCII describes an encoded list of characters that is terminated
	// with a NUL in its encoded form.
	TypeASCII Type = 2

	// TypeShort describes an encoded list of shorts.
	TypeShort Type = 3

	// TypeLong describes an encoded list of longs.
	TypeLong Type = 4

	// TypeRational describes an encoded list of rationals.
	TypeRational Type = 5

	// TypeUndefined describes an encoded value that has a complex/non-clearcut
	// interpretation.
	TypeUndefined Type = 7

	// TypeSignedShort describes an encoded list of signed shorts. (experimental)
	TypeSignedShort Type = 8

	// TypeSignedLong describes an encoded list of signed longs.
	TypeSignedLong Type = 9

	// TypeSignedRational describes an encoded list of signed rationals.
	TypeSignedRational Type = 10

	// TypeASCIINoNul is just a pseudo-type, for our own purposes.
	TypeASCIINoNul Type = 0xf0
)

// Tag sizes
const (
	TypeByteSize           = 1
	TypeASCIISize          = 1
	TypeASCIINoNulSize     = 1
	TypeShortSize          = 2
	TypeLongSize           = 4
	TypeRationalSize       = 8
	TypeSignedLongSize     = 4
	TypeSignedRationalSize = 8
)

// Size returns the size of one atomic unit of the type.
func (tagType Type) Size() uint32 {
	switch tagType {
	case TypeByte:
		return TypeByteSize
	case TypeASCII, TypeASCIINoNul:
		return TypeASCIISize
	case TypeShort:
		return TypeShortSize
	case TypeLong:
		return TypeLongSize
	case TypeRational:
		return TypeRationalSize
	case TypeSignedLong:
		return TypeSignedLongSize
	case TypeSignedRational:
		return TypeSignedRationalSize
	default:
		return 0
		//panic(fmt.Errorf("can not determine tag-value size for type (%d): [%s]", tagType, tagType.String()))
	}
}

// IsValid returns true if tagType is a valid type.
func (tagType Type) IsValid() bool {
	return tagType == TypeShort ||
		tagType == TypeLong ||
		tagType == TypeRational ||
		tagType == TypeByte ||
		tagType == TypeASCII ||
		tagType == TypeASCIINoNul ||
		tagType == TypeSignedLong ||
		tagType == TypeSignedRational ||
		tagType == TypeUndefined
}

// String returns the name of the Tag Type
func (tagType Type) String() string {
	switch tagType {
	case TypeByte:
		return "BYTE"
	case TypeASCII:
		return "ASCII"
	case TypeASCIINoNul:
		return "_ASCII_NO_NUL"
	case TypeShort:
		return "SHORT"
	case TypeLong:
		return "LONG"
	case TypeRational:
		return "RATIONAL"
	case TypeSignedLong:
		return "SLONG"
	case TypeSignedRational:
		return "SRATIONAL"
	case TypeUndefined:
		return "UNDEFINED"
	}
	return "UnknownType"
}

// NewTagType returns a new TagType and returns an error
// if the tag type cannot be determined
func NewTagType(raw uint16) (Type, error) {
	tagType := Type(raw)
	if !tagType.IsValid() {
		return 0, ErrTagTypeNotValid
	}
	return tagType, nil
}
