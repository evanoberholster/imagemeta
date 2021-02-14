package xmp

// Reader is an XMP Reader interface
type Reader interface {
	Discard(n int) (discarded int, err error)
	Peek(n int) ([]byte, error)
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
	ReadSlice(delim byte) (line []byte, err error)
	UnreadByte() error
}
