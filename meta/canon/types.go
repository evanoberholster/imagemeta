package canon

// CameraSettings is Canon Makernote Camera Settings
// TODO: Incomplete
type CameraSettings struct {
	Macromode         bool            // [1]
	SelfTimer         bool            // [2]
	ContinuousDrive   ContinuousDrive // [5]
	FocusMode         FocusMode       // [7]
	MeteringMode      MeteringMode    // [17]
	FocusRange        FocusRange      // [18]
	CanonExposureMode ExposureMode    // [20]
	//MaxFocalLength    int16           // [23]
	//MinFocalLength    int16           // [24]
	//FocalUnits        int16                       // [25]
	//FocusContinuous   CanonFocusContinous  // [32]
	//SpotMeteringMode  bool                        // [39]
	AESetting AESetting // [33]
}

// ShotInfo is Canon Makernote Shot Information
// TODO: Incomplete
type ShotInfo struct {
	CameraTemperature      int16         // [12] 	CameraTemperature 	int16s 	(newer EOS models only)
	FlashExposureComp      int16         // [15] 	FlashExposureComp 	int16s
	AutoExposureBracketing int16         // [16] 	AutoExposureBracketing 	int16s
	AEBBracketValue        int16         // [17] 	AEBBracketValue 	int16s
	SelfTimer              int16         // 29 	SelfTimer2 	int16s
	FocusDistance          FocusDistance // 19 	FocusDistanceUpper 	int16u // 20 	FocusDistanceLower 	int16u
}

// FileInfo is Canon Makernote File Information
type FileInfo struct {
	FocusDistance     FocusDistance // 20 	FocusDistanceUpper 	int16u // 21 	FocusDistanceLower 	int16u
	BracketMode       BracketMode   // 3 	BracketMode 	int16s
	BracketValue      int16         // 4 	BracketValue 	int16s
	BracketShotNumber int16         // 5 	BracketShotNumber 	int16s
	LiveViewShooting  bool          // 19 	LiveViewShooting 	int16s (bool)
}

// AFInfo is Canon Makernote Autofocus Information
type AFInfo struct {
	AFAreaMode    AFAreaMode
	NumAFPoints   uint16
	ValidAFPoints uint16
	AFPoints      []AFPoint
	InFocus       []int
	Selected      []int
}

// AFPoint is an Auto Focus Point
type AFPoint [4]int16

// NewAFPoint returns a new AFPoint from
// width, height, x-axis coord and y-axis coord
func NewAFPoint(w, h, x, y int16) AFPoint {
	return AFPoint{w, h, x, y}
}
