package exif

import "github.com/evanoberholster/imagemeta/meta"

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
	ShutterSpeed() (meta.ShutterSpeed, error)
}
