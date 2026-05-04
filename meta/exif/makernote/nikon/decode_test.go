package nikon

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestVersionStringBinaryDigits(t *testing.T) {
	t.Parallel()

	raw := []byte{2, 1, 1, 0}
	if got := VersionString(raw); got != "2110" {
		t.Fatalf("VersionString(%v) = %q, want %q", raw, got, "2110")
	}
}

func TestBitsetIndices(t *testing.T) {
	t.Parallel()

	raw := []byte{0b00000101, 0b00001010}
	got := BitsetIndices(raw)
	want := []int{0, 2, 9, 11}
	if len(got) != len(want) {
		t.Fatalf("BitsetIndices len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("BitsetIndices[%d] = %d, want %d (%v)", i, got[i], want[i], got)
		}
	}
}

func TestDecodeLensFormatting(t *testing.T) {
	t.Parallel()

	raw := make([]byte, 32)
	utils.BigEndian.PutUint32(raw[0:4], 100)
	utils.BigEndian.PutUint32(raw[4:8], 1)
	utils.BigEndian.PutUint32(raw[8:12], 400)
	utils.BigEndian.PutUint32(raw[12:16], 1)
	utils.BigEndian.PutUint32(raw[16:20], 9)
	utils.BigEndian.PutUint32(raw[20:24], 2)
	utils.BigEndian.PutUint32(raw[24:28], 28)
	utils.BigEndian.PutUint32(raw[28:32], 5)

	if got := DecodeLens(raw, utils.BigEndian); got != "100 400 4.5 5.6" {
		t.Fatalf("DecodeLens() = %q, want %q", got, "100 400 4.5 5.6")
	}
}

func TestFileInfoPrefersLE(t *testing.T) {
	t.Parallel()

	if !fileInfoPrefersLE("NIKON D750") {
		t.Fatal("fileInfoPrefersLE(NIKON D750) = false, want true")
	}
	if fileInfoPrefersLE("NIKON Z 9") {
		t.Fatal("fileInfoPrefersLE(NIKON Z 9) = true, want false")
	}
}
