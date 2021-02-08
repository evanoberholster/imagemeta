package metadata

import "io"

type Decoder struct {
	xmpDecodeFn  DecodeFn
	exifDecodeFn DecodeFn
}

type DecodeFn func(r io.Reader) error
