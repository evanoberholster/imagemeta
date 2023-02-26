package isobmff

import (
	"github.com/pkg/errors"
)

// readHdlr reads an "hdlr" box
func readHdlr(b *box) (ht hdlrType, err error) {
	if !b.isType(typeHdlr) {
		return hdlrUnknown, errors.Wrapf(ErrWrongBoxType, "Box %s", b.boxType)
	}
	if err = b.readFlags(); err != nil {
		return hdlrUnknown, err
	}
	buf, err := b.Peek(b.remain)
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

func hdlrFromBuf(buf []byte) hdlrType {
	switch string(buf[:4]) {
	case "pict":
		return hdlrPict
	case "vide":
		return hdlrVide
	case "meta":
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
