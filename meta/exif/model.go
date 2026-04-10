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
	Time         TimeTags
	IFD0         IFD0Tag
	ExifIFD      ExifIFDTags
	IFD1         ImageIFD
	IFD2         ImageIFD
	DNG          DNGTags
	MakerNote    makernote.Info
	CameraSerial string
	CameraMakeID makernote.CameraMake
	ImageType    imagetype.ImageType
}

// IFD0Tag groups tags from the primary image IFD (IFD0).
type IFD0Tag struct {
	ModifyDate       time.Time
	XResolution      tag.RationalU
	YResolution      tag.RationalU
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

	Compression    meta.Compression
	Orientation    meta.Orientation
	ResolutionUnit meta.ResolutionUnit
	Rating         uint16
	RatingPercent  uint16

	exifIfdPointer    uint32
	gpsIfdPointer     uint32
	subIFDOffsets     [8]uint32
	subIFDOffsetCount uint8
}

// ExifIFDTags groups tags from the ExifIFD directory.
// Field names and tag IDs follow ExifTool EXIF tag naming where applicable:
// https://exiftool.org/TagNames/EXIF.html
type ExifIFDTags struct {
	LensInfo *LensInfo // 0xa432 LensInfo (LensSpecification)
	// TODO: ExifVersion and FlashpixVersion are UNDEFINED in spec.
	// Keep canonical text form for parity with exiftool output.
	ExifVersion     string // 0x9000 ExifVersion
	FlashpixVersion string // 0xa000 FlashpixVersion

	FocalPlaneXResolution  tag.RationalU // 0xa20e FocalPlaneXResolution
	FocalPlaneYResolution  tag.RationalU // 0xa20f FocalPlaneYResolution
	CompressedBitsPerPixel tag.RationalU // 0x9102 CompressedBitsPerPixel
	SubjectDistance        tag.RationalU // 0x9206 SubjectDistance
	ExposureIndex          tag.RationalU // 0xa215 ExposureIndex

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
	InteropIFDPointer        uint32               // 0xa005 InteropOffset
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

	DigitalZoomRatio        tag.RationalU // 0xa404 DigitalZoomRatio
	SubjectArea             [4]uint16     // 0x9214 SubjectArea
	ComponentsConfiguration [4]byte       // 0x9101 ComponentsConfiguration
	// Time tags from ExifIFD are normalized into Exif.Time for a single source of truth.
}

// TIFFEPTags groups TIFF-EP extension tags commonly present in RAW containers.
type TIFFEPTags struct{}

// DNGTags groups Adobe DNG extension tags.
type DNGTags struct {
	DNGVersion         [8]byte
	DNGBackwardVersion [8]byte

	CameraModel         string
	OriginalRawFileName string
	ProfileName         string

	DNGVersionCount         uint8
	DNGBackwardVersionCount uint8

	BestQualityScale tag.RationalU
	AdobeData        DNGAdobeData
}

// DNGAdobeData stores selected information from IFD0 tag 0xc634
// (DNGAdobeData / Adobe private data).
//
// ExifTool parses this as an "Adobe\0" record stream. We currently model the
// overall record count plus the Adobe-mutated maker-note record details needed
// to rebase and parse MakN data.
type DNGAdobeData struct {
	RecordCount             uint8
	MakerNoteOriginalOffset uint32
	MakerNoteRecordLength   uint32
}

