package canon

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate msgp

// MacroMode represents Canon macro mode settings
type MacroMode uint8

const (
	MacroModeMacro  MacroMode = 1
	MacroModeNormal MacroMode = 2
)

func (m MacroMode) String() string {
	switch m {
	case MacroModeMacro:
		return "Macro"
	case MacroModeNormal:
		return "Normal"
	default:
		return "Unknown"
	}
}

func (m MacroMode) Valid() error {
	switch m {
	case MacroModeMacro, MacroModeNormal:
		return nil
	default:
		return fmt.Errorf("invalid macro mode: %d", m)
	}
}

// SelfTimer represents Canon self-timer settings
type SelfTimer uint16

const (
	SelfTimerOff    SelfTimer = 0
	SelfTimerCustom SelfTimer = 0x4000
)

func (st SelfTimer) String() string {
	if st == 0 {
		return "Off"
	}
	seconds := float64(st&0xfff) / 10.0
	if st&SelfTimerCustom != 0 {
		return fmt.Sprintf("%.1f s, Custom", seconds)
	}
	return fmt.Sprintf("%.1f s", seconds)
}

func ParseSelfTimer(s string) (SelfTimer, error) {
	if strings.HasPrefix(strings.ToLower(s), "off") {
		return SelfTimerOff, nil
	}

	// Remove "s", "sec" and handle Custom flag
	s = strings.TrimSpace(s)
	isCustom := strings.Contains(strings.ToLower(s), "custom")
	s = strings.TrimSuffix(s, " s")
	s = strings.TrimSuffix(s, " sec")
	s = strings.TrimSuffix(s, ", Custom")

	seconds, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid self-timer value: %s", s)
	}

	val := SelfTimer(seconds * 10)
	if isCustom {
		val |= SelfTimerCustom
	}
	return val, nil
}

// CanonFlashMode represents Canon flash mode settings
type CanonFlashMode int8

const (
	FlashModeOff        CanonFlashMode = 0
	FlashModeAuto       CanonFlashMode = 1
	FlashModeOn         CanonFlashMode = 2
	FlashModeRedEye     CanonFlashMode = 3
	FlashModeSlowSync   CanonFlashMode = 4
	FlashModeRedEyeAuto CanonFlashMode = 5
	FlashModeRedEyeOn   CanonFlashMode = 6
	FlashModeExternal   CanonFlashMode = 16
)

// String returns the string representation of the flash mode
func (f CanonFlashMode) String() string {
	switch f {
	case FlashModeOff:
		return "Off"
	case FlashModeAuto:
		return "Auto"
	case FlashModeOn:
		return "On"
	case FlashModeRedEye:
		return "Red-eye reduction"
	case FlashModeSlowSync:
		return "Slow-sync"
	case FlashModeRedEyeAuto:
		return "Red-eye reduction (Auto)"
	case FlashModeRedEyeOn:
		return "Red-eye reduction (On)"
	case FlashModeExternal:
		return "External flash"
	default:
		return "Unknown"
	}
}

// Valid efficiently validates flash mode values using range checks
func (f CanonFlashMode) Valid() error {
	// Check continuous range 0-6 and special case 16
	if (f >= FlashModeOff && f <= FlashModeRedEyeOn) || f == FlashModeExternal {
		return nil
	}
	return fmt.Errorf("invalid flash mode: %d", f)
}

// ContinuousDrive is part of the CanonCameraSettings field
//
//	0:  "Single",
//	1:  "Continuous",
//	2:  "Movie",
//	3:  "Continuous, Speed Priority",
//	4:  "Continuous, Low",
//	5:  "Continuous, High",
//	6:  "Silent Single",
//	7:  "Unknown", // not defined
//	8:  "Unknown", // not defined
//	9:  "Single, Silent",
//	10: "Continuous, Silent",
type ContinuousDrive int16

const (
	strContinuousDriveString = "SingleContinuousMovieContinuous, Speed PriorityContinuous, LowContinuous, HighSilent SingleUnknownUnknownSingle, SilentContinuous, Silent"
)

var (
	strContinuousDriveDist = []int{0, 6, 16, 21, 47, 62, 78, 91, 98, 105, 119, 137}
)

// String returns the string representation of the ContinuousDrive value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the ContinuousDrive value is out of range, it returns "Unknown".
func (ccd ContinuousDrive) String() string {
	if int(ccd) < len(strContinuousDriveDist)-1 {
		return strContinuousDriveString[strContinuousDriveDist[ccd]:strContinuousDriveDist[ccd+1]]
	}
	return "Unknown"
}

// stringToContinuousDriveMap maps string representations to ContinuousDrive values.
var stringToContinuousDriveMap = map[string]ContinuousDrive{
	"Single":                     0,
	"Continuous":                 1,
	"Movie":                      2,
	"Continuous, Speed Priority": 3,
	"Continuous, Low":            4,
	"Continuous, High":           5,
	"Silent Single":              6,
	"Unknown":                    7,
	"Single, Silent":             9,
	"Continuous, Silent":         10,
}

// StringToContinuousDrive converts a string to a ContinuousDrive value.
// It returns the ContinuousDrive value and a boolean indicating whether the conversion was successful.
func StringToContinuousDrive(s string) (ContinuousDrive, bool) {
	cd, ok := stringToContinuousDriveMap[s]
	return cd, ok
}

// CanonFocusMode represents Canon focus mode settings
type CanonFocusMode uint16

const (
	FocusModeOneShotAF     CanonFocusMode = 0
	FocusModeAIServoAF     CanonFocusMode = 1
	FocusModeAIFocusAF     CanonFocusMode = 2
	FocusModeManual3       CanonFocusMode = 3
	FocusModeSingle        CanonFocusMode = 4
	FocusModeContinuous    CanonFocusMode = 5
	FocusModeManual6       CanonFocusMode = 6
	FocusModePanFocus      CanonFocusMode = 16
	FocusModeOneShotAFLive CanonFocusMode = 256
	FocusModeAIServoAFLive CanonFocusMode = 257
	FocusModeAIFocusAFLive CanonFocusMode = 258
	FocusModeMovieSnap     CanonFocusMode = 512
	FocusModeMovieServo    CanonFocusMode = 519

	focusModeMaxValue = FocusModeMovieServo
)

