package mknote

import (
	"encoding/binary"
	"errors"
	"io"
)

// Errors
var (
	ErrNikonMkNote = errors.New("err makernote is not a Nikon makernote")
)

// NikonMkNoteHeader parses the Nikon Makernote from reader and returns byteOrder and error
// TODO: Exhaustatively Test and refactor
func NikonMkNoteHeader(reader io.Reader) (byteOrder binary.ByteOrder, err error) {
	// Nikon Makernotes header is 18 bytes. Move Reader up necessary bytes
	mknoteHeader := [18]byte{}
	if n, err := reader.Read(mknoteHeader[:]); n < 18 || err != nil {
		return nil, ErrNikonMkNote
	}
	// Nikon makernote header starts with "Nikon" with the first 5 bytes
	if isNikonMkNoteHeaderBytes(mknoteHeader[:5]) {
		if byteOrder := binaryOrder(mknoteHeader[10:14]); byteOrder != nil {
			return byteOrder, nil
		}
	}

	return nil, ErrNikonMkNote
}

// Nikon Makernote Header
// isNikonMkNoteHeaderBytes represents "Nikon" the first 5 bytes of the
func isNikonMkNoteHeaderBytes(buf []byte) bool {
	return buf[0] == 'N' &&
		buf[1] == 'i' &&
		buf[2] == 'k' &&
		buf[3] == 'o' &&
		buf[4] == 'n'
}

// BinaryOrder returns the binary.ByteOrder for a Tiff Header based
// on 4 bytes from the buf.
//
// Good reference:
// CIPA DC-008-2016; JEITA CP-3451D
// -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
func binaryOrder(buf []byte) binary.ByteOrder {
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
