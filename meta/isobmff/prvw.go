package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

type prvwBox struct {
	Size   uint32
	Width  uint16
	Height uint16
}

// readPreview parses Canon PRVW preview metadata and streams JPEG bytes via callback.
func (r *Reader) readPreview(b *box) (err error) {
	inner, err := r.createPRVWBox(b)
	if err != nil {
		return fmt.Errorf("readPreview: %w", err)
	}

	prvw, err := parsePreviewBox(&inner)
	if err != nil {
		return fmt.Errorf("parsePreviewBox: %w", err)
	}

	if r.previewImageReader != nil {
		header := meta.PreviewHeader{
			Size:      prvw.Size,
			Width:     prvw.Width,
			Height:    prvw.Height,
			ImageType: imagetype.ImageJPEG,
			Source:    meta.PreviewSourcePRVW,
		}
		if err = r.emitPreviewPayload(&inner, prvw.Size, header, metadataKindPRVW); err != nil {
			return err
		}
	}

	return inner.close()
}

// createPRVWBox creates an inner view over the PRVW box.
// Some CR3 UUID payloads include an 8-byte prefix before the PRVW header.
func (r *Reader) createPRVWBox(b *box) (inner box, err error) {
	inner, err = buildPRVWInnerBox(b)
	if err == nil {
		return inner, nil
	}
	if err != ErrWrongBoxType {
		return inner, err
	}

	if b.remain < 16 {
		return inner, fmt.Errorf("readPRVWBoxDiscard: %w", ErrBufLength)
	}
	if _, err = b.Discard(8); err != nil {
		return inner, fmt.Errorf("readPRVWBoxDiscard: %w", ErrBufLength)
	}
	return buildPRVWInnerBox(b)
}

// buildPRVWInnerBox validates the next box as PRVW and returns a bounded view.
func buildPRVWInnerBox(b *box) (inner box, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return inner, fmt.Errorf("readPRVWBoxPeek: %w", ErrBufLength)
	}

	inner.reader = b.reader
	inner.outer = b
	inner.offset = b.offset + b.size - b.remain
	inner.size = int64(bmffEndian.Uint32(buf[:4]))
	inner.remain = inner.size
	inner.boxType = boxTypeFromBuf(buf[4:8])
	if inner.boxType != typePRVW {
		return inner, ErrWrongBoxType
	}
	if inner.size < 8 || inner.remain > b.remain {
		return inner, ErrBufLength
	}

	return inner, nil
}

// parsePreviewBox reads PRVW dimensions and encoded preview size fields.
func parsePreviewBox(b *box) (prvw prvwBox, err error) {
	if !b.isType(typePRVW) {
		return prvw, ErrWrongBoxType
	}

	buf, err := b.Peek(24)
	if err != nil {
		return prvw, fmt.Errorf("parsePreviewBoxPeek: %w", ErrBufLength)
	}

	prvw.Width = bmffEndian.Uint16(buf[14:16])
	prvw.Height = bmffEndian.Uint16(buf[16:18])
	prvw.Size = bmffEndian.Uint32(buf[20:24])

	_, err = b.Discard(24)
	if err != nil {
		return prvw, fmt.Errorf("parsePreviewBoxDiscard: %w", ErrBufLength)
	}
	if uint64(prvw.Size) > uint64(b.remain) {
		return prvw, ErrRemainLengthInsufficient
	}

	return prvw, nil
}

func (r *Reader) emitPreviewPayload(b *box, size uint32, header meta.PreviewHeader, kind metadataKind) error {
	if size == 0 {
		return nil
	}
	if uint64(size) > uint64(b.remain) {
		return ErrRemainLengthInsufficient
	}

	payloadSize := int64(size)
	payload := box{
		reader:  b.reader,
		outer:   b,
		boxType: b.boxType,
		offset:  b.offset + b.size - b.remain,
		size:    payloadSize,
		remain:  payloadSize,
	}
	if err := r.callPreviewReader(&payload, header, kind, payload.remain); err != nil {
		return err
	}
	return payload.close()
}
