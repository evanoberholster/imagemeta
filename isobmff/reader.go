package isobmff

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// Constants
const (
	minBufReaderSize = 4096 // 4Kb
)

var (
	// bmffEndian ISOBMFF always uses BigEndian byteorder.
	// Can use either byteorder for Exif Information inside the ISOBMFF file
	bmffEndian = binary.BigEndian

	// crxEndian values are in BigEndian.
	crxEndian = binary.BigEndian
)

// readerPool for buffer
var readerPool = sync.Pool{
	New: func() interface{} { return bufio.NewReaderSize(nil, minBufReaderSize) },
}

// Reader is a ISO BMFF reader
type Reader struct {
	//br         bufReader
	//brand      Brand
	br *bufio.Reader

	ftyp FileTypeBox

	ExifReader func(r io.Reader, h meta.ExifHeader) error
	XMPReader  func(r io.Reader) error

	offset int
	rPool  bool
}

// NewReader returns a new bmff.Reader
func NewReader(r io.Reader) Reader {
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < minBufReaderSize {
		br = readerPool.Get().(*bufio.Reader)
		br.Reset(r)
		return Reader{br: br, rPool: true}
	}
	return Reader{br: br}
}

func (r *Reader) peek(n int) ([]byte, error) {
	return r.br.Peek(n)
}

func (r *Reader) discard(n int) (err error) {
	n, err = r.br.Discard(n)
	r.offset += n
	return err
}

func (r *Reader) reset(newReader io.Reader) {
	r.Close()
	*r = NewReader(newReader)
}

// Close the Reader. Returns the underlying bufio.Reader to the reader pool.
func (r *Reader) Close() {
	if r.rPool {
		readerPool.Put(r.br)
	}
}

// readBox reads an ISOBMFF box
func (r *Reader) readBox() (b box, err error) {
	// Read box size and box type
	var buf []byte
	if buf, err = r.peek(8); err != nil {
		return b, errors.Wrap(ErrBufLength, "readBox")
	}
	b.reader = r
	b.size = int64(bmffEndian.Uint32(buf[:4]))
	b.remain = int(b.size)
	b.boxType = boxTypeFromBuf(buf[4:8])
	b.offset = r.offset

	switch b.size {
	case 1:
		// 1 means it's actually a 64-bit size, after the type.
		if buf, err = r.peek(16); err != nil {
			return b, errors.Wrap(ErrBufLength, "readBox")
		}
		b.size = int64(bmffEndian.Uint32(buf[8:16]))
		if b.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return b, errors.Wrapf(errLargeBox, "readBox '%s'", b.boxType)
		}
		b.remain = int(b.size)
		return b, b.Discard(16)
		//case 0:
		// 0 means unknown & to read to end of file. No more boxes.
		// r.noMoreBoxes = true
		// TODO: error
	}
	return b, b.Discard(8)
}
