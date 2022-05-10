package imagehash

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
	"github.com/nfnt/resize"
)

//
func BenchmarkPHash(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Phash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = NewPHash(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Phash-Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = NewPHashFast(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("PHash-Fast-Pool", func(b *testing.B) {
		wp := transforms.StartWorkerPool(64)
		defer wp.Close()
		for i := 0; i < b.N; i++ {
			_, err = NewPHashFastPool(resized, wp)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkParralell(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("Phash", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err = NewPHash(resized)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("Phash-Fast", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err = NewPHashFast(resized)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	wp := transforms.StartWorkerPool(64)
	defer wp.Close()
	b.Run("PHash-Fast-Pool", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err = NewPHashFastPool(resized, wp)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

//
func TestPhash(t *testing.T) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		t.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)

	pixelsFast := pixelsPool.Get().(*[]float64)
	defer pixelsPool.Put(pixelsFast)
	transforms.Rgb2GrayFast(resized, pixelsFast)

	for j := 0; j < len(pixels[0]); j++ {
		if (*pixelsFast)[j] != pixels[0][j] {
			t.Errorf("Pixels wanted %0.6f got %0.6f", pixels[0][j], (*pixelsFast)[j])
		}
	}

	dct := transforms.DCT2D(pixels, 64, 64)
	transforms.DCT2DFast(pixelsFast)

	for i := 0; i < len(dct); i++ {
		if (*pixelsFast)[i] != dct[0][i] {
			t.Errorf("DCT wanted %0.6f got %0.6f", dct[0][i], (*pixelsFast)[i])
		}
	}
}
