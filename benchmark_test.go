package exiftool

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

const testFilename = "testImages/ARW.exif"

var (
	dir        = "../test/img/"
	benchmarks = []struct {
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
		//{".JPG/NoExif", "20.jpg"},
	}

	sampleDir     = "samples/"
	testFilenames = []string{
		"honor20.jpg",
		"hero6.jpg",
		"1.CR2",
		"3.CR2",
		"350D.CR2",
		"XT1.CR2",
		"60D.CR2",
		"6D.CR2",
		"7D.CR2",
		"90D.cr3",
		"2.CR3",
		"1.CR3",
		"1.jpg",
		"2.jpg",
		"1.NEF",
		"2.NEF",
		"3.NEF",
		"1.ARW",
		"2.ARW",
		"4.RW2",
		"hero6.gpr",
	}
)

func testOpenFile() *os.File {
	f, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}
	f.Seek(0, 0)
	return f
}

// BenchmarkSearchExifHeader200-8   	  254524	      4106 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkSearchExifHeader200-8   	  245546	      4206 ns/op	    4128 B/op	       2 allocs/op

func BenchmarkSearchImageType200(b *testing.B) {
	f := testOpenFile()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		if _, err := SearchImageType(f); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearchExifHeader100(b *testing.B) {
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.Seek(0, 0)
				if _, err := SearchExifHeader(f); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

//BenchmarkSearchExifHeader100/.CR2/GPS-8         	  365108	      2984 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.CR2/7D-8          	  403268	      2967 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.CR3-8   	  197685	      6220 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.JPG/GPS-8         	  374775	      3255 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.HEIC-8  	    5858	    196891 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.GoPro/6-8         	  372853	      3261 ns/op	    4096 B/op	       1 allocs/op
//BenchmarkSearchExifHeader100/.JPG/NoExif-8      	     475	   2673707 ns/op	    4096 B/op	       1 allocs/op

func BenchmarkSearchExifHeader200(b *testing.B) {
	f := testOpenFile()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		if _, err := SearchExifHeader(f); err != nil {
			//b.Fatal(err)
		}

	}
}

func BenchmarkParseExif100(b *testing.B) {
	for _, bm := range testFilenames {
		b.Run(bm, func(b *testing.B) {
			f, err := os.Open(sampleDir + bm)
			if err != nil {
				panic(err)
			}
			eh, err := SearchExifHeader(f)
			if err != nil {
				b.Fatal(err)
			}
			f.Seek(0, 0)
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.Seek(0, 0)
				_, err = eh.ParseExif(cb)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

//BenchmarkParseExif100/.CR2/GPS-8         	   23035	     54230 ns/op	    9310 B/op	      56 allocs/op
//BenchmarkParseExif100/.CR2/7D-8          	   24091	     48405 ns/op	    8957 B/op	      54 allocs/op
//BenchmarkParseExif100/.CR3-8   	  176527	      6746 ns/op	     901 B/op	      14 allocs/op
//BenchmarkParseExif100/.JPG/GPS-8         	   47270	     25754 ns/op	    5123 B/op	      32 allocs/op
//BenchmarkParseExif100/.HEIC-8  	   50145	     25194 ns/op	    4882 B/op	      29 allocs/op
//BenchmarkParseExif100/.GoPro/6-8         	   54031	     22543 ns/op	    3782 B/op	      28 allocs/op
//BenchmarkParseExif100/.NEF/Nikon-8       	   22287	     54464 ns/op	   12417 B/op	      59 allocs/op
//BenchmarkParseExif100/.ARW/Sony-8        	   28357	     42439 ns/op	    7671 B/op	      53 allocs/op
//BenchmarkParseExif100/.DNG/Adobe-8       	   13603	     87055 ns/op	   18494 B/op	      87 allocs/op

// 2.CR2
//BenchmarkParseExif200-8   	   25021	     47577 ns/op	    9434 B/op	      60 allocs/op
//BenchmarkParseExif200-8   	   25111	     49063 ns/op	    9314 B/op	      56 allocs/op
func BenchmarkParseExif200(b *testing.B) {
	f := testOpenFile()
	b.ReportAllocs()
	b.ResetTimer()
	f.Seek(0, 0)

	eh, err := SearchExifHeader(f)
	if err != nil {
		b.Fatal(err)
	}
	f.Seek(0, 0)

	//cb := bufra.NewBufReaderAt(f, 256*1024)
	buf, _ := ioutil.ReadAll(f)
	cb := bytes.NewReader(buf)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = eh.ParseExif(cb)
		if err != nil {
			b.Fatal(err)
		}
	}
}
