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
	GPSAltitudeRef       tag.ID = 0x0005
	GPSAltitude          tag.ID = 0x0006
	GPSTimeStamp         tag.ID = 0x0007
	GPSSatellites        tag.ID = 0x0008
	GPSStatus            tag.ID = 0x0009
	GPSMeasureMode       tag.ID = 0x000a
	GPSDOP               tag.ID = 0x000b
	GPSSpeedRef          tag.ID = 0x000c
	GPSSpeed             tag.ID = 0x000d
	GPSTrackRef          tag.ID = 0x000e
	GPSTrack             tag.ID = 0x000f
	GPSImgDirectionRef   tag.ID = 0x0010
	GPSImgDirection      tag.ID = 0x0011
	GPSMapDatum          tag.ID = 0x0012
	GPSDestLatitudeRef   tag.ID = 0x0013
	GPSDestLatitude      tag.ID = 0x0014
	GPSDestLongitudeRef  tag.ID = 0x0015
	GPSDestLongitude     tag.ID = 0x0016
	GPSDestBearingRef    tag.ID = 0x0017
	GPSDestBearing       tag.ID = 0x0018
	GPSDestDistanceRef   tag.ID = 0x0019
	GPSDestDistance      tag.ID = 0x001a
	GPSProcessingMethod  tag.ID = 0x001b
	GPSAreaInformation   tag.ID = 0x001c
	GPSDateStamp         tag.ID = 0x001d
	GPSDifferential      tag.ID = 0x001e
	GPSHPositioningError tag.ID = 0x001f
)

// TypeGPSLatitudeRef represents the GPS Latitude Reference types
type TypeGPSLatitudeRef uint8

// ExifTool will also accept a number when writing GPSLatitudeRef, positive for north latitudes or negative for south, or a string containing N, North, S or South.
const (
	GPSLatitudeRegUnknown TypeGPSLatitudeRef = 0
	GPSLatitudeRefNorth   TypeGPSLatitudeRef = 'N'
	GPSLatitudeRefSouth   TypeGPSLatitudeRef = 'S'
)

// TypeGPSLongitudeRef represents the GPS Longitude Reference types
type TypeGPSLongitudeRef uint8

// ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West.
const (
	GPSLongitudeRefUnknown TypeGPSLongitudeRef = 0
	GPSLongitudeRefEast    TypeGPSLongitudeRef = 'E'
	GPSLongitudeRefWest    TypeGPSLongitudeRef = 'W'
)

// TypeGPSDestBearingRef represents the GPS Destination Bearing Reference types
type TypeGPSDestBearingRef uint8

// GPS Destination Bearing Reference types
const (
	GPSDestBearingRefUnknown   TypeGPSDestBearingRef = 0
	GPSDestBearingRefMagNorth  TypeGPSDestBearingRef = 'M'
	GPSDestBearingRefTrueNorth TypeGPSDestBearingRef = 'T'
)

// GPSDestDistanceRef
type TypeGPSDestDistanceRef uint8

// GPS Destination Distance Reference types
const (
	GPSDestDistanceRefUnknown TypeGPSDestDistanceRef = 0
	GPSDestDistanceRefK       TypeGPSDestDistanceRef = 'K'
	GPSDestDistanceRefM       TypeGPSDestDistanceRef = 'M'
	GPSDestDistanceRefNM      TypeGPSDestDistanceRef = 'N'
)

// GPSAltitudeRef
type TypeGPSAltitudeRef uint8

const (
	GPSAltitudeRefAbove TypeGPSAltitudeRef = 0
	GPSAltitudeRefBelow TypeGPSAltitudeRef = 1
)

// GPSSpeedRef
type TypeGPSSpeedRef uint8

// GPS Speed Reference types
const (
	GPSSpeedRefUnknown TypeGPSSpeedRef = 0
	GPSSpeedRefK       TypeGPSSpeedRef = 'K'
	GPSSpeedRefM       TypeGPSSpeedRef = 'M'
	GPSSpeedRefN       TypeGPSSpeedRef = 'N'
)

// TypeGPSStaus
type TypeGPSStaus uint8

// GPS Status
const (
	GPSStatusUnknown TypeGPSStaus = 0
	GPSStatusA       TypeGPSStaus = 'A'
	GPSStatusV       TypeGPSStaus = 'V'
)

