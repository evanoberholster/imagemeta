package meta

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var (
	dir            = "../../test/img/"
	benchmarksJPEG = []struct {
		name     string
		fileName string
	}{
		{"1.jpg", "1.jpg"},
		{"2.jpg", "2.jpg"},
		{"3.jpg", "3.jpg"},
		{"10.jpg", "10.jpg"},
		{"13.jpg", "13.jpg"},
		{"14.jpg", "14.jpg"},
		{"16.jpg", "16.jpg"},
		{"17.jpg", "17.jpg"},
		{"20.jpg/NoExif", "20.jpg"},
		{"21.jpeg", "21.jpeg"},
		{"24.jpg", "24.jpg"},
		{"123.jpg", "123.jpg"},
	}
)

func BenchmarkScanJPEG200(b *testing.B) {
	for _, bm := range benchmarksJPEG {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cb.Seek(0, 0)
				br := bufio.NewReader(cb)
				if _, err := ScanJPEG(br); err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

// BenchmarkScanJPEG200/1.jpg-8         	  659306	      1606 ns/op	    5984 B/op	       4 allocs/op
// BenchmarkScanJPEG200/2.jpg-8         	  737552	      1730 ns/op	    5984 B/op	       4 allocs/op
// BenchmarkScanJPEG200/3.jpg-8         	  599175	      2095 ns/op	    5984 B/op	       4 allocs/op
// BenchmarkScanJPEG200/10.jpg-8        	  293434	      3818 ns/op	   28768 B/op	       4 allocs/op
// BenchmarkScanJPEG200/13.jpg-8        	 1306766	       914 ns/op	    4192 B/op	       2 allocs/op
// BenchmarkScanJPEG200/14.jpg-8        	 1274205	       939 ns/op	    4192 B/op	       2 allocs/op
// BenchmarkScanJPEG200/16.jpg-8        	  257604	      4253 ns/op	   28768 B/op	       4 allocs/op
// BenchmarkScanJPEG200/17.jpg-8        	  229893	      5284 ns/op	   31328 B/op	       4 allocs/op
// BenchmarkScanJPEG200/20.jpg/NoExif-8 	 1605872	       746 ns/op	    4192 B/op	       2 allocs/op
// BenchmarkScanJPEG200/21.jpeg-8       	  787837	      1566 ns/op	   10336 B/op	       4 allocs/op
// BenchmarkScanJPEG200/24.jpg-8        	  770265	      1550 ns/op	   10336 B/op	       4 allocs/op
// BenchmarkScanJPEG200/123.jpg-8       	 1633209	       737 ns/op	    4192 B/op	       2 allocs/op

var (
	benchmarksTiff = []struct {
		name     string
		fileName string
	}{
		{".CR2/GPS", "2.CR2"},
		{".CR2/7D", "7D2.CR2"},
		{".CR3", "1.CR3"},
		{".JPG/GPS", "17.jpg"},
		{".HEIC", "1.heic"},
		{".GoPro/6", "hero6.jpg"},
		{".NEF/Nikon", "2.NEF"},
		{".ARW/Sony", "2.ARW"},
		{".DNG/Adobe", "1.DNG"},
		{".JPG/NoExif", "20.jpg"},
	}
)

func BenchmarkScanTiff200(b *testing.B) {
	for _, bm := range benchmarksTiff {
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
				br := bufio.NewReader(f)
				if _, err := ScanTiff(br); err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkScanTiff200/.CR2/GPS-8         	  522862	      2255 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.CR2/7D-8          	  544670	      2238 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.CR3-8             	  258202	      4594 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.JPG/GPS-8         	  505470	      2420 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.HEIC-8            	    7010	    169910 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.GoPro/6-8         	  477766	      2389 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.NEF/Nikon-8       	  521940	      2256 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.ARW/Sony-8        	  541036	      2241 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.DNG/Adobe-8       	  549294	      2314 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkScanTiff200/.JPG/NoExif-8      	     520	   2351142 ns/op	    4096 B/op	       1 allocs/op
