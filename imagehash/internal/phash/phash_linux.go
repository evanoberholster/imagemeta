//go:build linux && amd64

package phash

import (
	"image"

	cpu "github.com/klauspost/cpuid/v2"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

func init() {
	FlagUseASM = cpu.CPU.Supports(cpu.AVX, cpu.AVX2, cpu.SSE, cpu.SSE2, cpu.SSE4)
	if FlagUseASM {
		ForwardDCT256 = asm.ForwardDCT256
		ForwardDCT64 = asm.ForwardDCT64
		YCbCrToGray = AsmYCbCrToGray
	}
}

// AsmYCbCrToGray adapts asm.YCbCrToGray to the YCbCrToGray signature.
func AsmYCbCrToGray(c *image.YCbCr, pixels []float32) {
	asm.YCbCrToGray(pixels,
		c.Rect.Min.X, c.Rect.Min.Y, c.Rect.Max.X, c.Rect.Max.Y,
		c.Y, c.Cb, c.Cr, c.YStride, c.CStride)
}
