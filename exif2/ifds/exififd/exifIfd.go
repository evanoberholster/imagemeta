// Package exififd provides types for "RootIfd/ExifIfd"
package exififd

import (
	"github.com/evanoberholster/imagemeta/exif2/tag"
)

func TagString(id tag.ID) string {
	name, ok := TagIDMap[id]
	if !ok {
		return id.String()
	}
	return name
}

// TODO: TagTypeMap is a Map of tag.ID to default tag,Type
// var TagTypeMap = map[tag.ID]tag.Type{}

// TagIDMap is a Map of tag.ID to string for the ExifIfd tags
var TagIDMap = map[tag.ID]string{
	Acceleration:               "Acceleration",
	AmbientTemperature:         "AmbientTemperature",
	ApertureValue:              "ApertureValue",
	BodySerialNumber:           "BodySerialNumber",
	BrightnessValue:            "BrightnessValue",
	CameraElevationAngle:       "CameraElevationAngle",
	CameraOwnerName:            "CameraOwnerName",
	CFAPattern:                 "CFAPattern",
	ColorSpace:                 "ColorSpace",
	ComponentsConfiguration:    "ComponentsConfiguration",
	CompositeImage:             "CompositeImage",
	CompositeImageCount:        "CompositeImageCount",
	CompositeImageExposureTime: "CompositeImageExposureTime",
	CompressedBitsPerPixel:     "CompressedBitsPerPixel",
	Contrast:                   "Contrast",
	CustomRendered:             "CustomRendered",
	DateTimeDigitized:          "DateTimeDigitized",
	DateTimeOriginal:           "DateTimeOriginal",
	DeviceSettingDescription:   "DeviceSettingDescription",
	DigitalZoomRatio:           "DigitalZoomRatio",
	ExifVersion:                "ExifVersion",
	ExposureBiasValue:          "ExposureBiasValue",
	ExposureIndex:              "ExposureIndex",
	ExposureMode:               "ExposureMode",
	ExposureProgram:            "ExposureProgram",
	ExposureTime:               "ExposureTime",
	FileSource:                 "FileSource",
	Flash:                      "Flash",
	FlashEnergy:                "FlashEnergy",
	FlashpixVersion:            "FlashpixVersion",
	FNumber:                    "FNumber",
	FocalLength:                "FocalLength",
	FocalLengthIn35mmFilm:      "FocalLengthIn35mmFilm",
	FocalPlaneResolutionUnit:   "FocalPlaneResolutionUnit",
	FocalPlaneXResolution:      "FocalPlaneXResolution",
	FocalPlaneYResolution:      "FocalPlaneYResolution",
	GainControl:                "GainControl",
	Gamma:                      "Gamma",
	GooglePlusUploadCode:       "GooglePlusUploadCode",
	Humidity:                   "Humidity",
	ImageHistory:               "ImageHistory",
	ImageNumber:                "ImageNumber",
	ImageUniqueID:              "ImageUniqueID",
	InteroperabilityTag:        "InteroperabilityTag",
	PhotographicSensitivity:    "PhotographicSensitivity",
	ISOSpeed:                   "ISOSpeed",
	ISOSpeedLatitudeyyy:        "ISOSpeedLatitudeyyy",
	ISOSpeedLatitudezzz:        "ISOSpeedLatitudezzz",
	LensMake:                   "LensMake",
	LensModel:                  "LensModel",
	LensSerialNumber:           "LensSerialNumber",
	LensSpecification:          "LensSpecification",
	LightSource:                "LightSource",
	MakerNote:                  "MakerNote",
	MaxApertureValue:           "MaxApertureValue",
	MeteringMode:               "MeteringMode",
	OECF:                       "OECF",
	OffsetTime:                 "OffsetTime",
	OffsetTimeDigitized:        "OffsetTimeDigitized",
	OffsetTimeOriginal:         "OffsetTimeOriginal",
	PixelXDimension:            "PixelXDimension",
	PixelYDimension:            "PixelYDimension",
	Pressure:                   "Pressure",
	RecommendedExposureIndex:   "RecommendedExposureIndex",
	RelatedSoundFile:           "RelatedSoundFile",
	Saturation:                 "Saturation",
	SceneCaptureType:           "SceneCaptureType",
	SceneType:                  "SceneType",
	SecurityClassification:     "SecurityClassification",
	SelfTimerMode:              "SelfTimerMode",
	SensingMethod:              "SensingMethod",
	SensitivityType:            "SensitivityType",
	Sharpness:                  "Sharpness",
	ShutterSpeedValue:          "ShutterSpeedValue",
	SpatialFrequencyResponse:   "SpatialFrequencyResponse",
	SpectralSensitivity:        "SpectralSensitivity",
	StandardOutputSensitivity:  "StandardOutputSensitivity",
	SubjectArea:                "SubjectArea",
	SubjectDistance:            "SubjectDistance",
	SubjectDistanceRange:       "SubjectDistanceRange",
	SubjectLocation:            "SubjectLocation",
	SubSecTime:                 "SubSecTime",
	SubSecTimeDigitized:        "SubSecTimeDigitized",
	SubSecTimeOriginal:         "SubSecTimeOriginal",
	TimeZoneOffset:             "TimeZoneOffset",
	UserComment:                "UserComment",
	WaterDepth:                 "WaterDepth",
	WhiteBalance:               "WhiteBalance",
}

