//go:build arm64

package imagehash

import "github.com/evanoberholster/imagemeta/imagehash/internal/asm"

// blurRow uses the NEON kernel on arm64. width is fixed at 64, which the
// kernel relies on.
func blurRow(out, lr, lg, lb []float32) {
	asm.BlurRow(out, lr, lg, lb, xvaluesT[:])
}
