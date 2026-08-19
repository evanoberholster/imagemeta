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
	// scanHeaderLength is the number of bytes to read while scanning image headers.
	scanHeaderLength = 64
)

// Scan reads from the reader and returns a fileType based on
// underlying rules. Returns ImageUnknown and ErrImageTypeNotFound if fileType was not
// identified.
func Scan(r io.Reader) (fileType FileType, err error) {
	// Parse Header for a FileType
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < scanHeaderLength {
		br = bufio.NewReaderSize(r, scanHeaderLength)
	}
	return ScanBuf(br)
}

// ScanBuf peeks at a bufio.Reader and returns a fileType based on
// underlying rules. Returns ImageUnknown and ErrImageTypeNotFound if fileType was not
// identified.
func ScanBuf(br *bufio.Reader) (fileType FileType, err error) {
	var buf []byte

	// Peek into the bufio.Reader for the length of scanHeaderLength bytes
	if buf, err = br.Peek(scanHeaderLength); err != nil {
		if len(buf) == 0 {
			return ImageUnknown, err
		}
		return detectShortBuffer(buf)
	}

	return Buf(buf[:])
}

// ReadAt reads from the reader at the given offset and returns a fileType based on
// underlying rules. Returns ImageUnknown and an error if fileType was not
// identified.
func ReadAt(r io.ReaderAt) (fileType FileType, err error) {
	buf := [scanHeaderLength]byte{}
	n, err := r.ReadAt(buf[:], 0)
	if err != nil {
		if n == 0 {
			return ImageUnknown, err
		}
		return detectShortBuffer(buf[:n])
	}

	return Buf(buf[:])
}

func detectShortBuffer(buf []byte) (FileType, error) {
	if len(buf) >= 2 && buf[0] == 0xFF && buf[1] == 0xD8 {
		return ImageJPEG, nil
	}
	if len(buf) >= 4 {
		switch {
		case buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x00:
			return ImageTiff, nil
		case buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x00 && buf[3] == 0x2A:
			return ImageTiff, nil
		// BigTIFF (magic 43), allowed for DNG since spec 1.7.
		case buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2B && buf[3] == 0x00:
			return ImageTiff, nil
		case buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x00 && buf[3] == 0x2B:
			return ImageTiff, nil
		}
	}
	if len(buf) >= 8 {
		if isPNG(buf) {
			return ImagePNG, nil
		}
		if isCRW(buf) {
			return ImageCRW, nil
		}
	}
	if len(buf) >= 12 {
		if it := isobmffSubtype(buf); it != ImageUnknown {
			return it, nil
		}
		if isWebP(buf) {
			return ImageWebP, nil
		}
	}
	if len(buf) >= 6 && isGIF(buf) {
		return ImageGIF, nil
	}
	if len(buf) >= 2 && isBMP(buf) {
		return ImageBMP, nil
	}
	return ImageUnknown, ErrImageTypeNotFound
}

// Buf parses a []byte for image magic numbers that identify the file type.
// If []byte is less than scanHeaderLength returns ImageUnknown and ErrDataLength
// If fileType was not identified returns ImageUnknown and ErrImageTypeNotFound
func Buf(buf []byte) (fileType FileType, err error) {
	if len(buf) < scanHeaderLength {
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
	switch buf[0] {
	case 0x00:
		// ICO/CUR Header
		if isICO(buf) {
			return ImageICO
		}
		if isCUR(buf) {
			return ImageCUR
		}

		// JPEG2000 Header
		if isJPEG2000(buf) {
			return ImageJP2K
		}

		// JPEG XL Header (container)
		if isJXL(buf) {
			return ImageJXL
		}

		// ISOBMFF Header
		if it := isobmffSubtype(buf); it != ImageUnknown {
			return it
		}
	case 0xFF:
		// JPEG Header
		if isJPEG(buf) {
			return ImageJPEG
		}
		// JPEG XL Header (codestream)
		if isJXL(buf) {
			return ImageJXL
		}
	case 0x8A:
		if isMNG(buf) {
			return ImageMNG
		}
	case 0x8B:
		if isJNG(buf) {
			return ImageJNG
		}
	case 0x49, 0x4D:
		// Canon CRW Header
		if isCRW(buf) {
			return ImageCRW
		}
		// TIFF secondary subtype detection
		if it := tiffSecondarySubtype(buf); it != ImageUnknown {
			return it
		}
		// Tiff Header
		if isTiff(buf) {
			return ImageTiff
		}
	case 'D':
		if isDDS(buf) {
			return ImageDDS
		}
	case 'v':
		if isEXR(buf) {
			return ImageEXR
		}
	case 'S':
		if isDPX(buf) {
			return ImageDPX
		}
		if isFITS(buf) {
			return ImageFITS
		}
	case 'X':
		if isDPX(buf) {
			return ImageDPX
		}
	case 'F':
		if isRAF(buf) {
			return ImageRAF
		}
		if isFLIF(buf) {
			return ImageFLIF
		}
	case 'g':
		if isXCF(buf) {
			return ImageXCF
		}
	case 'B':
		if isBMP(buf) {
			return ImageBMP
		}
		if isBPG(buf) {
			return ImageBPG
		}
	case '#':
		if isHDR(buf) {
			return ImageHDR
		}
	case 'A':
		if isDJVU(buf) {
			return ImageDJVU
		}
	case 0x89:
		if isPNG(buf) {
			return ImagePNG
		}
	case '8':
		if isPSD(buf) {
			return ImagePSD
		}
	case 'R':
		if isWebP(buf) {
			return ImageWebP
		}
	case '<', ' ', '\t', '\n', '\r', 0xEF:
		if isXMP(buf) {
			return ImageXMP
		}
		if isSVG(buf) {
			return ImageSVG
		}
	case 'G':
		if isGIF(buf) {
			return ImageGIF
		}
	case 'P':
		// Netpbm family (PBM/PGM/PPM/PAM)
		if it, ok := netpbmType(buf); ok {
			return it
		}
	}

	return ImageUnknown
}
