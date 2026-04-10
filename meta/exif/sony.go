package exif

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parseSonyTag(t tag.Entry) bool {
	dst := r.sonyMakerNote()
	switch t.ID {
	case sony.Rating:
		dst.Rating = r.parseSonyUint32(t)
	case sony.Contrast:
		dst.Contrast = r.parseSonyInt32(t)
	case sony.Saturation:
		dst.Saturation = r.parseSonyInt32(t)
	case sony.Sharpness:
		dst.Sharpness = r.parseSonyInt32(t)
	case sony.CreativeStyle:
		dst.CreativeStyle = r.parseSonyText(t)
	case sony.DynamicRangeOptimizer:
		dst.DynamicRangeOptimizer = r.parseSonyUint32(t)
	case sony.ImageStabilization:
		dst.ImageStabilization = r.parseSonyUint32(t)
	case sony.ColorMode:
		dst.ColorMode = r.parseSonyUint32(t)
	case sony.Quality:
		dst.Quality = r.parseSonyUint32(t)
	case sony.Quality2:
		dst.Quality2 = r.parseSonyU16Pair(t)
	case sony.WhiteBalance:
		dst.WhiteBalance = r.parseSonyUint32(t)
	case sony.WhiteBalanceFineTune:
		dst.WhiteBalanceFineTune = r.parseSonyInt32(t)
	case sony.FlashExposureComp:
		dst.FlashExposureComp = r.parseSonySignedRationalValue(t)
	case sony.Teleconverter:
		dst.Teleconverter = r.parseSonyUint32(t)
	case sony.SonyModelID:
		dst.SonyModelID = uint16(r.parseSonyUint32(t))
	case sony.LensType:
		dst.LensType = r.parseSonyUint32(t)
	default:
		return false
	}
	return true
}

func (r *Reader) parseSonyText(t tag.Entry) string {
	return strings.TrimSpace(r.parseStringAllowUndefined(t))
}

func (r *Reader) parseSonyUint32(t tag.Entry) uint32 {
	return r.parseMakerNoteUint32(t)
}

func (r *Reader) parseSonyInt32(t tag.Entry) int32 {
	switch t.Type {
	case tag.TypeSignedLong:
		var dst [1]int32
		if n := r.parseInt32List(t, dst[:]); n > 0 {
			return dst[0]
		}
	case tag.TypeSignedShort, tag.TypeShort:
		return int32(r.parseSonyInt16(t))
	case tag.TypeLong:
		return int32(r.parseSonyUint32(t))
	}
	return 0
}

func (r *Reader) parseSonyInt16(t tag.Entry) int16 {
	return r.parseMakerNoteInt16(t)
}

func (r *Reader) parseSonyU16Pair(t tag.Entry) [2]uint16 {
	var dst [2]uint16
	r.parseUint16List(t, dst[:])
	return dst
}

func (r *Reader) parseSonySignedRationalValue(t tag.Entry) float64 {
	var raw [2]int32
	if r.parseRationalSList(t, raw[:]) == 0 || raw[1] == 0 {
		return 0
	}
	return float64(raw[0]) / float64(raw[1])
}
