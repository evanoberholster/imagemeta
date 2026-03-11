package exif

import (
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
	IFD0         IFD0Tags
	ExifIFD      ExifIFDTags
	IFD1         ImageIFD
	IFD2         ImageIFD
	PanasonicRaw PanasonicRawTags
	DNG          DNGTags
	MakerNote    makernote.Info
	CameraMakeID makernote.CameraMake
	CameraSerial string
	ImageType    imagetype.ImageType

	ifdBitset    [8]uint64
	highTagIDs   [128]uint16
	highTagCount uint8
}

const exifBitsetMaxTagID uint16 = (8 * 64) - 1

// HasTagParsed reports whether a tag ID has been parsed into typed fields.
func (e Exif) HasTagParsed(tagID uint16) bool {
	if tagID <= exifBitsetMaxTagID {
		word := tagID >> 6
		mask := uint64(1) << (tagID & 63)
		return (e.ifdBitset[word] & mask) != 0
	}
	n := min(int(e.highTagCount), len(e.highTagIDs))
	for i := range n {
		if e.highTagIDs[i] == tagID {
			return true
		}
	}
	return false
}

// TagParsedBitset returns the internal parsed-tag bitset.
func (e Exif) TagParsedBitset() [8]uint64 {
	return e.ifdBitset
}

// markTagParsed marks a tag ID as parsed in the EXIF-level bitset.
func (e *Exif) markTagParsed(tagID uint16) {
	if tagID <= exifBitsetMaxTagID {
		word := tagID >> 6
		e.ifdBitset[word] |= uint64(1) << (tagID & 63)
		return
	}
	n := min(int(e.highTagCount), len(e.highTagIDs))
	for i := range n {
		if e.highTagIDs[i] == tagID {
			return
		}
	}
	if n >= len(e.highTagIDs) {
		return
	}
	e.highTagIDs[n] = tagID
	e.highTagCount++
}

// IFD0Tags groups tags from the primary image IFD (IFD0).
type IFD0Tags struct {
	XResolution   tag.RationalU
	YResolution   tag.RationalU
	BitsPerSample [8]uint16
	SubIFDOffsets [8]uint32

	Make             string
	Model            string
	Artist           string
	Copyright        string
	Software         string
	ImageDescription string
	// TODO: ApplicationNotes (UNDEFINED) may contain large payloads (for example XMP).
	// Parsing is currently disabled to avoid unnecessary allocations.

	StripOffsets    uint32
	StripByteCounts uint32
	ThumbnailOffset uint32
	ThumbnailLength uint32
	TileWidth       uint32
	TileLength      uint32
	TileOffsets     uint32
	TileByteCounts  uint32
	SubfileType     uint32
	SR2Private      uint32
	RowsPerStrip    uint32
	ExifIFDPointer  uint32
	GPSIFDPointer   uint32
	ImageWidth      uint32
	ImageHeight     uint32

	Compression         meta.Compression
	Orientation         meta.Orientation
	PlanarConfiguration uint16
	ResolutionUnit      uint16
	BitsPerSampleCount  uint8
	SubIFDOffsetCount   uint8

	DateTimeOriginal time.Time
	ModifyDate       time.Time

	// TODO: PrintIM payload is UNDEFINED and may include mixed binary/text content.
	PrintIM string
}

