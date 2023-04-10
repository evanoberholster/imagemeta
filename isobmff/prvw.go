package isobmff

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

type PRVWBox struct {
	Size   uint32
	Width  uint16
	Height uint16
}

func (r *Reader) readPreview(b *box) (err error) {
	inner, err := r.createPRVWBox(b)
	if err != nil {
		return errors.Wrapf(err, "ReadPRVWBox")
	}

	r.prvw, err = parsePreviewBox(&inner)
	if err != nil {
		return errors.Wrapf(err, "parsePreviewBox")
	}

	if r.PreviewImageReader != nil {
		if err = r.PreviewImageReader(&inner, meta.PreviewHeader(r.prvw)); err != nil {
			if logLevelError() {
				logError().Object("box", inner).Err(err).Send()
			}
		}
	}

	return inner.close()
}

func (r *Reader) createPRVWBox(b *box) (inner box, err error) {
	_, err = b.Discard(8)
	if err != nil {
		return inner, errors.Wrap(ErrBufLength, "readPRVWBoxDiscard")
	}

	buf, err := b.Peek(8)
	if err != nil {
		return inner, errors.Wrap(ErrBufLength, "readPRVWBoxPeek")
	}

	inner.reader = b.reader
	inner.outer = b
	inner.offset = int(b.size) - b.remain + b.offset
	inner.size = int64(bmffEndian.Uint32(buf[:4]))
	inner.remain = int(inner.size)
	inner.boxType = boxTypeFromBuf(buf[4:8])

	return inner, nil
}

func parsePreviewBox(b *box) (prvw PRVWBox, err error) {
	if !b.isType(typePRVW) {
		return prvw, ErrWrongBoxType
	}

	buf, err := b.Peek(24)
	if err != nil {
		return prvw, errors.Wrap(ErrBufLength, "parsePreviewBoxPeek")
	}

	prvw.Width = bmffEndian.Uint16(buf[14:16])
	prvw.Height = bmffEndian.Uint16(buf[16:18])
	prvw.Size = bmffEndian.Uint32(buf[20:24])

	_, err = b.Discard(24)
	if err != nil {
		return prvw, errors.Wrap(ErrBufLength, "parsePreviewBoxDiscard")
	}

	return prvw, nil
}
