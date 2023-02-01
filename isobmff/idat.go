package isobmff

import (
	"github.com/rs/zerolog"
)

func readIdat(b *box) (i idat, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return
	}
	i = idat{
		width:  bmffEndian.Uint16(buf[4:6]),
		height: bmffEndian.Uint16(buf[6:8])}
	if logLevelInfo() {
		logInfoBox(b).Object("idat", i).Send()
	}

	return i, b.close()
}

// ItemData is an "idat" box

// idat
type idat struct {
	width, height uint16
}

// MarshalZerologObject is a zerolog interface for logging
func (i idat) MarshalZerologObject(e *zerolog.Event) {
	e.Uint16("width", i.width).Uint16("height", i.height)
}