// PanasonicRawTags groups Panasonic RW2/RWL specific root-IFD tags.
// These tags are not part of standard EXIF/TIFF IFD0 and are modeled separately.
type PanasonicRawTags struct {
	Version [4]byte // 0x0001 PanasonicRawVersion

	RawDataOffset    uint32 // 0x0118 RawDataOffset
	JpgFromRawOffset uint32 // 0x002e JpgFromRaw (offset only)
	JpgFromRawLength uint32 // 0x002e JpgFromRaw (length only)
	ISO              uint32 // 0x0017/0x0037 ISO
	// TODO: PanasonicTitle fields are UNDEFINED and may contain mixed encodings.
	// Keep parsed printable strings for parity with exiftool dumps.
	Title  string // 0xc6d2 PanasonicTitle
	Title2 string // 0xc6d3 PanasonicTitle2

	SensorWidth   uint16 // 0x0002 SensorWidth
	SensorHeight  uint16 // 0x0003 SensorHeight
	BitsPerSample uint16 // 0x000a BitsPerSample
	Compression   uint16 // 0x000b Compression
	RawFormat     uint16 // 0x002d RawFormat
	CropTop       uint16 // 0x0121 CropTop
	CropLeft      uint16 // 0x0122 CropLeft
	CropBottom    uint16 // 0x0123 CropBottom
	CropRight     uint16 // 0x0124 CropRight
}

// ImageIFD stores the core image-bearing tags from non-primary root IFDs.
type ImageIFD struct {
	XResolution    tag.RationalU
	YResolution    tag.RationalU
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
	ModifyDate       time.Time
}

