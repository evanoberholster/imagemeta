package bmff

import (
	"bufio"
	"os"
	"testing"
)

var (
	dir        = "../../test/img/"
	benchmarks = []struct {
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
		//{"iPhone 11", "iPhone11Pro.heic"},
		{"iPhone 12", "iPhone12.heic"},
	} //
)

func BenchmarkReadBox100(b *testing.B) {
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
				br := bufio.NewReader(f)
				b.StartTimer()
				bmr := NewReader(br)
				_, err := bmr.ReadFtypBox()
				if err != nil {
					b.Fatal(err)
				}

				_, err = bmr.ReadMetaBox()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkReadBoxGo100(b *testing.B) {
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
				//a := heif.Open(f)
				//b.StartTimer()
				//
				//_, _ = a.EXIF()
				//if err != nil {
				//	b.Fatal(err)
				//}
			}
		})
	}
}

// GoMedia
// BenchmarkReadBoxGoMedia100/1         	    2144	    603678 ns/op	  479009 B/op	    1229 allocs/op
// BenchmarkReadBoxGoMedia100/2         	   14245	     73283 ns/op	   41185 B/op	     116 allocs/op
// BenchmarkReadBoxGoMedia100/3         	    1886	    784331 ns/op	  737339 B/op	    1919 allocs/op
// BenchmarkReadBoxGoMedia100/10        	    2090	    491453 ns/op	  472849 B/op	    1294 allocs/op
// BenchmarkReadBoxGoMedia100/d         	     594	   2141637 ns/op	 2224648 B/op	    5868 allocs/op
// BenchmarkReadBoxGoMedia100/Canon_R6  	    7159	    199500 ns/op	  178028 B/op	     549 allocs/op
// BenchmarkReadBoxGoMedia100/iPhone_12 	    2265	    506897 ns/op	  480784 B/op	    1332 allocs/op

// Optimized
// BenchmarkReadBox100/1         	   38785	     38055 ns/op	    7648 B/op	      77 allocs/op
// BenchmarkReadBox100/2         	   92409	     12364 ns/op	    2352 B/op	      41 allocs/op
// BenchmarkReadBox100/3         	   28192	     44652 ns/op	   12464 B/op	     117 allocs/op
// BenchmarkReadBox100/10        	   33835	     34452 ns/op	    7744 B/op	      81 allocs/op
// BenchmarkReadBox100/d         	    9351	    232209 ns/op	   38864 B/op	     325 allocs/op
// BenchmarkReadBox100/Canon_R6  	   86648	     18291 ns/op	    2592 B/op	      47 allocs/op
// BenchmarkReadBox100/iPhone_12 	   35826	     32786 ns/op	    8064 B/op	      85 allocs/op

// Latest
// BenchmarkReadBox100/1         	   37662	     29445 ns/op	    7536 B/op	      74 allocs/op
// BenchmarkReadBox100/2         	  109060	     11420 ns/op	    2256 B/op	      38 allocs/op
// BenchmarkReadBox100/3         	   23271	     43007 ns/op	   12368 B/op	     114 allocs/op
// BenchmarkReadBox100/10        	   40869	     29462 ns/op	    7648 B/op	      78 allocs/op
// BenchmarkReadBox100/d         	    8838	    122178 ns/op	   38768 B/op	     322 allocs/op
// BenchmarkReadBox100/Canon_R6  	   82357	     12993 ns/op	    2496 B/op	      44 allocs/op
// BenchmarkReadBox100/iPhone_12 	   42121	     31625 ns/op	    7968 B/op	      82 allocs/op
