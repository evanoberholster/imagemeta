package imagetype

import (
	"bufio"
	"errors"
	"io"
)

var (
	// ErrImageTypeNotFound is an error that represents an imagetype not being found.
	ErrImageTypeNotFound = errors.New("error imagetype not found")
)

const (
	// searchHeaderLength is the number of bytes to read while searching for an Image Header
	searchHeaderLength = 24
)

// Scan reads from the reader and returns an imageType based on
// underlying rules. Returns ErrImageTypeNotFound if imageType was not
// identified.
func Scan(reader io.Reader) (imageType ImageType, err error) {
	// Parse Header for an ImageType
	br := bufio.NewReaderSize(reader, searchHeaderLength)
	return ScanBuf(br)
}

// ScanBuf peeks at a bufio.Reader and returns an imageType based on
// underlying rules. Returns ErrImageTypeNotFound if imageType was not
// identified.
func ScanBuf(br *bufio.Reader) (imageType ImageType, err error) {
	var buf []byte

	// Peek into the bufio.Reader for the length of searchHeaderLength bytes
	if buf, err = br.Peek(searchHeaderLength); err != nil {
		return ImageUnknown, err
	}

	// Parse Header for an ImageType
	imageType = parseBuffer(buf)

	// Check if ImageType is Unknown
	if imageType == ImageUnknown {
		err = ErrImageTypeNotFound
	}
	return
}

// parseBuffer parses the []byte for image magic numbers
// that identify the imagetype. Returns an ImageType. Returns ImageUnknown
// when imagetype was not identified.
func parseBuffer(buf []byte) ImageType {
	if len(buf) < searchHeaderLength {
		return ImageUnknown
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
