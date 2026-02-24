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

var (
	itemTypeHvc1FourCC = fourCCFromString("hvc1")
	itemTypeExifFourCC = fourCCFromString("Exif")
	itemTypeAv01FourCC = fourCCFromString("av01")
	itemTypeGridFourCC = fourCCFromString("grid")
	itemTypeInfeFourCC = fourCCFromString("infe")
	itemTypeMimeFourCC = fourCCFromString("mime")
	itemTypeURIFourCC  = fourCCFromString("uri ")
)

// itemType from Buffer. always should be 4 bytes.
func itemTypeFromBuf(buf []byte) itemType {
	if len(buf) < 4 {
		return itemTypeUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case itemTypeHvc1FourCC:
		return itemTypeHvc1
	case itemTypeExifFourCC:
		return itemTypeExif
	case itemTypeAv01FourCC:
		return itemTypeAv01
	case itemTypeGridFourCC:
		return itemTypeGrid
	case itemTypeInfeFourCC:
		return itemTypeInfe
	case itemTypeMimeFourCC:
		return itemTypeMime
	case itemTypeURIFourCC:
		return itemTypeURI
	default:
		return itemTypeUnknown
	}
}
