package xmp

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// Flash represents exif:Flash
// Based on: https://exiftool.org/TagNames/XMP.html
type Flash struct {
	Fired      bool
	Mode       uint8
	RedEyeMode bool
	Function   bool
	Return     uint8
}

//func (xmpFlash *Flash) parse(p property) (err error) {
//	switch p.Name() {
//	case xmpns.Fired:
//		xmpFlash.Fired = parseBool(p.val)
//	case xmpns.Return:
//		xmpFlash.Return = uint8(parseUint(p.val))
//	case xmpns.Mode:
//		xmpFlash.Mode = uint8(parseUint(p.val))
//	case xmpns.Function:
//		xmpFlash.Function = parseBool(p.val)
//	case xmpns.RedEyeMode:
//		xmpFlash.RedEyeMode = parseBool(p.val)
//	default:
//		return ErrPropertyNotSet
//	}
//	return
//}

// Exif attributes of an XMP Packet.
//
//	Exif 2.21 or later: xmlns:exifEX="http://cipa.jp/exif/1.0/"
//	Exif 2.2 or earlier: xmlns:exif="http://ns.adobe.com/exif/1.0/"
//
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#exif
type Exif struct {
	ExifVersion      string
	PixelXDimension  uint32
	PixelYDimension  uint32
	DateTimeOriginal time.Time
	CreateDate       time.Time // Exif:DateTimeDigitized
	ExposureTime     meta.ExposureTime
	ExposureProgram  meta.ExposureProgram
	ExposureMode     meta.ExposureMode
	ExposureBias     meta.ExposureBias
	ISOSpeedRatings  uint32
	Flash            Flash
	MeteringMode     meta.MeteringMode
	Aperture         meta.Aperture
	FocalLength      meta.FocalLength
	SubjectDistance  float32
	GPSLatitude      float64
	GPSLongitude     float64
	GPSAltitude      float32
	GPSTimestamp     time.Time
}

func (exif *Exif) parse(p property) (err error) {
	switch p.Name() {
	case xmpns.PixelXDimension:
		exif.PixelXDimension = parseUint32(p.Value())
	case xmpns.PixelYDimension:
		exif.PixelYDimension = parseUint32(p.Value())
	case xmpns.DateTimeOriginal:
		exif.DateTimeOriginal, err = parseDate(p.Value())
	case xmpns.ExposureTime:
		n, d := parseRational(p.Value())
		exif.ExposureTime = meta.ExposureTime(float32(n) / float32(d))
	case xmpns.ExposureProgram:
		exif.ExposureProgram = meta.ExposureProgram(parseUint8(p.Value()))
	case xmpns.ExposureMode:
		exif.ExposureMode = meta.NewExposureMode(parseUint8(p.Value()))
	case xmpns.ExposureBiasValue:
		err = exif.ExposureBias.UnmarshalText(p.Value())
	case xmpns.FocalLength:
		// TODO: needs verification
		n, d := parseRational(p.Value())
		exif.FocalLength = meta.NewFocalLength(n, d)
	case xmpns.SubjectDistance:
		// TODO: needs verification
		n, d := parseRational(p.Value())
		exif.SubjectDistance = float32(float32(n) / float32(d))
	case xmpns.MeteringMode:
		exif.MeteringMode = meta.NewMeteringMode(uint16(parseUint8(p.Value())))
	case xmpns.FNumber:
		// TODO: needs verification
		n, d := parseRational(p.Value())
		exif.Aperture = meta.NewAperture(n, d)
	case xmpns.ISOSpeedRatings:
		exif.ISOSpeedRatings = parseUint32(p.val)
	case xmpns.GPSLatitude:
		exif.GPSLatitude = parseFloat64(p.Value())
	case xmpns.GPSLongitude:
		exif.GPSLongitude = parseFloat64(p.Value())
	case xmpns.GPSAltitude:
		exif.GPSAltitude = float32(parseFloat64(p.Value()))
	//case xmpns.Flash:
	default:
		return ErrPropertyNotSet
	}
	return
}

// Aux attributes of an XMP Packet. These are Adobe-defined auxiliary EXIF tags.
// This implmentation is incomplete and based on https://exiftool.org/TagNames/XMP.html#aux
type Aux struct {
	SerialNumber             string
	LensInfo                 string
	Lens                     string
	LensID                   uint32
	LensSerialNumber         string
	ImageNumber              uint16
	ApproximateFocusDistance string            // rational
	FlashCompensation        meta.ExposureBias // rational
	Firmware                 string
}

func (aux *Aux) parse(p property) (err error) {
	switch p.Name() {
	case xmpns.FlashCompensation:
		err = aux.FlashCompensation.UnmarshalText(p.Value())
	case xmpns.ImageNumber:
		aux.ImageNumber = uint16(parseUint(p.Value()))
	case xmpns.SerialNumber:
		aux.SerialNumber = parseString(p.Value())
	case xmpns.Lens:
		aux.Lens = parseString(p.Value())
	case xmpns.LensInfo:
		aux.LensInfo = parseString(p.Value())
	case xmpns.LensID:
		aux.LensID = parseUint32(p.Value())
	case xmpns.LensSerialNumber:
		aux.LensSerialNumber = parseString(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return
}
