package imagehash

import (
	"errors"
	"fmt"
	"image"
	"math"
)

const (
	xComponents   = 4
	yComponents   = 4
	width, height = 64, 64
)

// ErrBlurHashSize is returned when the image is not 64x64.
var ErrBlurHashSize = errors.New("blurhash requires a 64x64 image")

func init() {
	for i := range channelToLinear {
		channelToLinear[i] = srgbToLinear(i)
	}
	for xc := 0; xc < xComponents; xc++ {
		for x := 0; x < width; x++ {
			xvalues32[x+width*xc] = float32(math.Cos(math.Pi * float64(xc) * float64(x) / float64(width)))
		}
	}
	for yc := 0; yc < yComponents; yc++ {
		for y := 0; y < height; y++ {
			yvalues[y+height*yc] = math.Cos(math.Pi * float64(yc) * float64(y) / float64(height))
		}
	}
	for x := 0; x < width; x++ {
		for xc := 0; xc < xComponents; xc++ {
			xvaluesT[x*xComponents+xc] = xvalues32[x+width*xc]
		}
	}
}

// EncodeBlurHashFast encodes a 64x64 image as a BlurHash string.
func EncodeBlurHashFast(img image.Image) (string, error) {
	if img == nil {
		return "", ErrImageObject
	}
	size := img.Bounds().Size()
	if size.X != width || size.Y != height {
		return "", fmt.Errorf("%w: got %dx%d", ErrBlurHashSize, size.X, size.Y)
	}

	var b blur

	// Size Flag
	b.encode((xComponents-1)+(yComponents-1)*9, 1)

	var factors [xComponents * yComponents * 3]float64
	multiplyBasisFunction(img, factors[:])

	var maximumValue float64
	var quantisedMaximumValue int
	var acCount = xComponents*yComponents - 1
	if acCount > 0 {
		var actualMaximumValue float64
		for i := 0; i < acCount*3; i++ {
			actualMaximumValue = math.Max(math.Abs(factors[i+3]), actualMaximumValue) //nolint:gosec // G602: acCount*3+3 == len(factors).
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

var (
	channelToLinear [256]float64
	yvalues         = [yComponents * height]float64{}
	// xvalues32 holds the x basis functions as float32 for the basis kernel.
	xvalues32 = [xComponents * width]float32{}
	// xvaluesT is xvalues32 transposed so all components for one column are
	// contiguous: xvaluesT[x*xComponents+xc] == xvalues32[x+width*xc].
	xvaluesT = [width * xComponents]float32{}
)

// blur is a blurhash base83 encoder
type blur struct {
	b [4 + 2*xComponents*yComponents]byte
	p int
}

func (b blur) String() string {
	return string(b.b[:])
}

const (
	characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"
)

// pow83 holds 83^n for n in [0,4] so base-83 encoding avoids math.Pow.
var pow83 = [5]int{1, 83, 83 * 83, 83 * 83 * 83, 83 * 83 * 83 * 83}

func (b *blur) encode(value, length int) {
	divisor := pow83[length-1]
	for i := 0; i < length; i++ {
		b.b[b.p] = characters[(value/divisor)%83]
		b.p++
		divisor /= 83
	}
}

func encodeDC(r, g, b float64) int {
	return (linearToSRGB(r) << 16) + (linearToSRGB(g) << 8) + linearToSRGB(b)
}

func encodeAC(r, g, b, maximumValue float64) int {
	return quantAC(r, maximumValue)*19*19 + quantAC(g, maximumValue)*19 + quantAC(b, maximumValue)
}

func quantAC(f, maximumValue float64) int {
	return int(math.Max(0, math.Min(18, math.Floor(signedSqrt(f/maximumValue)*9+9.5))))
}

// srgbToLinear converts an 8-bit sRGB channel to linear light.
func srgbToLinear(value int) float64 {
	v := float64(value) / 255
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// linearToSRGB converts a linear light value to an 8-bit sRGB channel.
func linearToSRGB(value float64) int {
	v := math.Max(0, math.Min(1, value))
	if v <= 0.0031308 {
		return int(v*12.92*255 + 0.5)
	}
	return int((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5)
}

// signedSqrt returns sign(value) * sqrt(|value|).
func signedSqrt(value float64) float64 {
	return math.Copysign(math.Sqrt(math.Abs(value)), value)
}
