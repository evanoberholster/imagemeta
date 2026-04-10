package exif

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	hoursToSeconds   = 60 * minutesToSeconds
	minutesToSeconds = 60
)

// parseASCIIValueBytes parses ASCII-ish tag payload into a trimmed byte slice.
// Returned bytes reference parser buffers and are only valid until next read.
func (r *Reader) parseASCIIValueBytes(t tag.Entry) []byte {
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		return trimNULBuffer(r.state.buf[:t.Size()])
	}
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return nil
	}
	buf, _, err := r.readTagBytes(t, uint32(len(r.state.buf)))
	if err != nil {
		return nil
	}
	return trimNULBuffer(buf)
}

// parseString parses the requested value from EXIF metadata.
func (r *Reader) parseString(t tag.Entry) string {
	buf := r.parseASCIIValueBytes(t)
	if len(buf) == 0 {
		return ""
	}
	return string(buf)
}

// parseStringTrimRightSpaceNewline parses ASCII EXIF text and removes trailing
// spaces and line breaks before converting to string.
func (r *Reader) parseStringTrimRightSpaceNewline(t tag.Entry) string {
	buf := trimRightSpaceNewline(r.parseASCIIValueBytes(t))
	if len(buf) == 0 {
		return ""
	}
	return string(buf)
}

// parseStringAllowUndefined parses the requested value from EXIF metadata.
func (r *Reader) parseStringAllowUndefined(t tag.Entry) string {
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
		buf := bytes.TrimSpace(r.parseASCIIValueBytes(t))
		if len(buf) == 0 {
			return ""
		}
		return string(buf)
	case tag.TypeUndefined:
	default:
		return ""
	}
	buf := r.parseUndefinedBytes(t, 512)
	if len(buf) == 0 {
		return ""
	}
	// exiftool text dumps render non-printable bytes as '.'. Mirror that for
	// parity checks while still preserving raw bytes separately.
	trimmed := trimNULBuffer(buf)
	if len(trimmed) == 0 {
		return ""
	}
	if len(trimmed) <= 512 {
		var out [512]byte
		for i := 0; i < len(trimmed); i++ {
			b := trimmed[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
				continue
			}
			out[i] = '.'
		}
		sanitized := bytes.TrimSpace(out[:len(trimmed)])
		if len(sanitized) == 0 {
			return ""
		}
		return string(sanitized)
	}
	out := make([]byte, len(trimmed))
	for i := 0; i < len(trimmed); i++ {
		b := trimmed[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
			continue
		}
		out[i] = '.'
	}
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return ""
	}
	return string(out)
}

// parseDisplayString keeps a printable representation that is close to exiftool
// text dumps by converting non-printable bytes to '.' and preserving trailing dots.
func (r *Reader) parseDisplayString(t tag.Entry, maxBytes uint32) string {
	var buf []byte
	switch {
	case t.IsEmbedded():
		if maxBytes == 0 {
			return ""
		}
		n := t.Size()
		if n > maxBytes {
			n = maxBytes
		}
		t.EmbeddedValue(r.state.buf[:4])
		buf = r.state.buf[:n]
	case t.IsType(tag.TypeUndefined), t.IsType(tag.TypeByte), t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		if maxBytes == 0 {
			return ""
		}
		var err error
		buf, _, err = r.readTagBytes(t, maxBytes)
		if err != nil || len(buf) == 0 {
			return ""
		}
	default:
		return ""
	}

	allPrintable := true
	for i := 0; i < len(buf); i++ {
		if buf[i] < 0x20 || buf[i] > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return strings.TrimRight(string(buf), " \t\r\n")
	}

	if len(buf) <= 512 {
		var out [512]byte
		for i := 0; i < len(buf); i++ {
			b := buf[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
				continue
			}
			out[i] = '.'
		}
		return strings.TrimRight(string(out[:len(buf)]), " \t\r\n")
	}

	out := make([]byte, len(buf))
	for i := 0; i < len(buf); i++ {
		b := buf[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
			continue
		}
		out[i] = '.'
	}
	return strings.TrimRight(string(out), " \t\r\n")
}

