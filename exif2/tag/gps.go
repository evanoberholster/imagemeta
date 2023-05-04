package tag

import (
	"fmt"
	"math"
	"time"
)

// UnknownDate is unsed when a Date is unknown or incorrectly formed
var UnknownDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

// GPSVersionID is the GPS Version ID
type GPSVersionID [4]byte

// FromBytes parses the GPSVersionID from TagValue
func (vi *GPSVersionID) FromBytes(val TagValue) error {
	if len(val.Buf) >= 1 {
		copy(vi[:], val.Buf[:])
	}
	return nil
}

func (vi *GPSVersionID) String() string {
	if vi[0] != 0 {
		return fmt.Sprintf("%d.%d.%d.%d", vi[0], vi[1], vi[2], vi[3])
	}
	return ""
}

// GPSMapDatum is the GPS Map Datum
type GPSMapDatum string

// GPSMapDatum values
const (
	GPSMapDatumUnknown = "Unknown"
	GPSMapDatumWGS84   = "WGS-84"
	GPSMapDatumTokyo   = "Tokyo"
)

// FromBytes parses the GPSMapDatum from TagValue
func (md *GPSMapDatum) FromBytes(val TagValue) error {
	buf := trimNULBuffer(val.Buf)
	switch string(buf) {
	case GPSMapDatumUnknown, "UNKNOWN":
		*md = GPSMapDatumUnknown
	case GPSMapDatumTokyo, "TOKYO":
		*md = GPSMapDatumTokyo
	case GPSMapDatumWGS84, "WGS 84", "WGS84":
		*md = GPSMapDatumWGS84
	case "":
	default:
		*md = GPSMapDatum(buf)
	}
	return nil
}

// GPSSatellites are the GPS satellites used for measurement
type GPSSatellites string

// Most common number of GPSSatellites used
const (
	GPSSatellites0  = "0"
	GPSSatellites4  = "4"
	GPSSatellites5  = "5"
	GPSSatellites6  = "6"
	GPSSatellites7  = "7"
	GPSSatellites8  = "8"
	GPSSatellites9  = "9"
	GPSSatellites10 = "10"
	GPSSatellites11 = "11"
	GPSSatellites12 = "12"
	GPSSatellites13 = "13"
	GPSSatellites14 = "14"
	GPSSatellites15 = "15"
	GPSSatellites16 = "16"
)

// FromBytes parses the GPSSatellites from TagValue
func (s *GPSSatellites) FromBytes(val TagValue) error {
	buf := trimNULBuffer(val.Buf)
	switch string(buf) {
	case GPSSatellites4:
		*s = GPSSatellites4
	case GPSSatellites5:
		*s = GPSSatellites5
	case GPSSatellites6:
		*s = GPSSatellites6
	case GPSSatellites7:
		*s = GPSSatellites7
	case GPSSatellites8:
		*s = GPSSatellites8
	case GPSSatellites9:
		*s = GPSSatellites9
	case GPSSatellites10:
		*s = GPSSatellites10
	case GPSSatellites11:
		*s = GPSSatellites11
	case GPSSatellites12:
		*s = GPSSatellites12
	case GPSSatellites13:
		*s = GPSSatellites13
	case GPSSatellites14:
		*s = GPSSatellites14
	case GPSSatellites15:
		*s = GPSSatellites15
	case GPSSatellites16:
		*s = GPSSatellites16
	case "":
	default:
		*s = GPSSatellites(buf)
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

// FromBytes parses the GPSProcessingMethod from TagValue
func (pm *GPSProcessingMethod) FromBytes(val TagValue) error {
	buf := trimNULBuffer(val.Buf)
	switch string(buf) {
	case GPSProcessingMethodGPS:
		*pm = GPSProcessingMethodGPS
	case GPSProcessingMethodCellID:
		*pm = GPSProcessingMethodCellID
	case GPSProcessingMethodWLAN:
		*pm = GPSProcessingMethodWLAN
	case GPSProcessingMethodManual:
		*pm = GPSProcessingMethodManual
	case "":
	default:
		*pm = GPSProcessingMethod(buf)
	}
	return nil
}

// GPSLatitudeRef represents the GPS Latitude Reference types
type GPSLatitudeRef uint8

// FromBytes parses the GPSLatitudeRef from TagValue
func (lr *GPSLatitudeRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*lr = GPSLatitudeRef(val.Buf[0])
		}
	}
	return nil
}

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

// FromBytes parses the GPSLongitudeRef from TagValue
func (lr *GPSLongitudeRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*lr = GPSLongitudeRef(val.Buf[0])
		}
	}
	return nil
}

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

// FromBytes parses the GPSDestBearingRef from TagValue
func (dbr *GPSDestBearingRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*dbr = GPSDestBearingRef(val.Buf[0])
		}
	}
	return nil
}

