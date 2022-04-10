package bmff

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

// BenchmarkCrx10-12    	  345901	      4140 ns/op	    4592 B/op	       2 allocs/op
func BenchmarkCrx10(b *testing.B) {
	f, err := os.Open("../../test/samples/CanonR6_1.CR3")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(buf)

	for i := 0; i < b.N; i++ {
		r.Seek(0, 0)
		b.ReportAllocs()

		bmr := NewReader(r)
		if err != nil {
			b.Fatal(err)
		}

		_, err = bmr.ReadFtypBox()
		if err != nil {
			b.Fatal(err)
		}

		_, err = bmr.ReadCrxMoovBox()
		if err != nil {
			b.Fatal(err)
		}
	}

}
