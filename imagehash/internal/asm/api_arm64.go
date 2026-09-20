//go:build arm64

package asm

// DCT2DHash64 computes the PHash64 2D DCT. It transforms input in place and
// returns the top-left 8x8 block.
func DCT2DHash64(input []float32) [64]float32 { return asmDCT2DHash64(input) }

// YCbCrToGray converts a 4:4:4 YCbCr image to grayscale float32 pixels.
func YCbCrToGray(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int) {
	asmYCbCrToGray(pixels, minX, minY, maxX, maxY, sY, sCb, sCr, yStride, cStride)
}

// RGBAToGray converts an RGBA image to grayscale float32 pixels.
func RGBAToGray(pixels []float32, pix []uint8, stride, minX, minY, width, height int) {
	asmRGBAtoGray(pixels, pix, stride, minX, minY, width, height)
}

// BlurRow computes the per-row BlurHash basis sums for a 64-wide image row.
func BlurRow(out, lr, lg, lb, xvaluesT []float32) {
	asmBlurRow(out, lr, lg, lb, xvaluesT)
}

// PDQLuma444Row computes PDQ luminance for len(out) 1:1-mapped YCbCr pixels.
// len(out) must be a positive multiple of 8.
func PDQLuma444Row(out []float32, y, cb, cr []uint8) {
	asmPDQLuma444Row(out, y, cb, cr)
}

// PDQLuma420Row computes PDQ luminance for len(out) Y pixels sharing
// len(out)/2 Cb and Cr pixels. len(out) must be a positive multiple of 8.
func PDQLuma420Row(out []float32, y, cb, cr []uint8) {
	asmPDQLuma420Row(out, y, cb, cr)
}
