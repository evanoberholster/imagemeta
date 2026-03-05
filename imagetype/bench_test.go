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
		//{".CR3/R6", "canonR6.cr3"},
		{".JPG/GPS", "17.jpg"},
		{".HEIC", "1.heic"},
		//{".HEIC/iPhone11", "iPhone11Pro.heic"},
		//{".HEIC/iPhone12", "iPhone12.heic"},
		//{".AVIF", "image1.avif"},
		{".WEBP", "4.webp"},
		{".GoPro/6", "hero6.jpg"},
		{".NEF/Nikon", "2.NEF"},
		{".ARW/Sony", "2.ARW"},
		{".DNG/Adobe", "1.dng"},
		{".JPG/NoExif", "20.jpg"},
		{".JXL/Exif", "1.jxl"},
	}
)

func BenchmarkScan(b *testing.B) {
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
				if _, err := f.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				if _, err := Scan(f); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

//cpu: AMD Ryzen 9 7950X 16-Core Processor
//BenchmarkScan/.CR2/GPS-2         	  897670	      1146 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.CR2/7D-2          	 1000000	      1140 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.CR3-2             	 1000000	      1144 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.JPG/GPS-2         	 1000000	      1113 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.HEIC-2            	 1000000	      1226 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.WEBP-2            	 1000000	      1093 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.GoPro/6-2         	 1000000	      1112 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.NEF/Nikon-2       	  989836	      1123 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.ARW/Sony-2        	 1000000	      1152 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.DNG/Adobe-2       	 1000000	      1154 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.JPG/NoExif-2      	 1000000	      1113 ns/op	      64 B/op	       1 allocs/op
//BenchmarkScan/.JXL/Exif-2        	 1000000	      1129 ns/op	      64 B/op	       1 allocs/op
