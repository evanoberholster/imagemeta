package bmff

import (
	"fmt"
	"io"
)

// ItemType
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
	//FullBox
	size    int64
	Version uint8
	Count   uint16
	Flags   uint32

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
	fb, err := readFullBox(outer)
	if err != nil {
		return nil, err
	}
	ib := ItemInfoBox{
		size:    fb.size,
		Version: fb.Version,
		Flags:   fb.Flags}
	boxr := fb.newReader(fb.r.remain)

	ib.Count, err = boxr.br.readUint16()
	if err != nil {
		return
	}
	ib.ItemInfos = make([]ItemInfoEntry, 0, int(ib.Count))

	var inner box
	for fb.r.remain > 4 {
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
		ib.setBox(p)

		fb.r.remain -= inner.size
		//fmt.Println(p.(ItemInfoEntry), ib.r.remain, inner.r.remain, boxr.br.remain)
	}
	//fmt.Println(int(ib.r.remain))
	//ib.r.discard(int(ib.r.remain))
	if !br.ok() {
		return FullBox{}, br.err
	}
	return ib, nil
}

// ItemInfoEntry represents an "infe" box.
//
// TODO: currently only parses Version 2 boxes.
type ItemInfoEntry struct {
	Flags           uint32 // 24 bits
	ItemID          uint16
	ProtectionIndex uint16

	//Name string

	// If Type == "mime":
	ContentType     string
	ContentEncoding string

	// If Type == "uri ":
	ItemURIType string

	size     int16
	Version  uint8
	ItemType ItemType // always 4 bytes
}

func (infe ItemInfoEntry) Size() int64 {
	return int64(infe.size)
}

func (infe ItemInfoEntry) Type() BoxType {
	return TypeInfe
}

func (infe ItemInfoEntry) String() string {
	return fmt.Sprintf("ItemInfoEntry (\"infe\"), Version: %d, Flags: %d, ItemID: %d, ProtectionIndex: %d, ItemType: %s, size: %d", infe.Version, infe.Flags, infe.ItemID, infe.ProtectionIndex, infe.ItemType, infe.size)
}

func newItemInfoEntry(fb FullBox) ItemInfoEntry {
	return ItemInfoEntry{Version: fb.Version, Flags: fb.Flags}
}

func parseItemInfoEntry(outer *box, br bufReader) (Box, error) {
	fb, err := readFullBox(outer)
	if err != nil {
		return nil, err
	}
	if fb.Version != 2 {
		return nil, fmt.Errorf("TODO: found version %d infe box. Only 2 is supported now.", fb.Version)
	}
	ie := ItemInfoEntry{
		Version: fb.Version,
		Flags:   fb.Flags,
		size:    int16(fb.size),
	}

	ie.ItemID, _ = fb.r.readUint16()
	ie.ProtectionIndex, _ = fb.r.readUint16()
	if !br.ok() {
		return nil, br.err
	}
	ie.ItemType, err = fb.r.readItemType()
	if err != nil {
		return ie, fb.r.discard(int(fb.r.remain))
	}

	switch ie.ItemType {
	case ItemTypeMime:
		ie.ContentType, _ = fb.r.readString()
		if fb.r.anyRemain() {
			ie.ContentEncoding, _ = fb.r.readString()
		}
	case ItemTypeURI:
		ie.ItemURIType, _ = fb.r.readString()
	}
	//fb.r.discard(int(fb.r.remain))
	if !fb.r.ok() {
		return nil, fb.r.err
	}
	return ie, nil
}
