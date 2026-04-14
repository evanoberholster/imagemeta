package nikon

import (
	"github.com/evanoberholster/imagemeta/meta/utils"
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
