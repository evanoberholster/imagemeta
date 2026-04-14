package makernote

import "github.com/evanoberholster/imagemeta/meta"

// Selected Apple maker-note tags.
//
// Reference: https://exiftool.org/TagNames/Apple.html
const (
	TagAppleMakerNoteVersion uint16 = 0x0001
	TagAppleRunTime          uint16 = 0x0003
	TagAppleAETarget         uint16 = 0x0005
	TagAppleAEAverage        uint16 = 0x0006
	TagAppleAFStable         uint16 = 0x0007
	TagAppleBurstUUID        uint16 = 0x000b
	TagAppleAEStable         uint16 = 0x000d
	TagAppleOISMode          uint16 = 0x000f
	TagAppleContentID        uint16 = 0x0011
	TagAppleImageCaptureType uint16 = 0x0014
	TagAppleImageUniqueID    uint16 = 0x0015
)

// Apple contains selected Apple maker-note fields.
type Apple struct {
	RunTime           string
	BurstUUID         meta.UUID
	ImageUniqueID     meta.UUID
	ContentIdentifier string

	MakerNoteVersion int32
	AETarget         int32
	AEAverage        int32
	OISMode          int32
	ImageCaptureType int32
	AEStable         bool
	AFStable         bool
}
