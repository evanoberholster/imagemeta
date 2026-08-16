// Package sony provides data types and functions for representing Sony Camera
// Makernote values.
package sony

import "strings"

//go:generate stringer -type=WhiteBalanceValue,WhiteBalance0xB054Value,QualityValue,SceneModeValue,DynamicRangeOptimizerValue,FocusModeValue,FocusMode0xB042Value,AFAreaModeSettingValue,MeteringModeValue,FlashModeValue,ReleaseModeValue,PictureEffectValue,CreativeStyleValue,LongExposureNoiseReductionValue,HighISONoiseReductionValue,AntiBlurValue,IntelligentAutoValue,ExposureModeValue,ColorModeValue,AFAreaModeValue,ZoneMatchingValue,ImageStabilizationValue,SoftSkinEffectValue,MacroValue,AFIlluminatorValue,JPEGQualityValue,RAWFileTypeValue,PrioritySetInAWBValue,FlashActionValue,AFTrackingValue,MultiFrameNREffectValue,HDRValue -linecomment -output=sony_string.go

// WhiteBalanceValue describes the white-balance setting stored in Sony
// maker-note tag 0x0115 (int32u).
//
//	0x0 = Auto
//	0x1 = Color Temperature/Color Filter
//	0x10 = Daylight
//	0x20 = Cloudy
//	0x30 = Shade
//	0x40 = Tungsten
//	0x50 = Flash
//	0x60 = Fluorescent
//	0x70 = Custom
//	0x80 = Underwater
type WhiteBalanceValue uint32

const (
	WBAuto        WhiteBalanceValue = 0x00 // Auto
	WBColorTemp   WhiteBalanceValue = 0x01 // Color Temperature/Color Filter
	WBDaylight    WhiteBalanceValue = 0x10 // Daylight
	WBCloudy      WhiteBalanceValue = 0x20 // Cloudy
	WBShade       WhiteBalanceValue = 0x30 // Shade
	WBTungsten    WhiteBalanceValue = 0x40 // Tungsten
	WBFlash       WhiteBalanceValue = 0x50 // Flash
	WBFluorescent WhiteBalanceValue = 0x60 // Fluorescent
	WBCustom      WhiteBalanceValue = 0x70 // Custom
	WBUnderwater  WhiteBalanceValue = 0x80 // Underwater
)

// WhiteBalance0xB054Value describes the white-balance setting stored in Sony
// maker-note tag 0xb054 (int16u).
//
//	0 = Auto
//	4 = Custom
//	5 = Daylight
//	6 = Cloudy
//	7 = Cool White Fluorescent
//	8 = Day White Fluorescent
//	9 = Daylight Fluorescent
//	10 = Incandescent2
//	11 = Warm White Fluorescent
//	14 = Incandescent
//	15 = Flash
//	17 = Underwater 1 (Blue Water)
//	18 = Underwater 2 (Green Water)
type WhiteBalance0xB054Value uint16

const (
	WB54Auto           WhiteBalance0xB054Value = 0  // Auto
	WB54Custom         WhiteBalance0xB054Value = 4  // Custom
	WB54Daylight       WhiteBalance0xB054Value = 5  // Daylight
	WB54Cloudy         WhiteBalance0xB054Value = 6  // Cloudy
	WB54CoolWhiteFluor WhiteBalance0xB054Value = 7  // Cool White Fluorescent
	WB54DayWhiteFluor  WhiteBalance0xB054Value = 8  // Day White Fluorescent
	WB54DaylightFluor  WhiteBalance0xB054Value = 9  // Daylight Fluorescent
	WB54Incandescent2  WhiteBalance0xB054Value = 10 // Incandescent2
	WB54WarmWhiteFluor WhiteBalance0xB054Value = 11 // Warm White Fluorescent
	WB54Incandescent   WhiteBalance0xB054Value = 14 // Incandescent
	WB54Flash          WhiteBalance0xB054Value = 15 // Flash
	WB54Underwater1    WhiteBalance0xB054Value = 17 // Underwater 1 (Blue Water)
	WB54Underwater2    WhiteBalance0xB054Value = 18 // Underwater 2 (Green Water)
)

// QualityValue describes the JPEG/RAW quality stored in Sony maker-note tag
// 0x0102 (int32u).
//
//	0 = RAW
//	1 = Super Fine
//	2 = Fine
//	3 = Standard
//	4 = Economy
//	5 = Extra Fine
//	6 = RAW + JPEG/HEIF
//	7 = Compressed RAW
//	8 = Compressed RAW + JPEG
//	9 = Light
//	4294967295 = n/a
type QualityValue uint32

const (
	QualRAW               QualityValue = 0          // RAW
	QualSuperFine         QualityValue = 1          // Super Fine
	QualFine              QualityValue = 2          // Fine
	QualStandard          QualityValue = 3          // Standard
	QualEconomy           QualityValue = 4          // Economy
	QualExtraFine         QualityValue = 5          // Extra Fine
	QualRAWJPEG           QualityValue = 6          // RAW + JPEG/HEIF
	QualCompressedRAW     QualityValue = 7          // Compressed RAW
	QualCompressedRAWJPEG QualityValue = 8          // Compressed RAW + JPEG
	QualLight             QualityValue = 9          // Light
	QualNA                QualityValue = 0xFFFFFFFF // n/a
)

// SceneModeValue describes the Scene Mode stored in Sony maker-note tag 0xb023
// (int32u).
//
//	0 = Standard
//	1 = Portrait
//	2 = Text
//	3 = Night Scene
//	4 = Sunset
//	5 = Sports
//	6 = Landscape
//	7 = Night Portrait
//	8 = Macro
//	9 = Super Macro
//	16 = Auto
//	17 = Night View/Portrait
//	18 = Sweep Panorama
//	19 = Handheld Night Shot
//	20 = Anti Motion Blur
//	21 = Cont. Priority AE
//	22 = Auto+
//	23 = 3D Sweep Panorama
//	24 = Superior Auto
//	25 = High Sensitivity
//	26 = Fireworks
//	27 = Food
//	28 = Pet
//	33 = HDR
//	65535 = n/a
type SceneModeValue uint32

