// Package jpeg reads metadata information (Exif and XMP) from a JPEG Image.
package jpeg

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/evanoberholster/imagemeta/xmp"
)

// Errors
var (
	ErrNoExif       = meta.ErrNoExif
	ErrNoJPEGMarker = errors.New("no JPEG Marker")
	ErrEndOfImage   = errors.New("end of Image")
)

// Metadata from a JPEG file
type Metadata struct {
	// Decode Functions for EXIF and XMP metadata
	ExifDecodeFn exif.DecodeFn
	XmpDecodeFn  xmp.DecodeFn
	// SOF Header and Tiff Header
	sofHeader
	ExifHeader exif.Header
	XmpHeader  xmp.Header

	// Reader
	br        *bufio.Reader
	discarded uint32
	pos       uint8
}

// ScanJPEG scans a reader for JPEG Image markers. xmpDecodeFn and exifDecodeFn are run at their respective
// positions during the scan. Returns Metadata.
//
// Returns the error ErrNoJPEGMarker if a JPEG SOF was not found.
func ScanJPEG(r *bufio.Reader, xmpDecodeFn xmp.DecodeFn, exifDecodeFn exif.DecodeFn) (m Metadata, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	m = newMetadata(r, xmpDecodeFn, exifDecodeFn)

	var buf []byte
	for {
		if buf, err = m.br.Peek(16); err != nil {
			err = ErrNoJPEGMarker
			return
		}

		if !isMarkerFirstByte(buf) {
			if err = m.discard(1); err != nil {
				return
			}
			continue
		}
		if isSOIMarker(buf) {
			m.pos++
			//fmt.Println("SOI:", m.discarded, m.pos)
			if err = m.discard(2); err != nil {
				return
			}
			continue
		}
		if m.pos > 0 {
			if err := m.scanMarkers(buf); err == nil {
				continue
			}
		}

		break
	}
	if !m.ExifHeader.IsValid() {
		err = ErrNoExif
		return
	}
	return
}

func (m *Metadata) scanMarkers(buf []byte) (err error) {
	switch buf[1] {
	case markerSOF0, markerSOF1,
		markerSOF2, markerSOF3,
		markerSOF5, markerSOF6,
		markerSOF7, markerSOF9,
		markerSOF10:
		return m.readSOF(buf)
	case markerDHT:
		// Artificial End Of Image for DHT Marker.
		// This is done to improve performance.
		if m.pos == 1 {
			return ErrEndOfImage
		}
		// Ignore DHT Markers
		return m.ignoreMarker(buf)
	case markerSOI:
		m.pos++
		return m.discard(2)
	case markerEOI:
		m.pos--
		// Return EndOfImage
		if m.pos == 1 {
			return ErrEndOfImage
		}
		return m.discard(2)
	case markerDQT:
		// Ignore DQT Markers
		return m.ignoreMarker(buf)
	case markerDRI:
		return m.discard(6)
	case markerAPP0:
		return m.ignoreMarker(buf)
	case markerAPP2:
		if isICCProfilePrefix(buf) {
			// Ignore ICC Profile Marker
			return m.ignoreMarker(buf)
		}
		return m.ignoreMarker(buf)
	case markerAPP7, markerAPP8,
		markerAPP9, markerAPP10:
		return m.ignoreMarker(buf)
	case markerAPP13:
		if isPhotoshopPrefix(buf) {
			// Ignore Photoshop Profile Marker
			return m.ignoreMarker(buf)
		}
		return m.ignoreMarker(buf)
	case markerAPP14:
		return m.ignoreMarker(buf)
	case markerAPP1:
		return m.readAPP1(buf)
	}
	//fmt.Println(m.discarded)
	return m.discard(1)
}

// Size returns the width and height of the JPEG Image
func (m Metadata) Size() (width, height uint16) {
	return m.width, m.height
}

// newMetadata creates a New metadata object from an io.Reader
func newMetadata(reader *bufio.Reader, xmpDecodeFn xmp.DecodeFn, exifDecodeFn exif.DecodeFn) Metadata {
	jm := Metadata{
		br:        reader,
		discarded: 0,
	}
	jm.XmpDecodeFn = xmpDecodeFn
	jm.ExifDecodeFn = exifDecodeFn
	return jm
}

