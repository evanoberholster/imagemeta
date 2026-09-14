// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

// DihedralHashes generates 8 quantized hash representations by applying
// dihedral transformations to a 16x16 DCT array.
//
// Sign conventions are taken directly from the C++ reference (pdqhashing.cpp):
//
//	orig      rot90     rot180    rot270
//	noxpose   xpose     noxpose   xpose
//	+ + + +   - + - +   + - + -   - - - -
//	+ + + +   - + - +   - + - +   + + + +
//	+ + + +   - + - +   + - + -   - - - -
//	+ + + +   - + - +   - + - +   + + + +
//
//	flipx     flipy     flipplus  flipminus
//	noxpose   noxpose   xpose     xpose
//	- - - -   - + - +   + + + +   + - + -
//	+ + + +   - + - +   + + + +   - + - +
//	- - - -   - + - +   + + + +   + - + -
//	+ + + +   - + - +   + + + +   - + - +
//
// For transposing transforms the C++ writes B[j][i] = ±A[i][j].
func DihedralHashes(dct []float32) ([8][16]uint16, error) {
	if len(dct) != 16*16 {
		return [8][16]uint16{}, ErrDihedralDCTSize
	}

	var out [8][16]uint16
	var buf [16 * 16]float32

	for k := 0; k < 8; k++ {
		for i := 0; i < 16; i++ {
			for j := 0; j < 16; j++ {
				buf[i*16+j] = dihedralValue(k, dct, i, j)
			}
		}

		median, err := TorbenMedian(buf[:])
		if err != nil {
			return [8][16]uint16{}, err
		}
		h, err := Quantize(buf[:], median)
		if err != nil {
			return [8][16]uint16{}, err
		}
		out[k] = h
	}

	return out, nil
}

// dihedralValue returns element (i, j) of the k-th dihedral transform of dct.
// Keeping this as a switch (rather than a table of closures) keeps the hot
// path allocation free.
func dihedralValue(k int, dct []float32, i, j int) float32 {
	switch k {
	case 0: // original
		return dct[i*16+j]
	case 1: // rotate 90: transpose, negate even-row output
		if i&1 != 0 {
			return dct[j*16+i]
		}
		return -dct[j*16+i]
	case 2: // rotate 180: no transpose, negate where (i+j) is odd
		if (i+j)&1 != 0 {
			return -dct[i*16+j]
		}
		return dct[i*16+j]
	case 3: // rotate 270: transpose, negate even-col output
		if j&1 != 0 {
			return dct[j*16+i]
		}
		return -dct[j*16+i]
	case 4: // flip X: no transpose, negate even-row output
		if i&1 != 0 {
			return dct[i*16+j]
		}
		return -dct[i*16+j]
	case 5: // flip Y: no transpose, negate even-col output
		if j&1 != 0 {
			return dct[i*16+j]
		}
		return -dct[i*16+j]
	case 6: // flip +diagonal: transpose only
		return dct[j*16+i]
	default: // 7: flip -diagonal: transpose, negate where (i+j) is odd
		if (i+j)&1 != 0 {
			return -dct[j*16+i]
		}
		return dct[j*16+i]
	}
}
