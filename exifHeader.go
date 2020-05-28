package exiftool

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNoExif no exif information was found
	ErrNoExif = errors.New("No exif in object")
)

const (
	// Tiff Header Length is 8 bytes
	tiffHeaderLength = 8
)

var (
	// TiffBigEndianSignature is the Tiff Signature for BigEndian encoded images
	TiffBigEndianSignature = []byte{0x4d, 0x4d, 0x00, 0x2a}

	// TiffLittleEndianSignature is the Tiff Signature for LittleEndian encoded images
	TiffLittleEndianSignature = []byte{0x49, 0x49, 0x2a, 0x00}

	// JPEGStartOfImageMarker is the JPEG Start of Image Marker.
	// JPEG SOI Marker
	JPEGStartOfImageMarker = []byte{0xff, 0xd8}

	// JPEGEndofImageMarker is the JPEG End of Image Marker.
	// JPEG EOI Marker
	JPEGEndofImageMarker = []byte{0xff, 0xd9}

	// PNGImageSignature is the marker for the start of a PNG Image.
	// 4 Bytes
	PNGImageSignature = []byte{0x89, 0x50, 0x4E, 0x47}
)

// ExifHeader contains the byte Order, first Ifd Offset,
// tiff Header offset and Image type for the parsing
// of Exif information.
type ExifHeader struct {
	byteOrder        binary.ByteOrder
	firstIfdOffset   uint32
	tiffHeaderOffset uint32
	imageType        ImageType
}

func (eh ExifHeader) String() string {
	str, _ := eh.imageType.MarshalText()
	if eh.byteOrder == binary.BigEndian {
		return fmt.Sprintf("ExifHeader: BigEndian | Tiff offset:  0x%04x | IFD offset: 0x%04x | %s", eh.tiffHeaderOffset, eh.firstIfdOffset, str)
	} else if eh.byteOrder == binary.LittleEndian {
		return fmt.Sprintf("ExifHeader: LittleEndian | Tiff offset: 0x%04x | IFD offset: 0x%04x | %s", eh.tiffHeaderOffset, eh.firstIfdOffset, str)
	}
	return fmt.Sprintf("ExifHeader: empty | %s", str)
}

// SearchExifHeader searches an io.Reader for a LittleEndian Tiff Header or a BigEndian Tiff Header
func SearchExifHeader(reader io.Reader) (eh ExifHeader, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
			//err = log.Wrap(state.(error))
		}
	}()

	// Search for the beginning of the EXIF information. The EXIF is near the
	// beginning of most JPEGs, so this likely doesn't have a high cost (at
	// least, again, with JPEGs).
	br := bufio.NewReader(reader)
	discarded := 0

	var window []byte

	for {
		window, err = br.Peek(tiffHeaderLength)
		if err != nil {
			if err == io.EOF {
				err = ErrNoExif
				return
			}
			panic(err)
		}
		if len(window) < 8 {
			//log.Warningf(nil, "Not enough data for EXIF header: (%d)", len(data))
			err = ErrNoExif
			return
		}
		byteOrder := parseExifHeader(window)
		if byteOrder == nil {

			// Exif not identified. Move forward by one byte.
			if _, err := br.Discard(1); err != nil {
				panic(err)
			}

			discarded++

			continue
		}

		// Found
		eh.byteOrder = byteOrder
		eh.firstIfdOffset = eh.byteOrder.Uint32(window[4:8])
		eh.tiffHeaderOffset = uint32(discarded)
		break
	}

	return eh, nil
}

func parseExifHeader(data []byte) binary.ByteOrder {
	// Good reference:
	//
	//      CIPA DC-008-2016; JEITA CP-3451D
	//      -> http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
	if IsTiffBigEndian(data[:4]) {
		return binary.BigEndian
	} else if IsTiffLittleEndian(data[:4]) {
		return binary.LittleEndian
	}
	return nil
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func IsTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func IsTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}
