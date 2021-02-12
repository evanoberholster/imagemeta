package bmff

import (
	"encoding/binary"
	"fmt"
)

// ItemLocationBox is a "iloc" box
type ItemLocationBox struct {
	Items []ItemLocationBoxEntry
	Flags
	ItemCount uint16

	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// Size returns the size of the ItemLocationBox
func (iloc ItemLocationBox) Size() int64 {
	return 0 // TODO: int64(mb.size)
}

// Type returns TypeIloc
func (iloc ItemLocationBox) Type() BoxType {
	return TypeIloc
}

func (iloc ItemLocationBox) String() string {
	return fmt.Sprintf("iloc | ItemCount: %d, Flags: %d, Version: %d, OffsetSize: %d, LengthSize: %d, BaseOffsetSize: %d, indexSize: %d", iloc.ItemCount, iloc.Flags.Flags(), iloc.Flags.Version(), iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize)
}

func parseIloc(outer *box) (Box, error) {
	return parseItemLocationBox(outer)
}

func parseItemLocationBox(outer *box) (ilb ItemLocationBox, err error) {
	flags, err := outer.r.readFlags()
	if err != nil {
		return
	}
	ilb = ItemLocationBox{
		Flags: flags,
	}
	buf, err := outer.r.Peek(4)
	if err != nil {
		return
	}
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if flags.Version() > 0 { // version 1
		ilb.indexSize = buf[1] & 15
	}

	ilb.ItemCount = binary.BigEndian.Uint16(buf[2:4])
	outer.r.discard(4)
	ilb.Items = make([]ItemLocationBoxEntry, 0, ilb.ItemCount)

	if Debug {
		fmt.Println(ilb)
	}
	for i := 0; outer.r.ok() && outer.r.anyRemain() && i < int(ilb.ItemCount); i++ {

		var ent ItemLocationBoxEntry
		ent.ItemID, _ = outer.r.readUint16()

		if flags.Version() > 0 { // version 1
			cmeth, _ := outer.r.readUint16()
			ent.ConstructionMethod = byte(cmeth & 15)
		}
		ent.DataReferenceIndex, _ = outer.r.readUint16()

		// Adjust for baseOffset per issue "https://github.com/go4org/go4/issues/47" thanks to petercgrant
		if outer.r.ok() && ilb.baseOffsetSize > 0 {
			ent.BaseOffset, _ = outer.r.readUintN(ilb.baseOffsetSize * 8)
			//outer.r.discard(int(ilb.baseOffsetSize) / 8)
		}

		// ExtentCount
		ent.ExtentCount, _ = outer.r.readUint16()
		for j := 0; outer.r.ok() && j < int(ent.ExtentCount); j++ {
			var ol OffsetLength
			ol.Offset, _ = outer.r.readUintN(ilb.offsetSize * 8)
			ol.Length, _ = outer.r.readUintN(ilb.lengthSize * 8)
			if outer.r.err != nil {
				err = outer.r.err
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
			fmt.Println(ent, outer.r.remain)
		}
	}
	if !outer.r.ok() {
		err = outer.r.err
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
