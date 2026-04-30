package exif

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// Exif is the parsed EXIF result for the new meta/exif parser.
type Exif struct {
	GPS          GPSInfo
	IFD0         IFD0Tag
	ExifIFD      ExifIFDTags
	IFD1         *ImageIFD
	IFD2         *ImageIFD
	DNG          *makernote.DNG
	MakerNote    makernote.Info
	CameraSerial string
	CameraMakeID makernote.CameraMake
	ImageType    imagetype.ImageType
}

// IFD0Tag groups tags from the primary image IFD (IFD0).
type IFD0Tag struct {
	IFD0TimeTags
	XResolution      float64
	YResolution      float64
	Make             string
	Model            string
	Artist           string
	Copyright        string
	Software         string
	ImageDescription string
	// TODO: ApplicationNotes (UNDEFINED) may contain large payloads (for example XMP).
	// Parsing is currently disabled to avoid unnecessary allocations.

	ImageWidth      uint32
	ImageHeight     uint32
	ImageOffset     uint32
	ImageLength     uint32
	ThumbnailOffset uint32
	ThumbnailLength uint32
	SubfileType     meta.SubfileType

	Compression               meta.Compression
	Orientation               meta.Orientation
	ResolutionUnit            meta.ResolutionUnit
	PhotometricInterpretation uint16
	SamplesPerPixel           uint16
	PlanarConfiguration       uint16
	BitsPerSample             [4]uint16
	WhitePoint                [2]float64
	PrimaryChromaticities     [6]float64
	YCbCrCoefficients         [3]float64
	YCbCrPositioning          uint16
	Rating                    uint16

	exifIfdPointer    uint32
	gpsIfdPointer     uint32
	subIFDOffsets     [8]uint32
	subIFDOffsetCount uint8
}

// ExifIFDTags groups tags from the ExifIFD directory.
// Field names and tag IDs follow ExifTool EXIF tag naming where applicable:
// https://exiftool.org/TagNames/EXIF.html
type ExifIFDTags struct {
	ExifIFDTimeTags
	LensInfo LensInfo // 0xa432 LensInfo (LensSpecification)
	// TODO: ExifVersion and FlashpixVersion are UNDEFINED in spec.
	// Keep canonical text form for parity with exiftool output.
	ExifVersion     ExifVersion // 0x9000 ExifVersion
	FlashpixVersion string      // 0xa000 FlashpixVersion

	FocalPlaneXResolution  float64 // 0xa20e FocalPlaneXResolution
	FocalPlaneYResolution  float64 // 0xa20f FocalPlaneYResolution
	CompressedBitsPerPixel float64 // 0x9102 CompressedBitsPerPixel
	SubjectDistance        float64 // 0x9206 SubjectDistance
	ExposureIndex          float64 // 0xa215 ExposureIndex

	LensMake         string // 0xa433 LensMake
	LensModel        string // 0xa434 LensModel
	LensSerial       string // 0xa435 LensSerialNumber
	CameraOwnerName  string // 0xa430 OwnerName (CameraOwnerName)
	BodySerialNumber string // 0xa431 SerialNumber (BodySerialNumber)
	UserComment      string // 0x9286 UserComment

	ExposureTime             meta.ExposureTime    // 0x829a ExposureTime
	FocalLength              meta.FocalLength     // 0x920a FocalLength
	FocalLengthIn35mmFormat  meta.FocalLength     // 0xa405 FocalLengthIn35mmFormat
	FNumber                  meta.Aperture        // 0x829d FNumber
	ApertureValue            meta.Aperture        // 0x9202 ApertureValue converted from APEX to F-number
	MaxApertureValue         meta.Aperture        // 0x9205 MaxApertureValue converted from APEX to F-number
	ShutterSpeedValue        meta.ShutterSpeed    // 0x9201 ShutterSpeedValue converted from APEX to seconds
	BrightnessValue          float32              // 0x9203 BrightnessValue
	ExposureProgram          meta.ExposureProgram // 0x8822 ExposureProgram
	ExposureBias             meta.ExposureBias    // 0x9204 ExposureCompensation (ExposureBiasValue)
	ExposureMode             meta.ExposureMode    // 0xa402 ExposureMode
	MeteringMode             meta.MeteringMode    // 0x9207 MeteringMode
	Flash                    meta.Flash           // 0x9209 Flash
	ISOSpeedRatings          uint32               // 0x8827 ISO / ISOSpeedRatings
	RecommendedExposureIndex uint32               // 0x8832 RecommendedExposureIndex
	PixelXDimension          uint32               // 0xa002 ExifImageWidth
	PixelYDimension          uint32               // 0xa003 ExifImageHeight
	FocalPlaneResolutionUnit meta.ResolutionUnit  // 0xa210 FocalPlaneResolutionUnit
	ColorSpace               uint16               // 0xa001 ColorSpace
	LightSource              uint16               // 0x9208 LightSource
	CustomRendered           uint16               // 0xa401 CustomRendered
	WhiteBalance             uint16               // 0xa403 WhiteBalance
	SensingMethod            uint16               // 0xa217 SensingMethod
	FileSource               uint16               // 0xa300 FileSource
	SceneType                uint16               // 0xa301 SceneType
	GainControl              uint16               // 0xa407 GainControl
	Contrast                 uint16               // 0xa408 Contrast
	Saturation               uint16               // 0xa409 Saturation
	Sharpness                uint16               // 0xa40a Sharpness
	SubjectDistanceRange     uint16               // 0xa40c SubjectDistanceRange
	SceneCaptureType         uint16               // 0xa406 SceneCaptureType
	CompositeImage           uint16               // 0xa460 CompositeImage
	SensitivityType          uint16               // 0x8830 SensitivityType

	DigitalZoomRatio        float32   // 0xa404 DigitalZoomRatio
	SubjectArea             [4]uint16 // 0x9214 SubjectArea
	ComponentsConfiguration [4]byte   // 0x9101 ComponentsConfiguration
	Gamma                   float64   // 0xa500 Gamma
	InteropIndex            string    // 0x0001 InteropIndex
	InteropVersion          string    // 0x0002 InteropVersion
	RelatedImageWidth       uint32    // 0x1001 RelatedImageWidth
	RelatedImageHeight      uint32    // 0x1002 RelatedImageHeight

	focalPlaneXResolutionState unsignedRationalState
	focalPlaneYResolutionState unsignedRationalState
	subjectDistanceState       unsignedRationalState
	exposureIndexState         unsignedRationalState

	// IfdPointers
	interopIFDPointer uint32 // 0xa005 InteropOffset
}

