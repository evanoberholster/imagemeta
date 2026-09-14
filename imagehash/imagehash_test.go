package imagehash

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func TestHashImageSizeErrors(t *testing.T) {
	t.Parallel()

	if _, err := NewPHash64(nil); !errors.Is(err, ErrImageObject) {
		t.Errorf("NewPHash64(nil) err = %v, want ErrImageObject", err)
	}
	for _, size := range []image.Point{{100, 100}, {64, 128}, {128, 128}, {63, 63}} {
		img := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
		if _, err := NewPHash64(img); !errors.Is(err, ErrImageSize) {
			t.Errorf("NewPHash64(%dx%d) err = %v, want ErrImageSize", size.X, size.Y, err)
		}
		if _, err := NewAHash(img); !errors.Is(err, ErrImageSize) {
			t.Errorf("NewAHash(%dx%d) err = %v, want ErrImageSize", size.X, size.Y, err)
		}
	}
	if _, err := NewPHash256(image.NewRGBA(image.Rect(0, 0, 128, 128))); !errors.Is(err, ErrImageSize) {
		t.Errorf("NewPHash256(128x128) err = %v, want ErrImageSize", err)
	}
	if _, err := NewPHash64(image.NewRGBA(image.Rect(0, 0, 64, 64))); err != nil {
		t.Errorf("NewPHash64(64x64) err = %v", err)
	}
	if _, err := NewAHash(image.NewRGBA(image.Rect(0, 0, 8, 8))); err != nil {
		t.Errorf("NewAHash(8x8) err = %v", err)
	}
	if _, err := EncodeBlurHashFast(image.NewRGBA(image.Rect(0, 0, 64, 128))); !errors.Is(err, ErrBlurHashSize) {
		t.Errorf("EncodeBlurHashFast(64x128) err = %v, want ErrBlurHashSize", err)
	}
}

func TestPHashAltAliases(t *testing.T) {
	f, err := os.Open("../assets/JPEG.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bilinear)

	p, err := NewPHash64(resized)
	if err != nil {
		t.Fatal(err)
	}
	pa, err := NewPHash64Alt(resized)
	if err != nil {
		t.Fatal(err)
	}
	if p != pa {
		t.Errorf("NewPHash64 %v != NewPHash64Alt %v", p, pa)
	}
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
	if bh != "UM6KU]e9MyY5ysiwR*b^Qmi_j=i_IAV@nOsp" {
		t.Error(bh)
	}
}

func TestImageHash(t *testing.T) {
	hashTests := []struct {
		filename string
		phash64  string
		phash256 string
	}{
		{"../assets/a1.jpg", "p:bea4c5c322be8ccc", "p:be49a436c5b6c3fe2292bf068c1ecc0788f470e530f83841fe6dc1cec1869bee"},
		{"../assets/a2.jpg", "p:c3d83c65c5d3962c", "p:c33398cd3c6e651ac59cd2e792332cb18d6cdb32729364cd858db3b948e8cc66"},
		{"../assets/JPEG.jpg", "p:93b3071c583cf4d6", "p:9330933407f21cc758c7389bf438c631586a3b970f1c878c6c9f7ce091cbd6db"},
		{"../assets/NoExif.jpg", "p:94b463946b9c6f94", "p:9496b4b6636c94966b6994976b699496636cb4b66b689c926b49949794976b68"},
	}
	for _, h := range hashTests {
		f, err := os.Open(h.filename)
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
		resized := resize.Resize(256, 256, img, resize.Bilinear)

		p256Alt, err := NewPHash256Alt(resized)
		if err != nil {
			t.Fatal(err)
		}
		p256, err := NewPHash256(resized)
		if err != nil {
			t.Fatal(err)
		}

		resized = resize.Resize(64, 64, img, resize.Bilinear)
		p64Alt, err := NewPHash64Alt(resized)
		if err != nil {
			t.Fatal(err)
		}
		p64, err := NewPHash64(resized)
		if err != nil {
			t.Fatal(err)
		}
		if h.phash256 != p256Alt.String() && h.phash256 != p256.String() {
			t.Errorf("expected \t%s, got \t%s and \t%s", h.phash256, p256Alt, p256)
		}
		if h.phash64 != p64Alt.String() && h.phash64 != p64.String() {
			t.Errorf("expected \t%s, got \t%s and \t%s", h.phash64, p64Alt, p64)
		}
	}

}

func BenchmarkPHash64(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			b.Error(err)
		}
	}()
	buf, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	resized, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	resized = resize.Resize(64, 64, resized, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err = NewPHash64(resized); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FastAlt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err = NewPHash64Alt(resized); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Fast-Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				if _, err = NewPHash64Alt(resized); err != nil {
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
	defer func() {
		if err = f.Close(); err != nil {
			b.Error(err)
		}
	}()
	buf, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	resized, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	resized = resize.Resize(256, 256, resized, resize.Bicubic)
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err = NewPHash256(resized); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FastAlt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err = NewPHash256Alt(resized); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Fast-Parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				if _, err = NewPHash256Alt(resized); err != nil {
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
