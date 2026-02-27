package isobmff

import (
	"fmt"

	"github.com/rs/zerolog"
)

// readHdlr reads an "hdlr" box
func readHdlr(b *box) (ht hdlrType, err error) {
	if !b.isType(typeHdlr) {
		return hdlrUnknown, fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if err = b.readFlags(); err != nil {
		return hdlrUnknown, err
	}

	if b.remain < 8 {
		return hdlrUnknown, fmt.Errorf("readHdlr: %w", ErrBufLength)
	}

	buf, err := b.Peek(8)
	if err != nil {
		return hdlrUnknown, err
	}
	ht = hdlrFromBuf(buf[4:8])
	if logLevelInfo() {
		logInfoBox(b).Str("hdlr", ht.String()).Send()
	}
	return ht, b.close()
}

// hdlrType

// hdlrType always 4 bytes;
// Handler; usually "pict" for HEIF images
type hdlrType uint8

// hdlr types
const (
	hdlrUnknown hdlrType = iota
	hdlrPict
	hdlrVide
	hdlrMeta
)

// String is a stringer interface for hdlrType
func (ht hdlrType) String() string {
	if str, ok := hdlrStringMap[ht]; ok {
		return str
	}
	return "nnnn"
}

var (
	hdlrPictFourCC = fourCCFromString("pict")
	hdlrVideFourCC = fourCCFromString("vide")
	hdlrMetaFourCC = fourCCFromString("meta")
)

func hdlrFromBuf(buf []byte) hdlrType {
	if len(buf) < 4 {
		return hdlrUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case hdlrPictFourCC:
		return hdlrPict
	case hdlrVideFourCC:
		return hdlrVide
	case hdlrMetaFourCC:
		return hdlrMeta
	default:
		return hdlrUnknown
	}
}

var hdlrStringMap = map[hdlrType]string{
	hdlrPict: "pict",
	hdlrVide: "vide",
	hdlrMeta: "meta",
}

// pitmID is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
// 0 represents not set.

func readPitm(b *box) (id itemID, err error) {
	if err = b.readFlags(); err != nil {
		return invalidItemID, err
	}
	switch b.flags.version() {
	case 0:
		v, readErr := b.readUint16()
		if readErr != nil {
			return invalidItemID, readErr
		}
		id = itemID(v)
	case 1:
		v, readErr := b.readUint32()
		if readErr != nil {
			return invalidItemID, readErr
		}
		id = itemID(v)
	default:
		return invalidItemID, fmt.Errorf("readPitm: unsupported version %d", b.flags.version())
	}
	if logLevelInfo() {
		logInfoBox(b).Uint32("ptim", uint32(id)).Send()
	}
	return id, b.close()
}

// itemID
type itemID uint32

const invalidItemID itemID = 0

type itemReference struct {
	referenceType boxType
	fromID        itemID
	toID          itemID
}

type itemProperty struct {
	boxType boxType
	width   uint32
	height  uint32
}

type itemPropertyLink struct {
	itemID        itemID
	propertyIndex uint16
	essential     bool
}

const maxStoredItemGraphEntries = 4096

func (r *Reader) addItemReference(refType boxType, from, to itemID) {
	if len(r.heic.references) >= maxStoredItemGraphEntries {
		return
	}
	r.heic.references = append(r.heic.references, itemReference{
		referenceType: refType,
		fromID:        from,
		toID:          to,
	})
}

func (r *Reader) addItemProperty(prop itemProperty) {
	if len(r.heic.properties) >= maxStoredItemGraphEntries {
		return
	}
	r.heic.properties = append(r.heic.properties, prop)
}

func (r *Reader) addItemPropertyLink(link itemPropertyLink) {
	if len(r.heic.propertyLinks) >= maxStoredItemGraphEntries {
		return
	}
	r.heic.propertyLinks = append(r.heic.propertyLinks, link)
}

// readIref parses iref single-item-type reference records used by HEIF/JXL.
func (r *Reader) readIref(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return
	}
	var itemIDSize uint8
	switch b.flags.version() {
	case 0:
		itemIDSize = 2
	case 1:
		itemIDSize = 4
	default:
		return fmt.Errorf("readIref: unsupported version %d", b.flags.version())
	}
	if logLevelInfo() {
		logInfoBox(b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		if isSupportedItemReferenceType(inner.boxType) {
			err = r.readIrefEntry(&inner, itemIDSize)
		}
		if logLevelInfo() {
			logInfoBox(&inner).Send()
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}

func isSupportedItemReferenceType(bt boxType) bool {
	switch bt {
	case typeCdsc, typeThmb, typeDimg, typeAuxl, typeIovl:
		return true
	default:
		return false
	}
}

func readItemIDBySize(b *box, size uint8) (itemID, error) {
	switch size {
	case 2:
		v, err := b.readUint16()
		return itemID(v), err
	case 4:
		v, err := b.readUint32()
		return itemID(v), err
	default:
		return invalidItemID, ErrUnsupportedFieldSize
	}
}

func (r *Reader) readIrefEntry(b *box, itemIDSize uint8) error {
	fromID, err := readItemIDBySize(b, itemIDSize)
	if err != nil {
		return err
	}
	refCount, err := b.readUint16()
	if err != nil {
		return err
	}
	for i := 0; i < int(refCount); i++ {
		toID, readErr := readItemIDBySize(b, itemIDSize)
		if readErr != nil {
			return readErr
		}
		r.addItemReference(b.boxType, fromID, toID)
	}
	return nil
}

// readIprp walks item property boxes and parses ipco/ipma payloads.
func (r *Reader) readIprp(b *box) (err error) {
	if logLevelInfo() {
		logInfoBox(b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeIpma:
			err = r.readIpma(&inner)
		case typeIpco:
			err = r.readIpco(&inner)
		default:
			if logLevelInfo() {
				logInfoBox(&inner).Send()
			}
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}

// readIpco reads item properties and stores property order/index metadata.
func (r *Reader) readIpco(b *box) (err error) {
	if logLevelInfo() {
		logInfoBox(b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		prop := itemProperty{boxType: inner.boxType}
		if inner.boxType == typeIspe {
			prop.width, prop.height, err = readIspeProperty(&inner)
		}
		if err == nil {
			r.addItemProperty(prop)
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}

func readIspeProperty(b *box) (width, height uint32, err error) {
	if err = b.readFlags(); err != nil {
		return 0, 0, err
	}
	if width, err = b.readUint32(); err != nil {
		return 0, 0, err
	}
	if height, err = b.readUint32(); err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

// readIpma parses item->property links from the association table.
func (r *Reader) readIpma(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return err
	}
	count, err := b.readUint32()
	if err != nil {
		return err
	}
	extendedIndex := b.flags.flags()&1 != 0
	if logLevelInfo() {
		logInfoBox(b).Uint32("entries", count).Send()
	}
	for i := uint32(0); i < count; i++ {
		var id itemID
		if b.flags.version() < 1 {
			v, readErr := b.readUint16()
			if readErr != nil {
				return readErr
			}
			id = itemID(v)
		} else {
			v, readErr := b.readUint32()
			if readErr != nil {
				return readErr
			}
			id = itemID(v)
		}
		v, readErr := b.readUintN(1)
		if readErr != nil {
			return readErr
		}
		associationCount := int(v)

		for j := 0; j < associationCount; j++ {
			var raw uint16
			if extendedIndex {
				raw, err = b.readUint16()
				if err != nil {
					return err
				}
			} else {
				v8, readErr := b.readUintN(1)
				if readErr != nil {
					return readErr
				}
				raw = uint16(v8)
			}

			essential := false
			propertyIndex := uint16(0)
			if extendedIndex {
				essential = raw&0x8000 != 0
				propertyIndex = raw & 0x7fff
			} else {
				essential = raw&0x0080 != 0
				propertyIndex = raw & 0x007f
			}
			if propertyIndex == 0 {
				continue
			}
			r.addItemPropertyLink(itemPropertyLink{
				itemID:        id,
				propertyIndex: propertyIndex,
				essential:     essential,
			})
		}
	}
	return nil
}

// readIdat reads width/height from Canon idat payload used by CR3 metadata tracks.
func readIdat(b *box) (i idat, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return
	}
	i = idat{
		width:  bmffEndian.Uint16(buf[4:6]),
		height: bmffEndian.Uint16(buf[6:8])}
	if logLevelInfo() {
		logInfoBox(b).Object("idat", i).Send()
	}

	return i, b.close()
}

// ItemData is an "idat" box

// idat
type idat struct {
	width, height uint16
}

// MarshalZerologObject is a zerolog interface for logging
func (i idat) MarshalZerologObject(e *zerolog.Event) {
	e.Uint16("width", i.width).Uint16("height", i.height)
}
