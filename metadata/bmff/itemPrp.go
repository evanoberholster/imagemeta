package bmff

import (
	"fmt"
	"io"
)

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

func (iprp *ItemPropertiesBox) setBox(b Box) {
	switch v := b.(type) {
	case ItemPropertyContainerBox:
		iprp.PropertyContainer = v
	case ItemPropertyAssociation:
		iprp.Associations = append(iprp.Associations, v)
	default:
		if Debug {
			fmt.Printf("(iprp) Notset %T", v)
		}
	}
}

func parseItemPropertiesBox(outer *box) (b Box, err error) {
	ip := ItemPropertiesBox{}

	// New Reader
	boxr := outer.newReader(outer.r.remain)
	var inner box
	for outer.r.remain > 4 {
		// Read Box
		if inner, err = boxr.readBox(); err != nil {
			return ip, err
		}

		if inner.boxType == TypeIpco { // Read ItemPropertyContainerBox
			ip.PropertyContainer, err = parseItemPropertyContainerBox(&inner)
			if err != nil {
				return ip, err
			}
		} else if inner.boxType == TypeIpma { // Read ItemPropertyAssociation
			ipma, err := parseItemPropertyAssociation(&inner)
			if err != nil {
				return ipma, err
			}
			ip.Associations = append(ip.Associations, ipma)
		} else {
			if Debug {
				fmt.Printf("(iprp) unexpected Type: %s, Size: %d", inner.Type(), inner.size)
			}
		}
		if inner.r.remain > 0 {
			inner.r.discard(inner.r.remain)
		}
		outer.r.remain -= int(inner.size)
	}
	boxr.br.discard(outer.r.remain)
	return ip, nil
}

// ItemPropertyContainerBox is an ISOBMFF "ipco" box
type ItemPropertyContainerBox struct {
	//*box
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
	// New Reader
	boxr := outer.newReader(outer.r.remain)
	var p Box
	var inner box
	for outer.r.remain > 4 {
		inner, err = boxr.readBox()
		if err != nil {
			if err == io.EOF {
				return ipc, nil
			}
			boxr.br.err = err
			return ipc, err
		}
		p, err = inner.Parse()
		if Debug {
			fmt.Printf("(ipco) %T %s ", p, p)
			fmt.Printf("\t[ Outer: %d, Size: %d, Inner: %d ]", outer.r.remain, inner.size, inner.r.remain)
			if err != nil {
				fmt.Printf("error: %s", err)
			}
			fmt.Printf("\n")
		}
		ipc.Properties = append(ipc.Properties, p)
		inner.r.discard(int(inner.r.remain))
		outer.r.remain -= int(inner.size)
	}
	outer.r.discard(int(outer.r.remain))
	return ipc, nil
}

// ItemPropertyAssociation is an ISOBMFF "ipma" box
type ItemPropertyAssociation struct {
	Flags      Flags
	EntryCount uint32
	Entries    []ItemPropertyAssociationItem
}

// Size returns the size of the ItemPropertyAssociation
func (ipma ItemPropertyAssociation) Size() int64 {
	return 0
}

// Type returns TypeIpma
func (ipma ItemPropertyAssociation) Type() BoxType {
	return TypeIpma
}

func parseIpma(outer *box) (Box, error) {
	return parseItemPropertyAssociation(outer)
}

func parseItemPropertyAssociation(outer *box) (ipa ItemPropertyAssociation, err error) {
	ipa.Flags, err = outer.r.readFlags()
	if err != nil {
		return ipa, err
	}
	ipa.EntryCount, _ = outer.r.readUint32()

	// Entries
	ipa.Entries = make([]ItemPropertyAssociationItem, 0, ipa.EntryCount)

	for i := uint32(0); i < ipa.EntryCount && outer.r.ok(); i++ {
		var itemID uint32
		if ipa.Flags.Version() < 1 {
			itemID16, _ := outer.r.readUint16()
			itemID = uint32(itemID16)
		} else {
			itemID, _ = outer.r.readUint32()
		}
		assocCount, _ := outer.r.readUint8()
		ipai := ItemPropertyAssociationItem{
			ItemID:            itemID,
			AssociationsCount: int(assocCount),
			Associations:      make([]ItemProperty, 0, assocCount),
		}
		for j := 0; j < int(assocCount) && outer.r.ok(); j++ {
			first, _ := outer.r.readUint8()
			essential := first&(1<<7) != 0
			first &^= byte(1 << 7)

			var index uint16
			if ipa.Flags.Flags()&1 != 0 {
				second, _ := outer.r.readUint8()
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
	if !outer.r.ok() {
		return ipa, outer.r.err
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
	flags, err := outer.r.readFlags()
	if err != nil {
		return nil, err
	}
	w, _ := outer.r.readUint32()
	h, err := outer.r.readUint32()
	if err != nil {
		return nil, err
	}
	return ImageSpatialExtentsProperty{
		Flags: flags,
		W:     w,
		H:     h,
	}, nil
}

// ImageRotation is a ISOBMFF "irot" rotation property.
// Represents the Image Rotation Angle.
// 1 means 90 degrees counter-clockwise, 2 means 180 counter-clockwise
type ImageRotation uint8

// Type returns TypeIrot
func (irot ImageRotation) Type() BoxType {
	return TypeIrot
}

func (irot ImageRotation) String() string {
	switch irot {
	case 0:
		return fmt.Sprintf("(irot) No Rotation")
	case 1:
		return fmt.Sprintf("(irot) Angle: 90° Counter-Clockwise")
	case 2:
		return fmt.Sprintf("(irot) Angle: 180° Counter-Clockwise")
	case 3:
		return fmt.Sprintf("(irot) Angle: 270° Counter-Clockwise")
	default:
		return fmt.Sprintf("(irot) Unknown Angle: %d", irot)
	}
}

func parseImageRotation(outer *box) (Box, error) {
	v, err := outer.r.readUint8()
	if err != nil {
		return nil, err
	}
	return ImageRotation(v & 3), nil
}