// ImageIFD stores the core image-bearing tags from non-primary root IFDs.
type ImageIFD struct {
	IFD0TimeTags
	XResolution    float64
	YResolution    float64
	ResolutionUnit meta.ResolutionUnit
	Compression    meta.Compression
	Orientation    meta.Orientation

	SubfileType meta.SubfileType
	ImageOffset uint32
	ImageLength uint32
	ImageWidth  uint32
	ImageHeight uint32

	Make             string
	Model            string
	Software         string
	ImageDescription string
}

// LensInfo stores ExifIFD LensSpecification as four rationals.
type LensInfo struct {
	MinFocalLength        tag.RationalU
	MaxFocalLength        tag.RationalU
	MaxApertureAtMinFocal tag.RationalU
	MaxApertureAtMaxFocal tag.RationalU
}

// ExifToolDigitalZoomRatio returns DigitalZoomRatio formatted like ExifTool.
func (t ExifIFDTags) ExifToolDigitalZoomRatio() string {
	return exifToolUnsignedRationalValue(t.DigitalZoomRatio)
}

// ExifToolFocalPlaneXResolution returns FocalPlaneXResolution formatted like ExifTool.
func (t ExifIFDTags) ExifToolFocalPlaneXResolution() string {
	return exifToolUnsignedFloat64(t.FocalPlaneXResolution, t.focalPlaneXResolutionState)
}

// ExifToolFocalPlaneYResolution returns FocalPlaneYResolution formatted like ExifTool.
func (t ExifIFDTags) ExifToolFocalPlaneYResolution() string {
	return exifToolUnsignedFloat64(t.FocalPlaneYResolution, t.focalPlaneYResolutionState)
}

// ExifToolFocalPlaneResolutionUnit returns FocalPlaneResolutionUnit formatted like ExifTool.
func (t ExifIFDTags) ExifToolFocalPlaneResolutionUnit() string {
	if t.FocalPlaneResolutionUnit == 0 {
		return ""
	}
	return t.FocalPlaneResolutionUnit.String()
}

// ExifToolExposureIndex returns ExposureIndex formatted like ExifTool.
func (t ExifIFDTags) ExifToolExposureIndex() string {
	return exifToolUnsignedFloat64(t.ExposureIndex, t.exposureIndexState)
}

// ExifToolSubjectDistance returns SubjectDistance formatted like ExifTool.
func (t ExifIFDTags) ExifToolSubjectDistance() string {
	return exifToolDistanceFloat64(t.SubjectDistance, t.subjectDistanceState)
}

func (l LensInfo) String() string {
	return strings.Join([]string{
		l.lensInfoPart(l.MinFocalLength),
		l.lensInfoPart(l.MaxFocalLength),
		l.lensInfoPart(l.MaxApertureAtMinFocal),
		l.lensInfoPart(l.MaxApertureAtMaxFocal),
	}, " ")
}

