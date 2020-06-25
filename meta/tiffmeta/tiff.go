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
	ErrNoExif = errors.New("No exif found")
)

const (
	// Tiff Header Length is 8 bytes
	tiffHeaderLength = 8
)

// Scan searches an io.Reader for a LittleEndian Tiff Header or a BigEndian Tiff Header
// and returns the TiffHeader
func Scan(reader io.Reader) (h Header, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	// Search for the beginning of the EXIF information. The EXIF is near the
	// beginning of most Image files, so this likely doesn't have a high cost.
	br := bufio.NewReader(reader)
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
		return NewHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0), nil
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

	if IsTiffBigEndian(buf[:4]) {
		return binary.BigEndian
	} else if IsTiffLittleEndian(buf[:4]) {
		return binary.LittleEndian
	}
	return nil
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func IsTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func IsTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}
