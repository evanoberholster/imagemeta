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
	// Function indicates whether the flash function is present.
	Function bool
	// Mode is the EXIF flash mode code.
	Mode uint8
	// RedEyeMode indicates whether red-eye reduction was enabled.
	RedEyeMode bool
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
	ExifVersion      string // Often denfined as 0221 -> Exif 2.21
	PixelXDimension  uint32
	PixelYDimension  uint32
	DateTime         time.Time
	DateTimeOriginal time.Time
	CreateDate       time.Time // Exif:DateTimeDigitized
	ExposureTime     meta.Rational32
	ExposureProgram  meta.ExposureProgram
	ExposureMode     meta.ExposureMode
	ExposureBias     meta.ExposureBias
	ISOSpeedRatings  uint32
	Flash            Flash
	MeteringMode     meta.MeteringMode
	Aperture         meta.Aperture
	FocalLength      meta.FocalLength
	SubjectDistance  float32
	GPS              meta.GPS
	// ApertureValue is converted from EXIF APEX units to an f-number
	// to match ExifTool output.
	ApertureValue                meta.Aperture
	BrightnessValue              float64
	CameraOwnerName              string
	BodySerialNumber             string
	ColorSpace                   uint16
	ComponentsConfiguration      string
	CompressedBitsPerPixel       float64
	CustomRendered               uint8
	DigitalZoomRatio             float64
	FileSource                   uint8
	FlashpixVersion              string
	FocalLengthIn35mmFilm        uint16
	FocalPlaneResolutionUnit     uint8
	FocalPlaneXResolution        float64
	FocalPlaneYResolution        float64
	ExposureIndex                float64
	StandardOutputSensitivity    uint32
	ISOSpeed                     uint32
	ISOSpeedLatitudeyyy          uint16
	ISOSpeedLatitudezzz          uint16
	FlashEnergy                  float64
	Acceleration                 float64
	AmbientTemperature           float64
	Humidity                     float64
	Pressure                     float64
	WaterDepth                   float64
	Gamma                        float64
	CameraElevationAngle         float64
	GainControl                  uint8
	CompImageImagesPerSequence   uint32
	CompImageMaxExposureAll      float64
	CompImageMaxExposureUsed     float64
	CompImageMinExposureAll      float64
	CompImageMinExposureUsed     float64
	CompImageNumSequences        uint32
	CompImageSumExposureAll      float64
	CompImageSumExposureUsed     float64
	CompImageTotalExposurePeriod float64
	CompImageValues              string
	CompositeImage               uint8
	CompositeImageCount          uint16
	CompositeImageExposureTimes  string
	ImageUniqueID                string
	ImageTitle                   string
	ImageEditor                  string
	ImageEditingSoftware         string
	MetadataEditingSoftware      string
	RAWDevelopingSoftware        string
	Photographer                 string
	OwnerName                    string
	InteroperabilityIndex        string
	LightSource                  uint16
	MakerNote                    string
	// MaxApertureValue is converted from EXIF APEX units to an f-number
	// to match ExifTool output.
	MaxApertureValue                 meta.Aperture
	NativeDigest                     string
	OECF                             string
	OECFColumns                      uint16
	OECFNames                        string
	OECFRows                         uint16
	OECFValues                       string
	PhotometricInterpretation        uint16
	RecommendedExposureIndex         uint32
	RelatedSoundFile                 string
	SamplesPerPixel                  uint16
	SensingMethod                    uint16
	Saturation                       uint8
	Contrast                         uint8
	SceneCaptureType                 uint8
	SceneType                        uint8
	SensitivityType                  uint16
	Sharpness                        uint8
	SpatialFrequencyResponse         string
	SpatialFrequencyResponseColumns  uint16
	SpatialFrequencyResponseNames    string
	SpatialFrequencyResponseRows     uint16
	SpatialFrequencyResponseValues   string
	SpectralSensitivity              string
	SubjectArea                      string
	SubjectDistanceRange             uint16
	SubjectLocation                  string
	ShutterSpeedValue                meta.Rational32
	CFAPattern                       string
	CFAPatternColumns                uint16
	CFAPatternRows                   uint16
	CFAPatternValues                 string
	DeviceSettingDescription         string
	DeviceSettingDescriptionColumns  uint16
	DeviceSettingDescriptionRows     uint16
	DeviceSettingDescriptionSettings string
	WhiteBalance                     uint8
	UserComment                      string
	CameraFirmware                   string
	LensMake                         string
	LensModel                        string
	LensInfo                         string
	LensSerialNumber                 string
	SubsecTime                       string
	SubsecTimeDigitized              string
	SubsecTimeOriginal               string
}

