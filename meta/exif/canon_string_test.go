package exif

import "testing"

func TestCanonString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  []byte
		want string
	}{
		{
			name: "ascii preserved",
			raw:  []byte("EF70-200"),
			want: "EF70-200",
		},
		{
			name: "nul terminated",
			raw:  []byte{'E', 'F', 0, 'X'},
			want: "EF",
		},
		{
			name: "non printable converted and edge dots trimmed",
			raw:  []byte{0x01, 'A', 'B', 0x7f},
			want: "AB",
		},
		{
			name: "inner non printable kept as dot",
			raw:  []byte{'A', 0x01, 'B'},
			want: "A.B",
		},
		{
			name: "space and dot trimmed at edges",
			raw:  []byte("  .EF  "),
			want: "EF",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := canonString(tt.raw); got != tt.want {
				t.Fatalf("canonString(%v) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
