package tag

import "fmt"

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

// ExposureBiasValue
type ExposureBiasValue SRational

// FocalLength
type FocalLength Rational

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