// parseDisplayStringTrimRightSpaceNewline is like parseDisplayString but trims
// trailing spaces and line breaks on the byte slice before converting to string.
func (r *Reader) parseDisplayStringTrimRightSpaceNewline(t tag.Entry, maxBytes uint32) string {
	var buf []byte
	switch {
	case t.IsEmbedded():
		if maxBytes == 0 {
			return ""
		}
		n := t.Size()
		if n > maxBytes {
			n = maxBytes
		}
		t.EmbeddedValue(r.state.buf[:4])
		buf = r.state.buf[:n]
	case t.IsType(tag.TypeUndefined), t.IsType(tag.TypeByte), t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		if maxBytes == 0 {
			return ""
		}
		var err error
		buf, _, err = r.readTagBytes(t, maxBytes)
		if err != nil || len(buf) == 0 {
			return ""
		}
	default:
		return ""
	}

	buf = trimRightSpaceNewline(buf)
	if len(buf) == 0 {
		return ""
	}

	allPrintable := true
	for i := 0; i < len(buf); i++ {
		if buf[i] < 0x20 || buf[i] > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return string(buf)
	}

	if len(buf) <= 512 {
		var out [512]byte
		for i := 0; i < len(buf); i++ {
			b := buf[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
				continue
			}
			out[i] = '.'
		}
		outBuf := trimRightSpaceNewline(out[:len(buf)])
		if len(outBuf) == 0 {
			return ""
		}
		return string(outBuf)
	}

	out := make([]byte, len(buf))
	for i := 0; i < len(buf); i++ {
		b := buf[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
			continue
		}
		out[i] = '.'
	}
	out = trimRightSpaceNewline(out)
	if len(out) == 0 {
		return ""
	}
	return string(out)
}

func trimRightSpaceNewline(buf []byte) []byte {
	end := len(buf)
	if end > 0 && buf[end-1] == 0 {
		end--
	}
	if end > 0 && buf[end-1] == 0 {
		end--
	}
	for end > 0 {
		switch buf[end-1] {
		case ' ', '\n', '\r':
			end--
		default:
			return buf[:end]
		}
	}
	return buf[:0]
}

// parseUndefinedBytes parses the requested UNDEFINED value from EXIF metadata.
// Returned bytes reference parser buffers and are only valid until next read.
func (r *Reader) parseUndefinedBytes(t tag.Entry, maxBytes uint32) []byte {
	if !t.IsType(tag.TypeUndefined) {
		return nil
	}
	return r.parseOpaqueBytes(t, maxBytes)
}

// parseOpaqueBytes parses raw bytes for byte-like EXIF metadata without copying.
// Returned bytes reference parser buffers and are only valid until next read.
// Tag must be of type TypeUndefined, TypeByte, TypeASCII, or TypeASCIINoNul.
func (r *Reader) parseOpaqueBytes(t tag.Entry, maxBytes uint32) []byte {
	switch t.Type {
	case tag.TypeUndefined, tag.TypeByte, tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return nil
	}
	if maxBytes == 0 {
		return nil
	}
	if t.IsEmbedded() {
		n := min(t.Size(), maxBytes)
		t.EmbeddedValue(r.state.buf[:4])
		return r.state.buf[:n]
	}
	buf, _, err := r.readTagBytes(t, maxBytes)
	if err != nil || len(buf) == 0 {
		return nil
	}
	return buf
}

// parseUint16 parses the requested value from EXIF metadata.
func (r *Reader) parseUint16(t tag.Entry) uint16 {
	if !t.IsEmbedded() {
		return 0
	}
	switch t.Type {
	case tag.TypeShort:
		return t.EmbeddedShort()
	case tag.TypeLong:
		return uint16(t.EmbeddedLong())
	default:
		return 0
	}
}

