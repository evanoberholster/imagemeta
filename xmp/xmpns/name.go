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

	About // about
	Action
	AlreadyApplied
	Alt
	AltTapeName
	AltTimecode
	ApertureValue
	ApproximateFocusDistance
	Bag
	BitsPerSample
	BodySerialNumber
	BrightnessValue
	CameraOwnerName
	Changed
	ColorMode
	ColorSpace
	ComponentsConfiguration
	CompressedBitsPerPixel
	Compression
	Contrast
	CreateDate
	Creator
	CreatorTool
	CustomRendered
	DateCreated
	DateTimeDigitized
	DateTimeOriginal
	Dc
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
	FocalLengthIn35mmFilm
	FocalPlaneResolutionUnit
	FocalPlaneXResolution
	FocalPlaneYResolution
	Format // format
	Function
	GainControl
	GPSAltitude
	GPSAltitudeRef
	GPSDifferential
	GPSLatitude
	GPSLongitude
	GPSMapDatum
	GPSStatus
	GPSTimeStamp
	GPSVersionID
	H
	HierarchicalSubject
	History
	ICCProfile
	ImageDescription
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
	PreservedFileName
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
	Software
	SoftwareAgent
	StartTimecode
	StDim
	Subject
	SubjectDistance
	TapeName
	Temperature
	TimeValue
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
	VideoFieldOrder
	VideoFrameRate
	VideoFrameSize
	VideoPixelAspectRatio
	VideoPixelDepth
	W
	When
	WhiteBalance
	Xap
	XmpDM
	XmpMeta
	XResolution
	YCbCrPositioning
	YResolution
)

// mapNameString returns Name's value as a string
var mapNameString = map[Name]string{
	UnknownPropertyName:       "Unknown",
	About:                     "about",
	Action:                    "action",
	AlreadyApplied:            "AlreadyApplied",
	Alt:                       "Alt",
	AltTapeName:               "altTapeName",
	AltTimecode:               "altTimecode",
	ApertureValue:             "ApertureValue",
	ApproximateFocusDistance:  "ApproximateFocusDistance",
	Bag:                       "Bag",
	BitsPerSample:             "BitsPerSample",
	BodySerialNumber:          "BodySerialNumber",
	BrightnessValue:           "BrightnessValue",
	CameraOwnerName:           "CameraOwnerName",
	Changed:                   "changed",
	ColorMode:                 "ColorMode",
	ColorSpace:                "ColorSpace",
	ComponentsConfiguration:   "ComponentsConfiguration",
	CompressedBitsPerPixel:    "CompressedBitsPerPixel",
	Compression:               "Compression",
	Contrast:                  "Contrast",
	CreateDate:                "CreateDate",
	Creator:                   "creator",
	CreatorTool:               "CreatorTool",
	CustomRendered:            "CustomRendered",
	DateCreated:               "DateCreated",
	DateTimeDigitized:         "DateTimeDigitized",
	DateTimeOriginal:          "DateTimeOriginal",
	Dc:                        "dc",
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
	FocalLengthIn35mmFilm:     "FocalLengthIn35mmFilm",
	FocalPlaneResolutionUnit:  "FocalPlaneResolutionUnit",
	FocalPlaneXResolution:     "FocalPlaneXResolution",
	FocalPlaneYResolution:     "FocalPlaneYResolution",
	Format:                    "format",
	Function:                  "Function",
	GainControl:               "GainControl",
	GPSAltitude:               "GPSAltitude",
	GPSAltitudeRef:            "GPSAltitudeRef",
	GPSDifferential:           "GPSDifferential",
	GPSLatitude:               "GPSLatitude",
	GPSLongitude:              "GPSLongitude",
	GPSMapDatum:               "GPSMapDatum",
	GPSStatus:                 "GPSStatus",
	GPSTimeStamp:              "GPSTimeStamp",
	GPSVersionID:              "GPSVersionID",
	H:                         "h",
	HierarchicalSubject:       "hierarchicalSubject",
	History:                   "History",
	ICCProfile:                "ICCProfile",
	ImageDescription:          "ImageDescription",
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
	PreservedFileName:         "PreservedFileName",
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
	Software:                  "Software",
	SoftwareAgent:             "softwareAgent",
	StartTimecode:             "startTimecode",
	StDim:                     "stDim",
	Subject:                   "subject",
	SubjectDistance:           "SubjectDistance",
	TapeName:                  "tapeName",
	Temperature:               "Temperature",
	TimeValue:                 "timeValue",
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
	VideoFieldOrder:           "videoFieldOrder",
	VideoFrameRate:            "videoFrameRate",
	VideoFrameSize:            "videoFrameSize",
	VideoPixelAspectRatio:     "videoPixelAspectRatio",
	VideoPixelDepth:           "videoPixelDepth",
	W:                         "w",
	When:                      "when",
	WhiteBalance:              "WhiteBalance",
	Xap:                       "xap",
	XmpDM:                     "xmpDM",
	XmpMeta:                   "xmpmeta",
	XResolution:               "XResolution",
	YCbCrPositioning:          "YCbCrPositioning",
	YResolution:               "YResolution",
}

