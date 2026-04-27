package sony

import "github.com/evanoberholster/imagemeta/meta"

// SonyAFStatus15 stores the 18-point Sony AF status matrix.
type SonyAFStatus15 struct {
	UpperLeft        int16
	Left             int16
	LowerLeft        int16
	FarLeft          int16
	TopHorizontal    int16
	NearRight        int16
	CenterHorizontal int16
	NearLeft         int16
	BottomHorizontal int16
	TopVertical      int16
	CenterVertical   int16
	BottomVertical   int16
	FarRight         int16
	UpperRight       int16
	Right            int16
	LowerRight       int16
	UpperMiddle      int16
	LowerMiddle      int16
}

// SonyCameraInfo3 stores maker-note tag 0x0010 (CameraInfo3).
//
// Field comments use the byte offset for this table.
type SonyCameraInfo3 struct {
	LensSpec             string         // [0x0000] LensSpec (8 bytes)
	FocalLength          uint16         // [0x000e] FocalLength
	FocalLengthTeleZoom  uint16         // [0x0010] FocalLengthTeleZoom
	FocusStatus          int16          // [0x0019] FocusStatus
	AFPointSelected      uint8          // [0x001c] AFPointSelected
	FocusMode            uint8          // [0x001d] FocusMode
	AFPoint              uint8          // [0x0020] AFPoint
	AFStatusActiveSensor int16          // [0x0021] AFStatusActiveSensor
	AFStatus15           SonyAFStatus15 // [0x0023] AFStatus15 (18-point matrix)
}

// SonyCameraInfo2 stores maker-note tag 0x0010 legacy layouts.
//
// Field comments use the byte offset for this table.
type SonyCameraInfo2 struct {
	LensSpec                 string // [0x0000] LensSpec (8 bytes)
	AFPointSelected          uint8  // [0x0014] AFPointSelected
	FocusModeSetting         uint8  // [0x0015] FocusModeSetting
	AFPoint                  uint8  // [0x0018] AFPoint
	AFStatusActiveSensor     int16  // [0x001b] AFStatusActiveSensor
	AFStatusTopRight         int16  // [0x001d] AFStatusTopRight
	AFStatusBottomRight      int16  // [0x001f] AFStatusBottomRight
	AFStatusBottom           int16  // [0x0021] AFStatusBottom
	AFStatusMiddleHorizontal int16  // [0x0023] AFStatusMiddleHorizontal
	AFStatusCenterVertical   int16  // [0x0025] AFStatusCenterVertical
	AFStatusTop              int16  // [0x0027] AFStatusTop
	AFStatusTopLeft          int16  // [0x0029] AFStatusTopLeft
	AFStatusBottomLeft       int16  // [0x002b] AFStatusBottomLeft
	AFStatusLeft             int16  // [0x002d] AFStatusLeft
	AFStatusCenterHorizontal int16  // [0x002f] AFStatusCenterHorizontal
	AFStatusRight            int16  // [0x0031] AFStatusRight
}

// SonyFocusInfo stores maker-note tag 0x0020 legacy layouts.
//
// Field comments use the byte offset for this table.
type SonyFocusInfo struct {
	DriveMode2                   uint8 // [0x0e] DriveMode2
	Rotation                     uint8 // [0x10] Rotation
	ImageStabilizationSetting    uint8 // [0x14] ImageStabilizationSetting
	DynamicRangeOptimizerMode    uint8 // [0x15] DynamicRangeOptimizerMode
	BracketShotNumber            uint8 // [0x2b] BracketShotNumber
	WhiteBalanceBracketing       uint8 // [0x2c] WhiteBalanceBracketing
	BracketShotNumber2           uint8 // [0x2d] BracketShotNumber2
	DynamicRangeOptimizerBracket uint8 // [0x2e] DynamicRangeOptimizerBracket
	ExposureBracketShotNumber    uint8 // [0x2f] ExposureBracketShotNumber
	ExposureProgram              uint8 // [0x3f] ExposureProgram
	CreativeStyle                uint8 // [0x41] CreativeStyle
	ISOSetting                   uint8 // [0x6d] ISOSetting
	ISO                          uint8 // [0x6f] ISO
	DynamicRangeOptimizerMode2   uint8 // [0x77] DynamicRangeOptimizerMode2
	DynamicRangeOptimizerLevel   uint8 // [0x79] DynamicRangeOptimizerLevel
	FocusPosition                uint8 // [0x9bb] FocusPosition
}

// SonyFaceInfo stores the visible Sony face count.
type SonyFaceInfo struct {
	FacesDetected uint16
}

