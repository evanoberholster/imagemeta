package jpeg

import (
	"bytes"
	"io"
)

type extendedXMP struct {
	size   uint32
	chunks map[uint32][]byte
}

func (jr *jpegReader) readExtendedXMP() error {
	payloadLen := int(jr.size) - 2
	if payloadLen < xmpExtHeaderLen {
		jr.ignoreMarker()
		return jr.err
	}
	if jr.XMPReader == nil {
		jr.ignoreMarker()
		return jr.err
	}

	if err := jr.discard(4); err != nil {
		return err
	}

	buf := make([]byte, payloadLen)
	n, err := io.ReadFull(jr.br, buf)
	jr.discarded += uint32(n)
	if err != nil {
		return err
	}
	if !bytes.HasPrefix(buf, []byte(xmpPrefixExt)) {
		return nil
	}

	guid := string(buf[len(xmpPrefixExt) : len(xmpPrefixExt)+32])
	fullSize := jpegEndian.Uint32(buf[67:71])
	chunkOffset := jpegEndian.Uint32(buf[71:75])
	chunk := buf[xmpExtHeaderLen:]

	if fullSize == 0 || fullSize > maxExtendedXMP {
		return nil
	}
	if uint64(chunkOffset)+uint64(len(chunk)) > uint64(fullSize) {
		return nil
	}

	if jr.extendedXMP == nil {
		jr.extendedXMP = make(map[string]*extendedXMP)
	}
	ext := jr.extendedXMP[guid]
	if ext == nil {
		ext = &extendedXMP{size: fullSize, chunks: make(map[uint32][]byte)}
		jr.extendedXMP[guid] = ext
	}
	if ext.size != fullSize {
		return nil
	}
	ext.chunks[chunkOffset] = append(ext.chunks[chunkOffset][:0], chunk...)
	return nil
}

func (jr *jpegReader) processExtendedXMP() error {
	if jr.XMPReader == nil || len(jr.extendedXMP) == 0 {
		return nil
	}

	for _, ext := range jr.extendedXMP {
		if ext == nil || ext.size == 0 || ext.size > maxExtendedXMP {
			continue
		}

		assembled := make([]byte, 0, ext.size)
		for offset := uint32(0); offset < ext.size; {
			chunk, ok := ext.chunks[offset]
			if !ok || len(chunk) == 0 {
				assembled = nil
				break
			}
			if uint64(offset)+uint64(len(chunk)) > uint64(ext.size) {
				assembled = nil
				break
			}
			assembled = append(assembled, chunk...)
			offset += uint32(len(chunk))
		}
		if uint32(len(assembled)) != ext.size {
			continue
		}
		if err := jr.XMPReader(bytes.NewReader(assembled)); err != nil {
			return err
		}
	}

	jr.extendedXMP = nil
	return nil
}
