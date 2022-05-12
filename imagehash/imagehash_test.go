package imagehash

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
	"github.com/nfnt/resize"
)

//
func BenchmarkPHash64(b *testing.B) {
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
	b.Run("Regular", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = NewPHash(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = NewPHashFast(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err = NewPHash(resized)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("Fast-Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err = NewPHashFast(resized)
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

	p1, _ := NewPHash(resized)
	p2, _ := NewPHashFast(resized)
	if p1 != p2 {
		t.Errorf("PHash should equal PHashFast, wanted %v, got %v", p2, p1)
		for j := 0; j < len(pixels[0]); j++ {
			if (*pixelsFast)[j] != pixels[0][j] {
				t.Errorf("Pixels wanted %0.6f got %0.6f", pixels[0][j], (*pixelsFast)[j])
			}
		}

		dct := transforms.DCT2D(pixels, 64, 64)
		transforms.DCT2DFast(pixelsFast)

		for j := 0; j < len(dct); j++ {
			for i := 0; i < len(dct); i++ {
				if (*pixelsFast)[i+j*len(dct)] != dct[j][i] {
					t.Errorf("DCT wanted %0.8f got %0.8f", dct[j][i], (*pixelsFast)[i+j*len(dct)])
				}
			}
		}
	}

}