// MarshalText emits an ExifTool -n style representation like "24 70 2.8 4".
func (l LensInfo) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (LensInfo) lensInfoPart(v tag.RationalU) string {
	if v.Denominator == 0 {
		return "undef"
	}
	f := v.Float64()
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func exifToolUnsignedFloat64(v float64, state unsignedRationalState) string {
	if state == rationalStateMissing && v != 0 {
		state = rationalStateValue
	}
	switch state {
	case rationalStateMissing:
		return ""
	case rationalStateInfinite:
		return "inf"
	default:
		return exifToolFloatString(v)
	}
}

func exifToolUnsignedRationalValue(v float32) string {
	return exifToolFloatString(float64(v))
}

func exifToolDistanceFloat64(v float64, state unsignedRationalState) string {
	if state == rationalStateMissing && v != 0 {
		state = rationalStateValue
	}
	switch state {
	case rationalStateMissing:
		return ""
	case rationalStateInfinite:
		return "undef"
	default:
		return exifToolFloatString(v) + " m"
	}
}

func exifToolFloatString(v float64) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'g', 10, 64)
}

// IFD0TimeTags contains time values parsed from IFD0-style DateTime tags.
type IFD0TimeTags struct {
	ModifyDate    time.Time
	OffsetTime    *time.Location
	SubSecTime    uint16
	ifdBitset     uint16
	modifyDateRaw time.Time
	subSecTimeRaw string
}

// ExifIFDTimeTags contains time values parsed from ExifIFD.
type ExifIFDTimeTags struct {
	DateTimeOriginal    time.Time
	CreateDate          time.Time
	OffsetTimeOriginal  *time.Location
	OffsetTimeDigitized *time.Location
	SubSecTimeOriginal  uint16
	SubSecTimeDigitized uint16
	ifdBitset           uint16
	dateTimeOriginalRaw time.Time
	createDateRaw       time.Time
	subSecTimeOrigRaw   string
	subSecTimeDigitRaw  string
}

type ifd0TimeJSON struct {
	ModifyDate time.Time `json:"ModifyDate"`
	OffsetTime *string   `json:"OffsetTime"`
	SubSecTime uint16    `json:"SubSecTime"`
}

type exifIFDTimeJSON struct {
	DateTimeOriginal    time.Time `json:"DateTimeOriginal"`
	CreateDate          time.Time `json:"CreateDate"`
	OffsetTimeOriginal  *string   `json:"OffsetTimeOriginal"`
	OffsetTimeDigitized *string   `json:"OffsetTimeDigitized"`
	SubSecTimeOriginal  uint16    `json:"SubSecTimeOriginal"`
	SubSecTimeDigitized uint16    `json:"SubSecTimeDigitized"`
}

type ifd0JSON struct {
	ModifyDate time.Time `json:"ModifyDate"`
	OffsetTime *string   `json:"OffsetTime"`
	SubSecTime uint16    `json:"SubSecTime"`

	XResolution               float64             `json:"XResolution"`
	YResolution               float64             `json:"YResolution"`
	Make                      string              `json:"Make"`
	Model                     string              `json:"Model"`
	Artist                    string              `json:"Artist"`
	Copyright                 string              `json:"Copyright"`
	Software                  string              `json:"Software"`
	ImageDescription          string              `json:"ImageDescription"`
	ImageWidth                uint32              `json:"ImageWidth"`
	ImageHeight               uint32              `json:"ImageHeight"`
	ImageOffset               uint32              `json:"ImageOffset"`
	ImageLength               uint32              `json:"ImageLength"`
	ThumbnailOffset           uint32              `json:"ThumbnailOffset"`
	ThumbnailLength           uint32              `json:"ThumbnailLength"`
	SubfileType               meta.SubfileType    `json:"SubfileType"`
	Compression               meta.Compression    `json:"Compression"`
	Orientation               meta.Orientation    `json:"Orientation"`
	ResolutionUnit            meta.ResolutionUnit `json:"ResolutionUnit"`
	PhotometricInterpretation uint16              `json:"PhotometricInterpretation"`
	SamplesPerPixel           uint16              `json:"SamplesPerPixel"`
	PlanarConfiguration       uint16              `json:"PlanarConfiguration"`
	BitsPerSample             [4]uint16           `json:"BitsPerSample"`
	WhitePoint                [2]float64          `json:"WhitePoint"`
	PrimaryChromaticities     [6]float64          `json:"PrimaryChromaticities"`
	YCbCrCoefficients         [3]float64          `json:"YCbCrCoefficients"`
	YCbCrPositioning          uint16              `json:"YCbCrPositioning"`
	Rating                    uint16              `json:"Rating"`
}