// Concatenated string containing all focus mode names
const focusModeNames = "One-shot AF\x00" +
	"AI Servo AF\x00" +
	"AI Focus AF\x00" +
	"Manual Focus (3)\x00" +
	"Single\x00" +
	"Continuous\x00" +
	"Manual Focus (6)\x00" +
	"Pan Focus\x00" +
	"One-shot AF (Live View)\x00" +
	"AI Servo AF (Live View)\x00" +
	"AI Focus AF (Live View)\x00" +
	"Movie Snap Focus\x00" +
	"Movie Servo AF\x00" +
	"Unknown"

// Array of start/end indices into focusModeNames
var focusModeIndices = [focusModeMaxValue + 1][2]uint16{
	FocusModeOneShotAF:     {0, 10},
	FocusModeAIServoAF:     {11, 22},
	FocusModeAIFocusAF:     {23, 34},
	FocusModeManual3:       {35, 50},
	FocusModeSingle:        {51, 57},
	FocusModeContinuous:    {58, 68},
	FocusModeManual6:       {69, 84},
	FocusModePanFocus:      {85, 94},
	FocusModeOneShotAFLive: {95, 116},
	FocusModeAIServoAFLive: {117, 138},
	FocusModeAIFocusAFLive: {139, 160},
	FocusModeMovieSnap:     {161, 176},
	FocusModeMovieServo:    {177, 191},
}

// String returns the string representation of the focus mode
func (f CanonFocusMode) String() string {
	if f > focusModeMaxValue {
		return focusModeNames[192:] // "Unknown"
	}
	if indices := focusModeIndices[f]; indices[0] < indices[1] {
		return focusModeNames[indices[0]:indices[1]]
	}
	return focusModeNames[192:] // "Unknown"
}

// CanonRecordMode represents Canon image recording modes
type CanonRecordMode int16

const (
	RecordModeJPEG    CanonRecordMode = 1
	RecordModeCRWTHM  CanonRecordMode = 2  // 300D, etc
	RecordModeAVITHM  CanonRecordMode = 3  // 30D
	RecordModeTIF     CanonRecordMode = 4  // 1Ds (unconfirmed)
	RecordModeTIFJPEG CanonRecordMode = 5  // 1D (unconfirmed)
	RecordModeCR2     CanonRecordMode = 6  // 1D,30D,350D
	RecordModeCR2JPEG CanonRecordMode = 7  // S30
	RecordModeMOV     CanonRecordMode = 9  // S95 MOV
	RecordModeMP4     CanonRecordMode = 10 // SX280 MP4
	RecordModeCRM     CanonRecordMode = 11 // C200 CRM
	RecordModeCR3     CanonRecordMode = 12 // EOS R
	RecordModeCR3JPEG CanonRecordMode = 13 // EOS R
	RecordModeHIF     CanonRecordMode = 14 // NC
	RecordModeCR3HIF  CanonRecordMode = 15 // 1DXmkIII
)

// Single concatenated string containing all mode names
const recordModeNames = "JPEGCRW+THMAVI+THMTIFTIF+JPEGCR2CR2+JPEGMOVMP4CRMCR3CR3+JPEGHIFCR3+HIF"

var (
	// Index array for efficient string slicing
	recordModeIndices = []int{0, 4, 11, 18, 21, 29, 32, 40, 43, 46, 49, 52, 55, 63, 66, 74}

	// recordModeMap maps string representations of record modes to their corresponding CanonRecordMode values.
	recordModeMap = map[string]CanonRecordMode{
		"JPEG":     RecordModeJPEG,
		"CRW+THM":  RecordModeCRWTHM,
		"AVI+THM":  RecordModeAVITHM,
		"TIF":      RecordModeTIF,
		"TIF+JPEG": RecordModeTIFJPEG,
		"CR2":      RecordModeCR2,
		"CR2+JPEG": RecordModeCR2JPEG,
		"MOV":      RecordModeMOV,
		"MP4":      RecordModeMP4,
		"CRM":      RecordModeCRM,
		"CR3":      RecordModeCR3,
		"CR3+JPEG": RecordModeCR3JPEG,
		"HIF":      RecordModeHIF,
		"CR3+HIF":  RecordModeCR3HIF,
	}
)

// String returns the string representation of the record mode
func (r CanonRecordMode) String() string {
	if r < 1 || int(r) >= len(recordModeIndices) {
		return "Unknown"
	}
	return recordModeNames[recordModeIndices[r-1]:recordModeIndices[r]]
}

// Valid checks if the record mode is valid
func (r CanonRecordMode) Valid() error {
	switch r {
	case RecordModeJPEG, RecordModeCRWTHM, RecordModeAVITHM,
		RecordModeTIF, RecordModeTIFJPEG, RecordModeCR2,
		RecordModeCR2JPEG, RecordModeMOV, RecordModeMP4,
		RecordModeCRM, RecordModeCR3, RecordModeCR3JPEG,
		RecordModeHIF, RecordModeCR3HIF:
		return nil
	}
	return fmt.Errorf("invalid record mode: %d", r)
}

// ParseRecordMode converts a string to CanonRecordMode using a map
func ParseRecordMode(s string) (CanonRecordMode, error) {

	if mode, ok := recordModeMap[s]; ok {
		return mode, nil
	}
	return 0, fmt.Errorf("unknown record mode: %s", s)
}

// MeteringMode is part of the CanonCameraSettings field
//
//	0: "Default",
//	1: "Spot",
//	2: "Average",
//	3: "Evaluative",
//	4: "Partial",
//	5: "Center-weighted average",
type MeteringMode int16

const (
	strCanonMeteringModeString = "DefaultSpotAverageEvaluativePartialCenter-weighted average"
)

var (
	strCanonMeteringModeDist = []int{0, 7, 11, 18, 28, 35, 58}

	// stringToMeteringModeMap maps string representations to MeteringMode values.
	stringToMeteringModeMap = map[string]MeteringMode{
		"Default":                 0,
		"Spot":                    1,
		"Average":                 2,
		"Evaluative":              3,
		"Partial":                 4,
		"Center-weighted average": 5,
	}
)

// String returns the string representation of the MeteringMode value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the MeteringMode value is out of range, it returns "Unknown".
func (mm MeteringMode) String() string {
	if int(mm) < len(strCanonMeteringModeDist)-1 {
		return strCanonMeteringModeString[strCanonMeteringModeDist[mm]:strCanonMeteringModeDist[mm+1]]
	}
	return "Unknown"
}

// StringToMeteringMode converts a string to a MeteringMode value.
// It returns the MeteringMode value and a boolean indicating whether the conversion was successful.
func StringToMeteringMode(s string) (MeteringMode, bool) {
	mm, ok := stringToMeteringModeMap[s]
	return mm, ok
}

