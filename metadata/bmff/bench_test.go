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

// BenchmarkReadBox100/1         	  622712	      1926 ns/op	     176 B/op	       3 allocs/op
// BenchmarkReadBox100/2         	  607876	      1969 ns/op	     176 B/op	       3 allocs/op
// BenchmarkReadBox100/3         	  653284	      1924 ns/op	     176 B/op	       3 allocs/op
// BenchmarkReadBox100/10        	  535334	      2050 ns/op	     176 B/op	       3 allocs/op
// BenchmarkReadBox100/d         	  688734	      1883 ns/op	     176 B/op	       3 allocs/op
// BenchmarkReadBox100/Canon_R6  	  678283	      1743 ns/op	     176 B/op	       3 allocs/op

// BenchmarkReadBox100/1         	   50086	     23542 ns/op	   11408 B/op	      70 allocs/op
// BenchmarkReadBox100/2         	  289993	      4054 ns/op	    1200 B/op	      13 allocs/op
// BenchmarkReadBox100/3         	   37408	     32088 ns/op	   18704 B/op	     112 allocs/op
// BenchmarkReadBox100/10        	   54189	     24448 ns/op	   11632 B/op	      73 allocs/op
// BenchmarkReadBox100/d         	   12717	     90201 ns/op	   56978 B/op	     312 allocs/op
// BenchmarkReadBox100/Canon_R6  	  391518	      3842 ns/op	    1104 B/op	      12 allocs/op

// BenchmarkReadBox100/1         	   48442	     22721 ns/op	    7920 B/op	      66 allocs/op
// BenchmarkReadBox100/2         	  263370	      4366 ns/op	    1136 B/op	      11 allocs/op
// BenchmarkReadBox100/3         	   39302	     35117 ns/op	   13520 B/op	     108 allocs/op
// BenchmarkReadBox100/10        	   52827	     22508 ns/op	    8336 B/op	      69 allocs/op
// BenchmarkReadBox100/d         	   13585	     89125 ns/op	   39249 B/op	     308 allocs/op
// BenchmarkReadBox100/Canon_R6  	  266624	      3997 ns/op	    1040 B/op	      10 allocs/op

//BenchmarkReadBox100/1         	   47042	     24266 ns/op	    7920 B/op	      72 allocs/op
//BenchmarkReadBox100/2         	    7963	    145803 ns/op	    2096 B/op	      22 allocs/op
//BenchmarkReadBox100/3         	   34891	     35104 ns/op	   13520 B/op	     114 allocs/op
//BenchmarkReadBox100/10        	   47623	     22224 ns/op	    8336 B/op	      75 allocs/op
//BenchmarkReadBox100/d         	   10000	    104304 ns/op	   39249 B/op	     314 allocs/op
//BenchmarkReadBox100/Canon_R6  	  105330	     10634 ns/op	    2752 B/op	      35 allocs/op
