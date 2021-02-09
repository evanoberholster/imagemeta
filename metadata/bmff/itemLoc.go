package bmff

import (
	"encoding/binary"
	"fmt"
)

// ItemLocationBox is a "iloc" box
type ItemLocationBox struct {
	Flags

	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4

	ItemCount uint16
	Items     []ItemLocationBoxEntry
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
	return fmt.Sprintf("(Box) iloc | ItemCount: %d, Flags: %d, Version: %d, OffsetSize: %d, LengthSize: %d, BaseOffsetSize: %d, indexSize: %d", iloc.ItemCount, iloc.Flags.Flags(), iloc.Flags.Version(), iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize)
}

func parseItemLocationBox(outer *box) (Box, error) {
	flags, err := outer.r.readFlags()
	if err != nil {
		return nil, err
	}
	ilb := &ItemLocationBox{
		Flags: flags,
	}
	buf, err := outer.r.Peek(4)
	if err != nil {
		return nil, err
	}
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if flags.Version() > 0 { // version 1
		ilb.indexSize = buf[1] & 15
	}

	ilb.ItemCount = binary.BigEndian.Uint16(buf[2:4])
	outer.r.discard(4)

	if Debug {
		fmt.Println(ilb)
	}
	for i := 0; outer.r.ok() && outer.r.anyRemain() && i < int(ilb.ItemCount); i++ {

		var ent ItemLocationBoxEntry
		ent.ItemID, _ = outer.r.readUint16()

		if ilb.baseOffsetSize > 0 {
			outer.r.discard(int(ilb.baseOffsetSize))
		}
		//p2, _ := outer.r.Peek(16)
		////outer.r.discard(4)
		//fmt.Println(p2)
		if flags.Version() > 0 { // version 1
			cmeth, _ := outer.r.readUint16()
			ent.ConstructionMethod = byte(cmeth & 15)
		}
		ent.DataReferenceIndex, _ = outer.r.readUint16()
		if outer.r.ok() && ilb.baseOffsetSize > 0 {
			outer.r.discard(int(ilb.baseOffsetSize) / 8)
		}
		ent.ExtentCount, _ = outer.r.readUint16()
		for j := 0; outer.r.ok() && j < int(ent.ExtentCount); j++ {
			var ol OffsetLength
			ol.Offset, _ = outer.r.readUintN(ilb.offsetSize * 8)
			ol.Length, _ = outer.r.readUintN(ilb.lengthSize * 8)
			if outer.r.err != nil {
				return nil, outer.r.err
			}
			if j == 0 {
				ent.FirstExtent = ol
				continue
			}
			ent.Extents = append(ent.Extents, ol)
		}
		//fmt.Println(ent, outer.r.remain)
		ilb.Items = append(ilb.Items, ent)
	}
	if !outer.r.ok() {
		return nil, outer.r.err
	}
	return ilb, nil
}

// ItemLocationBoxEntry is not a box
type ItemLocationBoxEntry struct {
	ItemID             uint16
	ConstructionMethod uint8 // actually uint4
	DataReferenceIndex uint16
	BaseOffset         uint64 // uint32 or uint64, depending on encoding
	ExtentCount        uint16
	FirstExtent        OffsetLength
	Extents            []OffsetLength
}

func (ilbe ItemLocationBoxEntry) String() string {
	return fmt.Sprintf("ItemID: %d, ConstructionMethod: %d, DataReferenceIndex: %d, BaseOffset: %d, ExtentCount: %d, FirstExtent: %s", ilbe.ItemID, ilbe.ConstructionMethod, ilbe.DataReferenceIndex, ilbe.BaseOffset, ilbe.ExtentCount, ilbe.FirstExtent)
}

// OffsetLength contains an offset and length
type OffsetLength struct {
	Offset, Length uint64
}

func (ol OffsetLength) String() string {
	return fmt.Sprintf("{ Offset: %d, Length: %d }", ol.Offset, ol.Length)
}