// mapStringName returns string's value as a Name
var mapStringName = map[string]Name{
	"about":                     About,
	"action":                    Action,
	"AlreadyApplied":            AlreadyApplied,
	"Alt":                       Alt,
	"altTapeName":               AltTapeName,
	"altTimecode":               AltTimecode,
	"ApertureValue":             ApertureValue,
	"ApproximateFocusDistance":  ApproximateFocusDistance,
	"Bag":                       Bag,
	"BitsPerSample":             BitsPerSample,
	"BodySerialNumber":          BodySerialNumber,
	"BrightnessValue":           BrightnessValue,
	"CameraOwnerName":           CameraOwnerName,
	"changed":                   Changed,
	"ColorMode":                 ColorMode,
	"ColorSpace":                ColorSpace,
	"ComponentsConfiguration":   ComponentsConfiguration,
	"CompressedBitsPerPixel":    CompressedBitsPerPixel,
	"Compression":               Compression,
	"Contrast":                  Contrast,
	"CreateDate":                CreateDate,
	"creator":                   Creator,
	"CreatorTool":               CreatorTool,
	"CustomRendered":            CustomRendered,
	"DateCreated":               DateCreated,
	"DateTimeDigitized":         DateTimeDigitized,
	"DateTimeOriginal":          DateTimeOriginal,
	"dc":                        Dc,
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
	"FocalLengthIn35mmFilm":     FocalLengthIn35mmFilm,
	"FocalPlaneResolutionUnit":  FocalPlaneResolutionUnit,
	"FocalPlaneXResolution":     FocalPlaneXResolution,
	"FocalPlaneYResolution":     FocalPlaneYResolution,
	"format":                    Format,
	"Function":                  Function,
	"GainControl":               GainControl,
	"GPSAltitude":               GPSAltitude,
	"GPSAltitudeRef":            GPSAltitudeRef,
	"GPSDifferential":           GPSDifferential,
	"GPSLatitude":               GPSLatitude,
	"GPSLongitude":              GPSLongitude,
	"GPSMapDatum":               GPSMapDatum,
	"GPSStatus":                 GPSStatus,
	"GPSTimeStamp":              GPSTimeStamp,
	"GPSVersionID":              GPSVersionID,
	"h":                         H,
	"hierarchicalSubject":       HierarchicalSubject,
	"History":                   History,
	"ICCProfile":                ICCProfile,
	"ImageDescription":          ImageDescription,
	"ImageLength":               ImageLength,
	"ImageNumber":               ImageNumber,
	"ImageWidth":                ImageWidth,
	"instanceID":                InstanceID,
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
	"PreservedFileName":         PreservedFileName,
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
	"Software":                  Software,
	"softwareAgent":             SoftwareAgent,
	"startTimecode":             StartTimecode,
	"stDim":                     StDim,
	"subject":                   Subject,
	"SubjectDistance":           SubjectDistance,
	"tapeName":                  TapeName,
	"Temperature":               Temperature,
	"timeValue":                 TimeValue,
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
	"videoFieldOrder":           VideoFieldOrder,
	"videoFrameRate":            VideoFrameRate,
	"videoFrameSize":            VideoFrameSize,
	"videoPixelAspectRatio":     VideoPixelAspectRatio,
	"videoPixelDepth":           VideoPixelDepth,
	"w":                         W,
	"when":                      When,
	"WhiteBalance":              WhiteBalance,
	"xap":                       Xap,
	"xmpDM":                     XmpDM,
	"xmpmeta":                   XmpMeta,
	"XResolution":               XResolution,
	"YCbCrPositioning":          YCbCrPositioning,
	"YResolution":               YResolution,
}
