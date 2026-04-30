package canon

import (
	"math"

	"github.com/evanoberholster/imagemeta/meta"
)

// ---------------------------------------------------------------------------
// Low-level byte extraction helpers
// ---------------------------------------------------------------------------

func byteAt(buf []byte, off int) uint8 {
	if off < 0 || off >= len(buf) {
		return 0
	}
	return buf[off]
}

func u16LEAt(buf []byte, off int) uint16 {
	if off < 0 || off+1 >= len(buf) {
		return 0
	}
	return uint16(buf[off]) | uint16(buf[off+1])<<8
}

func u16BEAt(buf []byte, off int) uint16 {
	if off < 0 || off+1 >= len(buf) {
		return 0
	}
	return uint16(buf[off])<<8 | uint16(buf[off+1])
}

func u32LEAt(buf []byte, off int) uint32 {
	if off < 0 || off+3 >= len(buf) {
		return 0
	}
	return uint32(buf[off]) | uint32(buf[off+1])<<8 |
		uint32(buf[off+2])<<16 | uint32(buf[off+3])<<24
}

// ---------------------------------------------------------------------------
// Value conversion functions (matching ExifTool Canon.pm %ci* macros)
// ---------------------------------------------------------------------------

func ciFNumber(v uint8) meta.Aperture {
	if v == 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(float64(int(v)-8) / 16.0))
}

func ciExposureTime(v uint8) meta.ExposureTime {
	if v == 0 {
		return 0
	}
	return meta.ExposureTime(math.Exp2(4.0 * (1.0 - canonEV(float64(int16(v)-24)))))
}

func ciISO(v uint8) uint32 {
	if v == 0 {
		return 0
	}
	return uint32(math.Round(100.0 * math.Exp2(float64(v)/8.0-9.0)))
}

func ciTemperature(v uint8) int16 {
	if v == 0 {
		return 0
	}
	return int16(v) - 128
}

func ciMeasuredEV2(v uint8) float32 {
	if v == 0 {
		return 0
	}
	return float32(v)/8.0 - 6.0
}

func ciMacroMagnification(v uint8) float32 {
	if v == 0 {
		return 0
	}
	return float32(math.Exp((75.0 - float64(v)) * math.Ln2 * 3.0 / 40.0))
}

func ciFocalLength(v uint16) meta.FocalLength {
	return meta.FocalLength(float32(v))
}

func ciASCIIBytes(buf []byte, off, maxLen int) string {
	if off < 0 || maxLen <= 0 || off >= len(buf) {
		return ""
	}
	end := off + maxLen
	if end > len(buf) {
		end = len(buf)
	}
	slice := buf[off:end]
	for i, b := range slice {
		if b == 0 {
			return string(slice[:i])
		}
	}
	return string(slice)
}

func canonEV(val float64) float64 {
	if val >= 0 {
		return val
	}
	return -((-val) * 2.0 / 3.0)
}

// ---------------------------------------------------------------------------
// Exported helpers for use by the exif parser package.

// CIByteAt reads a byte at off or returns 0.
func CIByteAt(buf []byte, off int) uint8 { return byteAt(buf, off) }

// CIU16LEAt reads a uint16 LE at off or returns 0.
func CIU16LEAt(buf []byte, off int) uint16 { return u16LEAt(buf, off) }

// CIU16BEAt reads a uint16 BE at off or returns 0.
func CIU16BEAt(buf []byte, off int) uint16 { return u16BEAt(buf, off) }

// CIU32LEAt reads a uint32 LE at off or returns 0.
func CIU32LEAt(buf []byte, off int) uint32 { return u32LEAt(buf, off) }

// CIFNumber converts a camera-info FNumber byte to an Aperture.
func CIFNumber(v uint8) meta.Aperture { return ciFNumber(v) }

// CIExposureTime converts a camera-info ExposureTime byte to an ExposureTime.
func CIExposureTime(v uint8) meta.ExposureTime { return ciExposureTime(v) }

// CIISO converts a camera-info ISO byte to a uint32 ISO value.
func CIISO(v uint8) uint32 { return ciISO(v) }

// CITemperature converts a camera-info temperature byte to Celsius int16.
func CITemperature(v uint8) int16 { return ciTemperature(v) }

// CIMeasuredEV2 converts a camera-info MeasuredEV2 byte.
func CIMeasuredEV2(v uint8) float32 { return ciMeasuredEV2(v) }

// CIMacroMagnification converts a camera-info MacroMagnification byte.
func CIMacroMagnification(v uint8) float32 { return ciMacroMagnification(v) }

// CIFocalLength converts a camera-info focal-length uint16 to a FocalLength.
func CIFocalLength(v uint16) meta.FocalLength { return ciFocalLength(v) }

// CIAsciiBytes extracts a NUL-terminated ASCII string from buf at off with max length.
func CIAsciiBytes(buf []byte, off, maxLen int) string { return ciASCIIBytes(buf, off, maxLen) }
