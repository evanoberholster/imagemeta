package xmp

// Name is the XMP Property name
type Name uint8

func (n Name) String() string {
	if v, ok := mapNameString[n]; ok {
		return v
	}
	return mapNameString[UnknownPropertyName]
}

// IdentifyName returns the XMP Property Name correspondent to buf.
// If Property Name was not identified returns UnknownName.
func IdentifyName(buf []byte) Name {
	return identifyName(buf)
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
	DistortionCorrectionAlreadyApplied
	DerivedFrom
	DerivedFromDocumentID
	DerivedFromOriginalDocumentID
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
	FlashTag
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
	Good
	H
	HierarchicalSubject
	HueAdjustmentRed
	HueAdjustmentOrange
	HueAdjustmentYellow
	HueAdjustmentGreen
	HueAdjustmentAqua
	HueAdjustmentBlue
	HueAdjustmentPurple
	HueAdjustmentMagenta
	HistoryTag
	ICCProfile
	ImageDescription
	ImageLength
	ImageNumber
	ImageWidth
	InstanceID
	InteroperabilityIndex
	ISOSpeedRatings
	LateralChromaticAberrationCorrectionAlreadyApplied
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
	LuminanceAdjustmentRed
	LuminanceAdjustmentOrange
	LuminanceAdjustmentYellow
	LuminanceAdjustmentGreen
	LuminanceAdjustmentAqua
	LuminanceAdjustmentBlue
	LuminanceAdjustmentPurple
	LuminanceAdjustmentMagenta
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
	Parameters
	ParseType // parseType
	PhotometricInterpretation
	PhotographicSensitivity
	PixelXDimension
	PixelYDimension
	Pick
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
	SaturationAdjustmentRed
	SaturationAdjustmentOrange
	SaturationAdjustmentYellow
	SaturationAdjustmentGreen
	SaturationAdjustmentAqua
	SaturationAdjustmentBlue
	SaturationAdjustmentPurple
	SaturationAdjustmentMagenta
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
	SubsecTime
	SubsecTimeDigitized
	SubsecTimeOriginal
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
	VignetteCorrectionAlreadyApplied
	VideoFieldOrder
	VideoFrameRate
	VideoFrameSize
	VideoPixelAspectRatio
	VideoPixelDepth
	W
	When
	WeightedFlatSubject
	WhiteBalance
	Xap
	XMPToolkit
	XmpDM
	XmpMeta
	XResolution
	YCbCrPositioning
	YResolution
	GPSDOP
	GPSMeasureMode
	GPSSatellites
)

