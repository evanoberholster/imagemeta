package isobmff

import (
	"fmt"
	"strings"
)

func (r *Reader) readIinf(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return err
	}

	var count uint32
	switch b.flags.version() {
	case 0:
		c, readErr := b.readUint16()
		if readErr != nil {
			return readErr
		}
		count = uint32(c)
	case 1:
		c, readErr := b.readUint32()
		if readErr != nil {
			return readErr
		}
		count = c
	default:
		return fmt.Errorf("readIinf: unsupported version %d", b.flags.version())
	}

	if logLevelInfo() {
		logInfo().Object("box", b).Uint32("count", count).Send()
	}

	var (
		parsed uint32
		inner  box
		ok     bool
	)
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		if inner.boxType == typeInfe {
			parsed++
			if err = r.readInfe(&inner); err != nil {
				if logLevelError() {
					logError().Object("box", inner).Err(err).Send()
				}
				return err
			}
		}
		if err = inner.close(); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	if logLevelDebug() && parsed != count {
		logDebug().Object("box", b).Uint32("declared", count).Uint32("parsed", parsed).Msg("iinf entry count mismatch")
	}
	return b.close()
}

func (r *Reader) readInfe(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return err
	}

	var id itemID
	switch b.flags.version() {
	case 2:
		v, readErr := b.readUint16()
		if readErr != nil {
			return readErr
		}
		id = itemID(v)
	case 3:
		v, readErr := b.readUint32()
		if readErr != nil {
			return readErr
		}
		id = itemID(v)
	default:
		if logLevelDebug() {
			logDebug().Object("box", b).Uint8("version", b.flags.version()).Msg("skipping unsupported infe version")
		}
		return nil
	}

	protectionIndex, err := b.readUint16()
	if err != nil {
		return err
	}
	itemFourCC, err := b.readFourCC()
	if err != nil {
		return err
	}
	var itemTypeBuf [4]byte
	bmffEndian.PutUint32(itemTypeBuf[:], itemFourCC)
	itemType := itemTypeFromBuf(itemTypeBuf[:])

	// item_name
	if err = b.discardCString(maxBoxStringLength); err != nil {
		return err
	}

	var contentType string
	switch itemType {
	case itemTypeMime:
		contentType, err = b.readCString(maxBoxStringLength)
		if err != nil {
			return err
		}
		if isXMPMIMEType(contentType) {
			r.heic.xml.id = id
		}
	case itemTypeExif:
		r.heic.exif.id = id
	case itemTypeURI:
		// item_uri_type
		if err = b.discardCString(maxBoxStringLength); err != nil {
			return err
		}
	}

	if logLevelDebug() {
		ev := logDebug().
			Object("box", b).
			Object("flags", b.flags).
			Uint32("itemID", uint32(id)).
			Str("itemType", string(itemTypeBuf[:])).
			Uint16("idx", protectionIndex)
		if itemType == itemTypeMime {
			ev.Str("contentType", contentType)
		}
		ev.Send()
	}
	return nil
}

func isXMPMIMEType(contentType string) bool {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	switch ct {
	case "application/rdf+xml", "application/xml", "text/xml":
		return true
	}
	return strings.Contains(ct, "xmp") || strings.Contains(ct, "rdf+xml")
}
