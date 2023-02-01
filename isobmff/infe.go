package isobmff

import (
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func readInfe(b *box) (infe ItemInfoEntry, err error) {
	infeHeaderSize := 13
	buf, err := b.Peek(infeHeaderSize)
	if err != nil {
		return
	}
	b.readFlagsFromBuf(buf[:4])

	// Only support Infe version 2
	if b.flags.version() != 2 {
		if logLevelError() {
			logError().Object("box", b).Err(errors.Wrapf(ErrInfeVersionNotSupported, "found version %d infe box. Only 2 is supported now", b.flags.version())).Send()
		}
		err = b.close()
		return
	}
	var debugItemType string
	infe.itemID = itemID(bmffEndian.Uint16(buf[4:6]))
	infe.protectionIndex = bmffEndian.Uint16(buf[6:8])
	infe.itemType = itemTypeFromBuf(buf[8:12])

	if logLevelDebug() {
		debugItemType = string(buf[8:12])
	}

	// expect whitespace
	if buf[12] != '\x00' {
		if logLevelDebug() {
			logBoxExt(b, zerolog.DebugLevel).Str("itemType", debugItemType).Uint16("itemID", uint16(infe.itemID)).Uint16("idx", infe.protectionIndex).Msg("does't end on whitespace")
		}
		infeHeaderSize--
	}
	if infe.itemType == itemTypeMime {
		buf, err = b.Peek(b.remain)
		if err != nil {
			return
		}
		infe.contentType = imagetype.FromString(string(buf[infeHeaderSize : len(buf)-2]))
	}
	// TODO: implement URI type
	if logLevelDebug() {
		ev := logBoxExt(b, zerolog.DebugLevel).Str("itemType", debugItemType).Uint16("itemID", uint16(infe.itemID)).Uint16("idx", infe.protectionIndex)
		if infe.itemType == itemTypeMime {
			ev.Str("contentType", infe.contentType.String())
		}
		ev.Send()
	}
	return infe, b.close()
}

// ItemInfoEntry represents an "infe" box.
//
// TODO: currently only parses Version 2 boxes.
type ItemInfoEntry struct {
	itemID          itemID
	protectionIndex uint16
	itemType        itemType
	contentType     imagetype.ImageType

	// If Type == "uri ":
	//ItemURIType string
}

// ItemType

type itemType uint8

const (
	itemTypeUnknown itemType = iota
	itemTypeInfe
	itemTypeMime
	itemTypeURI
	itemTypeAv01
	itemTypeHvc1
	itemTypeGrid
	itemTypeExif
)

var mapItemType = map[string]itemType{
	"infe": itemTypeInfe,
	"mime": itemTypeMime,
	"uri ": itemTypeURI,
	"av01": itemTypeAv01,
	"hvc1": itemTypeHvc1,
	"grid": itemTypeGrid,
	"Exif": itemTypeExif,
}

var mapItemTypeString = map[itemType]string{
	itemTypeInfe: "infe",
	itemTypeMime: "mime",
	itemTypeURI:  "uri ",
	itemTypeAv01: "av01",
	itemTypeHvc1: "hvc1",
	itemTypeGrid: "grid",
	itemTypeExif: "Exif",
}

func itemTypeFromBuf(buf []byte) itemType {
	str := string(buf[:4])
	if str == "hvc1" {
		return itemTypeHvc1
	}
	switch str {
	case "Exif":
		return itemTypeExif
	case "av01":
		return itemTypeAv01
	case "grid":
		return itemTypeGrid
	case "infe":
		return itemTypeInfe
	case "mime":
		return itemTypeMime
	case "uri ":
		return itemTypeURI
	default:
		return itemTypeUnknown
	}
}
