package canon

import (
	"math"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// ByteAt reads a byte at off or returns 0.
func ByteAt(buf []byte, off int) uint8 {
	if off < 0 || off >= len(buf) {
		return 0
	}
	return buf[off]
}

// U16LEAt reads a uint16 LE at off or returns 0.
func U16LEAt(buf []byte, off int) uint16 {
	if off < 0 || off+1 >= len(buf) {
		return 0
	}
	return uint16(buf[off]) | uint16(buf[off+1])<<8
}

// U16BEAt reads a uint16 BE at off or returns 0.
func U16BEAt(buf []byte, off int) uint16 {
	if off < 0 || off+1 >= len(buf) {
		return 0
	}
	return uint16(buf[off])<<8 | uint16(buf[off+1])
}

// U32LEAt reads a uint32 LE at off or returns 0.
func U32LEAt(buf []byte, off int) uint32 {
	if off < 0 || off+3 >= len(buf) {
		return 0
	}
	return uint32(buf[off]) | uint32(buf[off+1])<<8 |
		uint32(buf[off+2])<<16 | uint32(buf[off+3])<<24
}

// CIFNumber converts a camera-info FNumber byte to an Aperture.
func CIFNumber(v uint8) meta.Aperture {
	if v == 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(float64(int(v)-8) / 16.0))
}

// CIExposureTime converts a camera-info ExposureTime byte to an ExposureTime.
func CIExposureTime(v uint8) meta.ExposureTime {
	if v == 0 {
		return 0
	}
	return meta.ExposureTime(math.Exp2(4.0 * (1.0 - canonEV(float64(int16(v)-24)))))
}

// CIISO converts a camera-info ISO byte to a uint32 ISO value.
func CIISO(v uint8) uint32 {
	if v == 0 {
		return 0
	}
	return uint32(math.Round(100.0 * math.Exp2(float64(v)/8.0-9.0)))
}

// CITemperature converts a camera-info temperature byte to Celsius int16.
func CITemperature(v uint8) int16 {
	if v == 0 {
		return 0
	}
	return int16(v) - 128
}

// CIMeasuredEV2 converts a camera-info MeasuredEV2 byte.
func CIMeasuredEV2(v uint8) float32 {
	if v == 0 {
		return 0
	}
	return float32(v)/8.0 - 6.0
}

// CIMacroMagnification converts a camera-info MacroMagnification byte.
func CIMacroMagnification(v uint8) float32 {
	if v == 0 {
		return 0
	}
	return float32(math.Exp((75.0 - float64(v)) * math.Ln2 * 3.0 / 40.0))
}

// CIFocalLength converts a camera-info focal-length uint16 to a FocalLength.
func CIFocalLength(v uint16) meta.FocalLength {
	return meta.FocalLength(float32(v))
}

// CIAsciiBytes extracts a NUL-terminated ASCII string from buf at off with max length.
func CIAsciiBytes(buf []byte, off, maxLen int) string {
	if off < 0 || maxLen <= 0 || off >= len(buf) {
		return ""
	}
	end := off + maxLen
	if end > len(buf) {
		end = len(buf)
	}
	slice := buf[off:end]
	for i, b := range slice {
		if b == 0 {
			return string(slice[:i])
		}
	}
	return string(slice)
}

