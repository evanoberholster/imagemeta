package isobmff

import (
	"fmt"

	"github.com/rs/zerolog"
)

// itemLocationBox is a "iloc" box
type itemLocationBox struct {
	items                                             []ilocEntry
	count                                             uint32
	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// MarshalZerologObject is a zerolog interface for logging
func (ilb itemLocationBox) MarshalZerologObject(e *zerolog.Event) {
	e.Array("entries", ilb).Uint32("items", ilb.count).Uint8("offsetSize", ilb.offsetSize).Uint8("lengthSize", ilb.lengthSize).Uint8("baseOffsetSize", ilb.baseOffsetSize).Uint8("indexSize", ilb.indexSize)
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
	e.Uint32("itemID", uint32(ie.id)).Object("extent", ie.firstExtent).Uint16("count", ie.count).Uint16("dri", ie.dataReferenceIndex).Uint8("cmeth", ie.constructionMethod)
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

	for i := uint32(0); i < ilb.count; i++ {
		var ent ilocEntry
		switch b.flags.version() {
		case 0, 1:
			v, readErr := b.readUint16()
			if readErr != nil {
				return readErr
			}
			ent.id = itemID(v)
		case 2:
			v, readErr := b.readUint32()
			if readErr != nil {
				return readErr
			}
			ent.id = itemID(v)
		default:
			return fmt.Errorf("readIloc: unsupported version %d", b.flags.version())
		}

		if b.flags.version() > 0 { // versions 1 and 2
			cmeth, readErr := b.readUint16()
			if readErr != nil {
				return readErr
			}
			ent.constructionMethod = uint8(cmeth & 0x0f)
		}
		ent.dataReferenceIndex, err = b.readUint16()
		if err != nil {
			return err
		}

		ent.baseOffset, err = b.readUintN(ilb.baseOffsetSize)
		if err != nil {
			return err
		}
		ent.count, err = b.readUint16()
		if err != nil {
			return err
		}

		for j := 0; j < int(ent.count); j++ {
			if b.flags.version() > 0 && ilb.indexSize > 0 {
				if _, err = b.readUintN(ilb.indexSize); err != nil {
					return err
				}
			}

			extentOffset, readErr := b.readUintN(ilb.offsetSize)
			if readErr != nil {
				return readErr
			}
			extentLength, readErr := b.readUintN(ilb.lengthSize)
			if readErr != nil {
				return readErr
			}

			if j == 0 {
				ent.firstExtent = offsetLength{
					offset: ent.baseOffset + extentOffset,
					length: extentLength,
				}
			}
		}
		if logLevelDebug() {
			logDebug().Object("entry", ent).Send()
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
	if err = b.readFlags(); err != nil {
		return ilb, err
	}

	buf, err := b.Peek(2)
	if err != nil {
		return ilb, fmt.Errorf("readIlocHeader: %w", ErrBufLength)
	}
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if b.flags.version() > 0 { // versions 1 and 2
		ilb.indexSize = buf[1] & 15
	}
	if _, err = b.Discard(2); err != nil {
		return ilb, err
	}

	switch b.flags.version() {
	case 0, 1:
		c, readErr := b.readUint16()
		if readErr != nil {
			return ilb, readErr
		}
		ilb.count = uint32(c)
	case 2:
		ilb.count, err = b.readUint32()
		if err != nil {
			return ilb, err
		}
	default:
		return ilb, fmt.Errorf("readIlocHeader: unsupported version %d", b.flags.version())
	}

	if logLevelInfo() {
		logInfoBox(b).Object("ItemLocation", ilb).Send()
	}
	return ilb, nil
}
