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

// blurRowGo is the portable implementation of blurRow.
func blurRowGo(out, lr, lg, lb []float32) {
	for c, lin := range [3][]float32{lr, lg, lb} {
		for xc := 0; xc < xComponents; xc++ {
			base := xc * width
			var acc float32
			for x := 0; x < width; x++ {
				// Round the product before accumulating so the compiler cannot
				// contract the multiply-add into an FMA. This keeps the result
				// bit-identical to the NEON kernel and across architectures.
				p := math.Float32frombits(math.Float32bits(lin[x] * xvalues32[base+x]))
				acc += p
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
		row := c.PixOffset(minX, minY+y)
		for x := 0; x < b.Dx(); x++ {
			p := row + x*4
			lr[x] = float32(channelToLinear[c.Pix[p]])
			lg[x] = float32(channelToLinear[c.Pix[p+1]])
			lb[x] = float32(channelToLinear[c.Pix[p+2]])
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

// finishBasis reduces the per-row partial sums over y into the DCT factors.
func finishBasis(rowSum []float32, size float64, factors []float64) {
	for yc := 0; yc < yComponents; yc++ {
		for xc := 0; xc < xComponents; xc++ {
			scale := 2 / size
			if xc == 0 && yc == 0 {
				scale = 1 / size
			}
			for c := 0; c < 3; c++ {
				var acc float64
				for y := 0; y < height; y++ {
					// Same FMA-contraction barrier as blurRowGo.
					p := math.Float64frombits(math.Float64bits(float64(rowSum[(y*3+c)*xComponents+xc]) * yvalues[y+height*yc]))
					acc += p
				}
				factors[c+xc*3+yc*3*xComponents] = acc * scale
			}
		}
	}
}
