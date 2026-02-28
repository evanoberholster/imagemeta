package xmp

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

// Flash represents the flattened subfields of the XMP exif:Flash structure.
// Based on https://exiftool.org/TagNames/XMP.html.
type Flash struct {
	// Fired indicates whether the flash fired.
	Fired bool
	// Mode is the EXIF flash mode code.
	Mode uint8
	// RedEyeMode indicates whether red-eye reduction was enabled.
	RedEyeMode bool
	// Function indicates whether the flash function is present.
	Function bool
	// Return is the EXIF flash return-detection code.
	Return uint8
}

// Exif stores decoded EXIF tags from an XMP packet.
//
//	Exif 2.21 or later: xmlns:exifEX="http://cipa.jp/exif/1.0/"
//	Exif 2.2 or earlier: xmlns:exif="http://ns.adobe.com/exif/1.0/"
//
// This implementation is based on https://exiftool.org/TagNames/XMP.html#exif.
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
	// ApertureValue is converted from EXIF APEX units to an f-number
	// to match ExifTool output.
	ApertureValue            meta.Aperture
	BrightnessValue          float64
	CameraOwnerName          string
	BodySerialNumber         string
	ColorSpace               uint16
	ComponentsConfiguration  string
	CompressedBitsPerPixel   float64
	CustomRendered           uint8
	DigitalZoomRatio         float64
	FileSource               uint8
	FlashpixVersion          string
	FocalLengthIn35mmFilm    uint16
	FocalPlaneResolutionUnit uint16
	FocalPlaneXResolution    float64
	FocalPlaneYResolution    float64
	GainControl              uint8
	GPSAltitudeRef           uint8
	GPSDifferential          uint8
	GPSMapDatum              string
	GPSStatus                string
	GPSDOP                   float64
	GPSMeasureMode           string
	GPSSatellites            string
	GPSVersionID             string
	InteroperabilityIndex    string
	LightSource              uint16
	// MaxApertureValue is converted from EXIF APEX units to an f-number
	// to match ExifTool output.
	MaxApertureValue          meta.Aperture
	PhotometricInterpretation uint16
	RecommendedExposureIndex  uint32
	SamplesPerPixel           uint16
	Saturation                uint8
	SceneCaptureType          uint8
	SceneType                 uint8
	SensitivityType           uint16
	Sharpness                 uint8
	// ShutterSpeedValue is converted from EXIF APEX units to seconds
	// to match ExifTool output.
	ShutterSpeedValue   float64
	WhiteBalance        uint8
	UserComment         string
	LensModel           string
	LensInfo            string
	LensSerialNumber    string
	SubsecTime          string
	SubsecTimeDigitized string
	SubsecTimeOriginal  string
}

