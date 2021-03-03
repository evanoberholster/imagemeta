package bmff

import (
	"fmt"
	"io"
)

// MoovBox is a 'moov' box
type MoovBox struct {
	size uint32
}

// Size returns the size of the MoovBox
func (moov MoovBox) Size() int64 {
	return int64(moov.size)
}

// Type returns TypeMoov
func (moov MoovBox) Type() BoxType {
	return TypeMoov
}

func parseMoov(outer *box) (Box, error) {
	if outer.boxType != TypeMoov {
		return nil, ErrWrongBoxType
	}
	return parseMetaBox(outer)
}

func parseMoovBox(outer *box) (moov MoovBox, err error) {
	moov = MoovBox{size: uint32(outer.size)}

	var inner box
	for outer.anyRemain() {
		inner, err = outer.readInnerBox()
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
		if inner.boxType == TypeUUID {
			uBox, err := parseUUIDBox(&inner)
			fmt.Println(uBox, err)
		}

		if err = outer.closeInnerBox(&inner); err != nil {
			// Log error
			if debugFlag {
				fmt.Println(err)
			}
			break
		}
	}
	err = outer.discard(outer.remain)
	return
}

func parseUUIDBox(outer *box) (b Box, err error) {
	if outer.boxType != TypeUUID {
		err = ErrWrongBoxType
		return
	}
	uuid, err := outer.readUUID()
	if err != nil {
		return
	}
	switch uuid {
	case CR3MetaBoxUUID:
		return parseCR3MetaBox(outer)
	}
	return
}
