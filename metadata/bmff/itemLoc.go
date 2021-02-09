package bmff

// ItemLocationBox is a "iloc" box
type ItemLocationBox struct {
	FullBox

	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4

	ItemCount uint16
	Items     []ItemLocationBoxEntry
}

// ItemLocationBoxEntry is not a box
type ItemLocationBoxEntry struct {
	ItemID             uint16
	ConstructionMethod uint8 // actually uint4
	DataReferenceIndex uint16
	BaseOffset         uint64 // uint32 or uint64, depending on encoding
	ExtentCount        uint16
	Extents            []OffsetLength
}

// OffsetLength contains an offset and length
type OffsetLength struct {
	Offset, Length uint64
}
