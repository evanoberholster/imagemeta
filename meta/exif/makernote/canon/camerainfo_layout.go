package canon

import "strings"

// CameraInfoLayout identifies the Canon CameraInfo payload layout.
type CameraInfoLayout uint8

const (
	CameraInfoLayoutUnknown CameraInfoLayout = iota
	CameraInfoLayout5D
	CameraInfoLayout5DmkII
	CameraInfoLayout5DmkIII
	CameraInfoLayout6D
	CameraInfoLayout7D
	CameraInfoLayout40D
	CameraInfoLayout50D
	CameraInfoLayout60D
	CameraInfoLayout70D
	CameraInfoLayout80D
	CameraInfoLayout450D
	CameraInfoLayout500D
	CameraInfoLayout550D
	CameraInfoLayout600D
	CameraInfoLayout650D
	CameraInfoLayout700D
	CameraInfoLayout750D
	CameraInfoLayout1000D
	CameraInfoLayoutPowerShot
	CameraInfoLayoutPowerShot2
	CameraInfoLayoutUnknown32
	CameraInfoLayoutR6
	CameraInfoLayoutR6m2
	CameraInfoLayoutR6m3
)

// CameraInfoSpec defines byte offsets for model-specific CameraInfo fields.
// An offset of -1 means the field is not present in that model's layout.
type CameraInfoSpec struct {
	FNumberOff               int
	ExposureTimeOff          int
	ISOOff                   int
	HighlightTonePriorityOff int
	FlashMeteringModeOff     int
	MeasuredEV2Off           int
	CameraTemperatureOff     int
	MacroMagnificationOff    int
	FocalLengthOff           int
	CameraOrientationOff     int
	WhiteBalanceOff          int
	ColorTemperatureOff      int
	LensTypeOff              int
	MinFocalLengthOff        int
	MaxFocalLengthOff        int
	JPEGQualityOff           int
	PictureStyleOff          int
	FirmwareVersionOff       int
	FirmwareVersionLen       int
	FileIndexOff             int
	DirectoryIndexOff        int
}

