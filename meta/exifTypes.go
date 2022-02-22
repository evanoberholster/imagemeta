package meta

import (
	"strconv"
)

//go:generate msgp

// FocalLength is a Focal Length expressed in millimeters.
type FocalLength float32

const (
	// FocalLength Suffix in millimeters
	sufFocalLength = "mm"
)

// NewFocalLength returns a new FocalLength by dividing
// "n" numerator and "d" demoninator
func NewFocalLength(n, d uint32) FocalLength {
	return FocalLength(float32(n) / float32(d))
}

func (fl FocalLength) String() string {
	return string(fl.toBytes())
}

func (fl FocalLength) toBytes() (buf []byte) {
	f := strconv.AppendFloat(buf, float64(fl), 'f', 2, 32)
	return append(f, sufFocalLength...)
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (fl FocalLength) MarshalText() (text []byte, err error) {
	return fl.toBytes(), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (fl *FocalLength) UnmarshalText(text []byte) (err error) {
	var f float64
	if len(text) > 0 {
		if text[len(text)-1] == sufFocalLength[1] && text[len(text)-2] == sufFocalLength[0] {
			text = text[:len(text)-2]
		}
		f, err = strconv.ParseFloat(string(text), 32)
		*fl = FocalLength(f)
		return
	}
	return nil
}

// Aperture contains the F-Number.
type Aperture float32

// NewAperture returns a new Aperture by dividing the
// "n" numerator over the "d" demoninator
func NewAperture(n uint32, d uint32) Aperture {
	return Aperture(float32(n) / float32(d))
}

// ParseString parses a string for an aperture value.
// ex: 1/100 or 300/100
func (aa *Aperture) ParseString(buf []byte) error {
	*aa = parseAperture(buf)
	return nil
}

func parseAperture(buf []byte) Aperture {
	// TODO: Improve parsing functionality
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			if i < len(buf)+1 {
				n := uint16(parseUint(buf[:i]))
				d := uint16(parseUint(buf[i+1:]))
				return Aperture(n / d)
			}
		}
	}
	return Aperture(0)
}

