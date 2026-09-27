package exif

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
)

// BenchmarkParseTestImages parses committed fixtures so parse performance is
// gated in CI without download_samples (see benchmark_test.go, which needs
// IMAGEMETA_BENCH_IMAGE_DIR). JPEG.jpg exercises the APP1 TIFF-header search
// path; CR2.exif and Heic.exif are extracted TIFF blobs; Unknown.exif covers
// the miss path and must keep returning ErrNoExif. (NoExif.jpg is excluded:
// its entropy data false-positives the header scan and parses empty.)
func BenchmarkParseTestImages(b *testing.B) {
	samples := []struct {
		name    string
		file    string
		wantErr error
	}{
		{name: "JPEG", file: "../../testImages/JPEG.jpg", wantErr: nil},
		{name: "CR2", file: "../../testImages/CR2.exif", wantErr: nil},
		{name: "HEIC", file: "../../testImages/Heic.exif", wantErr: nil},
		{name: "Unknown", file: "../../testImages/Unknown.exif", wantErr: meta.ErrNoExif},
	}

	for _, sample := range samples {
		data, err := os.ReadFile(sample.file)
		if err != nil {
			b.Fatalf("read %s: %v", sample.file, err)
		}
		if len(data) == 0 {
			b.Fatalf("empty fixture %s", sample.file)
		}

		b.Run(sample.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := Parse(bytes.NewReader(data))
				if sample.wantErr != nil {
					if !errors.Is(err, sample.wantErr) {
						b.Fatalf("parse %s: got %v, want %v", sample.file, err, sample.wantErr)
					}
					continue
				}
				if err != nil {
					b.Fatalf("parse %s: %v", sample.file, err)
				}
			}
		})
	}
}
