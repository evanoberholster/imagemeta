// Package canon provides data types and functions for representing Canon Exif and Makernote values
package canon

import (
	"encoding/binary"
	"math"
)

// CanonMakerNote is Canon Makernote data
// Based on Phil Harvey's exiftool
// Reference: https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/Canon.pm
type CanonMakerNote struct {
	CanonCameraSettings     CameraSettings  `json:"cameraSettings"`  // [1]
	CanonFocalLength        FocalLength     `json:"focalLength"`     // [2]
	CanonFlashInfo          bool            `json:"flashInfo"`       // [3]
	CanonShotInfo           ShotInfo        `json:"shotInfo"`        // [4]
	CanonPanorama           bool            `json:"panorama"`        // [5]
	CanonImageType          string          `json:"imageType"`       // [6]
	FirmwareVersion         string          `json:"firmwareVersion"` // [7]
	FileNumber              uint32          `json:"fileNumber"`      // [8]
	OwnerName               string          `json:"ownerName"`       // [9]
	SerialNumber            string          `json:"serialNumber"`    // [10]
	CameraInfo              CanonCameraInfo `json:"cameraInfo"`      // []
	CanonFileLength         uint32          `json:"fileLength"`      // []
	CustomFunctions         CustomFunctions `json:"customFunctions"` // [15]
	CanonModelID            CanonModelID    `json:"modelID"`         // []
	CanonAFInfo             AFInfo          `json:"afInfo"`          // []
	ThumbnailImageValidArea uint16          `json:"tivArea"`         // []
	SerialNumberFormat      uint32          `json:"serialFormat"`    // []
	SuperMacro              uint16          `json:"superMacro"`      // []
	DateStampMode           uint16          `json:"dateStampMode"`   // []
	MyColors                uint16          `json:"myColors"`        // []
	FirmwareRevision        uint32          `json:"firmwareRev"`     // []
	Categories              uint32          // []
	FaceDetect1             bool            `json:"faceDetect1"`       // []
	FaceDetect2             bool            `json:"faceDetect2"`       // []
	CanonAFInfo2            AFInfo          `json:"afInfo2"`           // []
	ContrastInfo            bool            `json:"contrastInfo"`      // []
	ImageUniqueID           string          `json:"imageUniqueID"`     // []
	WBInfo                  bool            `json:"wbInfo"`            // []
	FaceDetect3             bool            `json:"faceDetect3"`       // []
	Timeinfo                bool            `json:"timeInfo"`          // []
	BatteryType             string          `json:"batteryType"`       // []
	AFInfo3                 AFInfo          `json:"afInfo3"`           // []
	RawDataOffset           uint32          `json:"rawDataOffset"`     // []
	RawDataLength           uint32          `json:"rawDataLength"`     // []
	CanonFileInfo           FileInfo        `json:"fileInfo"`          // []
	AFPointsInFocus1D       AFPoint         `json:"afPointsInFocus1D"` // []
	AspectInfo              bool            `json:"aspectInfo"`        // []
	ColorTemperature        int16           `json:"colorTemperature"`  // []
	CanonFlags              uint32          `json:"flags"`             // []
	PreviewImageInfo        bool            `json:"previewImageInfo"`  // []
	SensorInfo              bool            `json:"sensorInfo"`        // []
	ColorInfo               bool            `json:"colorInfo"`         // []
	AFMicroAdj              bool            `json:"afMicroAdj"`        // []
	LensInfo                bool            `json:"lensInfo"`          // []
	AmbienceInfo            bool            `json:"ambienceInfo"`      // []
	Multiexposure           bool            `json:"multiexposure"`     // []
	FilterInfo              bool            `json:"filterInfo"`        // []
	AFConfig                bool            `json:"afConfig"`          // []
	LevelInfo               bool            `json:"levelInfo"`         // []

}

