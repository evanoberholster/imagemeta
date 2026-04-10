package exif

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

// makerNoteInfo returns the typed maker-note container from Exif.MakerNote.
func (r *Reader) makerNoteInfo() *makernote.Info {
	return &r.Exif.MakerNote
}

func (r *Reader) appleMakerNote() *makernote.Apple {
	info := r.makerNoteInfo()
	if info.Apple == nil {
		info.Apple = &makernote.Apple{}
	}
	return info.Apple
}

func (r *Reader) nikonMakerNote() *nikon.Nikon {
	info := r.makerNoteInfo()
	if info.Nikon == nil {
		info.Nikon = &nikon.Nikon{}
	}
	return info.Nikon
}

func (r *Reader) panasonicMakerNote() *panasonic.Panasonic {
	info := r.makerNoteInfo()
	if info.Panasonic == nil {
		info.Panasonic = &panasonic.Panasonic{}
	}
	return info.Panasonic
}

func (r *Reader) sonyMakerNote() *sony.Sony {
	info := r.makerNoteInfo()
	if info.Sony == nil {
		info.Sony = &sony.Sony{}
	}
	return info.Sony
}

// parseMakerNoteTag parses vendor-specific maker-note tags.
func (r *Reader) parseMakerNoteTag(t tag.Entry) bool {
	switch r.ensureMakerNoteMake() {
	case makernote.CameraMakeCanon:
		return r.parseCanonTag(t)
	case makernote.CameraMakeNikon:
		return r.parseNikonTag(t)
	case makernote.CameraMakePanasonic:
		return r.parsePanasonicTag(t)
	case makernote.CameraMakeSony:
		return r.parseSonyTag(t)
	case makernote.CameraMakeApple:
		return r.parseAppleMakerNoteTag(t)
	default:
		return false
	}
}

// parseAppleMakerNoteTag parses selected Apple maker-note tags.
func (r *Reader) parseAppleMakerNoteTag(t tag.Entry) bool {
	info := r.appleMakerNote()
	switch uint16(t.ID) {
	case makernote.TagAppleMakerNoteVersion:
		info.MakerNoteVersion = int32(r.parseUint32(t))
	case makernote.TagAppleRunTime:
		info.RunTime = r.parseDisplayString(t, 128)
	case makernote.TagAppleAETarget:
		info.AETarget = int32(r.parseUint32(t))
	case makernote.TagAppleAEAverage:
		info.AEAverage = int32(r.parseUint32(t))
	case makernote.TagAppleAFStable:
		info.AFStable = r.parseUint32(t) != 0
	case makernote.TagAppleBurstUUID:
		uuid := r.parseStringAllowUndefined(t)
		info.BurstUUID = meta.UUIDFromString(uuid)
	case makernote.TagAppleAEStable:
		info.AEStable = r.parseUint32(t) != 0
	case makernote.TagAppleOISMode:
		info.OISMode = int32(r.parseUint32(t))
	case makernote.TagAppleContentID:
		info.ContentIdentifier = r.parseStringAllowUndefined(t)
	case makernote.TagAppleImageCaptureType:
		info.ImageCaptureType = int32(r.parseUint32(t))
	case makernote.TagAppleImageUniqueID:
		uuid := r.parseStringAllowUndefined(t)
		info.ImageUniqueID = meta.UUIDFromString(uuid)
	default:
		return false
	}
	return true
}
