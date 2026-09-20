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

// ycbcrToLuma converts one YCbCr pixel to BT.601 luminance using the same
// integer JFIF ratios as image/color.YCbCrToRGB (inlined here to avoid the
// call and the uint8 round-trip) followed by the same float weights, so the
// result is bit-identical to the RGB round-trip. Clamping is preserved
// exactly: out-of-gamut channels saturate to 0/255 before weighting,
// matching Meta's reference behavior.
func ycbcrToLuma(y, cb, cr uint8) float32 {
	yy := int32(y) * 0x10101
	cb1 := int32(cb) - 128
	cr1 := int32(cr) - 128

	// Arithmetic shift then saturate: identical mapping to the branchless
	// select in image/color.YCbCrToRGB (in-range values pass through,
	// negatives saturate to 0, large positives to 255). Min/max form keeps
	// the data flow branch-free and mirrors future SIMD lanes exactly.
	r := yy + 91881*cr1
	if uint32(r)&0xff000000 == 0 {
		r >>= 16
	} else {
		r = ^(r >> 31)
	}
	g := yy - 22554*cb1 - 46802*cr1
	if uint32(g)&0xff000000 == 0 {
		g >>= 16
	} else {
		g = ^(g >> 31)
	}
	b := yy + 116130*cb1
	if uint32(b)&0xff000000 == 0 {
		b >>= 16
	} else {
		b = ^(b >> 31)
	}
	// r, g, b hold either a clamped [0,255] value or the -1 sentinel for
	// 255 (matching the uint8 conversion in YCbCrToRGB); masking resolves
	// the sentinel exactly.
	return 0.299*float32(r&0xff) + 0.587*float32(g&0xff) + 0.114*float32(b&0xff)
}

// toLuminanceDirect writes the luminance of src (w by h pixels) directly into
// out. Output is bit-identical to expanding src to RGBA and then calling
// toLuminanceInto: the YCbCr conversion inlines color.YCbCrToRGB's exact
// integer arithmetic (including clamping), and NRGBA/Gray fast paths
// replicate the generic At().RGBA() math.
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
	case *image.NRGBA:
		minX, minY := c.Rect.Min.X, c.Rect.Min.Y
		for y := 0; y < h; y++ {
			row := c.PixOffset(minX, minY+y)
			outRow := y * w
			for x := 0; x < w; x++ {
				p := row + x*4
				if a := c.Pix[p+3]; a == 0xff {
					// Opaque: identical to the generic path (v*257>>8 == v).
					out[outRow+x] = 0.299*float32(c.Pix[p]) +
						0.587*float32(c.Pix[p+1]) +
						0.114*float32(c.Pix[p+2])
				} else {
					// Translucent: replicate At().RGBA() premultiplication
					// exactly (rare; keeps the slow path for these pixels).
					r, g, b, _ := color.NRGBA{c.Pix[p], c.Pix[p+1], c.Pix[p+2], a}.RGBA()
					out[outRow+x] = 0.299*float32(r>>8) + 0.587*float32(g>>8) + 0.114*float32(b>>8)
				}
			}
		}
	case *image.Gray:
		minX, minY := c.Rect.Min.X, c.Rect.Min.Y
		for y := 0; y < h; y++ {
			row := c.PixOffset(minX, minY+y)
			outRow := y * w
			for x := 0; x < w; x++ {
				out[outRow+x] = float32(c.Pix[row+x])
			}
		}
	case *image.YCbCr:
		lumaYCbCr(c, out, w, h)
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

// lumaYCbCr writes the luminance of a YCbCr image. The chroma index math
// replicates image.YCbCr.COffset exactly (including odd origins).
func lumaYCbCr(c *image.YCbCr, out []float32, w, h int) {
	minX, minY := c.Rect.Min.X, c.Rect.Min.Y
	switch c.SubsampleRatio {
	case image.YCbCrSubsampleRatio444:
		for y := 0; y < h; y++ {
			yLine := c.Y[c.YOffset(minX, minY+y):]
			cLine := c.Cb[(y * c.CStride):]
			cLineR := c.Cr[(y * c.CStride):]
			outLine := out[y*w : (y+1)*w]
			for x := 0; x < w; x++ {
				outLine[x] = ycbcrToLuma(yLine[x], cLine[x], cLineR[x])
			}
		}
	case image.YCbCrSubsampleRatio422:
		cx0 := minX / 2
		for y := 0; y < h; y++ {
			yLine := c.Y[c.YOffset(minX, minY+y):]
			outLine := out[y*w : (y+1)*w]
			cRow := y * c.CStride
			for x := 0; x < w; x++ {
				ci := cRow + (minX+x)/2 - cx0
				outLine[x] = ycbcrToLuma(yLine[x], c.Cb[ci], c.Cr[ci])
			}
		}
	case image.YCbCrSubsampleRatio440:
		cy0 := minY / 2
		for y := 0; y < h; y++ {
			yLine := c.Y[c.YOffset(minX, minY+y):]
			outLine := out[y*w : (y+1)*w]
			cRow := ((minY+y)/2 - cy0) * c.CStride
			for x := 0; x < w; x++ {
				ci := cRow + x
				outLine[x] = ycbcrToLuma(yLine[x], c.Cb[ci], c.Cr[ci])
			}
		}
	case image.YCbCrSubsampleRatio420:
		cx0, cy0 := minX/2, minY/2
		for y := 0; y < h; y++ {
			yLine := c.Y[c.YOffset(minX, minY+y):]
			outLine := out[y*w : (y+1)*w]
			cRow := ((minY+y)/2-cy0)*c.CStride - cx0
			for x := 0; x < w; x++ {
				ci := cRow + (minX+x)/2
				outLine[x] = ycbcrToLuma(yLine[x], c.Cb[ci], c.Cr[ci])
			}
		}
	default:
		// Rare ratios: fused math with per-pixel COffset.
		for y := 0; y < h; y++ {
			yRow := c.YOffset(minX, minY+y)
			outRow := y * w
			for x := 0; x < w; x++ {
				ci := c.COffset(minX+x, minY+y)
				out[outRow+x] = ycbcrToLuma(c.Y[yRow+x], c.Cb[ci], c.Cr[ci])
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
