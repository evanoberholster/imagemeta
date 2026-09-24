package isobmff

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// TestReadBoxSizeZeroSeekable checks that a trailing zero-size box on a
// seekable source resolves to end of file instead of failing.
func TestReadBoxSizeZeroSeekable(t *testing.T) {
	t.Parallel()
	var data []byte
	data = append(data, 0, 0, 0, 20, 'f', 't', 'y', 'p')
	data = append(data, 'i', 's', 'o', 'm', 0, 0, 0, 0, 'i', 's', 'o', 'm')
	data = append(data, 0, 0, 0, 0, 'f', 'r', 'e', 'e')
	data = append(data, bytes.Repeat([]byte{0x63}, 8)...)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)
	if err := r.ReadFTYP(); err != nil {
		t.Fatal(err)
	}
	if r.ftyp.MajorBrand != brandIsom {
		t.Fatalf("major brand = %v, want isom", r.ftyp.MajorBrand)
	}
	b, err := r.readBox()
	if err != nil {
		t.Fatal(err)
	}
	if b.offset != 20 || !b.isType(typeFree) {
		t.Fatalf("box offset = %d type = %s, want 20 free", b.offset, b.boxType)
	}
	if b.size != 16 {
		t.Fatalf("box size = %d, want 16 (extends to EOF)", b.size)
	}
	if err := b.close(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.readBox(); !errors.Is(err, io.EOF) {
		t.Fatalf("readBox after EOF = %v, want io.EOF", err)
	}
}

// TestReadBoxSizeZeroStream checks that a zero-size box on a non-seekable
// source still fails with ErrBoxSizeZero.
func TestReadBoxSizeZeroStream(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	buf.Write([]byte{0, 0, 0, 0, 'f', 'r', 'e', 'e'})
	buf.Write(bytes.Repeat([]byte{0x63}, 8))

	r := NewReader(&buf, nil, nil, nil)
	t.Cleanup(r.Close)
	_, err := r.readBox()
	if !errors.Is(err, ErrBoxSizeZero) {
		t.Fatalf("readBox = %v, want ErrBoxSizeZero", err)
	}
}

// TestReadInnerBoxSizeZero checks that a nested zero-size box is bounded by
// its outer container instead of failing.
func TestReadInnerBoxSizeZero(t *testing.T) {
	t.Parallel()
	var data []byte
	data = append(data, 0, 0, 0, 0, 'f', 'r', 'e', 'e')
	data = append(data, bytes.Repeat([]byte{0x63}, 8)...)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)
	outer := box{reader: r, size: len(data), remain: len(data)}
	inner, ok, err := outer.readInnerBox()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("readInnerBox returned no box")
	}
	if !inner.isType(typeFree) {
		t.Fatalf("inner type = %s, want free", inner.boxType)
	}
	if inner.size != len(data) {
		t.Fatalf("inner size = %d, want %d (bounded by outer)", inner.size, len(data))
	}
	if err := inner.close(); err != nil {
		t.Fatal(err)
	}
	if outer.remain != 0 {
		t.Fatalf("outer remain = %d, want 0", outer.remain)
	}
}
