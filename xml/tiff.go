package xml

import (
	"encoding/xml"
	"fmt"
)

// Tiff attributes of an XMP Packet.
//   xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#tiff
type Tiff struct {
	Make             string
	Model            string
	Software         string
	Copyright        []string
	ImageDescription []string
	ImageWidth       uint16
	ImageLength      uint16
	Orientation      uint8
}

// Orientation represents image orientation
type Orientation uint8

//1 = Horizontal (normal)
//2 = Mirror horizontal
//3 = Rotate 180
//4 = Mirror vertical
//5 = Mirror horizontal and rotate 270 CW
//6 = Rotate 90 CW
//7 = Mirror horizontal and rotate 90 CW
//8 = Rotate 270 CW

// CRS is Camera Raw Settings. Photoshop Camera Raw namespace tags.
//	 xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#crs
type CRS struct {
	RawFileName string
}

func (crs *CRS) decodeAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "RawFileName":
		crs.RawFileName = attr.Value
	default:
		err = fmt.Errorf("unknown: %s: %s", attr.Name, attr.Value)
	}
	return err
}
