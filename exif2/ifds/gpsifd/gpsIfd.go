// Package gpsifd provides types for "RootIfd/GPSIfd"
package gpsifd

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

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
var tagIDString = [32]string{"GPSVersionID", "GPSLatitudeRef", "GPSLatitude", "GPSLongitudeRef", "GPSLongitude", "GPSAltitudeRef", "GPSAltitude", "GPSTimeStamp", "GPSSatellites", "GPSStatus", "GPSMeasureMode", "GPSDOP", "GPSSpeedRef", "GPSSpeed", "GPSTrackRef", "GPSTrack", "GPSImgDirectionRef", "GPSImgDirection", "GPSMapDatum", "GPSDestLatitudeRef", "GPSDestLatitude", "GPSDestLongitudeRef", "GPSDestLongitude", "GPSDestBearingRef", "GPSDestBearing", "GPSDestDistanceRef", "GPSDestDistance", "GPSProcessingMethod", "GPSAreaInformation", "GPSDateStamp", "GPSDifferential", "GPSHPositioningError"}

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
	GPSMapDatum          tag.GPSMapDatum         `tagid:"0x0012" json:"MapDatum,omitempty" msgpack:"0"`          //  string
	GPSSatellites        tag.GPSSatellites       `tagid:"0x0008" json:"Satellites,omitempty" msgpack:"1"`        //  string
	GPSProcessingMethod  tag.GPSProcessingMethod `tagid:"0x001b" json:"ProcessingMethod,omitempty" msgpack:"2"`  //  undef	(values of "GPS", "CELLID", "WLAN" or "MANUAL" by the EXIF spec.)
	GPSAreaInformation   []byte                  `tagid:"0x001c" json:"AreaInformation,omitempty" msgpack:"3"`   //  undef
	GPSVersionID         tag.GPSVersionID        `tagid:"0x0000" json:"VersionID,omitempty" msgpack:"4"`         //  int8u[4]
	GPSDateStamp         tag.GPSDateStamp        `tagid:"0x001d" json:"DateStamp,omitempty" msgpack:"5"`         //  string[11] (time is stripped off if present, after adjusting date/time to UTC if time includes a timezone. Format is YYYY:mm:dd)
	GPSLatitude          tag.GPSCoordinate       `tagid:"0x0002" json:"Latitude,omitempty" msgpack:"6"`          //  rational64u[3]
	GPSLongitude         tag.GPSCoordinate       `tagid:"0x0004" json:"Longitude,omitempty" msgpack:"7"`         //  rational64u[3]
	GPSDestLatitude      tag.GPSCoordinate       `tagid:"0x0014" json:"DestLatitude,omitempty" msgpack:"8"`      //  rational64u[3]
	GPSDestLongitude     tag.GPSCoordinate       `tagid:"0x0016" json:"DestLongitude,omitempty" msgpack:"9"`     //  rational64u[3]
	GPSTimeStamp         tag.GPSTimeStamp        `tagid:"0x0007" json:"TimeStamp,omitempty" msgpack:"a"`         //  rational64u[3] (UTC time of GPS fix)
	GPSAltitude          tag.Rational            `tagid:"0x0006" json:"Altitude,omitempty" msgpack:"b"`          //  rational64u
	GPSDOP               tag.Rational            `tagid:"0x000b" json:"DOP,omitempty" msgpack:"c"`               //  rational64u
	GPSSpeed             tag.Rational            `tagid:"0x000d" json:"Speed,omitempty" msgpack:"d"`             //  rational64u
	GPSTrack             tag.Rational            `tagid:"0x000f" json:"Track,omitempty" msgpack:"e"`             //  rational64u
	GPSImgDirection      tag.Rational            `tagid:"0x0011" json:"ImgDirection,omitempty" msgpack:"f"`      //  rational64u
	GPSDestBearing       tag.Rational            `tagid:"0x0018" json:"DestBearing,omitempty" msgpack:"g"`       //  rational64u
	GPSDestDistance      tag.Rational            `tagid:"0x001a" json:"DestDistance,omitempty" msgpack:"h"`      //  rational64u
	GPSHPositioningError tag.Rational            `tagid:"0x001f" json:"HPositioningError,omitempty" msgpack:"i"` //  rational64u
	GPSLatitudeRef       tag.GPSLatitudeRef      `tagid:"0x0001" json:"LatitudeRef,omitempty" msgpack:"j"`       //  string[2] 'E' = East  'W' = West
	GPSLongitudeRef      tag.GPSLongitudeRef     `tagid:"0x0003" json:"LongitudeRef,omitempty" msgpack:"k"`      //  string[2] (ExifTool will also accept a number when writing this tag, positive for east longitudes or negative for west, or a string containing E, East, W or West) 'E' = East 'W' = West
	GPSAltitudeRef       tag.GPSAltitudeRef      `tagid:"0x0005" json:"AltitudeRef,omitempty" msgpack:"l"`       //  int8u (ExifTool will also accept number when writing this tag, with negative numbers indicating below sea level) 0 = Above Sea Level 1 = Below Sea Level
	GPSStatus            tag.GPSStatus           `tagid:"0x0009" json:"Status,omitempty" msgpack:"m"`            //  string[2]	'A' = Measurement Active 'V' = Measurement Void
	GPSMeasureMode       tag.GPSMeasureMode      `tagid:"0x000a" json:"MeasureMode,omitempty" msgpack:"n"`       //  string[2]	2 = 2-Dimensional Measurement 3 = 3-Dimensional Measurement
	GPSSpeedRef          tag.GPSSpeedRef         `tagid:"0x000c" json:"SpeedRef,omitempty" msgpack:"o"`          //  string[2]	'K' = km/h 'M' = mph 'N' = knots
	GPSTrackRef          tag.GPSDestBearingRef   `tagid:"0x000e" json:"TrackRef,omitempty" msgpack:"p"`          //  string[2]	'M' = Magnetic North 'T' = True North
	GPSImgDirectionRef   tag.GPSDestBearingRef   `tagid:"0x0010" json:"ImgDirectionRef,omitempty" msgpack:"q"`   //  string[2]	'M' = Magnetic North  'T' = True North
	GPSDestLatitudeRef   tag.GPSLatitudeRef      `tagid:"0x0013" json:"DestLatitudeRef,omitempty" msgpack:"r"`   //  string[2] (tags 0x0013-0x001a used for subject location according to MWG 2.0)  'N' = North  'S' = South
	GPSDestLongitudeRef  tag.GPSLongitudeRef     `tagid:"0x0015" json:"DestLongitudeRef,omitempty" msgpack:"s"`  //  string[2]	'E' = East  'W' = West
	GPSDestBearingRef    tag.GPSDestBearingRef   `tagid:"0x0017" json:"DestBearingRef,omitempty" msgpack:"t"`    //  string[2]	'M' = Magnetic North  'T' = True North
	GPSDestDistanceRef   tag.GPSDestDistanceRef  `tagid:"0x0019" json:"DestDistanceRef,omitempty" msgpack:"u"`   //  string[2]	'K' = Kilometers 'M' = Miles 'N' = Nautical Miles
	GPSDifferential      tag.GPSDifferential     `tagid:"0x001e" json:"Differential,omitempty" msgpack:"v"`      //  int16u	0 = No Correction 1 = Differential Corrected
}

