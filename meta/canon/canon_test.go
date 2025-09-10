package canon

import (
	"testing"
)

func BenchmarkModelIDString(b *testing.B) {
	var m CanonModelID
	m = CanonModelID(250)
	for i := 0; i < b.N; i++ {
		_ = m.String()
	}
}

func BenchmarkCanonRFLensTypeString(b *testing.B) {
	var m CanonRFLensType
	m = CanonRFLensType(0x0010)
	for i := 0; i < b.N; i++ {
		_ = m.String()
	}
}

func BenchmarkCanonLensTypeString(b *testing.B) {
	var m CanonLensType
	m = CanonLensType(124)
	for i := 0; i < b.N; i++ {
		_ = m.String()
	}
}
