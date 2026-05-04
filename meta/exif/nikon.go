package exif

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parseNikonTag(t tag.Entry) bool {
	dst := r.nikonMakerNote()
	switch nikon.MakerNoteTag(t.ID) {
	case nikon.MakerNoteVersion:
		dst.MakerNoteVersion = r.parseNikonVersion(t)
	case nikon.ISO:
		dst.ISO = r.parseNikonISOValue(t)
	case nikon.ColorMode:
		dst.ColorMode = r.parseNikonText(t)
	case nikon.Quality:
		dst.Quality = r.parseNikonText(t)
	case nikon.WhiteBalance:
		dst.WhiteBalance = r.parseNikonText(t)
	case nikon.Sharpness:
		dst.Sharpness = r.parseNikonText(t)
	case nikon.FocusMode:
		dst.FocusMode = r.parseNikonText(t)
	case nikon.FlashSetting:
		dst.FlashSetting = r.parseNikonText(t)
	case nikon.FlashType:
		dst.FlashType = r.parseNikonText(t)
	case nikon.ISOSelection:
		dst.ISOSelection = r.parseNikonText(t)
	case nikon.ISOSetting:
		dst.ISOSetting = r.parseNikonISOValue(t)
	case nikon.SerialNumber:
		dst.SerialNumber = strings.TrimSpace(r.parseNikonText(t))
	case nikon.SerialNumber2:
		if dst.SerialNumber == "" {
			dst.SerialNumber = strings.TrimSpace(r.parseNikonText(t))
		}
	case nikon.ColorSpace:
		dst.ColorSpace = r.parseNikonUint16(t)
	case nikon.VRInfo:
		dst.VRInfo = r.parseNikonVRInfo(t)
	case nikon.ActiveDLighting:
		dst.ActiveDLighting = r.parseNikonUint16(t)
	case nikon.WorldTime:
		dst.WorldTime = r.parseNikonWorldTime(t)
	case nikon.ISOInfo:
		dst.ISOInfo = r.parseNikonISOInfo(t)
	case nikon.VignetteControl:
		dst.VignetteControl = r.parseNikonUint16(t)
	case nikon.ShutterMode:
		dst.ShutterMode = r.parseNikonUint16(t)
	case nikon.MechanicalShutterCount:
		dst.MechanicalShutterCount = r.parseNikonUint32(t)
	case nikon.ImageSizeRAW:
		dst.ImageSizeRAW = r.parseNikonUint16(t)
	case nikon.ColorTemperatureAuto:
		dst.ColorTemperatureAuto = r.parseNikonUint16(t)
	case nikon.LensType:
		dst.LensType = r.parseNikonUint8(t)
	case nikon.Lens:
		dst.Lens = r.parseNikonLens(t)
	case nikon.ManualFocusDistance:
		dst.ManualFocusDistance = r.parseRationalValue(t)
	case nikon.DigitalZoom:
		dst.DigitalZoom = r.parseRationalValue(t)
	case nikon.FlashMode:
		dst.FlashMode = r.parseNikonUint8(t)
	case nikon.AFInfo:
		dst.AFInfo = r.parseNikonAFInfo(t)
	case nikon.ShootingMode:
		dst.ShootingMode = r.parseNikonUint16(t)
	case nikon.LensFStops:
		dst.LensFStops = r.parseNikonLensFStops(t)
	case nikon.ImageCount:
		dst.ImageCount = r.parseNikonUint32(t)
	case nikon.DeletedImageCount:
		dst.DeletedImageCount = r.parseNikonUint32(t)
	case nikon.ShutterCount:
		dst.ShutterCount = r.parseNikonUint32(t)
	case nikon.PowerUpTime:
		dst.PowerUpTime = r.parseNikonPowerUpTime(t)
	case nikon.AFInfo2:
		dst.AFInfo2 = r.parseNikonAFInfo2(t)
	case nikon.FileInfo:
		dst.FileInfo = r.parseNikonFileInfo(t)
	case nikon.AFTune:
		dst.AFTune = r.parseNikonAFTune(t)
	case nikon.ShotInfo:
		dst.ShotInfo = r.parseNikonShotInfo(t)
	case nikon.SilentPhotography:
		dst.SilentPhotography = r.parseNikonUint32(t) != 0
	default:
		return false
	}
	return true
}

func (r *Reader) parseNikonUint16(t tag.Entry) uint16 {
	value := r.parseNikonUint32(t)
	converted, ok := meta.SafecastUint32ToUint16(value)
	if !ok {
		return 0
	}
	return converted
}

