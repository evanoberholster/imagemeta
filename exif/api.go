package exif

import (
	"errors"
	"math"
	"time"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/ifds/gpsifd"
	"github.com/evanoberholster/imagemeta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/golang/geo/s2"
)

var (
	// ErrGpsCoordsNotValid means that some part of the geographic data were unparseable.
	ErrGpsCoordsNotValid = errors.New("error GPS coordinates not valid")
	// ErrGPSRationalNotValid means that the rawCoordinates were not long enough.
	ErrGPSRationalNotValid = errors.New("error GPS Coords requires a raw-coordinate with exactly three rationals")
)

// CameraMake convenience func. "IFD" Make
func (e *Data) CameraMake() (make string) {
	return e.make
}

// CameraModel convenience func. "IFD" Model
func (e *Data) CameraModel() (model string) {
	return e.model
}

// Artist convenience func. "IFD" Artist
func (e *Data) Artist() (artist string, err error) {
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.Artist)
	if err != nil {
		return
	}
	return e.ParseASCIIValue(t)
}

// Copyright convenience func. "IFD" Copyright
func (e *Data) Copyright() (copyright string, err error) {
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.Copyright)
	if err != nil {
		return
	}
	return e.ParseASCIIValue(t)
}

// CameraSerial convenience func. "IFD/Exif" BodySerialNumber
func (e *Data) CameraSerial() (serial string, err error) {
	// BodySerialNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.BodySerialNumber)
	if err == nil {
		serial, err = e.ParseASCIIValue(t)
		return
	}

	// CameraSerialNumber
	t, err = e.GetTag(ifds.RootIFD, 0, ifds.CameraSerialNumber)
	if err == nil {
		serial, err = e.ParseASCIIValue(t)
		return
	}

	return
}

// DateTime returns a time.Time that corresponds with when it was created.
func (e *Data) DateTime() (time.Time, error) {
	// "IFD/Exif" DateTimeOriginal
	// "IFD/Exif" SubSecTimeOriginal
	// TODO: "IFD/Exif" OffsetTimeOriginal
	t1, err := e.GetTag(ifds.ExifIFD, 0, exififd.DateTimeOriginal)
	if err == nil {
		t2, _ := e.GetTag(ifds.ExifIFD, 0, exififd.SubSecTimeOriginal)
		return e.ParseTimeStamp(t1, t2, time.UTC)
	}

	// "IFD/Exif" DateTimeDigitized
	// "IFD/Exif" SubSecTimeDigitized
	// TODO: "IFD/Exif" OffsetTimeDigitized
	t1, err = e.GetTag(ifds.ExifIFD, 0, exififd.DateTimeDigitized)
	if err == nil {
		t2, _ := e.GetTag(ifds.ExifIFD, 0, exififd.SubSecTimeDigitized)
		return e.ParseTimeStamp(t1, t2, time.UTC)
	}
	return time.Time{}, ErrEmptyTag
}

// ModifyDate returns a time.Time that corresponds with when it was last modified.
func (e *Data) ModifyDate() (time.Time, error) {
	// "IFD" DateTime
	// "IFD/Exif" SubSecTime
	t1, err := e.GetTag(ifds.RootIFD, 0, ifds.DateTime)
	if err != nil {
		return time.Time{}, err
	}
	return e.ParseTimeStamp(t1, tag.Tag{}, time.UTC)
}

// LensMake convenience func. "IFD/Exif" LensMake
func (e *Data) LensMake() (make string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensMake)
	if err != nil {
		return
	}
	return e.ParseASCIIValue(t)
}

// LensModel convenience func. "IFD/Exif" LensModel
func (e *Data) LensModel() (model string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensModel)
	if err != nil {
		return
	}
	return e.ParseASCIIValue(t)
}

// LensSerial convenience func. "IFD/Exif" LensSerialNumber
func (e *Data) LensSerial() (serial string, err error) {
	// LensSerialNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensSerialNumber)
	if err == nil {
		serial, err = e.ParseASCIIValue(t)
		return
	}
	return
}