const (
	SceneStandard        SceneModeValue = 0     // Standard
	ScenePortrait        SceneModeValue = 1     // Portrait
	SceneText            SceneModeValue = 2     // Text
	SceneNight           SceneModeValue = 3     // Night Scene
	SceneSunset          SceneModeValue = 4     // Sunset
	SceneSports          SceneModeValue = 5     // Sports
	SceneLandscape       SceneModeValue = 6     // Landscape
	SceneNightPortrait   SceneModeValue = 7     // Night Portrait
	SceneMacro           SceneModeValue = 8     // Macro
	SceneSuperMacro      SceneModeValue = 9     // Super Macro
	SceneAuto            SceneModeValue = 16    // Auto
	SceneNightView       SceneModeValue = 17    // Night View/Portrait
	SceneSweepPanorama   SceneModeValue = 18    // Sweep Panorama
	SceneHandheldNight   SceneModeValue = 19    // Handheld Night Shot
	SceneAntiMotionBlur  SceneModeValue = 20    // Anti Motion Blur
	SceneContPriority    SceneModeValue = 21    // Cont. Priority AE
	SceneAutoPlus        SceneModeValue = 22    // Auto+
	Scene3DSweepPan      SceneModeValue = 23    // 3D Sweep Panorama
	SceneSuperiorAuto    SceneModeValue = 24    // Superior Auto
	SceneHighSensitivity SceneModeValue = 25    // High Sensitivity
	SceneFireworks       SceneModeValue = 26    // Fireworks
	SceneFood            SceneModeValue = 27    // Food
	ScenePet             SceneModeValue = 28    // Pet
	SceneHDR             SceneModeValue = 33    // HDR
	SceneNA              SceneModeValue = 65535 // n/a
)

// DynamicRangeOptimizerValue describes the DRO setting stored in Sony
// maker-note tag 0xb025 (int32u).
//
//	0 = Off
//	1 = Standard
//	2 = Advanced Auto
//	3 = Auto
//	8 = Advanced Lv1
//	9 = Advanced Lv2
//	10 = Advanced Lv3
//	11 = Advanced Lv4
//	12 = Advanced Lv5
//	16 = Lv1
//	17 = Lv2
//	18 = Lv3
//	19 = Lv4
//	20 = Lv5
type DynamicRangeOptimizerValue uint32

const (
	DROOff          DynamicRangeOptimizerValue = 0  // Off
	DROStandard     DynamicRangeOptimizerValue = 1  // Standard
	DROAdvancedAuto DynamicRangeOptimizerValue = 2  // Advanced Auto
	DROAuto         DynamicRangeOptimizerValue = 3  // Auto
	DROAdvancedLv1  DynamicRangeOptimizerValue = 8  // Advanced Lv1
	DROAdvancedLv2  DynamicRangeOptimizerValue = 9  // Advanced Lv2
	DROAdvancedLv3  DynamicRangeOptimizerValue = 10 // Advanced Lv3
	DROAdvancedLv4  DynamicRangeOptimizerValue = 11 // Advanced Lv4
	DROAdvancedLv5  DynamicRangeOptimizerValue = 12 // Advanced Lv5
	DROLv1          DynamicRangeOptimizerValue = 16 // Lv1
	DROLv2          DynamicRangeOptimizerValue = 17 // Lv2
	DROLv3          DynamicRangeOptimizerValue = 18 // Lv3
	DROLv4          DynamicRangeOptimizerValue = 19 // Lv4
	DROLv5          DynamicRangeOptimizerValue = 20 // Lv5
)

// FocusModeValue describes the focus mode stored in Sony maker-note tag 0x201b
// (int8u).
//
//	0 = Manual
//	2 = AF-S
//	3 = AF-C
//	4 = AF-A
//	6 = DMF
//	7 = AF-D
type FocusModeValue uint8

const (
	FocusManual FocusModeValue = 0 // Manual
	FocusAFS    FocusModeValue = 2 // AF-S
	FocusAFC    FocusModeValue = 3 // AF-C
	FocusAFA    FocusModeValue = 4 // AF-A
	FocusDMF    FocusModeValue = 6 // DMF
	FocusAFD    FocusModeValue = 7 // AF-D
)

// FocusMode0xB042Value describes the focus mode stored in Sony maker-note tag
// 0xb042 (int16u).
//
//	1 = AF-S
//	2 = AF-C
//	4 = Permanent-AF
//	65535 = n/a
type FocusMode0xB042Value uint16

const (
	FocusB042AFS       FocusMode0xB042Value = 1     // AF-S
	FocusB042AFC       FocusMode0xB042Value = 2     // AF-C
	FocusB042Permanent FocusMode0xB042Value = 4     // Permanent-AF
	FocusB042NA        FocusMode0xB042Value = 65535 // n/a
)

// AFAreaModeSettingValue describes the AF area mode stored in Sony maker-note
// tag 0x201c (int8u). Values vary by model family.
//
// Common:
//
//	0 = Wide
//	1 = Center
//	3 = Flexible Spot
//	4 = Local/Flexible Spot (LA-EA4)
//	8 = Zone
//	9 = Spot/Center (LA-EA4)
//	11 = Zone (NEX/ILCE)
//	12 = Expanded Flexible Spot
//	13 = Custom AF Area
type AFAreaModeSettingValue uint8

const (
	AFAWide          AFAreaModeSettingValue = 0  // Wide
	AFCenter         AFAreaModeSettingValue = 1  // Center
	AFFlexSpot       AFAreaModeSettingValue = 3  // Flexible Spot
	AFLocal          AFAreaModeSettingValue = 4  // Local
	AFZone           AFAreaModeSettingValue = 8  // Zone
	AFSpot           AFAreaModeSettingValue = 9  // Spot
	AFZoneNEX        AFAreaModeSettingValue = 11 // Zone (NEX/ILCE)
	AFExpandFlexSpot AFAreaModeSettingValue = 12 // Expanded Flexible Spot
	AFCustomAFArea   AFAreaModeSettingValue = 13 // Custom AF Area
)

