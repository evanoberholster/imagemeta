package exif

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	// Test Trim
	a := []byte{'a', 'b', 'c', 'd', '.', ' '}
	if !bytes.Equal(trim(a), a[:len(a)-1]) {
		t.Errorf("Trim should remove trailing spaces: expected %s got %s", a[:len(a)-1], trim(a))
	}
	a = []byte{' ', ' ', ' ', ' ', ' ', ' '}
	if len(trim(a)) != 0 {
		t.Errorf("Trim should remove trailing spaces: expected %d got %d", 0, len(trim(a)))
	}
}