// parseUint32 parses the requested value from EXIF metadata.
func (r *Reader) parseUint32(t tag.Entry) uint32 {
	if !t.IsEmbedded() {
		return 0
	}
	switch t.Type {
	case tag.TypeLong:
		return t.EmbeddedLong()
	case tag.TypeShort:
		return uint32(t.EmbeddedShort())
	default:
		return 0
	}
}

// parseRationalU parses the requested value from EXIF metadata.
func (r *Reader) parseRationalU(t tag.Entry) [2]uint32 {
	switch t.Type {
	case tag.TypeRational, tag.TypeSignedRational:
	default:
		return [2]uint32{}
	}
	buf, _, err := r.readTagBytes(t, 8)
	if err != nil || len(buf) < 8 {
		return [2]uint32{}
	}
	return [2]uint32{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
}

// parseRationalValue parses the requested value from EXIF metadata.
func (r *Reader) parseRationalValue(t tag.Entry) tag.RationalU {
	parts := r.parseRationalU(t)
	return tag.RationalU{Numerator: parts[0], Denominator: parts[1]}
}

// parseUint16List parses the requested value from EXIF metadata.
func (r *Reader) parseUint16List(t tag.Entry, dst []uint16) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeShort, tag.TypeSignedShort:
	default:
		return 0
	}
	n := min(int(t.UnitCount), len(dst))
	if t.IsEmbedded() {
		return t.EmbeddedShorts(dst[:n])
	}
	buf, _, err := r.readTagBytes(t, uint32(n*2))
	if err != nil {
		return 0
	}
	if got := len(buf) / 2; got < n {
		n = got
	}

	if t.ByteOrder == utils.BigEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+2 {
			dst[i] = binary.BigEndian.Uint16(buf[j:])
		}
		return n
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+2 {
		dst[i] = binary.LittleEndian.Uint16(buf[j:])
	}
	return n
}

// parseUndefinedUint16List parses uint16 values from UNDEFINED or SHORT payloads.
func (r *Reader) parseUndefinedUint16List(t tag.Entry, dst []uint16) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeShort:
		return r.parseUint16List(t, dst)
	case tag.TypeUndefined:
	default:
		return 0
	}

	n := int(t.UnitCount / 2)
	if n > len(dst) {
		n = len(dst)
	}
	if n == 0 {
		return 0
	}

	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		if n > 2 {
			n = 2
		}
		if t.ByteOrder == utils.BigEndian {
			for i, j := 0, 0; i < n; i, j = i+1, j+2 {
				dst[i] = binary.BigEndian.Uint16(r.state.buf[j:])
			}
			return n
		}
		for i, j := 0, 0; i < n; i, j = i+1, j+2 {
			dst[i] = binary.LittleEndian.Uint16(r.state.buf[j:])
		}
		return n
	}

	buf, _, err := r.readTagBytes(t, uint32(n*2))
	if err != nil {
		return 0
	}
	if got := len(buf) / 2; got < n {
		n = got
	}
	if t.ByteOrder == utils.BigEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+2 {
			dst[i] = binary.BigEndian.Uint16(buf[j:])
		}
		return n
	}
	for i, j := 0, 0; i < n; i, j = i+1, j+2 {
		dst[i] = binary.LittleEndian.Uint16(buf[j:])
	}
	return n
}

// parseUint32List parses the requested value from EXIF metadata.
func (r *Reader) parseUint32List(t tag.Entry, dst []uint32) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeLong, tag.TypeShort, tag.TypeIfd:
	default:
		return 0
	}
	n := min(int(t.UnitCount), len(dst))
	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeLong, tag.TypeIfd:
			dst[0] = t.EmbeddedLong()
		case tag.TypeShort:
			var shorts [2]uint16
			m := min(t.EmbeddedShorts(shorts[:]), n)
			dst[0] = uint32(shorts[0])
			if m == 2 {
				dst[1] = uint32(shorts[1])
			}
			return m
		}
		return 0
	}
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd:
		buf, _, err := r.readTagBytes(t, uint32(n*4))
		if err != nil {
			return 0
		}
		if got := len(buf) / 4; got < n {
			n = got
		}
		for i := 0; i < n; i++ {
			start := i * 4
			dst[i] = t.ByteOrder.Uint32(buf[start : start+4])
		}
	case tag.TypeShort:
		buf, _, err := r.readTagBytes(t, uint32(n*2))
		if err != nil {
			return 0
		}
		if got := len(buf) / 2; got < n {
			n = got
		}
		for i := 0; i < n; i++ {
			start := i * 2
			dst[i] = uint32(t.ByteOrder.Uint16(buf[start : start+2]))
		}
	default:
		return 0
	}
	return n
}

