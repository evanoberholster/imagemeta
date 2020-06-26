// Package exififd provides types for "RootIfd/ExifIfd"
package exififd

import "github.com/evanoberholster/exiftool/tag"

// TagIDMap is a Map of tag.ID to string for the ExifIfd tags
var TagIDMap = map[tag.ID]string{
	ExposureTime:              "ExposureTime",
	FNumber:                   "FNumber",
	ExposureProgram:           "ExposureProgram",
	SpectralSensitivity:       "SpectralSensitivity",
	ISOSpeedRatings:           "ISOSpeedRatings",
	OECF:                      "OECF",
	SensitivityType:           "SensitivityType",
	StandardOutputSensitivity: "StandardOutputSensitivity",
	RecommendedExposureIndex:  "RecommendedExposureIndex",
	ISOSpeed:                  "ISOSpeed",
	ISOSpeedLatitudeyyy:       "ISOSpeedLatitudeyyy",
	ISOSpeedLatitudezzz:       "ISOSpeedLatitudezzz",
	ExifVersion:               "ExifVersion",
	DateTimeOriginal:          "DateTimeOriginal",
	DateTimeDigitized:         "DateTimeDigitized",
	ComponentsConfiguration:   "ComponentsConfiguration",
	CompressedBitsPerPixel:    "CompressedBitsPerPixel",
	ShutterSpeedValue:         "ShutterSpeedValue",
	ApertureValue:             "ApertureValue",
	BrightnessValue:           "BrightnessValue",
	ExposureBiasValue:         "ExposureBiasValue",
	MaxApertureValue:          "MaxApertureValue",
	SubjectDistance:           "SubjectDistance",
	MeteringMode:              "MeteringMode",
	LightSource:               "LightSource",
	Flash:                     "Flash",
	FocalLength:               "FocalLength",
	SubjectArea:               "SubjectArea",
	MakerNote:                 "MakerNote",
	UserComment:               "UserComment",
	SubSecTime:                "SubSecTime",
	SubSecTimeOriginal:        "SubSecTimeOriginal",
	SubSecTimeDigitized:       "SubSecTimeDigitized",
	FlashpixVersion:           "FlashpixVersion",
	ColorSpace:                "ColorSpace",
	PixelXDimension:           "PixelXDimension",
	PixelYDimension:           "PixelYDimension",
	RelatedSoundFile:          "RelatedSoundFile",
	InteroperabilityTag:       "InteroperabilityTag",
	FlashEnergy:               "FlashEnergy",
	SpatialFrequencyResponse:  "SpatialFrequencyResponse",
	FocalPlaneXResolution:     "FocalPlaneXResolution",
	FocalPlaneYResolution:     "FocalPlaneYResolution",
	FocalPlaneResolutionUnit:  "FocalPlaneResolutionUnit",
	SubjectLocation:           "SubjectLocation",
	ExposureIndex:             "ExposureIndex",
	SensingMethod:             "SensingMethod",
	FileSource:                "FileSource",
	SceneType:                 "SceneType",
	CFAPattern:                "CFAPattern",
	CustomRendered:            "CustomRendered",
	ExposureMode:              "ExposureMode",
	WhiteBalance:              "WhiteBalance",
	DigitalZoomRatio:          "DigitalZoomRatio",
	FocalLengthIn35mmFilm:     "FocalLengthIn35mmFilm",
	SceneCaptureType:          "SceneCaptureType",
	GainControl:               "GainControl",
	Contrast:                  "Contrast",
	Saturation:                "Saturation",
	Sharpness:                 "Sharpness",
	DeviceSettingDescription:  "DeviceSettingDescription",
	SubjectDistanceRange:      "SubjectDistanceRange",
	ImageUniqueID:             "ImageUniqueID",
	CameraOwnerName:           "CameraOwnerName",
	BodySerialNumber:          "BodySerialNumber",
	LensSpecification:         "LensSpecification",
	LensMake:                  "LensMake",
	LensModel:                 "LensModel",
	LensSerialNumber:          "LensSerialNumber",
}

