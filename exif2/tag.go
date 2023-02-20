package exif2

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

// Tag is an Exif Tag (16 bytes)
type Tag struct {
	ValueOffset uint32          // 4 bytes
	UnitCount   uint32          // 4 bytes
	ID          tag.ID          // 2 bytes
	Type        tag.Type        // 1 byte
	Ifd         ifds.IfdType    // 1 byte
	IfdIndex    int8            // 1 byte
	ByteOrder   utils.ByteOrder // 1 byte
}

// NewTag returns a new Tag from tagID, tagType, unitCount, valueOffset and rawValueOffset.
// If tagType is Invalid returns ErrTagTypeNotValid
func NewTag(tagID tag.ID, tagType tag.Type, unitCount uint32, valueOffset uint32, ifd ifds.IfdType, ifdIndex int8, byteOrder utils.ByteOrder) Tag {
	return Tag{
		ID:          tagID,
		Type:        tagType,
		UnitCount:   unitCount,
		ValueOffset: valueOffset,
		Ifd:         ifd,
		IfdIndex:    ifdIndex,
		ByteOrder:   byteOrder,
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (t Tag) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("id", t.ID).Str("name", t.Name()).Stringer("type", t.Type).Stringer("ifd", t.Ifd).Uint32("units", t.UnitCount).Str("offset", fmt.Sprintf("0x%04x", t.ValueOffset))
}

func (t Tag) logTag(e *zerolog.Event) *zerolog.Event {
	t.MarshalZerologObject(e)
	return e
}

// Name returns the Tag name as a string
func (t Tag) Name() string {
	return t.Ifd.TagName(t.ID)
}

// EmbeddedValue fills the buf with the tag's embedded value, always <= 4 bytes
func (t Tag) EmbeddedValue(buf []byte) {
	t.ByteOrder.PutUint32(buf, t.ValueOffset)
}

// IsEmbedded checks if the Tag's value is embedded in the Tag.ValueOffset
func (t Tag) IsEmbedded() bool {
	return t.Size() <= 4 && t.Type != tag.TypeIfd
}

// Size returns the size of the Tag's value
func (t Tag) Size() uint32 {
	return uint32(t.Type.Size()) * uint32(t.UnitCount)
}

// IsIfd checks if the Tag's value is an IFD
func (t Tag) IsIfd() bool {
	return t.Type == tag.TypeIfd
}

// IsType returns true if tagType matches query Type
func (t Tag) IsType(tt tag.Type) bool {
	return t.Type == tt
}

// IsValid returns true if tagType is valid
func (t Tag) IsValid() bool {
	return t.Type.IsValid()
}

// childIfd returns the Ifd if it is a Child of the current Tag
// if it is not, it returns NullIFD
func (t Tag) childIfd() ifds.Ifd {
	switch t.Ifd {
	case ifds.IFD0: // IFD0 Children
		switch t.ID {
		case ifds.ExifTag:
			return ifds.NewIFD(t.ByteOrder, ifds.ExifIFD, t.IfdIndex, t.ValueOffset, 0)
		case ifds.GPSTag:
			return ifds.NewIFD(t.ByteOrder, ifds.GPSIFD, t.IfdIndex, t.ValueOffset, 0)
		}

	case ifds.ExifIFD: // ExifIfd Children
		switch t.ID {
		case exififd.MakerNote:
			return ifds.NewIFD(t.ByteOrder, ifds.MknoteIFD, t.IfdIndex, t.ValueOffset, 0)
		}
	case ifds.SubIfd0, ifds.SubIfd1, ifds.SubIfd2, ifds.SubIfd3, ifds.SubIfd4, ifds.SubIfd5:
		return ifds.NewIFD(t.ByteOrder, t.Ifd, t.IfdIndex, t.ValueOffset, 0)
	}

	return ifds.NewIFD(t.ByteOrder, ifds.NullIFD, t.IfdIndex, t.ValueOffset, 0)
}
