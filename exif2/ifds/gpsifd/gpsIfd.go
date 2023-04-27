// Package gpsifd provides types for "RootIfd/GPSIfd"
package gpsifd

import (
	"github.com/evanoberholster/imagemeta/exif2/tag"
)

// TagString returns the string representation of a tag.ID
func TagString(id tag.ID) string {
	if int(id) < len(tagIDString) {
		return tagIDString[id]
	}
	return id.String()
}

// TagIDMap is a Map of tag.ID to string for the GPSIfd tags
var tagIDString = []string{"GPSVersionID", "GPSLatitudeRef", "GPSLatitude", "GPSLongitudeRef", "GPSLongitude", "GPSAltitudeRef", "GPSAltitude", "GPSTimeStamp", "GPSSatellites", "GPSStatus", "GPSMeasureMode", "GPSDOP", "GPSSpeedRef", "GPSSpeed", "GPSTrackRef", "GPSTrack", "GPSImgDirectionRef", "GPSImgDirection", "GPSMapDatum", "GPSDestLatitudeRef", "GPSDestLatitude", "GPSDestLongitudeRef", "GPSDestLongitude", "GPSDestBearingRef", "GPSDestBearing", "GPSDestDistanceRef", "GPSDestDistance", "GPSProcessingMethod", "GPSAreaInformation", "GPSDateStamp", "GPSDifferential", "GPSHPositioningError"}

// GPSInfo Tags; GPSInfo Ifd
const (
	GPSVersionID         tag.ID = 0x0000
	GPSLatitudeRef       tag.ID = 0x0001
	GPSLatitude          tag.ID = 0x0002
	GPSLongitudeRef      tag.ID = 0x0003
	GPSLongitude         tag.ID = 0x0004
	GPSAltitudeRef       tag.ID = 0x0005 // Altitude is expressed as one RATIONAL value. The reference unit is meters.
	GPSAltitude          tag.ID = 0x0006
	GPSTimeStamp         tag.ID = 0x0007
	GPSSatellites        tag.ID = 0x0008
	GPSStatus            tag.ID = 0x0009
	GPSMeasureMode       tag.ID = 0x000a
	GPSDOP               tag.ID = 0x000b
	GPSSpeedRef          tag.ID = 0x000c
	GPSSpeed             tag.ID = 0x000d
	GPSTrackRef          tag.ID = 0x000e
	GPSTrack             tag.ID = 0x000f // Indicates the direction of GPS receiver movement. The range of values is from 0.00 to 359.99.
	GPSImgDirectionRef   tag.ID = 0x0010
	GPSImgDirection      tag.ID = 0x0011
	GPSMapDatum          tag.ID = 0x0012
	GPSDestLatitudeRef   tag.ID = 0x0013
	GPSDestLatitude      tag.ID = 0x0014
	GPSDestLongitudeRef  tag.ID = 0x0015
	GPSDestLongitude     tag.ID = 0x0016
	GPSDestBearingRef    tag.ID = 0x0017
	GPSDestBearing       tag.ID = 0x0018 // Indicates the bearing to the destination point. The range of values is from 0.00 to 359.99.
	GPSDestDistanceRef   tag.ID = 0x0019
	GPSDestDistance      tag.ID = 0x001a
	GPSProcessingMethod  tag.ID = 0x001b
	GPSAreaInformation   tag.ID = 0x001c
	GPSDateStamp         tag.ID = 0x001d
	GPSDifferential      tag.ID = 0x001e
	GPSHPositioningError tag.ID = 0x001f
)