// GPS Destination Bearing Reference types
const (
	GPSDestBearingRefUnknown   GPSDestBearingRef = 0
	GPSDestBearingRefMagNorth  GPSDestBearingRef = 'M'
	GPSDestBearingRefTrueNorth GPSDestBearingRef = 'T'
)

// GPSDestDistanceRef
type GPSDestDistanceRef uint8

// FromBytes parses the GPSDestDistanceRef from TagValue
func (ddr *GPSDestDistanceRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*ddr = GPSDestDistanceRef(val.Buf[0])
		}
	}
	return nil
}

// GPS Destination Distance Reference types
const (
	GPSDestDistanceRefUnknown GPSDestDistanceRef = 0
	GPSDestDistanceRefK       GPSDestDistanceRef = 'K'
	GPSDestDistanceRefM       GPSDestDistanceRef = 'M'
	GPSDestDistanceRefNM      GPSDestDistanceRef = 'N'
)

// GPSAltitudeRef
type GPSAltitudeRef uint8

// FromBytes parses the GPSAltitudeRef from TagValue
func (ar *GPSAltitudeRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*ar = GPSAltitudeRef(val.Buf[0])
		}
	}
	return nil
}

const (
	GPSAltitudeRefAbove GPSAltitudeRef = 0
	GPSAltitudeRefBelow GPSAltitudeRef = 1
)

// GPSSpeedRef
type GPSSpeedRef uint8

// FromBytes parses the GPSSpeedRef from TagValue
func (sr *GPSSpeedRef) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*sr = GPSSpeedRef(val.Buf[0])
		}
	}
	return nil
}

// GPS Speed Reference types
const (
	GPSSpeedRefUnknown GPSSpeedRef = 0
	GPSSpeedRefK       GPSSpeedRef = 'K'
	GPSSpeedRefM       GPSSpeedRef = 'M'
	GPSSpeedRefN       GPSSpeedRef = 'N'
)

// GPSStatus
type GPSStatus uint8

// FromBytes parses the GPSStatus from TagValue
func (s *GPSStatus) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*s = GPSStatus(val.Buf[0])
		}
	}
	return nil
}

// GPS Status
const (
	GPSStatusUnknown GPSStatus = 0
	GPSStatusA       GPSStatus = 'A'
	GPSStatusV       GPSStatus = 'V'
)

// GPSMeasureMode
type GPSMeasureMode uint8

// FromBytes parses the GPSMeasureMode from TagValue
func (mm *GPSMeasureMode) FromBytes(val TagValue) error {
	if len(val.Buf) > 0 {
		switch val.Type {
		case TypeASCII, TypeByte:
			*mm = GPSMeasureMode(val.Buf[0])
		}
	}
	return nil
}

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

// FromBytes parses a GPSCoordinate from TagValue
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

// GPSTimeStamp is a hour:min:sec value (UTC time of GPS fix)
type GPSTimeStamp [3]Rational

// FromBytes parses a GPSTimeStamp from TagValue
func (ts *GPSTimeStamp) FromBytes(val TagValue) error {
	switch val.Type {
	case TypeRational, TypeSignedRational: // Some cameras write tag out of spec using signed rational. We accept that too.
		if val.UnitCount == 3 && len(val.Buf) >= 24 {
			*ts = GPSTimeStamp{
				Rational{val.ByteOrder.Uint32(val.Buf[:4]), val.ByteOrder.Uint32(val.Buf[4:8])},
				Rational{val.ByteOrder.Uint32(val.Buf[8:12]), val.ByteOrder.Uint32(val.Buf[12:16])},
				Rational{val.ByteOrder.Uint32(val.Buf[16:20]), val.ByteOrder.Uint32(val.Buf[20:24])}}
			return nil
		}
	}
	return nil
}

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

// FromBytes parses a GPSDateStamp from TagValue
// (time is stripped off if present, after adjusting date/time to UTC if time includes a timezone. Format is YYYY:mm:dd)
func (ds *GPSDateStamp) FromBytes(val TagValue) error {
	if val.Type.Is(TypeASCII) {
		if len(val.Buf) == 11 && val.Buf[4] == ':' && val.Buf[7] == ':' {
			*ds = GPSDateStamp{uint16(parseStrUint(val.Buf[0:4])), uint8(parseStrUint(val.Buf[5:7])), uint8(parseStrUint(val.Buf[8:10]))}
			return nil
		}
		// check recieved value
		if len(val.Buf) > 19 && val.Buf[4] == ':' && val.Buf[7] == ':' && val.Buf[10] == ' ' &&
			val.Buf[13] == ':' && val.Buf[16] == ':' {
			*ds = GPSDateStamp{uint16(parseStrUint(val.Buf[0:4])), uint8(parseStrUint(val.Buf[5:7])), uint8(parseStrUint(val.Buf[8:10]))}
			return nil
		}
	}
	return nil
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
