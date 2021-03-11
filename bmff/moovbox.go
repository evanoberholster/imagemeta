package bmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta"
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

func (b *box) parseMoovBox() (moov MoovBox, err error) {
	moov.size = uint32(b.size)

	var inner box
	for b.anyRemain() {
		inner, err = b.readInnerBox()
		if err != nil {
			return
		}
		if inner.boxType == TypeUUID {
			uBox, err := inner.parseUUIDBox()
			fmt.Println(uBox, err)
		}

		if err = b.closeInnerBox(&inner); err != nil {
			break
		}
	}
	err = b.discard(b.remain)
	return
}

// UUIDBox is a special type of Box that contains a uuid
type UUIDBox struct {
	uuid meta.UUID
}

func (uuidBox UUIDBox) String() string {
	return fmt.Sprintf("uuid | %s\t", uuidBox.uuid.String())
}

// Type returns TypeUUID
func (uuidBox UUIDBox) Type() BoxType {
	return TypeUUID
}

func (b *box) parseUUIDBox() (Box, error) {
	if b.boxType != TypeUUID {
		return nil, ErrWrongBoxType
	}
	uuid, err := b.readUUID()
	if err != nil {
		return nil, err
	}
	switch uuid {
	case CR3MetaBoxUUID:
		return parseCR3MetaBox(b)
	}
	return nil, nil
}
