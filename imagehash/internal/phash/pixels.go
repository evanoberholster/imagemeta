// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package phash contains the transforms used by the PHash perceptual hash.
package phash

import (
	"image"
)

var (
	// FlagUseASM determines if ASM can be used on the current CPU
	FlagUseASM    = false
	ForwardDCT64  = forwardDCT64
	ForwardDCT256 = forwardDCT256
	YCbCrToGray   = yCbCrToGray
	RGBAtoGray    = rgbaToGray
)

// ImageToGray converts an image to a gray scale array.
func ImageToGray(img image.Image, pixels *[]float32) {
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		return
	}
	switch c := img.(type) {
	case *image.YCbCr:
		YCbCrToGray(c, *pixels)
	case *image.RGBA:
		RGBAtoGray(c, *pixels)
	default:
		imageToGrayDefault(c, *pixels)
	}
}

// imageToGrayDefault uses the image.Image interface to produce grayscale pixels
func imageToGrayDefault(img image.Image, pixels []float32) {
	b := img.Bounds()
	s := b.Dx()
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixelToGray(img.At(b.Min.X+j, b.Min.Y+i).RGBA())
		}
	}
}

// pixelToGray converts a 16-bit RGBA pixel to a grayscale value based on luminosity
func pixelToGray(r, g, b, a uint32) float32 {
	return float32(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256))
}

// yCbCrToGray converts an *image.YCbCr to an array of grayscale pixels by
// indexing the Y/Cb/Cr planes directly.
func yCbCrToGray(img *image.YCbCr, pixels []float32) {
	s := img.Rect.Dx()
	minX, minY := img.Rect.Min.X, img.Rect.Min.Y
	for y := 0; y < s; y++ {
		yRow := img.YOffset(minX, minY+y)
		base := y * s
		for x := 0; x < s; x++ {
			ci := img.COffset(minX+x, minY+y)

			yy := img.Y[yRow+x]
			cb := img.Cb[ci]
			cr := img.Cr[ci]

			yy1 := int32(yy) * 0x10101
			cb1 := int32(cb) - 128
			cr1 := int32(cr) - 128

			r := yy1 + 91881*cr1
			g := yy1 - 22554*cb1 - 46802*cr1
			b := yy1 + 116130*cb1

			pixels[base+x] = float32(0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b>>8))
		}
	}
}

// rgbaToGray converts an *image.RGBA to an array of grayscale pixels by
// indexing the pixel buffer directly.
func rgbaToGray(img *image.RGBA, pixels []float32) {
	s := img.Rect.Dx()
	minX, minY := img.Rect.Min.X, img.Rect.Min.Y
	for i := 0; i < s; i++ {
		row := img.PixOffset(minX, minY+i)
		base := i * s
		for j := 0; j < s; j++ {
			p := row + j*4
			pixels[base+j] = pixelToGray8(img.Pix[p], img.Pix[p+1], img.Pix[p+2])
		}
	}
}

// pixelToGray8 converts an 8-bit RGB pixel to a grayscale value. It is
// equivalent to pixelToGray applied to the corresponding 16-bit values.
func pixelToGray8(r, g, b uint8) float32 {
	return float32(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}
