package isobmff

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
