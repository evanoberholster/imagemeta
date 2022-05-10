// Copyright 2022 Evan Oberholster
// Copyright 2017 The goimagehash Authors.
// All rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagehash

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"sync"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
)

//go:generate msgp

// Phash is a 64bit Perception Hash
type Phash uint64

// Ahash is a 64bit Average Hash
type Ahash uint64

const (
	NilPhash    Phash = 0
	NilAhash    Ahash = 0
	LengthPHash       = 8
)

var (
	ErrImageObject = errors.New("image object can not be nil")

	encodeFn = binary.LittleEndian.PutUint64
	decodeFn = binary.LittleEndian.Uint64
)

// NewPHash is a Perception Hash function returns a hash computation of phash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func NewPHash(img image.Image) (phash Phash, err error) {
	if img == nil {
		return NilPhash, ErrImageObject
	}

	pixels := transforms.Rgb2Gray(img)
	dct := transforms.DCT2D(pixels, 64, 64)
	flattens := transforms.FlattenPixels(dct, 8, 8)
	median := transforms.MedianOfPixels(flattens)

	for idx, p := range flattens {
		if p > median {
			phash |= 1 << uint(len(flattens)-idx-1) // leftShiftSet
		}
	}
	return phash, nil
}

var pixelsPool = sync.Pool{
	New: func() interface{} { return make([]float64, 4096) },
}

// NewPHash is a Perception Hash function returns a hash computation of phash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func NewPHashFast(img image.Image) (phash Phash, err error) {
	if img == nil {
		return NilPhash, ErrImageObject
	}

	pixels := pixelsPool.Get().([]float64)
	defer pixelsPool.Put(pixels)

	transforms.Rgb2Gray_new(img, pixels)
	transforms.DCT2D_new(pixels)

	median := transforms.MedianOfPixels_new(pixels)

	for idx, p := range pixels {
		if p > median {
			phash |= 1 << uint(len(pixels)-idx-1) // leftShiftSet
		}
	}
	return phash, nil
}

// NewAHash is an Average Hash fuction that returns a hash computation of average hash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func NewAHash(img image.Image) (ahash Ahash, err error) {
	if img == nil {
		return NilAhash, ErrImageObject
	}

	// Create 64bits hash.
	//resized := resize.Resize(8, 8, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(img)
	flattens := transforms.FlattenPixels(pixels, 8, 8)
	avg := transforms.MeanOfPixels(flattens)

	for idx, p := range flattens {
		if p > avg {
			ahash |= 1 << uint(len(flattens)-idx-1)
		}
	}

	return ahash, nil
}

// Distance between Phash values
func (ph Phash) Distance(hash Phash) uint8 {
	return uint8(popcnt(uint64(ph) ^ uint64(hash)))
}

func (ph Phash) String() string {
	return fmt.Sprintf("p:%016x", uint64(ph))
}

func (ph Phash) Encode(buf []byte) {
	encodeFn(buf[:LengthPHash], uint64(ph))
}

func (ph *Phash) Decode(buf []byte) {
	*ph = Phash(decodeFn(buf[:LengthPHash]))
}
