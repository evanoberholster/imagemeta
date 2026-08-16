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

// MaxApertureFromCode converts a Canon CameraSettings aperture code to f-number.
func MaxApertureFromCode(raw uint16) meta.Aperture {
	code := meta.SafecastUint16ToInt16Bits(raw)
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

// ShotCameraTemperature converts a Canon ShotInfo raw camera temperature
// code. Returns the temperature in Celsius, or 0 if not available.
func ShotCameraTemperature(raw int16, modelID CanonCameraModel) int16 {
	if raw == 0 || !ModelIsEOS(modelID) || ModelUsesLegacyShotInfo(modelID) {
		return 0
	}
	return raw - 128
}

// ModelIsEOS reports whether the CanonCameraModel is an EOS body.
func ModelIsEOS(modelID CanonCameraModel) bool {
	switch modelID {
	case CanonModelEOSD30, CanonModelEOSD60, CanonModelEOSM3, CanonModelEOSM10,
		CanonModelEOSM5, CanonModelEOSM100, CanonModelEOSM6, CanonModelEOSM50,
		CanonModelEOSC50, CanonModelEOSC300, CanonModelEOSC200,
		CanonModelEOS1D, CanonModelEOS1DS, CanonModelEOS10D,
		CanonModelEOS1DMarkIII, CanonModelEOSDigitalRebel,
		CanonModelEOS1DMarkII, CanonModelEOS20D, CanonModelEOSDigitalRebelXSi,
		CanonModelEOS1DsMarkII, CanonModelEOSDigitalRebelXT, CanonModelEOS40D,
		CanonModelEOS5D, CanonModelEOS1DsMarkIII, CanonModelEOS5DMarkII,
		CanonModelEOS1DMarkIIN, CanonModelEOS30D, CanonModelEOSDigitalRebelXTi,
		CanonModelEOS7D, CanonModelEOSRebelT1i, CanonModelEOSRebelXS,
		CanonModelEOS50D, CanonModelEOS1DX, CanonModelEOSRebelT2i,
		CanonModelEOS1DMarkIV, CanonModelEOS5DMarkIII, CanonModelEOSRebelT3i,
		CanonModelEOS60D, CanonModelEOSRebelT3, CanonModelEOS7DMarkII,
		CanonModelEOSRebelT4i, CanonModelEOS6D, CanonModelEOS1DC,
		CanonModelEOS70D, CanonModelEOSRebelT5i, CanonModelEOSRebelT5,
		CanonModelEOS1DXMarkII, CanonModelEOSM, CanonModelEOS80D,
		CanonModelEOSM2, CanonModelEOSRebelSL1, CanonModelEOSRebelT6s,
		CanonModelEOS5DMarkIV, CanonModelEOS5DS, CanonModelEOSRebelT6i,
		CanonModelEOS5DSR, CanonModelEOSRebelT6, CanonModelEOSRebelT7i,
		CanonModelEOS6DMarkII, CanonModelEOS77D, CanonModelEOSRebelSL2,
		CanonModelEOSR5, CanonModelEOSRebelT100, CanonModelEOSR,
		CanonModelEOS1DXMarkIII, CanonModelEOSRebelT7, CanonModelEOSRP,
		CanonModelEOSRebelT8i, CanonModelEOSSL3, CanonModelEOS90D,
		CanonModelEOSR3, CanonModelEOSR6, CanonModelEOSR7,
		CanonModelEOSR10, CanonModelEOSM50MarkII, CanonModelEOSR50,
		CanonModelEOSR6MarkII, CanonModelEOSR8, CanonModelEOSR1,
		CanonModelEOSR5MarkII, CanonModelEOSR100, CanonModelEOSR50V,
		CanonModelEOSR6MarkIII, CanonModelEOSD2000C, CanonModelEOSD6000C:
		return true
	default:
		return false
	}
}

// ModelUsesLegacyShotInfo reports whether the camera uses the legacy ShotInfo format.
func ModelUsesLegacyShotInfo(modelID CanonCameraModel) bool {
	switch modelID {
	case CanonModelEOS1D, CanonModelEOS1DS, CanonModelEOSD30, CanonModelEOSD60:
		return true
	default:
		return false
	}
}

// ModelUsesLegacyShutterCount reports whether the camera stores shutter count
// in the legacy FileInfo format.
func ModelUsesLegacyShutterCount(modelID CanonCameraModel) bool {
	switch modelID {
	case CanonModelEOS1D, CanonModelEOS1DS:
		return true
	default:
		return false
	}
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
