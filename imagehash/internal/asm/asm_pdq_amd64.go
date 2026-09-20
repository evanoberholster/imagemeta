//go:build amd64

package asm

// asmPDQLuma444Row computes PDQ luminance for len(out) 1:1-mapped YCbCr
// pixels. len(out) must be a positive multiple of 8; y, cb and cr must each
// hold at least len(out) bytes.
//
//go:noescape
func asmPDQLuma444Row(out []float32, y, cb, cr []uint8)

// asmPDQLuma420Row computes PDQ luminance for len(out) Y pixels sharing
// len(out)/2 Cb and Cr pixels (horizontal 2:1 subsampling, even origin).
// len(out) must be a positive multiple of 8; y must hold at least len(out)
// bytes and cb/cr at least len(out)/2 bytes each.
//
//go:noescape
func asmPDQLuma420Row(out []float32, y, cb, cr []uint8)
