package imagemeta

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
	XMP() string
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
func ScanBuf2(br *bufio.Reader, t imagetype.ImageType, xmpDecodeFn DecodeFn) (m Metadata, err error) {
	switch t {
	case imagetype.ImageJPEG:
		if m, err = ScanJPEG(br, xmpDecodeFn, nil); err != nil {
			err = ErrNoExif
		}
		return
	case imagetype.ImageWebP:
		// Need to implement
		err = ErrMetadataNotSupported
		return
	case imagetype.ImageXMP:
		err = xmpDecodeFn(br)
		return
	default:
		m, err = ScanTiff(br)
		if err == ErrNoExif {
			err = ErrNoExif
		}
		return
	}
}
