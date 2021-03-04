package bmff

import (
	"fmt"
)

// HandlerType always 4 bytes; usually "pict" for HEIF images.
type HandlerType uint8

const (
	handlerUnknown HandlerType = iota
	handlerPict
	handlerVide
	handlerMeta
)

func (ht HandlerType) String() string {
	switch ht {
	case handlerPict:
		return "pict"
	case handlerVide:
		return "vide"
	case handlerMeta:
		return "meta"
	default:
		return "nnnn"
	}
}

func handler(buf []byte) HandlerType {
	if isHandler(buf, "pict") {
		return handlerPict
	}
	if isHandler(buf, "vide") {
		return handlerVide
	}
	if isHandler(buf, "meta") {
		return handlerMeta
	}
	return handlerUnknown
}

func isHandler(buf []byte, str string) bool {
	return buf[0] == str[0] && buf[1] == str[1] && buf[2] == str[2] && buf[3] == str[3]
}

// HandlerBox is a "hdlr" box.
//
// Handler box hdlr tells the metadata type. For HEIF, it is always ’pict’.
type HandlerBox struct {
	Flags Flags
	size  uint32
	//Name        string
	HandlerType HandlerType
}

// Type returns TypeHdlr.
func (hdlr HandlerBox) Type() BoxType {
	return TypeHdlr
}

func parseHdlr(outer *box) (Box, error) {
	return outer.parseHandlerBox()
}

func (b *box) parseHandlerBox() (hdlr HandlerBox, err error) {
	hdlr.size = uint32(b.size)
	buf, err := b.Peek(24)
	if err != nil {
		return
	}
	if err = b.discard(24); err != nil {
		return
	}

	hdlr.Flags = Flags(heicByteOrder.Uint32(buf[:4]))
	hdlr.HandlerType = handler(buf[8:12])
	if hdlr.HandlerType == handlerUnknown {
		err = fmt.Errorf("error Handler type unknown: %s", string(buf[8:12]))
		return
	}
	//hdlr.Name, _ = outer.readString()
	return hdlr, b.discard(b.remain)
}

// ItemTypeReferenceBox is an "iref" box.
//
// Item Reference box iref enables creating directional links from an item to one or several other items.
// Item references are extensively used by HEIF. For instance, thumbnail images are recognized from a thumbnail
// type reference which links from the thumbnail image to the master image.
type ItemTypeReferenceBox struct {
	Flags Flags
	size  uint32
}

// Type returns TypeIref.
func (iref ItemTypeReferenceBox) Type() BoxType {
	return TypeIref
}

func parseIref(outer *box) (Box, error) {
	return outer.parseItemTypeReferenceBox()
}
func (b *box) parseItemTypeReferenceBox() (iref ItemTypeReferenceBox, err error) {
	iref.size = uint32(b.size)
	iref.Flags, err = b.readFlags()
	if err != nil {
		return
	}
	var inner box
	for b.anyRemain() {
		if inner, err = b.readInnerBox(); err != nil {
			// TODO: write error
			break
		}
		// dimg -> derived image
		// thmb -> thumbnail
		// cdsc -> context description ref / exif
		//fmt.Println(inner, outer.r.remain)

		if err = b.closeInnerBox(&inner); err != nil {
			break
		}
	}
	return iref, b.discard(b.remain)
}

// ImageRotation is an "irot" - image rotation property.
// Represents the Image Rotation Angle at 90 degree intervals.
type ImageRotation uint8

// Type returns TypeIrot
func (irot ImageRotation) Type() BoxType {
	return TypeIrot
}

func (irot ImageRotation) String() string {
	if irot == 0 {
		return "(irot) No Rotation"
	}
	if irot >= 1 && irot <= 3 {
		return fmt.Sprintf("(irot) Angle: %d° Counter-Clockwise", irot*90)
	}
	return fmt.Sprintf("(irot) Unknown Angle: %d", irot)
}

func parseIrot(outer *box) (Box, error) {
	return outer.parseImageRotation()
}

func (b *box) parseImageRotation() (ImageRotation, error) {
	v, err := b.readUint8()
	return ImageRotation(v & 3), err
}

// PrimaryItemBox is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
type PrimaryItemBox struct {
	Flags  Flags
	ItemID uint16
}

func (pitm PrimaryItemBox) String() string {
	return fmt.Sprintf("pitm | ItemID: %d, Flags: %d, Version: %d ", pitm.ItemID, pitm.Flags.Flags(), pitm.Flags.Version())
}

// Type returns TypePitm
func (pitm PrimaryItemBox) Type() BoxType {
	return TypePitm
}

func parsePitm(outer *box) (Box, error) {
	return outer.parsePrimaryItemBox()
}

func (b *box) parsePrimaryItemBox() (pitm PrimaryItemBox, err error) {
	pitm.Flags, err = b.readFlags()
	if err != nil {
		return
	}
	pitm.ItemID, err = b.readUint16()
	if err != nil {
		return
	}
	return pitm, err
}
