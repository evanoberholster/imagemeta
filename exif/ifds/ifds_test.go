package ifds

import (
	"testing"

	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	key := NewKey(RootIFD, 1, TileWidth)

	assert.Equal(t, key, Key(0x1010142), "Ifd Key")
}

func TestIfdString(t *testing.T) {
	testIfds := []struct {
		ifd     IFD
		str     string
		rootTag tag.ID
		rootIFD IFD
		exifTag tag.ID
		exifIFD IFD
	}{
		{RootIFD, "Ifd", 0, NullIFD, 0, NullIFD},
		{SubIFD, "Ifd/SubIfd", SubIFDs, SubIFD, 0, NullIFD},
		{ExifIFD, "Ifd/Exif", ExifTag, ExifIFD, 0, NullIFD},
		{GPSIFD, "Ifd/GPS", GPSTag, GPSIFD, 0, NullIFD},
		{IopIFD, "Ifd/Iop", 0, NullIFD, 0, NullIFD},
		{MknoteIFD, "Ifd/Exif/Makernote", exififd.MakerNote, NullIFD, exififd.MakerNote, MknoteIFD},
		{DNGAdobeDataIFD, "Ifd/DNGAdobeData", 0, NullIFD, 0, NullIFD},
		{NullIFD, "Unknown", 0, NullIFD, 0, NullIFD},
	}

	for i, v := range testIfds {
		assert.Equal(t, v.ifd.String(), v.str)
		ta := tag.Tag{}
		ta.ID = v.rootTag

		assert.Equal(t, v.rootIFD, RootIFD.IsChildIfd(ta), "RootIfd Children: %v")

		ta.ID = v.exifTag
		assert.Equal(t, v.exifIFD, ExifIFD.IsChildIfd(ta), "ExifIfd Children: %d", i)
	}
	assert.Equal(t, RootIFD.TagName(ExifTag), "ExifTag")
	assert.Equal(t, ExifIFD.TagName(exififd.ApertureValue), "ApertureValue")
	assert.Equal(t, GPSIFD.TagName(gpsifd.GPSAltitude), "GPSAltitude")
	assert.Equal(t, MknoteIFD.TagName(mknote.BatteryType), "0x0038")
}
