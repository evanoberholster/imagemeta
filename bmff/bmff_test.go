package bmff

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFtyp = []struct {
	name     string
	ftyp     FileTypeBox
	metaJSON string
	filename string
	remain   int
	err      error
	assert   bool
}{
	{"Sample 1", FileTypeBox{MajorBrand: brandHeic, MinorVersion: "", Compatible: [6]Brand{brandMif1, brandHeic}}, "", "samples/1.sample", 30, nil, true},
	{"Sample 2", FileTypeBox{MajorBrand: brandMif1, MinorVersion: "", Compatible: [6]Brand{brandMif1, brandHeic}}, "", "samples/2.sample", 30, nil, true},
	{"Sample 3", FileTypeBox{MajorBrand: brandHeic, MinorVersion: "", Compatible: [6]Brand{brandHeic, brandMif1}}, "", "samples/3.sample", 20, nil, true},
	{"iPhone 11", FileTypeBox{MajorBrand: brandHeic, MinorVersion: "", Compatible: [6]Brand{brandMif1, brandMiaf, brandMiHB, brandHeic}}, "", "samples/iPhone11.sample", 20, nil, true},
	{"iPhone 12", FileTypeBox{MajorBrand: brandHeic, MinorVersion: "", Compatible: [6]Brand{brandMif1, brandMiHE, brandMiaf, brandMiHB, brandHeic}}, "", "samples/iPhone12.sample", 20, nil, true},
	//{"Avif Ftyp", FileTypeBox{MajorBrand: brandAvif, MinorVersion: "", Compatible: [6]Brand{}}, []byte("uri  "), 20, nil, true},
	//{"Cr3 Ftyp", FileTypeBox{MajorBrand: brandCrx, MinorVersion: "", Compatible: [6]Brand{}}, []byte("av01 "), 20, nil, true},
}

//func TestGenSamples(t *testing.T) {
//	testFilename := "../../test/img/iPhone12.heic"
//	f, err := os.Open(testFilename)
//	if err != nil {
//		panic(err)
//	}
//	defer func() {
//		err = f.Close()
//		if err != nil {
//			panic(err)
//		}
//	}()
//	buf, err := ioutil.ReadAll(f)
//	buf = buf[:3642]
//	dat, err := os.Create("samples/8.sample")
//	if err != nil {
//		panic(err)
//	}
//	defer func() {
//		err = dat.Close()
//		if err != nil {
//			panic(err)
//		}
//	}()
//	if _, err := dat.Write(buf); err != nil {
//		err = f.Close()
//		panic(err)
//	}
//}

func TestParseMeta(t *testing.T) {
	for _, v := range testFtyp {
		f, err := os.Open(v.filename)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}()
		bmr := NewReader(f)
		ftyp, err := bmr.ReadFtypBox()
		if err != v.err {
			t.Errorf("Error: (%s), %v", v.name, err)
		}
		if v.assert {
			assert.Equalf(t, v.ftyp, ftyp, "error message: %s", v.name)
		}
		_, err = bmr.ReadMetaBox()
		if err != v.err {
			t.Errorf("Error: (%s), %v", v.name, err)
		}

	}
}

func TestFileTypeBox(t *testing.T) {
	ftyp := FileTypeBox{}

	if ftyp.Type() != TypeFtyp {
		t.Errorf("Expected TypeFtyp")
	}
	if cleanString([]byte("abcd")) != "abcd" {
		t.Errorf("Expected complete string")
	}

	if brandAvif.String() != "avif" {
		t.Errorf("Brand Avif String Test Error: got %v, expected %v", brandAvif.String(), "avif")
	}

	if brandUnknown.String() != "nnnn" {
		t.Errorf("Brand Unknown String Test Error: got %v, expected %v", brandUnknown.String(), "nnnn")
	}

}

func TestParseHandler(t *testing.T) {
	expected := HandlerBox{
		Flags:       Flags(0),
		size:        34,
		HandlerType: handlerPict,
	}
	data := []byte{0, 0, 0, 34, 104, 100, 108, 114, 0, 0, 0, 0, 0, 0, 0, 0, 112, 105, 99, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	outer := newTestBox(data)
	inner, err := outer.readInnerBox()
	if err != nil {
		t.Errorf("Error: (%s), %v", "Handler", err)
	}
	if expected.Type() != inner.boxType {
		t.Errorf("Handler Box Test Error: got %v, expected %v", inner.boxType.String(), TypeHdlr.String())
	}
	hdlr, err := parseHdlr(&inner)
	if err != nil {
		t.Errorf("Error: (%s), %v", "Handler", err)
	}
	if a, ok := hdlr.(HandlerBox); ok {
		assert.Equalf(t, expected, a, "Handler Box")
	}

	copy(data[16:20], "nnnn")
	outer = newTestBox(data)
	inner, err = outer.readInnerBox()
	if err != nil {
		t.Errorf("Error: (%s), %v", "Handler", err)
	}
	_, err = parseHdlr(&inner)
	if err.Error() != "error Handler type unknown: nnnn" {
		t.Errorf("Error: (%s), %v", "Handler Box", err)
	}

	if handler([]byte("abcd")).String() != handlerUnknown.String() {
		t.Errorf("Handler Box Type Error: got %v, expected %v", handler([]byte("abcd")).String(), "nnnn")
	}
	if handler([]byte("pict")).String() != handlerPict.String() {
		t.Errorf("Handler Box Type Error: got %v, expected %v", handler([]byte("pict")).String(), "pict")
	}
	if handler([]byte("vide")).String() != handlerVide.String() {
		t.Errorf("Handler Box Type Error: got %v, expected %v", handler([]byte("vide")).String(), "vide")
	}
	if handler([]byte("meta")).String() != handlerMeta.String() {
		t.Errorf("Handler Box Type Error: got %v, expected %v", handler([]byte("meta")).String(), "meta")
	}
}

func TestItemType(t *testing.T) {
	it := itemType([]byte("infe"))
	if it.String() != "infe" {
		t.Errorf("Item Type Error: got %v, expected %v", it.String(), "infe")
	}
	// Error for itemType that is too long
	if itemType([]byte("infe123")).String() != ItemTypeUnknown.String() {
		t.Errorf("Item Type Error: got %v, expected %v", itemType([]byte("infe123")), ItemTypeUnknown)
	}
}
