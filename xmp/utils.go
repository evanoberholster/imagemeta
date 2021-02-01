package xmp

import (
	"bufio"
	"bytes"
	"time"
)

func readUntilByte(br *bufio.Reader, end byte) (n int, err error) {
	var b byte
	for {
		b, err = br.ReadByte()
		if err != nil {
			return
		}
		n++
		if b == end {
			return
		}
	}
}

func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter || buf[i] == markerGt {
			return buf[:i], buf[i+1:]
		}
	}
	return nil, nil
}

// parseDate
func parseDate(buf []byte) (t time.Time, err error) {
	str := string(buf)
	if t, err = time.Parse("2006-01-02T15:04:05Z07:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.00", str); err != nil {
			t, err = time.Parse("2006-01-02T15:04:05", str)
		}
	}
	return
}

//func parseUint32(s string) uint32 {
//	u64, err := strconv.ParseUint(s, 10, 32)
//	if err != nil {
//		return 0
//	}
//	return uint32(u64)
//}

// parseInt parses a []byte of a string representation of an int32 value and returns the value
func parseInt(buf []byte) (i int32) {
	var neg bool
	if buf[0] == '-' {
		buf = buf[1:]
		neg = true
	}
	for j := 0; j < len(buf); j++ {
		i *= 10
		i += int32(buf[j] - '0')
	}
	if neg {
		i *= -1
	}
	return
}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

// parseBool
func parseBool(buf []byte) bool {
	return bytes.EqualFold(buf, []byte("True"))
}

// parseString parses a []byte and returns a string
func parseString(buf []byte) string {
	return string(buf)
}
