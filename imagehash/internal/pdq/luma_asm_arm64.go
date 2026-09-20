//go:build arm64

package pdq

import (
	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

// init enables the NEON luminance kernels on ARM64.
func init() {
	useASMLuma = true
	pdqLuma444Row = asm.PDQLuma444Row
	pdqLuma420Row = asm.PDQLuma420Row
}
