package xmp

// Name is the local-name portion of an XMP property (for example "CreateDate").
type Name uint16

// TagName is retained as a compatibility alias for Name.
type TagName = Name

func (n Name) String() string {
	if v, ok := mapNameString[n]; ok {
		return v
	}
	return mapNameString[UnknownPropertyName]
}

// IdentifyName resolves an XML local-name token to its internal tag identifier.
// Unknown names resolve to UnknownPropertyName.
// The provided slice may be lowercased in-place as part of the fast lookup path.
func IdentifyName(buf []byte) Name {
	return identifyName(buf)
}

// IdentifyTagName is retained as a compatibility wrapper around IdentifyName.
func IdentifyTagName(buf []byte) Name {
	return IdentifyName(buf)
}

// Names
const (
	UnknownPropertyName Name = iota

	About // about
	Acceleration
	Action
	AlreadyApplied
	Alt
	AltTapeName
	AltTimecode
	AmbientTemperature
	ApertureValue
	ApproximateFocusDistance
	Artist
	Bag
	BitsPerSample
	BodySerialNumber
	BrightnessValue
	CFAPattern
	CFAPatternColumns
	CFAPatternRows
	CFAPatternValues
	CameraElevationAngle
	CameraFirmware
	CameraOwnerName
	Changed
	ColorMode
	ColorSpace
	CompImageImagesPerSequence
	CompImageMaxExposureAll
	CompImageMaxExposureUsed
	CompImageMinExposureAll
	CompImageMinExposureUsed
	CompImageNumSequences
	CompImageSumExposureAll
	CompImageSumExposureUsed
	CompImageTotalExposurePeriod
	CompImageValues
	CompositeImage
	CompositeImageCount
	CompositeImageExposureTimes
	ComponentsConfiguration
	CompressedBitsPerPixel
	Compression
	Copyright
	Contrast
	Contributor
	Coverage
	CreateDate
	Creator
	CreatorTool
	CustomRendered
	Date
	DateCreated
	DateTime
	DateTimeDigitized
	DateTimeOriginal
	Dc
	DerivedFrom
	DerivedFromDocumentID
	DerivedFromOriginalDocumentID
	Description
	DeviceSettingDescription
	DeviceSettingDescriptionColumns
	DeviceSettingDescriptionRows
	DeviceSettingDescriptionSettings
	DigitalZoomRatio
	DistortionCorrectionAlreadyApplied
	DocumentID
	EmbeddedXMPDigest
	ExifVersion
	ExposureBiasValue
	ExposureIndex
	ExposureMode
	ExposureProgram
	ExposureTime
	FNumber
	FileSource
	Fired
	Firmware
	FlashCompensation
	FlashEnergy
	FlashTag
	FlashpixVersion
	FocalLength
	FocalLengthIn35mmFilm
	FocalPlaneResolutionUnit
	FocalPlaneXResolution
	FocalPlaneYResolution
	Format // format
	Function
	GPSAltitude
	GPSAltitudeRef
	GPSAreaInformation
	GPSDOP
	GPSDestBearing
	GPSDestBearingRef
	GPSDestDistance
	GPSDestDistanceRef
	GPSDestLatitude
	GPSDestLongitude
	GPSDifferential
	GPSHPositioningError
	GPSImgDirection
	GPSImgDirectionRef
	GPSLatitude
	GPSLongitude
	GPSMapDatum
	GPSMeasureMode
	GPSProcessingMethod
	GPSSatellites
	GPSSpeed
	GPSSpeedRef
	GPSStatus
	GPSTimeStamp
	GPSTrack
	GPSTrackRef
	GPSVersionID
	Gamma
	GainControl
	Good
	H
	HierarchicalSubject
	HistoryTag
	Humidity
	HueAdjustmentAqua
	HueAdjustmentBlue
	HueAdjustmentGreen
	HueAdjustmentMagenta
	HueAdjustmentOrange
	HueAdjustmentPurple
	HueAdjustmentRed
	HueAdjustmentYellow
	ICCProfile
	ISOSpeed
	ISOSpeedLatitudeyyy
	ISOSpeedLatitudezzz
	ISOSpeedRatings
	Identifier
	ImageEditingSoftware
	ImageEditor
	ImageDescription
	ImageLength
	ImageNumber
	ImageTitle
	ImageUniqueID
	ImageWidth
	InstanceID
	InteroperabilityIndex
	Label
	Lang
	Language
	LateralChromaticAberrationCorrectionAlreadyApplied
	LegacyIPTCDigest
	Lens
	LensID
	LensInfo
	LensMake
	LensModel
	LensSerialNumber
	Li
	LightSource
	LuminanceAdjustmentAqua
	LuminanceAdjustmentBlue
	LuminanceAdjustmentGreen
	LuminanceAdjustmentMagenta
	LuminanceAdjustmentOrange
	LuminanceAdjustmentPurple
	LuminanceAdjustmentRed
	LuminanceAdjustmentYellow
	Make
	MakerNote
	MaxApertureValue
	MetadataDate
	MetadataEditingSoftware
	MeteringMode
	Mode
	Model
	ModifyDate
	NativeDigest
	OECF
	OECFColumns
	OECFNames
	OECFRows
	OECFValues
	OwnerName
	Orientation
	OriginalDocumentID
	Parameters
	ParseType // parseType
	PhotographicSensitivity
	PhotometricInterpretation
	Pick
	PixelXDimension
	PixelYDimension
	Photographer
	PlanarConfiguration
	Pressure
	PrimaryChromaticities
	PreservedFileName
	Publisher
	RAWDevelopingSoftware
	RDF
	Rating
	RawFileName
	ReferenceBlackWhite
	RecommendedExposureIndex
	RedEyeMode
	RelatedSoundFile
	Relation
	ResolutionUnit
	Return
	Rights
	SamplesPerPixel
	Saturation
	SaturationAdjustmentAqua
	SaturationAdjustmentBlue
	SaturationAdjustmentGreen
	SaturationAdjustmentMagenta
	SaturationAdjustmentOrange
	SaturationAdjustmentPurple
	SaturationAdjustmentRed
	SaturationAdjustmentYellow
	SceneCaptureType
	SceneType
	SensingMethod
	SensitivityType
	Seq
	SerialNumber
	Sharpness
	ShutterSpeedValue
	SidecarForExtension
	Software
	SoftwareAgent
	Source
	SpatialFrequencyResponse
	SpatialFrequencyResponseColumns
	SpatialFrequencyResponseNames
	SpatialFrequencyResponseRows
	SpatialFrequencyResponseValues
	SpectralSensitivity
	StandardOutputSensitivity
	StDim
	StartTimecode
	Subject
	SubjectArea
	SubjectDistance
	SubjectDistanceRange
	SubjectLocation
	SubsecTime
	SubsecTimeDigitized
	SubsecTimeOriginal
	TapeName
	Temperature
	TimeValue
	Title
	TransferFunction
	ToneCurve
	ToneCurveBlue
	ToneCurveGreen
	ToneCurvePV2012
	ToneCurvePV2012Blue
	ToneCurvePV2012Green
	ToneCurvePV2012Red
	ToneCurveRed
	Type
	UserComment
	VideoFieldOrder
	VideoFrameRate
	VideoFrameSize
	VideoPixelAspectRatio
	VideoPixelDepth
	VignetteCorrectionAlreadyApplied
	W
	WeightedFlatSubject
	WaterDepth
	When
	WhiteBalance
	WhitePoint
	XMPToolkit
	XResolution
	Xap
	XmpDM
	XmpMeta
	YCbCrCoefficients
	YCbCrPositioning
	YCbCrSubSampling
	YResolution

	// Known tags intentionally not decoded into structs yet.
	AutoLateralCA
	Blacks2012
	CameraProfile
	CameraProfileDigest
	Clarity2012
	ColorNoiseReduction
	ColorNoiseReductionDetail
	ColorNoiseReductionSmoothness
	Contrast2012
	ConvertToGrayscale
	DefringeGreenAmount
	DefringeGreenHueHi
	DefringeGreenHueLo
	DefringePurpleAmount
	DefringePurpleHueHi
	DefringePurpleHueLo
	Dehaze
	Exposure2012
	GrainAmount
	GrainFrequency
	GrainSeed
	GrainSize
	HasCrop
	HasSettings
	Highlights2012
	LensManualDistortionAmount
	LensProfileChromaticAberrationScale
	LensProfileDigest
	LensProfileDistortionScale
	LensProfileEnable
	LensProfileFilename
	LensProfileName
	LensProfileSetup
	LensProfileVignettingScale
	LookName
	LuminanceNoiseReductionContrast
	LuminanceNoiseReductionDetail
	LuminanceSmoothing
	OverrideLookVignette
	ParametricDarks
	ParametricHighlightSplit
	ParametricHighlights
	ParametricLights
	ParametricMidtoneSplit
	ParametricShadowSplit
	ParametricShadows
	PerspectiveAspect
	PerspectiveHorizontal
	PerspectiveRotate
	PerspectiveScale
	PerspectiveUpright
	PerspectiveVertical
	PerspectiveX
	PerspectiveY
	PostCropVignetteAmount
	PostCropVignetteFeather
	PostCropVignetteHighlightContrast
	PostCropVignetteMidpoint
	PostCropVignetteRoundness
	PostCropVignetteStyle
	ProcessVersion
	ShadowTint
	Shadows2012
	SharpenDetail
	SharpenEdgeMasking
	SharpenRadius
	SplitToningBalance
	SplitToningHighlightHue
	SplitToningHighlightSaturation
	SplitToningShadowHue
	SplitToningShadowSaturation
	Tint
	ToneCurveName
	ToneCurveName2012
	ToneMapStrength
	UprightCenterMode
	UprightCenterNormX
	UprightCenterNormY
	UprightFocalLength35mm
	UprightFocalMode
	UprightFourSegmentsCount
	UprightPreview
	UprightTransformCount
	UprightVersion
	Version
	Vibrance
	VignetteAmount
	Whites2012
	DerivedFromInstanceID
	RegionAppliedToDimensionsH
	RegionAppliedToDimensionsUnit
	RegionAppliedToDimensionsW
	RegionAreaH
	RegionAreaUnit
	RegionAreaW
	RegionAreaX
	RegionAreaY
	RegionExtensionsAngleInfoRoll
	RegionExtensionsAngleInfoYaw
	RegionExtensionsConfidenceLevel
	RegionExtensionsFaceID
	RegionExtensionsTimeStamp
	RegionTypeTag
	AppliedToDimensions
	AreaTag
	ExtensionsTag
	NameTag
	RegionListTag
	RegionsTag
	RoleTag
	Unit
	X
	Y
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
	Copyright:                          "Copyright",
	CreateDate:                         "CreateDate",
	Creator:                            "creator",
	CreatorTool:                        "CreatorTool",
	CustomRendered:                     "CustomRendered",
	DateCreated:                        "DateCreated",
	DateTime:                           "DateTime",
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
	Label:                               "Label",
	Lang:                                "lang",
	LegacyIPTCDigest:                    "LegacyIPTCDigest",
	Lens:                                "Lens",
	LensID:                              "LensID",
	LensInfo:                            "LensInfo",
	LensModel:                           "LensModel",
	LensSerialNumber:                    "LensSerialNumber",
	Li:                                  "li",
	LightSource:                         "LightSource",
	LuminanceAdjustmentRed:              "LuminanceAdjustmentRed",
	LuminanceAdjustmentOrange:           "LuminanceAdjustmentOrange",
	LuminanceAdjustmentYellow:           "LuminanceAdjustmentYellow",
	LuminanceAdjustmentGreen:            "LuminanceAdjustmentGreen",
	LuminanceAdjustmentAqua:             "LuminanceAdjustmentAqua",
	LuminanceAdjustmentBlue:             "LuminanceAdjustmentBlue",
	LuminanceAdjustmentPurple:           "LuminanceAdjustmentPurple",
	LuminanceAdjustmentMagenta:          "LuminanceAdjustmentMagenta",
	Make:                                "Make",
	MaxApertureValue:                    "MaxApertureValue",
	MetadataDate:                        "MetadataDate",
	MeteringMode:                        "MeteringMode",
	Mode:                                "Mode",
	Model:                               "Model",
	ModifyDate:                          "ModifyDate",
	NativeDigest:                        "NativeDigest",
	Orientation:                         "Orientation",
	OriginalDocumentID:                  "OriginalDocumentID",
	Parameters:                          "parameters",
	ParseType:                           "parseType",
	PhotometricInterpretation:           "PhotometricInterpretation",
	PhotographicSensitivity:             "PhotographicSensitivity",
	PixelXDimension:                     "PixelXDimension",
	PixelYDimension:                     "PixelYDimension",
	Pick:                                "pick",
	PlanarConfiguration:                 "PlanarConfiguration",
	PreservedFileName:                   "PreservedFileName",
	Rating:                              "Rating",
	RawFileName:                         "RawFileName",
	RDF:                                 "RDF",
	RecommendedExposureIndex:            "RecommendedExposureIndex",
	RedEyeMode:                          "RedEyeMode",
	ResolutionUnit:                      "ResolutionUnit",
	Return:                              "Return",
	Rights:                              "rights",
	SamplesPerPixel:                     "SamplesPerPixel",
	Saturation:                          "Saturation",
	SaturationAdjustmentRed:             "SaturationAdjustmentRed",
	SaturationAdjustmentOrange:          "SaturationAdjustmentOrange",
	SaturationAdjustmentYellow:          "SaturationAdjustmentYellow",
	SaturationAdjustmentGreen:           "SaturationAdjustmentGreen",
	SaturationAdjustmentAqua:            "SaturationAdjustmentAqua",
	SaturationAdjustmentBlue:            "SaturationAdjustmentBlue",
	SaturationAdjustmentPurple:          "SaturationAdjustmentPurple",
	SaturationAdjustmentMagenta:         "SaturationAdjustmentMagenta",
	SceneCaptureType:                    "SceneCaptureType",
	SceneType:                           "SceneType",
	SensitivityType:                     "SensitivityType",
	Seq:                                 "Seq",
	SerialNumber:                        "SerialNumber",
	Sharpness:                           "Sharpness",
	ShutterSpeedValue:                   "ShutterSpeedValue",
	SidecarForExtension:                 "SidecarForExtension",
	Software:                            "Software",
	SoftwareAgent:                       "softwareAgent",
	StartTimecode:                       "startTimecode",
	StDim:                               "stDim",
	Subject:                             "subject",
	SubjectDistance:                     "SubjectDistance",
	SubsecTime:                          "SubsecTime",
	SubsecTimeDigitized:                 "SubsecTimeDigitized",
	SubsecTimeOriginal:                  "SubsecTimeOriginal",
	TapeName:                            "tapeName",
	Temperature:                         "Temperature",
	TimeValue:                           "timeValue",
	Title:                               "Title",
	ToneCurve:                           "ToneCurve",
	ToneCurveBlue:                       "ToneCurveBlue",
	ToneCurveGreen:                      "ToneCurveGreen",
	ToneCurvePV2012:                     "ToneCurvePV2012",
	ToneCurvePV2012Blue:                 "ToneCurvePV2012Blue",
	ToneCurvePV2012Green:                "ToneCurvePV2012Green",
	ToneCurvePV2012Red:                  "ToneCurvePV2012Red",
	ToneCurveRed:                        "ToneCurveRed",
	UserComment:                         "UserComment",
	VignetteCorrectionAlreadyApplied:    "VignetteCorrectionAlreadyApplied",
	VideoFieldOrder:                     "videoFieldOrder",
	VideoFrameRate:                      "videoFrameRate",
	VideoFrameSize:                      "videoFrameSize",
	VideoPixelAspectRatio:               "videoPixelAspectRatio",
	VideoPixelDepth:                     "videoPixelDepth",
	W:                                   "w",
	When:                                "when",
	WeightedFlatSubject:                 "weightedFlatSubject",
	WhiteBalance:                        "WhiteBalance",
	Xap:                                 "xap",
	XMPToolkit:                          "xmptk",
	XmpDM:                               "xmpDM",
	XmpMeta:                             "xmpmeta",
	XResolution:                         "XResolution",
	YCbCrPositioning:                    "YCbCrPositioning",
	YResolution:                         "YResolution",
	AutoLateralCA:                       "AutoLateralCA",
	Blacks2012:                          "Blacks2012",
	CameraProfile:                       "CameraProfile",
	CameraProfileDigest:                 "CameraProfileDigest",
	Clarity2012:                         "Clarity2012",
	ColorNoiseReduction:                 "ColorNoiseReduction",
	ColorNoiseReductionDetail:           "ColorNoiseReductionDetail",
	ColorNoiseReductionSmoothness:       "ColorNoiseReductionSmoothness",
	Contrast2012:                        "Contrast2012",
	ConvertToGrayscale:                  "ConvertToGrayscale",
	DefringeGreenAmount:                 "DefringeGreenAmount",
	DefringeGreenHueHi:                  "DefringeGreenHueHi",
	DefringeGreenHueLo:                  "DefringeGreenHueLo",
	DefringePurpleAmount:                "DefringePurpleAmount",
	DefringePurpleHueHi:                 "DefringePurpleHueHi",
	DefringePurpleHueLo:                 "DefringePurpleHueLo",
	Dehaze:                              "Dehaze",
	Exposure2012:                        "Exposure2012",
	GrainAmount:                         "GrainAmount",
	GrainFrequency:                      "GrainFrequency",
	GrainSeed:                           "GrainSeed",
	GrainSize:                           "GrainSize",
	HasCrop:                             "HasCrop",
	HasSettings:                         "HasSettings",
	Highlights2012:                      "Highlights2012",
	LensManualDistortionAmount:          "LensManualDistortionAmount",
	LensProfileChromaticAberrationScale: "LensProfileChromaticAberrationScale",
	LensProfileDigest:                   "LensProfileDigest",
	LensProfileDistortionScale:          "LensProfileDistortionScale",
	LensProfileEnable:                   "LensProfileEnable",
	LensProfileFilename:                 "LensProfileFilename",
	LensProfileName:                     "LensProfileName",
	LensProfileSetup:                    "LensProfileSetup",
	LensProfileVignettingScale:          "LensProfileVignettingScale",
	LookName:                            "LookName",
	LuminanceNoiseReductionContrast:     "LuminanceNoiseReductionContrast",
	LuminanceNoiseReductionDetail:       "LuminanceNoiseReductionDetail",
	LuminanceSmoothing:                  "LuminanceSmoothing",
	OverrideLookVignette:                "OverrideLookVignette",
	ParametricDarks:                     "ParametricDarks",
	ParametricHighlightSplit:            "ParametricHighlightSplit",
	ParametricHighlights:                "ParametricHighlights",
	ParametricLights:                    "ParametricLights",
	ParametricMidtoneSplit:              "ParametricMidtoneSplit",
	ParametricShadowSplit:               "ParametricShadowSplit",
	ParametricShadows:                   "ParametricShadows",
	PerspectiveAspect:                   "PerspectiveAspect",
	PerspectiveHorizontal:               "PerspectiveHorizontal",
	PerspectiveRotate:                   "PerspectiveRotate",
	PerspectiveScale:                    "PerspectiveScale",
	PerspectiveUpright:                  "PerspectiveUpright",
	PerspectiveVertical:                 "PerspectiveVertical",
	PerspectiveX:                        "PerspectiveX",
	PerspectiveY:                        "PerspectiveY",
	PostCropVignetteAmount:              "PostCropVignetteAmount",
	PostCropVignetteFeather:             "PostCropVignetteFeather",
	PostCropVignetteHighlightContrast:   "PostCropVignetteHighlightContrast",
	PostCropVignetteMidpoint:            "PostCropVignetteMidpoint",
	PostCropVignetteRoundness:           "PostCropVignetteRoundness",
	PostCropVignetteStyle:               "PostCropVignetteStyle",
	ProcessVersion:                      "ProcessVersion",
	ShadowTint:                          "ShadowTint",
	Shadows2012:                         "Shadows2012",
	SharpenDetail:                       "SharpenDetail",
	SharpenEdgeMasking:                  "SharpenEdgeMasking",
	SharpenRadius:                       "SharpenRadius",
	SplitToningBalance:                  "SplitToningBalance",
	SplitToningHighlightHue:             "SplitToningHighlightHue",
	SplitToningHighlightSaturation:      "SplitToningHighlightSaturation",
	SplitToningShadowHue:                "SplitToningShadowHue",
	SplitToningShadowSaturation:         "SplitToningShadowSaturation",
	Tint:                                "Tint",
	ToneCurveName:                       "ToneCurveName",
	ToneCurveName2012:                   "ToneCurveName2012",
	ToneMapStrength:                     "ToneMapStrength",
	UprightCenterMode:                   "UprightCenterMode",
	UprightCenterNormX:                  "UprightCenterNormX",
	UprightCenterNormY:                  "UprightCenterNormY",
	UprightFocalLength35mm:              "UprightFocalLength35mm",
	UprightFocalMode:                    "UprightFocalMode",
	UprightFourSegmentsCount:            "UprightFourSegmentsCount",
	UprightPreview:                      "UprightPreview",
	UprightTransformCount:               "UprightTransformCount",
	UprightVersion:                      "UprightVersion",
	Version:                             "Version",
	Vibrance:                            "Vibrance",
	VignetteAmount:                      "VignetteAmount",
	Whites2012:                          "Whites2012",
	DerivedFromInstanceID:               "DerivedFromInstanceID",
	RegionAppliedToDimensionsH:          "RegionAppliedToDimensionsH",
	RegionAppliedToDimensionsUnit:       "RegionAppliedToDimensionsUnit",
	RegionAppliedToDimensionsW:          "RegionAppliedToDimensionsW",
	RegionAreaH:                         "RegionAreaH",
	RegionAreaUnit:                      "RegionAreaUnit",
	RegionAreaW:                         "RegionAreaW",
	RegionAreaX:                         "RegionAreaX",
	RegionAreaY:                         "RegionAreaY",
	RegionExtensionsAngleInfoRoll:       "RegionExtensionsAngleInfoRoll",
	RegionExtensionsAngleInfoYaw:        "RegionExtensionsAngleInfoYaw",
	RegionExtensionsConfidenceLevel:     "RegionExtensionsConfidenceLevel",
	RegionExtensionsFaceID:              "RegionExtensionsFaceID",
	RegionExtensionsTimeStamp:           "RegionExtensionsTimeStamp",
	RegionTypeTag:                       "RegionType",
	AppliedToDimensions:                 "AppliedToDimensions",
	AreaTag:                             "Area",
	ExtensionsTag:                       "Extensions",
	NameTag:                             "Name",
	RegionListTag:                       "RegionList",
	RegionsTag:                          "Regions",
	RoleTag:                             "Role",
	Unit:                                "unit",
	X:                                   "x",
	Y:                                   "y",
	GPSDOP:                              "GPSDOP",
	GPSMeasureMode:                      "GPSMeasureMode",
	GPSSatellites:                       "GPSSatellites",
	Contributor:                         "contributor",
	Coverage:                            "coverage",
	Date:                                "date",
	Identifier:                          "identifier",
	Language:                            "language",
	Publisher:                           "publisher",
	Relation:                            "relation",
	Source:                              "source",
	Type:                                "type",
	CFAPattern:                          "CFAPattern",
	CFAPatternColumns:                   "CFAPatternColumns",
	CFAPatternRows:                      "CFAPatternRows",
	CFAPatternValues:                    "CFAPatternValues",
	DeviceSettingDescription:            "DeviceSettingDescription",
	DeviceSettingDescriptionColumns:     "DeviceSettingDescriptionColumns",
	DeviceSettingDescriptionRows:        "DeviceSettingDescriptionRows",
	DeviceSettingDescriptionSettings:    "DeviceSettingDescriptionSettings",
	ExposureIndex:                       "ExposureIndex",
	FlashEnergy:                         "FlashEnergy",
	GPSAreaInformation:                  "GPSAreaInformation",
	GPSDestBearing:                      "GPSDestBearing",
	GPSDestBearingRef:                   "GPSDestBearingRef",
	GPSDestDistance:                     "GPSDestDistance",
	GPSDestDistanceRef:                  "GPSDestDistanceRef",
	GPSDestLatitude:                     "GPSDestLatitude",
	GPSDestLongitude:                    "GPSDestLongitude",
	GPSHPositioningError:                "GPSHPositioningError",
	GPSImgDirection:                     "GPSImgDirection",
	GPSImgDirectionRef:                  "GPSImgDirectionRef",
	GPSProcessingMethod:                 "GPSProcessingMethod",
	GPSSpeed:                            "GPSSpeed",
	GPSSpeedRef:                         "GPSSpeedRef",
	GPSTrack:                            "GPSTrack",
	GPSTrackRef:                         "GPSTrackRef",
	ImageUniqueID:                       "ImageUniqueID",
	MakerNote:                           "MakerNote",
	OECF:                                "OECF",
	OECFColumns:                         "OECFColumns",
	OECFNames:                           "OECFNames",
	OECFRows:                            "OECFRows",
	OECFValues:                          "OECFValues",
	RelatedSoundFile:                    "RelatedSoundFile",
	SensingMethod:                       "SensingMethod",
	SpatialFrequencyResponse:            "SpatialFrequencyResponse",
	SpatialFrequencyResponseColumns:     "SpatialFrequencyResponseColumns",
	SpatialFrequencyResponseNames:       "SpatialFrequencyResponseNames",
	SpatialFrequencyResponseRows:        "SpatialFrequencyResponseRows",
	SpatialFrequencyResponseValues:      "SpatialFrequencyResponseValues",
	SpectralSensitivity:                 "SpectralSensitivity",
	SubjectArea:                         "SubjectArea",
	SubjectDistanceRange:                "SubjectDistanceRange",
	SubjectLocation:                     "SubjectLocation",
	Acceleration:                        "Acceleration",
	AmbientTemperature:                  "AmbientTemperature",
	Artist:                              "Artist",
	CameraElevationAngle:                "CameraElevationAngle",
	CameraFirmware:                      "CameraFirmware",
	CompImageImagesPerSequence:          "CompImageImagesPerSequence",
	CompImageMaxExposureAll:             "CompImageMaxExposureAll",
	CompImageMaxExposureUsed:            "CompImageMaxExposureUsed",
	CompImageMinExposureAll:             "CompImageMinExposureAll",
	CompImageMinExposureUsed:            "CompImageMinExposureUsed",
	CompImageNumSequences:               "CompImageNumSequences",
	CompImageSumExposureAll:             "CompImageSumExposureAll",
	CompImageSumExposureUsed:            "CompImageSumExposureUsed",
	CompImageTotalExposurePeriod:        "CompImageTotalExposurePeriod",
	CompImageValues:                     "CompImageValues",
	CompositeImage:                      "CompositeImage",
	CompositeImageCount:                 "CompositeImageCount",
	CompositeImageExposureTimes:         "CompositeImageExposureTimes",
	Gamma:                               "Gamma",
	Humidity:                            "Humidity",
	ISOSpeed:                            "ISOSpeed",
	ISOSpeedLatitudeyyy:                 "ISOSpeedLatitudeyyy",
	ISOSpeedLatitudezzz:                 "ISOSpeedLatitudezzz",
	ImageEditingSoftware:                "ImageEditingSoftware",
	ImageEditor:                         "ImageEditor",
	ImageTitle:                          "ImageTitle",
	LensMake:                            "LensMake",
	MetadataEditingSoftware:             "MetadataEditingSoftware",
	OwnerName:                           "OwnerName",
	Photographer:                        "Photographer",
	Pressure:                            "Pressure",
	PrimaryChromaticities:               "PrimaryChromaticities",
	RAWDevelopingSoftware:               "RAWDevelopingSoftware",
	ReferenceBlackWhite:                 "ReferenceBlackWhite",
	StandardOutputSensitivity:           "StandardOutputSensitivity",
	TransferFunction:                    "TransferFunction",
	WaterDepth:                          "WaterDepth",
	WhitePoint:                          "WhitePoint",
	YCbCrCoefficients:                   "YCbCrCoefficients",
	YCbCrSubSampling:                    "YCbCrSubSampling",
}
