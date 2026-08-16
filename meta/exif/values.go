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
	maxBytes, ok := meta.SafecastIntToUint32(len(r.state.buf))
	if !ok {
		return nil
	}
	buf, err := r.readTagBytes(t, maxBytes)
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
	return string(trimAtNUL(buf))
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
		buf := trimASCIIWhitespace(trimAtNUL(r.parseASCIIValueBytes(t)))
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
		sanitized := trimASCIIWhitespace(out[:len(trimmed)])
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
	out = trimASCIIWhitespace(out)
	if len(out) == 0 {
		return ""
	}
	return string(out)
}

type ExifVersion [4]byte

var (
	ExifVersionUnknown = ExifVersion{}
	ExifVersion0200    = ExifVersion{'0', '2', '0', '0'}
	ExifVersion0210    = ExifVersion{'0', '2', '1', '0'}
	ExifVersion0220    = ExifVersion{'0', '2', '2', '0'}
	ExifVersion0221    = ExifVersion{'0', '2', '2', '1'}
	ExifVersion0230    = ExifVersion{'0', '2', '3', '0'}
	ExifVersion0231    = ExifVersion{'0', '2', '3', '1'}
	ExifVersion0232    = ExifVersion{'0', '2', '3', '2'}
	ExifVersion0300    = ExifVersion{'0', '3', '0', '0'}
)

func (v ExifVersion) String() string {
	switch v {
	case ExifVersion0200:
		return "0200"
	case ExifVersion0210:
		return "0210"
	case ExifVersion0220:
		return "0220"
	case ExifVersion0221:
		return "0221"
	case ExifVersion0230:
		return "0230"
	case ExifVersion0231:
		return "0231"
	case ExifVersion0232:
		return "0232"
	case ExifVersion0300:
		return "0300"
	default:
		key := [...]byte{v[0], v[1], v[2], v[3]}
		return string(key[:])
	}
}

// parseExifVersion parses EXIF version bytes without allocating.
//
// ExifVersion is a fixed 4-byte value in practice, so we can decode the common
// values directly from the raw bytes and return string literals.
func (r *Reader) parseExifVersion(t tag.Entry) ExifVersion {
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
		if buf := trimNULBuffer(r.parseASCIIValueBytes(t)); len(buf) == 4 {
			return [4]byte{buf[0], buf[1], buf[2], buf[3]}
		}

	case tag.TypeUndefined:
		if buf := trimNULBuffer(r.parseUndefinedBytes(t, 4)); len(buf) == 4 {
			return [4]byte{buf[0], buf[1], buf[2], buf[3]}
		}
	}
	return ExifVersionUnknown
}

type displayTrimMode uint8

const (
	displayTrimRightSpaces displayTrimMode = iota
	displayTrimRightSpaceNewlineMode
)

type exifTextEncoding uint8

const (
	exifTextEncodingUnknown exifTextEncoding = iota
	exifTextEncodingASCII
	exifTextEncodingUnicode
	exifTextEncodingJIS
)

// parseDisplayString keeps a printable representation that is close to exiftool
// text dumps by converting non-printable bytes to '.' and preserving trailing dots.
func (r *Reader) parseDisplayString(t tag.Entry, maxBytes uint32) string {
	return r.parseDisplayStringWithMode(t, maxBytes, displayTrimRightSpaces)
}

// parseDisplayStringTrimRightSpaceNewline is like parseDisplayString but trims
// trailing spaces and line breaks on the byte slice before converting to string.
func (r *Reader) parseDisplayStringTrimRightSpaceNewline(t tag.Entry, maxBytes uint32) string {
	return r.parseDisplayStringWithMode(t, maxBytes, displayTrimRightSpaceNewlineMode)
}

