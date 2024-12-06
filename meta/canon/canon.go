// Package canon provides data types and functions for representing Canon Exif and Makernote values
package canon

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
