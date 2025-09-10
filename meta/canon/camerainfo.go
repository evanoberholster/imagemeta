package canon

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Information in the "CameraInfo" records is tricky to decode because the
// encodings are very different than in other Canon records (even sometimes
// switching endianness between values within a single camera), plus there is
// considerable variation in format from model to model. The first table below
// lists CameraInfo tags for the 1D and 1DS.

//go:generate msgp

// CanonCameraInfo defines the interface for unmarshaling Canon CameraInfo maker note data.
// This data is highly model-specific, with different layouts and encodings for various camera models.
type CanonCameraInfo interface {
	UnmarshalBinary(data []byte) error
}

// CameraInfo1D represents camera information specific to the EOS-1D model
type CameraInfo1D struct {
	cameraModelID      CanonModelID   // Internal model ID, used for parsing variants.
	ExposureTime       CIExposureTime `json:"exposureTime"`
	FocalLength        FocalLength    `json:"focalLength"`
	LensType           CanonLensType  `json:"lensType"`
	MinFocalLength     CIFocalLength  `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength  `json:"maxFocalLength"`
	SharpnessFrequency uint8          `json:"sharpnessFrequency"`
	Sharpness          int8           `json:"sharpness"`
	WhiteBalance       WhiteBalance   `json:"whiteBalance"`
	ColorTemperature   uint16         `json:"colorTemperature"`
	PictureStyle       PictureStyle   `json:"pictureStyle"`
}

// CameraInfo1DmkII represents camera information specific to the EOS-1D Mark II model
type CameraInfo1DmkII struct {
	cameraModelID    CanonModelID   // Internal model ID, used for parsing variants.
	ExposureTime     CIExposureTime `json:"exposureTime"`
	FocalLength      FocalLength    `json:"focalLength"`
	LensType         CanonLensType  `json:"lensType"`
	MinFocalLength   CIFocalLength  `json:"minFocalLength"`
	MaxFocalLength   CIFocalLength  `json:"maxFocalLength"`
	WhiteBalance     WhiteBalance   `json:"whiteBalance"`
	ColorTemperature uint16         `json:"colorTemperature"`
	CanonImageSize   uint16         `json:"canonImageSize"`
	JPEGQuality      uint16         `json:"jpegQuality"`
	PictureStyle     PictureStyle   `json:"pictureStyle"`
	Saturation       int8           `json:"saturation"`
	ColorTone        int8           `json:"colorTone"`
	Sharpness        int8           `json:"sharpness"`
	Contrast         int8           `json:"contrast"`
	ISO              string         `json:"iso"`
}

// CameraInfo1DmkIIN represents camera information specific to the EOS-1D Mark II N model
type CameraInfo1DmkIIN struct {
	cameraModelID    CanonModelID   // Internal model ID, used for parsing variants.
	ExposureTime     CIExposureTime `json:"exposureTime"`
	FocalLength      FocalLength    `json:"focalLength"`
	LensType         CanonLensType  `json:"lensType"`
	MinFocalLength   CIFocalLength  `json:"minFocalLength"`
	MaxFocalLength   CIFocalLength  `json:"maxFocalLength"`
	WhiteBalance     WhiteBalance   `json:"whiteBalance"`
	ColorTemperature uint16         `json:"colorTemperature"`
	PictureStyle     PictureStyle   `json:"pictureStyle"`
	Sharpness        int8           `json:"sharpness"`
	Contrast         int8           `json:"contrast"`
	Saturation       int8           `json:"saturation"`
	ColorTone        int8           `json:"colorTone"`
	ISO              string         `json:"iso"`
}

// CameraInfo1DmkIII represents camera information specific to the EOS-1D Mark III and EOS-1Ds Mark III models
type CameraInfo1DmkIII struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	MacroMagnification uint8             `json:"macroMagnification"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	ShutterCount       uint32            `json:"shutterCount"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	TimeStamp          uint32            `json:"timeStamp"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo1DmkIV represents camera information for the EOS 1D Mark IV.
// Indices shown are for firmware version 1.0.2, but they may be different for other firmware versions.
type CameraInfo1DmkIV struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	MeasuredEV2           int16             `json:"measuredEV2"`
	MeasuredEV3           int16             `json:"measuredEV3"` // In some firmware versions
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	CameraPictureStyle    PictureStyle      `json:"cameraPictureStyle"`
	HighISONoiseReduction uint8             `json:"highISONoiseReduction"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo1DX represents camera information specific to the EOS-1D X model
type CameraInfo1DX struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo5D represents camera information specific to the EOS 5D model
type CameraInfo5D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	LensType           CanonLensType     `json:"lensType"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	MacroMagnification uint8             `json:"macroMagnification"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocalLength        FocalLength       `json:"focalLength"`
	AFPointsInFocus5D  AFPointsInFocus5D `json:"afPointsInFocus5D"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareRevision   string            `json:"firmwareRevision"`
	ShortOwnerName     string            `json:"shortOwnerName"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	FileIndex          uint16            `json:"fileIndex"`
	TimeStamp          uint32            `json:"timeStamp"`
	PictureStyleInfo   PictureStyleInfo5D
}

// PictureStyleInfo5D represents custom picture style information for the EOS 5D.
type PictureStyleInfo5D struct {
	ContrastStandard       int8         `json:"contrastStandard"`
	ContrastPortrait       int8         `json:"contrastPortrait"`
	ContrastLandscape      int8         `json:"contrastLandscape"`
	ContrastNeutral        int8         `json:"contrastNeutral"`
	ContrastFaithful       int8         `json:"contrastFaithful"`
	ContrastMonochrome     int8         `json:"contrastMonochrome"`
	ContrastUserDef1       int8         `json:"contrastUserDef1"`
	ContrastUserDef2       int8         `json:"contrastUserDef2"`
	ContrastUserDef3       int8         `json:"contrastUserDef3"`
	SharpnessStandard      uint8        `json:"sharpnessStandard"`
	SharpnessPortrait      uint8        `json:"sharpnessPortrait"`
	SharpnessLandscape     uint8        `json:"sharpnessLandscape"`
	SharpnessNeutral       uint8        `json:"sharpnessNeutral"`
	SharpnessFaithful      uint8        `json:"sharpnessFaithful"`
	SharpnessMonochrome    uint8        `json:"sharpnessMonochrome"`
	SharpnessUserDef1      uint8        `json:"sharpnessUserDef1"`
	SharpnessUserDef2      uint8        `json:"sharpnessUserDef2"`
	SharpnessUserDef3      uint8        `json:"sharpnessUserDef3"`
	SaturationStandard     int8         `json:"saturationStandard"`
	SaturationPortrait     int8         `json:"saturationPortrait"`
	SaturationLandscape    int8         `json:"saturationLandscape"`
	SaturationNeutral      int8         `json:"saturationNeutral"`
	SaturationFaithful     int8         `json:"saturationFaithful"`
	FilterEffectMonochrome FilterEffect `json:"filterEffectMonochrome"`
}

// CameraInfo5DmkII represents camera information specific to the EOS 5D Mark II model
type CameraInfo5DmkII struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	MacroMagnification    uint8             `json:"macroMagnification"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	PictureStyle          PictureStyle      `json:"pictureStyle"`
	HighISONoiseReduction uint8             `json:"highISONoiseReduction"`
	AutoLightingOptimizer uint8             `json:"autoLightingOptimizer"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo5DmkIII represents camera information specific to the EOS 5D Mark III model
type CameraInfo5DmkIII struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo6D represents camera information specific to the EOS 6D model
type CameraInfo6D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo70D represents camera information for the EOS 70D.
type CameraInfo70D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo80D represents camera information for the EOS 80D.
type CameraInfo80D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
}

// CameraInfo7D represents camera information for the EOS 7D.
// Indices shown are for firmware versions 1.0.x.
type CameraInfo7D struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	MeasuredEV2           int16             `json:"measuredEV2"`
	MeasuredEV            int16             `json:"measuredEV"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	CameraPictureStyle    PictureStyle      `json:"cameraPictureStyle"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo40D represents camera information for the EOS 40D.
type CameraInfo40D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	FlashMeteringMode  FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	MacroMagnification uint8             `json:"macroMagnification"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
	LensModel          string            `json:"lensModel"`
}

