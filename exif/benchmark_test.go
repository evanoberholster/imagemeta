package exif

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

var (
	sampleDir     = "../../test/img/"
	testFilenames = []string{
		//"honor20.jpg",
		//"hero6.jpg",
		"1.CR2",
		"3.CR2",
		//"350D.CR2",
		//"XT1.CR2",
		"60D.CR2",
		"6D.CR2",
		"7D.CR2",
		"90D.cr3",
		"2.CR3",
		"1.CR3",
		"1.jpg",
		"2.jpg",
		//"1.NEF",
		//"2.NEF",
		//"3.NEF",
		//"1.ARW",
		//"2.ARW",
		//"4.RW2",
		//"hero6.gpr",
		//"4.webp",
		//"20.jpg",
	}
	testFilenames2 = []string{
		"1.CR2",
	}
)

func BenchmarkScanExif100(b *testing.B) {
	for _, bm := range testFilenames {
		b.Run(bm, func(b *testing.B) {
			f, err := os.Open(sampleDir + bm)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = ScanExif(cb)
				if err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkParseExif100/honor20.jpg-8         	   41763	     26779 ns/op	    8791 B/op	      38 allocs/op
//BenchmarkParseExif100/hero6.jpg-8           	   56798	     21478 ns/op	    8036 B/op	      31 allocs/op
//BenchmarkParseExif100/1.CR2-8               	   25717	     43625 ns/op	   13215 B/op	      57 allocs/op
//BenchmarkParseExif100/3.CR2-8               	   31881	     37863 ns/op	   12144 B/op	      53 allocs/op
//BenchmarkParseExif100/350D.CR2-8            	   39760	     30335 ns/op	   10445 B/op	      46 allocs/op
//BenchmarkParseExif100/XT1.CR2-8             	   39951	     30473 ns/op	   10443 B/op	      46 allocs/op
//BenchmarkParseExif100/60D.CR2-8             	   28388	     41794 ns/op	   12592 B/op	      52 allocs/op
//BenchmarkParseExif100/6D.CR2-8              	   26862	     43931 ns/op	   13182 B/op	      57 allocs/op
//BenchmarkParseExif100/7D.CR2-8              	   27742	     43892 ns/op	   13215 B/op	      57 allocs/op
//BenchmarkParseExif100/90D.cr3-8             	  152034	      8243 ns/op	    5158 B/op	      17 allocs/op
//BenchmarkParseExif100/2.CR3-8               	  148568	      8020 ns/op	    5157 B/op	      17 allocs/op
//BenchmarkParseExif100/1.CR3-8               	  152665	      8190 ns/op	    5158 B/op	      17 allocs/op
//BenchmarkParseExif100/1.jpg-8               	   65166	     18525 ns/op	    6755 B/op	      30 allocs/op
//BenchmarkParseExif100/2.jpg-8               	   52927	     22747 ns/op	    8421 B/op	      33 allocs/op
//BenchmarkParseExif100/1.NEF-8               	   23182	     47926 ns/op	   13600 B/op	      61 allocs/op
//BenchmarkParseExif100/2.NEF-8               	   22272	     51129 ns/op	   16674 B/op	      63 allocs/op
//BenchmarkParseExif100/3.NEF-8               	   20660	     58233 ns/op	   17008 B/op	      67 allocs/op
//BenchmarkParseExif100/1.ARW-8               	   31310	     38229 ns/op	   11930 B/op	      56 allocs/op
//BenchmarkParseExif100/2.ARW-8               	   30970	     38146 ns/op	   11932 B/op	      57 allocs/op
//BenchmarkParseExif100/4.RW2-8               	   36390	     32900 ns/op	    8199 B/op	      31 allocs/op
//BenchmarkParseExif100/hero6.gpr-8           	   33552	     35745 ns/op	   13603 B/op	      39 allocs/op
//BenchmarkParseExif100/4.webp-8              	     826	   1536155 ns/op	    4432 B/op	       5 allocs/op
//BenchmarkParseExif100/20.jpg-8              	     472	   2536924 ns/op	    4432 B/op	       5 allocs/op

func BenchmarkParseExif100(b *testing.B) {
	for _, bm := range testFilenames2 {
		b.Run(bm, func(b *testing.B) {
			f, err := os.Open(sampleDir + bm)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = ScanExif(cb)
				if err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

// BenchmarkParseExif100/1.CR2         	   31756	     34805 ns/op	   13706 B/op	      38 allocs/op

// BenchmarkParseExif100/1.CR2         	   32517	     37431 ns/op	   13706 B/op	      38 allocs/op

// BenchmarkParseExif100/1.CR2         	   34333	     35437 ns/op	    9148 B/op	      56 allocs/op

// BenchmarkParseExif100/1.CR2         	   31680	     39503 ns/op	   13214 B/op	      57 allocs/op

func BenchmarkIfdKey(b *testing.B) {
	tm := make(ifds.TagMap)
	tagTest := tag.NewTag(ifds.ActiveArea, tag.TypeASCII, 32, 32, tag.RawValueOffset([4]byte{1, 1, 1, 1}))
	b.Run("Key/Set", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := ifds.NewKey(ifds.RootIFD, 0, tag.ID(uint8(i%255)))
			tm[key] = tagTest
		}
	})
	b.Run("Key/Get", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := ifds.NewKey(ifds.RootIFD, 0, tag.ID(uint8(i%255)))
			_ = key
			_ = tm[key]
		}
	})
}