// MeteringModeValue describes the metering mode stored in Sony maker-note tag
// 0x202c (int16u).
//
//	0x100 = Multi-segment
//	0x200 = Center-weighted average
//	0x301 = Spot (Standard)
//	0x302 = Spot (Large)
//	0x400 = Average
//	0x500 = Highlight
type MeteringModeValue uint16

const (
	MeterMultiSegment   MeteringModeValue = 0x100 // Multi-segment
	MeterCenterWeighted MeteringModeValue = 0x200 // Center-weighted average
	MeterSpotStandard   MeteringModeValue = 0x301 // Spot (Standard)
	MeterSpotLarge      MeteringModeValue = 0x302 // Spot (Large)
	MeterAverage        MeteringModeValue = 0x400 // Average
	MeterHighlight      MeteringModeValue = 0x500 // Highlight
)

// FlashModeValue describes the flash mode stored in the Sony CameraSettings3
// table (int8u).
//
//	1 = Flash Off
//	16 = Autoflash
//	17 = Fill-flash
//	18 = Slow Sync
//	19 = Rear Sync
//	20 = Wireless
type FlashModeValue uint8

const (
	FlashOff      FlashModeValue = 1  // Flash Off
	FlashAuto     FlashModeValue = 16 // Autoflash
	FlashFill     FlashModeValue = 17 // Fill-flash
	FlashSlow     FlashModeValue = 18 // Slow Sync
	FlashRear     FlashModeValue = 19 // Rear Sync
	FlashWireless FlashModeValue = 20 // Wireless
)

// ReleaseModeValue describes the release/burst mode stored in Sony maker-note
// tag 0xb049 (int16u).
//
//	0 = Normal
//	2 = Continuous
//	5 = Exposure Bracketing
//	6 = White Balance Bracketing
//	8 = DRO Bracketing
//	65535 = n/a
type ReleaseModeValue uint16

const (
	ReleaseNormal     ReleaseModeValue = 0     // Normal
	ReleaseContinuous ReleaseModeValue = 2     // Continuous
	ReleaseExpBkt     ReleaseModeValue = 5     // Exposure Bracketing
	ReleaseWBBkt      ReleaseModeValue = 6     // White Balance Bracketing
	ReleaseDROBkt     ReleaseModeValue = 8     // DRO Bracketing
	ReleaseNA         ReleaseModeValue = 65535 // n/a
)

// PictureEffectValue describes the picture effect stored in Sony maker-note
// tag 0x200e (int16u).
//
//	0 = Off
//	1 = Toy Camera
//	2 = Pop Color
//	3 = Posterization
//	4 = Posterization B/W
//	5 = Retro Photo
//	6 = Soft High Key
//	7 = Partial Color (red)
//	8 = Partial Color (green)
//	9 = Partial Color (blue)
//	10 = Partial Color (yellow)
//	13 = High Contrast Monochrome
//	16 = Toy Camera (normal)
//	17 = Toy Camera (cool)
//	18 = Toy Camera (warm)
//	19 = Toy Camera (green)
//	20 = Toy Camera (magenta)
//	32 = Soft Focus (low)
//	33 = Soft Focus
//	34 = Soft Focus (high)
//	48 = Miniature (auto)
//	49 = Miniature (top)
//	50 = Miniature (middle horizontal)
//	51 = Miniature (bottom)
//	52 = Miniature (left)
//	53 = Miniature (middle vertical)
//	54 = Miniature (right)
//	64 = HDR Painting (low)
//	65 = HDR Painting
//	66 = HDR Painting (high)
//	80 = Rich-tone Monochrome
//	97 = Water Color
//	98 = Water Color 2
//	112 = Illustration (low)
//	113 = Illustration
//	114 = Illustration (high)
type PictureEffectValue uint16

const (
	PEOff          PictureEffectValue = 0   // Off
	PEToyCamera    PictureEffectValue = 1   // Toy Camera
	PEPopColor     PictureEffectValue = 2   // Pop Color
	PEPoster       PictureEffectValue = 3   // Posterization
	PEPosterBW     PictureEffectValue = 4   // Posterization B/W
	PERetro        PictureEffectValue = 5   // Retro Photo
	PESoftHighKey  PictureEffectValue = 6   // Soft High Key
	PERed          PictureEffectValue = 7   // Partial Color (red)
	PEGreen        PictureEffectValue = 8   // Partial Color (green)
	PEBlue         PictureEffectValue = 9   // Partial Color (blue)
	PEYellow       PictureEffectValue = 10  // Partial Color (yellow)
	PEHighContBW   PictureEffectValue = 13  // High Contrast Monochrome
	PEToyNormal    PictureEffectValue = 16  // Toy Camera (normal)
	PEToyCool      PictureEffectValue = 17  // Toy Camera (cool)
	PEToyWarm      PictureEffectValue = 18  // Toy Camera (warm)
	PEToyGreen     PictureEffectValue = 19  // Toy Camera (green)
	PEToyMagenta   PictureEffectValue = 20  // Toy Camera (magenta)
	PESoftLow      PictureEffectValue = 32  // Soft Focus (low)
	PESoftFocus    PictureEffectValue = 33  // Soft Focus
	PESoftHigh     PictureEffectValue = 34  // Soft Focus (high)
	PEMiniAuto     PictureEffectValue = 48  // Miniature (auto)
	PEMiniTop      PictureEffectValue = 49  // Miniature (top)
	PEMiniMidH     PictureEffectValue = 50  // Miniature (middle horizontal)
	PEMiniBottom   PictureEffectValue = 51  // Miniature (bottom)
	PEMiniLeft     PictureEffectValue = 52  // Miniature (left)
	PEMiniMidV     PictureEffectValue = 53  // Miniature (middle vertical)
	PEMiniRight    PictureEffectValue = 54  // Miniature (right)
	PEHDRPaintLow  PictureEffectValue = 64  // HDR Painting (low)
	PEHDRPaint     PictureEffectValue = 65  // HDR Painting
	PEHDRPaintHigh PictureEffectValue = 66  // HDR Painting (high)
	PERichTone     PictureEffectValue = 80  // Rich-tone Monochrome
	PEWaterColor   PictureEffectValue = 97  // Water Color
	PEWaterColor2  PictureEffectValue = 98  // Water Color 2
	PEIllustLow    PictureEffectValue = 112 // Illustration (low)
	PEIllustrate   PictureEffectValue = 113 // Illustration
	PEIllustHigh   PictureEffectValue = 114 // Illustration (high)
)

