package exif

import (
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

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
func (g GPSInfo) VersionID() string {
	switch g.versionID {
	case [4]byte{2, 0, 0, 0}:
		return "2.0.0.0"
	case [4]byte{2, 1, 0, 0}:
		return "2.1.0.0"
	case [4]byte{2, 2, 0, 0}:
		return "2.2.0.0"
	case [4]byte{2, 3, 0, 0}:
		return "2.3.0.0"
	default:
		s := [...]byte{g.versionID[0], '.', g.versionID[1], '.', g.versionID[2], '.', g.versionID[3]}
		return string(s[:])
	}
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
