package bmff

import (
	"fmt"
	"io"
)

// ItemType is
// always 4 bytes
type ItemType uint8

// ItemTypes
const (
	ItemTypeUnknown ItemType = iota
	ItemTypeInfe
	ItemTypeMime
	ItemTypeURI
	ItemTypeHvc1
	ItemTypeGrid
	ItemTypeExif
)

func (it ItemType) String() string {
	return mapItemTypeString[it]
}

func itemType(buf []byte) ItemType {
	if buf[0] == 'h' {
		if buf[1] == 'v' && buf[2] == 'c' && buf[3] == '1' {
			return ItemTypeHvc1
		}
	}
	t, found := mapStringItemType[string(buf)]
	if found {
		return t
	}
	if Debug {
		fmt.Printf("Unknown Item Type: %s\n", buf)
	}
	return ItemTypeUnknown
}

var mapItemTypeString = map[ItemType]string{
	ItemTypeInfe: "infe",
	ItemTypeMime: "mime",
	ItemTypeURI:  "uri ",
	ItemTypeHvc1: "hvc1",
	ItemTypeGrid: "grid",
	ItemTypeExif: "Exif",
}

var mapStringItemType = map[string]ItemType{
	"infe": ItemTypeInfe,
	"mime": ItemTypeMime,
	"uri ": ItemTypeURI,
	"hvc1": ItemTypeHvc1,
	"grid": ItemTypeGrid,
	"Exif": ItemTypeExif,
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

func (iinf *ItemInfoBox) setBox(b Box) error {
	if infe, ok := b.(ItemInfoEntry); ok {
		switch infe.ItemType {
		case ItemTypeExif:
			iinf.Exif = infe
		case ItemTypeMime:
			iinf.XMP = infe
		default:
			iinf.ItemInfos = append(iinf.ItemInfos, infe)
		}
	}
	return nil
}

func (iinf ItemInfoBox) Size() int64 {
	return int64(iinf.size)
}

func (iinf ItemInfoBox) Type() BoxType {
	return TypeIinf
}

func parseItemInfoBox(outer *box, br bufReader) (b Box, err error) {
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
		inner, err = boxr.ReadBox()
		if err != nil {
			if err == io.EOF {
				return ib, nil
			}
			boxr.br.err = err
			return ib, err
		}

		p, err := inner.Parse()
		if err != nil {
			boxr.br.discard(int(inner.r.remain))
			return ib, fmt.Errorf("error parsing ItemInfoEntry in ItemInfoBox: %v", err)
		}
		if err = ib.setBox(p); err != nil {
			boxr.br.discard(int(inner.r.remain))
		}

		outer.r.remain -= inner.size
		//fmt.Println(p.(ItemInfoEntry), outer.r.remain, inner.r.remain, boxr.br.remain)
	}
	//fmt.Println(int(ib.r.remain))
	//boxr.br.discard(int(fb.r.remain))
	if !br.ok() {
		return FullBox{}, br.err
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
	ContentType     string
	ContentEncoding string

	// If Type == "uri ":
	ItemURIType string

	size     int16
	ItemType ItemType
}

func (infe ItemInfoEntry) Size() int64 {
	return int64(infe.size)
}

func (infe ItemInfoEntry) Type() BoxType {
	return TypeInfe
}

func (infe ItemInfoEntry) String() string {
	return fmt.Sprintf("ItemInfoEntry (\"infe\"), Version: %d, Flags: %d, ItemID: %d, ProtectionIndex: %d, ItemType: %s, size: %d", infe.Flags.Version(), infe.Flags.Flags(), infe.ItemID, infe.ProtectionIndex, infe.ItemType, infe.size)
}

func parseItemInfoEntry(outer *box, br bufReader) (Box, error) {
	// Read Flags
	flags, err := outer.r.readFlags()
	if err != nil {
		return nil, err
	}
	if flags.Version() != 2 {
		return nil, fmt.Errorf("TODO: found version %d infe box. Only 2 is supported now.", flags.Version())
	}

	// New ItemInfoEntry
	ie := ItemInfoEntry{
		Flags: flags,
		size:  int16(outer.size),
	}

	ie.ItemID, _ = outer.r.readUint16()
	ie.ProtectionIndex, _ = outer.r.readUint16()
	if !br.ok() {
		return nil, br.err
	}
	ie.ItemType, err = outer.r.readItemType()
	if err != nil {
		return ie, outer.r.discard(int(outer.r.remain))
	}

	switch ie.ItemType {
	case ItemTypeMime:
		ie.ContentType, _ = outer.r.readString()
		if outer.r.anyRemain() {
			ie.ContentEncoding, _ = outer.r.readString()
		}
	case ItemTypeURI:
		ie.ItemURIType, _ = outer.r.readString()
	}
	if !outer.r.ok() {
		return nil, outer.r.err
	}
	return ie, nil
}
