package canon

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func TestTagCanonString(t *testing.T) {
	if got := TagCanonString(tag.ID(RawDataLength)); got != "RawDataLength" {
		t.Fatalf("TagCanonString(RawDataLength) = %q, want %q", got, "RawDataLength")
	}

	if got := TagCanonString(tag.ID(CanonLightingOpt)); got != "LightingOpt" {
		t.Fatalf("TagCanonString(CanonLightingOpt) = %q, want %q", got, "LightingOpt")
	}

	if got := TagCanonString(tag.ID(CanonLensInfo)); got != "LensInfo" {
		t.Fatalf("TagCanonString(CanonLensInfo) = %q, want %q", got, "LensInfo")
	}
}

func TestCanonLensInfoTagID(t *testing.T) {
	if CanonLensInfo != 0x4019 {
		t.Fatalf("CanonLensInfo = 0x%04x, want 0x4019", uint16(CanonLensInfo))
	}
}

func TestTagCanonStringUnknownFallback(t *testing.T) {
	id := tag.ID(0x1234)
	if got := TagCanonString(id); got != id.String() {
		t.Fatalf("TagCanonString(0x1234) = %q, want %q", got, id.String())
	}
}
