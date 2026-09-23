package jpeg

import (
	"bytes"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// readAPP0 handles APP0 JFIF/JFXX markers.
func (jr *jpegReader) readAPP0() {
	// Is JFIF Marker
	jfif := isJFIFPrefix(jr.buf)
	if jfif || isJFIFPrefixExt(jr.buf) {
		jr.logMarker("APP0 JFIF")
		if jr.metadata != nil && jfif {
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
			jr.metadata.MPF, jr.err = parseMPF(payload, jr.offset+8)
			return
		}
	}
	jr.ignoreMarker()
}

// readAPP13 handles APP13 Photoshop markers.
func (jr *jpegReader) readAPP13() {
	if !isPhotoshopPrefix(jr.buf) {
		jr.ignoreMarker()
		return
	}
	jr.logMarker("APP13 Photoshop")
	if jr.metadata == nil {
		jr.ignoreMarker()
		return
	}
	compact, thumbLen, thumbSeen := jr.selectPhotoshopPayload()
	p, iptc, err := parsePhotoshop(compact)
	if err != nil {
		jr.err = err
		return
	}
	if thumbSeen {
		if p == nil {
			p = &Photoshop{}
		}
		p.PhotoshopThumbnailLength = thumbLen
	}
	jr.metadata.Photoshop, jr.metadata.IPTC = p, iptc
}

// wantPhotoshopResource reports whether a resource ID is parsed (its data
// kept) or skipped. 0x040c thumbnails only contribute their length, handled
// separately by the walker.
func wantPhotoshopResource(id uint16) bool {
	switch id {
	case 0x03ed, 0x0404, 0x0406, 0x040a, 0x040b, 0x040c, 0x0425:
		return true
	}
	return false
}

// selectPhotoshopPayload walks the resource directory and compacts only the
// parsed resources, skipping bulky payloads such as thumbnails in place. It
// returns the compact payload plus the thumbnail length recorded separately.
// Declared sizes are bounded by the segment so a lying size cannot consume
// the next marker; truncation keeps already-parsed resources.
func (jr *jpegReader) selectPhotoshopPayload() (compact []byte, thumbLen uint32, thumbSeen bool) {
	budget := int(jr.size) + 2 - 4 - len(photoshopPrefix)
	if err := jr.discard(4 + len(photoshopPrefix)); err != nil {
		return nil, 0, false
	}
	// Pre-size for the prefix plus typical small resources; IPTC-heavy
	// segments grow once more instead of per resource.
	compact = make([]byte, 0, 256)
	compact = append(compact, photoshopPrefix...)
	for budget > 0 {
		head, err := jr.peek(7)
		if err != nil {
			break
		}
		if !bytes.Equal(head[0:4], []byte("8BIM")) && !bytes.Equal(head[0:4], []byte("8B64")) {
			break
		}
		id := jpegEndian.Uint16(head[4:6])
		nameBytes := 1 + int(head[6])
		if nameBytes&1 != 0 {
			nameBytes++
		}
		headerLen := 4 + 2 + nameBytes + 4
		if headerLen > budget {
			break
		}
		head, err = jr.peek(headerLen)
		if err != nil {
			break
		}
		size := int(jpegEndian.Uint32(head[headerLen-4:]))
		dataPad := size & 1
		if size+dataPad > budget-headerLen {
			break
		}
		if err := jr.discard(headerLen); err != nil {
			break
		}
		budget -= headerLen
		if id == 0x040c {
			thumbLen, thumbSeen = thumbnailLength(size), true
			if err := jr.discard(size + dataPad); err != nil {
				break
			}
			budget -= size + dataPad
			continue
		}
		if !wantPhotoshopResource(id) {
			if err := jr.discard(size + dataPad); err != nil {
				break
			}
			budget -= size + dataPad
			continue
		}
		data := make([]byte, size)
		if _, err := io.ReadFull(jr.br, data); err != nil {
			break
		}
		if delta, ok := meta.SafecastIntToUint32(size); ok {
			jr.discarded += delta
		}
		compact = append(compact, head...)
		compact = append(compact, data...)
		if dataPad == 1 {
			compact = append(compact, 0)
			if err := jr.discard(1); err != nil {
				break
			}
			budget -= size + dataPad
			continue
		}
		budget -= size
	}
	return compact, thumbLen, thumbSeen
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
	if delta, ok := meta.SafecastIntToUint32(consumed); ok {
		jr.discarded += delta
	}
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
	exifLength, ok := meta.SafecastIntToUint32(remain)
	if !ok {
		return io.ErrUnexpectedEOF
	}
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
		jr.foundExif = true
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
