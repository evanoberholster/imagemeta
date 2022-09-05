package ifds

import "github.com/evanoberholster/imagemeta/exif2/tag"

// TagString returns the string representation of a tag.ID
func TagString(id tag.ID) string {
	name, ok := RootIfdTagIDMap[id]
	if !ok {
		return id.String()
	}
	return name
}

// RootIfdTagIDMap is a Map of tag.ID to string for the RootIfd tags
var RootIfdTagIDMap = map[tag.ID]string{
	ProcessingSoftware:          "ProcessingSoftware",
	NewSubfileType:              "NewSubfileType",
	SubfileType:                 "SubfileType",
	ImageWidth:                  "ImageWidth",
	ImageLength:                 "ImageLength",
	BitsPerSample:               "BitsPerSample",
	Compression:                 "Compression",
	PhotometricInterpretation:   "PhotometricInterpretation",
	Thresholding:                "Thresholding",
	CellWidth:                   "CellWidth",
	CellLength:                  "CellLength",
	FillOrder:                   "FillOrder",
	DocumentName:                "DocumentName",
	ImageDescription:            "ImageDescription",
	Make:                        "Make",
	Model:                       "Model",
	StripOffsets:                "StripOffsets",
	Orientation:                 "Orientation",
	SamplesPerPixel:             "SamplesPerPixel",
	RowsPerStrip:                "RowsPerStrip",
	StripByteCounts:             "StripByteCounts",
	XResolution:                 "XResolution",
	YResolution:                 "YResolution",
	PlanarConfiguration:         "PlanarConfiguration",
	GrayResponseUnit:            "GrayResponseUnit",
	GrayResponseCurve:           "GrayResponseCurve",
	T4Options:                   "T4Options",
	T6Options:                   "T6Options",
	ResolutionUnit:              "ResolutionUnit",
	PageNumber:                  "PageNumber",
	TransferFunction:            "TransferFunction",
	Software:                    "Software",
	DateTime:                    "DateTime",
	Artist:                      "Artist",
	HostComputer:                "HostComputer",
	Predictor:                   "Predictor",
	WhitePoint:                  "WhitePoint",
	PrimaryChromaticities:       "PrimaryChromaticities",
	ColorMap:                    "ColorMap",
	HalftoneHints:               "HalftoneHints",
	TileWidth:                   "TileWidth",
	TileLength:                  "TileLength",
	TileOffsets:                 "TileOffsets",
	TileByteCounts:              "TileByteCounts",
	SubIFDs:                     "SubIFDs",
	InkSet:                      "InkSet",
	InkNames:                    "InkNames",
	NumberOfInks:                "NumberOfInks",
	DotRange:                    "DotRange",
	TargetPrinter:               "TargetPrinter",
	ExtraSamples:                "ExtraSamples",
	SampleFormat:                "SampleFormat",
	SMinSampleValue:             "SMinSampleValue",
	SMaxSampleValue:             "SMaxSampleValue",
	TransferRange:               "TransferRange",
	ClipPath:                    "ClipPath",
	XClipPathUnits:              "XClipPathUnits",
	YClipPathUnits:              "YClipPathUnits",
	Indexed:                     "Indexed",
	JPEGTables:                  "JPEGTables",
	OPIProxy:                    "OPIProxy",
	JPEGProc:                    "JPEGProc",
	JPEGInterchangeFormat:       "JPEGInterchangeFormat",
	JPEGInterchangeFormatLength: "JPEGInterchangeFormatLength",
	JPEGRestartInterval:         "JPEGRestartInterval",
	JPEGLosslessPredictors:      "JPEGLosslessPredictors",
	JPEGPointTransforms:         "JPEGPointTransforms",
	JPEGQTables:                 "JPEGQTables",
	JPEGDCTables:                "JPEGDCTables",
	JPEGACTables:                "JPEGACTables",
	YCbCrCoefficients:           "YCbCrCoefficients",
	YCbCrSubSampling:            "YCbCrSubSampling",
	YCbCrPositioning:            "YCbCrPositioning",
	ReferenceBlackWhite:         "ReferenceBlackWhite",
	XMLPacket:                   "XMLPacket",
	Rating:                      "Rating",
	RatingPercent:               "RatingPercent",
	ImageID:                     "ImageID",
	CFARepeatPatternDim:         "CFARepeatPatternDim",
	CFAPattern:                  "CFAPattern",
	BatteryLevel:                "BatteryLevel",
	Copyright:                   "Copyright",
	ExposureTime:                "ExposureTime",
	FNumber:                     "FNumber",
	IPTCNAA:                     "IPTCNAA",
	ImageResources:              "ImageResources",
	ExifTag:                     "ExifTag",
	InterColorProfile:           "InterColorProfile",
	ExposureProgram:             "ExposureProgram",
	SpectralSensitivity:         "SpectralSensitivity",
	GPSTag:                      "GPSTag",
	ISOSpeedRatings:             "ISOSpeedRatings",
	OECF:                        "OECF",
	Interlace:                   "Interlace",
	SensitivityType:             "SensitivityType",
	TimeZoneOffset:              "TimeZoneOffset",
	SelfTimerMode:               "SelfTimerMode",
	RecommendedExposureIndex:    "RecommendedExposureIndex",
	DateTimeOriginal:            "DateTimeOriginal",
	DateTimeDigitized:           "DateTimeDigitized",
	CompressedBitsPerPixel:      "CompressedBitsPerPixel",
	ShutterSpeedValue:           "ShutterSpeedValue",
	ApertureValue:               "ApertureValue",
	BrightnessValue:             "BrightnessValue",
	ExposureBiasValue:           "ExposureBiasValue",
	MaxApertureValue:            "MaxApertureValue",
	SubjectDistance:             "SubjectDistance",
	MeteringMode:                "MeteringMode",
	LightSource:                 "LightSource",
	Flash:                       "Flash",
	FocalLength:                 "FocalLength",
	FlashEnergy:                 "FlashEnergy",
	SpatialFrequencyResponse:    "SpatialFrequencyResponse",
	Noise:                       "Noise",
	FocalPlaneXResolution:       "FocalPlaneXResolution",
	FocalPlaneYResolution:       "FocalPlaneYResolution",
	FocalPlaneResolutionUnit:    "FocalPlaneResolutionUnit",
	ImageNumber:                 "ImageNumber",
	SecurityClassification:      "SecurityClassification",
	ImageHistory:                "ImageHistory",
	SubjectLocation:             "SubjectLocation",
	ExposureIndex:               "ExposureIndex",
	TIFFEPStandardID:            "TIFFEPStandardID",
	SensingMethod:               "SensingMethod",
	XPTitle:                     "XPTitle",
	XPComment:                   "XPComment",
	XPAuthor:                    "XPAuthor",
	XPKeywords:                  "XPKeywords",
	XPSubject:                   "XPSubject",
	PrintImageMatching:          "PrintImageMatching",
	DNGVersion:                  "DNGVersion",
	DNGBackwardVersion:          "DNGBackwardVersion",
	UniqueCameraModel:           "UniqueCameraModel",
	LocalizedCameraModel:        "LocalizedCameraModel",
	CFAPlaneColor:               "CFAPlaneColor",
	CFALayout:                   "CFALayout",
	LinearizationTable:          "LinearizationTable",
	BlackLevelRepeatDim:         "BlackLevelRepeatDim",
	BlackLevel:                  "BlackLevel",
	BlackLevelDeltaH:            "BlackLevelDeltaH",
	BlackLevelDeltaV:            "BlackLevelDeltaV",
	WhiteLevel:                  "WhiteLevel",
	DefaultScale:                "DefaultScale",
	DefaultCropOrigin:           "DefaultCropOrigin",
	DefaultCropSize:             "DefaultCropSize",
	ColorMatrix1:                "ColorMatrix1",
	ColorMatrix2:                "ColorMatrix2",
	CameraCalibration1:          "CameraCalibration1",
	CameraCalibration2:          "CameraCalibration2",
	ReductionMatrix1:            "ReductionMatrix1",
	ReductionMatrix2:            "ReductionMatrix2",
	AnalogBalance:               "AnalogBalance",
	AsShotNeutral:               "AsShotNeutral",
	AsShotWhiteXY:               "AsShotWhiteXY",
	BaselineExposure:            "BaselineExposure",
	BaselineNoise:               "BaselineNoise",
	BaselineSharpness:           "BaselineSharpness",
	BayerGreenSplit:             "BayerGreenSplit",
	LinearResponseLimit:         "LinearResponseLimit",
	CameraSerialNumber:          "CameraSerialNumber",
	LensInfo:                    "LensInfo",
	ChromaBlurRadius:            "ChromaBlurRadius",
	AntiAliasStrength:           "AntiAliasStrength",
	ShadowScale:                 "ShadowScale",
	DNGPrivateData:              "DNGPrivateData",
	MakerNoteSafety:             "MakerNoteSafety",
	CalibrationIlluminant1:      "CalibrationIlluminant1",
	CalibrationIlluminant2:      "CalibrationIlluminant2",
	BestQualityScale:            "BestQualityScale",
	RawDataUniqueID:             "RawDataUniqueID",
	OriginalRawFileName:         "OriginalRawFileName",
	OriginalRawFileData:         "OriginalRawFileData",
	ActiveArea:                  "ActiveArea",
	MaskedAreas:                 "MaskedAreas",
	AsShotICCProfile:            "AsShotICCProfile",
	AsShotPreProfileMatrix:      "AsShotPreProfileMatrix",
	CurrentICCProfile:           "CurrentICCProfile",
	CurrentPreProfileMatrix:     "CurrentPreProfileMatrix",
	ColorimetricReference:       "ColorimetricReference",
	CameraCalibrationSignature:  "CameraCalibrationSignature",
	ProfileCalibrationSignature: "ProfileCalibrationSignature",
	AsShotProfileName:           "AsShotProfileName",
	NoiseReductionApplied:       "NoiseReductionApplied",
	ProfileName:                 "ProfileName",
	ProfileHueSatMapDims:        "ProfileHueSatMapDims",
	ProfileHueSatMapData1:       "ProfileHueSatMapData1",
	ProfileHueSatMapData2:       "ProfileHueSatMapData2",
	ProfileToneCurve:            "ProfileToneCurve",
	ProfileEmbedPolicy:          "ProfileEmbedPolicy",
	ProfileCopyright:            "ProfileCopyright",
	ForwardMatrix1:              "ForwardMatrix1",
	ForwardMatrix2:              "ForwardMatrix2",
	PreviewApplicationName:      "PreviewApplicationName",
	PreviewApplicationVersion:   "PreviewApplicationVersion",
	PreviewSettingsName:         "PreviewSettingsName",
	PreviewSettingsDigest:       "PreviewSettingsDigest",
	PreviewColorSpace:           "PreviewColorSpace",
	PreviewDateTime:             "PreviewDateTime",
	RawImageDigest:              "RawImageDigest",
	OriginalRawFileDigest:       "OriginalRawFileDigest",
	SubTileBlockSize:            "SubTileBlockSize",
	RowInterleaveFactor:         "RowInterleaveFactor",
	ProfileLookTableDims:        "ProfileLookTableDims",
	ProfileLookTableData:        "ProfileLookTableData",
	OpcodeList1:                 "OpcodeList1",
	OpcodeList2:                 "OpcodeList2",
	OpcodeList3:                 "OpcodeList3",
	NoiseProfile:                "NoiseProfile",
	CacheVersion:                "CacheVersion",
}