// LongExposureNoiseReductionValue describes the Long Exposure NR setting
// stored in Sony maker-note tag 0x2008 (int32u).
//
//	0x0 = Off
//	0x1 = On (unused)
//	0x10001 = On (dark subtracted)
//	0xffff0000 = Off (65535)
//	0xffff0001 = On (65535)
//	0xffffffff = n/a
type LongExposureNoiseReductionValue uint32

const (
	LENROff       LongExposureNoiseReductionValue = 0x0        // Off
	LENROnUnused  LongExposureNoiseReductionValue = 0x1        // On (unused)
	LENROnDarkSub LongExposureNoiseReductionValue = 0x10001    // On (dark subtracted)
	LENROff65535  LongExposureNoiseReductionValue = 0xFFFF0000 // Off (65535)
	LENROn65535   LongExposureNoiseReductionValue = 0xFFFF0001 // On (65535)
	LENRNA        LongExposureNoiseReductionValue = 0xFFFFFFFF // n/a
)

// HighISONoiseReductionValue describes the High ISO NR setting stored in Sony
// maker-note tag 0x2009 (int16u).
//
//	0 = Off
//	1 = Low
//	2 = Normal
//	3 = High
//	256 = Auto
//	65535 = n/a
type HighISONoiseReductionValue uint16

const (
	HNROff    HighISONoiseReductionValue = 0     // Off
	HNRLow    HighISONoiseReductionValue = 1     // Low
	HNRNormal HighISONoiseReductionValue = 2     // Normal
	HNRHigh   HighISONoiseReductionValue = 3     // High
	HNRAuto   HighISONoiseReductionValue = 256   // Auto
	HNRNA     HighISONoiseReductionValue = 65535 // n/a
)

// AntiBlurValue describes the anti-blur setting stored in Sony maker-note tag
// 0xb04b (int16u).
//
//	0 = Off
//	1 = On (Continuous)
//	2 = On (Shooting)
//	65535 = n/a
type AntiBlurValue uint16

const (
	AntiBlurOff        AntiBlurValue = 0     // Off
	AntiBlurContinuous AntiBlurValue = 1     // On (Continuous)
	AntiBlurShooting   AntiBlurValue = 2     // On (Shooting)
	AntiBlurNA         AntiBlurValue = 65535 // n/a
)

// IntelligentAutoValue describes the Intelligent Auto setting stored in Sony
// maker-note tag 0xb052 (int16u).
//
//	0 = Off
//	1 = On
//	2 = Advanced
type IntelligentAutoValue uint16

const (
	IAOff      IntelligentAutoValue = 0 // Off
	IAOn       IntelligentAutoValue = 1 // On
	IAAdvanced IntelligentAutoValue = 2 // Advanced
)

// ExposureModeValue describes the exposure mode stored in Sony maker-note tag
// 0xb041 (int16u).
//
//	0 = Program AE
//	1 = Portrait
//	2 = Beach
//	3 = Sports
//	4 = Snow
//	5 = Landscape
//	6 = Auto
//	7 = Aperture-priority AE
//	8 = Shutter speed priority AE
//	9 = Night Scene / Twilight
//	10 = Hi-Speed Shutter
//	11 = Twilight Portrait
//	12 = Soft Snap/Portrait
//	13 = Fireworks
//	14 = Smile Shutter
//	15 = Manual
//	18 = High Sensitivity
//	19 = Macro
//	20 = Advanced Sports Shooting
//	29 = Underwater
//	33 = Food
//	34 = Sweep Panorama
//	35 = Handheld Night Shot
//	36 = Anti Motion Blur
//	37 = Pet
//	38 = Backlight Correction HDR
//	39 = Superior Auto
//	40 = Background Defocus
//	41 = Soft Skin
//	42 = 3D Image
//	65535 = n/a
type ExposureModeValue uint16

const (
	ExpProgramAE        ExposureModeValue = 0     // Program AE
	ExpPortrait         ExposureModeValue = 1     // Portrait
	ExpBeach            ExposureModeValue = 2     // Beach
	ExpSports           ExposureModeValue = 3     // Sports
	ExpSnow             ExposureModeValue = 4     // Snow
	ExpLandscape        ExposureModeValue = 5     // Landscape
	ExpAuto             ExposureModeValue = 6     // Auto
	ExpApertureAE       ExposureModeValue = 7     // Aperture-priority AE
	ExpShutterAE        ExposureModeValue = 8     // Shutter speed priority AE
	ExpNightScene       ExposureModeValue = 9     // Night Scene / Twilight
	ExpHiSpeedShutter   ExposureModeValue = 10    // Hi-Speed Shutter
	ExpTwilightPortrait ExposureModeValue = 11    // Twilight Portrait
	ExpSoftSnap         ExposureModeValue = 12    // Soft Snap/Portrait
	ExpFireworks        ExposureModeValue = 13    // Fireworks
	ExpSmileShutter     ExposureModeValue = 14    // Smile Shutter
	ExpManual           ExposureModeValue = 15    // Manual
	ExpHighSensitivity  ExposureModeValue = 18    // High Sensitivity
	ExpMacro            ExposureModeValue = 19    // Macro
	ExpAdvSports        ExposureModeValue = 20    // Advanced Sports Shooting
	ExpUnderwater       ExposureModeValue = 29    // Underwater
	ExpFood             ExposureModeValue = 33    // Food
	ExpSweepPanorama    ExposureModeValue = 34    // Sweep Panorama
	ExpHandheldNight    ExposureModeValue = 35    // Handheld Night Shot
	ExpAntiMotionBlur   ExposureModeValue = 36    // Anti Motion Blur
	ExpPet              ExposureModeValue = 37    // Pet
	ExpBacklightHDR     ExposureModeValue = 38    // Backlight Correction HDR
	ExpSuperiorAuto     ExposureModeValue = 39    // Superior Auto
	ExpBgDefocus        ExposureModeValue = 40    // Background Defocus
	ExpSoftSkin         ExposureModeValue = 41    // Soft Skin
	Exp3DImage          ExposureModeValue = 42    // 3D Image
	ExpNA               ExposureModeValue = 65535 // n/a
)

