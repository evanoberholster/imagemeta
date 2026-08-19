package tag

import (
	"errors"
	"strconv"
)

var (
	ErrTagTypeNotValid = errors.New("exif tag type not valid")
)

// ID is the uint16 representation of an EXIF tag identifier.
type ID uint16

func (id ID) String() string {
	return hexUint16Lower(uint16(id))
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

	// BigTIFF (DNG 1.7) 64-bit types.
	TypeLong8       Type = 16
	TypeSignedLong8 Type = 17
	TypeIfd8        Type = 18

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
	TypeLong8Size          = 8
	TypeSignedLong8Size    = 8
	TypeIfd8Size           = 8
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
	TypeLong8:          1,
	TypeSignedLong8:    1,
	TypeIfd8:           1,
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
	case TypeLong8:
		return TypeLong8Size
	case TypeSignedLong8:
		return TypeSignedLong8Size
	case TypeIfd8:
		return TypeIfd8Size
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
	case TypeLong8:
		return "LONG8"
	case TypeSignedLong8:
		return "SLONG8"
	case TypeIfd8:
		return "IFD8"
	default:
		return "UNKNOWN"
	}
}

func (tt Type) IsValid() bool {
	return typeIsValidLookup[uint8(tt)] != 0
}

// RationalU stores an unsigned rational number.
type RationalU struct {
	Numerator   uint32
	Denominator uint32
}

// Float64 converts the rational value into a float64.
func (r RationalU) Float64() float64 {
	if r.Denominator == 0 {
		return 0
	}
	return float64(r.Numerator) / float64(r.Denominator)
}

// UsesIFDType reports whether an ID acts as an IFD pointer in directoryType.
func UsesIFDType(directoryType IfdType, id ID) bool {
	switch id {
	case TagExifIFDPointer, TagGPSIFDPointer:
		return directoryType == IFD0
	case TagMakerNote, TagInteropIFDPointer:
		return directoryType == ExifIFD
	default:
		return false
	}
}

// NormalizeType resolves parser pseudo-types for known IFD pointer tags.
func NormalizeType(directoryType IfdType, id ID, typ Type) Type {
	if (typ.Is(TypeLong) || typ.Is(TypeLong8) || typ.Is(TypeIfd8) || typ.Is(TypeUndefined)) && UsesIFDType(directoryType, id) {
		return TypeIfd
	}
	return typ
}

const hexDigitsLower = "0123456789abcdef"
const hexDigitsUpper = "0123456789ABCDEF"

func hexUint16(v uint16, digits string) string {
	var out [6]byte
	out[0] = '0'
	out[1] = 'x'
	out[2] = digits[(v>>12)&0xF]
	out[3] = digits[(v>>8)&0xF]
	out[4] = digits[(v>>4)&0xF]
	out[5] = digits[v&0xF]
	return string(out[:])
}

func hexUint16Lower(v uint16) string {
	return hexUint16(v, hexDigitsLower)
}

// HexUint16Upper returns 0x-prefixed, fixed-width uppercase hex.
func HexUint16Upper(v uint16) string {
	return hexUint16(v, hexDigitsUpper)
}

func hexUint32LowerMinWidth(v uint32, minWidth int) string {
	digits := strconv.FormatUint(uint64(v), 16)
	if len(digits) >= minWidth {
		return "0x" + digits
	}

	n := 2 + minWidth
	out := make([]byte, n)
	out[0] = '0'
	out[1] = 'x'
	pad := minWidth - len(digits)
	for i := 0; i < pad; i++ {
		out[2+i] = '0'
	}
	copy(out[2+pad:], digits)
	return string(out)
}

// HexUint32LowerMinWidth returns 0x-prefixed lowercase hex padded to minWidth.
func HexUint32LowerMinWidth(v uint32, minWidth int) string {
	return hexUint32LowerMinWidth(v, minWidth)
}
