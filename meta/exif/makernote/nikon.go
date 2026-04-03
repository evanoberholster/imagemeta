package makernote

import "github.com/evanoberholster/imagemeta/meta/utils"

// Selected Nikon maker-note tags (https://exiftool.org/TagNames/Nikon.html).
//
// Intentionally excluded:
// - LensType value decoding
const (
	TagNikonVersion      uint16 = 0x0001
	TagNikonISOSetting   uint16 = 0x0002
	TagNikonColorMode    uint16 = 0x0003
	TagNikonQuality      uint16 = 0x0004
	TagNikonWhiteBalance uint16 = 0x0005
	TagNikonSharpness    uint16 = 0x0006
	TagNikonFocusMode    uint16 = 0x0007
	TagNikonFlashSetting uint16 = 0x0008
	TagNikonFlashType    uint16 = 0x0009
	TagNikonISOSelection uint16 = 0x000f
	TagNikonSerialNumber uint16 = 0x001d
	TagNikonLens         uint16 = 0x0084
)

const nikonHeaderLength = 18

// HasNikonHeader reports whether the maker-note payload starts with a Nikon label.
func HasNikonHeader(buf []byte) bool {
	return len(buf) >= 5 && string(buf[:5]) == "Nikon"
}

// ParseNikonHeader parses Nikon maker-note TIFF info from the leading header.
func ParseNikonHeader(buf []byte) (bo utils.ByteOrder, ifdRelOffset uint32, ok bool) {
	if len(buf) < nikonHeaderLength || !HasNikonHeader(buf) {
		return utils.UnknownEndian, 0, false
	}
	bo = utils.BinaryOrder(buf[10:14])
	if bo == utils.UnknownEndian {
		return utils.UnknownEndian, 0, false
	}
	return bo, bo.Uint32(buf[14:18]), true
}
