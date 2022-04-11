package bmff

import (
	"bufio"
	"os"
	"testing"
)

var (
	dir = "../../test/img/"
	//dir2       = "../../test/samples/"
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
		{"iPhone 11", "iPhone11Pro.heic"},
		{"iPhone 12", "iPhone12.heic"},
	} //
)

//func parseDir(fn func(f *os.File) error) error {
//	files, err := os.ReadDir(dir2)
//	if err != nil {
//		return err
//	}
//	for _, f := range files {
//		if f.IsDir() {
//			continue
//		}
//		info, err := f.Info()
//		if err != nil {
//			return err
//		}
//		if path.Ext(info.Name()) == ".CR2" || path.Ext(info.Name()) == ".CR3" {
//			continue
//		}
//		f2, err := os.Open(dir2 + info.Name())
//		if err != nil {
//			return err
//		}
//		defer f2.Close()
//		if err = fn(f2); err != nil {
//			return errors.WithMessage(err, info.Name())
//		}
//		fmt.Println(dir2 + info.Name())
//	}
//	return nil
//}

//func TestBMFF(t *testing.T) {
//
//	err := parseDir(func(f *os.File) error {
//		br := bufio.NewReader(f)
//		bmr := NewReader(br)
//		ftyp, err := bmr.ReadFtypBox()
//		if err != nil {
//			return err
//		}
//
//		m, err := bmr.ReadMetaBox()
//		if err != nil {
//			return err
//		}
//		fmt.Println(ftyp, m)
//		return nil
//
//	})
//	fmt.Println(err)
//	t.Error("Hello2")
//}

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
				_, err := bmr.ReadFtypBox()
				if err != nil {
					b.Fatal(err)
				}

				_, err = bmr.ReadMetaBox()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkReadBoxGo100(b *testing.B) {
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
				//a := heif.Open(f)
				//b.StartTimer()
				//
				//_, _ = a.EXIF()
				//if err != nil {
				//	b.Fatal(err)
				//}
			}
		})
	}
}

// GoMedia
// BenchmarkReadBoxGoMedia100/1         	    2144	    603678 ns/op	  479009 B/op	    1229 allocs/op
// BenchmarkReadBoxGoMedia100/2         	   14245	     73283 ns/op	   41185 B/op	     116 allocs/op
// BenchmarkReadBoxGoMedia100/3         	    1886	    784331 ns/op	  737339 B/op	    1919 allocs/op
// BenchmarkReadBoxGoMedia100/10        	    2090	    491453 ns/op	  472849 B/op	    1294 allocs/op
// BenchmarkReadBoxGoMedia100/d         	     594	   2141637 ns/op	 2224648 B/op	    5868 allocs/op
// BenchmarkReadBoxGoMedia100/Canon_R6  	    7159	    199500 ns/op	  178028 B/op	     549 allocs/op
// BenchmarkReadBoxGoMedia100/iPhone_12 	    2265	    506897 ns/op	  480784 B/op	    1332 allocs/op

// Optimized
// BenchmarkReadBox100/1         	   46464	     26109 ns/op	    7248 B/op	      64 allocs/op
// BenchmarkReadBox100/2         	  133093	      9789 ns/op	    1744 B/op	      29 allocs/op
// BenchmarkReadBox100/3         	   26997	     41058 ns/op	   12192 B/op	     104 allocs/op
// BenchmarkReadBox100/10        	   41847	     29207 ns/op	    7344 B/op	      68 allocs/op
// BenchmarkReadBox100/d         	    9414	    118801 ns/op	   38976 B/op	     312 allocs/op
// BenchmarkReadBox100/Canon_R6  	  120914	     11273 ns/op	    1952 B/op	      33 allocs/op
// BenchmarkReadBox100/iPhone_11 	   45565	     27445 ns/op	    7312 B/op	      67 allocs/op
// BenchmarkReadBox100/iPhone_12 	   42534	     29444 ns/op	    7664 B/op	      72 allocs/op

// Latest
// BenchmarkReadBox100/1         	   50270	     25748 ns/op	    7248 B/op	      65 allocs/op
// BenchmarkReadBox100/2         	  134457	      9425 ns/op	    1744 B/op	      30 allocs/op
// BenchmarkReadBox100/3         	   30159	     38505 ns/op	   12192 B/op	     105 allocs/op
// BenchmarkReadBox100/10        	   48042	     26569 ns/op	    7344 B/op	      69 allocs/op
// BenchmarkReadBox100/d         	   10000	    108446 ns/op	   38976 B/op	     313 allocs/op
// BenchmarkReadBox100/Canon_R6  	  104004	     10748 ns/op	    1952 B/op	      34 allocs/op
// BenchmarkReadBox100/iPhone_12 	   39205	     30556 ns/op	    7664 B/op	      73 allocs/op
