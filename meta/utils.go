package meta

// parseInt parses a []byte of a string representation of an int64 value and returns the value
func parseInt(buf []byte) (i int64) {
	var neg bool
	if buf[0] == '-' {
		buf = buf[1:]
		neg = true
	}
	i = int64(parseUint(buf))
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

// parseUntil parses a []byte and splits the []byte at delimiter.
// Returns a and b without delimiter present.
func parseUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter {
			a = buf[:i]
			if i < len(buf)+1 {
				b = buf[i+1:]
				return
			}
			return
		}
	}
	return buf, nil
}
