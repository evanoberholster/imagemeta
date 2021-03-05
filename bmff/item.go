package bmff

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

var (
	// HeicByteOrder is BigEndian
	heicByteOrder = binary.BigEndian

	// ErrItemNotFound is returned as en error when an Item was not Found.
	ErrItemNotFound = errors.New("item not found")
	// ErrInfeVersionNotSupported is returned when an infe box with an unsupported was found.
	ErrInfeVersionNotSupported = errors.New("infe box version not supported")
)

// ItemType is always 4 bytes
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
	//size int64
	//Flags Flags
	//Count uint16
	ItemInfos []ItemInfoEntry
}

func (iinf ItemInfoBox) String() string {
	return fmt.Sprintf("(iinf) | ItemCount: %d", len(iinf.ItemInfos))
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
	for i := (len(iinf.ItemInfos) - 1); i >= 0; i-- {
		if iinf.ItemInfos[i].ItemID == itemID {
			return iinf.ItemInfos[i], nil
		}
	}
	return ItemInfoEntry{}, ErrItemNotFound
}

func parseIinf(outer *box) (Box, error) {
	return outer.parseItemInfoBox()
}

func (b *box) parseItemInfoBox() (iinf ItemInfoBox, err error) {
	// Read IinfHeader [4]Flags, [2]ItemCount
	buf, err := b.peek(6)
	if err != nil {
		err = errors.Wrap(err, "ParseItemInfoBox")
		return
	}
	_ = b.discard(6)
	// Read Flags
	flags := Flags(heicByteOrder.Uint32(buf[:4]))
	// Read Item count
	count := int(heicByteOrder.Uint16(buf[4:6]))
	iinf.ItemInfos = make([]ItemInfoEntry, int(count))
	if debugFlag {
		traceBoxWithFlags(iinf, *b, flags)
	}

	var inner box
	var infe ItemInfoEntry
	for i := 0; i < count && b.anyRemain(); i++ {
		if inner, err = b.readInnerBox(); err != nil {
			err = errors.Wrap(err, "ParseItemInfoBox (inner)")
			return
		}
		switch inner.Type() {
		case TypeInfe:
			if infe, err = inner.parseItemInfoEntry(); err != nil {
				if debugFlag {
					log.Debug("(infe) error parsing ItemInfoEntry: %s, %s", infe, err.Error())
				}
				err = errors.Wrap(err, "ParseItemInfoBox (infe)")
				return
			}
			iinf.ItemInfos[i] = infe
		default:
			if debugFlag {
				log.Debug("(infe) Unknown BoxType: %s", inner.Type())
			}
		}
		if err = b.closeInnerBox(&inner); err != nil {
			err = errors.Wrap(err, "ParseItemInfoBox (close)")
			return
		}
	}
	return iinf, b.discard(b.remain)
}

// ItemInfoEntry represents an "infe" box.
//
// TODO: currently only parses Version 2 boxes.
type ItemInfoEntry struct {
	//Flags           Flags
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
	return fmt.Sprintf(" \tItemInfoEntry: ItemID:%d, ProtectionIndex:%d, ItemType:%s", infe.ItemID, infe.ProtectionIndex, infe.ItemType)
}

func parseInfe(outer *box) (Box, error) {
	return outer.parseItemInfoEntry()
}

func (b *box) parseItemInfoEntry() (ie ItemInfoEntry, err error) {
	// Read ItemInfoEntry: [4]flags, [2]ItemID, [2]ProtectionIndex, [5]ItemType
	infeHeaderSize := 13
	buf, err := b.peek(infeHeaderSize)
	if err != nil {
		err = errors.Wrap(err, "ParseItemInfoEntry")
		return
	}
	flags := Flags(heicByteOrder.Uint32(buf[:4]))
	if flags.Version() != 2 {
		err = errors.Wrapf(ErrInfeVersionNotSupported, "found version %d infe box. Only 2 is supported now", flags.Version())
		return
	}
	ie.ItemID = heicByteOrder.Uint16(buf[4:6])
	ie.ProtectionIndex = heicByteOrder.Uint16(buf[6:8])
	ie.size = int16(b.size)
	ie.ItemType = itemType(buf[8:12])
	if buf[12] != '\x00' {
		// Read until whitespace
		if debugFlag {
			traceBoxWithMsg(*b, fmt.Sprintf("\t'%s' doesn't end on whitespace. %s", ie.ItemType, flags))
		}
		infeHeaderSize--
	}
	if err = b.discard(infeHeaderSize); err != nil {
		err = errors.Wrap(err, "ParseItemInfoEntry")
		return
	}

	switch ie.ItemType {
	case ItemTypeMime:
		_, _ = b.readString()
		if b.anyRemain() {
			_, _ = b.readString()
		}
		//ie.ContentType, _ = outer.r.readString()
		//if outer.r.anyRemain() {
		//	ie.ContentEncoding, _ = outer.r.readString()
		//}
	case ItemTypeURI:
		_, _ = b.readString()
		//ie.ItemURIType, _ = outer.r.readString()
	}
	if debugFlag {
		traceBoxWithFlags(ie, *b, flags)
	}
	return
}