// Per-model CameraInfo payload layouts matching ExifTool Canon.pm offsets.
var (
	CameraInfoSpecLayout5D = CameraInfoSpec{
		FNumberOff:            0x03,
		ExposureTimeOff:       0x04,
		ISOOff:                0x06,
		CameraTemperatureOff:  0x17,
		MacroMagnificationOff: 0x1b,
		FocalLengthOff:        0x28,
		CameraOrientationOff:  0x27,
		WhiteBalanceOff:       0x54,
		ColorTemperatureOff:   0x58,
		PictureStyleOff:       0x6c,
		LensTypeOff:           0x97,
		MinFocalLengthOff:     0x93,
		MaxFocalLengthOff:     0x95,
		FirmwareVersionOff:    0xa4, FirmwareVersionLen: 8,
		FileIndexOff:      0xd0,
		DirectoryIndexOff: 0xcc,
	}

	CameraInfoSpecLayout5DmkII = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		MacroMagnificationOff:    0x1b,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x31,
		WhiteBalanceOff:          0x6f,
		ColorTemperatureOff:      0x73,
		PictureStyleOff:          0xa7,
		LensTypeOff:              0xe6,
		MinFocalLengthOff:        0xe8,
		MaxFocalLengthOff:        0xea,
		FirmwareVersionOff:       0x17e, FirmwareVersionLen: 6,
		FileIndexOff:      0x1bb,
		DirectoryIndexOff: 0x1c7,
	}

	CameraInfoSpecLayout5DmkIII = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x7d,
		WhiteBalanceOff:      0xbc,
		ColorTemperatureOff:  0xc0,
		PictureStyleOff:      0xf4,
		LensTypeOff:          0x153,
		MinFocalLengthOff:    0x155,
		MaxFocalLengthOff:    0x157,
		FirmwareVersionOff:   0x23c, FirmwareVersionLen: 6,
		FileIndexOff:      0x28c,
		DirectoryIndexOff: 0x298,
	}

	CameraInfoSpecLayout6D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x83,
		WhiteBalanceOff:      0xc2,
		ColorTemperatureOff:  0xc6,
		PictureStyleOff:      0xfa,
		LensTypeOff:          0x161,
		MinFocalLengthOff:    0x163,
		MaxFocalLengthOff:    0x165,
		FirmwareVersionOff:   0x256, FirmwareVersionLen: 6,
		FileIndexOff:      0x2aa,
		DirectoryIndexOff: 0x2b6,
	}

	CameraInfoSpecLayout7D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		MeasuredEV2Off:           0x08,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x35,
		WhiteBalanceOff:          0x77,
		ColorTemperatureOff:      0x7b,
		PictureStyleOff:          0xaf,
		LensTypeOff:              0x112,
		MinFocalLengthOff:        0x114,
		MaxFocalLengthOff:        0x116,
		FirmwareVersionOff:       0x1ac, FirmwareVersionLen: 6,
		FileIndexOff:      0x1eb,
		DirectoryIndexOff: 0x1f7,
	}

	CameraInfoSpecLayout40D = CameraInfoSpec{
		FNumberOff:            0x03,
		ExposureTimeOff:       0x04,
		ISOOff:                0x06,
		FlashMeteringModeOff:  0x15,
		CameraTemperatureOff:  0x18,
		MacroMagnificationOff: 0x1b,
		FocalLengthOff:        0x1d,
		CameraOrientationOff:  0x30,
		WhiteBalanceOff:       0x6f,
		ColorTemperatureOff:   0x73,
		LensTypeOff:           0xd6,
		MinFocalLengthOff:     0xd8,
		MaxFocalLengthOff:     0xda,
		FirmwareVersionOff:    0xff, FirmwareVersionLen: 6,
		FileIndexOff:      0x133,
		DirectoryIndexOff: 0x13f,
	}

	CameraInfoSpecLayout50D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x31,
		WhiteBalanceOff:          0x6f,
		ColorTemperatureOff:      0x73,
		PictureStyleOff:          0xa7,
		LensTypeOff:              0xea,
		MinFocalLengthOff:        0xec,
		MaxFocalLengthOff:        0xee,
		FirmwareVersionOff:       0x15e, FirmwareVersionLen: 6,
		FileIndexOff:      0x19b,
		DirectoryIndexOff: 0x1a7,
	}

	CameraInfoSpecLayout60D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x38,
		WhiteBalanceOff:          0x7b,
		ColorTemperatureOff:      0x7f,
		PictureStyleOff:          0xb3,
		LensTypeOff:              0xea,
		MinFocalLengthOff:        0xec,
		MaxFocalLengthOff:        0xee,
		FirmwareVersionOff:       0x19b, FirmwareVersionLen: 6,
		FileIndexOff:      0x1db,
		DirectoryIndexOff: 0x1e7,
	}

	CameraInfoSpecLayout70D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x84,
		ColorTemperatureOff:  0xc7,
		LensTypeOff:          0x166,
		MinFocalLengthOff:    0x168,
		MaxFocalLengthOff:    0x16a,
		FirmwareVersionOff:   0x25e, FirmwareVersionLen: 6,
		FileIndexOff:      0x2b3,
		DirectoryIndexOff: 0x2bf,
	}

	CameraInfoSpecLayout80D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x96,
		ColorTemperatureOff:  0x13a,
		LensTypeOff:          0x189,
		MinFocalLengthOff:    0x18b,
		MaxFocalLengthOff:    0x18d,
		FirmwareVersionOff:   0x45a, FirmwareVersionLen: 6,
		FileIndexOff:      0x4ae,
		DirectoryIndexOff: 0x4ba,
	}

	CameraInfoSpecLayout450D = CameraInfoSpec{
		FNumberOff:            0x03,
		ExposureTimeOff:       0x04,
		ISOOff:                0x06,
		FlashMeteringModeOff:  0x15,
		CameraTemperatureOff:  0x18,
		MacroMagnificationOff: 0x1b,
		FocalLengthOff:        0x1d,
		CameraOrientationOff:  0x30,
		WhiteBalanceOff:       0x6f,
		ColorTemperatureOff:   0x73,
		LensTypeOff:           0xde,
		FirmwareVersionOff:    0x107, FirmwareVersionLen: 6,
		FileIndexOff:      0x13f,
		DirectoryIndexOff: 0x133,
	}

	CameraInfoSpecLayout500D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x31,
		WhiteBalanceOff:          0x73,
		ColorTemperatureOff:      0x77,
		PictureStyleOff:          0xab,
		LensTypeOff:              0xf6,
		MinFocalLengthOff:        0xf8,
		MaxFocalLengthOff:        0xfa,
		FirmwareVersionOff:       0x190, FirmwareVersionLen: 6,
		FileIndexOff:      0x1d3,
		DirectoryIndexOff: 0x1df,
	}

	CameraInfoSpecLayout550D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x35,
		WhiteBalanceOff:          0x78,
		ColorTemperatureOff:      0x7c,
		PictureStyleOff:          0xb0,
		LensTypeOff:              0xff,
		MinFocalLengthOff:        0x101,
		MaxFocalLengthOff:        0x103,
		FirmwareVersionOff:       0x1a4, FirmwareVersionLen: 6,
		FileIndexOff:      0x1e4,
		DirectoryIndexOff: 0x1f0,
	}

	CameraInfoSpecLayout600D = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x38,
		WhiteBalanceOff:          0x7b,
		ColorTemperatureOff:      0x7f,
		PictureStyleOff:          0xb3,
		LensTypeOff:              0xea,
		MinFocalLengthOff:        0xec,
		MaxFocalLengthOff:        0xee,
		FirmwareVersionOff:       0x19b, FirmwareVersionLen: 6,
		FileIndexOff:      0x1db,
		DirectoryIndexOff: 0x1e7,
	}

	CameraInfoSpecLayout650D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x7d,
		WhiteBalanceOff:      0xbc,
		ColorTemperatureOff:  0xc0,
		PictureStyleOff:      0xf4,
		LensTypeOff:          0x127,
		MinFocalLengthOff:    0x129,
		MaxFocalLengthOff:    0x12b,
		FirmwareVersionOff:   0x21b, FirmwareVersionLen: 6,
		FileIndexOff:      0x270,
		DirectoryIndexOff: 0x27c,
	}

	CameraInfoSpecLayout700D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x7d,
		WhiteBalanceOff:      0xbc,
		ColorTemperatureOff:  0xc0,
		PictureStyleOff:      0xf4,
		LensTypeOff:          0x127,
		MinFocalLengthOff:    0x129,
		MaxFocalLengthOff:    0x12b,
		FirmwareVersionOff:   0x220, FirmwareVersionLen: 6,
		FileIndexOff:      0x274,
		DirectoryIndexOff: 0x280,
	}

	CameraInfoSpecLayout750D = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x96,
		WhiteBalanceOff:      0x131,
		ColorTemperatureOff:  0x135,
		PictureStyleOff:      0x169,
		LensTypeOff:          0x184,
		MinFocalLengthOff:    0x186,
		MaxFocalLengthOff:    0x188,
		FirmwareVersionOff:   0x43d, FirmwareVersionLen: 6,
	}

	CameraInfoSpecLayout1000D = CameraInfoSpec{
		FNumberOff:            0x03,
		ExposureTimeOff:       0x04,
		ISOOff:                0x06,
		FlashMeteringModeOff:  0x15,
		CameraTemperatureOff:  0x18,
		MacroMagnificationOff: 0x1b,
		FocalLengthOff:        0x1d,
		CameraOrientationOff:  0x30,
		WhiteBalanceOff:       0x6f,
		ColorTemperatureOff:   0x73,
		LensTypeOff:           0xe2,
		MinFocalLengthOff:     0xe4,
		MaxFocalLengthOff:     0xe6,
		FirmwareVersionOff:    0x10b, FirmwareVersionLen: 6,
		FileIndexOff:      0x143,
		DirectoryIndexOff: 0x137,
	}
)

