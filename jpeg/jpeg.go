// Copyright (c) 2018-2023 Evan Oberholster. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// Package jpeg reads metadata information (Exif and XMP) from a JPEG Image.
package jpeg

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// Errors
var (
	ErrNoExif       = meta.ErrNoExif
	ErrNoJPEGMarker = errors.New("no JPEG Marker")
	ErrEndOfImage   = errors.New("end of Image")
)

const (
	bufferSize int = 4 * 1024 // 4Kb
)

type jpegReader struct {
	ExifReader func(r io.Reader, h meta.ExifHeader) error
	XMPReader  func(r io.Reader) error

	// Reader
	br  *bufio.Reader
	err error

	// SOF Header
	sofHeader

	// Marker
	buf    []byte
	offset uint32
	size   uint16
	marker markerType

	// Reader
	pos       uint8
	discarded uint32
}

var bufferPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, bufferSize) },
}

// ScanJPEG scans a reader for JPEG Image markers. exifReader and xmpReader are run at their respective
// positions during the scan. Returns en error.
//
// Returns the error ErrNoJPEGMarker if a JPEG SOF was not found.
func ScanJPEG(r io.Reader, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	var localBuffer bool
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() <= bufferSize {
		localBuffer = true
		br = bufferPool.Get().(*bufio.Reader)
		br.Reset(r)
	}

	jr := &jpegReader{br: br, ExifReader: exifReader, XMPReader: xmpReader}

	defer func() {
		if localBuffer {
			bufferPool.Put(jr.br)
		}
	}()

	for jr.nextMarker() {
		switch jr.marker {
		case markerSOF0, markerSOF1,
			markerSOF2, markerSOF3,
			markerSOF5, markerSOF6,
			markerSOF7, markerSOF9,
			markerSOF10:
			if logInfo() {
				jr.logMarker("")
			}
			jr.err = jr.readSOF(jr.buf)
		case markerDHT:
			if logInfo() {
				jr.logMarker("")
			}
			// Artificial End Of Image for DHT Marker.
			// This is done to improve performance.
			if jr.pos == 1 {
				return nil
			}
			// Ignore DHT Markers
			jr.ignoreMarker()
		case markerSOI:
			if logInfo() {
				jr.logMarker("")
			}
			jr.pos++
			jr.err = jr.discard(2)
		case markerEOI:
			if logInfo() {
				jr.logMarker("")
			}
			jr.pos--
			// Return EndOfImage
			if jr.pos == 1 {
				return ErrEndOfImage
			}
			jr.err = jr.discard(2)
		case markerDQT:
			// Ignore DQT Markers and close parsing
			// Stop parsing at DQT Markers
			//return nil
			if logInfo() {
				jr.logMarker("")
			}
			jr.ignoreMarker()
			return nil
		case markerDRI:
			return jr.discard(6)
		case markerAPP0:
			jr.readAPP0()
		case markerAPP1:
			jr.readAPP1()
		case markerAPP2:
			jr.readAPP2()
		case markerAPP13:
			jr.readAPP13()
		default:
			if logInfo() {
				jr.logMarker("")
			}
			jr.ignoreMarker()
		}
	}
	return jr.err
}

func (jr *jpegReader) nextMarker() bool {
	for jr.err == nil {
		if jr.buf, jr.err = jr.peek(64); jr.err != nil {
			jr.err = ErrNoJPEGMarker
			return false
		}
		if !isMarkerFirstByte(jr.buf) {
			var i int
			for i = 0; i < 64; i++ {
				if isMarkerFirstByte(jr.buf[i:]) {
					break
				}
			}
			jr.err = jr.discard(i)
			continue
		}

		if isSOIMarker(jr.buf) {
			jr.pos++
			jr.err = jr.discard(2)
			continue
		}
		if jr.pos > 0 {
			jr.offset = jr.discarded
			jr.size = jpegEndian.Uint16(jr.buf[2:4])
			jr.marker = markerType(jr.buf[1])
			return true
		}
	}
	return false
}

// peek returns the next n bytes without advancing the unerlying bufio.Reader
func (jr *jpegReader) peek(n int) ([]byte, error) {
	return jr.br.Peek(n)
}

// discard adds to m.discarded and discards from the underlying bufio.Reader
func (jr *jpegReader) discard(i int) (err error) {
	if i == 0 {
		return
	}
	i, err = jr.br.Discard(i)
	jr.discarded += uint32(i)
	return
}

// readAPP0
func (jr *jpegReader) readAPP0() {
	// Is JFIF Marker
	if isJFIFPrefix(jr.buf) || isJFIFPrefixExt(jr.buf) {
		if logInfo() {
			jr.logMarker("APP0 JFIF")
		}
	}
	jr.ignoreMarker()
}

