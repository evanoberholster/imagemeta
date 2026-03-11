package exif

import (
	"bytes"
	"testing"

	metacanon "github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var (
	sinkHotpathAFInfo metacanon.AFInfo
	sinkHotpathString string
	sinkHotpathInt32  int32
)

func benchmarkAFInfo2Words(numAFPoints int) []uint16 {
	if numAFPoints <= 0 {
		numAFPoints = 1
	}
	maskWords := canonBitWordCount(numAFPoints)
	// AFInfo2 layout: fixed(8) + width/height/x/y + inFocus mask + selected mask.
	total := 8 + (4 * numAFPoints) + (2 * maskWords)
	words := make([]uint16, total)

	words[1] = 2 // AFAreaMode
	words[2] = uint16(numAFPoints)
	words[3] = uint16(numAFPoints)
	words[4] = 6000 // CanonImageWidth
	words[5] = 4000 // CanonImageHeight
	words[6] = 6000 // AFImageWidth
	words[7] = 4000 // AFImageHeight

	widthStart := 8
	heightStart := widthStart + numAFPoints
	xStart := heightStart + numAFPoints
	yStart := xStart + numAFPoints
	inFocusStart := yStart + numAFPoints
	selectedStart := inFocusStart + maskWords

	for i := 0; i < numAFPoints; i++ {
		words[widthStart+i] = 30
		words[heightStart+i] = 20
		words[xStart+i] = uint16(i * 5)
		words[yStart+i] = uint16(i * 3)
	}
	for i := 0; i < maskWords; i++ {
		words[inFocusStart+i] = 0x5555
		words[selectedStart+i] = 0x1111
	}

	return words
}

func benchmarkUint32WordsToBytes(words []uint32, bo utils.ByteOrder) []byte {
	out := make([]byte, len(words)*4)
	for i := range words {
		bo.PutUint32(out[i*4:], words[i])
	}
	return out
}

func BenchmarkFillCanonAFInfo2(b *testing.B) {
	words := benchmarkAFInfo2Words(105)
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonAFInfo2),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	run := func(b *testing.B, opt ReaderOption) {
		r := NewReader(Logger, opt)
		defer r.Close()
		var br bytes.Reader
		var dst metacanon.AFInfo

		b.ReportAllocs()
		b.SetBytes(int64(len(raw)))
		for i := 0; i < b.N; i++ {
			br.Reset(raw)
			r.Reset(&br)
			r.Exif.IFD0.Model = "Canon EOS R6"
			dst = r.parseCanonAFInfo2(entry)
		}
		sinkHotpathAFInfo = dst
	}

	b.Run("All", func(b *testing.B) {
		run(b, WithAFInfoDecodeOptions(AFInfoDecodeAll))
	})

	b.Run("BitsetsOnly", func(b *testing.B) {
		run(b, WithAFInfoDecodeOptions(AFInfoDecodeInFocus|AFInfoDecodeSelected))
	})
}

func BenchmarkParseStringAllowUndefined(b *testing.B) {
	b.Run("ASCII", func(b *testing.B) {
		raw := []byte(" Canon EOS R6 Mark II ")
		entry := tag.NewEntry(
			tag.TagModel,
			tag.TypeASCII,
			uint32(len(raw)),
			0,
			ifd.IFD0,
			0,
			utils.LittleEndian,
		)

		r := NewReader(Logger)
		defer r.Close()
		var br bytes.Reader

		b.ReportAllocs()
		b.SetBytes(int64(len(raw)))
		for i := 0; i < b.N; i++ {
			br.Reset(raw)
			r.Reset(&br)
			sinkHotpathString = r.parseStringAllowUndefined(entry)
		}
	})

	b.Run("Undefined", func(b *testing.B) {
		raw := []byte{0, 'C', 'a', 'n', 'o', 'n', 0x01, 'R', '6', ' ', 0, 0xff}
		entry := tag.NewEntry(
			tag.TagMakerNote,
			tag.TypeUndefined,
			uint32(len(raw)),
			0,
			ifd.ExifIFD,
			0,
			utils.LittleEndian,
		)

		r := NewReader(Logger)
		defer r.Close()
		var br bytes.Reader

		b.ReportAllocs()
		b.SetBytes(int64(len(raw)))
		for i := 0; i < b.N; i++ {
			br.Reset(raw)
			r.Reset(&br)
			sinkHotpathString = r.parseStringAllowUndefined(entry)
		}
	})
}

func BenchmarkParseCanonInt32List(b *testing.B) {
	words := make([]uint32, 64)
	for i := range words {
		words[i] = uint32(i*7 + 3)
	}
	raw := benchmarkUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(0x0024),
		tag.TypeLong,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	var dst [64]int32

	b.ReportAllocs()
	b.SetBytes(int64(len(raw)))
	for i := 0; i < b.N; i++ {
		br.Reset(raw)
		r.Reset(&br)
		n := r.parseCanonInt32List(entry, dst[:])
		if n > 0 {
			sinkHotpathInt32 = dst[n-1]
		}
	}
}
