package jpegmeta

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var (
	dir        = "../../../test/img/"
	benchmarks = []struct {
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

func BenchmarkScan200(b *testing.B) {
	for _, bm := range benchmarks {
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
				cb.Seek(0, 0)
				if _, err := Scan(cb); err != nil {
					if err != ErrEndOfImage {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkScan200/1.jpg-8         	  800500	      1424 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/2.jpg-8         	  811891	      1475 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/3.jpg-8         	  619838	      1955 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/10.jpg-8        	  341607	      3350 ns/op	   28768 B/op	       4 allocs/op
//BenchmarkScan200/13.jpg-8        	 1422186	       842 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/14.jpg-8        	 1327370	       907 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/16.jpg-8        	  314200	      3740 ns/op	   28768 B/op	       4 allocs/op
//BenchmarkScan200/17.jpg-8        	  280066	      4283 ns/op	   31328 B/op	       4 allocs/op
//BenchmarkScan200/20.jpg/NoExif-8 	 1851289	       651 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/21.jpeg-8       	  833792	      1304 ns/op	   10336 B/op	       4 allocs/op
//BenchmarkScan200/24.jpg-8        	  868519	      1360 ns/op	   10336 B/op	       4 allocs/op
//BenchmarkScan200/123.jpg-8       	 1876675	       643 ns/op	    4192 B/op	       2 allocs/op
