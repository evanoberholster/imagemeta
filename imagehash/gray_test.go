package imagehash

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/transforms32"
	"github.com/nfnt/resize"
)

func TestGreyPixels(t *testing.T) {
	f, err := os.Open("../assets/JPEG.jpg")
	if err != nil {
		t.Fatal(err)
	}
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

	pixels := pixelsPool32.Get().(*[]float32)
	defer pixelsPool32.Put(pixels)

	pixels2 := pixelsPool32.Get().(*[]float32)
	defer pixelsPool32.Put(pixels2)

	yCbCr := img.(*image.YCbCr)
	count := transforms32.AsmYCbCrToGray8(*pixels,
		yCbCr.Rect.Min.X, yCbCr.Rect.Min.Y, yCbCr.Rect.Max.X, yCbCr.Rect.Max.Y,
		yCbCr.Y, yCbCr.Cb, yCbCr.Cr, yCbCr.YStride, yCbCr.CStride,
	)

	transforms32.PixelYCnCRGray32(yCbCr, *pixels2)

	fmt.Println(count)
	for i := 0; i < len(*pixels); i++ {
		if math.RoundToEven(float64((*pixels)[i])) != math.RoundToEven(float64((*pixels2)[i])) {
			t.Error("error", i)
			t.Error("ASM:\t", (*pixels)[i:i+64])
			t.Error("FN:\t", (*pixels2)[i:i+64])
			t.Error("ASM:\t", (*pixels)[i+64:i+256])
			t.Error("FN:\t", (*pixels2)[i+64:i+256])
			break
		}
	}
	t.Error("ASM:\t", (*pixels)[:32])
	t.Error("FN:\t", (*pixels2)[:32])

}

func BenchmarkGreyPixels(b *testing.B) {
	f, err := os.Open("../assets/JPEG.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	img = resize.Resize(64, 64, img, resize.Bilinear)

	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 64 {
		err = errors.New("error image size incompatible. PHash requires 64x64 image")
		b.Error(err)
		return
	}

	pixels := pixelsPool32.Get().(*[]float32)
	defer pixelsPool32.Put(pixels)

	pixels2 := pixelsPool32.Get().(*[]float32)
	defer pixelsPool32.Put(pixels2)

	yCbCr := img.(*image.YCbCr)

	b.Run("FN", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			transforms32.PixelYCnCRGray32(yCbCr, *pixels2)
		}
	})

	b.Run("ASM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			transforms32.AsmYCbCrToGray8(*pixels,
				yCbCr.Rect.Min.X, yCbCr.Rect.Min.Y, yCbCr.Rect.Max.X, yCbCr.Rect.Max.Y,
				yCbCr.Y, yCbCr.Cb, yCbCr.Cr, yCbCr.YStride, yCbCr.CStride,
			)
		}
	})
}
