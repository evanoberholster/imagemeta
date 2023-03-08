package transforms32

import (
	"math/rand"
	"testing"
)

func TestASM16(t *testing.T) {
	source := make([]float32, 64)
	input := make([]float32, len(source))
	input2 := make([]float32, len(source))
	for i := 0; i < 16; i++ {
		source[i] = rand.Float32()
	}

	copy(input2, source)
	ForwardDCT64(input2)

	copy(input, source)
	asmForwardDCT64(input)

	for i := 0; i < len(source); i++ {
		if input[i] != input2[i] {
			t.Error("error", i)
			t.Error(input)
			t.Error(input2)
			break
		}
	}

}
