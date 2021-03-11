package exif

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/meta"
)

// Errors
var (
	ErrNikonMkNote = errors.New("err makernote is not a Nikon makernote")
)

// NikonMkNoteHeader parses the Nikon Makernote from reader and returns byteOrder and error
func NikonMkNoteHeader(reader io.Reader) (byteOrder binary.ByteOrder, err error) {
	// Nikon Makernotes header is 18 bytes. Move Reader up necessary bytes
	mknoteHeader := [18]byte{}
	if n, err := reader.Read(mknoteHeader[:]); n < 18 || err != nil {
		return nil, ErrNikonMkNote
	}
	// Nikon makernote header starts with "Nikon" with the first 5 bytes
	if mknote.IsNikonMkNoteHeaderBytes(mknoteHeader[:5]) {
		// Exif Header
		if byteOrder := meta.BinaryOrder(mknoteHeader[10:14]); byteOrder != nil {
			return byteOrder, nil
		}
	}

	return nil, ErrNikonMkNote
}
