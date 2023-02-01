package isobmff

import (
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// readIinf reads an "iinf" box
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

// IinfBox represents an "iinf" box.
// size int64
// Flags Flags
// Count uint16
type IinfBox []ItemInfoEntry

func (r *Reader) readIinfFast(b *box) (iinf IinfBox, err error) {
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
	if err = r.readInfeFast(b); err != nil && logLevelError() {
		logError().Object("box", b).Err(err).Send()
	}

	return iinf, b.close()
}

func (r *Reader) readInfeFast(b *box) (err error) {
	buf, err := b.Peek(b.remain)
	if err != nil {
		return
	}
	offset := b.offset + int(b.size) - b.remain

	for i := 0; i < len(buf); {
		infeFastHeaderSize := 21

		var contentType imagetype.ImageType
		size := int(bmffEndian.Uint32(buf[i : i+4]))
		boxType := boxTypeFromBuf(buf[i+4 : i+8])
		flags := flags(bmffEndian.Uint32(buf[i+8 : i+12]))

		if boxType != typeInfe {
			i += size
			continue
		}
		// Only support Infe version 2
		if flags.version() != 2 {
			if logLevelError() {
				logError().Object("box", b).Err(errors.Wrapf(ErrInfeVersionNotSupported, "found version %d infe box. Only 2 is supported now", flags.version())).Send()
			}
			i += size
			continue
		}

		itemID := itemID(bmffEndian.Uint16(buf[i+12 : i+14]))
		itemType := itemTypeFromBuf(buf[i+16 : i+20])
		// expect whitespace
		if buf[i+20] != '\x00' {
			if logLevelDebug() {
				logBoxExt(b, zerolog.DebugLevel).Str("itemType", string(buf[i+16:i+20])).Uint16("itemID", uint16(itemID)).Msg("does't end on whitespace")
			}
			infeFastHeaderSize--
		}
		switch itemType {
		case itemTypeMime:
			contentType = imagetype.FromString(string(buf[i+infeFastHeaderSize : i+size-1]))
			r.heic.xml.id = itemID
		case itemTypeExif:
			r.heic.exif.id = itemID
		}
		if logLevelDebug() {
			protectionIndex := bmffEndian.Uint16(buf[i+14 : i+16])
			ev := logDebug().Str("BoxType", boxType.String()).Object("flags", flags).Uint16("itemID", uint16(itemID)).Str("itemType", string(buf[i+16:i+20])).Int("offset", i+offset).Int("size", size).Uint16("idx", protectionIndex)
			if itemType == itemTypeMime {
				ev.Str("contentType", contentType.String())
			}
			ev.Send()
		}
		i += size

	}
	return b.close()
}
