// Package canon provides data types and functions for representing Canon Camera Makernote values
package canon

import "time"

//go:generate msgp
//go:generate stringer -type=MacroMode,Quality,CanonFlashMode,ContinuousDrive,FocusMode,RecordMode,CanonImageSize,EasyMode,DigitalZoom,CameraISO,MeteringMode,FocusRange,ExposureMode,FlashModel,FocusContinuous,AESetting,ImageStabilization,SpotMeteringMode,PhotoEffect,ManualFlashOutput,SRAWQuality,FocusBracketing,HDRPQ,BracketMode,OnOffAuto,FilterEffect,ToningEffect,ShutterMode,RawJpgQuality,RawJpgSize,TimeZoneCity,DaylightSavings,AFAreaMode -linecomment -output=canon_string.go

// ContinuousDrive is part of the CanonCameraSettings field
//
//	0:  "Single",
//	1:  "Continuous",
//	2:  "Movie",
//	3:  "Continuous, Speed Priority",
//	4:  "Continuous, Low",
//	5:  "Continuous, High",
//	6:  "Silent Single",
//	7:  "Continuous",
//	8:  "Single",
//	9:  "Single, Silent",
//	10: "Continuous, Silent",
type ContinuousDrive int16

const (
	ContinuousDriveSingle                  ContinuousDrive = 0  // Single
	ContinuousDriveContinuous              ContinuousDrive = 1  // Continuous
	ContinuousDriveMovie                   ContinuousDrive = 2  // Movie
	ContinuousDriveContinuousSpeedPriority ContinuousDrive = 3  // Continuous, Speed Priority
	ContinuousDriveContinuousLow           ContinuousDrive = 4  // Continuous, Low
	ContinuousDriveContinuousHigh          ContinuousDrive = 5  // Continuous, High
	ContinuousDriveSilentSingle            ContinuousDrive = 6  // Silent Single
	ContinuousDriveContinuousAlt           ContinuousDrive = 7  // Continuous
	ContinuousDriveSingleAlt               ContinuousDrive = 8  // Single
	ContinuousDriveSingleSilent            ContinuousDrive = 9  // Single, Silent
	ContinuousDriveContinuousSilent        ContinuousDrive = 10 // Continuous, Silent
	ContinuousDriveUnknown7                ContinuousDrive = 7  // Continuous
	ContinuousDriveUnknown8                ContinuousDrive = 8  // Single
)

// FocusMode is part of the CanonCameraSettings field
//
//	0:   "One-shot AF",
//	1:   "AI Servo AF",
//	2:   "AI Focus AF",
//	3:   "Manual Focus",
//	4:   "Single",
//	5:   "Continuous",
//	6:   "Manual Focus",
//	16:  "Pan Focus",
//	256: "AF + MF",
//	258: "MF",
//	512: "Movie Snap Focus",
//	519: "Movie Servo AF",
type FocusMode int16

const (
	FocusModeOneShotAF      FocusMode = 0   // One-shot AF
	FocusModeAIServoAF      FocusMode = 1   // AI Servo AF
	FocusModeAIFocusAF      FocusMode = 2   // AI Focus AF
	FocusModeManualFocus    FocusMode = 3   // Manual Focus
	FocusModeSingle         FocusMode = 4   // Single
	FocusModeContinuous     FocusMode = 5   // Continuous
	FocusModeManualFocusAlt FocusMode = 6   // Manual Focus
	FocusModePanFocus       FocusMode = 16  // Pan Focus
	FocusModeAFPlusMF       FocusMode = 256 // AF + MF
	FocusModeMF             FocusMode = 258 // MF
	FocusModeMovieSnapFocus FocusMode = 512 // Movie Snap Focus
	FocusModeMovieServoAF   FocusMode = 519 // Movie Servo AF
)

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
	MeteringModeDefault               MeteringMode = 0 // Default
	MeteringModeSpot                  MeteringMode = 1 // Spot
	MeteringModeAverage               MeteringMode = 2 // Average
	MeteringModeEvaluative            MeteringMode = 3 // Evaluative
	MeteringModePartial               MeteringMode = 4 // Partial
	MeteringModeCenterWeightedAverage MeteringMode = 5 // Center-weighted average
)

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
	FocusRangeManual      FocusRange = 0  // Manual
	FocusRangeAuto        FocusRange = 1  // Auto
	FocusRangeNotKnown    FocusRange = 2  // Not Known
	FocusRangeMacro       FocusRange = 3  // Macro
	FocusRangeVeryClose   FocusRange = 4  // Very Close
	FocusRangeClose       FocusRange = 5  // Close
	FocusRangeMiddleRange FocusRange = 6  // Middle Range
	FocusRangeFarRange    FocusRange = 7  // Far Range
	FocusRangePanFocus    FocusRange = 8  // Pan Focus
	FocusRangeSuperMacro  FocusRange = 9  // Super Macro
	FocusRangeInfinity    FocusRange = 10 // Infinity
)

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
//	9: "Manual (in movie mode)",
type ExposureMode int16

