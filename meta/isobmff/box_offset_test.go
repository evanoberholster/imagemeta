package isobmff

import (
	"bytes"
	"testing"
)

// TestBoxReadAdvancesAbsoluteOffset checks that consuming payload bytes
// through box.Read advances the reader's absolute offset exactly like
// Discard does. Callbacks receive boxes as io.Readers, so a missing update
// silently corrupts every later absolute offset (mdat extent resolution,
// preview payloads, logs).
func TestBoxReadAdvancesAbsoluteOffset(t *testing.T) {
	t.Parallel()
	var data []byte
	data = append(data, 0, 0, 0, 24, 'f', 't', 'y', 'p')
	data = append(data, bytes.Repeat([]byte{0x61}, 16)...)
	data = append(data, 0, 0, 0, 16, 'f', 'r', 'e', 'e')
	data = append(data, bytes.Repeat([]byte{0x62}, 8)...)
	data = append(data, 0, 0, 0, 16, 'f', 'r', 'e', 'e')
	data = append(data, bytes.Repeat([]byte{0x63}, 8)...)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)
	if err := r.ReadFTYP(); err != nil {
		t.Fatal(err)
	}
	b, err := r.readBox()
	if err != nil {
		t.Fatal(err)
	}
	tmp := make([]byte, 4)
	if _, err = b.Read(tmp); err != nil {
		t.Fatal(err)
	}
	if r.offset != 36 {
		t.Fatalf("r.offset = %d, want 36 after reading 4 payload bytes", r.offset)
	}
	if err = b.close(); err != nil {
		t.Fatal(err)
	}
	next, err := r.readBox()
	if err != nil {
		t.Fatal(err)
	}
	if next.offset != 40 || !next.isType(typeFree) {
		t.Fatalf("next box offset = %d type = %s, want 40 free", next.offset, next.boxType)
	}
	if err = next.close(); err != nil {
		t.Fatal(err)
	}
}

// TestReadInnerBoxNestedOffsets checks that nested child boxes resolve
// exact absolute offsets however deep the nesting. The offset must be
// captured before the child header is consumed.
func TestReadInnerBoxNestedOffsets(t *testing.T) {
	t.Parallel()
	var data []byte
	data = append(data, 0, 0, 0, 28, 'm', 'o', 'o', 'v')
	data = append(data, 0, 0, 0, 20, 't', 'r', 'a', 'k')
	data = append(data, 0, 0, 0, 12, 'f', 'r', 'e', 'e')
	data = append(data, bytes.Repeat([]byte{0x63}, 4)...)

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)
	outer, err := r.readBox()
	if err != nil {
		t.Fatal(err)
	}
	if outer.offset != 0 {
		t.Fatalf("outer offset = %d, want 0", outer.offset)
	}
	mid, ok, err := outer.readInnerBox()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("readInnerBox returned no box")
	}
	if mid.offset != 8 || !mid.isType(typeTrak) {
		t.Fatalf("mid offset = %d type = %s, want 8 trak", mid.offset, mid.boxType)
	}
	leaf, ok, err := mid.readInnerBox()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("readInnerBox returned no nested box")
	}
	if leaf.offset != 16 || !leaf.isType(typeFree) {
		t.Fatalf("leaf offset = %d type = %s, want 16 free", leaf.offset, leaf.boxType)
	}
	if err = leaf.close(); err != nil {
		t.Fatal(err)
	}
	if err = mid.close(); err != nil {
		t.Fatal(err)
	}
	if err = outer.close(); err != nil {
		t.Fatal(err)
	}
}
