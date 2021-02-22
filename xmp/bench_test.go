package xmp

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var (
	dir = "test/jpeg.xmp"
)

// BenchmarkXMPRead200 	   27654	     45272 ns/op	    6784 B/op	      16 allocs/op
// BenchmarkXMPRead200 	   28819	     42524 ns/op	    6368 B/op	       8 allocs/op
// BenchmarkXMPRead200 	   28201	     42819 ns/op	    6240 B/op	       2 allocs/op
// BenchmarkXMPRead200 	   33976	     34644 ns/op	    6240 B/op	       2 allocs/op

// 7D MKII
// BenchmarkXMPRead200 	   22694	     52323 ns/op	    7248 B/op	      29 allocs/op
// BenchmarkXMPRead200 	   23542	     50447 ns/op	    7248 B/op	      29 allocs/op

// BenchmarkXMPRead 	   47311	     26398 ns/op	    2304 B/op	      44 allocs/op
func BenchmarkXMPRead(b *testing.B) {
	f, err := os.Open(dir) //+ "6D.xmp")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	a, _ := ioutil.ReadAll(f)
	r2 := bytes.NewReader(a)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r2.Seek(0, 0)
		b.StartTimer()

		if _, err := ParseXmp(r2); err != nil {
			if err != io.EOF {
				b.Fatal(err)
			}
		}
	}

}
