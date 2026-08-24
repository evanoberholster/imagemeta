package jpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/meta"
)

const (
	bufferSize           int = 4 * 1024        // 4Kb
	maxMetadataScanBytes     = 2 * 1024 * 1024 // 2 MiB metadata scan budget
)

type jpegReader struct {
	ctx context.Context

	ExifReader func(r io.Reader, h meta.ExifHeader) error
	XMPReader  func(r io.Reader) error

	// Reader
	br       *bufio.Reader
	readerAt io.ReaderAt
	err      error

	// SOF Header
	sofHeader

	// Marker
	buf    []byte
	offset uint32
	size   uint16
	marker markerType

	// Reader
	pos       uint8
	discarded uint32

	extendedXMP map[string]*extendedXMP
	metadata    *Metadata
	foundExif   bool
}

var bufferPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, bufferSize) },
}

func scanJPEG(ctx context.Context, r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return scanJPEGWithMetadata(ctx, r, readerAt, exifReader, xmpReader, nil)
}

func scanMetadata(r io.Reader, readerAt io.ReaderAt) (m Metadata, err error) {
	err = scanJPEGWithMetadata(context.Background(), r, readerAt, nil, nil, &m)
	if finishErr := m.finish(); err == nil {
		err = finishErr
	}
	return m, err
}

