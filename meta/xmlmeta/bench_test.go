package xmlmeta

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"testing"

	"github.com/evanoberholster/exiftool/meta/jpegmeta"
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
		{"test.xmp", "test.xmp"},
	}
)

func BenchmarkXMP200(b *testing.B) {
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			m, err := jpegmeta.Scan(f)
			if err != nil {
				b.Fatal(err)
			}
			str := m.XML()
			//str = xmlfmt.FormatXML(m.XML(), "\t", "  ")
			//str := strings.Replace(m.XML(), "\n", "", -1)
			//str = strings.Replace(str, "   ", "", -1)

			rXML := bytes.NewReader([]byte(str))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rXML.Seek(0, 0)
				_, err = Walk(rXML)
				if err != nil {
					if err != io.EOF {
						b.Fatal(err)
					}
				}
			}
		})
		b.Run(bm.name+"/Walk", func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			m, err := jpegmeta.Scan(f)
			if err != nil {
				b.Fatal(err)
			}
			//rXML := bytes.NewReader()
			buf := []byte(m.XML())
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var decoded XMPPacket2
				err = xml.Unmarshal(buf, &decoded)
				if err != nil {
					if err != io.EOF {
						b.Fatal(err)
					}
				}
			}
		})
	}
}
