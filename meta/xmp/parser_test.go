package xmp

import (
	"math"
	"testing"
)

func TestParseUintBoundaries(t *testing.T) {
	if got := parseUint8([]byte("255")); got != math.MaxUint8 {
		t.Fatalf("parseUint8(255) = %d, want %d", got, math.MaxUint8)
	}
	if got := parseUint8([]byte("256")); got != 0 {
		t.Fatalf("parseUint8(256) = %d, want 0", got)
	}

	if got := parseUint32([]byte("4294967295")); got != math.MaxUint32 {
		t.Fatalf("parseUint32(max) = %d, want %d", got, uint32(math.MaxUint32))
	}
	if got := parseUint32([]byte("4294967296")); got != 0 {
		t.Fatalf("parseUint32(overflow) = %d, want 0", got)
	}

	if got := parseUint8([]byte("+1")); got != 0 {
		t.Fatalf("parseUint8(+1) = %d, want 0", got)
	}
	if got := parseUint32([]byte("12x")); got != 0 {
		t.Fatalf("parseUint32(12x) = %d, want 0", got)
	}
}

func TestParseApexConversions(t *testing.T) {
	if got := parseApexAperture([]byte("8/1")); math.Abs(got-16.0) > 0.0001 {
		t.Fatalf("parseApexAperture(8/1) = %.8f, want 16.0", got)
	}

	if got := parseApexShutterSpeed([]byte("-378512/1000000")); math.Abs(got-1.30000033948) > 0.0001 {
		t.Fatalf("parseApexShutterSpeed(-378512/1000000) = %.8f, want 1.30000033948", got)
	}
}
