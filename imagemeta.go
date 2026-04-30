// Package imagemeta provides functions for parsing and extracting Metadata from Images.
// Different image types such as JPEG, Camera Raw, DNG, TIFF, HEIF, and AVIF.
package imagemeta

import (
	"bufio"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/evanoberholster/imagemeta/meta/jpeg"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/png"
	"github.com/evanoberholster/imagemeta/preview"
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

func getPooledReader() (*bufio.Reader, error) {
	rr, ok := readerPool.Get().(*bufio.Reader)
	if !ok || rr == nil {
		return nil, errors.New("readerPool returned non-*bufio.Reader")
	}
	return rr, nil
}

func Decode(r io.ReadSeeker) (exif.Exif, error) {
	rr, err := getPooledReader()
	if err != nil {
		return exif.Exif{}, err
	}
	rr.Reset(r)
	defer readerPool.Put(rr)

	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif.Exif{}, err
	}
	ir.Exif.ImageType = it
	switch it {
	case imagetype.ImageJPEG:
		if err = jpeg.ScanJPEGWithSource(rr, r, ir.DecodeJPEGIfd, nil); err != nil {
			return exif.Exif{}, err
		}
	case imagetype.ImageCRW:
		return DecodeCRW(r)
	case imagetype.ImageCR2, imagetype.ImageTiff, imagetype.ImagePanaRAW, imagetype.ImageDNG, imagetype.ImageNEF, imagetype.ImageARW:
		header, err := exif.ScanTiffHeader(rr, it)
		if err != nil {
			return exif.Exif{}, err
		}
		if err := ir.DecodeTiff(rr, header); err != nil {
			return ir.Exif, err
		}
	case imagetype.ImageCR3, imagetype.ImageAVIF, imagetype.ImageJXL, imagetype.ImageHEIF, imagetype.ImageHEIC:
		bmr := isobmff.NewReader(rr, ir.DecodeIfdAppend, nil, nil)
		defer bmr.Close()
		if err := bmr.ReadFTYP(); err != nil {
			return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
		}
		if err := bmr.ReadMetadataUntilEOF(); err != nil {
			return ir.Exif, err
		}
	case imagetype.ImagePNG:
		if _, err = r.Seek(0, io.SeekStart); err != nil {
			return exif.Exif{}, err
		}
		header, err := png.ScanPngHeader(r)
		if err != nil {
			return exif.Exif{}, err
		}
		if err := ir.DecodeTiff(r, header); err != nil {
			return ir.Exif, err
		}
	default:
		header, err := exif.ScanTiffHeader(rr, it)
		if err != nil {
			return exif.Exif{}, err
		}
		if err := ir.DecodeTiff(rr, header); err != nil {
			return ir.Exif, err
		}
	}

	return ir.Exif, nil
}

// DecodeCR3 decodes a CR3 file from an io.Reader returning Exif or an error.
func DecodeCR3(r io.ReadSeeker) (exif.Exif, error) {
	rr, err := getPooledReader()
	if err != nil {
		return exif.Exif{}, err
	}
	defer readerPool.Put(rr)
	rr.Reset(r)

	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	bmr := isobmff.NewReader(rr, ir.DecodeIfdAppend, nil, nil)
	defer bmr.Close()
	if err := bmr.ReadFTYP(); err != nil {
		return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
	}
	if err := bmr.ReadMetadataUntilEOF(); err != nil {
		return ir.Exif, err
	}
	return ir.Exif, nil
}

// DecodeTiff decodes a Tiff/DNG file from an io.Reader returning Exif or an error.
func DecodeTiff(r io.ReadSeeker) (exif.Exif, error) {
	rr, err := getPooledReader()
	if err != nil {
		return exif.Exif{}, err
	}
	rr.Reset(r)
	defer readerPool.Put(rr)

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif.Exif{}, err
	}
	header, err := exif.ScanTiffHeader(rr, it)
	if err != nil {
		return exif.Exif{}, err
	}
	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	if err := ir.DecodeTiff(rr, header); err != nil {
		return ir.Exif, err
	}
	return ir.Exif, nil
}

