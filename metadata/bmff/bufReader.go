package bmff

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Errors
var (
	ErrItemTypeWS      = errors.New("bufReader error: itemType doesn't end on whitespace")
	ErrBufReaderLength = errors.New("bufReader error: infufficient length")
)

// bufReader adds some HEIF/BMFF-specific methods around a *bufio.Reader.
type bufReader struct {
	*bufio.Reader
	err    error
	remain int
}

// ok reports whether all previous reads have been error-free.
func (br *bufReader) ok() bool { return br.err == nil }

func (br *bufReader) anyRemain() bool {
	return br.remain > 0 && br.ok()
}

func (br *bufReader) discard(n int) error {
	_, err := br.Discard(n)
	br.remain -= n
	if err != nil {
		br.err = fmt.Errorf("bufReader discard error: %v", err)
		return br.err
	}
	return err
}

func (br *bufReader) readString() (string, error) {
	if br.err != nil {
		return "", br.err
	}
	s0, err := br.ReadString(0)
	if err != nil {
		br.err = err
		return "", err
	}
	br.remain -= len(s0)
	if s0[len(s0)-1] == '\x00' {
		s0 = s0[:len(s0)-1]
		return string(s0), nil
	}
	s := strings.TrimSuffix(s0, "\x00")
	if len(s) == len(s0) {
		err = fmt.Errorf("unexpected non-null terminated string")
		br.err = err
		return "", err
	}
	return s, nil
}

func (br *bufReader) readUint8() (uint8, error) {
	if !br.anyRemain() {
		return 0, ErrBufReaderLength
	}
	v, err := br.ReadByte()
	if err != nil {
		br.err = err
		return 0, err
	}
	br.remain-- // remove 1 remaining byte
	return v, nil
}

func (br *bufReader) readUint16() (uint16, error) {
	if br.err != nil {
		return 0, br.err
	}
	if br.remain < 2 {
		return 0, ErrBufReaderLength
	}
	buf, err := br.Peek(2)
	if err != nil {
		br.err = err
		return 0, err
	}
	v := binary.BigEndian.Uint16(buf[:2])
	return v, br.discard(2)
}

func (br *bufReader) readUint32() (uint32, error) {
	if br.err != nil {
		return 0, br.err
	}
	if br.remain < 4 {
		return 0, ErrBufReaderLength
	}
	buf, err := br.Peek(4)
	if err != nil {
		br.err = err
		return 0, err
	}
	v := binary.BigEndian.Uint32(buf[:4])
	return v, br.discard(4)
}

func (br *bufReader) readUintN(bits uint8) (uint64, error) {
	if br.err != nil {
		return 0, br.err
	}
	if br.remain < int(bits/8) {
		return 0, ErrBufReaderLength
	}
	if bits == 0 {
		return 0, nil
	}
	nbyte := bits / 8
	buf, err := br.Peek(int(nbyte))
	if err != nil {
		br.err = err
		return 0, err
	}
	defer br.discard(int(nbyte))
	switch bits {
	case 8:
		return uint64(buf[0]), nil
	case 16:
		return uint64(binary.BigEndian.Uint16(buf[:2])), nil
	case 32:
		return uint64(binary.BigEndian.Uint32(buf[:4])), nil
	case 64:
		return binary.BigEndian.Uint64(buf[:8]), nil
	default:
		br.err = fmt.Errorf("invalid uintn read size")
		return 0, br.err
	}
}

func (br *bufReader) readBrand() (b Brand, err error) {
	if br.err != nil {
		err = br.err
		return
	}
	if br.remain < 4 {
		br.err = ErrBufReaderLength
		return brandUnknown, br.err
	}
	var buf []byte
	if buf, err = br.Peek(4); err != nil {
		return
	}
	return brand(buf[:4]), br.discard(4)
}

func (br *bufReader) readItemType() (it ItemType, err error) {
	if br.remain < 4 {
		br.err = ErrBufReaderLength
		return ItemTypeUnknown, br.err
	}
	buf, err := br.Peek(5)
	if err != nil {
		return ItemTypeUnknown, err
	}

	it = itemType(buf[:5])
	if buf[4] != '\x00' {
		// Read until whitespace
		//br.err = ErrItemTypeWS // errors.New("bufReader error: itemType doesn't end on whitespace")
		return it, br.discard(4)
	}

	return it, br.discard(5)
}

func (br *bufReader) readFlags() (f Flags, err error) {
	if br.remain < 4 {
		err = ErrBufReaderLength
		br.err = err
	}
	// Parse Flags from a FullBox header.
	buf, err := br.Peek(4)
	if err != nil {
		return f, fmt.Errorf("failed to read 4 bytes of Flags: %v", err)
	}

	f = Flags(binary.BigEndian.Uint32(buf[:4]))

	return f, br.discard(4)
}

func (br *bufReader) readInnerBox() (b box, err error) {
	b = box{bufReader: *br}

	// Read box size and box type
	var buf []byte
	if buf, err = b.Peek(8); err != nil {
		return b, err
	}
	b.size = int64(binary.BigEndian.Uint32(buf[:4]))
	b.boxType = boxType(buf[4:8])

	if err = b.discard(8); err != nil {
		return
	}

	var remain int
	switch b.size {
	case 1:
		// 1 means it's actually a 64-bit size, after the type.
		if buf, err = b.Peek(8); err != nil {
			return b, err
		}
		b.size = int64(binary.BigEndian.Uint64(buf[:8]))
		if b.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return b, fmt.Errorf("unexpectedly large box %q", b.boxType)
		}
		remain = int(b.size - 2*4 - 8)
		if err = b.discard(8); err != nil {
			// TODO: write error message
			return
		}
	case 0:
		// 0 means unknown & to read to end of file. No more boxes.
		// r.noMoreBoxes = true
		// TODO: error
	default:
		remain = int(b.size - 2*4)
	}
	b.remain = remain
	return b, nil
}
