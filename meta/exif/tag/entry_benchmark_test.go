package tag

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var isEmbeddedBenchSink uint64
var embeddedValueBenchSink uint64

var embeddedBenchEntries = [...]Entry{
	NewEntry(TagMake, TypeASCII, 4, 0, ifd.IFD0, 0, utils.LittleEndian),               // true
	NewEntry(TagMake, TypeASCII, 8, 0, ifd.IFD0, 0, utils.LittleEndian),               // false
	NewEntry(TagDateTime, TypeShort, 2, 0, ifd.IFD0, 0, utils.LittleEndian),           // true
	NewEntry(TagDateTime, TypeShort, 3, 0, ifd.IFD0, 0, utils.LittleEndian),           // false
	NewEntry(TagExposureTime, TypeLong, 1, 0, ifd.ExifIFD, 0, utils.LittleEndian),     // true
	NewEntry(TagExposureTime, TypeLong, 2, 0, ifd.ExifIFD, 0, utils.LittleEndian),     // false
	NewEntry(TagInteropIFDPointer, TypeIfd, 1, 0, ifd.ExifIFD, 0, utils.LittleEndian), // false
	NewEntry(TagMake, TypeASCIINoNul, 4, 0, ifd.IFD0, 0, utils.LittleEndian),          // true
	NewEntry(TagMake, TypeASCIINoNul, 5, 0, ifd.IFD0, 0, utils.LittleEndian),          // false
}

const (
	benchEmbeddedType1Word0 uint64 = (uint64(1) << uint(TypeByte)) |
		(uint64(1) << uint(TypeASCII)) |
		(uint64(1) << uint(TypeUndefined))
	benchEmbeddedType2Word0 uint64 = (uint64(1) << uint(TypeShort)) |
		(uint64(1) << uint(TypeSignedShort))
	benchEmbeddedType4Word0 uint64 = (uint64(1) << uint(TypeLong)) |
		(uint64(1) << uint(TypeSignedLong)) |
		(uint64(1) << uint(TypeFloat))
	benchEmbeddedType1Word3 uint64 = (uint64(1) << (uint(TypeASCIINoNul) - 192))
)

func isEmbeddedBitsetOld(e Entry) bool {
	typ := uint8(e.Type)
	if typ < 64 {
		mask := uint64(1) << typ
		if benchEmbeddedType1Word0&mask != 0 {
			return e.UnitCount <= 4
		}
		if benchEmbeddedType2Word0&mask != 0 {
			return e.UnitCount <= 2
		}
		if benchEmbeddedType4Word0&mask != 0 {
			return e.UnitCount <= 1
		}
		return false
	}
	if typ >= 192 {
		mask := uint64(1) << (typ - 192)
		return (benchEmbeddedType1Word3&mask) != 0 && e.UnitCount <= 4
	}
	return false
}

func isEmbeddedLegacy(e Entry) bool {
	return e.Size() <= 4 && e.Type != TypeIfd
}

func BenchmarkEntryIsEmbeddedLookup(b *testing.B) {
	entries := embeddedBenchEntries[:]
	var hit uint64
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if entries[idx].IsEmbedded() {
			hit++
		}
		idx++
		if idx == len(entries) {
			idx = 0
		}
	}
	isEmbeddedBenchSink = hit
}

func BenchmarkEntryIsEmbeddedBitsetOld(b *testing.B) {
	entries := embeddedBenchEntries[:]
	var hit uint64
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isEmbeddedBitsetOld(entries[idx]) {
			hit++
		}
		idx++
		if idx == len(entries) {
			idx = 0
		}
	}
	isEmbeddedBenchSink = hit
}

func BenchmarkEntryIsEmbeddedLegacy(b *testing.B) {
	entries := embeddedBenchEntries[:]
	var hit uint64
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if isEmbeddedLegacy(entries[idx]) {
			hit++
		}
		idx++
		if idx == len(entries) {
			idx = 0
		}
	}
	isEmbeddedBenchSink = hit
}

func embeddedShortRoundTrip(e Entry) uint16 {
	var buf [4]byte
	e.EmbeddedValue(buf[:])
	return e.ByteOrder.Uint16(buf[:2])
}

func embeddedLongRoundTrip(e Entry) uint32 {
	var buf [4]byte
	e.EmbeddedValue(buf[:])
	return e.ByteOrder.Uint32(buf[:4])
}

func BenchmarkEntryEmbeddedShortDirectLE(b *testing.B) {
	e := NewEntry(TagMake, TypeShort, 1, 0x04030201, ifd.IFD0, 0, utils.LittleEndian)
	var sum uint64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += uint64(e.EmbeddedShort())
	}
	embeddedValueBenchSink = sum
}

func BenchmarkEntryEmbeddedShortRoundTripLE(b *testing.B) {
	e := NewEntry(TagMake, TypeShort, 1, 0x04030201, ifd.IFD0, 0, utils.LittleEndian)
	var sum uint64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += uint64(embeddedShortRoundTrip(e))
	}
	embeddedValueBenchSink = sum
}

func BenchmarkEntryEmbeddedLongDirectLE(b *testing.B) {
	e := NewEntry(TagMake, TypeLong, 1, 0x04030201, ifd.IFD0, 0, utils.LittleEndian)
	var sum uint64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += uint64(e.EmbeddedLong())
	}
	embeddedValueBenchSink = sum
}

func BenchmarkEntryEmbeddedLongRoundTripLE(b *testing.B) {
	e := NewEntry(TagMake, TypeLong, 1, 0x04030201, ifd.IFD0, 0, utils.LittleEndian)
	var sum uint64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += uint64(embeddedLongRoundTrip(e))
	}
	embeddedValueBenchSink = sum
}
