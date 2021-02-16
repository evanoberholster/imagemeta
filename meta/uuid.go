package meta

import (
	uuid "github.com/satori/go.uuid"
)

// UUID is a [16]byte Universally Unique Identifier (UUID).
// Based on github.com/satori/go.uuid
type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
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
