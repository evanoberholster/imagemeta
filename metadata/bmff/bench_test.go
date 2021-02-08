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
