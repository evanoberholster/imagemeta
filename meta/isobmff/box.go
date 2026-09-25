package isobmff

import (
	"bytes"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// box is a bounded view over an ISOBMFF box payload.
// Nested boxes share the same underlying Reader and enforce limits via remain.
//
//	size is the total declared box size including its header.
//	remain counts the still-unread bytes including any unconsumed header;
//	it only shrinks as the box (or its children) advance.
//	offset is the absolute stream offset of the box start.
type box struct {
	size    int
	remain  int
	offset  int64
	flags   flags
	boxType boxType
	outer   *box
	reader  *Reader
}

const maxIntValue = int(^uint(0) >> 1)

// isType reports whether the box has the expected type.
func (b *box) isType(bt boxType) bool { return b.boxType == bt }

// Peek returns bytes without advancing the read position.
// Access is constrained to the current box bounds.
func (b *box) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrBufLength
	}
	if b.remain >= n {
		if b.outer != nil {
			return b.outer.Peek(n)
		}
		return b.reader.peek(n)
	}
	return nil, ErrRemainLengthInsufficient
}

// consume returns the next n bytes and advances past them in one step.
// It fuses the Peek/Discard pair used by fixed-width field reads into a
// single bounds check and a single delegation hop.
//
// Only the leaf remain is checked: boxes are processed strictly in order
// (a child is always closed before its siblings proceed), so an active
// child never outlives its ancestors' bounds. adjust() still propagates
// the consumption up the whole chain.
func (b *box) consume(n int) ([]byte, error) {
	if n < 0 || b.remain < n {
		if n < 0 {
			return nil, ErrBufLength
		}
		return nil, ErrRemainLengthInsufficient
	}
	buf, err := b.reader.consume(n)
	if err != nil {
		return nil, err
	}
	b.adjust(n)
	return buf, nil
}

// Discard advances by n bytes, bounded by the current box.
func (b *box) Discard(n int) (int, error) {
	if n < 0 {
		return 0, ErrBufLength
	}
	if n == 0 {
		return 0, nil
	}
	if b.remain >= n {
		var (
			discarded int
			err       error
		)
		if b.outer != nil {
			discarded, err = b.outer.Discard(n)
		} else {
			discarded, err = b.reader.discard(n)
		}
		b.remain -= discarded
		return discarded, err
	}
	return 0, ErrRemainLengthInsufficient
}

// Read copies bytes from the underlying reader while respecting box bounds.
// Consumed bytes advance the reader's absolute offset exactly like Discard,
// so later boxes resolve correct absolute positions after callbacks Read.
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
	if n > 0 {
		b.reader.offset += int64(n)
	}
	// adjust propagates to outer boxes; reading straight from the shared
	// buffered reader (instead of delegating through outer) saves a hop
	// per call with identical accounting.
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
	return discardBoxBytes(b, b.remain)
}

