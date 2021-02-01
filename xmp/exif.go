package xmp

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/evanoberholster/image-meta/xmp/xmpns"
)

// XMPFlash represents exif:Flash
// Based on: https://exiftool.org/TagNames/XMP.html
type XMPFlash struct {
	Fired      bool
	Mode       uint8
	RedEyeMode bool
	Function   bool
	Return     uint8
}

func (xmpFlash *XMPFlash) parse(p property) (err error) {
	switch p.Name() {
	case xmpns.Fired:
		xmpFlash.Fired = parseBool(p.val)
	case xmpns.Return:
		xmpFlash.Return = uint8(parseUint(p.val))
	case xmpns.Mode:
		xmpFlash.Mode = uint8(parseUint(p.val))
	case xmpns.Function:
		xmpFlash.Function = parseBool(p.val)
	case xmpns.RedEyeMode:
		xmpFlash.RedEyeMode = parseBool(p.val)
	default:
		return ErrPropertyNotSet
	}
	return
}

// Exif attributes of an XMP Packet.
//	 Exif 2.21 or later: xmlns:exifEX="http://cipa.jp/exif/1.0/"
//	 Exif 2.2 or earlier: xmlns:exif="http://ns.adobe.com/exif/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#exif
type Exif struct {
	ExifVersion       string
	PixelXDimension   uint32
	PixelYDimension   uint32
	DateTimeOriginal  time.Time
	CreateDate        time.Time // Exif:DateTimeDigitized
	ExposureTime      string
	ExposureMode      uint8
	ShutterSpeedValue string
	ExposureProgram   string
	ISOSpeedRatings   uint32
	Flash             XMPFlash
}

func (exif *Exif) decode(tag Tag) (err error) {
	if tag.parent.Name() == xmpns.Flash {
		return exif.Flash.parse(tag.property)
	}
	switch tag.Name() {
	case xmpns.ISOSpeedRatings:
		exif.ISOSpeedRatings = uint32(parseUint(tag.val))
	case xmpns.Flash:
		var attr Attribute
		for tag.nextAttr() {
			attr, _ = tag.attr()
			exif.Flash.parse(attr.property)
		}
	default:
		return ErrPropertyNotSet
	}
	return
}

func (exif *Exif) decodeAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "ExifVersion":
		exif.ExifVersion = attr.Value
	case "PixelXDimension":
		exif.PixelXDimension = parseUint32(attr.Value)
	case "PixelYDimension":
		exif.PixelYDimension = parseUint32(attr.Value)
	case "ExposureMode":
		exif.ExposureMode = uint8(parseUint32(attr.Value))
	default:
		err = fmt.Errorf("unknown: %s: %s", attr.Name, attr.Value)
	}
	return err
}

// Aux attributes of an XMP Packet. These are Adobe-defined auxiliary EXIF tags.
// This implmentation is incomplete and based on https://exiftool.org/TagNames/XMP.html#aux
type Aux struct {
	SerialNumber             string
	LensInfo                 string
	Lens                     string
	LensID                   uint32
	LensSerialNumber         string
	ImageNumber              uint16 // string
	ApproximateFocusDistance string // rational
	FlashCompensation        string // rational
	Firmware                 string
}

func (aux *Aux) decode(p property) (err error) {
	switch p.Name() {
	case xmpns.SerialNumber:
		aux.SerialNumber = parseString(p.val)
	case xmpns.Lens:
		aux.Lens = parseString(p.val)
	case xmpns.LensInfo:
		aux.LensInfo = parseString(p.val)
	case xmpns.LensID:
		aux.LensID = uint32(parseUint(p.val))
	case xmpns.LensSerialNumber:
		aux.LensSerialNumber = parseString(p.val)
	default:
		return ErrPropertyNotSet
	}
	return
}
