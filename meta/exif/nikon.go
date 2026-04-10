package exif

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	metanikon "github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func (r *Reader) parseNikonTag(t tag.Entry) bool {
	dst := r.nikonMakerNote()
	switch uint16(t.ID) {
	case makernote.TagNikonMakerNoteVersion:
		dst.MakerNoteVersion = r.parseNikonVersion(t)
	case makernote.TagNikonISO:
		dst.ISO = r.parseNikonISOValue(t)
	case makernote.TagNikonColorMode:
		dst.ColorMode = r.parseNikonText(t)
	case makernote.TagNikonQuality:
		dst.Quality = r.parseNikonText(t)
	case makernote.TagNikonWhiteBalance:
		dst.WhiteBalance = r.parseNikonText(t)
	case makernote.TagNikonSharpness:
		dst.Sharpness = r.parseNikonText(t)
	case makernote.TagNikonFocusMode:
		dst.FocusMode = r.parseNikonText(t)
	case makernote.TagNikonFlashSetting:
		dst.FlashSetting = r.parseNikonText(t)
	case makernote.TagNikonFlashType:
		dst.FlashType = r.parseNikonText(t)
	case makernote.TagNikonISOSelection:
		dst.ISOSelection = r.parseNikonText(t)
	case makernote.TagNikonISOSetting:
		dst.ISOSetting = r.parseNikonISOValue(t)
	case makernote.TagNikonSerialNumber:
		dst.SerialNumber = strings.TrimSpace(r.parseNikonText(t))
	case makernote.TagNikonSerialNumber2:
		if dst.SerialNumber == "" {
			dst.SerialNumber = strings.TrimSpace(r.parseNikonText(t))
		}
	case makernote.TagNikonColorSpace:
		dst.ColorSpace = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonVRInfo:
		dst.VRInfo = r.parseNikonVRInfo(t)
	case makernote.TagNikonActiveDLighting:
		dst.ActiveDLighting = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonWorldTime:
		dst.WorldTime = r.parseNikonWorldTime(t)
	case makernote.TagNikonISOInfo:
		dst.ISOInfo = r.parseNikonISOInfo(t)
	case makernote.TagNikonVignetteControl:
		dst.VignetteControl = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonShutterMode:
		dst.ShutterMode = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonMechanicalShutterCount:
		dst.MechanicalShutterCount = r.parseNikonUint32(t)
	case makernote.TagNikonImageSizeRAW:
		dst.ImageSizeRAW = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonColorTemperatureAuto:
		dst.ColorTemperatureAuto = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonLensType:
		dst.LensType = uint8(r.parseNikonUint32(t))
	case makernote.TagNikonLens:
		dst.Lens = r.parseNikonLens(t)
	case makernote.TagNikonManualFocusDistance:
		dst.ManualFocusDistance = r.parseRationalValue(t)
	case makernote.TagNikonDigitalZoom:
		dst.DigitalZoom = r.parseRationalValue(t)
	case makernote.TagNikonFlashMode:
		dst.FlashMode = uint8(r.parseNikonUint32(t))
	case makernote.TagNikonAFInfo:
		dst.AFInfo = r.parseNikonAFInfo(t)
	case makernote.TagNikonShootingMode:
		dst.ShootingMode = uint16(r.parseNikonUint32(t))
	case makernote.TagNikonLensFStops:
		dst.LensFStops = r.parseNikonLensFStops(t)
	case makernote.TagNikonImageCount:
		dst.ImageCount = r.parseNikonUint32(t)
	case makernote.TagNikonDeletedImageCount:
		dst.DeletedImageCount = r.parseNikonUint32(t)
	case makernote.TagNikonShutterCount:
		dst.ShutterCount = r.parseNikonUint32(t)
	case makernote.TagNikonPowerUpTime:
		dst.PowerUpTime = r.parseNikonPowerUpTime(t)
	case makernote.TagNikonAFInfo2:
		dst.AFInfo2 = r.parseNikonAFInfo2(t)
	case makernote.TagNikonFileInfo:
		dst.FileInfo = r.parseNikonFileInfo(t)
	case makernote.TagNikonAFTune:
		dst.AFTune = r.parseNikonAFTune(t)
	case makernote.TagNikonSilentPhotography:
		dst.SilentPhotography = r.parseNikonUint32(t) != 0
	default:
		return false
	}
	return true
}

func (r *Reader) parseNikonVersion(t tag.Entry) string {
	return nikonVersionString(r.parseOpaqueBytes(t, min(t.Size(), 8)))
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
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd:
		return r.parseUint32(t)
	case tag.TypeShort, tag.TypeSignedShort:
		var v [1]uint16
		if r.parseUint16List(t, v[:]) > 0 {
			return uint32(v[0])
		}
	case tag.TypeByte, tag.TypeUndefined, tag.TypeASCII, tag.TypeASCIINoNul:
		var v [1]byte
		if r.parseByteList(t, v[:]) > 0 {
			return uint32(v[0])
		}
	}
	return 0
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

