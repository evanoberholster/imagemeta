package tag

import (
	"fmt"
)

// ExifVersion is the Exif Version
type ExifVersion [4]byte

// FromBytes parses the ExifVersion from TagValue
func (ev *ExifVersion) FromBytes(val TagValue) error {
	if len(val.Buf) >= 1 {
		copy(ev[:], val.Buf[:])
	}

	return nil
}

func (ev *ExifVersion) String() string {
	if ev[0] != 0 {
		return string(ev[:])
	}
	return ""
}

// ComponentsConfiguration is the Components Configuration array
type ComponentsConfiguration [4]byte

// FromBytes parses the Components Configuration from TagValue
func (cc *ComponentsConfiguration) FromBytes(val TagValue) error {
	if len(val.Buf) >= 1 {
		copy(cc[:], val.Buf[:])
	}
	return nil
}

// Flash is the Flash type
type Flash uint8

// ImageUnique is a 128-bit fixed length identifier
type ImageUnique [16]byte

// ShutterSpeedValue is the Shutter Speed Value
type ShutterSpeedValue SRational

// ApertureVale
type ApertureValue Rational

// BrightnessValue
type BrightnessValue SRational

// FromBytes parses FocalLength from TagValue
func (bv *BrightnessValue) FromBytes(val TagValue) error {
	return nil
}

// ExposureBiasValue
type ExposureBiasValue SRational

// FromBytes parses FocalLength from TagValue
func (ebv *ExposureBiasValue) FromBytes(val TagValue) error {
	return nil
}

// FocalLength
type FocalLength Rational

// FromBytes parses FocalLength from TagValue
func (fl *FocalLength) FromBytes(val TagValue) error {
	l := len(val.Buf)
	switch val.Type {
	case TypeShort:
		if l == 2 {
			*fl = FocalLength{uint32(val.ByteOrder.Uint16(val.Buf)), 1}
		}
	case TypeLong:
		*fl = FocalLength{val.ByteOrder.Uint32(val.Buf), 1}
	case TypeRational, TypeSignedRational:
		switch l {
		case 4:
			*fl = FocalLength{val.ByteOrder.Uint32(val.Buf[:4]), 1}
		case 8:
			*fl = FocalLength{val.ByteOrder.Uint32(val.Buf[:4]), val.ByteOrder.Uint32(val.Buf[4:8])}
		}
	}
	return nil
}

// Float returns the Focal Length as a Float
func (tfl FocalLength) Float() float64 {
	return Rational(tfl).Float()
}

func (tfl FocalLength) String() string {
	return fmt.Sprintf("%0.2fmm", tfl.Float())
}

// LensSpecification
type LensSpecification [4]Rational

func (tls LensSpecification) String() string {
	return fmt.Sprintf("%0.1f-%0.1fmm f/%0.2f - %0.2f", tls[0].Float(), tls[1].Float(), tls[2].Float(), tls[3].Float())
}

// ExposureProgram is the Exif Exposure Program
type ExposureProgram uint16

// Exposure Program types
const (
	ExposureProgamNotDefined ExposureProgram = 0 + iota
	ExposureProgramManual
	ExposureProgramNormalProgram
	ExposureProgramAperturePriority
	ExposureProgramShutterPriority
	ExposureProgramCreativeProgram
	ExposureProgramActionProgram
	ExposureProgramPotraitMode
	ExposureProgramLandscapeMode
)

var strTypeExposureProgram = [9]string{"Not Defined", "Manual", "Normal Program", "Aperture Priority", "Shutter Priority", "Creative Program", "Action Program", "Potrait Mode", "Landscape Mode"}

func (tep ExposureProgram) String() string {
	if int(tep) < len(strTypeExposureProgram) {
		return strTypeExposureProgram[tep]
	}
	return strTypeExposureProgram[0]
}

// MeteringMode represents the Metering Mode
type MeteringMode uint16

// MeteringMode types
const (
	MeteringModeUnknown MeteringMode = 0 + iota
	MeteringModeAverage
	MeteringModeCenterWeightedAverage
	MeteringModeSpot
	MeteringModeMultiSpot
	MeteringModePattern
	MeteringModePartial
	MeteringModeOther MeteringMode = 255
)

var strTypeMeteringMode = [7]string{"Unknown", "Average", "Center Weighted Average", "Spot", "MultiSpot", "Pattern", "Partial"}

func (tmm MeteringMode) String() string {
	if int(tmm) < len(strTypeMeteringMode) {
		return strTypeMeteringMode[tmm]
	}
	return "Other"
}

