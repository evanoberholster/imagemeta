package xmpns

import "fmt"

// Name is the XMP Property name
type Name uint8

func (n Name) String() string {
	return fmt.Sprintf(mapNameString[n])
}

// IdentifyName returns the XMP Property Name correspondent to buf.
// If Property Name was not identified returns UnknownName.
func IdentifyName(buf []byte) (n Name) {
	return mapStringName[string(buf)]
}

// Names
const (
	UnknownPropertyName Name = iota

	VideoFrameRate
	VideoPixelDepth
	VideoPixelAspectRatio
	VideoFieldOrder
	TapeName
	AltTapeName
	About // about
	Action
	AlreadyApplied
	Alt
	ApertureValue
	ApproximateFocusDistance
	Bag
	BodySerialNumber
	CameraOwnerName
	Changed
	ColorMode
	ColorSpace
	ComponentsConfiguration
	Compression
	Contrast
	CreateDate
	Creator
	CreatorTool
	CustomRendered
	DateCreated
	DateTimeDigitized
	DateTimeOriginal
	DerivedFrom
	Description
	DigitalZoomRatio
	DocumentID
	EmbeddedXMPDigest
	ExifVersion
	ExposureBiasValue
	ExposureMode
	ExposureProgram
	ExposureTime
	FileSource
	Fired
	Firmware
	Flash
	FlashCompensation
	FlashpixVersion
	FNumber
	FocalLength
	FocalPlaneResolutionUnit
	FocalPlaneXResolution
	FocalPlaneYResolution
	Format // format
	Function
	GainControl
	GPSAltitude
	GPSAltitudeRef
	GPSLatitude
	GPSLongitude
	GPSMapDatum
	GPSTimeStamp
	GPSVersionID
	HierarchicalSubject
	History
	ICCProfile
	ImageLength
	ImageNumber
	ImageWidth
	InstanceID
	InteroperabilityIndex
	ISOSpeedRatings
	Label
	Lang
	LegacyIPTCDigest
	Lens
	LensID
	LensInfo
	LensModel
	LensSerialNumber
	Li
	LightSource
	Make
	MaxApertureValue
	MetadataDate
	MeteringMode
	Mode
	Model
	ModifyDate
	NativeDigest
	Orientation
	OriginalDocumentID
	ParseType // parseType
	PhotometricInterpretation
	PixelXDimension
	PixelYDimension
	PlanarConfiguration
	Rating
	RawFileName
	RDF
	RecommendedExposureIndex
	RedEyeMode
	ResolutionUnit
	Return
	Rights
	SamplesPerPixel
	Saturation
	SceneCaptureType
	SceneType
	SensitivityType
	Seq
	SerialNumber
	Sharpness
	ShutterSpeedValue
	SidecarForExtension
	SoftwareAgent
	Subject
	SubjectDistance
	Temperature
	Title
	ToneCurve
	ToneCurveBlue
	ToneCurveGreen
	ToneCurvePV2012
	ToneCurvePV2012Blue
	ToneCurvePV2012Green
	ToneCurvePV2012Red
	ToneCurveRed
	UserComment
	When
	WhiteBalance
	XmpMeta
	XResolution
	YResolution

	XmpDM
	StDim
	Xap
	Dc
	VideoFrameSize
	W
	H
	StartTimecode
	TimeValue
	AltTimecode
)