func (gps GPSInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version         string    `json:"version,omitempty"`
		Latitude        float64   `json:"latitude"`
		Longitude       float64   `json:"longitude"`
		Altitude        float32   `json:"altitude"`
		PositionalError float64   `json:"error,omitempty"`
		DateTime        time.Time `json:"datetime,omitempty"`
		Satellites      string    `json:"satellites,omitempty"`
	}{
		Version:         gps.GPSVersionID.String(),
		Latitude:        gps.Latitude(),
		Longitude:       gps.Longitude(),
		Altitude:        float32(gps.Altitude()),
		PositionalError: math.Floor(gps.GPSHPositioningError.Float()*1000) / 1000,
		DateTime:        gps.DateTime(),
		Satellites:      string(gps.GPSSatellites),
	})
}

func (gps GPSInfo) String() string {
	sb := strings.Builder{}
	sb.WriteString("GPSInfo: \n")
	sb.WriteString(fmt.Sprintf("Coordinates: \t%f,%f\n", gps.Latitude(), gps.Longitude()))
	sb.WriteString(fmt.Sprintf("Altitude: \t%f\n", gps.Altitude()))
	sb.WriteString(fmt.Sprintf("DateTime: \t%s\n", gps.DateTime()))
	sb.WriteString(fmt.Sprintf("Satellites: \t%s\n", gps.GPSSatellites))
	return sb.String()
}

// Latitude returns GPSInfo Latitude as a float64.
func (gps GPSInfo) Latitude() float64 {
	return gps.GPSLatitudeRef.Adjust(gps.GPSLatitude)
}

// Longitude returns GPSInfo Longitude as a float64.
func (gps GPSInfo) Longitude() float64 {
	return gps.GPSLongitudeRef.Adjust(gps.GPSLongitude)
}

// Altitude returns GPSInfo Altitude as a float64.
func (gps GPSInfo) Altitude() float64 {
	res := gps.GPSAltitude.Float()
	if gps.GPSAltitudeRef == tag.GPSAltitudeRefBelow {
		res *= -1
	}
	return math.Floor(res*1000) / 1000
}

// DateTime returns GPSInfo DateStamp and TimeStamp as a time.Time at UTC.
func (gps GPSInfo) DateTime() time.Time {
	if gps.GPSDateStamp.IsNil() {
		return tag.UnknownDate
	}
	hour, min, sec := gps.GPSTimeStamp.HourMinSec()
	return time.Date(int(gps.GPSDateStamp.Year), time.Month(gps.GPSDateStamp.Month), int(gps.GPSDateStamp.Day), int(hour), int(min), int(sec), 0, time.UTC)
}

// DestLatitude returns GPSInfo Destination Latitude as a float64.
func (gps GPSInfo) DestLatitude() float64 {
	return gps.GPSDestLatitudeRef.Adjust(gps.GPSDestLatitude)
}

// DestLongitude returns GPSInfo Destination Longitude as a float64.
func (gps GPSInfo) DestLongitude() float64 {
	return gps.GPSDestLongitudeRef.Adjust(gps.GPSDestLongitude)
}
