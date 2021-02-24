// Package tiff reads Tiff Header metadata information from image files before being processed by exif package
package tiff

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 16
)

// Scan searches an io.Reader for a LittleEndian or BigEndian Tiff Header
// and returns the TiffHeader
func Scan(r io.Reader) (Header, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReaderSize(r, 64)
	}
	return scan(br)
}

// scan searchs for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func scan(br *bufio.Reader) (header Header, err error) {
	discarded := 0

	var buf []byte

	for {
		if buf, err = br.Peek(TiffHeaderLength); err != nil {
			err = meta.ErrNoExif
			return
		}

		byteOrder := BinaryOrder(buf)
		if byteOrder == nil {
			// Exif not identified. Move forward by one byte.
			if buf[1] == 0x49 || buf[1] == 0x4d {
				_, _ = br.Discard(1)
				discarded++
				continue
			}
			_, _ = br.Discard(2)
			discarded += 2
			continue
		}

		// Found Tiff Header
		firstIfdOffset := byteOrder.Uint32(buf[4:8])
		tiffHeaderOffset := uint32(discarded)
		return NewHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0, imagetype.ImageTiff), nil
	}
}

// BinaryOrder returns the binary.ByteOrder for a Tiff Header
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func BinaryOrder(buf []byte) binary.ByteOrder {
	if isTiffBigEndian(buf[:4]) {
		return binary.BigEndian
	}
	if isTiffLittleEndian(buf[:4]) {
		return binary.LittleEndian
	}
	return nil
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func isTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func isTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}
