package apple

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"time"
	"unicode/utf16"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

const (
	// MakerNoteHeaderLength is the fixed Apple prefix length before the IFD tag table.
	MakerNoteHeaderLength = 26

	TagAppleMakerNoteVersion           uint16 = 0x0001
	TagAppleRunTime                    uint16 = 0x0003
	TagAppleAEStable                   uint16 = 0x0004
	TagAppleAETarget                   uint16 = 0x0005
	TagAppleAEAverage                  uint16 = 0x0006
	TagAppleAFStable                   uint16 = 0x0007
	TagAppleAccelerationVector         uint16 = 0x0008
	TagAppleHDRImageType               uint16 = 0x000a
	TagAppleBurstUUID                  uint16 = 0x000b
	TagAppleFocusDistanceRange         uint16 = 0x000c
	TagAppleOISMode                    uint16 = 0x000f
	TagAppleContentID                  uint16 = 0x0011
	TagAppleImageCaptureType           uint16 = 0x0014
	TagAppleImageUniqueID              uint16 = 0x0015
	TagAppleLivePhotoVideoIndex        uint16 = 0x0017
	TagAppleImageProcessingFlags       uint16 = 0x0019
	TagAppleQualityHint                uint16 = 0x001a
	TagAppleLuminanceNoiseAmplitude    uint16 = 0x001d
	TagApplePhotosAppFeatureFlags      uint16 = 0x001f
	TagAppleImageCaptureRequestID      uint16 = 0x0020
	TagAppleHDRHeadroom                uint16 = 0x0021
	TagAppleAFPerformance              uint16 = 0x0023
	TagAppleSceneFlags                 uint16 = 0x0025
	TagAppleSignalToNoiseRatioType     uint16 = 0x0026
	TagAppleSignalToNoiseRatio         uint16 = 0x0027
	TagApplePhotoIdentifier            uint16 = 0x002b
	TagAppleColorTemperature           uint16 = 0x002d
	TagAppleCameraType                 uint16 = 0x002e
	TagAppleFocusPosition              uint16 = 0x002f
	TagAppleHDRGain                    uint16 = 0x0030
	TagAppleAFMeasuredDepth            uint16 = 0x0038
	TagAppleAFConfidence               uint16 = 0x003d
	TagAppleColorCorrectionMatrix      uint16 = 0x003e
	TagAppleGreenGhostMitigationStatus uint16 = 0x003f
	TagAppleSemanticStyle              uint16 = 0x0040
	TagAppleSemanticStyleRenderingVer  uint16 = 0x0041
	TagAppleSemanticStylePreset        uint16 = 0x0042
	TagAppleUnknown004e                uint16 = 0x004e
	TagAppleUnknown004f                uint16 = 0x004f
	TagAppleUnknown0054                uint16 = 0x0054
	TagAppleUnknown005a                uint16 = 0x005a
)

// AppleRunTime contains the Apple PLIST-format CMTime structure from tag 0x0003.
type AppleRunTime struct {
	Flags uint32
	Value uint64
	Scale uint64
	Epoch uint64
}

// Duration returns the runtime as a time.Duration when the scale is non-zero.
func (r AppleRunTime) Duration() (time.Duration, bool) {
	if r.Scale == 0 {
		return 0, false
	}
	ns := (r.Value * uint64(time.Second)) / r.Scale
	nsI64, ok := meta.SafecastUint64ToInt64(ns)
	if !ok {
		return 0, false
	}
	return time.Duration(nsI64), true
}

// String renders the runtime in the same numeric form ExifTool reports via its
// subdirectory fields.
func (r AppleRunTime) String() string {
	if d, ok := r.Duration(); ok {
		return strconv.FormatUint(r.Value, 10) + "/" + strconv.FormatUint(r.Scale, 10) + " (" + d.String() + ")"
	}
	return strconv.FormatUint(r.Value, 10) + "/" + strconv.FormatUint(r.Scale, 10)
}

