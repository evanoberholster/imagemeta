package imagemeta

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var (
	dir = "../test/img/"
)

var (
	benchmarkCR3 = []struct {
		name     string
		fileName string
	}{
		{".CR3/CanonM6MkII", "CanonM6MkII_1.CR3"},
		{".CR3/CanonM50MkII", "CanonM50MkII_1.CR3"},
		{".CR3/CanonM200", "CanonM200_1.CR3"},
		{".CR3/CanonSL3", "CanonSL3_1.CR3"},
		{".CR3/Canon1DXMkIII", "Canon1DXMkIII.CR3"},
		{".CR3/CanonR", "CanonR_1.CR3"},
		{".CR3/CanonRP", "CanonRP_1.CR3"},
		{".CR3/CanonR3", "CanonR3_1.CR3"},
		{".CR3/CanonR5", "CanonR5_1.CR3"},
		{".CR3/CanonR6", "CanonR6_1.CR3"},
		{".CR3/CanonR7", "CanonR7_1.CR3"},
		{".CR3/CanonR10", "CanonR10_1.CR3"},
		{".CR3/Canon90D", "Canon90D_1.CR3"},
	}
)

func BenchmarkImageMetaCR3(b *testing.B) {
	for _, bm := range benchmarkCR3 {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			defer f.Close()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r.Seek(0, 0)
				if _, err = DecodeCR3(r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

var (
	benchmarksTiff = []struct {
		name     string
		fileName string
	}{
		{".CR2/6D", "2.CR2"},
		{".CR2/7DMkII", "7D2.CR2"},
		//	{".HEIC/CanonR5", "CanonR5_1.HIF"},
		//	{".HEIC/CanonR6", "CanonR6_1.HIF"},
		{".JPG/GPS", "17.jpg"},
		{".JPG/GoPro6", "hero6.jpg"},
		//	{".HEIC", "1.heic"},
		//	{".HEIC/iPhone11", "iPhone11.heic"},
		//	{".HEIC/iPhone12", "iPhone12.heic"},
		//	{".HEIC/iPhone13", "iPhone13.heic"},
		//	//{".AVIF", "image1.avif"},
		{".NEF/D300S", "1.NEF"},
		{".NEF/D700", "2.NEF"},
		{".NEF/D7100", "3.NEF"},
		{".RW2/Panasonic", "4.RW2"},
		{".ARW/Sony", "2.ARW"},
		//{".WEBP/Webp", "4.webp"},
		{".DNG/Adobe", "1.dng"},
		{".JPG/NoExif", "20.jpg"},
	}
)

var (
	benchmarksHeif = []struct {
		name     string
		fileName string
	}{
		{".HEIC", "1.heic"},
		{".HEIC/CanonR5", "CanonR5_1.HIF"},
		{".HEIC/CanonR6", "CanonR6_1.HIF"},
		{".HEIC/iPhone11", "iPhone11.heic"},
		{".HEIC/iPhone12", "iPhone12.heic"},
		{".HEIC/iPhone13", "iPhone13.heic"},
	}
)

var (
	benchmarksJPEG = []struct {
		name     string
		fileName string
	}{

		{".JPEG/1", "1.jpg"},
		{".JPEG/2", "2.jpg"},
		{".JPEG/3", "3.jpg"},
		{".JPEG/4", "4.jpg"},
		{".JPEG/5", "5.jpg"},
		{".JPEG/6", "6.jpg"},
		{".JPEG/7", "7.jpg"},
		{".JPEG/8", "8.jpg"},
		{".JPEG/9", "9.jpg"},
		{".JPEG/10", "10.jpg"},
		{".JPEG/11", "11.jpg"},
		{".JPG/GPS", "17.jpg"},
		{".JPEG/Honor", "honor20.jpg"},
		{".JPG/GoPro6", "hero6.jpg"},
		{".JPG/iPhoneXR", "xr.jpg"},
		{".JPG/NoExif", "20.jpg"},
	}
)

// BenchmarkImageMeta/.CR2/GPS-12         	   78484	     17255 ns/op	   10218 B/op	      22 allocs/op
// BenchmarkImageMeta/.CR2/7D-12          	   53593	     18817 ns/op	   10216 B/op	      21 allocs/op
// BenchmarkImageMeta/.CR3-12             	   87937	     15560 ns/op	    9236 B/op	      21 allocs/op
// BenchmarkImageMeta/.JPG/GPS-12         	  110964	     10347 ns/op	     280 B/op	       4 allocs/op
// BenchmarkImageMeta/.JPG/GoPro6-12      	  164203	      7023 ns/op	     280 B/op	       4 allocs/op
// BenchmarkImageMeta/.NEF/Nikon-12       	   58136	     23389 ns/op	   10241 B/op	      23 allocs/op
// BenchmarkImageMeta/.NEF/Nikon#01-12    	   49773	     23771 ns/op	   10243 B/op	      23 allocs/op
// BenchmarkImageMeta/.RW2/Panasonic-12   	   51008	     20251 ns/op	    4556 B/op	      15 allocs/op

func BenchmarkTiff(b *testing.B) {
	for _, bm := range benchmarksTiff {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r.Seek(0, 0)
				_, err = DecodeTiff(r)
				if err != nil {
					if err != ErrNoExif {
						b.Error(err)
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

func BenchmarkJPEG(b *testing.B) {
	for _, bm := range benchmarksJPEG {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r.Seek(0, 0)
				_, err = DecodeJPEG(r)
				if err != nil {
					if err != ErrNoExif {
						b.Error(err)
					}
				}
			}
		})
	}
}

func BenchmarkHeif(b *testing.B) {
	for _, bm := range benchmarksHeif {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r.Seek(0, 0)
				_, err = DecodeHeif(r)
				if err != nil {
					if err != ErrNoExif {
						b.Error(err)
					}
				}
			}
		})
	}
}
