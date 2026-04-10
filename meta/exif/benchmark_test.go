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
	name string
	glob string
}

func benchmarkParseSamples(b *testing.B, samples []benchSample, opts ...ReaderOption) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	for _, sample := range samples {
		sample := sample

		path, err := firstMatch(filepath.Join(benchDir, sample.glob))
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
		{name: "JPG", glob: "*.jpg"},
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

// firstMatch returns the first path that matches the provided glob.
func firstMatch(pattern string) (string, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", nil
	}
	return paths[0], nil
}

// BenchmarkParseFormats/Canon_EOS_6D/CR2-2   	  146330	      7737 ns/op	2766194.87 MB/s	    1544 B/op	      18 allocs/op
// BenchmarkParseFormats/Canon_EOS_R/CR3-2    	   94710	     11723 ns/op	2701806.39 MB/s	    3731 B/op	      20 allocs/op
// BenchmarkParseFormats/HERO6_Black/GPR-2    	  328575	      3688 ns/op	1208736.22 MB/s	     184 B/op	       8 allocs/op
// BenchmarkParseFormats/NIKON_D300S/NEF-2    	  150327	      7483 ns/op	1827959.84 MB/s	     736 B/op	      18 allocs/op
// BenchmarkParseFormats/JPG-2                	 1323674	       882.9 ns/op	1369075.04 MB/s	      68 B/op	       3 allocs/op
// BenchmarkParseFormats/Canon_EOS_R6/JXL-2   	  479952	      2724 ns/op	140208.68 MB/s	     232 B/op	      10 allocs/op
// BenchmarkParseFormats/iPhone_8/HEI-2       	   51920	     22360 ns/op	25725.17 MB/s	     417 B/op	      14 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_6D/CR2-2         	  137852	      8158 ns/op	2623566.14 MB/s	    1368 B/op	      17 allocs/op
// BenchmarkParseFormatsAFInfoBitsetsOnly/Canon_EOS_R/CR3-2          	  110773	     10201 ns/op	3104768.74 MB/s	    1424 B/op	      19 allocs/op
