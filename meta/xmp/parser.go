package xmp

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

func (xmp *XMP) parser(p property) (err error) {
	if len(p.Value()) == 0 {
		return
	}
	switch p.Namespace() {
	case XMLnsNS:
		return // Null operation
	case ExifNS, ExifEXNS:
		err = parseNamespace(&xmp.Exif, p, func(v *Exif, prop property) error { return v.parse(prop) })
	case AuxNS:
		err = parseNamespace(&xmp.Aux, p, func(v *Aux, prop property) error { return v.parse(prop) })
	case DcNS:
		err = parseNamespace(&xmp.DC, p, func(v *DublinCore, prop property) error { return v.parse(prop) })
	case XNS, XmpNS, XapNS:
		err = xmp.Basic.parse(p)
	case TiffNS:
		err = parseNamespace(&xmp.Tiff, p, func(v *Tiff, prop property) error { return v.parse(prop) })
	case CrsNS:
		err = parseNamespace(&xmp.CRS, p, func(v *CRS, prop property) error { return v.parse(prop) })
	case PhotoshopNS:
		err = parseNamespace(&xmp.Photoshop, p, func(v *Photoshop, prop property) error { return v.parse(prop) })
	case XmpMMNS, XapMMNS, StEvtNS, StRefNS:
		err = parseNamespace(&xmp.MM, p, func(v *XMPMM, prop property) error { return v.parse(prop) })
	case XmpDMNS:
		err = parseNamespace(&xmp.DynamicMedia, p, func(v *DynamicMedia, prop property) error { return v.parse(prop) })
	case LrNS:
		err = parseNamespace(&xmp.Lightroom, p, func(v *Lightroom, prop property) error { return v.parse(prop) })
	default:
		//fmt.Println(p, ns)
		return
	}
	if err != nil {
		// The decoder is intentionally permissive:
		// unknown/unhandled properties must not fail packet parsing.
		// In debug mode, surface non-ErrPropertyNotSet parse failures.
		if DebugMode && !errors.Is(err, ErrPropertyNotSet) {
			fmt.Println("XMP parse warning:", err, p)
		}
		return nil
	}

	return
}

// parseNamespace lazily allocates a namespace only when parsing succeeds.
func parseNamespace[T any](dst **T, p property, parse func(*T, property) error) error {
	if *dst != nil {
		return parse(*dst, p)
	}

	var v T
	if err := parse(&v, p); err != nil {
		return err
	}
	*dst = &v
	return nil
}

// parseDate parses a Date and returns a time.Time or an error
func parseDate(buf []byte) (t time.Time, err error) {
	str := string(buf)
	if t, err = time.Parse(time.RFC3339Nano, str); err == nil {
		return t, nil
	}
	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.00", str); err != nil {
			return time.Parse("2006-01-02T15:04:05", str)
		}
	}
	return
}

// parseUUID parses a UUID and returns a meta.UUID
func parseUUID(buf []byte) (uuid meta.UUID) {
	if _, b := readUntil(buf, ':'); len(b) > 0 {
		buf = b
	}
	err := uuid.UnmarshalText(buf)
	if err != nil {
		if DebugMode {
			fmt.Println("Parse UUID error: ", err)
		}
	}
	return
}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	if len(buf) == 0 {
		return 0
	}
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		if c < '0' || c > '9' {
			return 0
		}
		d := uint64(c - '0')
		if u > (math.MaxUint64-d)/10 {
			return 0
		}
		u = u*10 + d
	}
	return
}

// parseUint32 parses a []byte of a string representation of a uint32 value and returns the value.
// If the value is larger than uint32 returns 0.
func parseUint32(buf []byte) (u uint32) {
	if i := parseUint(buf); i <= math.MaxUint32 {
		return uint32(i)
	}
	return 0
}

// parseUint8 parses a []byte of a string representation of a uint8 value and returns the value.
// If the value is larger than uint8 returns 0.
func parseUint8(buf []byte) (u uint8) {
	if i := parseUint(buf); i <= math.MaxUint8 {
		return uint8(i)
	}
	return 0
}

// parseInt32 parses a []byte string representation of an int32 value.
// If the value is invalid or out of range it returns 0.
func parseInt32(buf []byte) (v int32) {
	if len(buf) == 0 {
		return 0
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
		return 0
	}

	n := int64(0)
	for ; i < len(buf); i++ {
		c := buf[i]
		if c < '0' || c > '9' {
			return 0
		}
		n = (n * 10) + int64(c-'0')
		if n > math.MaxInt32 {
			return 0
		}
	}

	n *= sign
	if n < math.MinInt32 || n > math.MaxInt32 {
		return 0
	}
	return int32(n)
}

// parseInt16 parses a []byte string representation of an int16 value.
// If the value is invalid or out of range it returns 0.
func parseInt16(buf []byte) int16 {
	v := parseInt32(buf)
	if v < math.MinInt16 || v > math.MaxInt16 {
		return 0
	}
	return int16(v)
}

// parseFloat64 parses a []byte of a string representation of a float64 value and returns the value
func parseFloat64(buf []byte) (f float64) {
	f, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return 0.0
	}
	return
}

// parseString parses a []byte and returns a string
func parseString(buf []byte) string {
	return string(buf)
}

// parseRational separates a string into a fraction.
// With "n" as the numerator and "d" as the denominator.
// TODO: Improve parsing functionality
func parseRational(buf []byte) (n uint32, d uint32) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			n = uint32(parseUint(buf[:i]))
			d = uint32(parseUint(buf[i+1:]))
			if d == 0 {
				return 0, 1
			}
			return
		}
	}
	return uint32(parseUint(buf)), 1
}

func parseRationalFloat64(buf []byte) float64 {
	f, _ := parseRationalFloat64OK(buf)
	return f
}

func parseRationalFloat64OK(buf []byte) (float64, bool) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			n, nErr := strconv.ParseFloat(string(buf[:i]), 64)
			d, dErr := strconv.ParseFloat(string(buf[i+1:]), 64)
			if nErr == nil && dErr == nil && d != 0 {
				return n / d, true
			}
			break
		}
	}
	f, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func parseApexAperture(buf []byte) float64 {
	av, ok := parseRationalFloat64OK(buf)
	if !ok {
		return 0
	}
	return math.Pow(2, av/2)
}

func parseApexShutterSpeed(buf []byte) float64 {
	tv, ok := parseRationalFloat64OK(buf)
	if !ok {
		return 0
	}
	return math.Pow(2, -tv)
}

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
	return eqFoldString(buf, "true") || eqFoldString(buf, "yes") || eqFoldString(buf, "on")
}

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

func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter || buf[i] == '>' {
			return buf[:i], buf[i+1:]
		}
	}
	return buf, nil
}

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