// CameraInfo50D represents camera information for the EOS 50D.
// Indices shown are for firmware versions 1.0.x.
type CameraInfo50D struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	PictureStyle          PictureStyle      `json:"pictureStyle"`
	HighISONoiseReduction uint8             `json:"highISONoiseReduction"`
	AutoLightingOptimizer uint8             `json:"autoLightingOptimizer"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo60D represents camera information for the EOS 60D and 1200D.
type CameraInfo60D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo450D represents camera information for the EOS 450D.
type CameraInfo450D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	MacroMagnification uint8             `json:"macroMagnification"`
	FlashMeteringMode  FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	OwnerName          string            `json:"ownerName"`
	FileIndex          uint32            `json:"fileIndex"`
	LensModel          string            `json:"lensModel"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo500D represents camera information for the EOS 500D.
type CameraInfo500D struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	PictureStyle          PictureStyle      `json:"pictureStyle"`
	HighISONoiseReduction uint8             `json:"highISONoiseReduction"`
	AutoLightingOptimizer uint8             `json:"autoLightingOptimizer"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo550D represents camera information for the EOS 550D.
type CameraInfo550D struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo600D represents camera information for the EOS 600D and 1100D.
type CameraInfo600D struct {
	cameraModelID         CanonModelID      // Internal model ID, used for parsing variants.
	FNumber               FNumber           `json:"fNumber"`
	ExposureTime          CIExposureTime    `json:"exposureTime"`
	ISO                   ISO               `json:"iso"`
	HighlightTonePriority uint8             `json:"highlightTonePriority"`
	FlashMeteringMode     FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature     CameraTemperature `json:"cameraTemperature"`
	FocalLength           FocalLength       `json:"focalLength"`
	CameraOrientation     CameraOrientation `json:"orientation"`
	FocusDistanceUpper    FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower    FocalLength       `json:"focusDistanceLower"`
	WhiteBalance          WhiteBalance      `json:"whiteBalance"`
	ColorTemperature      uint16            `json:"colorTemperature"`
	PictureStyle          PictureStyle      `json:"pictureStyle"`
	HighISONoiseReduction uint8             `json:"highISONoiseReduction"`
	AutoLightingOptimizer uint8             `json:"autoLightingOptimizer"`
	LensType              CanonLensType     `json:"lensType"`
	MinFocalLength        CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength        CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion       string            `json:"firmwareVersion"`
	FileIndex             uint32            `json:"fileIndex"`
	DirectoryIndex        uint32            `json:"directoryIndex"`
	PictureStyleInfo      PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo650D represents camera information for the EOS 650D.
type CameraInfo650D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	PictureStyle       PictureStyle      `json:"pictureStyle"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfo700D represents camera information for the EOS 700D.
type CameraInfo700D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`            // 0x03
	ExposureTime       CIExposureTime    `json:"exposureTime"`       // 0x04
	ISO                ISO               `json:"iso"`                // 0x06
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`  // 0x1b
	FocalLength        FocalLength       `json:"focalLength"`        // 0x23
	CameraOrientation  CameraOrientation `json:"orientation"`        // 0x7d
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"` // 0x8c
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"` // 0x8e
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`       // 0xbc
	ColorTemperature   uint16            `json:"colorTemperature"`   // 0xc0
	PictureStyle       PictureStyle      `json:"pictureStyle"`       // 0xf4
	LensType           CanonLensType     `json:"lensType"`           // 0x127
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`     // 0x129
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`     // 0x12b
	FirmwareVersion    string            `json:"firmwareVersion"`    // 0x220
	FileIndex          uint32            `json:"fileIndex"`          // 0x274
	DirectoryIndex     uint32            `json:"directoryIndex"`     // 0x280
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`   // 0x390
}

