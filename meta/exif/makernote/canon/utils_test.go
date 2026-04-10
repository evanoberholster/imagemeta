package canon

import "testing"

func s16u(v int16) uint16 {
	return uint16(v)
}

func TestPointsInFocusShortInput(t *testing.T) {
	_, _, err := PointsInFocus([]uint16{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for short AF input")
	}
}

func TestPointsInFocusUnknownCount(t *testing.T) {
	af := make([]uint16, 4)
	af[0] = 99

	_, _, err := PointsInFocus(af)
	if err == nil {
		t.Fatal("expected error for unsupported NumAFPoints")
	}
}

func TestPointsInFocusBoundsCheck(t *testing.T) {
	af := make([]uint16, 10)
	af[0] = 7 // requires more than 10 words

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
	af[0] = 7 // requires more than 10 words

	if got := ParseAFPoints(af); got != nil {
		t.Fatalf("expected nil for truncated AF payload, got %v", got)
	}
}

func TestPointsInFocusAFInfo2Decode(t *testing.T) {
	af := make([]uint16, 38)
	af[2] = 7 // NumAFPoints (AFInfo2)
	af[3] = 5 // ValidAFPoints; should not affect mask offsets

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

func TestPointsInFocusAFInfo2NonEOSShapeSkipsSelected(t *testing.T) {
	af := make([]uint16, 40)
	af[2] = 7
	af[3] = 7

	// seq 12 in-focus mask
	af[36] = 1 << 5
	// seq 13 unknown[count+1] followed by seq 14 primary point in non-EOS AFInfo2
	af[37] = 1 << 2
	af[39] = 6

	inFocus, selected, err := PointsInFocus(af)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inFocus) != 1 || inFocus[0] != 5 {
		t.Fatalf("unexpected inFocus bits: %v", inFocus)
	}
	if selected != nil {
		t.Fatalf("expected nil selected bits for non-EOS AFInfo2 shape, got %v", selected)
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

func TestParseAFPointsAFInfo2UsesNumAFPoints(t *testing.T) {
	af := make([]uint16, 38)
	af[2] = 7
	af[3] = 5
	af[4] = 100
	af[5] = 80
	for i := 0; i < 7; i++ {
		af[8+i] = 4
		af[15+i] = 6
		af[22+i] = uint16(i)
		af[29+i] = uint16(i)
	}

	got := ParseAFPoints(af)
	if len(got) != 7 {
		t.Fatalf("len(ParseAFPoints) = %d, want 7", len(got))
	}
	if got[0] != NewAFPoint(4, 6, 48, 37) {
		t.Fatalf("first AF point = %v, want [4 6 48 37]", got[0])
	}
}

func TestParseAFPointsLegacyAFInfo(t *testing.T) {
	af := make([]uint16, 19)
	af[0] = 5
	af[6] = 10
	af[7] = 12
	af[8] = s16u(-2)
	af[9] = s16u(-1)
	af[10] = 0
	af[11] = 1
	af[12] = 2
	af[13] = 3
	af[14] = 2
	af[15] = 1
	af[16] = 0
	af[17] = s16u(-1)

	got := ParseAFPoints(af)
	want := []AFPoint{
		NewAFPoint(10, 12, -2, 3),
		NewAFPoint(10, 12, -1, 2),
		NewAFPoint(10, 12, 0, 1),
		NewAFPoint(10, 12, 1, 0),
		NewAFPoint(10, 12, 2, -1),
	}
	if len(got) != len(want) {
		t.Fatalf("len(ParseAFPoints) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseAFPoints[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestParseAFAreaAFInfo2(t *testing.T) {
	af := make([]uint16, 38)
	af[2] = 7
	af[3] = 7
	for i := 0; i < 7; i++ {
		af[8+i] = 4
		af[15+i] = 6
		af[22+i] = uint16(i)
		af[29+i] = uint16(i + 10)
	}

	got := ParseAFArea(af)
	if len(got) != 7 {
		t.Fatalf("len(ParseAFArea) = %d, want 7", len(got))
	}
	if got[0] != NewAFPoint(4, 6, 0, 10) {
		t.Fatalf("first AFArea point = %v, want [4 6 0 10]", got[0])
	}
}
