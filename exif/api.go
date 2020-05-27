package exif

import (
	"strings"

	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/ifds/exififd"
)

// CameraMake convenience func. "IFD" Make
func (e *Exif) CameraMake() (make string) {
	return e.Make
}

// CameraModel convenience func. "IFD" Model
func (e *Exif) CameraModel() (model string) {
	return e.Model
}

// Artist convenience func. "IFD" Artist
func (e *Exif) Artist() (artist string, err error) {
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.Artist)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// Copyright convenience func. "IFD" Copyright
func (e *Exif) Copyright() (copyright string, err error) {
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.Copyright)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// CameraSerial convenience func. "IFD/Exif" BodySerialNumber
func (e *Exif) CameraSerial() (serial string, err error) {
	// BodySerialNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.BodySerialNumber)
	if err == nil {
		serial, err = t.ASCIIValue(e.exifReader)
		return
	}

	// CameraSerialNumber
	t, err = e.GetTag(ifds.RootIFD, 0, ifds.CameraSerialNumber)
	if err == nil {
		serial, err = t.ASCIIValue(e.exifReader)
		return
	}

	return
}

// LensMake convenience func. "IFD/Exif" LensMake
func (e *Exif) LensMake() (make string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensMake)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// LensModel convenience func. "IFD/Exif" LensModel
func (e *Exif) LensModel() (model string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensModel)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// LensSerial convenience func. "IFD/Exif" LensSerialNumber
func (e *Exif) LensSerial() (serial string, err error) {
	// LensSerialNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensSerialNumber)
	if err == nil {
		serial, err = t.ASCIIValue(e.exifReader)
		return
	}
	return
}

// Dimensions convenience func. "IFD" Dimensions
func (e *Exif) Dimensions() (width, height uint16, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.PixelXDimension)
	if err == nil {
		width, err = t.Uint16Value(e.exifReader)
		if err == nil {
			if t, err = e.GetTag(ifds.ExifIFD, 0, exififd.PixelYDimension); err == nil {
				height, err = t.Uint16Value(e.exifReader)
				return
			}
		}
	}

	t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageWidth)
	if err == nil {
		width, err = t.Uint16Value(e.exifReader)
		if err == nil {
			if t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageLength); err == nil {
				height, err = t.Uint16Value(e.exifReader)
				return
			}
		}
	}

	return 0, 0, ErrEmptyTag
}

// XMLPacket convenience func. "IFD" XMLPacket
// WIP
func (e *Exif) XMLPacket() (str string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.XMLPacket)
	if err != nil {
		return
	}
	str, err = t.ASCIIValue(e.exifReader)
	str = strings.Replace(str, "\n", "", -1)
	return strings.Replace(str, "   ", "", -1), nil
}

// ExposureProgram convenience func. "IFD/Exif" ExposureProgram
func (e *Exif) ExposureProgram() (ExposureMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureProgram)
	if err != nil {
		return 0, err
	}
	ep, err := t.Uint16Value(e.exifReader)
	if err != nil {
		return 0, err
	}
	return ExposureMode(ep), err
}

// MeteringMode convenience func. "IFD/Exif" MeteringMode
func (e *Exif) MeteringMode() (MeteringMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.MeteringMode)
	if err != nil {
		return 0, err
	}
	mm, err := t.Uint16Value(e.exifReader)
	if err != nil {
		return 0, err
	}
	return MeteringMode(mm), err
}

// ShutterSpeed convenience func. "IFD/Exif" ExposureTime
func (e *Exif) ShutterSpeed() (ShutterSpeed, error) {
	// ShutterSpeedValue
	// ExposureTime
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureTime)
	if err != nil {
		return ShutterSpeed{0, 0}, err
	}

	ss, err := t.RationalValues(e.exifReader)
	if err != nil {
		return ShutterSpeed{0, 0}, err
	}
	return ShutterSpeed{ss[0].Numerator, ss[0].Denominator}, err
}

// Aperture convenience func. "IFD/Exif" FNumber
func (e *Exif) Aperture() (float32, error) {
	// ApertureValue
	// FNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FNumber)
	if err != nil {
		return 0.0, err
	}

	ap, err := t.RationalValues(e.exifReader)
	if err != nil {
		return 0.0, err
	}
	return float32(ap[0].Numerator) / float32(ap[0].Denominator), nil
}

// FocalLength convenience func. "IFD/Exif" FocalLength
// Lens Focal Length in mm
func (e *Exif) FocalLength() (fl FocalLength, err error) {
	// FocalLength
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLength)
	if err == nil {
		rats, err := t.RationalValues(e.exifReader)
		if err == nil {
			fl = FocalLength(float32(rats[0].Numerator) / float32(rats[0].Denominator))
			if fl > 0.0 {
				return fl, nil
			}
		}
	}
	return 0.0, ErrEmptyTag
}

// FocalLengthIn35mmFilm convenience func. "IFD/Exif" FocalLengthIn35mmFilm
// Lens Focal Length Equivalent for 35mm sensor in mm
func (e *Exif) FocalLengthIn35mmFilm() (fl FocalLength, err error) {
	// FocalLengthIn35mmFilm
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLengthIn35mmFilm)
	if err == nil {
		rats, err := t.RationalValues(e.exifReader)
		if err == nil {
			fl = FocalLength(float32(rats[0].Numerator) / float32(rats[0].Denominator))
			if fl > 0.0 {
				return fl, nil
			}
		}
	}
	return 0.0, ErrEmptyTag
}

// ISOSpeed convenience func. "IFD/Exif" ISOSpeed
func (e *Exif) ISOSpeed() (iso int, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ISOSpeedRatings)
	if err != nil {
		return 0, err
	}
	i, err := t.Uint16Value(e.exifReader)
	if err != nil {
		return 0, err
	}

	return int(i), err
}

// Flash convenience func. "IFD/Exif" Flash
func (e *Exif) Flash() (FlashMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.Flash)
	if err != nil {
		return 0, err
	}
	f, err := t.Uint16Value(e.exifReader)
	if err != nil {
		return 0, err
	}
	return FlashMode(f), err
}

// Orientation convenience func. "IFD" Orientation
// TODO: Add Orientation Function
func (e *Exif) Orientation() (string, error) {
	// Orientation
	return "", nil
}

// ExposureBias convenience func. "IFD/Exif" ExposureBiasValue
// TODO: Add ExposureBias Function
func (e *Exif) ExposureBias() (string, error) {
	// ExposureBiasValue
	return "", nil
}
