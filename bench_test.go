package imagemeta

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
)

var (
	dir = "../test/img/"
)

var (
	benchmarksTiff = []struct {
		name     string
		fileName string
	}{
		//{".CR2/60D", "60D.CR2"},
		{".CR2/GPS", "2.CR2"},
		{".CR2/7D", "7D2.CR2"},
		{".CR3", "1.CR3"},
		//{".CR3/90D", "90D.cr3"},
		//{".CR3/R6", "canonR6.cr3"},
		//{".JPG/GPS", "17.jpg"},
		//{".JPF/GoPro6", "hero6.jpg"},
		//{".HEIC", "1.heic"},
		//{".HEIC/CanonR5", "canonR5.hif"},
		//{".HEIC/CanonR6", "canonR6.hif"},
		//{".HEIC/iPhone11", "iPhone11Pro.heic"},
		//{".HEIC/iPhone12", "iPhone12.heic"},
		//{".AVIF", "image1.avif"},
		//{".NEF/Nikon", "1.NEF"},
		//{".NEF/Nikon", "2.NEF"},
		//{".RW2/Panasonic", "4.RW2"},
		//{".ARW/Sony", "2.ARW"},
		//{".WEBP/Webp", "4.webp"},
		//{".DNG/Adobe", "1.dng"},
		//{".JPG/NoExif", "20.jpg"},
	}
)

func BenchmarkImagemeta100(b *testing.B) {
	for _, bm := range benchmarksTiff {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			buf, err := ioutil.ReadAll(f)
			r := bytes.NewReader(buf)
			defer f.Close()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r.Seek(0, 0)

				exifDecodeFn := func(r io.Reader, m *meta.Metadata) error {
					b.StopTimer()
					b.StartTimer()
					exif.ParseExif(f, m.ExifHeader)
					return nil
				}
				xmpDecodeFn := func(r io.Reader, m *meta.Metadata) error {
					b.StopTimer()
					b.StartTimer()
					return err
				}
				b.StartTimer()
				_, err := NewMetadata(r, xmpDecodeFn, exifDecodeFn)
				if err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkImagemeta100/.JPG/GPS         	  209788	     10444 ns/op	    4352 B/op	       5 allocs/op
//BenchmarkImagemeta100/.HEIC            	   42448	     31201 ns/op	   15728 B/op	      69 allocs/op
//BenchmarkImagemeta100/.HEIC/iPhone11   	   32762	     33636 ns/op	   15792 B/op	      72 allocs/op
//BenchmarkImagemeta100/.HEIC/iPhone12   	   36987	     33406 ns/op	   16144 B/op	      77 allocs/op
//BenchmarkImagemeta100/.GoPro/6         	  325867	      4101 ns/op	    4352 B/op	       5 allocs/op
//BenchmarkImagemeta100/.JPG/NoExif      	  663433	      1898 ns/op	    4288 B/op	       3 allocs/op
//

func BenchmarkScanTiff100(b *testing.B) {
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
				if _, err := tiff.ScanTiffHeader(br, imagetype.ImageTiff); err != nil {

					if err != ErrNoExif {
						b.Error(err)
						//	//b.Fatal(err)
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
