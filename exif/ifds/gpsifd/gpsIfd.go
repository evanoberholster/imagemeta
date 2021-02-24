// Package gpsifd provides types for "RootIfd/GPSIfd"
package gpsifd

import "github.com/evanoberholster/imagemeta/exif/tag"

// TagIDMap is a Map of tag.ID to string for the GPSIfd tags
var TagIDMap = map[tag.ID]string{
	GPSAltitude:         "GPSAltitude",
	GPSAltitudeRef:      "GPSAltitudeRef",
	GPSAreaInformation:  "GPSAreaInformation",
	GPSDateStamp:        "GPSDateStamp",
	GPSDestBearing:      "GPSDestBearing",
	GPSDestBearingRef:   "GPSDestBearingRef",
	GPSDestDistance:     "GPSDestDistance",
	GPSDestDistanceRef:  "GPSDestDistanceRef",
	GPSDestLatitude:     "GPSDestLatitude",
	GPSDestLatitudeRef:  "GPSDestLatitudeRef",
	GPSDestLongitude:    "GPSDestLongitude",
	GPSDestLongitudeRef: "GPSDestLongitudeRef",
	GPSDifferential:     "GPSDifferential",
	GPSDOP:              "GPSDOP",
	GPSImgDirection:     "GPSImgDirection",
	GPSImgDirectionRef:  "GPSImgDirectionRef",
	GPSLatitude:         "GPSLatitude",
	GPSLatitudeRef:      "GPSLatitudeRef",
	GPSLongitude:        "GPSLongitude",
	GPSLongitudeRef:     "GPSLongitudeRef",
	GPSMapDatum:         "GPSMapDatum",
	GPSMeasureMode:      "GPSMeasureMode",
	GPSProcessingMethod: "GPSProcessingMethod",
	GPSSatellites:       "GPSSatellites",
	GPSSpeed:            "GPSSpeed",
	GPSSpeedRef:         "GPSSpeedRef",
	GPSStatus:           "GPSStatus",
	GPSTimeStamp:        "GPSTimeStamp",
	GPSTrack:            "GPSTrack",
	GPSTrackRef:         "GPSTrackRef",
	GPSVersionID:        "GPSVersionID",
}

// GPSInfo Tags; GPSInfo Ifd
const (
	GPSVersionID        tag.ID = 0x0000
	GPSLatitudeRef      tag.ID = 0x0001
	GPSLatitude         tag.ID = 0x0002
	GPSLongitudeRef     tag.ID = 0x0003
	GPSLongitude        tag.ID = 0x0004
	GPSAltitudeRef      tag.ID = 0x0005
	GPSAltitude         tag.ID = 0x0006
	GPSTimeStamp        tag.ID = 0x0007
	GPSSatellites       tag.ID = 0x0008
	GPSStatus           tag.ID = 0x0009
	GPSMeasureMode      tag.ID = 0x000a
	GPSDOP              tag.ID = 0x000b
	GPSSpeedRef         tag.ID = 0x000c
	GPSSpeed            tag.ID = 0x000d
	GPSTrackRef         tag.ID = 0x000e
	GPSTrack            tag.ID = 0x000f
	GPSImgDirectionRef  tag.ID = 0x0010
	GPSImgDirection     tag.ID = 0x0011
	GPSMapDatum         tag.ID = 0x0012
	GPSDestLatitudeRef  tag.ID = 0x0013
	GPSDestLatitude     tag.ID = 0x0014
	GPSDestLongitudeRef tag.ID = 0x0015
	GPSDestLongitude    tag.ID = 0x0016
	GPSDestBearingRef   tag.ID = 0x0017
	GPSDestBearing      tag.ID = 0x0018
	GPSDestDistanceRef  tag.ID = 0x0019
	GPSDestDistance     tag.ID = 0x001a
	GPSProcessingMethod tag.ID = 0x001b
	GPSAreaInformation  tag.ID = 0x001c
	GPSDateStamp        tag.ID = 0x001d
	GPSDifferential     tag.ID = 0x001e
)
