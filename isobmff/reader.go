package isobmff

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Constants
const (
	bufReaderSize        = 64 * 1024
	seekDiscardThreshold = 64 * 1024
)

var (
	// bmffEndian ISOBMFF always uses BigEndian byteorder.
	// Can use either byteorder for Exif Information inside the ISOBMFF file
	bmffEndian = binary.BigEndian

	// crxEndian values are in BigEndian.
	crxEndian = binary.BigEndian
)

var readerPool = sync.Pool{
	New: func() any { return bufio.NewReaderSize(nil, bufReaderSize) },
}

// Reader is a ISO BMFF reader
type Reader struct {
	br *bufio.Reader

	ftyp FileTypeBox
	prvw PRVWBox
	heic HeicMeta

	exifReader         ExifReader
	xmpReader          XMPReader
	previewImageReader PreviewImageReader

	offset int

	source               io.Reader
	seeker               io.Seeker
	readerPool           *sync.Pool
	discardSeekThreshold int
}

// NewReader returns a new bmff.Reader
func NewReader(r io.Reader, exifReader ExifReader, xmpReader XMPReader, previewImageReader PreviewImageReader) Reader {
	// Reuse caller-provided bufio.Reader when large enough to avoid stacking buffers.
	if br, ok := r.(*bufio.Reader); ok && br.Size() >= bufReaderSize {
		reader := newReaderWithBufio(br, r, nil)
		reader.exifReader = exifReader
		reader.xmpReader = xmpReader
		reader.previewImageReader = previewImageReader
		return reader
	}

	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	reader := newReaderWithBufio(br, r, &readerPool)
	reader.exifReader = exifReader
	reader.xmpReader = xmpReader
	reader.previewImageReader = previewImageReader
	return reader
}

func newReaderWithBufio(br *bufio.Reader, source io.Reader, pool *sync.Pool) Reader {
	reader := Reader{
		br:                   br,
		source:               source,
		readerPool:           pool,
		discardSeekThreshold: seekDiscardThreshold,
	}
	if seeker, ok := source.(io.Seeker); ok {
		reader.seeker = seeker
	}
	return reader
}

func (r *Reader) peek(n int) ([]byte, error) {
	return r.br.Peek(n)
}

func (r *Reader) discard(n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}

	if r.seeker != nil && n >= r.discardSeekThreshold {
		discarded, err := r.discardWithSeek(n)
		r.offset += discarded
		return discarded, err
	}

	discarded, err := r.br.Discard(n)
	r.offset += discarded
	return discarded, err
}

func (r *Reader) discardWithSeek(n int) (discarded int, err error) {
	buffered := r.br.Buffered()
	if buffered > 0 {
		if buffered > n {
			buffered = n
		}
		discarded, err = r.br.Discard(buffered)
		n -= discarded
		if err != nil || n == 0 {
			return discarded, err
		}
	}

	// For large skips on seekable sources, avoid read-and-throw-away loops.
	if n >= r.discardSeekThreshold {
		if _, err = r.seeker.Seek(int64(n), io.SeekCurrent); err == nil {
			r.br.Reset(r.source)
			return discarded + n, nil
		}
	}

	skipped, err := r.br.Discard(n)
	return discarded + skipped, err
}

func (r *Reader) reset(newReader io.Reader) {
	exifReader := r.exifReader
	xmpReader := r.xmpReader
	previewImageReader := r.previewImageReader
	r.Close()
	*r = NewReader(newReader, exifReader, xmpReader, previewImageReader)
}

// Close the Reader. Returns the underlying bufio.Reader to the reader pool.
func (r *Reader) Close() {
	if r.readerPool != nil && r.br != nil {
		// Clear references to allow GC of the previous source quickly.
		r.br.Reset(nil)
		r.readerPool.Put(r.br)
	}
	r.br = nil
	r.source = nil
	r.seeker = nil
	r.readerPool = nil
}

// readBox reads an ISOBMFF box
func (r *Reader) readBox() (b box, err error) {
	// Read box size and box type (8-byte header)
	buf, err := r.peek(8)
	if err != nil {
		if err == io.EOF {
			return b, io.EOF
		}
		return b, fmt.Errorf("readBox: %w", ErrBufLength)
	}

	b.reader = r
	b.size = int64(bmffEndian.Uint32(buf[:4]))
	b.boxType = boxTypeFromBuf(buf[4:8])
	b.offset = r.offset

	headerSize := 8
	if b.size == 1 {
		buf, err = r.peek(16)
		if err != nil {
			return b, fmt.Errorf("readBox: %w", ErrBufLength)
		}

		// 1 means it's actually a 64-bit size, after the type.
		b.size = int64(bmffEndian.Uint64(buf[8:16]))
		if b.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return b, fmt.Errorf("readBox '%s': %w", b.boxType, errLargeBox)
		}
		headerSize = 16
	}
	if b.size < int64(headerSize) {
		return b, fmt.Errorf("readBox invalid size %d for '%s': %w", b.size, b.boxType, ErrBufLength)
	}

	maxInt := int64(^uint(0) >> 1)
	if b.size > maxInt {
		return b, fmt.Errorf("readBox '%s': %w", b.boxType, errLargeBox)
	}

	b.remain = int(b.size)
	_, err = b.Discard(headerSize)
	return b, err
}
