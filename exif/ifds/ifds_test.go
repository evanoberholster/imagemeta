package ifds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: write tests

func TestKey(t *testing.T) {
	key := NewKey(RootIFD, 1, TileWidth)

	assert.Equal(t, key, Key(0x1010142), "Ifd Key")

}
