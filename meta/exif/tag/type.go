package tag

import (
	"errors"
	"fmt"
)

var (
	ErrTagTypeNotValid = errors.New("exif tag type not valid")
)

// ID is the uint16 representation of an EXIF tag identifier.
type ID uint16

func (id ID) String() string {
	return fmt.Sprintf("0x%04x", uint16(id))
}

// Type is the EXIF field type.
type Type uint8

const (
	TypeUnknown        Type = 0
	TypeByte           Type = 1
	TypeASCII          Type = 2
	TypeShort          Type = 3
	TypeLong           Type = 4
	TypeRational       Type = 5
	TypeUndefined      Type = 7
	TypeSignedShort    Type = 8
	TypeSignedLong     Type = 9
	TypeSignedRational Type = 10
	TypeFloat          Type = 11
	TypeDouble         Type = 12

	// Pseudo-types used by parser internals.
	TypeASCIINoNul Type = 0xf0
	TypeIfd        Type = 0xf1
)

const (
	TypeByteSize           = 1
	TypeASCIISize          = 1
	TypeASCIINoNulSize     = 1
	TypeShortSize          = 2
	TypeLongSize           = 4
	TypeRationalSize       = 8
	TypeSignedLongSize     = 4
	TypeSignedRationalSize = 8
	TypeFloatSize          = 4
	TypeDoubleSize         = 8
	TypeIfdSize            = 4
)

var typeIsValidLookup = [256]uint8{
	TypeByte:           1,
	TypeASCII:          1,
	TypeShort:          1,
	TypeLong:           1,
	TypeRational:       1,
	TypeUndefined:      1,
	TypeSignedShort:    1,
	TypeSignedLong:     1,
	TypeSignedRational: 1,
	TypeFloat:          1,
	TypeDouble:         1,
	TypeASCIINoNul:     1,
	TypeIfd:            1,
}

func (tt Type) Is(t Type) bool {
	return tt == t
}

// Size returns the size of one atomic unit for this type.
func (tt Type) Size() uint8 {
	switch tt {
	case TypeByte:
		return TypeByteSize
	case TypeASCII:
		return TypeASCIISize
	case TypeShort:
		return TypeShortSize
	case TypeLong:
		return TypeLongSize
	case TypeRational:
		return TypeRationalSize
	case TypeUndefined:
		return TypeByteSize
	case TypeSignedShort:
		return TypeShortSize
	case TypeSignedLong:
		return TypeSignedLongSize
	case TypeSignedRational:
		return TypeSignedRationalSize
	case TypeFloat:
		return TypeFloatSize
	case TypeDouble:
		return TypeDoubleSize
	case TypeASCIINoNul:
		return TypeASCIINoNulSize
	case TypeIfd:
		return TypeIfdSize
	default:
		return 0
	}
}

func (tt Type) String() string {
	switch tt {
	case TypeByte:
		return "BYTE"
	case TypeASCII:
		return "ASCII"
	case TypeShort:
		return "SHORT"
	case TypeLong:
		return "LONG"
	case TypeRational:
		return "RATIONAL"
	case TypeUndefined:
		return "UNDEFINED"
	case TypeSignedShort:
		return "SSHORT"
	case TypeSignedLong:
		return "SLONG"
	case TypeSignedRational:
		return "SRATIONAL"
	case TypeFloat:
		return "FLOAT"
	case TypeDouble:
		return "DOUBLE"
	case TypeASCIINoNul:
		return "_ASCII_NO_NUL"
	case TypeIfd:
		return "IFD"
	default:
		return "UNKNOWN"
	}
}

func (tt Type) IsValid() bool {
	return typeIsValidLookup[uint8(tt)] != 0
}
