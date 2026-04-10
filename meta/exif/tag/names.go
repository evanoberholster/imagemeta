package tag

const (
	// Pseudo tag used by parser to model TIFF next-IFD chain.
	TagNextIFD ID = 0xffff

	// Root IFD tags (core set).
	TagSubfileType                  ID = 0x00fe
	TagImageWidth                   ID = 0x0100
	TagImageLength                  ID = 0x0101
	TagBitsPerSample                ID = 0x0102
	TagCompression                  ID = 0x0103
	TagPhotometricInterpretation    ID = 0x0106
	TagImageDescription             ID = 0x010e
	TagMake                         ID = 0x010f
	TagModel                        ID = 0x0110
	TagStripOffsets                 ID = 0x0111
	TagOrientation                  ID = 0x0112
	TagSamplesPerPixel              ID = 0x0115
	TagRowsPerStrip                 ID = 0x0116
	TagStripByteCounts              ID = 0x0117
	TagMinSampleValue               ID = 0x0118
	TagXResolution                  ID = 0x011a
	TagYResolution                  ID = 0x011b
	TagPlanarConfiguration          ID = 0x011c
	TagResolutionUnit               ID = 0x0128
	TagThumbnailOffset              ID = 0x0201
	TagThumbnailLength              ID = 0x0202
	TagReferenceBlackWhite          ID = 0x0214
	TagSoftware                     ID = 0x0131
	TagDateTime                     ID = 0x0132
	TagArtist                       ID = 0x013b
	TagTileWidth                    ID = 0x0142
	TagTileLength                   ID = 0x0143
	TagTileOffsets                  ID = 0x0144
	TagTileByteCounts               ID = 0x0145
	TagSubIFDs                      ID = 0x014a
	TagApplicationNotes             ID = 0x02bc
	TagCFARepeatPatternDim          ID = 0x828d
	TagCFAPattern2                  ID = 0x828e
	TagCopyright                    ID = 0x8298
	TagExifIFDPointer               ID = 0x8769
	TagGPSIFDPointer                ID = 0x8825
	TagDNGVersion                   ID = 0xc612
	TagDNGBackwardVersion           ID = 0xc613
	TagUniqueCameraModel            ID = 0xc614
	TagLocalizedCameraModel         ID = 0xc615
	TagCFAPlaneColor                ID = 0xc616
	TagCFALayout                    ID = 0xc617
	TagBlackLevelRepeatDim          ID = 0xc619
	TagBlackLevel                   ID = 0xc61a
	TagWhiteLevel                   ID = 0xc61d
	TagDefaultScale                 ID = 0xc61e
	TagDefaultCropOrigin            ID = 0xc61f
	TagDefaultCropSize              ID = 0xc620
	TagColorMatrix1                 ID = 0xc621
	TagColorMatrix2                 ID = 0xc622
	TagAnalogBalance                ID = 0xc627
	TagAsShotNeutral                ID = 0xc628
	TagBaselineExposure             ID = 0xc62a
	TagBaselineNoise                ID = 0xc62b
	TagBaselineSharpness            ID = 0xc62c
	TagBayerGreenSplit              ID = 0xc62d
	TagLinearResponseLimit          ID = 0xc62e
	TagCameraSerial                 ID = 0xc62f
	TagAntiAliasStrength            ID = 0xc632
	TagShadowScale                  ID = 0xc633
	TagCalibrationIlluminant1       ID = 0xc65a
	TagCalibrationIlluminant2       ID = 0xc65b
	TagBestQualityScale             ID = 0xc65c
	TagOriginalRawFileName          ID = 0xc68b
	TagActiveArea                   ID = 0xc68d
	TagProfileName                  ID = 0xc6f8
	TagProfileEmbedPolicy           ID = 0xc6fd
	TagOpcodeList2                  ID = 0xc741
	TagNoiseProfile                 ID = 0xc761
	TagNewRawImageDigest            ID = 0xc7a7
	TagDefaultUserCrop              ID = 0xc7b5
	TagPrintIM                      ID = 0xc4a5
	TagSR2Private                   ID = 0xc634
	TagPanasonicTitle               ID = 0xc6d2
	TagPanasonicTitle2              ID = 0xc6d3
	TagPanasonicRawVersion          ID = 0x0001
	TagPanasonicSensorWidth         ID = 0x0002
	TagPanasonicSensorHeight        ID = 0x0003
	TagPanasonicSensorTopBorder     ID = 0x0004
	TagPanasonicSensorLeftBorder    ID = 0x0005
	TagPanasonicSensorBottomBorder  ID = 0x0006
	TagPanasonicSensorRightBorder   ID = 0x0007
	TagPanasonicSamplesPerPixel     ID = 0x0008
	TagPanasonicCFAPattern          ID = 0x0009
	TagPanasonicBitsPerSample       ID = 0x000a
	TagPanasonicCompression         ID = 0x000b
	TagPanasonicUnknown000D         ID = 0x000d
	TagPanasonicLinearityLimitRed   ID = 0x000e
	TagPanasonicLinearityLimitGreen ID = 0x000f
	TagPanasonicLinearityLimitBlue  ID = 0x0010
	TagPanasonicISO                 ID = 0x0017
	TagPanasonicHighISOMultRed      ID = 0x0018
	TagPanasonicHighISOMultGreen    ID = 0x0019
	TagPanasonicHighISOMultBlue     ID = 0x001a
	TagNoiseReductionParams         ID = 0x001b
	TagPanasonicBlackLevelRed       ID = 0x001c
	TagPanasonicBlackLevelGreen     ID = 0x001d
	TagPanasonicBlackLevelBlue      ID = 0x001e
	TagPanasonicWBRedLevel          ID = 0x0024
	TagPanasonicWBGreenLevel        ID = 0x0025
	TagPanasonicWBBlueLevel         ID = 0x0026
	TagWBInfo2                      ID = 0x0027
	TagPanasonicRawFormat           ID = 0x002d
	TagJpgFromRaw                   ID = 0x002e
	TagPanasonicISOHighPrecision    ID = 0x0037
	TagPanasonicRawDataOffset       ID = 0x0118
	TagPanasonicDistortionInfo      ID = 0x0119
	TagPanasonicUnknown011A         ID = 0x011a
	TagPanasonicCropTop             ID = 0x0121
	TagPanasonicCropLeft            ID = 0x0122
	TagPanasonicCropBottom          ID = 0x0123
	TagPanasonicCropRight           ID = 0x0124

	// Exif IFD tags (core set).
	TagExposureTime             ID = 0x829a
	TagFNumber                  ID = 0x829d
	TagExposureProgram          ID = 0x8822
	TagISOSpeedRatings          ID = 0x8827
	TagSensitivityType          ID = 0x8830
	TagRecommendedExposureIndex ID = 0x8832
	TagExifVersion              ID = 0x9000
	TagDateTimeOriginal         ID = 0x9003
	TagDateTimeDigitized        ID = 0x9004
	TagOffsetTime               ID = 0x9010
	TagOffsetTimeOriginal       ID = 0x9011
	TagOffsetTimeDigitized      ID = 0x9012
	TagComponentsConfiguration  ID = 0x9101
	TagCompressedBitsPerPixel   ID = 0x9102
	TagShutterSpeedValue        ID = 0x9201
	TagApertureValue            ID = 0x9202
	TagBrightnessValue          ID = 0x9203
	TagExposureBiasValue        ID = 0x9204
	TagMaxApertureValue         ID = 0x9205
	TagSubjectDistance          ID = 0x9206
	TagMeteringMode             ID = 0x9207
	TagLightSource              ID = 0x9208
	TagFlash                    ID = 0x9209
	TagFocalLength              ID = 0x920a
	TagSubjectArea              ID = 0x9214
	TagTIFFEPStandardID         ID = 0x9216
	TagMakerNote                ID = 0x927c
	TagUserComment              ID = 0x9286
	TagSubSecTime               ID = 0x9290
	TagSubSecTimeOriginal       ID = 0x9291
	TagSubSecTimeDigitized      ID = 0x9292
	TagFlashpixVersion          ID = 0xa000
	TagColorSpace               ID = 0xa001
	TagPixelXDimension          ID = 0xa002
	TagPixelYDimension          ID = 0xa003
	TagInteropIFDPointer        ID = 0xa005
	TagFocalPlaneXResolution    ID = 0xa20e
	TagFocalPlaneYResolution    ID = 0xa20f
	TagFocalPlaneResolutionUnit ID = 0xa210
	TagExposureIndex            ID = 0xa215
	TagSensingMethod            ID = 0xa217
	TagFileSource               ID = 0xa300
	TagSceneType                ID = 0xa301
	TagCFAPattern               ID = 0xa302
	TagCustomRendered           ID = 0xa401
	TagExposureMode             ID = 0xa402
	TagWhiteBalance             ID = 0xa403
	TagDigitalZoomRatio         ID = 0xa404
	TagFocalLengthIn35mmFilm    ID = 0xa405
	TagSceneCaptureType         ID = 0xa406
	TagGainControl              ID = 0xa407
	TagContrast                 ID = 0xa408
	TagSaturation               ID = 0xa409
	TagSharpness                ID = 0xa40a
	TagDeviceSettingDescription ID = 0xa40b
	TagSubjectDistanceRange     ID = 0xa40c
	TagCompositeImage           ID = 0xa460
	TagCameraOwnerName          ID = 0xa430
	TagBodySerialNumber         ID = 0xa431
	TagLensSpecification        ID = 0xa432
	TagLensMake                 ID = 0xa433
	TagLensModel                ID = 0xa434
	TagLensSerialNumber         ID = 0xa435

	// GPS IFD tags (core set).
	TagGPSVersionID         ID = 0x0000
	TagGPSLatitudeRef       ID = 0x0001
	TagGPSLatitude          ID = 0x0002
	TagGPSLongitudeRef      ID = 0x0003
	TagGPSLongitude         ID = 0x0004
	TagGPSAltitudeRef       ID = 0x0005
	TagGPSAltitude          ID = 0x0006
	TagGPSTimeStamp         ID = 0x0007
	TagGPSSatellites        ID = 0x0008
	TagGPSStatus            ID = 0x0009
	TagGPSMeasureMode       ID = 0x000a
	TagGPSDOP               ID = 0x000b
	TagGPSSpeedRef          ID = 0x000c
	TagGPSSpeed             ID = 0x000d
	TagGPSTrackRef          ID = 0x000e
	TagGPSTrack             ID = 0x000f
	TagGPSImgDirectionRef   ID = 0x0010
	TagGPSImgDirection      ID = 0x0011
	TagGPSMapDatum          ID = 0x0012
	TagGPSDestLatitudeRef   ID = 0x0013
	TagGPSDestLatitude      ID = 0x0014
	TagGPSDestLongitudeRef  ID = 0x0015
	TagGPSDestLongitude     ID = 0x0016
	TagGPSDestBearingRef    ID = 0x0017
	TagGPSDestBearing       ID = 0x0018
	TagGPSDestDistanceRef   ID = 0x0019
	TagGPSDestDistance      ID = 0x001a
	TagGPSProcessingMethod  ID = 0x001b
	TagGPSAreaInformation   ID = 0x001c
	TagGPSDateStamp         ID = 0x001d
	TagGPSDifferential      ID = 0x001e
	TagGPSHPositioningError ID = 0x001f
)