// FocusRange is part of the CanonCameraSettings field
//
//	0:  "Manual",
//	1:  "Auto",
//	2:  "Not Known",
//	3:  "Macro",
//	4:  "Very Close",
//	5:  "Close",
//	6:  "Middle Range",
//	7:  "Far Range",
//	8:  "Pan Focus",
//	9:  "Super Macro",
//	10: "Infinity",
type FocusRange int16

const (
	strCanonFocusRangeString = "ManualAutoNot KnownMacroVery CloseCloseMiddle RangeFar RangePan FocusSuper MacroInfinity"
)

var (
	strCanonFocusRangeDist = []int{0, 6, 10, 19, 24, 34, 39, 51, 60, 70, 81}
)

var stringToFocusRangeMap = map[string]FocusRange{
	"Manual":       0,
	"Auto":         1,
	"Not Known":    2,
	"Macro":        3,
	"Very Close":   4,
	"Close":        5,
	"Middle Range": 6,
	"Far Range":    7,
	"Pan Focus":    8,
	"Super Macro":  9,
	"Infinity":     10,
}

// String returns the string representation of the FocusRange value.
// If the FocusRange value is out of range, it returns "Unknown".
func (fr FocusRange) String() string {
	if int(fr) < len(strCanonFocusRangeDist)-1 {
		return strCanonFocusRangeString[strCanonFocusRangeDist[fr]:strCanonFocusRangeDist[fr+1]]
	}
	return "Unknown"
}

// StringToFocusRange converts a string to a FocusRange value.
func StringToFocusRange(s string) (FocusRange, bool) {
	fr, ok := stringToFocusRangeMap[s]
	return fr, ok
}

// ExposureMode is part of the CanonCameraSettings field
//
//	0: "Easy",
//	1: "Program AE",
//	2: "Shutter speed priority AE",
//	3: "Aperture-priority AE",
//	4: "Manual",
//	5: "Depth-of-field AE",
//	6: "M-Dep",
//	7: "Bulb",
//	8: "Flexible-priority AE",
type ExposureMode int16

const (
	strCanonExposureModeString = "EasyProgram AEShutter speed priority AEAperture-priority AEManualDepth-of-field AEM-DepBulbFlexible-priority AE"
)

var (
	strCanonExposureModeDist = []int{0, 4, 14, 39, 60, 66, 84, 89, 93, 113}
)

// String returns the string representation of the ExposureMode value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the ExposureMode value is out of range, it returns "Unknown".
func (em ExposureMode) String() string {
	if int(em) < len(strCanonExposureModeDist)-1 {
		return strCanonExposureModeString[strCanonExposureModeDist[em]:strCanonExposureModeDist[em+1]]
	}
	return "Unknown"
}

// stringToExposureModeMap maps string representations to ExposureMode values.
var stringToExposureModeMap = map[string]ExposureMode{
	"Easy":                      0,
	"Program AE":                1,
	"Shutter speed priority AE": 2,
	"Aperture-priority AE":      3,
	"Manual":                    4,
	"Depth-of-field AE":         5,
	"M-Dep":                     6,
	"Bulb":                      7,
	"Flexible-priority AE":      8,
}

// StringToExposureMode converts a string to an ExposureMode value.
// It returns the ExposureMode value and a boolean indicating whether the conversion was successful.
func StringToExposureMode(s string) (ExposureMode, bool) {
	em, ok := stringToExposureModeMap[s]
	return em, ok
}

// FocusDistance -
type FocusDistance [2]int16

// NewFocusDistance creates a new FocusDistance with the upper
// and lower limits
func NewFocusDistance(upper, lower uint16) FocusDistance {
	return FocusDistance{int16(upper), int16(lower)}
}

// BracketMode - Canon Makernote Backet Mode
//
//	0: "Off",
//	1: "AEB",
//	2: "FEB",
//	3: "ISO",
//	4: "WB",
type BracketMode int16

// Active - returns true if BracketMode is On
func (bm BracketMode) Active() bool {
	return bm != 0
}

const (
	strCanonBracketModeString = "OffAEBFEBISOWB"
)

var (
	strCanonBracketModeDist = []int{0, 3, 6, 9, 12, 14}
)

// String returns the string representation of the BracketMode value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the BracketMode value is out of range, it returns "Unknown".
func (bm BracketMode) String() string {
	if int(bm) < len(strCanonBracketModeDist)-1 {
		return strCanonBracketModeString[strCanonBracketModeDist[bm]:strCanonBracketModeDist[bm+1]]
	}
	return "Unknown"
}

// stringToBracketModeMap maps string representations to BracketMode values.
var stringToBracketModeMap = map[string]BracketMode{
	"Off": 0,
	"AEB": 1,
	"FEB": 2,
	"ISO": 3,
	"WB":  4,
}

// StringToBracketMode converts a string to a BracketMode value.
// It returns the BracketMode value and a boolean indicating whether the conversion was successful.
func StringToBracketMode(s string) (BracketMode, bool) {
	bm, ok := stringToBracketModeMap[s]
	return bm, ok
}

// AESetting - Canon Makernote AutoExposure Setting
//
//	0: "Normal AE",
//	1: "Exposure Compensation",
//	2: "AE Lock",
//	3: "AE Lock + Exposure Compensation",
//	4: "No AE",
type AESetting int16

const (
	strCanonAESettingString = "Normal AEExposure CompensationAE LockAE Lock + Exposure CompensationNo AE"
)

var (
	strCanonAESettingDist = []int{0, 9, 30, 37, 72, 77}
)

// String returns the string representation of the AESetting value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the AESetting value is out of range, it returns "Unknown".
func (ae AESetting) String() string {
	if int(ae) < len(strCanonAESettingDist)-1 {
		return strCanonAESettingString[strCanonAESettingDist[ae]:strCanonAESettingDist[ae+1]]
	}
	return "Unknown"
}

// stringToAESettingMap maps string representations to AESetting values.
var stringToAESettingMap = map[string]AESetting{
	"Normal AE":                       0,
	"Exposure Compensation":           1,
	"AE Lock":                         2,
	"AE Lock + Exposure Compensation": 3,
	"No AE":                           4,
}

// StringToAESetting converts a string to an AESetting value.
// It returns the AESetting value and a boolean indicating whether the conversion was successful.
func StringToAESetting(s string) (AESetting, bool) {
	ae, ok := stringToAESettingMap[s]
	return ae, ok
}

