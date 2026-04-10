package exif

import (
	"strconv"
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote/panasonic"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parsePanasonicTag(t tag.Entry) bool {
	dst := r.panasonicMakerNote()
	switch panasonic.MakerNoteTag(t.ID) {
	case panasonic.ImageQuality:
		dst.ImageQuality = uint16(r.parsePanasonicUint32(t))
	case panasonic.FirmwareVersion:
		dst.FirmwareVersion = r.parsePanasonicVersionString(t)
	case panasonic.WhiteBalance:
		dst.WhiteBalance = uint16(r.parsePanasonicUint32(t))
	case panasonic.FocusMode:
		dst.FocusMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.AFAreaMode:
		dst.AFAreaMode = r.parsePanasonicAFAreaMode(t)
	case panasonic.ImageStabilization:
		dst.ImageStabilization = uint16(r.parsePanasonicUint32(t))
	case panasonic.MacroMode:
		dst.MacroMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.ShootingMode:
		dst.ShootingMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.Audio:
		dst.Audio = uint16(r.parsePanasonicUint32(t))
	case panasonic.WhiteBalanceBias:
		dst.WhiteBalanceBias = r.parsePanasonicThirdStops(t)
	case panasonic.FlashBias:
		dst.FlashBias = r.parsePanasonicThirdStops(t)
	case panasonic.PanasonicExifVersion:
		dst.PanasonicExifVersion = r.parsePanasonicText(t)
	case panasonic.ColorEffect:
		dst.ColorEffect = uint16(r.parsePanasonicUint32(t))
	case panasonic.TimeSincePowerOn:
		dst.TimeSincePowerOn = float64(r.parsePanasonicUint32(t)) / 100
	case panasonic.BurstMode:
		dst.BurstMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.SequenceNumber:
		dst.SequenceNumber = r.parsePanasonicUint32(t)
	case panasonic.ContrastMode:
		dst.ContrastMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.NoiseReduction:
		dst.NoiseReduction = uint16(r.parsePanasonicUint32(t))
	case panasonic.SelfTimer:
		dst.SelfTimer = uint16(r.parsePanasonicUint32(t))
	case panasonic.Rotation:
		dst.Rotation = uint16(r.parsePanasonicUint32(t))
	case panasonic.TravelDay:
		dst.TravelDay = uint16(r.parsePanasonicUint32(t))
	case panasonic.BatteryLevel:
		dst.BatteryLevel = uint16(r.parsePanasonicUint32(t))
	case panasonic.TextStamp, panasonic.TextStamp2, panasonic.TextStamp3:
		dst.TextStamp = uint16(r.parsePanasonicUint32(t))
	case panasonic.PanasonicImageWidth:
		dst.PanasonicImageWidth = r.parsePanasonicUint32(t)
	case panasonic.PanasonicImageHeight:
		dst.PanasonicImageHeight = r.parsePanasonicUint32(t)
	case panasonic.MakerNoteVersion:
		dst.MakerNoteVersion = r.parsePanasonicText(t)
	case panasonic.SceneMode:
		dst.SceneMode = uint16(r.parsePanasonicUint32(t))
	case panasonic.WBRedLevel:
		dst.WBRedLevel = uint16(r.parsePanasonicUint32(t))
	case panasonic.WBGreenLevel:
		dst.WBGreenLevel = uint16(r.parsePanasonicUint32(t))
	case panasonic.WBBlueLevel:
		dst.WBBlueLevel = uint16(r.parsePanasonicUint32(t))
	default:
		return false
	}
	return true
}

func (r *Reader) parsePanasonicUint32(t tag.Entry) uint32 {
	return r.parseMakerNoteUint32(t)
}

func (r *Reader) parsePanasonicInt16(t tag.Entry) int16 {
	return r.parseMakerNoteInt16(t)
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