// ExifIFD TagIDs
const (
	ExposureTime               tag.ID = 0x829a
	FNumber                    tag.ID = 0x829d
	ExposureProgram            tag.ID = 0x8822
	SpectralSensitivity        tag.ID = 0x8824
	PhotographicSensitivity    tag.ID = 0x8827
	OECF                       tag.ID = 0x8828
	TimeZoneOffset             tag.ID = 0x882a // int16s[n]
	SelfTimerMode              tag.ID = 0x882b // int16u
	SensitivityType            tag.ID = 0x8830
	StandardOutputSensitivity  tag.ID = 0x8831
	RecommendedExposureIndex   tag.ID = 0x8832
	ISOSpeed                   tag.ID = 0x8833
	ISOSpeedLatitudeyyy        tag.ID = 0x8834
	ISOSpeedLatitudezzz        tag.ID = 0x8835
	ExifVersion                tag.ID = 0x9000
	DateTimeOriginal           tag.ID = 0x9003
	DateTimeDigitized          tag.ID = 0x9004
	GooglePlusUploadCode       tag.ID = 0x9009 // undef[n]
	OffsetTime                 tag.ID = 0x9010
	OffsetTimeOriginal         tag.ID = 0x9011
	OffsetTimeDigitized        tag.ID = 0x9012
	ComponentsConfiguration    tag.ID = 0x9101 // undef[4]!:
	CompressedBitsPerPixel     tag.ID = 0x9102
	ShutterSpeedValue          tag.ID = 0x9201
	ApertureValue              tag.ID = 0x9202
	BrightnessValue            tag.ID = 0x9203
	ExposureBiasValue          tag.ID = 0x9204
	MaxApertureValue           tag.ID = 0x9205
	SubjectDistance            tag.ID = 0x9206
	MeteringMode               tag.ID = 0x9207
	LightSource                tag.ID = 0x9208
	Flash                      tag.ID = 0x9209
	FocalLength                tag.ID = 0x920a
	ImageNumber                tag.ID = 0x9211 // int32u
	SecurityClassification     tag.ID = 0x9212 // string
	ImageHistory               tag.ID = 0x9213 // string
	SubjectArea                tag.ID = 0x9214 // int16u[n]
	MakerNote                  tag.ID = 0x927c
	UserComment                tag.ID = 0x9286
	SubSecTime                 tag.ID = 0x9290 // fractional seconds for ModifyDate
	SubSecTimeOriginal         tag.ID = 0x9291 // fractional seconds for DateTimeOriginal
	SubSecTimeDigitized        tag.ID = 0x9292 // fractional seconds for CreateDate
	AmbientTemperature         tag.ID = 0x9400 // rational64s
	Humidity                   tag.ID = 0x9401 // rational64u
	Pressure                   tag.ID = 0x9402 // rational64u
	WaterDepth                 tag.ID = 0x9403 // rational64s
	Acceleration               tag.ID = 0x9404 // rational64u
	CameraElevationAngle       tag.ID = 0x9405 // rational64s
	FlashpixVersion            tag.ID = 0xa000
	ColorSpace                 tag.ID = 0xa001
	PixelXDimension            tag.ID = 0xa002 //  ExifImageWidth   int16u:
	PixelYDimension            tag.ID = 0xa003 //  ExifImageHeight  int16u:
	RelatedSoundFile           tag.ID = 0xa004
	InteroperabilityTag        tag.ID = 0xa005
	FlashEnergy                tag.ID = 0xa20b
	SpatialFrequencyResponse   tag.ID = 0xa20c
	FocalPlaneXResolution      tag.ID = 0xa20e // rational64u
	FocalPlaneYResolution      tag.ID = 0xa20f // rational64u
	FocalPlaneResolutionUnit   tag.ID = 0xa210 // int16u
	SubjectLocation            tag.ID = 0xa214 // int16u[2]
	ExposureIndex              tag.ID = 0xa215 // rational64u
	SensingMethod              tag.ID = 0xa217
	FileSource                 tag.ID = 0xa300
	SceneType                  tag.ID = 0xa301
	CFAPattern                 tag.ID = 0xa302
	CustomRendered             tag.ID = 0xa401
	ExposureMode               tag.ID = 0xa402
	WhiteBalance               tag.ID = 0xa403
	DigitalZoomRatio           tag.ID = 0xa404
	FocalLengthIn35mmFilm      tag.ID = 0xa405
	SceneCaptureType           tag.ID = 0xa406
	GainControl                tag.ID = 0xa407
	Contrast                   tag.ID = 0xa408
	Saturation                 tag.ID = 0xa409
	Sharpness                  tag.ID = 0xa40a
	DeviceSettingDescription   tag.ID = 0xa40b
	SubjectDistanceRange       tag.ID = 0xa40c
	ImageUniqueID              tag.ID = 0xa420
	CameraOwnerName            tag.ID = 0xa430
	BodySerialNumber           tag.ID = 0xa431
	LensSpecification          tag.ID = 0xa432 // rational64u[4]
	LensMake                   tag.ID = 0xa433
	LensModel                  tag.ID = 0xa434
	LensSerialNumber           tag.ID = 0xa435
	CompositeImage             tag.ID = 0xa460
	CompositeImageCount        tag.ID = 0xa461 // int16u[2]
	CompositeImageExposureTime tag.ID = 0xa462 // undef
	Gamma                      tag.ID = 0xa500 // rational64u
)

