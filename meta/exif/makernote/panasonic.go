package makernote

const (
	TagPanasonicImageQuality         uint16 = 0x0001
	TagPanasonicFirmwareVersion      uint16 = 0x0002
	TagPanasonicWhiteBalance         uint16 = 0x0003
	TagPanasonicFocusMode            uint16 = 0x0007
	TagPanasonicAFAreaMode           uint16 = 0x000f
	TagPanasonicImageStabilization   uint16 = 0x001a
	TagPanasonicMacroMode            uint16 = 0x001c
	TagPanasonicShootingMode         uint16 = 0x001f
	TagPanasonicAudio                uint16 = 0x0020
	TagPanasonicWhiteBalanceBias     uint16 = 0x0023
	TagPanasonicFlashBias            uint16 = 0x0024
	TagPanasonicInternalSerialNumber uint16 = 0x0025
	TagPanasonicExifVersion          uint16 = 0x0026
	TagPanasonicColorEffect          uint16 = 0x0028
	TagPanasonicTimeSincePowerOn     uint16 = 0x0029
	TagPanasonicBurstMode            uint16 = 0x002a
	TagPanasonicSequenceNumber       uint16 = 0x002b
	TagPanasonicContrastMode         uint16 = 0x002c
	TagPanasonicNoiseReduction       uint16 = 0x002d
	TagPanasonicSelfTimer            uint16 = 0x002e
	TagPanasonicRotation             uint16 = 0x0030
	TagPanasonicTravelDay            uint16 = 0x0036
	TagPanasonicBatteryLevel         uint16 = 0x0038
	TagPanasonicTextStamp            uint16 = 0x003b
	TagPanasonicImageWidth           uint16 = 0x004b
	TagPanasonicImageHeight          uint16 = 0x004c
	TagPanasonicAFPointPosition      uint16 = 0x004d
	TagPanasonicMakerNoteVersion     uint16 = 0x8000
	TagPanasonicSceneMode            uint16 = 0x8001
	TagPanasonicWBRedLevel           uint16 = 0x8004
	TagPanasonicWBGreenLevel         uint16 = 0x8005
	TagPanasonicWBBlueLevel          uint16 = 0x8006
	TagPanasonicTextStamp2           uint16 = 0x8008
	TagPanasonicTextStamp3           uint16 = 0x8009
)

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
