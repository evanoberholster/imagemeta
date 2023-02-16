// Package tag provides types and functions for decoding Exif Tags
package tag

import (
	"errors"
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
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

// Tag is an Exif Tag (16 bytes)
type Tag struct {
	ValueOffset uint32          // 4 bytes
	UnitCount   uint32          // 4 bytes
	ID          ID              // 2 bytes
	TagType     Type            // 1 byte
	Ifd         uint8           // 1 byte
	IfdIndex    int8            // 1 byte
	ByteOrder   utils.ByteOrder // 1 byte
}

// NewTag returns a new Tag from tagID, tagType, unitCount, valueOffset and rawValueOffset.
// If tagType is Invalid returns ErrTagTypeNotValid
func NewTag(tagID ID, tagType Type, unitCount uint32, valueOffset uint32, ifd uint8, ifdIndex int8, byteOrder utils.ByteOrder) (Tag, error) {
	t := Tag{
		ID:          tagID,
		TagType:     tagType,
		UnitCount:   unitCount,
		ValueOffset: valueOffset,
		Ifd:         ifd,
		IfdIndex:    ifdIndex,
		ByteOrder:   byteOrder,
	}
	if !tagType.IsValid() {
		return t, ErrTagTypeNotValid
	}
	return t, nil
}

// MarshalZerologObject is a zerolog interface for logging
func (t Tag) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", t.ID.String()).Str("type", t.TagType.String()).Uint32("count", t.UnitCount).Str("offset", fmt.Sprintf("0x%04x", t.ValueOffset))
}

// String is the Stringer interface for Tag
func (t Tag) String() string {
	return fmt.Sprintf("%s\t | %s | Size: %d", t.ID, t.TagType, t.UnitCount)
}

// EmbeddedValue fills the buf with the tag's embedded value, always <= 4 bytes
func (t Tag) EmbeddedValue(buf []byte) {
	t.ByteOrder.PutUint32(buf, uint32(t.ValueOffset))
}

// IsEmbedded checks if the Tag's value is embedded in the Tag.ValueOffset
func (t Tag) IsEmbedded() bool {
	return t.Size() <= 4 && t.TagType != TypeIfd
}

// IsIfd checks if the Tag's value is an IFD
func (t Tag) IsIfd() bool {
	return t.TagType == TypeIfd
}

// Size returns the size of the Tag's value
func (t Tag) Size() uint32 {
	return uint32(t.TagType.Size()) * uint32(t.UnitCount)
}

// Type returns the type of Tag
func (t Tag) Type() Type {
	return t.TagType
}

// IsType returns true if tagType matches query Type
func (t Tag) IsType(ty Type) bool {
	return t.TagType == ty
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
		tt == TypeSignedLong ||
		tt == TypeSignedRational ||
		tt == TypeFloat ||
		tt == TypeDouble ||
		tt == TypeUndefined ||
		tt == TypeIfd
}
