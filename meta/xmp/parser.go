package xmp

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/evanoberholster/imagemeta/meta"
)

// parser routes one parsed XMP property to its namespace-specific decoder.
//
// Parsing is intentionally permissive. Unknown or unhandled fields never fail
// packet decoding. In debug mode, non-ErrPropertyNotSet parse errors are
// logged to help diagnose unsupported input.
func (xmp *XMP) parser(p property, debug bool) (err error) {
	if len(p.Value()) == 0 {
		return
	}
	switch p.Namespace() {
	case XMLnsNS:
		return // Null operation
	case ExifNS, ExifEXNS:
		err = xmp.Exif.parse(p)
	case AuxNS:
		err = xmp.Aux.parse(p)
	case DcNS:
		err = xmp.DC.parse(p)
	case XNS, XmpNS, XapNS:
		err = xmp.Basic.parse(p)
	case TiffNS:
		err = xmp.Tiff.parse(p)
	case CrsNS:
		err = xmp.CRS.parse(p)
	case PhotoshopNS:
		err = xmp.Photoshop.parse(p)
	case XmpMMNS, XapMMNS, StEvtNS, StRefNS:
		err = xmp.MM.parse(p)
	case XmpDMNS:
		err = xmp.DynamicMedia.parse(p)
	case LrNS:
		err = xmp.Lightroom.parse(p)
	case MwgRSNS, StDimNS, StAreaNS, AppleFiNS:
		err = xmp.Regions.parse(p)
	default:
		//fmt.Println(p, ns)
		return
	}
	xmp.markParsed(p.Namespace())
	if err != nil {
		// The decoder is intentionally permissive:
		// unknown/unhandled properties must not fail packet parsing.
		// In debug mode, surface non-ErrPropertyNotSet parse failures.
		if debug && !errors.Is(err, ErrPropertyNotSet) {
			fmt.Println("XMP parse warning:", err, p)
		}
		return nil
	}

	return
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// parseUUID parses a UUID value, optionally skipping a leading prefix up to ':'.
func parseUUID(buf []byte) (uuid meta.UUID) {
	if _, b := readUntil(buf, ':'); len(b) > 0 {
		buf = b
	}
	_ = uuid.UnmarshalText(buf)
	return
}

// parseUintWithLimit parses an unsigned decimal integer with overflow bounds.
func parseUintWithLimit(buf []byte, max uint64) (uint64, bool) {
	if len(buf) == 0 {
		return 0, false
	}
	var u uint64
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		d := uint64(c - '0')
		if u > (max-d)/10 {
			return 0, false
		}
		u = u*10 + d
	}
	return u, true
}

// parseUint parses an unsigned decimal uint64 value.
// Invalid or out-of-range values return 0.
func parseUint(buf []byte) uint64 {
	u, ok := parseUintWithLimit(buf, math.MaxUint64)
	if !ok {
		return 0
	}
	return u
}

// parseUint32 parses an unsigned decimal uint32 value.
// Invalid or out-of-range values return 0.
func parseUint32(buf []byte) uint32 {
	u, ok := parseUintWithLimit(buf, math.MaxUint32)
	if !ok {
		return 0
	}
	return uint32(u)
}

// parseUint16 parses an unsigned decimal uint16 value.
// Invalid or out-of-range values return 0.
func parseUint16(buf []byte) uint16 {
	u, ok := parseUintWithLimit(buf, math.MaxUint16)
	if !ok {
		return 0
	}
	return uint16(u)
}

// parseUint8 parses an unsigned decimal uint8 value.
// Invalid or out-of-range values return 0.
func parseUint8(buf []byte) uint8 {
	u, ok := parseUintWithLimit(buf, math.MaxUint8)
	if !ok {
		return 0
	}
	return uint8(u)
}