// Apple contains selected Apple maker-note fields.
type Apple struct {
	MakerNoteVersion           int32
	RunTime                    AppleRunTime
	AEStable                   bool
	AETarget                   int32
	AEAverage                  int32
	AFStable                   bool
	AccelerationVector         [3]float64
	HDRImageType               int32
	BurstUUID                  meta.UUID
	FocusDistanceRange         [2]float64
	OISMode                    int32
	ContentIdentifier          string
	ImageCaptureType           int32
	ImageUniqueID              meta.UUID
	LivePhotoVideoIndex        uint64
	ImageProcessingFlags       uint32
	QualityHint                string
	LuminanceNoiseAmplitude    float64
	PhotosAppFeatureFlags      uint32
	ImageCaptureRequestID      string
	HDRHeadroom                float64
	AFPerformance              [2]int32
	SceneFlags                 uint32
	SignalToNoiseRatioType     int32
	SignalToNoiseRatio         float64
	PhotoIdentifier            string
	ColorTemperature           int32
	CameraType                 int32
	FocusPosition              int32
	HDRGain                    float64
	AFMeasuredDepth            int32
	AFConfidence               int32
	ColorCorrectionMatrix      string
	GreenGhostMitigationStatus int32
	SemanticStyle              string
	SemanticStyleRenderingVer  string
	SemanticStylePreset        string
	Apple_0x004e               string
	Apple_0x004f               string
	Apple_0x0054               string
	Apple_0x005a               string
}

