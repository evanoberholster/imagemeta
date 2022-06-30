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
	if w != h && w != pHashSize {
		return
	}
	switch c := colorImg.(type) {
	case *image.YCbCr:
		rgb2GrayYCbCR(c, *pixels, w)
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
			pixels[j+(i*s)] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// rgb2GrayYCbCR uses *image.YCbCr which is signifiantly faster than the image.Image interface.
func rgb2GrayYCbCR(colorImg *image.YCbCr, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[j+(i*s)] = pixel2Gray(colorImg.YCbCrAt(j, i).RGBA())
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

// FlattenPixelsFast function flattens pixels array from DCT2D into [64]float array.
func FlattenPixelsFast(pixels *[]float64) []float64 {
	flattens := [pHashSize]float64{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			flattens[8*i+j] = (*pixels)[64*i+j]
		}
	}
	return flattens[:]
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