func scanJPEGWithMetadata(ctx context.Context, r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error, metadata *Metadata) (err error) {
	defer func() {
		if state := recover(); state != nil {
			if recoveredErr, ok := state.(error); ok {
				err = recoveredErr
				return
			}
			err = fmt.Errorf("jpeg panic: %v", state)
		}
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	var localBuffer bool
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < bufferSize {
		localBuffer = true
		pooled, pooledOK := bufferPool.Get().(*bufio.Reader)
		if !pooledOK || pooled == nil {
			return fmt.Errorf("bufferPool returned non-*bufio.Reader")
		}
		br = pooled
		br.Reset(r)
	}

	jr := &jpegReader{ctx: ctx, br: br, readerAt: readerAt, ExifReader: exifReader, XMPReader: xmpReader, metadata: metadata}

	defer func() {
		if localBuffer {
			jr.br.Reset(nil)
			bufferPool.Put(jr.br)
		}
	}()

	for {
		if !jr.abortIfContextDone() {
			break
		}
		if !jr.nextMarker() {
			break
		}
		switch {
		case isSOFMarker(jr.marker):
			jr.readSOFMarker()
		case isAPPMarker(jr.marker):
			jr.readAPPMarker()
		default:
			switch jr.marker {
			case markerSOS:
				if err = jr.processExtendedXMP(); err != nil {
					return err
				}
				jr.logMarker("")
				return nil
			case markerDHT:
				jr.logMarker("")
				// Ignore DHT Markers
				jr.ignoreMarker()
			case markerSOI:
				jr.logMarker("")
				jr.pos++
				jr.err = jr.discard(2)
			case markerEOI:
				jr.logMarker("")
				if err = jr.processExtendedXMP(); err != nil {
					return err
				}
				jr.pos--
				if jr.err = jr.discard(2); jr.err != nil {
					return jr.err
				}
				return nil
			case markerDQT:
				jr.logMarker("")
				jr.ignoreMarker()
			case markerDRI:
				jr.err = jr.discard(6)
			default: // unknown marker
				jr.logMarker("")
				jr.ignoreMarker()
			}
		}
		// When only EXIF is requested, return immediately after decoding it.
		// This avoids walking potentially malformed trailing marker streams.
		if jr.err == nil && jr.foundExif && jr.ExifReader != nil && jr.XMPReader == nil && jr.metadata == nil {
			return nil
		}
	}
	if jr.err != nil {
		return jr.err
	}
	return jr.processExtendedXMP()
}

func (jr *jpegReader) abortIfContextDone() bool {
	if jr.ctx.Err() != nil {
		jr.err = jr.ctx.Err()
		return false
	}
	return true
}

func (jr *jpegReader) abortIfScanLimitExceeded() bool {
	if jr.discarded > maxMetadataScanBytes {
		jr.err = ErrMetadataScanLimit
		return false
	}
	return true
}

func (jr *jpegReader) nextMarker() bool {
	for jr.err == nil {
		if !jr.abortIfContextDone() {
			return false
		}
		if !jr.abortIfScanLimitExceeded() {
			return false
		}
		if jr.buf, jr.err = jr.peek(2); jr.err != nil {
			jr.err = ErrNoJPEGMarker
			return false
		}
		if !isMarkerFirstByte(jr.buf) {
			scanLen := 64
			if jr.buf, jr.err = jr.peek(scanLen); jr.err != nil && len(jr.buf) == 0 {
				jr.err = ErrNoJPEGMarker
				return false
			}
			var i int
			for i = 0; i < len(jr.buf); i++ {
				if isMarkerFirstByte(jr.buf[i:]) {
					break
				}
			}
			if i == len(jr.buf) {
				if i == 0 {
					jr.err = ErrNoJPEGMarker
					return false
				}
				// Keep the final byte in case it is the 0xff marker prefix.
				i--
			}
			jr.err = jr.discard(i)
			if jr.err != nil {
				return false
			}
			if !jr.abortIfScanLimitExceeded() {
				return false
			}
			continue
		}

		if isSOIMarker(jr.buf) {
			jr.pos++
			jr.err = jr.discard(2)
			if jr.err != nil {
				return false
			}
			if !jr.abortIfScanLimitExceeded() {
				return false
			}
			continue
		}
		if jr.pos > 0 {
			jr.offset = jr.discarded
			jr.marker = markerType(jr.buf[1])
			if markerHasNoLength(jr.marker) {
				jr.size = 0
				jr.buf = jr.buf[:2]
				return true
			}
			if jr.buf, jr.err = jr.peek(4); jr.err != nil {
				jr.err = ErrNoJPEGMarker
				return false
			}
			jr.size = jpegEndian.Uint16(jr.buf[2:4])
			if jr.size < 2 {
				// Invalid marker payload length. Attempt a byte-wise resync instead
				// of trusting this segment length and desynchronizing further.
				jr.err = jr.discard(1)
				if jr.err != nil {
					return false
				}
				if !jr.abortIfScanLimitExceeded() {
					return false
				}
				continue
			}
			peekLen := int(jr.size) + 2
			if peekLen > 64 {
				peekLen = 64
			}
			if peekLen < 4 {
				peekLen = 4
			}
			if jr.buf, jr.err = jr.peek(peekLen); jr.err != nil {
				jr.err = ErrNoJPEGMarker
				return false
			}
			return true
		}

		// Leading 0xFF without SOI: advance one byte to avoid infinite Peek loop.
		jr.err = jr.discard(1)
		if jr.err != nil {
			return false
		}
		if !jr.abortIfScanLimitExceeded() {
			return false
		}
	}
	return false
}

// peek returns the next n bytes without advancing the underlying bufio.Reader.
func (jr *jpegReader) peek(n int) ([]byte, error) {
	return jr.br.Peek(n)
}

// discard adds to m.discarded and discards from the underlying bufio.Reader
func (jr *jpegReader) discard(i int) (err error) {
	if i == 0 {
		return
	}
	i, err = jr.br.Discard(i)
	if delta, ok := meta.SafecastIntToUint32(i); ok {
		jr.discarded += delta
	}
	return
}

// readSOFMarker reads a JPEG Start of file with the uint16
// width, height, and components of the JPEG image.
func (jr *jpegReader) readSOFMarker() {
	precision := jr.buf[4]
	height := jpegEndian.Uint16(jr.buf[5:7])
	width := jpegEndian.Uint16(jr.buf[7:9])
	comp := jr.buf[9]
	if jr.pos == 1 {
		jr.sofHeader = sofHeader{height, width, comp}
	}
	if jr.metadata != nil {
		jr.metadata.SOF = SOF{
			Marker:          jr.marker.String(),
			EncodingProcess: uint8(jr.marker - markerSOF0),
			BitsPerSample:   precision,
			Width:           width,
			Height:          height,
			ColorComponents: comp,
		}
	}
	jr.err = jr.discard(int(jr.size) + 2)
}

// sofHeader contains height, width and number of components.
type sofHeader struct {
	height     uint16
	width      uint16
	components uint8
}

// ignoreMarker discards the marker size
func (jr *jpegReader) ignoreMarker() {
	jr.err = jr.discard(int(jr.size) + 2)
}

func (jr *jpegReader) readSegmentPayload() ([]byte, error) {
	payloadLen := int(jr.size) - 2
	if payloadLen < 0 {
		return nil, io.ErrUnexpectedEOF
	}
	payload := make([]byte, payloadLen)
	if jr.readerAt != nil {
		n, err := jr.readerAt.ReadAt(payload, int64(jr.offset)+4)
		if err != nil && n != payloadLen {
			return nil, err
		}
		if err := jr.discard(int(jr.size) + 2); err != nil {
			return nil, err
		}
		return payload, nil
	}
	if err := jr.discard(4); err != nil {
		return nil, err
	}
	n, err := io.ReadFull(jr.br, payload)
	if delta, ok := meta.SafecastIntToUint32(n); ok {
		jr.discarded += delta
	}
	if err != nil {
		return nil, err
	}
	return payload, nil
}