func (r *Reader) parseNikonUint8(t tag.Entry) uint8 {
	value := r.parseNikonUint32(t)
	converted, ok := meta.SafecastUint32ToUint8(value)
	if !ok {
		return 0
	}
	return converted
}

func (r *Reader) parseNikonVersion(t tag.Entry) string {
	return nikon.VersionString(r.parseOpaqueBytes(t, min(t.Size(), 8)))
}

func (r *Reader) parseNikonText(t tag.Entry) string {
	switch t.Type {
	case tag.TypeByte:
		var raw [1]byte
		if r.parseByteList(t, raw[:]) == 0 {
			return ""
		}
		return strconv.FormatUint(uint64(raw[0]), 10)
	case tag.TypeASCII, tag.TypeASCIINoNul, tag.TypeUndefined:
		raw := r.parseOpaqueBytes(t, min(t.Size(), 512))
		if len(raw) == 0 {
			return ""
		}
		raw = trimNULBuffer(raw)
		if len(raw) == 0 {
			return ""
		}
		if i := bytes.IndexByte(raw, 0); i >= 0 {
			raw = raw[:i]
		}
		if len(raw) == 0 {
			return ""
		}
		return strings.TrimSpace(string(raw))
	case tag.TypeShort, tag.TypeLong:
		return strconv.FormatUint(uint64(r.parseNikonUint32(t)), 10)
	default:
		return ""
	}
}

func (r *Reader) parseNikonUint32(t tag.Entry) uint32 {
	return r.parseMakerNoteUint32(t)
}

func (r *Reader) parseNikonISOValue(t tag.Entry) uint32 {
	var iso [2]uint16
	if n := r.parseUint16List(t, iso[:]); n == 2 {
		switch iso[0] {
		case 0, 1:
			return uint32(iso[1])
		default:
			if iso[1] != 0 {
				return uint32(iso[1])
			}
			return uint32(iso[0])
		}
	}
	if n := r.parseUint16List(t, iso[:1]); n == 1 {
		return uint32(iso[0])
	}
	return r.parseNikonUint32(t)
}

func (r *Reader) parseNikonVRInfo(t tag.Entry) nikon.NikonVRInfo {
	return nikon.DecodeVRInfo(r.parseOpaqueBytes(t, min(t.Size(), 16)))
}

func (r *Reader) parseNikonWorldTime(t tag.Entry) nikon.NikonWorldTime {
	return nikon.DecodeWorldTime(r.parseOpaqueBytes(t, min(t.Size(), 8)), t.ByteOrder)
}

func (r *Reader) parseNikonISOInfo(t tag.Entry) nikon.NikonISOInfo {
	return nikon.DecodeISOInfo(r.parseOpaqueBytes(t, min(t.Size(), 16)))
}

func (r *Reader) parseNikonLens(t tag.Entry) string {
	if !t.IsType(tag.TypeRational) {
		return r.parseNikonText(t)
	}
	raw, err := r.readTagBytes(t, min(t.Size(), 32))
	if err != nil || len(raw) < 32 {
		return ""
	}
	return nikon.DecodeLens(raw, t.ByteOrder)
}

func (r *Reader) parseNikonLensFStops(t tag.Entry) float64 {
	return nikon.DecodeLensFStops(r.parseOpaqueBytes(t, min(t.Size(), 4)))
}

func (r *Reader) parseNikonPowerUpTime(t tag.Entry) time.Time {
	return nikon.DecodePowerUpTime(r.parseOpaqueBytes(t, min(t.Size(), 16)), t.ByteOrder)
}

func (r *Reader) parseNikonAFInfo(t tag.Entry) nikon.NikonAFInfo {
	return nikon.DecodeAFInfo(r.parseOpaqueBytes(t, min(t.Size(), 16)), r.Exif.IFD0.Model)
}

func (r *Reader) parseNikonAFInfo2(t tag.Entry) nikon.NikonAFInfo2 {
	return nikon.DecodeAFInfo2(r.parseOpaqueBytes(t, min(t.Size(), 512)), r.Exif.IFD0.Model, t.ByteOrder)
}

func (r *Reader) parseNikonFileInfo(t tag.Entry) nikon.NikonFileInfo {
	return nikon.DecodeFileInfo(r.parseOpaqueBytes(t, min(t.Size(), 16)), r.Exif.IFD0.Model)
}

func (r *Reader) parseNikonAFTune(t tag.Entry) nikon.NikonAFTune {
	return nikon.DecodeAFTune(r.parseOpaqueBytes(t, min(t.Size(), 8)))
}

func (r *Reader) parseNikonShotInfo(t tag.Entry) nikon.NikonShotInfo {
	return nikon.DecodeShotInfo(r.parseOpaqueBytes(t, min(t.Size(), 16384)))
}
