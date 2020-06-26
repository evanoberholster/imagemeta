package tiffmeta

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

// Errors
var (
	// ErrNoExif no exif information was found
	ErrNoExif = errors.New("no exif found")
)

const (
	// Tiff Header Length is 8 bytes
	tiffHeaderLength = 16
)

// Scan searches an io.Reader for a LittleEndian Tiff Header or a BigEndian Tiff Header
// and returns the TiffHeader
func Scan(reader io.Reader) (m Metadata, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	br := bufio.NewReader(reader)
	return scan(br)
}

func ScanBuf(reader *bufio.Reader) (m Metadata, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	return scan(reader)
}

// Search for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func scan(br *bufio.Reader) (m Metadata, err error) {
	m = Metadata{}
	discarded := 0

	var buf []byte

	for {
		buf, err = br.Peek(tiffHeaderLength)
		if err != nil {
			if err == io.EOF {
				err = ErrNoExif
				break
			}
			panic(err)
		}
		if len(buf) < 8 {
			err = ErrNoExif
			break
		}

		byteOrder := BinaryOrder(buf)
		if byteOrder == nil {
			// Exif not identified. Move forward by one byte.
			if _, err = br.Discard(1); err != nil {
				panic(err)
			}

			discarded++
			continue
		}

		// Found
		firstIfdOffset := byteOrder.Uint32(buf[4:8])
		tiffHeaderOffset := uint32(discarded)
		m.Header = NewHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0)
		return
	}

	return
}

// BinaryOrder returns the binary.ByteOrder for a Tiff Header
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func BinaryOrder(buf []byte) binary.ByteOrder {
	if len(buf) < 4 {
		return nil
	}

	if isTiffBigEndian(buf[:4]) {
		return binary.BigEndian
	} else if isTiffLittleEndian(buf[:4]) {
		return binary.LittleEndian
	}
	return nil
}

// isTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func isTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// isTiffBigEndian checks the buf for the TiffBigEndianSignature
func isTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}
