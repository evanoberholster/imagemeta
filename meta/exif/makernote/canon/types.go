package canon

import "github.com/evanoberholster/imagemeta/meta"

// CameraSettings stores Canon CameraSettings values from MakerNote tag
// 0x0001 (CanonCameraSettings).
//
// Field comments use ExifTool sequence indices for this table.
type CameraSettings struct {
	MacroMode          MacroMode          // [1] MacroMode
	SelfTimer          int16              // [2] SelfTimer (deci-seconds; custom bit may be set)
	Quality            Quality            // [3] Quality
	CanonFlashMode     CanonFlashMode     // [4] CanonFlashMode
	ContinuousDrive    ContinuousDrive    // [5] ContinuousDrive
	FocusMode          FocusMode          // [7] FocusMode
	RecordMode         RecordMode         // [9] RecordMode
	CanonImageSize     CanonImageSize     // [10] CanonImageSize
	EasyMode           EasyMode           // [11] EasyMode
	DigitalZoom        DigitalZoom        // [12] DigitalZoom
	Contrast           int16              // [13] Contrast
	Saturation         int16              // [14] Saturation
	Sharpness          int16              // [15] Sharpness
	CameraISO          CameraISO          // [16] CameraISO
	MeteringMode       MeteringMode       // [17] MeteringMode
	FocusRange         FocusRange         // [18] FocusRange
	AFPoint            uint16             // [19] AFPoint
	CanonExposureMode  ExposureMode       // [20] CanonExposureMode
	LensType           CanonLensType      // [22] LensType
	MaxFocalLength     uint16             // [23] MaxFocalLength
	MinFocalLength     uint16             // [24] MinFocalLength
	FocalUnits         uint16             // [25] FocalUnits
	MaxAperture        meta.Aperture      // [26] MaxAperture
	MinAperture        meta.Aperture      // [27] MinAperture
	FlashModel         FlashModel         // [28] FlashModel
	FlashActivity      FlashModel         // [28] FlashActivity (backward-compatible alias)
	FlashBits          uint16             // [29] FlashBits
	FocusContinuous    FocusContinuous    // [32] FocusContinuous
	AESetting          AESetting          // [33] AESetting
	ImageStabilization ImageStabilization // [34] ImageStabilization
	DisplayAperture    meta.Aperture      // [35] DisplayAperture (x10)
	ZoomSourceWidth    uint16             // [36] ZoomSourceWidth
	ZoomTargetWidth    uint16             // [37] ZoomTargetWidth
	SpotMeteringMode   SpotMeteringMode   // [39] SpotMeteringMode
	PhotoEffect        PhotoEffect        // [40] PhotoEffect
	ManualFlashOutput  ManualFlashOutput  // [41] ManualFlashOutput
	ColorTone          int16              // [42] ColorTone
	SRAWQuality        SRAWQuality        // [46] SRAWQuality
	FocusBracketing    FocusBracketing    // [50] FocusBracketing
	Clarity            int16              // [51] Clarity (EOS R models)
	HDRPQ              HDRPQ              // [52] HDR_PQ
}

