// Package tag provides types and functions for decoding Exif Tags
package tag

import (
	"encoding/binary"
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

// Offset for reading data
type Offset uint32

// RawValueOffset are the 4 bytes after the valueOffset
type RawValueOffset [4]byte

// Tag is an Exif Tag
type Tag struct {
	ValueOffset uint32
	UnitCount   uint16 // 4 bytes
	TagID       ID     // 2 bytes
	TagType     Type   // 1 byte
}

// TagReader interface
// Implements the io.ReaderAt interface as well as a ByteOrder interface
type TagReader interface {
	ByteOrder() binary.ByteOrder
	ReadAt(p []byte, off int64) (n int, err error)
}

// NewTag returns a new Tag from tagID, tagType, unitCount, valueOffset and rawValueOffset
func NewTag(tagID ID, tagType Type, unitCount uint32, valueOffset uint32) Tag {
	return Tag{
		TagID:       tagID,
		TagType:     tagType,
		UnitCount:   uint16(unitCount),
		ValueOffset: valueOffset,
	}
}

func (t Tag) String() string {
	return fmt.Sprintf("0x%04x \t | %s ", t.TagID, t.TagType)
}

// IsEmbedded checks if the Tag's value is embedded in the Tag.ValueOffset
func (t Tag) IsEmbedded() bool {
	return t.Size() <= 4
}

// Size returns the size of the Tag's value
func (t Tag) Size() int {
	return int(t.TagType.Size() * uint32(t.UnitCount))
}
