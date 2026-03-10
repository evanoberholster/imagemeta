package makernote

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestIdentifyCameraMakeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want CameraMake
	}{
		{name: "canon", in: "Canon", want: CameraMakeCanon},
		{name: "nikon", in: "NIKON CORPORATION", want: CameraMakeNikon},
		{name: "apple", in: "Apple", want: CameraMakeApple},
		{name: "panasonic", in: "Panasonic", want: CameraMakePanasonic},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IdentifyCameraMakeString(tt.in); got != tt.want {
				t.Fatalf("IdentifyCameraMakeString(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseNikonHeader(t *testing.T) {
	t.Parallel()

	var hdr [18]byte
	copy(hdr[:5], []byte("Nikon"))
	hdr[10] = 'M'
	hdr[11] = 'M'
	hdr[12] = 0x00
	hdr[13] = 0x2a
	hdr[17] = 0x08

	bo, off, ok := ParseNikonHeader(hdr[:])
	if !ok {
		t.Fatal("ParseNikonHeader should succeed")
	}
	if bo != utils.BigEndian {
		t.Fatalf("byte order = %v, want BigEndian", bo)
	}
	if off != 8 {
		t.Fatalf("ifd offset = %d, want 8", off)
	}
}

func TestMakerNoteBitset(t *testing.T) {
	t.Parallel()

	var mn Info
	mn.MarkTagParsed(9)
	mn.MarkTagParsed(9)
	if !mn.HasTagParsed(9) {
		t.Fatal("HasTagParsed(9) should be true")
	}
	if mn.HasTagParsed(10) {
		t.Fatal("HasTagParsed(10) should be false")
	}
}
