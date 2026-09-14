//go:build arm64

package asm

//go:noescape
func asmBlurRow(out, lr, lg, lb, xvaluesT []float32)
