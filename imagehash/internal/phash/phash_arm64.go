//go:build arm64

package phash

import (
	"image"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// init enables the NEON kernels on ARM64.
func init() {
	FlagUseASM = true
	YCbCrToGray = yCbCrToGrayASM
	RGBAtoGray = rgbaToGrayASM
}

// yCbCrToGrayASM uses the NEON kernel for 4:4:4 images whose width is a
// multiple of 8, and falls back to the portable implementation otherwise.
func yCbCrToGrayASM(img *image.YCbCr, pixels []float32) {
	if img.SubsampleRatio == image.YCbCrSubsampleRatio444 && img.Rect.Dx()%8 == 0 {
		asm.YCbCrToGray(pixels,
			img.Rect.Min.X, img.Rect.Min.Y, img.Rect.Max.X, img.Rect.Max.Y,
			img.Y, img.Cb, img.Cr, img.YStride, img.CStride)
		return
	}
	yCbCrToGray(img, pixels)
}

// rgbaToGrayASM uses the NEON kernel when the width is a multiple of 8, and
// falls back to the portable implementation otherwise.
func rgbaToGrayASM(img *image.RGBA, pixels []float32) {
	bounds := img.Bounds()
	if bounds.Dx()%8 == 0 {
		asm.RGBAToGray(pixels, img.Pix, img.Stride, bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
		return
	}
	rgbaToGray(img, pixels)
}