// parseInt32List parses the requested value from EXIF metadata.
func (r *Reader) parseInt32List(t tag.Entry, dst []int32) int {
	if len(dst) == 0 || !t.IsType(tag.TypeSignedLong) || t.UnitCount == 0 {
		return 0
	}
	n := min(int(t.UnitCount), len(dst))
	if t.IsEmbedded() {
		dst[0] = int32(t.EmbeddedLong())
		return 1
	}
	buf, _, err := r.readTagBytes(t, uint32(n*4))
	if err != nil {
		return 0
	}
	if got := len(buf) / 4; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 4
		dst[i] = int32(t.ByteOrder.Uint32(buf[start : start+4]))
	}
	return n
}

// parseByteList parses the requested value from EXIF metadata.
func (r *Reader) parseByteList(t tag.Entry, dst []byte) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeByte, tag.TypeUndefined, tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return 0
	}

	n := min(int(t.UnitCount), len(dst))
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		copy(dst[:n], r.state.buf[:n])
		return n
	}
	buf, _, err := r.readTagBytes(t, uint32(n))
	if err != nil {
		return 0
	}
	if len(buf) < n {
		n = len(buf)
	}
	copy(dst[:n], buf[:n])
	return n
}

// parseRationalSList parses the requested value from EXIF metadata.
func (r *Reader) parseRationalSList(t tag.Entry, dst []int32) int {
	if len(dst) < 2 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeSignedRational, tag.TypeRational:
	default:
		return 0
	}
	n := min(int(t.UnitCount), len(dst)/2)
	if n == 0 {
		return 0
	}
	buf, _, err := r.readTagBytes(t, uint32(n*8))
	if err != nil {
		return 0
	}
	if got := len(buf) / 8; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 8
		dst[i*2] = int32(t.ByteOrder.Uint32(buf[start : start+4]))
		dst[i*2+1] = int32(t.ByteOrder.Uint32(buf[start+4 : start+8]))
	}
	return n
}

// parseDate parses the requested value from EXIF metadata.
func (r *Reader) parseDate(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 19 {
		return time.Time{}
	}
	// YYYY:MM:DD HH:MM:SS
	if buf[4] != ':' || buf[7] != ':' || buf[10] != ' ' || buf[13] != ':' || buf[16] != ':' {
		return time.Time{}
	}
	return time.Date(
		int(parseStrUint(buf[0:4])),
		time.Month(parseStrUint(buf[5:7])),
		int(parseStrUint(buf[8:10])),
		int(parseStrUint(buf[11:13])),
		int(parseStrUint(buf[14:16])),
		int(parseStrUint(buf[17:19])),
		0,
		time.UTC,
	)
}

// parseOffsetTime parses the requested value from EXIF metadata.
func (r *Reader) parseOffsetTime(t tag.Entry) *time.Location {
	if !t.IsType(tag.TypeASCII) {
		return nil
	}
	buf, _, err := r.readTagBytes(t, 8)
	if err != nil || len(buf) < 6 {
		return nil
	}
	if buf[3] != ':' {
		return nil
	}
	offset := int(parseStrUint(buf[1:3]))*hoursToSeconds + int(parseStrUint(buf[4:6]))*minutesToSeconds
	switch buf[0] {
	case '-':
		return getLocation(int32(offset*-1), buf[:6])
	case '+':
		return getLocation(int32(offset), buf[:6])
	default:
		return nil
	}
}

