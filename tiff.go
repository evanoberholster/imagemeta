package imagemeta

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/tiff"
)

// DecodeTiff decodes a Tiff/DNG file from an io.Reader returning Exif or an error.
func DecodeTiff(r io.ReadSeeker) (exif2.Exif, error) {
	rr := readerPool.Get().(*bufio.Reader)
	rr.Reset(r)
	defer readerPool.Put(rr)

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif2.Exif{}, err
	}
	header, err := tiff.ScanTiffHeader(rr, it)
	if err != nil {
		return exif2.Exif{}, err
	}
	ir := exif2.NewIfdReader(rr)
	defer ir.Close()

	if err := ir.DecodeTiff(r, header); err != nil {
		return ir.Exif, err
	}
	return ir.Exif, nil
	//return exif2.DecodeHeader(r, moov.Meta.Exif[0], moov.Meta.Exif[1], moov.Meta.Exif[3])
}

// DecodeCR2 decodes a CR2 file from an io.Reader returning Exif or an error.
func DecodeCR2(r io.ReadSeeker) (exif2.Exif, error) {
	return DecodeTiff(r)
}
