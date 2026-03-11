package makernote

import metacanon "github.com/evanoberholster/imagemeta/meta/canon"

// Canon contains selected Canon maker-note fields.
//
// Parsing lives in meta/exif/canon.go. This package keeps the parsed container
// model shared by the Exif result.
type Canon struct {
	ImageType            string
	FirmwareVersion      string
	OwnerName            string
	ImageUniqueID        string
	LensModel            string
	InternalSerialNumber string
	BatteryType          string

	FileNumber   uint32
	SerialNumber uint32
	ModelID      uint32
	ColorSpace   uint16

	// Structured Canon maker-note tables (ExifTool Canon.pm mappings).
	CameraSettings             metacanon.CameraSettings
	FocalLength                metacanon.FocalLengthInfo
	FlashInfo                  metacanon.FlashInfo
	ShotInfo                   metacanon.ShotInfo
	FileInfo                   metacanon.FileInfo
	TimeInfo                   metacanon.CanonTimeInfo
	AFInfo                     metacanon.AFInfo
	FaceDetect1                metacanon.FaceDetect1Info
	FaceDetect2                metacanon.FaceDetect2Info
	FaceDetect3                metacanon.FaceDetect3Info
	AspectInfo                 metacanon.AspectInfo
	ProcessingInfo             metacanon.ProcessingInfo
	CustomPictureStyleFileName string
	AFMicroAdj                 metacanon.AFMicroAdjInfo
	LensInfo                   metacanon.LensInfoForService
	MultiExp                   metacanon.MultiExpInfo
	HDRInfo                    metacanon.HDRInfo
	PreviewImageInfo           metacanon.PreviewImageInfo
	SensorInfo                 metacanon.SensorInfo
	AFConfig                   metacanon.AFConfig
	RawBurstModeRoll           metacanon.RawBurstInfo
	LightingOpt                metacanon.LightingOptInfo
}
