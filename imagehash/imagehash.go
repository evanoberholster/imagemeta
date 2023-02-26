// Package imagehash processes a Perception hash and Average hash from an image.
package imagehash

// Copyright 2022 Evan Oberholster
// Copyright 2017 The goimagehash Authors.
// All rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"sync"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
)

//go:generate msgp

// NewPHash64 is a Perception Hash function returns a hash computation of phash.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash64(img image.Image) (phash PHash64, err error) {
	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 64 {
		err = errors.New("error image size incompatible. PHash requires 64x64 image")
		return
	}

	pixels := pixelsPool64.Get().(*[]float64)

	transforms.Rgb2GrayFast(img, pixels)
	flattens := transforms.DCT2DHash64(pixels)
	//flattens := transforms.FlattenPixelsHash64(pixels)
	pixelsPool64.Put(pixels)

	median := transforms.MedianOfPixels64(flattens[:])

	for idx, p := range flattens {
		if p > median {
			phash |= 1 << uint(len(flattens)-idx-1) // leftShiftSet
		}
	}
	return phash, nil
}

// NewPHash256 is a Perception Hash function returns a 256bit hash computation
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash256(img image.Image) (phash PHash256, err error) {
	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 256 {
		err = errors.New("error image size incompatible. PHash256 requires 256x256 image")
		return
	}

	pixels := pixelsPool256.Get().(*[]float64)

	transforms.Rgb2GrayFast(img, pixels)
	flattens := transforms.DCT2DHash256(pixels)
	//flattens := transforms.FlattenPixelsHash256(pixels)
	pixelsPool256.Put(pixels)

	median := transforms.MedianOfPixels256(flattens[:])

	for idx, p := range flattens {
		indexOfArray := idx / 64
		if p > median {
			phash[indexOfArray] |= 1 << uint(64-idx%64-1) // leftShiftSet
		}
	}

	return phash, nil
}

// NewAHash is an Average Hash fuction that returns a hash computation of average hash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func NewAHash(img image.Image) (ahash Ahash, err error) {
	if img == nil {
		err = ErrImageObject
		return
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

// Pixel Pools

// Pixel pool 64bit
var pixelsPool64 = sync.Pool{
	New: func() interface{} {
		p := make([]float64, 4096)
		return &p
	},
}

// Pixel pool 256bit
var pixelsPool256 = sync.Pool{
	New: func() interface{} {
		p := make([]float64, 65536)
		return &p
	},
}

// Variables
var (
	ErrImageObject = errors.New("image object can not be nil")

	encodeFn = binary.LittleEndian.PutUint64
	decodeFn = binary.LittleEndian.Uint64
)

// Ahash is a 64bit Average Hash
type Ahash uint64

// Phash is a type alias for PHash64
type Phash = PHash64

// PHash64 is a 64bit Perception Hash
type PHash64 uint64

// Distance between Phash values
func (ph PHash64) Distance(hash PHash64) uint8 {
	return uint8(popcnt(uint64(ph) ^ uint64(hash)))
}

func (ph PHash64) String() string {
	return fmt.Sprintf("p:%016x", uint64(ph))
}

func (ph PHash64) Encode(dst []byte) {
	encodeFn(dst[:8], uint64(ph))
}

func (ph *PHash64) Decode(src []byte) {
	*ph = PHash64(decodeFn(src[:8]))
}

// PHash256 is a 256bit Perception Hash
type PHash256 [4]uint64

// Distance between Phash values
func (ph PHash256) Distance(hash PHash256) uint {
	var i uint
	i += uint(popcnt(ph[0] ^ hash[0]))
	i += uint(popcnt(ph[1] ^ hash[1]))
	i += uint(popcnt(ph[2] ^ hash[2]))
	i += uint(popcnt(ph[3] ^ hash[3]))
	return i
}

func (ph PHash256) String() string {
	return fmt.Sprintf("p:%016x%016x%016x%016x", uint64(ph[0]), uint64(ph[1]), uint64(ph[2]), uint64(ph[3]))
}

func (ph PHash256) Encode(buf []byte) {
	encodeFn(buf[:8], ph[0])
	encodeFn(buf[8*1:], ph[1])
	encodeFn(buf[8*2:], ph[2])
	encodeFn(buf[8*3:], ph[3])
}

func (ph *PHash256) Decode(buf []byte) {
	ph[0] = decodeFn(buf[:8])
	ph[1] = decodeFn(buf[8*1:])
	ph[2] = decodeFn(buf[8*2:])
	ph[3] = decodeFn(buf[8*3:])
}