type exifIFDJSON struct {
	DateTimeOriginal    time.Time `json:"DateTimeOriginal"`
	CreateDate          time.Time `json:"CreateDate"`
	OffsetTimeOriginal  *string   `json:"OffsetTimeOriginal"`
	OffsetTimeDigitized *string   `json:"OffsetTimeDigitized"`
	SubSecTimeOriginal  uint16    `json:"SubSecTimeOriginal"`
	SubSecTimeDigitized uint16    `json:"SubSecTimeDigitized"`

	LensInfo LensInfo `json:"LensInfo"`

	ExifVersion              string               `json:"ExifVersion"`
	FlashpixVersion          string               `json:"FlashpixVersion"`
	FocalPlaneXResolution    float64              `json:"FocalPlaneXResolution"`
	FocalPlaneYResolution    float64              `json:"FocalPlaneYResolution"`
	CompressedBitsPerPixel   float64              `json:"CompressedBitsPerPixel"`
	SubjectDistance          float64              `json:"SubjectDistance"`
	ExposureIndex            float64              `json:"ExposureIndex"`
	LensMake                 string               `json:"LensMake"`
	LensModel                string               `json:"LensModel"`
	LensSerial               string               `json:"LensSerial"`
	CameraOwnerName          string               `json:"CameraOwnerName"`
	BodySerialNumber         string               `json:"BodySerialNumber"`
	UserComment              string               `json:"UserComment"`
	ExposureTime             meta.ExposureTime    `json:"ExposureTime"`
	FocalLength              meta.FocalLength     `json:"FocalLength"`
	FocalLengthIn35mmFormat  meta.FocalLength     `json:"FocalLengthIn35mmFormat"`
	FNumber                  meta.Aperture        `json:"FNumber"`
	ApertureValue            meta.Aperture        `json:"ApertureValue"`
	MaxApertureValue         meta.Aperture        `json:"MaxApertureValue"`
	ShutterSpeedValue        meta.ShutterSpeed    `json:"ShutterSpeedValue"`
	BrightnessValue          float32              `json:"BrightnessValue"`
	ExposureProgram          meta.ExposureProgram `json:"ExposureProgram"`
	ExposureBias             meta.ExposureBias    `json:"ExposureBias"`
	ExposureMode             meta.ExposureMode    `json:"ExposureMode"`
	MeteringMode             meta.MeteringMode    `json:"MeteringMode"`
	Flash                    meta.Flash           `json:"Flash"`
	ISOSpeedRatings          uint32               `json:"ISOSpeedRatings"`
	RecommendedExposureIndex uint32               `json:"RecommendedExposureIndex"`
	PixelXDimension          uint32               `json:"PixelXDimension"`
	PixelYDimension          uint32               `json:"PixelYDimension"`
	FocalPlaneResolutionUnit meta.ResolutionUnit  `json:"FocalPlaneResolutionUnit"`
	ColorSpace               uint16               `json:"ColorSpace"`
	LightSource              uint16               `json:"LightSource"`
	CustomRendered           uint16               `json:"CustomRendered"`
	WhiteBalance             uint16               `json:"WhiteBalance"`
	SensingMethod            uint16               `json:"SensingMethod"`
	FileSource               uint16               `json:"FileSource"`
	SceneType                uint16               `json:"SceneType"`
	GainControl              uint16               `json:"GainControl"`
	Contrast                 uint16               `json:"Contrast"`
	Saturation               uint16               `json:"Saturation"`
	Sharpness                uint16               `json:"Sharpness"`
	SubjectDistanceRange     uint16               `json:"SubjectDistanceRange"`
	SceneCaptureType         uint16               `json:"SceneCaptureType"`
	CompositeImage           uint16               `json:"CompositeImage"`
	SensitivityType          uint16               `json:"SensitivityType"`
	DigitalZoomRatio         float32              `json:"DigitalZoomRatio"`
	SubjectArea              [4]uint16            `json:"SubjectArea"`
	ComponentsConfiguration  [4]byte              `json:"ComponentsConfiguration"`
	Gamma                    float64              `json:"Gamma"`
	InteropIndex             string               `json:"InteropIndex"`
	InteropVersion           string               `json:"InteropVersion"`
	RelatedImageWidth        uint32               `json:"RelatedImageWidth"`
	RelatedImageHeight       uint32               `json:"RelatedImageHeight"`
}

type imageIFDJSON struct {
	ModifyDate time.Time `json:"ModifyDate"`
	OffsetTime *string   `json:"OffsetTime"`
	SubSecTime uint16    `json:"SubSecTime"`

	XResolution      float64             `json:"XResolution"`
	YResolution      float64             `json:"YResolution"`
	ResolutionUnit   meta.ResolutionUnit `json:"ResolutionUnit"`
	Compression      meta.Compression    `json:"Compression"`
	Orientation      meta.Orientation    `json:"Orientation"`
	SubfileType      meta.SubfileType    `json:"SubfileType"`
	ImageOffset      uint32              `json:"ImageOffset"`
	ImageLength      uint32              `json:"ImageLength"`
	ImageWidth       uint32              `json:"ImageWidth"`
	ImageHeight      uint32              `json:"ImageHeight"`
	Make             string              `json:"Make"`
	Model            string              `json:"Model"`
	Software         string              `json:"Software"`
	ImageDescription string              `json:"ImageDescription"`
}

