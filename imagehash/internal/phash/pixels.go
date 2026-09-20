// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package phash contains the transforms used by the PHash perceptual hash.
package phash

import (
	"image"
	"image/color"
)

var (
	// FlagUseASM determines if ASM can be used on the current CPU
	FlagUseASM    = false
	ForwardDCT64  = forwardDCT64
	ForwardDCT256 = forwardDCT256
	YCbCrToGray   = yCbCrToGray
	RGBAtoGray    = rgbaToGray
	NRGBAtoGray   = nrgbaToGray
	GrayToGray    = grayToGray
)

// ImageToGray converts an image to a gray scale array. pixels must have room
// for Dx*Dx floats.
func ImageToGray(img image.Image, pixels []float32) {
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		return
	}
	switch c := img.(type) {
	case *image.YCbCr:
		YCbCrToGray(c, pixels)
	case *image.RGBA:
		RGBAtoGray(c, pixels)
	case *image.NRGBA:
		NRGBAtoGray(c, pixels)
	case *image.Gray:
		GrayToGray(c, pixels)
	default:
		imageToGrayDefault(c, pixels)
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

// pixelToGray converts a 16-bit RGBA pixel to a grayscale value using the same
// integer luminosity weights as every other pathway.
func pixelToGray(r, g, b, a uint32) float32 {
	r8, g8, b8 := r/257, g/257, b/257
	return float32(299*r8+587*g8+114*b8) / 1000
}

// yCbCrToGray converts an *image.YCbCr to an array of grayscale pixels by
// indexing the Y/Cb/Cr planes directly.
//
// The 16-bit YCbCr-to-RGB conversion is evaluated with 32-bit integer
// arithmetic and then reduced to grayscale with the integer luminosity weights
// (299*r + 587*g + 114*b) / 1000. Integer-only arithmetic keeps every
// architecture's result bit-identical (no float multiply-add fusion).
func yCbCrToGray(img *image.YCbCr, pixels []float32) {
	s := img.Rect.Dx()
	minX, minY := img.Rect.Min.X, img.Rect.Min.Y
	for y := 0; y < s; y++ {
		yRow := img.YOffset(minX, minY+y)
		base := y * s
		for x := 0; x < s; x++ {
			ci := img.COffset(minX+x, minY+y)

			yy := int32(img.Y[yRow+x]) * 0x10101
			cb := int32(img.Cb[ci]) - 128
			cr := int32(img.Cr[ci]) - 128

			red := (yy + 91881*cr) >> 8
			green := (yy - 22554*cb - 46802*cr) >> 8
			blue := (yy + 116130*cb) >> 8

			pixels[base+x] = float32(299*red+587*green+114*blue) / 1000
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

// pixelToGray8 converts an 8-bit RGB pixel to a grayscale value using the same
// integer luminosity weights as every other pathway.
func pixelToGray8(r, g, b uint8) float32 {
	return float32(299*int32(r)+587*int32(g)+114*int32(b)) / 1000
}

// nrgbaToGray converts an *image.NRGBA to grayscale by indexing the pixel
// buffer directly. Opaque pixels use the integer luminosity weights
// directly; translucent pixels replicate At().RGBA() premultiplication
// exactly, matching the previous generic-path behavior.
func nrgbaToGray(img *image.NRGBA, pixels []float32) {
	s := img.Rect.Dx()
	minX, minY := img.Rect.Min.X, img.Rect.Min.Y
	for i := 0; i < s; i++ {
		row := img.PixOffset(minX, minY+i)
		base := i * s
		for j := 0; j < s; j++ {
			p := row + j*4
			if alpha := img.Pix[p+3]; alpha == 0xff {
				pixels[base+j] = pixelToGray8(img.Pix[p], img.Pix[p+1], img.Pix[p+2])
			} else {
				r, g, b, a := color.NRGBA{R: img.Pix[p], G: img.Pix[p+1], B: img.Pix[p+2], A: alpha}.RGBA()
				pixels[base+j] = pixelToGray(r, g, b, a)
			}
		}
	}
}

// grayToGray converts an *image.Gray to grayscale by copying the luma plane.
// A gray value v equals luminosity( v,v,v ), so no weighting is needed.
func grayToGray(img *image.Gray, pixels []float32) {
	s := img.Rect.Dx()
	minX, minY := img.Rect.Min.X, img.Rect.Min.Y
	for i := 0; i < s; i++ {
		row := img.PixOffset(minX, minY+i)
		base := i * s
		for j := 0; j < s; j++ {
			pixels[base+j] = float32(img.Pix[row+j])
		}
	}
}
