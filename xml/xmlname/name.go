package xmlname

import "fmt"

// Name is an interface representing (NS) Namespace or TagName
type Name interface {
	String() string
}

// TagName is an XML Name
type TagName uint8

func (n TagName) String() string {
	return fmt.Sprintf(mapTagNameString[n])
}

// IdentifyTagName returns the (TagName) XML Tag Name correspondent to buf.
// If Tag Name was not identified returns UnknownName.
func IdentifyTagName(buf []byte) (n TagName) {
	return mapStringTagName[string(buf)]
}

// Names
const (
	UnknownTagName TagName = iota
	About                  // about
	Action
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
	Contrast
	CreateDate
	Creator
	CustomRendered
	DateCreated
	DateTimeDigitized
	DateTimeOriginal
	Description
	DocumentID
	EmbeddedXMPDigest
	ExifVersion
	ExposureBiasValue
	ExposureMode
	ExposureProgram
	ExposureTime
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
	Lang
	LegacyIPTCDigest
	Lens
	LensID
	LensInfo
	LensModel
	LensSerialNumber
	Li
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
	PixelXDimension
	PixelYDimension
	Rating
	RawFileName
	RDF
	RecommendedExposureIndex
	RedEyeMode
	ResolutionUnit
	Return
	Rights
	SceneCaptureType
	SensitivityType
	Seq
	SerialNumber
	ShutterSpeedValue
	SidecarForExtension
	SoftwareAgent
	Subject
	SubjectDistance
	Temperature
	ToneCurve
	UserComment
	When
	WhiteBalance
	XmpMeta
	XResolution
	YResolution
)

// mapTagNameString returns Name's value as a string
var mapTagNameString = map[TagName]string{
	UnknownTagName:           "Unknown",
	About:                    "about",
	Action:                   "action",
	Alt:                      "Alt",
	ApertureValue:            "ApertureValue",
	ApproximateFocusDistance: "ApproximateFocusDistance",
	Bag:                      "Bag",
	BodySerialNumber:         "BodySerialNumber",
	CameraOwnerName:          "CameraOwnerName",
	Changed:                  "changed",
	ColorMode:                "ColorMode",
	ColorSpace:               "ColorSpace",
	ComponentsConfiguration:  "ComponentsConfiguration",
	Contrast:                 "Contrast",
	CreateDate:               "CreateDate",
	Creator:                  "creator",
	CustomRendered:           "CustomRendered",
	DateCreated:              "DateCreated",
	DateTimeDigitized:        "DateTimeDigitized",
	DateTimeOriginal:         "DateTimeOriginal",
	Description:              "Description",
	DocumentID:               "DocumentID",
	EmbeddedXMPDigest:        "EmbeddedXMPDigest",
	ExifVersion:              "ExifVersion",
	ExposureBiasValue:        "ExposureBiasValue",
	ExposureMode:             "ExposureMode",
	ExposureProgram:          "ExposureProgram",
	ExposureTime:             "ExposureTime",
	Fired:                    "Fired",
	Firmware:                 "Firmware",
	Flash:                    "Flash",
	FlashCompensation:        "FlashCompensation",
	FlashpixVersion:          "FlashpixVersion",
	FNumber:                  "FNumber",
	FocalLength:              "FocalLength",
	FocalPlaneResolutionUnit: "FocalPlaneResolutionUnit",
	FocalPlaneXResolution:    "FocalPlaneXResolution",
	FocalPlaneYResolution:    "FocalPlaneYResolution",
	Format:                   "format",
	Function:                 "Function",
	GPSAltitude:              "GPSAltitude",
	GPSAltitudeRef:           "GPSAltitudeRef",
	GPSLatitude:              "GPSLatitude",
	GPSLongitude:             "GPSLongitude",
	GPSMapDatum:              "GPSMapDatum",
	GPSTimeStamp:             "GPSTimeStamp",
	GPSVersionID:             "GPSVersionID",
	HierarchicalSubject:      "hierarchicalSubject",
	History:                  "History",
	ICCProfile:               "ICCProfile",
	ImageLength:              "ImageLength",
	ImageNumber:              "ImageNumber",
	ImageWidth:               "ImageWidth",
	InstanceID:               "InstanceID",
	InteroperabilityIndex:    "InteroperabilityIndex",
	ISOSpeedRatings:          "ISOSpeedRatings",
	Lang:                     "lang",
	LegacyIPTCDigest:         "LegacyIPTCDigest",
	Lens:                     "Lens",
	LensID:                   "LensID",
	LensInfo:                 "LensInfo",
	LensModel:                "LensModel",
	LensSerialNumber:         "LensSerialNumber",
	Li:                       "li",
	Make:                     "Make",
	MaxApertureValue:         "MaxApertureValue",
	MetadataDate:             "MetadataDate",
	MeteringMode:             "MeteringMode",
	Mode:                     "Mode",
	Model:                    "Model",
	ModifyDate:               "ModifyDate",
	NativeDigest:             "NativeDigest",
	Orientation:              "Orientation",
	OriginalDocumentID:       "OriginalDocumentID",
	ParseType:                "parseType",
	PixelXDimension:          "PixelXDimension",
	PixelYDimension:          "PixelYDimension",
	Rating:                   "Rating",
	RawFileName:              "RawFileName",
	RDF:                      "RDF",
	RecommendedExposureIndex: "RecommendedExposureIndex",
	RedEyeMode:               "RedEyeMode",
	ResolutionUnit:           "ResolutionUnit",
	Return:                   "Return",
	Rights:                   "rights",
	SceneCaptureType:         "SceneCaptureType",
	SensitivityType:          "SensitivityType",
	Seq:                      "Seq",
	SerialNumber:             "SerialNumber",
	ShutterSpeedValue:        "ShutterSpeedValue",
	SidecarForExtension:      "SidecarForExtension",
	SoftwareAgent:            "softwareAgent",
	Subject:                  "subject",
	SubjectDistance:          "SubjectDistance",
	Temperature:              "Temperature",
	ToneCurve:                "ToneCurve",
	UserComment:              "UserComment",
	When:                     "when",
	WhiteBalance:             "WhiteBalance",
	XmpMeta:                  "xmpmeta",
	XResolution:              "XResolution",
	YResolution:              "YResolution",
}

