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
	ItemTypeUnknown = itemType([]byte{0, 0, 0, 0})
)

func (it ItemType) String() string {
	return string(it[:])
}

func itemType(buf []byte) ItemType {
	if len(buf) != 4 {
		return ItemType([4]byte{0, 0, 0, 0})
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

	//Exif      ItemInfoEntry
	//XMP       ItemInfoEntry
	ItemInfos []ItemInfoEntry
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

func parseIinf(outer *box) (Box, error) {
	return parseItemInfoBox(outer)
}

func parseItemInfoBox(outer *box) (iinf ItemInfoBox, err error) {
	// Read Flags
	flags, err := outer.readFlags()
	if err != nil {
		return
	}
	// Read Item count
	count, err := outer.readUint16()
	if err != nil {
		return
	}

	// New ItemInfoBox
	iinf = ItemInfoBox{
		size:      outer.size,
		Flags:     flags,
		Count:     count,
		ItemInfos: make([]ItemInfoEntry, 0, int(count))}

	var inner box
	for outer.remain > 0 {
		inner, err = outer.readInnerBox()
		if err != nil {
			break
		}

		if inner.Type() == TypeInfe {
			var infe ItemInfoEntry
			if infe, err = parseItemInfoEntry(&inner); err != nil {
				inner.discard(inner.remain)
			}
			iinf.ItemInfos = append(iinf.ItemInfos, infe)
			if Debug {
				fmt.Println(infe, outer.remain, inner.remain, outer.size)
			}
		} else {
			// Error here Box Unknown
			err = inner.discard(inner.remain)
		}

		if err != nil {
			if Debug {
				err = fmt.Errorf("error parsing ItemInfoEntry in ItemInfoBox: %v", err)
				fmt.Println(err)
			}
		}

		outer.remain -= int(inner.size)
		inner.discard(inner.remain)

	}
	if Debug {
		fmt.Println(iinf, outer.remain)
	}
	err = outer.discard(outer.remain)
	return iinf, err
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

	size     int16
	ItemType ItemType
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
	flags, err := outer.readFlags()
	if err != nil {
		return ie, err
	}
	if flags.Version() != 2 {
		return ie, fmt.Errorf("TODO: found version %d infe box. Only 2 is supported now", flags.Version())
	}

	// New ItemInfoEntry
	ie = ItemInfoEntry{
		Flags: flags,
		size:  int16(outer.size),
	}

	ie.ItemID, _ = outer.readUint16()
	ie.ProtectionIndex, _ = outer.readUint16()
	if !outer.ok() {
		return ie, outer.err
	}
	ie.ItemType, err = outer.readItemType()
	if err != nil {
		return ie, outer.discard(outer.remain)
	}

	switch ie.ItemType {
	case ItemTypeMime:
		_, _ = outer.readString()
		if outer.anyRemain() {
			_, _ = outer.readString()
		}
		//ie.ContentType, _ = outer.r.readString()
		//if outer.r.anyRemain() {
		//	ie.ContentEncoding, _ = outer.r.readString()
		//}
	case ItemTypeURI:
		_, _ = outer.readString()
		//ie.ItemURIType, _ = outer.r.readString()
	}
	if !outer.ok() {
		return ie, outer.err
	}
	return ie, nil
}
