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
