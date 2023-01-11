package imagemeta

import (
	"bufio"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/pkg/errors"
)

// readerPool for buffer
var readerPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, 4*1024) },
}

// DecodeCR3 decodes a CR3 file from an io.Reader returning Exif or an error.
func DecodeCR3(r io.ReadSeeker) (exif2.Exif, error) {
	rr := readerPool.Get().(*bufio.Reader)
	rr.Reset(r)
	defer readerPool.Put(rr)

	bmr := bmff.NewReader(rr)
	ir := exif2.NewIfdReader(nil)
	defer ir.Close()

	bmr.ExifReader = ir.DecodeCR3Ifd
	ftyp, err := bmr.ReadFtypBox()
	if err != nil {
		return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
	}
	moov, err := bmr.ReadCrxMoovBox()
	if err != nil {
		return ir.Exif, errors.Wrapf(err, "ReadCrxMoovBox")
	}
	// Set ImageType to CR3
	ir.Exif.ImageType = imagetype.ImageCR3
	_ = ftyp
	_ = moov
	//fmt.Println(moov.Meta.XMP())
	//fmt.Println(moov)
	return ir.Exif, nil
	//return exif2.DecodeHeader(r, moov.Meta.Exif[0], moov.Meta.Exif[1], moov.Meta.Exif[3])
}