// parseIntWithLimit parses a signed decimal integer constrained to [min, max].
func parseIntWithLimit(buf []byte, min, max int64) (int64, bool) {
	if len(buf) == 0 {
		return 0, false
	}

	sign := int64(1)
	i := 0
	switch buf[0] {
	case '-':
		sign = -1
		i = 1
	case '+':
		i = 1
	}
	if i >= len(buf) {
		return 0, false
	}

	absLimit := max
	absMin := -min
	// math.MinInt64 cannot be negated safely; clamp for bound checks.
	if min == math.MinInt64 {
		absMin = math.MaxInt64
	}
	if absMin > absLimit {
		absLimit = absMin
	}

	var n int64
	for ; i < len(buf); i++ {
		c := buf[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		d := int64(c - '0')
		if n > (absLimit-d)/10 {
			return 0, false
		}
		n = (n * 10) + d
	}

	n *= sign
	if n < min || n > max {
		return 0, false
	}
	return n, true
}

// parseInt32 parses a []byte string representation of an int32 value.
// If the value is invalid or out of range it returns 0.
func parseInt32(buf []byte) int32 {
	n, ok := parseIntWithLimit(buf, math.MinInt32, math.MaxInt32)
	if !ok {
		return 0
	}
	return int32(n)
}

// parseInt16 parses a []byte string representation of an int16 value.
// If the value is invalid or out of range it returns 0.
func parseInt16(buf []byte) int16 {
	n, ok := parseIntWithLimit(buf, math.MinInt16, math.MaxInt16)
	if !ok {
		return 0
	}
	return int16(n)
}

// parseFloat64 parses a decimal float64 value.
// Invalid input returns 0.
func parseFloat64(buf []byte) float64 {
	f, ok := parseFloat64OK(buf)
	if !ok {
		return 0.0
	}
	return f
}

// parseFloat64OK parses a decimal float and reports parse success.
func parseFloat64OK(buf []byte) (float64, bool) {
	f, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// parseString converts a byte slice to string.
func parseString(buf []byte) string {
	return string(buf)
}

// parseRational parses a rational number formatted as "numerator/denominator".
// Invalid or out-of-range components are coerced to 0, and denominator defaults to 1.
func parseRational(buf []byte) (n uint32, d uint32) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			n = parseUint32(buf[:i])
			d = parseUint32(buf[i+1:])
			if d == 0 {
				return 0, 1
			}
			return
		}
	}
	return parseUint32(buf), 1
}

// parseRational32 parses a signed rational number formatted as
// "numerator/denominator". Invalid values return 0/1.
func parseRational32(buf []byte) (n int32, d int32) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			n = parseInt32(buf[:i])
			d = parseInt32(buf[i+1:])
			if d == 0 {
				return 0, 1
			}
			return n, d
		}
	}
	return parseInt32(buf), 1
}

// parseRationalFloat64 parses either "n/d" or plain decimal numeric input.
func parseRationalFloat64(buf []byte) float64 {
	f, _ := parseRationalFloat64OK(buf)
	return f
}

// parseRationalFloat64OK parses either "n/d" or plain decimal numeric input
// and reports parse success.
func parseRationalFloat64OK(buf []byte) (float64, bool) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			n, nOK := parseFloat64OK(buf[:i])
			d, dOK := parseFloat64OK(buf[i+1:])
			if nOK && dOK && d != 0 {
				return n / d, true
			}
			break
		}
	}
	return parseFloat64OK(buf)
}

// parseApexAperture converts an APEX aperture value to an f-number.
func parseApexAperture(buf []byte) float64 {
	av, ok := parseRationalFloat64OK(buf)
	if !ok {
		return 0
	}
	return math.Pow(2, av/2)
}

// parseApexShutterSpeed converts an APEX shutter speed value to seconds.
func parseApexShutterSpeed(buf []byte) float64 {
	tv, ok := parseRationalFloat64OK(buf)
	if !ok {
		return 0
	}
	return math.Pow(2, -tv)
}

// parseBool parses common XMP boolean encodings.
//
// NOTE: buf is mutated in place via lowercase normalization.
func parseBool(buf []byte) bool {
	if len(buf) == 0 {
		return false
	}
	switch buf[0] {
	case '1':
		return true
	case '0':
		return false
	}

	toLowercaseBufInPlace(buf)
	switch string(buf) {
	case "on":
		return true
	case "yes":
		return true
	case "true":
		return true
	default: // off, no, false
		return false
	}
}

