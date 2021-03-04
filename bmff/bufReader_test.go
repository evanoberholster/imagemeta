package bmff

import (
	"bufio"
	"bytes"
	"io"

	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/stretchr/testify/assert"
)

func newTestBox(data []byte) box {
	bmr := NewReader(bytes.NewReader(data))
	bmr.br.remain = len(data)
	return box{bufReader: bmr.br, size: int64(len(data))}
}

var testReadBoxData = []struct {
	name   string
	box    box
	data   []byte
	err    error
	assert bool
}{
	{"Normal ftyp box", box{size: 120, bufReader: bufReader{remain: 112, offset: 8}, boxType: TypeFtyp}, []byte{0, 0, 0, 120, 'f', 't', 'y', 'p', 0, 0, 0, 0, 0, 0, 0}, nil, true},
	{"Large Box (64-bit) meta box", box{size: 61563, bufReader: bufReader{remain: 61547, offset: 16}, boxType: TypeMeta}, []byte{0, 0, 0, 1, 'm', 'e', 't', 'a', 0, 0, 0, 0, 0, 0, 240, 123}, nil, true},
	{"Large Box (64-bit) uint64", box{}, []byte{0, 0, 0, 1, 'm', 'e', 't', 'a', 255, 0, 0, 0, 0, 0, 240, 123}, errLargeBox, false},
	{"Box too short", box{}, []byte{0, 0, 0, 120, 'm', 'e', 't'}, ErrBufLength, false},
	{"Large Box too short", box{}, []byte{0, 0, 0, 1, 'm', 'e', 't', 'a'}, ErrBufLength, false},
}

func TestBufReaderReadInnerBox(t *testing.T) {
	for _, v := range testReadBoxData {
		outer := newTestBox(v.data)
		inner, err := outer.readInnerBox()
		assert.ErrorIs(t, err, v.err, v.name)

		inner.bufReader.Reader = nil
		if v.assert {
			assert.Equalf(t, v.box, inner, "error message: %s", v.name)
		}
	}
}