// Dimensions convenience func. "IFD" Dimensions
func (e *Data) Dimensions() (dimensions meta.Dimensions, err error) {
	if e.width > 0 && e.height > 0 {
		return meta.NewDimensions(uint32(e.width), uint32(e.height)), nil
	}
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.PixelXDimension)
	if err == nil {
		e.width, err = e.ParseUint16Value(t)
		if err == nil {
			if t, err = e.GetTag(ifds.ExifIFD, 0, exififd.PixelYDimension); err == nil {
				e.height, err = e.ParseUint16Value(t)
				return meta.NewDimensions(uint32(e.width), uint32(e.height)), err
			}
		}
	}

	t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageWidth)
	if err == nil {
		e.width, err = e.ParseUint16Value(t)
		if err == nil {
			if t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageLength); err == nil {
				e.height, err = e.ParseUint16Value(t)
				return meta.NewDimensions(uint32(e.width), uint32(e.height)), err
			}
		}
	}

	return meta.Dimensions(0), ErrEmptyTag
}

// ExposureProgram convenience func. "IFD/Exif" ExposureProgram
func (e *Data) ExposureProgram() (meta.ExposureProgram, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureProgram)
	if err != nil {
		return 0, err
	}
	ep, err := e.ParseUint16Value(t)
	if err != nil {
		return 0, err
	}
	return meta.ExposureProgram(ep), err
}

// ExposureMode convenience func. "IFD/Exif" ExposureMode
func (e *Data) ExposureMode() (meta.ExposureMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureMode)
	if err != nil {
		return 0, err
	}
	em, err := e.ParseUint16Value(t)
	if err != nil {
		return 0, err
	}
	return meta.NewExposureMode(uint8(em)), err
}

// ExposureBias convenience func. "IFD/Exif" ExposureBiasValue
// TODO: Add ExposureBias Function (Incomplete)
func (e *Data) ExposureBias() (meta.ExposureBias, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureBiasValue)
	if err != nil {
		return meta.ExposureBias(0), err
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return meta.ExposureBias(0), err
	}

	return meta.NewExposureBias(int16(n), int16(d)), nil
}

// MeteringMode convenience func. "IFD/Exif" MeteringMode
func (e *Data) MeteringMode() (meta.MeteringMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.MeteringMode)
	if err != nil {
		return 0, err
	}
	mm, err := e.ParseUint16Value(t)
	if err != nil {
		return 0, err
	}
	return meta.NewMeteringMode(uint8(mm)), err
}

// ShutterSpeed convenience func. "IFD/Exif" ExposureTime
func (e *Data) ShutterSpeed() (meta.ShutterSpeed, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureTime)
	if err != nil {
		return meta.ShutterSpeed{}, err
	}
	num, denom, err := e.ParseRationalValue(t)
	if err != nil {
		return meta.ShutterSpeed{}, err
	}
	return meta.NewShutterSpeed(uint16(num), uint16(denom)), err
}

// ExposureValue convenience func. "IFD/Exif" ShutterSpeedValue
func (e *Data) ExposureValue() (ev float32, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ShutterSpeedValue)
	if err != nil {
		return
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return
	}
	tv := -1 * math.Log2(float64(int32(n))/float64(int32(d)))

	t, err = e.GetTag(ifds.ExifIFD, 0, exififd.ApertureValue)
	if err != nil {
		return 0.0, err
	}
	n1, d2, err := e.ParseRationalValue(t)
	if err != nil {
		return
	}
	av := 2 * math.Log2(float64(n1)/float64(d2))
	return float32(av + tv), nil
}

// Aperture convenience func. "IFD/Exif" FNumber
func (e *Data) Aperture() (meta.Aperture, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FNumber)
	if err != nil {
		return meta.Aperture(0), err
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return meta.Aperture(0), err
	}
	return meta.NewAperture(n, d), nil
}

// FocalLength convenience func. "IFD/Exif" FocalLength
// Lens Focal Length in mm
func (e *Data) FocalLength() (fl meta.FocalLength, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLength)
	if err != nil {
		return
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return
	}
	return meta.NewFocalLength(n, d), nil
}

// FocalLengthIn35mmFilm convenience func. "IFD/Exif" FocalLengthIn35mmFilm
// Lens Focal Length Equivalent for 35mm sensor in mm
func (e *Data) FocalLengthIn35mmFilm() (fl meta.FocalLength, err error) {
	// FocalLengthIn35mmFilm
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLengthIn35mmFilm)
	if err != nil {
		return
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return
	}
	return meta.NewFocalLength(n, d), nil
}