var (
	exifTextASCII   = []byte("ASCII\x00\x00\x00")
	exifTextUnicode = []byte("UNICODE\x00")
	exifTextJIS     = []byte("JIS\x00\x00\x00\x00\x00")
	exifTextEmpty   = []byte{0, 0, 0, 0, 0, 0, 0, 0}
)

// parseExifUserComment decodes EXIF UserComment (0x9286) using the EXIF text
// header conventions that ExifTool also recognizes.
func (r *Reader) parseExifUserComment(t tag.Entry) string {
	switch t.Type {
	case tag.TypeUndefined, tag.TypeByte, tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return ""
	}
	buf := r.parseOpaqueBytes(t, min(t.Size(), 4096))
	if len(buf) == 0 {
		return ""
	}
	if len(buf) < 8 {
		return exifASCIIText(bytes.TrimSpace(trimNULBuffer(buf)))
	}

	header := buf[:8]
	payload := buf[8:]
	switch {
	case bytes.Equal(header, exifTextASCII), bytes.Equal(header, exifTextJIS), bytes.Equal(header, exifTextEmpty):
		return exifASCIIText(payload)
	case bytes.Equal(header, exifTextUnicode):
		return exifUTF16Text(payload, t.ByteOrder)
	default:
		// Fall back to best-effort text decoding for malformed headers.
		return exifASCIIText(buf)
	}
}

// parseAperture parses the requested value from EXIF metadata.
func (r *Reader) parseAperture(t tag.Entry) meta.Aperture {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return meta.Aperture(float32(rat[0]) / float32(rat[1]))
}

// parseApexAperture parses an APEX aperture value and converts it to an F-number.
func (r *Reader) parseApexAperture(t tag.Entry) meta.Aperture {
	switch {
	case t.IsType(tag.TypeSignedRational):
		var rat [2]int32
		if r.parseRationalSList(t, rat[:]) == 0 || rat[1] == 0 {
			return 0
		}
		return apexApertureToFNumber(float64(rat[0]) / float64(rat[1]))
	case t.IsType(tag.TypeRational):
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return 0
		}
		return apexApertureToFNumber(float64(rat[0]) / float64(rat[1]))
	default:
		return 0
	}
}

// parseExposureTime parses the requested value from EXIF metadata.
func (r *Reader) parseExposureTime(t tag.Entry) meta.ExposureTime {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return meta.ExposureTime(float32(rat[0]) / float32(rat[1]))
}

// parseShutterSpeed parses the requested value from EXIF metadata.
func (r *Reader) parseShutterSpeed(t tag.Entry) meta.ShutterSpeed {
	return apexShutterSpeedToSeconds(r.parseSignedRationalFloat64(t))
}

// parseSignedRationalFloat32 parses a rational (signed or unsigned) as float32.
func (r *Reader) parseSignedRationalFloat32(t tag.Entry) float32 {
	return float32(r.parseSignedRationalFloat64(t))
}

// parseSignedRationalFloat64 parses a rational (signed or unsigned) as float64.
func (r *Reader) parseSignedRationalFloat64(t tag.Entry) float64 {
	switch {
	case t.IsType(tag.TypeSignedRational):
		var rat [2]int32
		if r.parseRationalSList(t, rat[:]) == 0 {
			return 0
		}
		if rat[1] == 0 {
			return 0
		}
		return float64(rat[0]) / float64(rat[1])
	case t.IsType(tag.TypeRational):
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return 0
		}
		return float64(rat[0]) / float64(rat[1])
	default:
		return 0
	}
}

