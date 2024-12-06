package canon

import (
	"fmt"
	"math"
)

// CanonEv converts Canon hex-based EV (modulo 0x20) to real number
// https://github.com/exiftool/exiftool/lib/Image/ExifTool/Canon.pm
// Examples:
//
//	0x00 -> 0
//	0x0c -> 0.33333
//	0x10 -> 0.5
//	0x14 -> 0.66666
//	0x20 -> 1
func CanonEv(val float32) float32 {
	// Handle negative values
	sign := float32(1.0)
	if val < 0 {
		val = -val
		sign = -1.0
	}

	// Extract fraction (lower 5 bits)
	frac := float32(int(val) & 0x1f)

	// Remove fraction from value
	val -= frac

	// Convert 1/3 and 2/3 codes
	switch int(frac) {
	case 0x0c:
		frac = float32(0x20) / 3.0
	case 0x14:
		frac = float32(0x40) / 3.0
	}

	// Calculate final value
	return sign * (val + frac) / float32(0x20)
}

// CanonEvInv converts a number to Canon hex-based EV (modulo 0x20)
// This is the inverse of CanonEv
func CanonEvInv(num float32) float32 {
	// Handle negative values
	sign := 1.0
	if num < 0 {
		num = -num
		sign = -1.0
	}

	// Split into integer and fraction parts
	val := math.Floor(float64(num))
	frac := float64(num) - val

	// Convert fraction to Canon EV code
	var fracCode float64
	switch {
	case math.Abs(frac-0.33) < 0.05:
		fracCode = 0x0c
	case math.Abs(frac-0.67) < 0.05:
		fracCode = 0x14
	default:
		fracCode = math.Floor(frac*float64(0x20) + 0.5)
	}

	// Calculate final value
	return float32(sign * (val*float64(0x20) + fracCode))
}

// TempConv - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool/Canon.pm
func TempConv(val uint16) int16 {
	if val == 0 {
		return 0
	}
	return int16(val) - 128
}

// PointsInFocus returns AFPoints that are in focus and AFPoints that are selected
func PointsInFocus(af []uint16) (inFocus []int, selected []int, err error) {
	validPoints := int(af[3])
	var count int
	// NumAFPoints may be 7, 9, 11, 19, 31, 45 or 61, depending on the camera model.
	switch validPoints {
	case 7:
		count = 1 // 1
	case 9, 11:
		count = 1 // 1
	case 19, 31:
		count = 2 // 2
	case 45:
		count = 3 // 3
	case 61:
		count = 4 // 4
	case 65:
		count = 5 // 5
	case 1053:
		count = 66
	default:
		panic(fmt.Errorf("error parsing AFPoints from Canon Makernote. Expected 7, 9, 11, 19, 31, 45 or 61 got %d", validPoints))
	}
	off := 8 + (validPoints * 4)
	inFocus = decodeBits(af[off:off+count], 16)
	selected = decodeBits(af[off+count:off+count+count], 16)
	return
}

// decodeBits - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool.pm
func decodeBits(vals []uint16, bits int) (list []int) {
	var num int
	var n int
	for _, a := range vals {
		for i := 0; i < bits; i++ {
			n = i + num
			if a&(1<<uint(i)) > 0 {
				list = append(list, n)
			}
		}
		num += bits
	}
	return
}

// ParseAFPoints returns []AFPoint
func ParseAFPoints(af []uint16) (afPoints []AFPoint) {
	validPoints := int(af[3])
	// AFPoints
	afPoints = make([]AFPoint, validPoints)
	xAdjust := int16(af[4] / 2) // Adjust x-axis
	yAdjust := int16(af[5] / 2) // Adjust y-axis

	for i := 0; i < validPoints; i++ { // Start at an offset of 8
		offset := 8 + i
		w := int16(af[offset])
		h := int16(af[offset+validPoints])
		x := int16(af[offset+(2*validPoints)]) + xAdjust - (w / 2)
		y := int16(af[offset+(3*validPoints)]) + yAdjust - (h / 2)
		afPoints[i] = NewAFPoint(w, h, x, y)
	}
	return
}