const (
	ExposureModeEasy                 ExposureMode = 0 // Easy
	ExposureModeProgramAE            ExposureMode = 1 // Program AE
	ExposureModeShutterSpeedPriority ExposureMode = 2 // Shutter speed priority AE
	ExposureModeAperturePriority     ExposureMode = 3 // Aperture-priority AE
	ExposureModeManual               ExposureMode = 4 // Manual
	ExposureModeDepthOfFieldAE       ExposureMode = 5 // Depth-of-field AE
	ExposureModeMDep                 ExposureMode = 6 // M-Dep
	ExposureModeBulb                 ExposureMode = 7 // Bulb
	ExposureModeFlexiblePriorityAE   ExposureMode = 8 // Flexible-priority AE
	ExposureModeManualMovie          ExposureMode = 9 // Manual (in movie mode)
)

// MacroMode is part of the CanonCameraSettings field.
//
//	1: "Macro"
//	2: "Normal"
type MacroMode int16

const (
	MacroModeMacro  MacroMode = 1 // Macro
	MacroModeNormal MacroMode = 2 // Normal
)

// Quality is part of the CanonCameraSettings field.
//
//	2: "Normal"
//	3: "Fine"
//	4: "RAW"
//	5: "Superfine"
type Quality int16

const (
	QualityNormal    Quality = 2 // Normal
	QualityFine      Quality = 3 // Fine
	QualityRAW       Quality = 4 // RAW
	QualitySuperfine Quality = 5 // Superfine
)

// CanonFlashMode is part of the CanonCameraSettings field.
//
//	0:  "Off"
//	1:  "Auto"
//	2:  "On"
//	3:  "Red-eye reduction"
//	4:  "Slow-sync"
//	5:  "Auto + red-eye"
//	6:  "On + red-eye"
//	16: "External flash"
type CanonFlashMode int16

const (
	CanonFlashModeOff           CanonFlashMode = 0  // Off
	CanonFlashModeAuto          CanonFlashMode = 1  // Auto
	CanonFlashModeOn            CanonFlashMode = 2  // On
	CanonFlashModeRedEye        CanonFlashMode = 3  // Red-eye reduction
	CanonFlashModeSlowSync      CanonFlashMode = 4  // Slow-sync
	CanonFlashModeAutoRedEye    CanonFlashMode = 5  // Auto + red-eye
	CanonFlashModeOnRedEye      CanonFlashMode = 6  // On + red-eye
	CanonFlashModeExternalFlash CanonFlashMode = 16 // External flash
)

// RecordMode is part of the CanonCameraSettings field.
//
//	1:  "JPEG"
//	2:  "CRW+THM"
//	3:  "AVI+THM"
//	4:  "TIF"
//	5:  "TIF+JPEG"
//	6:  "CR2"
//	7:  "CR2+JPEG"
//	9:  "MOV"
//	10: "MP4"
type RecordMode int16

const (
	RecordModeJPEG    RecordMode = 1  // JPEG
	RecordModeCRWTHM  RecordMode = 2  // CRW+THM
	RecordModeAVITHM  RecordMode = 3  // AVI+THM
	RecordModeTIF     RecordMode = 4  // TIF
	RecordModeTIFJPEG RecordMode = 5  // TIF+JPEG
	RecordModeCR2     RecordMode = 6  // CR2
	RecordModeCR2JPEG RecordMode = 7  // CR2+JPEG
	RecordModeMOV     RecordMode = 9  // MOV
	RecordModeMP4     RecordMode = 10 // MP4
)

// CanonImageSize is part of the CanonCameraSettings field.
//
//	0: "Large"
//	1: "Medium"
//	2: "Small"
//	5: "Medium 1"
//	6: "Medium 2"
//	7: "Medium 3"
type CanonImageSize int16

const (
	CanonImageSizeLarge   CanonImageSize = 0 // Large
	CanonImageSizeMedium  CanonImageSize = 1 // Medium
	CanonImageSizeSmall   CanonImageSize = 2 // Small
	CanonImageSizeMedium1 CanonImageSize = 5 // Medium 1
	CanonImageSizeMedium2 CanonImageSize = 6 // Medium 2
	CanonImageSizeMedium3 CanonImageSize = 7 // Medium 3
)

// EasyMode is part of the CanonCameraSettings field.
type EasyMode int16

