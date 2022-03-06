package exif

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// Inital Log Level
var LogLevel = LogLevelError

const (
	LogLevelNone uint8 = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func checkLogLevel(l uint8) bool {
	return LogLevel >= l
}

func isError() bool {
	return checkLogLevel(LogLevelError)
}

func isWarn() bool {
	return checkLogLevel(LogLevelWarn)
}

func isInfo() bool {
	return checkLogLevel(LogLevelInfo)
}

func isDebug() bool {
	return checkLogLevel(LogLevelDebug)
}

func logTagInfo(ifd ifds.Ifd, t tag.Tag, offset uint32) {
	fmt.Printf("Tag: %s\t Offset: x%.4x\t Name: %s\n", t, offset, ifd.TagName(t.ID))
}

func logIfdError() {

}

func logIfdInfo(ifd ifds.Ifd, tagCount uint16, offset uint32) {
	fmt.Println(ifd, tagCount, offset)
}