// ShotInfo stores selected Canon ShotInfo values from MakerNote tag 0x0004
// (CanonShotInfo).
//
// Field comments use ExifTool sequence indices for this table.
type ShotInfo struct {
	AutoISO                 int16             // [1] AutoISO raw code
	AutoISOValue            float32           // [1] ExifTool-converted AutoISO
	ActualISO               float32           // [1-2] ExifTool: BaseISO * AutoISO / 100
	BaseISO                 int16             // [2] BaseISO raw code
	BaseISOValue            float32           // [2] ExifTool-converted BaseISO
	MeasuredEV              int16             // [3] MeasuredEV
	TargetAperture          int16             // [4] TargetAperture
	TargetApertureValue     meta.Aperture     // [4] ExifTool CanonEv-converted aperture
	TargetExposureTime      int16             // [5] TargetExposureTime
	TargetExposureTimeValue meta.ExposureTime // [5] ExifTool-converted exposure time
	ExposureCompensation    int16             // [6] ExposureCompensation
	WhiteBalance            WhiteBalance      // [7] WhiteBalance
	SlowShutter             SlowShutter       // [8] SlowShutter
	SequenceNumber          int16             // [9] SequenceNumber
	OpticalZoomCode         int16             // [10] OpticalZoomCode
	CameraTemperature       int16             // [12] CameraTemperature raw code
	CameraTemperatureC      int16             // [12] ExifTool-converted Celsius
	FlashGuideNumber        int16             // [13] FlashGuideNumber raw code
	FlashGuideNumberMeters  float32           // [13] ExifTool-converted meters
	AFPointsInFocus         uint16            // [14] AFPointsInFocus bitset
	FlashExposureComp       int16             // [15] FlashExposureCompensation
	AutoExposureBracketing  int16             // [16] AutoExposureBracketing
	AEBBracketValue         int16             // [17] AEBBracketValue
	ControlMode             int16             // [18] ControlMode
	FocusDistance           FocusDistance     // [19-20] FocusDistanceUpper/Lower
	FNumber                 int16             // [21] FNumber raw code
	FNumberValue            meta.Aperture     // [21] ExifTool CanonEv-converted aperture
	ExposureTime            int16             // [22] ExposureTime raw code
	ExposureTimeValue       meta.ExposureTime // [22] ExifTool-converted exposure time
	MeasuredEV2             int16             // [23] MeasuredEV2
	BulbDuration            int16             // [24] BulbDuration (deci-seconds)
	CameraType              CameraType        // [26] CameraType
	AutoRotate              AutoRotate        // [27] AutoRotate
	NDFilter                NDFilter          // [28] NDFilter
	SelfTimer2              int16             // [29] SelfTimer2 (deci-seconds)
	FlashOutput             int16             // [33] FlashOutput
}

// FileInfo stores selected Canon FileInfo values from MakerNote tag 0x0093
// (CanonFileInfo).
//
// Field comments use ExifTool sequence indices for this table.
type FileInfo struct {
	FileNumber                  uint32          // [1] FileNumber
	ShutterCount                uint32          // [1] ShutterCount (model-dependent variant)
	BracketMode                 BracketMode     // [3] BracketMode
	BracketValue                int16           // [4] BracketValue
	BracketShotNumber           int16           // [5] BracketShotNumber
	RawJpgQuality               RawJpgQuality   // [6] RawJpgQuality
	RawJpgSize                  RawJpgSize      // [7] RawJpgSize
	LongExposureNoiseReduction2 OnOffAuto       // [8] LongExposureNoiseReduction2
	WBBracketMode               int16           // [9] WBBracketMode
	WBBracketValueAB            int16           // [12] WBBracketValueAB
	WBBracketValueGM            int16           // [13] WBBracketValueGM
	FilterEffect                FilterEffect    // [14] FilterEffect
	ToningEffect                ToningEffect    // [15] ToningEffect
	MacroMagnification          int16           // [16] MacroMagnification
	LiveViewShooting            OnOffAuto       // [19] LiveViewShooting
	FocusDistance               FocusDistance   // [20-21] FocusDistanceUpper/Lower
	ShutterMode                 ShutterMode     // [23] ShutterMode
	FlashExposureLock           OnOffAuto       // [25] FlashExposureLock
	AntiFlicker                 OnOffAuto       // [32] AntiFlicker
	RFLensType                  CanonRFLensType // [0x3d] RFLensType
}

// CanonTimeInfo stores selected Canon TimeInfo values from MakerNote tag 0x0035
// (TimeInfo).
//
// Field comments use ExifTool sequence indices for this table.
type CanonTimeInfo struct {
	TimeZone        int32           // [1] TimeZone
	TimeZoneCity    TimeZoneCity    // [2] TimeZoneCity
	DaylightSavings DaylightSavings // [3] DaylightSavings
}

// AFInfo stores selected Canon autofocus record values from MakerNote tags
// 0x0012 (AFInfo), 0x0026 (AFInfo2), and 0x003c (AFInfo3).
type AFInfo struct {
	Source AFInfoSource // source Canon maker-note table

	AFAreaMode       AFAreaMode // AFInfo2 [1]
	NumAFPoints      uint16     // AFInfo [0] / AFInfo2 [2]
	ValidAFPoints    uint16     // AFInfo [1] / AFInfo2 [3]
	CanonImageWidth  uint16     // AFInfo [2] / AFInfo2 [4]
	CanonImageHeight uint16     // AFInfo [3] / AFInfo2 [5]
	AFImageWidth     uint16     // AFInfo [4] / AFInfo2 [6]
	AFImageHeight    uint16     // AFInfo [5] / AFInfo2 [7]
	AFAreaWidth      uint16     // AFInfo [6]
	AFAreaHeight     uint16     // AFInfo [7]

	// AFArea mirrors ExifTool's raw Canon AF area tuples: width, height,
	// x-position, and y-position in Canon AF coordinates.
	AFArea []AFPoint

	AFPointsInFocusBits  []int  // AFInfo [10] / AFInfo2 [12]
	AFPointsSelectedBits []int  // AFInfo2 [13]
	PrimaryAFPoint       uint16 // AFInfo [11/12] / AFInfo2 [14]

	// AFPoints is a derived convenience view for drawing AF rectangles.
	// AFInfo2/AFInfo3 converts Canon center-based coordinates to image-space
	// top-left rectangles. Legacy AFInfo exposes the raw tuples unchanged.
	AFPoints []AFPoint
}

