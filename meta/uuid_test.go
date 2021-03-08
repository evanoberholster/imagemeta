package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	u := UUID{}
	assert.Equal(t, u, NilUUID, NilUUID)

	u2 := UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	u3, err := UUIDFromBytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	assert.ErrorIs(t, err, nil)

	u, err = UUIDFromBytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	assert.Error(t, err, assert.AnError)

	assert.Equal(t, u2, u3)
	str, _ := u2.MarshalText()

	assert.ErrorIs(t, u.UnmarshalText(str), nil)
	assert.Equal(t, u, u2)

	assert.Equal(t, str, []byte(u.String()), u.String())

	u = UUID{}
	buf, err := u2.MarshalBinary()
	assert.ErrorIs(t, err, nil)

	assert.ErrorIs(t, u.UnmarshalBinary(buf), nil)
	assert.Equal(t, u, u2)

	u = UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

}
