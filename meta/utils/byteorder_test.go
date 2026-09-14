package utils

import "testing"

func TestBinaryOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		buf  []byte
		want ByteOrder
	}{
		{name: "classic little endian", buf: []byte{'I', 'I', '*', 0}, want: LittleEndian},
		{name: "classic big endian", buf: []byte{'M', 'M', 0, '*'}, want: BigEndian},
		{name: "BigTIFF little endian", buf: []byte{'I', 'I', '+', 0}, want: LittleEndian},
		{name: "BigTIFF big endian", buf: []byte{'M', 'M', 0, '+'}, want: BigEndian},
		{name: "invalid magic", buf: []byte{'I', 'I', 0, '*'}, want: UnknownEndian},
		{name: "short", buf: []byte{'I', 'I', '*'}, want: UnknownEndian},
		{name: "empty", want: UnknownEndian},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := BinaryOrder(tt.buf); got != tt.want {
				t.Fatalf("BinaryOrder(%v) = %v, want %v", tt.buf, got, tt.want)
			}
		})
	}
}
