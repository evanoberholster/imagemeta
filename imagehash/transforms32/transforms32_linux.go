//go:build linux && amd64

package transforms32

import (
	"image"

	cpu "github.com/klauspost/cpuid/v2"
)

func init() {
	FlagUseASM = cpu.CPU.Supports(cpu.AVX, cpu.AVX2, cpu.SSE, cpu.SSE2, cpu.SSE4)
	if FlagUseASM {
		ForwardDCT256 = asmForwardDCT256
		ForwardDCT64 = asmForwardDCT64
		YCbCrToGray = AsmYCbCrToGray
	}
}

func AsmYCbCrToGray(c *image.YCbCr, pixels []float32) {
	asmYCbCrToGray(pixels,
		c.Rect.Min.X, c.Rect.Min.Y, c.Rect.Max.X, c.Rect.Max.Y,
		c.Y, c.Cb, c.Cr, c.YStride, c.CStride)
}