// ParseRunTime decodes the Apple RunTime binary plist into a typed CMTime struct.
func ParseRunTime(raw []byte) (AppleRunTime, bool) {
	if len(raw) < 32 {
		return AppleRunTime{}, false
	}
	if i := bytes.Index(raw, []byte("timescale")); i >= 0 {
		return parseRunTimeKeyed(raw, i)
	}
	trailer := raw[len(raw)-32:]
	offsetIntSize := int(trailer[6])
	objRefSize := int(trailer[7])
	numObjects, convOK := meta.SafecastUint64ToInt(binary.BigEndian.Uint64(trailer[8:16]))
	if !convOK {
		return AppleRunTime{}, false
	}
	topObject, convOK := meta.SafecastUint64ToInt(binary.BigEndian.Uint64(trailer[16:24]))
	if !convOK {
		return AppleRunTime{}, false
	}
	offsetTableOffset, convOK := meta.SafecastUint64ToInt(binary.BigEndian.Uint64(trailer[24:32]))
	if !convOK {
		return AppleRunTime{}, false
	}
	if offsetIntSize <= 0 || objRefSize <= 0 || numObjects <= 0 || topObject < 0 || topObject >= numObjects {
		return AppleRunTime{}, false
	}
	if offsetTableOffset < 0 || offsetTableOffset+numObjects*offsetIntSize > len(raw)-32 {
		return AppleRunTime{}, false
	}

	offsets := make([]int, numObjects)
	for i := 0; i < numObjects; i++ {
		off := offsetTableOffset + i*offsetIntSize
		v, ok := readUnsigned(raw[off:], offsetIntSize)
		if !ok || v > uint64(len(raw)) {
			return AppleRunTime{}, false
		}
		offset, convOK := meta.SafecastUint64ToInt(v)
		if !convOK {
			return AppleRunTime{}, false
		}
		offsets[i] = offset
	}

	cache := make(map[int]any, numObjects)
	visiting := make(map[int]bool)
	var parseObject func(int) (any, bool)
	parseObject = func(idx int) (any, bool) {
		if v, ok := cache[idx]; ok {
			return v, true
		}
		if idx < 0 || idx >= numObjects {
			return nil, false
		}
		// Reject cyclic object references (e.g. an array or dict that refers
		// back to itself), which would otherwise recurse until the goroutine
		// stack is exhausted.
		if visiting[idx] {
			return nil, false
		}
		visiting[idx] = true
		defer delete(visiting, idx)
		off := offsets[idx]
		if off < 0 || off >= len(raw)-32 {
			return nil, false
		}
		marker := raw[off]
		kind := marker >> 4
		low := marker & 0x0f
		switch kind {
		case 0x0:
			switch low {
			case 0x8:
				return false, true
			case 0x9:
				return true, true
			default:
				return nil, true
			}
		case 0x1:
			size := 1 << low
			if size <= 0 || size > 8 || off+1+size > len(raw) {
				return nil, false
			}
			v, ok := readUnsigned(raw[off+1:], size)
			if !ok {
				return nil, false
			}
			cache[idx] = v
			return v, true
		case 0x2:
			size := 1 << low
			if size == 4 {
				if off+5 > len(raw) {
					return nil, false
				}
				v := math.Float32frombits(binary.BigEndian.Uint32(raw[off+1 : off+5]))
				cache[idx] = v
				return v, true
			}
			if size == 8 {
				if off+9 > len(raw) {
					return nil, false
				}
				v := math.Float64frombits(binary.BigEndian.Uint64(raw[off+1 : off+9]))
				cache[idx] = v
				return v, true
			}
			return nil, false
		case 0x3:
			if off+9 > len(raw) {
				return nil, false
			}
			seconds := math.Float64frombits(binary.BigEndian.Uint64(raw[off+1 : off+9]))
			v := applePlistEpoch.Add(time.Duration(seconds * float64(time.Second)))
			cache[idx] = v
			return v, true
		case 0x4:
			length, headerLen, ok := plistLength(raw, off, low)
			if !ok || off+headerLen+length > len(raw) {
				return nil, false
			}
			v := append([]byte(nil), raw[off+headerLen:off+headerLen+length]...)
			cache[idx] = v
			return v, true
		case 0x5:
			length, headerLen, ok := plistLength(raw, off, low)
			if !ok || off+headerLen+length > len(raw) {
				return nil, false
			}
			v := string(raw[off+headerLen : off+headerLen+length])
			cache[idx] = v
			return v, true
		case 0x6:
			length, headerLen, ok := plistLength(raw, off, low)
			if !ok || off+headerLen+length*2 > len(raw) {
				return nil, false
			}
			v := decodeUTF16BE(raw[off+headerLen : off+headerLen+length*2])
			cache[idx] = v
			return v, true
		case 0x8:
			size := int(low) + 1
			if size <= 0 || off+1+size > len(raw) {
				return nil, false
			}
			v, ok := readUnsigned(raw[off+1:], size)
			if !ok {
				return nil, false
			}
			cache[idx] = v
			return v, true
		case 0xA:
			length, headerLen, ok := plistLength(raw, off, low)
			if !ok || off+headerLen+length*objRefSize > len(raw) {
				return nil, false
			}
			items := make([]any, 0, length)
			base := off + headerLen
			for i := 0; i < length; i++ {
				ref, ok := readUnsigned(raw[base+i*objRefSize:], objRefSize)
				if !ok {
					return nil, false
				}
				refInt, ok := meta.SafecastUint64ToInt(ref)
				if !ok {
					return nil, false
				}
				item, ok := parseObject(refInt)
				if !ok {
					return nil, false
				}
				items = append(items, item)
			}
			cache[idx] = items
			return items, true
		case 0xD:
			length, headerLen, ok := plistLength(raw, off, low)
			if !ok || off+headerLen+length*objRefSize*2 > len(raw) {
				return nil, false
			}
			items := make(map[string]any, length)
			base := off + headerLen
			keysBase := base
			valsBase := base + length*objRefSize
			for i := 0; i < length; i++ {
				keyRef, ok := readUnsigned(raw[keysBase+i*objRefSize:], objRefSize)
				if !ok {
					return nil, false
				}
				keyRefInt, ok := meta.SafecastUint64ToInt(keyRef)
				if !ok {
					return nil, false
				}
				keyVal, ok := parseObject(keyRefInt)
				if !ok {
					return nil, false
				}
				key, ok := keyVal.(string)
				if !ok {
					continue
				}
				valRef, ok := readUnsigned(raw[valsBase+i*objRefSize:], objRefSize)
				if !ok {
					return nil, false
				}
				valRefInt, ok := meta.SafecastUint64ToInt(valRef)
				if !ok {
					return nil, false
				}
				val, ok := parseObject(valRefInt)
				if !ok {
					return nil, false
				}
				items[key] = val
			}
			cache[idx] = items
			return items, true
		default:
			return nil, false
		}
	}

	root, ok := parseObject(topObject)
	if !ok {
		return AppleRunTime{}, false
	}
	dict, ok := root.(map[string]any)
	if !ok {
		return AppleRunTime{}, false
	}

	var out AppleRunTime
	if v, ok := plistUint64(dict["flags"]); ok {
		if vv, ok := meta.SafecastUint64ToUint32(v); ok {
			out.Flags = vv
		}
	}
	if v, ok := plistUint64(dict["value"]); ok {
		out.Value = v
	}
	if v, ok := plistUint64(dict["timescale"]); ok {
		out.Scale = v
	}
	if v, ok := plistUint64(dict["epoch"]); ok {
		out.Epoch = v
	}
	return out, true
}

