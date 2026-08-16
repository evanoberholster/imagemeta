package utils

// BufferedReader is compatible with *bufio.Reader.
type BufferedReader interface {
	Peek(n int) ([]byte, error)
	Discard(n int) (discarded int, err error)
	Read(p []byte) (n int, err error)
}
