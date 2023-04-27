package tag

// GPSLatitudeRef represents the GPS Latitude Reference types
type GPSLatitudeRef uint8

// ExifTool will also accept a number when writing GPSLatitudeRef, positive for north latitudes or negative for south, or a string containing N, North, S or South.
const (
	GPSLatitudeRegUnknown GPSLatitudeRef = 0
	GPSLatitudeRefNorth   GPSLatitudeRef = 'N'
	GPSLatitudeRefSouth   GPSLatitudeRef = 'S'
)

// GPSLongitudeRef represents the GPS Longitude Reference types
type GPSLongitudeRef uint8

// ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West.
const (
	GPSLongitudeRefUnknown GPSLongitudeRef = 0
	GPSLongitudeRefEast    GPSLongitudeRef = 'E'
	GPSLongitudeRefWest    GPSLongitudeRef = 'W'
)

// TypeGPSDestBearingRef represents the GPS Destination Bearing Reference types
type GPSDestBearingRef uint8

// GPS Destination Bearing Reference types
const (
	GPSDestBearingRefUnknown   GPSDestBearingRef = 0
	GPSDestBearingRefMagNorth  GPSDestBearingRef = 'M'
	GPSDestBearingRefTrueNorth GPSDestBearingRef = 'T'
)

// GPSDestDistanceRef
type GPSDestDistanceRef uint8

// GPS Destination Distance Reference types
const (
	GPSDestDistanceRefUnknown GPSDestDistanceRef = 0
	GPSDestDistanceRefK       GPSDestDistanceRef = 'K'
	GPSDestDistanceRefM       GPSDestDistanceRef = 'M'
	GPSDestDistanceRefNM      GPSDestDistanceRef = 'N'
)

// GPSAltitudeRef
type GPSAltitudeRef uint8

const (
	GPSAltitudeRefAbove GPSAltitudeRef = 0
	GPSAltitudeRefBelow GPSAltitudeRef = 1
)

// GPSSpeedRef
type GPSSpeedRef uint8

// GPS Speed Reference types
const (
	GPSSpeedRefUnknown GPSSpeedRef = 0
	GPSSpeedRefK       GPSSpeedRef = 'K'
	GPSSpeedRefM       GPSSpeedRef = 'M'
	GPSSpeedRefN       GPSSpeedRef = 'N'
)

// GPSStatus
type GPSStatus uint8

// GPS Status
const (
	GPSStatusUnknown GPSStatus = 0
	GPSStatusA       GPSStatus = 'A'
	GPSStatusV       GPSStatus = 'V'
)

// GPSMeasureMode
type GPSMeasureMode uint8

// GPS Measure Mode
const (
	GPSMeasureModeUnknown GPSMeasureMode = 0
	GPSMeasureMode2       GPSMeasureMode = '2'
	GPSMeasureMode3       GPSMeasureMode = '3'
)
