// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import "fmt"

const outSize = 64

// JaroszFilter applies a 2-pass filter to src and returns
// a downsampled 64x64 result.
func JaroszFilter(src []float32, numRows, numCols int) ([]float32, error) {
	a := make([]float32, numRows*numCols)
	b := make([]float32, numRows*numCols)
	out := make([]float32, outSize*outSize)
	if err := jaroszFilterInto(src, a, b, numRows, numCols, out); err != nil {
		return nil, err
	}
	return out, nil
}

// jaroszFilterInto applies the four Jarosz box passes to src and writes the
// decimated 64x64 result into dst, using a and b as scratch buffers. All three
// slices are provided by the caller so the hot path can reuse pooled memory.
func jaroszFilterInto(src, a, b []float32, numRows, numCols int, dst []float32) error {
	if numRows <= 0 || numCols <= 0 {
		return fmt.Errorf(
			"pdq: jarosz: dimensions must be positive, got %dx%d",
			numRows, numCols,
		)
	}
	if len(src) != numRows*numCols {
		return fmt.Errorf(
			"pdq: jarosz: src length %d does not match %dx%d=%d",
			len(src), numRows, numCols, numRows*numCols,
		)
	}
	if len(a) < numRows*numCols || len(b) < numRows*numCols || len(dst) < outSize*outSize {
		return fmt.Errorf(
			"pdq: jarosz: scratch buffers too small for %dx%d",
			numRows, numCols,
		)
	}

	windowRows := jaroszWindowSize(numRows)
	windowCols := jaroszWindowSize(numCols)

	boxAlongRows(src, a, numRows, numCols, windowCols)
	boxAlongCols(a, b, numRows, numCols, windowRows)

	boxAlongRows(b, a, numRows, numCols, windowCols)
	boxAlongCols(a, b, numRows, numCols, windowRows)

	decimateInto(b, numRows, numCols, dst)
	return nil
}

// jaroszWindowSize computes the 1D box-filter window size for a single pass.
func jaroszWindowSize(oldDimension int) int {
	return (oldDimension + 2*outSize - 1) / (2 * outSize)
}

// decimateInto samples the center pixel of each output block, matching the
// reference decimateFloat implementation:
//
//	ini = int(((outi + 0.5) * inNumRows) / outNumRows)
func decimateInto(src []float32, numRows, numCols int, dst []float32) {
	for outi := 0; outi < outSize; outi++ {
		ini := int((float64(outi) + 0.5) * float64(numRows) / float64(outSize))
		srcRow := src[ini*numCols:]
		outRow := dst[outi*outSize : (outi+1)*outSize]
		for outj := 0; outj < outSize; outj++ {
			inj := int((float64(outj) + 0.5) * float64(numCols) / float64(outSize))
			outRow[outj] = srcRow[inj]
		}
	}
}

func boxAlongRows(src, dst []float32, numRows, numCols, windowSize int) {
	for y := range numRows {
		box1D(src[y*numCols:], dst[y*numCols:], numCols, 1, windowSize)
	}
}

func boxAlongCols(src, dst []float32, numRows, numCols, windowSize int) {
	if numCols <= maxStreamCols {
		boxColsStream(src, dst, numRows, numCols, windowSize)
		return
	}
	for x := range numCols {
		box1D(src[x:], dst[x:], numRows, numCols, windowSize)
	}
}

// maxStreamCols bounds the stack scratch used by boxColsStream. Larger
// inputs fall back to the strided boxAlongCols loop.
const maxStreamCols = ImageSize

// boxColsStream is boxAlongCols with the loop nest interchanged: the filter
// phase advances outermost while the inner loop streams across columns.
// Every column sees the exact same operation sequence as box1D (bit-identical
// output), but all memory access is row-sequential and the inner loop carries
// no dependency across columns, letting the compiler vectorize it.
func boxColsStream(src, dst []float32, numRows, numCols, windowSize int) {
	halfWindowSize := (windowSize + 2) / 2
	var sumsBuf [maxStreamCols]float32
	sums := sumsBuf[:numCols]

	li, ri, oi, currentWindowSize := 0, 0, 0, 0

	// accumulate without writing
	for range halfWindowSize - 1 {
		row := src[ri*numCols : ri*numCols+numCols]
		for x := 0; x < numCols; x++ {
			sums[x] += row[x]
		}
		currentWindowSize++
		ri++
	}

	// write with growing window
	for range windowSize - halfWindowSize + 1 {
		row := src[ri*numCols : ri*numCols+numCols]
		out := dst[oi*numCols : oi*numCols+numCols]
		currentWindowSize++
		for x := 0; x < numCols; x++ {
			sums[x] += row[x]
			out[x] = sums[x] / float32(currentWindowSize)
		}
		ri++
		oi++
	}

	// write with full window (add right, subtract left). The two updates
	// stay separate statements to match box1D's association exactly:
	// (sum + add) - sub is not the same float as sum + (add - sub).
	for range numRows - windowSize {
		add := src[ri*numCols : ri*numCols+numCols]
		sub := src[li*numCols : li*numCols+numCols]
		out := dst[oi*numCols : oi*numCols+numCols]
		for x := 0; x < numCols; x++ {
			sums[x] += add[x]
			sums[x] -= sub[x]
			out[x] = sums[x] / float32(currentWindowSize)
		}
		li++
		ri++
		oi++
	}

	// write with shrinking window
	for range halfWindowSize - 1 {
		sub := src[li*numCols : li*numCols+numCols]
		out := dst[oi*numCols : oi*numCols+numCols]
		for x := 0; x < numCols; x++ {
			sums[x] -= sub[x]
		}
		currentWindowSize--
		for x := 0; x < numCols; x++ {
			out[x] = sums[x] / float32(currentWindowSize)
		}
		li++
		oi++
	}
}

// box1D implements the 4-phase sliding box filter from the reference.
// stride is 1 for rows, numCols for columns.
func box1D(in, out []float32, length, stride, fullWindowSize int) {
	halfWindowSize := (fullWindowSize + 2) / 2
	li, ri, oi, currentWindowSize := 0, 0, 0, 0

	var sum float32

	// accumulate without writing
	for range halfWindowSize - 1 {
		sum += in[ri]
		currentWindowSize++
		ri += stride
	}

	// write with growing window
	for range fullWindowSize - halfWindowSize + 1 {
		sum += in[ri]
		currentWindowSize++
		out[oi] = sum / float32(currentWindowSize)
		ri += stride
		oi += stride
	}

	// write with full window (add right, subtract left)
	for range length - fullWindowSize {
		sum += in[ri]
		sum -= in[li]
		out[oi] = sum / float32(currentWindowSize)
		li += stride
		ri += stride
		oi += stride
	}

	// write with shrinking window
	for range halfWindowSize - 1 {
		sum -= in[li]
		currentWindowSize--
		out[oi] = sum / float32(currentWindowSize)
		li += stride
		oi += stride
	}
}
