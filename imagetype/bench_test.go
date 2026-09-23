package imagetype

import (
	"bytes"
	"os"
	"testing"
)

// scanBenchRecords mirrors the record layout in test.dat (see gen/main.go):
// consecutive 64-byte headers, one per name, in order.
var scanBenchRecords = []string{
	".CRW",
	".CR2/GPS", ".CR2/7D",
	".CR3",
	".JPG/GPS", ".JPG/NoExif", ".JPG/GoPro",
	".JPEG",
	".HEIC/iPhone", ".HEIC/Conv", ".HEIC/Alt",
	".WEBP",
	".GPR/GoPro",
	".NEF/Nikon",
	".ARW/Sony",
	".DNG/Adobe",
	".PNG",
	".RW2",
	".XMP",
	".PSD",
	".JP2/JPEG2000",
	".BMP",
}

func BenchmarkScan(b *testing.B) {
	data, err := os.ReadFile("test.dat")
	if err != nil {
		b.Fatal(err)
	}
	if len(data) != len(scanBenchRecords)*scanHeaderLength {
		b.Fatalf("test.dat holds %d bytes, want %d records", len(data), len(scanBenchRecords))
	}
	for i, name := range scanBenchRecords {
		rec := data[i*scanHeaderLength : (i+1)*scanHeaderLength]
		b.Run(name, func(b *testing.B) {
			r := bytes.NewReader(rec)
			b.SetBytes(int64(len(rec)))
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				if _, err := r.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				if _, err := Scan(r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
