// Package png reads PNG Header metadata information from image files before being processed by exif package
package png

import (
	"encoding/binary"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	pngSignatureValue uint64 = 0x89504e470d0a1a0a
	exifChunkType     uint32 = 0x65584966
	pngChunkCRCSize          = 4
)

type scanBuffer [16]byte

var scanBufferPool = sync.Pool{
	New: func() any { return new(scanBuffer) },
}

func acquireScanBuffer() *scanBuffer {
	buf, ok := scanBufferPool.Get().(*scanBuffer)
	if !ok || buf == nil {
		return new(scanBuffer)
	}
	return buf
}

func releaseScanBuffer(buf *scanBuffer) {
	if buf != nil {
		scanBufferPool.Put(buf)
	}
}

func ScanPngHeader(r io.ReadSeeker) (header meta.ExifHeader, err error) {
	buf := acquireScanBuffer()
	defer releaseScanBuffer(buf)

	if _, err = io.ReadFull(r, buf[:8]); err != nil {
		return header, err
	}
	if binary.BigEndian.Uint64(buf[:8]) != pngSignatureValue {
		return header, meta.ErrNoExif
	}

	for {
		if _, err = io.ReadFull(r, buf[:8]); err != nil {
			break
		}

		length := binary.BigEndian.Uint32(buf[:4])
		if binary.BigEndian.Uint32(buf[4:8]) == exifChunkType {
			offset, seekErr := r.Seek(0, io.SeekCurrent)
			if seekErr != nil {
				return header, seekErr
			}
			return scanExifChunkHeader(r, buf, offset, length)
		}

		// Skip the chunk payload and its CRC. Convert before addition to avoid
		// wrapping a maximum uint32 chunk length.
		if _, err = r.Seek(int64(length)+pngChunkCRCSize, io.SeekCurrent); err != nil {
			return header, err
		}
	}

	return header, meta.ErrNoExif
}

// scanExifChunkHeader parses the TIFF header at the start of a PNG eXIf chunk
// and restores the reader to the chunk payload for the EXIF decoder.
func scanExifChunkHeader(r io.ReadSeeker, buf *scanBuffer, offset int64, length uint32) (meta.ExifHeader, error) {
	if length < 8 {
		return meta.ExifHeader{}, meta.ErrNoExif
	}
	if _, err := io.ReadFull(r, buf[:8]); err != nil {
		return meta.ExifHeader{}, err
	}

	byteOrder := utils.BinaryOrder(buf[:8])
	if byteOrder == utils.UnknownEndian {
		return meta.ExifHeader{}, meta.ErrNoExif
	}

	bigTiff := utils.IsBigTiff(buf[:8])
	var firstIFDOffset uint32
	if bigTiff {
		if length < 16 || byteOrder.Uint16(buf[4:6]) != 8 || byteOrder.Uint16(buf[6:8]) != 0 {
			return meta.ExifHeader{}, meta.ErrNoExif
		}
		if _, err := io.ReadFull(r, buf[8:16]); err != nil {
			return meta.ExifHeader{}, err
		}
		firstIFDOffset64 := byteOrder.Uint64(buf[8:16])
		if firstIFDOffset64 > uint64(^uint32(0)) {
			return meta.ExifHeader{}, meta.ErrNoExif
		}
		firstIFDOffset = uint32(firstIFDOffset64)
	} else {
		firstIFDOffset = byteOrder.Uint32(buf[4:8])
	}

	minOffset := uint32(8)
	if bigTiff {
		minOffset = 16
	}
	if firstIFDOffset < minOffset || firstIFDOffset >= length {
		return meta.ExifHeader{}, meta.ErrNoExif
	}

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return meta.ExifHeader{}, err
	}
	offset32, ok := meta.SafecastInt64ToUint32(offset)
	if !ok {
		return meta.ExifHeader{}, meta.ErrNoExif
	}

	header := meta.NewExifHeader(byteOrder, firstIFDOffset, offset32, length, imagetype.ImagePNG)
	header.BigTIFF = bigTiff
	return header, nil
}
