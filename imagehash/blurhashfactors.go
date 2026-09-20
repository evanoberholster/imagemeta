package imagehash

import (
	"image"
	"image/color"
	"math"
)

// multiplyBasisFunction computes the BlurHash DCT factors for img. The basis
// functions are separable, so the sum over the image is evaluated as two
// passes: first over x for every row, then over y for every (xc, yc) component.
func multiplyBasisFunction(img image.Image, factors []float64) {
	size := float64(width * height)
	var rowSum [height * 3 * xComponents]float32
	var lr, lg, lb [width]float32
	var out [3 * xComponents]float32

	for y := 0; y < height; y++ {
		extractRow(img, y, lr[:], lg[:], lb[:])
		blurRow(out[:], lr[:], lg[:], lb[:])
		copy(rowSum[y*3*xComponents:], out[:])
	}

	finishBasis(rowSum[:], size, factors)
}

// noFMA32 rounds v to float32 precision, which stops the compiler contracting
// a multiply-add into an FMA. This keeps the result bit-identical to the NEON
// kernel and across architectures.
func noFMA32(v float32) float32 {
	return math.Float32frombits(math.Float32bits(v))
}

// noFMA64 is the float64 counterpart of noFMA32.
func noFMA64(v float64) float64 {
	return math.Float64frombits(math.Float64bits(v))
}

// blurRowGo is the portable implementation of blurRow.
func blurRowGo(out, lr, lg, lb []float32) {
	for c, lin := range [3][]float32{lr, lg, lb} {
		for xc := 0; xc < xComponents; xc++ {
			base := xc * width
			var acc float32
			for x := 0; x < width; x++ {
				acc += noFMA32(lin[x] * xvalues32[base+x])
			}
			out[c*xComponents+xc] = acc
		}
	}
}

// extractRow writes the linear-light RGB values of one image row (row index y
// relative to the image bounds) into lr, lg and lb.
func extractRow(img image.Image, y int, lr, lg, lb []float32) {
	b := img.Bounds()
	minX, minY := b.Min.X, b.Min.Y
	switch c := img.(type) {
	case *image.YCbCr:
		yRow := c.YOffset(minX, minY+y)
		for x := 0; x < b.Dx(); x++ {
			ci := c.COffset(minX+x, minY+y)
			r, g, bl := color.YCbCrToRGB(c.Y[yRow+x], c.Cb[ci], c.Cr[ci])
			lr[x] = float32(channelToLinear[r])
			lg[x] = float32(channelToLinear[g])
			lb[x] = float32(channelToLinear[bl])
		}
	case *image.RGBA:
		extractRowPacked(c.Pix, c.PixOffset(minX, minY+y), b.Dx(), lr, lg, lb)
	case *image.NRGBA:
		extractRowNRGBA(c, minX, minY+y, b.Dx(), lr, lg, lb)
	case *image.Gray:
		row := c.PixOffset(minX, minY+y)
		for x := 0; x < b.Dx(); x++ {
			lin := float32(channelToLinear[c.Pix[row+x]])
			lr[x], lg[x], lb[x] = lin, lin, lin
		}
	default:
		for x := 0; x < b.Dx(); x++ {
			rt, gt, bt, _ := img.At(minX+x, minY+y).RGBA()
			lr[x] = float32(channelToLinear[rt>>8])
			lg[x] = float32(channelToLinear[gt>>8])
			lb[x] = float32(channelToLinear[bt>>8])
		}
	}
}

// extractRowPacked writes one row of packed 4-byte RGBA-order pixels as
// linear-light values.
func extractRowPacked(pix []uint8, row, dx int, lr, lg, lb []float32) {
	for x := 0; x < dx; x++ {
		p := row + x*4
		lr[x] = float32(channelToLinear[pix[p]])
		lg[x] = float32(channelToLinear[pix[p+1]])
		lb[x] = float32(channelToLinear[pix[p+2]])
	}
}

// extractRowNRGBA writes one row of an NRGBA image as linear-light values.
// Opaque pixels share the packed layout; translucent pixels fall back to the
// generic path to preserve At().RGBA() premultiplication exactly.
func extractRowNRGBA(c *image.NRGBA, minX, y, dx int, lr, lg, lb []float32) {
	row := c.PixOffset(minX, y)
	for x := 0; x < dx; x++ {
		p := row + x*4
		if c.Pix[p+3] == 0xff {
			lr[x] = float32(channelToLinear[c.Pix[p]])
			lg[x] = float32(channelToLinear[c.Pix[p+1]])
			lb[x] = float32(channelToLinear[c.Pix[p+2]])
			continue
		}
		rt, gt, bt, _ := c.At(minX+x, y).RGBA()
		lr[x] = float32(channelToLinear[rt>>8])
		lg[x] = float32(channelToLinear[gt>>8])
		lb[x] = float32(channelToLinear[bt>>8])
	}
}

// finishBasis reduces the per-row partial sums over y into the DCT factors.
func finishBasis(rowSum []float32, size float64, factors []float64) {
	dcScale, acScale := 1/size, 2/size
	for yc := 0; yc < yComponents; yc++ {
		for xc := 0; xc < xComponents; xc++ {
			scale := acScale
			if xc == 0 && yc == 0 {
				scale = dcScale
			}
			for c := 0; c < 3; c++ {
				var acc float64
				for y := 0; y < height; y++ {
					acc += noFMA64(float64(rowSum[(y*3+c)*xComponents+xc]) * yvalues[y+height*yc])
				}
				factors[c+xc*3+yc*3*xComponents] = acc * scale
			}
		}
	}
}
