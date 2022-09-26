package imagemeta

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2"
)

// DecodeHeif decodes a Heif file from an io.Reader returning Exif or an error.
// Needs improvement
func DecodeHeif(r io.ReadSeeker) (exif2.Exif, error) {
	return DecodeTiff(r)
}
