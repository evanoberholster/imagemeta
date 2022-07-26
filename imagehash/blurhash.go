package imagehash

import (
	"errors"
	"image"
	"math"
	"unicode/utf8"

	"github.com/evanoberholster/imagemeta/imagehash/transforms"
)

const (
	xComponents   = 4
	yComponents   = 4
	width, height = 64, 64
)

func init() {
	initLinearTable(channelToLinear[:])
	initStaticBlurHashValues()
}

func EncodeBlurHashFast(img image.Image) (string, error) {
	if xComponents < 1 || xComponents > 9 {
		return "", errors.New("error invalid number of x components")
	}
	if yComponents < 1 || yComponents > 9 {
		return "", errors.New("error invalid number of y components")
	}
	if img.Bounds().Max.Y != height && img.Bounds().Max.X != width {
		return "", errors.New("error invalid image size")
	}

	b := newBlur(4 + 2*xComponents*yComponents)

	// Size Flag
	b.encode((xComponents-1)+(yComponents-1)*9, 1)

	factors := [xComponents * yComponents * 3]float64{}
	//factors := make([]float64, y*x*3)
	multiplyBasisFunction(img, factors[:])

	var maximumValue float64
	var quantisedMaximumValue int
	var acCount = xComponents*yComponents - 1
	if acCount > 0 {
		var actualMaximumValue float64
		for i := 0; i < acCount*3; i++ {
			actualMaximumValue = math.Max(math.Abs(factors[i+3]), actualMaximumValue)
		}
		quantisedMaximumValue = int(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = (float64(quantisedMaximumValue) + 1) / 166
	} else {
		maximumValue = 1
	}

	// Quantised max AC component
	b.encode(quantisedMaximumValue, 1)

	// DC value
	b.encode(encodeDC(factors[0], factors[1], factors[2]), 4)

	// AC values
	for i := 0; i < acCount; i++ {
		b.encode(encodeAC(factors[3+(i*3+0)], factors[3+(i*3+1)], factors[3+(i*3+2)], maximumValue), 2)
	}

	return b.String(), nil
}

func multiplyBasisFunction(img image.Image, factors []float64) {
	switch c := img.(type) {
	case *image.YCbCr:
		factorsYCbCR(c, factors)
	case *image.RGBA:
		factorsRGBA(c, factors)
	default:
		factorsDefault(c, factors)
	}
}

var (
	channelToLinear [256]float64
	xvalues         = [xComponents * width]float64{}
	yvalues         = [yComponents * height]float64{}
)

func initLinearTable(table []float64) {
	for i := range table {
		channelToLinear[i] = transforms.SRGBToLinear(i)
	}
}

func initStaticBlurHashValues() {
	for xc := 0; xc < xComponents; xc++ {
		for x := 0; x < width; x++ {
			xvalues[x+width*xc] = math.Cos(math.Pi * float64(xc) * float64(x) / float64(width))
		}
	}

	for yc := 0; yc < yComponents; yc++ {
		for y := 0; y < height; y++ {
			yvalues[y+height*yc] = math.Cos(math.Pi * float64(yc) * float64(y) / float64(height))
		}
	}
}

// blur is a blurhash base83 encoder
type blur struct {
	b [4 + 2*xComponents*yComponents]byte
	p int
}

func newBlur(size int) *blur {
	return &blur{}
}

func (b blur) String() string {
	return string(b.b[:])
}

const (
	characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"
)

func (b *blur) encode(value, length int) {
	divisor := int(math.Pow(83, float64(length))) / 83
	for i := 0; i < length; i++ {
		b.p += utf8.EncodeRune(b.b[b.p:], rune(characters[(value/divisor)%83]))
		divisor /= 83
	}
}

func encodeDC(r, g, b float64) int {
	return (transforms.LinearTosRGB(r) << 16) + (transforms.LinearTosRGB(g) << 8) + transforms.LinearTosRGB(b)
}

func encodeAC(r, g, b, maximumValue float64) int {
	quant := func(f float64) int {
		return int(math.Max(0, math.Min(18, math.Floor(transforms.SignPow(f/maximumValue, 0.5)*9+9.5))))
	}
	return quant(r)*19*19 + quant(g)*19 + quant(b)
}
