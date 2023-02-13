package heic

import (
	"io"

	"github.com/evanoberholster/imagemeta/exif2"
)

func Decode(r io.ReadSeeker) (exif2.Exif, error) {

	//
	//rr := readerPool.Get().(*bufio.Reader)
	//rr.Reset(r)
	//defer readerPool.Put(rr)
	//
	//bmr := bmff.NewReader(rr)
	//ir := exif2.NewIfdReader(nil)
	//defer ir.Close()
	//
	//bmr.ExifReader = ir.DecodeCR3Ifd
	//ftyp, err := bmr.ReadFtypBox()
	//moov, err := bmr.ReadCrxMoovBox()
	//if err != nil {
	//	return exif2.Exif{}, errors.Wrapf(err, "ReadCrxMoovBox")
	//}
	//_ = ftyp
	//_ = moov
	////fmt.Println(moov)
	return exif2.Exif{}, nil
	//return exif2.DecodeHeader(r, moov.Meta.Exif[0], moov.Meta.Exif[1], moov.Meta.Exif[3])
}
