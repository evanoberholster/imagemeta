package meta

import "io"

type Decoder struct {
	XMPDecodeFn  DecodeFn
	ExifDecodeFn DecodeFn
}

type DecodeFn func(r io.Reader) error
