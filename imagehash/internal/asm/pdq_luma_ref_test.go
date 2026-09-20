package asm

// pdqLumaRef implements ycbcrToLuma in plain Go for cross-checking the SIMD
// luminance kernels.
func pdqLumaRef(y, cb, cr uint8) float32 {
	yy := int32(y) * 0x10101
	cb1 := int32(cb) - 128
	cr1 := int32(cr) - 128
	r := min(max((yy+91881*cr1)>>16, 0), 255)
	g := min(max((yy-22554*cb1-46802*cr1)>>16, 0), 255)
	b := min(max((yy+116130*cb1)>>16, 0), 255)
	return 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
}