// ExifIFD TagIDs
const (
	ExposureTime              tag.ID = 0x829a
	FNumber                   tag.ID = 0x829d
	ExposureProgram           tag.ID = 0x8822
	SpectralSensitivity       tag.ID = 0x8824
	ISOSpeedRatings           tag.ID = 0x8827
	OECF                      tag.ID = 0x8828
	SensitivityType           tag.ID = 0x8830
	StandardOutputSensitivity tag.ID = 0x8831
	RecommendedExposureIndex  tag.ID = 0x8832
	ISOSpeed                  tag.ID = 0x8833
	ISOSpeedLatitudeyyy       tag.ID = 0x8834
	ISOSpeedLatitudezzz       tag.ID = 0x8835
	ExifVersion               tag.ID = 0x9000
	DateTimeOriginal          tag.ID = 0x9003
	DateTimeDigitized         tag.ID = 0x9004
	ComponentsConfiguration   tag.ID = 0x9101
	CompressedBitsPerPixel    tag.ID = 0x9102
	ShutterSpeedValue         tag.ID = 0x9201
	ApertureValue             tag.ID = 0x9202
	BrightnessValue           tag.ID = 0x9203
	ExposureBiasValue         tag.ID = 0x9204
	MaxApertureValue          tag.ID = 0x9205
	SubjectDistance           tag.ID = 0x9206
	MeteringMode              tag.ID = 0x9207
	LightSource               tag.ID = 0x9208
	Flash                     tag.ID = 0x9209
	FocalLength               tag.ID = 0x920a
	SubjectArea               tag.ID = 0x9214
	MakerNote                 tag.ID = 0x927c
	UserComment               tag.ID = 0x9286
	SubSecTime                tag.ID = 0x9290 // fractional seconds for ModifyDate
	SubSecTimeOriginal        tag.ID = 0x9291 // fractional seconds for DateTimeOriginal
	SubSecTimeDigitized       tag.ID = 0x9292 // fractional seconds for CreateDate
	FlashpixVersion           tag.ID = 0xa000
	ColorSpace                tag.ID = 0xa001
	PixelXDimension           tag.ID = 0xa002
	PixelYDimension           tag.ID = 0xa003
	RelatedSoundFile          tag.ID = 0xa004
	InteroperabilityTag       tag.ID = 0xa005
	FlashEnergy               tag.ID = 0xa20b
	SpatialFrequencyResponse  tag.ID = 0xa20c
	FocalPlaneXResolution     tag.ID = 0xa20e
	FocalPlaneYResolution     tag.ID = 0xa20f
	FocalPlaneResolutionUnit  tag.ID = 0xa210
	SubjectLocation           tag.ID = 0xa214
	ExposureIndex             tag.ID = 0xa215
	SensingMethod             tag.ID = 0xa217
	FileSource                tag.ID = 0xa300
	SceneType                 tag.ID = 0xa301
	CFAPattern                tag.ID = 0xa302
	CustomRendered            tag.ID = 0xa401
	ExposureMode              tag.ID = 0xa402
	WhiteBalance              tag.ID = 0xa403
	DigitalZoomRatio          tag.ID = 0xa404
	FocalLengthIn35mmFilm     tag.ID = 0xa405
	SceneCaptureType          tag.ID = 0xa406
	GainControl               tag.ID = 0xa407
	Contrast                  tag.ID = 0xa408
	Saturation                tag.ID = 0xa409
	Sharpness                 tag.ID = 0xa40a
	DeviceSettingDescription  tag.ID = 0xa40b
	SubjectDistanceRange      tag.ID = 0xa40c
	ImageUniqueID             tag.ID = 0xa420
	CameraOwnerName           tag.ID = 0xa430
	BodySerialNumber          tag.ID = 0xa431
	LensSpecification         tag.ID = 0xa432
	LensMake                  tag.ID = 0xa433
	LensModel                 tag.ID = 0xa434
	LensSerialNumber          tag.ID = 0xa435
)
