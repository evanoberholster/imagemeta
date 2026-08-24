//go:build !amd64

package transforms32

import "image"

// asmDCT2DHash64 is the pure Go fallback for non-amd64 architectures.
func asmDCT2DHash64(input []float32) [64]float32 {
	var flattens [64]float32
	if len(input) != 64*64 {
		panic("Incorrect forward transform size")
	}
	for i := 0; i < 64; i++ { // height
		forwardDCT64(input[i*64 : 64*i+64])
	}

	var row [64]float32
	for i := 0; i < 8; i++ { // width
		for j := 0; j < 64; j++ {
			row[j] = (input)[64*j+i]
		}
		forwardDCT64(row[:])
		for j := 0; j < 8; j++ {
			flattens[8*j+i] = row[j]
		}
	}
	return flattens
}

// asmForwardDCT64 is the fallback for non-amd64 architectures.
func asmForwardDCT64(input []float32) {
	forwardDCT64(input)
}

// asmForwardDCT256 is the fallback for non-amd64 architectures.
func asmForwardDCT256(input []float32) {
	forwardDCT256(input)
}

// asmYCbCrToGray is the fallback for non-amd64 architectures.
func asmYCbCrToGray(pixels []float32, minX int, minY int, maxX int, maxY int, sY []uint8, sCb []uint8, sCr []uint8, yStride int, cStride int) {
	c := &image.YCbCr{
		Y:       sY,
		Cb:      sCb,
		Cr:      sCr,
		YStride: yStride,
		CStride: cStride,
		Rect:    image.Rect(minX, minY, maxX, maxY),
	}
	yCbCrToGrayAlt(c, pixels)
}