var rootNames = map[ID]string{
	TagNextIFD:                      "NextIFD",
	TagSubfileType:                  "SubfileType",
	TagImageWidth:                   "ImageWidth",
	TagImageLength:                  "ImageLength",
	TagBitsPerSample:                "BitsPerSample",
	TagCompression:                  "Compression",
	TagPhotometricInterpretation:    "PhotometricInterpretation",
	TagXResolution:                  "XResolution",
	TagYResolution:                  "YResolution",
	TagSamplesPerPixel:              "SamplesPerPixel",
	TagRowsPerStrip:                 "RowsPerStrip",
	TagPlanarConfiguration:          "PlanarConfiguration",
	TagResolutionUnit:               "ResolutionUnit",
	TagThumbnailOffset:              "ThumbnailOffset",
	TagThumbnailLength:              "ThumbnailLength",
	TagStripOffsets:                 "StripOffsets",
	TagStripByteCounts:              "StripByteCounts",
	TagMinSampleValue:               "MinSampleValue",
	TagReferenceBlackWhite:          "ReferenceBlackWhite",
	TagOrientation:                  "Orientation",
	TagSoftware:                     "Software",
	TagDateTime:                     "DateTime",
	TagDateTimeOriginal:             "DateTimeOriginal",
	TagArtist:                       "Artist",
	TagTileWidth:                    "TileWidth",
	TagTileLength:                   "TileLength",
	TagTileOffsets:                  "TileOffsets",
	TagTileByteCounts:               "TileByteCounts",
	TagSubIFDs:                      "SubIFDs",
	TagTIFFEPStandardID:             "TIFF-EPStandardID",
	TagApplicationNotes:             "ApplicationNotes",
	TagCFARepeatPatternDim:          "CFARepeatPatternDim",
	TagCFAPattern2:                  "CFAPattern2",
	TagImageDescription:             "ImageDescription",
	TagMake:                         "Make",
	TagModel:                        "Model",
	TagCopyright:                    "Copyright",
	TagExifIFDPointer:               "ExifOffset",
	TagGPSIFDPointer:                "GPSInfo",
	TagDNGVersion:                   "DNGVersion",
	TagDNGBackwardVersion:           "DNGBackwardVersion",
	TagUniqueCameraModel:            "UniqueCameraModel",
	TagLocalizedCameraModel:         "LocalizedCameraModel",
	TagCFAPlaneColor:                "CFAPlaneColor",
	TagCFALayout:                    "CFALayout",
	TagBlackLevelRepeatDim:          "BlackLevelRepeatDim",
	TagBlackLevel:                   "BlackLevel",
	TagWhiteLevel:                   "WhiteLevel",
	TagDefaultScale:                 "DefaultScale",
	TagDefaultCropOrigin:            "DefaultCropOrigin",
	TagDefaultCropSize:              "DefaultCropSize",
	TagColorMatrix1:                 "ColorMatrix1",
	TagColorMatrix2:                 "ColorMatrix2",
	TagAnalogBalance:                "AnalogBalance",
	TagAsShotNeutral:                "AsShotNeutral",
	TagBaselineExposure:             "BaselineExposure",
	TagBaselineNoise:                "BaselineNoise",
	TagBaselineSharpness:            "BaselineSharpness",
	TagBayerGreenSplit:              "BayerGreenSplit",
	TagLinearResponseLimit:          "LinearResponseLimit",
	TagCameraSerial:                 "CameraSerialNumber",
	TagAntiAliasStrength:            "AntiAliasStrength",
	TagShadowScale:                  "ShadowScale",
	TagCalibrationIlluminant1:       "CalibrationIlluminant1",
	TagCalibrationIlluminant2:       "CalibrationIlluminant2",
	TagBestQualityScale:             "BestQualityScale",
	TagOriginalRawFileName:          "OriginalRawFileName",
	TagActiveArea:                   "ActiveArea",
	TagProfileName:                  "ProfileName",
	TagProfileEmbedPolicy:           "ProfileEmbedPolicy",
	TagOpcodeList2:                  "OpcodeList2",
	TagNoiseProfile:                 "NoiseProfile",
	TagNewRawImageDigest:            "NewRawImageDigest",
	TagDefaultUserCrop:              "DefaultUserCrop",
	TagPrintIM:                      "PrintIM",
	TagSR2Private:                   "SR2Private",
	TagPanasonicTitle:               "PanasonicTitle",
	TagPanasonicTitle2:              "PanasonicTitle2",
	TagPanasonicRawVersion:          "PanasonicRawVersion",
	TagPanasonicSensorWidth:         "SensorWidth",
	TagPanasonicSensorHeight:        "SensorHeight",
	TagPanasonicSensorTopBorder:     "SensorTopBorder",
	TagPanasonicSensorLeftBorder:    "SensorLeftBorder",
	TagPanasonicSensorBottomBorder:  "SensorBottomBorder",
	TagPanasonicSensorRightBorder:   "SensorRightBorder",
	TagPanasonicSamplesPerPixel:     "SamplesPerPixel",
	TagPanasonicCFAPattern:          "CFAPattern",
	TagPanasonicBitsPerSample:       "BitsPerSample",
	TagPanasonicCompression:         "Compression",
	TagPanasonicUnknown000D:         "PanasonicRaw_0x000d",
	TagPanasonicLinearityLimitRed:   "LinearityLimitRed",
	TagPanasonicLinearityLimitGreen: "LinearityLimitGreen",
	TagPanasonicLinearityLimitBlue:  "LinearityLimitBlue",
	TagPanasonicISO:                 "ISO",
	TagPanasonicHighISOMultRed:      "HighISOMultiplierRed",
	TagPanasonicHighISOMultGreen:    "HighISOMultiplierGreen",
	TagPanasonicHighISOMultBlue:     "HighISOMultiplierBlue",
	TagNoiseReductionParams:         "NoiseReductionParams",
	TagPanasonicBlackLevelRed:       "BlackLevelRed",
	TagPanasonicBlackLevelGreen:     "BlackLevelGreen",
	TagPanasonicBlackLevelBlue:      "BlackLevelBlue",
	TagPanasonicWBRedLevel:          "WBRedLevel",
	TagPanasonicWBGreenLevel:        "WBGreenLevel",
	TagPanasonicWBBlueLevel:         "WBBlueLevel",
	TagWBInfo2:                      "WBInfo2",
	TagPanasonicRawFormat:           "RawFormat",
	TagJpgFromRaw:                   "JpgFromRaw",
	TagPanasonicISOHighPrecision:    "ISO",
	TagPanasonicDistortionInfo:      "DistortionInfo",
	TagPanasonicCropTop:             "CropTop",
	TagPanasonicCropLeft:            "CropLeft",
	TagPanasonicCropBottom:          "CropBottom",
	TagPanasonicCropRight:           "CropRight",
}