// mapStringTagName returns string's value as a Name
var mapStringTagName = map[string]TagName{
	"about":                    About,
	"action":                   Action,
	"Alt":                      Alt,
	"ApertureValue":            ApertureValue,
	"ApproximateFocusDistance": ApproximateFocusDistance,
	"Bag":                      Bag,
	"BodySerialNumber":         BodySerialNumber,
	"CameraOwnerName":          CameraOwnerName,
	"changed":                  Changed,
	"ColorMode":                ColorMode,
	"ColorSpace":               ColorSpace,
	"ComponentsConfiguration":  ComponentsConfiguration,
	"Contrast":                 Contrast,
	"CreateDate":               CreateDate,
	"creator":                  Creator,
	"CustomRendered":           CustomRendered,
	"DateCreated":              DateCreated,
	"DateTimeDigitized":        DateTimeDigitized,
	"DateTimeOriginal":         DateTimeOriginal,
	"Description":              Description,
	"DocumentID":               DocumentID,
	"EmbeddedXMPDigest":        EmbeddedXMPDigest,
	"ExifVersion":              ExifVersion,
	"ExposureBiasValue":        ExposureBiasValue,
	"ExposureMode":             ExposureMode,
	"ExposureProgram":          ExposureProgram,
	"ExposureTime":             ExposureTime,
	"Fired":                    Fired,
	"Firmware":                 Firmware,
	"Flash":                    Flash,
	"FlashCompensation":        FlashCompensation,
	"FlashpixVersion":          FlashpixVersion,
	"FNumber":                  FNumber,
	"FocalLength":              FocalLength,
	"FocalPlaneResolutionUnit": FocalPlaneResolutionUnit,
	"FocalPlaneXResolution":    FocalPlaneXResolution,
	"FocalPlaneYResolution":    FocalPlaneYResolution,
	"format":                   Format,
	"Function":                 Function,
	"GPSAltitude":              GPSAltitude,
	"GPSAltitudeRef":           GPSAltitudeRef,
	"GPSLatitude":              GPSLatitude,
	"GPSLongitude":             GPSLongitude,
	"GPSMapDatum":              GPSMapDatum,
	"GPSTimeStamp":             GPSTimeStamp,
	"GPSVersionID":             GPSVersionID,
	"hierarchicalSubject":      HierarchicalSubject,
	"History":                  History,
	"ICCProfile":               ICCProfile,
	"ImageLength":              ImageLength,
	"ImageNumber":              ImageNumber,
	"ImageWidth":               ImageWidth,
	"InstanceID":               InstanceID,
	"InteroperabilityIndex":    InteroperabilityIndex,
	"ISOSpeedRatings":          ISOSpeedRatings,
	"lang":                     Lang,
	"LegacyIPTCDigest":         LegacyIPTCDigest,
	"Lens":                     Lens,
	"LensID":                   LensID,
	"LensInfo":                 LensInfo,
	"LensModel":                LensModel,
	"LensSerialNumber":         LensSerialNumber,
	"li":                       Li,
	"Make":                     Make,
	"MaxApertureValue":         MaxApertureValue,
	"MetadataDate":             MetadataDate,
	"MeteringMode":             MeteringMode,
	"Mode":                     Mode,
	"Model":                    Model,
	"ModifyDate":               ModifyDate,
	"NativeDigest":             NativeDigest,
	"Orientation":              Orientation,
	"OriginalDocumentID":       OriginalDocumentID,
	"parseType":                ParseType,
	"PixelXDimension":          PixelXDimension,
	"PixelYDimension":          PixelYDimension,
	"Rating":                   Rating,
	"RawFileName":              RawFileName,
	"RDF":                      RDF,
	"RecommendedExposureIndex": RecommendedExposureIndex,
	"RedEyeMode":               RedEyeMode,
	"ResolutionUnit":           ResolutionUnit,
	"Return":                   Return,
	"rights":                   Rights,
	"SceneCaptureType":         SceneCaptureType,
	"SensitivityType":          SensitivityType,
	"Seq":                      Seq,
	"SerialNumber":             SerialNumber,
	"ShutterSpeedValue":        ShutterSpeedValue,
	"SidecarForExtension":      SidecarForExtension,
	"softwareAgent":            SoftwareAgent,
	"subject":                  Subject,
	"SubjectDistance":          SubjectDistance,
	"Temperature":              Temperature,
	"ToneCurve":                ToneCurve,
	"Unknown":                  UnknownTagName,
	"UserComment":              UserComment,
	"when":                     When,
	"WhiteBalance":             WhiteBalance,
	"xmpmeta":                  XmpMeta,
	"XResolution":              XResolution,
	"YResolution":              YResolution,
}