func canonEV(val float64) float64 {
	if val >= 0 {
		return val
	}
	return -((-val) * 2.0 / 3.0)
}

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
	CameraInfoLayout1D
	CameraInfoLayout1DmkII
	CameraInfoLayout1DmkIIN
	CameraInfoLayout1DmkIII
	CameraInfoLayout1DmkIV
	CameraInfoLayout1DX
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

	CameraInfoSpecLayout1D = CameraInfoSpec{
		ExposureTimeOff:     0x04,
		FocalLengthOff:      0x0a,
		LensTypeOff:         0x0d,
		MinFocalLengthOff:   0x0e,
		MaxFocalLengthOff:   0x10,
		WhiteBalanceOff:     0x44,
		ColorTemperatureOff: 0x48,
		PictureStyleOff:     0x4b,
	}

	CameraInfoSpecLayout1DmkII = CameraInfoSpec{
		ExposureTimeOff:     0x04,
		FocalLengthOff:      0x09,
		LensTypeOff:         0x0c,
		MinFocalLengthOff:   0x11,
		MaxFocalLengthOff:   0x13,
		WhiteBalanceOff:     0x36,
		ColorTemperatureOff: 0x37,
		JPEGQualityOff:      0x66,
		PictureStyleOff:     0x6c,
		ISOOff:              0x75,
	}

	CameraInfoSpecLayout1DmkIIN = CameraInfoSpec{
		ExposureTimeOff:     0x04,
		FocalLengthOff:      0x09,
		LensTypeOff:         0x0c,
		MinFocalLengthOff:   0x11,
		MaxFocalLengthOff:   0x13,
		WhiteBalanceOff:     0x36,
		ColorTemperatureOff: 0x37,
		JPEGQualityOff:      0x72,
		PictureStyleOff:     0x73,
		ISOOff:              0x75,
	}

	CameraInfoSpecLayout1DmkIII = CameraInfoSpec{
		FNumberOff:            0x03,
		ExposureTimeOff:       0x04,
		ISOOff:                0x06,
		CameraTemperatureOff:  0x18,
		MacroMagnificationOff: 0x1b,
		FocalLengthOff:        0x1d,
		CameraOrientationOff:  0x30,
		WhiteBalanceOff:       0x5e,
		ColorTemperatureOff:   0x62,
		PictureStyleOff:       0x86,
		LensTypeOff:           0x111,
		MinFocalLengthOff:     0x113,
		MaxFocalLengthOff:     0x115,
		FirmwareVersionOff:    0x136, FirmwareVersionLen: 6,
		FileIndexOff:      0x172,
		DirectoryIndexOff: 0x17e,
	}

	CameraInfoSpecLayout1DmkIV = CameraInfoSpec{
		FNumberOff:               0x03,
		ExposureTimeOff:          0x04,
		ISOOff:                   0x06,
		HighlightTonePriorityOff: 0x07,
		MeasuredEV2Off:           0x09,
		FlashMeteringModeOff:     0x15,
		CameraTemperatureOff:     0x19,
		FocalLengthOff:           0x1e,
		CameraOrientationOff:     0x35,
		WhiteBalanceOff:          0x78,
		ColorTemperatureOff:      0x7c,
		LensTypeOff:              0x14f,
		MinFocalLengthOff:        0x151,
		MaxFocalLengthOff:        0x153,
		FirmwareVersionOff:       0x1ed, FirmwareVersionLen: 6,
		FileIndexOff:      0x22c,
		DirectoryIndexOff: 0x238,
	}

	CameraInfoSpecLayout1DX = CameraInfoSpec{
		FNumberOff:           0x03,
		ExposureTimeOff:      0x04,
		ISOOff:               0x06,
		CameraTemperatureOff: 0x1b,
		FocalLengthOff:       0x23,
		CameraOrientationOff: 0x7d,
		WhiteBalanceOff:      0xbc,
		ColorTemperatureOff:  0xc0,
		PictureStyleOff:      0xf4,
		LensTypeOff:          0x1a7,
		MinFocalLengthOff:    0x1a9,
		MaxFocalLengthOff:    0x1ab,
		FirmwareVersionOff:   0x280, FirmwareVersionLen: 6,
		FileIndexOff:      0x2d0,
		DirectoryIndexOff: 0x2dc,
	}
)

// CameraInfoLayoutForModelName returns the CameraInfo layout for a camera
// model name string returned from the MakerNote CameraInfo payload.
//
// Deprecated: use CameraInfoLayoutForModelID. Model name matching is no longer
// necessary — all affected models have unique ModelID values.
func CameraInfoLayoutForModelName(model string) (CameraInfoLayout, bool) {
	return CameraInfoLayoutUnknown, false
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
	case CanonModelEOS1D:
		return CameraInfoLayout1D, true
	case CanonModelEOS1DS:
		return CameraInfoLayout1D, true
	case CanonModelEOS1DMarkII, CanonModelEOS1DsMarkII:
		return CameraInfoLayout1DmkII, true
	case CanonModelEOS1DMarkIIN:
		return CameraInfoLayout1DmkIIN, true
	case CanonModelEOS1DMarkIII, CanonModelEOS1DsMarkIII:
		return CameraInfoLayout1DmkIII, true
	case CanonModelEOS1DMarkIV:
		return CameraInfoLayout1DmkIV, true
	case CanonModelEOS1DX, CanonModelEOS1DC, CanonModelEOS1DXMarkII, CanonModelEOS1DXMarkIII:
		return CameraInfoLayout1DX, true
	default:
		return CameraInfoLayoutUnknown, false
	}
}

