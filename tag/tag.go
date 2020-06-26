// Package tag provides types and functions for decoding Exif Tags
package tag

import (
	"encoding/binary"
	"fmt"
)

// ID is the uint16 representation of an IFD tag
type ID uint16

// Offset for reading data
type Offset uint32

// RawValueOffset are the 4 bytes after the valueOffset
type RawValueOffset [4]byte

// Tag is an Exif Tag
type Tag struct {
	UnitCount      uint32         // 4 bytes
	valueOffset    uint32         // 4 bytes
	rawValueOffset RawValueOffset // 4 bytes
	TagID          ID             // 2 bytes
	TagType        Type           // 1 byte
}

// TagReader interface
// Implements the io.ReaderAt interface as well as a ByteOrder interface
type TagReader interface {
	ByteOrder() binary.ByteOrder
	ReadAt(p []byte, off int64) (n int, err error)
}

// NewTag creates a new Tag from
func NewTag(tagID ID, tagType Type, unitCount uint32, valueOffset uint32, rawValueOffset RawValueOffset) Tag {
	return Tag{
		TagID:          tagID,
		TagType:        tagType,
		UnitCount:      unitCount,
		valueOffset:    valueOffset,
		rawValueOffset: rawValueOffset,
	}
}

// Offset returns the tag's valueOffset as a uint32 Offset
func (t Tag) Offset() uint32 {
	return t.valueOffset
}

func (t Tag) String() string {
	return fmt.Sprintf("0x%04x \t | %s ", t.TagID, t.TagType)
}
