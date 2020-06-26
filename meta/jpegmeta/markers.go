package jpegmeta

// Markers buf[1]
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

	xmpPrefixLength  = 29
	exifPrefixLength = 8
)

// isSOIMarker returns true if the first 2 bytes match an SOI marker
func isSOIMarker(buf []byte) bool {
	return buf[0] == markerFirstByte &&
		buf[1] == markerSOI
}

func isMarkerFirstByte(buf []byte) bool {
	return buf[0] == markerFirstByte
}

// PhotoshopPrefix returns true if
// buf[4:14] equals []byte{"Photoshop 3.0\000"},
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
// buf[4:14] equals []byte{},
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
// buf[4:15] equals byte{"http://ns.adobe.com/xap/1.0/\000"},
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
// buf[4:9] equals byte{'E', 'x', 'i', 'f', 0, 0},
// buf[0:2] is AppMarker, buf[2:4] is HeaderLength
func isJpegExifPrefix(buf []byte) bool {
	return buf[4] == 0x45 &&
		buf[5] == 0x78 &&
		buf[6] == 0x69 &&
		buf[7] == 0x66 &&
		buf[8] == 0x00 &&
		buf[9] == 0x00
}
