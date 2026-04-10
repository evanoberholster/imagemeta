package sony

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

// MakerNoteTag is a Sony maker-note tag ID from ExifTool's Sony Main table.
type MakerNoteTag tag.ID

// Selected Sony maker-note Main-table tag IDs.
const (
	Quality               MakerNoteTag = 0x0102 // Quality
	FlashExposureComp     MakerNoteTag = 0x0104 // FlashExposureComp
	Teleconverter         MakerNoteTag = 0x0105 // Teleconverter
	WhiteBalanceFineTune  MakerNoteTag = 0x0112 // WhiteBalanceFineTune
	CameraSettings        MakerNoteTag = 0x0114 // CameraSettings / CameraSettings3
	WhiteBalance          MakerNoteTag = 0x0115 // WhiteBalance
	Rating                MakerNoteTag = 0x2002 // Rating
	Contrast              MakerNoteTag = 0x2004 // Contrast
	Saturation            MakerNoteTag = 0x2005 // Saturation
	Sharpness             MakerNoteTag = 0x2006 // Sharpness
	Quality2              MakerNoteTag = 0x202e // Quality
	SonyModelID           MakerNoteTag = 0xb001 // SonyModelID
	CreativeStyle         MakerNoteTag = 0xb020 // CreativeStyle
	DynamicRangeOptimizer MakerNoteTag = 0xb025 // DynamicRangeOptimizer
	ImageStabilization    MakerNoteTag = 0xb026 // ImageStabilization
	LensType              MakerNoteTag = 0xb027 // LensType
	ColorMode             MakerNoteTag = 0xb029 // ColorMode
)