func parseRunTimeKeyed(raw []byte, timescaleIdx int) (AppleRunTime, bool) {
	var out AppleRunTime
	pos := timescaleIdx + len("timescale")
	if v, n, ok := readPlistInlineUintAt(raw, pos); ok {
		if vv, ok := meta.SafecastUint64ToUint32(v); ok {
			out.Flags = vv
		}
		pos += n
	}
	if v, n, ok := readPlistInlineUintAt(raw, pos); ok {
		out.Value = v
		pos += n
	}
	if v, n, ok := readPlistInlineUintAt(raw, pos); ok {
		out.Epoch = v
		pos += n
	}
	if v, _, ok := readPlistInlineUintAt(raw, pos); ok {
		out.Scale = v
	}
	if out.Scale == 0 && out.Value == 0 && out.Flags == 0 && out.Epoch == 0 {
		return AppleRunTime{}, false
	}
	return out, true
}

func readPlistInlineUintAt(raw []byte, pos int) (uint64, int, bool) {
	if pos >= len(raw) {
		return 0, 0, false
	}
	marker := raw[pos]
	if marker>>4 != 0x1 {
		return 0, 0, false
	}
	size := 1 << (marker & 0x0f)
	if size <= 0 || size > 8 || pos+1+size > len(raw) {
		return 0, 0, false
	}
	v, ok := readUnsigned(raw[pos+1:], size)
	if !ok {
		return 0, 0, false
	}
	return v, 1 + size, true
}

// ParseFloat64List decodes rational/float Apple maker-note values into dst.
func ParseFloat64List(raw []byte, order utils.ByteOrder, typ tag.Type, unitCount uint32, dst []float64) int {
	if len(dst) == 0 || unitCount == 0 {
		return 0
	}
	switch typ {
	case tag.TypeRational, tag.TypeSignedRational:
		n := min(int(unitCount), len(dst))
		if len(raw) < n*8 {
			n = len(raw) / 8
		}
		for i := 0; i < n; i++ {
			start := i * 8
			if typ == tag.TypeSignedRational {
				num := meta.SafecastUint32ToInt32Bits(order.Uint32(raw[start : start+4]))
				den := meta.SafecastUint32ToInt32Bits(order.Uint32(raw[start+4 : start+8]))
				if den == 0 {
					dst[i] = 0
					continue
				}
				dst[i] = float64(num) / float64(den)
				continue
			}
			num := order.Uint32(raw[start : start+4])
			den := order.Uint32(raw[start+4 : start+8])
			if den == 0 {
				dst[i] = 0
				continue
			}
			dst[i] = float64(num) / float64(den)
		}
		return n
	case tag.TypeFloat:
		n := min(int(unitCount), len(dst))
		if len(raw) < n*4 {
			n = len(raw) / 4
		}
		for i := 0; i < n; i++ {
			start := i * 4
			dst[i] = float64(math.Float32frombits(order.Uint32(raw[start : start+4])))
		}
		return n
	case tag.TypeDouble:
		n := min(int(unitCount), len(dst))
		if len(raw) < n*8 {
			n = len(raw) / 8
		}
		for i := 0; i < n; i++ {
			start := i * 8
			dst[i] = math.Float64frombits(order.Uint64(raw[start : start+8]))
		}
		return n
	default:
		return 0
	}
}

// ParseInt32List decodes integer Apple maker-note values into dst.
func ParseInt32List(raw []byte, order utils.ByteOrder, typ tag.Type, unitCount uint32, dst []int32) int {
	if len(dst) == 0 || unitCount == 0 {
		return 0
	}
	switch typ {
	case tag.TypeShort, tag.TypeLong, tag.TypeSignedShort, tag.TypeSignedLong:
	default:
		return 0
	}
	n := min(int(unitCount), len(dst))
	switch typ {
	case tag.TypeShort:
		if len(raw) < n*2 {
			n = len(raw) / 2
		}
		for i := 0; i < n; i++ {
			start := i * 2
			dst[i] = int32(order.Uint16(raw[start : start+2]))
		}
	case tag.TypeLong:
		if len(raw) < n*4 {
			n = len(raw) / 4
		}
		for i := 0; i < n; i++ {
			start := i * 4
			if v, ok := meta.SafecastUint32ToInt32(order.Uint32(raw[start : start+4])); ok {
				dst[i] = v
			}
		}
	case tag.TypeSignedShort:
		if len(raw) < n*2 {
			n = len(raw) / 2
		}
		for i := 0; i < n; i++ {
			start := i * 2
			dst[i] = int32(meta.SafecastUint16ToInt16Bits(order.Uint16(raw[start : start+2])))
		}
	case tag.TypeSignedLong:
		if len(raw) < n*4 {
			n = len(raw) / 4
		}
		for i := 0; i < n; i++ {
			start := i * 4
			dst[i] = meta.SafecastUint32ToInt32Bits(order.Uint32(raw[start : start+4]))
		}
	}
	return n
}

