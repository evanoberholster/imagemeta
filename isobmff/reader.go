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

	goalExif bool
	goalXMP  bool
	goalTHMB bool
	goalPRVW bool

	haveExif bool
	haveXMP  bool
	haveTHMB bool
	havePRVW bool

	stopAfterMetadata bool
	goalsInitialized  bool

	offset int

	source               io.Reader
	seeker               io.Seeker
	readerPool           *sync.Pool
	discardSeekThreshold int
}

// NewReader returns a new bmff.Reader
func NewReader(r io.Reader, exifReader ExifReader, xmpReader XMPReader, previewImageReader PreviewImageReader) *Reader {
	reader := newReader(r)
	reader.exifReader = exifReader
	reader.xmpReader = xmpReader
	reader.previewImageReader = previewImageReader
	return &reader
}

func newReader(r io.Reader) Reader {
	// Reuse caller-provided bufio.Reader when large enough to avoid stacking buffers.
	if br, ok := r.(*bufio.Reader); ok && br.Size() >= bufReaderSize {
		reader := newReaderWithBufio(br, r, nil)
		return reader
	}

	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	reader := newReaderWithBufio(br, r, &readerPool)
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

func (r *Reader) reset(newSource io.Reader) {
	exifReader := r.exifReader
	xmpReader := r.xmpReader
	previewImageReader := r.previewImageReader
	r.Close()
	*r = newReader(newSource)
	r.exifReader = exifReader
	r.xmpReader = xmpReader
	r.previewImageReader = previewImageReader
}

func (r *Reader) initMetadataGoals() {
	r.goalExif = r.exifReader != nil
	r.goalXMP = r.xmpReader != nil
	r.goalTHMB = false
	r.goalPRVW = false
	if r.previewImageReader != nil {
		if r.ftyp.MajorBrand == brandCrx {
			r.goalTHMB = true
			r.goalPRVW = true
		} else {
			r.goalPRVW = true
		}
	}

	r.haveExif = false
	r.haveXMP = false
	r.haveTHMB = false
	r.havePRVW = false
	r.stopAfterMetadata = false
	r.goalsInitialized = true
}

func (r *Reader) metadataGoalsSatisfied() bool {
	if r.goalExif && !r.haveExif {
		return false
	}
	if r.goalXMP && !r.haveXMP {
		return false
	}
	if r.goalTHMB && !r.haveTHMB {
		return false
	}
	if r.goalPRVW && !r.havePRVW {
		return false
	}
	return true
}

func (r *Reader) hasMetadataGoals() bool {
	return r.goalExif || r.goalXMP || r.goalTHMB || r.goalPRVW
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

	size, boxType, err := parseBoxSizeAndType(buf)
	if err != nil {
		return b, err
	}
	headerSize := 8
	if size == 1 {
		buf, err = r.peek(16)
		if err != nil {
			return b, fmt.Errorf("readBox: %w", ErrBufLength)
		}
		size, err = parseExtendedBoxSize(buf, boxType)
		if err != nil {
			return b, err
		}
		headerSize = 16
	}
	if err = validateBoxSize(size, headerSize, boxType); err != nil {
		return b, err
	}

	b.reader = r
	b.size = size
	b.boxType = boxType
	b.offset = r.offset

	b.remain = int(b.size)
	_, err = b.Discard(headerSize)
	return b, err
}
