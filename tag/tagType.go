package tag

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrTagTypeNotValid = errors.New("Tag type not valid")
)

// Type is the type of Tag
type Type uint8

// Rational is a rational value
type Rational struct {
	Numerator   uint32
	Denominator uint32
}

// SignedRational is a signed rational value
type SignedRational struct {
	Numerator   int32
	Denominator int32
}

// TagTypes defined
// Copied from dsoprea/go-exif
const (
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
		panic(fmt.Errorf("Can not determine tag-value size for type (%d): [%s]", tagType, tagType.String()))
	}
}

// IsValid returns true if tagType is a valid type.
func (tagType Type) IsValid() bool {
	return tagType == TypeByte ||
		tagType == TypeASCII ||
		tagType == TypeASCIINoNul ||
		tagType == TypeShort ||
		tagType == TypeLong ||
		tagType == TypeRational ||
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

// TypeFromRaw returns the Type of the Tag or panics
// if the tag type cannot be determined
func TypeFromRaw(tagTypeRaw uint16) Type {
	tagType := Type(tagTypeRaw)
	if tagType.IsValid() == false {
		panic(ErrTagTypeNotValid)
	}
	return tagType
}
