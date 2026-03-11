package canon

import "fmt"

// Ev - ported from Phil Harvey's exiftool
// Updated May-10-2020
// https://github.com/exiftool/exiftool/lib/Image/ExifTool/Canon.pm
func Ev(val int16) int16 {
	var sign int16
	if val < 0 {
		val = -val
		sign = -1
	} else {
		sign = 1
	}
	frac := val & 0x1f
	val -= frac
	// Convert 1/3 and 2/3 codes
	switch frac {
	case 0x0c:
		frac = 0x20 / 3
	case 0x14:
		frac = 0x40 / 3
	}
	return sign * (val + frac) / 0x20
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
	if len(af) < 4 {
		return nil, nil, fmt.Errorf("canon: af data too short: got %d words, need at least 4", len(af))
	}

	// AFInfo2 layout (0x0026/0x003c): NumAFPoints is sequence index 2, stored at af[3].
	validPoints := int(af[3])
	if validPoints > 0 {
		count := (validPoints + 15) / 16
		off := 8 + (validPoints * 4)
		end := off + count
		if off >= 0 && end >= off && end <= len(af) {
			inFocus = decodeBits(af[off:end], 16)
			// AFPointsSelected is EOS-only in ExifTool, but decode if present to preserve API behavior.
			if end+count <= len(af) {
				selected = decodeBits(af[end:end+count], 16)
			}
			return
		}
	}

	// AFInfo layout (0x0012): NumAFPoints is sequence index 0, stored at af[0].
	validPoints = int(af[0])
	if validPoints <= 0 {
		return nil, nil, fmt.Errorf("canon: unexpected NumAFPoints %d", validPoints)
	}
	count := (validPoints + 15) / 16
	off := 8 + (validPoints * 2)
	end := off + count
	if off < 0 || end < off || end > len(af) {
		return nil, nil, fmt.Errorf(
			"canon: af data too short for points-in-focus: got %d words, need %d",
			len(af),
			end,
		)
	}
	inFocus = decodeBits(af[off:end], 16)
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
	if len(af) < 6 {
		return nil
	}

	validPoints := int(af[3])
	if validPoints <= 0 {
		return nil
	}

	// Requires width/height/x/y blocks of NumAFPoints each, starting at offset 8.
	required := 8 + (validPoints * 4)
	if required < 8 || required > len(af) {
		return nil
	}

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
