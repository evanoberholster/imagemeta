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
		err = inner.discard(inner.remain)
	}
	err = outer.discard(outer.remain)
	return ip, err
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
	//boxr := outer.newReader(outer.r.remain)
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
		err = inner.discard(inner.remain)
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
