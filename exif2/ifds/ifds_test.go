package ifds

import (
	"testing"

	"github.com/evanoberholster/imagemeta/exif2/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif2/tag"
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
		ifd := NewIFD(tag.LittleEndian, v.ifdType, 0, 0)

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

		// Ifd ChildIfd test

		childIFDtest(t, ifd, NewIFD(ifd.ByteOrder, ExifIFD, 0, 0), IFD0, ExifTag, true)
		childIFDtest(t, ifd, NewIFD(ifd.ByteOrder, GPSIFD, 0, 0), IFD0, GPSTag, true)
		childIFDtest(t, ifd, NewIFD(ifd.ByteOrder, SubIFD, 0, 0), IFD0, SubIFDs, true)

		childIFDtest(t, ifd, NewIFD(ifd.ByteOrder, MknoteIFD, 0, 0), ExifIFD, exififd.MakerNote, true)
		childIFDtest(t, ifd, NewIFD(ifd.ByteOrder, v.exifIFD, 0, 0), NullIFD, ExifTag, false)
	}
}

func TestValidIfd(t *testing.T) {
	if IfdType(100).String() != NullIFD.String() {
		t.Errorf("Incorrect IFD String, wanted %s got %s", NullIFD.String(), IfdType(100).String())
	}

}

func childIFDtest(t *testing.T, ifd Ifd, childIfd Ifd, testType IfdType, id tag.ID, a bool) {
	if ifd.IsType(testType) {
		if cIfd := ifd.ChildIfd(tag.Tag{ID: id}); cIfd != childIfd && a {
			t.Errorf("Incorrect Ifd: \"%s\" ChildIFD: \"%s\", wanted \"%s\" got \"%s\"", ifd.Type, cIfd.Type, cIfd, childIfd)
		}
	}
}

func tagTest(t *testing.T, ifd Ifd, testType IfdType, id tag.ID, tagName string) {
	if ifd.IsType(testType) {
		if ifd.TagName(id) != tagName {
			t.Errorf("%s, Expected \"%s\" got \"%s\"", ifd.String(), tagName, ifd.TagName(id))
		}
	}
}
