// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"math"
	"testing"
)

func TestJaroszFilter_WrongSize(t *testing.T) {
	_, err := JaroszFilter(make([]float32, 100), 20, 20)
	if err != nil {
		return // expected
	}
	_, err = JaroszFilter(make([]float32, 100), 5, 5) // 25 != 100
	if err == nil {
		t.Fatal("expected error when src length does not match dimensions")
	}
}

func TestJaroszFilter_OutputSize(t *testing.T) {
	src := make([]float32, ImageSize*ImageSize)
	out, err := JaroszFilter(src, ImageSize, ImageSize)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 64*64 {
		t.Errorf("output length = %d, want %d", len(out), 64*64)
	}
}

func TestJaroszFilter_SmallImageOutputSize(t *testing.T) {
	// Small images (< 512) must also produce 64x64 output.
	src := make([]float32, 100*80)
	out, err := JaroszFilter(src, 100, 80)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 64*64 {
		t.Errorf("output length = %d, want %d", len(out), 64*64)
	}
}

func TestJaroszFilter_UniformInput(t *testing.T) {
	const value = 128.0
	src := make([]float32, ImageSize*ImageSize)
	for i := range src {
		src[i] = value
	}
	out, err := JaroszFilter(src, ImageSize, ImageSize)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range out {
		if math.Abs(float64(v-value)) > 1e-3 {
			t.Errorf("out[%d] = %v, want %v", i, v, value)
			break
		}
	}
}

func TestJaroszFilter_OutputInInputRange(t *testing.T) {
	src := make([]float32, ImageSize*ImageSize)
	for i := range src {
		src[i] = float32(i % 256)
	}
	out, err := JaroszFilter(src, ImageSize, ImageSize)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range out {
		if v < -1e-3 || v > 255+1e-3 {
			t.Errorf("out[%d] = %v out of range [0, 255]", i, v)
			break
		}
	}
}

func TestJaroszWindowSize(t *testing.T) {
	cases := []struct {
		dim  int
		want int
	}{
		{512, 4},
		{256, 2},
		{128, 1},
		{64, 1},
		{32, 1},
	}
	for _, tc := range cases {
		if got := jaroszWindowSize(tc.dim); got != tc.want {
			t.Errorf("jaroszWindowSize(%d) = %d, want %d", tc.dim, got, tc.want)
		}
	}
}

func TestBox1D_WriteCount(t *testing.T) {
	in := make([]float32, 16)
	out := make([]float32, 16)
	for i := range in {
		in[i] = float32(i + 1)
	}
	box1D(in, out, 16, 1, 4)
	for i, v := range out {
		if v == 0 {
			t.Errorf("out[%d] was not written (zero)", i)
		}
	}
}
