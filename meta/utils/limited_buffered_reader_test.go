package utils

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestLimitedBufferedReaderReadBounded(t *testing.T) {
	src := bufio.NewReader(bytes.NewReader([]byte("abcdef")))
	lr := NewLimitedBufferedReader(src, 3)

	buf := make([]byte, 8)
	n, err := lr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error = %v", err)
	}
	if got := string(buf[:n]); got != "abc" {
		t.Fatalf("Read data = %q, want %q", got, "abc")
	}

	n, err = lr.Read(buf)
	if n != 0 || err != io.EOF {
		t.Fatalf("second Read = (%d, %v), want (0, EOF)", n, err)
	}
}

func TestLimitedBufferedReaderPeekDiscardBounded(t *testing.T) {
	src := bufio.NewReader(bytes.NewReader([]byte("abcdef")))
	lr := NewLimitedBufferedReader(src, 4)

	peek, err := lr.Peek(6)
	if err != nil {
		t.Fatalf("Peek error = %v", err)
	}
	if got := string(peek); got != "abcd" {
		t.Fatalf("Peek data = %q, want %q", got, "abcd")
	}

	discarded, err := lr.Discard(6)
	if discarded != 4 || err != io.EOF {
		t.Fatalf("Discard = (%d, %v), want (4, EOF)", discarded, err)
	}
}
