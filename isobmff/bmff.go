package isobmff

import (
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

type ExifReader func(r io.Reader, h meta.ExifHeader) error