// LensInfo stores ExifIFD LensSpecification as four rationals.
type LensInfo struct {
	MinFocalLength        tag.RationalU
	MaxFocalLength        tag.RationalU
	MaxApertureAtMinFocal tag.RationalU
	MaxApertureAtMaxFocal tag.RationalU
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

// TimeTags contains parsed EXIF time values.
type TimeTags struct {
	ModifyDate          time.Time
	DateTimeOriginal    time.Time
	CreateDate          time.Time
	OffsetTime          *time.Location
	OffsetTimeOriginal  *time.Location
	OffsetTimeDigitized *time.Location
	SubSecTime          uint16
	SubSecTimeOriginal  uint16
	SubSecTimeDigitized uint16
	ifdBitset           uint16
}

type timeTagsJSON struct {
	ModifyDate          time.Time `json:"ModifyDate"`
	DateTimeOriginal    time.Time `json:"DateTimeOriginal"`
	CreateDate          time.Time `json:"CreateDate"`
	OffsetTime          *string   `json:"OffsetTime"`
	OffsetTimeOriginal  *string   `json:"OffsetTimeOriginal"`
	OffsetTimeDigitized *string   `json:"OffsetTimeDigitized"`
	SubSecTime          uint16    `json:"SubSecTime"`
	SubSecTimeOriginal  uint16    `json:"SubSecTimeOriginal"`
	SubSecTimeDigitized uint16    `json:"SubSecTimeDigitized"`
}

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

// timeTagBit maps supported time tag IDs to TimeTags bitset positions.
func timeTagBit(tagID tag.ID) uint16 {
	switch tagID {
	case tag.TagDateTime:
		return timeTagBitModifyDate
	case tag.TagDateTimeOriginal:
		return timeTagBitDateTimeOriginal
	case tag.TagDateTimeDigitized:
		return timeTagBitCreateDate
	case tag.TagSubSecTime:
		return timeTagBitSubSecTime
	case tag.TagSubSecTimeOriginal:
		return timeTagBitSubSecTimeOriginal
	case tag.TagSubSecTimeDigitized:
		return timeTagBitSubSecTimeDigitized
	case tag.TagOffsetTime:
		return timeTagBitOffsetTime
	case tag.TagOffsetTimeOriginal:
		return timeTagBitOffsetTimeOriginal
	case tag.TagOffsetTimeDigitized:
		return timeTagBitOffsetTimeDigitized
	default:
		return 0
	}
}

// markTagParsed marks a time tag as parsed in the TimeTags bitset.
func (t *TimeTags) markTagParsed(tagID tag.ID) {
	if bit := timeTagBit(tagID); bit != 0 {
		t.ifdBitset |= bit
	}
}

// HasTagParsed reports whether a recognized EXIF time tag has been parsed.
func (t TimeTags) HasTagParsed(tagID tag.ID) bool {
	bit := timeTagBit(tagID)
	return bit != 0 && (t.ifdBitset&bit) != 0
}

// TagParsedBitset returns the parsed EXIF time-tag bitset.
func (t TimeTags) TagParsedBitset() uint16 {
	return t.ifdBitset
}

// MarshalJSON preserves EXIF-style UTC offset strings instead of serializing
// time.Location internals as empty objects.
func (t TimeTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(timeTagsJSON{
		ModifyDate:          t.ModifyDate,
		DateTimeOriginal:    t.DateTimeOriginal,
		CreateDate:          t.CreateDate,
		OffsetTime:          offsetTimeStringPtr(t.OffsetTime),
		OffsetTimeOriginal:  offsetTimeStringPtr(t.OffsetTimeOriginal),
		OffsetTimeDigitized: offsetTimeStringPtr(t.OffsetTimeDigitized),
		SubSecTime:          t.SubSecTime,
		SubSecTimeOriginal:  t.SubSecTimeOriginal,
		SubSecTimeDigitized: t.SubSecTimeDigitized,
	})
}

// GetModifyDate returns the computed or normalized value.
func (t TimeTags) GetModifyDate() time.Time {
	return applyTimeParts(t.ModifyDate, t.SubSecTime, t.OffsetTime)
}

// GetDateTimeOriginal returns the computed or normalized value.
func (t TimeTags) GetDateTimeOriginal() time.Time {
	return applyTimeParts(t.DateTimeOriginal, t.SubSecTimeOriginal, t.OffsetTimeOriginal)
}

// GetCreateDate returns the computed or normalized value.
func (t TimeTags) GetCreateDate() time.Time {
	return applyTimeParts(t.CreateDate, t.SubSecTimeDigitized, t.OffsetTimeDigitized)
}

// GetSelectedDate returns the computed or normalized value.
func (t TimeTags) GetSelectedDate() time.Time {
	if d := t.GetDateTimeOriginal(); !d.IsZero() {
		return d
	}
	if d := t.GetCreateDate(); !d.IsZero() {
		return d
	}
	return t.GetModifyDate()
}

func offsetTimeStringPtr(loc *time.Location) *string {
	if loc == nil {
		return nil
	}
	s := loc.String()
	return &s
}

// HasSubSecTime reports whether the requested parsed value is present.
func (t TimeTags) HasSubSecTime() bool {
	return t.HasTagParsed(tag.TagSubSecTime)
}

// HasSubSecTimeOriginal reports whether the requested parsed value is present.
func (t TimeTags) HasSubSecTimeOriginal() bool {
	return t.HasTagParsed(tag.TagSubSecTimeOriginal)
}

// HasSubSecTimeDigitized reports whether the requested parsed value is present.
func (t TimeTags) HasSubSecTimeDigitized() bool {
	return t.HasTagParsed(tag.TagSubSecTimeDigitized)
}

// applyTimeParts applies normalization logic to parsed time components.
func applyTimeParts(date time.Time, subsec uint16, tz *time.Location) time.Time {
	if date.IsZero() {
		return date
	}
	date = date.Add(time.Duration(subsec) * time.Millisecond)
	if tz == nil {
		return date
	}
	date = date.In(tz)
	_, offset := date.Zone()
	return date.Add(time.Duration(offset) * -time.Second)
}

// ModifyDate returns the normalized top-level modify date.
func (e Exif) ModifyDate() time.Time {
	return e.Time.GetModifyDate()
}

// OriginalDate returns the normalized top-level original date.
func (e Exif) OriginalDate() time.Time {
	return e.Time.GetDateTimeOriginal()
}

// DigitizedDate returns the normalized top-level digitized date.
func (e Exif) DigitizedDate() time.Time {
	return e.Time.GetCreateDate()
}

// SelectedDate returns the preferred normalized date from original, create, then modify.
func (e Exif) SelectedDate() time.Time {
	return e.Time.GetSelectedDate()
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
