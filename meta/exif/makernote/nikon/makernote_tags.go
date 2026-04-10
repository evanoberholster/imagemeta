package nikon

import "github.com/evanoberholster/imagemeta/meta/exif/tag"

// MakerNoteTag is a Nikon maker-note tag ID from ExifTool's Nikon Main table.
//
// Sources:
// - ExifTool Nikon tag names: https://www.exiftool.org/TagNames/Nikon.html
// - /usr/share/perl5/Image/ExifTool/Nikon.pm
// - %Image::ExifTool::Nikon::Main
//
// ExifTool conditionally assigns multiple public names to some Nikon maker-note
// tag IDs depending on payload version or model. For those cases imagemeta
// keeps a canonical tag constant plus selected aliases that map to the same tag
// ID.
type MakerNoteTag tag.ID

// Nikon maker-note Main-table tag IDs.
const (
	MakerNoteVersion          MakerNoteTag = 0x0001 // MakerNoteVersion
	ISO                       MakerNoteTag = 0x0002 // ISO
	ColorMode                 MakerNoteTag = 0x0003 // ColorMode
	Quality                   MakerNoteTag = 0x0004 // Quality
	WhiteBalance              MakerNoteTag = 0x0005 // WhiteBalance
	Sharpness                 MakerNoteTag = 0x0006 // Sharpness
	FocusMode                 MakerNoteTag = 0x0007 // FocusMode
	FlashSetting              MakerNoteTag = 0x0008 // FlashSetting
	FlashType                 MakerNoteTag = 0x0009 // FlashType
	WhiteBalanceFineTune      MakerNoteTag = 0x000b // WhiteBalanceFineTune
	WBRBLevels                MakerNoteTag = 0x000c // WB_RBLevels
	ProgramShift              MakerNoteTag = 0x000d // ProgramShift
	ExposureDifference        MakerNoteTag = 0x000e // ExposureDifference
	ISOSelection              MakerNoteTag = 0x000f // ISOSelection
	DataDump                  MakerNoteTag = 0x0010 // DataDump
	PreviewIFD                MakerNoteTag = 0x0011 // PreviewIFD
	FlashExposureComp         MakerNoteTag = 0x0012 // FlashExposureComp
	ISOSetting                MakerNoteTag = 0x0013 // ISOSetting
	ColorBalanceA             MakerNoteTag = 0x0014 // ColorBalanceA
	NRWData                   MakerNoteTag = 0x0014 // NRWData
	ImageBoundary             MakerNoteTag = 0x0016 // ImageBoundary
	ExternalFlashExposureComp MakerNoteTag = 0x0017 // ExternalFlashExposureComp
	FlashExposureBracketValue MakerNoteTag = 0x0018 // FlashExposureBracketValue
	ExposureBracketValue      MakerNoteTag = 0x0019 // ExposureBracketValue
	ImageProcessing           MakerNoteTag = 0x001a // ImageProcessing
	CropHiSpeed               MakerNoteTag = 0x001b // CropHiSpeed
	ExposureTuning            MakerNoteTag = 0x001c // ExposureTuning
	SerialNumber              MakerNoteTag = 0x001d // SerialNumber
	ColorSpace                MakerNoteTag = 0x001e // ColorSpace
	VRInfo                    MakerNoteTag = 0x001f // VRInfo
	ImageAuthentication       MakerNoteTag = 0x0020 // ImageAuthentication
	FaceDetect                MakerNoteTag = 0x0021 // FaceDetect
	ActiveDLighting           MakerNoteTag = 0x0022 // ActiveD-Lighting
	PictureControlData        MakerNoteTag = 0x0023 // PictureControlData
	WorldTime                 MakerNoteTag = 0x0024 // WorldTime
	ISOInfo                   MakerNoteTag = 0x0025 // ISOInfo
	VignetteControl           MakerNoteTag = 0x002a // VignetteControl
	DistortInfo               MakerNoteTag = 0x002b // DistortInfo
	UnknownInfo               MakerNoteTag = 0x002c // UnknownInfo
	UnknownInfo2              MakerNoteTag = 0x0032 // UnknownInfo2
	ShutterMode               MakerNoteTag = 0x0034 // ShutterMode
	HDRInfo                   MakerNoteTag = 0x0035 // HDRInfo
	HDRInfo2                  MakerNoteTag = 0x0035 // HDRInfo2
	MechanicalShutterCount    MakerNoteTag = 0x0037 // MechanicalShutterCount
	LocationInfo              MakerNoteTag = 0x0039 // LocationInfo
	BlackLevel                MakerNoteTag = 0x003d // BlackLevel
	ImageSizeRAW              MakerNoteTag = 0x003e // ImageSizeRAW
	WhiteBalanceFineTune2     MakerNoteTag = 0x003f // WhiteBalanceFineTune
	JPGCompression            MakerNoteTag = 0x0044 // JPGCompression
	CropArea                  MakerNoteTag = 0x0045 // CropArea
	NikonSettings             MakerNoteTag = 0x004e // NikonSettings
	ColorTemperatureAuto      MakerNoteTag = 0x004f // ColorTemperatureAuto
	MakerNotes0x51            MakerNoteTag = 0x0051 // MakerNotes0x51
	MakerNotes0x56            MakerNoteTag = 0x0056 // MakerNotes0x56

	ImageAdjustment       MakerNoteTag = 0x0080 // ImageAdjustment
	ToneComp              MakerNoteTag = 0x0081 // ToneComp
	AuxiliaryLens         MakerNoteTag = 0x0082 // AuxiliaryLens
	LensType              MakerNoteTag = 0x0083 // LensType
	Lens                  MakerNoteTag = 0x0084 // Lens
	ManualFocusDistance   MakerNoteTag = 0x0085 // ManualFocusDistance
	DigitalZoom           MakerNoteTag = 0x0086 // DigitalZoom
	FlashMode             MakerNoteTag = 0x0087 // FlashMode
	AFInfo                MakerNoteTag = 0x0088 // AFInfo
	ShootingMode          MakerNoteTag = 0x0089 // ShootingMode
	LensFStops            MakerNoteTag = 0x008b // LensFStops
	ContrastCurve         MakerNoteTag = 0x008c // ContrastCurve
	ColorHue              MakerNoteTag = 0x008d // ColorHue
	SceneMode             MakerNoteTag = 0x008f // SceneMode
	LightSource           MakerNoteTag = 0x0090 // LightSource
	ShotInfo              MakerNoteTag = 0x0091 // ShotInfo*
	HueAdjustment         MakerNoteTag = 0x0092 // HueAdjustment
	NEFCompression        MakerNoteTag = 0x0093 // NEFCompression
	SaturationAdj         MakerNoteTag = 0x0094 // SaturationAdj
	NoiseReduction        MakerNoteTag = 0x0095 // NoiseReduction
	NEFLinearizationTable MakerNoteTag = 0x0096 // NEFLinearizationTable
	ColorBalance          MakerNoteTag = 0x0097 // ColorBalance*
	LensData              MakerNoteTag = 0x0098 // LensData*
	RawImageCenter        MakerNoteTag = 0x0099 // RawImageCenter
	SensorPixelSize       MakerNoteTag = 0x009a // SensorPixelSize
	SceneAssist           MakerNoteTag = 0x009c // SceneAssist
	DateStampMode         MakerNoteTag = 0x009d // DateStampMode
	RetouchHistory        MakerNoteTag = 0x009e // RetouchHistory
	SerialNumber2         MakerNoteTag = 0x00a0 // SerialNumber
	ImageDataSize         MakerNoteTag = 0x00a2 // ImageDataSize
	ImageCount            MakerNoteTag = 0x00a5 // ImageCount
	DeletedImageCount     MakerNoteTag = 0x00a6 // DeletedImageCount
	ShutterCount          MakerNoteTag = 0x00a7 // ShutterCount
	FlashInfo             MakerNoteTag = 0x00a8 // FlashInfo*
	ImageOptimization     MakerNoteTag = 0x00a9 // ImageOptimization
	Saturation            MakerNoteTag = 0x00aa // Saturation
	VariProgram           MakerNoteTag = 0x00ab // VariProgram
	ImageStabilization    MakerNoteTag = 0x00ac // ImageStabilization
	AFResponse            MakerNoteTag = 0x00ad // AFResponse
	MultiExposure         MakerNoteTag = 0x00b0 // MultiExposure
	MultiExposure2        MakerNoteTag = 0x00b0 // MultiExposure2
	HighISONoiseReduction MakerNoteTag = 0x00b1 // HighISONoiseReduction
	ToningEffect          MakerNoteTag = 0x00b3 // ToningEffect
	PowerUpTime           MakerNoteTag = 0x00b6 // PowerUpTime
	AFInfo2               MakerNoteTag = 0x00b7 // AFInfo2
	FileInfo              MakerNoteTag = 0x00b8 // FileInfo
	AFTune                MakerNoteTag = 0x00b9 // AFTune
	RetouchInfo           MakerNoteTag = 0x00bb // RetouchInfo
	PictureControlData2   MakerNoteTag = 0x00bd // PictureControlData
	SilentPhotography     MakerNoteTag = 0x00bf // SilentPhotography
	BarometerInfo         MakerNoteTag = 0x00c3 // BarometerInfo

	PrintIM                  MakerNoteTag = 0x0e00 // PrintIM
	NikonCaptureData         MakerNoteTag = 0x0e01 // NikonCaptureData
	NikonCaptureVersion      MakerNoteTag = 0x0e09 // NikonCaptureVersion
	NikonCaptureOffsets      MakerNoteTag = 0x0e0e // NikonCaptureOffsets
	NikonScanIFD             MakerNoteTag = 0x0e10 // NikonScanIFD
	NikonCaptureEditVersions MakerNoteTag = 0x0e13 // NikonCaptureEditVersions
	NikonICCProfile          MakerNoteTag = 0x0e1d // NikonICCProfile
	NikonCaptureOutput       MakerNoteTag = 0x0e1e // NikonCaptureOutput
	NEFBitDepth              MakerNoteTag = 0x0e22 // NEFBitDepth
)
