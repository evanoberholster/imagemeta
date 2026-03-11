package makernote

// Info contains parsed maker-note values for supported vendors.
type Info struct {
	Make  CameraMake
	Apple Apple
	Canon Canon
	Nikon Nikon

	parsedTagCount uint8
	parsedTagIDs   [128]uint16
}

// HasTagParsed reports whether a maker-note tag ID was parsed.
func (i Info) HasTagParsed(tagID uint16) bool {
	n := int(i.parsedTagCount)
	if n > len(i.parsedTagIDs) {
		n = len(i.parsedTagIDs)
	}
	for idx := 0; idx < n; idx++ {
		if i.parsedTagIDs[idx] == tagID {
			return true
		}
	}
	return false
}

// MarkTagParsed records a maker-note tag ID as parsed.
func (i *Info) MarkTagParsed(tagID uint16) {
	n := int(i.parsedTagCount)
	if n > len(i.parsedTagIDs) {
		n = len(i.parsedTagIDs)
	}
	for idx := 0; idx < n; idx++ {
		if i.parsedTagIDs[idx] == tagID {
			return
		}
	}
	if n >= len(i.parsedTagIDs) {
		return
	}
	i.parsedTagIDs[n] = tagID
	i.parsedTagCount++
}

// MergeParsedTags merges parsed-tag markers from src into i.
func (i *Info) MergeParsedTags(src Info) {
	n := int(src.parsedTagCount)
	if n > len(src.parsedTagIDs) {
		n = len(src.parsedTagIDs)
	}
	for idx := 0; idx < n; idx++ {
		i.MarkTagParsed(src.parsedTagIDs[idx])
	}
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
