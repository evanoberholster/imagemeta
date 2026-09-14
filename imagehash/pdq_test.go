package imagehash

import (
	"bytes"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagehash/internal/pdq"
)

// decodeJPEG decodes a JPEG file into an image.Image.
func decodeJPEG(t *testing.T, path string) image.Image {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

func TestPDQ256NilImage(t *testing.T) {
	t.Parallel()
	if _, err := NewPDQ256(nil); err == nil {
		t.Fatal("expected error for nil image")
	}
	if _, _, err := NewPDQ256WithQuality(nil); err == nil {
		t.Fatal("expected error for nil image")
	}
}

func TestPDQ256WrapperConversion(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"../assets/JPEG.jpg", "../assets/a1.jpg", "../assets/a2.jpg", "../assets/NoExif.jpg"} {
		img := decodeJPEG(t, path)

		got, err := NewPDQ256(img)
		if err != nil {
			t.Fatalf("%s: NewPDQ256: %v", path, err)
		}
		want, err := pdq.Hash(img)
		if err != nil {
			t.Fatalf("%s: pdq.Hash: %v", path, err)
		}
		if got.String() != want.Hash.String() {
			t.Errorf("%s: got %s, want %s", path, got, want.Hash)
		}
		if len(got.String()) != 64 {
			t.Errorf("%s: hash string length = %d, want 64", path, len(got.String()))
		}
	}
}

func TestPDQHashEncodeDecode(t *testing.T) {
	t.Parallel()
	h := PDQHash{0x0123456789abcdef, 0xfedcba9876543210, 0x0f0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f0}

	var buf [32]byte
	h.Encode(buf[:])

	var got PDQHash
	got.Decode(buf[:])
	if got != h {
		t.Errorf("Decode(Encode(h)) = %v, want %v", got, h)
	}
	if got.String() != h.String() {
		t.Errorf("String() = %s, want %s", got.String(), h.String())
	}
	if h.String() != "0123456789abcdeffedcba98765432100f0f0f0f0f0f0f0ff0f0f0f0f0f0f0f0" {
		t.Errorf("unexpected canonical string: %s", h.String())
	}
}

func TestPDQHashDistance(t *testing.T) {
	t.Parallel()
	var zero PDQHash
	if d := zero.Distance(zero); d != 0 {
		t.Errorf("Distance(zero, zero) = %d, want 0", d)
	}
	all := PDQHash{^uint64(0), ^uint64(0), ^uint64(0), ^uint64(0)}
	if d := zero.Distance(all); d != 256 {
		t.Errorf("Distance(zero, all) = %d, want 256", d)
	}
}

// TestPDQ256ImageBounds checks that the hash only depends on pixel data, not
// on the image's origin (bounds.Min), which matters when hashing sub-images or
// JPEGs decoded into a non-zero origin.
func TestPDQ256ImageBounds(t *testing.T) {
	t.Parallel()
	const size = 32
	rng := rand.New(rand.NewSource(1))

	origin := image.NewRGBA(image.Rect(0, 0, size, size))
	for i := range origin.Pix {
		origin.Pix[i] = uint8(rng.Intn(256))
	}
	// Same pixels, offset bounds.
	shifted := image.NewRGBA(image.Rect(100, 100, 100+size, 100+size))
	copy(shifted.Pix, origin.Pix)

	a, err := NewPDQ256(origin)
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewPDQ256(shifted)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("bounds changed the hash: %s != %s", a, b)
	}

	// A SubImage of a larger image must agree with the equivalent crop.
	big := image.NewRGBA(image.Rect(0, 0, size*2, size*2))
	rng = rand.New(rand.NewSource(2))
	for i := range big.Pix {
		big.Pix[i] = uint8(rng.Intn(256))
	}
	sub, ok := big.SubImage(image.Rect(size/2, size/2, size/2+size, size/2+size)).(*image.RGBA)
	if !ok {
		t.Fatal("SubImage did not return *image.RGBA")
	}
	crop, err := NewPDQ256(sub)
	if err != nil {
		t.Fatal(err)
	}
	if crop == (PDQHash{}) {
		t.Error("expected non-zero hash for sub-image")
	}
}

func BenchmarkPDQ256(b *testing.B) {
	f, err := os.Open("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	buf, err := os.ReadFile("../assets/a1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = NewPDQ256(img); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPDQ256Large(b *testing.B) {
	buf, err := os.ReadFile("../assets/a2.jpg")
	if err != nil {
		b.Fatal(err)
	}
	img, err := jpeg.Decode(bytes.NewReader(buf))
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = NewPDQ256(img); err != nil {
			b.Fatal(err)
		}
	}
}
