package meta

import (
	"strconv"
)

// FocalLength is a Focal Length expressed in tenth of a mm.
type FocalLength float32

var (
	sufFocalLength = []byte{'m', 'm'}
)

func (fl FocalLength) String() string {
	return string(fl.toBytes())
}

func (fl FocalLength) toBytes() (buf []byte) {
	f := strconv.AppendFloat(buf, float64(fl), 'f', 2, 32)
	//f = append(f, sufFocalLength...)
	buf = make([]byte, len(f)+2)
	copy(buf[len(buf)-2:], sufFocalLength)
	copy(buf[:len(buf)-2], f)
	return
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
	return err
}

// MeteringMode - Mode in which the image was metered.
type MeteringMode uint8

// NewMeteringMode returns a MeteringMode from the given uint8
func NewMeteringMode(meteringMode uint8) MeteringMode {
	if meteringMode < 7 || meteringMode == 255 {
		return MeteringMode(meteringMode)
	}
	return 0
}

// String - Return Metering Mode as a string
//
// MeteringMode values
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
func (mm MeteringMode) String() string {
	switch mm {
	case 0:
		return "Unknown"
	case 1:
		return "Average"
	case 2:
		return "Center-weighted average"
	case 3:
		return "Spot"
	case 4:
		return "Multi-spot"
	case 5:
		return "Multi-segment"
	case 6:
		return "Partial"
	case 255:
		return "Other"
	}
	// TODO: Create debug catch
	return "Unknown"
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (mm MeteringMode) MarshalText() (text []byte, err error) {
	return []byte(mm.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (mm *MeteringMode) UnmarshalText(text []byte) (err error) {
	*mm = mapStringMeteringMode[string(text)]
	return nil
}

var mapStringMeteringMode = map[string]MeteringMode{
	"Unknown":                 0,
	"Average":                 1,
	"Center-weighted average": 2,
	"Spot":                    3,
	"Multi-spot":              4,
	"Multi-segment":           5,
	"Partial":                 6,
	"Other":                   255,
}

// ExposureMode is the mode in which the Exposure was taken.
type ExposureMode uint8

// mapExposureModeString -
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (07/02/2021)
var mapExposureModeString = map[ExposureMode]string{
	0: "Auto",
	1: "Manual",
	2: "Auto bracket",
}

// String returns an ExposureMode as a string
func (em ExposureMode) String() string {
	return mapExposureModeString[em]
}

// ExposureProgram is the program in which the image was taken.
type ExposureProgram uint8

// mapExposureProgramString -
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
var mapExposureProgramString = map[ExposureProgram]string{
	0: "Not Defined",
	1: "Manual",
	2: "Program AE",
	3: "Aperture-priority AE",
	4: "Shutter speed priority AE",
	5: "Creative (Slow speed)",
	6: "Action (High speed)",
	7: "Portrait",
	8: "Landscape",
	9: "Bulb",
}

// String returns an ExposureProgram as a string
func (ep ExposureProgram) String() string {
	return mapExposureProgramString[ep]
}

// FlashMode - Mode in which a Flash was used.
// (uint8) - value of FlashMode
type FlashMode uint8

// Flash Modes
const (
	NoFlash           FlashMode = 0
	FlashFired        FlashMode = 1
	FlashOffNotFired  FlashMode = 16
	FlashAutoNotFired FlashMode = 24
	FlashAutoFired    FlashMode = 25
)

// String - Return string for FlashMode
func (fm FlashMode) String() string {
	return flashValues[fm]
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (fm FlashMode) MarshalText() (text []byte, err error) {
	return strconv.AppendUint(text, uint64(fm), 10), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (fm *FlashMode) UnmarshalText(text []byte) (err error) {
	var i int
	i, err = strconv.Atoi(string(text))
	*fm = ParseFlashMode(uint8(i))
	return err
}

// Bool returns true if Flash was fired.
func (fm FlashMode) Bool() bool {
	switch fm {
	case FlashFired, FlashAutoFired:
		return true
	case NoFlash, FlashAutoNotFired, FlashOffNotFired:
		return false
	}
	return false
}

// ParseFlashMode returns the FlashMode from an Exif flashmode integer
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
func ParseFlashMode(m uint8) FlashMode {
	switch m {
	case 0: // NoFlash
		return NoFlash
	case 25, 29, 31, 89, 93, 95: // Auto, Fired
		return FlashAutoFired
	case 24, 88: // Auto, Did not Fire
		return FlashAutoNotFired
	case 1, 5, 7, 9, 13, 15, 65, 69, 71, 73, 77, 79: // On, Fired
		return FlashFired
	case 8, 16, 20, 48, 80: // Off, Did not Fire
		return FlashOffNotFired
	}
	return NoFlash
}

// flashValues -
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
var flashValues = map[FlashMode]string{
	0:  "No Flash",
	1:  "Fired",
	5:  "Fired, Return not detected",
	7:  "Fired, Return detected",
	8:  "On, Did not fire",
	9:  "On, Fired",
	13: "On, Return not detected",
	15: "On, Return detected",
	16: "Off, Did not fire",
	20: "Off, Did not fire, Return not detected",
	24: "Auto, Did not fire",
	25: "Auto, Fired",
	29: "Auto, Fired, Return not detected",
	31: "Auto, Fired, Return detected",
	32: "No flash function",
	48: "Off, No flash function",
	65: "Fired, Red-eye reduction",
	69: "Fired, Red-eye reduction, Return not detected",
	71: "Fired, Red-eye reduction, Return detected",
	73: "On, Red-eye reduction",
	77: "On, Red-eye reduction, Return not detected",
	79: "On, Red-eye reduction, Return detected",
	80: "Off, Red-eye reduction",
	88: "Auto, Did not fire, Red-eye reduction",
	89: "Auto, Fired, Red-eye reduction",
	93: "Auto, Fired, Red-eye reduction, Return not detected",
	95: "Auto, Fired, Red-eye reduction, Return detected",
}

// Aperture contains the F-Number.
// [0] Numerator [1] Denominator
type Aperture [2]uint

// ShutterSpeed contains the shutter speed in seconds.
// Limit to 1/2 and 1/3 stops
// [0] Numerator [1] Denominator
type ShutterSpeed [2]uint16

// parseShutterSpeed parses a ShutterSpeed time value from []byte.
// Example: For less than 1 second: (1/250)
// Example: For more than 1 second: (1.3)
func parseShutterSpeed(buf []byte) (ss ShutterSpeed) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == '/' {
			if i < len(buf)+1 {
				return ShutterSpeed{uint16(parseUint(buf[:i])), uint16(parseUint(buf[i+1:]))}
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
					return ShutterSpeed{uint16(ua), 1}
				}
				ua *= 10
				ua += ub
				return ShutterSpeed{uint16(ua), 10}
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

// ExposureBias - [0] Numerator [1] Denominator
type ExposureBias [2]int16

// String - String value of Exposure Bias
func (eb ExposureBias) String() string {
	buf, _ := eb.MarshalText()
	return string(buf)
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (eb ExposureBias) MarshalText() (text []byte, err error) {
	if eb[1] != 0 {
		text = strconv.AppendInt(text, int64(eb[0]), 10)
		text = append(text, '/')
		text = strconv.AppendInt(text, int64(eb[1]), 10)
		return
	}
	return []byte{'0', '/', '0'}, nil
}
