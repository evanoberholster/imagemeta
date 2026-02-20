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

// Scan reads from the reader and returns a fileType based on
// underlying rules. Returns ImageUnknown and ErrImageTypeNotFound if fileType was not
// identified.
func Scan(r io.Reader) (fileType FileType, err error) {
	// Parse Header for a FileType
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < searchHeaderLength {
		br = bufio.NewReaderSize(r, searchHeaderLength)
	}
	return ScanBuf(br)
}

// ScanBuf peeks at a bufio.Reader and returns a fileType based on
// underlying rules. Returns ImageUnknown and ErrImageTypeNotFound if fileType was not
// identified.
func ScanBuf(br *bufio.Reader) (fileType FileType, err error) {
	var buf []byte

	// Peek into the bufio.Reader for the length of searchHeaderLength bytes
	if buf, err = br.Peek(searchHeaderLength); err != nil {
		return ImageUnknown, err
	}

	return Buf(buf[:])
}

// ReadAt reads from the reader at the given offset and returns a fileType based on
// underlying rules. Returns ImageUnknown and an error if fileType was not
// identified.
func ReadAt(r io.ReaderAt) (fileType FileType, err error) {
	buf := [searchHeaderLength]byte{}
	if _, err = r.ReadAt(buf[:], 0); err != nil {
		return ImageUnknown, err
	}

	return Buf(buf[:])
}

// Buf parses a []byte for image magic numbers that identify the file type.
// If []byte is less than searchHeaderLength returns ImageUnknown and ErrDataLength
// If fileType was not identified returns ImageUnknown and ErrImageTypeNotFound
func Buf(buf []byte) (fileType FileType, err error) {
	if len(buf) < searchHeaderLength {
		return ImageUnknown, ErrDataLength
	}

	// Parse Header for a FileType
	fileType = parseBuffer(buf)

	// Check if fileType is Unknown
	if fileType == ImageUnknown {
		err = ErrImageTypeNotFound
	}
	return
}

// parseBuffer parses the []byte for image magic numbers
// that identify the file type. Returns a FileType. Returns ImageUnknown
// when file type was not identified.
func parseBuffer(buf []byte) FileType {
	// JPEG Header
	if isJPEG(buf) {
		return ImageJPEG
	}

	// JPEG2000 Header
	if isJPEG2000(buf) {
		return ImageJP2K
	}

	// JPEG XL Header (codestream/container)
	if isJXL(buf) {
		return ImageJXL
	}

	// MNG/JNG Header
	if isMNG(buf) {
		return ImageMNG
	}
	if isJNG(buf) {
		return ImageJNG
	}

	// Canon CRW Header
	if isCRW(buf) {
		return ImageCRW
	}

	// Canon CR2 Header
	if isCR2(buf) {
		return ImageCR2
	}

	// ISOBMFF Header
	if isFTYPBox(buf) {
		// Canon CR3 Header
		if isCR3(buf) {
			return ImageCR3
		}
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

	// ICO/CUR Header
	if isICO(buf) {
		return ImageICO
	}
	if isCUR(buf) {
		return ImageCUR
	}

	// DDS Header
	if isDDS(buf) {
		return ImageDDS
	}

	// OpenEXR Header
	if isEXR(buf) {
		return ImageEXR
	}

	// DPX Header
	if isDPX(buf) {
		return ImageDPX
	}

	// FITS Header
	if isFITS(buf) {
		return ImageFITS
	}

	// Fuji RAF Header
	if isRAF(buf) {
		return ImageRAF
	}

	// GIMP XCF Header
	if isXCF(buf) {
		return ImageXCF
	}

	// FLIF Header
	if isFLIF(buf) {
		return ImageFLIF
	}

	// BPG Header
	if isBPG(buf) {
		return ImageBPG
	}

	// Radiance HDR Header
	if isHDR(buf) {
		return ImageHDR
	}

	// DjVu Header
	if isDJVU(buf) {
		return ImageDJVU
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

	if isGIF(buf) {
		return ImageGIF
	}

	// Netpbm family (PBM/PGM/PPM/PAM)
	if it, ok := netpbmType(buf); ok {
		return it
	}

	return ImageUnknown
}