// ExposureMode represents the Exposure Mode
type ExposureMode uint16

// ExposureMode types
const (
	ExposureModeAuto ExposureMode = 0
	ExposureModeManual
	ExposureModeAutoBracket
)

// LightSource represents Exif Light Source types
type LightSource uint16

var strTypeLightSource = [25]string{
	"Unknown", "Daylight", "Fluorescent", "Tungsten (incandescent light)", "Flash", "Unknown", "Unknown", "Unknown", "Unknown",
	"Fine Weather",
	"Cloudy Weather",
	"Shade",
	"Daylight Fluorescent (D 5700 - 7100 K)",
	"Day White Fluorescent (N 4600 - 5500 K)",
	"Cool White Fluorescent (W 3800 - 4500 K)",
	"White Fluorescent (WW 3250 - 3800 K)",
	"Warm White Fluorescent (L 2600 - 3250 K)",
	"Standard Light A",
	"Standard Light B",
	"Standard Light C",
	"D55",
	"D65",
	"D75",
	"D50",
	"ISO Studio Tungsten"}

// LightSource types
const (
	LightSourceUnknown       LightSource = 0
	LightSourceDaylight      LightSource = 1
	LightSourceFluorescent   LightSource = 2
	LightSourceTungsten      LightSource = 3
	LightSourceFlash         LightSource = 4
	LightSourceFineWeather   LightSource = 9
	LightSourceCloudyWeather LightSource = 10
	LightSourceShade         LightSource = 11
	LightSourceDaylightFl    LightSource = 12
	LightSourceDayWhiteFl    LightSource = 13
	LightSourceCoolWhiteFl   LightSource = 14
	LightSourceWhiteFl       LightSource = 15
	LightSourceWarmWhiteFl   LightSource = 16
	LightSourceStandLightA   LightSource = 17
	LightSourceStandLightB   LightSource = 18
	LightSourceStandLightC   LightSource = 19
	LightSourceD55           LightSource = 20
	LightSourceD65           LightSource = 21
	LightSourceD75           LightSource = 22
	LightSourceD50           LightSource = 23
	LightSourceISOStudioTg   LightSource = 24
	LightSourceOther         LightSource = 255
)

// String implements the stringer interface for TypeLightSOurce
func (tls LightSource) String() string {
	if int(tls) < len(strTypeLightSource) {
		return strTypeLightSource[tls]
	}
	if tls == LightSourceOther {
		return "Other Light Source"
	}
	return strTypeLightSource[0]
}

// DateTime is an Exif DateTime value
type DateTime struct {
	Year  uint16
	Month uint8
	Day   uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}

// FromBytes parses a DateTime from TagValue
func (dt *DateTime) FromBytes(val TagValue) error {
	if val.Type.Is(TypeASCII) {
		// check recieved value
		buf := val.Buf
		if len(buf) >= 19 && buf[4] == ':' && buf[7] == ':' && buf[10] == ' ' &&
			buf[13] == ':' && buf[16] == ':' {
			*dt = DateTime{
				Year:  uint16(parseStrUint(buf[0:4])),
				Month: uint8(parseStrUint(buf[5:7])),
				Day:   uint8(parseStrUint(buf[8:10])),
				Hour:  uint8(parseStrUint(buf[11:13])),
				Min:   uint8(parseStrUint(buf[14:16])),
				Sec:   uint8(parseStrUint(buf[17:19])),
			}
		}
	}
	return nil
}

// OffsetTime is an Exif OffsetTime value
type OffsetTime struct {
	Hour int8
	Min  int8
}

// FromBytes parses an OffsetTime from TagValue
func (ot *OffsetTime) FromBytes(val TagValue) error {
	//if val.Type.Is(TypeASCII) {
	if len(val.Buf) == 7 && val.Buf[3] == ':' {
		hour := int(parseStrUint(val.Buf[1:3]))
		min := int(parseStrUint(val.Buf[4:6]))
		if val.Buf[0] == '-' {
			hour *= -1
		}
		*ot = OffsetTime{
			Hour: int8(hour),
			Min:  int8(min),
		}
	}
	//}
	return nil
}

// SubSecTime is an Exif SubSecTime value
type SubSecTime uint16

// FromBytes parses a SubSecTime from TagValue
func (sst *SubSecTime) FromBytes(val TagValue) error {
	switch val.Type {
	case TypeASCII, TypeASCIINoNul, TypeByte:
		*sst = SubSecTime(parseStrUint(val.Buf))
	}
	return nil
}
