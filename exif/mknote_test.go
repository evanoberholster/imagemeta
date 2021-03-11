package exif

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNikonMkNoteHeader(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		bo   binary.ByteOrder
		err  error
	}{
		{"LittleEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'I', 'I', 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, binary.LittleEndian, nil},
		{"BigEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, binary.BigEndian, nil},
		{"Other Nikon Header Error", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, ErrNikonMkNote},
		{"Short BigEndian Nikon Header", []byte{'N', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0}, nil, ErrNikonMkNote},
		{"Incorrect BigEndian Nikon Header", []byte{'C', 'i', 'k', 'o', 'n', 0, 0, 0, 0, 0, 'M', 'M', 0, 0x2a, 0, 0, 0, 0, 0, 0}, nil, ErrNikonMkNote},
	}

	for _, v := range tests {
		bo, err := NikonMkNoteHeader(bytes.NewReader(v.buf))
		assert.ErrorIs(t, v.err, err, v.name)
		assert.Equal(t, v.bo, bo, v.name)
	}
}
