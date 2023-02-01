package isobmff

import (
	"github.com/rs/zerolog"
)

// itemLocationBox is a "iloc" box
type itemLocationBox struct {
	items                                             []ilocEntry
	count                                             uint16
	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// MarshalZerologObject is a zerolog interface for logging
func (ilb itemLocationBox) MarshalZerologObject(e *zerolog.Event) {
	e.Array("entries", ilb).Uint16("items", ilb.count).Uint8("offsetSize", ilb.offsetSize).Uint8("lengthSize", ilb.lengthSize).Uint8("baseOffsetSize", ilb.baseOffsetSize).Uint8("indexSize", ilb.indexSize)
}

// MarshalZerologArray is a zerolog interface for logging
func (ilb itemLocationBox) MarshalZerologArray(a *zerolog.Array) {
	for i := 0; i < len(ilb.items); i++ {
		a.Object(ilb.items[i])
	}
}

// ilocEntry is not a box
type ilocEntry struct {
	//extents            []OffsetLength
	firstExtent        offsetLength
	id                 itemID
	baseOffset         uint64 // uint32 or uint64, depending on encoding
	count              uint16
	dataReferenceIndex uint16 // dri
	constructionMethod uint8  // cmeth actually uint4
}

// MarshalZerologObject is a zerolog interface for logging
func (ie ilocEntry) MarshalZerologObject(e *zerolog.Event) {
	e.Uint16("itemID", uint16(ie.id)).Object("extent", ie.firstExtent).Uint16("count", ie.count).Uint16("dri", ie.dataReferenceIndex).Uint8("cmeth", ie.constructionMethod)
}

// offsetLength contains an offset and length
type offsetLength struct {
	offset, length uint64
}

// MarshalZerologObject is a zerolog interface for logging
func (ol offsetLength) MarshalZerologObject(e *zerolog.Event) {
	e.Uint64("length", ol.length).Uint64("offset", ol.offset)
}

func (r *Reader) readIloc(b *box) (err error) {
	ilb, err := readIlocHeader(b)
	if err != nil {
		return err
	}

	buf, err := b.Peek(b.remain)
	if err != nil {
		return
	}

	if optionSpeed == 0 {
		ilb.items = make([]ilocEntry, 0, ilb.count)
	}

	for i := 0; i < len(buf); {
		var ent ilocEntry
		ent.id = itemID(bmffEndian.Uint16(buf[i : i+2]))
		i += 2

		if b.flags.version() > 0 { // version 1
			cmeth := bmffEndian.Uint16(buf[i : i+2])
			ent.constructionMethod = byte(cmeth & 15)
			i += 2
		}
		ent.dataReferenceIndex = bmffEndian.Uint16(buf[i : i+2])
		i += 2

		// Adjust for baseOffset per issue "https://github.com/go4org/go4/issues/47" thanks to petercgrant
		if ilb.baseOffsetSize > 0 {
			ent.baseOffset = uintN(ilb.baseOffsetSize, buf[i:i+int(ilb.baseOffsetSize)])
			i += int(ilb.baseOffsetSize)
		}
		ent.count = bmffEndian.Uint16(buf[i : i+2])
		i += 2

		for j := 0; j < int(ent.count); j++ {
			var ol offsetLength
			if j == 0 {
				ol.offset = uintN(ilb.offsetSize, buf[i:i+int(ilb.offsetSize)])
				i += int(ilb.offsetSize)
				ol.length = uintN(ilb.lengthSize, buf[i:i+int(ilb.lengthSize)])
				i += int(ilb.lengthSize)
				ent.firstExtent = ol
			}
		}
		if optionSpeed == 0 {
			ilb.items = append(ilb.items, ent)
		}
		if logLevelDebug() {
			logBoxExt(nil, zerolog.DebugLevel).Object("entry", ent).Send()
		}

		switch ent.id {
		case r.heic.exif.id:
			r.heic.exif.ol = ent.firstExtent
		case r.heic.xml.id:
			r.heic.xml.ol = ent.firstExtent
		}
		//if ent.ItemID == Exif...
	}
	return b.close()
}

func readIlocHeader(b *box) (ilb itemLocationBox, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		return ilb, err
	}
	b.readFlagsFromBuf(buf)
	buf = buf[4:]
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if b.flags.version() > 0 { // version 1
		ilb.indexSize = buf[1] & 15
	}
	ilb.count = bmffEndian.Uint16(buf[2:4])
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Object("ItemLocation", ilb).Send()
	}
	return ilb, b.Discard(8)
}

func uintN(size uint8, buf []byte) uint64 {
	switch size {
	case 1:
		return uint64(buf[0])
	case 2:
		return uint64(bmffEndian.Uint16(buf[:2]))
	case 4:
		return uint64(bmffEndian.Uint32(buf[:4]))
	case 8:
		return bmffEndian.Uint64(buf[:8])
	default:
		panic("error here")
	}
}
