package nikon

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// VersionString normalizes a Nikon MakerNote version byte-string.
func VersionString(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	if len(raw) > 4 {
		raw = raw[:4]
	}
	if raw[0] <= 0x09 {
		var buf [12]byte
		out := buf[:0]
		for _, v := range raw {
			out = strconv.AppendInt(out, int64(v), 10)
		}
		return string(out)
	}
	raw = trimNUL(raw)
	if len(raw) == 0 {
		return ""
	}
	return strings.TrimSpace(string(raw))
}

// BitsetIndices expands a byte bitset into a slice of set-bit indices.
func BitsetIndices(raw []byte) []int {
	if len(raw) == 0 {
		return nil
	}

	count := 0
	for _, v := range raw {
		count += bits.OnesCount8(v)
	}
	if count == 0 {
		return nil
	}

	out := make([]int, 0, count)
	for byteIndex, v := range raw {
		for v != 0 {
			bit := bits.TrailingZeros8(v)
			out = append(out, byteIndex*8+bit)
			v &^= 1 << bit
		}
	}
	return out
}

// ISOFromRaw converts a Nikon logarithmic ISO byte to a numeric value.
func ISOFromRaw(raw float64) float64 {
	if raw == 0 {
		return 0
	}
	return 100 * math.Exp((raw/12-5)*math.Ln2)
}

// ByteAt returns the byte at offset or 0 if out of bounds.
func ByteAt(raw []byte, off int) byte {
	if off < 0 || off >= len(raw) {
		return 0
	}
	return raw[off]
}

// U16At returns the uint16 at offset using the given byte order, or 0.
func U16At(raw []byte, off int, bo utils.ByteOrder) uint16 {
	if off < 0 || off+2 > len(raw) {
		return 0
	}
	return bo.Uint16(raw[off : off+2])
}

// RationalPart formats a single rational value as a string.
func RationalPart(v tag.RationalU) string {
	var buf [32]byte
	return string(appendRationalPart(buf[:0], v))
}

func appendRationalPart(dst []byte, v tag.RationalU) []byte {
	if v.Denominator == 0 {
		return append(dst, "undef"...)
	}
	f := v.Float64()
	if f == float64(int64(f)) {
		return strconv.AppendInt(dst, int64(f), 10)
	}
	return strconv.AppendFloat(dst, f, 'f', -1, 64)
}

// LegacyAFPoints decodes AF point indices from a schema-based bitmask.
func LegacyAFPoints(raw []byte, schema uint8, offset int) []int {
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
	return BitsetIndices(raw[offset : offset+size])
}

// AFInfo2HasSelectedMask reports whether the AF area mode uses a selected mask.
func AFInfo2HasSelectedMask(areaMode uint8) bool {
	return areaMode == 8 || areaMode == 9 || areaMode == 13
}

