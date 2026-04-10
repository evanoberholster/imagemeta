package exif

import (
	"bytes"
	"os"
	"path/filepath"
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
		b.Run(sample.name, func(b *testing.B) {
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
