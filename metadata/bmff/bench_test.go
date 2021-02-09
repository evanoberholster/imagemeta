package bmff

import (
	"bufio"
	"os"
	"testing"
)

var (
	dir        = "../../../test/img/"
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
	}
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
				_, err := bmr.ReadAndParseBox(TypeFtyp)
				if err != nil {
					b.Fatal(err)
				}
				//hm.setBox(p)

				_, err = bmr.ReadAndParseBox(TypeMeta)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkReadBox100/1         	   50419	     23937 ns/op	    7920 B/op	      72 allocs/op
// BenchmarkReadBox100/2         	    3604	    434948 ns/op	    2096 B/op	      22 allocs/op
// BenchmarkReadBox100/3         	   35198	     31366 ns/op	   13520 B/op	     114 allocs/op
// BenchmarkReadBox100/10        	   55358	     22515 ns/op	    8336 B/op	      75 allocs/op
// BenchmarkReadBox100/d         	   13720	     85422 ns/op	   39249 B/op	     314 allocs/op
// BenchmarkReadBox100/Canon_R6  	   92840	     12604 ns/op	    2752 B/op	      35 allocs/op

// BenchmarkReadBox100/1         	   43686	     24866 ns/op	    7776 B/op	      72 allocs/op
// BenchmarkReadBox100/2         	  128497	      9871 ns/op	    2464 B/op	      35 allocs/op
// BenchmarkReadBox100/3         	   30134	     37833 ns/op	   13376 B/op	     114 allocs/op
// BenchmarkReadBox100/10        	   52178	     25404 ns/op	    8192 B/op	      75 allocs/op
// BenchmarkReadBox100/d         	   12722	     92464 ns/op	   39105 B/op	     314 allocs/op
// BenchmarkReadBox100/Canon_R6  	  115146	     10804 ns/op	    2608 B/op	      35 allocs/op

// BenchmarkReadBox100/1         	   34952	     34002 ns/op	   15680 B/op	      78 allocs/op
// BenchmarkReadBox100/2         	  113184	     12728 ns/op	    3184 B/op	      38 allocs/op
// BenchmarkReadBox100/3         	   23497	     52218 ns/op	   29473 B/op	     121 allocs/op
// BenchmarkReadBox100/10        	   34365	     33384 ns/op	   16096 B/op	      81 allocs/op
// BenchmarkReadBox100/d         	    8767	    137222 ns/op	  104355 B/op	     323 allocs/op
// BenchmarkReadBox100/Canon_R6  	   88756	     13562 ns/op	    4384 B/op	      39 allocs/op

// BenchmarkReadBox100/1         	   34029	     35442 ns/op	   16592 B/op	     128 allocs/op
// BenchmarkReadBox100/2         	    1474	    857551 ns/op	 3655128 B/op	      65 allocs/op
// BenchmarkReadBox100/3         	   19843	     63085 ns/op	   31041 B/op	     212 allocs/op
// BenchmarkReadBox100/10        	   33301	     39005 ns/op	   17040 B/op	     133 allocs/op
// BenchmarkReadBox100/d         	    7278	    157690 ns/op	  109123 B/op	     614 allocs/op
// BenchmarkReadBox100/Canon_R6  	   96159	     13923 ns/op	    4640 B/op	      48 allocs/op
