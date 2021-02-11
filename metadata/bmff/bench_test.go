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
		//{"1", "1.heic"},
		//{"2", "2.heic"},
		//{"3", "3.heic"},
		//{"4", "4.heic"},
		//{"5", "5.heic"},
		//{"6", "6.heic"},
		//{"7", "7.heic"},
		//{"8", "8.heic"},
		//{"9", "9.heic"},
		//{"10", "10.heic"},
		{"d", "d.heic"},
		//{"Canon R6", "r6.HIF"},
		//{"iPhone 11", "iPhone11Pro.heic"},
		//{"iPhone 12", "iPhone12.heic"},
	} //
)

func BenchmarkBoxType100(b *testing.B) {

	buf := []byte{'i', 'i', 'n', 'f'}
	b.Run("BoxType1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = boxType(buf)
		}
	})
}

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

// BenchmarkReadBox100/1         	   28059	     41523 ns/op	   21808 B/op	     190 allocs/op
// BenchmarkReadBox100/2         	   76903	     17394 ns/op	    4256 B/op	      69 allocs/op
// BenchmarkReadBox100/3         	   21050	     59580 ns/op	   40497 B/op	     225 allocs/op
// BenchmarkReadBox100/10        	   28465	     43213 ns/op	   22353 B/op	     200 allocs/op
// BenchmarkReadBox100/d         	    6300	    184694 ns/op	  121668 B/op	     926 allocs/op
// BenchmarkReadBox100/Canon_R6  	   59498	     21389 ns/op	    5536 B/op	      81 allocs/op

// Optimized
// BenchmarkReadBox100/1         	   40927	     31914 ns/op	    7648 B/op	      77 allocs/op
// BenchmarkReadBox100/2         	  101738	     12087 ns/op	    2352 B/op	      41 allocs/op
// BenchmarkReadBox100/3         	   27927	     41614 ns/op	   12464 B/op	     117 allocs/op
// BenchmarkReadBox100/10        	   42190	     30124 ns/op	    7744 B/op	      81 allocs/op
// BenchmarkReadBox100/d         	   10000	    123032 ns/op	   38864 B/op	     325 allocs/op
// BenchmarkReadBox100/Canon_R6  	   99526	     14239 ns/op	    2592 B/op	      47 allocs/op

// Latest
// BenchmarkReadBox100/1         	   35226	     34449 ns/op	    8096 B/op	     126 allocs/op
// BenchmarkReadBox100/2         	  101174	     11631 ns/op	    2416 B/op	      42 allocs/op
// BenchmarkReadBox100/3         	   26181	     44527 ns/op	   12544 B/op	     119 allocs/op
// BenchmarkReadBox100/10        	   32229	     35686 ns/op	    8224 B/op	     133 allocs/op
// BenchmarkReadBox100/d         	    8998	    148732 ns/op	   41232 B/op	     614 allocs/op
// BenchmarkReadBox100/Canon_R6  	   74504	     15766 ns/op	    2752 B/op	      57 allocs/op

//BenchmarkReadBox100/1         	   36826	     35510 ns/op	    8288 B/op	     131 allocs/op
//BenchmarkReadBox100/2         	   94675	     13084 ns/op	    2608 B/op	      47 allocs/op
//BenchmarkReadBox100/3         	   24847	     47115 ns/op	   12736 B/op	     124 allocs/op
//BenchmarkReadBox100/10        	   35782	     40730 ns/op	    8416 B/op	     138 allocs/op
//BenchmarkReadBox100/d         	    6584	    159503 ns/op	   41425 B/op	     619 allocs/op
//BenchmarkReadBox100/Canon_R6  	   80612	     16805 ns/op	    2944 B/op	      62 allocs/op
