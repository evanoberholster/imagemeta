package meta

import "bytes"

// parseInt parses a []byte of a string representation of an int64 value and returns the value
//func parseInt(buf []byte) (i int64) {
//	var neg bool
//	if buf[0] == '-' {
//		buf = buf[1:]
//		neg = true
//	}
//	i = int64(parseUint(buf))
//	if neg {
//		i *= -1
//	}
//	return
//}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

// parseUntil parses a []byte and splits the []byte at delimiter.
// Returns a and b without delimiter present.
//func parseUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
//	for i := 0; i < len(buf); i++ {
//		if buf[i] == delimiter {
//			a = buf[:i]
//			if i < len(buf)+1 {
//				b = buf[i+1:]
//				return
//			}
//			return
//		}
//	}
//	return buf, nil
//}

var closeTagXMP = []byte("</x:xmpmeta>")

// CleanXMPSuffixWhiteSpace returns the same slice with the whitespace after "</x:xmpmeta>" removed.
func CleanXMPSuffixWhiteSpace(buf []byte) []byte {
	for i := len(buf) - 1; i > 12; i-- {
		if buf[i] == '>' && buf[i-1] == 'a' {
			// </x:xmpmeta>
			if bytes.Equal(closeTagXMP, buf[i-11:i+1]) {
				buf = buf[:i+1]
				return buf
			}
		}
	}
	return buf
}