// ColorModeValue describes the color mode stored in Sony maker-note tag
// 0xb029 (int32u).
//
//	0 = Standard
//	1 = Vivid
//	2 = Portrait
//	3 = Landscape
//	4 = Sunset
//	5 = Night View/Portrait
//	6 = B&W
//	7 = Adobe RGB
//	12 = Neutral
//	13 = Clear
//	14 = Deep
//	15 = Light
//	16 = Autumn Leaves
//	17 = Sepia
//	18 = FL
//	19 = Vivid 2
//	20 = IN
//	21 = SH
//	22 = FL2
//	23 = FL3
//	100 = Neutral
//	101 = Clear
//	102 = Deep
//	103 = Light
//	104 = Night View
//	105 = Autumn Leaves
//	255 = Off
//	4294967295 = n/a
type ColorModeValue uint32

const (
	CMStandard        ColorModeValue = 0          // Standard
	CMVivid           ColorModeValue = 1          // Vivid
	CMPortrait        ColorModeValue = 2          // Portrait
	CMLandscape       ColorModeValue = 3          // Landscape
	CMSunset          ColorModeValue = 4          // Sunset
	CMNightView       ColorModeValue = 5          // Night View/Portrait
	CMBAndW           ColorModeValue = 6          // B&W
	CMAdobeRGB        ColorModeValue = 7          // Adobe RGB
	CMNeutral         ColorModeValue = 12         // Neutral
	CMClear           ColorModeValue = 13         // Clear
	CMDeep            ColorModeValue = 14         // Deep
	CMLight           ColorModeValue = 15         // Light
	CMAutumnLeaves    ColorModeValue = 16         // Autumn Leaves
	CMSepia           ColorModeValue = 17         // Sepia
	CMFL              ColorModeValue = 18         // FL
	CMVivid2          ColorModeValue = 19         // Vivid 2
	CMIN              ColorModeValue = 20         // IN
	CMSH              ColorModeValue = 21         // SH
	CMFL2             ColorModeValue = 22         // FL2
	CMFL3             ColorModeValue = 23         // FL3
	CMNeutral100      ColorModeValue = 100        // Neutral
	CMClear100        ColorModeValue = 101        // Clear
	CMDeep100         ColorModeValue = 102        // Deep
	CMLight100        ColorModeValue = 103        // Light
	CMNightView100    ColorModeValue = 104        // Night View
	CMAutumnLeaves100 ColorModeValue = 105        // Autumn Leaves
	CMOff             ColorModeValue = 255        // Off
	CMNA              ColorModeValue = 0xFFFFFFFF // n/a
)

// AFAreaModeValue describes the AF area mode stored in Sony maker-note tag
// 0xb043 (int16u). Values vary by model generation.
//
// Older models:
//
//	0 = Default
//	1 = Multi
//	2 = Center
//	3 = Spot
//	4 = Flexible Spot
//	6 = Touch
//	14 = Tracking
//	15 = Face Tracking
//	65535 = n/a
type AFAreaModeValue uint16

const (
	AFAreaDefault   AFAreaModeValue = 0     // Default
	AFAreaMulti     AFAreaModeValue = 1     // Multi
	AFAreaCenter    AFAreaModeValue = 2     // Center
	AFAreaSpot      AFAreaModeValue = 3     // Spot
	AFAreaFlexSpot  AFAreaModeValue = 4     // Flexible Spot
	AFAreaTouch     AFAreaModeValue = 6     // Touch
	AFAreaTrack     AFAreaModeValue = 14    // Tracking
	AFAreaFaceTrack AFAreaModeValue = 15    // Face Tracking
	AFAreaNA        AFAreaModeValue = 65535 // n/a
)

// ZoneMatchingValue describes the zone matching stored in Sony maker-note tag
// 0xb024 (int32u).
//
//	0 = ISO Setting Used
//	1 = High Key
//	2 = Low Key
type ZoneMatchingValue uint32

const (
	ZoneISO  ZoneMatchingValue = 0 // ISO Setting Used
	ZoneHigh ZoneMatchingValue = 1 // High Key
	ZoneLow  ZoneMatchingValue = 2 // Low Key
)

// ImageStabilizationValue describes the image stabilization setting stored in
// Sony maker-note tag 0xb026 (int32u).
//
//	0 = Off
//	1 = On
//	4294967295 = n/a
type ImageStabilizationValue uint32

const (
	ISOff ImageStabilizationValue = 0          // Off
	ISOn  ImageStabilizationValue = 1          // On
	ISNA  ImageStabilizationValue = 0xFFFFFFFF // n/a
)

// SoftSkinEffectValue describes the soft skin effect setting stored in Sony
// maker-note tag 0x200f (int32u).
//
//	0 = Off
//	1 = Low
//	2 = Mid
//	3 = High
//	4294967295 = n/a
type SoftSkinEffectValue uint32