// GPSInfo struct as per Exiftool.
// These GPS tags are part of the EXIF standard, and are stored in a separate IFD within the EXIF information.
// Some GPS tags have values which are fixed-length strings. For these, the indicated string lengths include a null terminator.
// https://exiftool.org/TagNames/GPS.html
// https://web.archive.org/web/20190624045241if_/http://www.cipa.jp:80/std/documents/e/DC-008-Translation-2019-E.pdf
// Updated 04/09/2023
type GPSInfo struct {
	GPSMapDatum          string                 // 0x0012 string
	GPSSatellites        string                 // 0x0008 string
	GPSProcessingMethod  []byte                 // 0x001b undef	(values of "GPS", "CELLID", "WLAN" or "MANUAL" by the EXIF spec.)
	GPSAreaInformation   []byte                 // 0x001c undef
	GPSVersionID         [4]uint8               // 0x0000 int8u[4]
	GPSDateStamp         [11]byte               // 0x001d string[11] (time is stripped off if present, after adjusting date/time to UTC if time includes a timezone. Format is YYYY:mm:dd)
	GPSLatitude          [3]tag.Rational64      // 0x0002 rational64u[3]
	GPSLongitude         [3]tag.Rational64      // 0x0004 rational64u[3]
	GPSDestLatitude      [3]tag.Rational64      // 0x0014 rational64u[3]
	GPSDestLongitude     [3]tag.Rational64      // 0x0016 rational64u[3]
	GPSTimeStamp         [3]tag.Rational64      // 0x0007 rational64u[3] (UTC time of GPS fix)
	GPSAltitude          tag.Rational64         // 0x0006 rational64u
	GPSDOP               tag.Rational64         // 0x000b rational64u
	GPSSpeed             tag.Rational64         // 0x000d rational64u
	GPSTrack             tag.Rational64         // 0x000f rational64u
	GPSImgDirection      tag.Rational64         // 0x0011 rational64u
	GPSDestBearing       tag.Rational64         // 0x0018 rational64u
	GPSDestDistance      tag.Rational64         // 0x001a rational64u
	GPSHPositioningError tag.Rational64         // 0x001f rational64u
	GPSLatitudeRef       tag.GPSLatitudeRef     // 0x0001 string[2]
	GPSLongitudeRef      tag.GPSLongitudeRef    // 0x0003 string[2] (ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West) 'E' = East 'W' = West
	GPSAltitudeRef       tag.GPSAltitudeRef     // 0x0005 int8u (ExifTool will also accept number when writing this tag, with negative numbers indicating below sea level) 0 = Above Sea Level 1 = Below Sea Level
	GPSStatus            tag.GPSStatus          // 0x0009 string[2]	'A' = Measurement Active 'V' = Measurement Void
	GPSMeasureMode       tag.GPSMeasureMode     // 0x000a string[2]	2 = 2-Dimensional Measurement 3 = 3-Dimensional Measurement
	GPSSpeedRef          tag.GPSSpeedRef        // 0x000c string[2]	'K' = km/h 'M' = mph 'N' = knots
	GPSTrackRef          tag.GPSDestBearingRef  // 0x000e string[2]	'M' = Magnetic North 'T' = True North
	GPSImgDirectionRef   tag.GPSDestBearingRef  // 0x0010 string[2]	'M' = Magnetic North  'T' = True North
	GPSDestLatitudeRef   tag.GPSLatitudeRef     // 0x0013 string[2] (tags 0x0013-0x001a used for subject location according to MWG 2.0)  'N' = North  'S' = South
	GPSDestLongitudeRef  tag.GPSLongitudeRef    // 0x0015 string[2]	'E' = East  'W' = West
	GPSDestBearingRef    tag.GPSDestBearingRef  // 0x0017 string[2]	'M' = Magnetic North  'T' = True North
	GPSDestDistanceRef   tag.GPSDestDistanceRef // 0x0019 string[2]	'K' = Kilometers 'M' = Miles 'N' = Nautical Miles
	GPSDifferential      uint16                 // 0x001e int16u	0 = No Correction 1 = Differential Corrected
}

// Convenience functions

// Latitude returns GPSInfo Latitude as a float64.
func (gps GPSInfo) Latitude() float64 {
	coord := rational64ToCoord(gps.GPSLatitude)
	if gps.GPSLatitudeRef == tag.GPSLatitudeRefSouth {
		coord *= -1
	}
	return coord
}

// Longitude returns GPSInfo Longitude as a float64.
func (gps GPSInfo) Longitude() float64 {
	coord := rational64ToCoord(gps.GPSLongitude)
	if gps.GPSLongitudeRef == tag.GPSLongitudeRefWest {
		coord *= -1
	}
	return coord
}

// Altitude returns GPSInfo Altitude as a float64.
func (gps GPSInfo) Altitude() float64 {
	res := float64(gps.GPSAltitude[0]) / float64(gps.GPSAltitude[1])
	if gps.GPSAltitudeRef == tag.GPSAltitudeRefBelow {
		res *= -1
	}
	return res
}

// TODO: add GPSDatetime
// Time returns GPSInfo Time as a time.Time.
//func (gps GPSInfo) Time() time.Time {
//
//	return time.Now()
//}

// DestLatitude returns GPSInfo Destination Latitude as a float64.
func (gps GPSInfo) DestLatitude() float64 {
	coord := rational64ToCoord(gps.GPSDestLatitude)
	if gps.GPSDestLatitudeRef == tag.GPSLatitudeRefSouth {
		coord *= -1
	}
	return coord
}

// DestLongitude returns GPSInfo Destination Longitude as a float64.
func (gps GPSInfo) DestLongitude() float64 {
	coord := rational64ToCoord(gps.GPSDestLongitude)
	if gps.GPSDestLongitudeRef == tag.GPSLongitudeRefWest {
		coord *= -1
	}
	return coord
}

func rational64ToCoord(r [3]tag.Rational64) float64 {
	coord := float64(r[0][0]) / float64(r[0][1])
	coord += (float64(r[1][0]) / float64(r[1][1]) / 60.0)
	coord += (float64(r[2][0]) / float64(r[2][1]) / 3600.0)
	return coord
}
