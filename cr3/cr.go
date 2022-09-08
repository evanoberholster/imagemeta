package cr3

import (
	"io"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/pkg/errors"
)

func Decode(r io.ReadSeeker) (exif2.Exif, error) {

	//
	bmr := bmff.NewReader(r)
	ir := exif2.NewIfdReader(nil)
	defer ir.Close()

	bmr.ExifReader = ir.DecodeIfd
	ftyp, err := bmr.ReadFtypBox()
	moov, err := bmr.ReadCrxMoovBox()
	if err != nil {
		return exif2.Exif{}, errors.Wrapf(err, "ReadCrxMoovBox")
	}
	_ = ftyp
	_ = moov
	return ir.Exif(), nil
	//return exif2.DecodeHeader(r, moov.Meta.Exif[0], moov.Meta.Exif[1], moov.Meta.Exif[3])
}

//m.CrxMoov, err = bmr.ReadCrxMoovBox()
//if err != nil {
//	return errors.Wrapf(err, "ReadCrxMoovBox")
//}
//for i, header := range m.CrxMoov.Meta.Exif {
//	if i == 0 {
//		m.e, err = exif.ParseExif(m.mr, header)
//		if err != nil {
//			fmt.Println(err)
//		}
//		continue
//	}
//	if err = m.e.ParseIfd(header); err != nil {
//		return err
//	}
//}