// AFPointsInFocus5D represents the autofocus points for the EOS 5D.
type AFPointsInFocus5D uint16

const (
	AFPoint5DCenter     AFPointsInFocus5D = 1 << 0
	AFPoint5DTop        AFPointsInFocus5D = 1 << 1
	AFPoint5DBottom     AFPointsInFocus5D = 1 << 2
	AFPoint5DUpperLeft  AFPointsInFocus5D = 1 << 3
	AFPoint5DUpperRight AFPointsInFocus5D = 1 << 4
	AFPoint5DLowerLeft  AFPointsInFocus5D = 1 << 5
	AFPoint5DLowerRight AFPointsInFocus5D = 1 << 6
	AFPoint5DLeft       AFPointsInFocus5D = 1 << 7
	AFPoint5DRight      AFPointsInFocus5D = 1 << 8
	AFPoint5DAIServo1   AFPointsInFocus5D = 1 << 9
	AFPoint5DAIServo2   AFPointsInFocus5D = 1 << 10
)

// CameraInfo750D represents camera information for the EOS 750D and 760D.
type CameraInfo750D struct{}
type CameraInfo760D struct{}

// CameraInfo1000D represents camera information for the EOS 1000D.
type CameraInfo1000D struct {
	cameraModelID      CanonModelID      // Internal model ID, used for parsing variants.
	FNumber            FNumber           `json:"fNumber"`
	ExposureTime       CIExposureTime    `json:"exposureTime"`
	ISO                ISO               `json:"iso"`
	MacroMagnification uint8             `json:"macroMagnification"`
	FlashMeteringMode  FlashMeteringMode `json:"flashMeteringMode"`
	CameraTemperature  CameraTemperature `json:"cameraTemperature"`
	FocalLength        FocalLength       `json:"focalLength"`
	CameraOrientation  CameraOrientation `json:"orientation"`
	FocusDistanceUpper FocalLength       `json:"focusDistanceUpper"`
	FocusDistanceLower FocalLength       `json:"focusDistanceLower"`
	WhiteBalance       WhiteBalance      `json:"whiteBalance"`
	ColorTemperature   uint16            `json:"colorTemperature"`
	LensType           CanonLensType     `json:"lensType"`
	MinFocalLength     CIFocalLength     `json:"minFocalLength"`
	MaxFocalLength     CIFocalLength     `json:"maxFocalLength"`
	FirmwareVersion    string            `json:"firmwareVersion"`
	FileIndex          uint32            `json:"fileIndex"`
	LensModel          string            `json:"lensModel"`
	DirectoryIndex     uint32            `json:"directoryIndex"`
	PictureStyleInfo   PictureStyleInfo2 `json:"pictureStyleInfo"`
}

