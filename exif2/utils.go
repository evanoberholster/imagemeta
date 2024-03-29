package exif2

// static values
const (
	hoursToSeconds   = 60 * minutesToSeconds
	minutesToSeconds = 60
)

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

// trimNULBuffer removes trailing bytes from Buffer
func trimNULBuffer(buf []byte) []byte {
	for i := len(buf) - 1; i > 0; i-- {
		if buf[i] == 0 || buf[i] == ' ' || buf[i] == '\n' {
			continue
		}
		return buf[:i+1]
	}
	return nil
}

//func lowerCase(buf []byte) []byte {
//	for i := 0; i < len(buf); i++ {
//		a := buf[i]
//		if 'A' <= a && a <= 'Z' {
//			buf[i] += a
//		}
//	}
//	return buf
//}
