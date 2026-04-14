package exif

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"strconv"
	"testing"
	"unsafe"

	metacanon "github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
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

func benchmarkUint32WordsToAlignedBytesUnsafe(words []uint32, bo utils.ByteOrder) []byte {
	aligned := make([]uint32, len(words))
	copy(aligned, words)
	if bo == utils.BigEndian {
		for i := range aligned {
			aligned[i] = bits.ReverseBytes32(aligned[i])
		}
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(aligned))), len(aligned)*4)
}

func benchmarkParseInt32Current(bo utils.ByteOrder, buf []byte, dst []int32) int {
	n := len(buf) / 4
	if n > len(dst) {
		n = len(dst)
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+4 {
		dst[i] = int32(bo.Uint32(buf[j : j+4]))
	}
	return n
}

func benchmarkParseInt32Hoisted(bo utils.ByteOrder, buf []byte, dst []int32) int {
	n := len(buf) / 4
	if n > len(dst) {
		n = len(dst)
	}
	if bo == utils.BigEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+4 {
			dst[i] = int32(binary.BigEndian.Uint32(buf[j:]))
		}
		return n
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+4 {
		dst[i] = int32(binary.LittleEndian.Uint32(buf[j:]))
	}
	return n
}

func benchmarkParseInt32BinaryOrder(bo utils.ByteOrder, buf []byte, dst []int32) int {
	n := len(buf) / 4
	if n > len(dst) {
		n = len(dst)
	}
	var order binary.ByteOrder = binary.LittleEndian
	if bo == utils.BigEndian {
		order = binary.BigEndian
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+4 {
		dst[i] = int32(order.Uint32(buf[j:]))
	}
	return n
}

func benchmarkParseInt32Manual(bo utils.ByteOrder, buf []byte, dst []int32) int {
	n := len(buf) / 4
	if n > len(dst) {
		n = len(dst)
	}
	if bo == utils.BigEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+4 {
			dst[i] = int32(uint32(buf[j])<<24 |
				uint32(buf[j+1])<<16 |
				uint32(buf[j+2])<<8 |
				uint32(buf[j+3]))
		}
		return n
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+4 {
		dst[i] = int32(uint32(buf[j]) |
			uint32(buf[j+1])<<8 |
			uint32(buf[j+2])<<16 |
			uint32(buf[j+3])<<24)
	}
	return n
}

func benchmarkParseInt32Unsafe(bo utils.ByteOrder, buf []byte, dst []int32) int {
	n := len(buf) / 4
	if n > len(dst) {
		n = len(dst)
	}
	u32 := unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(buf))), n)
	if bo == utils.BigEndian {
		for i := 0; i < n; i++ {
			dst[i] = int32(bits.ReverseBytes32(u32[i]))
		}
		return n
	}
	for i := 0; i < n; i++ {
		dst[i] = int32(u32[i])
	}
	return n
}

func BenchmarkFillCanonAFInfo2(b *testing.B) {
	words := benchmarkAFInfo2Words(105)
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonAFInfo2),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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
			tag.IFD0,
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
			tag.ExifIFD,
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
		tag.MakerNoteIFD,
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

func BenchmarkInt32DecodeVariants(b *testing.B) {
	decoders := []struct {
		name string
		fn   func(utils.ByteOrder, []byte, []int32) int
	}{
		{name: "Current", fn: benchmarkParseInt32Current},
		{name: "Hoisted", fn: benchmarkParseInt32Hoisted},
		{name: "BinaryOrder", fn: benchmarkParseInt32BinaryOrder},
		{name: "Manual", fn: benchmarkParseInt32Manual},
		{name: "Unsafe", fn: benchmarkParseInt32Unsafe},
	}
	byteOrders := []struct {
		name string
		bo   utils.ByteOrder
	}{
		{name: "LE", bo: utils.LittleEndian},
		{name: "BE", bo: utils.BigEndian},
	}
	sizes := []int{10, 20, 100}

	for _, byteOrder := range byteOrders {
		for _, size := range sizes {
			words := make([]uint32, size)
			for i := range words {
				words[i] = uint32(i*7 + 3)
			}
			dst := make([]int32, size)
			for _, decoder := range decoders {
				raw := benchmarkUint32WordsToBytes(words, byteOrder.bo)
				if decoder.name == "Unsafe" {
					raw = benchmarkUint32WordsToAlignedBytesUnsafe(words, byteOrder.bo)
				}
				b.Run(byteOrder.name+"/"+decoder.name+"/n="+strconv.Itoa(size), func(b *testing.B) {
					b.ReportAllocs()
					b.SetBytes(int64(len(raw)))
					for i := 0; i < b.N; i++ {
						n := decoder.fn(byteOrder.bo, raw, dst)
						if n > 0 {
							sinkHotpathInt32 = dst[n-1]
						}
					}
				})
			}
		}
	}
}
