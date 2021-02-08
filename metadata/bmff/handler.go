package bmff

import (
	"bytes"
	"fmt"
)

// HandlerType always 4 bytes; usually "pict" for iOS Camera images
type HandlerType uint8

const (
	handlerUnknown HandlerType = iota
	handlerPict
)

func handler(buf []byte) HandlerType {
	if bytes.Equal(buf[:], []byte{'p', 'i', 'c', 't'}) {
		return handlerPict
	}
	if Debug {
		fmt.Println("Unknown Handler: Error", string(buf), buf)
	}
	return handlerUnknown
}

// HandlerBox is a "hdlr" box.
type HandlerBox struct {
	FullBox
	HandlerType HandlerType
	Name        string
}

func parseHandlerBox(gen *box, br bufReader) (Box, error) {
	fb, err := readFullBox(gen)
	if err != nil {
		return nil, err
	}
	hb := HandlerBox{
		FullBox: fb,
	}
	buf, err := fb.r.Peek(20)
	if err != nil {
		return nil, err
	}
	hb.HandlerType = handler(buf[4:8])
	fb.r.discard(20)

	hb.Name, _ = fb.r.readString()
	fb.r.discard(int(fb.r.remain))
	return hb, br.err
}
