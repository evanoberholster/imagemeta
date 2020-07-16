// Package meta provides a Metadata interface for interpreting metadata
// from different image types such as JPEG, TIFF, HEIF.
package meta

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/exiftool/imagetype"
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
		m, err = ScanJPEG(reader)
		if err != nil {
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