// CameraInfoPowerShot represents camera information for various PowerShot models.
type CameraInfoPowerShot struct {
	cameraModelID     CanonModelID   // Internal model ID, used for parsing variants.
	ISO               ISO            `json:"iso"`
	FNumber           FNumber        `json:"fNumber"`
	ExposureTime      CIExposureTime `json:"exposureTime"`
	Rotation          uint32         `json:"rotation"`
	CameraTemperature int32          `json:"cameraTemperature"` // This is int32, not CameraTemperature (uint8)
}

// CameraInfo tags for the EOS R5 and R6.
type CameraInfoR6 struct {
	cameraModelID CanonModelID // Internal model ID, used for parsing variants.
	ShutterCount  uint32       `json:"shutterCount"`
}

// CameraInfoR6m2 represents camera information for the EOS R6 Mark II, R8 and R50.
type CameraInfoR6m2 struct {
	cameraModelID CanonModelID // Internal model ID, used for parsing variants.
	ShutterCount  uint32       `json:"shutterCount"`
}

// CameraInfoG5XII represents camera information for the PowerShot G5 X Mark II.
type CameraInfoG5XII struct {
	cameraModelID  CanonModelID // Internal model ID, used for parsing variants.
	ShutterCount   uint32       `json:"shutterCount"`
	DirectoryIndex uint32       `json:"directoryIndex"`
	FileIndex      uint32       `json:"fileIndex"`
}

// CameraInfoPowerShot2 represents camera information for various PowerShot models.
type CameraInfoPowerShot2 struct {
	cameraModelID     CanonModelID      // Internal model ID, used for parsing variants.
	ISO               ISO               `json:"iso"`
	FNumber           FNumber           `json:"fNumber"`
	ExposureTime      CIExposureTime    `json:"exposureTime"`
	Rotation          uint32            `json:"rotation"`
	CameraTemperature CameraTemperature `json:"cameraTemperature"`
}

// CameraInfoUnknown32 represents camera information for unknown models with 32-bit integer format.
type CameraInfoUnknown32 struct {
	cameraModelID     CanonModelID      // Internal model ID, used for parsing variants.
	CameraTemperature CameraTemperature `json:"cameraTemperature"`
}

// CameraInfoUnknown16 represents camera information for unknown models with 16-bit integer format.
type CameraInfoUnknown16 struct {
	cameraModelID CanonModelID // Internal model ID, used for parsing variants.
}

// CameraInfoUnknown represents camera information for unknown models.
type CameraInfoUnknown struct {
	cameraModelID CanonModelID // Internal model ID, used for parsing variants.
}