const (
	EasyModeFullAuto            EasyMode = 0   // Full auto
	EasyModeManual              EasyMode = 1   // Manual
	EasyModeLandscape           EasyMode = 2   // Landscape
	EasyModeFastShutter         EasyMode = 3   // Fast shutter
	EasyModeSlowShutter         EasyMode = 4   // Slow shutter
	EasyModeNight               EasyMode = 5   // Night
	EasyModeBWP                 EasyMode = 7   // B&W
	EasyModeSepia               EasyMode = 8   // Sepia
	EasyModePortrait            EasyMode = 9   // Portrait
	EasyModeSports              EasyMode = 10  // Sports
	EasyModeMacroCloseUp        EasyMode = 11  // Macro / Close-up
	EasyModePanFocus            EasyMode = 19  // Pan focus
	EasyModeFoliage             EasyMode = 20  // Foliage
	EasyModeNightSnapshot       EasyMode = 23  // Night Snapshot
	EasyModeSuperMacro          EasyMode = 25  // Super Macro
	EasyModeLowLight            EasyMode = 26  // Low Light
	EasyModeStitchAssist        EasyMode = 28  // Stitch Assist
	EasyModeMovie               EasyMode = 29  // Movie
	EasyModeFireworks           EasyMode = 30  // Fireworks
	EasyModeLongShutter         EasyMode = 31  // Long Shutter
	EasyModeBeach               EasyMode = 33  // Beach
	EasyModeUnderwater          EasyMode = 34  // Underwater
	EasyModeSnow                EasyMode = 35  // Snow
	EasyModeIndoor              EasyMode = 36  // Indoor
	EasyModeKidsPets            EasyMode = 37  // Kids & Pets
	EasyModeNightPortrait       EasyMode = 38  // Night Portrait
	EasyModeShade               EasyMode = 39  // Shade
	EasyModeMyColors            EasyMode = 40  // My Colors
	EasyModeStillImage          EasyMode = 41  // Still Image
	EasyModeColorAccent         EasyMode = 42  // Color Accent
	EasyModeColorSwap           EasyMode = 43  // Color Swap
	EasyModeAquarium            EasyMode = 44  // Aquarium
	EasyModeISO3200             EasyMode = 45  // ISO 3200
	EasyModeISO6400             EasyMode = 46  // ISO 6400
	EasyModeCreativeLightEffect EasyMode = 47  // Creative Light Effect
	EasyModeEasy                EasyMode = 48  // Easy
	EasyModeQuickShot           EasyMode = 49  // Quick Shot
	EasyModeCreativeAuto        EasyMode = 50  // Creative Auto
	EasyModeZoomBlur            EasyMode = 52  // Zoom Blur
	EasyModeLowLight2           EasyMode = 53  // Low Light
	EasyModeNostalgic           EasyMode = 54  // Nostalgic
	EasyModeSuperVivid          EasyMode = 55  // Super Vivid
	EasyModePosterEffect        EasyMode = 56  // Poster Effect
	EasyModeFaceSelfTimer       EasyMode = 57  // Face Self-timer
	EasyModeSmile               EasyMode = 59  // Smile
	EasyModeWinkSelfTimer       EasyMode = 62  // Wink Self-timer
	EasyModeFisheyeEffect       EasyMode = 63  // Fisheye Effect
	EasyModeHandheldNightScene  EasyMode = 67  // Handheld Night Scene
	EasyModeHDRBacklightControl EasyMode = 83  // HDR Backlight Control
	EasyModeFood                EasyMode = 99  // Food
	EasyModeKids                EasyMode = 100 // Kids
	EasyModeSmoothSkin          EasyMode = 119 // Smooth Skin
	EasyModeHybridAuto          EasyMode = 257 // Hybrid Auto
	EasyModePowerShotWv         EasyMode = 258 // PowerShot wv
	EasyModePowerShotLv         EasyMode = 259 // PowerShot lv
	EasyModeCreativeShot        EasyMode = 260 // Creative Shot
	EasyModeSelfPortrait        EasyMode = 261 // Self Portrait
	EasyModeMovieDigest         EasyMode = 263 // Movie Digest
	EasyModeLiveViewControl     EasyMode = 264 // Live View Control
	EasyModeDiscreet            EasyMode = 265 // Discreet
	EasyModeBlurReduction       EasyMode = 266 // Blur Reduction
	EasyModeMonochrome          EasyMode = 267 // Monochrome
	EasyModeFisheyeEffect2      EasyMode = 268 // Fisheye Effect
	EasyModeWaterPaintingEffect EasyMode = 269 // Water Painting Effect
	EasyModeToyCameraEffect     EasyMode = 270 // Toy Camera Effect
	EasyModeMiniatureEffect     EasyMode = 271 // Miniature Effect
	EasyModeHDR                 EasyMode = 274 // HDR
)

// DigitalZoom is part of the CanonCameraSettings field.
//
//	0: "None"
//	1: "2x"
//	2: "4x"
//	3: "Other"
type DigitalZoom int16

