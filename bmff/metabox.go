package bmff

import (
	"fmt"
	"io"
	"strings"
)

// MetaBox is a 'meta' box
type MetaBox struct {
	size  uint32
	Flags Flags

	Handler    HandlerBox
	Primary    PrimaryItemBox
	ItemInfo   ItemInfoBox
	Properties ItemPropertiesBox
	Location   ItemLocationBox
	Children   []Box
}

// Size returns the size of the MetaBox
func (mb MetaBox) Size() int64 {
	return int64(mb.size)
}

// Type returns TypeMeta
func (mb MetaBox) Type() BoxType {
	return TypeMeta
}

func (mb MetaBox) String() string {
	var sb strings.Builder
	sb.WriteString("(Box) meta | Children: ")
	sb.WriteString(fmt.Sprint(len(mb.Children)))
	sb.WriteString("\n")
	sb.WriteString("\t")
	sb.WriteString(mb.Primary.String())
	sb.WriteString("\n")
	sb.WriteString("\t")
	sb.WriteString(mb.ItemInfo.String())
	sb.WriteString("\n")
	sb.WriteString("\t")
	sb.WriteString(mb.Properties.String())
	sb.WriteString("\n")
	sb.WriteString("\t")
	sb.WriteString(mb.Location.String())
	return sb.String()
}

func parseMeta(outer *box) (Box, error) {
	return parseMetaBox(outer)
}

func parseMetaBox(outer *box) (mb MetaBox, err error) {
	if outer.boxType != TypeMeta {
		err = ErrWrongBoxType
		return
	}
	mb = MetaBox{size: uint32(outer.size)}
	mb.Flags, err = outer.readFlags()
	if err != nil {
		return mb, err
	}

	var inner box
	for outer.anyRemain() {
		inner, err = outer.readBox()
		if err != nil {
			if err == io.EOF {
				return mb, nil
			}
			return mb, err
		}
		switch inner.boxType {
		case TypeIdat, TypeDinf, TypeUUID, TypeIref:
			// Do not parse

		//case TypeIref:
		//	_, err = inner.Parse()
		case TypePitm:
			mb.Primary, err = parsePrimaryItemBox(&inner)
		case TypeIinf:
			mb.ItemInfo, err = parseItemInfoBox(&inner)
		case TypeHdlr:
			mb.Handler, err = parseHandlerBox(&inner)
		case TypeIprp:
			mb.Properties, err = parseItemPropertiesBox(&inner)
		case TypeIloc:
			mb.Location, err = parseItemLocationBox(&inner)
		default:
			p, err := inner.Parse()
			if err == nil {
				mb.Children = append(mb.Children, p)
			}
		}
		if err != nil {
			if Debug {
				fmt.Println(err)
			}
		}
		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			if Debug {
				fmt.Println(err)
			}
			// TODO: improve error handling
			break
		}

		if Debug {
			fmt.Println(inner, outer.remain, inner.remain, inner.size)
		}
	}
	err = outer.discard(outer.remain)
	return mb, err
}
