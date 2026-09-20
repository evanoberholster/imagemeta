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
	r := min(max((yy+91881*cr1)>>16, 0), 255)
	g := min(max((yy-22554*cb1-46802*cr1)>>16, 0), 255)
	b := min(max((yy+116130*cb1)>>16, 0), 255)
	return 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
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

// useASMLuma reports whether a SIMD luminance kernel is available. It is
// enabled by architecture-specific files (arm64 always, amd64 when the CPU
// supports the required features).
var useASMLuma = false

// pdqLuma444Row converts kw 1:1-mapped pixels in a row; pdqLuma420Row
// converts kw Y pixels sharing kw/2 chroma pixels. kw is always a positive
// multiple of 8. Architecture-specific files redirect these to SIMD kernels;
// otherwise they run the portable implementation.
var (
	pdqLuma444Row = scalarLuma444Row
	pdqLuma420Row = scalarLuma420Row
)

func scalarLuma444Row(out []float32, y, cb, cr []uint8) {
	for i := range out {
		out[i] = ycbcrToLuma(y[i], cb[i], cr[i])
	}
}

func scalarLuma420Row(out []float32, y, cb, cr []uint8) {
	for i := range out {
		out[i] = ycbcrToLuma(y[i], cb[i/2], cr[i/2])
	}
}

// lumaASMWidth returns the leading pixel count of a row eligible for the SIMD
// kernels: a multiple of 8 with an even origin (so 420 chroma pairs align).
// It returns 0 when no kernel is available.
func lumaASMWidth(minX, w int) int {
	if !useASMLuma || minX%2 != 0 || w < 8 {
		return 0
	}
	return w / 8 * 8
}

// lumaYCbCr writes the luminance of a YCbCr image. The chroma index math
// replicates image.YCbCr.COffset exactly (including odd origins).
func lumaYCbCr(c *image.YCbCr, out []float32, w, h int) {
	minX, minY := c.Rect.Min.X, c.Rect.Min.Y
	kw := lumaASMWidth(minX, w)
	switch c.SubsampleRatio {
	case image.YCbCrSubsampleRatio444:
		for y := 0; y < h; y++ {
			yOff := c.YOffset(minX, minY+y)
			outLine := out[y*w : (y+1)*w]
			if kw > 0 {
				pdqLuma444Row(outLine[:kw], c.Y[yOff:yOff+kw],
					c.Cb[y*c.CStride:y*c.CStride+kw],
					c.Cr[y*c.CStride:y*c.CStride+kw])
			}
			yLine := c.Y[yOff:]
			cLine := c.Cb[(y * c.CStride):]
			cLineR := c.Cr[(y * c.CStride):]
			for x := kw; x < w; x++ {
				outLine[x] = ycbcrToLuma(yLine[x], cLine[x], cLineR[x])
			}
		}
	case image.YCbCrSubsampleRatio422:
		cx0 := minX / 2
		for y := 0; y < h; y++ {
			yOff := c.YOffset(minX, minY+y)
			outLine := out[y*w : (y+1)*w]
			cRow := y * c.CStride
			if kw > 0 {
				pdqLuma420Row(outLine[:kw], c.Y[yOff:yOff+kw],
					c.Cb[cRow:cRow+kw/2], c.Cr[cRow:cRow+kw/2])
			}
			yLine := c.Y[yOff:]
			for x := kw; x < w; x++ {
				ci := cRow + (minX+x)/2 - cx0
				outLine[x] = ycbcrToLuma(yLine[x], c.Cb[ci], c.Cr[ci])
			}
		}
	case image.YCbCrSubsampleRatio440:
		cy0 := minY / 2
		for y := 0; y < h; y++ {
			yOff := c.YOffset(minX, minY+y)
			outLine := out[y*w : (y+1)*w]
			cRow := ((minY+y)/2 - cy0) * c.CStride
			if kw > 0 {
				pdqLuma444Row(outLine[:kw], c.Y[yOff:yOff+kw],
					c.Cb[cRow:cRow+kw], c.Cr[cRow:cRow+kw])
			}
			yLine := c.Y[yOff:]
			for x := kw; x < w; x++ {
				ci := cRow + x
				outLine[x] = ycbcrToLuma(yLine[x], c.Cb[ci], c.Cr[ci])
			}
		}
	case image.YCbCrSubsampleRatio420:
		cx0, cy0 := minX/2, minY/2
		for y := 0; y < h; y++ {
			yOff := c.YOffset(minX, minY+y)
			outLine := out[y*w : (y+1)*w]
			cRow := ((minY+y)/2-cy0)*c.CStride - cx0
			if kw > 0 {
				pdqLuma420Row(outLine[:kw], c.Y[yOff:yOff+kw],
					c.Cb[cRow:cRow+kw/2], c.Cr[cRow:cRow+kw/2])
			}
			yLine := c.Y[yOff:]
			for x := kw; x < w; x++ {
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
