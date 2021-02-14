// Package tiff reads Tiff Image metadata information from Camera Raw, DNG, and Tiff Image files
package tiff

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 16
)

// ScanTiff searches an io.Reader for a LittleEndian or BigEndian Tiff Header
// and returns the TiffHeader
func ScanTiff(reader *bufio.Reader) (m TiffMetadata, err error) {
	return scanTIFF(reader)
}

// Search for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func scanTIFF(br *bufio.Reader) (tm TiffMetadata, err error) {
	//tm = TiffMetadata{}
	discarded := 0

	var buf []byte

	for {
		if buf, err = br.Peek(TiffHeaderLength); err != nil {
			if err == io.EOF {
				err = meta.ErrNoExif
				return
			}
			return
		}

		byteOrder := BinaryOrder(buf)
		if byteOrder == nil {
			// Exif not identified. Move forward by one byte.
			_, _ = br.Discard(1)
			discarded++
			continue
		}

		// Found Tiff Header
		firstIfdOffset := byteOrder.Uint32(buf[4:8])
		tiffHeaderOffset := uint32(discarded)
		tm.header = NewHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0)
		return
	}
}

// BinaryOrder returns the binary.ByteOrder for a Tiff Header
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func BinaryOrder(buf []byte) binary.ByteOrder {
	if len(buf) == 16 {
		if isTiffBigEndian(buf[:4]) {
			return binary.BigEndian
		}
		if isTiffLittleEndian(buf[:4]) {
			return binary.LittleEndian
		}
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
