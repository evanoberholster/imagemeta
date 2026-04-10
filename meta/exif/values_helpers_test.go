package exif

import (
	"math"
	"testing"
	"time"
)

func TestParseStrUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   []byte
		want uint
	}{
		{in: []byte("123"), want: 123},
		{in: []byte("a1b2c3"), want: 123},
		{in: []byte("00"), want: 0},
		{in: []byte(""), want: 0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.in), func(t *testing.T) {
			t.Parallel()
			if got := parseStrUint(tt.in); got != tt.want {
				t.Fatalf("parseStrUint(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestTrimNULBuffer(t *testing.T) {
	t.Parallel()

	if got := trimNULBuffer([]byte{'a', 'b', 0, ' ', '\n'}); string(got) != "ab" {
		t.Fatalf("trimNULBuffer() = %q, want %q", got, "ab")
	}
	if got := trimNULBuffer([]byte{0, ' ', '\n'}); got != nil {
		t.Fatalf("trimNULBuffer(all-trim) = %v, want nil", got)
	}
}

func TestTrimTrailingNULBytes(t *testing.T) {
	t.Parallel()

	if got := trimTrailingNULBytes([]byte{'a', 'b', 0, 0}); string(got) != "ab" {
		t.Fatalf("trimTrailingNULBytes() = %q, want %q", got, "ab")
	}
	if got := trimTrailingNULBytes([]byte{'a', 'b', ' ', 0}); string(got) != "ab " {
		t.Fatalf("trimTrailingNULBytes() = %q, want %q", got, "ab ")
	}
	if got := trimTrailingNULBytes([]byte{0, 0}); got != nil {
		t.Fatalf("trimTrailingNULBytes(all-nul) = %v, want nil", got)
	}
}

func TestTrimRightSpaceNewline(t *testing.T) {
	t.Parallel()

	if got := trimRightSpaceNewline([]byte("Artist Name   \n\n")); string(got) != "Artist Name" {
		t.Fatalf("trimRightSpaceNewline() = %q, want %q", got, "Artist Name")
	}
	if got := trimRightSpaceNewline([]byte("Artist Name\x00")); string(got) != "Artist Name" {
		t.Fatalf("trimRightSpaceNewline() = %q, want %q", got, "Artist Name")
	}
	if got := trimRightSpaceNewline([]byte("      \x00")); len(got) != 0 {
		t.Fatalf("trimRightSpaceNewline() = %q, want empty", got)
	}
	if got := trimRightSpaceNewline([]byte("Copyright\r\n")); string(got) != "Copyright" {
		t.Fatalf("trimRightSpaceNewline() = %q, want %q", got, "Copyright")
	}
	if got := trimRightSpaceNewline([]byte("")); len(got) != 0 {
		t.Fatalf("trimRightSpaceNewline(empty) = %q, want empty", got)
	}
}

func TestRationalDuration(t *testing.T) {
	t.Parallel()

	if got, want := rationalDuration(1, 2, time.Second), 500*time.Millisecond; got != want {
		t.Fatalf("rationalDuration(1/2 sec) = %v, want %v", got, want)
	}
	if got, want := rationalDuration(3, 2, time.Minute), 90*time.Second; got != want {
		t.Fatalf("rationalDuration(3/2 min) = %v, want %v", got, want)
	}
	if got := rationalDuration(1, 0, time.Second); got != 0 {
		t.Fatalf("rationalDuration(den=0) = %v, want 0", got)
	}
	if got := rationalDuration(math.MaxUint32, 1, time.Hour); got != time.Duration(math.MaxInt64) {
		t.Fatalf("rationalDuration(overflow) = %v, want %v", got, time.Duration(math.MaxInt64))
	}
}