// ISOSpeed convenience func. "IFD/Exif" ISOSpeed
func (e *Data) ISOSpeed() (iso uint32, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ISOSpeedRatings)
	if err != nil {
		return 0, err
	}
	i, err := e.ParseUint16Value(t)
	if err != nil {
		return 0, err
	}

	return uint32(i), err
}

// Flash convenience func. "IFD/Exif" Flash
func (e *Data) Flash() (meta.FlashMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.Flash)
	if err != nil {
		return 0, err
	}
	f, err := e.ParseUint16Value(t)
	if err != nil {
		return 0, err
	}
	return meta.NewFlashMode(uint8(f)), err
}

// Orientation convenience func. "IFD" Orientation
// TODO: Add Orientation Function
func (e *Data) Orientation() (string, error) {
	// Orientation
	return "", nil
}

// GPSCoords is a convenience func. that retrieves "IFD/GPS" GPSLatitude and GPSLongitude
func (e *Data) GPSCoords() (lat float64, lng float64, err error) {
	// Ref - "IFD/GPS" GPSLatitudeRef
	t1, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLatitudeRef)
	if err != nil {
		// Error here
		return
	}
	// Latitude - "IFD/GPS" GPSLatitude
	t2, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLatitude)
	if err != nil {
		// Error here
		return
	}
	lat, err = e.ParseGPSCoord(t1, t2)
	if err != nil {
		return
	}

	// Ref - "IFD/GPS" GPSLongitudeRef
	t1, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLongitudeRef)
	if err != nil {
		// Error here
		return
	}
	// Latitude - "IFD/GPS" GPSLongitude
	t2, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLongitude)
	if err != nil {
		// Error here
		return
	}

	lng, err = e.ParseGPSCoord(t1, t2)
	if err != nil {
		return
	}
	return
}

// GPSDate convenience func. for "IFD/GPS" GPSDateStamp and GPSTimeStamp.
// Indicates the time as UTC (Coordinated Universal Time).
// Optionally sets subsecond based on "IFD/Exif" SubSecTimeOriginal.
// Sets time zone to time.UTC if non-provided.
func (e *Data) GPSDate(tz *time.Location) (t time.Time, err error) {
	if tz == nil {
		tz = time.UTC
	}
	ds, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSDateStamp)
	if err != nil {
		return
	}
	ts, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSTimeStamp)
	if err != nil {
		return
	}
	// ignore error for SubSec
	subSec, _ := e.GetTag(ifds.ExifIFD, 0, exififd.SubSecTimeOriginal)
	return e.ParseGPSTimeStamp(ds, ts, subSec, tz)
}

// GPSAltitude convenience func. for "IFD/GPS" GPSAltitude and GPSAltitudeRef.
// Altitude is expressed as one RATIONAL value. The reference unit is meters.
func (e *Data) GPSAltitude() (alt float32, err error) {
	t, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSAltitude)
	if err != nil {
		return
	}
	n, d, err := e.ParseRationalValue(t)
	if err != nil {
		return
	}
	alt = float32(n) / float32(d)

	t, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSAltitudeRef)
	if t.Type() == tag.TypeByte && t.IsEmbedded() {
		e.er.byteOrder.PutUint32(e.er.rawBuffer[:4], t.ValueOffset)
		if e.er.rawBuffer[0] == 1 {
			alt *= -1
		}
	}

	return alt, err
}

// GPSCellID returns the S2 cellID of the geographic location on the earth.
// A convenience func. that retrieves "IFD/GPS" GPSLatitude and GPSLongitude
// and converts them into an S2 CellID and returns the CellID.
//
// If the CellID is not valid it returns ErrGpsCoordsNotValid.
func (e *Data) GPSCellID() (cellID s2.CellID, err error) {
	lat, lng, err := e.GPSCoords()
	if err != nil {
		return
	}

	latLng := s2.LatLngFromDegrees(lat, lng)
	cellID = s2.CellIDFromLatLng(latLng)

	if !cellID.IsValid() {
		err = ErrGpsCoordsNotValid
		return
	}

	return cellID, nil
}
