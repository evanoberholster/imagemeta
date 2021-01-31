package xml

import (
	"os"
	"testing"
)

var (
	dir = "test/samples/CanonEOS7DII.xmp"
)

// BenchmarkXMPRead200 	   28201	     42819 ns/op	    6240 B/op	       2 allocs/op
// BenchmarkXMPRead200 	   33976	     34644 ns/op	    6240 B/op	       2 allocs/op

// Walk
// BenchmarkXMPRead200 	    6062	    191794 ns/op	   23920 B/op	     425 allocs/op

// Read
// BenchmarkXMPRead200 	   63470	     19887 ns/op	    4096 B/op	       1 allocs/op - entire file
// BenchmarkXMPRead200 	   88051	     13714 ns/op	    4096 B/op	       1 allocs/op

// BenchmarkXMPRead200 	   91040	     12859 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkXMPRead200 	   97404	     11858 ns/op	    4128 B/op	       5 allocs/op
// BenchmarkXMPRead200 	  129619	      9809 ns/op	    4096 B/op	       1 allocs/op
// BenchmarkXMPRead200 	  121147	      9268 ns/op	    4112 B/op	       3 allocs/op
// BenchmarkXMPRead200 	  101613	     11894 ns/op	    5024 B/op	      62 allocs/op
// BenchmarkXMPRead200 	  103136	     10140 ns/op	    5024 B/op	      62 allocs/op
// BenchmarkXMPRead200 	   96708	     12042 ns/op	    5024 B/op	      62 allocs/op
func BenchmarkXMPRead200(b *testing.B) {
	f, err := os.Open(dir) //+ "6D.xmp")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f.Seek(0, 0)
		b.StartTimer()
		if _, err := Read(f); err != nil {
			b.Fatal(err)
		}
	}

}

// BenchmarkAttribute100/1         	 5594470	       228 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAttribute100/2         	 5680426	       208 ns/op	       0 B/op	       0 allocs/op

// BenchmarkAttribute100/1         	 5663703	       206 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAttribute100/2         	 6010375	       202 ns/op	       0 B/op	       0 allocs/op
func BenchmarkAttribute100(b *testing.B) {
	tag := Tag{
		raw: []byte("xmpMM:InstanceID=\"xmp.iid:f43404a9-81d4-4ea2-b1ce-e2ecf2b852e6\""),
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tag.raw = []byte("xmpMM:InstanceID=\"xmp.iid:f43404a9-81d4-4ea2-b1ce-e2ecf2b852e6\"")
			b.StartTimer()
			tag.nextAttr()
		}
	})
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tag.raw = []byte("xmpMM:InstanceID=\"xmp.iid:f43404a9-81d4-4ea2-b1ce-e2ecf2b852e6\"")
			b.StartTimer()
			//tag.readAttr2()
		}
	})

}
