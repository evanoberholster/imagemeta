package exif

import (
	"log"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

var (
	InfoLogger *log.Logger
)

func isInfo() bool {
	return InfoLogger != nil
}

func logTagInfo(ifd ifds.Ifd, t tag.Tag, offset uint32) {
	InfoLogger.Printf("Tag: %s\t Offset: x%.4x\t Name: %s\n", t, offset, ifd.TagName(t.ID))
}

func logIfdInfo(ifd ifds.Ifd, tagCount uint16, offset uint32) {
	InfoLogger.Printf("Ifd: %s\t Offset: x%.4x\t TagCount: %d\n", ifd, offset, tagCount)
}
