package bmff

import "fmt"

// ItemInfoBox represents an "iinf" box.
type ItemInfoBox struct {
	FullBox
	Count     uint16
	ItemInfos []*ItemInfoEntry
}

// ItemInfoEntry represents an "infe" box.
//
// TODO: currently only parses Version 2 boxes.
type ItemInfoEntry struct {
	FullBox

	ItemID          uint16
	ProtectionIndex uint16
	ItemType        string // always 4 bytes

	Name string

	// If Type == "mime":
	ContentType     string
	ContentEncoding string

	// If Type == "uri ":
	ItemURIType string
}

func parseItemInfoEntry(outer *box, br *bufReader) (Box, error) {
	fb, err := readFullBox(outer)
	if err != nil {
		return nil, err
	}
	ie := &ItemInfoEntry{FullBox: fb}
	if fb.Version != 2 {
		return nil, fmt.Errorf("TODO: found version %d infe box. Only 2 is supported now.", fb.Version)
	}

	ie.ItemID, _ = br.readUint16()
	ie.ProtectionIndex, _ = br.readUint16()
	if !br.ok() {
		return nil, br.err
	}
	buf, err := br.Peek(4)
	if err != nil {
		return nil, err
	}
	ie.ItemType = string(buf[:4])
	ie.Name, _ = br.readString()

	switch ie.ItemType {
	case "mime":
		ie.ContentType, _ = br.readString()
		if br.anyRemain() {
			ie.ContentEncoding, _ = br.readString()
		}
	case "uri ":
		ie.ItemURIType, _ = br.readString()
	}
	if !br.ok() {
		return nil, br.err
	}
	return ie, nil
}