// discard adds to m.discarded and discards from the underlying bufio.Reader
func (m *Metadata) discard(i int) (err error) {
	if i == 0 {
		return
	}
	i, err = m.br.Discard(i)
	m.discarded += uint32(i)
	return
}

// readAPP1
func (m *Metadata) readAPP1(buf []byte) (err error) {
	// APP1 XML Marker
	if isXMPPrefix(buf) {
		return m.readXMP(buf)
	}
	// APP1 Exif Marker
	if isJpegExifPrefix(buf) {
		return m.readExif(buf)
	}
	return nil
}

// readExif reads the Exif header/component with the addtached metadata
// ExifDecodeFn. If the function is nil it discards the exif length.
func (m *Metadata) readExif(buf []byte) (err error) {
	// Read the length of the Exif Information
	length := jpegByteOrder.Uint16(buf[2:4]) - exifPrefixLength

	// Discard App Marker bytes and Exif header bytes
	if err = m.discard(2 + exifPrefixLength); err != nil {
		return err
	}

	// Peek at TiffHeader information
	if buf, err = m.br.Peek(exifPrefixLength); err != nil {
		return err
	}

	// Create a TiffHeader from the Tiff directory ByteOrder, root IFD Offset,
	// the tiff Header Offset, and the length of the exif information.
	byteOrder := tiff.BinaryOrder(buf)
	firstIfdOffset := byteOrder.Uint32(buf[4:8])
	exifLength := uint32(length)

	// Set Tiff Header
	m.ExifHeader = exif.NewHeader(byteOrder, firstIfdOffset, m.discarded, exifLength, imagetype.ImageJPEG)

	//fmt.Println("Exif Tiff Header:", m.Header)
	// Read Exif Information
	if m.ExifDecodeFn != nil {
		r := io.LimitReader(m.br, int64(length))
		err = m.ExifDecodeFn(r, m.ExifHeader)
		if err != nil {
			return err
		}
		remain := r.(*io.LimitedReader).N
		return m.discard(int(remain))
	}
	// Discard Exif information bytes
	return m.discard(int(length))
}

// readXMP reads the Exif header/component with the addtached metadata
// XmpDecodeFn. If the function is nil it discards the exif length.
func (m *Metadata) readXMP(buf []byte) (err error) {
	// Read the length of the XMPHeader
	length := int(jpegByteOrder.Uint16(buf[2:4])) - 2 - xmpPrefixLength

	// Discard App Marker bytes and header length bytes
	if err = m.discard(4 + xmpPrefixLength); err != nil {
		return err
	}
	m.XmpHeader = xmp.NewHeader(m.discarded, uint32(length))

	// TODO: XMP Header (offset, length)
	// Use XML Decode Function if not nil
	if m.XmpDecodeFn != nil {
		r := io.LimitReader(m.br, int64(length))
		err = m.XmpDecodeFn(r, m.XmpHeader)
		if err != nil {
			return err
		}
		remain := r.(*io.LimitedReader).N
		return m.discard(int(remain))
	}

	// Discard Xmp information bytes
	return m.discard(int(length))
}

// readSOF reads a JPEG Start of file with the uint16
// width, height, and components of the JPEG image.
func (m *Metadata) readSOF(buf []byte) error {
	length := int(jpegByteOrder.Uint16(buf[2:4]))
	header := sofHeader{
		jpegByteOrder.Uint16(buf[5:7]),
		jpegByteOrder.Uint16(buf[7:9]),
		buf[9]}
	if m.pos == 1 {
		m.sofHeader = header
	}
	return m.discard(length + 2)
}

// ignoreMarker reads the Marker Header length and then
// discards the said marker and its header length
func (m *Metadata) ignoreMarker(buf []byte) error {
	// Read Marker Header Length
	length := int(jpegByteOrder.Uint16(buf[2:4]))

	// Discard Marker Header Length and Marker Length
	return m.discard(length + 2)
}