const (
	DigitalZoomNone  DigitalZoom = 0 // None
	DigitalZoom2x    DigitalZoom = 1 // 2x
	DigitalZoom4x    DigitalZoom = 2 // 4x
	DigitalZoomOther DigitalZoom = 3 // Other
)

// CameraISO is part of the CanonCameraSettings field.
//
//	14: "n/a"
//	15: "Auto"
//	16: "50"
//	17: "100"
//	18: "200"
//	19: "400"
type CameraISO int16

const (
	CameraISONA   CameraISO = 14 // n/a
	CameraISOAuto CameraISO = 15 // Auto
	CameraISO50   CameraISO = 16 // 50
	CameraISO100  CameraISO = 17 // 100
	CameraISO200  CameraISO = 18 // 200
	CameraISO400  CameraISO = 19 // 400
)

// FlashModel is part of the CanonCameraSettings field.
type FlashModel int16

const (
	FlashModelNone                       FlashModel = 0  // None
	FlashModelEXSpeedlite                FlashModel = 1  // EX Speedlite
	FlashModelSpeedlite550EX             FlashModel = 2  // 550EX
	FlashModelSpeedlite420EX             FlashModel = 3  // 420EX
	FlashModelMacroRingLiteMR14EX        FlashModel = 4  // MR-14EX
	FlashModelSpeedlite220EX             FlashModel = 5  // 220EX
	FlashModelSpeedlite380EX             FlashModel = 6  // 380EX
	FlashModelSpeedlite470EXAI           FlashModel = 7  // 470EX-AI
	FlashModelSpeedlite600EX             FlashModel = 9  // 600EX
	FlashModelSpeedliteTransmitterSTE3RT FlashModel = 10 // ST-E3-RT
	FlashModelMacroRingLite              FlashModel = 11 // Macro Ring Lite
	FlashModelSpeedlite90EX              FlashModel = 12 // Canon Speedlite 90EX
	FlashModelSpeedlite270EXII           FlashModel = 13 // Canon Speedlite 270EX II
	FlashModelSpeedlite320EX             FlashModel = 14 // Canon Speedlite 320EX
	FlashModelSpeedlite430EXII           FlashModel = 15 // Canon Speedlite 430EX II
	FlashModelSpeedlite580EXII           FlashModel = 16 // Canon Speedlite 580EX II
	FlashModelSpeedlite270EX             FlashModel = 17 // Canon Speedlite 270EX
	FlashModelSpeedlite430EXIIIRT        FlashModel = 18 // Canon Speedlite 430EX III-RT
	FlashModelSpeedlite600EXIIRT         FlashModel = 19 // Canon Speedlite 600EX II-RT
)

// FocusContinuous is part of the CanonCameraSettings field.
//
//	0: "Single"
//	1: "Continuous"
type FocusContinuous int16

const (
	FocusContinuousSingle     FocusContinuous = 0 // Single
	FocusContinuousContinuous FocusContinuous = 1 // Continuous
)

// ImageStabilization is part of the CanonCameraSettings field.
//
//	0: "Off"
//	1: "On"
//	2: "Shoot Only"
//	3: "Panning"
//	4: "Dynamic"
//	-1: "n/a"
type ImageStabilization int16

const (
	ImageStabilizationOff       ImageStabilization = 0   // Off
	ImageStabilizationOn        ImageStabilization = 1   // On
	ImageStabilizationShootOnly ImageStabilization = 2   // Shoot Only
	ImageStabilizationPanning   ImageStabilization = 3   // Panning
	ImageStabilizationDynamic   ImageStabilization = 4   // Dynamic
	ImageStabilizationViaLens   ImageStabilization = 256 // Image Stabilization via the Lens
	ImageStabilizationNA        ImageStabilization = -1  // n/a
)

// SpotMeteringMode is part of the CanonCameraSettings field.
//
//	0: "Center"
//	1: "AF Point"
type SpotMeteringMode int16

const (
	SpotMeteringModeCenter  SpotMeteringMode = 0 // Center
	SpotMeteringModeAFPoint SpotMeteringMode = 1 // AF Point
)

// PhotoEffect is part of the CanonCameraSettings field.
type PhotoEffect int16

const (
	PhotoEffectOff         PhotoEffect = 0   // Off
	PhotoEffectVivid       PhotoEffect = 1   // Vivid
	PhotoEffectNeutral     PhotoEffect = 2   // Neutral
	PhotoEffectSmooth      PhotoEffect = 3   // Smooth
	PhotoEffectSepia       PhotoEffect = 4   // Sepia
	PhotoEffectBAndW       PhotoEffect = 5   // B&W
	PhotoEffectCustom      PhotoEffect = 6   // Custom
	PhotoEffectMyColorData PhotoEffect = 7   // My Color Data
	PhotoEffectBAndWAlt    PhotoEffect = 100 // B&W
)