// CameraSettings is Canon Makernote Camera Settings
// Based on Phil Harvey's exiftool
// Canon camera settings (MakerNotes tag 0x01)
// BinaryData (keys are indices into the int16s array)
// Reference: https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/Canon.pm
type CameraSettings struct {
	ISO                uint32             `json:"iso"`
	SelfTimer          SelfTimer          `json:"selfTimer"`          // [2]
	Quality            CanonQuality       `json:"quality"`            // [3]
	ContinuousDrive    ContinuousDrive    `json:"continuousDrive"`    // [5]
	FocusMode          CanonFocusMode     `json:"focusMode"`          // [7]
	RecordMode         CanonRecordMode    `json:"recordMode"`         // [9]
	ImageSize          CanonImageSize     `json:"imageSize"`          // [10]
	EasyMode           CanonEasyMode      `json:"easyMode"`           // [11]
	DigitalZoom        int16              `json:"digitalZoom"`        // [12]
	Contrast           int16              `json:"contrast"`           // [13]
	Saturation         int16              `json:"saturation"`         // [14]
	Sharpness          int16              `json:"sharpness"`          // [15]
	MeteringMode       MeteringMode       `json:"meteringMode"`       // [17]
	FocusRange         FocusRange         `json:"focusRange"`         // [18]
	AFPoint            AFPointSetting     `json:"afPoint"`            // [19]
	ExposureMode       ExposureMode       `json:"exposureMode"`       // [20]
	LensType           uint16             `json:"lensType"`           // LensType // [22]
	MaxFocalLength     FocalLength        `json:"maxFocalLength"`     // [23]
	MinFocalLength     FocalLength        `json:"minFocalLength"`     // [24]
	FocalUnits         FocalUnits         `json:"focalUnits"`         // [25]
	FlashBits          FlashBits          `json:"flashBits"`          // [29]
	AESetting          AESetting          `json:"aeSetting"`          // [33]
	ImageStabilization ImageStabilization `json:"imageStabilization"` // [34]
	DisplayAperture    DisplayAperture    `json:"displayAperture"`    // [35]. Stored as uint16, divide by 10 to get f-stop
	ZoomSourceWidth    uint16             `json:"zoomSourceWidth"`    // [36]
	ZoomTargetWidth    uint16             `json:"zoomTargetWidth"`    // [37]
	SRAWQuality        SRAWQuality        `json:"srawQuality"`        // [43]
	Clarity            Clarity            `json:"clarity"`            // [44]
	MacroMode          MacroMode          `json:"macroMode"`          // [1]
	CanonFlashMode     CanonFlashMode     `json:"flashMode"`          // [4]
	FocusContinuous    FocusContinuous    `json:"focusContinuous"`    // [32]
	PhotoEffect        PhotoEffect        `json:"photoEffect"`        // [40]
	SpotMeteringMode   SpotMeteringMode   `json:"spotMeteringMode"`   // [39]

	// ManualFlashOutput uint8 // [41] TODO: Not implemented
	// ColorTone         uint8 // [42] TODO: Not implemented
}

// ShotInfo is Canon Makernote Shot Information
// Based on Phil Harvey's exiftool
// Canon shot information (MakerNotes tag 0x04)
// BinaryData (keys are indices into the int16s array)
// Reference: https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/Canon.pm
type ShotInfo struct {
	AutoISO                int16         `json:"autoISO"`                // 1
	BaseISO                int16         `json:"baseISO"`                // 2
	MeasuredEV             int16         `json:"measuredEV"`             // 3
	TargetAperture         int16         `json:"targetAperture"`         // 4
	TargetExposureTime     int16         `json:"targetExposureTime"`     // 5
	ExposureCompensation   int16         `json:"exposureCompensation"`   // 6
	WhiteBalance           WhiteBalance  `json:"whiteBalance"`           // 7
	SlowShutter            int16         `json:"slowShutter"`            // 8
	SequenceNumber         int16         `json:"sequenceNumber"`         // 9
	OpticalZoomCode        int16         `json:"opticalZoomCode"`        // 10
	CameraTemperature      int16         `json:"cameraTemperature"`      // 12
	FlashGuideNumber       int16         `json:"flashGuideNumber"`       // 13
	AFPointsInFocus        uint16        `json:"afPointsInFocus"`        // 14
	FlashExposureComp      int16         `json:"flashExposureComp"`      // 15
	AutoExposureBracketing int16         `json:"autoExposureBracketing"` // 16
	AEBBracketValue        int16         `json:"aebBracketValue"`        // 17
	ControlMode            int16         `json:"controlMode"`            // 18
	FocusDistance          FocusDistance `json:"focusDistance"`          // 19, 20
	FNumber                int16         `json:"fNumber"`                // 21
	ExposureTime           int16         `json:"exposureTime"`           // 22
	MeasuredEV2            int16         `json:"measuredEV2"`            // 23
	BulbDuration           int16         `json:"bulbDuration"`           // 24
	CameraType             int16         `json:"cameraType"`             // 26
	AutoRotate             int16         `json:"autoRotate"`             // 27
	NDFilter               int16         `json:"ndFilter"`               // 28
	SelfTimer2             int16         `json:"selfTimer2"`             // 29
	FlashOutput            int16         `json:"flashOutput"`            // 33
}