// SonyMoreSettings stores the nested 0x0001 block inside Sony MoreInfo.
type SonyMoreSettings struct {
	DriveMode2                   uint8
	ExposureProgram              uint8
	MeteringMode                 uint8
	DynamicRangeOptimizerSetting uint8
	DynamicRangeOptimizerLevel   uint8
	ColorSpace                   uint8
	CreativeStyleSetting         uint8
	ContrastSetting              int8
	SaturationSetting            int8
	SharpnessSetting             int8
	WhiteBalanceSetting          uint8
	ColorTemperatureSetting      uint8
	ColorCompensationFilterSet   int8
	FlashMode                    uint8
	LongExposureNoiseReduction   uint8
	HighISONoiseReduction        uint8
	FocusMode                    uint8
	MultiFrameNoiseReduction     uint8
	HDRSetting                   uint8
	HDRLevel                     uint8
	ViewingMode                  uint8
	FaceDetection                uint8
	CustomWB_RBLevels            [2]uint16
	BrightnessValue              uint8
	ExposureCompensationSet      uint8
	FlashExposureCompSet         uint8
	LiveViewAFMethod             uint8
	ISO                          uint8
	FNumber                      uint8
	ExposureTime                 uint8
	FocalLength2                 uint8
	ExposureCompensation2        int16
	FlashExposureCompSet2        int16
	Orientation2                 uint8
	FocusPosition2               uint8
	FlashAction                  uint8
	FocusMode2                   uint8
	FlashActionExternal          uint8
	FlashStatus                  uint8
}

// SonyMoreInfo stores maker-note tag 0x0020 modern layouts.
type SonyMoreInfo struct {
	MoreSettings           SonyMoreSettings
	FaceInfo               SonyFaceInfo
	ImageCount             uint32
	ShutterCount           uint32
	ShotNumberSincePowerUp uint32
}

// SonyCameraSettings stores maker-note tag 0x0114 legacy layouts.
//
// Field comments use the byte offset for the 280-byte layout.
type SonyCameraSettings struct {
	ExposureTime                  uint16   // [0x00] ExposureTime
	FNumber                       uint16   // [0x01] FNumber
	HighSpeedSync                 uint16   // [0x02] HighSpeedSync
	ExposureCompensationSet       uint16   // [0x03] ExposureCompensationSet
	DriveMode                     uint16   // [0x04] DriveMode
	WhiteBalanceSetting           uint16   // [0x05] WhiteBalanceSetting
	WhiteBalanceFineTune          int16    // [0x06] WhiteBalanceFineTune
	ColorTemperatureSet           uint16   // [0x07] ColorTemperatureSet
	ColorCompensationFilterSet    int16    // [0x08] ColorCompensationFilterSet
	CustomWB_RGBLevels            [3]uint16 // [0x09] CustomWB_RGBLevels
	ColorTemperatureCustom        uint16   // [0x0c] ColorTemperatureCustom
	ColorCompensationFilterCustom int16    // [0x0d] ColorCompensationFilterCustom
	WhiteBalance                  uint16   // [0x0f/0x0e] WhiteBalance
	FocusModeSetting              uint16   // [0x10/0x0f] FocusModeSetting
	AFAreaMode                    uint16   // [0x11/0x10] AFAreaMode
	AFPointSetting                uint16   // [0x12/0x11] AFPointSetting
	FlashMode                     uint16   // [0x13/0x00] FlashMode (280/332)
	FlashExposureCompSet          uint16   // [0x14/0x12] FlashExposureCompSet
	MeteringMode                  uint16   // [0x15/0x13] MeteringMode
	ISOSetting                    uint16   // [0x16/0x14] ISOSetting
	DynamicRangeOptimizerMode     uint16   // [0x18/0x16] DynamicRangeOptimizerMode
	DynamicRangeOptimizerLevel    uint16   // [0x19/0x17] DynamicRangeOptimizerLevel
	CreativeStyle                 uint16   // [0x1a/0x18] CreativeStyle
	ColorSpace                    uint16   // [0x1b/0x??] ColorSpace
	Sharpness                     int16    // [0x1c/0x19] Sharpness
	Contrast                      int16    // [0x1d/0x1a] Contrast
	Saturation                    int16    // [0x1e/0x1b] Saturation
	ZoneMatchingValue             uint16   // [0x1f] ZoneMatchingValue
	Brightness                    int16    // [0x22] Brightness
	FlashControl                  uint16   // [0x23/0x1f] FlashControl
	PrioritySetupShutterRelease   uint16   // [0x28] PrioritySetupShutterRelease
	AFIlluminator                 uint16   // [0x29] AFIlluminator
	AFWithShutter                 uint16   // [0x2a] AFWithShutter
	LongExposureNoiseReduction    uint16   // [0x2b/0x25] LongExposureNoiseReduction
	HighISONoiseReduction         uint16   // [0x2c/0x26] HighISONoiseReduction
	ImageStyle                    uint16   // [0x2d/0x27] ImageStyle
	FocusModeSwitch               uint16   // [0x2e] FocusModeSwitch
	ShutterSpeedSetting           uint16   // [0x2f/0x28] ShutterSpeedSetting
	ApertureSetting               uint16   // [0x30/0x29] ApertureSetting
	ExposureProgram               uint16   // [0x3c] ExposureProgram
	ImageStabilizationSetting     uint16   // [0x3d] ImageStabilizationSetting
	FlashAction                   uint16   // [0x3e] FlashAction
	Rotation                      uint16   // [0x3f] Rotation
	AELock                        uint16   // [0x40] AELock
	FlashAction2                  uint16   // [0x4c] FlashAction2
	FocusMode                     uint16   // [0x4d] FocusMode
	BatteryState                  uint16   // [0x50] BatteryState
	BatteryLevel                  uint16   // [0x51] BatteryLevel
	FocusStatus                   uint16   // [0x53] FocusStatus
	SonyImageSize                 uint16   // [0x54] SonyImageSize
	AspectRatio                   uint16   // [0x55] AspectRatio
	Quality                       uint16   // [0x56] Quality
	ExposureLevelIncrements       uint16   // [0x58] ExposureLevelIncrements
	RedEyeReduction               uint16   // [0x6a] RedEyeReduction
}

