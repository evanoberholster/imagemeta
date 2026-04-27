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

// gpsJSON is the JSON-serializable representation of GPSInfo.
type gpsJSON struct {
	Date          time.Time `json:"Date,omitempty"`
	Latitude      float64   `json:"Latitude,omitempty"`
	Longitude     float64   `json:"Longitude,omitempty"`
	Altitude      float32   `json:"Altitude,omitempty"`
	DestLatitude  float64   `json:"DestLatitude,omitempty"`
	DestLongitude float64   `json:"DestLongitude,omitempty"`
	DOP           float64   `json:"DOP,omitempty"`
	SpeedRef      string    `json:"SpeedRef,omitempty"`
	Speed         float64   `json:"Speed,omitempty"`
	TrackRef      string    `json:"TrackRef,omitempty"`
	Track         float64   `json:"Track,omitempty"`
	ImgDirRef     string    `json:"ImgDirectionRef,omitempty"`
	ImgDirection  float64   `json:"ImgDirection,omitempty"`
	DestBearRef   string    `json:"DestBearingRef,omitempty"`
	DestBearing   float64   `json:"DestBearing,omitempty"`
	DestDistRef   string    `json:"DestDistanceRef,omitempty"`
	DestDistance  float64   `json:"DestDistance,omitempty"`
	HPosErr       float64   `json:"HPositioningError,omitempty"`
	GPSVersion    string    `json:"GPSVersion,omitempty"`
	Satellites    string    `json:"Satellites,omitempty"`
	Status        string    `json:"Status,omitempty"`
	MeasureMode   string    `json:"MeasureMode,omitempty"`
	MapDatum      string    `json:"MapDatum,omitempty"`
	Differential  uint16    `json:"Differential,omitempty"`
}

// MarshalJSON implements json.Marshaler for GPSInfo.
func (g GPSInfo) MarshalJSON() ([]byte, error) {
	v := gpsJSON{
		Date:          g.date,
		Latitude:      g.Latitude(),
		Longitude:     g.Longitude(),
		Altitude:      g.Altitude(),
		DestLatitude:  g.DestLatitudeSigned(),
		DestLongitude: g.DestLongitudeSigned(),
		DOP:           g.dop,
		SpeedRef:      g.speedRef.String(),
		Speed:         g.speed,
		TrackRef:      g.trackRef.String(),
		Track:         g.track,
		ImgDirRef:     g.imgDirectionRef.String(),
		ImgDirection:  g.imgDirection,
		DestBearRef:   g.destBearingRef.String(),
		DestBearing:   g.destBearing,
		DestDistRef:   g.destDistanceRef.String(),
		DestDistance:  g.destDistance,
		HPosErr:       g.hPositioningError,
		GPSVersion:    g.GPSVersion.String(),
		Satellites:    g.satellites,
		Status:        g.status,
		MeasureMode:   g.measureMode,
		MapDatum:      g.mapDatum,
		Differential:  g.differential,
	}
	return json.Marshal(v)
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

// Latitude returns the raw latitude in decimal degrees (unsigned).
func (g GPSInfo) Latitude() float64 {
	return g.latitude
}

// Longitude returns the raw longitude in decimal degrees (unsigned).
func (g GPSInfo) Longitude() float64 {
	return g.longitude
}

// DestLatitudeSigned returns the signed destination latitude.
func (g GPSInfo) DestLatitudeSigned() float64 {
	if g.destLatitudeRef == tag.GPSRefSouth {
		return -1 * g.destLatitude
	}
	return g.destLatitude
}

// DestLatitude returns the raw destination latitude (unsigned).
func (g GPSInfo) DestLatitude() float64 {
	return g.destLatitude
}

// DestLongitudeSigned returns the signed destination longitude.
func (g GPSInfo) DestLongitudeSigned() float64 {
	if g.destLongitudeRef == tag.GPSRefWest {
		return -1 * g.destLongitude
	}
	return g.destLongitude
}

// DestLongitude returns the raw destination longitude (unsigned).
func (g GPSInfo) DestLongitude() float64 {
	return g.destLongitude
}

// Altitude returns the raw altitude value (unsigned).
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
	case tag.TagGPSLatitudeRef:
		gps.latitudeRef = r.parseGPSRef(t)
	case tag.TagGPSLongitudeRef:
		gps.longitudeRef = r.parseGPSRef(t)
	case tag.TagGPSDestLatitudeRef:
		gps.destLatitudeRef = r.parseGPSRef(t)
	case tag.TagGPSDestLongitudeRef:
		gps.destLongitudeRef = r.parseGPSRef(t)
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
		gps.altitude = r.parseGPSAltitude(t)
		if gps.altitudeRef == tag.GPSRefBelowSeaLevel {
			gps.altitude *= -1
		}
	case tag.TagGPSLatitude:
		gps.latitude = r.parseGPSCoord(t)
		if gps.latitudeRef == tag.GPSRefSouth {
			gps.latitude *= -1
		}
	case tag.TagGPSLongitude:
		gps.longitude = r.parseGPSCoord(t)
		if gps.longitudeRef == tag.GPSRefWest {
			gps.longitude *= -1
		}
	case tag.TagGPSDestLatitude:
		r.Exif.GPS.destLatitude = r.parseGPSCoord(t)
	case tag.TagGPSDestLongitude:
		r.Exif.GPS.destLongitude = r.parseGPSCoord(t)
	case tag.TagGPSSatellites:
		r.Exif.GPS.satellites = r.parseString(t)
	case tag.TagGPSStatus:
		r.Exif.GPS.status = r.parseString(t)
	case tag.TagGPSMeasureMode:
		r.Exif.GPS.measureMode = r.parseString(t)
	case tag.TagGPSMapDatum:
		r.Exif.GPS.mapDatum = r.parseString(t)
	case tag.TagGPSDOP:
		r.Exif.GPS.dop = r.parseRationalValue(t).Float64()
	case tag.TagGPSSpeed:
		gps.speed = r.parseRationalValue(t).Float64()
		if gps.speedRef == tag.GPSRefKnots {
			gps.speed *= 1.852 // Convert knots to km/h.
		}
		if gps.speedRef == tag.GPSRefMilesPerHour {
			gps.speed *= 1.60934 // Convert mph to km/h.
		}
	case tag.TagGPSTrack:
		r.Exif.GPS.track = r.parseRationalValue(t).Float64()

	case tag.TagGPSImgDirection:
		r.Exif.GPS.imgDirection = r.parseRationalValue(t).Float64()
	case tag.TagGPSDestBearing:
		r.Exif.GPS.destBearing = r.parseRationalValue(t).Float64()
	case tag.TagGPSDestDistance:
		gps.destDistance = r.parseRationalValue(t).Float64()
		if gps.destDistanceRef == tag.GPSRefMiles {
			gps.destDistance *= 1.60934 * 1000 // Convert miles to m.
		}
		if gps.destDistanceRef == tag.GPSRefKilometers {
			gps.destDistance *= 1000 // Convert km to m.
		}
		if gps.destDistanceRef == tag.GPSRefNauticalMiles {
			gps.destDistance *= 1852 // Convert nautical miles to m.
		}
	case tag.TagGPSHPositioningError:
		r.Exif.GPS.hPositioningError = r.parseRationalValue(t).Float64()
	case tag.TagGPSTimeStamp:
		r.Exif.GPS.setTime(r.parseGPSTimeStamp(t))
	case tag.TagGPSDateStamp:
		r.Exif.GPS.setDate(r.parseGPSDateStamp(t))
	default:
		return false
	}
	return true
}
