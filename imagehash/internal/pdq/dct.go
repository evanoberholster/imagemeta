// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import "math"

const (
	dctRows = 16
	dctCols = 64
)

var dctMatrix = func() [dctRows * dctCols]float32 {
	phaseStep := math.Pi / 128.0
	scale64 := 2.0 / float64(dctCols)
	scaleF := float32(math.Sqrt(scale64))

	var m [dctRows * dctCols]float32

	for i := range dctRows {
		fi := float64(i + 1)
		row := m[i*dctCols : (i+1)*dctCols]

		for j := range dctCols {
			row[j] = scaleF * float32(math.Cos(phaseStep*fi*float64(2*j+1)))
		}
	}

	return m
}()

// DCT64To16 performs a 2D Discrete Cosine Transform (DCT) to convert a 64x64 input into a reduced 16x16 output matrix.
func DCT64To16(input []float32) []float32 {
	out := make([]float32, dctRows*dctRows)
	var at [dctCols * dctCols]float32
	var t [dctRows * dctCols]float32
	dct64To16Into(input, out, &at, &t)
	return out
}

// dct64To16Into performs the 2D DCT from a 64x64 input to a 16x16 output,
// writing into out and using at and t as caller-provided scratch buffers.
func dct64To16Into(input, out []float32, at *[dctCols * dctCols]float32, t *[dctRows * dctCols]float32) {
	for i := 0; i < dctCols; i++ {
		for j, v := range input[i*dctCols : (i+1)*dctCols] {
			at[j*dctCols+i] = v
		}
	}

	for i := 0; i < dctRows; i++ {
		dRow := dctMatrix[i*dctCols : (i+1)*dctCols]
		for j := 0; j < dctCols; j++ {
			atRow := at[j*dctCols : (j+1)*dctCols]
			var sum float32
			for k := 0; k < dctCols; k++ {
				sum += dRow[k] * atRow[k]
			}

			t[i*dctCols+j] = sum
		}
	}

	for i := 0; i < dctRows; i++ {
		tRow := t[i*dctCols : (i+1)*dctCols]
		for j := 0; j < dctRows; j++ {
			dRow := dctMatrix[j*dctCols : (j+1)*dctCols]
			var sum float32
			for k := 0; k < dctCols; k++ {
				sum += dRow[k] * tRow[k]
			}

			out[i*dctRows+j] = sum
		}
	}
}
