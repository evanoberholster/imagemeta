package imagehash

import "image"

// factorsYCbCR uses *image.YCbCr to produce factors
func factorsYCbCR(img *image.YCbCr, factors []float64) {
	var factor float64
	var scale float64
	var lr, lg, lb float64

	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X
	size := float64(width * height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rt, gt, bt, _ := img.YCbCrAt(x, y).RGBA()
			lr = channelToLinear[rt>>8]
			lg = channelToLinear[gt>>8]
			lb = channelToLinear[bt>>8]

			for yc := 0; yc < yComponents; yc++ {
				for xc := 0; xc < xComponents; xc++ {

					if xc != 0 || yc != 0 {
						scale = 2 / size
					} else {
						scale = 1 / size
					}
					factor = xvalues[x+width*xc] * yvalues[y+height*yc] * scale
					factors[0+xc*3+yc*3*xComponents] += lr * factor
					factors[1+xc*3+yc*3*xComponents] += lg * factor
					factors[2+xc*3+yc*3*xComponents] += lb * factor
				}
			}
		}
	}
}

// factorsRGBA uses *image.RGBA to produce factors
func factorsRGBA(img *image.RGBA, factors []float64) {
	var factor float64
	var scale float64
	var lr, lg, lb float64

	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X
	size := float64(width * height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rt, gt, bt, _ := img.RGBAAt(x, y).RGBA()
			lr = channelToLinear[rt>>8]
			lg = channelToLinear[gt>>8]
			lb = channelToLinear[bt>>8]

			for yc := 0; yc < yComponents; yc++ {
				for xc := 0; xc < xComponents; xc++ {

					if xc != 0 || yc != 0 {
						scale = 2 / size
					} else {
						scale = 1 / size
					}
					factor = xvalues[x+width*xc] * yvalues[y+height*yc] * scale
					factors[0+xc*3+yc*3*xComponents] += lr * factor
					factors[1+xc*3+yc*3*xComponents] += lg * factor
					factors[2+xc*3+yc*3*xComponents] += lb * factor
				}
			}
		}
	}
}

// factorsDefault uses image.Image to produce factors
func factorsDefault(img image.Image, factors []float64) {
	var factor float64
	var scale float64
	var lr, lg, lb float64

	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X
	size := float64(width * height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rt, gt, bt, _ := img.At(x, y).RGBA()
			lr = channelToLinear[rt>>8]
			lg = channelToLinear[gt>>8]
			lb = channelToLinear[bt>>8]

			for yc := 0; yc < yComponents; yc++ {
				for xc := 0; xc < xComponents; xc++ {

					if xc != 0 || yc != 0 {
						scale = 2 / size
					} else {
						scale = 1 / size
					}
					factor = xvalues[x+width*xc] * yvalues[y+height*yc] * scale
					factors[0+xc*3+yc*3*xComponents] += lr * factor
					factors[1+xc*3+yc*3*xComponents] += lg * factor
					factors[2+xc*3+yc*3*xComponents] += lb * factor
				}
			}
		}
	}
}
