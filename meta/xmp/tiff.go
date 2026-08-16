package xmp

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

// Tiff attributes of an XMP Packet.
//
//	xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
//
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#tiff
type Tiff struct {
	Artist                    string
	BitsPerSample             string
	DateTime                  time.Time
	Make                      string // Camera Make
	Model                     string // Camera Model
	Software                  string
	Copyright                 []string
	ImageDescription          []string
	ImageWidth                uint16
	ImageLength               uint16
	NativeDigest              string
	Orientation               meta.Orientation
	Compression               uint16
	PhotometricInterpretation uint16
	PlanarConfiguration       uint16
	PrimaryChromaticities     string
	ReferenceBlackWhite       string
	ResolutionUnit            meta.ResolutionUnit
	SamplesPerPixel           uint16
	TransferFunction          string
	WhitePoint                string
	XResolution               float64
	YCbCrCoefficients         string
	YCbCrPositioning          uint16
	YCbCrSubSampling          string
	YResolution               float64
}

func (t *Tiff) parse(p property) error {
	switch p.Name() {
	case Artist:
		t.Artist = parseString(p.Value())
	case BitsPerSample:
		t.BitsPerSample = parseString(p.Value())
	case DateTime:
		var err error
		t.DateTime, err = parseDate(p.Value())
		return err
	case Make:
		t.Make = parseString(p.Value())
	case Model:
		t.Model = parseString(p.Value())
	case Software:
		t.Software = parseString(p.Value())
	case Copyright:
		t.Copyright = append(t.Copyright, parseString(p.Value()))
	case ImageDescription:
		t.ImageDescription = append(t.ImageDescription, parseString(p.Value()))
	case ImageWidth:
		t.ImageWidth = parseUint16(p.Value())
	case ImageLength:
		t.ImageLength = parseUint16(p.Value())
	case NativeDigest:
		t.NativeDigest = parseString(p.Value())
	case Orientation:
		t.Orientation = meta.Orientation(parseUint16(p.Value()))
	case Compression:
		t.Compression = parseUint16(p.Value())
	case PhotometricInterpretation:
		t.PhotometricInterpretation = parseUint16(p.Value())
	case PlanarConfiguration:
		t.PlanarConfiguration = parseUint16(p.Value())
	case PrimaryChromaticities:
		t.PrimaryChromaticities = parseString(p.Value())
	case ReferenceBlackWhite:
		t.ReferenceBlackWhite = parseString(p.Value())
	case ResolutionUnit:
		t.ResolutionUnit = meta.ResolutionUnit(parseUint16(p.Value()))
	case SamplesPerPixel:
		t.SamplesPerPixel = parseUint16(p.Value())
	case TransferFunction:
		t.TransferFunction = parseString(p.Value())
	case WhitePoint:
		t.WhitePoint = parseString(p.Value())
	case XResolution:
		t.XResolution = parseRationalFloat64(p.Value())
	case YCbCrCoefficients:
		t.YCbCrCoefficients = parseString(p.Value())
	case YCbCrPositioning:
		t.YCbCrPositioning = parseUint16(p.Value())
	case YCbCrSubSampling:
		t.YCbCrSubSampling = parseString(p.Value())
	case YResolution:
		t.YResolution = parseRationalFloat64(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return nil
}
