// Package tiff reads Tiff Header metadata information from image files before being processed by exif package
package tiff

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 32

	bufReaderSize = 32
)

// Scan scans a reader for Tiff Image markers then xmpDecodeFn and exifDecodeFn are run at their respective
// positions during the scan. Returns an error.
func Scan(r io.Reader, it imagetype.ImageType, exifFn func(r io.Reader, header meta.ExifHeader) error, xmpFn func(r io.Reader, header meta.XmpHeader) error) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	br := bufio.NewReaderSize(r, bufReaderSize)

	exifHeader, err := scan(br, it)
	if err != nil {
		return errors.Wrap(err, "Scan error")
	}
	return exifFn(br, exifHeader)
}

// ScanTiffHeader searches an io.Reader for a LittleEndian or BigEndian Tiff Header
// and returns the ExifHeader with the given imagetype.
func ScanTiffHeader(r io.Reader, it imagetype.ImageType) (meta.ExifHeader, error) {
	return scan(bufio.NewReaderSize(r, bufReaderSize), it)
}

// scan searches for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func scan(br *bufio.Reader, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	discarded := 0

	var buf []byte

	for {
		if buf, err = br.Peek(TiffHeaderLength); err != nil {
			err = meta.ErrNoExif
			return
		}
		if discarded == 0 {
			it, _ = imagetype.Buf(buf)
		}

		byteOrder := meta.BinaryOrder(buf)
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
		header = meta.NewExifHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0, it)
		header.FirstIfd = ifds.IFD0
		return header, nil
	}
}
