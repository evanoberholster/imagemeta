package imagetype

import (
	"bufio"
	"io"
)

const (
	// searchHeaderLength is the number of bytes to read while searching for an Image Header
	searchHeaderLength = 24
)

// Scan is a conveninence function for ScanBuf
func Scan(reader io.Reader) (imageType ImageType, err error) {
	// Parse Header for an ImageType
	br := bufio.NewReader(reader)
	return ScanBuf(br)
}

// ScanBuf -
// TODO: Documentation
func ScanBuf(reader *bufio.Reader) (imageType ImageType, err error) {
	// Parse Header for an ImageType
	imageType = parseHeader(reader)

	return
}

// parseHeader
func parseHeader(br *bufio.Reader) ImageType {
	buf, err := br.Peek(searchHeaderLength)
	if err != nil {
		return ImageUnknown
	}

	//if len(buf) < searchHeaderLength {
	//	panic(ErrDataLength)
	//}

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

	// ISOBMFF Header
	if isFTYPBox(buf) {
		// AVIF Header
		if isAVIF(buf) {
			return ImageAVIF
		}
		// Heif Header
		if isHeif(buf) {
			return ImageHEIF
		}
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
