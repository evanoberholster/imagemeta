package bmff

import (
	"fmt"
)

// ItemType is
// always 4 bytes
type ItemType [4]byte

// Common ItemTypes
var (
	ItemTypeInfe    = itemType([]byte("infe"))
	ItemTypeMime    = itemType([]byte("mime"))
	ItemTypeURI     = itemType([]byte("uri "))
	ItemTypeAv01    = itemType([]byte("av01"))
	ItemTypeHvc1    = itemType([]byte("hvc1"))
	ItemTypeGrid    = itemType([]byte("grid"))
	ItemTypeExif    = itemType([]byte("Exif"))
	ItemTypeUnknown = itemType([]byte("nnnn"))
)

func (it ItemType) String() string {
	return string(it[:])
}

func itemType(buf []byte) ItemType {
	if len(buf) != 4 {
		// Error
	}
	it := ItemType{}
	copy(it[:], buf[:4])
	return it
}

// ItemInfoBox represents an "iinf" box.
type ItemInfoBox struct {
	size  int64
	Flags Flags
	Count uint16

	Exif      ItemInfoEntry
	XMP       ItemInfoEntry
	ItemInfos []ItemInfoEntry
}

func (iinf *ItemInfoBox) setBox(infe ItemInfoEntry) error {
	switch infe.ItemType {
	case ItemTypeExif:
		iinf.Exif = infe
	case ItemTypeMime:
		iinf.XMP = infe
	default:
		iinf.ItemInfos = append(iinf.ItemInfos, infe)
	}
	return nil
}

func (iinf ItemInfoBox) String() string {
	return fmt.Sprintf("iinf | ItemCount: %d, Flags: %d, Version: %d", len(iinf.ItemInfos), iinf.Flags.Flags(), iinf.Flags.Version())
}

// Size returns the size of the ItemInfoBox
func (iinf ItemInfoBox) Size() int64 {
	return int64(iinf.size)
}

// Type returns TypeIinf
func (iinf ItemInfoBox) Type() BoxType {
	return TypeIinf
}

func parseItemInfoBox(outer *box) (b Box, err error) {
	// Read Flags
	flags, err := outer.r.readFlags()
	if err != nil {
		return nil, err
	}
	// Read Item count
	count, err := outer.r.readUint16()
	if err != nil {
		return
	}

	// New ItemInfoBox
	ib := ItemInfoBox{
		size:      outer.size,
		Flags:     flags,
		Count:     count,
		ItemInfos: make([]ItemInfoEntry, 0, int(count))}

	boxr := outer.newReader(outer.r.remain)

	var inner box
	for outer.r.remain > 4 {
		inner, err = boxr.readBox()
		if err != nil {
			boxr.br.err = err
			return ib, err
		}
		if inner.Type() == TypeInfe {
			ie, err := parseItemInfoEntry(&inner)
			if err != nil {
				boxr.br.discard(int(inner.r.remain))
			}
			ib.ItemInfos = append(ib.ItemInfos, ie)
			//if err = ib.setBox(ie); err != nil {
			//	boxr.br.discard(int(inner.r.remain))
			//}
		} else {
			// Error here
			boxr.br.discard(int(inner.r.remain))
		}
		if err != nil {
			if Debug {
				err = fmt.Errorf("error parsing ItemInfoEntry in ItemInfoBox: %v", err)
				fmt.Println(err)
			}
		}

		outer.r.remain -= int(inner.size)
		boxr.br.discard(int(inner.r.remain))
		if Debug {
			//fmt.Println(p.(ItemInfoEntry), outer.r.remain, inner.r.remain, boxr.br.remain)
		}
	}
	//fmt.Println(int(ib.r.remain))
	//boxr.br.discard(int(fb.r.remain))
	if !outer.r.ok() {
		return ib, outer.r.err
	}
	if Debug {
		fmt.Println(ib)
	}
	return ib, nil
}

// ItemInfoEntry represents an "infe" box.
//
// TODO: currently only parses Version 2 boxes.
type ItemInfoEntry struct {
	Flags           Flags
	ItemID          uint16
	ProtectionIndex uint16

	//Name string

	// If Type == "mime":
	//ContentType     string
	//ContentEncoding string

	// If Type == "uri ":
	//ItemURIType string

	ItemType ItemType
	//size     int16
}

// Type returns TypeInfe
func (infe ItemInfoEntry) Type() BoxType {
	return TypeInfe
}

func (infe ItemInfoEntry) String() string {
	return fmt.Sprintf(" \t ItemInfoEntry (\"infe\"), Version: %d, Flags: %d, ItemID: %d, ProtectionIndex: %d, ItemType: %s", infe.Flags.Version(), infe.Flags.Flags(), infe.ItemID, infe.ProtectionIndex, infe.ItemType)
}

func parseInfe(outer *box) (Box, error) {
	return parseItemInfoEntry(outer)
}

func parseItemInfoEntry(outer *box) (ie ItemInfoEntry, err error) {
	// Read Flags
	flags, err := outer.r.readFlags()
	if err != nil {
		return ie, err
	}
	if flags.Version() != 2 {
		return ie, fmt.Errorf("TODO: found version %d infe box. Only 2 is supported now", flags.Version())
	}

	// New ItemInfoEntry
	ie = ItemInfoEntry{
		Flags: flags,
		//size:  int16(outer.size),
	}

	ie.ItemID, _ = outer.r.readUint16()
	ie.ProtectionIndex, _ = outer.r.readUint16()
	if !outer.r.ok() {
		return ie, outer.r.err
	}
	ie.ItemType, err = outer.r.readItemType()
	if err != nil {
		return ie, outer.r.discard(int(outer.r.remain))
	}

	switch ie.ItemType {
	case ItemTypeMime:
		_, _ = outer.r.readString()
		if outer.r.anyRemain() {
			_, _ = outer.r.readString()
		}
		//ie.ContentType, _ = outer.r.readString()
		//if outer.r.anyRemain() {
		//	ie.ContentEncoding, _ = outer.r.readString()
		//}
	case ItemTypeURI:
		_, _ = outer.r.readString()
		//ie.ItemURIType, _ = outer.r.readString()
	}
	if !outer.r.ok() {
		return ie, outer.r.err
	}
	return ie, nil
}
