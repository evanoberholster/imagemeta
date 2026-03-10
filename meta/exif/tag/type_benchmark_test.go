package tag

import "testing"

var isValidBenchSink uint64

var isValidBenchTypesHot = [...]Type{
	TypeByte,
	TypeASCII,
	TypeShort,
	TypeLong,
	TypeRational,
	TypeUndefined,
	TypeSignedShort,
	TypeSignedLong,
	TypeSignedRational,
	TypeFloat,
	TypeDouble,
	TypeASCIINoNul,
	TypeIfd,
	TypeUnknown,
}

var isValidBenchTypesAll = func() [256]Type {
	var out [256]Type
	for i := range out {
		out[i] = Type(i)
	}
	return out
}()

var isValidLookupByType = func() [256]bool {
	var out [256]bool
	for _, t := range isValidBenchTypesHot[:12] {
		out[uint8(t)] = true
	}
	out[uint8(TypeASCIINoNul)] = true
	out[uint8(TypeIfd)] = true
	return out
}()

var isValidBitset256 = [4]uint64{
	0: (uint64(1) << uint(TypeByte)) |
		(uint64(1) << uint(TypeASCII)) |
		(uint64(1) << uint(TypeShort)) |
		(uint64(1) << uint(TypeLong)) |
		(uint64(1) << uint(TypeRational)) |
		(uint64(1) << uint(TypeUndefined)) |
		(uint64(1) << uint(TypeSignedShort)) |
		(uint64(1) << uint(TypeSignedLong)) |
		(uint64(1) << uint(TypeSignedRational)) |
		(uint64(1) << uint(TypeFloat)) |
		(uint64(1) << uint(TypeDouble)),
	3: (uint64(1) << (uint(TypeASCIINoNul) - 192)) |
		(uint64(1) << (uint(TypeIfd) - 192)),
}

func isValidLookup(t Type) bool {
	return isValidLookupByType[uint8(t)]
}

func isValidBitsetFull(t Type) bool {
	u := uint8(t)
	return (isValidBitset256[u>>6] & (uint64(1) << (u & 63))) != 0
}

func BenchmarkTypeIsValidCurrentHot(b *testing.B) {
	in := isValidBenchTypesHot[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if in[idx].IsValid() {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}

func BenchmarkTypeIsValidLookupHot(b *testing.B) {
	in := isValidBenchTypesHot[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isValidLookup(in[idx]) {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}

func BenchmarkTypeIsValidBitsetFullHot(b *testing.B) {
	in := isValidBenchTypesHot[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isValidBitsetFull(in[idx]) {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}

func BenchmarkTypeIsValidCurrentAll(b *testing.B) {
	in := isValidBenchTypesAll[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if in[idx].IsValid() {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}

func BenchmarkTypeIsValidLookupAll(b *testing.B) {
	in := isValidBenchTypesAll[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isValidLookup(in[idx]) {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}

func BenchmarkTypeIsValidBitsetFullAll(b *testing.B) {
	in := isValidBenchTypesAll[:]
	var hits uint64
	idx := 0
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isValidBitsetFull(in[idx]) {
			hits++
		}
		idx++
		if idx == len(in) {
			idx = 0
		}
	}
	isValidBenchSink = hits
}
