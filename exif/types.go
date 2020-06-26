package exif

import (
	"fmt"
	"strconv"
)

// FocalLength - Focal Length in which the image was captured
type FocalLength float32

func (fl FocalLength) String() string {
	return fmt.Sprintf("%.2fmm", fl)
}

// MeteringMode - Mode in which the image was metered.
type MeteringMode uint8

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
	return "Not Defined"
}

// MeteringModeValues -
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
var meteringModeValues = map[MeteringMode]string{
	0:   "Unknown",
	1:   "Average",
	2:   "Center-weighted average",
	3:   "Spot",
	4:   "Multi-spot",
	5:   "Multi-segment",
	6:   "Partial",
	255: "Other",
}

// ExposureMode - Mode in which the Exposure was taken.
type ExposureMode uint8

// String - Return Exposure Mode as a string
func (em ExposureMode) String() string {
	return exposureModeValues[em]
}

// ExposureModeValues -
// Derived from https://sno.phy.queensu.ca/~phil/exiftool/TagNames/EXIF.html (23/09/2019)
var exposureModeValues = map[ExposureMode]string{
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

// FlashMode - Mode in which a Flash was used.
// (uint8) - value of FlashMode
type FlashMode uint8

// String - Return string for FlashMode
func (fm FlashMode) String() string {
	return flashValues[fm]
}

// Bool - Returns true if Flash was fired
func (fm FlashMode) Bool() bool {
	return flashBoolValues[fm]
}

// flashBoolValues -
// (bool) - true if the flash was fired
var flashBoolValues = map[FlashMode]bool{
	0:  false,
	1:  true,
	5:  true,
	7:  true,
	8:  false,
	9:  true,
	13: true,
	15: true,
	16: false,
	20: false,
	24: false,
	25: true,
	29: true,
	31: true,
	32: false,
	48: false,
	65: true,
	69: true,
	71: true,
	73: true,
	77: true,
	79: true,
	80: false,
	88: false,
	89: true,
	93: true,
	95: true,
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

// ShutterSpeed - [0] Numerator [1] Denominator
type ShutterSpeed [2]uint32

// String - return a ShutterSpeed as a string
func (ss ShutterSpeed) String() string {
	if ss[1] == 0 {
		return strconv.Itoa(int(ss[0]))
	}
	if ss[0] == 0 {
		return "Unknown"
	}
	return strconv.Itoa(int(ss[0])) + "/" + strconv.Itoa(int(ss[1]))
	//return fmt.Sprintf("%d/%d", ss[0], ss[1])
}

// MarshalJSON - Custom Marshall JSON
func (ss ShutterSpeed) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ss.String() + "\""), nil
}

// ExposureBias - [0] Numerator [1] Denominator
type ExposureBias [2]int16

// String - String value of Exposure Bias
func (eb ExposureBias) String() string {
	return strconv.Itoa(int(eb[0])) + "/" + strconv.Itoa(int(eb[1]))
	//return fmt.Sprintf("%d/%d", eb[0], eb[1])
}

// MarshalJSON - Custom Marshall JSON
func (eb ExposureBias) MarshalJSON() ([]byte, error) {
	return []byte("\"" + eb.String() + "\""), nil
}