// parseGPSCoordinate parses XMP GPS coordinates in either decimal format or
// ExifTool-style "deg,minutesHemisphere" format (for example "11,57.1312N").
func parseGPSCoordinate(buf []byte) float64 {
	if len(buf) == 0 {
		return 0
	}

	// "11,57.1312N" -> 11 + 57.1312/60
	last := buf[len(buf)-1]
	sign := 1.0
	switch last {
	case 'N', 'E', 'n', 'e':
		buf = buf[:len(buf)-1]
	case 'S', 'W', 's', 'w':
		sign = -1.0
		buf = buf[:len(buf)-1]
	}

	for i := range buf {
		if buf[i] == ',' {
			deg := parseFloat64(buf[:i])
			minutes := parseFloat64(buf[i+1:])
			return sign * (deg + (minutes / 60))
		}
	}

	return sign * parseFloat64(buf)
}

// readUntil splits buf at the first delimiter (or '>'), returning both sides.
func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter || buf[i] == '>' {
			return buf[:i], buf[i+1:]
		}
	}
	return buf, nil
}

// decodeXMLEntities decodes the small XML entity subset used in XMP text
// values, including numeric entities.
func decodeXMLEntities(buf []byte) []byte {
	if len(buf) == 0 || bytes.IndexByte(buf, '&') < 0 {
		return buf
	}

	out := make([]byte, 0, len(buf))
	for i := 0; i < len(buf); i++ {
		if buf[i] != '&' {
			out = append(out, buf[i])
			continue
		}

		semi := -1
		for j := i + 1; j < len(buf); j++ {
			if buf[j] == ';' {
				semi = j
				break
			}
			if j-i > 12 {
				break
			}
		}
		if semi < 0 {
			out = append(out, buf[i])
			continue
		}

		entity := buf[i+1 : semi]
		switch {
		case bytes.Equal(entity, []byte("amp")):
			out = append(out, '&')
		case bytes.Equal(entity, []byte("lt")):
			out = append(out, '<')
		case bytes.Equal(entity, []byte("gt")):
			out = append(out, '>')
		case bytes.Equal(entity, []byte("quot")):
			out = append(out, '"')
		case bytes.Equal(entity, []byte("apos")):
			out = append(out, '\'')
		case len(entity) >= 2 && entity[0] == '#':
			if r, ok := parseXMLNumericEntity(entity[1:]); ok {
				out = appendRuneUTF8(out, r)
			} else {
				out = append(out, buf[i:semi+1]...)
			}
		default:
			out = append(out, buf[i:semi+1]...)
		}
		i = semi
	}
	return out
}

// parseXMLNumericEntity parses XML numeric entities without the leading "#".
// Both decimal and hexadecimal ("x...") forms are supported.
func parseXMLNumericEntity(buf []byte) (rune, bool) {
	if len(buf) == 0 {
		return 0, false
	}
	base := uint32(10)
	if len(buf) >= 2 && (buf[0] == 'x' || buf[0] == 'X') {
		base = 16
		buf = buf[1:]
		if len(buf) == 0 {
			return 0, false
		}
	}

	var v uint32
	for i := 0; i < len(buf); i++ {
		digit := uint32(0xFFFFFFFF)
		c := buf[i]
		switch {
		case '0' <= c && c <= '9':
			digit = uint32(c - '0')
		case base == 16 && 'a' <= c && c <= 'f':
			digit = uint32(c-'a') + 10
		case base == 16 && 'A' <= c && c <= 'F':
			digit = uint32(c-'A') + 10
		}
		if digit == 0xFFFFFFFF || digit >= base {
			return 0, false
		}
		v = v*base + digit
		if v > 0x10FFFF {
			return 0, false
		}
	}
	return rune(v), true
}

// appendRuneUTF8 appends r encoded as UTF-8 to dst.
func appendRuneUTF8(dst []byte, r rune) []byte {
	switch {
	case r <= 0x7F:
		return append(dst, byte(r))
	case r <= 0x7FF:
		dst = append(dst, 0xC0|byte(r>>6))
		dst = append(dst, 0x80|byte(r&0x3F))
		return dst
	case r <= 0xFFFF:
		dst = append(dst, 0xE0|byte(r>>12))
		dst = append(dst, 0x80|byte((r>>6)&0x3F))
		dst = append(dst, 0x80|byte(r&0x3F))
		return dst
	default:
		dst = append(dst, 0xF0|byte(r>>18))
		dst = append(dst, 0x80|byte((r>>12)&0x3F))
		dst = append(dst, 0x80|byte((r>>6)&0x3F))
		dst = append(dst, 0x80|byte(r&0x3F))
		return dst
	}
}