// CameraOrientation represents the orientation of the camera
type CameraOrientation uint8

const (
	OrientationHorizontal  CameraOrientation = 0
	OrientationRotate90CW  CameraOrientation = 1
	OrientationRotate270CW CameraOrientation = 2
)

// OrientationString returns human readable orientation
func (o CameraOrientation) OrientationString() string {
	switch o {
	case OrientationHorizontal:
		return "Horizontal (normal)"
	case OrientationRotate90CW:
		return "Rotate 90 CW"
	case OrientationRotate270CW:
		return "Rotate 270 CW"
	default:
		return "Unknown"
	}
}

// PictureStyleInfo2 represents custom picture style information for EOS models
type PictureStyleInfo2 struct {
	// Standard settings
	ContrastStandard     int32 `json:"contrastStandard"`     // 0x00
	SharpnessStandard    int32 `json:"sharpnessStandard"`    // 0x04
	SaturationStandard   int32 `json:"saturationStandard"`   // 0x08
	ColorToneStandard    int32 `json:"colorToneStandard"`    // 0x0c
	FilterEffectStandard int32 `json:"filterEffectStandard"` // 0x10
	ToningEffectStandard int32 `json:"toningEffectStandard"` // 0x14

	// Portrait settings
	ContrastPortrait     int32 `json:"contrastPortrait"`     // 0x18
	SharpnessPortrait    int32 `json:"sharpnessPortrait"`    // 0x1c
	SaturationPortrait   int32 `json:"saturationPortrait"`   // 0x20
	ColorTonePortrait    int32 `json:"colorTonePortrait"`    // 0x24
	FilterEffectPortrait int32 `json:"filterEffectPortrait"` // 0x28
	ToningEffectPortrait int32 `json:"toningEffectPortrait"` // 0x2c

	// Landscape settings
	ContrastLandscape     int32 `json:"contrastLandscape"`     // 0x30
	SharpnessLandscape    int32 `json:"sharpnessLandscape"`    // 0x34
	SaturationLandscape   int32 `json:"saturationLandscape"`   // 0x38
	ColorToneLandscape    int32 `json:"colorToneLandscape"`    // 0x3c
	FilterEffectLandscape int32 `json:"filterEffectLandscape"` // 0x40
	ToningEffectLandscape int32 `json:"toningEffectLandscape"` // 0x44

	// Neutral settings
	ContrastNeutral     int32 `json:"contrastNeutral"`     // 0x48
	SharpnessNeutral    int32 `json:"sharpnessNeutral"`    // 0x4c
	SaturationNeutral   int32 `json:"saturationNeutral"`   // 0x50
	ColorToneNeutral    int32 `json:"colorToneNeutral"`    // 0x54
	FilterEffectNeutral int32 `json:"filterEffectNeutral"` // 0x58
	ToningEffectNeutral int32 `json:"toningEffectNeutral"` // 0x5c

	// Faithful settings
	ContrastFaithful     int32 `json:"contrastFaithful"`     // 0x60
	SharpnessFaithful    int32 `json:"sharpnessFaithful"`    // 0x64
	SaturationFaithful   int32 `json:"saturationFaithful"`   // 0x68
	ColorToneFaithful    int32 `json:"colorToneFaithful"`    // 0x6c
	FilterEffectFaithful int32 `json:"filterEffectFaithful"` // 0x70
	ToningEffectFaithful int32 `json:"toningEffectFaithful"` // 0x74

	// Monochrome settings
	ContrastMonochrome     int32        `json:"contrastMonochrome"`     // 0x78
	SharpnessMonochrome    int32        `json:"sharpnessMonochrome"`    // 0x7c
	SaturationMonochrome   int32        `json:"saturationMonochrome"`   // 0x80
	ColorToneMonochrome    int32        `json:"colorToneMonochrome"`    // 0x84
	FilterEffectMonochrome FilterEffect `json:"filterEffectMonochrome"` // 0x88
	ToningEffectMonochrome ToningEffect `json:"toningEffectMonochrome"` // 0x8c

	// Auto settings
	ContrastAuto     int32        `json:"contrastAuto"`     // 0x90
	SharpnessAuto    int32        `json:"sharpnessAuto"`    // 0x94
	SaturationAuto   int32        `json:"saturationAuto"`   // 0x98
	ColorToneAuto    int32        `json:"colorToneAuto"`    // 0x9c
	FilterEffectAuto FilterEffect `json:"filterEffectAuto"` // 0xa0
	ToningEffectAuto ToningEffect `json:"toningEffectAuto"` // 0xa4

	// User Defined 1
	ContrastUserDef1     int32        `json:"contrastUserDef1"`     // 0xa8
	SharpnessUserDef1    int32        `json:"sharpnessUserDef1"`    // 0xac
	SaturationUserDef1   int32        `json:"saturationUserDef1"`   // 0xb0
	ColorToneUserDef1    int32        `json:"colorToneUserDef1"`    // 0xb4
	FilterEffectUserDef1 FilterEffect `json:"filterEffectUserDef1"` // 0xb8
	ToningEffectUserDef1 ToningEffect `json:"toningEffectUserDef1"` // 0xbc

	// User Defined 2
	ContrastUserDef2     int32        `json:"contrastUserDef2"`     // 0xc0
	SharpnessUserDef2    int32        `json:"sharpnessUserDef2"`    // 0xc4
	SaturationUserDef2   int32        `json:"saturationUserDef2"`   // 0xc8
	ColorToneUserDef2    int32        `json:"colorToneUserDef2"`    // 0xcc
	FilterEffectUserDef2 FilterEffect `json:"filterEffectUserDef2"` // 0xd0
	ToningEffectUserDef2 ToningEffect `json:"toningEffectUserDef2"` // 0xd4

	// User Defined 3
	ContrastUserDef3     int32        `json:"contrastUserDef3"`     // 0xd8
	SharpnessUserDef3    int32        `json:"sharpnessUserDef3"`    // 0xdc
	SaturationUserDef3   int32        `json:"saturationUserDef3"`   // 0xe0
	ColorToneUserDef3    int32        `json:"colorToneUserDef3"`    // 0xe4
	FilterEffectUserDef3 FilterEffect `json:"filterEffectUserDef3"` // 0xe8
	ToningEffectUserDef3 ToningEffect `json:"toningEffectUserDef3"` // 0xec

	// Base picture styles
	UserDef1PictureStyle uint16 `json:"userDef1PictureStyle"` // 0xf0
	UserDef2PictureStyle uint16 `json:"userDef2PictureStyle"` // 0xf2
	UserDef3PictureStyle uint16 `json:"userDef3PictureStyle"` // 0xf4
}