// ManualFlashOutput is part of the CanonCameraSettings field.
type ManualFlashOutput int16

const (
	ManualFlashOutputNA     ManualFlashOutput = 0   // n/a
	ManualFlashOutputFull   ManualFlashOutput = 128 // Full
	ManualFlashOutputMedium ManualFlashOutput = 129 // Medium
	ManualFlashOutputLow    ManualFlashOutput = 130 // Low
	ManualFlashOutputNAAlt  ManualFlashOutput = 131 // n/a
)

// SRAWQuality is part of the CanonCameraSettings field.
//
//	0: "n/a"
//	1: "sRAW1 (mRAW)"
//	2: "sRAW2 (sRAW)"
type SRAWQuality int16

const (
	SRAWQualityNA    SRAWQuality = 0 // n/a
	SRAWQualitySRAW1 SRAWQuality = 1 // sRAW1 (mRAW)
	SRAWQualitySRAW2 SRAWQuality = 2 // sRAW2 (sRAW)
)

// FocusBracketing is part of the CanonCameraSettings field.
//
//	0: "Off"
//	1: "On"
type FocusBracketing int16

const (
	FocusBracketingOff FocusBracketing = 0 // Off
	FocusBracketingOn  FocusBracketing = 1 // On
)

// HDRPQ is part of the CanonCameraSettings field.
//
//	0: "Off"
//	1: "On"
type HDRPQ int16

const (
	HDRPQOff HDRPQ = 0 // Off
	HDRPQOn  HDRPQ = 1 // On
)

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

const (
	BracketModeOff BracketMode = 0 // Off
	BracketModeAEB BracketMode = 1 // AEB
	BracketModeFEB BracketMode = 2 // FEB
	BracketModeISO BracketMode = 3 // ISO
	BracketModeWB  BracketMode = 4 // WB
)

// Active - returns true if BracketMode is On
func (bm BracketMode) Active() bool {
	return bm != 0
}

// OnOffAuto is used by Canon FileInfo fields that encode Off/On/Auto states.
type OnOffAuto uint16

const (
	OnOffAutoOff    OnOffAuto = 0 // Off
	OnOffAutoOn1D   OnOffAuto = 1 // On
	OnOffAutoOn     OnOffAuto = 3 // On
	OnOffAutoAuto   OnOffAuto = 4 // Auto
	OnOffAutoNotSet OnOffAuto = 5 // n/a
)

// FilterEffect - Canon FileInfo FilterEffect.
//
//	0: "n/a"
//	1: "None"
//	2: "Yellow"
//	3: "Orange"
//	4: "Red"
type FilterEffect uint16

const (
	FilterEffectNA     FilterEffect = 0 // n/a
	FilterEffectNone   FilterEffect = 1 // None
	FilterEffectYellow FilterEffect = 2 // Yellow
	FilterEffectOrange FilterEffect = 3 // Orange
	FilterEffectRed    FilterEffect = 4 // Red
)

// ToningEffect - Canon FileInfo ToningEffect.
//
//	0: "n/a"
//	1: "None"
//	2: "Sepia"
//	3: "Blue"
//	4: "Purple"
type ToningEffect uint16

const (
	ToningEffectNA     ToningEffect = 0 // n/a
	ToningEffectNone   ToningEffect = 1 // None
	ToningEffectSepia  ToningEffect = 2 // Sepia
	ToningEffectBlue   ToningEffect = 3 // Blue
	ToningEffectPurple ToningEffect = 4 // Purple
)

// ShutterMode - Canon FileInfo ShutterMode.
//
//	0: "Mechanical"
//	1: "Electronic"
//	2: "Electronic (first curtain)"
type ShutterMode uint16

const (
	ShutterModeMechanical             ShutterMode = 0 // Mechanical
	ShutterModeElectronic             ShutterMode = 1 // Electronic
	ShutterModeElectronicFirstCurtain ShutterMode = 2 // Electronic (first curtain)
)

// RawJpgQuality - Canon FileInfo RawJpgQuality.
//
//	-1:  "n/a"
//	1:   "Econom."
//	2:   "Normal"
//	3:   "Fine"
//	4:   "RAW"
//	5:   "Superfine"
//	7:   "CRAW"
//	130: "HEIF"
//	131: "HEIF10Bit"
type RawJpgQuality uint16

