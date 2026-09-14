//go:build arm64

package imagehash

import (
	"math"
	"math/rand"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// TestAsmBlurRow checks the NEON basis kernel against the portable
// implementation.
func TestAsmBlurRow(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for trial := 0; trial < 200; trial++ {
		var lr, lg, lb [width]float32
		for i := 0; i < width; i++ {
			lr[i] = rng.Float32() * 2
			lg[i] = rng.Float32() * 2
			lb[i] = rng.Float32() * 2
		}
		var got, want [3 * xComponents]float32
		asm.BlurRow(got[:], lr[:], lg[:], lb[:], xvaluesT[:])
		blurRowGo(want[:], lr[:], lg[:], lb[:])
		for i := range got {
			if d := math.Abs(float64(got[i] - want[i])); d > 1e-4 {
				t.Fatalf("trial %d: index %d got %v want %v (diff %v)", trial, i, got[i], want[i], d)
			}
		}
	}
}
