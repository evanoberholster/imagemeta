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
	ValueOffset uint32 // 4 bytes
	UnitCount   uint32 // 4 bytes
	ID          ID     // 2 bytes
	t           Type   // 1 byte
	Ifd         uint8  // 1 byte
}

// NewTag returns a new Tag from tagID, tagType, unitCount, valueOffset and rawValueOffset.
// If tagType is Invalid returns ErrTagTypeNotValid
func NewTag(tagID ID, tagType Type, unitCount uint32, valueOffset uint32, ifd uint8) (Tag, error) {
	if !tagType.IsValid() {
		return Tag{}, ErrTagTypeNotValid
	}
	// Special tags

	return Tag{
		ID:          tagID,
		t:           tagType,
		UnitCount:   unitCount,
		ValueOffset: valueOffset,
		Ifd:         ifd,
	}, nil
}

func (t Tag) String() string {
	return fmt.Sprintf("%s\t | %s ", t.ID, t.t)
}

// IsEmbedded checks if the Tag's value is embedded in the Tag.ValueOffset
func (t Tag) IsEmbedded() bool {
	return t.Size() <= 4
}

// IsIFD checks if the Tag's value is an IFD
func (t Tag) IsIfd() bool {
	return t.t == TypeIfd
}

// Size returns the size of the Tag's value
func (t Tag) Size() uint32 {
	return uint32(t.t.Size()) * uint32(t.UnitCount)
}

// Type returns the type of Tag
func (t Tag) Type() Type {
	return t.t
}

// IsType returns true if tagType matches query Type
func (tag Tag) IsType(t Type) bool {
	return tag.t == t
}

// Is returns true if tagType matches query Type
func (tagType Type) Is(t Type) bool {
	return tagType == t
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
	TypeIfdSize            = 4

	// TagType Stringer String
	_TagTypeStringerString = "UnknownBYTEASCIISHORTLONGRATIONALUnknownUNDEFINEDSSHORTSLONGSRATIONAL"
)

var (
	//Tag sizes
	_tagSize = [...]uint8{0, TypeByteSize, TypeASCIISize, TypeShortSize, TypeLongSize, TypeRationalSize, 0, 0, TypeShortSize, TypeSignedLongSize, TypeSignedRationalSize}

	// TagType Stringer Index
	_TagTypeStringerIndex = [...]uint8{0, 7, 11, 16, 21, 25, 33, 40, 49, 55, 60, 69}
)

// Size returns the size of one atomic unit of the type.
func (tagType Type) Size() uint8 {
	if int(tagType) < len(_tagSize) {
		return uint8(_tagSize[uint8(tagType)])
	}
	if tagType == TypeIfd {
		return TypeIfdSize
	}
	if tagType == TypeASCIINoNul {
		return TypeASCIINoNulSize
	}
	return 0
}

// String returns the name of the Tag Type
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
		tt == TypeSignedLong ||
		tt == TypeSignedRational ||
		tt == TypeUndefined ||
		tt == TypeIfd
}
