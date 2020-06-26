package jpegmeta

import (
	"bufio"
	"encoding/binary"

	"github.com/evanoberholster/exiftool/meta/tiffmeta"
)

// jpegByteOrder - JPEG uses a BigEndian Byte Order.
var jpegByteOrder = binary.BigEndian

// SOFHeader contains height, width and number of components.
type SOFHeader struct {
	//offset     uint32
	height     uint16
	width      uint16
	components uint8
}

// Metadata from JPEG files
type Metadata struct {
	sof        SOFHeader
	xml        string
	tiffHeader tiffmeta.Header
	Exif       []byte

	// Reader
	br        *bufio.Reader
	discarded uint32
	pos       uint8
}

// Size returns the width and height of the JPEG Image
func (m Metadata) Size() (width, height uint16) {
	return m.sof.width, m.sof.height
}

// XML returns the xml in the JPEG Image as a string
func (m Metadata) XML() string {
	return m.xml
}

// TiffHeader returns the tiffmeta.Header from the JPEG Image
func (m Metadata) TiffHeader() tiffmeta.Header {
	return m.tiffHeader
}

// newMetadata creates a New metadata object from an io.Reader
func newMetadata(reader *bufio.Reader) Metadata {
	return Metadata{
		br:        reader,
		discarded: 0,
	}
}

// readAPP1
// TODO: Documentation/Testing
func (m *Metadata) discard(i int) (err error) {
	// Exif not identified. Move forward by one byte.
	i, err = m.br.Discard(i)
	m.discarded += uint32(i)
	return
}

// readAPP1
// TODO: Documentation/Testing
func (m *Metadata) readAPP1(buf []byte) (err error) {
	// APP1 XML Marker
	if isXMPPrefix(buf) {
		// ReadXMP reads the XML header/component into metadata
		return m.readXMP(buf)
	}
	// APP1 Exif Marker
	if isJpegExifPrefix(buf) {
		// Read the length of the Exif Information
		length := jpegByteOrder.Uint16(buf[2:4]) - 8

		// Discard App Marker bytes and Exif header bytes
		m.discard(2 + 8)

		// Peek at TiffHeader information
		if buf, err = m.br.Peek(exifPrefixLength); err != nil {
			return err
		}

		// Create a TiffHeader from the Tiff directory ByteOrder, root IFD Offset,
		// the tiff Header Offset, and the length of the exif information.
		byteOrder := tiffmeta.BinaryOrder(buf)
		firstIfdOffset := byteOrder.Uint32(buf[4:8])
		exifLength := uint32(length)
		m.setTiffHeader(tiffmeta.NewHeader(byteOrder, firstIfdOffset, m.discarded, exifLength))

		//fmt.Println("Exif Tiff Header:", m.tiffHeader)

		// Discard Exif information bytes
		return m.discard(int(length))
	}
	return nil
}

// readXMP
// TODO: Documentation/Testing
func (m *Metadata) readXMP(buf []byte) (err error) {
	// Read the length of the XMPHeader
	length := int(jpegByteOrder.Uint16(buf[2:4])) - 2 - xmpPrefixLength

	// Discard App Marker bytes and header length bytes
	if err = m.discard(4 + xmpPrefixLength); err != nil {
		return err
	}

	//fmt.Println("XML Header:", m.discarded, m.pos, length)
	xmpBuf := make([]byte, length)
	for i := 0; i < length; {
		n, err := m.br.Read(xmpBuf[i:])
		if err != nil {
			return err
		}
		i += n
	}

	m.xml = string(xmpBuf)
	//str := strings.Replace(string(xmpBuf), "\n", "", -1)
	//m.XML = strings.Replace(str, "   ", "", -1)
	//m.XML = xmlfmt.FormatXML(string(buf), "\t", "  ")
	return nil
}

// readSOF
// TODO: Documentation/Testing
func (m *Metadata) readSOF(buf []byte) error {
	length := int(jpegByteOrder.Uint16(buf[2:4]))
	header := SOFHeader{
		//m.discarded,
		jpegByteOrder.Uint16(buf[5:7]),
		jpegByteOrder.Uint16(buf[7:9]),
		buf[9]}
	if m.pos == 1 {
		m.sof = header
	}
	return m.discard(length + 2)
}

func (m *Metadata) setTiffHeader(tiffHeader tiffmeta.Header) {
	m.tiffHeader = tiffHeader
}

// ignoreMarker reads the Marker Header length and then
// discards the said marker and its header length
func (m *Metadata) ignoreMarker(buf []byte) error {
	// Read Marker Header Length
	length := int(jpegByteOrder.Uint16(buf[2:4]))

	// Discard Marker Header Length and Marker Length
	return m.discard(length + 2)
}

//
// scanMarkersAlt is an alternative solution for scanning JPEG Markers
//
// NOT USED
func (m *Metadata) scanMarkersAlt(buf []byte) (err error) {
	if isSOIMarker(buf) {
		//fmt.Println("SOI:", jm.discarded, jm.pos)
		m.pos++
		return m.discard(2)
	} else if m.pos > 0 {
		if isDRIMarker(buf) {
			return m.discard(6)
		}
		if isSOFMarker(buf) {
			return m.readSOF(buf)
		}
		if isAPP0Marker(buf) {
			// Ignore APP0 Marker
			return m.ignoreMarker(buf)
		}
		if isAPP1Marker(buf) {
			return m.readAPP1(buf)
		}
		if isAPP2Marker(buf) {
			//if isICCProfilePrefix(buf) {
			// Ignore ICC Profile Marker
			//}
			return m.ignoreMarker(buf)
		}
		if isAPP13Marker(buf) {
			//if isPhotoshopPrefix(buf) {
			// Ignore Photoshop Profile Marker
			//}
			return m.ignoreMarker(buf)
		}
		if isAPP14Marker(buf) {
			return m.ignoreMarker(buf)
		}
		if isDQTMarker(buf) {
			// Ignore DQT Markers
			return m.ignoreMarker(buf)
		}
		if isDHTMarker(buf) {
			// Artifical End Of Image for DHT Marker.
			// This is done to improve performance.
			if m.pos == 1 {
				return ErrEndOfImage
			}
			// Ignore DHT Markers
			return m.ignoreMarker(buf)
		}
		// EOIMarker: End of Image Marker
		if isEOIMarker(buf) {
			m.pos--
			if m.pos == 0 {
				return ErrEndOfImage // Return EndOfImage
			}
			return m.discard(2)
		}
		if isMarker(buf, markerAPP7) {
			// Ignore APP7 Marker
			return m.ignoreMarker(buf)
		}
		if isMarker(buf, markerAPP8) {
			// Ignore APP8 Marker
			return m.ignoreMarker(buf)
		}
		if isMarker(buf, markerAPP9) {
			// Ignore APP9 Marker
			return m.ignoreMarker(buf)
		}
		if isMarker(buf, markerAPP10) {
			// Ignore APP10 Marker
			return m.ignoreMarker(buf)
		}
	}
	return m.discard(1)
}
