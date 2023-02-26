package ifds

import (
	"github.com/evanoberholster/imagemeta/exif2/tag"
)

// TagSubIfdString returns the string representation of a tag.ID when found on a SubIfd
func TagSubIfdString(id tag.ID, it IfdType) string {
	switch it {
	case SubIfd2:
		switch id {
		case 0x0111:
			return "JpgFromRawStart"
		case 0x0117:
			return "JpgFromRawLength"
		}
	}
	switch id {
	case 0x0111:
		return "PreviewImageStart"
	case 0x0117:
		return "PreviewImageLength"
	}
	return TagString(id)
}
