package isobmff

// itemType
type itemType uint8

// itemTypes
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

// itemType from Buffer. always should be 4 bytes.
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
