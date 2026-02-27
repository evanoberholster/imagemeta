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

const (
	fourCCAvci = uint32('a')<<24 | uint32('v')<<16 | uint32('c')<<8 | uint32('i')
	fourCCAvif = uint32('a')<<24 | uint32('v')<<16 | uint32('i')<<8 | uint32('f')
	fourCCCrx  = uint32('c')<<24 | uint32('r')<<16 | uint32('x')<<8 | uint32(' ')
	fourCCHeic = uint32('h')<<24 | uint32('e')<<16 | uint32('i')<<8 | uint32('c')
	fourCCHeim = uint32('h')<<24 | uint32('e')<<16 | uint32('i')<<8 | uint32('m')
	fourCCHeis = uint32('h')<<24 | uint32('e')<<16 | uint32('i')<<8 | uint32('s')
	fourCCHeix = uint32('h')<<24 | uint32('e')<<16 | uint32('i')<<8 | uint32('x')
	fourCCHevc = uint32('h')<<24 | uint32('e')<<16 | uint32('v')<<8 | uint32('c')
	fourCCHevm = uint32('h')<<24 | uint32('e')<<16 | uint32('v')<<8 | uint32('m')
	fourCCHevs = uint32('h')<<24 | uint32('e')<<16 | uint32('v')<<8 | uint32('s')
	fourCCHevx = uint32('h')<<24 | uint32('e')<<16 | uint32('v')<<8 | uint32('x')
	fourCCIso8 = uint32('i')<<24 | uint32('s')<<16 | uint32('o')<<8 | uint32('8')
	fourCCIsom = uint32('i')<<24 | uint32('s')<<16 | uint32('o')<<8 | uint32('m')
	fourCCM4A  = uint32('M')<<24 | uint32('4')<<16 | uint32('A')<<8 | uint32(' ')
	fourCCMA1B = uint32('M')<<24 | uint32('A')<<16 | uint32('1')<<8 | uint32('B')
	fourCCMeta = uint32('m')<<24 | uint32('e')<<16 | uint32('t')<<8 | uint32('a')
	fourCCMiaf = uint32('m')<<24 | uint32('i')<<16 | uint32('a')<<8 | uint32('f')
	fourCCMiAn = uint32('M')<<24 | uint32('i')<<16 | uint32('A')<<8 | uint32('n')
	fourCCMiBr = uint32('M')<<24 | uint32('i')<<16 | uint32('B')<<8 | uint32('r')
	fourCCMif1 = uint32('m')<<24 | uint32('i')<<16 | uint32('f')<<8 | uint32('1')
	fourCCMif2 = uint32('m')<<24 | uint32('i')<<16 | uint32('f')<<8 | uint32('2')
	fourCCMiHA = uint32('M')<<24 | uint32('i')<<16 | uint32('H')<<8 | uint32('A')
	fourCCMiHB = uint32('M')<<24 | uint32('i')<<16 | uint32('H')<<8 | uint32('B')
	fourCCMiHE = uint32('M')<<24 | uint32('i')<<16 | uint32('H')<<8 | uint32('E')
	fourCCMiPr = uint32('M')<<24 | uint32('i')<<16 | uint32('P')<<8 | uint32('r')
	fourCCMp41 = uint32('m')<<24 | uint32('p')<<16 | uint32('4')<<8 | uint32('1')
	fourCCMp42 = uint32('m')<<24 | uint32('p')<<16 | uint32('4')<<8 | uint32('2')
	fourCCMsf1 = uint32('m')<<24 | uint32('s')<<16 | uint32('f')<<8 | uint32('1')
)

func brandFromBufSwitch(buf []byte) brand {
	if len(buf) < 4 {
		return brandUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case fourCCAvci:
		return brandAvci
	case fourCCAvif:
		return brandAvif
	case fourCCCrx:
		return brandCrx
	case fourCCHeic:
		return brandHeic
	case fourCCHeim:
		return brandHeim
	case fourCCHeis:
		return brandHeis
	case fourCCHeix:
		return brandHeix
	case fourCCHevc:
		return brandHevc
	case fourCCHevm:
		return brandHevm
	case fourCCHevs:
		return brandHevs
	case fourCCHevx:
		return brandHevx
	case fourCCIso8:
		return brandIso8
	case fourCCIsom:
		return brandIsom
	case fourCCM4A:
		return brandM4A
	case fourCCMA1B:
		return brandMA1B
	case fourCCMeta:
		return brandMeta
	case fourCCMiaf:
		return brandMiaf
	case fourCCMiAn:
		return brandMiAn
	case fourCCMiBr:
		return brandMiBr
	case fourCCMif1:
		return brandMif1
	case fourCCMif2:
		return brandMif2
	case fourCCMiHA:
		return brandMiHA
	case fourCCMiHB:
		return brandMiHB
	case fourCCMiHE:
		return brandMiHE
	case fourCCMiPr:
		return brandMiPr
	case fourCCMp41:
		return brandMp41
	case fourCCMp42:
		return brandMp42
	case fourCCMsf1:
		return brandMsf1
	default:
		return brandUnknown
	}
}

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

func BenchmarkBrandLookupMapKnown(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsKnown, brandFromBuf)
}

func BenchmarkBrandLookupSwitchKnown(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsKnown, brandFromBufSwitch)
}

func BenchmarkBrandLookupMapMixed(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsMixed, brandFromBuf)
}

func BenchmarkBrandLookupSwitchMixed(b *testing.B) {
	benchmarkBrandLookup(b, benchmarkBrandInputsMixed, brandFromBufSwitch)
}

func BenchmarkParseFTYPBrandsMap(b *testing.B) {
	benchmarkParseFTYPBrands(b, brandFromBuf)
}

func BenchmarkParseFTYPBrandsSwitch(b *testing.B) {
	benchmarkParseFTYPBrands(b, brandFromBufSwitch)
}

//BenchmarkCR3Samples/01_EOS_R_EOS_R_CRAW_ISO_100-2         	  303460	      3776 ns/op	     384 B/op	       8 allocs/op
//BenchmarkCR3Samples/02_EOS_R6_EOS_R6_CRAW_ISO_100_crop_nodual-2         	  300764	      3796 ns/op	     432 B/op	       9 allocs/op
//BenchmarkCR3Samples/03_EOS_R7_443A0157-2                                	  319448	      3799 ns/op	     432 B/op	       9 allocs/op
//BenchmarkCR3Samples/04_EOS_M200_MG_2231-2                               	  298084	      3729 ns/op	     384 B/op	       8 allocs/op
//BenchmarkCR3Samples/05_EOS_M200_MG_2233-2                               	  310292	      3779 ns/op	     384 B/op	       8 allocs/op
//PASS