// ExifIFDTags groups tags from the ExifIFD directory.
// Field names and tag IDs follow ExifTool EXIF tag naming where applicable:
// https://exiftool.org/TagNames/EXIF.html
type ExifIFDTags struct {
	LensInfo [8]uint32 // 0xa432 LensInfo (LensSpecification)
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

	ExposureTime             meta.ExposureTime    // 0x829a ExposureTime
	FocalLength              meta.FocalLength     // 0x920a FocalLength
	FocalLengthIn35mmFormat  meta.FocalLength     // 0xa405 FocalLengthIn35mmFormat
	FNumber                  meta.Aperture        // 0x829d FNumber
	ApertureValue            meta.Aperture        // 0x9202 ApertureValue (APEX units)
	MaxApertureValue         meta.Aperture        // 0x9205 MaxApertureValue (APEX units)
	ShutterSpeedValue        meta.ShutterSpeed    // 0x9201 ShutterSpeedValue
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
	FocalPlaneResolutionUnit uint16               // 0xa210 FocalPlaneResolutionUnit
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
	// TODO: DeviceSettingDescription is UNDEFINED and often vendor-specific binary data.
	// Keep printable representation for parity with exiftool dumps.
	DeviceSettingDescription string // 0xa40b DeviceSettingDescription
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
	XResolution        tag.RationalU
	YResolution        tag.RationalU
	BitsPerSample      [8]uint16
	Compression        meta.Compression
	SubfileType        uint32
	ThumbnailOffset    uint32
	ThumbnailLength    uint32
	StripOffsets       uint32
	StripByteCounts    uint32
	ImageWidth         uint32
	ImageHeight        uint32
	Make               string
	Model              string
	Software           string
	ImageDescription   string
	ModifyDate         time.Time
	Orientation        meta.Orientation
	ResolutionUnit     uint16
	BitsPerSampleCount uint8
}

// LensInfo is equivalent to ExifIFD LensSpecification packed as rationals.
type LensInfo [8]uint32

// GPSInfo stores parsed GPS fields.
type GPSInfo struct {
	date              time.Time
	satellites        string
	status            string
	measureMode       string
	mapDatum          string
	latitude          float64
	longitude         float64
	destLatitude      float64
	destLongitude     float64
	ifdBitset         uint64
	dop               tag.RationalU
	speed             tag.RationalU
	track             tag.RationalU
	imgDirection      tag.RationalU
	destBearing       tag.RationalU
	destDistance      tag.RationalU
	hPositioningError tag.RationalU
	altitude          float32
	versionID         [4]byte
	differential      uint16
	speedRef          tag.GPSRef
	trackRef          tag.GPSRef
	imgDirectionRef   tag.GPSRef
	destLatitudeRef   tag.GPSRef
	destLongitudeRef  tag.GPSRef
	destBearingRef    tag.GPSRef
	destDistanceRef   tag.GPSRef
	latitudeRef       tag.GPSRef
	longitudeRef      tag.GPSRef
	altitudeRef       tag.GPSRef
}

// Date returns the combined GPS timestamp.
func (g GPSInfo) Date() time.Time {
	return g.GPSTimestamp()
}

// GPSTimestamp returns the combined GPS timestamp.
func (g GPSInfo) GPSTimestamp() time.Time {
	return g.date
}

// GPSTime returns the combined GPS timestamp.
// Deprecated: use GPSTimestamp.
func (g GPSInfo) GPSTime() time.Time {
	return g.GPSTimestamp()
}

// setDate sets the internal state value used during parsing.
func (g *GPSInfo) setDate(date time.Time) {
	if pending, ok := gpsPendingDelta(g.date); ok {
		g.date = date.Add(pending)
		return
	}
	g.date = date
}

// setTime sets the internal state value used during parsing.
func (g *GPSInfo) setTime(delta time.Duration) {
	if delta == 0 {
		return
	}
	if g.date.IsZero() {
		// Store pending GPS time without adding another struct field.
		g.date = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC).Add(delta)
		return
	}
	g.date = g.date.Add(delta)
}

// gpsPendingDelta extracts a pending GPS time offset encoded in sentinel date form.
func gpsPendingDelta(ts time.Time) (time.Duration, bool) {
	if ts.IsZero() {
		return 0, false
	}
	if ts.Year() != 1 || ts.Month() != time.January || ts.Day() != 1 {
		return 0, false
	}
	return time.Duration(ts.Hour())*time.Hour +
		time.Duration(ts.Minute())*time.Minute +
		time.Duration(ts.Second())*time.Second +
		time.Duration(ts.Nanosecond()), true
}

