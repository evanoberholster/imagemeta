// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"fmt"
	"image"
)

const (
	// DefaultMatchThreshold is the recommended Hamming distance threshold for
	// considering two PDQ hashes a match. Per the PDQ spec README:
	// "Distance Threshold to consider two hashes similar/matching: <=31"
	DefaultMatchThreshold = 31

	// DefaultQualityThreshold is the minimum quality score for a hash to be
	// considered reliable. Per the PDQ spec README:
	// "Quality Threshold where we recommend discarding hashes: <=49"
	DefaultQualityThreshold = 50
)

// Hash computes the PDQ perceptual hash of an image.
// Returns the primary hash, quality score (0–100), and all 8 dihedral hashes.
//
// All intermediate buffers are reused through a pooled workspace, so a Hash
// call performs no steady-state heap allocation.
func Hash(img image.Image) (Result, error) {
	if img == nil {
		return Result{}, ErrNilImage
	}

	ws, ok := workspacePool.Get().(*workspace)
	if !ok || ws == nil {
		ws = newWorkspace()
	}
	defer workspacePool.Put(ws)

	numRows, numCols := prepareImage(img, ws)
	n := numRows * numCols

	if err := jaroszFilterInto(ws.luma[:n], ws.a[:n], ws.b[:n], numRows, numCols, ws.filtered[:]); err != nil {
		return Result{}, fmt.Errorf("pdq: jarosz filter failed: %w", err)
	}

	quality, err := ComputeQuality(ws.filtered[:])
	if err != nil {
		return Result{}, fmt.Errorf("pdq: quality computation failed: %w", err)
	}

	dct64To16Into(ws.filtered[:], ws.dct[:], &ws.dctAt, &ws.dctT)
	dihedrals, err := DihedralHashes(ws.dct[:])
	if err != nil {
		return Result{}, fmt.Errorf("pdq: dihedral hashes failed: %w", err)
	}

	return Result{
		Hash:      packHash(dihedrals[0]),
		Quality:   quality,
		Dihedrals: dihedrals,
	}, nil
}

// packHash packs a [16]uint16 hash into Hash256 bytes.
// The PDQ reference (Python Hash256.__str__, C++ Hash256::format) outputs
// slots from index 15 down to 0, each as a big-endian 16-bit hex word.
// We match this by writing slot 15 into bytes 0-1, ..., slot 0 into bytes 30-31.
// Each uint16 is stored big-endian (high byte first) so that hex.EncodeToString
// produces the canonical PDQ hash string.
func packHash(hash [16]uint16) Hash256 {
	var out Hash256

	for i, v := range hash {
		pos := (15 - i) * 2
		out[pos] = byte(v >> 8) // high byte first (big-endian)
		out[pos+1] = byte(v)    //nolint:gosec // G115: deliberate low-byte extraction.
	}

	return out
}
