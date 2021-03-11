package mknote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNikon(t *testing.T) {
	v := []byte("Nikon")
	assert.True(t, IsNikonMkNoteHeaderBytes(v))
}
