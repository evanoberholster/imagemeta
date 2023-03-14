// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package transforms32 contains DCT transformations
package transforms32

import (
	"image"
	"math"

	cpu "github.com/klauspost/cpuid/v2"
)

var (
	// FlagUseASM determines if ASM can be used on the current CPU
	FlagUseASM    = false
	ForwardDCT64  = forwardDCT64
	ForwardDCT256 = forwardDCT256
)

func init() {
	FlagUseASM = cpu.CPU.Supports(cpu.AVX, cpu.AVX2, cpu.SSE, cpu.SSE2, cpu.SSE4)
	if FlagUseASM {
		ForwardDCT256 = asmForwardDCT256
		ForwardDCT64 = asmForwardDCT64
	}
}

// ImageToGray function converts an Image to a gray scale array.
func ImageToGray(img image.Image, pixels *[]float32) {
	bounds := img.Bounds()
	if bounds.Max.X-bounds.Min.X != bounds.Max.Y-bounds.Min.Y {
		return
	}
	switch c := img.(type) {
	case *image.YCbCr:
		if FlagUseASM {
			AsmYCbCrToGray(*pixels,
				c.Rect.Min.X, c.Rect.Min.Y, c.Rect.Max.X, c.Rect.Max.Y,
				c.Y, c.Cb, c.Cr, c.YStride, c.CStride)
		} else {
			YCbCrToGrayAlt(c, *pixels)
		}
	case *image.RGBA:
		rgbaToGray(c, *pixels)
	default:
		imageToGrayDefault(c, *pixels)
	}
}

// imageToGrayDefault uses the image.Image interface to produce grayscale pixels
func imageToGrayDefault(img image.Image, pixels []float32) {
	b := img.Bounds()
	s := b.Max.X - b.Min.X
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = float32(pixelToGray(img.At(j, i).RGBA()))
		}
	}
}

// pixel2Gray converts a pixel to grayscale value base on luminosity
func pixelToGray(r, g, b, a uint32) float64 {
	return 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
}

// YCbCrToGray uses *image.YCbCr which is signifiantly faster than the image.Image interface.
func YCbCrToGray(colorImg *image.YCbCr, pixels []float64) {
	s := colorImg.Rect.Dx()
	for i := 0; i < s; i++ {
		for j := 0; j < s; j += 4 {
			pixels[(i*s)+j+0] = pixelToGray(colorImg.YCbCrAt(j+0, i).RGBA())
			pixels[(i*s)+j+1] = pixelToGray(colorImg.YCbCrAt(j+1, i).RGBA())
			pixels[(i*s)+j+2] = pixelToGray(colorImg.YCbCrAt(j+2, i).RGBA())
			pixels[(i*s)+j+3] = pixelToGray(colorImg.YCbCrAt(j+3, i).RGBA())
		}
	}
}

// YCbCrToGrayAlt convers an *image.YCbCr to array of pixels.
func YCbCrToGrayAlt(img *image.YCbCr, pixels []float32) {
	s := img.Rect.Max.X - img.Rect.Min.X
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			yi := img.YOffset(x, y)
			ci := img.COffset(x, y)

			yy := img.Y[yi]
			cb := img.Cb[ci]
			cr := img.Cr[ci]

			yy1 := int32(yy) * 0x10101
			cb1 := int32(cb) - 128
			cr1 := int32(cr) - 128

			r := yy1 + 91881*cr1
			//if uint32(r)&0xff000000 == 0 {
			//	r >>= 16
			//} else {
			//	r = ^(r >> 31)
			//}

			g := yy1 - 22554*cb1 - 46802*cr1
			//if uint32(g)&0xff000000 == 0 {
			//	g >>= 16
			//} else {
			//	g = ^(g >> 31)
			//}

			b := yy1 + 116130*cb1
			//if uint32(b)&0xff000000 == 0 {
			//	b >>= 16
			//} else {
			//	b = ^(b >> 31)
			//}

			pixels[(y*s)+x] = float32(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b>>8))
		}
	}
}

// rgbaToGray uses *image.RGBA which is signifiantly faster than the image.Image interface.
func rgbaToGray(img *image.RGBA, pixels []float32) {
	s := img.Rect.Max.X - img.Rect.Min.X
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = float32(pixelToGray(img.At(j, i).RGBA()))
		}
	}
}

// FlattenPixels function flattens 2d array into 1d array.
func FlattenPixels32(pixels [][]float64, x int, y int) []float64 {
	flattens := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[i][j]
		}
	}
	return flattens
}

func LinearTosRGB32(value float64) int {
	v := math.Max(0, math.Min(1, value))
	if v <= 0.0031308 {
		return int(v*12.92*255 + 0.5)
	}
	return int((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5)
}

func SRGBToLinear32(value int) float64 {
	v := float64(value) / 255
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func SignPow32(value, exp float64) float64 {
	return math.Copysign(math.Pow(math.Abs(value), exp), value)
}
