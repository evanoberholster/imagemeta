package mknote

import (
	"testing"
)

func TestIsNikon(t *testing.T) {
	v := []byte("Nikon")
	if !IsNikonMkNoteHeaderBytes(v) {
		t.Errorf("Error identifying NikonMkNoteHeaderBytes")
	}
}

func TestString(t *testing.T) {
	if TagCanonString(CanonAFInfo) != "CanonAFInfo" {
		t.Errorf("Expected %s got %s", "CanonAFInfo", TagCanonString(CanonAFInfo))
	}
	if TagCanonString(0x1234) != "0x1234" {
		t.Errorf("Expected %s got %s", "0x1234", TagCanonString(0x1234))
	}
}