// AFAreaMode - Canon Autofocus Area Mode
//
//	0:  "Off (Manual Focus)",
//	1:  "AF Point Expansion (surround)",
//	2:  "Single-point AF",
//	4:  "Auto",
//	5:  "Face Detect AF",
//	6:  "Face + Tracking",
//	7:  "Zone AF",
//	8:  "AF Point Expansion (4 point)",
//	9:  "Spot AF",
//	10: "AF Point Expansion (8 point)",
//	11: "Flexizone Multi (49 point)",
//	12: "Flexizone Multi (9 point)",
//	13: "Flexizone Single",
//	14: "Large Zone AF",
type AFAreaMode int16

const (
	strCanonAFAreaModeString = "Off (Manual Focus)AF Point Expansion (surround)Single-point AFAutoFace Detect AFFace + TrackingZone AFAF Point Expansion (4 point)Spot AFAF Point Expansion (8 point)Flexizone Multi (49 point)Flexizone Multi (9 point)Flexizone SingleLarge Zone AF"
)

var (
	strCanonAFAreaModeDist = []int{0, 18, 47, 61, 65, 79, 94, 101, 127, 134, 160, 186, 202, 215}
)

// String returns the string representation of the AFAreaMode value.
// It uses efficient string slicing with an array of integers and a single concatenated string.
// If the AFAreaMode value is out of range, it returns "Unknown".
func (caf AFAreaMode) String() string {
	if int(caf) < len(strCanonAFAreaModeDist)-1 {
		return strCanonAFAreaModeString[strCanonAFAreaModeDist[caf]:strCanonAFAreaModeDist[caf+1]]
	}
	return "Unknown"
}

// stringToAFAreaModeMap maps string representations to AFAreaMode values.
var stringToAFAreaModeMap = map[string]AFAreaMode{
	"Off (Manual Focus)":            0,
	"AF Point Expansion (surround)": 1,
	"Single-point AF":               2,
	"Auto":                          4,
	"Face Detect AF":                5,
	"Face + Tracking":               6,
	"Zone AF":                       7,
	"AF Point Expansion (4 point)":  8,
	"Spot AF":                       9,
	"AF Point Expansion (8 point)":  10,
	"Flexizone Multi (49 point)":    11,
	"Flexizone Multi (9 point)":     12,
	"Flexizone Single":              13,
	"Large Zone AF":                 14,
}

// StringToAFAreaMode converts a string to an AFAreaMode value.
// It returns the AFAreaMode value and a boolean indicating whether the conversion was successful.
func StringToAFAreaMode(s string) (AFAreaMode, bool) {
	am, ok := stringToAFAreaModeMap[s]
	return am, ok
}

// PictureStyle represents a Canon picture style.
type PictureStyle uint16

const (
	PictureStyleNone           PictureStyle = 0x00
	PictureStyleStandard       PictureStyle = 0x01
	PictureStylePortrait       PictureStyle = 0x02
	PictureStyleHighSaturation PictureStyle = 0x03
	PictureStyleAdobeRGB       PictureStyle = 0x04
	PictureStyleLowSaturation  PictureStyle = 0x05
	PictureStyleCMSet1         PictureStyle = 0x06
	PictureStyleCMSet2         PictureStyle = 0x07
	PictureStyleUserDef1       PictureStyle = 0x21
	PictureStyleUserDef2       PictureStyle = 0x22
	PictureStyleUserDef3       PictureStyle = 0x23
	PictureStylePC1            PictureStyle = 0x41
	PictureStylePC2            PictureStyle = 0x42
	PictureStylePC3            PictureStyle = 0x43
	PictureStyleStandardAlt    PictureStyle = 0x81
	PictureStylePortraitAlt    PictureStyle = 0x82
	PictureStyleLandscape      PictureStyle = 0x83
	PictureStyleNeutral        PictureStyle = 0x84
	PictureStyleFaithful       PictureStyle = 0x85
	PictureStyleMonochrome     PictureStyle = 0x86
	PictureStyleAuto           PictureStyle = 0x87
	PictureStyleFineDetail     PictureStyle = 0x88
	PictureStyleNA             PictureStyle = 0xff
	PictureStyleNAAlt          PictureStyle = 0xffff
)

// pictureStyles maps PictureStyle values to their string representations.
var pictureStyles = map[PictureStyle]string{
	PictureStyleNone:           "None",
	PictureStyleStandard:       "Standard",
	PictureStylePortrait:       "Portrait",
	PictureStyleHighSaturation: "High Saturation",
	PictureStyleAdobeRGB:       "Adobe RGB",
	PictureStyleLowSaturation:  "Low Saturation",
	PictureStyleCMSet1:         "CM Set 1",
	PictureStyleCMSet2:         "CM Set 2",
	PictureStyleUserDef1:       "User Def. 1",
	PictureStyleUserDef2:       "User Def. 2",
	PictureStyleUserDef3:       "User Def. 3",
	PictureStylePC1:            "PC 1",
	PictureStylePC2:            "PC 2",
	PictureStylePC3:            "PC 3",
	PictureStyleStandardAlt:    "Standard",
	PictureStylePortraitAlt:    "Portrait",
	PictureStyleLandscape:      "Landscape",
	PictureStyleNeutral:        "Neutral",
	PictureStyleFaithful:       "Faithful",
	PictureStyleMonochrome:     "Monochrome",
	PictureStyleAuto:           "Auto",
	PictureStyleFineDetail:     "Fine Detail",
	PictureStyleNA:             "n/a",
	PictureStyleNAAlt:          "n/a",
}

// String returns the string representation of the PictureStyle value.
func (ps PictureStyle) String() string {
	if str, ok := pictureStyles[ps]; ok {
		return str
	}
	return "Unknown"
}

// IsValid checks if the PictureStyle value is valid.
func (ps PictureStyle) IsValid() bool {
	_, ok := pictureStyles[ps]
	return ok
}

// CanonQuality represents Canon image quality settings
type CanonQuality int16

const (
	QualityUnknown   CanonQuality = 0
	QualityEconomy   CanonQuality = 1
	QualityNormal    CanonQuality = 2
	QualityFine      CanonQuality = 3
	QualityRAW       CanonQuality = 4
	QualitySuperfine CanonQuality = 5
	QualityCRAW      CanonQuality = 7
	QualityLightRAW  CanonQuality = 130
	QualityStdRAW    CanonQuality = 131
)

