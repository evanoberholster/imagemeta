package meta

import (
	"testing"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/stretchr/testify/assert"
)

func TestBinaryOrder(t *testing.T) {
	buf := []byte{0, 0, 0, 0}
	bo := BinaryOrder(buf)
	if bo != UnknownEndian {
		t.Error("Binary Order for an empty buffer should be nil.")
	}

	buf = []byte{0x49, 0x49, 0x2a, 0}
	bo = BinaryOrder(buf)
	if bo != LittleEndian {
		t.Errorf("Binary Order expected %T got %T", LittleEndian, bo)
	}

	buf = []byte{0x4d, 0x4d, 0, 0x2a}
	bo = BinaryOrder(buf)
	if bo != BigEndian {
		t.Errorf("Binary Order expected %T got %T", BigEndian, bo)
	}
}

func TestXmpHeader(t *testing.T) {
	h1 := XmpHeader{1, 2}
	h2 := NewXMPHeader(1, 2)
	assert.Equal(t, h1, h2, "")
}

func TestExifHeader(t *testing.T) {
	h1 := ExifHeader{ByteOrder: BigEndian, FirstIfd: ifds.IFD0, FirstIfdOffset: 1234, TiffHeaderOffset: 16, ExifLength: 1024, ImageType: imagetype.ImagePNG}
	h2 := NewExifHeader(BigEndian, 1234, 16, 1024, imagetype.ImagePNG)
	h2.FirstIfd = ifds.IFD0

	assert.Equal(t, h1, h2, "")
	assert.True(t, h2.IsValid(), "IsValid")
}

func TestMetadata(t *testing.T) {
	m := Metadata{Dim: NewDimensions(1024, 768), It: imagetype.ImageDNG}

	assert.Equal(t, m.Dimensions(), NewDimensions(1024, 768))
	assert.Equal(t, m.ImageType(), imagetype.ImageDNG)

	// Aspect Ratio
	assert.Equal(t, m.Dim.AspectRatio(), float32(1024)/float32(768))
	assert.Equal(t, NewDimensions(0, 0).AspectRatio(), float32(0.0))

	// Orientation
	assert.Equal(t, int(m.Dim.Orientation()), 0)
	assert.Equal(t, int(NewDimensions(300, 400).Orientation()), 1)

	// Width and Height
	w, h := m.Dim.Size()
	assert.Equal(t, int(w), 1024)
	assert.Equal(t, int(h), 768)

	assert.NotEqual(t, m.Dim.String(), "")

}
