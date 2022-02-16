// Copyright (C) 2021 by Evan Oberholster
// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package meta

import (
	"bytes"
	"encoding/hex"

	"github.com/pkg/errors"
)

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}

	// ErrUUIDFormat is returned for an incorrect UUID format.
	ErrUUIDFormat = errors.New("uuid: incorrect UUID format")

	// ErrUUIDLength is returned for an incorrect UUID length.
	ErrUUIDLength = errors.New("uuid: incorrect UUID length")
)

// UUID is a 128 bits Universally Unique Identifier (UUID).
// Based on github.com/satori/go.uuid
type UUID [16]byte

// NilUUID is special form of UUID that is specified to have all
// 128 bits set to zero.
var NilUUID = UUID{}

// UUIDFromBytes returns UUID converted from raw byte slice input. It will return error if the slice isn't 16 bytes long.
func UUIDFromBytes(buf []byte) (UUID, error) {
	var u UUID
	err := u.UnmarshalBinary(buf)
	return UUID(u), err
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// String returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	str, _ := u.MarshalText()
	return string(str)
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (u UUID) MarshalText() (text []byte, err error) {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return buf, nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	return u.Bytes(), nil
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

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//   "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//   "6ba7b8109dad11d180b400c04fd430c8"
// ABNF for supported UUID text representation follows:
//   uuid := canonical | hashlike | braced | urn
//   plain := canonical | hashlike
//   canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//   hashlike := 12hexoct
//   braced := '{' plain '}'
//   urn := URN ':' UUID-NID ':' plain
//   URN := 'urn'
//   UUID-NID := 'uuid'
//   12hexoct := 6hexoct 6hexoct
//   6hexoct := 4hexoct 2hexoct
//   4hexoct := 2hexoct 2hexoct
//   2hexoct := hexoct hexoct
//   hexoct := hexdig hexdig
//   hexdig := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' |
//             'a' | 'b' | 'c' | 'd' | 'e' | 'f' |
//             'A' | 'B' | 'C' | 'D' | 'E' | 'F'
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	case 34, 38:
		return u.decodeBraced(text)
	case 41:
		fallthrough
	case 45:
		return u.decodeURN(text)
	default:
		return errors.Wrap(ErrUUIDLength, string(text))
	}
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return errors.Wrap(ErrUUIDFormat, string(t))
	}

	src := t[:]
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		if _, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup]); err != nil {
			*u = NilUUID
			return errors.Wrap(ErrUUIDFormat, err.Error())
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) (err error) {
	if _, err = hex.Decode(u[:], t[:]); err != nil {
		*u = NilUUID
		return errors.Wrap(ErrUUIDFormat, err.Error())
	}
	return
}

// decodeBraced decodes UUID string in format
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}" or in format
// "{6ba7b8109dad11d180b400c04fd430c8}".
func (u *UUID) decodeBraced(t []byte) (err error) {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return errors.Wrap(ErrUUIDFormat, string(t))
	}

	return u.decodePlain(t[1 : l-1])
}

// decodeURN decodes UUID string in format
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in format
// "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeURN(t []byte) (err error) {
	// t[:9] is urnUUIDPrefix
	if !bytes.Equal(t[:9], urnPrefix) {
		return errors.Wrap(ErrUUIDFormat, string(t))
	}

	return u.decodePlain(t[9:])
}

// decodePlain decodes UUID string in canonical format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in hash-like format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodePlain(t []byte) (err error) {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return errors.Wrap(ErrUUIDLength, string(t))
	}
}