var exifNames = map[ID]string{
	TagExposureTime:             "ExposureTime",
	TagFNumber:                  "FNumber",
	TagExposureProgram:          "ExposureProgram",
	TagISOSpeedRatings:          "ISOSpeedRatings",
	TagSensitivityType:          "SensitivityType",
	TagRecommendedExposureIndex: "RecommendedExposureIndex",
	TagExifVersion:              "ExifVersion",
	TagDateTimeOriginal:         "DateTimeOriginal",
	TagDateTimeDigitized:        "DateTimeDigitized",
	TagOffsetTime:               "OffsetTime",
	TagOffsetTimeOriginal:       "OffsetTimeOriginal",
	TagOffsetTimeDigitized:      "OffsetTimeDigitized",
	TagComponentsConfiguration:  "ComponentsConfiguration",
	TagCompressedBitsPerPixel:   "CompressedBitsPerPixel",
	TagShutterSpeedValue:        "ShutterSpeedValue",
	TagApertureValue:            "ApertureValue",
	TagBrightnessValue:          "BrightnessValue",
	TagExposureBiasValue:        "ExposureBiasValue",
	TagMaxApertureValue:         "MaxApertureValue",
	TagSubjectDistance:          "SubjectDistance",
	TagLightSource:              "LightSource",
	TagMeteringMode:             "MeteringMode",
	TagFlash:                    "Flash",
	TagFocalLength:              "FocalLength",
	TagSubjectArea:              "SubjectArea",
	TagMakerNote:                "MakerNote",
	TagUserComment:              "UserComment",
	TagSubSecTime:               "SubSecTime",
	TagSubSecTimeOriginal:       "SubSecTimeOriginal",
	TagSubSecTimeDigitized:      "SubSecTimeDigitized",
	TagFlashpixVersion:          "FlashpixVersion",
	TagColorSpace:               "ColorSpace",
	TagPixelXDimension:          "PixelXDimension",
	TagPixelYDimension:          "PixelYDimension",
	TagInteropIFDPointer:        "InteropOffset",
	TagFocalPlaneXResolution:    "FocalPlaneXResolution",
	TagFocalPlaneYResolution:    "FocalPlaneYResolution",
	TagFocalPlaneResolutionUnit: "FocalPlaneResolutionUnit",
	TagExposureIndex:            "ExposureIndex",
	TagSensingMethod:            "SensingMethod",
	TagFileSource:               "FileSource",
	TagSceneType:                "SceneType",
	TagCFAPattern:               "CFAPattern",
	TagCustomRendered:           "CustomRendered",
	TagExposureMode:             "ExposureMode",
	TagWhiteBalance:             "WhiteBalance",
	TagDigitalZoomRatio:         "DigitalZoomRatio",
	TagFocalLengthIn35mmFilm:    "FocalLengthIn35mmFilm",
	TagSceneCaptureType:         "SceneCaptureType",
	TagGainControl:              "GainControl",
	TagContrast:                 "Contrast",
	TagSaturation:               "Saturation",
	TagSharpness:                "Sharpness",
	TagDeviceSettingDescription: "DeviceSettingDescription",
	TagSubjectDistanceRange:     "SubjectDistanceRange",
	TagCompositeImage:           "CompositeImage",
	TagCameraOwnerName:          "CameraOwnerName",
	TagBodySerialNumber:         "BodySerialNumber",
	TagLensSpecification:        "LensSpecification",
	TagLensMake:                 "LensMake",
	TagLensModel:                "LensModel",
	TagLensSerialNumber:         "LensSerialNumber",
}

