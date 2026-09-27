package exif

import (
	"math"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// parseExifTag parses ExifIFD/SubIFD tags into typed model fields.
//
// Non-parsed ExifIFD/SubIFD tags are currently handled by falling through to
// the default case (`return false`) when there is no modeled parser mapping.
//
// TODO(exiftool-coverage): Expand ExifIFD parity against ExifTool
// lib/Image/ExifTool/Exif.pm (%Image::ExifTool::Exif::Main) and
// https://exiftool.org/TagNames/EXIF.html (ExifIFD rows).
// Keep large/opaque payloads opt-in and avoid allocations in hot paths.
func (r *Reader) parseExifTag(t tag.Entry) bool {
	exifIfd := &r.Exif.ExifIFD
	ifd0 := &r.Exif.IFD0
	switch t.ID {
	// Time/date.
	case tag.TagDateTimeOriginal:
		exifIfd.setDate(t.ID, r.parseDate(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagDateTimeDigitized:
		exifIfd.setDate(t.ID, r.parseDate(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagSubSecTime:
		value, raw := r.parseSubSecTimeParts(t)
		ifd0.setSubSec(raw, value)
		ifd0.markTagParsed(t.ID)
	case tag.TagSubSecTimeOriginal:
		value, raw := r.parseSubSecTimeParts(t)
		exifIfd.setSubSec(t.ID, raw, value)
		exifIfd.markTagParsed(t.ID)
	case tag.TagSubSecTimeDigitized:
		value, raw := r.parseSubSecTimeParts(t)
		exifIfd.setSubSec(t.ID, raw, value)
		exifIfd.markTagParsed(t.ID)
	case tag.TagOffsetTime:
		ifd0.setOffset(r.parseOffsetTime(t))
		ifd0.markTagParsed(t.ID)
	case tag.TagOffsetTimeOriginal:
		exifIfd.setOffset(t.ID, r.parseOffsetTime(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagOffsetTimeDigitized:
		exifIfd.setOffset(t.ID, r.parseOffsetTime(t))
		exifIfd.markTagParsed(t.ID)

	// Text/meta.
	// TODO(exiftool-coverage): Add selected ExifIFD text fields with bounded
	// parsing and explicit truncation limits where needed:
	// ImageUniqueID(0xa420), CompositeImageCount(0xa461), and
	// CompositeImageExposureTimes(0xa462) display summaries.
	case tag.TagExifVersion:
		exifIfd.ExifVersion = r.parseExifVersion(t)
	case tag.TagLensMake:
		exifIfd.LensMake = r.parseString(t)
	case tag.TagLensModel:
		exifIfd.LensModel = r.parseString(t)
	case tag.TagLensSerialNumber:
		exifIfd.LensSerial = r.parseString(t)
	case tag.TagCameraOwnerName:
		exifIfd.CameraOwnerName = r.parseString(t)
	case tag.TagBodySerialNumber:
		exifIfd.BodySerialNumber = r.parseString(t)
	case tag.TagUserComment:
		exifIfd.UserComment = r.parseExifUserComment(t)
	case tag.TagFlashpixVersion:
		exifIfd.FlashpixVersion = r.parseStringAllowUndefined(t)
	case tag.TagDeviceSettingDescription:
		// Not parsed.
		// The payload is often large and not needed in the hot parse path.

	// Image characterization.
	// TODO(exiftool-coverage): Consider lightweight modeling for ExifIFD image
	// characterization tags used by ExifTool:
	// OECF(0x8828), SpatialFrequencyResponse(0x920c), SubjectLocation(0xa214),
	// and CFAPattern(0xa302). Prefer compact fixed-size summaries over raw blobs.
	case tag.TagPixelXDimension:
		exifIfd.PixelXDimension = r.parseUint32(t)
	case tag.TagPixelYDimension:
		exifIfd.PixelYDimension = r.parseUint32(t)
	case tag.TagInteropIFDPointer:
		exifIfd.interopIFDPointer = r.parseUint32(t)
	case tag.TagInteropIndex:
		exifIfd.InteropIndex = r.parseStringAllowUndefined(t)
	case tag.TagInteropVersion:
		exifIfd.InteropVersion = r.parseStringAllowUndefined(t)
	case tag.TagRelatedImageWidth:
		exifIfd.RelatedImageWidth = r.parseUint32(t)
	case tag.TagRelatedImageHeight:
		exifIfd.RelatedImageHeight = r.parseUint32(t)
	case tag.TagColorSpace:
		exifIfd.ColorSpace = r.parseUint16(t)
	case tag.TagLensSpecification:
		exifIfd.LensInfo = r.parseLensInfo(t)
	case tag.TagComponentsConfiguration:
		r.parseByteList(t, exifIfd.ComponentsConfiguration[:])
	case tag.TagCompressedBitsPerPixel:
		exifIfd.CompressedBitsPerPixel = r.parseRationalValue(t).Float64()
	case tag.TagFocalPlaneXResolution:
		exifIfd.FocalPlaneXResolution, exifIfd.focalPlaneXResolutionState = r.parseUnsignedRationalFloat64(t)
	case tag.TagFocalPlaneYResolution:
		exifIfd.FocalPlaneYResolution, exifIfd.focalPlaneYResolutionState = r.parseUnsignedRationalFloat64(t)
	case tag.TagFocalPlaneResolutionUnit:
		exifIfd.FocalPlaneResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagSubjectArea:
		r.parseUint16List(t, exifIfd.SubjectArea[:])
	case tag.TagGamma:
		exifIfd.Gamma = r.parseRationalValue(t).Float64()

	// Exposure/optics.
	// TODO(exiftool-coverage): Add ISO-family ExifIFD fields documented by
	// ExifTool for EXIF 2.31+ parity:
	// StandardOutputSensitivity(0x8831), ISOSpeed(0x8833),
	// ISOSpeedLatitudeyyy(0x8834), ISOSpeedLatitudezzz(0x8835).
	// Keep parse paths branch-light and avoid allocations.
	case tag.TagExposureTime:
		exifIfd.ExposureTime = r.parseExposureTime(t)
	case tag.TagShutterSpeedValue:
		exifIfd.ShutterSpeedValue = r.parseShutterSpeed(t)
	case tag.TagFNumber:
		exifIfd.FNumber = r.parseAperture(t)
	case tag.TagApertureValue:
		exifIfd.ApertureValue = r.parseApexAperture(t)
		if exifIfd.FNumber == 0 && apertureIsFinite(exifIfd.ApertureValue) {
			exifIfd.FNumber = exifIfd.ApertureValue
		}
	case tag.TagMaxApertureValue:
		exifIfd.MaxApertureValue = r.parseApexAperture(t)
	case tag.TagSubjectDistance:
		exifIfd.SubjectDistance, exifIfd.subjectDistanceState = r.parseUnsignedRationalFloat64(t)
	case tag.TagBrightnessValue:
		exifIfd.BrightnessValue = r.parseSignedRationalFloat32(t)
	case tag.TagExposureProgram:
		exifIfd.ExposureProgram = meta.ExposureProgram(r.parseUint16(t))
	case tag.TagSensitivityType:
		exifIfd.SensitivityType = r.parseUint16(t)
	case tag.TagRecommendedExposureIndex:
		exifIfd.RecommendedExposureIndex = r.parseUint32(t)
	case tag.TagExposureBiasValue:
		exifIfd.ExposureBias = r.parseExposureBias(t)
	case tag.TagExposureMode:
		exifIfd.ExposureMode = meta.ExposureMode(r.parseUint16(t))
	case tag.TagMeteringMode:
		exifIfd.MeteringMode = meta.MeteringMode(r.parseUint16(t))
	case tag.TagLightSource:
		exifIfd.LightSource = r.parseUint16(t)
	case tag.TagISOSpeedRatings:
		exifIfd.ISOSpeedRatings = r.parseUint32(t)
	case tag.TagFlash:
		exifIfd.Flash = meta.Flash(r.parseUint16(t))
	case tag.TagFocalLength:
		exifIfd.FocalLength = r.parseFocalLength(t)
	case tag.TagFocalLengthIn35mmFilm:
		exifIfd.FocalLengthIn35mmFormat = r.parseFocalLength(t)
	case tag.TagExposureIndex:
		exifIfd.ExposureIndex, exifIfd.exposureIndexState = r.parseUnsignedRationalFloat64(t)

	// Capture/mode.
	// TODO(exiftool-coverage): Evaluate adding environmental capture tags
	// surfaced by ExifTool:
	// AmbientTemperature(0x9400), Humidity(0x9401), Pressure(0x9402),
	// WaterDepth(0x9403), Acceleration(0x9404), CameraElevationAngle(0x9405).
	// These are niche but useful for action cameras and drone media.
	case tag.TagSensingMethod:
		exifIfd.SensingMethod = r.parseUint16(t)
	case tag.TagFileSource:
		exifIfd.FileSource = r.parseSceneType(t)
	case tag.TagSceneType:
		exifIfd.SceneType = r.parseSceneType(t)
	case tag.TagCustomRendered:
		exifIfd.CustomRendered = r.parseUint16(t)
	case tag.TagWhiteBalance:
		exifIfd.WhiteBalance = r.parseUint16(t)
	case tag.TagDigitalZoomRatio:
		exifIfd.DigitalZoomRatio = float32(r.parseRationalValue(t).Float64())
	case tag.TagSceneCaptureType:
		exifIfd.SceneCaptureType = r.parseUint16(t)
	case tag.TagGainControl:
		exifIfd.GainControl = r.parseUint16(t)
	case tag.TagContrast:
		exifIfd.Contrast = r.parseUint16(t)
	case tag.TagSaturation:
		exifIfd.Saturation = r.parseUint16(t)
	case tag.TagSharpness:
		exifIfd.Sharpness = r.parseUint16(t)
	case tag.TagSubjectDistanceRange:
		exifIfd.SubjectDistanceRange = r.parseUint16(t)
	case tag.TagCompositeImage:
		exifIfd.CompositeImage = r.parseUint16(t)
	default:
		return false
	}
	return true
}

func apertureIsFinite(v meta.Aperture) bool {
	f := float64(v)
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

func apexApertureToFNumber(v float64) meta.Aperture {
	if v == 0 {
		return 0
	}
	// Some cameras write a large sentinel for "infinite" aperture. Mirror
	// ExifTool behavior and preserve that as +Inf instead of a huge float.
	if math.Abs(v) > 1024 {
		return meta.Aperture(math.Inf(1))
	}
	return meta.Aperture(math.Exp2(v * 0.5))
}

func apexShutterSpeedToSeconds(v float64) meta.ShutterSpeed {
	if math.IsNaN(v) || math.Abs(v) >= 100 {
		return 0
	}
	return meta.ShutterSpeed(math.Exp2(-v))
}
