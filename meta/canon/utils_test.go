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

func TestPointsInFocusAFInfo2Decode(t *testing.T) {
	af := make([]uint16, 38)
	af[3] = 7 // ValidAFPoints / NumAFPoints (AFInfo2)

	// AFInfo2 in-focus mask starts at 8 + 4*7 = 36.
	af[36] = (1 << 1) | (1 << 15)
	af[37] = 1 << 0 // selected mask

	inFocus, selected, err := PointsInFocus(af)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inFocus) != 2 || inFocus[0] != 1 || inFocus[1] != 15 {
		t.Fatalf("unexpected inFocus bits: %v", inFocus)
	}
	if len(selected) != 1 || selected[0] != 0 {
		t.Fatalf("unexpected selected bits: %v", selected)
	}
}

func TestPointsInFocusAFInfoDecode(t *testing.T) {
	af := make([]uint16, 24)
	af[0] = 5 // NumAFPoints (AFInfo)

	// AFInfo in-focus mask starts at 8 + 2*5 = 18.
	af[18] = (1 << 0) | (1 << 4)

	inFocus, selected, err := PointsInFocus(af)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inFocus) != 2 || inFocus[0] != 0 || inFocus[1] != 4 {
		t.Fatalf("unexpected inFocus bits: %v", inFocus)
	}
	if selected != nil {
		t.Fatalf("expected nil selected bits for AFInfo, got %v", selected)
	}
}
