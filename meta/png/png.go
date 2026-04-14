// Package png reads PNG Header metadata information from image files before being processed by exif package
package png

import (
	"encoding/binary"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func ScanPngHeader(r io.ReadSeeker) (header meta.ExifHeader, err error) {
	// 5.2 PNG signature
	const signature = "\x89PNG\r\n\x1a\n"

	// 5.3 Chunk layout
	const crcSize = 4

	// 8 is the size of both the signature and the chunk
	// id (4 bytes) + chunk length (4 bytes).
	// This is just a coincidence.
	buf := make([]byte, 8)

	var n int
	n, err = r.Read(buf)
	if err != nil {
		return
	}

	if n != len(signature) || string(buf) != signature {
		err = meta.ErrNoExif

		return
	}

	for {
		// 5.3 Chunk layout
		n, err = r.Read(buf)
		if err != nil {
			break
		}

		if n != len(buf) {
			break
		}

		length := binary.BigEndian.Uint32(buf[0:4])
		chunkType := string(buf[4:8])

		switch chunkType {
		case "eXIf":
			offset, _ := r.Seek(0, io.SeekCurrent)

			return meta.NewExifHeader(utils.BigEndian, 8, uint32(offset), length, imagetype.ImagePNG), nil

		default:
			// Discard the chunk length + CRC.
			_, err := r.Seek(int64(length+crcSize), io.SeekCurrent)
			if err != nil {
				return header, err
			}
		}
	}

	return header, meta.ErrNoExif
}
