//go:build linux && amd64

package pdq

import (
	cpu "github.com/klauspost/cpuid/v2"

	"github.com/evanoberholster/imagemeta/imagehash/internal/asm"
)

func init() {
	// The AVX2 kernels use separate multiply/add (matching the scalar
	// codegen on amd64), so FMA support is not required.
	if cpu.CPU.Supports(cpu.AVX, cpu.AVX2, cpu.SSE, cpu.SSE2, cpu.SSE4) {
		useASMLuma = true
		pdqLuma444Row = asm.PDQLuma444Row
		pdqLuma420Row = asm.PDQLuma420Row
	}
}
