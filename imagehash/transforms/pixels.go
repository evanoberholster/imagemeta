// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transforms

import (
	"image"
	"math"
)

// Rgb2GrayFast function converts RGB to a gray scale array.
func Rgb2GrayFast(colorImg image.Image, pixels *[]float64) {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	if w != h {
		return
	}
	switch c := colorImg.(type) {
	case *image.YCbCr:
		PixelYCnCRGray(c, *pixels)
	case *image.RGBA:
		rgb2GrayRGBA(c, *pixels, w)
	default:
		rgb2GrayDefault(c, *pixels, w)
	}
}

// pixel2Gray converts a pixel to grayscale value base on luminosity
func pixel2Gray(r, g, b, a uint32) float64 {
	return 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
}

// rgb2GrayDefault uses the image.Image interface
func rgb2GrayDefault(colorImg image.Image, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// YCbCR2Gray uses *image.YCbCr which is signifiantly faster than the image.Image interface.
func YCbCR2Gray(colorImg *image.YCbCr, pixels []float64) {
	s := colorImg.Rect.Dx()
	for i := 0; i < s; i++ {
		for j := 0; j < s; j += 4 {
			pixels[(i*s)+j+0] = pixel2Gray(colorImg.YCbCrAt(j+0, i).RGBA())
			pixels[(i*s)+j+1] = pixel2Gray(colorImg.YCbCrAt(j+1, i).RGBA())
			pixels[(i*s)+j+2] = pixel2Gray(colorImg.YCbCrAt(j+2, i).RGBA())
			pixels[(i*s)+j+3] = pixel2Gray(colorImg.YCbCrAt(j+3, i).RGBA())
		}
	}
}

func PixelYCnCRGray(img *image.YCbCr, pixels []float64) {
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

			pixels[(y*s)+x] = 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b>>8)
		}
	}
}

// OldYCbCR2Gray uses *image.YCbCr which is signifiantly faster than the image.Image interface.
func OldYCbCR2Gray(colorImg *image.YCbCr, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixel2Gray(colorImg.YCbCrAt(j, i).RGBA())
		}
	}
}

// rgb2GrayYCbCR uses *image.RGBA which is signifiantly faster than the image.Image interface.
func rgb2GrayRGBA(colorImg *image.RGBA, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// Rgb2Gray function converts RGB to a gray scale array.
func Rgb2Gray(colorImg image.Image) [][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			pixels[i][j] = lum
		}
	}

	return pixels
}

// FlattenPixels function flattens 2d array into 1d array.
func FlattenPixels(pixels [][]float64, x int, y int) []float64 {
	flattens := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[i][j]
		}
	}
	return flattens
}

func LinearTosRGB(value float64) int {
	v := math.Max(0, math.Min(1, value))
	if v <= 0.0031308 {
		return int(v*12.92*255 + 0.5)
	}
	return int((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5)
}

func SRGBToLinear(value int) float64 {
	v := float64(value) / 255
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func SignPow(value, exp float64) float64 {
	return math.Copysign(math.Pow(math.Abs(value), exp), value)
}
