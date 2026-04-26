package panasonic

const PanasonicMakerNotePrefixLength = 12

// HasPanasonicHeader reports whether the maker-note payload starts with a
// Panasonic label prefix.
func HasPanasonicHeader(buf []byte) bool {
	return len(buf) >= PanasonicMakerNotePrefixLength &&
		buf[0] == 'P' &&
		buf[1] == 'a' &&
		buf[2] == 'n' &&
		buf[3] == 'a' &&
		buf[4] == 's' &&
		buf[5] == 'o' &&
		buf[6] == 'n' &&
		buf[7] == 'i' &&
		buf[8] == 'c' &&
		buf[9] == 0 &&
		buf[10] == 0 &&
		buf[11] == 0
}

// TIFFEPTags groups TIFF-EP extension tags commonly present in RAW containers.
type TIFFEPTags struct{}

// PanasonicRawTags groups Panasonic RW2/RWL specific root-IFD tags.
// These tags are not part of standard EXIF/TIFF IFD0 and are modeled separately.
type PanasonicRawTags struct {
	Version [4]byte // 0x0001 PanasonicRawVersion

	RawDataOffset    uint32 // 0x0118 RawDataOffset
	JpgFromRawOffset uint32 // 0x002e JpgFromRaw (offset only)
	JpgFromRawLength uint32 // 0x002e JpgFromRaw (length only)
	ISO              uint32 // 0x0017/0x0037 ISO
	// TODO: PanasonicTitle fields are UNDEFINED and may contain mixed encodings.
	// Keep parsed printable strings for parity with exiftool dumps.
	Title  string // 0xc6d2 PanasonicTitle
	Title2 string // 0xc6d3 PanasonicTitle2

	SensorWidth   uint16 // 0x0002 SensorWidth
	SensorHeight  uint16 // 0x0003 SensorHeight
	BitsPerSample uint16 // 0x000a BitsPerSample
	Compression   uint16 // 0x000b Compression
	RawFormat     uint16 // 0x002d RawFormat
	CropTop       uint16 // 0x0121 CropTop
	CropLeft      uint16 // 0x0122 CropLeft
	CropBottom    uint16 // 0x0123 CropBottom
	CropRight     uint16 // 0x0124 CropRight
}
