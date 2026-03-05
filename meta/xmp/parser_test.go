package xmp

import (
	"math"
	"testing"
	"time"
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

	if got := parseUint16([]byte("65535")); got != math.MaxUint16 {
		t.Fatalf("parseUint16(65535) = %d, want %d", got, math.MaxUint16)
	}
	if got := parseUint16([]byte("65536")); got != 0 {
		t.Fatalf("parseUint16(65536) = %d, want 0", got)
	}
}

func TestParseIntBoundaries(t *testing.T) {
	if got := parseInt16([]byte("-32768")); got != math.MinInt16 {
		t.Fatalf("parseInt16(-32768) = %d, want %d", got, math.MinInt16)
	}
	if got := parseInt16([]byte("32767")); got != math.MaxInt16 {
		t.Fatalf("parseInt16(32767) = %d, want %d", got, math.MaxInt16)
	}
	if got := parseInt16([]byte("32768")); got != 0 {
		t.Fatalf("parseInt16(32768) = %d, want 0", got)
	}
	if got := parseInt32([]byte("+42")); got != 42 {
		t.Fatalf("parseInt32(+42) = %d, want 42", got)
	}
	if got := parseInt32([]byte("42x")); got != 0 {
		t.Fatalf("parseInt32(42x) = %d, want 0", got)
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

func TestParseRationalBounds(t *testing.T) {
	n, d := parseRational([]byte("4294967296/1"))
	if n != 0 || d != 1 {
		t.Fatalf("parseRational(overflow numerator) = %d/%d, want 0/1", n, d)
	}

	n, d = parseRational([]byte("1/4294967296"))
	if n != 0 || d != 1 {
		t.Fatalf("parseRational(overflow denominator) = %d/%d, want 0/1", n, d)
	}
}

func TestParseFastDateXMPFractional(t *testing.T) {
	got, err := parseFastDate([]byte("2021-01-10T17:30:57.00"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	want := time.Date(2021, time.January, 10, 17, 30, 57, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("parseFastDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
	}
}

func TestParseFastDateISOTimezoneAndFraction(t *testing.T) {
	got, err := parseFastDate([]byte("2021-01-10T17:30:57.123456789+02:30"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	want := time.Date(2021, time.January, 10, 17, 30, 57, 123456789, time.FixedZone("", 2*3600+30*60))
	if !got.Equal(want) {
		t.Fatalf("parseFastDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
	}
}

func TestParseFastDateExifTimezone(t *testing.T) {
	got, err := parseFastDate([]byte("2021:01:10 17:30:57-05:00"))
	if err != nil {
		t.Fatalf("parseFastDate returned error: %v", err)
	}

	want := time.Date(2021, time.January, 10, 17, 30, 57, 0, time.FixedZone("", -5*3600))
	if !got.Equal(want) {
		t.Fatalf("parseFastDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
	}
}

func TestParseFastDate_ISOWithoutFallback(t *testing.T) {
	orig := parseDateStringFallback
	defer func() { parseDateStringFallback = orig }()

	fallbackCalls := 0
	parseDateStringFallback = func(s string) (time.Time, error) {
		fallbackCalls++
		return orig(s)
	}

	tests := []struct {
		in         string
		want       time.Time
		wantOffset int
	}{
		{
			in:         "2003-02-04T08:06:56Z",
			want:       time.Date(2003, time.February, 4, 8, 6, 56, 0, time.UTC),
			wantOffset: 0,
		},
		{
			in:         "2007-08-16T11:57:04+01:00",
			want:       time.Date(2007, time.August, 16, 11, 57, 4, 0, time.FixedZone("", 1*3600)),
			wantOffset: 1 * 3600,
		},
		{
			in:         "2021-02-03T17:34:04+08:00",
			want:       time.Date(2021, time.February, 3, 17, 34, 4, 0, time.FixedZone("", 8*3600)),
			wantOffset: 8 * 3600,
		},
	}

	for _, tt := range tests {
		got, err := parseFastDate([]byte(tt.in))
		if err != nil {
			t.Fatalf("parseFastDate(%q) error: %v", tt.in, err)
		}
		if !got.Equal(tt.want) {
			t.Fatalf("parseFastDate(%q) = %s, want %s", tt.in, got.Format(time.RFC3339Nano), tt.want.Format(time.RFC3339Nano))
		}
		_, off := got.Zone()
		if off != tt.wantOffset {
			t.Fatalf("parseFastDate(%q) zone offset = %d, want %d", tt.in, off, tt.wantOffset)
		}
	}

	if fallbackCalls != 0 {
		t.Fatalf("parseDateString fallback called %d times, want 0", fallbackCalls)
	}
}

func TestParseDigitsEquivalentToParse2AndParse4(t *testing.T) {
	twoDigitCases := []string{
		"00",
		"01",
		"42",
		"99",
		"a1",
		"1a",
		"--",
	}
	for _, tc := range twoDigitCases {
		got, ok := parseDigits([]byte(tc))
		want, wantOK := parse2Digits(tc[0], tc[1])
		if ok != wantOK || got != want {
			t.Fatalf("parseDigits(%q) = (%d,%t), parse2Digits = (%d,%t)", tc, got, ok, want, wantOK)
		}
	}

	fourDigitCases := []string{
		"0000",
		"2024",
		"9999",
		"12a4",
		"a234",
		"123a",
	}
	for _, tc := range fourDigitCases {
		got, ok := parseDigits([]byte(tc))
		want, wantOK := parse4Digits(tc[0], tc[1], tc[2], tc[3])
		if ok != wantOK || got != want {
			t.Fatalf("parseDigits(%q) = (%d,%t), parse4Digits = (%d,%t)", tc, got, ok, want, wantOK)
		}
	}
}

func TestParseDigitsShortInput(t *testing.T) {
	got, ok := parseDigits([]byte("9"))
	if ok || got != 0 {
		t.Fatalf("parseDigits(\"9\") = (%d,%t), want (0,false)", got, ok)
	}
}
