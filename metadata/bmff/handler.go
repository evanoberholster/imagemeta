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
	Flags       Flags
	size        uint32
	Name        string
	HandlerType HandlerType
}

// Size returns the size of the HandlerBox
func (hdlr HandlerBox) Size() int64 {
	return int64(hdlr.size)
}

// Type returns TypeHdlr.
func (hdlr HandlerBox) Type() BoxType {
	return TypeHdlr
}

func parseHdlr(outer *box) (Box, error) {
	return parseHandlerBox(outer)
}

func parseHandlerBox(outer *box) (hdlr HandlerBox, err error) {
	hdlr.size = uint32(outer.size)
	if hdlr.Flags, err = outer.readFlags(); err != nil {
		return
	}
	buf, err := outer.Peek(20)
	if err != nil {
		return
	}
	hdlr.HandlerType = handler(buf[4:8])
	if err = outer.discard(20); err != nil {
		return
	}
	if hdlr.HandlerType == handlerUnknown {
		if Debug {
			fmt.Println("Unknown Handler: Error", string(buf), buf)
		}
		// err = Unknown Handler. Cancel Parsing file
	}
	hdlr.Name, _ = outer.readString()
	err = outer.discard(outer.remain)
	return hdlr, err
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
	return parseItemTypeReferenceBox(outer)
}
func parseItemTypeReferenceBox(outer *box) (iref ItemTypeReferenceBox, err error) {

	iref.size = uint32(outer.size)
	iref.Flags, err = outer.readFlags()
	if err != nil {
		return
	}
	var inner box
	for outer.anyRemain() {
		// Read Box
		if inner, err = outer.readInnerBox(); err != nil {
			// TODO: write error
			break
		}
		// dimg -> derived image
		// thmb -> thumbnail
		// cdsc -> context description ref / exif
		//fmt.Println(inner, outer.r.remain)
		if inner.remain > 0 {
			err = inner.discard(inner.remain)
		}
		outer.remain -= int(inner.size)
	}
	err = outer.discard(outer.remain)
	return
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

func parseImageRotation(outer *box) (Box, error) {
	v, err := outer.readUint8()
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
	return parsePrimaryItemBox(outer)
}

func parsePrimaryItemBox(outer *box) (pitm PrimaryItemBox, err error) {
	pitm.Flags, err = outer.readFlags()
	if err != nil {
		return
	}
	pitm.ItemID, err = outer.readUint16()
	if err != nil {
		return
	}
	return pitm, err
}
