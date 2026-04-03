package utils

import (
	"bufio"
	"io"
	"sync"
)

// BufioReaderPool reuses bufio.Reader instances with a fixed buffer size.
type BufioReaderPool struct {
	pool        sync.Pool
	size        int
	resetReader io.Reader
}

// NewBufioReaderPool creates a new bufio.Reader pool.
//
// resetReader is used when returning readers to the pool to clear source refs.
func NewBufioReaderPool(size int, resetReader io.Reader) *BufioReaderPool {
	p := &BufioReaderPool{
		size:        size,
		resetReader: resetReader,
	}
	p.pool.New = func() any {
		return bufio.NewReaderSize(resetReader, size)
	}
	return p
}

// Acquire returns a pooled reader reset to src.
func (p *BufioReaderPool) Acquire(src io.Reader) *bufio.Reader {
	if p == nil {
		return bufio.NewReader(src)
	}
	br, ok := p.pool.Get().(*bufio.Reader)
	if !ok || br == nil {
		return bufio.NewReaderSize(src, p.size)
	}
	br.Reset(src)
	return br
}

// Release resets br to resetReader and returns it to the pool.
func (p *BufioReaderPool) Release(br *bufio.Reader) {
	if p == nil || br == nil {
		return
	}
	br.Reset(p.resetReader)
	p.pool.Put(br)
}
