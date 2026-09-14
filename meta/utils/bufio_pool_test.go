package utils

import (
	"bytes"
	"testing"
)

func TestBufioReaderPoolAcquireRelease(t *testing.T) {
	t.Parallel()

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