type timeTagBitFunc func(tag.ID) uint16

const (
	timeTagBitModifyDate uint16 = 1 << iota
	timeTagBitDateTimeOriginal
	timeTagBitCreateDate
	timeTagBitSubSecTime
	timeTagBitSubSecTimeOriginal
	timeTagBitSubSecTimeDigitized
	timeTagBitOffsetTime
	timeTagBitOffsetTimeOriginal
	timeTagBitOffsetTimeDigitized
)

type unsignedRationalState uint8

const (
	rationalStateMissing unsignedRationalState = iota
	rationalStateValue
	rationalStateInfinite
)

// timeTagBit maps supported time tag IDs to TimeTags bitset positions.
func ifd0TimeTagBit(tagID tag.ID) uint16 {
	switch tagID {
	case tag.TagDateTime:
		return timeTagBitModifyDate
	case tag.TagSubSecTime:
		return timeTagBitSubSecTime
	case tag.TagOffsetTime:
		return timeTagBitOffsetTime
	default:
		return 0
	}
}

func exifIFDTimeTagBit(tagID tag.ID) uint16 {
	switch tagID {
	case tag.TagDateTimeOriginal:
		return timeTagBitDateTimeOriginal
	case tag.TagDateTimeDigitized:
		return timeTagBitCreateDate
	case tag.TagSubSecTimeOriginal:
		return timeTagBitSubSecTimeOriginal
	case tag.TagSubSecTimeDigitized:
		return timeTagBitSubSecTimeDigitized
	case tag.TagOffsetTimeOriginal:
		return timeTagBitOffsetTimeOriginal
	case tag.TagOffsetTimeDigitized:
		return timeTagBitOffsetTimeDigitized
	default:
		return 0
	}
}

// markTagParsed marks a time tag as parsed in the TimeTags bitset.
func (t *IFD0TimeTags) markTagParsed(tagID tag.ID) {
	markTimeTagParsed(&t.ifdBitset, tagID, ifd0TimeTagBit)
}

func (t *ExifIFDTimeTags) markTagParsed(tagID tag.ID) {
	markTimeTagParsed(&t.ifdBitset, tagID, exifIFDTimeTagBit)
}

// HasTagParsed reports whether a recognized EXIF time tag has been parsed.
func (t IFD0TimeTags) HasTagParsed(tagID tag.ID) bool {
	return hasTimeTagParsed(t.ifdBitset, tagID, ifd0TimeTagBit)
}

// HasTagParsed reports whether a recognized ExifIFD time tag has been parsed.
func (t ExifIFDTimeTags) HasTagParsed(tagID tag.ID) bool {
	return hasTimeTagParsed(t.ifdBitset, tagID, exifIFDTimeTagBit)
}

// TagParsedBitset returns the parsed EXIF time-tag bitset.
func (t IFD0TimeTags) TagParsedBitset() uint16 {
	return t.ifdBitset
}

// TagParsedBitset returns the parsed ExifIFD time-tag bitset.
func (t ExifIFDTimeTags) TagParsedBitset() uint16 {
	return t.ifdBitset
}

// MarshalJSON preserves EXIF-style UTC offset strings instead of serializing
// time.Location internals as empty objects.
func (t IFD0TimeTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(ifd0TimeJSON{
		ModifyDate: t.ModifyDate,
		OffsetTime: offsetTimeStringPtr(t.OffsetTime),
		SubSecTime: t.SubSecTime,
	})
}

// MarshalJSON merges flattened time fields with the rest of IFD0's exported fields.
func (t IFD0Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(ifd0JSON{
		ModifyDate:                t.ModifyDate,
		OffsetTime:                offsetTimeStringPtr(t.OffsetTime),
		SubSecTime:                t.SubSecTime,
		XResolution:               t.XResolution,
		YResolution:               t.YResolution,
		Make:                      t.Make,
		Model:                     t.Model,
		Artist:                    t.Artist,
		Copyright:                 t.Copyright,
		Software:                  t.Software,
		ImageDescription:          t.ImageDescription,
		ImageWidth:                t.ImageWidth,
		ImageHeight:               t.ImageHeight,
		ImageOffset:               t.ImageOffset,
		ImageLength:               t.ImageLength,
		ThumbnailOffset:           t.ThumbnailOffset,
		ThumbnailLength:           t.ThumbnailLength,
		SubfileType:               t.SubfileType,
		Compression:               t.Compression,
		Orientation:               t.Orientation,
		ResolutionUnit:            t.ResolutionUnit,
		PhotometricInterpretation: t.PhotometricInterpretation,
		SamplesPerPixel:           t.SamplesPerPixel,
		PlanarConfiguration:       t.PlanarConfiguration,
		BitsPerSample:             t.BitsPerSample,
		WhitePoint:                t.WhitePoint,
		PrimaryChromaticities:     t.PrimaryChromaticities,
		YCbCrCoefficients:         t.YCbCrCoefficients,
		YCbCrPositioning:          t.YCbCrPositioning,
		Rating:                    t.Rating,
	})
}

