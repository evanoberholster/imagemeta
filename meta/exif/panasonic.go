package exif

import (
	"strconv"
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parsePanasonicTag(t tag.Entry) bool {
	dst := r.panasonicMakerNote()
	switch uint16(t.ID) {
	case makernote.TagPanasonicImageQuality:
		dst.ImageQuality = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicFirmwareVersion:
		dst.FirmwareVersion = r.parsePanasonicVersionString(t)
	case makernote.TagPanasonicWhiteBalance:
		dst.WhiteBalance = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicFocusMode:
		dst.FocusMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicAFAreaMode:
		dst.AFAreaMode = r.parsePanasonicAFAreaMode(t)
	case makernote.TagPanasonicImageStabilization:
		dst.ImageStabilization = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicMacroMode:
		dst.MacroMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicShootingMode:
		dst.ShootingMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicAudio:
		dst.Audio = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicWhiteBalanceBias:
		dst.WhiteBalanceBias = r.parsePanasonicThirdStops(t)
	case makernote.TagPanasonicFlashBias:
		dst.FlashBias = r.parsePanasonicThirdStops(t)
	case makernote.TagPanasonicExifVersion:
		dst.PanasonicExifVersion = r.parsePanasonicText(t)
	case makernote.TagPanasonicColorEffect:
		dst.ColorEffect = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicTimeSincePowerOn:
		dst.TimeSincePowerOn = float64(r.parsePanasonicUint32(t)) / 100
	case makernote.TagPanasonicBurstMode:
		dst.BurstMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicSequenceNumber:
		dst.SequenceNumber = r.parsePanasonicUint32(t)
	case makernote.TagPanasonicContrastMode:
		dst.ContrastMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicNoiseReduction:
		dst.NoiseReduction = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicSelfTimer:
		dst.SelfTimer = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicRotation:
		dst.Rotation = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicTravelDay:
		dst.TravelDay = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicBatteryLevel:
		dst.BatteryLevel = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicTextStamp, makernote.TagPanasonicTextStamp2, makernote.TagPanasonicTextStamp3:
		dst.TextStamp = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicImageWidth:
		dst.PanasonicImageWidth = r.parsePanasonicUint32(t)
	case makernote.TagPanasonicImageHeight:
		dst.PanasonicImageHeight = r.parsePanasonicUint32(t)
	case makernote.TagPanasonicMakerNoteVersion:
		dst.MakerNoteVersion = r.parsePanasonicText(t)
	case makernote.TagPanasonicSceneMode:
		dst.SceneMode = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicWBRedLevel:
		dst.WBRedLevel = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicWBGreenLevel:
		dst.WBGreenLevel = uint16(r.parsePanasonicUint32(t))
	case makernote.TagPanasonicWBBlueLevel:
		dst.WBBlueLevel = uint16(r.parsePanasonicUint32(t))
	default:
		return false
	}
	return true
}

func (r *Reader) parsePanasonicUint32(t tag.Entry) uint32 {
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

func (r *Reader) parsePanasonicInt16(t tag.Entry) int16 {
	switch t.Type {
	case tag.TypeSignedShort, tag.TypeShort:
	default:
		return 0
	}
	var dst [1]uint16
	if r.parseUint16List(t, dst[:]) == 0 {
		return 0
	}
	return int16(dst[0])
}

func (r *Reader) parsePanasonicThirdStops(t tag.Entry) float64 {
	return float64(r.parsePanasonicInt16(t)) / 3
}

func (r *Reader) parsePanasonicText(t tag.Entry) string {
	raw := trimNULBuffer(r.parseOpaqueBytes(t, min(t.Size(), 128)))
	if len(raw) == 0 {
		return ""
	}
	return strings.TrimSpace(string(raw))
}

func (r *Reader) parsePanasonicVersionString(t tag.Entry) string {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) == 0 {
		return ""
	}
	allPrintable := true
	for _, b := range raw {
		if b >= 0x20 && b <= 0x7e {
			continue
		}
		allPrintable = false
		break
	}
	if allPrintable {
		return strings.TrimSpace(string(trimNULBuffer(raw)))
	}
	var b strings.Builder
	for i, v := range raw {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(int(v)))
	}
	return b.String()
}

func (r *Reader) parsePanasonicAFAreaMode(t tag.Entry) [2]uint8 {
	var raw [2]byte
	r.parseByteList(t, raw[:])
	return [2]uint8{raw[0], raw[1]}
}