// readAPP1
func (jr *jpegReader) readAPP1() {
	// APP1 Exif Marker
	if isExifPrefix(jr.buf) {
		if logInfo() {
			jr.logMarker("APP1 Exif")
		}
		jr.err = jr.readExif()
		return
	}

	// APP1 XMP Marker
	if isXMPPrefix(jr.buf) {
		if logInfo() {
			jr.logMarker("APP1 XMP")
		}
		jr.err = jr.readXMP()
		return
	}

	// APP1 XMP Extension marker (NOT SUPPORTED)
	if isXMPPrefixExt(jr.buf) {
		if logInfo() {
			jr.logMarker("APP1 XMP Extension")
		}
		// Ignore XMP Extension
	}
	jr.ignoreMarker()
}

// readAPP2
func (jr *jpegReader) readAPP2() {
	if isICCProfilePrefix(jr.buf) {
		if logInfo() {
			jr.logMarker("APP2 ICC Profile")
		}
		// Ignore ICC Profile Marker
	}
	jr.ignoreMarker()
}

// readAPP13
func (jr *jpegReader) readAPP13() {
	if isPhotoshopPrefix(jr.buf) {
		if logInfo() {
			jr.logMarker("APP13 Photoshop")
		}
		// Ignore Photoshop Profile Marker
	}
	jr.ignoreMarker()
}

// readExif reads the Exif header/component with the addtached metadata
// ExifDecodeFn. If the function is nil it discards the exif length.
func (jr *jpegReader) readExif() (err error) {
	var buf []byte
	// Read the length of the Exif Information
	remain := int(jr.size) - exifPrefixLength

	// Discard App Marker bytes and Exif header bytes
	if err = jr.discard(2 + exifPrefixLength); err != nil {
		return err
	}

	// Peek at TiffHeader information
	if buf, err = jr.peek(exifPrefixLength); err != nil {
		return err
	}

	// Read Exif
	if jr.ExifReader != nil {
		// Create a TiffHeader from the Tiff directory ByteOrder, root IFD Offset,
		// the tiff Header Offset, and the length of the exif information.
		byteOrder := utils.BinaryOrder(buf)
		firstIfdOffset := byteOrder.Uint32(buf[4:8])
		exifLength := uint32(remain)

		// Set Tiff Header
		exifHeader := meta.NewExifHeader(byteOrder, firstIfdOffset, jr.discarded, exifLength, imagetype.ImageJPEG)

		if err = jr.ExifReader(jr.br, exifHeader); err != nil {
			return err
		}
		// Discard remaining bytes
		remain = 0
	}

	// Discard remaining bytes
	return jr.discard(remain)
}

// readXMP reads the Exif header/component with the addtached metadata
// XmpDecodeFn. If the function is nil it discards the exif length.
func (jr *jpegReader) readXMP() (err error) {
	// Read the length of the XMPHeader
	remain := int(jr.size) - 2 - xmpPrefixLength

	// Discard App Marker bytes and header length bytes
	if err = jr.discard(4 + xmpPrefixLength); err != nil {
		return err
	}
	// Read XMP Decode Function here
	if jr.XMPReader != nil {
		r := io.LimitReader(jr.br, int64(remain))
		if err = jr.XMPReader(r); err != nil {
			return err
		}
		// Discard remaining bytes
		remain = int(r.(*io.LimitedReader).N)
	}
	// Discard remaining bytes
	return jr.discard(remain)
}

// readSOF reads a JPEG Start of file with the uint16
// width, height, and components of the JPEG image.
func (jr *jpegReader) readSOF(buf []byte) error {
	height := jpegEndian.Uint16(buf[5:7])
	width := jpegEndian.Uint16(buf[7:9])
	comp := uint8(buf[9])
	if jr.pos == 1 {
		jr.sofHeader = sofHeader{height, width, comp}
	}
	return jr.discard(int(jr.size) + 2)
}

// sofHeader contains height, width and number of components.
type sofHeader struct {
	height     uint16
	width      uint16
	components uint8
}

// ignoreMarker discards the marker size
func (jr *jpegReader) ignoreMarker() {
	jr.err = jr.discard(int(jr.size) + 2)
}

// markerType refers to the second byte of a JPEG Marker.
// The first is always 0xFF
type markerType uint8

