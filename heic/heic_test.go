package heic

import (
	"bufio"
	"os"
	"testing"
)

// TODO: write tests

// TODO: write benchmarks

var (
	dir            = "../../test/img/"
	benchmarksHeic = []struct {
		name     string
		fileName string
	}{
		{"1", "1.heic"},
		{"2", "2.heic"},
		{"3", "3.heic"},
		//{"4", "4.heic"},
		//{"5", "5.heic"},
		//{"6", "6.heic"},
		//{"7", "7.heic"},
		//{"8", "8.heic"},
		//{"9", "9.heic"},
		{"10", "10.heic"},
		{"d", "d.heic"},
		{"Canon R6", "r6.HIF"},
		{"iPhone 11", "iPhone11Pro.heic"},
		{"iPhone 12", "iPhone12.heic"},
	} //
)

func BenchmarkHeicDecodeExif(b *testing.B) {
	for _, bm := range benchmarksHeic {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				f.Seek(0, 0)
				br := bufio.NewReader(f)
				b.StartTimer()
				hm := NewMetadata(br)
				err = hm.GetMeta()
				if err != nil {
					b.Fatal(err)
				}
				_, _ = hm.DecodeExif(f)
			}
		})
	}
}
