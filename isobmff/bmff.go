package isobmff

import (
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

type ExifReader func(r io.Reader, h meta.ExifHeader) error
type XMPReader func(r io.Reader) error
type PreviewImageReader func(r io.Reader, h meta.PreviewHeader) error

const (
	optionSpeed uint8 = 1
)
