package isobmff

import (
	"errors"
	"io"
)

func newLimitedReader(r io.Reader, limit int) io.Reader {
	if limit <= 0 {
		return io.LimitReader(r, 0)
	}
	return &io.LimitedReader{R: r, N: int64(limit)}
}

func handleCallbackError(b *box, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	if logLevelError() {
		if b == nil {
			logError().Err(err).Send()
		} else {
			logError().Object("box", b).Err(err).Send()
		}
	}
	return nil
}
