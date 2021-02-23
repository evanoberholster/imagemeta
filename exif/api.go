package exif

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/meta"
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
	return t.ASCIIValue(e.exifReader)
}

// Copyright convenience func. "IFD" Copyright
func (e *Data) Copyright() (copyright string, err error) {
	t, err := e.GetTag(ifds.RootIFD, 0, ifds.Copyright)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// CameraSerial convenience func. "IFD/Exif" BodySerialNumber
func (e *Data) CameraSerial() (serial string, err error) {
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
func (e *Data) LensMake() (make string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensMake)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// LensModel convenience func. "IFD/Exif" LensModel
func (e *Data) LensModel() (model string, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensModel)
	if err != nil {
		return
	}
	return t.ASCIIValue(e.exifReader)
}

// LensSerial convenience func. "IFD/Exif" LensSerialNumber
func (e *Data) LensSerial() (serial string, err error) {
	// LensSerialNumber
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.LensSerialNumber)
	if err == nil {
		serial, err = t.ASCIIValue(e.exifReader)
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
		e.width, err = t.Uint16Value(e.exifReader)
		if err == nil {
			if t, err = e.GetTag(ifds.ExifIFD, 0, exififd.PixelYDimension); err == nil {
				e.height, err = t.Uint16Value(e.exifReader)
				return meta.NewDimensions(uint32(e.width), uint32(e.height)), err
			}
		}
	}

	t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageWidth)
	if err == nil {
		e.width, err = t.Uint16Value(e.exifReader)
		if err == nil {
			if t, err = e.GetTag(ifds.RootIFD, 0, ifds.ImageLength); err == nil {
				e.height, err = t.Uint16Value(e.exifReader)
				return meta.NewDimensions(uint32(e.width), uint32(e.height)), err
			}
		}
	}

	return meta.Dimensions(0), ErrEmptyTag
}

// XMLPacket convenience func. that returns XMP metadata
// from a JPEG image or XMP Packet from "IFD" XMLPacket.
// Whichever is present.
func (e *Data) XMLPacket() (str string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	//if len(e.XMP) > 0 {
	//	str = strings.Replace(string(e.XMP), "\n", "", -1)
	//	return strings.Replace(str, "   ", "", -1), nil
	//	//return xmlfmt.FormatXML(e.XMP, "\t", "  "), nil
	//}
	//
	//t, err := e.GetTag(ifds.RootIFD, 0, ifds.XMLPacket)
	//if err != nil {
	//	return
	//}
	//str, err = t.ASCIIValue(e.exifReader)
	//str = strings.Replace(str, "\n", "", -1)
	//return strings.Replace(str, "   ", "", -1), nil
	//return xmlfmt.FormatXML(str, "\t", "  "), nil
	return
}

// ExposureProgram convenience func. "IFD/Exif" ExposureProgram
func (e *Data) ExposureProgram() (meta.ExposureProgram, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ExposureProgram)
	if err != nil {
		return 0, err
	}
	ep, err := t.Uint16Value(e.exifReader)
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
	em, err := t.Uint16Value(e.exifReader)
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
	_, err = t.RationalValues(e.exifReader)
	if err != nil {
		return meta.ExposureBias(0), err
	}

	return meta.NewExposureBias(0, 0), nil
}

// MeteringMode convenience func. "IFD/Exif" MeteringMode
func (e *Data) MeteringMode() (meta.MeteringMode, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.MeteringMode)
	if err != nil {
		return 0, err
	}
	mm, err := t.Uint16Value(e.exifReader)
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

	ss, err := t.RationalValues(e.exifReader)
	if err != nil {
		return meta.ShutterSpeed{}, err
	}
	return meta.NewShutterSpeed(uint16(ss[0].Numerator), uint16(ss[0].Denominator)), err
}

// Aperture convenience func. "IFD/Exif" FNumber
func (e *Data) Aperture() (meta.Aperture, error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FNumber)
	if err != nil {
		return meta.Aperture(0), err
	}

	ap, err := t.RationalValues(e.exifReader)
	if err != nil {
		return meta.Aperture(0), err
	}
	return meta.NewAperture(ap[0].Numerator, ap[0].Denominator), nil
}

// FocalLength convenience func. "IFD/Exif" FocalLength
// Lens Focal Length in mm
func (e *Data) FocalLength() (fl meta.FocalLength, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLength)
	if err == nil {
		rats, err := t.RationalValues(e.exifReader)
		if err == nil {
			fl = meta.NewFocalLength(rats[0].Numerator, rats[0].Denominator)
			if fl > 0.0 {
				return fl, nil
			}
		}
	}
	return meta.FocalLength(0), ErrEmptyTag
}

// FocalLengthIn35mmFilm convenience func. "IFD/Exif" FocalLengthIn35mmFilm
// Lens Focal Length Equivalent for 35mm sensor in mm
func (e *Data) FocalLengthIn35mmFilm() (fl meta.FocalLength, err error) {
	// FocalLengthIn35mmFilm
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.FocalLengthIn35mmFilm)
	if err == nil {
		rats, err := t.RationalValues(e.exifReader)
		if err == nil {
			fl = meta.NewFocalLength(rats[0].Numerator, rats[0].Denominator)
			if fl > 0.0 {
				return fl, nil
			}
		}
	}
	return meta.FocalLength(0), ErrEmptyTag
}

// ISOSpeed convenience func. "IFD/Exif" ISOSpeed
func (e *Data) ISOSpeed() (iso uint32, err error) {
	t, err := e.GetTag(ifds.ExifIFD, 0, exififd.ISOSpeedRatings)
	if err != nil {
		return 0, err
	}
	fmt.Println(t)
	i, err := t.Uint16Value(e.exifReader)
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
	f, err := t.Uint16Value(e.exifReader)
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
