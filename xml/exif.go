package xml

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/evanoberholster/image-meta/xml/xmpns"
)

// XMPFlash represents exif:Flash
// Based on: https://exiftool.org/TagNames/XMP.html
type XMPFlash struct {
	Fired      bool
	Mode       uint8
	RedEyeMode bool
	Return     uint8
}

func (xmpFlash *XMPFlash) read(t Tag) (err error) {
	var attr Attribute
	if t.t == SoloTag {
		for t.nextAttr() {
			attr, _ = t.attr()
			switch attr.Name() {
			case xmpns.Fired:
				xmpFlash.Fired = parseBool(attr.value)
			case xmpns.Return:
				xmpFlash.Return = uint8(parseUint32(string(attr.value)))
			case xmpns.Mode:
				xmpFlash.Mode = uint8(parseUint32(string(attr.value)))
			case xmpns.Function:
				//xmpFlash.Function = parseBool(attr.value)
			case xmpns.RedEyeMode:
				xmpFlash.RedEyeMode = parseBool(attr.value)
			default:
			}
			_ = attr
			//fmt.Println(attr)
		}
	} else if t.t == StartTag {

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

func (exif *Exif) decodeElement(decoder *xml.Decoder, start *xml.StartElement) {
	switch start.Name.Local {
	case "ISOSpeedRatings":
		arr := decodeRDF(decoder, start)
		if len(arr) > 0 {
			exif.ISOSpeedRatings = parseUint32(arr[len(arr)-1])
		}
	// TODO: Add support for flash
	//case "Flash":
	default:
		if DebugMode {
			fmt.Println("My Name is: ", start.Name.Local)
		}
	}
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

func (aux *Aux) decodeAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "SerialNumber":
		aux.SerialNumber = attr.Value
	case "Lens":
		aux.Lens = attr.Value
	default:
		err = fmt.Errorf("unknown: %s: %s", attr.Name, attr.Value)
	}
	return err
}