// canonQualityMap maps quality values to their string representations
var canonQualityMap = map[CanonQuality]string{
	QualityEconomy:   "Economy",
	QualityNormal:    "Normal",
	QualityFine:      "Fine",
	QualityRAW:       "RAW",
	QualitySuperfine: "Superfine",
	QualityCRAW:      "CRAW",
	QualityLightRAW:  "Light (RAW)",
	QualityStdRAW:    "Standard (RAW)",
}

// String returns the string representation of the quality setting
func (q CanonQuality) String() string {
	if str, ok := canonQualityMap[q]; ok {
		return str
	}
	return "Unknown"
}

// CanonImageSize represents Canon image size settings
type CanonImageSize uint16

const (
	ImageSizeLarge          CanonImageSize = 0
	ImageSizeMedium         CanonImageSize = 1
	ImageSizeSmall          CanonImageSize = 2
	ImageSizeMedium1        CanonImageSize = 5
	ImageSizeMedium2        CanonImageSize = 6
	ImageSizeMedium3        CanonImageSize = 7
	ImageSizePostcard       CanonImageSize = 8
	ImageSizeWidescreen     CanonImageSize = 9
	ImageSizeMediumWide     CanonImageSize = 10
	ImageSizeSmall1         CanonImageSize = 14
	ImageSizeSmall2         CanonImageSize = 15
	ImageSizeSmall3         CanonImageSize = 16
	ImageSize640x480Movie   CanonImageSize = 128
	ImageSizeMediumMovie    CanonImageSize = 129
	ImageSizeSmallMovie     CanonImageSize = 130
	ImageSize1280x720Movie  CanonImageSize = 137
	ImageSize1920x1080Movie CanonImageSize = 142
	ImageSize4096x2160Movie CanonImageSize = 143
)

// canonImageSizeMap maps size values to their string representations
var canonImageSizeMap = map[CanonImageSize]string{
	ImageSizeLarge:          "Large",
	ImageSizeMedium:         "Medium",
	ImageSizeSmall:          "Small",
	ImageSizeMedium1:        "Medium 1",
	ImageSizeMedium2:        "Medium 2",
	ImageSizeMedium3:        "Medium 3",
	ImageSizePostcard:       "Postcard",
	ImageSizeWidescreen:     "Widescreen",
	ImageSizeMediumWide:     "Medium Widescreen",
	ImageSizeSmall1:         "Small 1",
	ImageSizeSmall2:         "Small 2",
	ImageSizeSmall3:         "Small 3",
	ImageSize640x480Movie:   "640x480 Movie",
	ImageSizeMediumMovie:    "Medium Movie",
	ImageSizeSmallMovie:     "Small Movie",
	ImageSize1280x720Movie:  "1280x720 Movie",
	ImageSize1920x1080Movie: "1920x1080 Movie",
	ImageSize4096x2160Movie: "4096x2160 Movie",
}

// String returns the string representation of the image size
func (s CanonImageSize) String() string {
	if str, ok := canonImageSizeMap[s]; ok {
		return str
	}
	return "Unknown"
}

// CanonEasyMode represents Canon camera easy shooting modes
// References:
//   - http://homepage3.nifty.com/kamisaka/makernote/makernote_canon.htm (kamisaka)
//   - http://www.burren.cx/david/canon.html (burren)
//   - Canon DPP 3.11.26 software
type CanonEasyMode uint16

const (
	EasyModeFullAuto CanonEasyMode = iota
	EasyModeManual
	EasyModeLandscape
	EasyModeFastShutter
	EasyModeSlowShutter
	EasyModeNight
	EasyModeGrayScale
	EasyModeSepia
	EasyModePortrait
	EasyModeSports
	EasyModeMacro
	EasyModeBlackWhite
	EasyModePanFocus
	EasyModeVivid
	EasyModeNeutral
	EasyModeFlashOff
	EasyModeLongShutter
	EasyModeSuperMacro
	EasyModeFoliage
	EasyModeIndoor
	EasyModeFireworks
	EasyModeBeach
	EasyModeUnderwater
	EasyModeSnow
	EasyModeKidsPets
	EasyModeNightSnapshot
	EasyModeDigitalMacro
	EasyModeMyColors
	EasyModeMovieSnap
	EasyModeSuperMacro2
	EasyModeColorAccent
	EasyModeColorSwap
	EasyModeAquarium
	EasyModeISO3200
	EasyModeISO6400
	EasyModeCreativeLightEffect
	EasyModeEasy
	EasyModeQuickShot
	EasyModeCreativeAuto
	EasyModeZoomBlur
	EasyModeLowLight
	EasyModeNostalgic
	EasyModeSuperVivid
	EasyModePosterEffect
	EasyModeFaceSelfTimer
	EasyModeSmile
	EasyModeWinkSelfTimer
	EasyModeFisheyeEffect
	EasyModeMiniatureEffect
	EasyModeHighSpeedBurst
	EasyModeBestImageSelection
	EasyModeHighDynamicRange
	EasyModeHandheldNightScene
	EasyModeMovieDigest
	EasyModeLiveViewControl
	EasyModeDiscreet
	EasyModeBlurReduction
	EasyModeMonochrome
	EasyModeToyCameraEffect
	EasyModeSceneIntelligentAuto
	EasyModeHighSpeedBurstHQ
	EasyModeSmoothSkin
	EasyModeSoftFocus
	EasyModeFood           = 68
	EasyModeHDRArtStandard = 84
	EasyModeHDRArtVivid    = 85
	EasyModeHDRArtBold     = 93
	EasyModeSpotlight      = 257
	EasyModeNight2         = 258
	EasyModeNightPlus      = 259
	EasyModeSuperNight     = 260
	EasyModeSunset         = 261
	EasyModeNightScene     = 263
	EasyModeSurface        = 264
	EasyModeLowLight2      = 265

	// Concatenated string containing all mode names without padding for efficient slicing
	strCanonEasyModeString = "Full auto" + // 0
		"Manual" +
		"Landscape" +
		"Fast shutter" +
		"Slow shutter" +
		"Night" +
		"Gray Scale" +
		"Sepia" +
		"Portrait" +
		"Sports" +
		"Macro" +
		"Black & White" + // 11
		"Pan focus" +
		"Vivid" +
		"Neutral" +
		"Flash Off" + // 15
		"Long Shutter" +
		"Super Macro" +
		"Foliage" +
		"Indoor" +
		"Fireworks" + // 20
		"Beach" +
		"Underwater" +
		"Snow" +
		"Kids & Pets" +
		"Night Snapshot" + // 25
		"Digital Macro" +
		"My Colors" +
		"Movie Snap" +
		"Super Macro 2" +
		"Color Accent" + // 30
		"Color Swap" +
		"Aquarium" +
		"ISO 3200" +
		"ISO 6400" +
		"Creative Light Effect" + // 35
		"Easy" +
		"Quick Shot" +
		"Creative Auto" +
		"Zoom Blur" +
		"Low Light" + // 40
		"Nostalgic" +
		"Super Vivid" +
		"Poster Effect" +
		"Face Self-timer" +
		"Smile" + // 45
		"Wink Self-timer" +
		"Fisheye Effect" +
		"Miniature Effect" +
		"High-speed Burst" +
		"Best Image Selection" + // 50
		"High Dynamic Range" +
		"Handheld Night Scene" +
		"Movie Digest" +
		"Live View Control" +
		"Discreet" + // 55
		"Blur Reduction" +
		"Monochrome" +
		"Toy Camera Effect" +
		"Scene Intelligent Auto" +
		"High-speed Burst HQ" + // 60
		"Smooth Skin" +
		"Soft Focus" + // 62
		"Food" + // 68
		"HDR Art Standard" + // 84
		"HDR Art Vivid" + // 85
		"HDR Art Bold" + // 93
		"Spotlight" + // 257
		"Night 2" +
		"Night+" +
		"Super Night" +
		"Sunset" + // 261
		"Night Scene" + // 263
		"Surface" +
		"Low Light 2" // 265
)

