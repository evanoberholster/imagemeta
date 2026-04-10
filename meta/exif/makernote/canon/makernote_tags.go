package canon

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

//go:generate stringer -type=MakerNoteTag -linecomment -output=makernote_tags_string.go

// MakerNoteTag is a Canon maker-note tag ID from ExifTool Canon Main table.
//
// Source:
// - ExifTool 13.25
// - /usr/share/perl5/Image/ExifTool/Canon.pm
// - %Image::ExifTool::Canon::Main
//
// Unknown Canon_0xNNNN tags below are observed in the local CR2/CR3 corpus and
// included so the Canon maker-note tag list is exhaustive for current samples.

type MakerNoteTag tag.ID

// Canon maker-note Main-table tag IDs.
const (
	CanonCameraSettings             MakerNoteTag = 0x0001 // CanonCameraSettings
	CanonFocalLength                MakerNoteTag = 0x0002 // CanonFocalLength
	CanonFlashInfo                  MakerNoteTag = 0x0003 // CanonFlashInfo
	CanonShotInfo                   MakerNoteTag = 0x0004 // CanonShotInfo
	CanonPanorama                   MakerNoteTag = 0x0005 // CanonPanorama
	CanonImageType                  MakerNoteTag = 0x0006 // CanonImageType
	CanonFirmwareVersion            MakerNoteTag = 0x0007 // CanonFirmwareVersion
	FileNumber                      MakerNoteTag = 0x0008 // FileNumber
	OwnerName                       MakerNoteTag = 0x0009 // OwnerName
	UnknownD30                      MakerNoteTag = 0x000a // UnknownD30
	SerialNumber                    MakerNoteTag = 0x000c // SerialNumber
	CanonCameraInfo                 MakerNoteTag = 0x000d // CanonCameraInfo
	CanonFileLength                 MakerNoteTag = 0x000e // CanonFileLength
	CustomFunctions                 MakerNoteTag = 0x000f // CustomFunctions
	CanonModelID                    MakerNoteTag = 0x0010 // CanonModelID
	MovieInfo                       MakerNoteTag = 0x0011 // MovieInfo
	CanonAFInfo                     MakerNoteTag = 0x0012 // CanonAFInfo
	ThumbnailImageValidArea         MakerNoteTag = 0x0013 // ThumbnailImageValidArea
	SerialNumberFormat              MakerNoteTag = 0x0015 // SerialNumberFormat
	Canon0x0019                     MakerNoteTag = 0x0019 // Canon_0x0019
	SuperMacro                      MakerNoteTag = 0x001a // SuperMacro
	DateStampMode                   MakerNoteTag = 0x001c // DateStampMode
	MyColors                        MakerNoteTag = 0x001d // MyColors
	FirmwareRevision                MakerNoteTag = 0x001e // FirmwareRevision
	Categories                      MakerNoteTag = 0x0023 // Categories
	FaceDetect1                     MakerNoteTag = 0x0024 // FaceDetect1
	FaceDetect2                     MakerNoteTag = 0x0025 // FaceDetect2
	CanonAFInfo2                    MakerNoteTag = 0x0026 // CanonAFInfo2
	ContrastInfo                    MakerNoteTag = 0x0027 // ContrastInfo
	ImageUniqueID                   MakerNoteTag = 0x0028 // ImageUniqueID
	WBInfo                          MakerNoteTag = 0x0029 // WBInfo
	FaceDetect3                     MakerNoteTag = 0x002f // FaceDetect3
	Canon0x0032                     MakerNoteTag = 0x0032 // Canon_0x0032
	Canon0x0033                     MakerNoteTag = 0x0033 // Canon_0x0033
	TimeInfo                        MakerNoteTag = 0x0035 // TimeInfo
	BatteryType                     MakerNoteTag = 0x0038 // BatteryType
	AFInfo3                         MakerNoteTag = 0x003c // AFInfo3
	Canon0x003d                     MakerNoteTag = 0x003d // Canon_0x003d
	Canon0x003f                     MakerNoteTag = 0x003f // Canon_0x003f
	RawDataOffset                   MakerNoteTag = 0x0081 // RawDataOffset
	RawDataLength                   MakerNoteTag = 0x0082 // RawDataLength
	OriginalDecisionDataOffset      MakerNoteTag = 0x0083 // OriginalDecisionDataOffset
	CustomFunctions1D               MakerNoteTag = 0x0090 // CustomFunctions1D
	PersonalFunctions               MakerNoteTag = 0x0091 // PersonalFunctions
	PersonalFunctionValues          MakerNoteTag = 0x0092 // PersonalFunctionValues
	CanonFileInfo                   MakerNoteTag = 0x0093 // CanonFileInfo
	AFPointsInFocus1D               MakerNoteTag = 0x0094 // AFPointsInFocus1D
	LensModel                       MakerNoteTag = 0x0095 // LensModel
	CanonSerialInfo                 MakerNoteTag = 0x0096 // SerialInfo
	CanonInternalSerialNumber       MakerNoteTag = 0x0096 // InternalSerialNumber
	CanonDustRemovalData            MakerNoteTag = 0x0097 // DustRemovalData
	CanonCropInfo                   MakerNoteTag = 0x0098 // CropInfo
	CanonCustomFunctions            MakerNoteTag = 0x0099 // CustomFunctions2
	CanonAspectInfo                 MakerNoteTag = 0x009a // AspectInfo
	CanonProcessingInfo             MakerNoteTag = 0x00a0 // ProcessingInfo
	CanonToneCurveTable             MakerNoteTag = 0x00a1 // ToneCurveTable
	CanonSharpnessTable             MakerNoteTag = 0x00a2 // SharpnessTable
	CanonSharpnessFreqTable         MakerNoteTag = 0x00a3 // SharpnessFreqTable
	CanonWhiteBalanceTable          MakerNoteTag = 0x00a4 // WhiteBalanceTable
	CanonColorBalance               MakerNoteTag = 0x00a9 // ColorBalance
	CanonMeasuredColor              MakerNoteTag = 0x00aa // MeasuredColor
	CanonColorTemperature           MakerNoteTag = 0x00ae // ColorTemperature
	CanonCanonFlags                 MakerNoteTag = 0x00b0 // CanonFlags
	CanonModifiedInfo               MakerNoteTag = 0x00b1 // ModifiedInfo
	CanonToneCurveMatching          MakerNoteTag = 0x00b2 // ToneCurveMatching
	CanonWhiteBalanceMatching       MakerNoteTag = 0x00b3 // WhiteBalanceMatching
	CanonColorSpace                 MakerNoteTag = 0x00b4 // ColorSpace
	CanonPreviewImageInfo           MakerNoteTag = 0x00b6 // PreviewImageInfo
	CanonVRDOffset                  MakerNoteTag = 0x00d0 // VRDOffset
	CanonSensorInfo                 MakerNoteTag = 0x00e0 // SensorInfo
	CanonColorData                  MakerNoteTag = 0x4001 // ColorData
	CanonCRWParam                   MakerNoteTag = 0x4002 // CRWParam
	CanonColorInfo                  MakerNoteTag = 0x4003 // ColorInfo
	CanonFlavor                     MakerNoteTag = 0x4005 // Flavor
	CanonPictureStyleUserDef        MakerNoteTag = 0x4008 // PictureStyleUserDef
	CanonPictureStylePC             MakerNoteTag = 0x4009 // PictureStylePC
	CanonCustomPictureStyleFileName MakerNoteTag = 0x4010 // CustomPictureStyleFileName
	Canon0x4011                     MakerNoteTag = 0x4011 // Canon_0x4011
	Canon0x4012                     MakerNoteTag = 0x4012 // Canon_0x4012
	CanonAFMicroAdj                 MakerNoteTag = 0x4013 // AFMicroAdj
	Canon0x4014                     MakerNoteTag = 0x4014 // Canon_0x4014
	CanonVignettingCorr             MakerNoteTag = 0x4015 // VignettingCorr
	CanonVignettingCorr2            MakerNoteTag = 0x4016 // VignettingCorr2
	Canon0x4017                     MakerNoteTag = 0x4017 // Canon_0x4017
	CanonLightingOpt                MakerNoteTag = 0x4018 // LightingOpt
	CanonLensInfo                   MakerNoteTag = 0x4019 // LensInfo
	CanonAmbienceInfo               MakerNoteTag = 0x4020 // AmbienceInfo
	CanonMultiExp                   MakerNoteTag = 0x4021 // MultiExp
	CanonFilterInfo                 MakerNoteTag = 0x4024 // FilterInfo
	CanonHDRInfo                    MakerNoteTag = 0x4025 // HDRInfo
	CanonLogInfo                    MakerNoteTag = 0x4026 // LogInfo
	Canon0x4027                     MakerNoteTag = 0x4027 // Canon_0x4027
	CanonAFConfig                   MakerNoteTag = 0x4028 // AFConfig
	Canon0x402c                     MakerNoteTag = 0x402c // Canon_0x402c
	Canon0x402e                     MakerNoteTag = 0x402e // Canon_0x402e
	Canon0x4031                     MakerNoteTag = 0x4031 // Canon_0x4031
	Canon0x4033                     MakerNoteTag = 0x4033 // Canon_0x4033
	Canon0x4035                     MakerNoteTag = 0x4035 // Canon_0x4035
	Canon0x4037                     MakerNoteTag = 0x4037 // Canon_0x4037
	Canon0x4039                     MakerNoteTag = 0x4039 // Canon_0x4039
	Canon0x403a                     MakerNoteTag = 0x403a // Canon_0x403a
	Canon0x403b                     MakerNoteTag = 0x403b // Canon_0x403b
	Canon0x403c                     MakerNoteTag = 0x403c // Canon_0x403c
	CanonRawBurstModeRoll           MakerNoteTag = 0x403f // RawBurstModeRoll
	Canon0x4040                     MakerNoteTag = 0x4040 // Canon_0x4040
	Canon0x4045                     MakerNoteTag = 0x4045 // Canon_0x4045
	Canon0x4049                     MakerNoteTag = 0x4049 // Canon_0x4049
	Canon0x404a                     MakerNoteTag = 0x404a // Canon_0x404a
	Canon0x404b                     MakerNoteTag = 0x404b // Canon_0x404b
	Canon0x404e                     MakerNoteTag = 0x404e // Canon_0x404e
	Canon0x404f                     MakerNoteTag = 0x404f // Canon_0x404f
	Canon0x4051                     MakerNoteTag = 0x4051 // Canon_0x4051
	Canon0x4053                     MakerNoteTag = 0x4053 // Canon_0x4053
	Canon0x4054                     MakerNoteTag = 0x4054 // Canon_0x4054
	Canon0x4055                     MakerNoteTag = 0x4055 // Canon_0x4055
	Canon0x4056                     MakerNoteTag = 0x4056 // Canon_0x4056
	Canon0x4058                     MakerNoteTag = 0x4058 // Canon_0x4058
	CanonLevelInfo                  MakerNoteTag = 0x4059 // LevelInfo
	Canon0x405b                     MakerNoteTag = 0x405b // Canon_0x405b
)