// UnmarshalBinary unmarshals the binary data into a ShotInfo struct.
// The data is an array of 16-bit signed integers.
func (si *ShotInfo) UnmarshalBinary(b []byte) error {
	const recordSize = 2 // Each record is an int16 (2 bytes)

	// Helper to safely read int16 from the byte slice
	readInt16 := func(index int) int16 {
		offset := index * recordSize
		if offset+recordSize > len(b) {
			return 0 // or handle error appropriately
		}
		return int16(binary.LittleEndian.Uint16(b[offset : offset+recordSize]))
	}

	// Helper to safely read uint16 from the byte slice
	readUint16 := func(index int) uint16 {
		offset := index * recordSize
		if offset+recordSize > len(b) {
			return 0 // or handle error appropriately
		}
		return binary.LittleEndian.Uint16(b[offset : offset+recordSize])
	}

	si.AutoISO = int16(math.Exp(float64(readInt16(1))/32.0*math.Log(2)) * 100)                       // 1. Notes: actual ISO used = BaseISO * AutoISO / 100
	si.BaseISO = int16(math.Exp(float64(readInt16(2))/32.0*math.Log(2)) * 100 / 32)                  // 2.
	si.MeasuredEV = readInt16(3)/32 + 5                                                              // 3. Notes: this is the Canon name for what could better be called MeasuredLV, and should be close to the calculated LightValue for a proper exposure with most models
	si.TargetAperture = int16(math.Exp(float64(CanonEv(float32(readInt16(4)))) * math.Log(2) / 2.0)) // 4.
	si.TargetExposureTime = int16(math.Exp(float64(-CanonEv(float32(readInt16(5)))) * math.Log(2)))  // 5.
	si.ExposureCompensation = int16(CanonEv(float32(readInt16(6))))                                  // 6.
	si.WhiteBalance = WhiteBalance(readInt16(7))                                                     // 7.
	si.SlowShutter = readInt16(8)                                                                    // 8.
	si.SequenceNumber = readInt16(9)                                                                 // 9. Description: Shot Number In Continuous Burst. Notes: valid only for some models (eg. not the 5DmkIII)
	si.OpticalZoomCode = readInt16(10)                                                               // 10. Notes: for many PowerShot models, a this is 0-6 for wide-tele zoom
	si.CameraTemperature = readInt16(12)                                                             // 12. Notes: newer EOS models only
	si.FlashGuideNumber = readInt16(13) / 32                                                         // 13.
	si.AFPointsInFocus = readUint16(14)                                                              // 14. Notes: used by D30, D60 and some PowerShot/Ixus models
	si.FlashExposureComp = int16(CanonEv(float32(readInt16(15))))                                    // 15. Description: Flash Exposure Compensation
	si.AutoExposureBracketing = readInt16(16)                                                        // 16.
	si.AEBBracketValue = int16(CanonEv(float32(readInt16(17))))                                      // 17.
	si.ControlMode = readInt16(18)                                                                   // 18.
	si.FocusDistance = FocusDistance{readInt16(19), readInt16(20)}                                   // 19, 20. Notes: FocusDistance tags are only extracted if FocusDistanceUpper is non-zero
	si.FNumber = int16(math.Exp(float64(CanonEv(float32(readInt16(21)))) * math.Log(2) / 2.0))       // 21.
	si.ExposureTime = int16(math.Exp(float64(-CanonEv(float32(readInt16(22)))) * math.Log(2)))       // 22.
	si.MeasuredEV2 = readInt16(23)/8 - 6                                                             // 23. Description: Measured EV 2
	si.BulbDuration = readInt16(24) / 10                                                             // 24.
	si.CameraType = readInt16(26)                                                                    // 26.
	si.AutoRotate = readInt16(27)                                                                    // 27.
	si.NDFilter = readInt16(28)                                                                      // 28.
	si.SelfTimer2 = readInt16(29)                                                                    // 29.
	si.FlashOutput = readInt16(33)                                                                   // 33. Notes: used only for PowerShot models, this has a maximum value of 500 for models like the A570IS

	return nil
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
	AFPoints      []AFPoint
	InFocus       []int
	Selected      []int
	AFAreaMode    AFAreaMode
	NumAFPoints   uint16
	ValidAFPoints uint16
}

// AFPoint is an Auto Focus Point
type AFPoint [4]int16

// NewAFPoint returns a new AFPoint from
// width, height, x-axis coord and y-axis coord
func NewAFPoint(w, h, x, y int16) AFPoint {
	return AFPoint{w, h, x, y}
}