var (
	// Map for string to CanonEasyMode conversion
	easyModeMap = map[string]CanonEasyMode{
		"Full auto": 0, "Manual": 1, "Landscape": 2, "Fast shutter": 3,
		"Slow shutter": 4, "Night": 5, "Gray Scale": 6, "Sepia": 7,
		"Portrait": 8, "Sports": 9, "Macro": 10, "Black & White": 11,
		"Pan focus": 12, "Vivid": 13, "Neutral": 14, "Flash Off": 15,
		"Long Shutter": 16, "Super Macro": 17, "Foliage": 18, "Indoor": 19,
		"Fireworks": 20, "Beach": 21, "Underwater": 22, "Snow": 23,
		"Kids & Pets": 24, "Night Snapshot": 25, "Digital Macro": 26,
		"My Colors": 27, "Movie Snap": 28, "Super Macro 2": 29,
		"Color Accent": 30, "Color Swap": 31, "Aquarium": 32,
		"ISO 3200": 33, "ISO 6400": 34, "Creative Light Effect": 35,
		"Easy": 36, "Quick Shot": 37, "Creative Auto": 38, "Zoom Blur": 39,
		"Low Light": 40, "Nostalgic": 41, "Super Vivid": 42, "Poster Effect": 43,
		"Face Self-timer": 44, "Smile": 45, "Wink Self-timer": 46,
		"Fisheye Effect": 47, "Miniature Effect": 48, "High-speed Burst": 49,
		"Best Image Selection": 50, "High Dynamic Range": 51,
		"Handheld Night Scene": 52, "Movie Digest": 53, "Live View Control": 54,
		"Discreet": 55, "Blur Reduction": 56, "Monochrome": 57,
		"Toy Camera Effect": 58, "Scene Intelligent Auto": 59,
		"High-speed Burst HQ": 60, "Smooth Skin": 61, "Soft Focus": 62,
		"Food": 68, "HDR Art Standard": 84, "HDR Art Vivid": 85,
		"HDR Art Bold": 93, "Spotlight": 257, "Night 2": 258, "Night+": 259,
		"Super Night": 260, "Sunset": 261, "Night Scene": 263, "Surface": 264,
		"Low Light 2": 265,
	}

	// Index array for efficient string slicing
	strCanonEasyModeDist = []int{0, 9, 15, 24, 36, 48, 53, 63, 68, 76, 82, 87, 100, 109, 114, 121, 130, 142, 152, 159, 165, 174, 179, 189, 193, 204, 218, 230, 239, 249, 262, 274, 284, 292, 300, 320, 324, 334, 347, 356, 365, 374, 385, 398, 412, 417, 422, 436, 451, 465, 484, 500, 519, 530, 545, 555, 568, 578, 595, 612, 623, 632, 636, 652, 665, 677, 686, 694, 700, 705, 716, 721, 730, 737, 744}
)

// ParseEasyMode converts a string to CanonEasyMode
func ParseEasyMode(s string) (CanonEasyMode, error) {
	if mode, ok := easyModeMap[s]; ok {
		return mode, nil
	}
	return 0, fmt.Errorf("unknown easy mode: %s", s)
}

// String returns the string representation of the easy mode
func (e CanonEasyMode) String() string {
	switch {
	// Handle 0-62 range
	case e <= 62:
		return strCanonEasyModeString[strCanonEasyModeDist[e]:strCanonEasyModeDist[e+1]]

	// Handle special cases
	case e == 68: // Food
		return strCanonEasyModeString[strCanonEasyModeDist[63]:strCanonEasyModeDist[64]]
	case e == 84: // HDR Art Standard
		return strCanonEasyModeString[strCanonEasyModeDist[64]:strCanonEasyModeDist[65]]
	case e == 85: // HDR Art Vivid
		return strCanonEasyModeString[strCanonEasyModeDist[65]:strCanonEasyModeDist[66]]
	case e == 93: // HDR Art Bold
		return strCanonEasyModeString[strCanonEasyModeDist[66]:strCanonEasyModeDist[67]]

	// Handle 257-265 range (excluding 262)
	case e >= 257 && e <= 261:
		idx := 67 + int(e-257)
		return strCanonEasyModeString[strCanonEasyModeDist[idx]:strCanonEasyModeDist[idx+1]]
	case e >= 263 && e <= 265:
		idx := 67 + int(e-257-1) // Subtract 1 to skip 262
		return strCanonEasyModeString[strCanonEasyModeDist[idx]:strCanonEasyModeDist[idx+1]]
	}

	return "Unknown"
}

// FocusContinuous represents Canon continuous focus settings
type FocusContinuous uint8

const (
	// Standard focus modes
	FocusContinuousSingle     FocusContinuous = 0
	FocusContinuousContinuous FocusContinuous = 1
	FocusContinuousManual     FocusContinuous = 8
)

