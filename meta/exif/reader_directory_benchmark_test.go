package exif

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
	"github.com/rs/zerolog"
)

var directoryHeaderParseSink uint32

type ifdHeaderSpec struct {
	id         tag.ID
	typ        tag.Type
	unitCount  uint32
	valueOrOff uint32
}

func buildIFDHeaderPayload(specs []ifdHeaderSpec) []byte {
	payload := make([]byte, len(specs)*12)
	for i := 0; i < len(specs); i++ {
		base := i * 12
		s := specs[i]
		utils.LittleEndian.PutUint16(payload[base:base+2], uint16(s.id))
		utils.LittleEndian.PutUint16(payload[base+2:base+4], uint16(s.typ))
		utils.LittleEndian.PutUint32(payload[base+4:base+8], s.unitCount)
		utils.LittleEndian.PutUint32(payload[base+8:base+12], s.valueOrOff)
	}
	return payload
}

// benchmarkDirectoryTagHeaderParse benchmarks directory-entry header parsing paths.
func benchmarkDirectoryTagHeaderParse(b *testing.B, fn func(*Reader, ifd.Directory, uint16) error, tagCount uint16) {
	specs := make([]ifdHeaderSpec, int(tagCount))
	for i := 0; i < int(tagCount); i++ {
		specs[i] = ifdHeaderSpec{
			id:         tag.TagImageWidth,
			typ:        tag.TypeLong,
			unitCount:  2,
			valueOrOff: uint32(4096 + i*8),
		}
	}
	payload := buildIFDHeaderPayload(specs)

	var raw bytes.Reader
	br := bufio.NewReaderSize(&raw, len(payload))
	r := NewReader(zerolog.Nop())
	defer r.Close()

	directory := ifd.New(utils.LittleEndian, ifd.IFD0, 0, 0, 0)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw.Reset(payload)
		br.Reset(&raw)
		r.Reset(br)
		if err := fn(r, directory, tagCount); err != nil {
			b.Fatal(err)
		}
		directoryHeaderParseSink += r.state.len
	}
}

// BenchmarkReadDirectoryTagHeaders compares per-entry reads against read-once parsing.
func BenchmarkReadDirectoryTagHeaders(b *testing.B) {
	const tagCount uint16 = 64

	b.Run("PerEntryRead12", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParse(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersPerEntry(d, n)
		}, tagCount)
	})

	b.Run("BulkReadOnce", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParse(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersBulk(d, n)
		}, tagCount)
	})

	b.Run("BulkTrustedInline", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParse(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersBulkTrusted(d, n)
		}, tagCount)
	})
}

// benchmarkDirectoryTagHeaderParseFromPayload benchmarks parsing on a fixed header payload.
func benchmarkDirectoryTagHeaderParseFromPayload(b *testing.B, fn func(*Reader, ifd.Directory, uint16) error, payload []byte) {
	var raw bytes.Reader
	br := bufio.NewReaderSize(&raw, len(payload))
	r := NewReader(zerolog.Nop())
	defer r.Close()

	directory := ifd.New(utils.LittleEndian, ifd.IFD0, 0, 0, 0)
	tagCount := uint16(len(payload) / 12)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw.Reset(payload)
		br.Reset(&raw)
		r.Reset(br)
		if err := fn(r, directory, tagCount); err != nil {
			b.Fatal(err)
		}
		directoryHeaderParseSink += r.state.len
	}
}

// BenchmarkReadDirectoryTagHeadersCR2IFD0Like models a production-like CR2 IFD0 entry distribution.
func BenchmarkReadDirectoryTagHeadersCR2IFD0Like(b *testing.B) {
	// Representative CR2 IFD0-like mix: short/long embedded scalars, rational/ascii offsets,
	// and Exif/GPS pointer tags (mapped to TypeIfd and treated as non-embedded).
	specs := []ifdHeaderSpec{
		{id: tag.TagSubfileType, typ: tag.TypeLong, unitCount: 1, valueOrOff: 0},
		{id: tag.TagImageWidth, typ: tag.TypeShort, unitCount: 1, valueOrOff: 5472},
		{id: tag.TagImageLength, typ: tag.TypeShort, unitCount: 1, valueOrOff: 3648},
		{id: tag.TagBitsPerSample, typ: tag.TypeShort, unitCount: 3, valueOrOff: 0x280},
		{id: tag.TagCompression, typ: tag.TypeShort, unitCount: 1, valueOrOff: 6},
		{id: tag.TagMake, typ: tag.TypeASCII, unitCount: 6, valueOrOff: 0x2A0},
		{id: tag.TagModel, typ: tag.TypeASCII, unitCount: 14, valueOrOff: 0x2A8},
		{id: tag.TagStripOffsets, typ: tag.TypeLong, unitCount: 1, valueOrOff: 0x1000},
		{id: tag.TagOrientation, typ: tag.TypeShort, unitCount: 1, valueOrOff: 1},
		{id: tag.TagSamplesPerPixel, typ: tag.TypeShort, unitCount: 1, valueOrOff: 3},
		{id: tag.TagRowsPerStrip, typ: tag.TypeLong, unitCount: 1, valueOrOff: 3650},
		{id: tag.TagStripByteCounts, typ: tag.TypeLong, unitCount: 1, valueOrOff: 909916},
		{id: tag.TagXResolution, typ: tag.TypeRational, unitCount: 1, valueOrOff: 0x2C0},
		{id: tag.TagYResolution, typ: tag.TypeRational, unitCount: 1, valueOrOff: 0x2C8},
		{id: tag.TagPlanarConfiguration, typ: tag.TypeShort, unitCount: 1, valueOrOff: 1},
		{id: tag.TagResolutionUnit, typ: tag.TypeShort, unitCount: 1, valueOrOff: 2},
		{id: tag.TagSoftware, typ: tag.TypeASCII, unitCount: 13, valueOrOff: 0x2D0},
		{id: tag.TagDateTime, typ: tag.TypeASCII, unitCount: 20, valueOrOff: 0x2E0},
		{id: tag.TagArtist, typ: tag.TypeASCII, unitCount: 8, valueOrOff: 0x300},
		{id: tag.TagSubIFDs, typ: tag.TypeLong, unitCount: 2, valueOrOff: 0x310},
		{id: tag.TagCopyright, typ: tag.TypeASCII, unitCount: 16, valueOrOff: 0x320},
		{id: tag.TagExifIFDPointer, typ: tag.TypeLong, unitCount: 1, valueOrOff: 0x400},
		{id: tag.TagGPSIFDPointer, typ: tag.TypeLong, unitCount: 1, valueOrOff: 0x600},
		{id: tag.TagThumbnailOffset, typ: tag.TypeLong, unitCount: 1, valueOrOff: 0x480},
		{id: tag.TagThumbnailLength, typ: tag.TypeLong, unitCount: 1, valueOrOff: 15589},
		{id: tag.TagPrintIM, typ: tag.TypeUndefined, unitCount: 62, valueOrOff: 0x380},
	}
	payload := buildIFDHeaderPayload(specs)

	b.Run("PerEntryRead12", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParseFromPayload(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersPerEntry(d, n)
		}, payload)
	})

	b.Run("BulkReadOnce", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParseFromPayload(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersBulk(d, n)
		}, payload)
	})

	b.Run("BulkTrustedInline", func(b *testing.B) {
		benchmarkDirectoryTagHeaderParseFromPayload(b, func(r *Reader, d ifd.Directory, n uint16) error {
			return r.parseDirectoryTagHeadersBulkTrusted(d, n)
		}, payload)
	})
}
