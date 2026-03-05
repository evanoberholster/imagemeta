package meta

import "time"

// GPSRef is a normalized reference/unit for GPS values.
type GPSRef uint8

// GPSRef values.
const (
	GPSRefUnknown GPSRef = iota
	GPSRefNorth
	GPSRefSouth
	GPSRefEast
	GPSRefWest
	GPSRefTrue
	GPSRefMagnetic
	GPSRefKilometers
	GPSRefMiles
	GPSRefKnots
	GPSRefAboveSeaLevel
	GPSRefBelowSeaLevel
)

// GPSCoordinate stores a coordinate value with an optional hemisphere reference.
// Value stores the absolute value in decimal degrees.
type GPSCoordinate struct {
	Value float64
	Ref   GPSRef
}

// Signed returns the signed decimal degree representation.
func (c GPSCoordinate) Signed() float64 {
	if c.Ref == GPSRefSouth || c.Ref == GPSRefWest {
		return -c.Value
	}
	return c.Value
}

// GPSMeasure stores a scalar GPS value with an optional reference/unit.
type GPSMeasure struct {
	Value float64
	Ref   GPSRef
}

// GPSAltitude stores altitude with an explicit altitude reference.
type GPSAltitude struct {
	Value float32
	Ref   GPSRef
}

// Signed returns altitude with altitude reference applied.
func (a GPSAltitude) Signed() float32 {
	if a.Ref == GPSRefBelowSeaLevel {
		return -a.Value
	}
	return a.Value
}

// GPS stores normalized GPS fields where base values and reference tags are combined.
type GPS struct {
	Latitude             GPSCoordinate
	Longitude            GPSCoordinate
	Altitude             GPSAltitude
	Time                 time.Time
	DestinationBearing   GPSMeasure
	DestinationDistance  GPSMeasure
	DestinationLatitude  GPSCoordinate
	DestinationLongitude GPSCoordinate
	ImageDirection       GPSMeasure
	Speed                GPSMeasure
	Track                GPSMeasure
	AreaInformation      string
	Differential         uint8
	HPositioningError    float64
	MapDatum             string
	ProcessingMethod     string
	Status               string
	DOP                  float64
	MeasureMode          string
	Satellites           string
	VersionID            string
}