// mapNameString returns Name's value as a string
var mapNameString = map[Name]string{
	UnknownPropertyName:       "Unknown",
	VideoFrameSize:            "videoFrameSize",
	W:                         "w",
	H:                         "h",
	StartTimecode:             "startTimecode",
	TimeValue:                 "timeValue",
	AltTimecode:               "altTimecode",
	XmpDM:                     "xmpDM",
	StDim:                     "stDim",
	Xap:                       "xap",
	Dc:                        "dc",
	VideoFrameRate:            "videoFrameRate",
	VideoPixelDepth:           "videoPixelDepth",
	VideoPixelAspectRatio:     "videoPixelAspectRatio",
	VideoFieldOrder:           "videoFieldOrder",
	TapeName:                  "tapeName",
	AltTapeName:               "altTapeName",
	About:                     "about",
	Action:                    "action",
	AlreadyApplied:            "AlreadyApplied",
	Alt:                       "Alt",
	ApertureValue:             "ApertureValue",
	ApproximateFocusDistance:  "ApproximateFocusDistance",
	Bag:                       "Bag",
	BodySerialNumber:          "BodySerialNumber",
	CameraOwnerName:           "CameraOwnerName",
	Changed:                   "changed",
	ColorMode:                 "ColorMode",
	ColorSpace:                "ColorSpace",
	ComponentsConfiguration:   "ComponentsConfiguration",
	Compression:               "Compression",
	Contrast:                  "Contrast",
	CreateDate:                "CreateDate",
	Creator:                   "creator",
	CreatorTool:               "CreatorTool",
	CustomRendered:            "CustomRendered",
	DateCreated:               "DateCreated",
	DateTimeDigitized:         "DateTimeDigitized",
	DateTimeOriginal:          "DateTimeOriginal",
	DerivedFrom:               "DerivedFrom",
	Description:               "Description",
	DigitalZoomRatio:          "DigitalZoomRatio",
	DocumentID:                "DocumentID",
	EmbeddedXMPDigest:         "EmbeddedXMPDigest",
	ExifVersion:               "ExifVersion",
	ExposureBiasValue:         "ExposureBiasValue",
	ExposureMode:              "ExposureMode",
	ExposureProgram:           "ExposureProgram",
	ExposureTime:              "ExposureTime",
	FileSource:                "FileSource",
	Fired:                     "Fired",
	Firmware:                  "Firmware",
	Flash:                     "Flash",
	FlashCompensation:         "FlashCompensation",
	FlashpixVersion:           "FlashpixVersion",
	FNumber:                   "FNumber",
	FocalLength:               "FocalLength",
	FocalPlaneResolutionUnit:  "FocalPlaneResolutionUnit",
	FocalPlaneXResolution:     "FocalPlaneXResolution",
	FocalPlaneYResolution:     "FocalPlaneYResolution",
	Format:                    "format",
	Function:                  "Function",
	GainControl:               "GainControl",
	GPSAltitude:               "GPSAltitude",
	GPSAltitudeRef:            "GPSAltitudeRef",
	GPSLatitude:               "GPSLatitude",
	GPSLongitude:              "GPSLongitude",
	GPSMapDatum:               "GPSMapDatum",
	GPSTimeStamp:              "GPSTimeStamp",
	GPSVersionID:              "GPSVersionID",
	HierarchicalSubject:       "hierarchicalSubject",
	History:                   "History",
	ICCProfile:                "ICCProfile",
	ImageLength:               "ImageLength",
	ImageNumber:               "ImageNumber",
	ImageWidth:                "ImageWidth",
	InstanceID:                "InstanceID",
	InteroperabilityIndex:     "InteroperabilityIndex",
	ISOSpeedRatings:           "ISOSpeedRatings",
	Label:                     "Label",
	Lang:                      "lang",
	LegacyIPTCDigest:          "LegacyIPTCDigest",
	Lens:                      "Lens",
	LensID:                    "LensID",
	LensInfo:                  "LensInfo",
	LensModel:                 "LensModel",
	LensSerialNumber:          "LensSerialNumber",
	Li:                        "li",
	LightSource:               "LightSource",
	Make:                      "Make",
	MaxApertureValue:          "MaxApertureValue",
	MetadataDate:              "MetadataDate",
	MeteringMode:              "MeteringMode",
	Mode:                      "Mode",
	Model:                     "Model",
	ModifyDate:                "ModifyDate",
	NativeDigest:              "NativeDigest",
	Orientation:               "Orientation",
	OriginalDocumentID:        "OriginalDocumentID",
	ParseType:                 "parseType",
	PhotometricInterpretation: "PhotometricInterpretation",
	PixelXDimension:           "PixelXDimension",
	PixelYDimension:           "PixelYDimension",
	PlanarConfiguration:       "PlanarConfiguration",
	Rating:                    "Rating",
	RawFileName:               "RawFileName",
	RDF:                       "RDF",
	RecommendedExposureIndex:  "RecommendedExposureIndex",
	RedEyeMode:                "RedEyeMode",
	ResolutionUnit:            "ResolutionUnit",
	Return:                    "Return",
	Rights:                    "rights",
	SamplesPerPixel:           "SamplesPerPixel",
	Saturation:                "Saturation",
	SceneCaptureType:          "SceneCaptureType",
	SceneType:                 "SceneType",
	SensitivityType:           "SensitivityType",
	Seq:                       "Seq",
	SerialNumber:              "SerialNumber",
	Sharpness:                 "Sharpness",
	ShutterSpeedValue:         "ShutterSpeedValue",
	SidecarForExtension:       "SidecarForExtension",
	SoftwareAgent:             "softwareAgent",
	Subject:                   "subject",
	SubjectDistance:           "SubjectDistance",
	Temperature:               "Temperature",
	Title:                     "Title",
	ToneCurve:                 "ToneCurve",
	ToneCurveBlue:             "ToneCurveBlue",
	ToneCurveGreen:            "ToneCurveGreen",
	ToneCurvePV2012:           "ToneCurvePV2012",
	ToneCurvePV2012Blue:       "ToneCurvePV2012Blue",
	ToneCurvePV2012Green:      "ToneCurvePV2012Green",
	ToneCurvePV2012Red:        "ToneCurvePV2012Red",
	ToneCurveRed:              "ToneCurveRed",
	UserComment:               "UserComment",
	When:                      "when",
	WhiteBalance:              "WhiteBalance",
	XmpMeta:                   "xmpmeta",
	XResolution:               "XResolution",
	YResolution:               "YResolution",
}

