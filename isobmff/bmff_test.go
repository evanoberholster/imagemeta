package isobmff

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	Logger = log.Level(zerolog.PanicLevel)
}

func BenchmarkCR3(b *testing.B) {
	dir := "../../test/img/"
	f, err := os.Open(dir + "/" + "CanonR6_1.CR3")
	//f, err := os.Open(dir + "/" + "CanonR6_1.HIF")
	//f, err := os.Open(dir + "/" + "iPhone13.heic")
	//f, err := os.Open("../cmd/3.CR3")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	buf, _ := io.ReadAll(f)
	reader := bytes.NewReader(buf)

	r := NewReader(reader)
	defer r.Close()
	for i := 0; i < b.N; i++ {
		if _, err = reader.Seek(0, 0); err != nil {
			b.Fatal(err)
		}
		r.reset(reader)
		if err := r.ReadFTYP(); err != nil {
			b.Fatal(err)
		}
		if err := r.ReadMetadata(); err != nil {
			b.Fatal(err)
		}
		if err := r.ReadMetadata(); err != nil {
			b.Fatal(err)
		}

	}

}

// BenchmarkCR3-12    	  600543	      2001 ns/op	    1168 B/op	       3 allocs/op
// BenchmarkCR3-12    	  420733	      2757 ns/op	    1606 B/op	       3 allocs/op
// BenchmarkCR3-12    	  781185	      1569 ns/op	     931 B/op	       3 allocs/op
// BenchmarkCR3-12    	  680792	      1539 ns/op	     332 B/op	       3 allocs/op
// BenchmarkCR3-12    	  728954	      1635 ns/op	     320 B/op	       3 allocs/op
// BenchmarkCR3-12    	  896595	      1279 ns/op	     143 B/op	       0 allocs/op
// BenchmarkCR3-12    	 1029416	      1160 ns/op	     124 B/op	       0 allocs/op
// BenchmarkCR3-12    	 1456377	       783.9 ns/op	      88 B/op	       0 allocs/op
// BenchmarkFTYP-12    	 9032180	       135.2 ns/op	      14 B/op	       0 allocs/op