// SonyCameraSettings3 stores maker-note tag 0x0114 modern 1-byte layouts.
//
// Field comments use the byte offset for this table.
type SonyCameraSettings3 struct {
	ShutterSpeedSetting          uint8    // [0x00] ShutterSpeedSetting
	ApertureSetting              uint8    // [0x01] ApertureSetting
	ISOSetting                   uint8    // [0x02] ISOSetting
	ExposureCompensationSet      uint8    // [0x03] ExposureCompensationSet
	DriveModeSetting             uint8    // [0x04] DriveModeSetting
	ExposureProgram              uint8    // [0x05] ExposureProgram
	FocusModeSetting             uint8    // [0x06] FocusModeSetting
	MeteringMode                 uint8    // [0x07] MeteringMode
	SonyImageSize                uint8    // [0x09] SonyImageSize
	AspectRatio                  uint8    // [0x0a] AspectRatio
	Quality                      uint8    // [0x0b] Quality
	DynamicRangeOptimizerSetting uint8    // [0x0c] DynamicRangeOptimizerSetting
	DynamicRangeOptimizerLevel   uint8    // [0x0d] DynamicRangeOptimizerLevel
	ColorSpace                   uint8    // [0x0e] ColorSpace
	CreativeStyleSetting         uint8    // [0x0f] CreativeStyleSetting
	ContrastSetting              int8     // [0x10] ContrastSetting
	SaturationSetting            int8     // [0x11] SaturationSetting
	SharpnessSetting             int8     // [0x12] SharpnessSetting
	WhiteBalanceSetting          uint8    // [0x16] WhiteBalanceSetting
	ColorTemperatureSetting      uint8    // [0x17] ColorTemperatureSetting
	ColorCompensationFilterSet   int8     // [0x18] ColorCompensationFilterSet
	CustomWB_RGBLevels           [3]uint16 // [0x19] CustomWB_RGBLevels
	FlashMode                    uint8    // [0x20] FlashMode
	FlashControl                 uint8    // [0x21] FlashControl
	FlashExposureCompSet         uint8    // [0x23] FlashExposureCompSet
	AFAreaMode                   uint8    // [0x24] AFAreaMode
	LongExposureNoiseReduction   uint8    // [0x25] LongExposureNoiseReduction
	HighISONoiseReduction        uint8    // [0x26] HighISONoiseReduction
	SmileShutterMode             uint8    // [0x27] SmileShutterMode
	RedEyeReduction              uint8    // [0x28] RedEyeReduction
	HDRSetting                   uint8    // [0x2d] HDRSetting
	HDRLevel                     uint8    // [0x2e] HDRLevel
	ViewingMode                  uint8    // [0x2f] ViewingMode
	FaceDetection                uint8    // [0x30] FaceDetection
	SmileShutter                 uint8    // [0x31] SmileShutter
	SweepPanoramaSize            uint8    // [0x32] SweepPanoramaSize
	SweepPanoramaDirection       uint8    // [0x33] SweepPanoramaDirection
	DriveMode                    uint8    // [0x34] DriveMode
	MultiFrameNoiseReduction     uint8    // [0x35] MultiFrameNoiseReduction
	LiveViewAFSetting            uint8    // [0x36] LiveViewAFSetting
	PanoramaSize3D               uint8    // [0x38] PanoramaSize3D
	AFButtonPressed              uint8    // [0x83] AFButtonPressed
	LiveViewMetering             uint8    // [0x84] LiveViewMetering
	ViewingMode2                 uint8    // [0x85] ViewingMode2
	AELock                       uint8    // [0x86] AELock
	FlashStatusBuiltIn           uint8    // [0x87] FlashStatusBuiltIn
	FlashStatusExternal          uint8    // [0x88] FlashStatusExternal
	LiveViewFocusMode            uint8    // [0x8b] LiveViewFocusMode
	LensMount                    uint8    // [0x99] LensMount
	SequenceNumber               uint8    // [0x10c] SequenceNumber
	FolderNumber                 uint16   // [0x114] FolderNumber
	ImageNumber                  uint16   // [0x114] ImageNumber
	ShotNumberSincePowerUp2      uint32   // [0x200] ShotNumberSincePowerUp2
}

