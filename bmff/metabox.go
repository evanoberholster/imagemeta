package bmff

import (
	"fmt"
	"strings"
)

// MetaBox is a 'meta' box
type MetaBox struct {
	//size  uint32
	Flags Flags

	Handler    HandlerBox
	Primary    PrimaryItemBox
	ItemInfo   ItemInfoBox
	Properties ItemPropertiesBox
	Location   ItemLocationBox
	Children   []Box
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
	mb.Flags, err = outer.readFlags()
	if err != nil {
		return mb, err
	}

	var inner box
	for outer.anyRemain() {
		inner, err = outer.readInnerBox()
		if err != nil {
			return mb, err
		}
		switch inner.boxType {
		case TypeIdat, TypeDinf, TypeUUID, TypeIref, TypeColr, TypeHvcC:
			// Do not parse

		//case TypeIref:
		//	_, err = inner.Parse()
		case TypePitm:
			mb.Primary, err = parsePrimaryItemBox(&inner)
		case TypeIinf:
			mb.ItemInfo, err = inner.parseItemInfoBox()
		case TypeHdlr:
			mb.Handler, err = parseHandlerBox(&inner)
		case TypeIprp:
			mb.Properties, err = inner.parseItemPropertiesBox()
		case TypeIloc:
			mb.Location, err = inner.parseItemLocationBox()
		default:
			//p, err := inner.Parse()
			//if err == nil {
			//	mb.Children = append(mb.Children, p)
			//}
		}
		if err != nil {
			return
		}
		if err = outer.closeInnerBox(&inner); err != nil {
			break
		}

		if debugFlag {
			fmt.Println(inner, outer.remain, inner.remain, inner.size)
		}
	}
	return mb, outer.discard(outer.remain)
}
