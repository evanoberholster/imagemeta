package exif

import (
	"encoding/json"
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

type GPSVersion [4]byte

var (
	GPSVersionUnknown = GPSVersion{}
	GPSVersion2000    = GPSVersion{2, 0, 0, 0}
	GPSVersion2100    = GPSVersion{2, 1, 0, 0}
	GPSVersion2200    = GPSVersion{2, 2, 0, 0}
	GPSVersion2300    = GPSVersion{2, 3, 0, 0}
)

func (v GPSVersion) String() string {
	switch v {
	case GPSVersion2000:
		return "2000"
	case GPSVersion2100:
		return "2100"
	case GPSVersion2200:
		return "2200"
	case GPSVersion2300:
		return "2300"
	case GPSVersionUnknown:
		return ""
	default:
		key := [4]byte{v[0], v[1], v[2], v[3]}
		return string(key[:])
	}
}

// GPSInfo stores parsed GPS fields.
type GPSInfo struct {
	date              time.Time
	satellites        string
	status            string
	measureMode       string
	mapDatum          string
	processingMethod  string
	areaInformation   string
	latitude          float64
	longitude         float64
	destLatitude      float64
	destLongitude     float64
	dop               float64
	speed             float64
	track             float64
	imgDirection      float64
	destBearing       float64
	destDistance      float64
	hPositioningError float64
	altitude          float32
	GPSVersion        GPSVersion
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

const (
	knotsToKPH       = 1.852
	milesToKilometer = 1.60934
	kmToMeter        = 1000.0
	nauticalToMeter  = 1852.0
)

// gpsJSON is the JSON-serializable representation of GPSInfo.
type gpsJSON struct {
	Date             time.Time `json:"Date,omitempty"`
	Latitude         float64   `json:"Latitude,omitempty"`
	Longitude        float64   `json:"Longitude,omitempty"`
	Altitude         float32   `json:"Altitude,omitempty"`
	DestLatitude     float64   `json:"DestLatitude,omitempty"`
	DestLongitude    float64   `json:"DestLongitude,omitempty"`
	DOP              float64   `json:"DOP,omitempty"`
	SpeedRef         string    `json:"SpeedRef,omitempty"`
	Speed            float64   `json:"Speed,omitempty"`
	TrackRef         string    `json:"TrackRef,omitempty"`
	Track            float64   `json:"Track,omitempty"`
	ImgDirRef        string    `json:"ImgDirectionRef,omitempty"`
	ImgDirection     float64   `json:"ImgDirection,omitempty"`
	DestBearRef      string    `json:"DestBearingRef,omitempty"`
	DestBearing      float64   `json:"DestBearing,omitempty"`
	DestDistRef      string    `json:"DestDistanceRef,omitempty"`
	DestDistance     float64   `json:"DestDistance,omitempty"`
	HPosErr          float64   `json:"HPositioningError,omitempty"`
	GPSVersion       string    `json:"GPSVersion,omitempty"`
	Satellites       string    `json:"Satellites,omitempty"`
	Status           string    `json:"Status,omitempty"`
	MeasureMode      string    `json:"MeasureMode,omitempty"`
	MapDatum         string    `json:"MapDatum,omitempty"`
	ProcessingMethod string    `json:"ProcessingMethod,omitempty"`
	AreaInformation  string    `json:"AreaInformation,omitempty"`
	Differential     uint16    `json:"Differential,omitempty"`
}

// MarshalJSON implements json.Marshaler for GPSInfo.
func (g GPSInfo) MarshalJSON() ([]byte, error) {
	v := gpsJSON{
		Date:             g.date,
		Latitude:         g.Latitude(),
		Longitude:        g.Longitude(),
		Altitude:         g.Altitude(),
		DestLatitude:     g.DestLatitudeSigned(),
		DestLongitude:    g.DestLongitudeSigned(),
		DOP:              g.dop,
		SpeedRef:         g.speedRef.String(),
		Speed:            g.speed,
		TrackRef:         g.trackRef.String(),
		Track:            g.track,
		ImgDirRef:        g.imgDirectionRef.String(),
		ImgDirection:     g.imgDirection,
		DestBearRef:      g.destBearingRef.String(),
		DestBearing:      g.destBearing,
		DestDistRef:      g.destDistanceRef.String(),
		DestDistance:     g.destDistance,
		HPosErr:          g.hPositioningError,
		GPSVersion:       g.GPSVersion.String(),
		Satellites:       g.satellites,
		Status:           g.status,
		MeasureMode:      g.measureMode,
		MapDatum:         g.mapDatum,
		ProcessingMethod: g.processingMethod,
		AreaInformation:  g.areaInformation,
		Differential:     g.differential,
	}
	return json.Marshal(v)
}

// GPSTimestamp returns the combined GPS timestamp.
func (g GPSInfo) GPSTimestamp() time.Time {
	return g.date
}

// applyDate merges GPSDateStamp into the combined GPS timestamp state.
func (g *GPSInfo) applyDate(date time.Time) {
	if pending, ok := gpsPendingDelta(g.date); ok {
		g.date = date.Add(pending)
		return
	}
	g.date = date
}

// applyTime merges GPSTimeStamp into the combined GPS timestamp state.
func (g *GPSInfo) applyTime(delta time.Duration) {
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
	return g.latitude
}

// Longitude returns the signed longitude in decimal degrees.
func (g GPSInfo) Longitude() float64 {
	return g.longitude
}

// DestLatitudeSigned returns the signed destination latitude.
func (g GPSInfo) DestLatitudeSigned() float64 {
	return signedByRef(abs64(g.destLatitude), g.destLatitudeRef, tag.GPSRefSouth)
}

// DestLatitude returns the signed destination latitude.
func (g GPSInfo) DestLatitude() float64 {
	return g.DestLatitudeSigned()
}

// DestLongitudeSigned returns the signed destination longitude.
func (g GPSInfo) DestLongitudeSigned() float64 {
	return signedByRef(abs64(g.destLongitude), g.destLongitudeRef, tag.GPSRefWest)
}

// DestLongitude returns the signed destination longitude.
func (g GPSInfo) DestLongitude() float64 {
	return g.DestLongitudeSigned()
}

// Altitude returns the signed altitude in meters.
func (g GPSInfo) Altitude() float32 {
	return g.altitude
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
func (g GPSInfo) DOP() float64 {
	return g.dop
}

// HPositioningError returns the GPSHPositioningError value.
func (g GPSInfo) HPositioningError() float64 {
	return g.hPositioningError
}

// MapDatum returns the GPSMapDatum field.
func (g GPSInfo) MapDatum() string {
	return g.mapDatum
}

// ProcessingMethod returns the GPSProcessingMethod field.
func (g GPSInfo) ProcessingMethod() string {
	return g.processingMethod
}

// AreaInformation returns the GPSAreaInformation field.
func (g GPSInfo) AreaInformation() string {
	return g.areaInformation
}

// Differential returns the GPSDifferential field.
func (g GPSInfo) Differential() uint16 {
	return g.differential
}

// SpeedWithRef returns the speed with its reference unit.
func (g GPSInfo) SpeedWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.speedRef.String(),
		Value: tag.RationalU{Numerator: uint32(g.speed), Denominator: 1},
	}
}

// TrackWithRef returns the track with its reference.
func (g GPSInfo) TrackWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.trackRef.String(),
		Value: tag.RationalU{Numerator: uint32(g.track), Denominator: 1},
	}
}