// RootIFD TagIDs
const (
	ProcessingSoftware          tag.ID = 0x000b
	NewSubfileType              tag.ID = 0x00fe
	SubfileType                 tag.ID = 0x00ff
	ImageWidth                  tag.ID = 0x0100
	ImageLength                 tag.ID = 0x0101
	BitsPerSample               tag.ID = 0x0102
	Compression                 tag.ID = 0x0103
	PhotometricInterpretation   tag.ID = 0x0106
	Thresholding                tag.ID = 0x0107
	CellWidth                   tag.ID = 0x0108
	CellLength                  tag.ID = 0x0109
	FillOrder                   tag.ID = 0x010a
	DocumentName                tag.ID = 0x010d
	ImageDescription            tag.ID = 0x010e
	Make                        tag.ID = 0x010f
	Model                       tag.ID = 0x0110
	StripOffsets                tag.ID = 0x0111
	Orientation                 tag.ID = 0x0112
	SamplesPerPixel             tag.ID = 0x0115
	RowsPerStrip                tag.ID = 0x0116
	StripByteCounts             tag.ID = 0x0117
	XResolution                 tag.ID = 0x011a
	YResolution                 tag.ID = 0x011b
	PlanarConfiguration         tag.ID = 0x011c
	GrayResponseUnit            tag.ID = 0x0122
	GrayResponseCurve           tag.ID = 0x0123
	T4Options                   tag.ID = 0x0124
	T6Options                   tag.ID = 0x0125
	ResolutionUnit              tag.ID = 0x0128
	PageNumber                  tag.ID = 0x0129
	TransferFunction            tag.ID = 0x012d
	Software                    tag.ID = 0x0131
	DateTime                    tag.ID = 0x0132
	Artist                      tag.ID = 0x013b
	HostComputer                tag.ID = 0x013c
	Predictor                   tag.ID = 0x013d
	WhitePoint                  tag.ID = 0x013e
	PrimaryChromaticities       tag.ID = 0x013f
	ColorMap                    tag.ID = 0x0140
	HalftoneHints               tag.ID = 0x0141
	TileWidth                   tag.ID = 0x0142
	TileLength                  tag.ID = 0x0143
	TileOffsets                 tag.ID = 0x0144
	TileByteCounts              tag.ID = 0x0145
	SubIFDs                     tag.ID = 0x014a
	InkSet                      tag.ID = 0x014c
	InkNames                    tag.ID = 0x014d
	NumberOfInks                tag.ID = 0x014e
	DotRange                    tag.ID = 0x0150
	TargetPrinter               tag.ID = 0x0151
	ExtraSamples                tag.ID = 0x0152
	SampleFormat                tag.ID = 0x0153
	SMinSampleValue             tag.ID = 0x0154
	SMaxSampleValue             tag.ID = 0x0155
	TransferRange               tag.ID = 0x0156
	ClipPath                    tag.ID = 0x0157
	XClipPathUnits              tag.ID = 0x0158
	YClipPathUnits              tag.ID = 0x0159
	Indexed                     tag.ID = 0x015a
	JPEGTables                  tag.ID = 0x015b
	OPIProxy                    tag.ID = 0x015f
	JPEGProc                    tag.ID = 0x0200
	JPEGInterchangeFormat       tag.ID = 0x0201
	JPEGInterchangeFormatLength tag.ID = 0x0202
	JPEGRestartInterval         tag.ID = 0x0203
	JPEGLosslessPredictors      tag.ID = 0x0205
	JPEGPointTransforms         tag.ID = 0x0206
	JPEGQTables                 tag.ID = 0x0207
	JPEGDCTables                tag.ID = 0x0208
	JPEGACTables                tag.ID = 0x0209
	YCbCrCoefficients           tag.ID = 0x0211
	YCbCrSubSampling            tag.ID = 0x0212
	YCbCrPositioning            tag.ID = 0x0213
	ReferenceBlackWhite         tag.ID = 0x0214
	XMLPacket                   tag.ID = 0x02bc
	Rating                      tag.ID = 0x4746
	RatingPercent               tag.ID = 0x4749
	ImageID                     tag.ID = 0x800d
	CFARepeatPatternDim         tag.ID = 0x828d
	CFAPattern                  tag.ID = 0x828e
	BatteryLevel                tag.ID = 0x828f
	Copyright                   tag.ID = 0x8298
	ExposureTime                tag.ID = 0x829a // IFD/EXIF and IFD
	FNumber                     tag.ID = 0x829d
	IPTCNAA                     tag.ID = 0x83bb
	ImageResources              tag.ID = 0x8649
	ExifTag                     tag.ID = 0x8769
	InterColorProfile           tag.ID = 0x8773
	ExposureProgram             tag.ID = 0x8822
	SpectralSensitivity         tag.ID = 0x8824
	GPSTag                      tag.ID = 0x8825
	ISOSpeedRatings             tag.ID = 0x8827
	OECF                        tag.ID = 0x8828
	Interlace                   tag.ID = 0x8829
	SensitivityType             tag.ID = 0x8830
	TimeZoneOffset              tag.ID = 0x882a
	SelfTimerMode               tag.ID = 0x882b
	RecommendedExposureIndex    tag.ID = 0x8832
	DateTimeOriginal            tag.ID = 0x9003
	DateTimeDigitized           tag.ID = 0x9004
	CompressedBitsPerPixel      tag.ID = 0x9102
	ShutterSpeedValue           tag.ID = 0x9201
	ApertureValue               tag.ID = 0x9202
	BrightnessValue             tag.ID = 0x9203
	ExposureBiasValue           tag.ID = 0x9204
	MaxApertureValue            tag.ID = 0x9205
	SubjectDistance             tag.ID = 0x9206
	MeteringMode                tag.ID = 0x9207
	LightSource                 tag.ID = 0x9208
	Flash                       tag.ID = 0x9209
	FocalLength                 tag.ID = 0x920a
	FlashEnergy                 tag.ID = 0x920b
	SpatialFrequencyResponse    tag.ID = 0x920c
	Noise                       tag.ID = 0x920d
	FocalPlaneXResolution       tag.ID = 0x920e
	FocalPlaneYResolution       tag.ID = 0x920f
	FocalPlaneResolutionUnit    tag.ID = 0x9210
	ImageNumber                 tag.ID = 0x9211
	SecurityClassification      tag.ID = 0x9212
	ImageHistory                tag.ID = 0x9213
	SubjectLocation             tag.ID = 0x9214
	ExposureIndex               tag.ID = 0x9215
	TIFFEPStandardID            tag.ID = 0x9216
	SensingMethod               tag.ID = 0x9217
	XPTitle                     tag.ID = 0x9c9b
	XPComment                   tag.ID = 0x9c9c
	XPAuthor                    tag.ID = 0x9c9d
	XPKeywords                  tag.ID = 0x9c9e
	XPSubject                   tag.ID = 0x9c9f
	PrintImageMatching          tag.ID = 0xc4a5
	DNGVersion                  tag.ID = 0xc612
	DNGBackwardVersion          tag.ID = 0xc613
	UniqueCameraModel           tag.ID = 0xc614
	LocalizedCameraModel        tag.ID = 0xc615
	CFAPlaneColor               tag.ID = 0xc616
	CFALayout                   tag.ID = 0xc617
	LinearizationTable          tag.ID = 0xc618
	BlackLevelRepeatDim         tag.ID = 0xc619
	BlackLevel                  tag.ID = 0xc61a
	BlackLevelDeltaH            tag.ID = 0xc61b
	BlackLevelDeltaV            tag.ID = 0xc61c
	WhiteLevel                  tag.ID = 0xc61d
	DefaultScale                tag.ID = 0xc61e
	DefaultCropOrigin           tag.ID = 0xc61f
	DefaultCropSize             tag.ID = 0xc620
	ColorMatrix1                tag.ID = 0xc621
	ColorMatrix2                tag.ID = 0xc622
	CameraCalibration1          tag.ID = 0xc623
	CameraCalibration2          tag.ID = 0xc624
	ReductionMatrix1            tag.ID = 0xc625
	ReductionMatrix2            tag.ID = 0xc626
	AnalogBalance               tag.ID = 0xc627
	AsShotNeutral               tag.ID = 0xc628
	AsShotWhiteXY               tag.ID = 0xc629
	BaselineExposure            tag.ID = 0xc62a
	BaselineNoise               tag.ID = 0xc62b
	BaselineSharpness           tag.ID = 0xc62c
	BayerGreenSplit             tag.ID = 0xc62d
	LinearResponseLimit         tag.ID = 0xc62e
	CameraSerialNumber          tag.ID = 0xc62f
	LensInfo                    tag.ID = 0xc630
	ChromaBlurRadius            tag.ID = 0xc631
	AntiAliasStrength           tag.ID = 0xc632
	ShadowScale                 tag.ID = 0xc633
	DNGPrivateData              tag.ID = 0xc634
	MakerNoteSafety             tag.ID = 0xc635
	CalibrationIlluminant1      tag.ID = 0xc65a
	CalibrationIlluminant2      tag.ID = 0xc65b
	BestQualityScale            tag.ID = 0xc65c
	RawDataUniqueID             tag.ID = 0xc65d
	OriginalRawFileName         tag.ID = 0xc68b
	OriginalRawFileData         tag.ID = 0xc68c
	ActiveArea                  tag.ID = 0xc68d
	MaskedAreas                 tag.ID = 0xc68e
	AsShotICCProfile            tag.ID = 0xc68f
	AsShotPreProfileMatrix      tag.ID = 0xc690
	CurrentICCProfile           tag.ID = 0xc691
	CurrentPreProfileMatrix     tag.ID = 0xc692
	ColorimetricReference       tag.ID = 0xc6bf
	CameraCalibrationSignature  tag.ID = 0xc6f3
	ProfileCalibrationSignature tag.ID = 0xc6f4
	AsShotProfileName           tag.ID = 0xc6f6
	NoiseReductionApplied       tag.ID = 0xc6f7
	ProfileName                 tag.ID = 0xc6f8
	ProfileHueSatMapDims        tag.ID = 0xc6f9
	ProfileHueSatMapData1       tag.ID = 0xc6fa
	ProfileHueSatMapData2       tag.ID = 0xc6fb
	ProfileToneCurve            tag.ID = 0xc6fc
	ProfileEmbedPolicy          tag.ID = 0xc6fd
	ProfileCopyright            tag.ID = 0xc6fe
	ForwardMatrix1              tag.ID = 0xc714
	ForwardMatrix2              tag.ID = 0xc715
	PreviewApplicationName      tag.ID = 0xc716
	PreviewApplicationVersion   tag.ID = 0xc717
	PreviewSettingsName         tag.ID = 0xc718
	PreviewSettingsDigest       tag.ID = 0xc719
	PreviewColorSpace           tag.ID = 0xc71a
	PreviewDateTime             tag.ID = 0xc71b
	RawImageDigest              tag.ID = 0xc71c
	OriginalRawFileDigest       tag.ID = 0xc71d
	SubTileBlockSize            tag.ID = 0xc71e
	RowInterleaveFactor         tag.ID = 0xc71f
	ProfileLookTableDims        tag.ID = 0xc725
	ProfileLookTableData        tag.ID = 0xc726
	OpcodeList1                 tag.ID = 0xc740
	OpcodeList2                 tag.ID = 0xc741
	OpcodeList3                 tag.ID = 0xc74e
	NoiseProfile                tag.ID = 0xc761
	CacheVersion                tag.ID = 0xc7aa
)
