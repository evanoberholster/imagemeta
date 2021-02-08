// Package metadata provides a Metadata interface for interpreting metadata
// from different image types such as JPEG, TIFF, HEIF, and AVIF.
package metadata

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
)

// Errors
var (
	ErrNoExif               = errors.New("error No Exif")
	ErrMetadataNotSupported = errors.New("metadata not supported for this imagetype")
)

// Metadata interface
type Metadata interface {
	Size() (width uint16, height uint16)
	Header() TiffHeader
	XML() string
}

// Scan -
func Scan(reader io.Reader, t imagetype.ImageType) (m Metadata, err error) {
	return ScanBuf(bufio.NewReader(reader), t)
}

// ScanBuf -
func ScanBuf(reader *bufio.Reader, t imagetype.ImageType) (m Metadata, err error) {
	switch t {
	case imagetype.ImageJPEG:
		if m, err = ScanJPEG(reader, nil, nil); err != nil {
			err = ErrNoExif
		}
		return
	case imagetype.ImageWebP:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageXMP:
		err = ErrMetadataNotSupported
		return
	default:
		m, err = ScanTiff(reader)
		if err == ErrNoExif {
			err = ErrNoExif
		}
		return
	}
}

// ScanBuf2 -
func ScanBuf2(reader *bufio.Reader, t imagetype.ImageType, xmpDecodeFn DecodeFn) (m Metadata, err error) {
	switch t {
	case imagetype.ImageJPEG:
		if m, err = ScanJPEG(reader, xmpDecodeFn, nil); err != nil {
			err = ErrNoExif
		}
		return
	case imagetype.ImageWebP:
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageXMP:
		err = ErrMetadataNotSupported
		return
	default:
		m, err = ScanTiff(reader)
		if err == ErrNoExif {
			err = ErrNoExif
		}
		return
	}
}
