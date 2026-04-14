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

	layout, ok := parseAFLayout(af)
	if !ok {
		return nil, nil, fmt.Errorf("canon: unsupported AFInfo payload layout")
	}

	inFocus = decodeBits(af[layout.inFocusStart:layout.inFocusStart+layout.maskWordCount], 16)
	if layout.selectedStart >= 0 {
		selectedEnd := layout.selectedStart + layout.maskWordCount
		if selectedEnd <= len(af) {
			selected = decodeBits(af[layout.selectedStart:selectedEnd], 16)
		}
	}
	return inFocus, selected, nil
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
	layout, ok := parseAFLayout(af)
	if !ok {
		return nil
	}

	switch layout.kind {
	case afLayoutLegacy:
		return parseLegacyAFArea(af, layout)
	case afLayoutInfo2:
		raw := parseAFInfo2AFArea(af, layout)
		if len(raw) == 0 {
			return nil
		}
		afPoints = make([]AFPoint, len(raw))
		xAdjust := int16(layout.canonImageWidth / 2)
		yAdjust := int16(layout.canonImageHeight / 2)
		for i := range raw {
			w, h, x, y := raw[i][0], raw[i][1], raw[i][2], raw[i][3]
			x += xAdjust - (w / 2)
			y += yAdjust - (h / 2)
			afPoints[i] = NewAFPoint(w, h, x, y)
		}
		return afPoints
	default:
		return nil
	}
}

// ParseAFArea returns the raw Canon AF area tuples matching ExifTool's
// width/height/x/y tables for both legacy AFInfo and AFInfo2/AFInfo3 records.
func ParseAFArea(af []uint16) []AFPoint {
	layout, ok := parseAFLayout(af)
	if !ok {
		return nil
	}
	switch layout.kind {
	case afLayoutLegacy:
		return parseLegacyAFArea(af, layout)
	case afLayoutInfo2:
		return parseAFInfo2AFArea(af, layout)
	default:
		return nil
	}
}

type afLayoutKind uint8

const (
	afLayoutUnknown afLayoutKind = iota
	afLayoutLegacy
	afLayoutInfo2
)

type afLayout struct {
	kind             afLayoutKind
	numPoints        int
	maskWordCount    int
	inFocusStart     int
	selectedStart    int
	canonImageWidth  uint16
	canonImageHeight uint16
	areaWidth        uint16
	areaHeight       uint16
}

func parseAFLayout(af []uint16) (afLayout, bool) {
	if layout, ok := parseAFInfo2Layout(af); ok {
		return layout, true
	}
	if layout, ok := parseLegacyAFInfoLayout(af); ok {
		return layout, true
	}
	return afLayout{}, false
}

func parseAFInfo2Layout(af []uint16) (afLayout, bool) {
	if len(af) < 8 {
		return afLayout{}, false
	}
	numPoints := int(af[2])
	if numPoints <= 0 {
		return afLayout{}, false
	}
	maskWordCount := bitWordCount(numPoints)
	inFocusStart := 8 + (numPoints * 4)
	inFocusEnd := inFocusStart + maskWordCount
	if inFocusStart < 8 || inFocusEnd < inFocusStart || inFocusEnd > len(af) {
		return afLayout{}, false
	}

	selectedStart := -1
	remaining := len(af) - inFocusEnd
	// ExifTool only exposes AFPointsSelected for EOS AFInfo2 records. Without
	// model context, only treat the record as having a selected-mask when the
	// payload length exactly matches the EOS layout.
	if remaining == maskWordCount {
		selectedStart = inFocusEnd
	}

	return afLayout{
		kind:             afLayoutInfo2,
		numPoints:        numPoints,
		maskWordCount:    maskWordCount,
		inFocusStart:     inFocusStart,
		selectedStart:    selectedStart,
		canonImageWidth:  af[4],
		canonImageHeight: af[5],
	}, true
}

func parseLegacyAFInfoLayout(af []uint16) (afLayout, bool) {
	if len(af) < 8 {
		return afLayout{}, false
	}
	numPoints := int(af[0])
	if numPoints <= 0 {
		return afLayout{}, false
	}
	maskWordCount := bitWordCount(numPoints)
	inFocusStart := 8 + (numPoints * 2)
	inFocusEnd := inFocusStart + maskWordCount
	if inFocusStart < 8 || inFocusEnd < inFocusStart || inFocusEnd > len(af) {
		return afLayout{}, false
	}
	return afLayout{
		kind:          afLayoutLegacy,
		numPoints:     numPoints,
		maskWordCount: maskWordCount,
		inFocusStart:  inFocusStart,
		selectedStart: -1,
		areaWidth:     af[6],
		areaHeight:    af[7],
	}, true
}

func parseLegacyAFArea(af []uint16, layout afLayout) []AFPoint {
	if layout.numPoints <= 0 {
		return nil
	}
	xStart := 8
	yStart := xStart + layout.numPoints
	if yStart+layout.numPoints > len(af) {
		return nil
	}
	out := make([]AFPoint, layout.numPoints)
	for i := 0; i < layout.numPoints; i++ {
		out[i] = NewAFPoint(
			int16(layout.areaWidth),
			int16(layout.areaHeight),
			int16(af[xStart+i]),
			int16(af[yStart+i]),
		)
	}
	return out
}

func parseAFInfo2AFArea(af []uint16, layout afLayout) []AFPoint {
	if layout.numPoints <= 0 {
		return nil
	}
	widthStart := 8
	heightStart := widthStart + layout.numPoints
	xStart := heightStart + layout.numPoints
	yStart := xStart + layout.numPoints
	if yStart+layout.numPoints > len(af) {
		return nil
	}
	out := make([]AFPoint, layout.numPoints)
	for i := 0; i < layout.numPoints; i++ {
		out[i] = NewAFPoint(
			int16(af[widthStart+i]),
			int16(af[heightStart+i]),
			int16(af[xStart+i]),
			int16(af[yStart+i]),
		)
	}
	return out
}

func bitWordCount(pointCount int) int {
	if pointCount <= 0 {
		return 0
	}
	return (pointCount + 15) / 16
}