// CameraInfoDecode parses a Canon CameraInfo byte payload using model-specific
// byte offsets. The spec should be obtained from CameraInfoSpecLayout* variables.
func CameraInfoDecode(buf []byte, spec CameraInfoSpec) CameraInfo {
	if len(buf) == 0 {
		return CameraInfo{}
	}
	dst := CameraInfo{
		FNumber:               CIFNumber(ByteAt(buf, spec.FNumberOff)),
		ExposureTime:          CIExposureTime(ByteAt(buf, spec.ExposureTimeOff)),
		ISO:                   CIISO(ByteAt(buf, spec.ISOOff)),
		HighlightTonePriority: int16(ByteAt(buf, spec.HighlightTonePriorityOff)),
		FlashMeteringMode:     int16(ByteAt(buf, spec.FlashMeteringModeOff)),
		MeasuredEV2:           CIMeasuredEV2(ByteAt(buf, spec.MeasuredEV2Off)),
		CameraTemperature:     CITemperature(ByteAt(buf, spec.CameraTemperatureOff)),
		MacroMagnification:    CIMacroMagnification(ByteAt(buf, spec.MacroMagnificationOff)),
		FocalLength:           CIFocalLength(U16BEAt(buf, spec.FocalLengthOff)),
		WhiteBalance:          NewWhiteBalanceFromRaw(U16LEAt(buf, spec.WhiteBalanceOff)),
		ColorTemperature:      U16LEAt(buf, spec.ColorTemperatureOff),
		LensType:              CanonLensType(U16BEAt(buf, spec.LensTypeOff)),
		MinFocalLength:        CIFocalLength(U16BEAt(buf, spec.MinFocalLengthOff)),
		MaxFocalLength:        CIFocalLength(U16BEAt(buf, spec.MaxFocalLengthOff)),
		JPEGQuality:           ByteAt(buf, spec.JPEGQualityOff),
		PictureStyle:          int16(ByteAt(buf, spec.PictureStyleOff)),
		FirmwareVersion:       CIAsciiBytes(buf, spec.FirmwareVersionOff, spec.FirmwareVersionLen),
		CameraOrientation:     ByteAt(buf, spec.CameraOrientationOff),
	}

	// 5DmkII alternate orientation fallback (ExifTool Canon.pm).
	if dst.CameraOrientation == 0 && spec.CameraOrientationOff == 0x36 {
		if alt := ByteAt(buf, 0x3a); alt != 0 {
			dst.CameraOrientation = alt
		}
	}

	if v := U32LEAt(buf, spec.FileIndexOff); v > 0 {
		dst.FileIndex = v + 1
	}
	if v := U32LEAt(buf, spec.DirectoryIndexOff); v > 0 {
		dst.DirectoryIndex = v - 1
	}
	return dst
}

// BatteryTypePayloadSize is the fixed MakerNote BatteryType tag payload length.
const BatteryTypePayloadSize = 76

// BatteryTypeHeaderLen is the number of leading bytes to skip.
const BatteryTypeHeaderLen = 4

// ParseBatteryType extracts the NUL-terminated battery model string from a
// raw 72-byte CameraMaker:BatteryType payload (tag 0x0038).
//
// string(payload) does not allocate in modern Go when used only for switch
// comparison and the source is not subsequently modified.
func ParseBatteryType(payload []byte) string {
	// Find first NUL byte. The 72-byte payload contains the model string
	// followed by NUL padding and a trailing 0x01 marker.
	i := 0
	for i < len(payload) && payload[i] != 0 {
		i++
	}
	if i == 0 {
		return ""
	}

	switch string(payload[:i]) {
	case "LP-E6":
		return "LP-E6"
	case "LP-E6N":
		return "LP-E6N"
	case "LP-E6NH":
		return "LP-E6NH"
	case "LP-E6P":
		return "LP-E6P"
	case "LP-E12":
		return "LP-E12"
	case "LP-E17":
		return "LP-E17"
	case "LP-E19":
		return "LP-E19"
	case "NB-13L":
		return "NB-13L"
	default:
		return string(payload[:i])
	}
}

// AFWordsBuffer returns an AF info word buffer capped at AFWordsMax.
func AFWordsBuffer(stack []uint16, unitCount uint32) ([]uint16, bool) {
	if unitCount == 0 {
		return stack[:0], false
	}
	wordCount := int(unitCount)
	truncated := false
	if unitCount > AFWordsMax {
		wordCount = AFWordsMax
		truncated = true
	}
	if wordCount <= len(stack) {
		return stack[:wordCount], truncated
	}
	return make([]uint16, wordCount), truncated
}