const (
	SSOff  SoftSkinEffectValue = 0          // Off
	SSLow  SoftSkinEffectValue = 1          // Low
	SSMid  SoftSkinEffectValue = 2          // Mid
	SSHigh SoftSkinEffectValue = 3          // High
	SSNA   SoftSkinEffectValue = 0xFFFFFFFF // n/a
)

// MacroValue describes the macro mode stored in Sony maker-note tag 0xb040
// (int16u).
//
//	0 = Off
//	1 = On
//	2 = Close Focus
//	65535 = n/a
type MacroValue uint16

const (
	MacroOff        MacroValue = 0     // Off
	MacroOn         MacroValue = 1     // On
	MacroCloseFocus MacroValue = 2     // Close Focus
	MacroNA         MacroValue = 65535 // n/a
)

// AFIlluminatorValue describes the AF illuminator setting stored in Sony
// maker-note tag 0xb044 (int16u).
//
//	0 = Off
//	1 = Auto
//	65535 = n/a
type AFIlluminatorValue uint16

const (
	AFIllumOff  AFIlluminatorValue = 0     // Off
	AFIllumAuto AFIlluminatorValue = 1     // Auto
	AFIllumNA   AFIlluminatorValue = 65535 // n/a
)

// JPEGQualityValue describes the JPEG quality stored in Sony maker-note tag
// 0xb047 (int16u).
//
//	0 = Standard
//	1 = Fine
//	2 = Extra Fine
//	65535 = n/a
type JPEGQualityValue uint16

const (
	JpegStandard  JPEGQualityValue = 0     // Standard
	JpegFine      JPEGQualityValue = 1     // Fine
	JpegExtraFine JPEGQualityValue = 2     // Extra Fine
	JpegNA        JPEGQualityValue = 65535 // n/a
)

// RAWFileTypeValue describes the RAW file type stored in Sony maker-note tag
// 0x2029 (int16u).
//
//	0 = Compressed RAW
//	1 = Uncompressed RAW
//	2 = Lossless Compressed RAW
//	3 = Compressed RAW 2
//	65535 = n/a
type RAWFileTypeValue uint16

const (
	RAWCompressed   RAWFileTypeValue = 0     // Compressed RAW
	RAWUncompressed RAWFileTypeValue = 1     // Uncompressed RAW
	RAWLosslessComp RAWFileTypeValue = 2     // Lossless Compressed RAW
	RAWCompressed2  RAWFileTypeValue = 3     // Compressed RAW 2
	RAWFileTypeNA   RAWFileTypeValue = 65535 // n/a
)

// PrioritySetInAWBValue describes the Priority Set in AWB setting stored in
// Sony maker-note tag 0x202b (int8u).
//
//	0 = Standard
//	1 = Ambience
//	2 = White
type PrioritySetInAWBValue uint8

const (
	AWBPriorityStandard PrioritySetInAWBValue = 0 // Standard
	AWBPriorityAmbience PrioritySetInAWBValue = 1 // Ambience
	AWBPriorityWhite    PrioritySetInAWBValue = 2 // White
)

// FlashActionValue describes the flash action stored in Sony maker-note tag
// 0x2017 (int32u).
//
//	0 = Did not fire
//	1 = Flash Fired
//	2 = External Flash Fired
//	3 = Wireless Controlled Flash Fired
type FlashActionValue uint32

const (
	FlashDidNotFire   FlashActionValue = 0 // Did not fire
	FlashFired        FlashActionValue = 1 // Flash Fired
	FlashExternal     FlashActionValue = 2 // External Flash Fired
	FlashWirelessCtrl FlashActionValue = 3 // Wireless Controlled Flash Fired
)

// AFTrackingValue describes the AF tracking mode stored in Sony maker-note tag
// 0x2021 (int8u).
//
//	0 = Off
//	1 = Face tracking
//	2 = Lock On AF
type AFTrackingValue uint8

const (
	AFTrackOff    AFTrackingValue = 0 // Off
	AFTrackFace   AFTrackingValue = 1 // Face tracking
	AFTrackLockOn AFTrackingValue = 2 // Lock On AF
)

// MultiFrameNREffectValue describes the Multi Frame NR Effect stored in Sony
// maker-note tag 0x2023 (int32u).
//
//	0 = Normal
//	1 = High
type MultiFrameNREffectValue uint32

const (
	MFNRNormal MultiFrameNREffectValue = 0 // Normal
	MFNRHigh   MultiFrameNREffectValue = 1 // High
)

// HDRValue describes the HDR setting stored in Sony maker-note tag 0x200a
// (int32u, stored as two int16u values).
//
// [Value 0] - HDR Auto/EV:
//
//	0x0 = Off
//	0x1 = Auto
//	0x10 = 1.0 EV
//	0x11 = 1.5 EV
//	0x12 = 2.0 EV
//	0x13 = 2.5 EV
//	0x14 = 3.0 EV
//	0x15 = 3.5 EV
//	0x16 = 4.0 EV
//	0x17 = 4.5 EV
//	0x18 = 5.0 EV
//	0x19 = 5.5 EV
//	0x1a = 6.0 EV
//
// [Value 1] - HDR Image State:
//
//	0x0 = Uncorrected image
//	0x1 = HDR image (good)
//	0x2 = HDR image (fail 1)
//	0x3 = HDR image (fail 2)
type HDRValue uint32

const (
	HDROff   HDRValue = 0x0  // Off
	HDRAuto  HDRValue = 0x1  // Auto
	HDR1_0EV HDRValue = 0x10 // 1.0 EV
	HDR1_5EV HDRValue = 0x11 // 1.5 EV
	HDR2_0EV HDRValue = 0x12 // 2.0 EV
	HDR2_5EV HDRValue = 0x13 // 2.5 EV
	HDR3_0EV HDRValue = 0x14 // 3.0 EV
	HDR3_5EV HDRValue = 0x15 // 3.5 EV
	HDR4_0EV HDRValue = 0x16 // 4.0 EV
	HDR4_5EV HDRValue = 0x17 // 4.5 EV
	HDR5_0EV HDRValue = 0x18 // 5.0 EV
	HDR5_5EV HDRValue = 0x19 // 5.5 EV
	HDR6_0EV HDRValue = 0x1a // 6.0 EV
)

