package xmlname

import "fmt"

// Name is an XML Name
type Name uint8

func (n Name) String() string {
	return fmt.Sprintf(MapNameString[n])
}

// Names
const (
	Unknown Name = iota
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
)

// MapNameString returns Name's value as a string
var MapNameString = map[Name]string{
	Unknown:                  "Unknown",
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
	Format:                   "Format",
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
}

// MapStringName returns string's value as a Name
var MapStringName = map[string]Name{
	"Unknown":                  Unknown,
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
	"Format":                   Format,
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
}
