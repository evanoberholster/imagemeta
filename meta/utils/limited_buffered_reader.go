package utils

import "io"

// LimitedBufferedReader bounds reads/discards/peeks to N bytes while preserving
// BufferedReader capabilities.
type LimitedBufferedReader struct {
	R BufferedReader
	N int
}

// NewLimitedBufferedReader returns a bounded buffered reader view.
func NewLimitedBufferedReader(r BufferedReader, limit int) *LimitedBufferedReader {
	if limit < 0 {
		limit = 0
	}
	return &LimitedBufferedReader{R: r, N: limit}
}

// Read reads up to len(p) bytes, bounded by remaining limit N.
func (l *LimitedBufferedReader) Read(p []byte) (int, error) {
	if l == nil || l.R == nil {
		return 0, io.EOF
	}
	if l.N <= 0 {
		return 0, io.EOF
	}
	if len(p) > l.N {
		p = p[:l.N]
	}
	n, err := l.R.Read(p)
	if n > 0 {
		l.N -= n
	}
	return n, err
}

// Peek returns the next n bytes without advancing, bounded by remaining N.
func (l *LimitedBufferedReader) Peek(n int) ([]byte, error) {
	if l == nil || l.R == nil {
		return nil, io.EOF
	}
	if n <= 0 {
		return l.R.Peek(0)
	}
	if l.N <= 0 {
		return nil, io.EOF
	}
	if n > l.N {
		n = l.N
	}
	return l.R.Peek(n)
}

// Discard skips n bytes (or fewer if limit is smaller), updating N.
func (l *LimitedBufferedReader) Discard(n int) (int, error) {
	if l == nil || l.R == nil {
		return 0, io.EOF
	}
	if n <= 0 {
		return 0, nil
	}
	if l.N <= 0 {
		return 0, io.EOF
	}
	need := n
	wantEOF := false
	if need > l.N {
		need = l.N
		wantEOF = true
	}
	discarded, err := l.R.Discard(need)
	if discarded > 0 {
		l.N -= discarded
	}
	if err != nil {
		return discarded, err
	}
	if wantEOF || discarded < need {
		return discarded, io.EOF
	}
	return discarded, nil
}
