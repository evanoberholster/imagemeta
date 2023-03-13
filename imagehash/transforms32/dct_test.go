package transforms32

import (
	"math/rand"
	"testing"
)

func TestDCT2DHash64(t *testing.T) {
	source := make([]float32, 4096)
	input := make([]float32, len(source))
	input2 := make([]float32, len(source))
	for i := 0; i < len(source); i++ {
		source[i] = rand.Float32()
	}

	copy(input, source)
	fl := asmDCT2DHash64(input)

	copy(input2, source)
	flattens := DCT2DHash64(input2)

	for i := 0; i < len(source); i++ {
		if input[i] != input2[i] {
			t.Error("error", i)
			t.Error(input[i:64])
			t.Error(input2[i:64])
			break
		}
	}
	for i := 0; i < len(flattens); i++ {
		if flattens[i] != fl[i] {
			t.Error("error", i)
			t.Error(flattens[:64])
			t.Error(fl[:64])
			break
		}
	}

}

func BenchmarkDCT2DHash64(b *testing.B) {
	source := make([]float32, 4096)
	for i := 0; i < len(source); i++ {
		source[i] = rand.Float32()
	}
	input := make([]float32, len(source))

	b.Run("ASM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			copy(input, source)
			asmDCT2DHash64(input)
		}
	})
	FlagUseASM = false
	b.Run("FN", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			copy(input, source)
			DCT2DHash64(input)
		}
	})
}
