//go:build arm64

package asm

// asmDCT2DHash64 is the NEON implementation of DCT2DHash64. It returns the
// top-left 8x8 DCT block and transforms input in place.
//
//go:noescape
func asmDCT2DHash64(input []float32) [64]float32
