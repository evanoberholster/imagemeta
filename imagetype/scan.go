package imagetype

import (
	"bufio"
	"io"
)

const (
	// searchHeaderLength is the number of bytes to read while searching for an Image Header
	searchHeaderLength = 16
)

// Scan -
// TODO: Documentation
func Scan(reader io.Reader) (imageType ImageType, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	br := bufio.NewReader(reader)
	if err != nil {
		panic(err)
	}

	// Parse Header for an ImageType
	imageType = parseHeader(br)

	return
}

func parseHeader(br *bufio.Reader) ImageType {
	buf, err := br.Peek(searchHeaderLength)
	if err != nil {
		if err == io.EOF {
			return ImageUnknown
		}
		panic(err)
	}

	if len(buf) < searchHeaderLength {
		panic(ErrDataLength)
	}

	// JPEG Header
	if isJPEG(buf) {
		return ImageJPEG
	}

	// JPEG2000 Header
	if isJPEG2000(buf) {
		return ImageJPEG
	}

	// Canon CRW Header
	if isCRW(buf) {
		return ImageCRW
	}

	// Canon CR2 Header
	if isCR2(buf) {
		return ImageCR2
	}

	// Canon CR3 Header
	if isCR3(buf) {
		return ImageCR3
	}

	// Heif Header
	if isHeif(buf) {
		return ImageHEIF
	}

	// Panasonic/Leica Raw Header
	if isRW2(buf) {
		return ImagePanaRAW
	}

	// Tiff Header
	if isTiff(buf) {
		return ImageTiff
	}

	// PNG Header
	if isPNG(buf) {
		return ImagePNG
	}

	// PSD Header
	if isPSD(buf) {
		return ImagePSD
	}

	// BMP Header
	if isBMP(buf) {
		return ImageBMP
	}

	// Webp Header
	if isWebP(buf) {
		return ImageWebP
	}

	// XMP file Header
	if isXMP(buf) {
		return ImageXMP
	}

	return ImageUnknown
}
