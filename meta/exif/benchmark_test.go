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
		sample := sample

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

// BenchmarkParseFormats/Canon_EOS_6D/CR2-2   	  141033	      8127 ns/op	2633377.29 MB/s	    1544 B/op	      18 allocs/op
// BenchmarkParseFormats/Canon_EOS_R/CR3-2    	   93792	     11411 ns/op	2775657.62 MB/s	    3729 B/op	      20 allocs/op
// BenchmarkParseFormats/HERO6_Black/GPR-2    	  328492	      3794 ns/op	1174922.25 MB/s	     184 B/op	       8 allocs/op
// BenchmarkParseFormats/NIKON_D300S/NEF-2    	  144343	      7871 ns/op	1737803.63 MB/s	     736 B/op	      18 allocs/op
// BenchmarkParseFormats/JPG-1-2              	 1316133	       964.6 ns/op	1253088.14 MB/s	      68 B/op	       3 allocs/op
// BenchmarkParseFormats/Canon_EOS_6D/JPG-2-2 	  420878	      2755 ns/op	1326597.08 MB/s	     208 B/op	       8 allocs/op
// BenchmarkParseFormats/Canon_EOS_R6/JXL-2   	  398570	      2909 ns/op	131304.96 MB/s	     232 B/op	      10 allocs/op
// BenchmarkParseFormats/iPhone_8/HEI-2       	   57021	     20159 ns/op	28533.38 MB/s	     419 B/op	      14 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_6D/CR2-2         	  146234	      7979 ns/op	2682219.44 MB/s	    1368 B/op	      17 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_R/CR3-2          	   94360	     11061 ns/op	2863413.30 MB/s	    1424 B/op	      19 allocs/op
