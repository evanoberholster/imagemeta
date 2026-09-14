// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"math/bits"
)

const (
	numRows = 64
	numCols = 64
	dims    = numRows * numCols
	scale   = float32(100.0 / 255.0)
)

func abs(x int) int {
	mask := x >> (bits.UintSize - 1)
	return (x ^ mask) - mask
}

// ComputeQuality calculates a quality score based on gradient differences in a 64x64 matrix of float32 values.
// Returns an error if the input matrix does not match the required dimensions.
func ComputeQuality(matrix []float32) (int, error) {
	if len(matrix) != dims {
		return 0, ErrMatrixLength
	}

	m := matrix[:dims]
	var gradientSum int

	// Vertical differences (63 x 64)
	for i := 0; i < numRows-1; i++ {
		row := m[i*numCols : i*numCols+numCols]
		next := m[(i+1)*numCols : (i+1)*numCols+numCols]
		for j := 0; j < numCols; j++ {
			gradientSum += abs(int((row[j] - next[j]) * scale))
		}
	}

	// Horizontal differences (64 x 63)
	for i := 0; i < numRows; i++ {
		row := m[i*numCols : i*numCols+numCols]
		for j := 0; j < numCols-1; j++ {
			gradientSum += abs(int((row[j] - row[j+1]) * scale))
		}
	}

	// 90 is pulled from Meta implementations in C++
	if quality := gradientSum / 90; quality < 100 {
		return quality, nil
	}

	return 100, nil
}
