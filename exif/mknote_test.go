package exif

import (
	"bytes"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/stretchr/testify/assert"
)

func TestNikonMkNoteHeader(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		bo   meta.ByteOrder
		err  error
	}{
		{"LittleEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'I', 'I', 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, meta.LittleEndian, nil},
		{"BigEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, meta.BigEndian, nil},
		{"Other Nikon Header Error", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, meta.UnknownEndian, ErrNikonMkNote},
		{"Short BigEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0}, meta.UnknownEndian, ErrNikonMkNote},
		{"Incorrect BigEndian Nikon Header", []byte{'C', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0}, meta.UnknownEndian, ErrNikonMkNote},
	}

	for _, v := range tests {
		bo, err := NikonMkNoteHeader(bytes.NewReader(v.buf))
		assert.ErrorIs(t, v.err, err, v.name)
		assert.Equal(t, v.bo, bo, v.name)
	}
}
