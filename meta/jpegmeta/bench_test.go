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
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkScan200/1.jpg-8         	  745874	      1562 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/2.jpg-8         	  714488	      1630 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/3.jpg-8         	  558403	      1940 ns/op	    5984 B/op	       4 allocs/op
//BenchmarkScan200/10.jpg-8        	  359642	      3430 ns/op	   28768 B/op	       4 allocs/op
//BenchmarkScan200/13.jpg-8        	 1320680	       908 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/14.jpg-8        	 1290163	       934 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/16.jpg-8        	  327699	      3665 ns/op	   28768 B/op	       4 allocs/op
//BenchmarkScan200/17.jpg-8        	  283201	      4253 ns/op	   31328 B/op	       4 allocs/op
//BenchmarkScan200/20.jpg/NoExif-8 	 1753628	       682 ns/op	    4192 B/op	       2 allocs/op
//BenchmarkScan200/21.jpeg-8       	  826738	      1328 ns/op	   10336 B/op	       4 allocs/op
//BenchmarkScan200/24.jpg-8        	  794427	      1392 ns/op	   10336 B/op	       4 allocs/op
//BenchmarkScan200/123.jpg-8       	 1745244	       687 ns/op	    4192 B/op	       2 allocs/op
