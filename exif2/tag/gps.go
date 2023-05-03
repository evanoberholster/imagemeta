package tag

import (
	"fmt"
	"math"
	"time"
)

var UnknownDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

// GPSVersionID is the GPS Version ID
type GPSVersionID [4]byte

// FromBytes parses the GPSVersionID from TagValue
func (vi *GPSVersionID) FromBytes(val TagValue) error {
	if len(val.Buf) == 4 {
		copy(vi[:], val.Buf[:4])
	}
	return nil
}

// GPSMapDatum is the GPS Map Datum
type GPSMapDatum string

const (
	GPSMapDatumUnknown = "Unknown"
	GPSMapDatumWGS84   = "WGS-84"
	GPSMapDatumTokyo   = "Tokyo"
)

// FromBytes the GPSMapDatum from TagValue
func (md *GPSMapDatum) FromBytes(val TagValue) error {
	switch string(val.Buf) {
	case GPSMapDatumUnknown, "UNKNOWN":
		*md = GPSMapDatumUnknown
	case GPSMapDatumTokyo, "TOKYO":
		*md = GPSMapDatumTokyo
	case GPSMapDatumWGS84, "WGS 84", "WGS84":
		*md = GPSMapDatumWGS84
	default:
	}
	return nil
}

// GPSProcessingMethod is the GPS Processing Method
type GPSProcessingMethod string

// GPS Processing Methods
const (
	GPSProcessingMethodGPS    = "GPS"
	GPSProcessingMethodCellID = "CELLID"
	GPSProcessingMethodWLAN   = "WLAN"
	GPSProcessingMethodManual = "Manual"
)

// NewGPSProcessingMethos returns the GPSProcessingMethod given []byte
func NewGPSProcessingMethod(buf []byte) GPSProcessingMethod {
	switch string(buf) {
	case GPSProcessingMethodGPS:
		return GPSProcessingMethodGPS

	case GPSProcessingMethodCellID:
		return GPSProcessingMethodCellID

	case GPSProcessingMethodWLAN:
		return GPSProcessingMethodWLAN

	case GPSProcessingMethodManual:
		return GPSProcessingMethodManual

	default:
		return GPSProcessingMethod(buf)
	}
}

// GPSLatitudeRef represents the GPS Latitude Reference types
type GPSLatitudeRef uint8

// ExifTool will also accept a number when writing GPSLatitudeRef, positive for north latitudes or negative for south, or a string containing N, North, S or South.
const (
	GPSLatitudeRegUnknown GPSLatitudeRef = 0
	GPSLatitudeRefNorth   GPSLatitudeRef = 'N'
	GPSLatitudeRefSouth   GPSLatitudeRef = 'S'
)

// Adjust retrurns the GPS Coordinate adjusted to GPSLatitudeRef
func (latR GPSLatitudeRef) Adjust(coord GPSCoordinate) float64 {
	res := coord.Float()
	if latR == GPSLatitudeRefSouth {
		res *= -1
	}
	if math.IsNaN(res) {
		return 0.0
	}
	return res
}

// GPSLongitudeRef represents the GPS Longitude Reference types
type GPSLongitudeRef uint8

// ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West.
const (
	GPSLongitudeRefUnknown GPSLongitudeRef = 0
	GPSLongitudeRefEast    GPSLongitudeRef = 'E'
	GPSLongitudeRefWest    GPSLongitudeRef = 'W'
)

// Adjust retrurns the GPS Coordinate adjusted to GPSLongitudeRef
func (lngR GPSLongitudeRef) Adjust(coord GPSCoordinate) float64 {
	res := coord.Float()
	if lngR == GPSLongitudeRefWest {
		res *= -1
	}
	if math.IsNaN(res) {
		return 0.0
	}
	return res
}

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

// GPSCoordinate
type GPSCoordinate [3]Rational

// NewGPSCoordinate returns a new GPSCoordinate
func NewGPSCoordinate(n0, d0, n1, d1, n2, d2 uint32) GPSCoordinate {
	return GPSCoordinate{
		Rational{n0, d0},
		Rational{n1, d1},
		Rational{n2, d2},
	}
}

func (c *GPSCoordinate) FromBytes(t TagValue) error {
	if t.UnitCount == 3 && len(t.Buf) >= 24 {
		switch t.Type {
		case TypeRational, TypeSignedRational: // Some cameras write tag out of spec using signed rational. We accept that too.
			*c = NewGPSCoordinate(
				t.ByteOrder.Uint32(t.Buf[:4]), t.ByteOrder.Uint32(t.Buf[4:8]),
				t.ByteOrder.Uint32(t.Buf[8:12]), t.ByteOrder.Uint32(t.Buf[12:16]),
				t.ByteOrder.Uint32(t.Buf[16:20]), t.ByteOrder.Uint32(t.Buf[20:24]))
		}
	}
	return nil
}

func (c GPSCoordinate) IsNil() bool {
	return c[0].Den() == 0 || c[1].Den() == 0
}

// Float returns a GPSCoordinate as a float64
func (c GPSCoordinate) Float() float64 {
	coord := float64(c[0][0]) / float64(c[0][1])
	coord += (float64(c[1][0]) / float64(c[1][1]) / 60.0)
	coord += (float64(c[2][0]) / float64(c[2][1]) / 3600.0)
	return coord
}

// GPSTimeStamp is a hour:min:sec value
type GPSTimeStamp [3]Rational

//func (ts GPSTimeStamp) String() string {
//	return fmt.Sprintf("%d:%d:%d", uint32(ts[0].Float()), uint32(ts[1].Float()), uint32(ts[2].Float()))
//}

// HourMinSec returns the GPSTimeStamp's Hour, Min, and Sec
func (ts GPSTimeStamp) HourMinSec() (hour int, min int, sec int) {
	return int(ts[0].Float()), int(ts[1].Float()), int(ts[2].Float())
}

// GPSDateStamp is a YYYYMMDD value
type GPSDateStamp struct {
	Year  uint16
	Month uint8
	Day   uint8
}

// NewGPSDateStamp returns a new GPSDateStamp
func NewGPSDateStamp(year uint16, month uint8, day uint8) GPSDateStamp {
	return GPSDateStamp{year, month, day}
}

func (ds GPSDateStamp) IsNil() bool {
	return !(ds.Year > 1000 && ds.Month > 0 && ds.Month < 13 && ds.Day > 0 && ds.Day < 31)
}
func (ds GPSDateStamp) String() string {
	return fmt.Sprintf("%d:%d:%d", ds.Year, ds.Month, ds.Day)
}

// GPSDifferential is the GPS Differential if it was corrected
type GPSDifferential uint16

const (
	GPSDifferentialNoCorrection   = 0
	GPSDifferentialWithCorrection = 1
)
