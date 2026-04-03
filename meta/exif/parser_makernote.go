package exif

import (
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// parseMakerNoteTag parses vendor-specific maker-note tags.
func (r *Reader) parseMakerNoteTag(t tag.Entry) bool {
	switch r.ensureMakerNoteMake() {
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
	return r.parseCanonTag(t)
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
