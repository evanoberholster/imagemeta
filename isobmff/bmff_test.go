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
	//dir := "../../test/img/"
	//f, err := os.Open(dir + "/" + "CanonR6_1.CR3")
	f, err := os.Open("../cmd/3.CR3")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, _ := io.ReadAll(f)
	reader := bytes.NewReader(buf)

	r := NewReader(reader)
	defer r.Close()
	for i := 0; i < b.N; i++ {
		reader.Seek(0, 0)
		r.reset(reader)
		if err := r.ReadFTYP(); err != nil {
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

func BenchmarkLookup(b *testing.B) {
	//str := mapBoxTypeString[typeCCDT]
	for i := 0; i < b.N; i++ {
		for str, _ := range mapStringBoxType {
			_ = mapStringBoxType[str]
			//_ = cnv2(str)
		}

	}
}

func BenchmarkLookup2(b *testing.B) {
	data := map[uint32]BoxType{}
	for bt, str := range mapBoxTypeString {
		data[cnv(str)] = bt
	}
	//str := mapBoxTypeString[typeCCDT]
	for i := 0; i < b.N; i++ {
		for str, _ := range mapStringBoxType {
			_ = data[cnv(str)]
		}
	}
}

func cnv(str string) uint32 {
	return uint32(str[3]) | uint32(str[2])<<8 | uint32(str[1])<<16 | uint32(str[0])<<24
}

func cnv2(str string) BoxType {
	if str[:4] == "infe" { // inital check for performance reasons
		return typeInfe
	}
	if b, ok := mapStringBoxType[str]; ok {
		return b
	}
	if logLevelError() {
		logErrorMsg("BoxType", "error BoxType '%s' unknown", []byte(str))
	}
	return typeUnknown
}

func TestBuild(t *testing.T) {
	build(mapStringBoxType)
	panic("test")
}