// ImgDirectionWithRef returns the image direction with its reference.
func (g GPSInfo) ImgDirectionWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.imgDirectionRef.String(),
		Value: tag.RationalU{Numerator: uint32(g.imgDirection), Denominator: 1},
	}
}

// DestBearingWithRef returns the destination bearing with its reference.
func (g GPSInfo) DestBearingWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.destBearingRef.String(),
		Value: tag.RationalU{Numerator: uint32(g.destBearing), Denominator: 1},
	}
}

// DestDistanceWithRef returns the destination distance with its reference unit.
func (g GPSInfo) DestDistanceWithRef() tag.GPSRationalRef[tag.RationalU] {
	return tag.GPSRationalRef[tag.RationalU]{
		Ref:   g.destDistanceRef.String(),
		Value: tag.RationalU{Numerator: uint32(g.destDistance), Denominator: 1},
	}
}

func signedByRef(value float64, ref tag.GPSRef, negativeRef tag.GPSRef) float64 {
	if ref == negativeRef {
		return -value
	}
	return value
}

func signedAltitude(value float32, ref tag.GPSRef) float32 {
	if ref == tag.GPSRefBelowSeaLevel {
		return -value
	}
	return value
}

func speedInKPH(value float64, ref tag.GPSRef) float64 {
	switch ref {
	case tag.GPSRefKnots:
		return value * knotsToKPH
	case tag.GPSRefMilesPerHour:
		return value * milesToKilometer
	default:
		return value
	}
}

func distanceInMeters(value float64, ref tag.GPSRef) float64 {
	switch ref {
	case tag.GPSRefMiles:
		return value * milesToKilometer * kmToMeter
	case tag.GPSRefKilometers:
		return value * kmToMeter
	case tag.GPSRefNauticalMiles:
		return value * nauticalToMeter
	default:
		return value
	}
}

