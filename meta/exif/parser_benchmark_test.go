package exif

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var parserBenchSink uint32

// BenchmarkParseTagDispatch benchmarks parser tag dispatch on representative IFD types.
func BenchmarkParseTagDispatch(b *testing.B) {
	b.Run("IFD0", func(b *testing.B) {
		r := &Reader{state: &state{}}
		entry := tag.NewEntry(tag.TagImageWidth, tag.TypeLong, 1, 4000, ifd.IFD0, 0, utils.LittleEndian)

		b.ReportAllocs()
		b.SetBytes(12)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.parseTag(entry)
		}
		parserBenchSink += r.Exif.IFD0.ImageWidth
	})

	b.Run("ExifIFD", func(b *testing.B) {
		r := &Reader{state: &state{}}
		entry := tag.NewEntry(tag.TagExposureProgram, tag.TypeShort, 1, 2, ifd.ExifIFD, 0, utils.LittleEndian)

		b.ReportAllocs()
		b.SetBytes(12)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.parseTag(entry)
		}
		parserBenchSink += uint32(r.Exif.ExifIFD.ExposureProgram)
	})

	b.Run("GPSIFD", func(b *testing.B) {
		r := &Reader{state: &state{}}
		entry := tag.NewEntry(tag.TagGPSDifferential, tag.TypeShort, 1, 1, ifd.GPSIFD, 0, utils.LittleEndian)

		b.ReportAllocs()
		b.SetBytes(12)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.parseTag(entry)
		}
		parserBenchSink += uint32(r.Exif.GPS.Differential())
	})
}

// BenchmarkParseSubIFDsEmbedded benchmarks parsing an embedded SubIFD pointer.
func BenchmarkParseSubIFDsEmbedded(b *testing.B) {
	r := &Reader{state: &state{}}
	entry := tag.NewEntry(tag.TagSubIFDs, tag.TypeLong, 1, 0x1234, ifd.IFD0, 0, utils.LittleEndian)

	b.ReportAllocs()
	b.SetBytes(12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.state.reset()
		r.Exif = Exif{}
		r.parseSubIFDs(entry)
		parserBenchSink += r.state.len
	}
}

// BenchmarkParseTagStreamedNonEmbedded benchmarks parsing tags that require streamed reads.
func BenchmarkParseTagStreamedNonEmbedded(b *testing.B) {
	var payload [64]byte
	// ExifIFD ExposureTime (RATIONAL) at offset 32: 1/200 second.
	utils.LittleEndian.PutUint32(payload[32:36], 1)
	utils.LittleEndian.PutUint32(payload[36:40], 200)

	var raw bytes.Reader
	br := bufio.NewReaderSize(&raw, 64)
	r := &Reader{
		reader:     br,
		state:      &state{},
		exifLength: uint32(len(payload)),
	}
	entry := tag.NewEntry(tag.TagExposureTime, tag.TypeRational, 1, 32, ifd.ExifIFD, 0, utils.LittleEndian)

	b.ReportAllocs()
	b.SetBytes(int64(entry.Size()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw.Reset(payload[:])
		br.Reset(&raw)
		r.po = 0
		r.Exif = Exif{}
		r.parseTag(entry)
		parserBenchSink += uint32(r.Exif.ExifIFD.ExposureTime * 1000)
	}
}
