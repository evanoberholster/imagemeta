package panasonic

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

// Panasonic contains the selected Panasonic maker-note fields currently
// decoded by imagemeta.
//
// The field set mirrors the subset of ExifTool's
// Image::ExifTool::Panasonic::Main table that imagemeta parses today.
type Panasonic struct {
	ImageQuality         uint16
	FirmwareVersion      string
	WhiteBalance         uint16
	FocusMode            uint16
	AFAreaMode           [2]uint8
	ImageStabilization   uint16
	MacroMode            uint16
	ShootingMode         uint16
	Audio                uint16
	WhiteBalanceBias     float64
	FlashBias            float64
	InternalSerialNumber string
	PanasonicExifVersion string
	ColorEffect          uint16
	TimeSincePowerOn     float64
	BurstMode            uint16
	SequenceNumber       uint32
	ContrastMode         uint16
	NoiseReduction       uint16
	SelfTimer            uint16
	Rotation             uint16
	TravelDay            uint16
	BatteryLevel         uint16
	TextStamp            uint16
	PanasonicImageWidth  uint32
	PanasonicImageHeight uint32
	AFPointPosition      [2]tag.RationalU
	MakerNoteVersion     string
	SceneMode            uint16
	WBRedLevel           uint16
	WBGreenLevel         uint16
	WBBlueLevel          uint16
}
