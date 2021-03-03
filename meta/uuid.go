package meta

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// NilUUID is a empty UUID. All zeros.
var NilUUID = UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// UUID is a [16]byte Universally Unique Identifier (UUID).
// Based on github.com/satori/go.uuid
type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// UUIDFromBytes returns UUID converted from raw byte slice input. It will return error if the slice isn't 16 bytes long.
func UUIDFromBytes(buf []byte) (UUID, error) {
	u, err := uuid.FromBytes(buf)
	return UUID(u), err
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (u UUID) MarshalText() (text []byte, err error) {
	return uuid.UUID(u).MarshalText()
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (u *UUID) UnmarshalText(text []byte) (err error) {
	uid, err := uuid.FromString(string(text))
	*u = UUID(uid)
	return err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 16 {
		err = errors.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}
