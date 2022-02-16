// Package tiff reads Tiff Header metadata information from image files before being processed by exif package
package tiff

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 16
)

// ScanTiff searches an io.Reader for a LittleEndian or BigEndian Tiff Header
// and returns the ExifHeader with the imagetype ImageTiff.
func ScanTiff(r io.Reader) (meta.ExifHeader, error) {
	return Scan(r, imagetype.ImageTiff)
}

// Scan searches an io.Reader for a LittleEndian or BigEndian Tiff Header
// and returns the ExifHeader with the given imagetype.
func Scan(r io.Reader, it imagetype.ImageType) (meta.ExifHeader, error) {
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < 32 {
		br = bufio.NewReaderSize(r, 32)
	}
	return scan(br, it)
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
		header.FirstIfd = ifds.RootIFD
		return header, nil
	}
}
