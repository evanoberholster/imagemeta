// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"encoding/hex"
)

// Hash256 is a 256-bit PDQ perceptual hash stored as 32 bytes.
type Hash256 [HashSize]byte

func (h Hash256) Bytes() []byte {
	b := make([]byte, HashSize)
	copy(b, h[:])

	return b
}

// String returns the hash as a 64-character lowercase hex string.
func (h Hash256) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash256) IsZero() bool {
	return h == Hash256{}
}

// Distance calculates the Hamming distance between the current 256-bit hash and another 256-bit hash.
func (h Hash256) Distance(otherHash Hash256) int {
	var distance int

	for i := range h {
		x := h[i] ^ otherHash[i]

		for x != 0 {
			distance += int(x & 1)
			x >>= 1
		}
	}

	return distance
}