// Markers refers to the second byte of a JPEG Marker.
// The first is always 0xFF
const (
	markerFirstByte = 0xFF

	// SOF Markers
	markerSOF0  = 0xC0
	markerSOF1  = 0xC1
	markerSOF2  = 0xC2
	markerSOF3  = 0xC3
	markerSOF5  = 0xC5
	markerSOF6  = 0xC6
	markerSOF7  = 0xC7
	markerSOF9  = 0xC9
	markerSOF10 = 0xCA
	markerSOF11 = 0xCB

	// Other Markers
	markerDHT = 0xC4
	markerSOI = 0xD8
	markerEOI = 0xD9
	markerDQT = 0xDB
	markerDRI = 0xDD

	// APP Markers
	markerAPP0  = 0xE0
	markerAPP1  = 0xE1
	markerAPP2  = 0xE2
	markerAPP7  = 0xE7
	markerAPP8  = 0xE8
	markerAPP9  = 0xE9
	markerAPP10 = 0xEA
	markerAPP13 = 0xED
	markerAPP14 = 0xEE
)

// Prefix lengths
const (
	xmpPrefixLength  = 29
	exifPrefixLength = 8
)

// jpegByteOrder JPEG always uses a BigEndian byteorder inside the JPEG image.
// Can use either byteorder for Exif Information inside the JPEG image.
var jpegByteOrder = binary.BigEndian

// sofHeader contains height, width and number of components.
type sofHeader struct {
	height     uint16
	width      uint16
	components uint8
}

// isSOIMarker returns true if the first 2 bytes match an SOI marker
func isSOIMarker(buf []byte) bool {
	return buf[0] == markerFirstByte &&
		buf[1] == markerSOI
}

func isMarkerFirstByte(buf []byte) bool {
	return buf[0] == markerFirstByte
}

// PhotoshopPrefix returns true if
// buf[4:14] equals "Photoshop 3.0\000",
// buf[0:2] is AppMarker, buf[2:4] is HeaderLength
func isPhotoshopPrefix(buf []byte) bool {
	return buf[4] == 0x50 &&
		buf[5] == 0x68 &&
		buf[6] == 0x6f &&
		buf[7] == 0x74 &&
		buf[8] == 0x6f &&
		buf[9] == 0x70 &&
		buf[10] == 0x20 &&
		buf[11] == 0x33 &&
		buf[12] == 0x2e &&
		buf[13] == 0x30 &&
		buf[14] == 0x00
}

// isICCProfilePrefix returns true if
// buf[4:14] equals []byte,
// buf[0:2] is AppMarker, buf[2:4] is HeaderLength
func isICCProfilePrefix(buf []byte) bool {
	return buf[4] == 0x49 &&
		buf[5] == 0x43 &&
		buf[6] == 0x43 &&
		buf[7] == 0x5f &&
		buf[8] == 0x50 &&
		buf[9] == 0x52 &&
		buf[10] == 0x4f &&
		buf[11] == 0x46 &&
		buf[12] == 0x49 &&
		buf[13] == 0x4c &&
		buf[14] == 0x45
}

// isXMPPrefix returns true if
// buf[4:15] equals "http://ns.adobe.com/xap/1.0/\000",
// buf[0:2] is AppMarker, buf[2:4] is HeaderLength
func isXMPPrefix(buf []byte) bool {
	return buf[4] == 0x68 &&
		buf[5] == 0x74 &&
		buf[6] == 0x74 &&
		buf[7] == 0x70 &&
		buf[8] == 0x3a &&
		buf[9] == 0x2f &&
		buf[10] == 0x2f &&
		buf[11] == 0x6e &&
		buf[12] == 0x73 &&
		buf[13] == 0x2e &&
		buf[14] == 0x61 &&
		buf[15] == 0x64
}

// isJpegExifPrefix returns true if
// buf[4:9] equals "Exif" and '0', '0',
// buf[0:2] is AppMarker, buf[2:4] is HeaderLength
func isJpegExifPrefix(buf []byte) bool {
	return buf[4] == 0x45 &&
		buf[5] == 0x78 &&
		buf[6] == 0x69 &&
		buf[7] == 0x66 &&
		buf[8] == 0x00 &&
		buf[9] == 0x00
}