// CreativeStyleValue describes the Creative Style setting stored in Sony
// maker-note CameraSettings3 (uint8), CameraSettings (uint16), and Tag9416
// (uint8) tables.
//
// The raw numeric values are not linear: multiple codes map to the same
// display name.  Name() normalises the known aliases.
//
//	0, 1,  16: "Standard"
//	2,    32: "Vivid"
//	3,    64: "Portrait"
//	4,    80: "Landscape"
//	5, 96, 160: "Sunset"
//	  6: "Night View/Portrait"
//	  8: "B&W"
//	  9: "Adobe RGB"
//	 11: "Neutral"
//	 12: "Clear"
//	 13: "Deep"
//	 14: "Light"
//	 15: "Autumn Leaves"
//	255: "Off"
type CreativeStyleValue uint16

const (
	CreativeStandard     CreativeStyleValue = 0   // Standard
	CreativeVivid        CreativeStyleValue = 2   // Vivid
	CreativePortrait     CreativeStyleValue = 3   // Portrait
	CreativeLandscape    CreativeStyleValue = 4   // Landscape
	CreativeSunset       CreativeStyleValue = 5   // Sunset
	CreativeNightView    CreativeStyleValue = 6   // Night View/Portrait
	CreativeBAndW        CreativeStyleValue = 8   // B&W
	CreativeAdobeRGB     CreativeStyleValue = 9   // Adobe RGB
	CreativeNeutral      CreativeStyleValue = 11  // Neutral
	CreativeClear        CreativeStyleValue = 12  // Clear
	CreativeDeep         CreativeStyleValue = 13  // Deep
	CreativeLight        CreativeStyleValue = 14  // Light
	CreativeAutumnLeaves CreativeStyleValue = 15  // Autumn Leaves
	CreativeOff          CreativeStyleValue = 255 // Off
)

// Name returns the display name for the creative style, handling the alias
// values (1→Standard, 16→Standard, 32→Vivid, 64→Portrait, 80→Landscape,
// 96/160→Sunset).  Returns empty string when the value is unrecognized.
func (c CreativeStyleValue) Name() string {
	switch c {
	case 0, 1, 16:
		return "Standard"
	case 2, 32:
		return "Vivid"
	case 3, 64:
		return "Portrait"
	case 4, 80:
		return "Landscape"
	case 5, 96, 160:
		return "Sunset"
	case 6:
		return "Night View/Portrait"
	case 8:
		return "B&W"
	case 9:
		return "Adobe RGB"
	case 11:
		return "Neutral"
	case 12:
		return "Clear"
	case 13:
		return "Deep"
	case 14:
		return "Light"
	case 15:
		return "Autumn Leaves"
	case 255:
		return "Off"
	default:
		return ""
	}
}

// UsesCameraInfo3 reports whether the camera model uses the CameraInfo3
// layout for maker-note tag 0x0010.
func UsesCameraInfo3(model string) bool {
	switch {
	case strings.HasPrefix(model, "SLT-"),
		strings.HasPrefix(model, "NEX-"),
		strings.HasPrefix(model, "ILC"),
		strings.HasPrefix(model, "ILME-"),
		strings.HasPrefix(model, "ZV-"):
		return true
	default:
		return false
	}
}

