package bmff

import (
	"bufio"
	"encoding/binary"
	"io"
	"strings"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// Errors
var (
	ErrBufLength                = errors.New("insufficient buffer length")
	ErrItemTypeWS               = errors.New("bufReader error: itemType doesn't end on whitespace")
	ErrRemainLengthInsufficient = errors.New("bufReader error: remain length insufficient")
	ErrBufReaderFlags           = errors.New("bufReader error: failed to read 4 bytes of Flags")
	ErrItemTypeLength           = errors.New("bufReader: insufficient itemType Length")
	errLargeBox                 = errors.New("unexpectedly large box")
	errUintSize                 = errors.New("invalid uintn read size")
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

// peek returns the next n bytes without advancing the reader.
// The bytes stop being valid at the next read call. If Peek returns fewer than n bytes,
// it also returns an error explaining why the read is short.
// The error is ErrBufferFull if n is larger than b's buffer size.
//
// peek is limited to the bufReaders "remain" bytes.
func (br *bufReader) peek(n int) ([]byte, error) {
	if br.remain < n {
		return nil, io.EOF
	}
	return br.Peek(n)
}

// discard skips the next n bytes, returning the number of bytes discarded.
//
// discard is limited to the bufReader's "remain" bytes.
func (br *bufReader) discard(n int) (err error) {
	if n == 0 {
		return nil
	}
	if br.remain < n {
		n = br.remain // limit discarded amount to remaining in bufReader
	}
	n, err = br.Discard(n)
	br.remain -= n
	if err != nil {
		return errors.Wrap(err, "bufReader discard")
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
		return "", errors.New("unexpected non-null terminated string")
	}
	return s, nil
}

func (br *bufReader) readUint8() (uint8, error) {
	if !br.anyRemain() {
		return 0, io.EOF
	}
	v, err := br.ReadByte()
	if err != nil {
		return 0, errors.Wrap(err, "readUint8")
	}
	br.remain-- // remove 1 remaining byte
	return v, nil
}

func (br *bufReader) readUint16() (uint16, error) {
	buf, err := br.peek(2)
	if err != nil {
		return 0, errors.Wrap(err, "readUint16")
	}
	v := binary.BigEndian.Uint16(buf[:2])
	return v, br.discard(2)
}

func (br *bufReader) readUint32() (uint32, error) {
	buf, err := br.peek(4)
	if err != nil {
		return 0, errors.Wrap(err, "readUint32")
	}
	v := binary.BigEndian.Uint32(buf[:4])
	return v, br.discard(4)
}

func (br *bufReader) readUintN(bits uint8) (uint64, error) {
	if bits == 0 {
		return 0, nil
	}

	nbyte := int(bits / 8)

	buf, err := br.peek(nbyte)
	if err != nil {
		return 0, errors.Wrap(err, "readUintN")
	}
	defer br.discard(nbyte)
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
		return 0, errors.Wrap(errUintSize, "readUintN")
	}
}

func (br *bufReader) readBrand() (b Brand, err error) {
	buf, err := br.peek(4)
	if err != nil {
		return brandUnknown, errors.Wrap(ErrBufLength, "readBrand")
	}
	return brand(buf[:4]), br.discard(4)
}

func (br *bufReader) readItemType() (it ItemType, err error) {
	buf, err := br.peek(5)
	if err != nil {
		return ItemTypeUnknown, errors.Wrap(ErrBufLength, "readItemType")
	}

	it = itemType(buf[:4])
	if buf[4] != '\x00' {
		// Read until whitespace
		//br.err = ErrItemTypeWS // errors.New("bufReader error: itemType doesn't end on whitespace")
		return it, br.discard(4)
	}

	return it, br.discard(5)
}

// readFlags reads the Flags from a FullBox header.
func (br *bufReader) readFlags() (f Flags, err error) {
	buf, err := br.peek(4)
	if err != nil {
		return f, errors.Wrap(ErrBufLength, "readFlags")
	}

	f = Flags(binary.BigEndian.Uint32(buf[:4]))

	return f, br.discard(4)
}

// readUUID reads a 16 byte UUID from the bufReader.
func (br *bufReader) readUUID() (u meta.UUID, err error) {
	buf, err := br.peek(16)
	if err != nil {
		return
	}
	u, err = meta.UUIDFromBytes(buf)
	if err != nil {
		err = errors.Wrap(err, "readUUID")
		return
	}
	return u, br.discard(16)
}

// readBox reads the next Box
func (br *bufReader) readBox() (b box, err error) {
	b = box{bufReader: *br}

	// Read box size and box type
	var buf []byte
	if buf, err = b.Peek(8); err != nil {
		return b, errors.Wrap(ErrBufLength, "readBox")
	}
	b.size = int64(binary.BigEndian.Uint32(buf[:4]))
	b.boxType = boxType(buf[4:8])

	switch b.size {
	case 1:
		// 1 means it's actually a 64-bit size, after the type.
		if buf, err = b.Peek(16); err != nil {
			return b, errors.Wrap(ErrBufLength, "readBox")
		}
		b.size = int64(binary.BigEndian.Uint64(buf[8:16]))
		if b.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return b, errors.Wrapf(errLargeBox, "box '%s'", b.boxType)
		}
		b.remain = int(b.size)
		return b, b.discard(16)
		//case 0:
		// 0 means unknown & to read to end of file. No more boxes.
		// r.noMoreBoxes = true
		// TODO: error
	}
	b.remain = int(b.size)
	return b, b.discard(8)
}
