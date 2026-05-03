package exif

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/apple"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) appleMakerNote() *apple.Apple {
	mknote := &r.Exif.MakerNote
	if mknote.Apple == nil {
		mknote.Apple = &apple.Apple{}
	}
	return mknote.Apple
}

func (r *Reader) dngMakerNote() *makernote.DNG {
	mknote := &r.Exif.MakerNote
	if mknote.DNG == nil {
		mknote.DNG = &makernote.DNG{}
	}
	return mknote.DNG
}

func (r *Reader) nikonMakerNote() *nikon.Nikon {
	mknote := &r.Exif.MakerNote
	if mknote.Nikon == nil {
		mknote.Nikon = &nikon.Nikon{}
	}
	return mknote.Nikon
}

func (r *Reader) panasonicMakerNote() *panasonic.Panasonic {
	mknote := &r.Exif.MakerNote
	if mknote.Panasonic == nil {
		mknote.Panasonic = &panasonic.Panasonic{}
	}
	return mknote.Panasonic
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
		parsed := r.parseSonyTag(t)
		if parsed {
			r.promoteSonyDerivedFields()
		}
		return parsed
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
	case apple.TagAppleMakerNoteVersion:
		info.MakerNoteVersion = r.parseAppleInt32(t)
	case apple.TagAppleRunTime:
		info.RunTime = r.parseAppleRunTime(t)
	case apple.TagAppleAEStable:
		info.AEStable = r.parseAppleInt32(t) != 0
	case apple.TagAppleAETarget:
		info.AETarget = r.parseAppleInt32(t)
	case apple.TagAppleAEAverage:
		info.AEAverage = r.parseAppleInt32(t)
	case apple.TagAppleAFStable:
		info.AFStable = r.parseAppleInt32(t) != 0
	case apple.TagAppleAccelerationVector:
		r.parseAppleFloat64List(t, info.AccelerationVector[:])
	case apple.TagAppleHDRImageType:
		info.HDRImageType = r.parseAppleInt32(t)
	case apple.TagAppleBurstUUID:
		uuid := r.parseStringAllowUndefined(t)
		info.BurstUUID = meta.UUIDFromString(uuid)
	case apple.TagAppleFocusDistanceRange:
		r.parseAppleFloat64List(t, info.FocusDistanceRange[:])
	case apple.TagAppleOISMode:
		info.OISMode = r.parseAppleInt32(t)
	case apple.TagAppleContentID:
		info.ContentIdentifier = r.parseStringAllowUndefined(t)
	case apple.TagAppleImageCaptureType:
		info.ImageCaptureType = r.parseAppleInt32(t)
	case apple.TagAppleImageUniqueID:
		uuid := r.parseStringAllowUndefined(t)
		info.ImageUniqueID = meta.UUIDFromString(uuid)
	case apple.TagAppleLivePhotoVideoIndex:
		info.LivePhotoVideoIndex = uint64(r.parseMakerNoteUint32(t))
	case apple.TagAppleImageProcessingFlags:
		info.ImageProcessingFlags = r.parseMakerNoteUint32(t)
	case apple.TagAppleQualityHint:
		info.QualityHint = r.parseAppleTextOrNumeric(t, 128)
	case apple.TagAppleLuminanceNoiseAmplitude:
		info.LuminanceNoiseAmplitude = r.parseAppleFloat64(t)
	case apple.TagApplePhotosAppFeatureFlags:
		info.PhotosAppFeatureFlags = r.parseMakerNoteUint32(t)
	case apple.TagAppleImageCaptureRequestID:
		info.ImageCaptureRequestID = r.parseAppleTextOrNumeric(t, 128)
	case apple.TagAppleHDRHeadroom:
		info.HDRHeadroom = r.parseAppleFloat64(t)
	case apple.TagAppleAFPerformance:
		r.parseAppleInt32List(t, info.AFPerformance[:])
	case apple.TagAppleSceneFlags:
		info.SceneFlags = r.parseMakerNoteUint32(t)
	case apple.TagAppleSignalToNoiseRatioType:
		info.SignalToNoiseRatioType = r.parseAppleInt32(t)
	case apple.TagAppleSignalToNoiseRatio:
		info.SignalToNoiseRatio = r.parseAppleFloat64(t)
	case apple.TagApplePhotoIdentifier:
		info.PhotoIdentifier = r.parseAppleTextOrNumeric(t, 128)
	case apple.TagAppleColorTemperature:
		info.ColorTemperature = r.parseAppleInt32(t)
	case apple.TagAppleCameraType:
		info.CameraType = r.parseAppleInt32(t)
	case apple.TagAppleFocusPosition:
		info.FocusPosition = r.parseAppleInt32(t)
	case apple.TagAppleHDRGain:
		info.HDRGain = r.parseAppleFloat64(t)
	case apple.TagAppleAFMeasuredDepth:
		info.AFMeasuredDepth = r.parseAppleInt32(t)
	case apple.TagAppleAFConfidence:
		info.AFConfidence = r.parseAppleInt32(t)
	case apple.TagAppleColorCorrectionMatrix:
		info.ColorCorrectionMatrix = r.parseAppleTextOrNumeric(t, 256)
	case apple.TagAppleGreenGhostMitigationStatus:
		info.GreenGhostMitigationStatus = r.parseAppleInt32(t)
	case apple.TagAppleSemanticStyle:
		info.SemanticStyle = r.parseAppleTextOrNumeric(t, 256)
	case apple.TagAppleSemanticStyleRenderingVer:
		info.SemanticStyleRenderingVer = r.parseAppleTextOrNumeric(t, 256)
	case apple.TagAppleSemanticStylePreset:
		info.SemanticStylePreset = r.parseAppleTextOrNumeric(t, 256)
	case apple.TagAppleUnknown004e:
		info.Apple_0x004e = r.parseDisplayString(t, 1024)
	case apple.TagAppleUnknown004f:
		info.Apple_0x004f = r.parseDisplayString(t, 1024)
	case apple.TagAppleUnknown0054:
		info.Apple_0x0054 = r.parseDisplayString(t, 1024)
	case apple.TagAppleUnknown005a:
		info.Apple_0x005a = r.parseDisplayString(t, 1024)
	default:
		return false
	}
	return true
}
