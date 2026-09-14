// Ported from github.com/MatthewSH/pdq (MIT License, Copyright (c) 2026 Matt Hatcher).
// See the LICENSE file in this directory.

package pdq

import (
	"image"
	"sync"
)

const maxPixels = ImageSize * ImageSize

// workspace holds the reusable scratch buffers for a single PDQ computation.
// Every buffer is sized for the largest processing image (ImageSize) and is
// sliced down to the actual dimensions. Reusing the workspace through a pool
// keeps the steady-state hot path allocation free, mirroring the fixed-size
// structure approach used by the PHash implementation in imagehash.
type workspace struct {
	rgba     *image.RGBA
	luma     []float32
	a        []float32
	b        []float32
	filtered [outSize * outSize]float32
	dct      [dctRows * dctRows]float32
	dctAt    [dctCols * dctCols]float32
	dctT     [dctRows * dctCols]float32
}

var workspacePool = sync.Pool{
	New: func() any { return newWorkspace() },
}

// newWorkspace allocates a workspace sized for the largest processing image.
func newWorkspace() *workspace {
	return &workspace{
		rgba: image.NewRGBA(image.Rect(0, 0, ImageSize, ImageSize)),
		luma: make([]float32, maxPixels),
		a:    make([]float32, maxPixels),
		b:    make([]float32, maxPixels),
	}
}