func (aa Aperture) String() string {
	if uint32(aa*100.0)%100 == 0 {
		return strconv.FormatFloat(float64(aa), 'f', 0, 32)
	}
	return strconv.FormatFloat(float64(aa), 'f', 2, 32)
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (aa Aperture) MarshalText() (text []byte, err error) {
	buf := make([]byte, 0, 4)
	buf = strconv.AppendFloat(buf, float64(aa), 'f', 2, 32)
	return buf, nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json.
func (aa *Aperture) UnmarshalText(text []byte) (err error) {
	fl, err := strconv.ParseFloat(string(text), 32)
	*aa = Aperture(fl)
	return err
}

// ShutterSpeed contains the shutter speed in seconds.
// Limit to 1/2 and 1/3 stops
// [0] Numerator [1] Denominator
type ShutterSpeed [2]uint32

// NewShutterSpeed creates a new ShutterSpeed with "n" as numerator and
// "d" as denominator
func NewShutterSpeed(n uint32, d uint32) ShutterSpeed {
	return ShutterSpeed{n, d}
}

// parseShutterSpeed parses a ShutterSpeed time value from []byte.
// Example: For less than 1 second: (1/250)
// Example: For more than 1 second: (1.3)
func parseShutterSpeed(buf []byte) (ss ShutterSpeed) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			if i < len(buf)+1 {
				return ShutterSpeed{uint32(parseUint(buf[:i])), uint32(parseUint(buf[i+1:]))}
			}
		}
		if buf[i] == '.' {
			if i < len(buf)+1 {
				ua := parseUint(buf[:i])
				b := buf[i+1:]
				if len(b) > 1 {
					b = b[:1]
				}
				ub := parseUint(b)
				if ub == 0 {
					return ShutterSpeed{uint32(ua), 1}
				}
				ua *= 10
				ua += ub
				return ShutterSpeed{uint32(ua), 10}
			}
		}
	}
	return
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (ss ShutterSpeed) MarshalText() (text []byte, err error) {
	if ss[1] != 0 {
		if ss[0] == 1 {
			if ss[1] == 1 {
				return []byte{'1', '.', '0'}, nil
			}
			n := 8
			if ss[1] < 100 {
				n = 4
			} else if ss[1] < 1000 {
				n = 5
			}
			text = make([]byte, 2, n)
			text[0] = '1'
			text[1] = '/'
			return strconv.AppendUint(text, uint64(ss[1]), 10), nil
		}
		if ss[0] > 1 {
			v := ss[0] / ss[1]
			r := ss[0] % ss[1]
			text = strconv.AppendUint(text, uint64(v), 10)
			text = append(text, '.')
			return strconv.AppendUint(text, uint64(r), 10), nil
		}
	}
	return []byte{'0'}, nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (ss *ShutterSpeed) UnmarshalText(text []byte) (err error) {
	*ss = parseShutterSpeed(text)
	return nil
}

// String returns a ShutterSpeed as a string
func (ss ShutterSpeed) String() string {
	buf, _ := ss.MarshalText()
	return string(buf)
}

// ExposureBias is the Exposure Bias of an image expressed as
// a positive or negative fraction.
// Bit1 = Sign
// Bit2-7 = Numerator
// Bit8-15 = Denominator
// Bit16 = Empty
type ExposureBias int16

const (
	exposureBiasZero = "0/0"
)

// NewExposureBias creates a new Exposure Bias from the provided
// "n" as numerator and "d" as denominator.
//
// To parse a string expressed as a positive or negative fraction.
// ex: "+1/3" or ex: "-2/3" use ExposureBias.UnmarshalText.
func NewExposureBias(n int16, d int16) ExposureBias {
	n = n << 8
	d = d << 8 >> 8
	return ExposureBias(n + d)
}

// String returns the value of Exposure Bias as a string
func (eb ExposureBias) String() string {
	buf, _ := eb.MarshalText()
	return string(buf)
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (eb ExposureBias) MarshalText() (text []byte, err error) {
	if eb == 0 {
		return unsafeGetBytes(exposureBiasZero), nil
	}
	text = make([]byte, 0, 5)
	if eb > 0 {
		text = append(text, '+')
	}
	text = strconv.AppendInt(text, int64(eb>>8), 10)
	text = append(text, '/')
	text = strconv.AppendUint(text, uint64(uint16(eb)<<8>>8), 10)
	return text, nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (eb *ExposureBias) UnmarshalText(text []byte) (err error) {
	if text[0] == '0' {
		return
	}
	for i := 0; i < len(text); i++ {
		if text[i] == '/' {
			if i < len(text)+1 {
				var n int16
				if text[0] == '+' {
					n = int16(parseUint(text[1:i]))
				} else if text[0] == '-' {
					n = int16(parseUint(text[1:i])) * -1
				} else {
					n = int16(parseUint(text[:i]))
				}
				n = n << 8
				n += int16(parseUint(text[i+1:]))
				*eb = ExposureBias(n)
				return err
			}
		}
	}
	return
}

// MeteringMode is the mode in which the image was metered.
type MeteringMode uint8

// Metering Modes
const (
	MeteringModeUnknown MeteringMode = iota
	MeteringModeAverage
	MeteringModeCenterWeightedAverage
	MeteringModeSpot
	MeteringModeMultispot
	MeteringModeMultisegment
	MeteringModePartial
	MeteringModeOther MeteringMode = 255

	// MeteringModeName
	_MeteringModeName = "UnknownAverageCenter-weighted averageSpotMulti-spotMulti-segmentPartial"
)

// MeteringMode values
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
var (
	_MeteringModeIndex    = [...]uint8{0, 7, 14, 37, 41, 51, 64, 71}
	mapStringMeteringMode = map[string]MeteringMode{
		"Unknown":                 MeteringModeUnknown,
		"Average":                 MeteringModeAverage,
		"Center-weighted average": MeteringModeCenterWeightedAverage,
		"Spot":                    MeteringModeSpot,
		"Multi-spot":              MeteringModeMultispot,
		"Multi-segment":           MeteringModeMultisegment,
		"Partial":                 MeteringModePartial,
		"Other":                   MeteringModeOther,
	}
)

// NewMeteringMode returns a MeteringMode from the given uint8
func NewMeteringMode(meteringMode uint8) MeteringMode {
	if meteringMode < 7 || meteringMode == 255 {
		return MeteringMode(meteringMode)
	}
	return 0
}

// String - Return Metering Mode as a string
//
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
func (mm MeteringMode) String() string {
	if int(mm) < len(_MeteringModeIndex)-1 {
		return _MeteringModeName[_MeteringModeIndex[mm]:_MeteringModeIndex[mm+1]]
	}
	if mm == MeteringModeOther {
		return "Other"
	}
	return _MeteringModeName[:_MeteringModeIndex[1]]
}

// MarshalJSON implements the JSONMarshaler interface that is
// used by encoding/json
func (mm MeteringMode) MarshalJSON() (buf []byte, err error) {
	return strconv.AppendUint(buf, uint64(mm), 10), nil
}

// UnmarshalJSON implements the JSONMarshaler interface that is
// used by encoding/json
func (mm *MeteringMode) UnmarshalJSON(buf []byte) error {
	v, err := strconv.ParseUint(string(buf), 10, 8)
	*mm = MeteringMode(v)
	return err
}

// MarshalText implements the TextMarshaler interface
func (mm MeteringMode) MarshalText() (text []byte, err error) {
	return unsafeGetBytes(mm.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (mm *MeteringMode) UnmarshalText(text []byte) (err error) {
	*mm = mapStringMeteringMode[string(text)]
	return nil
}

// ExposureMode is the mode in which the Exposure was taken.
//
//  Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (07/02/2021)
type ExposureMode uint8

// Exposure Modes
const (
	ExposureModeAuto ExposureMode = iota
	ExposureModeManual
	ExposureModeAutoBracket

	// ExposureMode Stringer
	_ExposureModeName = "AutoManualAuto bracket"
)

// Exposure Mode Values
var (
	_ExposureModeIndex    = [...]uint8{0, 4, 10, 22}
	mapStringExposureMode = map[string]ExposureMode{
		"Auto":         ExposureModeAuto,
		"Manual":       ExposureModeManual,
		"Auto bracket": ExposureModeAutoBracket,
	}
)

// NewExposureMode returns an ExposureMode from the given uint8
func NewExposureMode(em uint8) ExposureMode {
	if em <= 2 {
		return ExposureMode(em)
	}
	return 0
}

// String returns an ExposureMode as a string
func (em ExposureMode) String() string {
	if int(em) < len(_ExposureModeIndex)-1 {
		return _ExposureModeName[_ExposureModeIndex[em]:_ExposureModeIndex[em+1]]
	}
	return "Unknown"
}

// MarshalText implements the TextMarshaler interface
func (em ExposureMode) MarshalText() (text []byte, err error) {
	return unsafeGetBytes(em.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (em *ExposureMode) UnmarshalText(text []byte) (err error) {
	*em = mapStringExposureMode[string(text)]
	return nil
}

// ExposureProgram is the program in which the image was taken.
//
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
// 	0: "Not Defined",
// 	1: "Manual",
// 	2: "Program AE",
// 	3: "Aperture-priority AE",
// 	4: "Shutter speed priority AE",
// 	5: "Creative (Slow speed)",
// 	6: "Action (High speed)",
// 	7: "Portrait",
// 	8: "Landscape",
// 	9: "Bulb",
type ExposureProgram uint8

// Exposure Programs
const (
	ExposureProgramNotDefined ExposureProgram = iota
	ExposureProgramManual
	ExposureProgramProgramAE
	ExposureProgramAperturePriority
	ExposureProgramShutterSpeedPriority
	ExposureProgramCreative
	ExposureProgramAction
	ExposureProgramPortrait
	ExposureProgramLandscape
	ExposureProgramBulb

	// ExposureProgramName
	_ExposureProgramName = "Not DefinedManualProgram AEAperture-priority AEShutter speed priority AECreative (Slow speed)Action (High speed)PortraitLandscapeBulb"
)

// ExposureProgramIndex
var (
	_ExposureProgramIndex    = [...]uint8{0, 11, 17, 27, 47, 72, 93, 112, 120, 129, 133}
	mapStringExposureProgram = map[string]ExposureProgram{
		"Not Defined":               ExposureProgramNotDefined,
		"Manual":                    ExposureProgramManual,
		"Program AE":                ExposureProgramProgramAE,
		"Aperture-priority AE":      ExposureProgramAperturePriority,
		"Shutter speed priority AE": ExposureProgramShutterSpeedPriority,
		"Creative (Slow speed)":     ExposureProgramCreative,
		"Action (High speed)":       ExposureProgramAction,
		"Portrait":                  ExposureProgramPortrait,
		"Landscape":                 ExposureProgramLandscape,
		"Bulb":                      ExposureProgramBulb,
	}
)

// NewExposureProgram returns an ExposureProgram from the given uint8
func NewExposureProgram(ep uint8) ExposureProgram {
	if ep <= 9 {
		return ExposureProgram(ep)
	}
	return ExposureProgramNotDefined
}

// String returns an ExposureProgram as a string
func (ep ExposureProgram) String() string {
	if int(ep) < len(_ExposureProgramIndex)-1 {
		return _ExposureProgramName[_ExposureProgramIndex[ep]:_ExposureProgramIndex[ep+1]]
	}
	return ExposureProgramNotDefined.String()
}

// MarshalText implements the TextMarshaler interface
func (ep ExposureProgram) MarshalText() (text []byte, err error) {
	return unsafeGetBytes(ep.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (ep *ExposureProgram) UnmarshalText(text []byte) (err error) {
	*ep = mapStringExposureProgram[string(text)]
	return nil
}

// Flash is in bit format and represents the mode in which flash was used.
//
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
type Flash uint8

// Flashtypes
const (
	FlashNoFlash Flash = 0
	FlashFired   Flash = 1

	// FlashStringer Strings
	_FlashStringerStrings = "No FlashFiredFired, Return not detectedFired, Return detectedOn, Did not fireOn, FiredOn, Return not detectedOn, Return detectedOff, Did not fireOff, Did not fire, Return not detectedAuto, Did not fireAuto, FiredAuto, Fired, Return not detectedAuto, Fired, Return detectedNo flash functionOff, No flash functionFired, Red-eye reductionFired, Red-eye reduction, Return not detectedFired, Red-eye reduction, Return detectedOn, Red-eye reductionOn, Red-eye reduction, Return not detectedOn, Red-eye reduction, Return detectedOff, Red-eye reductionAuto, Did not fire, Red-eye reductionAuto, Fired, Red-eye reductionAuto, Fired, Red-eye reduction, Return not detectedAuto, Fired, Red-eye reduction, Return detected"
)

var (
	// FlashStringer Index
	_FlashStringerIndex = [...]uint16{0, 8,
		13, 13, 13, 13, 39, 39, 61, 77, 86, 86, 86, 86, 109, 109, 128, 145,
		145, 145, 145, 183, 183, 183, 183, 201, 212, 212, 212, 212, 244, 244,
		272, 289, 289, 289, 289, 289, 289, 289, 289, 289, 289, 289, 289, 289,
		289, 289, 289, 311, 311, 311, 311, 311, 311, 311, 311, 311, 311, 311,
		311, 311, 311, 311, 311, 311, 335, 335, 335, 335, 380, 380, 421, 421, 442,
		442, 442, 442, 484, 484, 522, 544, 544, 544, 544, 544, 544, 544, 544,
		581, 611, 611, 611, 611, 662, 662, 709}
)

// NewFlash returns a new Flash value
func NewFlash(f uint8) Flash {
	return Flash(f)
}

// String returns an ExposureProgram as a string
func (f Flash) String() string {
	if int(f) < len(_FlashStringerIndex)-1 {
		str := _FlashStringerStrings[_FlashStringerIndex[f]:_FlashStringerIndex[f+1]]
		if len(str) > 0 {
			return str
		}
	}
	return Flash(0).String()
}

// FlashMode is what mode the flash was used in
type FlashMode uint8

// FlashModes
const (
	FlashModeNone FlashMode = 0
	FlashNoReturn FlashMode = 4
	FlashReturn   FlashMode = 6
	FlashModeOn   FlashMode = 8
	FlashModeOff  FlashMode = 16
	FlashModeAuto FlashMode = 24
)

// Fired is bit 0, returns true if Flash was fired.
func (f Flash) Fired() bool {
	return 0b00000001&f == 0b00000001
}

// ReturnStatus is bits 1 and 2, returns 4 if "No Return" present and 6 if "Return" present.
// 	FlashNoReturn: 4
// 	FlashReturn:  6
func (f Flash) ReturnStatus() FlashMode {
	return FlashMode(0b00000110 & f)
}

// FlashFunction is bit 5, returns true if flash function was not present
func (f Flash) FlashFunction() bool {
	return 0b00100000&f == 0b00100000
}

// Mode is bits 3 and 4, returns 0 if "NoFlash", 8 if "On", 16 if "Off", and 24 if "Auto".
//  FlashModeNone: 0
// 	FlashModeOn: 8
// 	FlashModeOff: 16
// 	FlashModeAuto: 24
func (f Flash) Mode() FlashMode {
	return FlashMode(0b00011000 & f)
}

// Redeye is bit 6, returns true if "Red-eye reduction" was present
func (f Flash) Redeye() bool {
	return 0b01000000&f == 0b01000000
}

type Orientation uint8

const (
	OrientationHorizontal                Orientation = 1
	OrientationMirrorHorizontal          Orientation = 2
	OrientationRotate180                 Orientation = 3
	OrientationMirrorVertical            Orientation = 4
	OrientationMirrorHorizontalRotate270 Orientation = 5
	OrientationRotate90                  Orientation = 6
	OrientationMirrorHorizontalRotate90  Orientation = 7
	OrientationRotate270                 Orientation = 8
)

// String representation lifted from exiftool.
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (01/02/2022)
var orientationValues = map[Orientation]string{
	OrientationHorizontal:                "Horizontal (normal)",
	OrientationMirrorHorizontal:          "Mirror horizontal",
	OrientationRotate180:                 "Rotate 180",
	OrientationMirrorVertical:            "Mirror vertical",
	OrientationMirrorHorizontalRotate270: "Mirror horizontal and rotate 270 CW",
	OrientationRotate90:                  "Rotate 90 CW",
	OrientationMirrorHorizontalRotate90:  "Mirror horizontal and rotate 90 CW",
	OrientationRotate270:                 "Rotate 270 CW",
}

// String returns the value of Orientation as a string
func (o Orientation) String() string {
	str, ok := orientationValues[o]
	if ok {
		return str
	}
	return "Unknown"
}