const (
	RawJpgQualityNA        RawJpgQuality = 0xFFFF // n/a
	RawJpgQualityEconomy   RawJpgQuality = 1      // Econom.
	RawJpgQualityNormal    RawJpgQuality = 2      // Normal
	RawJpgQualityFine      RawJpgQuality = 3      // Fine
	RawJpgQualityRAW       RawJpgQuality = 4      // RAW
	RawJpgQualitySuperfine RawJpgQuality = 5      // Superfine
	RawJpgQualityCRAW      RawJpgQuality = 7      // CRAW
	RawJpgQualityHEIF      RawJpgQuality = 130    // HEIF
	RawJpgQualityHEIF10Bit RawJpgQuality = 131    // HEIF10Bit
)

// RawJpgSize - Canon FileInfo RawJpgSize.
//
//	-1:   "n/a"
//	0:    "Large"
//	1:    "Medium"
//	2:    "Small"
//	5:    "Medium 1"
//	6:    "Medium 2 (invalid size)"
//	7:    "Medium 3"
//	8:    "Postcard"
//	9:    "Widescreen"
//	10:   "Medium Widescreen"
//	14:   "Small 1"
//	15:   "Small 2"
//	16:   "Small 3"
//	128:  "5760x3840"
//	129:  "3840x2560"
//	130:  "1920x1280"
//	137:  "4096x2160"
//	142:  "5632x3168"
//	143:  "4864x3648"
type RawJpgSize uint16

const (
	RawJpgSizeNA                 RawJpgSize = 0xFFFF // n/a
	RawJpgSizeLarge              RawJpgSize = 0      // Large
	RawJpgSizeMedium             RawJpgSize = 1      // Medium
	RawJpgSizeSmall              RawJpgSize = 2      // Small
	RawJpgSizeMedium1            RawJpgSize = 5      // Medium 1
	RawJpgSizeMedium2InvalidSize RawJpgSize = 6      // Medium 2 (invalid size)
	RawJpgSizeMedium3            RawJpgSize = 7      // Medium 3
	RawJpgSizePostcard           RawJpgSize = 8      // Postcard
	RawJpgSizeWidescreen         RawJpgSize = 9      // Widescreen
	RawJpgSizeMediumWidescreen   RawJpgSize = 10     // Medium Widescreen
	RawJpgSizeSmall1             RawJpgSize = 14     // Small 1
	RawJpgSizeSmall2             RawJpgSize = 15     // Small 2
	RawJpgSizeSmall3             RawJpgSize = 16     // Small 3
	RawJpgSize5760x3840          RawJpgSize = 128    // 5760x3840
	RawJpgSize3840x2560          RawJpgSize = 129    // 3840x2560
	RawJpgSize1920x1280          RawJpgSize = 130    // 1920x1280
	RawJpgSize4096x2160          RawJpgSize = 137    // 4096x2160
	RawJpgSize5632x3168          RawJpgSize = 142    // 5632x3168
	RawJpgSize4864x3648          RawJpgSize = 143    // 4864x3648
)

// TimeZoneCity - Canon TimeInfo TimeZoneCity.
//
//	0:  "n/a"
//	1:  "Manual"
//	2:  "Time zone"
//	3:  "Athens"
//	4:  "Auckland"
//	5:  "Bangkok"
//	6:  "Beijing"
//	7:  "Caracas"
//	8:  "Casablanca"
//	9:  "Darwin"
//	10: "Edmonton"
//	11: "Honolulu"
//	12: "Kathmandu"
//	13: "London"
//	14: "New York"
//	15: "Samoa"
//	16: "Santiago"
//	17: "Tokyo"
//	18: "Abu Dhabi"
//	19: "Anchorage"
//	20: "Buenos Aires"
//	21: "Chicago"
//	22: "Denver"
//	23: "Dubai"
type TimeZoneCity int32

const (
	TimeZoneCityNA          TimeZoneCity = 0  // n/a
	TimeZoneCityManual      TimeZoneCity = 1  // Manual
	TimeZoneCityTimeZone    TimeZoneCity = 2  // Time zone
	TimeZoneCityAthens      TimeZoneCity = 3  // Athens
	TimeZoneCityAuckland    TimeZoneCity = 4  // Auckland
	TimeZoneCityBangkok     TimeZoneCity = 5  // Bangkok
	TimeZoneCityBeijing     TimeZoneCity = 6  // Beijing
	TimeZoneCityCaracas     TimeZoneCity = 7  // Caracas
	TimeZoneCityCasablanca  TimeZoneCity = 8  // Casablanca
	TimeZoneCityDarwin      TimeZoneCity = 9  // Darwin
	TimeZoneCityEdmonton    TimeZoneCity = 10 // Edmonton
	TimeZoneCityHonolulu    TimeZoneCity = 11 // Honolulu
	TimeZoneCityKathmandu   TimeZoneCity = 12 // Kathmandu
	TimeZoneCityLondon      TimeZoneCity = 13 // London
	TimeZoneCityNewYork     TimeZoneCity = 14 // New York
	TimeZoneCitySamoa       TimeZoneCity = 15 // Samoa
	TimeZoneCitySantiago    TimeZoneCity = 16 // Santiago
	TimeZoneCityTokyo       TimeZoneCity = 17 // Tokyo
	TimeZoneCityAbuDhabi    TimeZoneCity = 18 // Abu Dhabi
	TimeZoneCityAnchorage   TimeZoneCity = 19 // Anchorage
	TimeZoneCityBuenosAires TimeZoneCity = 20 // Buenos Aires
	TimeZoneCityChicago     TimeZoneCity = 21 // Chicago
	TimeZoneCityDenver      TimeZoneCity = 22 // Denver
	TimeZoneCityDubai       TimeZoneCity = 23 // Dubai
)