// mapNameString returns Name's value as a string
var mapNameString = map[Name]string{
	UnknownPropertyName:                "Unknown",
	About:                              "about",
	Action:                             "action",
	AlreadyApplied:                     "AlreadyApplied",
	Alt:                                "Alt",
	AltTapeName:                        "altTapeName",
	AltTimecode:                        "altTimecode",
	ApertureValue:                      "ApertureValue",
	ApproximateFocusDistance:           "ApproximateFocusDistance",
	Bag:                                "Bag",
	BitsPerSample:                      "BitsPerSample",
	BodySerialNumber:                   "BodySerialNumber",
	BrightnessValue:                    "BrightnessValue",
	CameraOwnerName:                    "CameraOwnerName",
	Changed:                            "changed",
	ColorMode:                          "ColorMode",
	ColorSpace:                         "ColorSpace",
	ComponentsConfiguration:            "ComponentsConfiguration",
	CompressedBitsPerPixel:             "CompressedBitsPerPixel",
	Compression:                        "Compression",
	Contrast:                           "Contrast",
	CreateDate:                         "CreateDate",
	Creator:                            "creator",
	CreatorTool:                        "CreatorTool",
	CustomRendered:                     "CustomRendered",
	DateCreated:                        "DateCreated",
	DateTimeDigitized:                  "DateTimeDigitized",
	DateTimeOriginal:                   "DateTimeOriginal",
	Dc:                                 "dc",
	DistortionCorrectionAlreadyApplied: "DistortionCorrectionAlreadyApplied",
	DerivedFrom:                        "DerivedFrom",
	DerivedFromDocumentID:              "DerivedFromDocumentID",
	DerivedFromOriginalDocumentID:      "DerivedFromOriginalDocumentID",
	Description:                        "Description",
	DigitalZoomRatio:                   "DigitalZoomRatio",
	DocumentID:                         "DocumentID",
	EmbeddedXMPDigest:                  "EmbeddedXMPDigest",
	ExifVersion:                        "ExifVersion",
	ExposureBiasValue:                  "ExposureBiasValue",
	ExposureMode:                       "ExposureMode",
	ExposureProgram:                    "ExposureProgram",
	ExposureTime:                       "ExposureTime",
	FileSource:                         "FileSource",
	Fired:                              "Fired",
	Firmware:                           "Firmware",
	FlashTag:                           "Flash",
	FlashCompensation:                  "FlashCompensation",
	FlashpixVersion:                    "FlashpixVersion",
	FNumber:                            "FNumber",
	FocalLength:                        "FocalLength",
	FocalLengthIn35mmFilm:              "FocalLengthIn35mmFilm",
	FocalPlaneResolutionUnit:           "FocalPlaneResolutionUnit",
	FocalPlaneXResolution:              "FocalPlaneXResolution",
	FocalPlaneYResolution:              "FocalPlaneYResolution",
	Format:                             "format",
	Function:                           "Function",
	GainControl:                        "GainControl",
	GPSAltitude:                        "GPSAltitude",
	GPSAltitudeRef:                     "GPSAltitudeRef",
	GPSDifferential:                    "GPSDifferential",
	GPSLatitude:                        "GPSLatitude",
	GPSLongitude:                       "GPSLongitude",
	GPSMapDatum:                        "GPSMapDatum",
	GPSStatus:                          "GPSStatus",
	GPSTimeStamp:                       "GPSTimeStamp",
	GPSVersionID:                       "GPSVersionID",
	Good:                               "good",
	H:                                  "h",
	HierarchicalSubject:                "hierarchicalSubject",
	HueAdjustmentRed:                   "HueAdjustmentRed",
	HueAdjustmentOrange:                "HueAdjustmentOrange",
	HueAdjustmentYellow:                "HueAdjustmentYellow",
	HueAdjustmentGreen:                 "HueAdjustmentGreen",
	HueAdjustmentAqua:                  "HueAdjustmentAqua",
	HueAdjustmentBlue:                  "HueAdjustmentBlue",
	HueAdjustmentPurple:                "HueAdjustmentPurple",
	HueAdjustmentMagenta:               "HueAdjustmentMagenta",
	HistoryTag:                         "History",
	ICCProfile:                         "ICCProfile",
	ImageDescription:                   "ImageDescription",
	ImageLength:                        "ImageLength",
	ImageNumber:                        "ImageNumber",
	ImageWidth:                         "ImageWidth",
	InstanceID:                         "InstanceID",
	InteroperabilityIndex:              "InteroperabilityIndex",
	ISOSpeedRatings:                    "ISOSpeedRatings",
	LateralChromaticAberrationCorrectionAlreadyApplied: "LateralChromaticAberrationCorrectionAlreadyApplied",
	Label:                            "Label",
	Lang:                             "lang",
	LegacyIPTCDigest:                 "LegacyIPTCDigest",
	Lens:                             "Lens",
	LensID:                           "LensID",
	LensInfo:                         "LensInfo",
	LensModel:                        "LensModel",
	LensSerialNumber:                 "LensSerialNumber",
	Li:                               "li",
	LightSource:                      "LightSource",
	LuminanceAdjustmentRed:           "LuminanceAdjustmentRed",
	LuminanceAdjustmentOrange:        "LuminanceAdjustmentOrange",
	LuminanceAdjustmentYellow:        "LuminanceAdjustmentYellow",
	LuminanceAdjustmentGreen:         "LuminanceAdjustmentGreen",
	LuminanceAdjustmentAqua:          "LuminanceAdjustmentAqua",
	LuminanceAdjustmentBlue:          "LuminanceAdjustmentBlue",
	LuminanceAdjustmentPurple:        "LuminanceAdjustmentPurple",
	LuminanceAdjustmentMagenta:       "LuminanceAdjustmentMagenta",
	Make:                             "Make",
	MaxApertureValue:                 "MaxApertureValue",
	MetadataDate:                     "MetadataDate",
	MeteringMode:                     "MeteringMode",
	Mode:                             "Mode",
	Model:                            "Model",
	ModifyDate:                       "ModifyDate",
	NativeDigest:                     "NativeDigest",
	Orientation:                      "Orientation",
	OriginalDocumentID:               "OriginalDocumentID",
	Parameters:                       "parameters",
	ParseType:                        "parseType",
	PhotometricInterpretation:        "PhotometricInterpretation",
	PhotographicSensitivity:          "PhotographicSensitivity",
	PixelXDimension:                  "PixelXDimension",
	PixelYDimension:                  "PixelYDimension",
	Pick:                             "pick",
	PlanarConfiguration:              "PlanarConfiguration",
	PreservedFileName:                "PreservedFileName",
	Rating:                           "Rating",
	RawFileName:                      "RawFileName",
	RDF:                              "RDF",
	RecommendedExposureIndex:         "RecommendedExposureIndex",
	RedEyeMode:                       "RedEyeMode",
	ResolutionUnit:                   "ResolutionUnit",
	Return:                           "Return",
	Rights:                           "rights",
	SamplesPerPixel:                  "SamplesPerPixel",
	Saturation:                       "Saturation",
	SaturationAdjustmentRed:          "SaturationAdjustmentRed",
	SaturationAdjustmentOrange:       "SaturationAdjustmentOrange",
	SaturationAdjustmentYellow:       "SaturationAdjustmentYellow",
	SaturationAdjustmentGreen:        "SaturationAdjustmentGreen",
	SaturationAdjustmentAqua:         "SaturationAdjustmentAqua",
	SaturationAdjustmentBlue:         "SaturationAdjustmentBlue",
	SaturationAdjustmentPurple:       "SaturationAdjustmentPurple",
	SaturationAdjustmentMagenta:      "SaturationAdjustmentMagenta",
	SceneCaptureType:                 "SceneCaptureType",
	SceneType:                        "SceneType",
	SensitivityType:                  "SensitivityType",
	Seq:                              "Seq",
	SerialNumber:                     "SerialNumber",
	Sharpness:                        "Sharpness",
	ShutterSpeedValue:                "ShutterSpeedValue",
	SidecarForExtension:              "SidecarForExtension",
	Software:                         "Software",
	SoftwareAgent:                    "softwareAgent",
	StartTimecode:                    "startTimecode",
	StDim:                            "stDim",
	Subject:                          "subject",
	SubjectDistance:                  "SubjectDistance",
	SubsecTime:                       "SubsecTime",
	SubsecTimeDigitized:              "SubsecTimeDigitized",
	SubsecTimeOriginal:               "SubsecTimeOriginal",
	TapeName:                         "tapeName",
	Temperature:                      "Temperature",
	TimeValue:                        "timeValue",
	Title:                            "Title",
	ToneCurve:                        "ToneCurve",
	ToneCurveBlue:                    "ToneCurveBlue",
	ToneCurveGreen:                   "ToneCurveGreen",
	ToneCurvePV2012:                  "ToneCurvePV2012",
	ToneCurvePV2012Blue:              "ToneCurvePV2012Blue",
	ToneCurvePV2012Green:             "ToneCurvePV2012Green",
	ToneCurvePV2012Red:               "ToneCurvePV2012Red",
	ToneCurveRed:                     "ToneCurveRed",
	UserComment:                      "UserComment",
	VignetteCorrectionAlreadyApplied: "VignetteCorrectionAlreadyApplied",
	VideoFieldOrder:                  "videoFieldOrder",
	VideoFrameRate:                   "videoFrameRate",
	VideoFrameSize:                   "videoFrameSize",
	VideoPixelAspectRatio:            "videoPixelAspectRatio",
	VideoPixelDepth:                  "videoPixelDepth",
	W:                                "w",
	When:                             "when",
	WeightedFlatSubject:              "weightedFlatSubject",
	WhiteBalance:                     "WhiteBalance",
	Xap:                              "xap",
	XMPToolkit:                       "xmptk",
	XmpDM:                            "xmpDM",
	XmpMeta:                          "xmpmeta",
	XResolution:                      "XResolution",
	YCbCrPositioning:                 "YCbCrPositioning",
	YResolution:                      "YResolution",
	GPSDOP:                           "GPSDOP",
	GPSMeasureMode:                   "GPSMeasureMode",
	GPSSatellites:                    "GPSSatellites",
}
