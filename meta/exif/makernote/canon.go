package makernote

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
)

// Canon contains selected Canon maker-note fields.
//
// Parsing lives in meta/exif/canon.go. This package keeps the parsed container
// model shared by the Exif result.
type Canon struct {
	ImageType            string    // 16 bytes
	FirmwareVersion      string    // 16 bytes
	OwnerName            string    // 16 bytes
	ImageUniqueID        meta.UUID // 16 bytes
	LensModel            string    // 16 bytes
	InternalSerialNumber string    // 16 bytes
	BatteryType          string    // 16 bytes

	FileNumber   uint32 // 4 bytes
	SerialNumber uint32 // 4 bytes
	ModelID      uint32 // 4 bytes
	ColorSpace   uint16 // 4 bytes in struct (2 data + 2 padding)

	// Structured Canon maker-note tables (ExifTool Canon.pm mappings).
	CanonCameraSettings        canon.CameraSettings     // 88 bytes
	CanonFocalLength           canon.FocalLengthInfo    // 8 bytes
	CanonShotInfo              canon.ShotInfo           // 100 bytes
	CanonFileInfo              canon.FileInfo           // 48 bytes
	TimeInfo                   canon.CanonTimeInfo      // 12 bytes
	AFInfo                     canon.AFInfo             // 128 bytes
	FaceDetect1                canon.FaceDetect1Info    // 42 bytes
	FaceDetect2                canon.FaceDetect2Info    // 2 bytes
	FaceDetect3                canon.FaceDetect3Info    // 4 bytes in struct (2 data + 2 padding)
	AspectInfo                 canon.AspectInfo         // 20 bytes
	ProcessingInfo             canon.ProcessingInfo     // 36 bytes in struct (30 data + 6 padding)
	CustomPictureStyleFileName string                   // 16 bytes
	AFMicroAdj                 canon.AFMicroAdjInfo     // 16 bytes in struct (12 data + 4 padding)
	LensInfo                   canon.LensInfoForService // 24 bytes
	MultiExp                   canon.MultiExpInfo       // 12 bytes
	HDRInfo                    canon.HDRInfo            // 8 bytes
	PreviewImageInfo           canon.PreviewImageInfo   // 20 bytes
	SensorInfo                 canon.SensorInfo         // 20 bytes
	AFConfig                   canon.AFConfig           // 84 bytes
	RawBurstModeRoll           canon.RawBurstInfo       // 8 bytes
	LightingOpt                canon.LightingOptInfo    // 32 bytes in struct (28 data + 4 padding)
}
