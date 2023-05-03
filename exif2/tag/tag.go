// Package tag provides types and functions for decoding Exif Tags
package tag

import (
	"errors"
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

// Errors
var (
	ErrEmptyTag        = errors.New("error empty tag")
	ErrNotEnoughData   = errors.New("error not enough data to parse tag")
	ErrTagTypeNotValid = errors.New("error tag type not valid")
)

// ID is the uint16 representation of an IFD tag
type ID uint16

// String is the Stringer interface for ID
func (id ID) String() string {
	return fmt.Sprintf("0x%04x", uint16(id))
}

// Is returns true if tagType matches query Type
func (tt Type) Is(t Type) bool {
	return tt == t
}

// Type is the type of Tag
type Type uint8

// TagTypes defined
// Copied from dsoprea/go-exif
const (
	// TypeUnknown is an unknown TagType.
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

	// TypeFloat describes an encoded float (float32).
	TypeFloat Type = 11

	// TypeDouble describes an emcoded double (uint64).
	TypeDouble Type = 12

	// PseudoTypes

	// TypeASCIINoNul is just a pseudo-type, for our own purposes.
	TypeASCIINoNul Type = 0xf0

	// TypeIfd is a pseudo-type, for our own purposes.
	TypeIfd Type = 0xf1
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
	TypeFloatSize          = 4
	TypeDoubleSize         = 8
	TypeIfdSize            = 4

	// TagType Stringer String
	_TagTypeStringerString = "UnknownBYTEASCIISHORTLONGRATIONALUnknownUNDEFINEDSSHORTSLONGSRATIONALFLOATDOUBLE"
)

var (
	//Tag sizes
	_tagSize = [256]uint8{
		0, TypeByteSize, TypeASCIISize, TypeShortSize, TypeLongSize, TypeRationalSize, 0, TypeByteSize, TypeShortSize, TypeSignedLongSize,
		TypeSignedRationalSize, TypeFloatSize, TypeDoubleSize, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, TypeASCIINoNulSize,
		TypeIfdSize, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}

	// TagType Stringer Index
	_TagTypeStringerIndex = [...]uint8{0, 7, 11, 16, 21, 25, 33, 40, 49, 55, 60, 69, 74, 80}
)

// Size returns the size of one atomic unit of the type.
func (tt Type) Size() uint8 {
	return _tagSize[uint8(tt)]
}

// String is the stringer interface for TagType
func (tt Type) String() string {
	if int(tt) < len(_TagTypeStringerIndex)-1 {
		return _TagTypeStringerString[_TagTypeStringerIndex[tt]:_TagTypeStringerIndex[tt+1]]
	}
	if tt == TypeIfd {
		return "IFD"
	}
	if tt == TypeASCIINoNul {
		return "_ASCII_NO_NUL"
	}
	return TypeUnknown.String()
}

// IsValid returns true if tagType is a valid type.
func (tt Type) IsValid() bool {
	return tt == TypeShort ||
		tt == TypeLong ||
		tt == TypeRational ||
		tt == TypeByte ||
		tt == TypeASCII ||
		tt == TypeASCIINoNul ||
		tt == TypeSignedShort ||
		tt == TypeSignedLong ||
		tt == TypeSignedRational ||
		tt == TypeFloat ||
		tt == TypeDouble ||
		tt == TypeUndefined ||
		tt == TypeIfd
}

// Rational is a rational value
type Rational [2]uint32

// Num returns the the Rational Numerator as a unit32
func (rat Rational) Num() uint32 {
	return rat[0]
}

// Den returns the Rational Denominator as a unit32
func (rat Rational) Den() uint32 {
	return rat[1]
}

// Float returns the Rational as a float64
func (rat Rational) Float() float64 {
	if rat.Den() == 0 {
		return 0.0
	}
	return float64(rat.Num()) / float64(rat.Den())
}

// SRational is a signed rational value
type SRational [2]int32

type TagValue struct {
	Buf       []byte          // 8 bytes
	UnitCount uint32          // 4 bytes
	ID        ID              // 2 bytes
	Type      Type            // 1 byte
	ByteOrder utils.ByteOrder // 1 byte
}

// parseStrUint parses a []byte of a string representation of a uint value and returns the value.
func parseStrUint(buf []byte) (u uint) {
	for i := 0; i < len(buf); i++ {
		if buf[i] >= '0' {
			u *= 10
			u += uint(buf[i] - '0')
		}
	}
	return
}

// trimNULBuffer removes trailing bytes from Buffer
func trimNULBuffer(buf []byte) []byte {
	for i := len(buf) - 1; i > 0; i-- {
		if buf[i] == 0 || buf[i] == ' ' || buf[i] == '\n' {
			continue
		}
		return buf[:i+1]
	}
	return nil
}
