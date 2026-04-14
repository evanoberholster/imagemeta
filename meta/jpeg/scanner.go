package jpeg

import (
	"bufio"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/meta"
)

const (
	bufferSize int = 4 * 1024 // 4Kb
)

type jpegReader struct {
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
}

var bufferPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, bufferSize) },
}

func scanJPEG(r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error) (err error) {
	return scanJPEGWithMetadata(r, readerAt, exifReader, xmpReader, nil)
}

func scanMetadata(r io.Reader, readerAt io.ReaderAt) (m Metadata, err error) {
	err = scanJPEGWithMetadata(r, readerAt, nil, nil, &m)
	if finishErr := m.finish(); err == nil {
		err = finishErr
	}
	return m, err
}

func scanJPEGWithMetadata(r io.Reader, readerAt io.ReaderAt, exifReader func(r io.Reader, header meta.ExifHeader) error, xmpReader func(r io.Reader) error, metadata *Metadata) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	var localBuffer bool
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < bufferSize {
		localBuffer = true
		br = bufferPool.Get().(*bufio.Reader)
		br.Reset(r)
	}

	jr := &jpegReader{br: br, readerAt: readerAt, ExifReader: exifReader, XMPReader: xmpReader, metadata: metadata}

	defer func() {
		if localBuffer {
			jr.br.Reset(nil)
			bufferPool.Put(jr.br)
		}
	}()

	for jr.nextMarker() {
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
				if logInfo() {
					jr.logMarker("")
				}
				return nil
			case markerDHT:
				if logInfo() {
					jr.logMarker("")
				}
				// Ignore DHT Markers
				jr.ignoreMarker()
			case markerSOI:
				if logInfo() {
					jr.logMarker("")
				}
				jr.pos++
				jr.err = jr.discard(2)
			case markerEOI:
				if logInfo() {
					jr.logMarker("")
				}
				if err = jr.processExtendedXMP(); err != nil {
					return err
				}
				jr.pos--
				if jr.err = jr.discard(2); jr.err != nil {
					return jr.err
				}
				return nil
			case markerDQT:
				if logInfo() {
					jr.logMarker("")
				}
				jr.ignoreMarker()
			case markerDRI:
				jr.err = jr.discard(6)
			default: // unknown marker
				if logInfo() {
					jr.logMarker("")
				}
				jr.ignoreMarker()
			}
		}
	}
	if jr.err != nil {
		return jr.err
	}
	return jr.processExtendedXMP()
}

func (jr *jpegReader) nextMarker() bool {
	for jr.err == nil {
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
			continue
		}

		if isSOIMarker(jr.buf) {
			jr.pos++
			jr.err = jr.discard(2)
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
	jr.discarded += uint32(i)
	return
}

// readSOFMarker reads a JPEG Start of file with the uint16
// width, height, and components of the JPEG image.
func (jr *jpegReader) readSOFMarker() {
	precision := uint8(jr.buf[4])
	height := jpegEndian.Uint16(jr.buf[5:7])
	width := jpegEndian.Uint16(jr.buf[7:9])
	comp := uint8(jr.buf[9])
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
	jr.discarded += uint32(n)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
