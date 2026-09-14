//go:build !amd64 && !arm64

package asm

import "errors"

// errUnsupported is reported if an assembly kernel is called on an
// architecture without an implementation. Callers only select these functions
// when a kernel is present, so this is a guard rather than an expected path.
var errUnsupported = errors.New("asm: no kernel for this architecture")

// DCT2DHash64 is a guard; see errUnsupported.
func DCT2DHash64(input []float32) [64]float32 { panic(errUnsupported) }

// ForwardDCT64 is a guard; see errUnsupported.
func ForwardDCT64(input []float32) { panic(errUnsupported) }

// ForwardDCT256 is a guard; see errUnsupported.
func ForwardDCT256(input []float32) { panic(errUnsupported) }

// YCbCrToGray is a guard; see errUnsupported.
func YCbCrToGray(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int) {
	panic(errUnsupported)
}

// RGBAToGray is a guard; see errUnsupported.
func RGBAToGray(pixels []float32, pix []uint8, stride, minX, minY, width, height int) {
	panic(errUnsupported)
}

// BlurRow is a guard; see errUnsupported.
func BlurRow(out, lr, lg, lb, xvaluesT []float32) { panic(errUnsupported) }
