package utils

import (
	"bytes"
	"io"
	"testing"
)

func TestBufioReaderPoolAcquireRelease(t *testing.T) {
	pool := NewBufioReaderPool(16, bytes.NewReader(nil))

	a := pool.Acquire(bytes.NewReader([]byte("abc")))
	buf, err := a.Peek(3)
	if err != nil {
		t.Fatalf("Peek error = %v", err)
	}
	if got := string(buf); got != "abc" {
		t.Fatalf("Peek data = %q, want %q", got, "abc")
	}
	pool.Release(a)

	b := pool.Acquire(bytes.NewReader([]byte("xy")))
	buf, err = b.Peek(2)
	if err != nil {
		t.Fatalf("Peek error = %v", err)
	}
	if got := string(buf); got != "xy" {
		t.Fatalf("Peek data = %q, want %q", got, "xy")
	}
	if _, err := b.Peek(3); err == nil && b.Buffered() < 3 {
		t.Fatal("expected short-buffer peek error after reset to new source")
	}
	pool.Release(b)
}

func TestBufioReaderPoolAcquireNilPool(t *testing.T) {
	var pool *BufioReaderPool
	br := pool.Acquire(bytes.NewReader([]byte("z")))
	buf := make([]byte, 1)
	n, err := br.Read(buf)
	if n != 1 || err != nil && err != io.EOF {
		t.Fatalf("Read = (%d, %v), want (1, nil/EOF)", n, err)
	}
	if got := string(buf[:n]); got != "z" {
		t.Fatalf("Read data = %q, want %q", got, "z")
	}
}