func (r *Reader) parseNikonVRInfo(t tag.Entry) metanikon.NikonVRInfo {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) == 0 {
		return metanikon.NikonVRInfo{}
	}
	var dst metanikon.NikonVRInfo
	if len(raw) >= 4 {
		dst.VRInfoVersion = nikonVersionString(raw[:4])
	}
	if len(raw) > 4 {
		dst.VibrationReduction = raw[4]
	}
	if len(raw) > 6 {
		dst.VRMode = raw[6]
	}
	if len(raw) > 8 {
		dst.VRType = raw[8]
	}
	return dst
}

func (r *Reader) parseNikonWorldTime(t tag.Entry) metanikon.NikonWorldTime {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 8))
	if len(raw) < 4 {
		return metanikon.NikonWorldTime{}
	}
	bo := nikonWorldTimeByteOrder(raw, t.ByteOrder)
	return metanikon.NikonWorldTime{
		TimeZone:          int16(bo.Uint16(raw[:2])),
		DaylightSavings:   nikonByteAt(raw, 2),
		DateDisplayFormat: nikonByteAt(raw, 3),
	}
}

func (r *Reader) parseNikonISOInfo(t tag.Entry) metanikon.NikonISOInfo {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) < 12 {
		return metanikon.NikonISOInfo{}
	}
	bo := utils.BigEndian
	dst := metanikon.NikonISOInfo{
		// ExifTool's Nikon::ISOInfo table defines ISO/ISO2 at offsets 0 and 6
		// as the raw logarithmic Nikon ISO byte, not a uint16.
		ISO:           nikonISOFromRaw(float64(nikonByteAt(raw, 0))),
		ISOExpansion:  nikonU16At(raw, 4, bo),
		ISO2:          nikonISOFromRaw(float64(nikonByteAt(raw, 6))),
		ISOExpansion2: nikonU16At(raw, 10, bo),
	}
	if dst.ISO == 0 && dst.ISOExpansion == 0 && dst.ISO2 == 0 && dst.ISOExpansion2 == 0 {
		return metanikon.NikonISOInfo{}
	}
	return dst
}

func (r *Reader) parseNikonLens(t tag.Entry) string {
	if !t.IsType(tag.TypeRational) {
		return r.parseNikonText(t)
	}
	raw, _, err := r.readTagBytes(t, min(t.Size(), 32))
	if err != nil || len(raw) < 32 {
		return ""
	}
	var parts [4]tag.RationalU
	for i := range parts {
		start := i * 8
		parts[i] = tag.RationalU{
			Numerator:   t.ByteOrder.Uint32(raw[start : start+4]),
			Denominator: t.ByteOrder.Uint32(raw[start+4 : start+8]),
		}
	}
	return strings.Join([]string{
		nikonRationalPart(parts[0]),
		nikonRationalPart(parts[1]),
		nikonRationalPart(parts[2]),
		nikonRationalPart(parts[3]),
	}, " ")
}

func (r *Reader) parseNikonLensFStops(t tag.Entry) float64 {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 4))
	if len(raw) < 3 || raw[2] == 0 {
		return 0
	}
	return float64(raw[0]) * (float64(raw[1]) / float64(raw[2]))
}

func (r *Reader) parseNikonPowerUpTime(t tag.Entry) time.Time {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) < 7 {
		return time.Time{}
	}
	year := int(t.ByteOrder.Uint16(raw[:2]))
	month := time.Month(raw[2])
	day := int(raw[3])
	hour := int(raw[4])
	minute := int(raw[5])
	second := int(raw[6])
	if year == 0 || month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}
	}
	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}

func (r *Reader) parseNikonAFInfo(t tag.Entry) metanikon.NikonAFInfo {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) < 4 {
		return metanikon.NikonAFInfo{}
	}
	bo := nikonAFInfoByteOrder(r.Exif.IFD0.Model)
	mask := bo.Uint16(raw[2:4])
	return metanikon.NikonAFInfo{
		AFAreaMode:          raw[0],
		AFPoint:             raw[1],
		AFPointsInFocusMask: mask,
		AFPointsInFocus:     nikonBitsetIndices(raw[2:4]),
	}
}

