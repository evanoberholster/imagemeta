//go:build !arm64

package imagehash

// blurRow is the portable implementation on architectures without a BlurHash
// kernel.
func blurRow(out, lr, lg, lb []float32) {
	blurRowGo(out, lr, lg, lb)
}
