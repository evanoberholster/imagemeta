// Package imagehash computes perceptual hashes and blur placeholders from an
// image: 64/256-bit pHash, 64-bit aHash, 256-bit PDQ and BlurHash strings.
//
// Binary Encode/Decode use little-endian for PHash64, PHash256 and Ahash, and
// big-endian for PDQHash (matching Meta's canonical hex). String encodes
// PHash64/PHash256 with a "p:" prefix and Ahash with an "a:" prefix, while
// PDQHash renders as plain 64-character hex.
package imagehash

// Copyright 2022 Evan Oberholster
// Copyright 2017 The goimagehash Authors.
// All rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"math/bits"
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

// Ahash is a 64bit Average Hash
type Ahash uint64

// PHash64 is a 64bit Perception Hash
type PHash64 uint64

// PHash256 is a 256bit Perception Hash
type PHash256 [4]uint64

// Phash is a type alias for PHash64
type Phash = PHash64

// NewPHash64 is a Perception Hash function. It returns a 64 bit hash of the
// image and requires a 64x64 image.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash64(img image.Image) (hash PHash64, err error) {
	if err = checkImageSize(img, phash64Side); err != nil {
		return 0, err
	}

	pixels, err := borrowPixels(&pixelsPool64)
	if err != nil {
		return 0, err
	}
	defer pixelsPool64.Put(pixels)

	phash.ImageToGray(img, *pixels)
	flattens := phash.DCT2DHash64(*pixels)
	median := phash.MedianOfPixels64(flattens[:])

	return PHash64(bits64(flattens[:], median)), nil
}

// NewPHash256 is a Perception Hash function. It returns a 256 bit hash of the
// image and requires a 256x256 image.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash256(img image.Image) (hash PHash256, err error) {
	if err = checkImageSize(img, phash256Side); err != nil {
		return PHash256{}, err
	}

	pixels, err := borrowPixels(&pixelsPool256)
	if err != nil {
		return PHash256{}, err
	}
	defer pixelsPool256.Put(pixels)

	phash.ImageToGray(img, *pixels)
	var flattens [256]float32
	phash.DCT2DHash256(*pixels, &flattens)
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

	// 8x8 needs only 64 floats (256 bytes): keep it on the stack instead of
	// borrowing a pool.
	var buf [ahashSide * ahashSide]float32
	pixels := buf[:]
	phash.ImageToGray(img, pixels)

	return Ahash(bits64(pixels, meanOfPixels(pixels))), nil
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

// Pixel pools hold reusable grayscale buffers. Buffers traffic as pointers:
// boxing a slice header into the pool would allocate 24 bytes per call,
// while a pointer is a single word. A failed assertion means a foreign value
// was Put back.
var (
	pixelsPool64 = sync.Pool{
		New: func() any {
			p := make([]float32, phash64Side*phash64Side)
			return &p
		},
	}
	pixelsPool256 = sync.Pool{
		New: func() any {
			p := make([]float32, phash256Side*phash256Side)
			return &p
		},
	}
)

// borrowPixels takes a grayscale buffer from pool; return it with
// pool.Put(p) when done.
func borrowPixels(pool *sync.Pool) (*[]float32, error) {
	p, ok := pool.Get().(*[]float32)
	if !ok || p == nil {
		return nil, ErrPixelPool
	}
	return p, nil
}

// bits64 packs values into a 64-bit hash, setting bit (63-idx) for each value
// above thresh.
func bits64(values []float32, thresh float32) uint64 {
	var hash uint64
	for idx, p := range values {
		if p > thresh {
			hash |= 1 << (63 - idx) // leftShiftSet
		}
	}
	return hash
}

func meanOfPixels(pixels []float32) float32 {
	var sum float32
	for _, p := range pixels {
		sum += p
	}
	return sum / float32(len(pixels))
}

// Distance between Phash values
func (ph PHash64) Distance(hash PHash64) uint8 {
	return uint8(bits.OnesCount64(uint64(ph) ^ uint64(hash))) //nolint:gosec // popcnt is bounded to [0,64].
}

func (ph PHash64) String() string {
	var raw [8]byte
	binary.BigEndian.PutUint64(raw[:], uint64(ph))
	var out [2 + 16]byte
	out[0], out[1] = 'p', ':'
	hex.Encode(out[2:], raw[:])
	return string(out[:])
}

func (ph PHash64) Encode(dst []byte) {
	if len(dst) < 8 {
		panic("imagehash: PHash64.Encode requires dst len >= 8")
	}
	binary.LittleEndian.PutUint64(dst[:8], uint64(ph))
}

func (ph *PHash64) Decode(src []byte) {
	if len(src) < 8 {
		panic("imagehash: PHash64.Decode requires src len >= 8")
	}
	*ph = PHash64(binary.LittleEndian.Uint64(src[:8]))
}

// Distance between Phash values
func (ph PHash256) Distance(hash PHash256) uint {
	return uint(
		bits.OnesCount64(ph[0]^hash[0]) +
			bits.OnesCount64(ph[1]^hash[1]) +
			bits.OnesCount64(ph[2]^hash[2]) +
			bits.OnesCount64(ph[3]^hash[3]))
}

func (ph PHash256) String() string {
	var raw [32]byte
	for i, w := range ph {
		binary.BigEndian.PutUint64(raw[i*8:], w)
	}
	var out [2 + 64]byte
	out[0], out[1] = 'p', ':'
	hex.Encode(out[2:], raw[:])
	return string(out[:])
}

func (ph PHash256) Encode(buf []byte) {
	if len(buf) < 32 {
		panic("imagehash: PHash256.Encode requires buf len >= 32")
	}
	binary.LittleEndian.PutUint64(buf[:8], ph[0])
	binary.LittleEndian.PutUint64(buf[8*1:], ph[1])
	binary.LittleEndian.PutUint64(buf[8*2:], ph[2])
	binary.LittleEndian.PutUint64(buf[8*3:], ph[3])
}

func (ph *PHash256) Decode(buf []byte) {
	if len(buf) < 32 {
		panic("imagehash: PHash256.Decode requires buf len >= 32")
	}
	ph[0] = binary.LittleEndian.Uint64(buf[:8])
	ph[1] = binary.LittleEndian.Uint64(buf[8*1:])
	ph[2] = binary.LittleEndian.Uint64(buf[8*2:])
	ph[3] = binary.LittleEndian.Uint64(buf[8*3:])
}
