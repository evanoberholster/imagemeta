package canon

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

//go:generate stringer -type=MakerNoteTag -linecomment -output=makernote_tags_string.go

// MakerNoteTag is a Canon maker-note tag ID from ExifTool Canon Main table.
//
// Source:
// - ExifTool 13.25
// - /usr/share/perl5/Image/ExifTool/Canon.pm
// - %Image::ExifTool::Canon::Main
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
	TimeInfo                        MakerNoteTag = 0x0035 // TimeInfo
	BatteryType                     MakerNoteTag = 0x0038 // BatteryType
	AFInfo3                         MakerNoteTag = 0x003c // AFInfo3
	RawDataOffset                   MakerNoteTag = 0x0081 // RawDataOffset
	RawDataLength                   MakerNoteTag = 0x0082 // RawDataLength
	OriginalDecisionDataOffset      MakerNoteTag = 0x0083 // OriginalDecisionDataOffset
	CustomFunctions1D               MakerNoteTag = 0x0090 // CustomFunctions1D
	PersonalFunctions               MakerNoteTag = 0x0091 // PersonalFunctions
	PersonalFunctionValues          MakerNoteTag = 0x0092 // PersonalFunctionValues
	CanonFileInfo                   MakerNoteTag = 0x0093 // CanonFileInfo
	AFPointsInFocus1D               MakerNoteTag = 0x0094 // AFPointsInFocus1D
	LensModel                       MakerNoteTag = 0x0095 // LensModel
	CanonInternalSerialNumber       MakerNoteTag = 0x0096 // InternalSerialNumber
	CanonDustRemovalData            MakerNoteTag = 0x0097 // DustRemovalData
	CanonCropInfo                   MakerNoteTag = 0x0098 // CropInfo
	CanonCustomFunctions            MakerNoteTag = 0x0099 // CustomFunctions2
	CanonAspectInfo                 MakerNoteTag = 0x009a // AspectInfo
	CanonProcessingInfo             MakerNoteTag = 0x00a0 // ProcessingInfo
	CanonColorBalance               MakerNoteTag = 0x00a9 // ColorBalance
	CanonMeasuredColor              MakerNoteTag = 0x00aa // MeasuredColor
	CanonColorTemperature           MakerNoteTag = 0x00ae // ColorTemperature
	CanonCanonFlags                 MakerNoteTag = 0x00b0 // CanonFlags
	CanonModifiedInfo               MakerNoteTag = 0x00b1 // ModifiedInfo
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
	CanonAFMicroAdj                 MakerNoteTag = 0x4013 // AFMicroAdj
	CanonVignettingCorr             MakerNoteTag = 0x4015 // VignettingCorr
	CanonVignettingCorr2            MakerNoteTag = 0x4016 // VignettingCorr2
	CanonLightingOpt                MakerNoteTag = 0x4018 // LightingOpt
	CanonLensInfo                   MakerNoteTag = 0x4019 // LensInfo
	CanonAmbienceInfo               MakerNoteTag = 0x4020 // AmbienceInfo
	CanonMultiExp                   MakerNoteTag = 0x4021 // MultiExp
	CanonFilterInfo                 MakerNoteTag = 0x4024 // FilterInfo
	CanonHDRInfo                    MakerNoteTag = 0x4025 // HDRInfo
	CanonLogInfo                    MakerNoteTag = 0x4026 // LogInfo
	CanonAFConfig                   MakerNoteTag = 0x4028 // AFConfig
	CanonRawBurstModeRoll           MakerNoteTag = 0x403f // RawBurstModeRoll
	CanonLevelInfo                  MakerNoteTag = 0x4059 // LevelInfo
)