func abs64(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// parseGPSTag parses GPS IFD tags into typed model fields.
//
// Non-parsed GPS tags are currently handled by falling through to
// the default path (`return false`) when there is no modeled parser mapping.
func (r *Reader) parseGPSTag(t tag.Entry) bool {
	gps := &r.Exif.GPS
	switch t.ID {
	case tag.TagGPSVersionID:
		r.parseByteList(t, gps.GPSVersion[:])
	case tag.TagGPSDifferential:
		gps.differential = r.parseUint16(t)
	case tag.TagGPSAltitudeRef:
		gps.altitudeRef = r.parseGPSRef(t)
		gps.altitude = signedAltitude(abs32(gps.altitude), gps.altitudeRef)
	case tag.TagGPSLatitudeRef:
		gps.latitudeRef = r.parseGPSRef(t)
		gps.latitude = signedByRef(abs64(gps.latitude), gps.latitudeRef, tag.GPSRefSouth)
	case tag.TagGPSLongitudeRef:
		gps.longitudeRef = r.parseGPSRef(t)
		gps.longitude = signedByRef(abs64(gps.longitude), gps.longitudeRef, tag.GPSRefWest)
	case tag.TagGPSDestLatitudeRef:
		gps.destLatitudeRef = r.parseGPSRef(t)
		gps.destLatitude = signedByRef(abs64(gps.destLatitude), gps.destLatitudeRef, tag.GPSRefSouth)
	case tag.TagGPSDestLongitudeRef:
		gps.destLongitudeRef = r.parseGPSRef(t)
		gps.destLongitude = signedByRef(abs64(gps.destLongitude), gps.destLongitudeRef, tag.GPSRefWest)
	case tag.TagGPSSpeedRef:
		gps.speedRef = r.parseGPSRef(t)
	case tag.TagGPSTrackRef:
		gps.trackRef = r.parseGPSRef(t)
	case tag.TagGPSImgDirectionRef:
		gps.imgDirectionRef = r.parseGPSRef(t)
	case tag.TagGPSDestBearingRef:
		gps.destBearingRef = r.parseGPSRef(t)
	case tag.TagGPSDestDistanceRef:
		gps.destDistanceRef = r.parseGPSRef(t)
	case tag.TagGPSAltitude:
		gps.altitude = signedAltitude(abs32(r.parseGPSAltitude(t)), gps.altitudeRef)
	case tag.TagGPSLatitude:
		gps.latitude = signedByRef(abs64(r.parseGPSCoord(t)), gps.latitudeRef, tag.GPSRefSouth)
	case tag.TagGPSLongitude:
		gps.longitude = signedByRef(abs64(r.parseGPSCoord(t)), gps.longitudeRef, tag.GPSRefWest)
	case tag.TagGPSDestLatitude:
		gps.destLatitude = signedByRef(abs64(r.parseGPSCoord(t)), gps.destLatitudeRef, tag.GPSRefSouth)
	case tag.TagGPSDestLongitude:
		gps.destLongitude = signedByRef(abs64(r.parseGPSCoord(t)), gps.destLongitudeRef, tag.GPSRefWest)
	case tag.TagGPSSatellites:
		gps.satellites = r.parseString(t)
	case tag.TagGPSStatus:
		gps.status = r.parseString(t)
	case tag.TagGPSMeasureMode:
		gps.measureMode = r.parseString(t)
	case tag.TagGPSMapDatum:
		gps.mapDatum = r.parseString(t)
	case tag.TagGPSProcessingMethod:
		gps.processingMethod = r.parseExifUserComment(t)
	case tag.TagGPSAreaInformation:
		gps.areaInformation = r.parseExifUserComment(t)
	case tag.TagGPSDOP:
		gps.dop = r.parseRationalValue(t).Float64()
	case tag.TagGPSSpeed:
		gps.speed = speedInKPH(r.parseRationalValue(t).Float64(), gps.speedRef)
	case tag.TagGPSTrack:
		gps.track = r.parseRationalValue(t).Float64()
	case tag.TagGPSImgDirection:
		gps.imgDirection = r.parseRationalValue(t).Float64()
	case tag.TagGPSDestBearing:
		gps.destBearing = r.parseRationalValue(t).Float64()
	case tag.TagGPSDestDistance:
		gps.destDistance = distanceInMeters(r.parseRationalValue(t).Float64(), gps.destDistanceRef)
	case tag.TagGPSHPositioningError:
		gps.hPositioningError = r.parseRationalValue(t).Float64()
	case tag.TagGPSTimeStamp:
		gps.applyTime(r.parseGPSTimeStamp(t))
	case tag.TagGPSDateStamp:
		gps.applyDate(r.parseGPSDateStamp(t))
	default:
		return false
	}
	return true
}

func (r *Reader) parseGPSRef(t tag.Entry) tag.GPSRef {
	first, ok := r.firstTagByte(t)
	if !ok {
		return tag.GPSRefUnknown
	}

	switch t.ID {
	case tag.TagGPSAltitudeRef:
		// ExifTool GPS docs define 0/2 as above and 1/3 as below sea level.
		if first == 1 || first == 3 {
			return tag.GPSRefBelowSeaLevel
		}
		if first == 0 || first == 2 {
			return tag.GPSRefAboveSeaLevel
		}
		return tag.GPSRefUnknown
	case tag.TagGPSLatitudeRef, tag.TagGPSDestLatitudeRef:
		switch first | 0x20 {
		case 's':
			return tag.GPSRefSouth
		case 'n':
			return tag.GPSRefNorth
		default:
			return tag.GPSRefUnknown
		}
	case tag.TagGPSLongitudeRef, tag.TagGPSDestLongitudeRef:
		switch first | 0x20 {
		case 'w':
			return tag.GPSRefWest
		case 'e':
			return tag.GPSRefEast
		default:
			return tag.GPSRefUnknown
		}
	case tag.TagGPSSpeedRef:
		switch first | 0x20 {
		case 'k':
			return tag.GPSRefKilometersPerHour
		case 'm':
			return tag.GPSRefMilesPerHour
		case 'n':
			return tag.GPSRefKnots
		default:
			return tag.GPSRefUnknown
		}
	case tag.TagGPSTrackRef, tag.TagGPSImgDirectionRef, tag.TagGPSDestBearingRef:
		switch first | 0x20 {
		case 't':
			return tag.GPSRefTrueDirection
		case 'm':
			return tag.GPSRefMagneticDirection
		default:
			return tag.GPSRefUnknown
		}
	case tag.TagGPSDestDistanceRef:
		switch first | 0x20 {
		case 'k':
			return tag.GPSRefKilometers
		case 'm':
			return tag.GPSRefMiles
		case 'n':
			return tag.GPSRefNauticalMiles
		default:
			return tag.GPSRefUnknown
		}
	default:
		return tag.GPSRefUnknown
	}
}

func (r *Reader) firstTagByte(t tag.Entry) (byte, bool) {
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		return r.state.buf[0], true
	}
	buf, _, err := r.readTagBytes(t, 1)
	if err != nil || len(buf) == 0 {
		return 0, false
	}
	return buf[0], true
}

