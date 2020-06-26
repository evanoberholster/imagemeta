package mknote

import (
	"encoding/binary"
	"errors"
	"io"
)

// Errors
var (
	ErrNikonMkNote = errors.New("makernote is not a Nikon makernote")
)

// NikonMkNoteHeader parses the Nikon Makernote from reader and returns byteOrder and error
// TODO: Exhaustatively Test and refactor
func NikonMkNoteHeader(reader io.Reader) (byteOrder binary.ByteOrder, err error) {
	// Nikon Makernotes header is 18 bytes. Move Reader up necessary bytes
	mknoteHeader := make([]byte, 18)
	if n, err := reader.Read(mknoteHeader); n < 18 || err != nil {
		err = ErrNikonMkNote
		return nil, err
	}
	// Nikon makernote header starts with "Nikon" with the first 5 bytes
	if isNikonMkNoteHeaderBytes(mknoteHeader[:5]) {
		if isTiffBigEndian(mknoteHeader[10:14]) {
			byteOrder = binary.BigEndian
			return byteOrder, nil
		} else if isTiffLittleEndian(mknoteHeader[10:14]) {
			byteOrder = binary.LittleEndian
			return byteOrder, nil
		}
	}

	err = ErrNikonMkNote
	return
}

// Nikon Makernote Header
// isNikonMkNoteHeaderBytes represents "Nikon" the first 5 bytes of the
func isNikonMkNoteHeaderBytes(buf []byte) bool {
	return buf[0] == 0x4e &&
		buf[1] == 0x69 &&
		buf[2] == 0x6b &&
		buf[3] == 0x6f &&
		buf[4] == 0x6e
	//nikonMkNoteHeaderBytes = []byte{0x4e, 0x69, 0x6b, 0x6f, 0x6e}
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
