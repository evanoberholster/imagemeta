//go:build arm64

package asm

// asmYCbCrToGray converts an *image.YCbCr with 4:4:4 chroma to grayscale
// float32 pixels. The caller must pass a width that is a multiple of 8.
//
//go:noescape
func asmYCbCrToGray(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)

// asmRGBAtoGray converts an *image.RGBA to grayscale float32 pixels. The caller
// must pass a width that is a multiple of 8.
//
//go:noescape
func asmRGBAtoGray(pixels []float32, pix []uint8, stride, minX, minY, width, height int)
