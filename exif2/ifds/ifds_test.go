package ifds

import (
	"testing"

	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/stretchr/testify/assert"
)

func TestIfdString(t *testing.T) {
	testIfds := []struct {
		ifdType IfdType
		str     string
		rootTag tag.ID
		rootIFD IfdType
		exifTag tag.ID
		exifIFD IfdType
		valid   bool
	}{
		{IFD0, "Ifd", 0, NullIFD, 0, NullIFD, true},
		{SubIFD, "Ifd/SubIfd", SubIFDs, SubIFD, 0, NullIFD, true},
		{ExifIFD, "Ifd/Exif", ExifTag, ExifIFD, 0, NullIFD, true},
		{GPSIFD, "Ifd/GPS", GPSTag, GPSIFD, 0, NullIFD, true},
		{IopIFD, "Ifd/Iop", 0, NullIFD, 0, NullIFD, true},
		{MknoteIFD, "Ifd/Exif/Makernote", exififd.MakerNote, NullIFD, exififd.MakerNote, MknoteIFD, true},
		{DNGAdobeDataIFD, "Ifd/DNGAdobeData", 0, NullIFD, 0, NullIFD, true},
		{NullIFD, "UnknownIfd", 0, NullIFD, 0, NullIFD, false},
		{255, "UnknownIfd", 0, NullIFD, 0, NullIFD, false},
	}

	for _, v := range testIfds {
		ifd := NewIFD(meta.LittleEndian, v.ifdType, 0, 0)

		// Ifd Valid
		if ifd.IsValid() != v.valid {
			t.Errorf("Expected %s valid (%t) got valid (%t)", v.ifdType, v.valid, ifd.IsValid())
		}

		// Ifd String
		if v.ifdType.String() != v.str {
			t.Errorf("Expected \"%s\" got \"%s\"", v.str, v.ifdType)
		}

		// Ifd testing

		if ifd.String() == "" {
			t.Errorf("Expected \"%s\" got \"%s\"", "Some text", ifd.String())
		}
		// Ifd Tagname test
		tagTest(t, ifd, IFD0, ExifTag, "ExifTag")
		tagTest(t, ifd, ExifIFD, exififd.ApertureValue, "ApertureValue")
		tagTest(t, ifd, GPSIFD, gpsifd.GPSAltitude, "GPSAltitude")
		tagTest(t, ifd, MknoteIFD, mknote.CanonAFInfo, "CanonAFInfo")
		tagTest(t, ifd, 255, ExifTag, "0x8769")

		ta := tag.Tag{}
		ta.ID = v.rootTag
	}
}

func TestValidIfd(t *testing.T) {
	if IfdType(100).String() != NullIFD.String() {
		t.Errorf("Incorrect IFD String, wanted %s got %s", NullIFD.String(), IfdType(100).String())
	}

}

func tagTest(t *testing.T, ifd Ifd, testType IfdType, id tag.ID, tagName string) {
	if ifd.IsType(testType) {
		if ifd.TagName(id) != tagName {
			t.Errorf("%s, Expected \"%s\" got \"%s\"", ifd.String(), tagName, ifd.TagName(id))
		}
	}
}

func TestTagString(t *testing.T) {
	tests := []struct {
		id  tag.ID
		str string
	}{

		{0, "0x0000"},
		{Make, "Make"},
		{Model, "Model"},
	}
	for _, v := range tests {
		if TagString(v.id) != v.str {
			t.Errorf("expected \"%s\" got \"%s\"", v.str, TagString(v.id))
		}
	}
}

func TestChildIfd(t *testing.T) {
	tests := []struct {
		t tag.Tag
		i Ifd
	}{
		{t: tag.Tag{}, i: Ifd{}},
		{t: tag.Tag{Ifd: uint8(IFD0), ID: ExifTag, ByteOrder: meta.LittleEndian, ValueOffset: 12345}, i: Ifd{Offset: 12345, ByteOrder: meta.LittleEndian, Type: ExifIFD}},
		{t: tag.Tag{Ifd: uint8(IFD0), ID: GPSTag, ByteOrder: meta.LittleEndian, ValueOffset: 23456}, i: Ifd{Offset: 23456, ByteOrder: meta.LittleEndian, Type: GPSIFD}},
		{t: tag.Tag{Ifd: uint8(IFD0), ID: SubIFDs, ByteOrder: meta.LittleEndian, ValueOffset: 112233}, i: Ifd{Offset: 112233, ByteOrder: meta.LittleEndian, Type: SubIFD}},
		{t: tag.Tag{Ifd: uint8(ExifIFD), ID: exififd.MakerNote, ByteOrder: meta.BigEndian, ValueOffset: 3456}, i: Ifd{Offset: 3456, ByteOrder: meta.BigEndian, Type: MknoteIFD}},
	}

	for _, test := range tests {
		assert.Equal(t, test.i, ChildIfd(test.t))
	}
}