// TypeGPSMeasureMode
type TypeGPSMeasureMode uint8

// GPS Measure Mode
const (
	GPSMeasureModeUnknown TypeGPSMeasureMode = 0
	GPSMeasureMode2       TypeGPSMeasureMode = '2'
	GPSMeasureMode3       TypeGPSMeasureMode = '3'
)

// GPSInfo struct as per Exiftool.
// These GPS tags are part of the EXIF standard, and are stored in a separate IFD within the EXIF information.
// Some GPS tags have values which are fixed-length strings. For these, the indicated string lengths include a null terminator.
// https://exiftool.org/TagNames/GPS.html
// Updated 04/09/2023
type GPSInfo struct {
	GPSVersionID         [4]uint8               // 0x0000 int8u[4]
	GPSLatitudeRef       TypeGPSLatitudeRef     // 0x0001 string[2]
	GPSLatitude          [3]tag.Rational64      // 0x0002 rational64u[3]
	GPSLongitudeRef      TypeGPSLongitudeRef    // 0x0003 string[2] (ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West) 'E' = East 'W' = West
	GPSLongitude         [3]tag.Rational64      // 0x0004 rational64u[3]
	GPSAltitudeRef       TypeGPSAltitudeRef     // 0x0005 int8u (ExifTool will also accept number when writing this tag, with negative numbers indicating below sea level) 0 = Above Sea Level 1 = Below Sea Level
	GPSAltitude          tag.Rational64         // 0x0006 rational64u
	GPSTimeStamp         [3]tag.Rational64      // 0x0007 rational64u[3] (UTC time of GPS fix)
	GPSSatellites        string                 // 0x0008 string
	GPSStatus            TypeGPSStaus           // 0x0009 string[2]	'A' = Measurement Active 'V' = Measurement Void
	GPSMeasureMode       TypeGPSMeasureMode     // 0x000a string[2]	2 = 2-Dimensional Measurement 3 = 3-Dimensional Measurement
	GPSDOP               tag.Rational64         // 0x000b rational64u
	GPSSpeedRef          TypeGPSSpeedRef        // 0x000c string[2]	'K' = km/h 'M' = mph 'N' = knots
	GPSSpeed             tag.Rational64         // 0x000d rational64u
	GPSTrackRef          TypeGPSDestBearingRef  // 0x000e string[2]	'M' = Magnetic North 'T' = True North
	GPSTrack             tag.Rational64         // 0x000f rational64u
	GPSImgDirectionRef   TypeGPSDestBearingRef  // 0x0010 string[2]	'M' = Magnetic North  'T' = True North
	GPSImgDirection      tag.Rational64         // 0x0011 rational64u
	GPSMapDatum          string                 // 0x0012 string
	GPSDestLatitudeRef   TypeGPSLatitudeRef     // 0x0013 string[2] (tags 0x0013-0x001a used for subject location according to MWG 2.0)  'N' = North  'S' = South
	GPSDestLatitude      [3]tag.Rational64      // 0x0014 rational64u[3]
	GPSDestLongitudeRef  TypeGPSLongitudeRef    // 0x0015 string[2]	'E' = East  'W' = West
	GPSDestLongitude     [3]tag.Rational64      // 0x0016 rational64u[3]
	GPSDestBearingRef    TypeGPSDestBearingRef  // 0x0017 string[2]	'M' = Magnetic North  'T' = True North
	GPSDestBearing       tag.Rational64         // 0x0018 rational64u
	GPSDestDistanceRef   TypeGPSDestDistanceRef // 0x0019 string[2]	'K' = Kilometers 'M' = Miles 'N' = Nautical Miles
	GPSDestDistance      tag.Rational64         // 0x001a rational64u
	GPSProcessingMethod  []byte                 // 0x001b undef	(values of "GPS", "CELLID", "WLAN" or "MANUAL" by the EXIF spec.)
	GPSAreaInformation   []byte                 // 0x001c undef
	GPSDateStamp         [11]byte               // 0x001d string[11] (time is stripped off if present, after adjusting date/time to UTC if time includes a timezone. Format is YYYY:mm:dd)
	GPSDifferential      uint16                 // 0x001e int16u	0 = No Correction 1 = Differential Corrected
	GPSHPositioningError tag.Rational64         // 0x001f rational64u
}
