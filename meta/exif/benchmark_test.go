package exif

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const defaultBenchImageDir = "/home/evanoberholster/go/src/github.com/evanoberholster/test/img"

type benchSample struct {
	name     string
	glob     string
	matchIdx int
}

func benchmarkParseSamples(b *testing.B, samples []benchSample, opts ...ReaderOption) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	for _, sample := range samples {

		path, err := nthMatch(filepath.Join(benchDir, sample.glob), sample.matchIdx)
		if err != nil {
			b.Fatalf("glob %q: %v", sample.glob, err)
		}
		if path == "" {
			b.Skipf("no sample found for %s in %s", sample.glob, benchDir)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			b.Fatalf("read %s: %v", path, err)
		}

		benchName := sample.name
		if parsed, parseErr := Parse(bytes.NewReader(data)); parseErr == nil {
			if model := strings.TrimSpace(parsed.IFD0.Model); model != "" {
				benchName = model + "/" + sample.name
			}
		}

		b.Run(benchName, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if len(opts) == 0 {
					_, err = Parse(bytes.NewReader(data))
				} else {
					_, err = ParseWithReaderOptions(bytes.NewReader(data), opts...)
				}
				if err != nil {
					b.Fatalf("parse %s: %v", path, err)
				}
			}
		})
	}
}

// BenchmarkParseFormats benchmarks EXIF parsing against representative RAW inputs.
func BenchmarkParseFormats(b *testing.B) {
	samples := []benchSample{
		{name: "CR2", glob: "*.CR2"},
		{name: "CR3", glob: "*.CR3"},
		{name: "GPR", glob: "*.GPR"},
		{name: "NEF", glob: "*.NEF"},
		{name: "JPG-1", glob: "*.jpg"},
		{name: "JPG-2", glob: "*.jpg", matchIdx: 1},
		{name: "JXL", glob: "*.jxl"},
		{name: "HEI", glob: "*.heic"},
		{name: "ARW", glob: "*.ARW"},
	}

	benchmarkParseSamples(b, samples)
}

// BenchmarkParseFormatsAFInfoBitsetsOnly benchmarks CR2/CR3 parsing while only
// decoding Canon AFInfo in-focus/selected bitsets.
func BenchmarkParseFormatsAFInfoBitsetsOnly(b *testing.B) {
	samples := []benchSample{
		{name: "CR2", glob: "*.CR2"},
		{name: "CR3", glob: "*.CR3"},
	}
	opts := WithAFInfoDecodeOptions(AFInfoDecodeInFocus | AFInfoDecodeSelected)

	benchmarkParseSamples(b, samples, opts)
}

// nthMatch returns the path at idx from the sorted paths matching the provided glob.
func nthMatch(pattern string, idx int) (string, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if idx < 0 || idx >= len(paths) {
		return "", nil
	}
	return paths[idx], nil
}

// BenchmarkParseFormats/Canon_EOS_6D/CR2-2   	  118957	      9075 ns/op	2358490.69 MB/s	    1592 B/op	      20 allocs/op
// BenchmarkParseFormats/Canon_EOS_R/CR3-2    	   93140	     12421 ns/op	2549865.16 MB/s	    3704 B/op	      21 allocs/op
// BenchmarkParseFormats/HERO6_Black/GPR-2    	  268288	      4212 ns/op	1058341.14 MB/s	     112 B/op	       4 allocs/op
// BenchmarkParseFormats/NIKON_D300S/NEF-2    	  133698	      8673 ns/op	1577125.12 MB/s	     752 B/op	      21 allocs/op
// BenchmarkParseFormats/JPG-1-2              	 1143925	       993.7 ns/op	1216424.78 MB/s	      64 B/op	       2 allocs/op
// BenchmarkParseFormats/Canon_EOS_6D/JPG-2-2 	  414728	      2884 ns/op	1267275.98 MB/s	     176 B/op	       8 allocs/op
// BenchmarkParseFormats/Canon_EOS_R6/JXL-2   	  374104	      3468 ns/op	110154.00 MB/s	     200 B/op	      10 allocs/op
// BenchmarkParseFormats/iPhone_8/HEI-2       	   58581	     20369 ns/op	28239.76 MB/s	     833 B/op	      16 allocs/op
// BenchmarkParseFormats/SLT-A55V/ARW-2       	  163225	      7035 ns/op	2459388.31 MB/s	    1472 B/op	      10 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_6D/CR2-2         	  116385	      8856 ns/op	2416727.39 MB/s	    1416 B/op	      19 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_R/CR3-2          	  110014	     10455 ns/op	3029547.42 MB/s	    1400 B/op	      20 allocs/op