// parseExposureBias parses the requested value from EXIF metadata.
func (r *Reader) parseExposureBias(t tag.Entry) meta.ExposureBias {
	switch {
	case t.IsType(tag.TypeSignedRational):
		var rat [2]int32
		if r.parseRationalSList(t, rat[:]) == 0 || rat[1] == 0 {
			return meta.NewExposureBias(0, 0)
		}
		n, d := rat[0], rat[1]
		if d < 0 {
			n = -n
			d = -d
		}
		return meta.NewExposureBias(int16(n), int16(d))
	case t.IsType(tag.TypeRational):
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return meta.NewExposureBias(0, 0)
		}
		return meta.NewExposureBias(int16(rat[0]), int16(rat[1]))
	default:
		return meta.NewExposureBias(0, 0)
	}
}

// parseFocalLength parses the requested value from EXIF metadata.
func (r *Reader) parseFocalLength(t tag.Entry) meta.FocalLength {
	switch t.Type {
	case tag.TypeShort, tag.TypeLong:
		return meta.NewFocalLength(r.parseUint32(t), 1)
	case tag.TypeRational, tag.TypeSignedRational:
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return 0
		}
		return meta.FocalLength(float32(rat[0]) / float32(rat[1]))
	default:
		return 0
	}
}

// parseLensInfo parses the requested value from EXIF metadata.
func (r *Reader) parseLensInfo(t tag.Entry) *LensInfo {
	if t.IsEmbedded() {
		return nil
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 32 {
		return nil
	}
	return &LensInfo{
		MinFocalLength:        tag.RationalU{Numerator: t.ByteOrder.Uint32(buf[:4]), Denominator: t.ByteOrder.Uint32(buf[4:8])},
		MaxFocalLength:        tag.RationalU{Numerator: t.ByteOrder.Uint32(buf[8:12]), Denominator: t.ByteOrder.Uint32(buf[12:16])},
		MaxApertureAtMinFocal: tag.RationalU{Numerator: t.ByteOrder.Uint32(buf[16:20]), Denominator: t.ByteOrder.Uint32(buf[20:24])},
		MaxApertureAtMaxFocal: tag.RationalU{Numerator: t.ByteOrder.Uint32(buf[24:28]), Denominator: t.ByteOrder.Uint32(buf[28:32])},
	}
}

// parseSceneType parses Exif SceneType from BYTE/UNDEFINED/ASCII or integer encodings.
func (r *Reader) parseSceneType(t tag.Entry) uint16 {
	switch {
	case t.IsType(tag.TypeShort), t.IsType(tag.TypeLong):
		return uint16(r.parseUint32(t))
	case t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		s := strings.TrimSpace(r.parseString(t))
		if s == "" {
			return 0
		}
		return uint16(parseStrUint([]byte(s)))
	case t.IsType(tag.TypeByte), t.IsType(tag.TypeUndefined):
		var b [4]byte
		if r.parseByteList(t, b[:]) == 0 {
			return 0
		}
		return uint16(b[0])
	default:
		return 0
	}
}

// parseGPSRef parses and normalizes GPS reference tags.
func (r *Reader) parseGPSRef(t tag.Entry) tag.GPSRef {
	first := byte(0)
	ok := false
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		first = r.state.buf[0]
		ok = true
	} else {
		buf, _, err := r.readTagBytes(t, 1)
		if err == nil && len(buf) > 0 {
			first = buf[0]
			ok = true
		}
	}
	if !ok {
		return tag.GPSRefUnknown
	}

	switch t.ID {
	case tag.TagGPSAltitudeRef:
		// ExifTool documents altitude ref values:
		// 0 above sea level, 1 below, 2 above sea-level reference, 3 below.
		if first == 1 || first == 3 {
			return tag.GPSRefBelowSeaLevel
		}
		if first == 0 || first == 2 {
			return tag.GPSRefAboveSeaLevel
		}
		return tag.GPSRefUnknown
	case tag.TagGPSLatitudeRef, tag.TagGPSDestLatitudeRef:
		if first == 'S' || first == 's' {
			return tag.GPSRefSouth
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefNorth
		}
		return tag.GPSRefUnknown
	case tag.TagGPSLongitudeRef, tag.TagGPSDestLongitudeRef:
		if first == 'W' || first == 'w' {
			return tag.GPSRefWest
		}
		if first == 'E' || first == 'e' {
			return tag.GPSRefEast
		}
		return tag.GPSRefUnknown
	case tag.TagGPSSpeedRef:
		if first == 'K' || first == 'k' {
			return tag.GPSRefKilometersPerHour
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMilesPerHour
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefKnots
		}
		return tag.GPSRefUnknown
	case tag.TagGPSTrackRef, tag.TagGPSImgDirectionRef, tag.TagGPSDestBearingRef:
		if first == 'T' || first == 't' {
			return tag.GPSRefTrueDirection
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMagneticDirection
		}
		return tag.GPSRefUnknown
	case tag.TagGPSDestDistanceRef:
		if first == 'K' || first == 'k' {
			return tag.GPSRefKilometers
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMiles
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefNauticalMiles
		}
		return tag.GPSRefUnknown
	}
	return tag.GPSRefUnknown
}