// ExifIfd is the Exif Ifd
// Based on Exif 2.23 Sepc https://web.archive.org/web/20190624045241if_/http://www.cipa.jp:80/std/documents/e/DC-008-Translation-2019-E.pdf
type ExifIfd struct {
	ExifVersion              [4]byte
	PixelXDimension          uint32
	PixelYDimension          uint32
	ComponentsConfiguration  [4]uint8
	UserComment              string
	DateTimeOriginal         [20]byte
	DateTimeDigitized        [20]byte
	OffsetTime               [7]byte
	OffsetTimeOriginal       [7]byte
	OffsetTimeDigitized      [7]byte
	SubsecTime               uint16
	SubsecTimeOriginal       uint16
	SubsecTimeDigitized      uint16
	ExposureTime             tag.Rational64
	FNumber                  tag.Rational64
	ExposureProgram          tag.ExposureProgram
	ISOSpeed                 uint32
	ShutterSpeedValue        tag.ShutterSpeedValue // APEX
	ApertureValue            tag.ApertureValue     // APEX
	BrightnessValue          tag.BrightnessValue   // APEX
	ExposureBiasValue        tag.ExposureBiasValue // APEX
	SubjectDistance          tag.Rational64
	LightSource              tag.LightSource
	Flash                    tag.Flash
	FocalLength              tag.FocalLength
	FocalPlaneXResolution    tag.Rational64
	FocalPlaneYResolution    tag.Rational64
	FocalPlaneResolutionUnit uint16
	ExposureMode             tag.ExposureMode
	DigitalZoomRatio         tag.Rational64
	FocalLengthIn35mmFilm    tag.FocalLength
	//GainControl              uint16 // values
	//Contrast                 uint16 // values
	//Saturation               uint16 // values
	//Sharpness                uint16 // values
	//
	ImageUnique       tag.ImageUnique // ASCII hexadecimal notation 33 bytes
	CameraOwnerName   string
	BodySerialNumber  string
	LensSpecification [4]tag.Rational64
	LensMake          string
	LensModel         string
	LensSerialNumber  string
	//other                    map[tag.ID]string
}
