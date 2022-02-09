package meta

import (
	"math/bits"
	"strconv"
)

//go:generate msgp

var (
	_ExposureMode_index    = [...]uint8{0, 4, 10, 22}
	_ExposureProgram_index = [...]uint8{0, 11, 17, 27, 47, 72, 93, 112, 120, 129, 133}
	_MeteringMode_index    = [...]uint8{0, 7, 14, 37, 41, 51, 64, 71}

	// FocalLength Suffix in millimeters
	sufFocalLength = "mm"

	exposureBiasZero = []byte("0/0")
)

const (
	_ExposureMode_name    = "AutoManualAuto bracket"
	_MeteringMode_name    = "UnknownAverageCenter-weighted averageSpotMulti-spotMulti-segmentPartial"
	_ExposureProgram_name = "Not DefinedManualProgram AEAperture-priority AEShutter speed priority AECreative (Slow speed)Action (High speed)PortraitLandscapeBulb"
)

// FocalLength is a Focal Length expressed in millimeters.
type FocalLength float32

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
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
func (mm MeteringMode) String() string {
	if int(mm) < len(_MeteringMode_index)-1 {
		return _MeteringMode_name[_MeteringMode_index[mm]:_MeteringMode_index[mm+1]]
	}
	if mm == 255 {
		return "Other"
	}
	return _MeteringMode_name[:_MeteringMode_index[1]]
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
	return []byte(mm.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (mm *MeteringMode) UnmarshalText(text []byte) (err error) {
	*mm = mapStringMeteringMode[string(text)]
	return nil
}

// MeteringMode values
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
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
//
//  Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (07/02/2021)
//	0: "Auto",
//	1: "Manual",
//	2: "Auto bracket",
type ExposureMode uint8

// NewExposureMode returns an ExposureMode from the given uint8
func NewExposureMode(em uint8) ExposureMode {
	if em <= 2 {
		return ExposureMode(em)
	}
	return 255
}

// String returns an ExposureMode as a string
func (em ExposureMode) String() string {
	if int(em) < len(_ExposureMode_index)-1 {
		return _ExposureMode_name[_ExposureMode_index[em]:_ExposureMode_index[em+1]]
	}
	return "Unknown"
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

// NewExposureProgram returns an ExposureProgram from the given uint8
func NewExposureProgram(ep uint8) ExposureProgram {
	if ep <= 9 {
		return ExposureProgram(ep)
	}
	return 255
}

// String returns an ExposureProgram as a string
func (ep ExposureProgram) String() string {
	if int(ep) < len(_ExposureProgram_index)-1 {
		return _ExposureProgram_name[_ExposureProgram_index[ep]:_ExposureProgram_index[ep+1]]
	}
	return "Unknown"
}

// FlashMode - Mode in which a Flash was used.
// (uint8) - value of FlashMode
//
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
type FlashMode uint8

// NewFlashMode returns a new FlashMode
func NewFlashMode(fm uint8) FlashMode {
	return FlashMode(fm)
}

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
	*fm = parseFlashMode(uint8(i))
	return err
}

const (
	flashFiredFlag      = 00000001
	flashNoFunctionFlag = 00000040
)

// Fired returns true if Flash was fired.
func (fm FlashMode) Fired() bool {
	return bits.OnesCount8(flashFiredFlag&uint8(fm)) == 1
}

// NoFunction returns true if Flash function wasn't present.
func (fm FlashMode) NoFunction() bool {
	return bits.OnesCount8(flashNoFunctionFlag&uint8(fm)) == 1
}

// parseFlashMode returns the FlashMode from an Exif flashmode integer
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html#Flash (23/09/2019)
func parseFlashMode(m uint8) FlashMode {
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
		return exposureBiasZero, nil
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
