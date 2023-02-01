package isobmff

import (
	"github.com/rs/zerolog"
)

// ItemTypeReferenceBox is an "iref" box.
//
// Item Reference box iref enables creating directional links from an item to one or several other items.
// Item references are extensively used by HEIF. For instance, thumbnail images are recognized from a thumbnail
// type reference which links from the thumbnail image to the master image.
// dimg -> derived image
// thmb -> thumbnail
// cdsc -> context description ref / exif
type ItemTypeReferenceBox struct{}

func readIref(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return
	}
	if logLevelInfo() {
		logBoxExt(b, zerolog.InfoLevel).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		//switch inner.boxType {
		//case typeDimg:
		//case typeThmb:
		//case typeCdsc:
		//}
		// TODO: implement reading iref boxes
		if logLevelInfo() {
			logBoxExt(&inner, zerolog.InfoLevel).Send()
		}
		if err = inner.close(); err != nil && logLevelError() {
			logBoxExt(&inner, zerolog.ErrorLevel).Err(err).Send()
			break
		}
	}
	return b.close()
}
