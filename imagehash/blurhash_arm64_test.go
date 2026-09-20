//go:build arm64

package imagehash

import (
	"math/rand"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// TestAsmBlurRow checks the NEON basis kernel against the portable
// implementation. Results must be bit-identical.
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
			if got[i] != want[i] {
				t.Fatalf("trial %d: index %d got %v want %v", trial, i, got[i], want[i])
			}
		}
	}
}
