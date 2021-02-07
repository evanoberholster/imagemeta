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
		{".JPG/GPS", "17.jpg"},
		{".HEIC", "1.heic"},
		{".AVIF", "image1.avif"},
		{".GoPro/6", "hero6.jpg"},
		{".NEF/Nikon", "2.NEF"},
		{".ARW/Sony", "2.ARW"},
		{".DNG/Adobe", "1.DNG"},
		{".JPG/NoExif", "20.jpg"},
	}
)

func BenchmarkScan200(b *testing.B) {
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
				f.Seek(0, 0)
				if _, err := Scan(f); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkScan200/.CR2/GPS-8         	  368462	      3166 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.CR2/7D-8          	  380390	      3101 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.CR3-8             	  389886	      3103 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.JPG/GPS-8         	  388550	      3047 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.HEIC-8            	  384488	      3126 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.GoPro/6-8         	  397183	      3059 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.NEF/Nikon-8       	  385065	      3175 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.ARW/Sony-8        	  379454	      3594 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.DNG/Adobe-8       	  295646	      3799 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkScan200/.JPG/NoExif-8      	  388678	      3774 ns/op	    4096 B/op	       1 allocs/op
