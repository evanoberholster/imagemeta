package makernote

const makerNoteBitsetMaxTagID uint16 = (8 * 64) - 1

// Makernote is the generic maker-note payload container type.
type Makernote any

// Info contains parsed maker-note values for supported vendors.
type Info struct {
	Make  CameraMake
	Apple Apple
	Canon Canon
	Nikon Nikon

	ifdBitset [8]uint64
}

// HasTagParsed reports whether a maker-note tag ID has been parsed.
func (m Info) HasTagParsed(tagID uint16) bool {
	if tagID > makerNoteBitsetMaxTagID {
		return false
	}
	word := tagID >> 6
	mask := uint64(1) << (tagID & 63)
	return (m.ifdBitset[word] & mask) != 0
}

// TagParsedBitset returns the parsed-tag bitset for maker notes.
func (m Info) TagParsedBitset() [8]uint64 {
	return m.ifdBitset
}

// MarkTagParsed marks a maker-note tag ID as parsed.
func (m *Info) MarkTagParsed(tagID uint16) {
	if tagID > makerNoteBitsetMaxTagID {
		return
	}
	word := tagID >> 6
	m.ifdBitset[word] |= uint64(1) << (tagID & 63)
}

// Nikon contains selected Nikon maker-note fields.
type Nikon struct {
	Version [8]byte

	Quality      string
	ColorMode    string
	WhiteBalance string
	Sharpness    string
	FocusMode    string
	FlashSetting string
	FlashType    string
	ISOSelection string
	SerialNumber string
	Lens         string

	VersionCount uint8
	ISOSetting   uint32
}

// Apple contains selected Apple maker-note fields.
type Apple struct {
	RunTime           string
	BurstUUID         string
	ContentIdentifier string
	ImageUniqueID     string

	MakerNoteVersion int32
	AETarget         int32
	AEAverage        int32
	OISMode          int32
	ImageCaptureType int32
	AEStable         bool
	AFStable         bool
}