// FilterEffect represents monochrome filter effects
type FilterEffect int32

const (
	FilterEffectNone   FilterEffect = 0
	FilterEffectYellow FilterEffect = 1
	FilterEffectOrange FilterEffect = 2
	FilterEffectRed    FilterEffect = 3
	FilterEffectGreen  FilterEffect = 4
	FilterEffectNA     FilterEffect = -559038737 // 0xdeadbeef
)

// String returns human readable filter effect
func (f FilterEffect) String() string {
	switch f {
	case FilterEffectNone:
		return "None"
	case FilterEffectYellow:
		return "Yellow"
	case FilterEffectOrange:
		return "Orange"
	case FilterEffectRed:
		return "Red"
	case FilterEffectGreen:
		return "Green"
	case FilterEffectNA:
		return "n/a"
	default:
		return "Unknown"
	}
}

// ToningEffect represents monochrome toning effects
type ToningEffect int32

const (
	ToningEffectNone   ToningEffect = 0
	ToningEffectSepia  ToningEffect = 1
	ToningEffectBlue   ToningEffect = 2
	ToningEffectPurple ToningEffect = 3
	ToningEffectGreen  ToningEffect = 4
	ToningEffectNA     ToningEffect = -559038737 // 0xdeadbeef
)

// String returns human readable toning effect
func (t ToningEffect) String() string {
	switch t {
	case ToningEffectNone:
		return "None"
	case ToningEffectSepia:
		return "Sepia"
	case ToningEffectBlue:
		return "Blue"
	case ToningEffectPurple:
		return "Purple"
	case ToningEffectGreen:
		return "Green"
	case ToningEffectNA:
		return "n/a"
	default:
		return "Unknown"
	}
}

// FNumber represents a camera f-stop value
type FNumber float32

// FromRaw converts a raw uint8 value to FNumber
func FNumberFromRaw(raw uint8) FNumber {
	if raw == 0 {
		return 0
	}
	// ValueConv: exp((val-8)/16*log(2))
	return FNumber(math.Exp(float64(raw-8) / 16.0 * math.Log(2)))
}

