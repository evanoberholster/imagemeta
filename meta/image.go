package meta

import "fmt"

//go:generate msgp

// Dimensions stores width and height
type Dimensions struct {
	Width  uint32
	Height uint32
}

// NewDimensions returns Dimensions with the given width and height
func NewDimensions(width, height uint32) Dimensions {
	if width == 0 || height == 0 {
		return Dimensions{}
	}
	return Dimensions{Width: width, Height: height}
}

func (d Dimensions) String() string {
	width, height := d.Size()
	return fmt.Sprintf("width: %d, height: %d", width, height)
}

// Size returns width and height from underlying dimensions
func (d Dimensions) Size() (width, height uint32) {
	return d.Width, d.Height
}

// AspectRatio calculates the Aspect ratio of the Image (Width/Height)
func (d Dimensions) AspectRatio() float32 {
	if d.Width == 0 {
		return 0.0
	}
	return float32(d.Width) / float32(d.Height)
}

// Orientation -
func (d Dimensions) Orientation() uint {
	if d.AspectRatio() < 1 {
		return 1
	}
	return 0
}
