package imagehash

import (
	"bytes"
	"image/jpeg"
	"io"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

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

	b.Run("Fast32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//resized, _ = jpeg.Decode(bytes.NewReader(buf))
			_, err = NewPHash64Alt(resized)
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
				_, err = NewPHash64Alt(resized)
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
	if bh != "UM6KU]e9MyY5ysiwR*b^Qmi_j=i_IAV@nOsp" {
		t.Error(bh)
	}
}

func TestImageHash(t *testing.T) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		t.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	resized := resize.Resize(64, 64, img, resize.Bilinear)
	p32, err := NewPHash64Alt(resized)
	if err != nil {
		t.Fatal(err)
	}
	p64, err := NewPHash64(resized)
	if err != nil {
		t.Fatal(err)
	}

	p256Alt, err := NewPHash256Alt(resized)
	if err != nil {
		t.Fatal(err)
	}
	p256, err := NewPHash256(resized)
	if err != nil {
		t.Fatal(err)
	}
	if p256.Distance(p256Alt) != 0 {
		t.Error(p32)
		t.Error(p64)
	}
	if p32.Distance(p64) != 0 {
		t.Error(p32)
		t.Error(p64)
	}
}
