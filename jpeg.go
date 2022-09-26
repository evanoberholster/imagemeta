package imagemeta

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/jpeg"
)

// DecodeJPEG decodes a JPEG file from an io.Reader returning Exif or an error.
func DecodeJPEG(r io.ReadSeeker) (exif2.Exif, error) {
	rr := readerPool.Get().(*bufio.Reader)
	rr.Reset(r)
	defer readerPool.Put(rr)

	ir := exif2.NewIfdReader(nil)
	defer ir.Close()

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif2.Exif{}, err
	}
	if it != imagetype.ImageJPEG {
		return exif2.Exif{}, ErrMetadataNotSupported
	}

	if err = jpeg.ScanJPEG(rr, ir.DecodeJPEGIfd, nil); err != nil {
		return exif2.Exif{}, err
	}
	ir.Exif.ImageType = it
	return ir.Exif, nil
}
