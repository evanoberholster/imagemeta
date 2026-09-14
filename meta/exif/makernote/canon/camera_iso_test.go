package canon

import "testing"

func TestCameraISOValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  int16
		want int64
	}{
		{name: "not applicable", raw: 0, want: 0},
		{name: "auto high", raw: 14, want: int64(CameraISOAutoHighSentinel)},
		{name: "auto", raw: 15, want: int64(CameraISOAutoSentinel)},
		{name: "ISO 100", raw: 17, want: 100},
		{name: "encoded ISO", raw: 0x4064, want: 100},
		{name: "sentinel", raw: 0x7fff, want: 0},
		{name: "encoded all bits set", raw: -1, want: 0x3fff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := CameraISOValue(tt.raw); got != tt.want {
				t.Fatalf("CameraISOValue(%d) = %d, want %d", tt.raw, got, tt.want)
			}
		})
	}
}

func TestNewResolvedCameraISOFromRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  int16
		want uint32
	}{
		{name: "auto high", raw: 14, want: CameraISOAutoHighSentinel},
		{name: "auto", raw: 15, want: CameraISOAutoSentinel},
		{name: "ISO 100", raw: 17, want: 100},
		{name: "encoded all bits set", raw: -1, want: 0x3fff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewResolvedCameraISOFromRaw(tt.raw); got != tt.want {
				t.Fatalf("NewResolvedCameraISOFromRaw(%d) = %d, want %d", tt.raw, got, tt.want)
			}
		})
	}
}
