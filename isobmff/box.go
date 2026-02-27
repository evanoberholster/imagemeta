package isobmff

import (
	"bytes"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// box is a bounded view over an ISOBMFF box payload.
// Nested boxes share the same underlying Reader and enforce limits via remain.
type box struct {
	size    int64
	remain  int
	offset  int
	flags   flags
	boxType boxType
	outer   *box
	reader  *Reader
}

// isType reports whether the box has the expected type.
func (b box) isType(bt boxType) bool { return b.boxType == bt }

// Peek returns bytes without advancing the read position.
// Access is constrained to the current box bounds.
func (b *box) Peek(n int) ([]byte, error) {
	if b.remain >= n {
		if b.outer != nil {
			return b.outer.Peek(n)
		}
		return b.reader.peek(n)
	}
	return nil, ErrRemainLengthInsufficient
}

// Discard advances by n bytes, bounded by the current box.
func (b *box) Discard(n int) (int, error) {
	if b.remain >= n {
		b.remain -= n
		if b.outer != nil {
			return b.outer.Discard(n)
		}
		return b.reader.discard(n)
	}
	return 0, ErrRemainLengthInsufficient
}

// Read copies bytes from the underlying reader while respecting box bounds.
func (b *box) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.remain == 0 {
		return 0, io.EOF
	}

	readLen := len(p)
	if readLen > b.remain {
		readLen = b.remain
	}

	n, err = b.reader.br.Read(p[:readLen])
	b.adjust(n)
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return n, err
}

func (b *box) adjust(n int) {
	if n > b.remain {
		n = b.remain
	}
	b.remain -= n
	if b.outer != nil {
		b.outer.adjust(n)
	}
}

// close discards any unread bytes in the current box.
func (b *box) close() error {
	if b.remain == 0 {
		return nil
	}
	_, err := b.Discard(b.remain)
	return err
}