// AFInfoSource identifies which Canon maker-note AF table produced AFInfo.
type AFInfoSource uint8

const (
	AFInfoSourceUnknown AFInfoSource = iota
	AFInfoSourceAFInfo
	AFInfoSourceAFInfo2
	AFInfoSourceAFInfo3
)

// FacePosition stores a face center point in FaceDetect frame coordinates.
type FacePosition struct {
	X int16
	Y int16
}

// FaceDetect1Info stores Canon MakerNote tag 0x0024 (FaceDetect1).
// Field comments use ExifTool sequence indices.
type FaceDetect1Info struct {
	FacesDetected       uint16          // [2] FacesDetected
	FaceDetectFrameSize [2]uint16       // [3] FaceDetectFrameSize
	FacePositions       [9]FacePosition // [8..25] Face1..Face9 position pairs
}

// FaceDetect2Info stores Canon MakerNote tag 0x0025 (FaceDetect2).
// Field comments use ExifTool sequence indices.
type FaceDetect2Info struct {
	FaceWidth     uint8 // [1] FaceWidth
	FacesDetected uint8 // [2] FacesDetected
}

// FaceDetect3Info stores Canon MakerNote tag 0x002f (FaceDetect3).
// Field comments use ExifTool sequence indices.
type FaceDetect3Info struct {
	FacesDetected uint16 // [3] FacesDetected
}

// AFPoint stores width, height, x, and y values for a Canon AF area tuple.
//
// In AFInfo.AFArea, x/y are the raw Canon AF coordinates. In AFInfo.AFPoints,
// x/y may be derived rectangle coordinates depending on the source table.
type AFPoint [4]int16

// NewAFPoint returns a new AFPoint from
// width, height, x-axis coord and y-axis coord
func NewAFPoint(w, h, x, y int16) AFPoint {
	return AFPoint{w, h, x, y}
}

// CanonRawPreviewLen is the maximum raw-byte preview stored for opaque
// maker-note blocks to avoid retaining large payloads.
const CanonRawPreviewLen = 64

// BlockPreview stores size and a short preview for opaque maker-note blocks.
type BlockPreview struct {
	Size         uint32
	Preview      [CanonRawPreviewLen]byte
	PreviewCount uint8
}

// FlashInfo stores Canon MakerNote tag 0x0003 (CanonFlashInfo).
//
// ExifTool currently treats this payload as unknown, so we retain only a
// bounded raw preview.
type FlashInfo struct {
	Raw BlockPreview
}

// PreviewImageInfo stores Canon MakerNote tag 0x00b6 (PreviewImageInfo).
// Field comments use ExifTool sequence indices.
//
// The first uint32 in this block is a size word and is intentionally omitted
// from this struct (ExifTool comments call this "PreviewImageInfoWords").
type PreviewImageInfo struct {
	PreviewQuality     Quality // [1] PreviewQuality
	PreviewImageLength uint32  // [2] PreviewImageLength
	PreviewImageWidth  uint32  // [3] PreviewImageWidth
	PreviewImageHeight uint32  // [4] PreviewImageHeight
	PreviewImageStart  uint32  // [5] PreviewImageStart
}

