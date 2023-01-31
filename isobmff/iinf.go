package isobmff

import (
	"github.com/rs/zerolog"
)

func (r *Reader) readIinf(b *box) (iinf IinfBox, err error) {
	if err = b.readFlags(); err != nil {
		return iinf, err
	}
	count, err := b.readUint16()
	if err != nil {
		return iinf, err
	}
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Uint16("count", count).Send()
	}
	if optionSpeed == 0 {
		iinf = make(IinfBox, 0, count)
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		var infe ItemInfoEntry
		switch inner.boxType {
		case typeInfe:
			infe, err = readInfe(&inner)
			if optionSpeed == 0 {
				iinf = append(iinf, infe)
			}
			if infe.itemType == itemTypeExif {
				r.heic.exif.id = infe.itemID
			}
			if infe.itemType == itemTypeMime {
				r.heic.xml.id = infe.itemID
			}
		default:
			if logLevelDebug() {
				logBoxExt(&inner, zerolog.DebugLevel).Send()
			}
		}
		if err != nil && logLevelError() {
			logBoxExt(&inner, zerolog.ErrorLevel).Err(err).Send()
		}
		if err = inner.close(); err != nil {
			logBoxExt(&inner, zerolog.ErrorLevel).Err(err).Send()
			break
		}
	}
	return iinf, b.close()
}

func itemTypeFromBuf(buf []byte) itemType {
	if it, ok := mapItemType[string(buf[:4])]; ok {
		return it
	}
	return itemTypeUnknown
}

// IinfBox represents an "iinf" box.
// size int64
// Flags Flags
// Count uint16
type IinfBox []ItemInfoEntry
