package meta

import (
	"bytes"
	"math"
	"math/bits"
)

// parseInt parses a []byte of a string representation of an int64 value and returns the value
//func parseInt(buf []byte) (i int64) {
//	var neg bool
//	if buf[0] == '-' {
//		buf = buf[1:]
//		neg = true
//	}
//	i = int64(parseUint(buf))
//	if neg {
//		i *= -1
//	}
//	return
//}

// parseUint parses a []byte of a string representation of a uint64 value and returns the value.
func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

// parseUntil parses a []byte and splits the []byte at delimiter.
// Returns a and b without delimiter present.
//func parseUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
//	for i := 0; i < len(buf); i++ {
//		if buf[i] == delimiter {
//			a = buf[:i]
//			if i < len(buf)+1 {
//				b = buf[i+1:]
//				return
//			}
//			return
//		}
//	}
//	return buf, nil
//}

var closeTagXMP = []byte("</x:xmpmeta>")

// CleanXMPSuffixWhiteSpace returns the same slice with the whitespace after "</x:xmpmeta>" removed.
func CleanXMPSuffixWhiteSpace(buf []byte) []byte {
	for i := len(buf) - 1; i > 12; i-- {
		if buf[i] == '>' && buf[i-1] == 'a' {
			// </x:xmpmeta>
			if bytes.Equal(closeTagXMP, buf[i-11:i+1]) {
				buf = buf[:i+1]
				return buf
			}
		}
	}
	return buf
}

func unsafeGetBytes(s string) (b []byte) {
	return []byte(s)
	//(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	//(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	//(*reflect.SliceHeader)(unsafe.Pointer(&b)).Len = len(s)
	//return
}

func SafecastIntToUint32(v int) (uint32, bool) {
	if v < 0 || v > math.MaxUint32 {
		return 0, false
	}
	return uint32(v), true
}

func SafecastIntToUint(v int) (uint, bool) {
	if v < 0 {
		return 0, false
	}
	return uint(v), true
}

func SafecastIntToUint8(v int) (uint8, bool) {
	if v < 0 || v > math.MaxUint8 {
		return 0, false
	}
	return uint8(v), true
}

func SafecastIntToInt32(v int) (int32, bool) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

func SafecastIntToInt8(v int) (int8, bool) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		return 0, false
	}
	return int8(v), true
}

func SafecastUintToInt(v uint) (int, bool) {
	if bits.UintSize == 32 {
		if v > uint(math.MaxInt32) {
			return 0, false
		}
		return int(v), true
	}
	if v > uint(math.MaxInt64) {
		return 0, false
	}
	return int(v), true
}

func SafecastUintToUint16(v uint) (uint16, bool) {
	if v > uint(math.MaxUint16) {
		return 0, false
	}
	return uint16(v), true
}

func SafecastUint32ToUint16(v uint32) (uint16, bool) {
	if v > math.MaxUint16 {
		return 0, false
	}
	return uint16(v), true
}

func SafecastUint32ToUint8(v uint32) (uint8, bool) {
	if v > math.MaxUint8 {
		return 0, false
	}
	return uint8(v), true
}

func SafecastUint16ToUint8(v uint16) (uint8, bool) {
	if v > math.MaxUint8 {
		return 0, false
	}
	return uint8(v), true
}

func SafecastUint32ToInt32(v uint32) (int32, bool) {
	if v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

func SafecastUint32ToInt16(v uint32) (int16, bool) {
	if v > math.MaxInt16 {
		return 0, false
	}
	return int16(v), true
}

func SafecastUint64ToUint32(v uint64) (uint32, bool) {
	if v > math.MaxUint32 {
		return 0, false
	}
	return uint32(v), true
}

func SafecastUint64ToUint16(v uint64) (uint16, bool) {
	if v > math.MaxUint16 {
		return 0, false
	}
	return uint16(v), true
}

func SafecastUint64ToInt16(v uint64) (int16, bool) {
	if v > math.MaxInt16 {
		return 0, false
	}
	return int16(v), true
}

func SafecastInt32ToInt16(v int32) (int16, bool) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		return 0, false
	}
	return int16(v), true
}

func SafecastInt32ToUint32(v int32) (uint32, bool) {
	if v < 0 {
		return 0, false
	}
	return uint32(v), true
}

func SafecastUint64ToInt(v uint64) (int, bool) {
	if bits.UintSize == 32 {
		if v > math.MaxInt32 {
			return 0, false
		}
		return int(v), true
	}
	if v > math.MaxInt64 {
		return 0, false
	}
	return int(v), true
}

func SafecastUint64ToInt64(v uint64) (int64, bool) {
	if v > math.MaxInt64 {
		return 0, false
	}
	return int64(v), true
}

func SafecastInt64ToUint64(v int64) (uint64, bool) {
	if v < 0 {
		return 0, false
	}
	return uint64(v), true
}

func SafecastInt64ToUint32(v int64) (uint32, bool) {
	if v < 0 || v > math.MaxUint32 {
		return 0, false
	}
	return uint32(v), true
}

func SafecastInt16ToUint16Bits(v int16) uint16 {
	//nolint:gosec // G115: deliberate two's-complement bit reinterpretation.
	return uint16(v)
}

func SafecastInt32ToUint32Bits(v int32) uint32 {
	//nolint:gosec // G115: deliberate two's-complement bit reinterpretation.
	return uint32(v)
}

func SafecastUint8ToInt8Bits(v uint8) int8 {
	//nolint:gosec // G115: deliberate two's-complement bit reinterpretation.
	return int8(v)
}

func SafecastUint16ToInt16Bits(v uint16) int16 {
	//nolint:gosec // G115: deliberate two's-complement bit reinterpretation.
	return int16(v)
}

func SafecastUint32ToInt32Bits(v uint32) int32 {
	//nolint:gosec // G115: deliberate two's-complement bit reinterpretation.
	return int32(v)
}