// SensorInfo stores Canon MakerNote tag 0x00e0 (SensorInfo/ImageAreaDesc).
// Field comments use ExifTool sequence indices.
type SensorInfo struct {
	SensorWidth           int16 // [1] SensorWidth
	SensorHeight          int16 // [2] SensorHeight
	SensorLeftBorder      int16 // [5] SensorLeftBorder
	SensorTopBorder       int16 // [6] SensorTopBorder
	SensorRightBorder     int16 // [7] SensorRightBorder
	SensorBottomBorder    int16 // [8] SensorBottomBorder
	BlackMaskLeftBorder   int16 // [9] BlackMaskLeftBorder
	BlackMaskTopBorder    int16 // [10] BlackMaskTopBorder
	BlackMaskRightBorder  int16 // [11] BlackMaskRightBorder
	BlackMaskBottomBorder int16 // [12] BlackMaskBottomBorder
}

// AFConfig stores Canon MakerNote tag 0x4028 (AFConfig/AFTabInfo).
// Field comments use ExifTool sequence indices.
type AFConfig struct {
	AFConfigTool              uint32 // [1] AFConfigTool (ExifTool ValueConv: +1)
	AFTrackingSensitivity     int32  // [2] AFTrackingSensitivity
	AFAccelDecelTracking      int32  // [3] AFAccelDecelTracking
	AFPointSwitching          int32  // [4] AFPointSwitching
	AIServoFirstImage         int32  // [5] AIServoFirstImage
	AIServoSecondImage        int32  // [6] AIServoSecondImage
	USMLensElectronicMF       int32  // [7] USMLensElectronicMF
	AFAssistBeam              int32  // [8] AFAssistBeam
	OneShotAFRelease          int32  // [9] OneShotAFRelease
	AutoAFPointSelEOSiTRAF    int32  // [10] AutoAFPointSelEOSiTRAF
	LensDriveWhenAFImpossible int32  // [11] LensDriveWhenAFImpossible
	SelectAFAreaSelectionMode uint32 // [12] SelectAFAreaSelectionMode (bitmask)
	AFAreaSelectionMethod     int32  // [13] AFAreaSelectionMethod
	OrientationLinkedAF       int32  // [14] OrientationLinkedAF
	ManualAFPointSelPattern   int32  // [15] ManualAFPointSelPattern
	AFPointDisplayDuringFocus int32  // [16] AFPointDisplayDuringFocus
	VFDisplayIllumination     int32  // [17] VFDisplayIllumination
	AFStatusViewfinder        int32  // [18] AFStatusViewfinder
	InitialAFPointInServo     int32  // [19] InitialAFPointInServo
	SubjectToDetect           int32  // [20] SubjectToDetect
	EyeDetection              int32  // [24] EyeDetection
}

// RawBurstInfo stores Canon MakerNote tag 0x403f (RawBurstModeRoll/RawBurstInfo).
// Field comments use ExifTool sequence indices.
type RawBurstInfo struct {
	RawBurstImageNum   uint32 // [1] RawBurstImageNum
	RawBurstImageCount uint32 // [2] RawBurstImageCount
}

// FocalLengthInfo stores Canon MakerNote tag 0x0002 (CanonFocalLength).
// Field comments use ExifTool sequence indices.
type FocalLengthInfo struct {
	FocalType       uint16 // [0] FocalType (1=fixed, 2=zoom)
	FocalLength     uint16 // [1] FocalLength (raw units)
	FocalPlaneXSize uint16 // [2] FocalPlaneXSize
	FocalPlaneYSize uint16 // [3] FocalPlaneYSize
}

// AspectInfo stores Canon MakerNote tag 0x009a (AspectInfo).
// Field comments use ExifTool sequence indices.
type AspectInfo struct {
	AspectRatio        uint32 // [0] AspectRatio
	CroppedImageWidth  uint32 // [1] CroppedImageWidth
	CroppedImageHeight uint32 // [2] CroppedImageHeight
	CroppedImageLeft   uint32 // [3] CroppedImageLeft
	CroppedImageTop    uint32 // [4] CroppedImageTop
}

// ProcessingInfo stores Canon MakerNote tag 0x00a0 (ProcessingInfo).
// Field comments use ExifTool sequence indices (FIRST_ENTRY=1).
type ProcessingInfo struct {
	ToneCurve            int16 // [1] ToneCurve
	Sharpness            int16 // [2] Sharpness
	SharpnessFrequency   int16 // [3] SharpnessFrequency
	SensorRedLevel       int16 // [4] SensorRedLevel
	SensorBlueLevel      int16 // [5] SensorBlueLevel
	WhiteBalanceRed      int16 // [6] WhiteBalanceRed
	WhiteBalanceBlue     int16 // [7] WhiteBalanceBlue
	WhiteBalance         int16 // [8] WhiteBalance
	ColorTemperature     int16 // [9] ColorTemperature
	PictureStyle         int16 // [10] PictureStyle
	DigitalGain          int16 // [11] DigitalGain
	WBShiftAB            int16 // [12] WBShiftAB
	WBShiftGM            int16 // [13] WBShiftGM
	UnsharpMaskFineness  int16 // [14] UnsharpMaskFineness
	UnsharpMaskThreshold int16 // [15] UnsharpMaskThreshold
}

