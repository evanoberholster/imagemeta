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
	Flags       Flags
	size        uint32
	HandlerType HandlerType
	Name        string
}

// Size returns the size of the HandlerBox
func (hdlr HandlerBox) Size() int64 {
	return int64(hdlr.size)
}

// Type returns TypeHdlr
func (hdlr HandlerBox) Type() BoxType {
	return TypeHdlr
}

func parseHandlerBox(outer *box) (Box, error) {
	flags, err := outer.r.readFlags()
	//fb, err := readFullBox(outer)
	if err != nil {
		return nil, err
	}
	hb := HandlerBox{
		Flags: flags,
		size:  uint32(outer.size),
		//FullBox: fb,
	}
	buf, err := outer.r.Peek(20)
	if err != nil {
		return nil, err
	}
	hb.HandlerType = handler(buf[4:8])
	outer.r.discard(20)

	hb.Name, _ = outer.r.readString()
	outer.r.discard(int(outer.r.remain))
	return hb, outer.r.err
}
