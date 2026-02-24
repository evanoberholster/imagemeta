package isobmff

import (
	"fmt"
)

// readHdlr reads an "hdlr" box
func readHdlr(b *box) (ht hdlrType, err error) {
	if !b.isType(typeHdlr) {
		return hdlrUnknown, fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if err = b.readFlags(); err != nil {
		return hdlrUnknown, err
	}

	if b.remain < 8 {
		return hdlrUnknown, fmt.Errorf("readHdlr: %w", ErrBufLength)
	}

	buf, err := b.Peek(8)
	if err != nil {
		return hdlrUnknown, err
	}
	ht = hdlrFromBuf(buf[4:8])
	if logLevelInfo() {
		logInfoBox(b).Str("hdlr", ht.String()).Send()
	}
	return ht, b.close()
}

// hdlrType

// hdlrType always 4 bytes;
// Handler; usually "pict" for HEIF images
type hdlrType uint8

// hdlr types
const (
	hdlrUnknown hdlrType = iota
	hdlrPict
	hdlrVide
	hdlrMeta
)

// String is a stringer interface for hdlrType
func (ht hdlrType) String() string {
	if str, ok := hdlrStringMap[ht]; ok {
		return str
	}
	return "nnnn"
}

var (
	hdlrPictFourCC = fourCCFromString("pict")
	hdlrVideFourCC = fourCCFromString("vide")
	hdlrMetaFourCC = fourCCFromString("meta")
)

func hdlrFromBuf(buf []byte) hdlrType {
	if len(buf) < 4 {
		return hdlrUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case hdlrPictFourCC:
		return hdlrPict
	case hdlrVideFourCC:
		return hdlrVide
	case hdlrMetaFourCC:
		return hdlrMeta
	default:
		return hdlrUnknown
	}
}

var hdlrStringMap = map[hdlrType]string{
	hdlrPict: "pict",
	hdlrVide: "vide",
	hdlrMeta: "meta",
}
