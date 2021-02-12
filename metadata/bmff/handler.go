package bmff

import (
	"fmt"
)

// HandlerType always 4 bytes; usually "pict" for HEIF images.
type HandlerType uint8

const (
	handlerUnknown HandlerType = iota
	handlerPict
)

func handler(buf []byte) HandlerType {
	if isBoxpict(buf) {
		return handlerPict
	}
	return handlerUnknown
}

func isBoxpict(buf []byte) bool {
	return buf[0] == 'p' && buf[1] == 'i' && buf[2] == 'c' && buf[3] == 't'
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
	if hdlr.Flags, err = outer.r.readFlags(); err != nil {
		return
	}
	buf, err := outer.r.Peek(20)
	if err != nil {
		return
	}
	hdlr.HandlerType = handler(buf[4:8])
	if err = outer.r.discard(20); err != nil {
		return
	}
	if hdlr.HandlerType == handlerUnknown {
		if Debug {
			fmt.Println("Unknown Handler: Error", string(buf), buf)
		}
		// err = Unknown Handler. Cancel Parsing file
	}
	hdlr.Name, _ = outer.r.readString()
	err = outer.r.discard(int(outer.r.remain))
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
	iref.Flags, err = outer.r.readFlags()
	if err != nil {
		return
	}
	boxr := outer.newReader(outer.r.remain)
	var inner box
	for outer.r.remain > 0 {
		// Read Box
		if inner, err = boxr.readBox(); err != nil {
			// TODO: write error
			break
		}
		// dimg -> derived image
		// thmb -> thumbnail
		// cdsc -> context description ref / exif
		//fmt.Println(inner, outer.r.remain)
		if inner.r.remain > 0 {
			inner.r.discard(inner.r.remain)
		}
		outer.r.remain -= int(inner.size)
	}
	outer.r.discard(outer.r.remain)
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
	v, err := outer.r.readUint8()
	return ImageRotation(v & 3), err
}

// PrimaryItemBox is a "pitm" box.
//
// Primary Item Reference pitm allows setting one image as the primary item.
type PrimaryItemBox struct {
	Flags  Flags
	ItemID uint16
}

// Size returns the size of the PrimaryItemBox
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
	pitm.Flags, err = outer.r.readFlags()
	if err != nil {
		return
	}
	pitm.ItemID, err = outer.r.readUint16()
	if err != nil {
		return
	}
	return pitm, err
}
