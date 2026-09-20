package imagehash

import (
	"image"
	"image/color"
	"testing"
)

func TestAhashRoundTrip(t *testing.T) {
	t.Parallel()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	h, err := NewAHash(img)
	if err != nil {
		t.Fatal(err)
	}
	var buf [8]byte
	h.Encode(buf[:])
	var got Ahash
	got.Decode(buf[:])
	if got != h {
		t.Errorf("Decode(Encode(h)) = %v, want %v", got, h)
	}
	parsed, err := ParseAhash(h.String())
	if err != nil {
		t.Fatal(err)
	}
	if parsed != h {
		t.Errorf("ParseAhash(String(h)) = %v, want %v", parsed, h)
	}
	text, err := h.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var fromText Ahash
	if err := fromText.UnmarshalText(text); err != nil {
		t.Fatal(err)
	}
	if fromText != h {
		t.Errorf("UnmarshalText(MarshalText(h)) = %v, want %v", fromText, h)
	}
	if d := h.Distance(h); d != 0 {
		t.Errorf("Distance(h, h) = %d, want 0", d)
	}
}

func TestPHashParseRoundTrip(t *testing.T) {
	t.Parallel()
	h64 := PHash64(0x0123456789abcdef)
	if parsed, err := ParsePHash64(h64.String()); err != nil || parsed != h64 {
		t.Errorf("ParsePHash64(String(h)) = %v, %v; want %v, nil", parsed, err, h64)
	}
	var buf64 [8]byte
	h64.Encode(buf64[:])
	var got64 PHash64
	got64.Decode(buf64[:])
	if got64 != h64 {
		t.Errorf("Decode(Encode(h)) = %v, want %v", got64, h64)
	}
	var text64 PHash64
	if err := text64.UnmarshalText([]byte(h64.String())); err != nil || text64 != h64 {
		t.Errorf("UnmarshalText(String(h)) = %v, %v; want %v, nil", text64, err, h64)
	}

	h256 := PHash256{1, 2, 3, 4}
	if parsed, err := ParsePHash256(h256.String()); err != nil || parsed != h256 {
		t.Errorf("ParsePHash256(String(h)) = %v, %v; want %v, nil", parsed, err, h256)
	}
	var buf256 [32]byte
	h256.Encode(buf256[:])
	var got256 PHash256
	got256.Decode(buf256[:])
	if got256 != h256 {
		t.Errorf("Decode(Encode(h)) = %v, want %v", got256, h256)
	}
}

func TestPDQParseRoundTrip(t *testing.T) {
	t.Parallel()
	h := PDQHash{0x0123456789abcdef, 0xfedcba9876543210, 0x0f0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f0}
	parsed, err := ParsePDQHash(h.String())
	if err != nil {
		t.Fatal(err)
	}
	if parsed != h {
		t.Errorf("ParsePDQHash(String(h)) = %v, want %v", parsed, h)
	}
	var text PDQHash
	textErr := text.UnmarshalText([]byte(h.String()))
	if textErr != nil || text != h {
		t.Errorf("UnmarshalText(String(h)) = %v, %v; want %v, nil", text, textErr, h)
	}
	var buf [32]byte
	h.Encode(buf[:])
	var got PDQHash
	got.Decode(buf[:])
	if got != h {
		t.Errorf("Decode(Encode(h)) = %v, want %v", got, h)
	}
	bts, err := h.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var msg PDQHash
	if _, err := msg.UnmarshalMsg(bts); err != nil {
		t.Fatal(err)
	}
	if msg != h {
		t.Errorf("UnmarshalMsg(MarshalMsg(h)) = %v, want %v", msg, h)
	}
}

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
	got, err := NewAHash(nrgba)
	if err != nil || got != want {
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
	gotGray, err := NewAHash(gray)
	if err != nil || gotGray != wantGray {
		t.Errorf("NewAHash(Gray) = %v, %v; want %v, nil", gotGray, err, wantGray)
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
	got, err := EncodeBlurHashFast(nrgba)
	if err != nil || got != want {
		t.Errorf("EncodeBlurHashFast(NRGBA) = %q, %v; want %q, nil", got, err, want)
	}
	wantGray, err := EncodeBlurHashFast(grayRGBA)
	if err != nil {
		t.Fatal(err)
	}
	gotGray, err := EncodeBlurHashFast(gray)
	if err != nil || gotGray != wantGray {
		t.Errorf("EncodeBlurHashFast(Gray) = %q, %v; want %q, nil", gotGray, err, wantGray)
	}
}

func TestEncodeBlurHashAlias(t *testing.T) {
	t.Parallel()
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	a, err := EncodeBlurHash(img)
	if err != nil {
		t.Fatal(err)
	}
	b, err := EncodeBlurHashFast(img)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("EncodeBlurHash = %q, EncodeBlurHashFast = %q", a, b)
	}
}

func TestEncodeShortPanics(t *testing.T) {
	t.Parallel()
	cases := map[string]func(){
		"PHash64.Encode":  func() { PHash64(1).Encode(nil) },
		"PHash64.Decode":  func() { var h PHash64; h.Decode(nil) },
		"PHash256.Encode": func() { PHash256{}.Encode(make([]byte, 31)) },
		"PHash256.Decode": func() { var h PHash256; h.Decode(make([]byte, 31)) },
		"Ahash.Encode":    func() { Ahash(1).Encode(nil) },
		"Ahash.Decode":    func() { var h Ahash; h.Decode(nil) },
		"PDQ.Encode":      func() { PDQHash{}.Encode(make([]byte, 31)) },
		"PDQ.Decode":      func() { var h PDQHash; h.Decode(make([]byte, 31)) },
	}
	for name, fn := range cases {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic on short buffer", name)
				}
			}()
			fn()
		}()
	}
}
