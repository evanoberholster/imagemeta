// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

// Quantize generates a 16-element hash by thresholding each value of a 16x16 input array against the provided median.
// Returns an error if the input array length is not 256.
func Quantize(input []float32, median float32) ([16]uint16, error) {
	if len(input) != 16*16 {
		return [16]uint16{}, ErrQuantizeLength
	}

	var hash [16]uint16

	for i := range 16 {
		var bits uint16
		row := input[i*16 : (i+1)*16]

		for j, v := range row {
			if v > median {
				bits |= 1 << j
			}
		}

		hash[i] = bits
	}

	return hash, nil
}
