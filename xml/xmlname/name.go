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
	RDF
	Description
	ModifyDate
	CreateDate
	MetadataDate
	Make
	Model
	Orientation
	ImageWidth
	ImageLength
	ExifVersion
	ExposureTime
	ShutterSpeedValue
	FNumber
	ApertureValue
	ExposureProgram
	SensitivityType
	RecommendedExposureIndex
	ExposureBiasValue
	MaxApertureValue
	MeteringMode
	FocalLength
	CustomRendered
	ExposureMode
	WhiteBalance
	SceneCaptureType
	FocalPlaneXResolution
	FocalPlaneYResolution
	FocalPlaneResolutionUnit
	DateTimeOriginal
	PixelXDimension
	PixelYDimension
	Format // format
	SerialNumber
	LensInfo
	Lens
	LensID
	LensSerialNumber
	ImageNumber
	ApproximateFocusDistance
	FlashCompensation
	Firmware
	LensModel
	DateCreated
	SidecarForExtension
	EmbeddedXMPDigest
	DocumentID
	OriginalDocumentID
	InstanceID
	RawFileName

	ISOSpeedRatings
	Seq
	Li
	Flash
	Fired
	Return
	Mode
	Function
	RedEyeMode
	Creator
	Rights
	Alt
	Lang
)

// mapTagNameString returns Name's value as a string
var mapTagNameString = map[TagName]string{
	UnknownTagName:           "Unknown",
	RDF:                      "RDF",
	Description:              "Description",
	ModifyDate:               "ModifyDate",
	CreateDate:               "CreateDate",
	MetadataDate:             "MetadataDate",
	Make:                     "Make",
	Model:                    "Model",
	Orientation:              "Orientation",
	ImageWidth:               "ImageWidth",
	ImageLength:              "ImageLength",
	ExifVersion:              "ExifVersion",
	ExposureTime:             "ExposureTime",
	ShutterSpeedValue:        "ShutterSpeedValue",
	FNumber:                  "FNumber",
	ApertureValue:            "ApertureValue",
	ExposureProgram:          "ExposureProgram",
	SensitivityType:          "SensitivityType",
	RecommendedExposureIndex: "RecommendedExposureIndex",
	ExposureBiasValue:        "ExposureBiasValue",
	MaxApertureValue:         "MaxApertureValue",
	MeteringMode:             "MeteringMode",
	FocalLength:              "FocalLength",
	CustomRendered:           "CustomRendered",
	ExposureMode:             "ExposureMode",
	WhiteBalance:             "WhiteBalance",
	SceneCaptureType:         "SceneCaptureType",
	FocalPlaneXResolution:    "FocalPlaneXResolution",
	FocalPlaneYResolution:    "FocalPlaneYResolution",
	FocalPlaneResolutionUnit: "FocalPlaneResolutionUnit",
	DateTimeOriginal:         "DateTimeOriginal",
	PixelXDimension:          "PixelXDimension",
	PixelYDimension:          "PixelYDimension",
	Format:                   "format",
	SerialNumber:             "SerialNumber",
	LensInfo:                 "LensInfo",
	Lens:                     "Lens",
	LensID:                   "LensID",
	LensSerialNumber:         "LensSerialNumber",
	ImageNumber:              "ImageNumber",
	ApproximateFocusDistance: "ApproximateFocusDistance",
	FlashCompensation:        "FlashCompensation",
	Firmware:                 "Firmware",
	LensModel:                "LensModel",
	DateCreated:              "DateCreated",
	SidecarForExtension:      "SidecarForExtension",
	EmbeddedXMPDigest:        "EmbeddedXMPDigest",
	DocumentID:               "DocumentID",
	OriginalDocumentID:       "OriginalDocumentID",
	InstanceID:               "InstanceID",
	RawFileName:              "RawFileName",
	ISOSpeedRatings:          "ISOSpeedRatings",
	Seq:                      "Seq",
	Li:                       "li",
	Flash:                    "Flash",
	Fired:                    "Fired",
	Return:                   "Return",
	Mode:                     "Mode",
	Function:                 "Function",
	RedEyeMode:               "RedEyeMode",
	Creator:                  "creator",
	Rights:                   "rights",
	Alt:                      "Alt",
	Lang:                     "lang",
}

// mapStringTagName returns string's value as a Name
var mapStringTagName = map[string]TagName{
	"Unknown":                  UnknownTagName,
	"RDF":                      RDF,
	"Description":              Description,
	"ModifyDate":               ModifyDate,
	"CreateDate":               CreateDate,
	"MetadataDate":             MetadataDate,
	"Make":                     Make,
	"Model":                    Model,
	"Orientation":              Orientation,
	"ImageWidth":               ImageWidth,
	"ImageLength":              ImageLength,
	"ExifVersion":              ExifVersion,
	"ExposureTime":             ExposureTime,
	"ShutterSpeedValue":        ShutterSpeedValue,
	"FNumber":                  FNumber,
	"ApertureValue":            ApertureValue,
	"ExposureProgram":          ExposureProgram,
	"SensitivityType":          SensitivityType,
	"RecommendedExposureIndex": RecommendedExposureIndex,
	"ExposureBiasValue":        ExposureBiasValue,
	"MaxApertureValue":         MaxApertureValue,
	"MeteringMode":             MeteringMode,
	"FocalLength":              FocalLength,
	"CustomRendered":           CustomRendered,
	"ExposureMode":             ExposureMode,
	"WhiteBalance":             WhiteBalance,
	"SceneCaptureType":         SceneCaptureType,
	"FocalPlaneXResolution":    FocalPlaneXResolution,
	"FocalPlaneYResolution":    FocalPlaneYResolution,
	"FocalPlaneResolutionUnit": FocalPlaneResolutionUnit,
	"DateTimeOriginal":         DateTimeOriginal,
	"PixelXDimension":          PixelXDimension,
	"PixelYDimension":          PixelYDimension,
	"format":                   Format,
	"SerialNumber":             SerialNumber,
	"LensInfo":                 LensInfo,
	"Lens":                     Lens,
	"LensID":                   LensID,
	"LensSerialNumber":         LensSerialNumber,
	"ImageNumber":              ImageNumber,
	"ApproximateFocusDistance": ApproximateFocusDistance,
	"FlashCompensation":        FlashCompensation,
	"Firmware":                 Firmware,
	"LensModel":                LensModel,
	"DateCreated":              DateCreated,
	"SidecarForExtension":      SidecarForExtension,
	"EmbeddedXMPDigest":        EmbeddedXMPDigest,
	"DocumentID":               DocumentID,
	"OriginalDocumentID":       OriginalDocumentID,
	"InstanceID":               InstanceID,
	"RawFileName":              RawFileName,
	"ISOSpeedRatings":          ISOSpeedRatings,
	"Seq":                      Seq,
	"li":                       Li,
	"Flash":                    Flash,
	"Fired":                    Fired,
	"Return":                   Return,
	"Mode":                     Mode,
	"Function":                 Function,
	"RedEyeMode":               RedEyeMode,
	"creator":                  Creator,
	"rights":                   Rights,
	"Alt":                      Alt,
	"lang":                     Lang,
}
