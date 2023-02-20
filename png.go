package imagemeta

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/png"
)

// DecodePng decodes a PNG file from an io.Reader returning Exif or an error.
func DecodePng(r io.ReadSeeker) (exif2.Exif, error) {
	header, err := png.ScanPngHeader(r)
	if err != nil {
		return exif2.Exif{}, err
	}

	ir := exif2.NewIfdReader(exif2.Logger)
	defer ir.Close()

	if err := ir.DecodeTiff(r, header); err != nil {
		return ir.Exif, err
	}

	return ir.Exif, nil
}
