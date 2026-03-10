package canon

// CameraSettings stores selected Canon CameraSettings values from MakerNote tag
// 0x0001 (CanonCameraSettings).
//
// Field comments use ExifTool sequence indices for this table.
type CameraSettings struct {
	MacroMode          int16           // [1] MacroMode
	SelfTimer          int16           // [2] SelfTimer (deci-seconds; custom bit may be set)
	Quality            int16           // [3] Quality
	CanonFlashMode     int16           // [4] CanonFlashMode
	ContinuousDrive    ContinuousDrive // [5] ContinuousDrive
	FocusMode          FocusMode       // [7] FocusMode
	RecordMode         int16           // [9] RecordMode
	CanonImageSize     int16           // [10] CanonImageSize
	EasyMode           int16           // [11] EasyMode
	DigitalZoom        int16           // [12] DigitalZoom
	Contrast           int16           // [13] Contrast
	Saturation         int16           // [14] Saturation
	Sharpness          int16           // [15] Sharpness
	CameraISO          int16           // [16] CameraISO
	MeteringMode       MeteringMode    // [17] MeteringMode
	FocusRange         FocusRange      // [18] FocusRange
	AFPoint            uint16          // [19] AFPoint
	CanonExposureMode  ExposureMode    // [20] CanonExposureMode
	LensType           uint16          // [22] LensType
	MaxFocalLength     uint16          // [23] MaxFocalLength
	MinFocalLength     uint16          // [24] MinFocalLength
	FocalUnits         uint16          // [25] FocalUnits
	MaxAperture        int16           // [26] MaxAperture
	MinAperture        int16           // [27] MinAperture
	FlashActivity      int16           // [28] FlashActivity
	FlashBits          uint16          // [29] FlashBits
	FocusContinuous    int16           // [32] FocusContinuous
	AESetting          AESetting       // [33] AESetting
	ImageStabilization int16           // [34] ImageStabilization
	DisplayAperture    uint16          // [35] DisplayAperture (x10)
	ZoomSourceWidth    uint16          // [36] ZoomSourceWidth
	ZoomTargetWidth    uint16          // [37] ZoomTargetWidth
	SpotMeteringMode   int16           // [39] SpotMeteringMode
	PhotoEffect        int16           // [40] PhotoEffect
	ManualFlashOutput  int16           // [41] ManualFlashOutput
	ColorTone          int16           // [42] ColorTone
	SRAWQuality        int16           // [46] SRAWQuality
	Clarity            int16           // [51] Clarity (EOS R models)
}

// ShotInfo stores selected Canon ShotInfo values from MakerNote tag 0x0004
// (CanonShotInfo).
//
// Field comments use ExifTool sequence indices for this table.
type ShotInfo struct {
	AutoISO                int16         // [1] AutoISO
	BaseISO                int16         // [2] BaseISO
	MeasuredEV             int16         // [3] MeasuredEV
	TargetAperture         int16         // [4] TargetAperture
	TargetExposureTime     int16         // [5] TargetExposureTime
	ExposureCompensation   int16         // [6] ExposureCompensation
	WhiteBalance           int16         // [7] WhiteBalance
	SlowShutter            int16         // [8] SlowShutter
	SequenceNumber         int16         // [9] SequenceNumber
	OpticalZoomCode        int16         // [10] OpticalZoomCode
	CameraTemperature      int16         // [12] CameraTemperature
	FlashGuideNumber       int16         // [13] FlashGuideNumber
	AFPointsInFocus        uint16        // [14] AFPointsInFocus bitset
	FlashExposureComp      int16         // [15] FlashExposureCompensation
	AutoExposureBracketing int16         // [16] AutoExposureBracketing
	AEBBracketValue        int16         // [17] AEBBracketValue
	ControlMode            int16         // [18] ControlMode
	FocusDistance          FocusDistance // [19-20] FocusDistanceUpper/Lower
	FNumber                int16         // [21] FNumber
	ExposureTime           int16         // [22] ExposureTime
	MeasuredEV2            int16         // [23] MeasuredEV2
	BulbDuration           int16         // [24] BulbDuration (deci-seconds)
	CameraType             int16         // [26] CameraType
	AutoRotate             int16         // [27] AutoRotate
	NDFilter               int16         // [28] NDFilter
	SelfTimer2             int16         // [29] SelfTimer2 (deci-seconds)
	FlashOutput            int16         // [33] FlashOutput
}

