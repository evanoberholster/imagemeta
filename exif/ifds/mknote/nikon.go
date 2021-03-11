package mknote

// IsNikonMkNoteHeaderBytes represents "Nikon" the first 5 bytes of the
func IsNikonMkNoteHeaderBytes(buf []byte) bool {
	return buf[0] == 'N' &&
		buf[1] == 'i' &&
		buf[2] == 'k' &&
		buf[3] == 'o' &&
		buf[4] == 'n'
}
