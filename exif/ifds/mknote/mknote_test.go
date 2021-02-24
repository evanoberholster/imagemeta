package mknote

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNikonMkNoteHeader(t *testing.T) {
	// Test LittleEndian Nikon Header
	buf := []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'I', 'I', 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bo, err := NikonMkNoteHeader(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("Did not expect error %v %v", err, bo)
	}
	if bo != binary.LittleEndian {
		t.Errorf("Expected LittleEndian got %v", bo)
	}

	// Test BigEndian Nikon Header
	buf = []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bo, err = NikonMkNoteHeader(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("Did not expect error %v %v", err, bo)
	}
	if bo != binary.BigEndian {
		t.Errorf("Expected BigEndian got %v", bo)
	}

	// Test Short BigEndian Nikon Header
	buf = []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0}
	bo, err = NikonMkNoteHeader(bytes.NewReader(buf))
	if err != ErrNikonMkNote {
		t.Errorf("Did not expect nil %v %v", err, bo)
	}
	if bo != nil {
		t.Errorf("Expected nil got %v", bo)
	}

	// Test Incorrect BigEndian Nikon Header
	buf = []byte{'C', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0}
	bo, err = NikonMkNoteHeader(bytes.NewReader(buf))
	if err != ErrNikonMkNote {
		t.Errorf("Did not expect nil %v %v", err, bo)
	}
	if bo != nil {
		t.Errorf("Expected nil got %v", bo)
	}
}