func (r *Reader) parseNikonAFInfo2(t tag.Entry) metanikon.NikonAFInfo2 {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 512))
	if len(raw) < 8 {
		return metanikon.NikonAFInfo2{}
	}
	version := nikonVersionString(raw[:4])
	dst := metanikon.NikonAFInfo2{
		AFInfo2Version:    version,
		AFDetectionMethod: raw[4],
		AFAreaMode:        raw[5],
	}

	switch {
	case strings.HasPrefix(version, "04"):
		dst.AFCoordinatesAvailable = raw[7]
		pointsLen := nikonAFInfo2V0400PointsLen(r.Exif.IFD0.Model)
		if pointsLen > 0 && len(raw) >= 10+pointsLen {
			dst.AFPointsUsed = nikonBitsetIndices(raw[10 : 10+pointsLen])
		}
		if len(raw) >= 0x42 {
			dst.AFImageWidth = nikonU16At(raw, 0x3e, t.ByteOrder)
			dst.AFImageHeight = nikonU16At(raw, 0x40, t.ByteOrder)
		}
		if len(raw) >= 0x48 {
			dst.AFAreaXPosition = nikonU16At(raw, 0x42, t.ByteOrder)
			dst.AFAreaYPosition = nikonU16At(raw, 0x44, t.ByteOrder)
			dst.AFAreaWidth = nikonU16At(raw, 0x46, t.ByteOrder)
			dst.AFAreaHeight = nikonU16At(raw, 0x48, t.ByteOrder)
		}
		if len(raw) > 0x4a {
			dst.FocusResult = raw[0x4a]
		}
	default:
		dst.FocusPointSchema = raw[6]
		dst.PrimaryAFPoint = raw[7]
		dst.AFPointsUsed = nikonLegacyAFPoints(raw, dst.FocusPointSchema, 8)
		if dst.AFDetectionMethod == 1 {
			if len(raw) >= 0x1c {
				dst.AFImageWidth = nikonU16At(raw, 0x10, t.ByteOrder)
				dst.AFImageHeight = nikonU16At(raw, 0x12, t.ByteOrder)
				dst.AFAreaXPosition = nikonU16At(raw, 0x14, t.ByteOrder)
				dst.AFAreaYPosition = nikonU16At(raw, 0x16, t.ByteOrder)
				dst.AFAreaWidth = nikonU16At(raw, 0x18, t.ByteOrder)
				dst.AFAreaHeight = nikonU16At(raw, 0x1a, t.ByteOrder)
			}
			if len(raw) > 0x1c {
				dst.ContrastDetectAFInFocus = raw[0x1c] != 0
			}
			return dst
		}
		switch dst.FocusPointSchema {
		case 1:
			if len(raw) >= 0x37 {
				dst.AFPointsInFocus = nikonBitsetIndices(raw[0x30:0x37])
			}
		case 7:
			if len(raw) >= 0x30+20 {
				dst.AFPointsInFocus = nikonBitsetIndices(raw[0x30 : 0x30+20])
			}
			if nikonAFInfo2HasSelectedMask(dst.AFAreaMode) && len(raw) >= 0x1c+20 {
				dst.AFPointsSelected = nikonBitsetIndices(raw[0x1c : 0x1c+20])
			}
			if len(raw) > 0x44 {
				dst.PrimaryAFPoint = raw[0x44]
			}
		}
	}

	return dst
}

func (r *Reader) parseNikonFileInfo(t tag.Entry) metanikon.NikonFileInfo {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 16))
	if len(raw) < 10 {
		return metanikon.NikonFileInfo{}
	}
	bo := nikonFileInfoByteOrder(raw, r.Exif.IFD0.Model)
	return metanikon.NikonFileInfo{
		FileInfoVersion:  nikonVersionString(raw[:4]),
		MemoryCardNumber: nikonU16At(raw, 4, bo),
		DirectoryNumber:  nikonU16At(raw, 6, bo),
		FileNumber:       nikonU16At(raw, 8, bo),
	}
}

func (r *Reader) parseNikonAFTune(t tag.Entry) metanikon.NikonAFTune {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 8))
	if len(raw) < 4 {
		return metanikon.NikonAFTune{}
	}
	return metanikon.NikonAFTune{
		AFFineTune:        raw[0],
		AFFineTuneIndex:   raw[1],
		AFFineTuneAdj:     int8(raw[2]),
		AFFineTuneAdjTele: int8(raw[3]),
	}
}

func nikonVersionString(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	if len(raw) > 4 {
		raw = raw[:4]
	}
	if raw[0] <= 0x09 {
		var b strings.Builder
		b.Grow(8)
		for i := range raw {
			b.WriteString(strconv.Itoa(int(raw[i])))
		}
		return b.String()
	}
	raw = trimNULBuffer(raw)
	if len(raw) == 0 {
		return ""
	}
	return strings.TrimSpace(string(raw))
}

