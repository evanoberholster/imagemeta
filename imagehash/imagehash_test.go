package imagehash

import (
	"bytes"
	"image/jpeg"
	"io"
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
	defer f.Close()
	buf, _ := io.ReadAll(f)
	resized, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	resized = resize.Resize(64, 64, resized, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()

	//b.Run("Regular", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		_, err = NewPHash(resized)
	//		if err != nil {
	//			b.Fatal(err)
	//		}
	//	}
	//})

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//resized, _ = jpeg.Decode(bytes.NewReader(buf))
			_, err = NewPHash64(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	//b.Run("Parallel", func(b *testing.B) {
	//	b.RunParallel(func(p *testing.PB) {
	//		for p.Next() {
	//			resized, _ = jpeg.Decode(bytes.NewReader(buf))
	//			_, err = NewPHash(resized)
	//			if err != nil {
	//				b.Fatal(err)
	//			}
	//		}
	//	})
	//})

	b.Run("Fast-Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				//resized, _ = jpeg.Decode(bytes.NewReader(buf))
				_, err = NewPHash64(resized)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

}

func BenchmarkPHash256(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	buf, _ := io.ReadAll(f)
	resized, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	resized = resize.Resize(256, 256, resized, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()

	//b.Run("Regular", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		_, err = NewPHash(resized)
	//		if err != nil {
	//			b.Fatal(err)
	//		}
	//	}
	//})

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//resized, _ = jpeg.Decode(bytes.NewReader(buf))
			_, err = NewPHash256(resized)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	//b.Run("Parallel", func(b *testing.B) {
	//	b.RunParallel(func(p *testing.PB) {
	//		for p.Next() {
	//			resized, _ = jpeg.Decode(bytes.NewReader(buf))
	//			_, err = NewPHash(resized)
	//			if err != nil {
	//				b.Fatal(err)
	//			}
	//		}
	//	})
	//})

	b.Run("Fast-Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				//resized, _ = jpeg.Decode(bytes.NewReader(buf))
				_, err = NewPHash256(resized)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

}

func BenchmarkBlurHash100(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}

	resized := resize.Resize(64, 64, img, resize.Bilinear)
	b.Run("BlurHash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bh, err := EncodeBlurHashFast(resized)
			if err != nil {
				b.Fatal(err)
			}
			_ = bh
		}
	})
	//b.Run("BlurHashFast", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		bh, err := EncodeBlurHashFast(resized)
	//		if err != nil {
	//			b.Fatal(err)
	//		}
	//		_ = bh
	//	}
	//})

}

func TestBlurHash(t *testing.T) {
	f, err := os.Open("../assets/JPEG.jpg")
	if err != nil {
		t.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bilinear)
	bh, err := EncodeBlurHashFast(resized)
	if err != nil {
		t.Fatal(err)
	}
	//t.Error(bh)
	_ = bh

	//  UcE:P7s;$-xt~qkCt9WV%3t7ayRjogs;RjWA
	//  UcE:P7s;$-xt~qkCt9WV%3t7ayRjogs;RjWAFAIL
}

// p:bea4c5c322be8ccc p:d8b3b0beb9112bc1
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

	pixelsFast := pixelsPool64.Get().(*[]float64)
	defer pixelsPool64.Put(pixelsFast)
	transforms.Rgb2GrayFast(resized, pixelsFast)

	resized2 := resize.Resize(256, 256, img, resize.Bilinear)
	p1, _ := NewPHash(resized)
	p2, _ := NewPHash64(resized)
	p3, _ := NewPHash(resized2)
	p4, _ := NewPHash256(resized2)

	t.Error(p1, p2)
	t.Error(p3, p4)
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
