package imagehash

import (
	"errors"
	"image"
	"sync"

	"github.com/evanoberholster/imagemeta/imagehash/transforms32"
)

// NewPHash64Alt is a Perception Hash function returns a hash computation of phash.
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash64Alt(img image.Image) (phash PHash64, err error) {
	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 64 {
		err = errors.New("error image size incompatible. PHash requires 64x64 image")
		return
	}

	pixels, ok := pixelsPool32.Get().(*[]float32)
	if !ok || pixels == nil {
		return 0, errors.New("pixelsPool32 returned non-*[]float32")
	}
	defer pixelsPool32.Put(pixels)

	transforms32.ImageToGray(img, pixels)
	flattens := transforms32.DCT2DHash64(*pixels)
	median := transforms32.MedianOfPixels64(flattens[:])

	for idx, p := range flattens {
		if p > median {
			phash |= 1 << (63 - idx) // leftShiftSet
		}
	}
	return phash, nil
}

// NewPHash256Alt is a Perception Hash function returns a 256bit hash computation
// Implementation follows: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Optimized for performance and reduced memory footprint.
func NewPHash256Alt(img image.Image) (phash PHash256, err error) {
	var size image.Point
	if img != nil {
		size = img.Bounds().Size()
	}
	if size.X != size.Y && size.X != 256 {
		err = errors.New("error image size incompatible. PHash256 requires 256x256 image")
		return
	}

	pixels, ok := pixelsPool256Alt.Get().(*[]float32)
	if !ok || pixels == nil {
		return PHash256{}, errors.New("pixelsPool256Alt returned non-*[]float32")
	}
	defer pixelsPool256Alt.Put(pixels)

	transforms32.ImageToGray(img, pixels)
	flattens := transforms32.DCT2DHash256(pixels)

	median := transforms32.MedianOfPixels256(flattens[:])

	for idx, p := range flattens {
		indexOfArray := idx / 64
		if p > median {
			phash[indexOfArray] |= 1 << (63 - idx%64) // leftShiftSet
		}
	}

	return phash, nil
}

// Pixel pool 64bit
var pixelsPool32 = sync.Pool{
	New: func() interface{} {
		p := make([]float32, 4096)
		return &p
	},
}

// Pixel pool 256bit Hash
var pixelsPool256Alt = sync.Pool{
	New: func() interface{} {
		p := make([]float32, 65536)
		return &p
	},
}