// String returns a human-readable representation of the focus mode
func (fc FocusContinuous) String() string {
	switch fc {
	case FocusContinuousSingle:
		return "Single"
	case FocusContinuousContinuous:
		return "Continuous"
	case FocusContinuousManual:
		return "Manual"
	default:
		return "Unknown"
	}
}

// Valid checks if the focus mode value is valid
func (fc FocusContinuous) Valid() error {
	switch fc {
	case FocusContinuousSingle, FocusContinuousContinuous, FocusContinuousManual:
		return nil
	}
	return fmt.Errorf("invalid focus continuous mode: %d", fc)
}

// ImageStabilization represents Canon image stabilization settings
type ImageStabilization int16

const (
	// ImageStabilizationUndef represents undefined value (-1)
	ImageStabilizationUndef ImageStabilization = -1

	// Standard range (0-4)
	ImageStabilizationOff ImageStabilization = iota
	ImageStabilizationOn
	ImageStabilizationShootOnly
	ImageStabilizationPanning
	ImageStabilizationDynamic

	// Extended range (256-260)
	ImageStabilizationOff2       ImageStabilization = 256
	ImageStabilizationOn2        ImageStabilization = 257
	ImageStabilizationShootOnly2 ImageStabilization = 258
	ImageStabilizationPanning2   ImageStabilization = 259
	ImageStabilizationDynamic2   ImageStabilization = 260
)

// Concatenated string containing all mode names
const strImageStabilizationString = "OffOnShoot OnlyPanningDynamicOff (2)On (2)Shoot Only (2)Panning (2)Dynamic (2)"

// Indices into the concatenated string for efficient slicing
var strImageStabilizationDist = []int{0, 3, 5, 15, 22, 29, 36, 42, 55, 66, 77}

// String returns a human-readable representation of the image stabilization mode
func (is ImageStabilization) String() string {
	switch {
	case is == ImageStabilizationUndef:
		return "Undefined"
	case is >= 0 && is <= 4:
		return strImageStabilizationString[strImageStabilizationDist[is]:strImageStabilizationDist[is+1]]
	case is >= 256 && is <= 260:
		idx := int(is - 256 + 5)
		return strImageStabilizationString[strImageStabilizationDist[idx]:strImageStabilizationDist[idx+1]]
	default:
		return "Unknown"
	}
}

// Valid checks if the image stabilization mode is valid
func (is ImageStabilization) Valid() error {
	switch {
	case is == ImageStabilizationUndef:
		return nil
	case is >= 0 && is <= 4:
		return nil
	case is >= 256 && is <= 260:
		return nil
	}
	return fmt.Errorf("invalid image stabilization mode: %d", is)
}

// SRAWQuality represents Canon sRAW quality settings
type SRAWQuality int16

const (
	SRAWQualityNA    SRAWQuality = iota // 0: n/a
	SRAWQualitySRAW1                    // 1: mRAW
	SRAWQualitySRAW2                    // 2: sRAW
)

// String returns a human-readable representation of the sRAW quality
func (sq SRAWQuality) String() string {
	switch sq {
	case SRAWQualityNA:
		return "n/a"
	case SRAWQualitySRAW1:
		return "sRAW1 (mRAW)"
	case SRAWQualitySRAW2:
		return "sRAW2 (sRAW)"
	default:
		return "Unknown"
	}
}

// Valid checks if the sRAW quality value is valid
func (sq SRAWQuality) Valid() error {
	if sq >= SRAWQualityNA && sq <= SRAWQualitySRAW2 {
		return nil
	}
	return fmt.Errorf("invalid sRAW quality: %d", sq)
}

// FlashBits represents Canon flash settings bitmask
type FlashBits uint16

const (
	FlashBitsNone     FlashBits = 0
	FlashBitsManual   FlashBits = 1 << 0  // bit 0
	FlashBitsTTL      FlashBits = 1 << 1  // bit 1
	FlashBitsATTL     FlashBits = 1 << 2  // bit 2
	FlashBitsETTL     FlashBits = 1 << 3  // bit 3
	FlashBitsFPSync   FlashBits = 1 << 4  // bit 4
	FlashBits2ndSync  FlashBits = 1 << 7  // bit 7
	FlashBitsFPUsed   FlashBits = 1 << 11 // bit 11
	FlashBitsBuiltIn  FlashBits = 1 << 13 // bit 13
	FlashBitsExternal FlashBits = 1 << 14 // bit 14
)

// Concatenated string of all flash bit names
const strFlashBitsString = " None ManualTTLA-TTLE-TTLFP sync enabled2nd-curtain sync usedFP sync usedBuilt-inExternal"

// Indices into the concatenated string for efficient slicing
var strFlashBitsDist = []int{0, 6, 12, 15, 20, 25, 39, 58, 69, 77, 85}

// String returns a human-readable representation of the flash bits without allocations
func (fb FlashBits) String() string {
	if fb == FlashBitsNone {
		return strFlashBitsString[strFlashBitsDist[0]:strFlashBitsDist[1]]
	}

	// Pre-allocate buffer for worst case: all flags + separators
	var buf [150]byte
	pos := 0
	needComma := false

	// Check each bit and write to buffer
	masks := [...]struct {
		bit   FlashBits
		start int
		end   int
	}{
		{FlashBitsManual, strFlashBitsDist[1], strFlashBitsDist[2]},
		{FlashBitsTTL, strFlashBitsDist[2], strFlashBitsDist[3]},
		{FlashBitsATTL, strFlashBitsDist[3], strFlashBitsDist[4]},
		{FlashBitsETTL, strFlashBitsDist[4], strFlashBitsDist[5]},
		{FlashBitsFPSync, strFlashBitsDist[5], strFlashBitsDist[6]},
		{FlashBits2ndSync, strFlashBitsDist[6], strFlashBitsDist[7]},
		{FlashBitsFPUsed, strFlashBitsDist[7], strFlashBitsDist[8]},
		{FlashBitsBuiltIn, strFlashBitsDist[8], strFlashBitsDist[9]},
		{FlashBitsExternal, strFlashBitsDist[9], strFlashBitsDist[10]},
	}

	for _, m := range masks {
		if fb&m.bit != 0 {
			if needComma {
				copy(buf[pos:], ", ")
				pos += 2
			}
			n := copy(buf[pos:], strFlashBitsString[m.start:m.end])
			pos += n
			needComma = true
		}
	}

	if pos == 0 {
		return "Unknown"
	}
	return string(buf[:pos])
}

