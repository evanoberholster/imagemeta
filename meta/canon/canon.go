// Package canon provides data types and functions for representing Canon Camera Makernote values
package canon

//go:generate msgp
//go:generate stringer -type=ContinuousDrive,FocusMode,MeteringMode,FocusRange,ExposureMode,BracketMode,AESetting,AFAreaMode -linecomment -output=canon_string.go

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
	ContinuousDriveSingle                  ContinuousDrive = 0  // Single
	ContinuousDriveContinuous              ContinuousDrive = 1  // Continuous
	ContinuousDriveMovie                   ContinuousDrive = 2  // Movie
	ContinuousDriveContinuousSpeedPriority ContinuousDrive = 3  // Continuous, Speed Priority
	ContinuousDriveContinuousLow           ContinuousDrive = 4  // Continuous, Low
	ContinuousDriveContinuousHigh          ContinuousDrive = 5  // Continuous, High
	ContinuousDriveSilentSingle            ContinuousDrive = 6  // Silent Single
	ContinuousDriveUnknown7                ContinuousDrive = 7  // Unknown
	ContinuousDriveUnknown8                ContinuousDrive = 8  // Unknown
	ContinuousDriveSingleSilent            ContinuousDrive = 9  // Single, Silent
	ContinuousDriveContinuousSilent        ContinuousDrive = 10 // Continuous, Silent
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

// AESetting - Canon Makernote AutoExposure Setting
//
//	0: "Normal AE",
//	1: "Exposure Compensation",
//	2: "AE Lock",
//	3: "AE Lock + Exposure Compensation",
//	4: "No AE",
type AESetting int16

const (
	AESettingNormalAE                       AESetting = 0 // Normal AE
	AESettingExposureCompensation           AESetting = 1 // Exposure Compensation
	AESettingAELock                         AESetting = 2 // AE Lock
	AESettingAELockWithExposureCompensation AESetting = 3 // AE Lock + Exposure Compensation
	AESettingNoAE                           AESetting = 4 // No AE
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
