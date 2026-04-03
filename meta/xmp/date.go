package xmp

import (
	"fmt"
	"time"
)

var parseDateStringFallback = parseDateString

const (
	// XMP/Exif timezone offsets are minute-based and bounded to +/-23:59.
	minCachedOffsetMinutes = -(23*60 + 59)
	maxCachedOffsetMinutes = 23*60 + 59
	cachedOffsetSpan       = maxCachedOffsetMinutes - minCachedOffsetMinutes + 1
)

// fixedOffsetByMinute caches fixed timezone locations keyed by minute offset.
var fixedOffsetByMinute = func() [cachedOffsetSpan]*time.Location {
	var zones [cachedOffsetSpan]*time.Location
	for minute := minCachedOffsetMinutes; minute <= maxCachedOffsetMinutes; minute++ {
		index := minute - minCachedOffsetMinutes
		if minute == 0 {
			zones[index] = time.UTC
			continue
		}
		zones[index] = time.FixedZone("", minute*60)
	}
	return zones
}()

// parseFastDate parses Exif/XMP timestamp shapes with low overhead.
func parseFastDate(data []byte) (time.Time, error) {
	return parseFastDateWithFallback(data, parseDateStringFallback)
}

// parseFastDateWithFallback parses Exif/XMP timestamp shapes and uses fallback
// when the fast path cannot validate the input.
func parseFastDateWithFallback(data []byte, fallback func(string) (time.Time, error)) (time.Time, error) {
	if fallback == nil {
		fallback = parseDateString
	}

	isISO, ok := identifyDateLayout(data)
	if !ok {
		return fallback(string(data))
	}

	year, month, day, hour, minute, sec, ok := parseDateClockParts(data)
	if !ok || !isValidDateTime(year, month, day, hour, minute, sec) {
		return fallback(string(data))
	}

	nsec := 0
	i := 19
	if i < len(data) && data[i] == '.' {
		var consumed int
		nsec, consumed, ok = parseFracNanos(data[i+1:])
		if !ok {
			return fallback(string(data))
		}
		i += 1 + consumed
	}

	loc := time.Local
	if isISO {
		loc = time.UTC
	}

	if i < len(data) {
		var consumed int
		loc, consumed, ok = parseTZOffset(data[i:])
		if !ok {
			return fallback(string(data))
		}
		i += consumed
	}

	if i != len(data) {
		return fallback(string(data))
	}

	return time.Date(year, time.Month(month), day, hour, minute, sec, nsec, loc), nil
}

// parseDateString parses date strings using common XMP layouts, from most
// specific to most permissive.
func parseDateString(s string) (t time.Time, err error) {
	if t, err = time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", s); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.00", s); err != nil {
			return time.Parse("2006-01-02T15:04:05", s)
		}
	}
	return
}

// identifyDateLayout reports whether data matches the fixed-width fast parser
// prefix. The returned bool indicates whether the shape is ISO (true) or Exif
// style (false).
func identifyDateLayout(data []byte) (isISO bool, ok bool) {
	if len(data) < 19 {
		return false, false
	}
	if data[13] != ':' || data[16] != ':' {
		return false, false
	}

	switch data[4] {
	case ':':
		if data[7] == ':' && data[10] == ' ' {
			return false, true
		}
	case '-':
		if data[7] == '-' && (data[10] == 'T' || data[10] == ' ') {
			return true, true
		}
	}
	return false, false
}

func parseDateClockParts(data []byte) (year, month, day, hour, minute, sec int, ok bool) {
	year, ok = parse4Digits(data[0], data[1], data[2], data[3])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	month, ok = parse2Digits(data[5], data[6])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	day, ok = parse2Digits(data[8], data[9])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	hour, ok = parse2Digits(data[11], data[12])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	minute, ok = parse2Digits(data[14], data[15])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	sec, ok = parse2Digits(data[17], data[18])
	if !ok {
		return 0, 0, 0, 0, 0, 0, false
	}

	return year, month, day, hour, minute, sec, true
}

func isValidDateTime(year, month, day, hour, minute, sec int) bool {
	if month < 1 || month > 12 || day < 1 || hour > 23 || minute > 59 || sec > 59 {
		return false
	}
	return day <= daysInMonth(year, month)
}

func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 0
	}
}

func isLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	}
	if year%100 != 0 {
		return true
	}
	return year%400 == 0
}

func parse2Digits(a, b byte) (int, bool) {
	if a < '0' || a > '9' || b < '0' || b > '9' {
		return 0, false
	}
	return int(a-'0')*10 + int(b-'0'), true
}

func parse4Digits(a, b, c, d byte) (int, bool) {
	if a < '0' || a > '9' || b < '0' || b > '9' || c < '0' || c > '9' || d < '0' || d > '9' {
		return 0, false
	}
	return int(a-'0')*1000 + int(b-'0')*100 + int(c-'0')*10 + int(d-'0'), true
}

func parseDigits(buf []byte) (int, bool) {
	switch len(buf) {
	case 2:
		return parse2Digits(buf[0], buf[1])
	case 4:
		return parse4Digits(buf[0], buf[1], buf[2], buf[3])
	default:
		return 0, false
	}
}