// ParseText converts Apple opaque payloads to a printable string.
func ParseText(raw []byte, maxBytes uint32) string {
	if maxBytes == 0 || len(raw) == 0 {
		return ""
	}
	if rawLen, ok := meta.SafecastIntToUint32(len(raw)); !ok || rawLen > maxBytes {
		raw = raw[:maxBytes]
	}
	raw = trimTrailingNULBytes(raw)
	if len(raw) == 0 {
		return ""
	}
	allPrintable := true
	for i := 0; i < len(raw); i++ {
		if raw[i] < 0x20 || raw[i] > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return string(bytes.TrimRight(raw, " \t\r\n"))
	}
	if len(raw) <= 512 {
		var out [512]byte
		for i := 0; i < len(raw); i++ {
			b := raw[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
			} else {
				out[i] = '.'
			}
		}
		return string(bytes.TrimRight(out[:len(raw)], " \t\r\n"))
	}
	out := make([]byte, len(raw))
	for i := 0; i < len(raw); i++ {
		b := raw[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
		} else {
			out[i] = '.'
		}
	}
	return string(bytes.TrimRight(out, " \t\r\n"))
}

// ParseMakerNoteVersion extracts the first Apple maker-note entry when it is
// the version tag. This is used as a compatibility fallback for files that
// place the version in the first slot of the Apple maker-note IFD.
func ParseMakerNoteVersion(raw []byte, order utils.ByteOrder) (int32, bool) {
	if len(raw) < 14 {
		return 0, false
	}
	if order.Uint16(raw[2:4]) != TagAppleMakerNoteVersion {
		return 0, false
	}
	if order.Uint16(raw[4:6]) != uint16(tag.TypeSignedLong) {
		return 0, false
	}
	if order.Uint32(raw[6:10]) != 1 {
		return 0, false
	}
	return meta.SafecastUint32ToInt32Bits(order.Uint32(raw[10:14])), true
}

var applePlistEpoch = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

func plistLength(raw []byte, off int, low byte) (length int, headerLen int, ok bool) {
	if low < 0x0f {
		return int(low), 1, true
	}
	if off+1 >= len(raw) {
		return 0, 0, false
	}
	marker := raw[off+1]
	if marker>>4 != 0x1 {
		return 0, 0, false
	}
	size := 1 << (marker & 0x0f)
	if size <= 0 || size > 8 || off+2+size > len(raw) {
		return 0, 0, false
	}
	v, ok := readUnsigned(raw[off+2:], size)
	if !ok {
		return 0, 0, false
	}
	length, ok = meta.SafecastUint64ToInt(v)
	if !ok {
		return 0, 0, false
	}
	return length, 2 + size, true
}

func readUnsigned(raw []byte, size int) (uint64, bool) {
	if size <= 0 || size > 8 || len(raw) < size {
		return 0, false
	}
	var v uint64
	for i := 0; i < size; i++ {
		v = (v << 8) | uint64(raw[i])
	}
	return v, true
}

func decodeUTF16BE(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	u16 := make([]uint16, len(raw)/2)
	for i := range u16 {
		u16[i] = binary.BigEndian.Uint16(raw[i*2:])
	}
	return string(utf16.Decode(u16))
}

func plistUint64(v any) (uint64, bool) {
	switch vv := v.(type) {
	case uint64:
		return vv, true
	case uint32:
		return uint64(vv), true
	case uint16:
		return uint64(vv), true
	case uint8:
		return uint64(vv), true
	case int64:
		if vv < 0 {
			return 0, false
		}
		return uint64(vv), true
	case int32:
		if vv < 0 {
			return 0, false
		}
		return uint64(vv), true
	case int:
		if vv < 0 {
			return 0, false
		}
		return uint64(vv), true
	case string:
		n, err := strconv.ParseUint(vv, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

func trimTrailingNULBytes(buf []byte) []byte {
	end := len(buf)
	for end > 0 && buf[end-1] == 0 {
		end--
	}
	return buf[:end]
}
