// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"

	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/isobmff"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/pkg/errors"
)

// Errors
var (
	ErrNoExif               = meta.ErrNoExif
	ErrNoExifDecodeFn       = errors.New("error no Exif Decode Func set")
	ErrNoXmpDecodeFn        = errors.New("error no Xmp Decode Func set")
	ErrImageTypeNotFound    = imagetype.ErrImageTypeNotFound
	ErrMetadataNotSupported = errors.New("error metadata reading not supported for this imagetype")
)

// readerPool for buffer
var readerPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, 4*1024) },
}

func Decode(r io.ReadSeeker) (exif2.Exif, error) {
	rr := readerPool.Get().(*bufio.Reader)
	rr.Reset(r)
	defer readerPool.Put(rr)

	ir := exif2.NewIfdReader(rr)
	defer ir.Close()

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif2.Exif{}, err
	}
	switch it {
	case imagetype.ImageJPEG:
		if err = jpeg.ScanJPEG(rr, ir.DecodeJPEGIfd, nil); err != nil {
			return exif2.Exif{}, err
		}
	case imagetype.ImageCR2, imagetype.ImageTiff:
		header, err := tiff.ScanTiffHeader(rr, it)
		if err != nil {
			return exif2.Exif{}, err
		}
		if err := ir.DecodeTiff(r, header); err != nil {
			return ir.Exif, err
		}
	case imagetype.ImageCR3:
		bmr := bmff.NewReader(rr)
		bmr.ExifReader = ir.DecodeIfd
		_, err := bmr.ReadFtypBox()
		if err != nil {
			return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
		}
		if _, err = bmr.ReadCrxMoovBox(); err != nil {
			return ir.Exif, errors.Wrapf(err, "ReadCrxMoovBox")
		}
	case imagetype.ImageHEIF:
		header, err := tiff.ScanTiffHeader(rr, it)
		if err != nil {
			return exif2.Exif{}, err
		}
		if err := ir.DecodeTiff(r, header); err != nil {
			return ir.Exif, err
		}
	default:
		return exif2.Exif{}, ErrMetadataNotSupported
	}
	ir.Exif.ImageType = it
	return ir.Exif, nil
}

// DecodeCR3 decodes a CR3 file from an io.Reader returning Exif or an error.
func DecodeCR3(r io.ReadSeeker) (exif2.Exif, error) {
	rr := readerPool.Get().(*bufio.Reader)
	rr.Reset(r)
	defer readerPool.Put(rr)

	ir := exif2.NewIfdReader(nil)
	defer ir.Close()

	bmr := isobmff.NewReader(rr)
	defer bmr.Close()
	bmr.ExifReader = ir.DecodeIfd
	if err := bmr.ReadFTYP(); err != nil {
		return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
	}
	if err := bmr.ReadMetadata(); err != nil {
		panic(err)
	}
	// Set ImageType to CR3
	ir.Exif.ImageType = imagetype.ImageCR3
	return ir.Exif, nil
	//return exif2.DecodeHeader(r, moov.Meta.Exif[0], moov.Meta.Exif[1], moov.Meta.Exif[3])
}

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
