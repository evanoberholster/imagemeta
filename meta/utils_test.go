package meta

import (
	"math"
	"testing"
)

func TestSafecastIntToUint32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value int
		want  uint32
		ok    bool
	}{
		{name: "negative", value: -1},
		{name: "zero", value: 0, want: 0, ok: true},
		{name: "positive", value: 123, want: 123, ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := SafecastIntToUint32(tt.value)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("SafecastIntToUint32(%d) = (%d, %t), want (%d, %t)", tt.value, got, ok, tt.want, tt.ok)
			}
		})
	}

	maxInt := int(^uint(0) >> 1)
	if uint64(maxInt) < math.MaxUint32 {
		return
	}

	maxUint32 := uint64(math.MaxUint32)
	got, ok := SafecastIntToUint32(int(maxUint32))
	if got != math.MaxUint32 || !ok {
		t.Fatalf("SafecastIntToUint32(MaxUint32) = (%d, %t), want (%d, true)", got, ok, uint32(math.MaxUint32))
	}

	got, ok = SafecastIntToUint32(int(maxUint32 + 1))
	if got != 0 || ok {
		t.Fatalf("SafecastIntToUint32(MaxUint32 + 1) = (%d, %t), want (0, false)", got, ok)
	}
}

func TestSafecastUintToInt(t *testing.T) {
	t.Parallel()

	maxInt := ^uint(0) >> 1
	tests := []struct {
		name  string
		value uint
		want  int
		ok    bool
	}{
		{name: "zero", value: 0, want: 0, ok: true},
		{name: "max int", value: maxInt, want: int(maxInt), ok: true},
		{name: "max int plus one", value: maxInt + 1},
		{name: "max uint", value: ^uint(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := SafecastUintToInt(tt.value)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("SafecastUintToInt(%d) = (%d, %t), want (%d, %t)", tt.value, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestSafecastUint64ToInt(t *testing.T) {
	t.Parallel()

	maxInt := uint64(^uint(0) >> 1)
	tests := []struct {
		name  string
		value uint64
		want  int
		ok    bool
	}{
		{name: "zero", value: 0, want: 0, ok: true},
		{name: "max int", value: maxInt, want: int(maxInt), ok: true},
		{name: "max int plus one", value: maxInt + 1},
		{name: "max uint64", value: math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := SafecastUint64ToInt(tt.value)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("SafecastUint64ToInt(%d) = (%d, %t), want (%d, %t)", tt.value, got, ok, tt.want, tt.ok)
			}
		})
	}
}
