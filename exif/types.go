package exif

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/golang/geo/s2"
)

// Exif is an interface representation of Exif Information
type Exif interface {
	// Aperture convenience func. "IFD/Exif" FNumber
	Aperture() (meta.Aperture, error)

	// Artist convenience func. "IFD" Artist
	Artist() (artist string, err error)

	// CameraSerial convenience func. "IFD/Exif" BodySerialNumber
	CameraSerial() (serial string, err error)

	// CameraMake convenience func. "IFD" Make
	CameraMake() (make string)

	// CameraModel convenience func. "IFD" Model
	CameraModel() (model string)

	// Copyright convenience func. "IFD" Copyright
	Copyright() (copyright string, err error)

	// Dimensions convenience func. "IFD" Dimensions
	Dimensions() (dimensions meta.Dimensions)

	// ExposureBias convenience func. "IFD/Exif" ExposureBiasValue
	ExposureBias() (meta.ExposureBias, error)

	// ExposureProgram convenience func. "IFD/Exif" ExposureProgram
	ExposureProgram() (meta.ExposureProgram, error)

	// ExposureMode convenience func. "IFD/Exif" ExposureMode
	ExposureMode() (meta.ExposureMode, error)

	// Flash convenience func. "IFD/Exif" Flash
	Flash() (meta.Flash, error)

	// FocalLength convenience func. "IFD/Exif" FocalLength
	// Lens Focal Length in mm
	FocalLength() (fl meta.FocalLength, err error)

	// FocalLengthIn35mmFilm convenience func. "IFD/Exif" FocalLengthIn35mmFilm
	// Lens Focal Length Equivalent for 35mm sensor in mm
	FocalLengthIn35mmFilm() (fl meta.FocalLength, err error)

	// Rating convenience func. "IFD/Rating" Rating
	Rating() (rating string, err error)

	// ISOSpeed convenience func. "IFD/Exif" ISOSpeed
	ISOSpeed() (iso uint32, err error)

	// LensMake convenience func. "IFD/Exif" LensMake
	LensMake() (make string, err error)

	// LensModel convenience func. "IFD/Exif" LensModel
	LensModel() (model string, err error)

	// LensSerial convenience func. "IFD/Exif" LensSerialNumber
	LensSerial() (serial string, err error)

	// MeteringMode convenience func. "IFD/Exif" MeteringMode
	MeteringMode() (meta.MeteringMode, error)

	// Orientation convenience func. "IFD" Orientation
	Orientation() meta.Orientation

	// ShutterSpeed convenience func. "IFD/Exif" ExposureTime
	ShutterSpeed() (meta.ExposureTime, error)

	// GPSCoords convenience func. "IFD/GPS" Latitude and Longitude
	GPSCoords() (lat float64, lng float64, err error)

	// GPSCellID convenience func. "IFD/GPS" Latitude and Longitude converted to S2 cellID
	GPSCellID() (cellID s2.CellID, err error)

	// DateTime returns a time.Time that corresponds with when it was created.
	// Since EXIF data does not contain any timezone information, you should
	// select a timezone using tz. If tz is nil UTC is assumed.
	DateTime(tz *time.Location) (tm time.Time, err error)

	// ModifyDate returns a time.Time that corresponds with when it was last modified.
	// Since EXIF data does not contain any timezone information, you should
	// select a timezone using tz. If tz is nil UTC is assumed.
	ModifyDate(tz *time.Location) (time.Time, error)

	// GPSDate convenience func. for "IFD/GPS" GPSDateStamp and GPSTimeStamp.
	// Indicates the time as UTC (Coordinated Universal Time).
	// Optionally sets subsecond based on "IFD/Exif" SubSecTimeOriginal.
	// Sets time zone to time.UTC if non-provided.
	GPSDate(tz *time.Location) (t time.Time, err error)

	// GPSAltitude convenience func. for "IFD/GPS" GPSAltitude and GPSAltitudeRef.
	// Altitude is expressed as one RATIONAL value. The reference unit is meters.
	GPSAltitude() (alt float32, err error)

	// ExposureValue convenience func. "IFD/Exif" ShutterSpeedValue
	ExposureValue() (ev float32, err error)

	// CanonCameraSettings convenience func. "IFD/Exif/Makernotes.Canon" CanonCameraSettings
	// Canon Camera Settings from the Makernote
	CanonCameraSettings() (canon.CameraSettings, error)

	// CanonFileInfo convenience func. "IFD/Exif/Makernotes.Canon" CanonFileInfo
	// Canon Camera File Info from the Makernote
	CanonFileInfo() (canon.FileInfo, error)

	// CanonShotInfo convenience func. "IFD/Exif/Makernotes.Canon" CanonShotInfo
	// Canon Camera Shot Info from the Makernote
	CanonShotInfo() (canon.ShotInfo, error)

	// CanonAFInfo convenience func. "IFD/Exif/Makernotes.Canon" CanonAFInfo
	// Canon Camera AutoFocus Information from the Makernote
	CanonAFInfo() (afInfo canon.AFInfo, err error)
}
