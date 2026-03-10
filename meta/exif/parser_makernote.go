package exif

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// parseMakerNoteTag parses vendor-specific maker-note tags.
func (r *Reader) parseMakerNoteTag(t tag.Entry) bool {
	switch r.makerNoteInfo().Make {
	case makernote.CameraMakeCanon:
		return r.parseCanonMakerNoteTag(t)
	case makernote.CameraMakeNikon:
		return r.parseNikonMakerNoteTag(t)
	case makernote.CameraMakeApple:
		return r.parseAppleMakerNoteTag(t)
	default:
		return false
	}
}

// parseCanonMakerNoteTag parses selected Canon maker-note tags.
func (r *Reader) parseCanonMakerNoteTag(t tag.Entry) bool {
	info := r.makerNoteInfo()
	dec := readerCanonValueDecoder{r: r}
	return makernote.ParseCanonTag(&info.Canon, t, dec)
}

// parseNikonMakerNoteTag parses selected Nikon maker-note tags.
func (r *Reader) parseNikonMakerNoteTag(t tag.Entry) bool {
	info := r.makerNoteInfo()
	switch uint16(t.ID) {
	case makernote.TagNikonVersion:
		info.Nikon.VersionCount = uint8(r.parseByteList(t, info.Nikon.Version[:]))
	case makernote.TagNikonISOSetting:
		info.Nikon.ISOSetting = r.parseNikonISOSetting(t)
	case makernote.TagNikonColorMode:
		info.Nikon.ColorMode = r.parseStringAllowUndefined(t)
	case makernote.TagNikonQuality:
		info.Nikon.Quality = r.parseStringAllowUndefined(t)
	case makernote.TagNikonWhiteBalance:
		info.Nikon.WhiteBalance = r.parseStringAllowUndefined(t)
	case makernote.TagNikonSharpness:
		info.Nikon.Sharpness = r.parseStringAllowUndefined(t)
	case makernote.TagNikonFocusMode:
		info.Nikon.FocusMode = r.parseStringAllowUndefined(t)
	case makernote.TagNikonFlashSetting:
		info.Nikon.FlashSetting = r.parseStringAllowUndefined(t)
	case makernote.TagNikonFlashType:
		info.Nikon.FlashType = r.parseStringAllowUndefined(t)
	case makernote.TagNikonISOSelection:
		info.Nikon.ISOSelection = r.parseStringAllowUndefined(t)
	case makernote.TagNikonSerialNumber:
		info.Nikon.SerialNumber = strings.TrimSpace(r.parseStringAllowUndefined(t))
	case makernote.TagNikonLens:
		info.Nikon.Lens = r.parseStringAllowUndefined(t)
	default:
		return false
	}
	return true
}

// parseAppleMakerNoteTag parses selected Apple maker-note tags.
func (r *Reader) parseAppleMakerNoteTag(t tag.Entry) bool {
	info := r.makerNoteInfo()
	switch uint16(t.ID) {
	case makernote.TagAppleMakerNoteVersion:
		info.Apple.MakerNoteVersion = int32(r.parseUint32(t))
	case makernote.TagAppleRunTime:
		info.Apple.RunTime = r.parseDisplayString(t, 128)
	case makernote.TagAppleAETarget:
		info.Apple.AETarget = int32(r.parseUint32(t))
	case makernote.TagAppleAEAverage:
		info.Apple.AEAverage = int32(r.parseUint32(t))
	case makernote.TagAppleAFStable:
		info.Apple.AFStable = r.parseUint32(t) != 0
	case makernote.TagAppleBurstUUID:
		info.Apple.BurstUUID = r.parseStringAllowUndefined(t)
	case makernote.TagAppleAEStable:
		info.Apple.AEStable = r.parseUint32(t) != 0
	case makernote.TagAppleOISMode:
		info.Apple.OISMode = int32(r.parseUint32(t))
	case makernote.TagAppleContentID:
		info.Apple.ContentIdentifier = r.parseStringAllowUndefined(t)
	case makernote.TagAppleImageCaptureType:
		info.Apple.ImageCaptureType = int32(r.parseUint32(t))
	case makernote.TagAppleImageUniqueID:
		info.Apple.ImageUniqueID = r.parseStringAllowUndefined(t)
	default:
		return false
	}
	return true
}

func (r *Reader) parseNikonISOSetting(t tag.Entry) uint32 {
	var iso [2]uint16
	if n := r.parseUint16List(t, iso[:]); n > 0 {
		return uint32(iso[0])
	}
	return r.parseUint32(t)
}

type readerCanonValueDecoder struct {
	r *Reader
}

func (d readerCanonValueDecoder) String(t tag.Entry) string {
	return d.r.parseStringAllowUndefined(t)
}

func (d readerCanonValueDecoder) Uint32(t tag.Entry) uint32 {
	return d.r.parseUint32(t)
}

func (d readerCanonValueDecoder) Uint16List(t tag.Entry, dst []uint16) int {
	if n := d.r.parseUint16List(t, dst); n > 0 {
		return n
	}
	return d.r.parseUndefinedUint16List(t, dst)
}

func (d readerCanonValueDecoder) Int16List(t tag.Entry, dst []int16) int {
	if n := d.r.parseInt16List(t, dst); n > 0 {
		return n
	}
	if !t.IsType(tag.TypeUndefined) || len(dst) == 0 {
		return 0
	}
	var raw [2048]uint16
	n := len(dst)
	if n > len(raw) {
		n = len(raw)
	}
	n = d.r.parseUndefinedUint16List(t, raw[:n])
	for i := 0; i < n; i++ {
		dst[i] = int16(raw[i])
	}
	return n
}
