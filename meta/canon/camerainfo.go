package canon

import (
	"encoding/binary"
	"fmt"
	"math"
)

//go:generate msgp

type CanonCameraInfo interface {
}

// CameraInfo6D represents camera information specific to the EOS 6D model
type CameraInfo6D struct {
	FNumber            FNumber           `json:"fNumber"`            // 0x03
	ExposureTime       CIExposureTime    `json:"exposureTime"`       // 0x04
	ISO                ISO               `json:"iso"`                // 0x06
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`  // 0x1b
	FocalLength        FocalLength       `json:"focalLength"`        // 0x23
	CameraOrientation  CameraOrientation `json:"orientation"`        // 0x83
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"` // 0x92
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"` // 0x94
	WhiteBalance       uint16            `json:"whiteBalance"`       // 0xc2
	ColorTemperature   uint16            `json:"colorTemperature"`   // 0xc6
	PictureStyle       uint8             `json:"pictureStyle"`       // 0xfa
	LensType           CanonLensType     `json:"lensType"`           // 0x161. Value is big-endian
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`     // 0x163
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`     // 0x165
	FirmwareVersion    string            `json:"firmwareVersion"`    // 0x256
	FileIndex          uint32            `json:"fileIndex"`          // 0x2aa
	DirectoryIndex     uint32            `json:"directoryIndex"`     // 0x2b6
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`   // 0x3c6
}

// CameraOrientation represents the orientation of the camera
type CameraOrientation uint8

const (
	OrientationHorizontal  CameraOrientation = 0
	OrientationRotate90CW  CameraOrientation = 1
	OrientationRotate270CW CameraOrientation = 2
)

// String returns human readable orientation
func (o CameraOrientation) OrientationString() string {
	switch o {
	case OrientationHorizontal:
		return "Horizontal (normal)"
	case OrientationRotate90CW:
		return "Rotate 90 CW"
	case OrientationRotate270CW:
		return "Rotate 270 CW"
	default:
		return "Unknown"
	}
}

// PictureStyleInfo2 represents custom picture style information for EOS models
type PictureStyleInfo2 struct {
	// Standard settings
	ContrastStandard     int32 `json:"contrastStandard"`     // 0x00
	SharpnessStandard    int32 `json:"sharpnessStandard"`    // 0x04
	SaturationStandard   int32 `json:"saturationStandard"`   // 0x08
	ColorToneStandard    int32 `json:"colorToneStandard"`    // 0x0c
	FilterEffectStandard int32 `json:"filterEffectStandard"` // 0x10
	ToningEffectStandard int32 `json:"toningEffectStandard"` // 0x14

	// Portrait settings
	ContrastPortrait     int32 `json:"contrastPortrait"`     // 0x18
	SharpnessPortrait    int32 `json:"sharpnessPortrait"`    // 0x1c
	SaturationPortrait   int32 `json:"saturationPortrait"`   // 0x20
	ColorTonePortrait    int32 `json:"colorTonePortrait"`    // 0x24
	FilterEffectPortrait int32 `json:"filterEffectPortrait"` // 0x28
	ToningEffectPortrait int32 `json:"toningEffectPortrait"` // 0x2c

	// Landscape settings
	ContrastLandscape     int32 `json:"contrastLandscape"`     // 0x30
	SharpnessLandscape    int32 `json:"sharpnessLandscape"`    // 0x34
	SaturationLandscape   int32 `json:"saturationLandscape"`   // 0x38
	ColorToneLandscape    int32 `json:"colorToneLandscape"`    // 0x3c
	FilterEffectLandscape int32 `json:"filterEffectLandscape"` // 0x40
	ToningEffectLandscape int32 `json:"toningEffectLandscape"` // 0x44

	// Neutral settings
	ContrastNeutral     int32 `json:"contrastNeutral"`     // 0x48
	SharpnessNeutral    int32 `json:"sharpnessNeutral"`    // 0x4c
	SaturationNeutral   int32 `json:"saturationNeutral"`   // 0x50
	ColorToneNeutral    int32 `json:"colorToneNeutral"`    // 0x54
	FilterEffectNeutral int32 `json:"filterEffectNeutral"` // 0x58
	ToningEffectNeutral int32 `json:"toningEffectNeutral"` // 0x5c

	// Faithful settings
	ContrastFaithful     int32 `json:"contrastFaithful"`     // 0x60
	SharpnessFaithful    int32 `json:"sharpnessFaithful"`    // 0x64
	SaturationFaithful   int32 `json:"saturationFaithful"`   // 0x68
	ColorToneFaithful    int32 `json:"colorToneFaithful"`    // 0x6c
	FilterEffectFaithful int32 `json:"filterEffectFaithful"` // 0x70
	ToningEffectFaithful int32 `json:"toningEffectFaithful"` // 0x74

	// Monochrome settings
	ContrastMonochrome     int32        `json:"contrastMonochrome"`     // 0x78
	SharpnessMonochrome    int32        `json:"sharpnessMonochrome"`    // 0x7c
	SaturationMonochrome   int32        `json:"saturationMonochrome"`   // 0x80
	ColorToneMonochrome    int32        `json:"colorToneMonochrome"`    // 0x84
	FilterEffectMonochrome FilterEffect `json:"filterEffectMonochrome"` // 0x88
	ToningEffectMonochrome ToningEffect `json:"toningEffectMonochrome"` // 0x8c

	// Auto settings
	ContrastAuto     int32        `json:"contrastAuto"`     // 0x90
	SharpnessAuto    int32        `json:"sharpnessAuto"`    // 0x94
	SaturationAuto   int32        `json:"saturationAuto"`   // 0x98
	ColorToneAuto    int32        `json:"colorToneAuto"`    // 0x9c
	FilterEffectAuto FilterEffect `json:"filterEffectAuto"` // 0xa0
	ToningEffectAuto ToningEffect `json:"toningEffectAuto"` // 0xa4

	// User Defined 1
	ContrastUserDef1     int32        `json:"contrastUserDef1"`     // 0xa8
	SharpnessUserDef1    int32        `json:"sharpnessUserDef1"`    // 0xac
	SaturationUserDef1   int32        `json:"saturationUserDef1"`   // 0xb0
	ColorToneUserDef1    int32        `json:"colorToneUserDef1"`    // 0xb4
	FilterEffectUserDef1 FilterEffect `json:"filterEffectUserDef1"` // 0xb8
	ToningEffectUserDef1 ToningEffect `json:"toningEffectUserDef1"` // 0xbc

	// User Defined 2
	ContrastUserDef2     int32        `json:"contrastUserDef2"`     // 0xc0
	SharpnessUserDef2    int32        `json:"sharpnessUserDef2"`    // 0xc4
	SaturationUserDef2   int32        `json:"saturationUserDef2"`   // 0xc8
	ColorToneUserDef2    int32        `json:"colorToneUserDef2"`    // 0xcc
	FilterEffectUserDef2 FilterEffect `json:"filterEffectUserDef2"` // 0xd0
	ToningEffectUserDef2 ToningEffect `json:"toningEffectUserDef2"` // 0xd4

	// User Defined 3
	ContrastUserDef3     int32        `json:"contrastUserDef3"`     // 0xd8
	SharpnessUserDef3    int32        `json:"sharpnessUserDef3"`    // 0xdc
	SaturationUserDef3   int32        `json:"saturationUserDef3"`   // 0xe0
	ColorToneUserDef3    int32        `json:"colorToneUserDef3"`    // 0xe4
	FilterEffectUserDef3 FilterEffect `json:"filterEffectUserDef3"` // 0xe8
	ToningEffectUserDef3 ToningEffect `json:"toningEffectUserDef3"` // 0xec

	// Base picture styles
	UserDef1PictureStyle uint16 `json:"userDef1PictureStyle"` // 0xf0
	UserDef2PictureStyle uint16 `json:"userDef2PictureStyle"` // 0xf2
	UserDef3PictureStyle uint16 `json:"userDef3PictureStyle"` // 0xf4
}

// FilterEffect represents monochrome filter effects
type FilterEffect int32

const (
	FilterEffectNone   FilterEffect = 0
	FilterEffectYellow FilterEffect = 1
	FilterEffectOrange FilterEffect = 2
	FilterEffectRed    FilterEffect = 3
	FilterEffectGreen  FilterEffect = 4
	FilterEffectNA     FilterEffect = -559038737 // 0xdeadbeef
)

// String returns human readable filter effect
func (f FilterEffect) String() string {
	switch f {
	case FilterEffectNone:
		return "None"
	case FilterEffectYellow:
		return "Yellow"
	case FilterEffectOrange:
		return "Orange"
	case FilterEffectRed:
		return "Red"
	case FilterEffectGreen:
		return "Green"
	case FilterEffectNA:
		return "n/a"
	default:
		return "Unknown"
	}
}

// ToningEffect represents monochrome toning effects
type ToningEffect int32

const (
	ToningEffectNone   ToningEffect = 0
	ToningEffectSepia  ToningEffect = 1
	ToningEffectBlue   ToningEffect = 2
	ToningEffectPurple ToningEffect = 3
	ToningEffectGreen  ToningEffect = 4
	ToningEffectNA     ToningEffect = -559038737 // 0xdeadbeef
)

// String returns human readable toning effect
func (t ToningEffect) String() string {
	switch t {
	case ToningEffectNone:
		return "None"
	case ToningEffectSepia:
		return "Sepia"
	case ToningEffectBlue:
		return "Blue"
	case ToningEffectPurple:
		return "Purple"
	case ToningEffectGreen:
		return "Green"
	case ToningEffectNA:
		return "n/a"
	default:
		return "Unknown"
	}
}

// FNumber represents a camera f-stop value
type FNumber float32

// FromRaw converts a raw uint8 value to FNumber
func FNumberFromRaw(raw uint8) FNumber {
	if raw == 0 {
		return 0
	}
	// ValueConv: exp((val-8)/16*log(2))
	return FNumber(math.Exp(float64(raw-8) / 16.0 * math.Log(2)))
}

// ToRaw converts an FNumber to raw uint8 value
func (f FNumber) ToRaw() uint8 {
	if f == 0 {
		return 0
	}
	// ValueConvInv: log(val)*16/log(2)+8
	return uint8(math.Log(float64(f))*16.0/math.Log(2) + 8)
}

// String returns formatted f-stop value with 2 significant digits
func (f FNumber) String() string {
	if f == 0 {
		return "undefined"
	}
	return fmt.Sprintf("%.2g", float64(f))
}

// ISO represents a camera ISO value
type ISO uint32

// ISOFromRaw converts a raw uint8 value from Camera Info ISO to ISO
func ISOFromRaw(raw uint8) ISO {
	// ValueConv: 100*exp((val/8-9)*log(2))
	return ISO(100 * math.Exp((float64(raw)/8-9)*math.Log(2)))
}

// ToRaw converts an ISO to raw uint8 Camera Info ISO value
func (i ISO) ToRaw() uint8 {
	// ValueConvInv: (log(val/100)/log(2)+9)*8
	return uint8((math.Log(float64(i)/100)/math.Log(2) + 9) * 8)
}

// String returns formatted ISO value
func (i ISO) String() string {
	return fmt.Sprintf("%.0f", float64(i))
}

// CameraTemperature represents the camera temperature in Celsius
type CameraTemperature uint8

// CameraTemperatureFromRaw converts a raw uint8 value to CameraTemperature
func CameraTemperatureFromRaw(raw uint8) CameraTemperature {
	// ValueConv: val - 128
	return CameraTemperature(int16(raw) - 128)
}

// ToRaw converts a CameraTemperature to raw uint8 value
func (t CameraTemperature) ToRaw() uint8 {
	// ValueConvInv: val + 128
	return uint8(int16(t) + 128)
}

// String returns formatted temperature with Celsius suffix
func (t CameraTemperature) String() string {
	return fmt.Sprintf("%d°C", int8(t))
}

// CIFocalLength represents a camera focal length in millimeters
type CIFocalLength uint16

// FromRaw converts a raw big-endian uint16 value to CIFocalLength
func CIFocalLengthFromRaw(raw uint16) CIFocalLength {
	if raw == 0 {
		return 0
	}
	// Convert from big-endian
	return CIFocalLength(binary.BigEndian.Uint16([]byte{byte(raw >> 8), byte(raw)}))
}

// ToRaw converts a CIFocalLength to raw big-endian uint16 value
func (f CIFocalLength) ToRaw() uint16 {
	if f == 0 {
		return 0
	}
	// Convert to big-endian
	return binary.BigEndian.Uint16([]byte{byte(uint16(f) >> 8), byte(uint16(f))})
}

// String returns formatted focal length with mm suffix
func (f CIFocalLength) String() string {
	if f == 0 {
		return "undefined"
	}
	return fmt.Sprintf("%d mm", uint16(f))
}

// CIExposureTime represents a camera exposure time value from camera info
type CIExposureTime float32

// FromRaw converts a raw uint8 value to CIExposureTime
func CIExposureTimeFromRaw(raw uint8) CIExposureTime {
	if raw == 0 {
		return 0
	}
	// ValueConv: exp(4*log(2)*(1-CanonEv(val-24)))
	return CIExposureTime(math.Exp(4 * math.Log(2) * (1 - float64(CanonEv(float32(raw-24))))))
}

// ToRaw converts a CIExposureTime to raw uint8 value
func (e CIExposureTime) ToRaw() uint8 {
	if e == 0 {
		return 0
	}
	// ValueConvInv: CanonEvInv(1-log(val)/(4*log(2)))+24
	return uint8(CanonEvInv(float32(1-math.Log(float64(e))/(4*math.Log(2))) + 24))
}

// String returns exposure time as a fraction or decimal
func (e CIExposureTime) String() string {
	if e == 0 {
		return "undefined"
	}

	// For times >= 1 second, show decimal
	if e >= 1 {
		return fmt.Sprintf("%.1f", float64(e))
	}

	// For times < 1 second, show as fraction (1/x)
	denominator := 1.0 / float64(e)
	return fmt.Sprintf("1/%.0f", denominator)
}