const (
	markerFirstByte markerType = 0xFF

	// SOF Markers
	markerSOF0  markerType = 0xC0
	markerSOF1  markerType = 0xC1
	markerSOF2  markerType = 0xC2
	markerSOF3  markerType = 0xC3
	markerSOF5  markerType = 0xC5
	markerSOF6  markerType = 0xC6
	markerSOF7  markerType = 0xC7
	markerSOF9  markerType = 0xC9
	markerSOF10 markerType = 0xCA
	markerSOF11 markerType = 0xCB

	// Other Markers
	markerDHT       markerType = 0xC4
	markerSOI       markerType = 0xD8
	markerEOI       markerType = 0xD9
	markerImageData markerType = 0xD9
	markerDQT       markerType = 0xDB
	markerDRI       markerType = 0xDD

	// APP Markers
	markerAPP0  markerType = 0xE0
	markerAPP1  markerType = 0xE1
	markerAPP2  markerType = 0xE2
	markerAPP7  markerType = 0xE7
	markerAPP8  markerType = 0xE8
	markerAPP9  markerType = 0xE9
	markerAPP10 markerType = 0xEA
	markerAPP13 markerType = 0xED
	markerAPP14 markerType = 0xEE

	// Prefixes for JPEG markers
	exifPrefix       = "Exif\000\000"
	jfifPrefix       = "JFIF\000"
	jfifPrefixExt    = "JFXX\000"
	iccPrefix        = "ICC_PROFILE"
	xmpPrefix        = "http://ns.adobe.com/xap/1.0/\000"
	xmpPrefixExt     = "http://ns.adobe.com/xmp/extension/"
	photoshopPrefix  = "Photoshop "
	exifPrefixLength = 8
	xmpPrefixLength  = 29
)

var (
	// jpegEndian JPEG and JFIF always use BigEndian byteorder.
	// Can use either byteorder for Exif Information inside the JPEG image.
	jpegEndian = binary.BigEndian

	mapMarkerTypeString = map[markerType]string{
		markerSOF0:  "SOF0",
		markerSOF1:  "SOF1",
		markerSOF2:  "SOF2",
		markerSOF3:  "SOF3",
		markerSOF5:  "SOF5",
		markerSOF6:  "SOF6",
		markerSOF7:  "SOF7",
		markerSOF9:  "SOF9",
		markerSOF10: "SOF10",
		markerSOF11: "SOF11",
		markerDHT:   "DHT",
		markerSOI:   "SOI",
		markerEOI:   "EOI",
		markerDQT:   "DQT",
		markerDRI:   "DRI",
		markerAPP0:  "APP0",
		markerAPP1:  "APP1",
		markerAPP2:  "APP2",
		markerAPP7:  "APP7",
		markerAPP8:  "APP8",
		markerAPP9:  "APP9",
		markerAPP10: "APP10",
		markerAPP13: "APP13",
		markerAPP14: "APP14",
	}
)

// String is a Stringer interface for markerType
// returns the string name of the second byte of a JPEG marker
func (mt markerType) String() string {
	str, ok := mapMarkerTypeString[mt]
	if ok {
		return str
	}
	return fmt.Sprintf("Unknown marker %x", uint8(mt))
}

// isSOIMarker returns true if the first 2 bytes match an SOI marker
func isSOIMarker(buf []byte) bool {
	return isMarkerFirstByte(buf) &&
		buf[1] == byte(markerSOI)
}

// isMarkerFirstByte returns true if the first byte matches a marker
func isMarkerFirstByte(buf []byte) bool {
	return buf[0] == byte(markerFirstByte)
}

// PhotoshopPrefix returns true if marker matches photoshopPrefix
func isPhotoshopPrefix(buf []byte) bool {
	return string(buf[4:14]) == photoshopPrefix
}

// isICCProfilePrefix returns true if marker matches iccPrefix
func isICCProfilePrefix(buf []byte) bool {
	return string(buf[4:15]) == iccPrefix
}

// isXMPPrefix returns true if marker matches xmpPrefix
func isXMPPrefix(buf []byte) bool {
	return string(buf[4:33]) == xmpPrefix
}

// isXMPPrefixExt returns true if marker matches xmpPrefixExt
func isXMPPrefixExt(buf []byte) bool {
	return string(buf[4:38]) == xmpPrefixExt
}

// isJpegExifPrefix returns true if marker matches exifPrefix
func isExifPrefix(buf []byte) bool {
	return string(buf[4:10]) == exifPrefix
}

// isJFIFPrefix returns true if marker matches jfifPrefix
func isJFIFPrefix(buf []byte) bool {
	return string(buf[4:9]) == jfifPrefix
}

// isJFIFPrefixExt returns true of marker matches "JFXX"
func isJFIFPrefixExt(buf []byte) bool {
	return string(buf[4:9]) == jfifPrefixExt
}
