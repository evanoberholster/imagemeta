package meta

import (
	"bytes"
	"errors"
	"testing"
)

func TestUUID(t *testing.T) {
	ub, err := UUIDFromBytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	if err != nil {
		t.Errorf("UUID Error wanted nil, got: %s", err)
	}
	testList := []struct {
		u1 UUID
		u2 UUID
		b  bool
	}{
		{UUID{}, NilUUID, true},
		{UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, ub, true},
	}

	for _, l := range testList {
		if !bytes.EqualFold(l.u1[:], l.u2[:]) {
			t.Errorf("Incorrect UUID wanted: %s got: %s ", l.u1, l.u2)
		}
	}

	_, err = UUIDFromBytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
	if !(err != nil) {
		t.Errorf("Incorrect UUID error got: %s wanted: nil", err)
	}

	if err = ub.decodePlain([]byte("hello")); !errors.Is(err, ErrUUIDLength) {
		t.Errorf("Incorrect wanted: %s error got: %s", ErrUUIDLength, err)
	}

	//assert.Equal(t, u2, u3)
	//str, _ := u2.MarshalText()
	//
	//assert.ErrorIs(t, u.UnmarshalText(str), nil)
	//assert.Equal(t, u, u2)
	//
	//assert.Equal(t, str, []byte(u.String()), u.String())
	//
	//u = UUID{}
	//buf, err := u2.MarshalBinary()
	//assert.ErrorIs(t, err, nil)
	//
	//assert.ErrorIs(t, u.UnmarshalBinary(buf), nil)
	//assert.Equal(t, u, u2)
	//
	//u = UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

}

func TestUUIDUnmarshalText(t *testing.T) {
	var err error
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
		u, u2 := UUID{}, UUID{}
		if err := u.UnmarshalText([]byte(v.str)); !errors.Is(err, v.err) {
			t.Errorf("error is %v", err)
		}

		if !bytes.EqualFold(v.uuid[:], u[:]) {
			t.Errorf("Incorrect UUID wanted: %s got: %s", v.uuid, u)
		}

		if err = u.UnmarshalText([]byte(v.uuid.String())); err != nil {
			t.Error(err)
		}
		if err = u2.UnmarshalText([]byte(v.str)); !errors.Is(err, v.err) {
			t.Error(err)
		}

		if u.String() != u2.String() {
			t.Errorf("Incorrect UUID wanted: %s got: %s", u2.String(), u.String())
		}
		if buf, err := u.MarshalBinary(); err != nil {
			if bytes.Equal(v.uuid[:], buf) {
				t.Errorf("Incorrect UUID wanted: %s got: %s", v.uuid, buf)
			}
		}

	}

}