// Latitude returns the signed latitude in decimal degrees.
func (g GPSInfo) Latitude() float64 {
	if g.latitudeRef == tag.GPSRefSouth {
		return -1 * g.latitude
	}
	return g.latitude
}

// Longitude returns the signed longitude in decimal degrees.
func (g GPSInfo) Longitude() float64 {
	if g.longitudeRef == tag.GPSRefWest {
		return -1 * g.longitude
	}
	return g.longitude
}

// DestLatitude returns the signed destination latitude in decimal degrees.
func (g GPSInfo) DestLatitude() float64 {
	if g.destLatitudeRef == tag.GPSRefSouth {
		return -1 * g.destLatitude
	}
	return g.destLatitude
}

// DestLongitude returns the signed destination longitude in decimal degrees.
func (g GPSInfo) DestLongitude() float64 {
	if g.destLongitudeRef == tag.GPSRefWest {
		return -1 * g.destLongitude
	}
	return g.destLongitude
}

// Altitude returns the signed altitude value.
func (g GPSInfo) Altitude() float32 {
	if g.altitudeRef == tag.GPSRefBelowSeaLevel {
		return -1 * g.altitude
	}
	return g.altitude
}

// VersionID returns the GPSVersionID tuple.
func (g GPSInfo) VersionID() [4]byte {
	return g.versionID
}

// Satellites returns the GPSSatellites field.
func (g GPSInfo) Satellites() string {
	return g.satellites
}

// Status returns the GPSStatus field.
func (g GPSInfo) Status() string {
	return g.status
}

// MeasureMode returns the GPSMeasureMode field.
func (g GPSInfo) MeasureMode() string {
	return g.measureMode
}

// DOP returns the parsed GPSDOP rational value.
func (g GPSInfo) DOP() tag.RationalU {
	return g.dop
}

// SpeedWithRef returns GPSSpeed together with GPSSpeedRef.
func (g GPSInfo) SpeedWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.speedRef.String(),
		Value: g.speed,
	}
}

// TrackWithRef returns GPSTrack together with GPSTrackRef.
func (g GPSInfo) TrackWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.trackRef.String(),
		Value: g.track,
	}
}

// ImgDirectionWithRef returns GPSImgDirection together with GPSImgDirectionRef.
func (g GPSInfo) ImgDirectionWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.imgDirectionRef.String(),
		Value: g.imgDirection,
	}
}

// DestBearingWithRef returns GPSDestBearing together with GPSDestBearingRef.
func (g GPSInfo) DestBearingWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.destBearingRef.String(),
		Value: g.destBearing,
	}
}

// DestDistanceWithRef returns GPSDestDistance together with GPSDestDistanceRef.
func (g GPSInfo) DestDistanceWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.destDistanceRef.String(),
		Value: g.destDistance,
	}
}

// HPositioningError returns the GPSHPositioningError value.
func (g GPSInfo) HPositioningError() tag.RationalU {
	return g.hPositioningError
}

// MapDatum returns the GPSMapDatum field.
func (g GPSInfo) MapDatum() string {
	return g.mapDatum
}

// Differential returns the GPSDifferential field.
func (g GPSInfo) Differential() uint16 {
	return g.differential
}

// HasTagParsed reports whether a GPS tag ID was parsed.
// GPS tags are addressed by their raw tag ID (for example 0x0002 for GPSLatitude).
func (g GPSInfo) HasTagParsed(tagID uint16) bool {
	if tagID > 63 {
		return false
	}
	return (g.ifdBitset & (uint64(1) << tagID)) != 0
}

// TagParsedBitset returns the GPS parsed-tag bitset.
func (g GPSInfo) TagParsedBitset() uint64 {
	return g.ifdBitset
}

// markTagParsed marks a GPS tag ID as parsed in the GPS-level bitset.
func (g *GPSInfo) markTagParsed(id tag.ID) {
	if id > 63 {
		return
	}
	g.ifdBitset |= uint64(1) << id
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
