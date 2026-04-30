package canon

import (
	"math"
	"strings"

	"github.com/evanoberholster/imagemeta/meta"
)

// CanonEV decodes Canon's hex-based EV codes using the int16 convention.
func CanonEV(code int16) float64 {
	val := int(code)
	sign := 1.0
	if val < 0 {
		val = -val
		sign = -1
	}
	frac := val & 0x1f
	base := val - frac
	fracEV := float64(frac)
	switch frac {
	case 0x0c:
		fracEV = 32.0 / 3.0
	case 0x14:
		fracEV = 64.0 / 3.0
	}
	return sign * (float64(base) + fracEV) / 32.0
}

// CameraSettingValue normalizes a Canon CameraSettings int16 value.
// Sentinel 0x7fff is mapped to 0; all other values pass through.
func CameraSettingValue(v int16) int16 {
	if v == 0x7fff {
		return 0
	}
	return v
}

// ClarityValue normalizes a Canon Clarity (index 51) int16 value.
func ClarityValue(v int16) int16 {
	if v == 0x7fff {
		return 0
	}
	return v
}

// MaxApertureFromCode converts a Canon CameraSettings aperture code to f-number.
func MaxApertureFromCode(raw uint16) meta.Aperture {
	code := int16(raw)
	if code <= 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(CanonEV(code) * 0.5))
}

// DisplayApertureFromCode converts DisplayAperture (seq 35) raw value / 10.
func DisplayApertureFromCode(raw uint16) meta.Aperture {
	if raw == 0 {
		return 0
	}
	return meta.Aperture(float32(raw) / 10.0)
}

// ShotISO converts a Canon ShotInfo ISO code to a float32 ISO value.
func ShotISO(code int16) float32 {
	if code == 0 {
		return 100
	}
	return float32(100.0 * math.Exp2(float64(code-160)/32.0))
}

// ShotMeasuredEV converts a Canon ShotInfo MeasuredEV code.
func ShotMeasuredEV(code int16) float32 {
	if code == 0 {
		return 0
	}
	return float32(CanonEV(code) + 5.0)
}

// ShotActualISO computes the actual ISO from autoISO and baseISO.
func ShotActualISO(autoISO, baseISO float32) float32 {
	if autoISO <= 0 || baseISO <= 0 {
		return 0
	}
	return (autoISO * baseISO) / 100.0
}

// ShotAperture converts a Canon ShotInfo aperture code.
func ShotAperture(code int16) meta.Aperture {
	if code == 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(CanonEV(code) * 0.5))
}

// ShotExposureTime converts a Canon ShotInfo exposure time code.
func ShotExposureTime(code int16, legacy20D350D bool) meta.ExposureTime {
	if code == 0 {
		return 0
	}
	if legacy20D350D {
		return meta.ExposureTime(math.Exp2(float64(code-640) / 32.0))
	}
	return meta.ExposureTime(math.Exp2(-CanonEV(code)))
}

// ShotExposureCompensation converts a Canon ShotInfo exposure compensation code.
func ShotExposureCompensation(code int16) float32 {
	if code == 0 {
		return 0
	}
	return float32(CanonEV(code))
}

// ShotFlashGuideNumber converts a Canon ShotInfo flash guide number code.
func ShotFlashGuideNumber(raw int16) float32 {
	if raw < 0 {
		return 0
	}
	return float32(raw) / 32.0
}

// ShotMeasuredEV2 converts a Canon ShotInfo MeasuredEV2 raw code.
func ShotMeasuredEV2(raw int16) float32 {
	if raw == 0 {
		return 0
	}
	return float32(raw)/8.0 - 6.0
}

// FocalPlaneSizeMM converts a Canon CameraInfo FocalPlane raw uint16 to mm.
func FocalPlaneSizeMM(raw uint16) float32 {
	if raw == 0 {
		return 0
	}
	return float32(raw) * 25.4 / 1000.0
}

// NormalizeFirmwareVersion strips the "Firmware Version " prefix from a string.
func NormalizeFirmwareVersion(s string) string {
	s = strings.TrimSpace(s)
	return strings.TrimPrefix(s, "Firmware Version ")
}

// HexBytes converts raw bytes to a lowercase hex string.
func HexBytes(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	const table = "0123456789abcdef"
	var out strings.Builder
	out.Grow(len(b) * 2)
	for i := range b {
		v := b[i]
		out.WriteByte(table[v>>4])
		out.WriteByte(table[v&0x0f])
	}
	return out.String()
}
