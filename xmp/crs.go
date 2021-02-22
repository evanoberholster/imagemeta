package xmp

import (
	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// Double represents 2 uint16 values
type Double [2]uint16

// CRS is Camera Raw Settings. Photoshop Camera Raw namespace tags.
//	 xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#crs
type CRS struct {
	RawFileName string
}

func (crs *CRS) parse(p property) (err error) {
	switch p.Name() {
	case xmpns.RawFileName:
		crs.RawFileName = parseString(p.Value())
	// Null Operation
	case xmpns.ToneCurve, xmpns.ToneCurveRed, xmpns.ToneCurveGreen, xmpns.ToneCurveBlue:
		return
	case xmpns.ToneCurvePV2012, xmpns.ToneCurvePV2012Red, xmpns.ToneCurvePV2012Green, xmpns.ToneCurvePV2012Blue:
		return
	}
	return ErrPropertyNotSet
}
