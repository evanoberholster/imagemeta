package mknote

import (
	"fmt"
	"testing"
)

func TestParseCanonModelID(t *testing.T) {
	//rawID := uint32(0x80000001)
	rawID2 := uint32(0x4007d77b)
	id := parseCanonModelID(rawID2)
	t.Error("hello", id)
	t.Error(fmt.Sprintf("0x%04x", rawID2), fmt.Sprintf("0x%04x", id))
}