// DecodeCR2 decodes a CR2 file from an io.Reader returning Exif or an error.
func DecodeCR2(r io.ReadSeeker) (exif.Exif, error) {
	return DecodeTiff(r)
}

// DecodeARW decodes a Sony ARW file from an io.Reader returning Exif or an error.
func DecodeARW(r io.ReadSeeker) (exif.Exif, error) {
	return DecodeTiff(r)
}

// DecodeCRW decodes a CRW file from an io.ReadSeeker returning Exif or an error.
func DecodeCRW(r io.ReadSeeker) (exif.Exif, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return exif.Exif{}, err
	}
	it, err := imagetype.Scan(r)
	if err != nil {
		return exif.Exif{}, err
	}
	if it != imagetype.ImageCRW {
		return exif.Exif{}, ErrMetadataNotSupported
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return exif.Exif{}, err
	}
	ciff, err := jpeg.ParseCIFFReader(r)
	if err != nil {
		return exif.Exif{}, err
	}
	return exif.FromCIFF(ciff, it), nil
}

// DecodeHeif decodes a Heif file from an io.Reader returning Exif or an error.
// Needs improvement
func DecodeHeif(r io.ReadSeeker) (exif.Exif, error) {
	rr, err := getPooledReader()
	if err != nil {
		return exif.Exif{}, err
	}
	defer readerPool.Put(rr)
	rr.Reset(r)

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif.Exif{}, err
	}
	if it != imagetype.ImageHEIC && it != imagetype.ImageHEIF {
		return exif.Exif{}, ErrMetadataNotSupported
	}

	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	bmr := isobmff.NewReader(rr, ir.DecodeIfdAppend, nil, nil)
	defer bmr.Close()
	if err := bmr.ReadFTYP(); err != nil {
		return ir.Exif, errors.Wrapf(err, "ReadFtypBox")
	}
	if err := bmr.ReadMetadataUntilEOF(); err != nil {
		return ir.Exif, err
	}
	return ir.Exif, nil
}

// DecodeJPEG decodes a JPEG file from an io.Reader returning Exif or an error.
func DecodeJPEG(r io.ReadSeeker) (exif.Exif, error) {
	rr, err := getPooledReader()
	if err != nil {
		return exif.Exif{}, err
	}
	rr.Reset(r)
	defer readerPool.Put(rr)

	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	it, err := imagetype.ScanBuf(rr)
	if err != nil {
		return exif.Exif{}, err
	}
	if it != imagetype.ImageJPEG {
		return exif.Exif{}, ErrMetadataNotSupported
	}

	if err = jpeg.ScanJPEGWithSource(rr, r, ir.DecodeJPEGIfd, nil); err != nil {
		return exif.Exif{}, err
	}
	ir.Exif.ImageType = it
	return ir.Exif, nil
}

// DecodePng decodes a PNG file from an io.Reader returning Exif or an error.
func DecodePng(r io.ReadSeeker) (exif.Exif, error) {
	header, err := png.ScanPngHeader(r)
	if err != nil {
		return exif.Exif{}, err
	}

	ir := exif.AcquirePooledReader(metalog.GetLogger())
	defer exif.ReleasePooledReader(ir)

	if err := ir.DecodeTiff(r, header); err != nil {
		return ir.Exif, err
	}

	return ir.Exif, nil
}

// PreviewCR3 previews a CR3 file from an io.Reader returning preview image binary or an error.
func PreviewCR3(r io.ReadSeeker) ([]byte, error) {
	rr, err := getPooledReader()
	if err != nil {
		return nil, err
	}
	defer readerPool.Put(rr)
	rr.Reset(r)

	pr := preview.NewPreviewReader(preview.Logger)

	bmr := isobmff.NewReader(rr, nil, nil, pr.RenderPreview)
	defer bmr.Close()

	if err := bmr.ReadFTYP(); err != nil {
		return nil, errors.Wrapf(err, "ReadFtypBox")
	}

	if err := bmr.ReadMetadataUntilEOF(); err != nil {
		return nil, err
	}

	return pr.PreviewImage, nil
}
