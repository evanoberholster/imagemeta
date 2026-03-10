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

// BenchmarkParseFormats benchmarks EXIF parsing against representative RAW inputs.
// It includes one sample each for CR2, CR3, GPR, and NEF.
func BenchmarkParseFormats(b *testing.B) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	samples := []benchSample{
		{name: "CR2", glob: "*.CR2"},
		{name: "CR3", glob: "*.CR3"},
		{name: "GPR", glob: "*.GPR"},
		{name: "NEF", glob: "*.NEF"},
		{name: "JPG", glob: "*.jpg"},
		{name: "JXL", glob: "*.jxl"},
		{name: "HEI", glob: "*.heic"},
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
				_, err = Parse(bytes.NewReader(data))
				if err != nil {
					b.Fatalf("parse %s: %v", path, err)
				}
			}
		})
	}
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

//BenchmarkParseFormats/CR2-2         	  147487	      7079 ns/op	3023240.75 MB/s	    8432 B/op	      31 allocs/op
//BenchmarkParseFormats/CR3-2         	   75324	     13956 ns/op	2269424.08 MB/s	   70449 B/op	      37 allocs/op
//BenchmarkParseFormats/GPR-2         	  157712	      7388 ns/op	603337.93 MB/s	    8121 B/op	      29 allocs/op
//BenchmarkParseFormats/NEF-2         	  135216	      8451 ns/op	1618579.51 MB/s	    8168 B/op	      28 allocs/op

// BenchmarkParseFormats/CR2-2         	  171204	      5913 ns/op	3619465.20 MB/s	    7032 B/op	      12 allocs/op
// BenchmarkParseFormats/CR3-2         	   97310	     12160 ns/op	2604718.82 MB/s	   69050 B/op	      19 allocs/op
// BenchmarkParseFormats/GPR-2         	  194019	      6852 ns/op	650491.41 MB/s	    7041 B/op	      13 allocs/op
// BenchmarkParseFormats/NEF-2         	  148620	      7716 ns/op	1772613.86 MB/s	    7032 B/op	      10 allocs/op

// BenchmarkParseFormats/CR2-2         	  154380	      6549 ns/op	3268284.34 MB/s	    6440 B/op	      17 allocs/op
// BenchmarkParseFormats/CR3-2         	   76944	     13071 ns/op	2423077.99 MB/s	   68459 B/op	      24 allocs/op
// BenchmarkParseFormats/GPR-2         	  194202	      6150 ns/op	724749.46 MB/s	    6145 B/op	      13 allocs/op
// BenchmarkParseFormats/NEF-2         	  149605	      7650 ns/op	1788114.05 MB/s	    6184 B/op	      11 allocs/op
// BenchmarkParseFormats/JPG-2         	  373689	      2782 ns/op	434417.57 MB/s	    6045 B/op	       8 allocs/op
// BenchmarkParseFormats/JXL-2         	  275858	      4464 ns/op	85572.98 MB/s	    6171 B/op	      15 allocs/op
// BenchmarkParseFormats/HEI-2         	   21391	     56614 ns/op	10160.06 MB/s	    6107 B/op	      14 allocs/op

//BenchmarkParseFormats/CR2-2         	  139138	      8569 ns/op	2497620.79 MB/s	    7472 B/op	      21 allocs/op
//BenchmarkParseFormats/CR3-2         	   92350	     12444 ns/op	2545128.77 MB/s	   69357 B/op	      24 allocs/op
//BenchmarkParseFormats/GPR-2         	  160053	      6836 ns/op	652073.37 MB/s	    7042 B/op	      13 allocs/op
//BenchmarkParseFormats/NEF-2         	  123918	      9664 ns/op	1415430.10 MB/s	    7097 B/op	      15 allocs/op
//BenchmarkParseFormats/JPG-2         	  471280	      2567 ns/op	470961.42 MB/s	    6944 B/op	       8 allocs/op
//BenchmarkParseFormats/JXL-2         	  303675	      4972 ns/op	76816.52 MB/s	    7070 B/op	      15 allocs/op
//BenchmarkParseFormats/HEI-2         	   21937	     55618 ns/op	10342.08 MB/s	    7008 B/op	      14 allocs/op

//BBenchmarkParseFormats/CR2-2         	  227608	      5670 ns/op	3774956.21 MB/s	     328 B/op	      16 allocs/op
//BBenchmarkParseFormats/CR3-2         	  183457	      5804 ns/op	5456934.55 MB/s	     680 B/op	      18 allocs/op
//BBenchmarkParseFormats/GPR-2         	  293424	      4104 ns/op	1086062.61 MB/s	     192 B/op	       9 allocs/op
//BBenchmarkParseFormats/NEF-2         	  191506	      6008 ns/op	2276591.18 MB/s	     192 B/op	       9 allocs/op
//BBenchmarkParseFormats/JPG-2         	 1237856	       935.5 ns/op	1292134.63 MB/s	      88 B/op	       5 allocs/op
//BBenchmarkParseFormats/JXL-2         	  478900	      2525 ns/op	151257.30 MB/s	     216 B/op	      11 allocs/op
//BBenchmarkParseFormats/HEI-2         	   21470	     54171 ns/op	10618.33 MB/s	     144 B/op	      10 allocs/op
