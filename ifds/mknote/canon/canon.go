package canon

import "fmt"

// CameraSettings - Canon Makernote Camera Settings
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

// ShotInfo - Canon Makernote Shot Information
// TODO: Incomplete
type ShotInfo struct {
	CameraTemperature      int16         // [12] 	CameraTemperature 	int16s 	(newer EOS models only)
	FlashExposureComp      int16         // [15] 	FlashExposureComp 	int16s
	AutoExposureBracketing int16         // [16] 	AutoExposureBracketing 	int16s
	AEBBracketValue        int16         // [17] 	AEBBracketValue 	int16s
	SelfTimer              int16         // 29 	SelfTimer2 	int16s
	FocusDistance          FocusDistance // 19 	FocusDistanceUpper 	int16u // 20 	FocusDistanceLower 	int16u
}

// FileInfo - Canon Makernote File Information
type FileInfo struct {
	FocusDistance     FocusDistance // 20 	FocusDistanceUpper 	int16u // 21 	FocusDistanceLower 	int16u
	BracketMode       BracketMode   // 3 	BracketMode 	int16s
	BracketValue      int16         // 4 	BracketValue 	int16s
	BracketShotNumber int16         // 5 	BracketShotNumber 	int16s
	LiveViewShooting  bool          // 19 	LiveViewShooting 	int16s (bool)
}

// AFInfo - Canon Makernote Autofocus Information
type AFInfo struct {
	AFAreaMode    AFAreaMode
	NumAFPoints   uint16
	ValidAFPoints uint16
	AFPoints      []AFPoint
	InFocus       []int
	Selected      []int
}

// AFPoint - AutoFocusPoint
type AFPoint [4]int16

// NewAFPoint - creates a new AFPoint from
// width, height, x-axis coord and y-axis coord
func NewAFPoint(w, h, x, y int16) AFPoint {
	return AFPoint{w, h, x, y}
}

// Ev - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool/Canon.pm
func Ev(val int16) int16 {
	var sign int16
	if val < 0 {
		val = -val
		sign = -1
	} else {
		sign = 1
	}
	frac := val & 0x1f
	val -= frac
	// Convert 1/3 and 2/3 codes
	if frac == 0x0c {
		frac = 0x20 / 3
	} else if frac == 0x14 {
		frac = 0x40 / 3
	}
	return sign * (val + frac) / 0x20
}

// TempConv - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool/Canon.pm
func TempConv(val uint16) int16 {
	if val == 0 {
		return 0
	}
	return int16(val) - 128
}

// PointsInFocus
func PointsInFocus(af []uint16) (inFocus []int, selected []int, err error) {
	validPoints := int(af[3])
	var count int
	// NumAFPoints may be 7, 9, 11, 19, 31, 45 or 61, depending on the camera model.
	switch validPoints {
	case 7:
		count = 1 // 1
	case 9, 11:
		count = 1 // 1
	case 19, 31:
		count = 2 // 2
	case 45:
		count = 3 // 3
	case 61:
		count = 4 // 4
	case 65:
		count = 5 // 5
	default:
		panic(fmt.Errorf("Error parsing AFPoints from Canon Makernote. Expected 7, 9, 11, 19, 31, 45 or 61 got %d", validPoints))
	}
	off := 8 + (validPoints * 4)
	inFocus = decodeBits(af[off:off+count], 16)
	selected = decodeBits(af[off+count:off+count+count], 16)
	return
}

// decodeBits - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool.pm
func decodeBits(vals []uint16, bits int) (list []int) {
	var num int
	var n int
	for _, a := range vals {
		for i := 0; i < bits; i++ {
			n = i + num
			if a&(1<<uint(i)) > 0 {
				list = append(list, n)
			}
		}
		num += bits
	}
	return
}

func ParseAFPoints(af []uint16) (afPoints []AFPoint) {
	validPoints := int(af[3])
	// AFPoints
	afPoints = make([]AFPoint, validPoints)
	xAdjust := int16(af[4] / 2) // Adjust x-axis
	yAdjust := int16(af[5] / 2) // Adjust y-axis

	for i := 0; i < validPoints; i++ { // Start at an offset of 8
		offset := 8 + i
		w := int16(af[offset])
		h := int16(af[offset+validPoints])
		x := int16(af[offset+(2*validPoints)]) + xAdjust - (w / 2)
		y := int16(af[offset+(3*validPoints)]) + yAdjust - (h / 2)
		afPoints[i] = NewAFPoint(w, h, x, y)
	}
	return
}