// FileInfo stores selected Canon FileInfo values from MakerNote tag 0x0093
// (CanonFileInfo).
//
// Field comments use ExifTool sequence indices for this table.
type FileInfo struct {
	FileNumber                  uint32        // [1] FileNumber
	ShutterCount                uint32        // [1] ShutterCount (model-dependent variant)
	BracketMode                 BracketMode   // [3] BracketMode
	BracketValue                int16         // [4] BracketValue
	BracketShotNumber           int16         // [5] BracketShotNumber
	RawJpgQuality               int16         // [6] RawJpgQuality
	RawJpgSize                  int16         // [7] RawJpgSize
	LongExposureNoiseReduction2 int16         // [8] LongExposureNoiseReduction2
	WBBracketMode               int16         // [9] WBBracketMode
	WBBracketValueAB            int16         // [12] WBBracketValueAB
	WBBracketValueGM            int16         // [13] WBBracketValueGM
	FilterEffect                int16         // [14] FilterEffect
	ToningEffect                int16         // [15] ToningEffect
	MacroMagnification          int16         // [16] MacroMagnification
	LiveViewShooting            bool          // [19] LiveViewShooting
	FocusDistance               FocusDistance // [20-21] FocusDistanceUpper/Lower
	ShutterMode                 int16         // [23] ShutterMode
	FlashExposureLock           bool          // [25] FlashExposureLock
	AntiFlicker                 bool          // [32] AntiFlicker
	RFLensType                  uint16        // [0x3d] RFLensType
}

// AFInfo stores selected Canon autofocus record values from MakerNote tags
// 0x0012 (AFInfo) and 0x0026 (AFInfo2).
type AFInfo struct {
	AFAreaMode       AFAreaMode // AFInfo2 [1]
	NumAFPoints      uint16     // AFInfo [0] / AFInfo2 [2]
	ValidAFPoints    uint16     // AFInfo [1] / AFInfo2 [3]
	CanonImageWidth  uint16     // AFInfo [2] / AFInfo2 [4]
	CanonImageHeight uint16     // AFInfo [3] / AFInfo2 [5]
	AFImageWidth     uint16     // AFInfo [4] / AFInfo2 [6]
	AFImageHeight    uint16     // AFInfo [5] / AFInfo2 [7]
	AFAreaWidth      uint16     // AFInfo [6]
	AFAreaHeight     uint16     // AFInfo [7]

	AFAreaWidths     []int16 // AFInfo2 [8]
	AFAreaHeights    []int16 // AFInfo2 [9]
	AFAreaXPositions []int16 // AFInfo [8] / AFInfo2 [10]
	AFAreaYPositions []int16 // AFInfo [9] / AFInfo2 [11]

	AFPointsInFocusBits  []int  // AFInfo [10] / AFInfo2 [12]
	AFPointsSelectedBits []int  // AFInfo2 [13]
	PrimaryAFPoint       uint16 // AFInfo [11/12] / AFInfo2 [14]

	AFPoints []AFPoint

	// Backward-compatible aliases used by older parser code paths.
	InFocus  []int
	Selected []int
}

// AFPoint is an Auto Focus Point
type AFPoint [4]int16

// NewAFPoint returns a new AFPoint from
// width, height, x-axis coord and y-axis coord
func NewAFPoint(w, h, x, y int16) AFPoint {
	return AFPoint{w, h, x, y}
}
