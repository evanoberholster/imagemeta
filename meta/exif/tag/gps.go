package tag

// GPSRationalRef combines a GPS rational value with its associated reference tag.
type GPSRationalRef[T any] struct {
	Ref   string
	Value T
}

// GPSRef stores normalized GPS reference values.
type GPSRef uint8

const (
	// GPSRefUnknown indicates an unsupported or invalid reference value.
	GPSRefUnknown GPSRef = iota
	// GPSRefNorth indicates north latitude.
	GPSRefNorth
	// GPSRefSouth indicates south latitude.
	GPSRefSouth
	// GPSRefEast indicates east longitude.
	GPSRefEast
	// GPSRefWest indicates west longitude.
	GPSRefWest
	// GPSRefAboveSeaLevel indicates altitude above sea level.
	GPSRefAboveSeaLevel
	// GPSRefBelowSeaLevel indicates altitude below sea level.
	GPSRefBelowSeaLevel
	// GPSRefKilometersPerHour indicates speed units in kilometers per hour.
	GPSRefKilometersPerHour
	// GPSRefMilesPerHour indicates speed units in miles per hour.
	GPSRefMilesPerHour
	// GPSRefKnots indicates speed units in knots.
	GPSRefKnots
	// GPSRefTrueDirection indicates true north reference.
	GPSRefTrueDirection
	// GPSRefMagneticDirection indicates magnetic north reference.
	GPSRefMagneticDirection
	// GPSRefKilometers indicates distance units in kilometers.
	GPSRefKilometers
	// GPSRefMiles indicates distance units in miles.
	GPSRefMiles
	// GPSRefNauticalMiles indicates distance units in nautical miles.
	GPSRefNauticalMiles
)

// String returns the EXIF string form of the GPS reference.
func (r GPSRef) String() string {
	switch r {
	case GPSRefNorth:
		return "N"
	case GPSRefSouth:
		return "S"
	case GPSRefEast:
		return "E"
	case GPSRefWest:
		return "W"
	case GPSRefAboveSeaLevel:
		return "0"
	case GPSRefBelowSeaLevel:
		return "1"
	case GPSRefKilometersPerHour:
		return "K"
	case GPSRefMilesPerHour:
		return "M"
	case GPSRefKnots:
		return "N"
	case GPSRefTrueDirection:
		return "T"
	case GPSRefMagneticDirection:
		return "M"
	case GPSRefKilometers:
		return "K"
	case GPSRefMiles:
		return "M"
	case GPSRefNauticalMiles:
		return "N"
	default:
		return ""
	}
}
