package exif

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/apple"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func (r *Reader) parseAppleRunTime(t tag.Entry) apple.AppleRunTime {
	raw := r.parseOpaqueBytes(t, 4096)
	if len(raw) == 0 {
		return apple.AppleRunTime{}
	}
	rt, ok := apple.ParseRunTime(raw)
	if !ok {
		return apple.AppleRunTime{}
	}
	return rt
}

func (r *Reader) parseAppleFloat64(t tag.Entry) float64 {
	switch t.Type {
	case tag.TypeRational, tag.TypeSignedRational:
		return r.parseSignedRationalFloat64(t)
	case tag.TypeFloat:
		if !t.IsEmbedded() {
			buf, err := r.readTagBytes(t, 4)
			if err != nil || len(buf) < 4 {
				return 0
			}
			if t.ByteOrder == utils.BigEndian {
				return float64(math.Float32frombits(binary.BigEndian.Uint32(buf[:4])))
			}
			return float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[:4])))
		}
		return float64(math.Float32frombits(t.EmbeddedLong()))
	case tag.TypeDouble:
		buf, err := r.readTagBytes(t, 8)
		if err != nil || len(buf) < 8 {
			return 0
		}
		if t.ByteOrder == utils.BigEndian {
			return math.Float64frombits(binary.BigEndian.Uint64(buf[:8]))
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(buf[:8]))
	case tag.TypeShort, tag.TypeLong:
		return float64(r.parseMakerNoteUint32(t))
	case tag.TypeSignedShort:
		return float64(r.parseMakerNoteInt16(t))
	case tag.TypeSignedLong:
		var raw [1]int32
		if r.parseInt32List(t, raw[:]) == 0 {
			return 0
		}
		return float64(raw[0])
	default:
		return 0
	}
}

func (r *Reader) parseAppleFloat64List(t tag.Entry, dst []float64) {
	if len(dst) == 0 || t.UnitCount == 0 {
		return
	}
	switch t.Type {
	case tag.TypeRational, tag.TypeSignedRational:
		n := min(int(t.UnitCount), len(dst))
		maxBytes, ok := meta.SafecastIntToUint32(n * 8)
		if !ok {
			return
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil {
			return
		}
		apple.ParseFloat64List(buf, t.ByteOrder, t.Type, t.UnitCount, dst)
		return
	case tag.TypeFloat:
		n := min(int(t.UnitCount), len(dst))
		maxBytes, ok := meta.SafecastIntToUint32(n * 4)
		if !ok {
			return
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil {
			return
		}
		apple.ParseFloat64List(buf, t.ByteOrder, t.Type, t.UnitCount, dst)
		return
	case tag.TypeDouble:
		n := min(int(t.UnitCount), len(dst))
		maxBytes, ok := meta.SafecastIntToUint32(n * 8)
		if !ok {
			return
		}
		buf, err := r.readTagBytes(t, maxBytes)
		if err != nil {
			return
		}
		apple.ParseFloat64List(buf, t.ByteOrder, t.Type, t.UnitCount, dst)
		return
	default:
		return
	}
}

func (r *Reader) parseAppleInt32List(t tag.Entry, dst []int32) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	switch t.Type {
	case tag.TypeShort, tag.TypeLong, tag.TypeSignedShort, tag.TypeSignedLong:
	default:
		return 0
	}
	n := min(int(t.UnitCount), len(dst))
	var count uint32
	switch t.Type {
	case tag.TypeShort, tag.TypeSignedShort:
		converted, ok := meta.SafecastIntToUint32(n * 2)
		if !ok {
			return 0
		}
		count = converted
	default:
		converted, ok := meta.SafecastIntToUint32(n * 4)
		if !ok {
			return 0
		}
		count = converted
	}
	buf, err := r.readTagBytes(t, count)
	if err != nil {
		return 0
	}
	return apple.ParseInt32List(buf, t.ByteOrder, t.Type, t.UnitCount, dst)
}

func (r *Reader) parseAppleInt32(t tag.Entry) int32 {
	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeShort:
			return int32(t.EmbeddedShort())
		case tag.TypeSignedShort:
			return int32(meta.SafecastUint16ToInt16Bits(t.EmbeddedShort()))
		case tag.TypeLong, tag.TypeSignedLong:
			return meta.SafecastUint32ToInt32Bits(t.EmbeddedLong())
		}
	}
	var dst [1]int32
	if r.parseAppleInt32List(t, dst[:]) == 0 {
		return 0
	}
	return dst[0]
}

func (r *Reader) parseAppleTextOrNumeric(t tag.Entry, maxBytes uint32) string {
	if s := r.parseDisplayString(t, maxBytes); s != "" {
		return s
	}
	switch t.Type {
	case tag.TypeShort, tag.TypeLong:
		return strconv.FormatUint(uint64(r.parseMakerNoteUint32(t)), 10)
	case tag.TypeSignedShort:
		return strconv.FormatInt(int64(r.parseMakerNoteInt16(t)), 10)
	case tag.TypeSignedLong:
		var raw [1]int32
		if r.parseInt32List(t, raw[:]) == 0 {
			return ""
		}
		return strconv.FormatInt(int64(raw[0]), 10)
	case tag.TypeRational, tag.TypeSignedRational, tag.TypeFloat, tag.TypeDouble:
		return strconv.FormatFloat(r.parseAppleFloat64(t), 'f', -1, 64)
	default:
		return ""
	}
}
