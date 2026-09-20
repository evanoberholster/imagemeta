package imagehash

import (
	"image"
	"image/color"
	"testing"
)

// TestPixelPathEquivalence checks that the NRGBA/Gray fast paths agree with
// the RGBA path for identical RGB data.
func TestPixelPathEquivalence(t *testing.T) {
	t.Parallel()
	rgba := image.NewRGBA(image.Rect(0, 0, 8, 8))
	nrgba := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	gray := image.NewGray(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			r, g, b := uint8(x*32), uint8(y*32), uint8((x+y)*16)
			rgba.SetRGBA(x, y, color.RGBA{r, g, b, 0xff})
			nrgba.SetNRGBA(x, y, color.NRGBA{r, g, b, 0xff})
			gray.SetGray(x, y, color.Gray{r})
		}
	}
	want, err := NewAHash(rgba)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := NewAHash(nrgba); err != nil || got != want {
		t.Errorf("NewAHash(NRGBA) = %v, %v; want %v, nil", got, err, want)
	}
	grayRGBA := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			v := gray.GrayAt(x, y).Y
			grayRGBA.SetRGBA(x, y, color.RGBA{v, v, v, 0xff})
		}
	}
	wantGray, err := NewAHash(grayRGBA)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := NewAHash(gray); err != nil || got != wantGray {
		t.Errorf("NewAHash(Gray) = %v, %v; want %v, nil", got, err, wantGray)
	}
}

// TestBlurHashPixelPathEquivalence checks that the NRGBA/Gray fast paths
// agree with the RGBA path for identical RGB data.
func TestBlurHashPixelPathEquivalence(t *testing.T) {
	t.Parallel()
	rgba := image.NewRGBA(image.Rect(0, 0, 64, 64))
	nrgba := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	gray := image.NewGray(image.Rect(0, 0, 64, 64))
	grayRGBA := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			r, g, b := uint8(x*4), uint8(y*4), uint8((x+y)*2)
			rgba.SetRGBA(x, y, color.RGBA{r, g, b, 0xff})
			nrgba.SetNRGBA(x, y, color.NRGBA{r, g, b, 0xff})
			gray.SetGray(x, y, color.Gray{r})
			grayRGBA.SetRGBA(x, y, color.RGBA{r, r, r, 0xff})
		}
	}
	want, err := EncodeBlurHashFast(rgba)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := EncodeBlurHashFast(nrgba); err != nil || got != want {
		t.Errorf("EncodeBlurHashFast(NRGBA) = %q, %v; want %q, nil", got, err, want)
	}
	wantGray, err := EncodeBlurHashFast(grayRGBA)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := EncodeBlurHashFast(gray); err != nil || got != wantGray {
		t.Errorf("EncodeBlurHashFast(Gray) = %q, %v; want %q, nil", got, err, wantGray)
	}
}
