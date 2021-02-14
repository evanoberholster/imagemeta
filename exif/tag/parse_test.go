package tag

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	if valueIsEmbbeded(16) {
		t.Errorf("ValueIsEmbbeded is true when equal or less than 4 bytes")
	}
	a := []byte{'a', 'b', 'c', 'd', '.', ' '}

	if !bytes.Equal(trim(a), a[:len(a)-1]) {
		t.Errorf("Trim should remove trailing spaces: expected %s got %s", trim(a), a[:len(a)-1])
	}
}
