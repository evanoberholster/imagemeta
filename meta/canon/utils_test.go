package canon

import "testing"

func TestPointsInFocusShortInput(t *testing.T) {
	_, _, err := PointsInFocus([]uint16{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for short AF input")
	}
}

func TestPointsInFocusUnknownCount(t *testing.T) {
	af := make([]uint16, 4)
	af[3] = 99 // unsupported NumAFPoints

	_, _, err := PointsInFocus(af)
	if err == nil {
		t.Fatal("expected error for unsupported NumAFPoints")
	}
}

func TestPointsInFocusBoundsCheck(t *testing.T) {
	af := make([]uint16, 10)
	af[3] = 7 // requires more than 10 words

	_, _, err := PointsInFocus(af)
	if err == nil {
		t.Fatal("expected bounds error for truncated AF payload")
	}
}

func TestParseAFPointsShortInput(t *testing.T) {
	if got := ParseAFPoints([]uint16{1, 2, 3}); got != nil {
		t.Fatalf("expected nil for short AF input, got %v", got)
	}
}

func TestParseAFPointsBoundsCheck(t *testing.T) {
	af := make([]uint16, 10)
	af[3] = 7 // requires more than 10 words

	if got := ParseAFPoints(af); got != nil {
		t.Fatalf("expected nil for truncated AF payload, got %v", got)
	}
}
