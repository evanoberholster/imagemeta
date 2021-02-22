package meta

import "fmt"

//go:generate msgp

// Dimensions stores width and height
type Dimensions uint64

// NewDimensions returns Dimensions with the given width and height
func NewDimensions(width, height uint32) Dimensions {
	if width == 0 || height == 0 {
		return Dimensions(0)
	}
	return Dimensions(uint64(width)<<32 + uint64(height))
}

func (d Dimensions) String() string {
	width, height := d.Size()
	return fmt.Sprintf("width: %d, height: %d", width, height)
}

// Size returns width and height from underlying dimensions
func (d Dimensions) Size() (width, height uint32) {
	height = uint32(uint64(d) << 32 >> 32)
	width = uint32(uint64(d) >> 32)
	return
}

// AspectRatio calculates the Aspect ratio of the Image (Width/Height)
func (d Dimensions) AspectRatio() float32 {
	if uint64(d) == 0 {
		return 0.0
	}
	width, height := d.Size()
	return float32(width) / float32(height)
}

// Orientation -
func (d Dimensions) Orientation() uint {
	if d.AspectRatio() < 1 {
		return 1
	}
	return 0
}