func (exif *Exif) parse(p property) (err error) {
	switch p.Name() {
	case ExifVersion:
		exif.ExifVersion = parseString(p.Value())
	case PixelXDimension:
		exif.PixelXDimension = parseUint32(p.Value())
	case PixelYDimension:
		exif.PixelYDimension = parseUint32(p.Value())
	case DateTime:
		err = parseDateWithSubseconds(&exif.DateTime, p.Value(), exif.SubsecTime)
	case DateTimeDigitized:
		err = parseDateWithSubseconds(&exif.CreateDate, p.Value(), exif.SubsecTimeDigitized)
	case DateTimeOriginal:
		err = parseDateWithSubseconds(&exif.DateTimeOriginal, p.Value(), exif.SubsecTimeOriginal)
	case ApertureValue:
		exif.ApertureValue = meta.Aperture(parseApexAperture(p.Value()))
	case ExposureTime:
		n, d := parseRational32(p.Value())
		exif.ExposureTime = meta.NewRational32(n, d)
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
		exif.FocalLengthIn35mmFilm = parseUint16(p.Value())
	case FocalPlaneResolutionUnit:
		exif.FocalPlaneResolutionUnit = parseUint8(p.Value())
	case FocalPlaneXResolution:
		exif.FocalPlaneXResolution = parseRationalFloat64(p.Value())
	case FocalPlaneYResolution:
		exif.FocalPlaneYResolution = parseRationalFloat64(p.Value())
	case ExposureIndex:
		exif.ExposureIndex = parseRationalFloat64(p.Value())
	case StandardOutputSensitivity:
		exif.StandardOutputSensitivity = parseUint32(p.Value())
	case ISOSpeed:
		exif.ISOSpeed = parseUint32(p.Value())
	case ISOSpeedLatitudeyyy:
		exif.ISOSpeedLatitudeyyy = parseUint16(p.Value())
	case ISOSpeedLatitudezzz:
		exif.ISOSpeedLatitudezzz = parseUint16(p.Value())
	case FlashEnergy:
		exif.FlashEnergy = parseRationalFloat64(p.Value())
	case Acceleration:
		exif.Acceleration = parseRationalFloat64(p.Value())
	case AmbientTemperature:
		exif.AmbientTemperature = parseRationalFloat64(p.Value())
	case Humidity:
		exif.Humidity = parseRationalFloat64(p.Value())
	case Pressure:
		exif.Pressure = parseRationalFloat64(p.Value())
	case WaterDepth:
		exif.WaterDepth = parseRationalFloat64(p.Value())
	case Gamma:
		exif.Gamma = parseRationalFloat64(p.Value())
	case CameraElevationAngle:
		exif.CameraElevationAngle = parseRationalFloat64(p.Value())
	case SubjectDistance:
		n, d := parseRational32(p.Value())
		exif.SubjectDistance = float32(n) / float32(d)
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
		exif.GPS.Latitude = parseGPSCoordinateWithReference(p.Value(), meta.GPSRefNorth, meta.GPSRefSouth)
	case GPSLongitude:
		exif.GPS.Longitude = parseGPSCoordinateWithReference(p.Value(), meta.GPSRefEast, meta.GPSRefWest)
	case GPSAltitude:
		exif.GPS.Altitude = parseGPSAltitude(p.Value(), exif.GPS.Altitude.Ref)
	case GPSTimeStamp:
		exif.GPS.Time, err = parseDate(p.Value())
	case GPSAltitudeRef:
		exif.GPS.Altitude.Ref = parseGPSRef(p.Value(), gpsRefKindAltitude)
	case GPSAreaInformation:
		exif.GPS.AreaInformation = parseString(p.Value())
	case GPSDestBearing:
		exif.GPS.DestinationBearing.Value = parseRationalFloat64(p.Value())
	case GPSDestBearingRef:
		exif.GPS.DestinationBearing.Ref = parseGPSRef(p.Value(), gpsRefKindDirection)
	case GPSDestDistance:
		exif.GPS.DestinationDistance.Value = parseRationalFloat64(p.Value())
	case GPSDestDistanceRef:
		exif.GPS.DestinationDistance.Ref = parseGPSRef(p.Value(), gpsRefKindDistance)
	case GPSDestLatitude:
		exif.GPS.DestinationLatitude = parseGPSCoordinateWithReference(p.Value(), meta.GPSRefNorth, meta.GPSRefSouth)
	case GPSDestLongitude:
		exif.GPS.DestinationLongitude = parseGPSCoordinateWithReference(p.Value(), meta.GPSRefEast, meta.GPSRefWest)
	case GPSDifferential:
		exif.GPS.Differential = parseUint8(p.Value())
	case GPSHPositioningError:
		exif.GPS.HPositioningError = parseRationalFloat64(p.Value())
	case GPSImgDirection:
		exif.GPS.ImageDirection.Value = parseRationalFloat64(p.Value())
	case GPSImgDirectionRef:
		exif.GPS.ImageDirection.Ref = parseGPSRef(p.Value(), gpsRefKindDirection)
	case GPSMapDatum:
		exif.GPS.MapDatum = parseString(p.Value())
	case GPSProcessingMethod:
		exif.GPS.ProcessingMethod = parseString(p.Value())
	case GPSSpeed:
		exif.GPS.Speed.Value = parseRationalFloat64(p.Value())
	case GPSSpeedRef:
		exif.GPS.Speed.Ref = parseGPSRef(p.Value(), gpsRefKindDistance)
	case GPSStatus:
		exif.GPS.Status = parseString(p.Value())
	case GPSDOP:
		exif.GPS.DOP = parseRationalFloat64(p.Value())
	case GPSMeasureMode:
		exif.GPS.MeasureMode = parseString(p.Value())
	case GPSSatellites:
		exif.GPS.Satellites = parseString(p.Value())
	case GPSTrack:
		exif.GPS.Track.Value = parseRationalFloat64(p.Value())
	case GPSTrackRef:
		exif.GPS.Track.Ref = parseGPSRef(p.Value(), gpsRefKindDirection)
	case GPSVersionID:
		exif.GPS.VersionID = parseString(p.Value())
	case ImageUniqueID:
		exif.ImageUniqueID = parseString(p.Value())
	case CameraOwnerName:
		exif.CameraOwnerName = parseString(p.Value())
	case BodySerialNumber:
		exif.BodySerialNumber = parseString(p.Value())
	case ColorSpace:
		exif.ColorSpace = parseUint16(p.Value())
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
	case CompImageImagesPerSequence:
		exif.CompImageImagesPerSequence = parseUint32(p.Value())
	case CompImageMaxExposureAll:
		exif.CompImageMaxExposureAll = parseRationalFloat64(p.Value())
	case CompImageMaxExposureUsed:
		exif.CompImageMaxExposureUsed = parseRationalFloat64(p.Value())
	case CompImageMinExposureAll:
		exif.CompImageMinExposureAll = parseRationalFloat64(p.Value())
	case CompImageMinExposureUsed:
		exif.CompImageMinExposureUsed = parseRationalFloat64(p.Value())
	case CompImageNumSequences:
		exif.CompImageNumSequences = parseUint32(p.Value())
	case CompImageSumExposureAll:
		exif.CompImageSumExposureAll = parseRationalFloat64(p.Value())
	case CompImageSumExposureUsed:
		exif.CompImageSumExposureUsed = parseRationalFloat64(p.Value())
	case CompImageTotalExposurePeriod:
		exif.CompImageTotalExposurePeriod = parseRationalFloat64(p.Value())
	case CompImageValues:
		exif.CompImageValues = parseString(p.Value())
	case CompositeImage:
		exif.CompositeImage = parseUint8(p.Value())
	case CompositeImageCount:
		exif.CompositeImageCount = parseUint16(p.Value())
	case CompositeImageExposureTimes:
		exif.CompositeImageExposureTimes = parseString(p.Value())
	case InteroperabilityIndex:
		exif.InteroperabilityIndex = parseString(p.Value())
	case LightSource:
		exif.LightSource = parseUint16(p.Value())
	case MakerNote:
		exif.MakerNote = parseString(p.Value())
	case MaxApertureValue:
		exif.MaxApertureValue = meta.Aperture(parseApexAperture(p.Value()))
	case NativeDigest:
		exif.NativeDigest = parseString(p.Value())
	case OECF:
		exif.OECF = parseString(p.Value())
	case OECFColumns:
		exif.OECFColumns = parseUint16(p.Value())
	case OECFNames:
		exif.OECFNames = parseString(p.Value())
	case OECFRows:
		exif.OECFRows = parseUint16(p.Value())
	case OECFValues:
		exif.OECFValues = parseString(p.Value())
	case PhotometricInterpretation:
		exif.PhotometricInterpretation = parseUint16(p.Value())
	case RecommendedExposureIndex:
		exif.RecommendedExposureIndex = parseUint32(p.Value())
	case RelatedSoundFile:
		exif.RelatedSoundFile = parseString(p.Value())
	case SamplesPerPixel:
		exif.SamplesPerPixel = parseUint16(p.Value())
	case SensingMethod:
		exif.SensingMethod = parseUint16(p.Value())
	case Saturation:
		exif.Saturation = parseUint8(p.Value())
	case Contrast:
		exif.Contrast = parseUint8(p.Value())
	case SceneCaptureType:
		exif.SceneCaptureType = parseUint8(p.Value())
	case SceneType:
		exif.SceneType = parseUint8(p.Value())
	case SensitivityType:
		exif.SensitivityType = parseUint16(p.Value())
	case Sharpness:
		exif.Sharpness = parseUint8(p.Value())
	case SpatialFrequencyResponse:
		exif.SpatialFrequencyResponse = parseString(p.Value())
	case SpatialFrequencyResponseColumns:
		exif.SpatialFrequencyResponseColumns = parseUint16(p.Value())
	case SpatialFrequencyResponseNames:
		exif.SpatialFrequencyResponseNames = parseString(p.Value())
	case SpatialFrequencyResponseRows:
		exif.SpatialFrequencyResponseRows = parseUint16(p.Value())
	case SpatialFrequencyResponseValues:
		exif.SpatialFrequencyResponseValues = parseString(p.Value())
	case SpectralSensitivity:
		exif.SpectralSensitivity = parseString(p.Value())
	case SubjectArea:
		exif.SubjectArea = parseString(p.Value())
	case SubjectDistanceRange:
		exif.SubjectDistanceRange = parseUint16(p.Value())
	case SubjectLocation:
		exif.SubjectLocation = parseString(p.Value())
	case ShutterSpeedValue:
		n, d := parseRational32(p.Value())
		exif.ShutterSpeedValue = meta.NewRational32(n, d)
	case CFAPattern:
		exif.CFAPattern = parseString(p.Value())
	case CFAPatternColumns:
		exif.CFAPatternColumns = parseUint16(p.Value())
	case CFAPatternRows:
		exif.CFAPatternRows = parseUint16(p.Value())
	case CFAPatternValues:
		exif.CFAPatternValues = parseString(p.Value())
	case DeviceSettingDescription:
		exif.DeviceSettingDescription = parseString(p.Value())
	case DeviceSettingDescriptionColumns:
		exif.DeviceSettingDescriptionColumns = parseUint16(p.Value())
	case DeviceSettingDescriptionRows:
		exif.DeviceSettingDescriptionRows = parseUint16(p.Value())
	case DeviceSettingDescriptionSettings:
		exif.DeviceSettingDescriptionSettings = parseString(p.Value())
	case WhiteBalance:
		exif.WhiteBalance = parseUint8(p.Value())
	case UserComment:
		exif.UserComment = parseString(p.Value())
	case CameraFirmware:
		exif.CameraFirmware = parseString(p.Value())
	case ImageTitle:
		exif.ImageTitle = parseString(p.Value())
	case ImageEditor:
		exif.ImageEditor = parseString(p.Value())
	case ImageEditingSoftware:
		exif.ImageEditingSoftware = parseString(p.Value())
	case MetadataEditingSoftware:
		exif.MetadataEditingSoftware = parseString(p.Value())
	case RAWDevelopingSoftware:
		exif.RAWDevelopingSoftware = parseString(p.Value())
	case Photographer:
		exif.Photographer = parseString(p.Value())
	case OwnerName:
		exif.OwnerName = parseString(p.Value())
	case LensMake:
		exif.LensMake = parseString(p.Value())
	case LensModel:
		exif.LensModel = parseString(p.Value())
	case LensInfo:
		exif.LensInfo = parseString(p.Value())
	case LensSerialNumber:
		exif.LensSerialNumber = parseString(p.Value())
	case SerialNumber:
		exif.BodySerialNumber = parseString(p.Value())
	case SubsecTime:
		parseSubsecondsField(&exif.SubsecTime, p.Value(), &exif.DateTime)
	case SubsecTimeDigitized:
		parseSubsecondsField(&exif.SubsecTimeDigitized, p.Value(), &exif.CreateDate)
	case SubsecTimeOriginal:
		parseSubsecondsField(&exif.SubsecTimeOriginal, p.Value(), &exif.DateTimeOriginal)
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
		aux.ImageNumber = parseUint16(p.Value())
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

func parseGPSCoordinateWithReference(buf []byte, positiveRef, negativeRef meta.GPSRef) meta.GPSCoordinate {
	ref := meta.GPSRefUnknown
	if len(buf) > 0 {
		switch buf[len(buf)-1] {
		case 'N', 'n', 'E', 'e':
			ref = positiveRef
		case 'S', 's', 'W', 'w':
			ref = negativeRef
		}
	}

	v := parseGPSCoordinate(buf)
	if v < 0 {
		v = -v
		if ref == meta.GPSRefUnknown {
			ref = negativeRef
		}
	}

	return meta.GPSCoordinate{Value: v, Ref: ref}
}

func parseGPSAltitude(buf []byte, currentRef meta.GPSRef) meta.GPSAltitude {
	v := float32(parseRationalFloat64(buf))
	ref := currentRef
	if v < 0 {
		v = -v
		if ref == meta.GPSRefUnknown {
			ref = meta.GPSRefBelowSeaLevel
		}
	}
	return meta.GPSAltitude{Value: v, Ref: ref}
}

type gpsRefKind uint8

const (
	gpsRefKindAltitude gpsRefKind = iota
	gpsRefKindDirection
	gpsRefKindDistance
)

func parseGPSRef(buf []byte, kind gpsRefKind) meta.GPSRef {
	if len(buf) == 0 {
		return meta.GPSRefUnknown
	}
	switch kind {
	case gpsRefKindAltitude:
		switch buf[0] {
		case '0':
			return meta.GPSRefAboveSeaLevel
		case '1':
			return meta.GPSRefBelowSeaLevel
		}
	case gpsRefKindDirection:
		switch buf[0] {
		case 'T', 't':
			return meta.GPSRefTrue
		case 'M', 'm':
			return meta.GPSRefMagnetic
		}
	case gpsRefKindDistance:
		switch buf[0] {
		case 'K', 'k':
			return meta.GPSRefKilometers
		case 'M', 'm':
			return meta.GPSRefMiles
		case 'N', 'n':
			return meta.GPSRefKnots
		}
	}
	return meta.GPSRefUnknown
}

func parseDateWithSubseconds(target *time.Time, buf []byte, subsec string) error {
	t, err := parseDate(buf)
	if err != nil {
		return err
	}
	*target = addSubseconds(t, subsec)
	return nil
}

func parseSubsecondsField(subsec *string, buf []byte, target *time.Time) {
	*subsec = parseString(buf)
	*target = addSubseconds(*target, *subsec)
}

func addSubseconds(base time.Time, subsec string) time.Time {
	if base.IsZero() || len(subsec) == 0 {
		return base
	}

	ns, ok := parseSubsecondNanoseconds(subsec)
	if !ok || ns == 0 {
		return base
	}

	return base.Add(time.Duration(ns) * time.Nanosecond)
}

func parseSubsecondNanoseconds(subsec string) (int64, bool) {
	var value int64
	digits := 0

	for i := 0; i < len(subsec); i++ {
		c := subsec[i]
		if c < '0' || c > '9' {
			if digits > 0 {
				break
			}
			continue
		}

		if digits < 9 {
			value = (value * 10) + int64(c-'0')
		}
		digits++
	}

	if digits == 0 {
		return 0, false
	}
	if digits > 9 {
		digits = 9
	}
	for digits < 9 {
		value *= 10
		digits++
	}

	return value, true
}
