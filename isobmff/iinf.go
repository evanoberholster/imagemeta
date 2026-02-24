package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/imagetype"
)

func (r *Reader) readIinf(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return
	}
	count, err := b.readUint16()
	if err != nil {
		return
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Uint16("count", count).Send()
	}
	if err = r.readInfe(b); err != nil && logLevelError() {
		logError().Object("box", b).Err(err).Send()
	}

	return b.close()
}

func (r *Reader) readInfe(b *box) (err error) {
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
				logError().Object("box", b).Err(fmt.Errorf("found version %d infe box. Only 2 is supported now: %w", flags.version(), ErrInfeVersionNotSupported)).Send()
			}
			i += size
			continue
		}

		itemID := itemID(bmffEndian.Uint16(buf[i+12 : i+14]))
		itemType := itemTypeFromBuf(buf[i+16 : i+20])
		// expect whitespace
		if buf[i+20] != '\x00' {
			if logLevelDebug() {
				logDebug().Object("box", b).Str("itemType", string(buf[i+16:i+20])).Uint16("itemID", uint16(itemID)).Msg("does't end on whitespace")
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
