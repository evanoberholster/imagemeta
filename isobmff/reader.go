package isobmff

import (
	"bufio"
	"encoding/binary"
	"errors"
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
)

type heicMeta struct {
	pitm itemID
	idat idat
	exif item
	xml  item

	idatData      offsetLength
	items         []itemInfo
	locations     []itemLocation
	references    []itemReference
	properties    []itemProperty
	propertyLinks []itemPropertyLink
	// irot
}

type item struct {
	id itemID
	ol offsetLength
}

type metadataKind uint8

const (
	metadataKindExif metadataKind = iota
	metadataKindXMP
	metadataKindTHMB
	metadataKindPRVW
	metadataKindCount
)

func goalBit(kind metadataKind) uint8 {
	return uint8(kind)
}

func haveBit(kind metadataKind) uint8 {
	return uint8(kind) + 4
}

func hasBit(flags uint8, bit uint8) bool {
	return flags&(1<<bit) != 0
}

func setBit(flags *uint8, bit uint8) {
	*flags |= 1 << bit
}

func clearBit(flags *uint8, bit uint8) {
	*flags &^= 1 << bit
}

var readerPool = sync.Pool{
	New: func() any { return bufio.NewReaderSize(nil, bufReaderSize) },
}

// Reader is a ISO BMFF reader
type Reader struct {
	source io.Reader
	seeker io.Seeker

	heic heicMeta
	ftyp fileTypeBox

	br                   *bufio.Reader
	exifReader           ExifReader
	xmpReader            XMPReader
	previewImageReader   PreviewImageReader
	pooledBufio          bool
	offset               int
	discardSeekThreshold int

	prvw prvwBox

	metadataFlags uint8

	stopAfterMetadata bool
	goalsInitialized  bool
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
		return newReaderWithBufio(br, r, false)
	}

	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	return newReaderWithBufio(br, r, true)
}

// newReaderWithBufio wires a Reader around an already configured bufio.Reader.
func newReaderWithBufio(br *bufio.Reader, source io.Reader, pooled bool) Reader {
	reader := Reader{
		br:                   br,
		source:               source,
		pooledBufio:          pooled,
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

// discard advances the stream and updates absolute offset.
// Large skips on seekable sources are delegated to discardWithSeek.
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

// discardWithSeek prefers Seek for large skips to avoid read-and-discard loops.
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

// reset reinitializes reader state for a new source while preserving callbacks.
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

// initMetadataGoals derives extraction goals from active callbacks and file type.
func (r *Reader) initMetadataGoals() {
	// Reset parsed metadata graph when starting a new file scan.
	r.heic = heicMeta{}
	r.prvw = prvwBox{}
	r.metadataFlags = 0

	r.setGoal(metadataKindExif, r.exifReader != nil)
	r.setGoal(metadataKindXMP, r.xmpReader != nil)
	r.setGoal(metadataKindTHMB, false)
	r.setGoal(metadataKindPRVW, false)
	if r.previewImageReader != nil {
		if r.ftyp.MajorBrand == brandCrx {
			r.setGoal(metadataKindTHMB, true)
			r.setGoal(metadataKindPRVW, true)
		} else {
			r.setGoal(metadataKindPRVW, true)
		}
	}

	r.stopAfterMetadata = false
	r.goalsInitialized = true
}

// metadataGoalsSatisfied reports whether all requested metadata callbacks have fired.
func (r *Reader) metadataGoalsSatisfied() bool {
	for kind := metadataKind(0); kind < metadataKindCount; kind++ {
		if r.hasGoal(kind) && !r.hasHave(kind) {
			return false
		}
	}
	return true
}

func (r *Reader) hasMetadataGoals() bool {
	for kind := metadataKind(0); kind < metadataKindCount; kind++ {
		if r.hasGoal(kind) {
			return true
		}
	}
	return false
}

func (r *Reader) hasGoal(kind metadataKind) bool {
	return hasBit(r.metadataFlags, goalBit(kind))
}

func (r *Reader) hasHave(kind metadataKind) bool {
	return hasBit(r.metadataFlags, haveBit(kind))
}

func (r *Reader) setGoal(kind metadataKind, enabled bool) {
	if enabled {
		setBit(&r.metadataFlags, goalBit(kind))
		return
	}
	clearBit(&r.metadataFlags, goalBit(kind))
}

func (r *Reader) setHave(kind metadataKind, enabled bool) {
	if enabled {
		setBit(&r.metadataFlags, haveBit(kind))
		return
	}
	clearBit(&r.metadataFlags, haveBit(kind))
}

// Close the Reader. Returns the underlying bufio.Reader to the reader pool.
func (r *Reader) Close() {
	if r.pooledBufio && r.br != nil {
		// Clear references to allow GC of the previous source quickly.
		r.br.Reset(nil)
		readerPool.Put(r.br)
	}
	r.br = nil
	r.source = nil
	r.seeker = nil
	r.pooledBufio = false
}

// readBox reads an ISOBMFF box
func (r *Reader) readBox() (b box, err error) {
	// Read box size and box type (8-byte header)
	buf, err := r.peek(8)
	if err != nil {
		if err == io.EOF {
			return b, io.EOF
		}
		return b, fmt.Errorf("readBox: failed to read header: %w", errors.Join(ErrBufLength, err))
	}

	size, boxType, err := parseBoxSizeAndType(buf)
	if err != nil {
		return b, err
	}
	headerSize := 8
	if size == 1 {
		buf, err = r.peek(16)
		if err != nil {
			return b, fmt.Errorf("readBox: failed to read extended header: %w", errors.Join(ErrBufLength, err))
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