// CameraInfoLayoutForModelName returns the CameraInfo layout for a camera
// model name string returned from the MakerNote CameraInfo payload.
func CameraInfoLayoutForModelName(model string) (CameraInfoLayout, bool) {
	switch {
	case strings.Contains(model, "Kiss X70"), strings.Contains(model, "Rebel T5"), strings.Contains(model, "1200D"):
		return CameraInfoLayout60D, true
	default:
		return CameraInfoLayoutUnknown, false
	}
}

// CameraInfoLayoutForModelID returns the CameraInfo layout for a Canon
// camera model identifier. The boolean is false when the model is unknown.
func CameraInfoLayoutForModelID(modelID CanonCameraModel) (CameraInfoLayout, bool) {
	switch modelID {
	case CanonModelEOS5D:
		return CameraInfoLayout5D, true
	case CanonModelEOS5DMarkII:
		return CameraInfoLayout5DmkII, true
	case CanonModelEOS5DMarkIII:
		return CameraInfoLayout5DmkIII, true
	case CanonModelEOS6D:
		return CameraInfoLayout6D, true
	case CanonModelEOS7D:
		return CameraInfoLayout7D, true
	case CanonModelEOS40D:
		return CameraInfoLayout40D, true
	case CanonModelEOS50D:
		return CameraInfoLayout50D, true
	case CanonModelEOS60D, CanonModelEOSRebelT5:
		return CameraInfoLayout60D, true
	case CanonModelEOS70D:
		return CameraInfoLayout70D, true
	case CanonModelEOS80D:
		return CameraInfoLayout80D, true
	case CanonModelEOSDigitalRebelXSi:
		return CameraInfoLayout450D, true
	case CanonModelEOSRebelT1i:
		return CameraInfoLayout500D, true
	case CanonModelEOSRebelT2i:
		return CameraInfoLayout550D, true
	case CanonModelEOSRebelT3i, CanonModelEOSRebelT3:
		return CameraInfoLayout600D, true
	case CanonModelEOSRebelT4i:
		return CameraInfoLayout650D, true
	case CanonModelEOSRebelT5i:
		return CameraInfoLayout700D, true
	case CanonModelEOSRebelT6i, CanonModelEOSRebelT6s:
		return CameraInfoLayout750D, true
	case CanonModelEOSRebelXS:
		return CameraInfoLayout1000D, true
	case CanonModelEOSR5, CanonModelEOSR6:
		return CameraInfoLayoutR6, true
	case CanonModelEOSR50, CanonModelEOSR6MarkII, CanonModelEOSR8:
		return CameraInfoLayoutR6m2, true
	case CanonModelEOSR6MarkIII:
		return CameraInfoLayoutR6m3, true
	default:
		return CameraInfoLayoutUnknown, false
	}
}