// MarshalJSON preserves EXIF-style UTC offset strings instead of serializing
// time.Location internals as empty objects.
func (t ExifIFDTimeTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(exifIFDTimeJSON{
		DateTimeOriginal:    t.DateTimeOriginal,
		CreateDate:          t.CreateDate,
		OffsetTimeOriginal:  offsetTimeStringPtr(t.OffsetTimeOriginal),
		OffsetTimeDigitized: offsetTimeStringPtr(t.OffsetTimeDigitized),
		SubSecTimeOriginal:  t.SubSecTimeOriginal,
		SubSecTimeDigitized: t.SubSecTimeDigitized,
	})
}

// MarshalJSON merges flattened time fields with the rest of ExifIFD's exported fields.
func (t ExifIFDTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(exifIFDJSON{
		DateTimeOriginal:         t.DateTimeOriginal,
		CreateDate:               t.CreateDate,
		OffsetTimeOriginal:       offsetTimeStringPtr(t.OffsetTimeOriginal),
		OffsetTimeDigitized:      offsetTimeStringPtr(t.OffsetTimeDigitized),
		SubSecTimeOriginal:       t.SubSecTimeOriginal,
		SubSecTimeDigitized:      t.SubSecTimeDigitized,
		LensInfo:                 t.LensInfo,
		ExifVersion:              t.ExifVersion.String(),
		FlashpixVersion:          t.FlashpixVersion,
		FocalPlaneXResolution:    t.FocalPlaneXResolution,
		FocalPlaneYResolution:    t.FocalPlaneYResolution,
		CompressedBitsPerPixel:   t.CompressedBitsPerPixel,
		SubjectDistance:          t.SubjectDistance,
		ExposureIndex:            t.ExposureIndex,
		LensMake:                 t.LensMake,
		LensModel:                t.LensModel,
		LensSerial:               t.LensSerial,
		CameraOwnerName:          t.CameraOwnerName,
		BodySerialNumber:         t.BodySerialNumber,
		UserComment:              t.UserComment,
		ExposureTime:             t.ExposureTime,
		FocalLength:              t.FocalLength,
		FocalLengthIn35mmFormat:  t.FocalLengthIn35mmFormat,
		FNumber:                  t.FNumber,
		ApertureValue:            t.ApertureValue,
		MaxApertureValue:         t.MaxApertureValue,
		ShutterSpeedValue:        t.ShutterSpeedValue,
		BrightnessValue:          t.BrightnessValue,
		ExposureProgram:          t.ExposureProgram,
		ExposureBias:             t.ExposureBias,
		ExposureMode:             t.ExposureMode,
		MeteringMode:             t.MeteringMode,
		Flash:                    t.Flash,
		ISOSpeedRatings:          t.ISOSpeedRatings,
		RecommendedExposureIndex: t.RecommendedExposureIndex,
		PixelXDimension:          t.PixelXDimension,
		PixelYDimension:          t.PixelYDimension,
		FocalPlaneResolutionUnit: t.FocalPlaneResolutionUnit,
		ColorSpace:               t.ColorSpace,
		LightSource:              t.LightSource,
		CustomRendered:           t.CustomRendered,
		WhiteBalance:             t.WhiteBalance,
		SensingMethod:            t.SensingMethod,
		FileSource:               t.FileSource,
		SceneType:                t.SceneType,
		GainControl:              t.GainControl,
		Contrast:                 t.Contrast,
		Saturation:               t.Saturation,
		Sharpness:                t.Sharpness,
		SubjectDistanceRange:     t.SubjectDistanceRange,
		SceneCaptureType:         t.SceneCaptureType,
		CompositeImage:           t.CompositeImage,
		SensitivityType:          t.SensitivityType,
		DigitalZoomRatio:         t.DigitalZoomRatio,
		SubjectArea:              t.SubjectArea,
		ComponentsConfiguration:  t.ComponentsConfiguration,
		Gamma:                    t.Gamma,
		InteropIndex:             t.InteropIndex,
		InteropVersion:           t.InteropVersion,
		RelatedImageWidth:        t.RelatedImageWidth,
		RelatedImageHeight:       t.RelatedImageHeight,
	})
}

// MarshalJSON merges flattened time fields with the rest of ImageIFD's exported fields.
func (t ImageIFD) MarshalJSON() ([]byte, error) {
	return json.Marshal(imageIFDJSON{
		ModifyDate:       t.ModifyDate,
		OffsetTime:       offsetTimeStringPtr(t.OffsetTime),
		SubSecTime:       t.SubSecTime,
		XResolution:      t.XResolution,
		YResolution:      t.YResolution,
		ResolutionUnit:   t.ResolutionUnit,
		Compression:      t.Compression,
		Orientation:      t.Orientation,
		SubfileType:      t.SubfileType,
		ImageOffset:      t.ImageOffset,
		ImageLength:      t.ImageLength,
		ImageWidth:       t.ImageWidth,
		ImageHeight:      t.ImageHeight,
		Make:             t.Make,
		Model:            t.Model,
		Software:         t.Software,
		ImageDescription: t.ImageDescription,
	})
}

