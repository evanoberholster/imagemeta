package xmp

import (
	"testing"
	"time"
)

func TestParseFastDateCachesMinuteOffsetLocations(t *testing.T) {
	// from meta/xmp/test/canon_eos_r6_cr3_lightroom_blue_label_baseline.xmp
	a, err := parseFastDate([]byte("2024-11-02T12:35:44.40-04:00"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	// from meta/xmp/test/canon_eos_r6_cr3_lightroom_custom_adjustments.xmp (same timezone offset)
	b, err := parseFastDate([]byte("2024-11-02T12:35:43.99-04:00"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	if a.Location() != b.Location() {
		t.Fatal("expected identical cached location pointers for -04:00")
	}

	_, offset := a.Zone()
	if offset != -4*3600 {
		t.Fatalf("zone offset = %d, want %d", offset, -4*3600)
	}
}

func TestParseFastDateUTCOffsetUsesTimeUTC(t *testing.T) {
	// from meta/xmp/test/jpeg.xmp
	got, err := parseFastDate([]byte("2003-02-04T08:06:56Z"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	if got.Location() != time.UTC {
		t.Fatal("expected time.UTC location pointer for Z offset")
	}
}

func TestParseFastDateWithFallbackHook(t *testing.T) {
	calls := 0
	_, err := parseFastDateWithFallback([]byte("not-a-date"), func(string) (time.Time, error) {
		calls++
		return time.Time{}, nil
	})
	if err != nil {
		t.Fatalf("parseFastDateWithFallback returned error: %v", err)
	}

	if calls != 1 {
		t.Fatalf("fallback called %d times, want 1", calls)
	}
}