// Sony contains the selected Sony maker-note fields currently decoded by
// imagemeta.
//
// The field set mirrors the subset of ExifTool's Image::ExifTool::Sony::Main,
// Sony::CameraSettings3, and Sony::Tag9050 tables that imagemeta parses today.
//
// Top-level fields are the most commonly accessed metadata. Sub-directory
// tables referenced by ExifTool are kept in embedded sub-structs for
// forward-compatibility and to preserve the full decoded data.
type Sony struct {
	Rating                      uint32    // 0x2002 Rating (int32u)
	Contrast                    int32     // 0x2004 Contrast (int32s)
	Saturation                  int32     // 0x2005 Saturation (int32s)
	Sharpness                   int32     // 0x2006 Sharpness (int32s)
	CreativeStyle               string    // 0xb020 CreativeStyle (string)
	DynamicRangeOptimizer       uint32    // 0xb025 DynamicRangeOptimizer (int32u)
	ImageStabilization          uint32    // 0xb026 ImageStabilization (int32u)
	ColorMode                   uint32    // 0xb029 ColorMode (int32u)
	Quality                     uint32    // 0x0102 Quality (int32u)
	Quality2                    [2]uint16 // 0x202e Quality2 (int16u[2])
	WhiteBalance                uint32    // 0x0115 WhiteBalance (int32u)
	WhiteBalanceFineTune        int32     // 0x0112 WhiteBalanceFineTune (int32s)
	FlashExposureComp           float64   // 0x0104 FlashExposureComp (rational64s)
	Teleconverter               uint32    // 0x0105 Teleconverter (int32u)
	SonyModelID                 uint16    // 0xb001 SonyModelID (int16u)
	LensType                    uint32    // 0xb027 LensType (int32u)
	LensSpec                    string    // 0xb02a LensSpec (undef[8])
	FileFormat                  [4]uint8  // 0xb000 FileFormat (int8u[4])
	ColorTemperature            uint32    // 0xb021 ColorTemperature (int32u)
	ColorCompensationFilter     int32     // 0xb022 ColorCompensationFilter (int32s)
	SceneMode                   uint32    // 0xb023 SceneMode (int32u)
	ZoneMatching                uint32    // 0xb024 ZoneMatching (int32u)
	FullImageSize               [2]uint32 // 0xb02b FullImageSize (int32u[2])
	PreviewImageSize            [2]uint32 // 0xb02c PreviewImageSize (int32u[2])
	ExposureMode                uint16    // 0xb041 ExposureMode (int16u)
	FocusMode0xB042             uint16    // 0xb042 FocusMode (int16u)
	AFAreaMode                  uint16    // 0xb043 AFAreaMode (int16u)
	AFIlluminator               uint16    // 0xb044 AFIlluminator (int16u)
	JPEGQuality                 uint16    // 0xb047 JPEGQuality (int16u)
	FlashLevel                  int16     // 0xb048 FlashLevel (int16s)
	ReleaseMode                 uint16    // 0xb049 ReleaseMode (int16u)
	SequenceNumber              uint16    // 0xb04a SequenceNumber (int16u)
	AntiBlur                    uint16    // 0xb04b Anti-Blur (int16u)
	AFTracking                  uint8     // 0x2021 AFTracking (int8u)
	DynamicRangeOptimizer0xB04F uint16    // 0xb04f DynamicRangeOptimizer (int16u)
	HighISONoiseReduction2      uint16    // 0xb050 HighISONoiseReduction2 (int16u)
	IntelligentAuto             uint16    // 0xb052 IntelligentAuto (int16u)
	WhiteBalance0xB054          uint16    // 0xb054 WhiteBalance (int16u)

	Brightness                    int32  // 0x2007 Brightness (int32s)
	LongExposureNoiseReduction    uint32 // 0x2008 LongExposureNoiseReduction (int32u)
	HighISONoiseReduction         uint16 // 0x2009 HighISONoiseReduction (int16u)
	HDR                           uint32 // 0x200a HDR (int32u)
	MultiFrameNoiseReduction      uint32 // 0x200b MultiFrameNoiseReduction (int32u)
	PictureEffect                 uint16 // 0x200e PictureEffect (int16u)
	SoftSkinEffect                uint32 // 0x200f SoftSkinEffect (int32u)
	VignettingCorrection          uint32 // 0x2011 VignettingCorrection (int32u)
	LateralChromaticAberration    uint32 // 0x2012 LateralChromaticAberration (int32u)
	DistortionCorrectionSetting   uint32 // 0x2013 DistortionCorrectionSetting (int32u)
	AutoPortraitFramed            uint16 // 0x2016 AutoPortraitFramed (int16u)
	FlashAction                   uint32 // 0x2017 FlashAction (int32u)
	ElectronicFrontCurtainShutter uint32 // 0x201a ElectronicFrontCurtainShutter (int32u)
	FocusMode                     uint8  // 0x201b FocusMode (int8u)
	AFAreaModeSetting             uint8  // 0x201c AFAreaModeSetting (int8u)
	AFPointSelected               uint8  // 0x201e AFPointSelected (int8u)
	MultiFrameNREffect            uint32 // 0x2023 MultiFrameNREffect (int32u)
	RAWFileType                   uint16 // 0x2029 RAWFileType (int16u)
	PrioritySetInAWB              uint8  // 0x202b PrioritySetInAWB (int8u)
	MeteringMode2                 uint16 // 0x202c MeteringMode2 (int16u)
	Macro                         uint16 // 0xb040 Macro (int16u)
	FocusMode0xB04E               uint16 // 0xb04e FocusMode (int16u)

	// Image processing controls (0x2014-0x2036).
	WBShiftABGM                [2]int32  // 0x2014 WBShiftAB_GM (int32s[2])
	FlexibleSpotPosition       [2]uint16 // 0x201d FlexibleSpotPosition (int16u[2])
	WBShiftABGMPrecise         [2]int32  // 0x2026 WBShiftAB_GM_Precise (int32s[2])
	FocusLocation              [4]uint16 // 0x2027 FocusLocation (int16u[4])
	VariableLowPassFilter      [2]uint16 // 0x2028 VariableLowPassFilter (int16u[2])
	ExposureStandardAdjustment float64   // 0x202d ExposureStandardAdjustment (rational64s)
	SerialNumber               string    // 0x2031 SerialNumber (string)
	Shadows                    int32     // 0x2032 Shadows (int32s)
	Highlights                 int32     // 0x2033 Highlights (int32s)
	Fade                       int32     // 0x2034 Fade (int32s)
	SharpnessRange             int32     // 0x2035 SharpnessRange (int32s)
	Clarity                    int32     // 0x2036 Clarity (int32s)
	FocusLocation2             [4]uint16 // 0x204a FocusLocation2 (int16u[4])

	// Embedded Sony maker-note sub-directory tables.

	CameraInfo2     SonyCameraInfo2     // 0x0010 CameraInfo (legacy)
	CameraInfo3     SonyCameraInfo3     // 0x0010 CameraInfo
	FocusInfo       SonyFocusInfo       // 0x0020 FocusInfo (legacy)
	MoreInfo        SonyMoreInfo        // 0x0020 MoreInfo
	CameraSettings  SonyCameraSettings  // 0x0114 CameraSettings (legacy)
	CameraSettings3 SonyCameraSettings3 // 0x0114 CameraSettings3
	ShotInfo        SonyShotInfo        // 0x3000 ShotInfo
	Tag9400         SonyTag9400         // 0x9400 Tag9400
	Tag9404         SonyTag9404         // 0x9404 Tag9404
	Tag9405         SonyTag9405         // 0x9405 Tag9405
	Tag9406         SonyTag9406         // 0x9406 Tag9406
	Tag940A         SonyTag940A         // 0x940a Tag940a
	Tag940C         SonyTag940C         // 0x940c Tag940c
	Tag2010         SonyTag2010         // 0x2010 Tag2010
	Tag202A         SonyTag202A         // 0x202a Tag202a
	HiddenInfo      SonyHiddenInfo      // 0x2044 HiddenInfo
	Tag9050         SonyTag9050         // 0x9050 Tag9050
	Tag9416         SonyTag9416         // 0x9416 Tag9416
	AFInfo          SonyAFInfo          // 0x940e AFInfo
}