// GetModifyDate returns the computed or normalized value.
func (t IFD0TimeTags) GetModifyDate() time.Time {
	return t.ModifyDate
}

// ExifToolModifyDate returns ModifyDate formatted like ExifTool's composite
// SubSecModifyDate output, using parsed subsecond digits and offset when present.
func (t IFD0TimeTags) ExifToolModifyDate() string {
	return formatExifToolDateTime(t.ModifyDate, t.subSecTimeRaw)
}

// GetDateTimeOriginal returns the computed or normalized value.
func (t ExifIFDTimeTags) GetDateTimeOriginal() time.Time {
	return t.DateTimeOriginal
}

// ExifToolDateTimeOriginal returns DateTimeOriginal formatted like ExifTool's
// composite SubSecDateTimeOriginal output.
func (t ExifIFDTimeTags) ExifToolDateTimeOriginal() string {
	return formatExifToolDateTime(t.DateTimeOriginal, t.subSecTimeOrigRaw)
}

// GetCreateDate returns the computed or normalized value.
func (t ExifIFDTimeTags) GetCreateDate() time.Time {
	return t.CreateDate
}

// ExifToolCreateDate returns CreateDate formatted like ExifTool's composite
// SubSecCreateDate output.
func (t ExifIFDTimeTags) ExifToolCreateDate() string {
	return formatExifToolDateTime(t.CreateDate, t.subSecTimeDigitRaw)
}

func offsetTimeStringPtr(loc *time.Location) *string {
	if loc == nil {
		return nil
	}
	s := loc.String()
	return &s
}

// HasSubSecTime reports whether the requested parsed value is present.
func (t IFD0TimeTags) HasSubSecTime() bool {
	return t.HasTagParsed(tag.TagSubSecTime)
}

// ExifToolSubSecTime returns the parsed SubSecTime digits.
func (t IFD0TimeTags) ExifToolSubSecTime() string {
	return exifToolSubSecString(t.subSecTimeRaw, t.SubSecTime)
}

// HasSubSecTimeOriginal reports whether the requested parsed value is present.
func (t ExifIFDTimeTags) HasSubSecTimeOriginal() bool {
	return t.HasTagParsed(tag.TagSubSecTimeOriginal)
}

// ExifToolSubSecTimeOriginal returns the parsed SubSecTimeOriginal digits.
func (t ExifIFDTimeTags) ExifToolSubSecTimeOriginal() string {
	return exifToolSubSecString(t.subSecTimeOrigRaw, t.SubSecTimeOriginal)
}

// HasSubSecTimeDigitized reports whether the requested parsed value is present.
func (t ExifIFDTimeTags) HasSubSecTimeDigitized() bool {
	return t.HasTagParsed(tag.TagSubSecTimeDigitized)
}

// ExifToolSubSecTimeDigitized returns the parsed SubSecTimeDigitized digits.
func (t ExifIFDTimeTags) ExifToolSubSecTimeDigitized() string {
	return exifToolSubSecString(t.subSecTimeDigitRaw, t.SubSecTimeDigitized)
}

func (t *IFD0TimeTags) setDate(date time.Time) {
	setTimeDate(&t.ModifyDate, &t.modifyDateRaw, date, t.subSecTimeRaw, t.OffsetTime)
}

func (t *IFD0TimeTags) setSubSec(raw string, value uint16) {
	setTimeSubSec(&t.ModifyDate, t.modifyDateRaw, &t.SubSecTime, &t.subSecTimeRaw, raw, value, t.OffsetTime)
}

func (t *IFD0TimeTags) setOffset(loc *time.Location) {
	setTimeOffset(&t.ModifyDate, t.modifyDateRaw, t.subSecTimeRaw, &t.OffsetTime, loc)
}

func (t *ExifIFDTimeTags) setDate(tagID tag.ID, date time.Time) {
	switch tagID {
	case tag.TagDateTimeOriginal:
		setTimeDate(&t.DateTimeOriginal, &t.dateTimeOriginalRaw, date, t.subSecTimeOrigRaw, t.OffsetTimeOriginal)
	case tag.TagDateTimeDigitized:
		setTimeDate(&t.CreateDate, &t.createDateRaw, date, t.subSecTimeDigitRaw, t.OffsetTimeDigitized)
	}
}

