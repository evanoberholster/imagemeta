// Code generated by command: go run asm.go -out asm.s -stubs stub.go. DO NOT EDIT.

package transforms32

// asmForwardDCT64 is a forward DCT transform for [64]float32
func asmForwardDCT64(input []float32)

// asmForwardDCT32 is a forward DCT transform for [32]float32
func asmForwardDCT32(input []float32)

// asmForwardDCT16 is a forward DCT transform for [16]float32
func asmForwardDCT16(input []float32)

// AsmYCbCrToGray is a forward DCT transform for []float32
func AsmYCbCrToGray(pixels []float32, minX int, minY int, maxX int, maxY int, sY []uint8, sCb []uint8, sCr []uint8, yStride int, cStride int) uint64

// AsmYCbCrToGray8 is a forward DCT transform for []float32
func AsmYCbCrToGray8(pixels []float32, minX int, minY int, maxX int, maxY int, sY []uint8, sCb []uint8, sCr []uint8, yStride int, cStride int) uint64
