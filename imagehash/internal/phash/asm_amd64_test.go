//go:build amd64

package phash

import (
	"image"
	"math/rand"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// TestAsmYCbCrToGrayAMD64 checks the x86 AVX2 kernel against the portable
// implementation. Results must be bit-identical. It also runs on darwin/amd64
// via Rosetta, which guards the hand-maintained asm_amd64.s (avo is not a
// build dependency in this module).
func TestAsmYCbCrToGrayAMD64(t *testing.T) {
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
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("trial %d: index %d asm %v != go %v", trial, i, got[i], want[i])
			}
		}
	}
}
