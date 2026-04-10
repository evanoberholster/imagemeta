package panasonic

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

// MakerNoteTag is a Panasonic maker-note tag ID from ExifTool's Panasonic Main
// table.
type MakerNoteTag tag.ID

// Selected Panasonic maker-note Main-table tag IDs.
const (
	ImageQuality         MakerNoteTag = 0x0001 // ImageQuality
	FirmwareVersion      MakerNoteTag = 0x0002 // FirmwareVersion
	WhiteBalance         MakerNoteTag = 0x0003 // WhiteBalance
	FocusMode            MakerNoteTag = 0x0007 // FocusMode
	AFAreaMode           MakerNoteTag = 0x000f // AFAreaMode
	ImageStabilization   MakerNoteTag = 0x001a // ImageStabilization
	MacroMode            MakerNoteTag = 0x001c // MacroMode
	ShootingMode         MakerNoteTag = 0x001f // ShootingMode
	Audio                MakerNoteTag = 0x0020 // Audio
	WhiteBalanceBias     MakerNoteTag = 0x0023 // WhiteBalanceBias
	FlashBias            MakerNoteTag = 0x0024 // FlashBias
	InternalSerialNumber MakerNoteTag = 0x0025 // InternalSerialNumber
	PanasonicExifVersion MakerNoteTag = 0x0026 // PanasonicExifVersion
	ColorEffect          MakerNoteTag = 0x0028 // ColorEffect
	TimeSincePowerOn     MakerNoteTag = 0x0029 // TimeSincePowerOn
	BurstMode            MakerNoteTag = 0x002a // BurstMode
	SequenceNumber       MakerNoteTag = 0x002b // SequenceNumber
	ContrastMode         MakerNoteTag = 0x002c // ContrastMode
	NoiseReduction       MakerNoteTag = 0x002d // NoiseReduction
	SelfTimer            MakerNoteTag = 0x002e // SelfTimer
	Rotation             MakerNoteTag = 0x0030 // Rotation
	TravelDay            MakerNoteTag = 0x0036 // TravelDay
	BatteryLevel         MakerNoteTag = 0x0038 // BatteryLevel
	TextStamp            MakerNoteTag = 0x003b // TextStamp
	PanasonicImageWidth  MakerNoteTag = 0x004b // PanasonicImageWidth
	PanasonicImageHeight MakerNoteTag = 0x004c // PanasonicImageHeight
	AFPointPosition      MakerNoteTag = 0x004d // AFPointPosition
	MakerNoteVersion     MakerNoteTag = 0x8000 // MakerNoteVersion
	SceneMode            MakerNoteTag = 0x8001 // SceneMode
	WBRedLevel           MakerNoteTag = 0x8004 // WBRedLevel
	WBGreenLevel         MakerNoteTag = 0x8005 // WBGreenLevel
	WBBlueLevel          MakerNoteTag = 0x8006 // WBBlueLevel
	TextStamp2           MakerNoteTag = 0x8008 // TextStamp
	TextStamp3           MakerNoteTag = 0x8009 // TextStamp
)
