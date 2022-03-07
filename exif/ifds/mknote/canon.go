// Package mknote provides functions and types for decoding Exif Makernote values
package mknote

import (
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// TagCanonString returns the string representation of a tag.ID for Canon Makernotes
func TagCanonString(id tag.ID) string {
	name, ok := TagCanonIDMap[id]
	if !ok {
		return id.String()
	}
	return name
}

// TODO: TagTypeMap is a Map of tag.ID to default tag,Type
// var TagTypeMap = map[tag.ID]tag.Type{}

// TagCanonIDMap is a Map of tag.ID to string for the CanonMakerNote tags
var TagCanonIDMap = map[tag.ID]string{
	CanonCameraSettings:             "CanonCameraSettings",
	CanonFocalLength:                "CanonFocalLength",
	CanonFlashInfo:                  "CanonFlashInfo",
	CanonShotInfo:                   "CanonShotInfo",
	CanonPanorama:                   "CanonPanorama",
	CanonImageType:                  "CanonImageType",
	CanonFirmwareVersion:            "CanonFirmwareVersion",
	FileNumber:                      "FileNumber",
	OwnerName:                       "OwnerName",
	UnknownD30:                      "UnknownD30",
	SerialNumber:                    "SerialNumber",
	CanonCameraInfo:                 "CanonCameraInfo",
	CanonFileLength:                 "CanonFileLength",
	CustomFunctions:                 "CustomFunctions",
	CanonModelID:                    "CanonModelID",
	MovieInfo:                       "MovieInfo",
	CanonAFInfo:                     "CanonAFInfo",
	ThumbnailImageValidArea:         "ThumbnailImageValidArea",
	SerialNumberFormat:              "SerialNumberFormat",
	SuperMacro:                      "SuperMacro",
	DateStampMode:                   "DateStampMode",
	MyColors:                        "MyColors",
	FirmwareRevision:                "FirmwareRevision",
	Categories:                      "Categories",
	FaceDetect1:                     "FaceDetect1",
	FaceDetect2:                     "FaceDetect2",
	CanonAFInfo2:                    "CanonAFInfo2",
	ContrastInfo:                    "ContrastInfo",
	ImageUniqueID:                   "ImageUniqueID",
	WBInfo:                          "WBInfo",
	FaceDetect3:                     "FaceDetect3",
	TimeInfo:                        "TimeInfo",
	BatteryType:                     "BatteryType",
	AFInfo3:                         "AFInfo3",
	RawDataOffset:                   "RawDataOffset",
	OriginalDecisionDataOffset:      "OriginalDecisionDataOffset",
	CustomFunctions1D:               "CustomFunctions1D",
	PersonalFunctions:               "PersonalFunctions",
	PersonalFunctionValues:          "PersonalFunctionValues",
	CanonFileInfo:                   "CanonFileInfo",
	AFPointsInFocus1D:               "AFPointsInFocus1D",
	LensModel:                       "LensModel",
	CanonInternalSerialNumber:       "CanonInternalSerialNumber",
	CanonDustRemovalData:            "CanonDustRemovalData",
	CanonCustomFunctions:            "CanonCustomFunctions",
	CanonAspectInfo:                 "CanonAspectInfo",
	CanonProcessingInfo:             "CanonProcessingInfo",
	CanonToneCurveTable:             "CanonToneCurveTable",
	CanonSharpnessTable:             "CanonSharpnessTable",
	CanonSharpnessFreqTable:         "CanonSharpnessFreqTable",
	CanonWhiteBalanceTable:          "CanonWhiteBalanceTable",
	CanonColorBalance:               "CanonColorBalance",
	CanonMeasuredColor:              "CanonMeasuredColor",
	CanonColorTemperature:           "CanonColorTemperature",
	CanonCanonFlags:                 "CanonCanonFlags",
	CanonModifiedInfo:               "CanonModifiedInfo",
	CanonToneCurveMatching:          "CanonToneCurveMatching",
	CanonWhiteBalanceMatching:       "CanonWhiteBalanceMatching",
	CanonColorSpace:                 "CanonColorSpace",
	Canon0x00b5:                     "Canon0x00b5",
	CanonPreviewImageInfo:           "CanonPreviewImageInfo",
	Canon0x00c0:                     "Canon0x00c0",
	Canon0x00c1:                     "Canon0x00c1",
	CanonVRDOffset:                  "CanonVRDOffset",
	CanonSensorInfo:                 "CanonSensorInfo",
	CanonAFInfoSize:                 "CanonAFInfoSize",
	CanonAFAreaMode:                 "CanonAFAreaMode",
	CanonAFNumPoints:                "CanonAFNumPoints",
	CanonAFValidPoints:              "CanonAFValidPoints",
	CanonAFCanonImageWidth:          "CanonAFCanonImageWidth",
	CanonAFCanonImageHeight:         "CanonAFCanonImageHeight",
	CanonAFImageWidth:               "CanonAFImageWidth",
	CanonAFImageHeight:              "CanonAFImageHeight",
	CanonAFAreaWidths:               "CanonAFAreaWidths",
	CanonAFAreaHeights:              "CanonAFAreaHeights",
	CanonAFXPositions:               "CanonAFXPositions",
	CanonAFYPositions:               "CanonAFYPositions",
	CanonAFPointsInFocus:            "CanonAFPointsInFocus",
	CanonAFPointsSelected:           "CanonAFPointsSelected",
	CanonAFPointsUnusable:           "CanonAFPointsUnusable",
	CanonColorData:                  "CanonColorData",
	CanonCRWParam:                   "CanonCRWParam",
	CanonColorInfo:                  "CanonColorInfo",
	CanonFlavor:                     "CanonFlavor",
	CanonPictureStyleUserDef:        "CanonPictureStyleUserDef",
	CanonCustomPictureStyleFileName: "CanonCustomPictureStyleFileName",
	CanonAFMicroAdj:                 "CanonAFMicroAdj",
	CanonVignettingCorr:             "CanonVignettingCorr",
	CanonVignettingCorr2:            "CanonVignettingCorr2",
	CanonLensInfo:                   "CanonLensInfo",
	CanonAmbienceInfo:               "CanonAmbienceInfo",
	CanonMultiExp:                   "CanonMultiExp",
	CanonFilterInfo:                 "CanonFilterInfo",
	CanonHDRInfo:                    "CanonHDRInfo",
	CanonAFConfig:                   "CanonAFConfig",
	CanonRawBurstModeRoll:           "CanonRawBurstModeRoll",
}

// CanonMKnoteIFD TagIDs
// Source: https://exiftool.org/TagNames/Canon.html on 8/05/2020
// Secondary Source: https://exiv2.org/tags-canon.html on Mar 6, 2022
const (
	CanonCameraSettings             tag.ID = 0x0001
	CanonFocalLength                tag.ID = 0x0002
	CanonFlashInfo                  tag.ID = 0x0003
	CanonShotInfo                   tag.ID = 0x0004
	CanonPanorama                   tag.ID = 0x0005
	CanonImageType                  tag.ID = 0x0006
	CanonFirmwareVersion            tag.ID = 0x0007
	FileNumber                      tag.ID = 0x0008
	OwnerName                       tag.ID = 0x0009
	UnknownD30                      tag.ID = 0x000a
	SerialNumber                    tag.ID = 0x000c
	CanonCameraInfo                 tag.ID = 0x000d // WIP
	CanonFileLength                 tag.ID = 0x000e // WIP
	CustomFunctions                 tag.ID = 0x000f // WIP
	CanonModelID                    tag.ID = 0x0010
	MovieInfo                       tag.ID = 0x0011 // WIP
	CanonAFInfo                     tag.ID = 0x0012
	ThumbnailImageValidArea         tag.ID = 0x0013 // WIP
	SerialNumberFormat              tag.ID = 0x0015 // WIP
	SuperMacro                      tag.ID = 0x001a // WIP
	DateStampMode                   tag.ID = 0x001c // WIP
	MyColors                        tag.ID = 0x001d // WIP
	FirmwareRevision                tag.ID = 0x001e // WIP
	Categories                      tag.ID = 0x0023 // WIP
	FaceDetect1                     tag.ID = 0x0024 // WIP
	FaceDetect2                     tag.ID = 0x0025 // WIP
	CanonAFInfo2                    tag.ID = 0x0026
	ContrastInfo                    tag.ID = 0x0027 // WIP
	ImageUniqueID                   tag.ID = 0x0028 // WIP
	WBInfo                          tag.ID = 0x0029 // WIP
	FaceDetect3                     tag.ID = 0x002f // WIP
	TimeInfo                        tag.ID = 0x0035
	BatteryType                     tag.ID = 0x0038 // WIP
	AFInfo3                         tag.ID = 0x003c // WIP
	RawDataOffset                   tag.ID = 0x0081 // WIP
	OriginalDecisionDataOffset      tag.ID = 0x0083 // WIP
	CustomFunctions1D               tag.ID = 0x0090 // WIP
	PersonalFunctions               tag.ID = 0x0091 // WIP
	PersonalFunctionValues          tag.ID = 0x0092 // WIP
	CanonFileInfo                   tag.ID = 0x0093
	AFPointsInFocus1D               tag.ID = 0x0094 // WIP
	LensModel                       tag.ID = 0x0095
	CanonInternalSerialNumber       tag.ID = 0x0096 //	ASCII
	CanonDustRemovalData            tag.ID = 0x0097 //	Ascii	Dust removal data
	CanonCustomFunctions            tag.ID = 0x0099 //	Short	Custom functions
	CanonAspectInfo                 tag.ID = 0x009a //	Short	AspectInfo
	CanonProcessingInfo             tag.ID = 0x00a0 //	Short	Processing info
	CanonToneCurveTable             tag.ID = 0x00a1 //	Short	ToneCurveTable
	CanonSharpnessTable             tag.ID = 0x00a2 //	Short	SharpnessTable
	CanonSharpnessFreqTable         tag.ID = 0x00a3 //	Short	SharpnessTable
	CanonWhiteBalanceTable          tag.ID = 0x00a4 //	Short	SharpnessTable
	CanonColorBalance               tag.ID = 0x00a9 //	Short	ColorBalance
	CanonMeasuredColor              tag.ID = 0x00aa //	Short	Measured color
	CanonColorTemperature           tag.ID = 0x00ae //	Short	ColorTemperature
	CanonCanonFlags                 tag.ID = 0x00b0 //	Short	CanonFlags
	CanonModifiedInfo               tag.ID = 0x00b1 //	Short	ModifiedInfo
	CanonToneCurveMatching          tag.ID = 0x00b2 //	Short	ToneCurveMatching
	CanonWhiteBalanceMatching       tag.ID = 0x00b3 //	Short	WhiteBalanceMatching
	CanonColorSpace                 tag.ID = 0x00b4 //	SShort	ColorSpace
	Canon0x00b5                     tag.ID = 0x00b5 //	Short	Unknown
	CanonPreviewImageInfo           tag.ID = 0x00b6 //	Short	PreviewImageInfo
	Canon0x00c0                     tag.ID = 0x00c0 //	Short	Unknown
	Canon0x00c1                     tag.ID = 0x00c1 //	Short	Unknown
	CanonVRDOffset                  tag.ID = 0x00d0 //	Long	VRD offset
	CanonSensorInfo                 tag.ID = 0x00e0 //	Short	Sensor info
	CanonAFInfoSize                 tag.ID = 0x2600 //	SShort	AF InfoSize
	CanonAFAreaMode                 tag.ID = 0x2601 //	SShort	AF Area Mode
	CanonAFNumPoints                tag.ID = 0x2602 //	SShort	AF NumPoints
	CanonAFValidPoints              tag.ID = 0x2603 //	SShort	AF ValidPoints
	CanonAFCanonImageWidth          tag.ID = 0x2604 //	SShort	AF ImageWidth
	CanonAFCanonImageHeight         tag.ID = 0x2605 //	SShort	AF ImageHeight
	CanonAFImageWidth               tag.ID = 0x2606 //	SShort	AF Width
	CanonAFImageHeight              tag.ID = 0x2607 //	SShort	AF Height
	CanonAFAreaWidths               tag.ID = 0x2608 //	SShort	AF Area Widths
	CanonAFAreaHeights              tag.ID = 0x2609 //	SShort	AF Area Heights
	CanonAFXPositions               tag.ID = 0x260a //	SShort	AF X Positions
	CanonAFYPositions               tag.ID = 0x260b //	SShort	AF Y Positions
	CanonAFPointsInFocus            tag.ID = 0x260c //	SShort	AF Points in Focus
	CanonAFPointsSelected           tag.ID = 0x260d //	SShort	AF Points Selected
	CanonAFPointsUnusable           tag.ID = 0x260e //	SShort	AF Points Unusable
	CanonColorData                  tag.ID = 0x4001 //	Short	Color data
	CanonCRWParam                   tag.ID = 0x4002 //	Short	CRWParam
	CanonColorInfo                  tag.ID = 0x4003 //	Short	ColorInfo
	CanonFlavor                     tag.ID = 0x4005 //	Short	Flavor
	CanonPictureStyleUserDef        tag.ID = 0x4008 //	Short	PictureStyleUserDef
	CanonCustomPictureStyleFileName tag.ID = 0x4010 //	Short	CustomPictureStyleFileName
	CanonAFMicroAdj                 tag.ID = 0x4013 //	Short	AFMicroAdj
	CanonVignettingCorr             tag.ID = 0x4015 //	Short	VignettingCorr
	CanonVignettingCorr2            tag.ID = 0x4016 //	Short	VignettingCorr2
	//CanonLightingOpt                tag.ID = 0x4018 //	Short	LightingOpt
	CanonLensInfo         tag.ID = 0x4018 //	Short	LensInfo
	CanonAmbienceInfo     tag.ID = 0x4020 //	Short	AmbienceInfo
	CanonMultiExp         tag.ID = 0x4021 //	Short	MultiExp
	CanonFilterInfo       tag.ID = 0x4024 //	Short	FilterInfo
	CanonHDRInfo          tag.ID = 0x4025 //	Short	HDRInfo
	CanonAFConfig         tag.ID = 0x4028 //	Short	AFConfig
	CanonRawBurstModeRoll tag.ID = 0x403f //	Short	RawBurstModeRoll
)