func parseFracNanos(buf []byte) (nsec int, consumed int, ok bool) {
	if len(buf) == 0 {
		return 0, 0, false
	}

	for consumed < len(buf) && consumed < 9 {
		c := buf[consumed]
		if c < '0' || c > '9' {
			break
		}
		nsec = (nsec * 10) + int(c-'0')
		consumed++
	}
	if consumed == 0 {
		return 0, 0, false
	}

	// Match RFC3339Nano precision in the fast path.
	if consumed < len(buf) && buf[consumed] >= '0' && buf[consumed] <= '9' {
		return 0, 0, false
	}

	for i := consumed; i < 9; i++ {
		nsec *= 10
	}
	return nsec, consumed, true
}

func parseTZOffset(buf []byte) (loc *time.Location, consumed int, ok bool) {
	if len(buf) == 0 {
		return nil, 0, false
	}

	switch buf[0] {
	case 'Z', 'z':
		return time.UTC, 1, true
	case '+', '-':
		if len(buf) >= 6 && buf[3] == ':' {
			hh, okH := parse2Digits(buf[1], buf[2])
			mm, okM := parse2Digits(buf[4], buf[5])
			if !okH || !okM || hh > 23 || mm > 59 {
				return nil, 0, false
			}
			seconds := (hh * 3600) + (mm * 60)
			if buf[0] == '-' {
				seconds = -seconds
			}
			return fixedOffsetZone(seconds), 6, true
		}
		if len(buf) >= 5 {
			hh, okH := parse2Digits(buf[1], buf[2])
			mm, okM := parse2Digits(buf[3], buf[4])
			if !okH || !okM || hh > 23 || mm > 59 {
				return nil, 0, false
			}
			seconds := (hh * 3600) + (mm * 60)
			if buf[0] == '-' {
				seconds = -seconds
			}
			return fixedOffsetZone(seconds), 5, true
		}
	}

	return nil, 0, false
}

func fixedOffsetZone(seconds int) *time.Location {
	if seconds == 0 {
		return time.UTC
	}

	if seconds%60 == 0 {
		minute := seconds / 60
		if minute >= minCachedOffsetMinutes && minute <= maxCachedOffsetMinutes {
			return fixedOffsetByMinute[minute-minCachedOffsetMinutes]
		}
	}

	return time.FixedZone("", seconds)
}

// parseDate parses XMP date values using CIPA DC-010-2017 date variants.
func parseDate(buf []byte) (time.Time, error) {
	trimmed := trimSpace(buf)
	if len(trimmed) == 0 {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	// Fast path for CIPA full timestamp forms.
	if isCIPADateTime(trimmed) {
		if t, err := parseFastDate(trimmed); err == nil {
			return t, nil
		}
	}

	// Fallback for CIPA reduced-precision and non-timezoned date-time forms.
	s := string(trimmed)
	for i := range cipaDateLayouts {
		if t, err := time.Parse(cipaDateLayouts[i], s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date: %q", s)
}

// cipaDateLayouts are the CIPA DC-010-2017 date variants accepted for XMP:
//   - YYYY
//   - YYYY-MM
//   - YYYY-MM-DD
//   - YYYY-MM-DDThh:mm[TZD]
//   - YYYY-MM-DDThh:mm:ss[TZD]
//   - YYYY-MM-DDThh:mm:ss.s[TZD]
//
// TZD is optional and represented as "Z" or "+/-hh:mm".
var cipaDateLayouts = [...]string{
	"2006",
	"2006-01",
	"2006-01-02",
	"2006-01-02T15:04",
	"2006-01-02T15:04Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05.999999999Z07:00",
}

// isCIPADateTime validates CIPA "YYYY-MM-DDThh:mm[:ss[.s]][TZD]" shape.
// Based on specification at: https://web.archive.org/web/20180921145139if_/http://www.cipa.jp:80/std/documents/e/DC-010-2017_E.pdf
func isCIPADateTime(buf []byte) bool {
	if len(buf) < len("2006-01-02T15:04") {
		return false
	}
	if !isDigit(buf[0]) || !isDigit(buf[1]) || !isDigit(buf[2]) || !isDigit(buf[3]) {
		return false
	}
	if buf[4] != '-' || !isDigit(buf[5]) || !isDigit(buf[6]) || buf[7] != '-' || !isDigit(buf[8]) || !isDigit(buf[9]) {
		return false
	}
	if buf[10] != 'T' {
		return false
	}
	if !isDigit(buf[11]) || !isDigit(buf[12]) || buf[13] != ':' || !isDigit(buf[14]) || !isDigit(buf[15]) {
		return false
	}

	i := 16
	// Optional :ss
	if i < len(buf) && buf[i] == ':' {
		if i+2 >= len(buf) || !isDigit(buf[i+1]) || !isDigit(buf[i+2]) {
			return false
		}
		i += 3
	}

	// Optional .frac (at least one digit)
	if i < len(buf) && buf[i] == '.' {
		i++
		if i >= len(buf) || !isDigit(buf[i]) {
			return false
		}
		for i < len(buf) && isDigit(buf[i]) {
			i++
		}
	}

	// No timezone suffix.
	if i == len(buf) {
		return true
	}

	// UTC designator.
	if buf[i] == 'Z' && i+1 == len(buf) {
		return true
	}

	// Numeric offset: +/-hh:mm
	if (buf[i] == '+' || buf[i] == '-') && i+5 < len(buf) && i+6 == len(buf) {
		return isDigit(buf[i+1]) && isDigit(buf[i+2]) && buf[i+3] == ':' && isDigit(buf[i+4]) && isDigit(buf[i+5])
	}

	return false
}
