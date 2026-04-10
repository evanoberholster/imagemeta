package panasonic

const PanasonicMakerNotePrefixLength = 12

// HasPanasonicHeader reports whether the maker-note payload starts with a
// Panasonic label prefix.
func HasPanasonicHeader(buf []byte) bool {
	return len(buf) >= PanasonicMakerNotePrefixLength &&
		buf[0] == 'P' &&
		buf[1] == 'a' &&
		buf[2] == 'n' &&
		buf[3] == 'a' &&
		buf[4] == 's' &&
		buf[5] == 'o' &&
		buf[6] == 'n' &&
		buf[7] == 'i' &&
		buf[8] == 'c' &&
		buf[9] == 0 &&
		buf[10] == 0 &&
		buf[11] == 0
}
