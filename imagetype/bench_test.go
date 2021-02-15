package imagetype

import (
	"os"
	"testing"
)

var (
	dir        = "../../test/img/"
	benchmarks = []struct {
		name     string
		fileName string
	}{
		{".CR2/GPS", "2.CR2"},
		{".CR2/7D", "7D2.CR2"},
		{".CR3", "1.CR3"},
		{".CR3/R6", "canonR6.cr3"},
		{".JPG/GPS", "17.jpg"},
		{".HEIC", "1.heic"},
		{".HEIC/iPhone11", "iPhone11Pro.heic"},
		{".HEIC/iPhone12", "iPhone12.heic"},
		{".AVIF", "image1.avif"},
		{".WEBP", "4.webp"},
		{".GoPro/6", "hero6.jpg"},
		{".NEF/Nikon", "2.NEF"},
		{".ARW/Sony", "2.ARW"},
		{".DNG/Adobe", "1.dng"},
		{".JPG/NoExif", "20.jpg"},
	}
)

func BenchmarkScan100(b *testing.B) {
	for _, bm := range benchmarks {
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
				b.StartTimer()
				if _, err := Scan(f); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkScan100/.CR2/GPS         	 1591376	       754 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.CR2/7D          	 1652364	       769 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.CR3             	 1496485	       761 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.JPG/GPS         	 1622518	       714 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.HEIC            	 1560817	       792 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.AVIF            	 1494444	       802 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.WEBP            	 1555882	       744 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.GoPro/6         	 1697720	       757 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.NEF/Nikon       	 1614928	       749 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.ARW/Sony        	 1508358	       823 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.DNG/Adobe       	 1610572	       728 ns/op	      32 B/op	       1 allocs/op
// BenchmarkScan100/.JPG/NoExif      	 1576730	       718 ns/op	      32 B/op	       1 allocs/op