// mapStringName returns string's value as a Name
var mapStringName = map[string]Name{
	"videoFrameSize":            VideoFrameSize,
	"w":                         W,
	"h":                         H,
	"startTimecode":             StartTimecode,
	"timeValue":                 TimeValue,
	"altTimecode":               AltTimecode,
	"xmpDM":                     XmpDM,
	"stDim":                     StDim,
	"xap":                       Xap,
	"dc":                        Dc,
	"videoFrameRate":            VideoFrameRate,
	"videoPixelDepth":           VideoPixelDepth,
	"videoPixelAspectRatio":     VideoPixelAspectRatio,
	"videoFieldOrder":           VideoFieldOrder,
	"tapeName":                  TapeName,
	"altTapeName":               AltTapeName,
	"about":                     About,
	"action":                    Action,
	"AlreadyApplied":            AlreadyApplied,
	"Alt":                       Alt,
	"ApertureValue":             ApertureValue,
	"ApproximateFocusDistance":  ApproximateFocusDistance,
	"Bag":                       Bag,
	"BodySerialNumber":          BodySerialNumber,
	"CameraOwnerName":           CameraOwnerName,
	"changed":                   Changed,
	"ColorMode":                 ColorMode,
	"ColorSpace":                ColorSpace,
	"ComponentsConfiguration":   ComponentsConfiguration,
	"Compression":               Compression,
	"Contrast":                  Contrast,
	"CreateDate":                CreateDate,
	"creator":                   Creator,
	"CreatorTool":               CreatorTool,
	"CustomRendered":            CustomRendered,
	"DateCreated":               DateCreated,
	"DateTimeDigitized":         DateTimeDigitized,
	"DateTimeOriginal":          DateTimeOriginal,
	"DerivedFrom":               DerivedFrom,
	"Description":               Description,
	"description":               Description,
	"DigitalZoomRatio":          DigitalZoomRatio,
	"DocumentID":                DocumentID,
	"EmbeddedXMPDigest":         EmbeddedXMPDigest,
	"ExifVersion":               ExifVersion,
	"ExposureBiasValue":         ExposureBiasValue,
	"ExposureMode":              ExposureMode,
	"ExposureProgram":           ExposureProgram,
	"ExposureTime":              ExposureTime,
	"FileSource":                FileSource,
	"Fired":                     Fired,
	"Firmware":                  Firmware,
	"Flash":                     Flash,
	"FlashCompensation":         FlashCompensation,
	"FlashpixVersion":           FlashpixVersion,
	"FNumber":                   FNumber,
	"FocalLength":               FocalLength,
	"FocalPlaneResolutionUnit":  FocalPlaneResolutionUnit,
	"FocalPlaneXResolution":     FocalPlaneXResolution,
	"FocalPlaneYResolution":     FocalPlaneYResolution,
	"format":                    Format,
	"Function":                  Function,
	"GainControl":               GainControl,
	"GPSAltitude":               GPSAltitude,
	"GPSAltitudeRef":            GPSAltitudeRef,
	"GPSLatitude":               GPSLatitude,
	"GPSLongitude":              GPSLongitude,
	"GPSMapDatum":               GPSMapDatum,
	"GPSTimeStamp":              GPSTimeStamp,
	"GPSVersionID":              GPSVersionID,
	"hierarchicalSubject":       HierarchicalSubject,
	"History":                   History,
	"ICCProfile":                ICCProfile,
	"ImageLength":               ImageLength,
	"ImageNumber":               ImageNumber,
	"ImageWidth":                ImageWidth,
	"InstanceID":                InstanceID,
	"InteroperabilityIndex":     InteroperabilityIndex,
	"ISOSpeedRatings":           ISOSpeedRatings,
	"Label":                     Label,
	"lang":                      Lang,
	"LegacyIPTCDigest":          LegacyIPTCDigest,
	"Lens":                      Lens,
	"LensID":                    LensID,
	"LensInfo":                  LensInfo,
	"LensModel":                 LensModel,
	"LensSerialNumber":          LensSerialNumber,
	"li":                        Li,
	"LightSource":               LightSource,
	"Make":                      Make,
	"MaxApertureValue":          MaxApertureValue,
	"MetadataDate":              MetadataDate,
	"MeteringMode":              MeteringMode,
	"Mode":                      Mode,
	"Model":                     Model,
	"ModifyDate":                ModifyDate,
	"NativeDigest":              NativeDigest,
	"Orientation":               Orientation,
	"OriginalDocumentID":        OriginalDocumentID,
	"parseType":                 ParseType,
	"PhotometricInterpretation": PhotometricInterpretation,
	"PixelXDimension":           PixelXDimension,
	"PixelYDimension":           PixelYDimension,
	"PlanarConfiguration":       PlanarConfiguration,
	"Rating":                    Rating,
	"RawFileName":               RawFileName,
	"rdf":                       RDF,
	"RDF":                       RDF,
	"RecommendedExposureIndex":  RecommendedExposureIndex,
	"RedEyeMode":                RedEyeMode,
	"ResolutionUnit":            ResolutionUnit,
	"Return":                    Return,
	"rights":                    Rights,
	"SamplesPerPixel":           SamplesPerPixel,
	"Saturation":                Saturation,
	"SceneCaptureType":          SceneCaptureType,
	"SceneType":                 SceneType,
	"SensitivityType":           SensitivityType,
	"Seq":                       Seq,
	"SerialNumber":              SerialNumber,
	"Sharpness":                 Sharpness,
	"ShutterSpeedValue":         ShutterSpeedValue,
	"SidecarForExtension":       SidecarForExtension,
	"softwareAgent":             SoftwareAgent,
	"subject":                   Subject,
	"SubjectDistance":           SubjectDistance,
	"Temperature":               Temperature,
	"Title":                     Title,
	"title":                     Title,
	"ToneCurve":                 ToneCurve,
	"ToneCurveBlue":             ToneCurveBlue,
	"ToneCurveGreen":            ToneCurveGreen,
	"ToneCurvePV2012":           ToneCurvePV2012,
	"ToneCurvePV2012Blue":       ToneCurvePV2012Blue,
	"ToneCurvePV2012Green":      ToneCurvePV2012Green,
	"ToneCurvePV2012Red":        ToneCurvePV2012Red,
	"ToneCurveRed":              ToneCurveRed,
	"Unknown":                   UnknownPropertyName,
	"UserComment":               UserComment,
	"when":                      When,
	"WhiteBalance":              WhiteBalance,
	"xmpmeta":                   XmpMeta,
	"XResolution":               XResolution,
	"YResolution":               YResolution,
}
