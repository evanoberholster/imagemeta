package makernote

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
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
		{name: "nikon nul", in: "NIKON CORPORATION\x00", want: CameraMakeNikon},
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

func TestIdentifyCameraMake(t *testing.T) {
	t.Parallel()

	raw := []byte("NIKON CORPORATION\x00")
	if got := IdentifyCameraMake(raw); got != CameraMakeNikon {
		t.Fatalf("IdentifyCameraMake(%q) = %v, want %v", raw, got, CameraMakeNikon)
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

	bo, off, ok := nikon.ParseNikonHeader(hdr[:])
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
