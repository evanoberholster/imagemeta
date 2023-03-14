//go:build linux && amd64

package transforms32

import (
	"errors"
	"image"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func TestGreyPixels(t *testing.T) {
	f, err := os.Open("../../assets/JPEG.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Error(err)
		}
	}()
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	img = resize.Resize(64, 64, img, resize.Bilinear)

	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 64 {
		err = errors.New("error image size incompatible. PHash requires 64x64 image")
		t.Error(err)
		return
	}
	pixels := make([]float32, 64*64)
	pixels2 := make([]float32, 64*64)
	for i := 0; i < len(pixels); i++ {
		pixels[i] = rand.Float32()
	}
	copy(pixels2, pixels)

	yCbCr := img.(*image.YCbCr)
	AsmYCbCrToGray(yCbCr, pixels)

	yCbCrToGrayAlt(yCbCr, pixels2)

	for i := 0; i < len(pixels); i++ {
		if math.Abs(float64((pixels2)[i])-float64((pixels)[i])) > 2.0 {
			t.Error("error", i)
			t.Error("ASM:\t", (pixels)[i:i+64])
			t.Error("FN:\t", (pixels2)[i:i+64])
			break
		}
	}
}

func BenchmarkGreyPixels(b *testing.B) {
	f, err := os.Open("../../assets/JPEG.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	img = resize.Resize(64, 64, img, resize.Bilinear)
	yCbCr := img.(*image.YCbCr)

	pixels := make([]float32, 64*64)
	pixels2 := make([]float32, 64*64)
	for i := 0; i < len(pixels); i++ {
		pixels[i] = rand.Float32()
	}
	copy(pixels2, pixels)

	b.Run("FN", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			yCbCrToGrayAlt(yCbCr, pixels2)
		}
	})

	b.Run("ASM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AsmYCbCrToGray(yCbCr, pixels2)
		}
	})
}
