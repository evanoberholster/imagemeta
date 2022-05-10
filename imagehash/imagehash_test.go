package imagehash

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
	"github.com/nfnt/resize"
)

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
	for i := 0; i < b.N; i++ {
		_, err = NewPHashFast(resized)
		if err != nil {
			b.Fatal(err)
		}
	}
}

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

	pixels_new := make([]float64, 4096)
	transforms.Rgb2Gray_new(resized, pixels_new)

	for j := 0; j < len(pixels[0]); j++ {
		if pixels_new[j] != pixels[0][j] {
			t.Errorf("Pixels wanted %0.6f got %0.6f", pixels[0][j], pixels_new[j])
		}
	}

	dct := transforms.DCT2D(pixels, 64, 64)
	transforms.DCT2D_new(pixels_new)

	for i := 0; i < len(dct); i++ {
		if pixels_new[i] != dct[0][i] {
			t.Errorf("DCT wanted %0.6f got %0.6f", dct[0][i], pixels_new[i])
		}
	}
}
