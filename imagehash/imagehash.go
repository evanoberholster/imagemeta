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

	"github.com/evanoberholster/imagemeta/imagehash/internal/phash"
)

//go:generate msgp

const (
	phash64Side  = 64
	phash256Side = 256
	ahashSide    = 8
)

// Errors returned by the hash constructors.
var (
	ErrImageObject = errors.New("image object can not be nil")
	ErrImageSize   = errors.New("image size incompatible with hash size")
	ErrPixelPool   = errors.New("pixel pool returned unexpected type")
)

// NewPHash64 is a Perception Hash function. It returns a 64 bit hash of the
// image and requires a 64x64 image.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash64(img image.Image) (hash PHash64, err error) {
	if err = checkImageSize(img, phash64Side); err != nil {
		return 0, err
	}

	pixels, ok := pixelsPool64.Get().(*[]float32)
	if !ok || pixels == nil {
		return 0, ErrPixelPool
	}
	defer pixelsPool64.Put(pixels)

	phash.ImageToGray(img, pixels)
	flattens := phash.DCT2DHash64(*pixels)
	median := phash.MedianOfPixels64(flattens[:])

	for idx, p := range flattens {
		if p > median {
			hash |= 1 << (63 - idx) // leftShiftSet
		}
	}
	return hash, nil
}

// NewPHash256 is a Perception Hash function. It returns a 256 bit hash of the
// image and requires a 256x256 image.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash256(img image.Image) (hash PHash256, err error) {
	if err = checkImageSize(img, phash256Side); err != nil {
		return PHash256{}, err
	}

	pixels, ok := pixelsPool256.Get().(*[]float32)
	if !ok || pixels == nil {
		return PHash256{}, ErrPixelPool
	}
	defer pixelsPool256.Put(pixels)

	phash.ImageToGray(img, pixels)
	var flattens [256]float32
	phash.DCT2DHash256(pixels, &flattens)
	median := phash.MedianOfPixels256(flattens[:])

	for idx, p := range flattens {
		if p > median {
			hash[idx/64] |= 1 << (63 - idx%64) // leftShiftSet
		}
	}
	return hash, nil
}

// NewAHash is an Average Hash function that returns a 64 bit hash of the image
// and requires an 8x8 image.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func NewAHash(img image.Image) (ahash Ahash, err error) {
	if err = checkImageSize(img, ahashSide); err != nil {
		return 0, err
	}

	pixels, ok := pixelsPool64.Get().(*[]float32)
	if !ok || pixels == nil {
		return 0, ErrPixelPool
	}
	defer pixelsPool64.Put(pixels)

	phash.ImageToGray(img, pixels)
	flattens := (*pixels)[:ahashSide*ahashSide]
	avg := meanOfPixels(flattens)

	for idx, p := range flattens {
		if p > avg {
			ahash |= 1 << (63 - idx) // leftShiftSet
		}
	}
	return ahash, nil
}

// NewPHash64Alt is retained for backwards compatibility and is equivalent to
// NewPHash64.
//
// Deprecated: use NewPHash64.
func NewPHash64Alt(img image.Image) (PHash64, error) { return NewPHash64(img) }

// NewPHash256Alt is retained for backwards compatibility and is equivalent to
// NewPHash256.
//
// Deprecated: use NewPHash256.
func NewPHash256Alt(img image.Image) (PHash256, error) { return NewPHash256(img) }

// checkImageSize returns ErrImageObject for a nil image and wraps ErrImageSize
// when the image is not exactly side x side.
func checkImageSize(img image.Image, side int) error {
	if img == nil {
		return ErrImageObject
	}
	size := img.Bounds().Size()
	if size.X != side || size.Y != side {
		return fmt.Errorf("%w: got %dx%d, want %dx%d", ErrImageSize, size.X, size.Y, side, side)
	}
	return nil
}

func meanOfPixels(pixels []float32) float32 {
	var sum float32
	for _, p := range pixels {
		sum += p
	}
	return sum / float32(len(pixels))
}

// Pixel pools
var (
	pixelsPool64 = sync.Pool{
		New: func() interface{} {
			p := make([]float32, phash64Side*phash64Side)
			return &p
		},
	}
	pixelsPool256 = sync.Pool{
		New: func() interface{} {
			p := make([]float32, phash256Side*phash256Side)
			return &p
		},
	}
)

// Variables
var (
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
	return uint8(popcnt(uint64(ph) ^ uint64(hash))) //nolint:gosec // popcnt is bounded to [0,64].
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
	return uint(
		popcnt(ph[0]^hash[0]) +
			popcnt(ph[1]^hash[1]) +
			popcnt(ph[2]^hash[2]) +
			popcnt(ph[3]^hash[3]))
}

func (ph PHash256) String() string {
	return fmt.Sprintf("p:%016x%016x%016x%016x", ph[0], ph[1], ph[2], ph[3])
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
