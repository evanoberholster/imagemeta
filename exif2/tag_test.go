package exif2

import (
	"bytes"
	"testing"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/stretchr/testify/assert"
)

func TestChildIfd(t *testing.T) {
	tests := []struct {
		t Tag
		i ifds.Ifd
	}{
		{t: Tag{}, i: ifds.Ifd{}},
		{t: Tag{Ifd: ifds.IFD0, ID: ifds.ExifTag, ByteOrder: utils.LittleEndian, ValueOffset: 12345}, i: ifds.Ifd{Offset: 12345, ByteOrder: utils.LittleEndian, Type: ifds.ExifIFD}},
		{t: Tag{Ifd: ifds.IFD0, ID: ifds.GPSTag, ByteOrder: utils.LittleEndian, ValueOffset: 23456}, i: ifds.Ifd{Offset: 23456, ByteOrder: utils.LittleEndian, Type: ifds.GPSIFD}},
		{t: Tag{Ifd: ifds.IFD0, ID: ifds.SubIFDs, ByteOrder: utils.LittleEndian, ValueOffset: 112233}, i: ifds.Ifd{Offset: 112233, ByteOrder: utils.LittleEndian, Type: ifds.NullIFD}},
		//{t: Tag{Ifd: ifds.ExifIFD, ID: exififd.MakerNote, ByteOrder: utils.BigEndian, ValueOffset: 3456}, i: ifds.Ifd{Offset: 3456, ByteOrder: utils.BigEndian, Type: ifds.MknoteIFD}},
	}

	for _, test := range tests {
		assert.Equal(t, test.i, test.t.childIfd())
	}
}

func TestTag(t *testing.T) {
	t2 := NewTag(tag.ID(0x0010), tag.TypeASCII, 16, 0x0002, 0, 0, 0)

	if t2.ID != tag.ID(0x0010) {
		t.Errorf("Incorrect Tag ID wanted 0x%04x got 0x%04x", tag.ID(0x0010), t2.ID)
	}
	if t2.Type != tag.TypeASCII || !t2.IsType(tag.TypeASCII) || !t2.Type.Is(tag.TypeASCII) {
		t.Errorf("Incorrect Tag Type wanted %s got %s", tag.TypeASCII, t2.Type)
	}
	if t2.UnitCount != 16 {
		t.Errorf("Incorrect Tag UnitCount wanted %d got %d", 16, t2.UnitCount)
	}
	if t2.ValueOffset != 0x0002 {
		t.Errorf("Incorrect Tag Offset wanted 0x%04x got 0x%04x", 0x0002, t2.ValueOffset)
	}
	if t2.IsEmbedded() {
		t.Errorf("ValueIsEmbedded is true when equal or less than 4 bytes")
	}
	if t2.Size() != 16 {
		t.Errorf("Incorrect Tag Size wanted %d got %d", 16, t2.Size())
	}

	if t2.ID.String() != "0x0010" {
		t.Errorf("Incorrect ID String wanted %v got %v", "0x0010", t2.ID.String())
	}

	t2 = NewTag(tag.ID(0x0010), 100, 16, 0x0002, 0, 0, 0)
	if t2.IsValid() {
		t.Error("incorrect wanted invalid")
	}
	t2.Type = tag.TypeIfd
	t2.UnitCount = 1
	if !t2.IsIfd() {
		t.Errorf("Incorrect Tag Type wanted %s got %s", tag.TypeIfd, t2.Type)
	}
	if t2.Type.String() != "IFD" {
		t.Errorf("Incorrect Tag String wanted %s got %s", "IFD", t2.Type.String())
	}
	if t2.Size() != tag.TypeIfdSize {
		t.Errorf("Incorrect Tag Size wanted %d got %d", tag.TypeIfdSize, t2.Size())
	}
	bufVal := [4]byte{2, 0, 0, 0}
	var buf [4]byte
	t2.EmbeddedValue(buf[:])
	if !bytes.Equal(buf[:], bufVal[:]) {
		t.Errorf("Incorrect Embedded Value wanted %d got %d", bufVal, buf)
	}
}