// SonyShotInfo stores maker-note tag 0x3000 (ShotInfo).
//
// Field comments use the byte offset for this table.
type SonyShotInfo struct {
	FaceInfoOffset  uint16 // [0x02] FaceInfoOffset
	SonyDateTime    string // [0x06] SonyDateTime (20 chars)
	SonyImageHeight uint16 // [0x1a] SonyImageHeight
	SonyImageWidth  uint16 // [0x1c] SonyImageWidth
	FacesDetected   uint16 // [0x30] FacesDetected
	FaceInfoLength  uint16 // [0x32] FaceInfoLength
	MetaVersion     string // [0x34] MetaVersion (16 chars)
}

// SonyTag9400 stores maker-note tag 0x9400.
type SonyTag9400 struct {
	SequenceImageNumber    uint32
	SequenceFileNumber     uint32
	ReleaseMode2           uint8
	ShotNumberSincePowerUp uint32
	SequenceLength         uint8
	CameraOrientation      uint8
	Quality2               uint8
	SonyImageHeight        uint16
	ModelReleaseYear       uint8
}

// SonyTag9404 stores maker-note tag 0x9404.
type SonyTag9404 struct {
	ExposureProgram uint8
	IntelligentAuto uint8
}

// SonyTag9405 stores maker-note tag 0x9405 (lens correction parameters).
//
// Field comments use the byte offset for this table.
type SonyTag9405 struct {
	DistortionCorrParamsPresent   uint8     // [0x0600] DistortionCorrParamsPresent
	DistortionCorrection          uint8     // [0x0601] DistortionCorrection
	LensFormat                    uint8     // [0x0603] LensFormat
	LensMount                     uint8     // [0x0604] LensMount
	LensType                      uint16    // [0x0608] LensType
	VignettingCorrParams          [16]int16 // [0x064a] VignettingCorrParams
	ChromaticAberrationCorrParams [32]int16 // [0x066a] ChromaticAberrationCorrParams
	DistortionCorrParams          [16]int16 // [0x06ca] DistortionCorrParams
}

// SonyTag9406 stores maker-note tag 0x9406.
type SonyTag9406 struct {
	BatteryTemperature uint8
	BatteryLevelGrip1  uint8
	BatteryLevel       uint8
	BatteryLevelGrip2  uint8
}

// SonyTag940A stores maker-note tag 0x940a.
type SonyTag940A struct {
	AFPointsSelected uint32
}

// SonyTag940C stores maker-note tag 0x940c.
type SonyTag940C struct {
	LensMount2          uint8
	LensType3           uint16
	CameraEMountVersion uint16
	LensEMountVersion   uint16
	LensFirmwareVersion uint16
}

// SonyTag2010 stores maker-note tag 0x2010.
type SonyTag2010 struct {
	SequenceImageNumber   uint32
	SequenceFileNumber    uint32
	ReleaseMode2          uint32
	SonyDateTime          string
	DynamicRangeOptimizer uint8
	ReleaseMode3          uint8
	SelfTimer             uint8
	FlashMode             uint8
	StopsAboveBaseISO     uint16
	BrightnessValue       uint16
	HDRSetting            uint8
	ExposureCompensation  int16
	PictureProfile        uint8
	PictureEffect2        uint8
	Quality2              uint8
	MeteringMode          uint8
	ExposureProgram       uint8
	WB_RGBLevels          [3]uint16
	SonyISO               uint16
}

