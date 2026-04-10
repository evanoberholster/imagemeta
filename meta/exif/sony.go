package exif

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parseSonyTag(t tag.Entry) bool {
	dst := r.sonyMakerNote()
	switch uint16(t.ID) {
	case makernote.TagSonyRating:
		dst.Rating = r.parseSonyUint32(t)
	case makernote.TagSonyContrast:
		dst.Contrast = r.parseSonyInt32(t)
	case makernote.TagSonySaturation:
		dst.Saturation = r.parseSonyInt32(t)
	case makernote.TagSonySharpness:
		dst.Sharpness = r.parseSonyInt32(t)
	case makernote.TagSonyCreativeStyle:
		dst.CreativeStyle = r.parseSonyText(t)
	case makernote.TagSonyDynamicRangeOptimizer:
		dst.DynamicRangeOptimizer = r.parseSonyUint32(t)
	case makernote.TagSonyImageStabilization:
		dst.ImageStabilization = r.parseSonyUint32(t)
	case makernote.TagSonyColorMode:
		dst.ColorMode = r.parseSonyUint32(t)
	case makernote.TagSonyQuality:
		dst.Quality = r.parseSonyUint32(t)
	case makernote.TagSonyQuality2:
		dst.Quality2 = r.parseSonyU16Pair(t)
	case makernote.TagSonyWhiteBalance:
		dst.WhiteBalance = r.parseSonyUint32(t)
	case makernote.TagSonyWhiteBalanceFineTune:
		dst.WhiteBalanceFineTune = r.parseSonyInt32(t)
	case makernote.TagSonyFlashExposureComp:
		dst.FlashExposureComp = r.parseSonySignedRationalValue(t)
	case makernote.TagSonyTeleconverter:
		dst.Teleconverter = r.parseSonyUint32(t)
	case makernote.TagSonyModelID:
		dst.SonyModelID = uint16(r.parseSonyUint32(t))
	case makernote.TagSonyLensType:
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
	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeLong, tag.TypeIfd:
			return t.EmbeddedLong()
		case tag.TypeShort:
			return uint32(t.EmbeddedShort())
		}
	}
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd, tag.TypeShort:
		var dst [2]uint32
		if n := r.parseUint32List(t, dst[:]); n > 0 {
			return dst[0]
		}
	case tag.TypeByte, tag.TypeUndefined, tag.TypeASCII, tag.TypeASCIINoNul:
		var dst [4]byte
		if n := r.parseByteList(t, dst[:]); n > 0 {
			return uint32(dst[0])
		}
	}
	return 0
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
	switch t.Type {
	case tag.TypeSignedShort, tag.TypeShort:
	default:
		return 0
	}
	var raw [1]uint16
	if r.parseUint16List(t, raw[:]) == 0 {
		return 0
	}
	return int16(raw[0])
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
