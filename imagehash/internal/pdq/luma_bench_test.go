package pdq

import (
	"image"
	"math/rand"
	"testing"
)

// BenchmarkLumaPath isolates the SIMD vs scalar luminance conversion on a
// 512x512 420 image. It toggles the dispatch flag; the caller restores it.
func BenchmarkLumaPath(b *testing.B) {
	rng := rand.New(rand.NewSource(3))
	mkimg := func() *image.YCbCr {
		y := make([]uint8, 512*512)
		cb := make([]uint8, 256*256)
		cr := make([]uint8, 256*256)
		for i := range y {
			y[i] = uint8(rng.Intn(256))
		}
		for i := range cb {
			cb[i] = uint8(rng.Intn(256))
			cr[i] = uint8(rng.Intn(256))
		}
		return &image.YCbCr{Y: y, Cb: cb, Cr: cr, YStride: 512, CStride: 256,
			SubsampleRatio: image.YCbCrSubsampleRatio420, Rect: image.Rect(0, 0, 512, 512)}
	}
	img := mkimg()
	out := make([]float32, 512*512)
	saved := useASMLuma

	b.Run("Scalar", func(b *testing.B) {
		useASMLuma = false
		defer func() { useASMLuma = saved }()
		b.SetBytes(512 * 512)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			toLuminanceDirect(img, out, 512, 512)
		}
	})
	b.Run("SIMD", func(b *testing.B) {
		useASMLuma = true
		defer func() { useASMLuma = saved }()
		b.SetBytes(512 * 512)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			toLuminanceDirect(img, out, 512, 512)
		}
	})
}
