package isobmff

import "testing"

var (
	brandSink    brand
	fileTypeSink fileTypeBox

	benchmarkBrandInputsKnown = [][]byte{
		[]byte("avci"),
		[]byte("avif"),
		[]byte("crx "),
		[]byte("heic"),
		[]byte("heim"),
		[]byte("heis"),
		[]byte("heix"),
		[]byte("hevc"),
		[]byte("hevm"),
		[]byte("hevs"),
		[]byte("hevx"),
		[]byte("iso8"),
		[]byte("isom"),
		[]byte("M4A "),
		[]byte("MA1B"),
		[]byte("meta"),
		[]byte("miaf"),
		[]byte("MiAn"),
		[]byte("MiBr"),
		[]byte("mif1"),
		[]byte("mif2"),
		[]byte("MiHA"),
		[]byte("MiHB"),
		[]byte("MiHE"),
		[]byte("MiPr"),
		[]byte("mp41"),
		[]byte("mp42"),
		[]byte("msf1"),
	}
	benchmarkBrandInputsMixed = [][]byte{
		[]byte("heic"),
		[]byte("mif1"),
		[]byte("zzzz"),
		[]byte("avif"),
		[]byte("xxxx"),
		[]byte("isom"),
		[]byte("????"),
		[]byte("mp42"),
	}

	// major + minor version + 8 compatible brands
	benchmarkFTYPBuf = []byte("heic0001mif1heicmiafiso8avifmp42msf1isom")
)

func parseFTYPBrands(buf []byte, lookup func([]byte) brand) (ftyp fileTypeBox) {
	ftyp.MajorBrand = lookup(buf[:4])
	copy(ftyp.MinorVersion[:4], buf[4:8])
	for i, compatibleBrand := 8, 0; i+4 <= len(buf) && compatibleBrand < maxBrandCount; compatibleBrand++ {
		ftyp.Compatible[compatibleBrand] = lookup(buf[i : i+4])
		i += 4
	}
	return ftyp
}

func benchmarkBrandLookup(b *testing.B, inputs [][]byte, lookup func([]byte) brand) {
	var out brand
	for i := 0; i < b.N; i++ {
		out = lookup(inputs[i%len(inputs)])
	}
	brandSink = out
}

func benchmarkParseFTYPBrands(b *testing.B, lookup func([]byte) brand) {
	var out fileTypeBox
	for i := 0; i < b.N; i++ {
		out = parseFTYPBrands(benchmarkFTYPBuf, lookup)
	}
	fileTypeSink = out
}

func BenchmarkBrandLookupKnown(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsKnown, brandFromBuf)
}

func BenchmarkBrandLookupMixed(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsMixed, brandFromBuf)
}

func BenchmarkParseFTYPBrands(b *testing.B) {
	benchmarkParseFTYPBrands(b, brandFromBuf)
}