func (t *ExifIFDTimeTags) setSubSec(tagID tag.ID, raw string, value uint16) {
	switch tagID {
	case tag.TagSubSecTimeOriginal:
		setTimeSubSec(&t.DateTimeOriginal, t.dateTimeOriginalRaw, &t.SubSecTimeOriginal, &t.subSecTimeOrigRaw, raw, value, t.OffsetTimeOriginal)
	case tag.TagSubSecTimeDigitized:
		setTimeSubSec(&t.CreateDate, t.createDateRaw, &t.SubSecTimeDigitized, &t.subSecTimeDigitRaw, raw, value, t.OffsetTimeDigitized)
	}
}

func (t *ExifIFDTimeTags) setOffset(tagID tag.ID, loc *time.Location) {
	switch tagID {
	case tag.TagOffsetTimeOriginal:
		setTimeOffset(&t.DateTimeOriginal, t.dateTimeOriginalRaw, t.subSecTimeOrigRaw, &t.OffsetTimeOriginal, loc)
	case tag.TagOffsetTimeDigitized:
		setTimeOffset(&t.CreateDate, t.createDateRaw, t.subSecTimeDigitRaw, &t.OffsetTimeDigitized, loc)
	}
}

func markTimeTagParsed(bitset *uint16, tagID tag.ID, bitFn timeTagBitFunc) {
	if bit := bitFn(tagID); bit != 0 {
		*bitset |= bit
	}
}

func hasTimeTagParsed(bitset uint16, tagID tag.ID, bitFn timeTagBitFunc) bool {
	bit := bitFn(tagID)
	return bit != 0 && (bitset&bit) != 0
}

func setTimeDate(value *time.Time, rawValue *time.Time, date time.Time, subSecRaw string, tz *time.Location) {
	*rawValue = date
	*value = applyTimeParts(date, subSecRaw, tz)
}

func setTimeSubSec(value *time.Time, rawValue time.Time, parsed *uint16, rawStore *string, raw string, parsedValue uint16, tz *time.Location) {
	*parsed = parsedValue
	*rawStore = raw
	*value = applyTimeParts(rawValue, raw, tz)
}

func setTimeOffset(value *time.Time, rawValue time.Time, subSecRaw string, offset **time.Location, loc *time.Location) {
	*offset = loc
	*value = applyTimeParts(rawValue, subSecRaw, loc)
}

func parseSubSecNanoseconds(raw string) int {
	if raw == "" {
		return 0
	}
	nanos := 0
	digits := 0
	for i := 0; i < len(raw) && digits < 9; i++ {
		b := raw[i]
		if b < '0' || b > '9' {
			continue
		}
		nanos = nanos*10 + int(b-'0')
		digits++
	}
	for digits < 9 {
		nanos *= 10
		digits++
	}
	return nanos
}

func exifToolSubSecString(raw string, parsed uint16) string {
	if raw != "" {
		return raw
	}
	if parsed == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(parsed), 10)
}

func formatExifToolDateTime(date time.Time, subsecRaw string) string {
	if date.IsZero() {
		return ""
	}
	s := date.Format("2006:01:02 15:04:05")
	if subsec := exifToolSubSecString(subsecRaw, 0); subsec != "" {
		s += "." + subsec
	}
	if zone := date.Location().String(); zone != "" && zone != "UTC" {
		s += zone
	}
	return s
}

// applyTimeParts applies normalization logic to parsed time components.
func applyTimeParts(date time.Time, subsecRaw string, tz *time.Location) time.Time {
	if date.IsZero() {
		return date
	}
	loc := time.UTC
	if tz != nil {
		loc = tz
	}
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		date.Hour(),
		date.Minute(),
		date.Second(),
		parseSubSecNanoseconds(subsecRaw),
		loc,
	)
}

// ModifyDate returns the normalized top-level modify date.
func (e Exif) ModifyDate() time.Time {
	return e.IFD0.GetModifyDate()
}

// OriginalDate returns the normalized top-level original date.
func (e Exif) OriginalDate() time.Time {
	return e.ExifIFD.GetDateTimeOriginal()
}

// DigitizedDate returns the normalized top-level digitized date.
func (e Exif) DigitizedDate() time.Time {
	return e.ExifIFD.GetCreateDate()
}

// SelectedDate returns the preferred normalized date from original, create, then modify.
func (e Exif) SelectedDate() time.Time {
	if d := e.ExifIFD.GetDateTimeOriginal(); !d.IsZero() {
		return d
	}
	if d := e.ExifIFD.GetCreateDate(); !d.IsZero() {
		return d
	}
	return e.IFD0.GetModifyDate()
}

// CameraMake returns the preferred camera make value.
//
// When the normalized make enum is known, it returns the enum display string.
// Otherwise it falls back to the raw parsed IFD0 Make text.
func (e Exif) CameraMake() string {
	if e.CameraMakeID != makernote.CameraMakeUnknown {
		return e.CameraMakeID.String()
	}
	return e.IFD0.Make
}