const AFWordsMax = 8192

// AFInfoSourceFromID maps a tag ID to an AF info source.
func AFInfoSourceFromID(id tag.ID) AFInfoSource {
	switch MakerNoteTag(id) {
	case CanonAFInfo:
		return AFInfoSourceAFInfo
	case CanonAFInfo2:
		return AFInfoSourceAFInfo2
	case AFInfo3:
		return AFInfoSourceAFInfo3
	default:
		return AFInfoSourceUnknown
	}
}

// ShouldReplaceAFInfo decides whether candidate should replace current AF info.
func ShouldReplaceAFInfo(current, candidate AFInfo) bool {
	curHas := afInfoHasData(current)
	candHas := afInfoHasData(candidate)
	switch {
	case candHas && !curHas:
		return true
	case !candHas && curHas:
		return false
	case !candHas && !curHas:
		return afInfoSourcePriority(candidate.Source) > afInfoSourcePriority(current.Source)
	}

	curScore := afInfoQualityScore(current)
	candScore := afInfoQualityScore(candidate)
	if candScore != curScore {
		return candScore > curScore
	}
	return afInfoSourcePriority(candidate.Source) > afInfoSourcePriority(current.Source)
}

func afInfoHasData(v AFInfo) bool {
	return v.NumAFPoints != 0 ||
		v.ValidAFPoints != 0 ||
		v.CanonImageWidth != 0 ||
		v.CanonImageHeight != 0 ||
		len(v.AFArea) != 0 ||
		len(v.AFPointsInFocusBits) != 0 ||
		len(v.AFPointsSelectedBits) != 0 ||
		v.PrimaryAFPoint != 0
}

func afInfoQualityScore(v AFInfo) int {
	score := int(v.NumAFPoints) + int(v.ValidAFPoints)
	score += len(v.AFArea)
	score += len(v.AFPoints)
	score += len(v.AFPointsInFocusBits)
	score += len(v.AFPointsSelectedBits)
	if v.CanonImageWidth != 0 && v.CanonImageHeight != 0 {
		score += 8
	}
	if v.AFImageWidth != 0 && v.AFImageHeight != 0 {
		score += 8
	}
	if v.AFAreaWidth != 0 || v.AFAreaHeight != 0 {
		score += 4
	}
	return score
}

func afInfoSourcePriority(source AFInfoSource) int {
	switch source {
	case AFInfoSourceAFInfo2:
		return 3
	case AFInfoSourceAFInfo3:
		return 2
	case AFInfoSourceAFInfo:
		return 1
	default:
		return 0
	}
}

// MaxAFPoints returns the known maximum AF points for the camera model.
// Returns 0 for unknown models.
func MaxAFPoints(modelID CanonCameraModel) int {
	switch modelID {
	case CanonModelEOS1DX, CanonModelEOS1DC, CanonModelEOS1DXMarkII, CanonModelEOS1DXMarkIII:
		return 191
	case CanonModelEOSR3:
		return 1053
	case CanonModelEOSR5, CanonModelEOSR5MarkII, CanonModelEOSR6, CanonModelEOSR6MarkII,
		CanonModelEOSR6MarkIII, CanonModelEOSR7, CanonModelEOSR8, CanonModelEOSR10,
		CanonModelEOSR50, CanonModelEOSR, CanonModelEOSRP, CanonModelEOSR1,
		CanonModelEOSR100, CanonModelEOSR50V:
		return 1053
	case CanonModelEOS5DMarkIV, CanonModelEOS5DS, CanonModelEOS5DSR:
		return 61
	case CanonModelEOS5DMarkIII, CanonModelEOS6DMarkII, CanonModelEOS7DMarkII:
		return 65
	case CanonModelEOS7D, CanonModelEOS70D, CanonModelEOS80D, CanonModelEOS90D:
		return 19
	case CanonModelEOS60D, CanonModelEOSRebelT3i, CanonModelEOSRebelT4i, CanonModelEOSRebelT5i:
		return 9
	case CanonModelEOS6D, CanonModelEOS5DMarkII:
		return 11
	case CanonModelEOS1DMarkIII, CanonModelEOS1DsMarkIII, CanonModelEOS1DMarkIV:
		return 45
	}
	return 45
}