// AFMicroAdjInfo stores Canon MakerNote tag 0x4013 (AFMicroAdj).
type AFMicroAdjInfo struct {
	Mode             int32 // [1] AFMicroAdjMode
	ValueNumerator   int32 // [2] AFMicroAdjValue numerator
	ValueDenominator int32 // [2] AFMicroAdjValue denominator
}

// LightingOptInfo stores Canon MakerNote tag 0x4018 (LightingOpt).
// Field comments use ExifTool sequence indices (FIRST_ENTRY=1).
type LightingOptInfo struct {
	PeripheralIlluminationCorr int32 // [1] PeripheralIlluminationCorr
	AutoLightingOptimizer      int32 // [2] AutoLightingOptimizer
	HighlightTonePriority      int32 // [3] HighlightTonePriority
	LongExposureNoiseReduction int32 // [4] LongExposureNoiseReduction
	HighISONoiseReduction      int32 // [5] HighISONoiseReduction
	DigitalLensOptimizer       int32 // [10] DigitalLensOptimizer
	DualPixelRaw               int32 // [11] DualPixelRaw
}

// LensInfoForService stores Canon MakerNote tag 0x4019 (LensInfo).
type LensInfoForService struct {
	LensSerialNumber string // ExifTool-style hex for first 5 bytes
	Raw              [5]byte
	RawCount         uint8
}

// MultiExpInfo stores Canon MakerNote tag 0x4021 (MultiExp).
type MultiExpInfo struct {
	MultiExposure        int32 // [1] MultiExposure
	MultiExposureControl int32 // [2] MultiExposureControl
	MultiExposureShots   int32 // [3] MultiExposureShots
}

// HDRInfo stores Canon MakerNote tag 0x4025 (HDRInfo).
type HDRInfo struct {
	HDR       int32 // [1] HDR
	HDREffect int32 // [2] HDREffect
}

// CanonCustomFunctionMaxEntries caps parsed custom-function entries per table.
const CanonCustomFunctionMaxEntries = 128

// CustomFunctionVariant identifies Canon custom-function table variants.
type CustomFunctionVariant string

// Canon custom-function table variants.
const (
	CustomFunctionVariantFunctions1D        CustomFunctionVariant = "Functions1D"
	CustomFunctionVariantFunctions5D        CustomFunctionVariant = "Functions5D"
	CustomFunctionVariantFunctions10D       CustomFunctionVariant = "Functions10D"
	CustomFunctionVariantFunctions20D       CustomFunctionVariant = "Functions20D"
	CustomFunctionVariantFunctions30D       CustomFunctionVariant = "Functions30D"
	CustomFunctionVariantFunctions350D      CustomFunctionVariant = "Functions350D"
	CustomFunctionVariantFunctions400D      CustomFunctionVariant = "Functions400D"
	CustomFunctionVariantFunctionsD30       CustomFunctionVariant = "FunctionsD30"
	CustomFunctionVariantFunctionsD60       CustomFunctionVariant = "FunctionsD60"
	CustomFunctionVariantFunctionsUnknown   CustomFunctionVariant = "FunctionsUnknown"
	CustomFunctionVariantPersonalFuncs      CustomFunctionVariant = "PersonalFuncs"
	CustomFunctionVariantPersonalFuncValues CustomFunctionVariant = "PersonalFuncValues"
	CustomFunctionVariantFunctions2         CustomFunctionVariant = "Functions2"
	CustomFunctionVariantNotDetermined      CustomFunctionVariant = ""
)

// CustomFunctionEntry stores one parsed Canon custom-function value.
type CustomFunctionEntry struct {
	ID    uint16
	Value int32
	Name  string
}

// CustomFunctionSet stores decoded Canon custom-function entries.
type CustomFunctionSet struct {
	Variant    CustomFunctionVariant
	EntryCount uint8
	Entries    [CanonCustomFunctionMaxEntries]CustomFunctionEntry
}