// ToRaw converts an FNumber to raw uint8 value
func (f FNumber) ToRaw() uint8 {
	if f == 0 {
		return 0
	}
	// ValueConvInv: log(val)*16/log(2)+8
	return uint8(math.Log(float64(f))*16.0/math.Log(2) + 8)
}

// String returns formatted f-stop value with 2 significant digits
func (f FNumber) String() string {
	if f == 0 {
		return "undefined"
	}
	return fmt.Sprintf("%.2g", float64(f))
}

// ISO represents a camera ISO value from CameraInfo.
type ISO uint32

// ISOFromRaw converts a raw uint8 value from Camera Info ISO to ISO
func ISOFromRaw(raw uint8) ISO {
	// ValueConv: 100*exp((val/8-9)*log(2))
	return ISO(100 * math.Exp((float64(raw)/8-9)*math.Log(2)))
}

// ToRaw converts an ISO to raw uint8 Camera Info ISO value
func (i ISO) ToRaw() uint8 {
	// ValueConvInv: (log(val/100)/log(2)+9)*8
	return uint8((math.Log(float64(i)/100)/math.Log(2) + 9) * 8)
}

// String returns formatted ISO value
func (i ISO) String() string {
	return fmt.Sprintf("%.0f", float64(i))
}

// CameraTemperature represents the camera temperature in Celsius from CameraInfo.
type CameraTemperature uint8

// CameraTemperatureFromRaw converts a raw uint8 value to CameraTemperature
func CameraTemperatureFromRaw(raw uint8) CameraTemperature {
	// ValueConv: val - 128
	return CameraTemperature(int16(raw) - 128)
}

// ToRaw converts a CameraTemperature to raw uint8 value
func (t CameraTemperature) ToRaw() uint8 {
	// ValueConvInv: val + 128
	return uint8(int16(t) + 128)
}

// String returns formatted temperature with Celsius suffix
func (t CameraTemperature) String() string {
	return fmt.Sprintf("%d°C", int8(t))
}

// CIFocalLength represents a camera focal length in millimeters from CameraInfo.
type CIFocalLength uint16

// FromRaw converts a raw big-endian uint16 value to CIFocalLength
func CIFocalLengthFromRaw(raw uint16) CIFocalLength {
	if raw == 0 {
		return 0
	}
	// Convert from big-endian
	return CIFocalLength(binary.BigEndian.Uint16([]byte{byte(raw >> 8), byte(raw)}))
}

// ToRaw converts a CIFocalLength to raw big-endian uint16 value
func (f CIFocalLength) ToRaw() uint16 {
	if f == 0 {
		return 0
	}
	// Convert to big-endian
	return binary.BigEndian.Uint16([]byte{byte(uint16(f) >> 8), byte(uint16(f))})
}

// String returns formatted focal length with mm suffix
func (f CIFocalLength) String() string {
	if f == 0 {
		return "undefined"
	}
	return fmt.Sprintf("%d mm", uint16(f))
}

// CIExposureTime represents a camera exposure time value from CameraInfo.
type CIExposureTime float32

// FromRaw converts a raw uint8 value to CIExposureTime
func CIExposureTimeFromRaw(raw uint8) CIExposureTime {
	if raw == 0 {
		return 0
	}
	// ValueConv: exp(4*log(2)*(1-CanonEv(val-24)))
	return CIExposureTime(math.Exp(4 * math.Log(2) * (1 - float64(CanonEv(float32(raw-24))))))
}

// ToRaw converts a CIExposureTime to raw uint8 value
func (e CIExposureTime) ToRaw() uint8 {
	if e == 0 {
		return 0
	}
	// ValueConvInv: CanonEvInv(1-log(val)/(4*log(2)))+24
	return uint8(CanonEvInv(float32(1-math.Log(float64(e))/(4*math.Log(2))) + 24))
}

// String returns exposure time as a fraction or decimal
func (e CIExposureTime) String() string {
	if e == 0 {
		return "undefined"
	}

	// For times >= 1 second, show decimal
	if e >= 1 {
		return fmt.Sprintf("%.1f", float64(e))
	}

	// For times < 1 second, show as fraction (1/x)
	denominator := 1.0 / float64(e)
	return fmt.Sprintf("1/%.0f", denominator)
}
