//go:build arm64

package asm

import (
	"math/rand"
	"testing"
)

// TestPDQLumaKernels checks the NEON luminance kernels against the portable
// implementation. Results must be bit-identical, including clamping
// extremes and chroma pairing.
func TestPDQLumaKernels(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(1))
	edges := []uint8{0, 1, 2, 127, 128, 129, 253, 254, 255}
	for _, n := range []int{8, 16, 64, 256} {
		y := make([]uint8, n)
		cb := make([]uint8, n)
		cr := make([]uint8, n)
		for trial := 0; trial < 50; trial++ {
			for i := range y {
				switch trial % 3 {
				case 0:
					y[i] = edges[rng.Intn(len(edges))]
					cb[i] = edges[rng.Intn(len(edges))]
					cr[i] = edges[rng.Intn(len(edges))]
				default:
					y[i] = uint8(rng.Intn(256))
					cb[i] = uint8(rng.Intn(256))
					cr[i] = uint8(rng.Intn(256))
				}
			}
			got := make([]float32, n)
			PDQLuma444Row(got, y, cb, cr)
			for i := range got {
				if want := pdqLumaRef(y[i], cb[i], cr[i]); got[i] != want {
					t.Fatalf("444 n=%d trial=%d px=%d: got %v want %v (y=%d cb=%d cr=%d)",
						n, trial, i, got[i], want, y[i], cb[i], cr[i])
				}
			}

			got420 := make([]float32, n)
			PDQLuma420Row(got420, y, cb[:n/2], cr[:n/2])
			for i := range got420 {
				if want := pdqLumaRef(y[i], cb[i/2], cr[i/2]); got420[i] != want {
					t.Fatalf("420 n=%d trial=%d px=%d: got %v want %v (y=%d cb=%d cr=%d)",
						n, trial, i, got420[i], want, y[i], cb[i/2], cr[i/2])
				}
			}
		}
	}
}
