// Code generated by command: go run asm.go -out asm_x86.s -stubs asm_x86.go. DO NOT EDIT.

package transforms32

// asmDCT2DHash64 function returns a result of DCT2D by using the seperable property.
// DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.
// Custom built by Evan Oberholster for Hash64. Returns flattened pixels
//
//go:noescape
func asmDCT2DHash64(input []float32) [64]float32

// asmForwardDCT64 is a forward DCT transform for [64]float32
func asmForwardDCT64(input []float32)

// asmForwardDCT256 is a forward DCT transform for [256]float32
func asmForwardDCT256(input []float32)

// AsmYCbCrToGra converts a YCbCr image to grayscale pixels.
// Converts using 8x SIMD instructions and requires AVX and AVX2
func AsmYCbCrToGray(pixels []float32, minX int, minY int, maxX int, maxY int, sY []uint8, sCb []uint8, sCr []uint8, yStride int, cStride int)