func nikonWorldTimeByteOrder(raw []byte, defaultBO utils.ByteOrder) utils.ByteOrder {
	if len(raw) < 4 {
		return defaultBO
	}
	littleTZ := int16(utils.LittleEndian.Uint16(raw[:2]))
	littleValid := nikonValidWorldTime(littleTZ, raw[2], raw[3])
	bigTZ := int16(utils.BigEndian.Uint16(raw[:2]))
	bigValid := nikonValidWorldTime(bigTZ, raw[2], raw[3])
	switch {
	case littleValid && !bigValid:
		return utils.LittleEndian
	case bigValid && !littleValid:
		return utils.BigEndian
	case defaultBO != utils.UnknownEndian:
		return defaultBO
	default:
		return utils.LittleEndian
	}
}

func nikonValidWorldTime(timeZone int16, daylightSavings, dateDisplayFormat byte) bool {
	return timeZone >= -14*60 && timeZone <= 14*60 && daylightSavings <= 1 && dateDisplayFormat <= 2
}

func nikonISOFromRaw(raw float64) float64 {
	if raw == 0 {
		return 0
	}
	return 100 * math.Exp((raw/12-5)*math.Ln2)
}

func nikonRationalPart(v tag.RationalU) string {
	if v.Denominator == 0 {
		return "undef"
	}
	f := v.Float64()
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func nikonAFInfoByteOrder(model string) utils.ByteOrder {
	model = strings.ToUpper(strings.TrimSpace(model))
	if strings.HasPrefix(model, "NIKON D") {
		return utils.BigEndian
	}
	return utils.LittleEndian
}

func nikonLegacyAFPoints(raw []byte, schema uint8, offset int) []int {
	size := 0
	switch schema {
	case 1:
		size = 7
	case 2:
		size = 2
	case 3:
		size = 5
	case 7:
		size = 20
	default:
		return nil
	}
	if len(raw) < offset+size {
		return nil
	}
	return nikonBitsetIndices(raw[offset : offset+size])
}

func nikonAFInfo2HasSelectedMask(areaMode uint8) bool {
	switch areaMode {
	case 8, 9, 13:
		return true
	default:
		return false
	}
}

func nikonAFInfo2V0400PointsLen(model string) int {
	model = strings.ToUpper(strings.TrimSpace(model))
	switch {
	case strings.Contains(model, "Z 8"), strings.Contains(model, "Z 9"):
		return 51
	case strings.Contains(model, "Z50_2"), strings.Contains(model, "Z 50_2"), strings.Contains(model, "Z50II"):
		return 29
	default:
		return 38
	}
}

func nikonBitsetIndices(raw []byte) []int {
	if len(raw) == 0 {
		return nil
	}
	out := make([]int, 0, len(raw))
	for byteIndex, v := range raw {
		if v == 0 {
			continue
		}
		for bit := 0; bit < 8; bit++ {
			if v&(1<<bit) == 0 {
				continue
			}
			out = append(out, byteIndex*8+bit)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func nikonFileInfoByteOrder(raw []byte, model string) utils.ByteOrder {
	if len(raw) < 10 {
		return utils.BigEndian
	}
	littleDir := nikonU16At(raw, 6, utils.LittleEndian)
	littleFile := nikonU16At(raw, 8, utils.LittleEndian)
	littleValid := nikonValidFileInfo(littleDir, littleFile)

	bigDir := nikonU16At(raw, 6, utils.BigEndian)
	bigFile := nikonU16At(raw, 8, utils.BigEndian)
	bigValid := nikonValidFileInfo(bigDir, bigFile)

	switch {
	case littleValid && !bigValid:
		return utils.LittleEndian
	case bigValid && !littleValid:
		return utils.BigEndian
	case nikonFileInfoPrefersLittleEndian(model):
		return utils.LittleEndian
	default:
		return utils.BigEndian
	}
}

func nikonValidFileInfo(directoryNumber, fileNumber uint16) bool {
	return ((directoryNumber >= 100 && directoryNumber <= 999) || directoryNumber == 99) && fileNumber <= 9999
}

func nikonFileInfoPrefersLittleEndian(model string) bool {
	model = strings.ToUpper(strings.TrimSpace(model))
	for _, token := range []string{"D4S", "D750", "D810", "D3300", "D5200", "D5300", "D5500", "D7100"} {
		if strings.Contains(model, token) {
			return true
		}
	}
	return false
}

func nikonByteAt(raw []byte, offset int) byte {
	if offset < 0 || offset >= len(raw) {
		return 0
	}
	return raw[offset]
}

func nikonU16At(raw []byte, offset int, byteOrder utils.ByteOrder) uint16 {
	if offset < 0 || offset+2 > len(raw) {
		return 0
	}
	return byteOrder.Uint16(raw[offset : offset+2])
}