// parseBoxSizeAndType parses the first 8 bytes of a BMFF box header:
// 32-bit size followed by 32-bit type (FourCC).
func parseBoxSizeAndType(buf []byte) (size int, bt boxType, err error) {
	if len(buf) < 8 {
		return 0, typeUnknown, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	size32 := bmffEndian.Uint32(buf[:4])
	if uint64(size32) > uint64(maxIntValue) {
		return 0, typeUnknown, errLargeBox
	}
	size = int(size32)
	bt = boxTypeFromBuf(buf[4:8])
	return size, bt, nil
}

// parseBoxSizeAndType consumes the next 8 bytes as a BMFF box header.
func (b *box) parseBoxSizeAndType() (size int, bt boxType, err error) {
	buf, err := b.consume(8)
	if err != nil {
		return 0, typeUnknown, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	return parseBoxSizeAndType(buf)
}

// parseExtendedBoxSize parses a BMFF "largesize" value: the 8 bytes
// following a size32 == 1 header.
func parseExtendedBoxSize(ext []byte, bt boxType) (int, error) {
	if len(ext) < 8 {
		return 0, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	size := bmffEndian.Uint64(ext[:8])
	if size > uint64(maxIntValue) {
		return 0, fmt.Errorf("readBox '%s': %w", bt, errLargeBox)
	}
	return int(size), nil
}

// parseExtendedBoxSize consumes a 16-byte BMFF extended-size header's
// 8-byte largesize suffix (the leading 8 bytes were already consumed).
func (b *box) parseExtendedBoxSize(bt boxType) (int, error) {
	ext, err := b.consume(8)
	if err != nil {
		return 0, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	return parseExtendedBoxSize(ext, bt)
}

// validateBoxSize ensures the declared box size is sane for this parser:
// it must include at least the header and fit in host int width.
func validateBoxSize(size int, headerSize int, bt boxType) error {
	if size < headerSize {
		return fmt.Errorf("readBox invalid size %d for '%s': %w", size, bt, ErrBufLength)
	}
	return nil
}

// readInnerBox reads the next child box header within the current container and
// returns a child view constrained to that box's byte range.
func (b *box) readInnerBox() (inner box, next bool, err error) {
	if b.remain < 8 {
		return inner, false, nil
	}
	size, bt, err := b.parseBoxSizeAndType()
	if err != nil {
		return inner, false, err
	}
	headerSize := 8
	if size == 1 {
		size, err = b.parseExtendedBoxSize(bt)
		if err != nil {
			return inner, false, err
		}
		headerSize = 16
	}
	if size == 0 {
		// Nested zero-size boxes are malformed per spec (0 means end of
		// file), so bound them by the outer container instead of failing.
		// The header is already consumed; the child takes all that remains.
		size = headerSize + b.remain
	}
	if err = validateBoxSize(size, headerSize, bt); err != nil {
		return inner, false, err
	}

	inner = box{
		reader:  b.reader,
		outer:   b,
		offset:  b.offset + int64(b.size-b.remain),
		size:    size,
		boxType: bt,
		remain:  size - headerSize,
	}
	return inner, true, nil
}

// readUint16 reads a big-endian uint16 and advances the box cursor.
func (b *box) readUint16() (uint16, error) {
	buf, err := b.consume(2)
	if err != nil {
		return 0, fmt.Errorf("readUint16: %w", ErrBufLength)
	}
	return bmffEndian.Uint16(buf), nil
}

// readUint32 reads a big-endian uint32 and advances the box cursor.
func (b *box) readUint32() (uint32, error) {
	buf, err := b.consume(4)
	if err != nil {
		return 0, fmt.Errorf("readUint32: %w", ErrBufLength)
	}
	return bmffEndian.Uint32(buf), nil
}

// readUint64 reads a big-endian uint64 and advances the box cursor.
func (b *box) readUint64() (uint64, error) {
	buf, err := b.consume(8)
	if err != nil {
		return 0, fmt.Errorf("readUint64: %w", ErrBufLength)
	}
	return bmffEndian.Uint64(buf), nil
}

// readUintN reads a 0/1/2/4/8-byte unsigned integer from the box.
func (b *box) readUintN(size uint8) (uint64, error) {
	switch size {
	case 0:
		return 0, nil
	case 1:
		buf, err := b.consume(1)
		if err != nil {
			return 0, fmt.Errorf("readUintN: %w", ErrBufLength)
		}
		return uint64(buf[0]), nil
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
		if chunk > boxStringReadChunk {
			chunk = boxStringReadChunk
		}
		buf, err := b.Peek(chunk)
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
		discarded += chunk
		if discarded > maxLen {
			return ErrBoxStringTooLong
		}
		if _, err = b.Discard(chunk); err != nil {
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
		if chunk > boxStringReadChunk {
			chunk = boxStringReadChunk
		}
		buf, err := b.Peek(chunk)
		if err != nil {
			return dst, fmt.Errorf("readCStringBytes: %w", ErrBufLength)
		}
		if idx := bytes.IndexByte(buf, 0); idx >= 0 {
			if len(dst)+idx > maxLen {
				return dst, ErrBoxStringTooLong
			}
			dst = append(dst, buf[:idx]...)
			_, err = b.Discard(idx + 1)
			return dst, err
		}
		if len(dst)+chunk > maxLen {
			return dst, ErrBoxStringTooLong
		}
		dst = append(dst, buf[:chunk]...)
		if _, err = b.Discard(chunk); err != nil {
			return dst, err
		}
	}
}

// readUUID reads a 16 byte UUID from the box.
func (b *box) readUUID() (u meta.UUID, err error) {
	buf, err := b.consume(16)
	if err != nil {
		return u, fmt.Errorf("readUUID: %w", ErrBufLength)
	}
	if err = u.UnmarshalBinary(buf); err != nil {
		return u, err
	}
	return u, nil
}

func discardBoxBytes(b *box, n int) error {
	for n > 0 {
		chunk := n
		discarded, err := b.Discard(chunk)
		n -= discarded
		if err != nil {
			return err
		}
		if discarded == 0 {
			return io.ErrNoProgress
		}
	}
	return nil
}