// Valid checks if the flash bits value contains valid flags
func (fb FlashBits) Valid() error {
	// Create mask of all valid bits
	const validBits = FlashBitsManual | FlashBitsTTL | FlashBitsATTL |
		FlashBitsETTL | FlashBitsFPSync | FlashBits2ndSync |
		FlashBitsFPUsed | FlashBitsBuiltIn | FlashBitsExternal

	// Check if any invalid bits are set
	if fb != FlashBitsNone && fb&^validBits != 0 {
		return fmt.Errorf("invalid flash bits: %d", fb)
	}
	return nil
}

// FocalUnits represents the conversion factor from raw focal length values to millimeters
type FocalUnits int16

// String returns the string representation of focal units in mm format
func (f FocalUnits) String() string {
	return fmt.Sprintf("%d/mm", int(f))
}

// Valid checks if the focal units value is valid (non-zero)
func (f FocalUnits) Valid() error {
	if f <= 0 {
		return fmt.Errorf("invalid focal units: %d", f)
	}
	return nil
}

// FocalLength represents the minimum focal length in millimeters
type FocalLength int16

// ToRaw converts millimeter value to raw units using focal units
func (f FocalLength) ToRaw(units FocalUnits) int16 {
	if units == 0 {
		units = 1
	}
	return int16(f) * int16(units)
}

// FromRaw converts raw value to millimeters using focal units
func FromRawMin(raw int16, units FocalUnits) FocalLength {
	if units == 0 {
		units = 1
	}
	return FocalLength(raw) / FocalLength(units)
}

// String returns the string representation in millimeters
func (f FocalLength) String() string {
	return fmt.Sprintf("%d mm", int(f))
}

// PhotoEffect represents Canon photo effect settings
type PhotoEffect uint8

const (
	PhotoEffectOff PhotoEffect = iota
	PhotoEffectVivid
	PhotoEffectNeutral
	PhotoEffectSmooth
	PhotoEffectSepia
	PhotoEffectBW
	PhotoEffectCustom

	PhotoEffectMyColorData PhotoEffect = 100
)

// Concatenated string of all effect names
const strPhotoEffectString = "OffVividNeutralSmoothSepiaB&WCustomMy Color Data"

// Indices into the concatenated string for efficient slicing
var strPhotoEffectDist = []int{0, 3, 8, 15, 21, 26, 30, 36, 49}

// String returns a human-readable representation of the photo effect
func (pe PhotoEffect) String() string {
	if pe == PhotoEffectMyColorData {
		return strPhotoEffectString[strPhotoEffectDist[7]:strPhotoEffectDist[8]]
	}
	if pe <= PhotoEffectCustom {
		return strPhotoEffectString[strPhotoEffectDist[pe]:strPhotoEffectDist[pe+1]]
	}
	return "Unknown"
}

// Valid checks if the photo effect value is valid
func (pe PhotoEffect) Valid() error {
	if (pe <= PhotoEffectCustom) || pe == PhotoEffectMyColorData {
		return nil
	}
	return fmt.Errorf("invalid photo effect: %d", pe)
}

// ParsePhotoEffect converts a string to PhotoEffect
func ParsePhotoEffect(s string) (PhotoEffect, error) {
	if s == "My Color Data" {
		return PhotoEffectMyColorData, nil
	}
	for i := PhotoEffectOff; i <= PhotoEffectCustom; i++ {
		if s == strPhotoEffectString[strPhotoEffectDist[i]:strPhotoEffectDist[i+1]] {
			return i, nil
		}
	}
	return 0, fmt.Errorf("unknown photo effect: %s", s)
}

// DisplayAperture represents the camera's display aperture value in f-stops
type DisplayAperture uint16

// FromRaw converts a raw integer value to DisplayAperture (divided by 10)
func DisplayApertureFromRaw(raw uint16) DisplayAperture {
	return DisplayAperture(raw)
}

// ToFloat32 converts DisplayAperture to float32 f-stop value
func (a DisplayAperture) ToFloat32() float32 {
	return float32(a) / 10.0
}

// String returns the string representation of the aperture value
func (a DisplayAperture) String() string {
	return fmt.Sprintf("f/%.1f", a.ToFloat32())
}

// SpotMeteringMode represents Canon spot metering modes
type SpotMeteringMode uint8

const (
	SpotMeteringModeCenter SpotMeteringMode = iota
	SpotMeteringModeAFPoint
)

// String returns a human-readable representation of the spot metering mode
func (sm SpotMeteringMode) String() string {
	switch sm {
	case SpotMeteringModeCenter:
		return "Center"
	case SpotMeteringModeAFPoint:
		return "AF Point"
	default:
		return "Unknown"
	}
}

// Clarity represents Canon clarity setting for EOS R models
type Clarity int16

// String returns the string representation of the clarity value
func (c Clarity) String() string {
	return fmt.Sprintf("%d", c)
}

// AFPointSetting represents Canon autofocus point settings
type AFPointSetting uint16

const (
	AFPointManualSelection AFPointSetting = 0x2005
	AFPointNoneMF          AFPointSetting = 0x3000
	AFPointAutoSelection   AFPointSetting = 0x3001
	AFPointRight           AFPointSetting = 0x3002
	AFPointCenter          AFPointSetting = 0x3003
	AFPointLeft            AFPointSetting = 0x3004
	AFPointAutoSelection2  AFPointSetting = 0x4001
	AFPointFaceDetect      AFPointSetting = 0x4006
)

// String returns a human-readable representation of the AF point
func (ap AFPointSetting) String() string {
	switch ap {
	case 0:
		return "Invalid"
	case AFPointManualSelection:
		return "Manual AF point selection"
	case AFPointNoneMF:
		return "None (MF)"
	case AFPointAutoSelection:
		return "Auto AF point selection"
	case AFPointRight:
		return "Right"
	case AFPointCenter:
		return "Center"
	case AFPointLeft:
		return "Left"
	case AFPointAutoSelection2:
		return "Auto AF point selection"
	case AFPointFaceDetect:
		return "Face Detect"
	default:
		return fmt.Sprintf("Unknown (0x%04X)", uint16(ap))
	}
}

// Valid checks if the AF point value is valid
func (ap AFPointSetting) Valid() error {
	switch ap {
	case 0, AFPointManualSelection, AFPointNoneMF, AFPointAutoSelection,
		AFPointRight, AFPointCenter, AFPointLeft,
		AFPointAutoSelection2, AFPointFaceDetect:
		return nil
	}
	return fmt.Errorf("invalid AF point: 0x%04X", uint16(ap))
}