// parseGPSCoord parses the requested value from EXIF metadata.
func (r *Reader) parseGPSCoord(t tag.Entry) float64 {
	if t.UnitCount != 3 {
		return 0
	}
	if !(t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	dNum := t.ByteOrder.Uint32(buf[:4])
	dDen := t.ByteOrder.Uint32(buf[4:8])
	mNum := t.ByteOrder.Uint32(buf[8:12])
	mDen := t.ByteOrder.Uint32(buf[12:16])
	sNum := t.ByteOrder.Uint32(buf[16:20])
	sDen := t.ByteOrder.Uint32(buf[20:24])
	if dDen == 0 || mDen == 0 || sDen == 0 {
		return 0
	}
	coord := (float64(dNum) / float64(dDen))
	coord += (float64(mNum) / float64(mDen) / 60.0)
	coord += (float64(sNum) / float64(sDen) / 3600.0)
	return coord
}

// parseGPSAltitude parses the requested value from EXIF metadata.
func (r *Reader) parseGPSAltitude(t tag.Entry) float32 {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return float32(rat[0]) / float32(rat[1])
}

// parseGPSTimeStamp parses the requested value from EXIF metadata.
func (r *Reader) parseGPSTimeStamp(t tag.Entry) time.Duration {
	if t.UnitCount != 3 || !t.IsType(tag.TypeRational) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	v := [6]uint32{
		t.ByteOrder.Uint32(buf[:4]),
		t.ByteOrder.Uint32(buf[4:8]),
		t.ByteOrder.Uint32(buf[8:12]),
		t.ByteOrder.Uint32(buf[12:16]),
		t.ByteOrder.Uint32(buf[16:20]),
		t.ByteOrder.Uint32(buf[20:24]),
	}
	return rationalDuration(v[0], v[1], time.Hour) +
		rationalDuration(v[2], v[3], time.Minute) +
		rationalDuration(v[4], v[5], time.Second)
}

// parseGPSDateStamp parses the requested value from EXIF metadata.
func (r *Reader) parseGPSDateStamp(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 10 {
		return time.Time{}
	}
	if len(buf) >= 19 && buf[10] == ' ' && buf[13] == ':' && buf[16] == ':' {
		return time.Date(
			int(parseStrUint(buf[0:4])),
			time.Month(parseStrUint(buf[5:7])),
			int(parseStrUint(buf[8:10])),
			int(parseStrUint(buf[11:13])),
			int(parseStrUint(buf[14:16])),
			int(parseStrUint(buf[17:19])),
			0,
			time.UTC,
		)
	}
	if buf[4] == ':' && buf[7] == ':' {
		return time.Date(
			int(parseStrUint(buf[0:4])),
			time.Month(parseStrUint(buf[5:7])),
			int(parseStrUint(buf[8:10])),
			0, 0, 0, 0, time.UTC,
		)
	}
	return time.Time{}
}

// TODO(simd): accelerate trim/ASCII scan and offset-date parsing with simd/archsimd
// when GOEXPERIMENT=simd integration is prioritized.