// ItemLocationBox is a "iloc" box
type ItemLocationBox struct {
	Items     []ItemLocationBoxEntry
	ItemCount uint16

	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// Type returns TypeIloc
func (iloc ItemLocationBox) Type() BoxType {
	return TypeIloc
}

func (iloc ItemLocationBox) String() string {
	return fmt.Sprintf("iloc | ItemCount:%d, OffsetSize:%d, LengthSize:%d, BaseOffsetSize:%d, indexSize:%d", iloc.ItemCount, iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize)
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
	return outer.parseItemLocationBox()
}

func (b *box) parseItemLocationBox() (ilb ItemLocationBox, err error) {
	buf, err := b.peek(8)
	if err != nil {
		err = errors.Wrap(err, "parseItemLocationBoc")
		return
	}
	flags := Flags(heicByteOrder.Uint32(buf[:4]))
	buf = buf[4:]
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if flags.Version() > 0 { // version 1
		ilb.indexSize = buf[1] & 15
	}

	ilb.ItemCount = heicByteOrder.Uint16(buf[2:4])
	if err = b.discard(8); err != nil {
		err = errors.Wrap(err, "parseItemLocationBoc")
		return
	}

	ilb.Items = make([]ItemLocationBoxEntry, 0, ilb.ItemCount)

	if debugFlag {
		traceBox(ilb, *b)
	}
	for i := 0; b.anyRemain() && i < int(ilb.ItemCount) && err == nil; i++ {
		var ent ItemLocationBoxEntry
		// Refactor for performance
		if buf, err = b.peek(6); err != nil {
			err = errors.Wrap(err, "ItemLocationBoxEntry")
			return
		}
		ent.ItemID = heicByteOrder.Uint16(buf[:2])

		if flags.Version() > 0 { // version 1
			cmeth := heicByteOrder.Uint16(buf[2:4])
			ent.ConstructionMethod = byte(cmeth & 15)
			buf = buf[2:]
		}
		ent.DataReferenceIndex = heicByteOrder.Uint16(buf[2:4])
		if err = b.discard(10 - len(buf)); err != nil {
			return
		}

		// Adjust for baseOffset per issue "https://github.com/go4org/go4/issues/47" thanks to petercgrant
		if ilb.baseOffsetSize > 0 && err == nil {
			ent.BaseOffset, err = b.readUintN(ilb.baseOffsetSize * 8)
			if err != nil {
				return
			}
			//outer.r.discard(int(ilb.baseOffsetSize) / 8)
		}

		// ExtentCount
		ent.ExtentCount, err = b.readUint16()
		for j := 0; j < int(ent.ExtentCount) && err == nil; j++ {
			var ol OffsetLength
			ol.Offset, err = b.readUintN(ilb.offsetSize * 8)
			if err != nil {
				break
			}
			ol.Length, err = b.readUintN(ilb.lengthSize * 8)
			if err != nil {
				break
			}
			if j == 0 {
				ent.FirstExtent = ol
				continue
			}
			ent.Extents = append(ent.Extents, ol)
		}

		ilb.Items = append(ilb.Items, ent)

		if debugFlag {
			log.Debug("%s", ent)
		}
	}
	return ilb, err
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
	return fmt.Sprintf("\t ItemID:%d, ConstructionMethod:%d, DataReferenceIndex:%d, BaseOffset:%d, ExtentCount:%d, FirstExtent:%s", ilbe.ItemID, ilbe.ConstructionMethod, ilbe.DataReferenceIndex, ilbe.BaseOffset, ilbe.ExtentCount, ilbe.FirstExtent)
}

// OffsetLength contains an offset and length
type OffsetLength struct {
	Offset, Length uint64
}

func (ol OffsetLength) String() string {
	return fmt.Sprintf("{Offset:%d, Length:%d}", ol.Offset, ol.Length)
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
	return outer.parseItemPropertiesBox()
}

func (b *box) parseItemPropertiesBox() (ip ItemPropertiesBox, err error) {
	var inner box
	for b.remain > 4 {
		// Read Box
		if inner, err = b.readInnerBox(); err != nil {
			err = errors.Wrap(err, "parseItemPropertiesBox")
			return
		}
		switch inner.Type() {
		case TypeIpco:
			// Parse ItemPropertyContainerBox
			ip.PropertyContainer, err = inner.parseItemPropertyContainerBox()
			if err != nil {
				err = errors.Wrap(err, "parseItemPropertiesBox")
				return
			}
		case TypeIpma:
			// Parse ItemPropertyAssociation
			var ipma ItemPropertyAssociation
			ipma, err = inner.parseItemPropertyAssociation()
			if err != nil {
				err = errors.Wrap(err, "parseItemPropertiesBox")
				return
			}
			ip.Associations = append(ip.Associations, ipma)
		default:
			if debugFlag {
				log.Debug("(iprp) Unexpected Box Type: %s, Size: %d", inner.Type(), inner.size)
				traceBox(inner, inner)
			}
		}

		if err = b.closeInnerBox(&inner); err != nil {
			err = errors.Wrap(err, "parseItemPropertiesBox")
			return
		}
	}
	return ip, b.discard(b.remain)
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
	return outer.parseItemPropertyContainerBox()
}

func (b *box) parseItemPropertyContainerBox() (ipc ItemPropertyContainerBox, err error) {
	var p Box
	var inner box
	for b.remain > 4 {
		if inner, err = b.readInnerBox(); err != nil {
			return
		}
		p, err = inner.Parse()
		if err != nil {
			break
		}
		if debugFlag {
			log.Debug("(ipco) %T %s \t", p, p)
			if err != nil {
				fmt.Printf("error: %s", err)
			}
		}
		ipc.Properties = append(ipc.Properties, p)

		if err = b.closeInnerBox(&inner); err != nil {
			break
		}
	}
	return ipc, b.discard(b.remain)
}

// ItemPropertyAssociation is an ISOBMFF "ipma" box
type ItemPropertyAssociation struct {
	//Flags      Flags
	//EntryCount uint32
	Entries []ItemPropertyAssociationItem
}

// Type returns TypeIpma
func (ipma ItemPropertyAssociation) Type() BoxType {
	return TypeIpma
}

func parseIpma(outer *box) (Box, error) {
	return outer.parseItemPropertyAssociation()
}

func (b *box) parseItemPropertyAssociation() (ipa ItemPropertyAssociation, err error) {
	buf, err := b.peek(8)
	if err != nil {
		err = errors.Wrap(err, "parseItemPropertyAssociation")
		return
	}
	_ = b.discard(8)

	flags := Flags(heicByteOrder.Uint32(buf[:4]))
	count := int(heicByteOrder.Uint32(buf[4:8]))

	// Entries
	//	ipa.EntryCount = uint32(count)
	ipa.Entries = make([]ItemPropertyAssociationItem, count)
	for i := 0; i < count && err == nil; i++ {
		var itemID uint32
		if flags.Version() < 1 {
			var itemID16 uint16
			itemID16, err = b.readUint16() //2
			itemID = uint32(itemID16)
		} else {
			itemID, err = b.readUint32() // 4
		}
		assocCount, err := b.readUint8() // 1
		ipai := ItemPropertyAssociationItem{
			ItemID: itemID,
			//AssociationsCount: uint32(assocCount),
			//Associations:      make([]ItemProperty, 0, assocCount),
		}
		var first uint8
		for j := 0; j < int(assocCount) && err == nil; j++ {
			first, err = b.readUint8()
			if err != nil {
				break
			}
			essential := first&(1<<7) != 0
			first &^= byte(1 << 7)

			var index uint16
			var second uint8
			if flags.Flags()&1 != 0 {
				second, err = b.readUint8()
				index = uint16(first)<<8 | uint16(second)
			} else {
				index = uint16(first)
			}
			if j < len(ipai.Associations) {
				ipai.Associations[j] = index
				_ = essential
			}
			//ipai.Associations = append(ipai.Associations, ItemProperty{
			//	Essential: essential,
			//	Index:     index,
			//})
		}
		ipa.Entries[i] = ipai
	}
	if debugFlag {
		traceBox(ipa, *b)
	}
	return ipa, nil
}

// ItemPropertyAssociationItem is not a box
type ItemPropertyAssociationItem struct {
	ItemID       uint32
	Associations [6]uint16
	//AssociationsCount uint32 // as declared
	//Associations      []ItemProperty // as parsed
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
	return fmt.Sprintf("(ispe) Image Width:%d, Height:%d", ispe.W, ispe.H)
}

// Type returns TypeIspe
func (ispe ImageSpatialExtentsProperty) Type() BoxType {
	return TypeIspe
}

func parseIspe(b *box) (Box, error) {
	return b.parseImageSpatialExtentsProperty()
}

func (b *box) parseImageSpatialExtentsProperty() (ispe ImageSpatialExtentsProperty, err error) {
	buf, err := b.peek(12)
	if err != nil {
		err = errors.Wrap(err, "parseImageSpatialExtentsProperty")
		return
	}
	return ImageSpatialExtentsProperty{
		Flags: Flags(heicByteOrder.Uint32(buf[:4])),
		W:     heicByteOrder.Uint32(buf[4:8]),
		H:     heicByteOrder.Uint32(buf[8:12]),
	}, b.discard(12)
}
