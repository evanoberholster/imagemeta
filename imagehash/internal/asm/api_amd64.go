//go:build amd64

package asm

// DCT2DHash64 computes the PHash64 2D DCT. It transforms input in place and
// returns the top-left 8x8 block.
func DCT2DHash64(input []float32) [64]float32 { return asmDCT2DHash64(input) }

// ForwardDCT64 applies a forward DCT to a 64 element slice in place.
func ForwardDCT64(input []float32) { asmForwardDCT64(input) }

// ForwardDCT256 applies a forward DCT to a 256 element slice in place.
func ForwardDCT256(input []float32) { asmForwardDCT256(input) }

// YCbCrToGray converts a YCbCr image to grayscale float32 pixels.
func YCbCrToGray(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int) {
	asmYCbCrToGray(pixels, minX, minY, maxX, maxY, sY, sCb, sCr, yStride, cStride)
}
