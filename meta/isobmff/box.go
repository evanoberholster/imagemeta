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
	remain  int64
	offset  int64
	flags   flags
	boxType boxType
	outer   *box
	reader  *Reader
}

const maxInt64Value = int64(^uint64(0) >> 1)

// isType reports whether the box has the expected type.
func (b box) isType(bt boxType) bool { return b.boxType == bt }

// Peek returns bytes without advancing the read position.
// Access is constrained to the current box bounds.
func (b *box) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrBufLength
	}
	if b.remain >= int64(n) {
		if b.outer != nil {
			return b.outer.Peek(n)
		}
		return b.reader.peek(n)
	}
	return nil, ErrRemainLengthInsufficient
}

// Discard advances by n bytes, bounded by the current box.
func (b *box) Discard(n int) (int, error) {
	if n < 0 {
		return 0, ErrBufLength
	}
	if n == 0 {
		return 0, nil
	}
	if b.remain >= int64(n) {
		var (
			discarded int
			err       error
		)
		if b.outer != nil {
			discarded, err = b.outer.Discard(n)
		} else {
			discarded, err = b.reader.discard(n)
		}
		b.remain -= int64(discarded)
		if b.remain < 0 {
			b.remain = 0
		}
		return discarded, err
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
	if int64(readLen) > b.remain {
		readLen = int(b.remain)
	}

	n, err = b.reader.br.Read(p[:readLen])
	b.adjust(int64(n))
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return n, err
}

func (b *box) adjust(n int64) {
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
	return discardBoxBytes(b, b.remain)
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

// parseBoxSizeAndType reads and parses the next 8 bytes as a BMFF box header.
func (b *box) parseBoxSizeAndType() (size int64, bt boxType, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return 0, typeUnknown, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	return parseBoxSizeAndType(buf)
}

// parseExtendedBoxSize parses a BMFF "largesize" header (size32 == 1),
// where bytes 8..15 contain the 64-bit box size.
func parseExtendedBoxSize(buf []byte, bt boxType) (int64, error) {
	if len(buf) < 16 {
		return 0, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	size := bmffEndian.Uint64(buf[8:16])
	if size > uint64(maxInt64Value) {
		return 0, fmt.Errorf("readBox '%s': %w", bt, errLargeBox)
	}
	return int64(size), nil
}

// parseExtendedBoxSize reads and parses a 16-byte BMFF extended-size header.
func (b *box) parseExtendedBoxSize(bt boxType) (int64, error) {
	buf, err := b.Peek(16)
	if err != nil {
		return 0, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	return parseExtendedBoxSize(buf, bt)
}

// validateBoxSize ensures the declared box size is sane for this parser:
// it must include at least the header and fit in host int width.
func validateBoxSize(size int64, headerSize int, bt boxType) error {
	if size < int64(headerSize) {
		return fmt.Errorf("readBox invalid size %d for '%s': %w", size, bt, ErrBufLength)
	}

	if size > maxInt64Value {
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
	size, boxType, err := b.parseBoxSizeAndType()
	if err != nil {
		return inner, false, err
	}
	headerSize := 8
	if size == 1 {
		size, err = b.parseExtendedBoxSize(boxType)
		if err != nil {
			return inner, false, err
		}
		headerSize = 16
	}
	if err = validateBoxSize(size, headerSize, boxType); err != nil {
		return inner, false, err
	}

	inner = box{
		reader:  b.reader,
		outer:   b,
		offset:  b.offset + b.size - b.remain,
		size:    size,
		boxType: boxType,
		remain:  size,
	}
	_, err = inner.Discard(headerSize)
	return inner, true, err
}

// readUint16 reads a big-endian uint16 and advances the box cursor.
func (b *box) readUint16() (uint16, error) {
	buf, err := b.Peek(2)
	if err != nil {
		return 0, fmt.Errorf("readUint16: %w", ErrBufLength)
	}
	_, err = b.Discard(2)
	return bmffEndian.Uint16(buf[:2]), err
}

// readUint32 reads a big-endian uint32 and advances the box cursor.
func (b *box) readUint32() (uint32, error) {
	buf, err := b.Peek(4)
	if err != nil {
		return 0, fmt.Errorf("readUint32: %w", ErrBufLength)
	}
	_, err = b.Discard(4)
	return bmffEndian.Uint32(buf[:4]), err
}

// readUint64 reads a big-endian uint64 and advances the box cursor.
func (b *box) readUint64() (uint64, error) {
	buf, err := b.Peek(8)
	if err != nil {
		return 0, fmt.Errorf("readUint64: %w", ErrBufLength)
	}
	_, err = b.Discard(8)
	return bmffEndian.Uint64(buf[:8]), err
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
		return b.readUint64()
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
		if chunk > int64(boxStringReadChunk) {
			chunk = int64(boxStringReadChunk)
		}
		buf, err := b.Peek(int(chunk))
		if err != nil {
			return fmt.Errorf("discardCString: %w", ErrBufLength)
		}
		if idx := bytes.IndexByte(buf, 0); idx >= 0 {
			if discarded+idx > maxLen {
				return ErrBoxStringTooLong
			}
			_, err = b.Discard(idx + 1)
			return err
		}
		discarded += int(chunk)
		if discarded > maxLen {
			return ErrBoxStringTooLong
		}
		if _, err = b.Discard(int(chunk)); err != nil {
			return err
		}
	}
}

// readCStringBytes reads bytes until a NUL terminator and appends into dst.
// Returned bytes are the same underlying slice as dst.
func (b *box) readCStringBytes(dst []byte, maxLen int) ([]byte, error) {
	if maxLen <= 0 {
		maxLen = maxBoxStringLength
	}
	dst = dst[:0]
	for {
		if b.remain == 0 {
			return dst, ErrBufLength
		}
		chunk := b.remain
		if chunk > int64(boxStringReadChunk) {
			chunk = int64(boxStringReadChunk)
		}
		buf, err := b.Peek(int(chunk))
		if err != nil {
			return dst, fmt.Errorf("readCString: %w", ErrBufLength)
		}
		if idx := bytes.IndexByte(buf, 0); idx >= 0 {
			if len(dst)+idx > maxLen {
				return dst, ErrBoxStringTooLong
			}
			dst = append(dst, buf[:idx]...)
			_, err = b.Discard(idx + 1)
			return dst, err
		}
		if len(dst)+int(chunk) > maxLen {
			return dst, ErrBoxStringTooLong
		}
		dst = append(dst, buf[:int(chunk)]...)
		if _, err = b.Discard(int(chunk)); err != nil {
			return dst, err
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

func discardBoxBytes(b *box, n int64) error {
	for n > 0 {
		chunk := n
		if chunk > int64(^uint(0)>>1) {
			chunk = int64(^uint(0) >> 1)
		}
		discarded, err := b.Discard(int(chunk))
		n -= int64(discarded)
		if err != nil {
			return err
		}
		if discarded == 0 {
			return io.ErrNoProgress
		}
	}
	return nil
}