func TestCloseInnerBox(t *testing.T) {
	data := []byte{0, 0, 0, 18, 'c', 'o', 'l', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	b := newTestBox(data)
	b.boxType = TypeMeta
	inner, err := b.readInnerBox()
	if err != nil {
		t.Error(err)
	}
	// expected offset 10, remain 8
	err = b.closeInnerBox(&inner)
	assert.ErrorIs(t, err, nil)
	//assert.ErrorAs(t, nil, err, "TestCloseInnerBox")
	expected := box{boxType: TypeMeta, size: 18}
	expected.bufReader = b.bufReader
	expected.bufReader.offset = 18
	expected.bufReader.remain = 0
	assert.Equal(t, expected, b, "TestCloseInnerBox")

	// Should discard 0, because remain is empty
	err = b.discard(10)
	assert.ErrorIs(t, err, nil)

	data = []byte{0, 0, 0, 18, 'c', 'o', 'l', 'r', 0, 0, 0, 0, 0, 0, 0, 0}
	b = newTestBox(data)
	b.boxType = TypeMeta
	inner, err = b.readInnerBox()
	if err != nil {
		t.Error(err)
	}
	// expected offset 10, remain 8
	err = b.closeInnerBox(&inner)
	assert.ErrorIs(t, err, io.EOF)

	// Should return err negative count
	err = b.discard(10)
	assert.ErrorIs(t, err, bufio.ErrNegativeCount)
}

func TestBufReaderReadUUID(t *testing.T) {
	UUIDtests := []struct {
		name     string
		buf      []byte
		expected meta.UUID
		err      error
	}{
		{"CR3 uuid", []byte{133, 192, 182, 135, 130, 15, 17, 224, 129, 17, 244, 206, 70, 43, 106, 72}, CR3MetaBoxUUID, nil},
		{"Err Length", []byte{0, 0, 0, 0, 0, 0}, meta.NilUUID, ErrBufLength},
		{"Err UUID", []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, meta.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, nil},
	}
	for _, v := range UUIDtests {
		b := newTestBox(v.buf)
		b.boxType = TypeUUID
		uuid, err := b.readUUID()
		assert.ErrorIs(t, err, v.err)
		assert.Equal(t, v.expected, uuid, "TestReadUUID")
	}
}

func TestBufReaderReadFlags(t *testing.T) {
	var expected Flags = 2 << 24 // Version 2
	expected++                   // Flags 1

	data := []byte{2, 0, 0, 1}
	outer := newTestBox(data)
	flags, err := outer.readFlags()
	if err != nil {
		t.Errorf("%v", err)
	}
	if expected.Version() != flags.Version() {
		t.Errorf("Flags Version Test Error: got %v, expected %v", flags.Version(), expected.Version())
	}
	if expected.Flags() != flags.Flags() {
		t.Errorf("Flags Flags Test Error: got %v, expected %v", flags.Flags(), expected.Flags())
	}
	// Error Length
	_, err = outer.readFlags()
	assert.ErrorIs(t, err, ErrBufLength)
}

//var testReadItemType = []struct {
//	name     string
//	itemType ItemType
//	data     []byte
//	remain   int
//	err      error
//	assert   bool
//}{
//	{"ItemTypeInfe", ItemTypeInfe, []byte("infe "), 5, nil, true},
//	{"ItemTypeMime", ItemTypeMime, []byte("mime "), 5, nil, true},
//	{"ItemTypeURI", ItemTypeURI, []byte("uri  "), 5, nil, true},
//	{"ItemTypeAv01", ItemTypeAv01, []byte("av01 "), 5, nil, true},
//	{"ItemTypeHvc1", ItemTypeHvc1, []byte("hvc1 "), 5, nil, true},
//	{"ItemTypeGrid", ItemTypeGrid, []byte("grid "), 5, nil, true},
//	{"ItemTypeExif", ItemTypeExif, []byte("Exif "), 5, nil, true},
//	{"ItemTypeInfe Type 2", ItemTypeInfe, append([]byte("infe"), 0), 5, nil, true},
//	{"ItemType too Short", ItemTypeUnknown, []byte{0, 0, 0}, 5, ErrBufLength, false},
//	{"ItemType remain too Short", ItemTypeUnknown, []byte{0, 0, 0}, 4, ErrBufLength, false},
//}

//func TestBufReaderReadItemType(t *testing.T) {
//	for _, v := range testReadItemType {
//		outer := newTestBox(v.data)
//		outer.remain = v.remain
//		itemType, err := outer.readItemType()
//		assert.ErrorIsf(t, err, v.err, v.name)
//		if v.assert {
//			assert.Equalf(t, v.itemType, itemType, "error message: %s", v.name)
//		}
//	}
//}

var testReadBrand = []struct {
	name   string
	brand  Brand
	data   []byte
	remain int
	err    error
	assert bool
}{
	{"brandAvif", brandAvif, []byte("avif "), 5, nil, true},
	{"brandHeic", brandHeic, []byte("heic "), 5, nil, true},
	{"brandCrx", brandCrx, []byte("crx "), 5, nil, true},
	{"brandMif1", brandMif1, []byte("mif1 "), 5, nil, true},
	{"brandUnknown", brandUnknown, []byte("nnna "), 5, nil, true},
	{"brand Error", brandUnknown, []byte("nnn"), 5, ErrBufLength, false},
	{"brand Error", brandUnknown, []byte("nnn"), 3, ErrBufLength, false},
}

func TestBufReaderReadBrand(t *testing.T) {
	for _, v := range testReadBrand {
		outer := newTestBox(v.data)
		outer.remain = v.remain
		brand, err := outer.readBrand()
		assert.ErrorIsf(t, err, v.err, v.name)
		if v.assert {
			assert.Equalf(t, v.brand, brand, "error message: %s", v.name)
		}
	}

}

var testReadUint8 = []struct {
	name   string
	data   []byte
	remain int
	i      uint8
	err    error
}{
	{"Normal", []byte{5}, 1, 5, nil},
	{"Error", []byte{}, 0, 0, io.EOF},
	{"Error", []byte{}, 1, 0, io.EOF},
}
var testReadUint16 = []struct {
	name   string
	data   []byte
	remain int
	i      uint16
	err    error
}{
	{"Normal", []byte{10, 12}, 2, 2572, nil},
	{"Error", []byte{}, 0, 0, io.EOF},
	{"Error", []byte{1}, 2, 0, io.EOF},
}
var testReadUint32 = []struct {
	name   string
	data   []byte
	remain int
	i      uint32
	err    error
}{
	{"Normal", []byte{0, 1, 15, 12}, 4, 69388, nil},
	{"Error", []byte{}, 0, 0, io.EOF},
	{"Error", []byte{1, 12, 34}, 4, 0, io.EOF},
}
var testReadUintN = []struct {
	name   string
	data   []byte
	remain int
	bits   uint8
	i      uint64
	err    error
}{
	{"UintN Remain Error", []byte{5, 0, 0, 10}, 3, 64, 0, io.EOF},
	{"UintN Remain Error", []byte{5, 0}, 4, 32, 0, io.EOF},
	{"UintN Invalid bit size Error", []byte{5, 0}, 2, 15, 0, errUintSize},
	{"UintN 0", []byte{5, 0, 0, 10}, 4, 0, 0, nil},
	{"UintN 8", []byte{5, 0, 0, 10}, 1, 8, 5, nil},
	{"UintN 16", []byte{0, 15, 0, 11}, 2, 16, 15, nil},
	{"UintN 32", []byte{0, 0, 10, 12}, 4, 32, 2572, nil},
	{"UintN 64", []byte{0, 0, 0, 0, 1, 1, 1, 1}, 8, 64, 16843009, nil},
}

func TestBufReaderReadUint(t *testing.T) {
	for _, v := range testReadUint8 {
		outer := newTestBox(v.data)
		outer.remain = v.remain
		i, err := outer.readUint8()
		assert.ErrorIsf(t, err, v.err, v.name)

		if i != v.i {
			t.Errorf("Uint8 Test Error: got %d, expected %d", i, v.i)
		}
	}
	for _, v := range testReadUint16 {
		outer := newTestBox(v.data)
		outer.remain = v.remain
		i, err := outer.readUint16()
		assert.ErrorIsf(t, err, v.err, v.name)

		if i != v.i {
			t.Errorf("Uint8 Test Error: got %d, expected %d", i, v.i)
		}
	}
	for _, v := range testReadUint32 {
		outer := newTestBox(v.data)
		outer.remain = v.remain
		i, err := outer.readUint32()
		assert.ErrorIsf(t, err, v.err, v.name)

		if i != v.i {
			t.Errorf("Uint8 Test Error: got %d, expected %d", i, v.i)
		}
	}
	for _, v := range testReadUintN {
		outer := newTestBox(v.data)
		outer.remain = v.remain
		i, err := outer.readUintN(v.bits)
		assert.ErrorIsf(t, err, v.err, v.name)

		if i != v.i {
			t.Errorf("UintN %d bits Test Error: got %d, expected %d", v.bits, i, v.i)
		}
	}
}

func TestBufReaderReadString(t *testing.T) {
}