// AFInfo2V0400PointsLen returns the AF point bitmask byte count for v0400+.
func AFInfo2V0400PointsLen(model string) int {
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

// AFInfoByteOrder returns the byte order for AFInfo (BE for D-series, LE for Z).
func AFInfoByteOrder(model string) utils.ByteOrder {
	model = strings.ToUpper(strings.TrimSpace(model))
	if strings.HasPrefix(model, "NIKON D") {
		return utils.BigEndian
	}
	return utils.LittleEndian
}

// WorldTimeByteOrder resolves the byte order for WorldTime heuristically.
func WorldTimeByteOrder(raw []byte, defaultBO utils.ByteOrder) utils.ByteOrder {
	if len(raw) < 4 {
		return defaultBO
	}
	littleTZ := int16(utils.LittleEndian.Uint16(raw[:2]))
	bigTZ := int16(utils.BigEndian.Uint16(raw[:2]))
	littleValid := worldTimeValid(littleTZ, raw[2], raw[3])
	bigValid := worldTimeValid(bigTZ, raw[2], raw[3])
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

func worldTimeValid(tz int16, ds, df byte) bool {
	return tz >= -14*60 && tz <= 14*60 && ds <= 1 && df <= 2
}

// FileInfoByteOrder resolves the byte order for FileInfo heuristically.
func FileInfoByteOrder(raw []byte, model string) utils.ByteOrder {
	if len(raw) < 10 {
		return utils.BigEndian
	}
	littleDir := U16At(raw, 6, utils.LittleEndian)
	littleFile := U16At(raw, 8, utils.LittleEndian)
	bigDir := U16At(raw, 6, utils.BigEndian)
	bigFile := U16At(raw, 8, utils.BigEndian)

	littleValid := fileInfoValid(littleDir, littleFile)
	bigValid := fileInfoValid(bigDir, bigFile)

	switch {
	case littleValid && !bigValid:
		return utils.LittleEndian
	case bigValid && !littleValid:
		return utils.BigEndian
	case fileInfoPrefersLE(model):
		return utils.LittleEndian
	default:
		return utils.BigEndian
	}
}

func fileInfoValid(dir, file uint16) bool {
	return ((dir >= 100 && dir <= 999) || dir == 99) && file <= 9999
}

func fileInfoPrefersLE(model string) bool {
	model = strings.ToUpper(strings.TrimSpace(model))
	switch {
	case strings.Contains(model, "D4S"),
		strings.Contains(model, "D750"),
		strings.Contains(model, "D810"),
		strings.Contains(model, "D3300"),
		strings.Contains(model, "D5200"),
		strings.Contains(model, "D5300"),
		strings.Contains(model, "D5500"),
		strings.Contains(model, "D7100"):
		return true
	default:
		return false
	}
}

// ----- struct decoders (follow Canon.Decode* pattern) -----

// DecodeVRInfo decodes a Nikon VRInfo payload (tag 0x001f).
func DecodeVRInfo(raw []byte) NikonVRInfo {
	if len(raw) == 0 {
		return NikonVRInfo{}
	}
	var dst NikonVRInfo
	if len(raw) >= 4 {
		dst.VRInfoVersion = VersionString(raw[:4])
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

// DecodeWorldTime decodes a Nikon WorldTime payload (tag 0x0024).
func DecodeWorldTime(raw []byte, defaultBO utils.ByteOrder) NikonWorldTime {
	if len(raw) < 4 {
		return NikonWorldTime{}
	}
	bo := WorldTimeByteOrder(raw, defaultBO)
	return NikonWorldTime{
		TimeZone:          int16(bo.Uint16(raw[:2])),
		DaylightSavings:   ByteAt(raw, 2),
		DateDisplayFormat: ByteAt(raw, 3),
	}
}

// DecodeISOInfo decodes a Nikon ISOInfo payload (tag 0x0025).
func DecodeISOInfo(raw []byte) NikonISOInfo {
	if len(raw) < 12 {
		return NikonISOInfo{}
	}
	bo := utils.BigEndian
	dst := NikonISOInfo{
		ISO:           ISOFromRaw(float64(ByteAt(raw, 0))),
		ISOExpansion:  U16At(raw, 4, bo),
		ISO2:          ISOFromRaw(float64(ByteAt(raw, 6))),
		ISOExpansion2: U16At(raw, 10, bo),
	}
	if dst.ISO == 0 && dst.ISOExpansion == 0 && dst.ISO2 == 0 && dst.ISOExpansion2 == 0 {
		return NikonISOInfo{}
	}
	return dst
}

// DecodeLens decodes a Nikon Lens rational payload (tag 0x0084).
func DecodeLens(raw []byte, bo utils.ByteOrder) string {
	if len(raw) < 32 {
		return ""
	}

	var out [128]byte
	buf := out[:0]
	for i := 0; i < 4; i++ {
		if i > 0 {
			buf = append(buf, ' ')
		}
		start := i * 8
		v := tag.RationalU{
			Numerator:   bo.Uint32(raw[start : start+4]),
			Denominator: bo.Uint32(raw[start+4 : start+8]),
		}
		buf = appendRationalPart(buf, v)
	}
	return string(buf)
}

// DecodeLensFStops decodes a Nikon LensFStops payload (tag 0x008b).
func DecodeLensFStops(raw []byte) float64 {
	if len(raw) < 3 || raw[2] == 0 {
		return 0
	}
	return float64(raw[0]) * (float64(raw[1]) / float64(raw[2]))
}

// DecodePowerUpTime decodes a Nikon PowerUpTime payload (tag 0x00b6).
func DecodePowerUpTime(raw []byte, bo utils.ByteOrder) time.Time {
	if len(raw) < 7 {
		return time.Time{}
	}
	year := int(bo.Uint16(raw[:2]))
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

// DecodeAFInfo decodes a Nikon AFInfo payload (tag 0x0088).
func DecodeAFInfo(raw []byte, model string) NikonAFInfo {
	if len(raw) < 4 {
		return NikonAFInfo{}
	}
	bo := AFInfoByteOrder(model)
	mask := bo.Uint16(raw[2:4])
	return NikonAFInfo{
		AFAreaMode:          raw[0],
		AFPoint:             raw[1],
		AFPointsInFocusMask: mask,
		AFPointsInFocus:     BitsetIndices(raw[2:4]),
	}
}

// DecodeAFInfo2 decodes a Nikon AFInfo2 payload (tag 0x00b7).
func DecodeAFInfo2(raw []byte, model string, bo utils.ByteOrder) NikonAFInfo2 {
	if len(raw) < 8 {
		return NikonAFInfo2{}
	}
	version := VersionString(raw[:4])
	dst := NikonAFInfo2{
		AFInfo2Version:    version,
		AFDetectionMethod: raw[4],
		AFAreaMode:        raw[5],
	}

	switch {
	case strings.HasPrefix(version, "04"):
		dst.AFCoordinatesAvailable = raw[7]
		pointsLen := AFInfo2V0400PointsLen(model)
		if pointsLen > 0 && len(raw) >= 10+pointsLen {
			dst.AFPointsUsed = BitsetIndices(raw[10 : 10+pointsLen])
		}
		if len(raw) >= 0x42 {
			dst.AFImageWidth = U16At(raw, 0x3e, bo)
			dst.AFImageHeight = U16At(raw, 0x40, bo)
		}
		if len(raw) >= 0x48 {
			dst.AFAreaXPosition = U16At(raw, 0x42, bo)
			dst.AFAreaYPosition = U16At(raw, 0x44, bo)
			dst.AFAreaWidth = U16At(raw, 0x46, bo)
			dst.AFAreaHeight = U16At(raw, 0x48, bo)
		}
		if len(raw) > 0x4a {
			dst.FocusResult = raw[0x4a]
		}
	default:
		dst.FocusPointSchema = raw[6]
		dst.PrimaryAFPoint = raw[7]
		dst.AFPointsUsed = LegacyAFPoints(raw, dst.FocusPointSchema, 8)
		if dst.AFDetectionMethod == 1 {
			if len(raw) >= 0x1c {
				dst.AFImageWidth = U16At(raw, 0x10, bo)
				dst.AFImageHeight = U16At(raw, 0x12, bo)
				dst.AFAreaXPosition = U16At(raw, 0x14, bo)
				dst.AFAreaYPosition = U16At(raw, 0x16, bo)
				dst.AFAreaWidth = U16At(raw, 0x18, bo)
				dst.AFAreaHeight = U16At(raw, 0x1a, bo)
			}
			if len(raw) > 0x1c {
				dst.ContrastDetectAFInFocus = raw[0x1c] != 0
			}
			return dst
		}
		switch dst.FocusPointSchema {
		case 1:
			if len(raw) >= 0x37 {
				dst.AFPointsInFocus = BitsetIndices(raw[0x30:0x37])
			}
		case 7:
			if len(raw) >= 0x30+20 {
				dst.AFPointsInFocus = BitsetIndices(raw[0x30 : 0x30+20])
			}
			if AFInfo2HasSelectedMask(dst.AFAreaMode) && len(raw) >= 0x1c+20 {
				dst.AFPointsSelected = BitsetIndices(raw[0x1c : 0x1c+20])
			}
			if len(raw) > 0x44 {
				dst.PrimaryAFPoint = raw[0x44]
			}
		}
	}

	return dst
}

// DecodeFileInfo decodes a Nikon FileInfo payload (tag 0x00b8).
func DecodeFileInfo(raw []byte, model string) NikonFileInfo {
	if len(raw) < 10 {
		return NikonFileInfo{}
	}
	bo := FileInfoByteOrder(raw, model)
	return NikonFileInfo{
		FileInfoVersion:  VersionString(raw[:4]),
		MemoryCardNumber: U16At(raw, 4, bo),
		DirectoryNumber:  U16At(raw, 6, bo),
		FileNumber:       U16At(raw, 8, bo),
	}
}

// DecodeAFTune decodes a Nikon AFTune payload (tag 0x00b9).
func DecodeAFTune(raw []byte) NikonAFTune {
	if len(raw) < 4 {
		return NikonAFTune{}
	}
	return NikonAFTune{
		AFFineTune:        raw[0],
		AFFineTuneIndex:   raw[1],
		AFFineTuneAdj:     int8(raw[2]),
		AFFineTuneAdjTele: int8(raw[3]),
	}
}

func trimNUL(buf []byte) []byte {
	for i, b := range buf {
		if b == 0 {
			return buf[:i]
		}
	}
	return buf
}
