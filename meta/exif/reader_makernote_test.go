package exif

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestParseMakerNoteTIFFPrefix(t *testing.T) {
	t.Parallel()

	tiffLE := [8]byte{'I', 'I', '*', 0, 8, 0, 0, 0}
	bo, off, ok := parseMakerNoteTIFFPrefix(tiffLE[:])
	if !ok {
		t.Fatal("parseMakerNoteTIFFPrefix returned ok=false for valid TIFF header")
	}
	if bo != utils.LittleEndian {
		t.Fatalf("byteOrder = %v, want %v", bo, utils.LittleEndian)
	}
	if off != 8 {
		t.Fatalf("ifdRelOffset = %d, want 8", off)
	}

	canonPrefix := [8]byte{'C', 'a', 'n', 'o', 'n', 0, 0, 0}
	if _, _, ok = parseMakerNoteTIFFPrefix(canonPrefix[:]); ok {
		t.Fatal("parseMakerNoteTIFFPrefix returned ok=true for non-TIFF prefix")
	}
}

func TestIsCanonMakerNotePrefix(t *testing.T) {
	t.Parallel()

	if !isCanonMakerNotePrefix([]byte{'C', 'a', 'n', 'o', 'n', 0, 0, 0}) {
		t.Fatal("isCanonMakerNotePrefix returned false for Canon prefix")
	}
	if isCanonMakerNotePrefix([]byte{'I', 'I', '*', 0, 8, 0, 0, 0}) {
		t.Fatal("isCanonMakerNotePrefix returned true for TIFF header")
	}
}
