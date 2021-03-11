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

func TestUUIDUnmarshalText(t *testing.T) {
	//   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	//   "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	//   "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	//   "6ba7b8109dad11d180b400c04fd430c8"

	tests := []struct {
		uuid UUID
		str  string
		err  error
	}{
		// Errors
		{NilUUID, "6ba7b810-9dad-11d1-80b4-00c04fd430", ErrUUIDFormat},
		{NilUUID, "urn::6ba7b810-9dad-11d1-80b4-00c04fd430c8", ErrUUIDFormat},
		{NilUUID, "6ba7b810-9dad-11d1-80b4-00c04fd430xx", ErrUUIDFormat},
		{NilUUID, "6ba7b810-9dad-11d1-80b4b00c04fd430xx", ErrUUIDFormat},
		{NilUUID, "6ba7b8109dad11d1804b00c04fd430xx", ErrUUIDFormat},
		{NilUUID, "{6ba7b810-9dad-11d1-80b4-00c04xx430c8}", ErrUUIDFormat},
		{NilUUID, "{6ba7b810-9dad-11d1-80b4-00c04xx430c8)", ErrUUIDFormat},
		{NilUUID, "{6ba7b8109dad11d180b400c04fdxx0c8}", ErrUUIDFormat},
		{NilUUID, "{6ba7b8109dad11d180b400c04fdxx0c8abcddddabc}", ErrUUIDLength},
		{NilUUID, "{6ba7b8109dad11d180b400}", ErrUUIDLength},
		//
		{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x0, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, "{6ba7b8109dad11d180b400c04fd430c8}", nil},
		{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x0, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil},
		{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x0, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}", nil},
		{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x0, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil},
		{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x0, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, "6ba7b8109dad11d180b400c04fd430c8", nil},
	}
	for _, v := range tests {
		u := UUID{}
		err := u.UnmarshalText([]byte(v.str))
		if assert.ErrorIs(t, err, v.err) {
			assert.Equal(t, v.uuid, u)
		}
	}

}