// parseBoxSizeAndType parses the first 8 bytes of a BMFF box header:
// 32-bit size followed by 32-bit type (FourCC).
func parseBoxSizeAndType(buf []byte) (size int64, bt boxType, err error) {
	if len(buf) < 8 {
		return 0, typeUnknown, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	size = int64(bmffEndian.Uint32(buf[:4]))
	bt = boxTypeFromBuf(buf[4:8])
	return size, bt, nil
}

// parseExtendedBoxSize parses a BMFF "largesize" header (size32 == 1),
// where bytes 8..15 contain the 64-bit box size.
func parseExtendedBoxSize(buf []byte, bt boxType) (int64, error) {
	if len(buf) < 16 {
		return 0, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	maxInt := uint64(^uint(0) >> 1)
	size := bmffEndian.Uint64(buf[8:16])
	if size > maxInt {
		return 0, fmt.Errorf("readBox '%s': %w", bt, errLargeBox)
	}
	return int64(size), nil
}

// validateBoxSize ensures the declared box size is sane for this parser:
// it must include at least the header and fit in host int width.
func validateBoxSize(size int64, headerSize int, bt boxType) error {
	if size < int64(headerSize) {
		return fmt.Errorf("readBox invalid size %d for '%s': %w", size, bt, ErrBufLength)
	}

	maxInt := int64(^uint(0) >> 1)
	if size > maxInt {
		return fmt.Errorf("readBox '%s': %w", bt, errLargeBox)
	}
	return nil
}

// readInnerBox reads the next child box header within the current container and
// returns a child view constrained to that box's byte range.
func (b *box) readInnerBox() (inner box, next bool, err error) {
	if b.remain < 8 {
		return inner, false, nil
	}
	buf, err := b.Peek(8)
	if err != nil {
		return inner, false, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	size, boxType, err := parseBoxSizeAndType(buf)
	if err != nil {
		return inner, false, err
	}
	headerSize := 8
	if size == 1 {
		buf, err = b.Peek(16)
		if err != nil {
			return inner, false, fmt.Errorf("readBox: %w", ErrBufLength)
		}
		size, err = parseExtendedBoxSize(buf, boxType)
		if err != nil {
			return inner, false, err
		}
		headerSize = 16
	}
	if err = validateBoxSize(size, headerSize, boxType); err != nil {
		return inner, false, err
	}

	inner.reader = b.reader
	inner.outer = b
	inner.offset = int(b.size) - b.remain + b.offset
	inner.size = size
	inner.boxType = boxType

	inner.remain = int(inner.size)
	_, err = inner.Discard(headerSize)
	return inner, true, err
}

// readUint16 from box
func (b *box) readUint16() (uint16, error) {
	buf, err := b.Peek(2)
	if err != nil {
		return 0, fmt.Errorf("readUint16: %w", ErrBufLength)
	}
	_, err = b.Discard(2)
	return bmffEndian.Uint16(buf[:2]), err
}

// readUint32 from box
func (b *box) readUint32() (uint32, error) {
	buf, err := b.Peek(4)
	if err != nil {
		return 0, fmt.Errorf("readUint32: %w", ErrBufLength)
	}
	_, err = b.Discard(4)
	return bmffEndian.Uint32(buf[:4]), err
}

// readUintN reads a 0/1/2/4/8-byte unsigned integer from the box.
func (b *box) readUintN(size uint8) (uint64, error) {
	switch size {
	case 0:
		return 0, nil
	case 1:
		buf, err := b.Peek(1)
		if err != nil {
			return 0, fmt.Errorf("readUintN: %w", ErrBufLength)
		}
		_, err = b.Discard(1)
		return uint64(buf[0]), err
	case 2:
		v, err := b.readUint16()
		return uint64(v), err
	case 4:
		v, err := b.readUint32()
		return uint64(v), err
	case 8:
		buf, err := b.Peek(8)
		if err != nil {
			return 0, fmt.Errorf("readUintN: %w", ErrBufLength)
		}
		_, err = b.Discard(8)
		return bmffEndian.Uint64(buf[:8]), err
	default:
		return 0, ErrUnsupportedFieldSize
	}
}

func (b *box) readFourCC() (uint32, error) {
	return b.readUint32()
}

// discardCString discards bytes until a NUL terminator is reached.
func (b *box) discardCString(maxLen int) error {
	if maxLen <= 0 {
		maxLen = maxBoxStringLength
	}
	discarded := 0
	for {
		if b.remain == 0 {
			return ErrBufLength
		}
		chunk := b.remain
		if chunk > boxStringReadChunk {
			chunk = boxStringReadChunk
		}
		buf, err := b.Peek(chunk)
		if err != nil {
			return fmt.Errorf("discardCString: %w", ErrBufLength)
		}
		if idx := bytes.IndexByte(buf, 0); idx >= 0 {
			_, err = b.Discard(idx + 1)
			return err
		}
		discarded += chunk
		if discarded > maxLen {
			return ErrBoxStringTooLong
		}
		if _, err = b.Discard(chunk); err != nil {
			return err
		}
	}
}

// readCString reads bytes until a NUL terminator is reached.
func (b *box) readCString(maxLen int) (string, error) {
	if maxLen <= 0 {
		maxLen = maxBoxStringLength
	}
	var out []byte
	for {
		if b.remain == 0 {
			return "", ErrBufLength
		}
		chunk := b.remain
		if chunk > boxStringReadChunk {
			chunk = boxStringReadChunk
		}
		buf, err := b.Peek(chunk)
		if err != nil {
			return "", fmt.Errorf("readCString: %w", ErrBufLength)
		}
		if idx := bytes.IndexByte(buf, 0); idx >= 0 {
			if len(out)+idx > maxLen {
				return "", ErrBoxStringTooLong
			}
			if out == nil {
				_, err = b.Discard(idx + 1)
				return string(buf[:idx]), err
			}
			out = append(out, buf[:idx]...)
			_, err = b.Discard(idx + 1)
			return string(out), err
		}
		if len(out)+chunk > maxLen {
			return "", ErrBoxStringTooLong
		}
		if out == nil {
			out = make([]byte, 0, chunk*2)
		}
		out = append(out, buf[:chunk]...)
		if _, err = b.Discard(chunk); err != nil {
			return "", err
		}
	}
}

// readUUID reads a 16 byte UUID from the box.
func (b *box) readUUID() (u meta.UUID, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		return u, fmt.Errorf("readUUID: %w", ErrBufLength)
	}
	if err = u.UnmarshalBinary(buf); err != nil {
		return u, err
	}
	_, err = b.Discard(16)
	return u, err
}
