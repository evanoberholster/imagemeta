package exif

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/mknote"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// Errors
var (
	ErrNikonMkNote = errors.New("error makernote is not a Nikon makernote")
)

const (
	// Length of Nikon Makernote Header in bytes
	lengthMkNoteHeaderNikon = 18
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

// isNikonMkNoteHeader parses the Nikon Makernote from reader and returns byteOrder and error
func (r *reader) isNikonMkNoteHeader(ifd ifds.Ifd) (ifds.Ifd, binary.ByteOrder, error) {
	// Nikon Makernotes header is 18 bytes. Move Reader up necessary bytes
	mknoteHeader, err := r.ReadBufferAt(lengthMkNoteHeaderNikon, int(ifd.Offset))
	if err != nil {
		return ifd, nil, errors.Wrapf(err, "error NikonMkNoteHeader at IFD %s", ifd.String())
	}
	// Nikon makernote header starts with "Nikon" with the first 5 bytes
	if mknote.IsNikonMkNoteHeaderBytes(mknoteHeader[:5]) {
		// Exif Header
		if byteOrder := meta.BinaryOrder(mknoteHeader[10:14]); byteOrder != nil {
			ifd.Offset += lengthMkNoteHeaderNikon
			return ifd, byteOrder, nil
		}
	}

	return ifd, nil, ErrNikonMkNote
}

func (r *reader) parseMknoteIFD(e *Data, ifd ifds.Ifd) (ifds.Ifd, binary.ByteOrder) {
	if e.make == "" {
		return ifd, nil
	}
	//make := strings.ToUpper(e.make)
	if e.make == "Canon" {
		// Canon Makernotes do not have a Makernote Header
		// offset 0
		// ByteOrder is the same as RootIfd
		return ifd, r.byteOrder
	}
	if e.make == "Nikon" || e.make == "NIKON CORPORATION" {
		// Nikon v3 maker note is a self-contained Ifd
		// (offsets are relative to the start of the maker note)

		if ifd, byteOrder, err := r.isNikonMkNoteHeader(ifd); err == nil {
			// update imagetype
			if e.imageType == imagetype.ImageTiff {
				e.imageType = imagetype.ImageNEF
			}
			return ifd, byteOrder
		}
		return ifd, nil
	}
	if e.make == "Sony" {
		if e.imageType == imagetype.ImageTiff {
			e.imageType = imagetype.ImageARW
		}
		return ifd, r.byteOrder
	}

	return ifd, nil
}
