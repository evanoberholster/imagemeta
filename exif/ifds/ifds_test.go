package ifds

import (
	"testing"

	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

func TestKey(t *testing.T) {
	tests := []struct {
		tagID    tag.ID
		ifdType  IfdType
		ifdIndex uint8
	}{
		{TileLength, RootIFD, 1},
		{TileByteCounts, ExifIFD, 2},
		{CacheVersion, GPSIFD, 3},
		{OpcodeList3, MknoteIFD, 4},
		{OpcodeList2, RootIFD, 5},
	}

	for _, v := range tests {
		key := NewKey(v.ifdType, v.ifdIndex, v.tagID)
		ifdType, ifdIndex, tagID := key.Val()
		key2 := NewKey(ifdType, ifdIndex, tagID)

		if key != key2 {
			t.Errorf("Expected Key %d got %d", key, key2)
		}
		if v.tagID != tagID {
			t.Errorf("Expected TagID %s got %s", v.tagID, tagID)
		}
		if v.ifdType != ifdType {
			t.Errorf("Expected IfdType %d got %d", key, key2)
		}
		if v.ifdIndex != ifdIndex {
			t.Errorf("Expected IfdIndex %d got %d", v.ifdIndex, ifdIndex)
		}

		if !ifdType.IsValid() {
			t.Errorf("Expected %s, got %t", "true", ifdType.IsValid())
		}
	}
}

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
		{RootIFD, "Ifd", 0, NullIFD, 0, NullIFD, true},
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
		ifd := NewIFD(v.ifdType, 0, 0)

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
		tagTest(t, ifd, RootIFD, ExifTag, "ExifTag")
		tagTest(t, ifd, ExifIFD, exififd.ApertureValue, "ApertureValue")
		tagTest(t, ifd, GPSIFD, gpsifd.GPSAltitude, "GPSAltitude")
		tagTest(t, ifd, MknoteIFD, mknote.CanonAFInfo, "CanonAFInfo")
		tagTest(t, ifd, 255, ExifTag, "0x8769")

		ta := tag.Tag{}
		ta.ID = v.rootTag

		// Ifd ChildIfd test

		childIFDtest(t, ifd, RootIFD, ExifTag, true)
		childIFDtest(t, ifd, RootIFD, GPSTag, true)
		childIFDtest(t, ifd, RootIFD, SubIFDs, true)

		childIFDtest(t, ifd, ExifIFD, exififd.MakerNote, true)
		childIFDtest(t, ifd, NullIFD, ExifTag, false)
	}
}

func TestValidIfd(t *testing.T) {
	if IfdType(100).String() != NullIFD.String() {
		t.Errorf("Incorrect IFD String, wanted %s got %s", NullIFD.String(), IfdType(100).String())
	}

}

func childIFDtest(t *testing.T, ifd Ifd, testType IfdType, id tag.ID, a bool) {
	if ifd.IsType(testType) {
		if cIfd, b := ifd.IsChildIfd(tag.Tag{ID: id}); b != a {
			t.Errorf("Incorrect Ifd: \"%s\" ChildIFD: \"%s\", wanted \"%t\" got \"%t\"", ifd.Type, cIfd.Type, a, b)
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
