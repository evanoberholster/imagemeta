package imagemeta

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/evanoberholster/imagemeta/tiff"
)

var (
	dir            = "../test/img/"
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

func BenchmarkScanJPEG100(b *testing.B) {
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
				b.StopTimer()
				cb.Seek(0, 0)
				br := bufio.NewReader(cb)
				b.StartTimer()
				if _, err := jpeg.ScanJPEG(br, nil, nil); err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

// BenchmarkScanJPEG100/1.jpg         	  251211	      4848 ns/op	   19200 B/op	       2 allocs/op
// BenchmarkScanJPEG100/2.jpg         	  393896	      2940 ns/op	    6144 B/op	       2 allocs/op
// BenchmarkScanJPEG100/3.jpg         	  461436	      2223 ns/op	    1792 B/op	       2 allocs/op
// BenchmarkScanJPEG100/10.jpg        	  255290	      4926 ns/op	   20480 B/op	       2 allocs/op
// BenchmarkScanJPEG100/13.jpg        	 2233306	       567 ns/op	       0 B/op	       0 allocs/op
// BenchmarkScanJPEG100/14.jpg        	 1882832	       593 ns/op	       0 B/op	       0 allocs/op
// BenchmarkScanJPEG100/16.jpg        	  242924	      5573 ns/op	   19200 B/op	       2 allocs/op
// BenchmarkScanJPEG100/17.jpg        	  211317	      6126 ns/op	   23040 B/op	       2 allocs/op
// BenchmarkScanJPEG100/20.jpg/NoExif 	 4142763	       307 ns/op	       0 B/op	       0 allocs/op
// BenchmarkScanJPEG100/21.jpeg       	  556494	      1958 ns/op	    3840 B/op	       2 allocs/op
// BenchmarkScanJPEG100/24.jpg        	  554869	      2195 ns/op	    4096 B/op	       2 allocs/op
// BenchmarkScanJPEG100/123.jpg       	 3695756	       299 ns/op	       0 B/op	       0 allocs/op

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
		{".WEBP/Webp", "4.webp"},
		{".DNG/Adobe", "1.dng"},
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
				b.StopTimer()
				f.Seek(0, 0)
				br := bufio.NewReader(f)
				b.StartTimer()
				if _, err := tiff.ScanTiff(br); err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkScanTiff200/.CR2/GPS         	 1422090	       831 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.CR2/7D          	 1672982	       724 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.CR3             	  300817	      4007 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.JPG/GPS         	 1371778	       924 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.HEIC            	   19898	     67681 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.GoPro/6         	 1362571	       926 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.NEF/Nikon       	 1599162	       758 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.ARW/Sony        	 1687218	       693 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.WEBP/Webp       	     621	   1740838 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.DNG/Adobe       	 1743273	       689 ns/op	       0 B/op	       0 allocs/op
//BenchmarkScanTiff200/.JPG/NoExif      	     398	   3031075 ns/op	       0 B/op	       0 allocs/op