func (r *Reader) parseDisplayStringWithMode(t tag.Entry, maxBytes uint32, mode displayTrimMode) string {
	buf := r.parseDisplayBytes(t, maxBytes)
	if len(buf) == 0 {
		return ""
	}
	if mode == displayTrimRightSpaceNewlineMode {
		buf = trimRightSpaceNewline(buf)
		if len(buf) == 0 {
			return ""
		}
	}

	allPrintable := true
	for i := 0; i < len(buf); i++ {
		if buf[i] < 0x20 || buf[i] > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		if mode == displayTrimRightSpaceNewlineMode {
			return string(buf)
		}
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
		if mode == displayTrimRightSpaceNewlineMode {
			return string(out[:len(buf)])
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
	if mode == displayTrimRightSpaceNewlineMode {
		return string(out)
	}
	return strings.TrimRight(string(out), " \t\r\n")
}

func (r *Reader) parseDisplayBytes(t tag.Entry, maxBytes uint32) []byte {
	switch {
	case t.IsEmbedded():
		if maxBytes == 0 {
			return nil
		}
		n := t.Size()
		if n > maxBytes {
			n = maxBytes
		}
		t.EmbeddedValue(r.state.buf[:4])
		return r.state.buf[:n]
	case t.IsType(tag.TypeUndefined), t.IsType(tag.TypeByte), t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		if maxBytes == 0 {
			return nil
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil || len(buf) == 0 {
			return nil
		}
		return buf
	default:
		return nil
	}
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
	buf, err := r.readTagBytes(t, maxBytes)
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
		value, ok := meta.SafecastUint32ToUint16(t.EmbeddedLong())
		if !ok {
			return 0
		}
		return value
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

// parseMakerNoteUint32 parses a maker-note integer value from common numeric and
// byte-like representations used by multiple vendors.
func (r *Reader) parseMakerNoteUint32(t tag.Entry) uint32 {
	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeLong, tag.TypeIfd:
			return t.EmbeddedLong()
		case tag.TypeShort, tag.TypeSignedShort:
			return uint32(t.EmbeddedShort())
		}
	}
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd, tag.TypeShort:
		var dst [2]uint32
		if n := r.parseUint32List(t, dst[:]); n > 0 {
			return dst[0]
		}
	case tag.TypeSignedShort:
		var dst [1]uint16
		if n := r.parseUint16List(t, dst[:]); n > 0 {
			return uint32(dst[0])
		}
	case tag.TypeByte, tag.TypeUndefined, tag.TypeASCII, tag.TypeASCIINoNul:
		var dst [4]byte
		if n := r.parseByteList(t, dst[:]); n > 0 {
			return uint32(dst[0])
		}
	}
	return 0
}

func (r *Reader) parseMakerNoteInt16(t tag.Entry) int16 {
	switch t.Type {
	case tag.TypeSignedShort, tag.TypeShort:
	default:
		return 0
	}
	var raw [1]uint16
	if r.parseUint16List(t, raw[:]) == 0 {
		return 0
	}
	return meta.SafecastUint16ToInt16Bits(raw[0])
}

// parseFirstUint32 parses the first uint32-compatible value from EXIF metadata.
func (r *Reader) parseFirstUint32(t tag.Entry) uint32 {
	var raw [1]uint32
	if r.parseUint32List(t, raw[:]) == 0 {
		return 0
	}
	return raw[0]
}

// parseRationalU parses the requested value from EXIF metadata.
func (r *Reader) parseRationalU(t tag.Entry) [2]uint32 {
	switch t.Type {
	case tag.TypeRational, tag.TypeSignedRational:
	default:
		return [2]uint32{}
	}
	buf, err := r.readTagBytes(t, 8)
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

func (r *Reader) parseUnsignedRationalFloat64(t tag.Entry) (float64, unsignedRationalState) {
	rat := r.parseRationalU(t)
	switch {
	case rat[1] == 0 && rat[0] == 0:
		return 0, rationalStateMissing
	case rat[1] == 0:
		return 0, rationalStateInfinite
	default:
		return float64(rat[0]) / float64(rat[1]), rationalStateValue
	}
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
	maxBytes, ok := meta.SafecastIntToUint32(n * 2)
	if !ok {
		return 0
	}
	buf, err := r.readTagBytes(t, maxBytes)
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
		maxBytes, ok := meta.SafecastIntToUint32(n * 4)
		if !ok {
			return 0
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil {
			return 0
		}
		if got := len(buf) / 4; got < n {
			n = got
		}
		if t.ByteOrder == utils.LittleEndian {
			for i, j := 0, 0; i < n; i, j = i+1, j+4 {
				dst[i] = binary.LittleEndian.Uint32(buf[j:])
			}
		} else {
			for i := 0; i < n; i++ {
				start := i * 4
				dst[i] = t.ByteOrder.Uint32(buf[start : start+4])
			}
		}
	case tag.TypeShort:
		maxBytes, ok := meta.SafecastIntToUint32(n * 2)
		if !ok {
			return 0
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil {
			return 0
		}
		if got := len(buf) / 2; got < n {
			n = got
		}
		if t.ByteOrder == utils.LittleEndian {
			for i, j := 0, 0; i < n; i, j = i+1, j+2 {
				dst[i] = uint32(binary.LittleEndian.Uint16(buf[j:]))
			}
		} else {
			for i := 0; i < n; i++ {
				start := i * 2
				dst[i] = uint32(t.ByteOrder.Uint16(buf[start : start+2]))
			}
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
		dst[0] = meta.SafecastUint32ToInt32Bits(t.EmbeddedLong())
		return 1
	}
	maxBytes, ok := meta.SafecastIntToUint32(n * 4)
	if !ok {
		return 0
	}
	buf, err := r.readTagBytes(t, maxBytes)
	if err != nil {
		return 0
	}
	if got := len(buf) / 4; got < n {
		n = got
	}
	if t.ByteOrder == utils.LittleEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+4 {
			dst[i] = meta.SafecastUint32ToInt32Bits(binary.LittleEndian.Uint32(buf[j:]))
		}
	} else {
		for i := 0; i < n; i++ {
			start := i * 4
			dst[i] = meta.SafecastUint32ToInt32Bits(t.ByteOrder.Uint32(buf[start : start+4]))
		}
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
	maxBytes, ok := meta.SafecastIntToUint32(n)
	if !ok {
		return 0
	}
	buf, err := r.readTagBytes(t, maxBytes)
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
	maxBytes, ok := meta.SafecastIntToUint32(n * 8)
	if !ok {
		return 0
	}
	buf, err := r.readTagBytes(t, maxBytes)
	if err != nil {
		return 0
	}
	if got := len(buf) / 8; got < n {
		n = got
	}
	if t.ByteOrder == utils.LittleEndian {
		for i, j := 0, 0; i < n; i, j = i+1, j+8 {
			dst[i*2] = meta.SafecastUint32ToInt32Bits(binary.LittleEndian.Uint32(buf[j:]))
			dst[i*2+1] = meta.SafecastUint32ToInt32Bits(binary.LittleEndian.Uint32(buf[j+4:]))
		}
	} else {
		for i := 0; i < n; i++ {
			start := i * 8
			dst[i*2] = meta.SafecastUint32ToInt32Bits(t.ByteOrder.Uint32(buf[start : start+4]))
			dst[i*2+1] = meta.SafecastUint32ToInt32Bits(t.ByteOrder.Uint32(buf[start+4 : start+8]))
		}
	}
	return n
}

// parseDate parses the requested value from EXIF metadata.
func (r *Reader) parseDate(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 19 {
		return time.Time{}
	}
	// YYYY:MM:DD HH:MM:SS
	if buf[4] != ':' || buf[7] != ':' || buf[10] != ' ' || buf[13] != ':' || buf[16] != ':' {
		return time.Time{}
	}
	year, ok := meta.SafecastUintToInt(parseStrUint(buf[0:4]))
	if !ok {
		return time.Time{}
	}
	month, ok := meta.SafecastUintToInt(parseStrUint(buf[5:7]))
	if !ok {
		return time.Time{}
	}
	day, ok := meta.SafecastUintToInt(parseStrUint(buf[8:10]))
	if !ok {
		return time.Time{}
	}
	hour, ok := meta.SafecastUintToInt(parseStrUint(buf[11:13]))
	if !ok {
		return time.Time{}
	}
	minute, ok := meta.SafecastUintToInt(parseStrUint(buf[14:16]))
	if !ok {
		return time.Time{}
	}
	second, ok := meta.SafecastUintToInt(parseStrUint(buf[17:19]))
	if !ok {
		return time.Time{}
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}

// parseOffsetTime parses the requested value from EXIF metadata.
func (r *Reader) parseOffsetTime(t tag.Entry) *time.Location {
	if !t.IsType(tag.TypeASCII) {
		return nil
	}
	buf, err := r.readTagBytes(t, 8)
	if err != nil || len(buf) < 6 {
		return nil
	}
	if buf[3] != ':' {
		return nil
	}
	hours, ok := meta.SafecastUintToInt(parseStrUint(buf[1:3]))
	if !ok {
		return nil
	}
	minutes, ok := meta.SafecastUintToInt(parseStrUint(buf[4:6]))
	if !ok {
		return nil
	}
	offset := hours*hoursToSeconds + minutes*minutesToSeconds
	offset32, ok := meta.SafecastIntToInt32(offset)
	if !ok {
		return nil
	}
	switch buf[0] {
	case '-':
		return getLocation(-offset32, buf[:6])
	case '+':
		return getLocation(offset32, buf[:6])
	default:
		return nil
	}
}

var (
	exifTextASCII   = []byte("ASCII\x00\x00\x00")
	exifTextUnicode = []byte("UNICODE\x00")
	exifTextJIS     = []byte("JIS\x00\x00\x00\x00\x00")
)

func isNULOrSpace(b byte) bool {
	return b == 0 || b == ' '
}

func classifyExifTextHeader(header []byte) exifTextEncoding {
	if len(header) < 8 {
		return exifTextEncodingUnknown
	}

	// Common canonical headers fast-path.
	switch {
	case bytes.Equal(header, exifTextASCII):
		return exifTextEncodingASCII
	case bytes.Equal(header, exifTextUnicode):
		return exifTextEncodingUnicode
	case bytes.Equal(header, exifTextJIS):
		return exifTextEncodingJIS
	}

	// ExifTool compatibility for malformed-but-common padding:
	// ASCII and empty IDs padded with NUL/spaces.
	if bytes.Equal(header[:5], []byte("ASCII")) {
		if isNULOrSpace(header[5]) && isNULOrSpace(header[6]) && isNULOrSpace(header[7]) {
			return exifTextEncodingASCII
		}
		return exifTextEncodingUnknown
	}
	if bytes.Equal(header[:7], []byte("UNICODE")) {
		// Match ExifTool behavior: uppercase UNICODE only.
		if isNULOrSpace(header[7]) {
			return exifTextEncodingUnicode
		}
		return exifTextEncodingUnknown
	}
	if bytes.Equal(header[:3], []byte("JIS")) {
		for i := 3; i < 8; i++ {
			if !isNULOrSpace(header[i]) {
				return exifTextEncodingUnknown
			}
		}
		return exifTextEncodingJIS
	}
	for i := 0; i < 8; i++ {
		if !isNULOrSpace(header[i]) {
			return exifTextEncodingUnknown
		}
	}
	return exifTextEncodingASCII
}

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
		return exifASCIIText(trimASCIIWhitespace(trimNULBuffer(buf)))
	}

	header := buf[:8]
	payload := buf[8:]
	switch classifyExifTextHeader(header) {
	case exifTextEncodingASCII, exifTextEncodingJIS:
		return exifASCIIText(payload)
	case exifTextEncodingUnicode:
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
		n16, ok := meta.SafecastInt32ToInt16(n)
		if !ok {
			return meta.NewExposureBias(0, 0)
		}
		d16, ok := meta.SafecastInt32ToInt16(d)
		if !ok {
			return meta.NewExposureBias(0, 0)
		}
		return meta.NewExposureBias(n16, d16)
	case t.IsType(tag.TypeRational):
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return meta.NewExposureBias(0, 0)
		}
		n16, ok := meta.SafecastUint32ToInt16(rat[0])
		if !ok {
			return meta.NewExposureBias(0, 0)
		}
		d16, ok := meta.SafecastUint32ToInt16(rat[1])
		if !ok {
			return meta.NewExposureBias(0, 0)
		}
		return meta.NewExposureBias(n16, d16)
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
func (r *Reader) parseLensInfo(t tag.Entry) LensInfo {
	if t.IsEmbedded() {
		return LensInfo{}
	}
	buf, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 32 {
		return LensInfo{}
	}
	return LensInfo{
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
		value, ok := meta.SafecastUint32ToUint16(r.parseUint32(t))
		if !ok {
			return 0
		}
		return value
	case t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		s := strings.TrimSpace(r.parseString(t))
		if s == "" {
			return 0
		}
		value, ok := meta.SafecastUintToUint16(parseStrUint([]byte(s)))
		if !ok {
			return 0
		}
		return value
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

// TODO(simd): accelerate trim/ASCII scan and offset-date parsing with simd/archsimd
// when GOEXPERIMENT=simd integration is prioritized.
