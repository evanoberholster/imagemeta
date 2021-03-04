package meta

import (
	"encoding/binary"
	"testing"
)

func TestBinaryOrder(t *testing.T) {
	buf := []byte{0, 0, 0, 0}
	bo := BinaryOrder(buf)
	if bo != nil {
		t.Error("Binary Order for an empty buffer should be nil.")
	}

	buf = []byte{0x49, 0x49, 0x2a, 0}
	bo = BinaryOrder(buf)
	if bo != binary.LittleEndian {
		t.Errorf("Binary Order expected %T got %T", binary.LittleEndian, bo)
	}

	buf = []byte{0x4d, 0x4d, 0, 0x2a}
	bo = BinaryOrder(buf)
	if bo != binary.BigEndian {
		t.Errorf("Binary Order expected %T got %T", binary.BigEndian, bo)
	}
}