// DaylightSavings - Canon TimeInfo DaylightSavings.
//
//	0:  "Off"
//	60: "On"
type DaylightSavings int32

const (
	DaylightSavingsOff DaylightSavings = 0  // Off
	DaylightSavingsOn  DaylightSavings = 60 // On
)

// GMTOffset holds timezone offset values in hours.
//
// HoursFromGMT is the standard offset from GMT.
// DSTHours is the daylight-savings adjustment to add when DST is active.
type GMTOffset struct {
	HoursFromGMT float64
	DSTHours     float64
}

// Hours returns the effective GMT offset in hours.
func (o GMTOffset) Hours(daylightSavings bool) float64 {
	if daylightSavings {
		return o.HoursFromGMT + o.DSTHours
	}
	return o.HoursFromGMT
}

// EffectiveGMTHours is kept for compatibility.
//
// Deprecated: use Hours.
func (o GMTOffset) EffectiveGMTHours(daylightSavings bool) float64 {
	return o.Hours(daylightSavings)
}

// cityGMTOffsets maps Canon timezone city enum values to GMT offsets.
var cityGMTOffsets = map[TimeZoneCity]GMTOffset{
	TimeZoneCityNA:          {HoursFromGMT: 0, DSTHours: 0},
	TimeZoneCityManual:      {HoursFromGMT: 0, DSTHours: 0},
	TimeZoneCityTimeZone:    {HoursFromGMT: 0, DSTHours: 0},
	TimeZoneCityAthens:      {HoursFromGMT: 2, DSTHours: 1},
	TimeZoneCityAuckland:    {HoursFromGMT: 12, DSTHours: 1},
	TimeZoneCityBangkok:     {HoursFromGMT: 7, DSTHours: 0},
	TimeZoneCityBeijing:     {HoursFromGMT: 8, DSTHours: 0},
	TimeZoneCityCaracas:     {HoursFromGMT: -4, DSTHours: 0},
	TimeZoneCityCasablanca:  {HoursFromGMT: 1, DSTHours: 0},
	TimeZoneCityDarwin:      {HoursFromGMT: 9.5, DSTHours: 0},
	TimeZoneCityEdmonton:    {HoursFromGMT: -7, DSTHours: 1},
	TimeZoneCityHonolulu:    {HoursFromGMT: -10, DSTHours: 0},
	TimeZoneCityKathmandu:   {HoursFromGMT: 5.75, DSTHours: 0},
	TimeZoneCityLondon:      {HoursFromGMT: 0, DSTHours: 1},
	TimeZoneCityNewYork:     {HoursFromGMT: -5, DSTHours: 1},
	TimeZoneCitySamoa:       {HoursFromGMT: 13, DSTHours: 0},
	TimeZoneCitySantiago:    {HoursFromGMT: -4, DSTHours: 1},
	TimeZoneCityTokyo:       {HoursFromGMT: 9, DSTHours: 0},
	TimeZoneCityAbuDhabi:    {HoursFromGMT: 4, DSTHours: 0},
	TimeZoneCityAnchorage:   {HoursFromGMT: -9, DSTHours: 1},
	TimeZoneCityBuenosAires: {HoursFromGMT: -3, DSTHours: 0},
	TimeZoneCityChicago:     {HoursFromGMT: -6, DSTHours: 1},
	TimeZoneCityDenver:      {HoursFromGMT: -7, DSTHours: 1},
	TimeZoneCityDubai:       {HoursFromGMT: 4, DSTHours: 0},
}

// GMTOffsetForCity returns the GMT offset config for a Canon timezone city.
func GMTOffsetForCity(city TimeZoneCity) (GMTOffset, bool) {
	offset, ok := cityGMTOffsets[city]
	return offset, ok
}

// GMTHoursForCity returns the effective GMT offset in hours for a city.
func GMTHoursForCity(city TimeZoneCity, daylightSavings bool) (float64, bool) {
	offset, ok := cityGMTOffsets[city]
	if !ok {
		return 0, false
	}
	return offset.Hours(daylightSavings), true
}

