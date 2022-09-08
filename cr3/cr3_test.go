package cr3

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func BenchmarkExif(b *testing.B) {
	//f, err := os.Open("../cmd/test.CR2")
	//f, err := os.Open("../testImages/CR2.exif")
	//f, err := os.Open("../../test/img/14.JPG")
	//f, err := os.Open("../cmd/IMG_3001.jpeg")
	//f, err := os.Open("../testImages/Heic.exif")
	f, err := os.Open("../cmd/1.CR3")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	buf, _ := io.ReadAll(f)
	r := bytes.NewReader(buf)

	for i := 0; i < b.N; i++ {
		r.Seek(0, 0)
		if _, err = Decode(r); err != nil {
			b.Fatal(err)
		}
	}
}

// 1.CR3
// BenchmarkExif-12    	  133880	      8810 ns/op	    5481 B/op	       8 allocs/op

// BenchmarkExif-12    	   94231	     10635 ns/op	    6850 B/op	      15 allocs/op

// 2.CR3
// BenchmarkExif-12    	  146803	      7716 ns/op	    5977 B/op	       7 allocs/op

// 3.CR3
// BenchmarkExif-12    	  137854	      7947 ns/op	    8727 B/op	       7 allocs/op
