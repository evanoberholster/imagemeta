// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"errors"
	"math/bits"
	"testing"
)

func makeDCT() []float32 {
	dct := make([]float32, 256)
	for i := range dct {
		row, col := i/16, i%16
		dct[i] = float32(row*3-col*2) + float32(row+col+1)*0.1
	}
	return dct
}

func TestDihedralHashes_WrongSize(t *testing.T) {
	_, err := DihedralHashes(make([]float32, 100))
	if !errors.Is(err, ErrDihedralDCTSize) {
		t.Fatalf("expected ErrDihedralDCTSize, got %v", err)
	}
}

func TestDihedralHashes_OriginalMatchesDirect(t *testing.T) {
	dct := makeDCT()
	dihedrals, err := DihedralHashes(dct)
	if err != nil {
		t.Fatal(err)
	}

	median, err := TorbenMedian(dct)
	if err != nil {
		t.Fatal(err)
	}
	direct, err := Quantize(dct, median)
	if err != nil {
		t.Fatal(err)
	}

	if dihedrals[0] != direct {
		t.Errorf("dihedral[0] != direct quantize of input\ngot:  %v\nwant: %v", dihedrals[0], direct)
	}
}

func TestDihedralHashes_EightResults(t *testing.T) {
	dihedrals, err := DihedralHashes(makeDCT())
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, h := range dihedrals {
		if h != ([16]uint16{}) {
			count++
		}
	}
	if count != len(dihedrals) {
		t.Errorf("got %d non-empty dihedrals, want %d", count, len(dihedrals))
	}
}

func TestDihedralHashes_AllDistinct(t *testing.T) {
	dihedrals, err := DihedralHashes(makeDCT())
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[[16]uint16]int)
	for i, h := range dihedrals {
		if prev, ok := seen[h]; ok {
			t.Errorf("dihedral[%d] == dihedral[%d]: %v", i, prev, h)
		}
		seen[h] = i
	}
}

func TestDihedralHashes_FlipPlusIsTranspose(t *testing.T) {
	dct := makeDCT()
	dihedrals, err := DihedralHashes(dct)
	if err != nil {
		t.Fatal(err)
	}

	orig := dihedrals[0]
	flipplus := dihedrals[6]

	for i := range 16 {
		for j := range 16 {
			origBit := (orig[i] >> j) & 1
			transposedBit := (flipplus[j] >> i) & 1
			if origBit != transposedBit {
				t.Errorf("transpose mismatch at (%d,%d): orig=%d flipplus=%d", i, j, origBit, transposedBit)
				return
			}
		}
	}
}

func TestDihedralHashes_HammingDistance(t *testing.T) {
	dct := make([]float32, 256)
	for i := range dct {
		dct[i] = float32(i) * 0.5
	}
	dihedrals, err := DihedralHashes(dct)
	if err != nil {
		t.Fatal(err)
	}
	orig := dihedrals[0]
	for k := 1; k < 8; k++ {
		var dist int
		for i := range 16 {
			dist += bits.OnesCount16(orig[i] ^ dihedrals[k][i])
		}
		t.Logf("dihedral[0] vs dihedral[%d]: hamming=%d", k, dist)
	}
}