var gpsNames = map[ID]string{
	TagGPSVersionID:         "GPSVersionID",
	TagGPSLatitudeRef:       "GPSLatitudeRef",
	TagGPSLatitude:          "GPSLatitude",
	TagGPSLongitudeRef:      "GPSLongitudeRef",
	TagGPSLongitude:         "GPSLongitude",
	TagGPSAltitudeRef:       "GPSAltitudeRef",
	TagGPSAltitude:          "GPSAltitude",
	TagGPSTimeStamp:         "GPSTimeStamp",
	TagGPSSatellites:        "GPSSatellites",
	TagGPSStatus:            "GPSStatus",
	TagGPSMeasureMode:       "GPSMeasureMode",
	TagGPSDOP:               "GPSDOP",
	TagGPSSpeedRef:          "GPSSpeedRef",
	TagGPSSpeed:             "GPSSpeed",
	TagGPSTrackRef:          "GPSTrackRef",
	TagGPSTrack:             "GPSTrack",
	TagGPSImgDirectionRef:   "GPSImgDirectionRef",
	TagGPSImgDirection:      "GPSImgDirection",
	TagGPSMapDatum:          "GPSMapDatum",
	TagGPSDestLatitudeRef:   "GPSDestLatitudeRef",
	TagGPSDestLatitude:      "GPSDestLatitude",
	TagGPSDestLongitudeRef:  "GPSDestLongitudeRef",
	TagGPSDestLongitude:     "GPSDestLongitude",
	TagGPSDestBearingRef:    "GPSDestBearingRef",
	TagGPSDestBearing:       "GPSDestBearing",
	TagGPSDestDistanceRef:   "GPSDestDistanceRef",
	TagGPSDestDistance:      "GPSDestDistance",
	TagGPSProcessingMethod:  "GPSProcessingMethod",
	TagGPSAreaInformation:   "GPSAreaInformation",
	TagGPSDateStamp:         "GPSDateStamp",
	TagGPSDifferential:      "GPSDifferential",
	TagGPSHPositioningError: "GPSHPositioningError",
}

// NameFor returns the known tag name for an IFD/tag pair.
func NameFor(directoryType IfdType, id ID) string {
	switch directoryType {
	case IFD0, IFD1, IFD2:
		if name, ok := rootNames[id]; ok {
			return name
		}
	case ExifIFD, SubIFD0, SubIFD1, SubIFD2, SubIFD3, SubIFD4, SubIFD5, SubIFD6, SubIFD7:
		if name, ok := exifNames[id]; ok {
			return name
		}
	case GPSIFD:
		if name, ok := gpsNames[id]; ok {
			return name
		}
	}
	return id.String()
}
