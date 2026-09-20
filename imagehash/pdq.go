package imagehash

// Copyright 2024 Evan Oberholster
// All rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"math/bits"

	"github.com/evanoberholster/imagemeta/imagehash/internal/pdq"
)

// PDQHash is a 256-bit PDQ (Perceptual Difference Quantization) perceptual
// hash. It is stored as four uint64 words in big-endian order so that String
// produces Meta's canonical 64-character hex representation.
type PDQHash [4]uint64

// PDQMatchThreshold is Meta's recommended maximum Hamming distance for
// considering two PDQ hashes a match.
const PDQMatchThreshold = pdq.DefaultMatchThreshold

// PDQQualityThreshold is Meta's recommended minimum quality score for a PDQ
// hash to be considered reliable.
const PDQQualityThreshold = pdq.DefaultQualityThreshold

// NewPDQ256 computes the 256-bit PDQ perceptual hash of img using Meta's PDQ
// algorithm. img may be any size; the algorithm handles downscaling internally.
func NewPDQ256(img image.Image) (PDQHash, error) {
	hash, _, err := NewPDQ256WithQuality(img)
	return hash, err
}

// NewPDQ256WithQuality is like NewPDQ256 but also returns the PDQ quality
// score (0-100). Per Meta, hashes with a quality at or below
// PDQQualityThreshold are unreliable.
func NewPDQ256WithQuality(img image.Image) (PDQHash, int, error) {
	if img == nil {
		return PDQHash{}, 0, ErrImageObject
	}
	result, err := pdq.Hash(img)
	if err != nil {
		return PDQHash{}, 0, err
	}
	return pdqHashFromBytes(result.Hash), result.Quality, nil
}

// Distance returns the Hamming distance between two PDQ hashes.
func (h PDQHash) Distance(other PDQHash) uint {
	return uint(
		bits.OnesCount64(h[0]^other[0]) +
			bits.OnesCount64(h[1]^other[1]) +
			bits.OnesCount64(h[2]^other[2]) +
			bits.OnesCount64(h[3]^other[3]))
}

// String returns the hash as a 64-character lowercase hex string, matching
// Meta's canonical PDQ representation.
func (h PDQHash) String() string {
	var raw [32]byte
	for i, w := range h {
		binary.BigEndian.PutUint64(raw[i*8:], w)
	}
	var out [64]byte
	hex.Encode(out[:], raw[:])
	return string(out[:])
}

// Encode writes the big-endian 32-byte representation of the hash to dst.
func (h PDQHash) Encode(dst []byte) {
	if len(dst) < 32 {
		panic("imagehash: PDQHash.Encode requires dst len >= 32")
	}
	binary.BigEndian.PutUint64(dst[:8], h[0])
	binary.BigEndian.PutUint64(dst[8*1:], h[1])
	binary.BigEndian.PutUint64(dst[8*2:], h[2])
	binary.BigEndian.PutUint64(dst[8*3:], h[3])
}

// Decode reads the big-endian 32-byte representation of the hash from src.
func (h *PDQHash) Decode(src []byte) {
	if len(src) < 32 {
		panic("imagehash: PDQHash.Decode requires src len >= 32")
	}
	h[0] = binary.BigEndian.Uint64(src[:8])
	h[1] = binary.BigEndian.Uint64(src[8*1:])
	h[2] = binary.BigEndian.Uint64(src[8*2:])
	h[3] = binary.BigEndian.Uint64(src[8*3:])
}

// pdqHashFromBytes converts the 32-byte big-endian PDQ representation into a
// PDQHash.
func pdqHashFromBytes(h [32]byte) PDQHash {
	return PDQHash{
		binary.BigEndian.Uint64(h[0:8]),
		binary.BigEndian.Uint64(h[8:16]),
		binary.BigEndian.Uint64(h[16:24]),
		binary.BigEndian.Uint64(h[24:32]),
	}
}

// pdqHashFromHex parses a 64-character hex PDQ hash.
func pdqHashFromHex(s string) (PDQHash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PDQHash{}, err
	}
	if len(b) != 32 {
		return PDQHash{}, fmt.Errorf("imagehash: pdq hash must be 32 bytes, got %d", len(b))
	}
	return pdqHashFromBytes([32]byte(b)), nil
}
