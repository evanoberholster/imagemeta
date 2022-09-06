package exif2

// parseStrUint parses a []byte of a string representation of a uint value and returns the value.
func parseStrUint(buf []byte) (u uint) {
	for i := 0; i < len(buf); i++ {
		if buf[i] >= '0' {
			u *= 10
			u += uint(buf[i] - '0')
		}
	}
	return
}

func trimNULString(buf []byte) string {
	for i := len(buf) - 1; i > 0; i-- {
		if buf[i] == 0 || buf[i] == ' ' {
			continue
		}
		return string(buf[:i+1])
	}
	return ""
}