func (r *Reader) parseGPSCoord(t tag.Entry) float64 {
	if t.UnitCount != 3 {
		return 0
	}
	if !(t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	dNum := t.ByteOrder.Uint32(buf[:4])
	dDen := t.ByteOrder.Uint32(buf[4:8])
	mNum := t.ByteOrder.Uint32(buf[8:12])
	mDen := t.ByteOrder.Uint32(buf[12:16])
	sNum := t.ByteOrder.Uint32(buf[16:20])
	sDen := t.ByteOrder.Uint32(buf[20:24])
	if dDen == 0 || mDen == 0 || sDen == 0 {
		return 0
	}

	deg := float64(dNum) / float64(dDen)
	min := float64(mNum) / float64(mDen)
	sec := float64(sNum) / float64(sDen)
	return deg + min/60.0 + sec/3600.0
}

func (r *Reader) parseGPSAltitude(t tag.Entry) float32 {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return float32(rat[0]) / float32(rat[1])
}

func (r *Reader) parseGPSTimeStamp(t tag.Entry) time.Duration {
	if t.UnitCount != 3 || !t.IsType(tag.TypeRational) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	v := [6]uint32{
		t.ByteOrder.Uint32(buf[:4]),
		t.ByteOrder.Uint32(buf[4:8]),
		t.ByteOrder.Uint32(buf[8:12]),
		t.ByteOrder.Uint32(buf[12:16]),
		t.ByteOrder.Uint32(buf[16:20]),
		t.ByteOrder.Uint32(buf[20:24]),
	}
	return rationalDuration(v[0], v[1], time.Hour) +
		rationalDuration(v[2], v[3], time.Minute) +
		rationalDuration(v[4], v[5], time.Second)
}

func (r *Reader) parseGPSDateStamp(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 10 {
		return time.Time{}
	}
	sepA := buf[4]
	sepB := buf[7]
	if !((sepA == ':' && sepB == ':') || (sepA == 0 && sepB == 0)) {
		return time.Time{}
	}

	year := int(parseStrUint(buf[0:4]))
	month := time.Month(parseStrUint(buf[5:7]))
	day := int(parseStrUint(buf[8:10]))
	if len(buf) >= 19 && buf[10] == ' ' && buf[13] == ':' && buf[16] == ':' {
		return time.Date(
			year,
			month,
			day,
			int(parseStrUint(buf[11:13])),
			int(parseStrUint(buf[14:16])),
			int(parseStrUint(buf[17:19])),
			0,
			time.UTC,
		)
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