func (exif *Exif) parse(p property) (err error) {
	switch p.Name() {
	case ExifVersion:
		exif.ExifVersion = parseString(p.Value())
	case PixelXDimension:
		exif.PixelXDimension = parseUint32(p.Value())
	case PixelYDimension:
		exif.PixelYDimension = parseUint32(p.Value())
	case DateTimeDigitized:
		exif.CreateDate, err = parseDate(p.Value())
	case DateTimeOriginal:
		exif.DateTimeOriginal, err = parseDate(p.Value())
	case ApertureValue:
		exif.ApertureValue = meta.Aperture(parseApexAperture(p.Value()))
	case ExposureTime:
		n, d := parseRational(p.Value())
		exif.ExposureTime = meta.ExposureTime(float32(n) / float32(d))
	case ExposureProgram:
		exif.ExposureProgram = meta.ExposureProgram(parseUint8(p.Value()))
	case ExposureMode:
		exif.ExposureMode = meta.NewExposureMode(parseUint8(p.Value()))
	case ExposureBiasValue:
		err = exif.ExposureBias.UnmarshalText(p.Value())
	case BrightnessValue:
		exif.BrightnessValue = parseRationalFloat64(p.Value())
	case FocalLength:
		n, d := parseRational(p.Value())
		exif.FocalLength = meta.NewFocalLength(n, d)
	case FocalLengthIn35mmFilm:
		exif.FocalLengthIn35mmFilm = uint16(parseUint(p.Value()))
	case FocalPlaneResolutionUnit:
		exif.FocalPlaneResolutionUnit = uint16(parseUint(p.Value()))
	case FocalPlaneXResolution:
		exif.FocalPlaneXResolution = parseRationalFloat64(p.Value())
	case FocalPlaneYResolution:
		exif.FocalPlaneYResolution = parseRationalFloat64(p.Value())
	case SubjectDistance:
		n, d := parseRational(p.Value())
		exif.SubjectDistance = float32(float32(n) / float32(d))
	case MeteringMode:
		exif.MeteringMode = meta.NewMeteringMode(uint16(parseUint8(p.Value())))
	case FNumber:
		n, d := parseRational(p.Value())
		exif.Aperture = meta.NewAperture(n, d)
	case ISOSpeedRatings:
		exif.ISOSpeedRatings = parseUint32(p.Value())
	case PhotographicSensitivity:
		exif.ISOSpeedRatings = parseUint32(p.Value())
	case GPSLatitude:
		exif.GPSLatitude = parseGPSCoordinate(p.Value())
	case GPSLongitude:
		exif.GPSLongitude = parseGPSCoordinate(p.Value())
	case GPSAltitude:
		exif.GPSAltitude = float32(parseRationalFloat64(p.Value()))
	case GPSTimeStamp:
		exif.GPSTimestamp, err = parseDate(p.Value())
	case GPSAltitudeRef:
		exif.GPSAltitudeRef = parseUint8(p.Value())
	case GPSDifferential:
		exif.GPSDifferential = parseUint8(p.Value())
	case GPSMapDatum:
		exif.GPSMapDatum = parseString(p.Value())
	case GPSStatus:
		exif.GPSStatus = parseString(p.Value())
	case GPSDOP:
		exif.GPSDOP = parseRationalFloat64(p.Value())
	case GPSMeasureMode:
		exif.GPSMeasureMode = parseString(p.Value())
	case GPSSatellites:
		exif.GPSSatellites = parseString(p.Value())
	case GPSVersionID:
		exif.GPSVersionID = parseString(p.Value())
	case CameraOwnerName:
		exif.CameraOwnerName = parseString(p.Value())
	case BodySerialNumber:
		exif.BodySerialNumber = parseString(p.Value())
	case ColorSpace:
		exif.ColorSpace = uint16(parseUint(p.Value()))
	case ComponentsConfiguration:
		exif.ComponentsConfiguration = parseString(p.Value())
	case CompressedBitsPerPixel:
		exif.CompressedBitsPerPixel = parseRationalFloat64(p.Value())
	case CustomRendered:
		exif.CustomRendered = parseUint8(p.Value())
	case DigitalZoomRatio:
		exif.DigitalZoomRatio = parseRationalFloat64(p.Value())
	case FileSource:
		exif.FileSource = parseUint8(p.Value())
	case FlashpixVersion:
		exif.FlashpixVersion = parseString(p.Value())
	case GainControl:
		exif.GainControl = parseUint8(p.Value())
	case InteroperabilityIndex:
		exif.InteroperabilityIndex = parseString(p.Value())
	case LightSource:
		exif.LightSource = uint16(parseUint(p.Value()))
	case MaxApertureValue:
		exif.MaxApertureValue = meta.Aperture(parseApexAperture(p.Value()))
	case PhotometricInterpretation:
		exif.PhotometricInterpretation = uint16(parseUint(p.Value()))
	case RecommendedExposureIndex:
		exif.RecommendedExposureIndex = parseUint32(p.Value())
	case SamplesPerPixel:
		exif.SamplesPerPixel = uint16(parseUint(p.Value()))
	case Saturation:
		exif.Saturation = parseUint8(p.Value())
	case SceneCaptureType:
		exif.SceneCaptureType = parseUint8(p.Value())
	case SceneType:
		exif.SceneType = parseUint8(p.Value())
	case SensitivityType:
		exif.SensitivityType = uint16(parseUint(p.Value()))
	case Sharpness:
		exif.Sharpness = parseUint8(p.Value())
	case ShutterSpeedValue:
		exif.ShutterSpeedValue = parseApexShutterSpeed(p.Value())
	case WhiteBalance:
		exif.WhiteBalance = parseUint8(p.Value())
	case UserComment:
		exif.UserComment = parseString(p.Value())
	case LensModel:
		exif.LensModel = parseString(p.Value())
	case LensInfo:
		exif.LensInfo = parseString(p.Value())
	case LensSerialNumber:
		exif.LensSerialNumber = parseString(p.Value())
	case SerialNumber:
		exif.BodySerialNumber = parseString(p.Value())
	case SubsecTime:
		exif.SubsecTime = parseString(p.Value())
	case SubsecTimeDigitized:
		exif.SubsecTimeDigitized = parseString(p.Value())
	case SubsecTimeOriginal:
		exif.SubsecTimeOriginal = parseString(p.Value())
	case Fired:
		exif.Flash.Fired = parseBool(p.Value())
	case Return:
		exif.Flash.Return = parseUint8(p.Value())
	case Mode:
		exif.Flash.Mode = parseUint8(p.Value())
	case Function:
		exif.Flash.Function = parseBool(p.Value())
	case RedEyeMode:
		exif.Flash.RedEyeMode = parseBool(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return
}

// Aux attributes of an XMP Packet. These are Adobe-defined auxiliary EXIF tags.
// This implementation is based on https://exiftool.org/TagNames/XMP.html#aux.
type Aux struct {
	// SerialNumber is camera serial number.
	SerialNumber string
	// LensInfo stores the lens range descriptor.
	LensInfo string
	// Lens stores the lens model string.
	Lens string
	// LensID is vendor-specific numeric lens identifier.
	LensID uint32
	// LensSerialNumber is the lens serial number.
	LensSerialNumber string
	// ImageNumber is the camera image counter.
	ImageNumber uint16
	// ApproximateFocusDistance stores the raw rational distance token.
	ApproximateFocusDistance string
	// FlashCompensation is stored as EXIF rational exposure bias.
	FlashCompensation meta.ExposureBias
	// Firmware stores camera firmware version.
	Firmware string
	// DistortionCorrectionAlreadyApplied indicates geometric correction state.
	DistortionCorrectionAlreadyApplied bool
	// LateralChromaticAberrationCorrectionAlreadyApplied indicates CA correction state.
	LateralChromaticAberrationCorrectionAlreadyApplied bool
	// VignetteCorrectionAlreadyApplied indicates vignette correction state.
	VignetteCorrectionAlreadyApplied bool
}

func (aux *Aux) parse(p property) (err error) {
	switch p.Name() {
	case ApproximateFocusDistance:
		aux.ApproximateFocusDistance = parseString(p.Value())
	case FlashCompensation:
		err = aux.FlashCompensation.UnmarshalText(p.Value())
	case ImageNumber:
		aux.ImageNumber = uint16(parseUint(p.Value()))
	case SerialNumber:
		aux.SerialNumber = parseString(p.Value())
	case Lens:
		aux.Lens = parseString(p.Value())
	case LensInfo:
		aux.LensInfo = parseString(p.Value())
	case LensID:
		aux.LensID = parseUint32(p.Value())
	case LensSerialNumber:
		aux.LensSerialNumber = parseString(p.Value())
	case Firmware:
		aux.Firmware = parseString(p.Value())
	case DistortionCorrectionAlreadyApplied:
		aux.DistortionCorrectionAlreadyApplied = parseBool(p.Value())
	case LateralChromaticAberrationCorrectionAlreadyApplied:
		aux.LateralChromaticAberrationCorrectionAlreadyApplied = parseBool(p.Value())
	case VignetteCorrectionAlreadyApplied:
		aux.VignetteCorrectionAlreadyApplied = parseBool(p.Value())
	default:
		return ErrPropertyNotSet
	}
	return
}