// EffectiveGMTHoursForCity is kept for compatibility.
//
// Deprecated: use GMTHoursForCity.
func EffectiveGMTHoursForCity(city TimeZoneCity, daylightSavings bool) (float64, bool) {
	return GMTHoursForCity(city, daylightSavings)
}

// cityIANA maps Canon timezone city enum values to IANA timezone names.
var cityIANA = [...]string{
	TimeZoneCityNA:          "", // n/a
	TimeZoneCityManual:      "", // manual offset/config, not a city
	TimeZoneCityTimeZone:    "", // generic "time zone", not a city
	TimeZoneCityAthens:      "Europe/Athens",
	TimeZoneCityAuckland:    "Pacific/Auckland",
	TimeZoneCityBangkok:     "Asia/Bangkok",
	TimeZoneCityBeijing:     "Asia/Shanghai", // Beijing uses China standard time
	TimeZoneCityCaracas:     "America/Caracas",
	TimeZoneCityCasablanca:  "Africa/Casablanca",
	TimeZoneCityDarwin:      "Australia/Darwin",
	TimeZoneCityEdmonton:    "America/Edmonton",
	TimeZoneCityHonolulu:    "Pacific/Honolulu",
	TimeZoneCityKathmandu:   "Asia/Kathmandu",
	TimeZoneCityLondon:      "Europe/London",
	TimeZoneCityNewYork:     "America/New_York",
	TimeZoneCitySamoa:       "Pacific/Apia", // country of Samoa
	TimeZoneCitySantiago:    "America/Santiago",
	TimeZoneCityTokyo:       "Asia/Tokyo",
	TimeZoneCityAbuDhabi:    "Asia/Dubai", // UAE
	TimeZoneCityAnchorage:   "America/Anchorage",
	TimeZoneCityBuenosAires: "America/Argentina/Buenos_Aires",
	TimeZoneCityChicago:     "America/Chicago",
	TimeZoneCityDenver:      "America/Denver",
	TimeZoneCityDubai:       "Asia/Dubai",
}

var cityLocations = buildCityLocations()

func buildCityLocations() []*time.Location {
	const cityCount = int(TimeZoneCityDubai) + 1
	locations := make([]*time.Location, cityCount)

	for city, name := range cityIANA {
		if name == "" {
			continue
		}

		loc, err := time.LoadLocation(name)
		if err != nil {
			continue
		}
		locations[int(city)] = loc
	}

	return locations
}

// LocationForCity returns the cached IANA location for a Canon timezone city.
func LocationForCity(city TimeZoneCity) (*time.Location, bool) {
	idx := int(city)
	if idx < 0 || idx >= len(cityLocations) {
		return nil, false
	}

	loc := cityLocations[idx]
	if loc == nil {
		return nil, false
	}

	return loc, true
}

// AESetting - Canon Makernote AutoExposure Setting
//
//	0: "Normal AE",
//	1: "Exposure Compensation",
//	2: "AE Lock",
//	3: "AE Lock + Exposure Compensation",
//	4: "No AE",
//	5: "Pattern",
type AESetting int16

const (
	AESettingNormalAE                       AESetting = 0 // Normal AE
	AESettingExposureCompensation           AESetting = 1 // Exposure Compensation
	AESettingAELock                         AESetting = 2 // AE Lock
	AESettingAELockWithExposureCompensation AESetting = 3 // AE Lock + Exposure Compensation
	AESettingNoAE                           AESetting = 4 // No AE
	AESettingPattern                        AESetting = 5 // Pattern
)

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
	AFAreaModeOffManualFocus    AFAreaMode = 0  // Off (Manual Focus)
	AFAreaModeAFPointExpansion  AFAreaMode = 1  // AF Point Expansion (surround)
	AFAreaModeSinglePointAF     AFAreaMode = 2  // Single-point AF
	AFAreaModeAuto              AFAreaMode = 4  // Auto
	AFAreaModeFaceDetectAF      AFAreaMode = 5  // Face Detect AF
	AFAreaModeFaceTracking      AFAreaMode = 6  // Face + Tracking
	AFAreaModeZoneAF            AFAreaMode = 7  // Zone AF
	AFAreaModeAFPointExpansion4 AFAreaMode = 8  // AF Point Expansion (4 point)
	AFAreaModeSpotAF            AFAreaMode = 9  // Spot AF
	AFAreaModeAFPointExpansion8 AFAreaMode = 10 // AF Point Expansion (8 point)
	AFAreaModeFlexizoneMulti49  AFAreaMode = 11 // Flexizone Multi (49 point)
	AFAreaModeFlexizoneMulti9   AFAreaMode = 12 // Flexizone Multi (9 point)
	AFAreaModeFlexizoneSingle   AFAreaMode = 13 // Flexizone Single
	AFAreaModeLargeZoneAF       AFAreaMode = 14 // Large Zone AF
)
