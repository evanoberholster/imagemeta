package jpeg

import (
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// readAPP0 handles APP0 JFIF/JFXX markers.
func (jr *jpegReader) readAPP0() {
	// Is JFIF Marker
	if isJFIFPrefix(jr.buf) || isJFIFPrefixExt(jr.buf) {
		jr.logMarker("APP0 JFIF")
		if jr.metadata != nil && isJFIFPrefix(jr.buf) {
			payload, ok := jr.readMetadataPayload()
			if !ok {
				return
			}
			jr.metadata.JFIF, jr.err = parseJFIF(payload)
			return
		}
	}
	if jr.metadata != nil {
		payload, ok := jr.readMetadataPayload()
		if !ok {
			return
		}
		if isCIFFPayload(payload) {
			jr.logMarker("APP0 CIFF")
			jr.metadata.CIFF, jr.err = ParseCIFF(payload)
		}
		return
	}
	jr.ignoreMarker()
}

// readAPP1 handles APP1 Exif, XMP, and extended XMP markers.
func (jr *jpegReader) readAPP1() {
	// APP1 Exif Marker
	if isExifPrefix(jr.buf) {
		jr.logMarker("APP1 Exif")
		jr.err = jr.readExif()
		return
	}

	// APP1 XMP Marker
	if isXMPPrefix(jr.buf) {
		jr.logMarker("APP1 XMP")
		jr.err = jr.readXMP()
		return
	}

	if isXMPPrefixExt(jr.buf) {
		jr.logMarker("APP1 XMP Extension")
		jr.err = jr.readExtendedXMP()
		return
	}
	jr.ignoreMarker()
}

// readAPP2 handles APP2 markers.
func (jr *jpegReader) readAPP2() {
	if isICCProfilePrefix(jr.buf) {
		jr.logMarker("APP2 ICC Profile")
		if jr.metadata != nil {
			payload, ok := jr.readMetadataPayload()
			if !ok {
				return
			}
			jr.err = jr.metadata.addICCChunk(payload)
			return
		}
	}
	if isMPFPrefix(jr.buf) {
		jr.logMarker("APP2 MPF")
		if jr.metadata != nil {
			payload, ok := jr.readMetadataPayload()
			if !ok {
				return
			}
			jr.metadata.MPF, jr.err = parseMPF(payload, uint32(jr.offset)+8)
			return
		}
	}
	jr.ignoreMarker()
}

// readAPP13 handles APP13 Photoshop markers.
func (jr *jpegReader) readAPP13() {
	if isPhotoshopPrefix(jr.buf) {
		jr.logMarker("APP13 Photoshop")
		if jr.metadata != nil {
			payload, ok := jr.readMetadataPayload()
			if !ok {
				return
			}
			jr.metadata.Photoshop, jr.metadata.IPTC, jr.err = parsePhotoshop(payload)
			return
		}
	}
	jr.ignoreMarker()
}

func (jr *jpegReader) readAPP14() {
	if isAdobePrefix(jr.buf) {
		jr.logMarker("APP14 Adobe")
		if jr.metadata != nil {
			payload, ok := jr.readMetadataPayload()
			if !ok {
				return
			}
			jr.metadata.Adobe, jr.err = parseAdobe(payload)
			return
		}
	}
	jr.ignoreMarker()
}

// readAPPMarker reads an APP JPEG Marker
func (jr *jpegReader) readAPPMarker() {
	switch jr.marker {
	case markerAPP0:
		jr.readAPP0()
	case markerAPP1:
		jr.readAPP1()
	case markerAPP2:
		jr.readAPP2()
	case markerAPP13:
		jr.readAPP13()
	case markerAPP14:
		jr.readAPP14()
	default:
		jr.logMarker("")
		jr.ignoreMarker()
	}
}

func (jr *jpegReader) readMetadataPayload() ([]byte, bool) {
	payload, err := jr.readSegmentPayload()
	if err != nil {
		jr.err = err
		return nil, false
	}
	return payload, true
}

func (jr *jpegReader) readCallbackPayload(remain int, cb func(io.Reader) error) (int, error) {
	if cb == nil {
		return remain, nil
	}
	if jr.readerAt != nil {
		sr := io.NewSectionReader(jr.readerAt, int64(jr.discarded), int64(remain))
		return remain, cb(sr)
	}
	lr := utils.NewLimitedBufferedReader(jr.br, remain)
	err := cb(lr)
	consumed := remain - lr.N
	jr.discarded += uint32(consumed)
	return lr.N, err
}

// readExif reads the Exif header/component with the attached metadata
// ExifDecodeFn. If the function is nil it discards the Exif segment.
func (jr *jpegReader) readExif() (err error) {
	// Read the length of the Exif Information
	remain := int(jr.size) - exifPrefixLength
	if remain < tiffHeaderLength {
		return io.ErrUnexpectedEOF
	}

	// Discard App Marker bytes and Exif header bytes
	if err = jr.discard(2 + exifPrefixLength); err != nil {
		return err
	}

	// Peek at TiffHeader information
	buf, err := jr.peek(tiffHeaderLength)
	if err != nil {
		return err
	}
	byteOrder := utils.BinaryOrder(buf)
	firstIfdOffset := byteOrder.Uint32(buf[4:8])
	exifLength := uint32(remain)
	exifHeader := meta.NewExifHeader(byteOrder, firstIfdOffset, jr.discarded, exifLength, imagetype.ImageJPEG)

	// Read Exif
	if jr.ExifReader != nil {
		payloadLength := remain
		remain, err = jr.readCallbackPayload(remain, func(r io.Reader) error {
			return jr.ExifReader(r, exifHeader)
		})
		if err != nil {
			return err
		}
		jr.logDecodedItem("exif", payloadLength)
	}

	// Discard remaining bytes
	return jr.discard(remain)
}

// readXMP reads the XMP packet with the attached metadata XMPDecodeFn.
// If the function is nil it discards the XMP segment.
func (jr *jpegReader) readXMP() (err error) {
	// Read the length of the XMPHeader
	remain := int(jr.size) - 2 - xmpPrefixLength
	if remain < 0 {
		return io.ErrUnexpectedEOF
	}

	// Discard App Marker bytes and header length bytes
	if err = jr.discard(4 + xmpPrefixLength); err != nil {
		return err
	}
	// Read XMP Decode Function here
	if jr.XMPReader != nil {
		payloadLength := remain
		remain, err = jr.readCallbackPayload(remain, jr.XMPReader)
		if err != nil {
			return err
		}
		jr.logDecodedItem("xmp", payloadLength)
	}
	// Discard remaining bytes
	return jr.discard(remain)
}
