// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"math"
	"math/rand"
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

// TestBoxColsStreamMatchesStrided proves the streaming column filter is
// bit-identical to the strided loop over varied shapes, windows and data
// (including clamping extremes).
func TestBoxColsStreamMatchesStrided(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(7))
	for _, dims := range [][2]int{{1, 1}, {8, 8}, {64, 64}, {100, 80}, {80, 100}, {512, 512}, {7, 511}, {511, 7}} {
		rows, cols := dims[0], dims[1]
		src := make([]float32, rows*cols)
		for i := range src {
			switch i % 4 {
			case 0:
				src[i] = float32(rng.Intn(256))
			case 1:
				src[i] = rng.Float32() * 255
			case 2:
				src[i] = 0
			default:
				src[i] = 255
			}
		}
		for _, window := range []int{1, 2, 4, 8} {
			if window > rows {
				continue
			}
			want := make([]float32, rows*cols)
			got := make([]float32, rows*cols)
			// Strided reference: one box1D per column (the pre-streaming
			// algorithm boxAlongCols used unconditionally).
			for x := 0; x < cols; x++ {
				box1D(src[x:], want[x:], rows, cols, window)
			}
			boxColsStream(src, got, rows, cols, window)
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("dims=%dx%d window=%d: out[%d] = %v, want %v, src %v",
						rows, cols, window, i, got[i], want[i], src[i])
				}
			}
		}
	}
}

// TestBoxColsWideFallback exercises the strided fallback for inputs wider
// than the streaming scratch buffer.
func TestBoxColsWideFallback(t *testing.T) {
	t.Parallel()
	const rows, cols = 8, maxStreamCols + 8
	src := make([]float32, rows*cols)
	for i := range src {
		src[i] = float32(i % 251)
	}
	dst := make([]float32, rows*cols)
	boxAlongCols(src, dst, rows, cols, 4)
	for i, v := range dst {
		if v == 0 {
			t.Fatalf("dst[%d] was not written (zero)", i)
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
