package imagemeta

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/evanoberholster/imagemeta/meta/tiff"
	"github.com/rs/zerolog"
)

var (
	dir = "../test/img/"
)

var (
	benchmarkCR3 = []struct {
		name     string
		fileName string
	}{
		{".CR3/CanonM6MkII", "CanonM6MkII_1.CR3"},
		{".CR3/CanonM50MkII", "CanonM50MkII_1.CR3"},
		{".CR3/CanonM200", "CanonM200_1.CR3"},
		{".CR3/CanonSL3", "CanonSL3_1.CR3"},
		{".CR3/Canon1DXMkIII", "Canon1DXMkIII.CR3"},
		{".CR3/CanonR", "CanonR_1.CR3"},
		{".CR3/CanonRP", "CanonRP_1.CR3"},
		{".CR3/CanonR3", "CanonR3_1.CR3"},
		{".CR3/CanonR5", "CanonR5_1.CR3"},
		{".CR3/CanonR6", "CanonR6_1.CR3"},
		{".CR3/CanonR7", "CanonR7_1.CR3"},
		{".CR3/CanonR10", "CanonR10_1.CR3"},
		{".CR3/Canon90D", "Canon90D_1.CR3"},
	}
)

func BenchmarkCR3(b *testing.B) {
	for _, bm := range benchmarkCR3 {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				panic(err)
			}
			buf, err := io.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			defer f.Close()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err = r.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				if _, err = DecodeCR3(r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

var (
	benchmarksTiff = []struct {
		name     string
		fileName string
	}{
		{".CR2/6D", "2.CR2"},
		{".CR2/7DMkII", "7D2.CR2"},
		{".JPG/GPS", "17.jpg"},
		{".JPG/GoPro6", "hero6.jpg"},
		{".NEF/D300S", "1.NEF"},
		{".NEF/D700", "2.NEF"},
		{".NEF/D7100", "3.NEF"},
		{".RW2/Panasonic", "4.RW2"},
		{".ARW/Sony", "2.ARW"},
		{".DNG/Adobe", "1.dng"},
		{".JPG/NoExif", "20.jpg"},
	}
)

var (
	benchmarksHeif = []struct {
		name     string
		fileName string
	}{
		{".HEIC", "1.heic"},
		{".HEIC/CanonR5", "CanonR5_1.HIF"},
		{".HEIC/CanonR6", "CanonR6_1.HIF"},
		{".HEIC/iPhone11", "iPhone11.heic"},
		{".HEIC/iPhone12", "iPhone12.heic"},
		{".HEIC/iPhone13", "iPhone13.heic"},
	}
)

var (
	benchmarksJPEG = []struct {
		name     string
		fileName string
	}{
		{".JPEG/1", "1.jpg"},
		{".JPEG/2", "2.jpg"},
		{".JPEG/3", "3.jpg"},
		{".JPEG/4", "4.jpg"},
		{".JPEG/5", "5.jpg"},
		{".JPEG/6", "6.jpg"},
		{".JPEG/7", "7.jpg"},
		{".JPEG/8", "8.jpg"},
		{".JPEG/9", "9.jpg"},
		{".JPEG/10", "10.jpg"},
		{".JPEG/11", "11.jpg"},
		{".JPG/GPS", "17.jpg"},
		{".JPEG/Honor", "honor20.jpg"},
		{".JPG/GoPro6", "hero6.jpg"},
		{".JPG/iPhoneXR", "xr.jpg"},
		{".JPG/NoExif", "20.jpg"},
	}
)

func BenchmarkTiff(b *testing.B) {
	for _, bm := range benchmarksTiff {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := io.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err = r.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				if _, err = DecodeTiff(r); err != nil {
					if err != ErrNoExif {
						b.Error(err)
					}
				}
			}
		})
	}
}

func BenchmarkJPEG(b *testing.B) {
	for _, bm := range benchmarksJPEG {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := io.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err = r.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				_, err = DecodeJPEG(r)
				if err != nil {
					if err != ErrNoExif {
						b.Error(err)
					}
				}
			}
		})
	}
}

func BenchmarkHeif(b *testing.B) {
	for _, bm := range benchmarksHeif {
		b.Run(bm.name, func(b *testing.B) {
			f, err := os.Open(dir + bm.fileName)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			buf, err := io.ReadAll(f)
			if err != nil {
				b.Fatal(err)
			}
			r := bytes.NewReader(buf)

			b.Run("Tiff", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err = r.Seek(0, 0); err != nil {
						b.Fatal(err)
					}
					rr := readerPool.Get().(*bufio.Reader)
					rr.Reset(r)

					ir := exif2.NewIfdReader(zerolog.Logger{})

					it, err := imagetype.ScanBuf(rr)
					if err != nil {
						b.Fatal(err)
					}
					header, err := tiff.ScanTiffHeader(rr, it)
					if err != nil {
						b.Fatal(err)
					}
					if err := ir.DecodeTiff(r, header); err != nil {
						b.Fatal(err)
					}
					ir.Close()
					readerPool.Put(rr)
				}
			})
			b.Run("BMFF", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err = r.Seek(0, 0); err != nil {
						b.Fatal(err)
					}
					ir := exif2.NewIfdReader(zerolog.Logger{})

					br := isobmff.NewReader(r, ir.DecodeIfd, nil, nil)
					if err := br.ReadFTYP(); err != nil {
						panic(err)
					}
					if err := readMetadataAllowBenchmarkEOF(br); err != nil {
						panic(err)
					}
					if err := readMetadataAllowBenchmarkEOF(br); err != nil {
						panic(err)
					}
					ir.Close()
					br.Close()
				}
			})
		})
	}
}

func readMetadataAllowBenchmarkEOF(r *isobmff.Reader) error {
	err := r.ReadMetadata()
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