// SonyTag202A stores maker-note tag 0x202a.
type SonyTag202A struct {
	FocalPlaneAFPointsUsed uint8
}

// SonyHiddenInfo stores maker-note tag 0x2044.
type SonyHiddenInfo struct {
	HiddenDataOffset uint32
	HiddenDataLength uint32
}

// SonyTag9050 stores maker-note tag 0x9050.
//
// Field comments use the byte offset for this table.
type SonyTag9050 struct {
	SonyMaxAperture             uint8     // [0x00] SonyMaxAperture
	SonyMinAperture             uint8     // [0x01] SonyMinAperture
	Shutter                     [3]uint16 // [0x20] Shutter (3 values)
	FlashStatus                 uint8     // [0x31] FlashStatus
	ShutterCount                uint32    // [0x32] ShutterCount
	SonyExposureTime            uint16    // [0x3a] SonyExposureTime
	SonyFNumber                 uint16    // [0x3c] SonyFNumber
	ReleaseMode2                uint8     // [0x3f] ReleaseMode2
	InternalSerialNumber        [5]uint8  // [0xf0] InternalSerialNumber
	LensMount                   uint8     // [0x105] LensMount
	LensFormat                  uint8     // [0x106] LensFormat
	LensType                    uint16    // [0x109] LensType
	DistortionCorrParamsPresent uint8     // [0x10b] DistortionCorrParamsPresent
	LensSpecFeatures            string    // [0x115] LensSpecFeatures
	ShutterCount3               uint32    // [0x1bd] ShutterCount3
}

// SonyTag9416 stores maker-note tag 0x9416.
//
// Field comments use the byte offset for this table.
type SonyTag9416 struct {
	SonyISO                       uint16    // [0x04] SonyISO
	StopsAboveBaseISO             uint16    // [0x06] StopsAboveBaseISO
	SonyExposureTime2             uint16    // [0x0a] SonyExposureTime2
	ExposureTime                  meta.ExposureTime // [0x0c] ExposureTime (rational32u)
	SonyFNumber2                  uint16    // [0x10] SonyFNumber2
	SonyMaxApertureValue          uint16    // [0x12] SonyMaxApertureValue
	SequenceImageNumber           uint32    // [0x1d] SequenceImageNumber
	ReleaseMode2                  uint8     // [0x2b] ReleaseMode2
	InternalSerialNumber          [6]uint8  // [0x38] InternalSerialNumber
	ExposureProgram               uint8     // [0x35] ExposureProgram
	CreativeStyle                 uint8     // [0x37] CreativeStyle
	LensMount                     uint8     // [0x48] LensMount
	LensFormat                    uint8     // [0x49] LensFormat
	LensType2                     uint16    // [0x4b] LensType2
	DistortionCorrParams          [16]int16 // [0x4f] DistortionCorrParams
	PictureProfile                uint8     // [0x70] PictureProfile
	FocalLength                   uint16    // [0x71] FocalLength
	MinFocalLength                uint16    // [0x73] MinFocalLength
	MaxFocalLength                uint16    // [0x75] MaxFocalLength
	VignettingCorrParams          [32]int16 // [0x89d] VignettingCorrParams
	APSCSizeCapture               uint8     // [0x8e5] APSCSizeCapture
	ChromaticAberrationCorrParams [32]int16 // [0x945] ChromaticAberrationCorrParams
}

// SonyAFInfo stores maker-note tag 0x940e (AFInfo).
//
// Field comments use the byte offset for this table.
type SonyAFInfo struct {
	AFType                  uint8         // [0x02] AFType
	AFStatusActiveSensor    int16         // [0x04] AFStatusActiveSensor
	AFPoint                 uint8         // [0x07] AFPoint
	AFPointInFocus          uint8         // [0x08] AFPointInFocus
	AFPointAtShutterRelease uint8         // [0x09] AFPointAtShutterRelease
	AFAreaMode              uint8         // [0x0a] AFAreaMode
	FocusMode               uint8         // [0x0b] FocusMode
	AFStatus15              SonyAFStatus15 // [0x11] AFStatus15 (18-point matrix)
	AFPointsUsed            uint32        // [0x16e] AFPointsUsed
	AFMicroAdj              int8          // [0x17d] AFMicroAdj
	ExposureProgram         uint8         // [0x17e] ExposureProgram
}
