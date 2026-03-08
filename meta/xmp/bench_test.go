package xmp

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var benchmarkPackets = []string{
	"test/acr_sidecar.xmp",
	"test/dng_embedded.xmp",
	"test/lightroom_sidecar.xmp",
}

func benchmarkParseFromFile(b *testing.B, packet string, parseFn func(io.Reader) error) {
	b.Helper()

	data, err := os.ReadFile(packet)
	if err != nil {
		b.Fatalf("read %s: %v", packet, err)
	}

	r := bytes.NewReader(data)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err = r.Seek(0, io.SeekStart); err != nil {
			b.Fatalf("seek %s: %v", packet, err)
		}

		if err := parseFn(r); err != nil && err != io.EOF {
			b.Fatalf("parse %s: %v", packet, err)
		}
	}
}

func BenchmarkParseXmp(b *testing.B) {
	for _, packet := range benchmarkPackets {
		packet := packet
		b.Run(filepath.Base(packet), func(b *testing.B) {
			benchmarkParseFromFile(b, packet, func(r io.Reader) error {
				_, err := ParseXmp(r)
				return err
			})
		})
	}
}

func BenchmarkParseAuto(b *testing.B) {
	for _, packet := range benchmarkPackets {
		packet := packet
		b.Run(filepath.Base(packet), func(b *testing.B) {
			benchmarkParseFromFile(b, packet, func(r io.Reader) error {
				_, err := Parse(r)
				return err
			})
		})
	}
}

// BenchmarkParseXmp/acr_sidecar.xmp-2         	  462160	      2317 ns/op	     112 B/op	       4 allocs/op
// BenchmarkParseXmp/dng_embedded.xmp-2        	  249444	      4751 ns/op	     393 B/op	      20 allocs/op
// BenchmarkParseXmp/lightroom_sidecar.xmp-2   	   44100	     25342 ns/op	    2577 B/op	      76 allocs/op

// BenchmarkParseAuto/acr_sidecar.xmp-2        	  445069	      2551 ns/op	      88 B/op	       3 allocs/op
// BenchmarkParseAuto/dng_embedded.xmp-2       	  242296	      4633 ns/op	     393 B/op	      19 allocs/op
// BenchmarkParseAuto/lightroom_sidecar.xmp-2      50900	     26083 ns/op	    1721 B/op	      59 allocs/op
