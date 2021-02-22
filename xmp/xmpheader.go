package xmp

import "io"

// Header is an XMP header in an image file.
// Contains Offset and Length of XMP metadata.
type Header struct {
	Offset, Length uint32
}

// NewHeader returns a new xmp.Header with an offset
// and length of where to read XMP metadata.
func NewHeader(offset, length uint32) Header {
	return Header{offset, length}
}

// DecodeFn is a function for decoding Xmp Metadata
type DecodeFn func(r io.Reader, header Header) (err error)
