// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"math"
	"testing"
)

func TestDCT64To16_OutputSize(t *testing.T) {
	input := make([]float32, 64*64)
	out := DCT64To16(input)
	if len(out) != 16*16 {
		t.Errorf("output length = %d, want 256", len(out))
	}
}

func TestDCT64To16_ZeroInput(t *testing.T) {
	out := DCT64To16(make([]float32, 64*64))
	for i, v := range out {
		if v != 0 {
			t.Errorf("out[%d] = %v, want 0 for zero input", i, v)
			break
		}
	}
}

func TestDCT64To16_Linearity(t *testing.T) {
	input := make([]float32, 64*64)
	for i := range input {
		input[i] = float32(i%64+1) * float32(i/64+1)
	}

	out1 := DCT64To16(input)

	scaled := make([]float32, 64*64)
	const k = 3.0
	for i, v := range input {
		scaled[i] = v * k
	}
	out2 := DCT64To16(scaled)

	for i := range out1 {
		want := out1[i] * k
		if math.Abs(float64(out2[i]-want)) > 1e-2 {
			t.Errorf("out[%d]: got %v, want %v (linearity check)", i, out2[i], want)
			break
		}
	}
}

func TestDCT64To16_UniformInput(t *testing.T) {
	input := make([]float32, 64*64)
	for i := range input {
		input[i] = 100.0
	}
	out := DCT64To16(input)
	if out[0] == 0 {
		t.Error("DCT of uniform input: out[0] should be nonzero (DC energy)")
	}

	for i := 1; i < len(out); i++ {
		if math.Abs(float64(out[i])) > 1e-2 {
			t.Errorf("out[%d] = %v, want ~0 for uniform input", i, out[i])
			break
		}
	}
}
