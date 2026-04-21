package meta

import (
	"bytes"
	"strings"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestBinaryOrder(t *testing.T) {
	buf := []byte{0, 0, 0, 0}
	bo := utils.BinaryOrder(buf)
	if bo != utils.UnknownEndian {
		t.Error("Binary Order for an empty buffer should be nil.")
	}

	buf = []byte{0x49, 0x49, 0x2a, 0}
	bo = utils.BinaryOrder(buf)
	if bo != utils.LittleEndian {
		t.Errorf("Binary Order expected %T got %T", utils.LittleEndian, bo)
	}

	buf = []byte{0x4d, 0x4d, 0, 0x2a}
	bo = utils.BinaryOrder(buf)
	if bo != utils.BigEndian {
		t.Errorf("Binary Order expected %T got %T", utils.BigEndian, bo)
	}
}

func TestXmpHeader(t *testing.T) {
	h1 := XmpHeader{1, 2}
	h2 := NewXMPHeader(1, 2)
	assert.Equal(t, h1, h2, "")
}

func TestExifHeader(t *testing.T) {
	h1 := ExifHeader{ByteOrder: utils.BigEndian, FirstIfd: tag.IFD0, FirstIfdOffset: 1234, TiffHeaderOffset: 16, ExifLength: 1024, ImageType: imagetype.ImagePNG}
	h2 := NewExifHeader(utils.BigEndian, 1234, 16, 1024, imagetype.ImagePNG)
	h2.FirstIfd = tag.IFD0

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

func TestHeaderMarshalZerologObjectUsesLowerCamelKeys(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	logger.Info().
		Object("exifHeader", ExifHeader{
			ByteOrder:        utils.LittleEndian,
			FirstIfd:         tag.IFD0,
			FirstIfdOffset:   8,
			TiffHeaderOffset: 12,
			ExifLength:       4096,
			ImageType:        imagetype.ImageJPEG,
		}).
		Object("xmpHeader", XPacketHeader{
			Offset:       42,
			Length:       128,
			HasXPacketPI: true,
			HasXMPMeta:   true,
		}).
		Object("previewHeader", PreviewHeader{
			Size:      256,
			Width:     3,
			Height:    2,
			ImageType: imagetype.ImageJPEG,
			Source:    PreviewSourcePRVW,
		}).
		Msg("test")

	out := buf.String()
	for _, want := range []string{
		`"firstIfd":"IFD0"`,
		`"firstIfdOffset":8`,
		`"tiffHeaderOffset":12`,
		`"exifLength":4096`,
		`"byteOrder":"LittleEndian"`,
		`"imageType":"image/jpeg"`,
		`"offset":42`,
		`"length":128`,
		`"hasXPacketPI":true`,
		`"hasXMPMeta":true`,
		`"size":256`,
		`"width":3`,
		`"height":2`,
		`"source":"PRVW"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
	for _, unwanted := range []string{`"FirstIfd"`, `"ExifLength"`, `"Endian"`, `"ImageType"`} {
		if strings.Contains(out, unwanted) {
			t.Fatalf("log output unexpectedly contains legacy key %s: %q", unwanted, out)
		}
	}
}
