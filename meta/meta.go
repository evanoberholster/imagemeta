package meta

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/exiftool/imagetype"
	"github.com/evanoberholster/exiftool/meta/jpegmeta"
	"github.com/evanoberholster/exiftool/meta/tiffmeta"
)

// Errors
var (
	ErrNoExif               = errors.New("Error No Exif")
	ErrMetadataNotSupported = errors.New("Metadata not supported for this imagetype")
)

// Metadata interface
type Metadata interface {
	Size() (width uint16, height uint16)
	TiffHeader() tiffmeta.Header
	XML() string
}

func Scan(reader io.Reader, t imagetype.ImageType) (m Metadata, err error) {
	switch t {
	case imagetype.ImageJPEG:
		m, err = jpegmeta.Scan(reader)
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
		m, err = tiffmeta.Scan(reader)
		if err == tiffmeta.ErrNoExif {
			err = ErrNoExif
		}
		return
	}
}

func ScanBuf(reader *bufio.Reader, t imagetype.ImageType) (m Metadata, err error) {
	switch t {
	case imagetype.ImageJPEG:
		m, err = jpegmeta.ScanBuf(reader)
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
		m, err = tiffmeta.ScanBuf(reader)
		if err == tiffmeta.ErrNoExif {
			err = ErrNoExif
		}
		return
	}
}
