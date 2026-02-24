package isobmff

import (
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

type ExifReader func(r io.Reader, h meta.ExifHeader) error
type XMPReader func(r io.Reader, h XPacketHeader) error
type PreviewImageReader func(r io.Reader, h meta.PreviewHeader) error

type XPacketHeader struct {
	Offset       uint64
	Length       uint32
	HasXPacketPI bool
	HasXMPMeta   bool
}
