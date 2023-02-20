// Package tiff reads Tiff Header metadata information from image files before being processed by exif package
package tiff

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 32

	bufReaderSize = 32
)

// ScanTiffHeader searches for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func ScanTiffHeader(r io.Reader, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok && br.Size() > bufReaderSize {
		br = bufio.NewReader(r)
	}
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

		byteOrder := utils.BinaryOrder(buf)
		if byteOrder == utils.UnknownEndian {
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
