package isobmff

import (
	"github.com/rs/zerolog"
)

func readIinf(b *box) (iinf IinfBox, err error) {
	if err = b.readFlags(); err != nil {
		return iinf, err
	}
	count, err := b.readUint16()
	if err != nil {
		return iinf, err
	}
	if logLevelInfo() {
		logInfoBoxExt(b, zerolog.InfoLevel).Uint16("count", count).Send()
	}
	iinf = make(IinfBox, count)
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeInfe:
			_, err = readInfe(&inner)
		default:
			if logLevelDebug() {
				logInfoBoxExt(&inner, zerolog.DebugLevel).Send()
			}
			inner.close()
		}
		if err != nil {
			return
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
