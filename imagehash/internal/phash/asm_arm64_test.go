//go:build arm64

package phash

import (
	"image"
	"math"
	"math/rand"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// withoutASM runs fn with the assembly path disabled and restores the flag.
func withoutASM(fn func()) {
	orig := FlagUseASM
	FlagUseASM = false
	defer func() { FlagUseASM = orig }()
	fn()
}

// TestAsmDCT2DHash64 verifies the NEON implementation produces the same result
// as the portable Go implementation, both for the returned hash and for the
// in-place transformed input.
func TestAsmDCT2DHash64(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for trial := 0; trial < 200; trial++ {
		src := make([]float32, 64*64)
		for i := range src {
			src[i] = rng.Float32()*200 - 100
		}
		asmIn := append([]float32(nil), src...)
		goIn := append([]float32(nil), src...)

		got := asm.DCT2DHash64(asmIn)

		var want [64]float32
		withoutASM(func() { want = DCT2DHash64(goIn) })

		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("hash mismatch at %d: got %v want %v", i, got[i], want[i])
			}
			if asmIn[i] != goIn[i] {
				t.Fatalf("in-place mismatch at %d: got %v want %v", i, asmIn[i], goIn[i])
			}
		}
	}
}

func TestAsmDCT2DHash64Zeros(t *testing.T) {
	in := make([]float32, 64*64)
	got := asm.DCT2DHash64(in)
	for i, v := range got {
		if v != 0 {
			t.Fatalf("expected zero at %d, got %v", i, v)
		}
	}
}

func benchSource() []float32 {
	src := make([]float32, 64*64)
	for i := range src {
		src[i] = float32(math.Sin(float64(i)))
	}
	return src
}

func BenchmarkAsmDCT2DHash64(b *testing.B) {
	src := benchSource()
	buf := make([]float32, len(src))
	b.SetBytes(int64(len(src) * 4))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, src)
		_ = asm.DCT2DHash64(buf)
	}
}

func BenchmarkGoDCT2DHash64(b *testing.B) {
	src := benchSource()
	buf := make([]float32, len(src))
	b.SetBytes(int64(len(src) * 4))
	b.ReportAllocs()
	withoutASM(func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			copy(buf, src)
			_ = DCT2DHash64(buf)
		}
	})
}

func floatAbs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}

func TestAsmYCbCrToGray(t *testing.T) {
	const s = 64
	for trial := 0; trial < 50; trial++ {
		rng := rand.New(rand.NewSource(int64(trial)))
		img := image.NewYCbCr(image.Rect(0, 0, s, s), image.YCbCrSubsampleRatio444)
		for i := range img.Y {
			img.Y[i] = uint8(rng.Intn(256))
		}
		for i := range img.Cb {
			img.Cb[i] = uint8(rng.Intn(256))
		}
		for i := range img.Cr {
			img.Cr[i] = uint8(rng.Intn(256))
		}
		got := make([]float32, s*s)
		want := make([]float32, s*s)
		asm.YCbCrToGray(got, img.Rect.Min.X, img.Rect.Min.Y, img.Rect.Max.X, img.Rect.Max.Y,
			img.Y, img.Cb, img.Cr, img.YStride, img.CStride)
		yCbCrToGray(img, want)
		maxd := float32(0)
		for i := range got {
			if d := floatAbs(got[i] - want[i]); d > maxd {
				maxd = d
			}
		}
		if maxd > 2.0 {
			t.Fatalf("trial %d: max grey diff %v", trial, maxd)
		}
	}
}

func TestAsmRGBAtoGray(t *testing.T) {
	const s = 64
	rng := rand.New(rand.NewSource(11))
	img := image.NewRGBA(image.Rect(0, 0, s, s))
	for i := range img.Pix {
		img.Pix[i] = uint8(rng.Intn(256))
	}
	got := make([]float32, s*s)
	want := make([]float32, s*s)
	asm.RGBAToGray(got, img.Pix, img.Stride, img.Rect.Min.X, img.Rect.Min.Y, img.Rect.Dx(), img.Rect.Dy())
	rgbaToGray(img, want)
	maxd := float32(0)
	for i := range got {
		if d := floatAbs(got[i] - want[i]); d > maxd {
			maxd = d
		}
	}
	if maxd > 0.05 {
		t.Fatalf("max grey diff %v", maxd)
	}
}
