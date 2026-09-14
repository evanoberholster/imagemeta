// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// prepareImage conditionally downscales src to at most ImageSize in either
// dimension, writes its luminance into ws.luma, and returns the processing
// dimensions.
//
// Small images are processed at their native size and converted straight to
// luminance, avoiding an intermediate RGBA buffer. Only images larger than
// ImageSize in either dimension are resized in RGB space before luminance is
// taken.
//
// Downscaling uses draw.ApproxBiLinear rather than the higher-quality
// draw.BiLinear: the Jarosz box filter that follows removes the extra
// high-frequency detail, so Meta's reference vectors are unchanged while the
// resize avoids BiLinear's large per-call temporary buffers.
func prepareImage(src image.Image, ws *workspace) (numRows, numCols int) {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w > ImageSize || h > ImageSize {
		draw.ApproxBiLinear.Scale(ws.rgba, image.Rect(0, 0, ImageSize, ImageSize), src, bounds, draw.Src, nil)
		toLuminanceInto(ws.rgba, ws.luma[:ImageSize*ImageSize], ImageSize, ImageSize)
		return ImageSize, ImageSize
	}

	toLuminanceDirect(src, ws.luma[:w*h], w, h)
	return h, w
}

// toLuminanceDirect writes the luminance of src (w by h pixels) directly into
// out. The value is byte-for-byte equivalent to expanding src to RGBA and then
// calling toLuminanceInto.
func toLuminanceDirect(src image.Image, out []float32, w, h int) {
	switch c := src.(type) {
	case *image.RGBA:
		minX, minY := c.Rect.Min.X, c.Rect.Min.Y
		for y := 0; y < h; y++ {
			row := c.PixOffset(minX, minY+y)
			outRow := y * w
			for x := 0; x < w; x++ {
				p := row + x*4
				out[outRow+x] = 0.299*float32(c.Pix[p]) +
					0.587*float32(c.Pix[p+1]) +
					0.114*float32(c.Pix[p+2])
			}
		}
	case *image.YCbCr:
		minX, minY := c.Rect.Min.X, c.Rect.Min.Y
		for y := 0; y < h; y++ {
			yRow := c.YOffset(minX, minY+y)
			outRow := y * w
			for x := 0; x < w; x++ {
				ci := c.COffset(minX+x, minY+y)
				r, g, b := color.YCbCrToRGB(c.Y[yRow+x], c.Cb[ci], c.Cr[ci])
				out[outRow+x] = 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
			}
		}
	default:
		bounds := src.Bounds()
		for y := 0; y < h; y++ {
			outRow := y * w
			for x := 0; x < w; x++ {
				r, g, b, _ := src.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
				out[outRow+x] = 0.299*float32(r>>8) + 0.587*float32(g>>8) + 0.114*float32(b>>8)
			}
		}
	}
}

// toLuminanceInto writes the luminance of the first numCols pixels of numRows
// rows of src into out.
func toLuminanceInto(src *image.RGBA, out []float32, numRows, numCols int) {
	pixels := src.Pix
	stride := src.Stride

	for y := 0; y < numRows; y++ {
		rowOffset := y * stride
		outOffset := y * numCols
		for x := 0; x < numCols; x++ {
			i := rowOffset + x*4
			out[outOffset+x] = 0.299*float32(pixels[i]) + 0.587*float32(pixels[i+1]) + 0.114*float32(pixels[i+2])
		}
	}
}
