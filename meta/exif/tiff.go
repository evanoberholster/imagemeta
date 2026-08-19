// Package tiff reads Tiff Header metadata information from image files before being processed by exif package
package exif

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	// TiffHeaderLength is 8 bytes
	TiffHeaderLength = 32

	// maxTiffOffset caps file offsets to 32 bits. BigTIFF offsets are 64-bit on
	// the wire, but imagemeta stores them as uint32; a container that places an
	// IFD, value or preview beyond 4 GiB is reported as unsupported (ErrNoExif)
	// rather than parsed with a truncated offset.
	maxTiffOffset = 0xFFFFFFFF
)

// ScanTiffHeader searches for the beginning of the EXIF information. The EXIF is near the
// beginning of most Image files, so this likely doesn't have a high cost.
func ScanTiffHeader(r io.Reader, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
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
			// Keep the caller-supplied type unless we can positively detect a better one.
			if detected, detectErr := imagetype.Buf(buf); detectErr == nil && !detected.IsUnknown() {
				it = detected
			}
		}

		byteOrder := utils.BinaryOrder(buf)
		if byteOrder == utils.UnknownEndian {
			// Exif not identified. Move forward by one byte.
			if buf[1] == 0x49 || buf[1] == 0x4d {
				if _, err = br.Discard(1); err != nil {
					return header, err
				}
				discarded++
				continue
			}
			if _, err = br.Discard(2); err != nil {
				return header, err
			}
			discarded += 2
			continue
		}

		// Found Tiff Header
		tiffHeaderOffset, ok := meta.SafecastIntToUint32(discarded)
		if !ok {
			return header, meta.ErrNoExif
		}

		bigTiff := utils.IsBigTiff(buf)
		var firstIfdOffset uint32
		if bigTiff {
			// BigTIFF header: bytes 4-5 offset size (always 8), 6-7 constant 0,
			// 8-15 the 8-byte IFD0 offset. Offsets beyond 4 GiB are unsupported.
			if byteOrder.Uint16(buf[4:6]) != 8 || byteOrder.Uint16(buf[6:8]) != 0 {
				return header, meta.ErrNoExif
			}
			firstIfdOffset64 := byteOrder.Uint64(buf[8:16])
			if firstIfdOffset64 > maxTiffOffset {
				return header, meta.ErrNoExif
			}
			firstIfdOffset = uint32(firstIfdOffset64)
		} else {
			firstIfdOffset = byteOrder.Uint32(buf[4:8])
		}

		header = meta.NewExifHeader(byteOrder, firstIfdOffset, tiffHeaderOffset, 0, it)
		header.FirstIfd = tag.IFD0
		header.BigTIFF = bigTiff
		return header, nil
	}
}
