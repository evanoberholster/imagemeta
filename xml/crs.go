package xml

import (
	"github.com/evanoberholster/image-meta/xml/xmpns"
)

// Double represents 2 uint16 values
type Double [2]uint16

// CRS is Camera Raw Settings. Photoshop Camera Raw namespace tags.
//	 xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#crs
type CRS struct {
	RawFileName string
}

func (crs *CRS) decode(p property) (err error) {
	switch p.Name() {
	case xmpns.RawFileName:
		crs.RawFileName = parseString(p.val)
	// Null Operation
	case xmpns.ToneCurve, xmpns.ToneCurveRed, xmpns.ToneCurveBlue, xmpns.ToneCurveGreen:
		return
	}
	return ErrPropertyNotSet
}
