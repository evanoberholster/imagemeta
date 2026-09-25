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
