package bmff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrItemNotFound is returned as en error when an Item was not Found.
	ErrItemNotFound = errors.New("item not found")
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
	return ItemType{buf[0], buf[1], buf[2], buf[3]}
}

// ItemInfoBox represents an "iinf" box.
type ItemInfoBox struct {
	size  int64
	Flags Flags
	Count uint16

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

// LastItemByType returns the Last Item of a given Type.
// If the ItemType is not found, the returned error is ErrItemNotFound.
func (iinf ItemInfoBox) LastItemByType(itemType ItemType) (ItemInfoEntry, error) {
	for i := len(iinf.ItemInfos) - 1; i >= 0; i-- {
		if iinf.ItemInfos[i].ItemType == itemType {
			return iinf.ItemInfos[i], nil
		}
	}
	return ItemInfoEntry{}, ErrItemNotFound
}

// ItemByID returns the Item of a given ID.
// If the ItemID is not found, the returned error is ErrItemNotFound.
func (iinf ItemInfoBox) ItemByID(itemID uint16) (ItemInfoEntry, error) {
	for i := len(iinf.ItemInfos) - 1; i >= 0; i-- {
		if iinf.ItemInfos[i].ItemID == itemID {
			return iinf.ItemInfos[i], nil
		}
	}
	return ItemInfoEntry{}, ErrItemNotFound
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
				err = inner.discard(inner.remain)
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
		if err = inner.discard(inner.remain); err != nil {
			break
		}

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

// ItemLocationBox is a "iloc" box
type ItemLocationBox struct {
	Items []ItemLocationBoxEntry
	Flags
	ItemCount uint16

	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// Type returns TypeIloc
func (iloc ItemLocationBox) Type() BoxType {
	return TypeIloc
}

func (iloc ItemLocationBox) String() string {
	return fmt.Sprintf("iloc | ItemCount: %d, Flags: %d, Version: %d, OffsetSize: %d, LengthSize: %d, BaseOffsetSize: %d, indexSize: %d", iloc.ItemCount, iloc.Flags.Flags(), iloc.Flags.Version(), iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize)
}

// EntryByID returns the last Item Location Box Entry for the given ID. If
// the ID is not found returns ErrItemNotFound.
func (iloc ItemLocationBox) EntryByID(id uint16) (ItemLocationBoxEntry, error) {
	for i := len(iloc.Items) - 1; i >= 0; i-- {
		if iloc.Items[i].ItemID == id {
			return iloc.Items[i], nil
		}
	}
	return ItemLocationBoxEntry{}, ErrItemNotFound
}

func parseIloc(outer *box) (Box, error) {
	return parseItemLocationBox(outer)
}

func parseItemLocationBox(outer *box) (ilb ItemLocationBox, err error) {
	ilb.Flags, err = outer.readFlags()
	if err != nil {
		return
	}

	buf, err := outer.Peek(4)
	if err != nil {
		// TODO: Write error handling
		return
	}
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if ilb.Flags.Version() > 0 { // version 1
		ilb.indexSize = buf[1] & 15
	}

	ilb.ItemCount = binary.BigEndian.Uint16(buf[2:4])
	if err = outer.discard(4); err != nil {
		// TODO: Write error handling
		return
	}

	ilb.Items = make([]ItemLocationBoxEntry, 0, ilb.ItemCount)

	if Debug {
		fmt.Println(ilb)
	}
	for i := 0; outer.anyRemain() && i < int(ilb.ItemCount); i++ {
		var ent ItemLocationBoxEntry
		ent.ItemID, _ = outer.readUint16()

		if ilb.Flags.Version() > 0 { // version 1
			cmeth, _ := outer.readUint16()
			ent.ConstructionMethod = byte(cmeth & 15)
		}
		ent.DataReferenceIndex, _ = outer.readUint16()

		// Adjust for baseOffset per issue "https://github.com/go4org/go4/issues/47" thanks to petercgrant
		if outer.ok() && ilb.baseOffsetSize > 0 {
			ent.BaseOffset, _ = outer.readUintN(ilb.baseOffsetSize * 8)
			//outer.r.discard(int(ilb.baseOffsetSize) / 8)
		}

		// ExtentCount
		ent.ExtentCount, _ = outer.readUint16()
		for j := 0; outer.ok() && j < int(ent.ExtentCount); j++ {
			var ol OffsetLength
			ol.Offset, _ = outer.readUintN(ilb.offsetSize * 8)
			ol.Length, _ = outer.readUintN(ilb.lengthSize * 8)
			if outer.err != nil {
				err = outer.err
				return
			}
			if j == 0 {
				ent.FirstExtent = ol
				continue
			}
			ent.Extents = append(ent.Extents, ol)
		}

		ilb.Items = append(ilb.Items, ent)

		if Debug {
			fmt.Println(ent, outer.remain)
		}
	}
	if !outer.ok() {
		err = outer.err
		return
	}
	return ilb, nil
}

// ItemLocationBoxEntry is not a box
type ItemLocationBoxEntry struct {
	Extents            []OffsetLength
	FirstExtent        OffsetLength
	BaseOffset         uint64 // uint32 or uint64, depending on encoding
	ItemID             uint16
	ExtentCount        uint16
	DataReferenceIndex uint16
	ConstructionMethod uint8 // actually uint4
}

func (ilbe ItemLocationBoxEntry) String() string {
	return fmt.Sprintf("\t ItemID: %d, ConstructionMethod: %d, DataReferenceIndex: %d, BaseOffset: %d, ExtentCount: %d, FirstExtent: %s", ilbe.ItemID, ilbe.ConstructionMethod, ilbe.DataReferenceIndex, ilbe.BaseOffset, ilbe.ExtentCount, ilbe.FirstExtent)
}

// OffsetLength contains an offset and length
type OffsetLength struct {
	Offset, Length uint64
}

func (ol OffsetLength) String() string {
	return fmt.Sprintf("{ Offset: %d, Length: %d }", ol.Offset, ol.Length)
}

// ItemPropertiesBox is an ISOBMFF "iprp" box
type ItemPropertiesBox struct {
	PropertyContainer ItemPropertyContainerBox
	Associations      []ItemPropertyAssociation // at least 1
}

// Type returns TypeIprp
func (iprp ItemPropertiesBox) Type() BoxType {
	return TypeIprp
}

func (iprp ItemPropertiesBox) String() string {
	return fmt.Sprintf("iprp | Properties: %d, Associations: %d", len(iprp.PropertyContainer.Properties), len(iprp.Associations))
}

func parseIprp(outer *box) (Box, error) {
	return parseItemPropertiesBox(outer)
}

func parseItemPropertiesBox(outer *box) (ip ItemPropertiesBox, err error) {
	// New Reader
	//boxr := outer.newReader(outer.remain)
	var inner box
	for outer.remain > 4 {
		// Read Box
		if inner, err = outer.readInnerBox(); err != nil {
			// TODO: write error
			break
		}

		if inner.boxType == TypeIpco { // Read ItemPropertyContainerBox
			ip.PropertyContainer, err = parseItemPropertyContainerBox(&inner)
			if err != nil {
				// TODO: write error
				break
			}
		} else if inner.boxType == TypeIpma { // Read ItemPropertyAssociation
			ipma, err := parseItemPropertyAssociation(&inner)
			if err != nil {
				// TODO: write error
				break
			}
			ip.Associations = append(ip.Associations, ipma)
		} else {
			if Debug {
				fmt.Printf("(iprp) Unexpected Box Type: %s, Size: %d", inner.Type(), inner.size)
			}
		}

		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			break
		}
	}
	err = outer.discard(outer.remain)
	return ip, err
}

// ItemPropertyContainerBox is an ISOBMFF "ipco" box
type ItemPropertyContainerBox struct {
	Properties []Box // of ItemProperty or ItemFullProperty
}

// Type returns TypeIpco
func (ipco ItemPropertyContainerBox) Type() BoxType {
	return TypeIpco
}

func parseIpco(outer *box) (Box, error) {
	return parseItemPropertyContainerBox(outer)
}

func parseItemPropertyContainerBox(outer *box) (ipc ItemPropertyContainerBox, err error) {
	var p Box
	var inner box
	for outer.remain > 4 {
		inner, err = outer.readInnerBox()
		if err != nil {
			if err == io.EOF {
				return ipc, nil
			}
			outer.err = err
			return ipc, err
		}
		p, err = inner.Parse()
		if Debug {
			fmt.Printf("(ipco) %T %s ", p, p)
			fmt.Printf("\t[ Outer: %d, Size: %d, Inner: %d ]", outer.remain, inner.size, inner.remain)
			if err != nil {
				fmt.Printf("error: %s", err)
			}
			fmt.Printf("\n")
		}
		ipc.Properties = append(ipc.Properties, p)

		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			break
		}
	}
	err = outer.discard(outer.remain)
	return ipc, err
}

// ItemPropertyAssociation is an ISOBMFF "ipma" box
type ItemPropertyAssociation struct {
	Flags      Flags
	EntryCount uint32
	Entries    []ItemPropertyAssociationItem
}

// Type returns TypeIpma
func (ipma ItemPropertyAssociation) Type() BoxType {
	return TypeIpma
}

func parseIpma(outer *box) (Box, error) {
	return parseItemPropertyAssociation(outer)
}

func parseItemPropertyAssociation(outer *box) (ipa ItemPropertyAssociation, err error) {
	ipa.Flags, err = outer.readFlags()
	if err != nil {
		return
	}
	ipa.EntryCount, err = outer.readUint32()
	if err != nil {
		// TODO: Error handling
		return
	}

	// Entries
	ipa.Entries = make([]ItemPropertyAssociationItem, 0, ipa.EntryCount)

	for i := uint32(0); i < ipa.EntryCount && outer.ok(); i++ {
		var itemID uint32
		if ipa.Flags.Version() < 1 {
			itemID16, _ := outer.readUint16()
			itemID = uint32(itemID16)
		} else {
			itemID, _ = outer.readUint32()
		}
		assocCount, _ := outer.readUint8()
		ipai := ItemPropertyAssociationItem{
			ItemID:            itemID,
			AssociationsCount: int(assocCount),
			Associations:      make([]ItemProperty, 0, assocCount),
		}
		for j := 0; j < int(assocCount) && outer.ok(); j++ {
			first, _ := outer.readUint8()
			essential := first&(1<<7) != 0
			first &^= byte(1 << 7)

			var index uint16
			if ipa.Flags.Flags()&1 != 0 {
				second, _ := outer.readUint8()
				index = uint16(first)<<8 | uint16(second)
			} else {
				index = uint16(first)
			}
			ipai.Associations = append(ipai.Associations, ItemProperty{
				Essential: essential,
				Index:     index,
			})
		}
		ipa.Entries = append(ipa.Entries, ipai)
	}
	if !outer.ok() {
		return ipa, outer.err
	}
	if Debug {
		fmt.Println(ipa)
	}
	return ipa, nil
}

// ItemPropertyAssociationItem is not a box
type ItemPropertyAssociationItem struct {
	ItemID            uint32
	AssociationsCount int            // as declared
	Associations      []ItemProperty // as parsed
}

// ItemProperty is not a box
type ItemProperty struct {
	Essential bool
	Index     uint16
}

// ImageSpatialExtentsProperty is an "ispe" Property
type ImageSpatialExtentsProperty struct {
	Flags
	W uint32
	H uint32
}

func (ispe ImageSpatialExtentsProperty) String() string {
	return fmt.Sprintf("(ispe) Image Width: %d, Height: %d", ispe.W, ispe.H)
}

// Type returns TypeIspe
func (ispe ImageSpatialExtentsProperty) Type() BoxType {
	return TypeIspe
}

func parseImageSpatialExtentsProperty(outer *box) (Box, error) {
	flags, err := outer.readFlags()
	if err != nil {
		return nil, err
	}
	w, _ := outer.readUint32()
	h, err := outer.readUint32()
	if err != nil {
		return nil, err
	}
	return ImageSpatialExtentsProperty{
		Flags: flags,
		W:     w,
		H:     h,
	}, nil
}
