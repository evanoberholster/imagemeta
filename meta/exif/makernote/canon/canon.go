// Package canon provides data types and functions for representing Canon Camera Makernote values
package canon

import (
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

//go:generate msgp
//go:generate stringer -type=MacroMode,Quality,CanonFlashMode,ContinuousDrive,FocusMode,RecordMode,CanonImageSize,EasyMode,DigitalZoom,MeteringMode,FocusRange,ExposureMode,FlashModel,FocusContinuous,AESetting,ImageStabilization,SpotMeteringMode,PhotoEffect,ManualFlashOutput,SRAWQuality,FocusBracketing,HDRPQ,BracketMode,OnOffAuto,FilterEffect,ToningEffect,ShutterMode,RawJpgQuality,RawJpgSize,TimeZoneCity,DaylightSavings,AFAreaMode -linecomment -output=canon_string.go

// Canon contains selected Canon maker-note fields.
//
// Parsing lives in meta/exif/canon.go. This package keeps the parsed container
// model shared by the Exif result.
// Canon maker-note tag IDs from ExifTool Canon::Main are defined in
// makernote_tags.go.
//
// TODO(canon): Continue Canon maker-note parity work against ExifTool, including
// CanonCustom, VRD, sensor/camera-temperature variants, preview offsets and
// newer model/lens fields.
//
//msgp:ignore Canon
type Canon struct {
	ImageType            string    // 16 bytes
	FirmwareVersion      string    // 16 bytes
	OwnerName            string    // 16 bytes
	ImageUniqueID        meta.UUID // 16 bytes
	LensModel            string    // 16 bytes
	InternalSerialNumber string    // 16 bytes
	BatteryType          string    // 16 bytes

	FileNumber   uint32 // 4 bytes
	SerialNumber uint32 // 4 bytes
	ModelID      uint32 // 4 bytes
	ColorSpace   uint16 // 4 bytes in struct (2 data + 2 padding)
	// ColorTemperature is Canon maker-note main tag 0x00ae, distinct from
	// model-specific CameraInfo/ProcessingInfo color-temperature fields.
	ColorTemperature uint16

	// Structured Canon maker-note tables (ExifTool Canon.pm mappings).
	CanonCameraSettings        CameraSettings     // 88 bytes
	CameraInfo                 CameraInfo         // selected Canon CameraInfo fields
	CanonFocalLength           FocalLengthInfo    // 8 bytes
	CanonShotInfo              ShotInfo           // 100 bytes
	CanonFileInfo              FileInfo           // 48 bytes
	TimeInfo                   CanonTimeInfo      // 12 bytes
	AFInfo                     AFInfo             // 128 bytes
	FaceDetect                 FaceDetectInfo     // combined FaceDetect1+2+3
	AspectInfo                 AspectInfo         // 20 bytes
	ProcessingInfo             ProcessingInfo     // 36 bytes in struct (30 data + 6 padding)
	CustomPictureStyleFileName string             // 16 bytes
	AFMicroAdj                 AFMicroAdjInfo     // 16 bytes in struct (12 data + 4 padding)
	LensInfo                   LensInfoForService // 24 bytes
	MultiExp                   MultiExpInfo       // 12 bytes
	HDRInfo                    HDRInfo            // 8 bytes
	PreviewImageInfo           PreviewImageInfo   // 20 bytes
	SensorInfo                 SensorInfo         // 20 bytes
	AFConfig                   AFConfig           // 84 bytes
	RawBurstModeRoll           RawBurstInfo       // 8 bytes
	LightingOpt                LightingOptInfo    // 32 bytes in struct (28 data + 4 padding)
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
//	8:  "Continuous, High+",
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
	ContinuousDriveContinuousHighPlus      ContinuousDrive = 8  // Continuous, High+
	ContinuousDriveSingleSilent            ContinuousDrive = 9  // Single, Silent
	ContinuousDriveContinuousSilent        ContinuousDrive = 10 // Continuous, Silent
)

// FocusMode is part of the CanonCameraSettings field
//
//	0:   "One-shot AF",
//	1:   "AI Servo AF",
//	2:   "AI Focus AF",
//	3:   "Manual Focus (3)",
//	4:   "Single",
//	5:   "Continuous",
//	6:   "Manual Focus (6)",
//	16:  "Pan Focus",
//	256: "One-shot AF (Live View)",
//	257: "AI Servo AF (Live View)",
//	258: "AI Focus AF (Live View)",
//	512: "Movie Snap Focus",
//	519: "Movie Servo AF",
type FocusMode int16

const (
	FocusModeOneShotAF         FocusMode = 0   // One-shot AF
	FocusModeAIServoAF         FocusMode = 1   // AI Servo AF
	FocusModeAIFocusAF         FocusMode = 2   // AI Focus AF
	FocusModeManualFocus3      FocusMode = 3   // Manual Focus (3)
	FocusModeSingle            FocusMode = 4   // Single
	FocusModeContinuous        FocusMode = 5   // Continuous
	FocusModeManualFocus6      FocusMode = 6   // Manual Focus (6)
	FocusModePanFocus          FocusMode = 16  // Pan Focus
	FocusModeOneShotAFLiveView FocusMode = 256 // One-shot AF (Live View)
	FocusModeAIServoAFLiveView FocusMode = 257 // AI Servo AF (Live View)
	FocusModeAIFocusAFLiveView FocusMode = 258 // AI Focus AF (Live View)
	FocusModeMovieSnapFocus    FocusMode = 512 // Movie Snap Focus
	FocusModeMovieServoAF      FocusMode = 519 // Movie Servo AF
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
//	-1: "n/a"
//	1:  "Economy"
//	2:  "Normal"
//	3:  "Fine"
//	4:  "RAW"
//	5:  "Superfine"
//	7:  "CRAW"
type Quality int16

const (
	QualityNA        Quality = -1 // n/a
	QualityEconomy   Quality = 1  // Economy
	QualityNormal    Quality = 2  // Normal
	QualityFine      Quality = 3  // Fine
	QualityRAW       Quality = 4  // RAW
	QualitySuperfine Quality = 5  // Superfine
	QualityCRAW      Quality = 7  // CRAW
)

// CanonFlashMode is part of the CanonCameraSettings field.
//
//	-1: "n/a"
//	0:  "Off"
//	1:  "Auto"
//	2:  "On"
//	3:  "Red-eye reduction"
//	4:  "Slow-sync"
//	5:  "Red-eye reduction (Auto)"
//	6:  "Red-eye reduction (On)"
//	16: "External flash"
type CanonFlashMode int16

const (
	CanonFlashModeNA            CanonFlashMode = -1 // n/a
	CanonFlashModeOff           CanonFlashMode = 0  // Off
	CanonFlashModeAuto          CanonFlashMode = 1  // Auto
	CanonFlashModeOn            CanonFlashMode = 2  // On
	CanonFlashModeRedEye        CanonFlashMode = 3  // Red-eye reduction
	CanonFlashModeSlowSync      CanonFlashMode = 4  // Slow-sync
	CanonFlashModeAutoRedEye    CanonFlashMode = 5  // Red-eye reduction (Auto)
	CanonFlashModeOnRedEye      CanonFlashMode = 6  // Red-eye reduction (On)
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
//	11: "CRM"
//	12: "CR3"
//	13: "CR3+JPEG"
//	14: "HIF"
//	15: "CR3+HIF"
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
	RecordModeCRM     RecordMode = 11 // CRM
	RecordModeCR3     RecordMode = 12 // CR3
	RecordModeCR3JPEG RecordMode = 13 // CR3+JPEG
	RecordModeHIF     RecordMode = 14 // HIF
	RecordModeCR3HIF  RecordMode = 15 // CR3+HIF
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
	CanonImageSizeNA      CanonImageSize = -1 // n/a
	CanonImageSizeLarge   CanonImageSize = 0  // Large
	CanonImageSizeMedium  CanonImageSize = 1  // Medium
	CanonImageSizeSmall   CanonImageSize = 2  // Small
	CanonImageSizeMedium1 CanonImageSize = 5  // Medium 1
	CanonImageSizeMedium2 CanonImageSize = 6  // Medium 2
	CanonImageSizeMedium3 CanonImageSize = 7  // Medium 3
)

// EasyMode is part of the CanonCameraSettings field.
type EasyMode int16

const (
	EasyModeFullAuto             EasyMode = 0   // Full auto
	EasyModeManual               EasyMode = 1   // Manual
	EasyModeLandscape            EasyMode = 2   // Landscape
	EasyModeFastShutter          EasyMode = 3   // Fast shutter
	EasyModeSlowShutter          EasyMode = 4   // Slow shutter
	EasyModeNight                EasyMode = 5   // Night
	EasyModeGrayScale            EasyMode = 6   // Gray Scale
	EasyModeSepia                EasyMode = 7   // Sepia
	EasyModePortrait             EasyMode = 8   // Portrait
	EasyModeSports               EasyMode = 9   // Sports
	EasyModeMacro                EasyMode = 10  // Macro
	EasyModeBlackAndWhite        EasyMode = 11  // Black & White
	EasyModePanFocus             EasyMode = 12  // Pan focus
	EasyModeVivid                EasyMode = 13  // Vivid
	EasyModeNeutral              EasyMode = 14  // Neutral
	EasyModeFlashOff             EasyMode = 15  // Flash Off
	EasyModeLongShutter          EasyMode = 16  // Long Shutter
	EasyModeSuperMacro           EasyMode = 17  // Super Macro
	EasyModeFoliage              EasyMode = 18  // Foliage
	EasyModeIndoor               EasyMode = 19  // Indoor
	EasyModeFireworks            EasyMode = 20  // Fireworks
	EasyModeBeach                EasyMode = 21  // Beach
	EasyModeUnderwater           EasyMode = 22  // Underwater
	EasyModeSnow                 EasyMode = 23  // Snow
	EasyModeKidsPets             EasyMode = 24  // Kids & Pets
	EasyModeNightSnapshot        EasyMode = 25  // Night Snapshot
	EasyModeDigitalMacro         EasyMode = 26  // Digital Macro
	EasyModeMyColors             EasyMode = 27  // My Colors
	EasyModeMovieSnap            EasyMode = 28  // Movie Snap
	EasyModeSuperMacro2          EasyMode = 29  // Super Macro 2
	EasyModeColorAccent          EasyMode = 30  // Color Accent
	EasyModeColorSwap            EasyMode = 31  // Color Swap
	EasyModeAquarium             EasyMode = 32  // Aquarium
	EasyModeISO3200              EasyMode = 33  // ISO 3200
	EasyModeISO6400              EasyMode = 34  // ISO 6400
	EasyModeCreativeLightEffect  EasyMode = 35  // Creative Light Effect
	EasyModeEasy                 EasyMode = 36  // Easy
	EasyModeQuickShot            EasyMode = 37  // Quick Shot
	EasyModeCreativeAuto         EasyMode = 38  // Creative Auto
	EasyModeZoomBlur             EasyMode = 39  // Zoom Blur
	EasyModeLowLight             EasyMode = 40  // Low Light
	EasyModeNostalgic            EasyMode = 41  // Nostalgic
	EasyModeSuperVivid           EasyMode = 42  // Super Vivid
	EasyModePosterEffect         EasyMode = 43  // Poster Effect
	EasyModeFaceSelfTimer        EasyMode = 44  // Face Self-timer
	EasyModeSmile                EasyMode = 45  // Smile
	EasyModeWinkSelfTimer        EasyMode = 46  // Wink Self-timer
	EasyModeFisheyeEffect        EasyMode = 47  // Fisheye Effect
	EasyModeMiniatureEffect      EasyMode = 48  // Miniature Effect
	EasyModeHighSpeedBurst       EasyMode = 49  // High-speed Burst
	EasyModeBestImageSelection   EasyMode = 50  // Best Image Selection
	EasyModeHighDynamicRange     EasyMode = 51  // High Dynamic Range
	EasyModeHandheldNightScene   EasyMode = 52  // Handheld Night Scene
	EasyModeMovieDigest          EasyMode = 53  // Movie Digest
	EasyModeLiveViewControl      EasyMode = 54  // Live View Control
	EasyModeDiscreet             EasyMode = 55  // Discreet
	EasyModeBlurReduction        EasyMode = 56  // Blur Reduction
	EasyModeMonochrome           EasyMode = 57  // Monochrome
	EasyModeToyCameraEffect      EasyMode = 58  // Toy Camera Effect
	EasyModeSceneIntelligentAuto EasyMode = 59  // Scene Intelligent Auto
	EasyModeHighSpeedBurstHQ     EasyMode = 60  // High-speed Burst HQ
	EasyModeSmoothSkin           EasyMode = 61  // Smooth Skin
	EasyModeSoftFocus            EasyMode = 62  // Soft Focus
	EasyModeFood                 EasyMode = 68  // Food
	EasyModeHDRArtStandard       EasyMode = 84  // HDR Art Standard
	EasyModeHDRArtVivid          EasyMode = 85  // HDR Art Vivid
	EasyModeHDRArtBold           EasyMode = 93  // HDR Art Bold
	EasyModeSpotlight            EasyMode = 257 // Spotlight
	EasyModeNight2               EasyMode = 258 // Night 2
	EasyModeNightPlus            EasyMode = 259 // Night+
	EasyModeSuperNight           EasyMode = 260 // Super Night
	EasyModeSunset               EasyMode = 261 // Sunset
	EasyModeNightScene           EasyMode = 263 // Night Scene
	EasyModeSurface              EasyMode = 264 // Surface
	EasyModeLowLight2            EasyMode = 265 // Low Light 2
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

// CameraISO stores the raw CameraSettings index-16 value.
//
// For newer EOS models the value is an enum (0=n/a, 14=Auto High, 15=Auto,
// 16..20=50..800). For older PowerShot models the value encodes the actual ISO
// directly (either in the lower 14 bits with bit 0x4000 set, or as a raw ISO
// number).
//
// Use CameraISOValue to resolve the raw stored value into an ISO number.
//
//	0:  "n/a"
//	14: "Auto High"
//	15: "Auto"
//	16: "50"
//	17: "100"
//	18: "200"
//	19: "400"
//	20: "800"
type CameraISO int16

const (
	CameraISONA       CameraISO = 0  // n/a
	CameraISOAutoHigh CameraISO = 14 // Auto High
	CameraISOAuto     CameraISO = 15 // Auto
	CameraISO50       CameraISO = 16 // 50
	CameraISO100      CameraISO = 17 // 100
	CameraISO200      CameraISO = 18 // 200
	CameraISO400      CameraISO = 19 // 400
	CameraISO800      CameraISO = 20 // 800
)

const (
	CameraISOAutoSentinel     = math.MaxUint32     // resolved ISO sentinel for "Auto"
	CameraISOAutoHighSentinel = math.MaxUint32 - 1 // resolved ISO sentinel for "Auto High"
)

// NewCameraISO returns a validated CameraISO enum, or 0 for unknown values.
func NewCameraISO(v int16) CameraISO {
	switch CameraISO(v) {
	case CameraISONA,
		CameraISOAutoHigh,
		CameraISOAuto,
		CameraISO50,
		CameraISO100,
		CameraISO200,
		CameraISO400,
		CameraISO800:
		return CameraISO(v)
	default:
		return 0
	}
}

// NewCameraISOFromRaw decodes a raw uint16 wire value into CameraISO.
func NewCameraISOFromRaw(raw uint16) CameraISO {
	return NewCameraISO(meta.SafecastUint16ToInt16Bits(raw))
}

// CameraISOValue resolves the raw CameraISO value using ExifTool-style logic.
//
// Returns the resolved ISO value, or a sentinel (CameraISOAutoSentinel /
// CameraISOAutoHighSentinel) for non-numeric modes. Returns 0 for n/a.
func CameraISOValue(raw int16) int64 {
	switch {
	case raw == 0x7fff:
		return 0
	case raw&0x4000 != 0:
		return int64(raw & 0x3fff)
	}
	switch CameraISO(raw) {
	case 0:
		return 0
	case 14:
		return CameraISOAutoHighSentinel
	case 15:
		return CameraISOAutoSentinel
	case 16:
		return 50
	case 17:
		return 100
	case 18:
		return 200
	case 19:
		return 400
	case 20:
		return 800
	default:
		return int64(raw)
	}
}

// NewResolvedCameraISOFromRaw resolves a raw CameraISO value to ExifTool-style output.
// Returns 0 if conversion to uint32 fails.
func NewResolvedCameraISOFromRaw(raw int16) uint32 {
	value, ok := meta.SafecastInt64ToUint32(CameraISOValue(raw))
	if !ok {
		return 0
	}
	return value
}

// CameraISOString returns the ExifTool-style display string for a resolved CameraISO value.
func CameraISOString(v uint32) string {
	switch v {
	case 0:
		return "n/a"
	case CameraISOAutoSentinel:
		return "Auto"
	case CameraISOAutoHighSentinel:
		return "Auto High"
	default:
		return strconv.FormatUint(uint64(v), 10)
	}
}

// String returns the display string for the raw CameraISO enum value.
func (i CameraISO) String() string {
	return CameraISOString(NewResolvedCameraISOFromRaw(int16(i)))
}

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
//	8: "Manual"
type FocusContinuous int16

const (
	FocusContinuousSingle     FocusContinuous = 0 // Single
	FocusContinuousContinuous FocusContinuous = 1 // Continuous
	FocusContinuousManual     FocusContinuous = 8 // Manual
)

// ImageStabilization is part of the CanonCameraSettings field.
//
//	0:   "Off"
//	1:   "On"
//	2:   "Shoot Only"
//	3:   "Panning"
//	4:   "Dynamic"
//	256: "Off (2)"
//	257: "On (2)"
//	258: "Shoot Only (2)"
//	259: "Panning (2)"
//	260: "Dynamic (2)"
//	-1:  "n/a"
type ImageStabilization int16

const (
	ImageStabilizationOff        ImageStabilization = 0   // Off
	ImageStabilizationOn         ImageStabilization = 1   // On
	ImageStabilizationShootOnly  ImageStabilization = 2   // Shoot Only
	ImageStabilizationPanning    ImageStabilization = 3   // Panning
	ImageStabilizationDynamic    ImageStabilization = 4   // Dynamic
	ImageStabilizationOff2       ImageStabilization = 256 // Off (2)
	ImageStabilizationOn2        ImageStabilization = 257 // On (2)
	ImageStabilizationShootOnly2 ImageStabilization = 258 // Shoot Only (2)
	ImageStabilizationPanning2   ImageStabilization = 259 // Panning (2)
	ImageStabilizationDynamic2   ImageStabilization = 260 // Dynamic (2)
	ImageStabilizationNA         ImageStabilization = -1  // n/a
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
	PhotoEffectMyColorData PhotoEffect = 100 // My Color Data
)

// ManualFlashOutput is part of the CanonCameraSettings field.
type ManualFlashOutput int16

const (
	ManualFlashOutputNA0  ManualFlashOutput = 0      // n/a
	ManualFlashOutputFull ManualFlashOutput = 0x500  // Full
	ManualFlashOutputMed  ManualFlashOutput = 0x502  // Medium
	ManualFlashOutputLow  ManualFlashOutput = 0x504  // Low
	ManualFlashOutputNA   ManualFlashOutput = 0x7fff // n/a
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
//	0: "Disable"
//	1: "Enable"
type FocusBracketing int16

const (
	FocusBracketingDisable FocusBracketing = 0 // Disable
	FocusBracketingEnable  FocusBracketing = 1 // Enable
)

// HDRPQ is part of the CanonCameraSettings field.
//
//	-1: "n/a"
//	0:  "Off"
//	1:  "On"
type HDRPQ int16

const (
	HDRPQNA  HDRPQ = -1 // n/a
	HDRPQOff HDRPQ = 0  // Off
	HDRPQOn  HDRPQ = 1  // On
)

// FocusDistance -
type FocusDistance [2]int16

// NewFocusDistance creates a new FocusDistance with the upper
// and lower limits
func NewFocusDistance(upper, lower uint16) FocusDistance {
	return FocusDistance{meta.SafecastUint16ToInt16Bits(upper), meta.SafecastUint16ToInt16Bits(lower)}
}

// UpperRaw returns the raw upper focus-distance code as stored by Canon.
func (fd FocusDistance) UpperRaw() uint16 {
	return meta.SafecastInt16ToUint16Bits(fd[0])
}

// LowerRaw returns the raw lower focus-distance code as stored by Canon.
func (fd FocusDistance) LowerRaw() uint16 {
	return meta.SafecastInt16ToUint16Bits(fd[1])
}

// UpperMeters converts the upper focus-distance code into meters like ExifTool.
func (fd FocusDistance) UpperMeters() float32 {
	return focusDistanceMeters(fd.UpperRaw())
}

// LowerMeters converts the lower focus-distance code into meters like ExifTool.
func (fd FocusDistance) LowerMeters() float32 {
	return focusDistanceMeters(fd.LowerRaw())
}

func focusDistanceMeters(raw uint16) float32 {
	switch raw {
	case 0:
		return 0
	case 0xffff:
		return float32(math.Inf(1))
	default:
		return float32(raw) / 100.0
	}
}

// WhiteBalance is part of the CanonShotInfo field.
type WhiteBalance int16

const (
	WhiteBalanceAuto                 WhiteBalance = 0  // Auto
	WhiteBalanceDaylight             WhiteBalance = 1  // Daylight
	WhiteBalanceCloudy               WhiteBalance = 2  // Cloudy
	WhiteBalanceTungsten             WhiteBalance = 3  // Tungsten
	WhiteBalanceFluorescent          WhiteBalance = 4  // Fluorescent
	WhiteBalanceFlash                WhiteBalance = 5  // Flash
	WhiteBalanceCustom               WhiteBalance = 6  // Custom
	WhiteBalanceBlackAndWhite        WhiteBalance = 7  // B&W
	WhiteBalanceShade                WhiteBalance = 8  // Shade
	WhiteBalanceManualTemperature    WhiteBalance = 9  // Manual Temperature (Kelvin)
	WhiteBalancePCSet1               WhiteBalance = 10 // PC Set 1
	WhiteBalancePCSet2               WhiteBalance = 11 // PC Set 2
	WhiteBalancePCSet3               WhiteBalance = 12 // PC Set 3
	WhiteBalanceDaylightFluorescent  WhiteBalance = 14 // Daylight Fluorescent
	WhiteBalanceCustom2              WhiteBalance = 15 // Custom 2
	WhiteBalanceUnderwater           WhiteBalance = 16 // Underwater
	WhiteBalanceCustom3              WhiteBalance = 18 // Custom 3
	WhiteBalancePCSet4               WhiteBalance = 19 // PC Set 4
	WhiteBalancePCSet5               WhiteBalance = 20 // PC Set 5
	WhiteBalanceAutoAmbiencePriority WhiteBalance = 21 // Auto (ambience priority)
	WhiteBalanceAutoWhitePriority    WhiteBalance = 23 // Auto (white priority)
)

// NewWhiteBalance returns a validated WhiteBalance, or 0 for unknown values.
func NewWhiteBalance(v int16) WhiteBalance {
	wb := WhiteBalance(v)
	switch wb {
	case WhiteBalanceAuto,
		WhiteBalanceDaylight,
		WhiteBalanceCloudy,
		WhiteBalanceTungsten,
		WhiteBalanceFluorescent,
		WhiteBalanceFlash,
		WhiteBalanceCustom,
		WhiteBalanceBlackAndWhite,
		WhiteBalanceShade,
		WhiteBalanceManualTemperature,
		WhiteBalancePCSet1,
		WhiteBalancePCSet2,
		WhiteBalancePCSet3,
		WhiteBalanceDaylightFluorescent,
		WhiteBalanceCustom2,
		WhiteBalanceUnderwater,
		WhiteBalanceCustom3,
		WhiteBalancePCSet4,
		WhiteBalancePCSet5,
		WhiteBalanceAutoAmbiencePriority,
		WhiteBalanceAutoWhitePriority:
		return wb
	default:
		return WhiteBalanceAuto
	}
}

// NewWhiteBalanceFromRaw decodes a raw uint16 wire value into WhiteBalance.
func NewWhiteBalanceFromRaw(raw uint16) WhiteBalance {
	return NewWhiteBalance(meta.SafecastUint16ToInt16Bits(raw))
}

// SlowShutter is part of the CanonShotInfo field.
type SlowShutter int16

const (
	SlowShutterNA         SlowShutter = -1 // n/a
	SlowShutterOff        SlowShutter = 0  // Off
	SlowShutterNightScene SlowShutter = 1  // Night Scene
	SlowShutterOn         SlowShutter = 2  // On
	SlowShutterNone       SlowShutter = 3  // None
)

// CameraType is part of the CanonShotInfo field.
type CameraType int16

const (
	CameraTypeEOSHighEnd CameraType = 248 // EOS High-end
	CameraTypeCompact    CameraType = 250 // Compact
	CameraTypeEOSMid     CameraType = 252 // EOS Mid-range
	CameraTypeDV         CameraType = 255 // DV Camera
)

// AutoRotate is part of the CanonShotInfo field.
type AutoRotate int16

const (
	AutoRotateNone        AutoRotate = 0 // None
	AutoRotateRotate90CW  AutoRotate = 1 // Rotate 90 CW
	AutoRotateRotate180   AutoRotate = 2 // Rotate 180
	AutoRotateRotate270CW AutoRotate = 3 // Rotate 270 CW
)

// NDFilter is part of the CanonShotInfo field.
type NDFilter int16

const (
	NDFilterOff NDFilter = 0 // Off
	NDFilterOn  NDFilter = 1 // On
)

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
//	1: "Electronic First Curtain"
//	2: "Electronic"
type ShutterMode uint16

const (
	ShutterModeMechanical             ShutterMode = 0 // Mechanical
	ShutterModeElectronicFirstCurtain ShutterMode = 1 // Electronic First Curtain
	ShutterModeElectronic             ShutterMode = 2 // Electronic
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
// Indexed by TimeZoneCity (0–23).
var cityGMTOffsets = [24]GMTOffset{
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
	if int(city) < 0 || int(city) >= len(cityGMTOffsets) {
		return GMTOffset{}, false
	}
	return cityGMTOffsets[city], true
}

// GMTHoursForCity returns the effective GMT offset in hours for a city.
func GMTHoursForCity(city TimeZoneCity, daylightSavings bool) (float64, bool) {
	if int(city) < 0 || int(city) >= len(cityGMTOffsets) {
		return 0, false
	}
	return cityGMTOffsets[city].Hours(daylightSavings), true
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
	TimeZoneCityKathmandu:   "Asia/Katmandu",
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

var cityLocationsOnce sync.Once
var cityLocations []*time.Location

func buildCityLocations() {
	const cityCount = int(TimeZoneCityDubai) + 1
	cityLocations = make([]*time.Location, cityCount)

	for city, name := range cityIANA {
		if name == "" {
			continue
		}

		loc, err := time.LoadLocation(name)
		if err != nil {
			continue
		}
		cityLocations[city] = loc
	}
}

// LocationForCity returns the cached IANA location for a Canon timezone city.
func LocationForCity(city TimeZoneCity) (*time.Location, bool) {
	cityLocationsOnce.Do(buildCityLocations)
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
//	3: "AE Lock + Exposure Comp.",
//	4: "No AE",
type AESetting int16

const (
	AESettingNormalAE                       AESetting = 0 // Normal AE
	AESettingExposureCompensation           AESetting = 1 // Exposure Compensation
	AESettingAELock                         AESetting = 2 // AE Lock
	AESettingAELockWithExposureCompensation AESetting = 3 // AE Lock + Exposure Comp.
	AESettingNoAE                           AESetting = 4 // No AE
)

// AFPointSetting is the stored value for Canon CameraSettings AFPoint (seq 19).
type AFPointSetting uint16

// ExifTool Canon.pm AFPoint PrintConv values.
const (
	AFPointManualSelection AFPointSetting = 0x2005 // Manual AF point selection
	AFPointNoneMF          AFPointSetting = 0x3000 // None (MF)
	AFPointAutoSelection   AFPointSetting = 0x3001 // Auto AF point selection
	AFPointRight           AFPointSetting = 0x3002 // Right
	AFPointCenter          AFPointSetting = 0x3003 // Center
	AFPointLeft            AFPointSetting = 0x3004 // Left
	AFPointAutoAlt         AFPointSetting = 0x4001 // Auto AF point selection
	AFPointFaceDetect      AFPointSetting = 0x4006 // Face Detect
)

var afPointLabels = map[AFPointSetting]string{
	AFPointManualSelection: "Manual AF point selection",
	AFPointNoneMF:          "None (MF)",
	AFPointAutoSelection:   "Auto AF point selection",
	AFPointRight:           "Right",
	AFPointCenter:          "Center",
	AFPointLeft:            "Left",
	AFPointAutoAlt:         "Auto AF point selection",
	AFPointFaceDetect:      "Face Detect",
}

// AFPointString returns the ExifTool-style display string for an AFPoint value.
func AFPointString(v uint16) string {
	if s, ok := afPointLabels[AFPointSetting(v)]; ok {
		return s
	}
	return ""
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

func NewAFAreaModeFromRAW(raw uint16) AFAreaMode {
	return AFAreaMode(meta.SafecastUint16ToInt16Bits(raw))
}

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
	AFAreaModeLargeZoneAFVert   AFAreaMode = 16 // Large Zone AF (vertical)
	AFAreaModeLargeZoneAFHoriz  AFAreaMode = 17 // Large Zone AF (horizontal)
	AFAreaModeFlexibleZoneAF1   AFAreaMode = 19 // Flexible Zone AF 1
	AFAreaModeFlexibleZoneAF2   AFAreaMode = 20 // Flexible Zone AF 2
	AFAreaModeFlexibleZoneAF3   AFAreaMode = 21 // Flexible Zone AF 3
	AFAreaModeWholeAreaAF       AFAreaMode = 22 // Whole Area AF
)
