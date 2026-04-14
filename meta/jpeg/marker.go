package jpeg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
	markerSOF8  markerType = 0xC8
	markerSOF9  markerType = 0xC9
	markerSOF10 markerType = 0xCA
	markerSOF11 markerType = 0xCB

	// Other Markers
	markerDHT       markerType = 0xC4
	markerSOI       markerType = 0xD8
	markerEOI       markerType = 0xD9
	markerSOS       markerType = 0xDA
	markerImageData markerType = 0xDA
	markerDQT       markerType = 0xDB
	markerDRI       markerType = 0xDD

	// APP Markers
	markerAPP0  markerType = 0xE0
	markerAPP1  markerType = 0xE1
	markerAPP2  markerType = 0xE2
	markerAPP3  markerType = 0xE3
	markerAPP4  markerType = 0xE4
	markerAPP5  markerType = 0xE5
	markerAPP6  markerType = 0xE6
	markerAPP7  markerType = 0xE7
	markerAPP8  markerType = 0xE8
	markerAPP9  markerType = 0xE9
	markerAPP10 markerType = 0xEA
	markerAPP11 markerType = 0xEB
	markerAPP12 markerType = 0xEC
	markerAPP13 markerType = 0xED
	markerAPP14 markerType = 0xEE
	markerAPP15 markerType = 0xEF

	// Prefixes for JPEG markers
	exifPrefix       = "Exif\000\000"
	jfifPrefix       = "JFIF\000"
	jfifPrefixExt    = "JFXX\000"
	iccPrefix        = "ICC_PROFILE"
	mpfPrefix        = "MPF\000"
	adobePrefix      = "Adobe"
	xmpPrefix        = "http://ns.adobe.com/xap/1.0/\000"
	xmpPrefixExt     = "http://ns.adobe.com/xmp/extension/\000"
	photoshopPrefix  = "Photoshop 3.0\000"
	exifPrefixLength = 8
	xmpPrefixLength  = 29
	tiffHeaderLength = 8
	xmpExtHeaderLen  = 75
	maxExtendedXMP   = 64 * 1024 * 1024
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
		markerSOF8:  "SOF8",
		markerSOF9:  "SOF9",
		markerSOF10: "SOF10",
		markerSOF11: "SOF11",
		markerDHT:   "DHT",
		markerSOI:   "SOI",
		markerEOI:   "EOI",
		markerSOS:   "SOS",
		markerDQT:   "DQT",
		markerDRI:   "DRI",
		markerAPP0:  "APP0",
		markerAPP1:  "APP1",
		markerAPP2:  "APP2",
		markerAPP3:  "APP3",
		markerAPP4:  "APP4",
		markerAPP5:  "APP5",
		markerAPP6:  "APP6",
		markerAPP7:  "APP7",
		markerAPP8:  "APP8",
		markerAPP9:  "APP9",
		markerAPP10: "APP10",
		markerAPP11: "APP11",
		markerAPP12: "APP12",
		markerAPP13: "APP13",
		markerAPP14: "APP14",
		markerAPP15: "APP15",
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

func isSOFMarker(marker markerType) bool {
	return (marker&0xF0) == 0xC0 && (marker == markerSOF0 || marker&0x03 != 0)
}

func isAPPMarker(marker markerType) bool {
	return marker >= markerAPP0 && marker <= markerAPP15
}

func markerHasNoLength(marker markerType) bool {
	return marker == markerSOI ||
		marker == markerEOI ||
		marker == 0x01 ||
		(marker >= 0xD0 && marker <= 0xD7)
}

// PhotoshopPrefix returns true if marker matches photoshopPrefix
func isPhotoshopPrefix(buf []byte) bool {
	return len(buf) >= 18 && bytes.Equal(buf[4:18], []byte(photoshopPrefix))
}

// isICCProfilePrefix returns true if marker matches iccPrefix
func isICCProfilePrefix(buf []byte) bool {
	return len(buf) >= 15 && bytes.Equal(buf[4:15], []byte(iccPrefix))
}

func isMPFPrefix(buf []byte) bool {
	return len(buf) >= 8 && bytes.Equal(buf[4:8], []byte(mpfPrefix))
}

func isAdobePrefix(buf []byte) bool {
	return len(buf) >= 9 && bytes.Equal(buf[4:9], []byte(adobePrefix))
}

// isXMPPrefix returns true if marker matches xmpPrefix
func isXMPPrefix(buf []byte) bool {
	return len(buf) >= 33 && bytes.Equal(buf[4:33], []byte(xmpPrefix))
}

// isXMPPrefixExt returns true if marker matches xmpPrefixExt
func isXMPPrefixExt(buf []byte) bool {
	return len(buf) >= 39 && bytes.Equal(buf[4:39], []byte(xmpPrefixExt))
}

// isJpegExifPrefix returns true if marker matches exifPrefix
func isExifPrefix(buf []byte) bool {
	return len(buf) >= 10 && bytes.Equal(buf[4:10], []byte(exifPrefix))
}

// isJFIFPrefix returns true if marker matches jfifPrefix
func isJFIFPrefix(buf []byte) bool {
	return len(buf) >= 9 && bytes.Equal(buf[4:9], []byte(jfifPrefix))
}

// isJFIFPrefixExt returns true of marker matches "JFXX"
func isJFIFPrefixExt(buf []byte) bool {
	return len(buf) >= 9 && bytes.Equal(buf[4:9], []byte(jfifPrefixExt))
}

func isJFIFPayload(buf []byte) bool {
	return len(buf) >= 5 && bytes.Equal(buf[:5], []byte(jfifPrefix))
}

func isCIFFPayload(buf []byte) bool {
	return len(buf) >= 14 &&
		(bytes.Equal(buf[:2], []byte("II")) || bytes.Equal(buf[:2], []byte("MM"))) &&
		bytes.Equal(buf[6:14], []byte("HEAPJPGM"))
}

func isICCPayload(buf []byte) bool {
	return len(buf) >= 14 && bytes.Equal(buf[:12], []byte(iccPrefix+"\000"))
}

func isMPFPayload(buf []byte) bool {
	return len(buf) >= 4 && bytes.Equal(buf[:4], []byte(mpfPrefix))
}

func isPhotoshopPayload(buf []byte) bool {
	return len(buf) >= len(photoshopPrefix) && bytes.Equal(buf[:len(photoshopPrefix)], []byte(photoshopPrefix))
}

func isAdobePayload(buf []byte) bool {
	return len(buf) >= 5 && bytes.Equal(buf[:5], []byte(adobePrefix))
}
