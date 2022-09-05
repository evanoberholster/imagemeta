package exif2

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
	f, err := os.Open("../cmd/IMG_3001.jpeg")
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("IMG1.CR3")
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

// BenchmarkExif-12    	   23187	     59066 ns/op	    2754 B/op	      29 allocs/op

// BenchmarkExif-12    	  108328	     10992 ns/op	    3934 B/op	      29 allocs/op

// test.CR2
// BenchmarkExif-12    	  171722	      7073 ns/op	    1319 B/op	      27 allocs/op
