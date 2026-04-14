package canon

import (
	"fmt"
	"math/bits"
)

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

// AFInfo2DecodeConfig controls optional AFInfo2/AFInfo3 decode work.
type AFInfo2DecodeConfig struct {
	Source         AFInfoSource
	EOS            bool
	AFInfo3        bool
	DecodeCoords   bool
	DecodePoints   bool
	DecodeInFocus  bool
	DecodeSelected bool
}

// BitWordCount returns the number of uint16 bitset words needed for pointCount.
func BitWordCount(pointCount int) int {
	return bitWordCount(pointCount)
}

// DecodeAFInfo decodes Canon AFInfo (0x0012) words into AFInfo.
func DecodeAFInfo(af []uint16, eos bool, afInfoCount int) AFInfo {
	n := len(af)
	dst := AFInfo{
		Source:           AFInfoSourceAFInfo,
		NumAFPoints:      u16At(af, n, 0),
		ValidAFPoints:    u16At(af, n, 1),
		CanonImageWidth:  u16At(af, n, 2),
		CanonImageHeight: u16At(af, n, 3),
		AFImageWidth:     u16At(af, n, 4),
		AFImageHeight:    u16At(af, n, 5),
		AFAreaWidth:      u16At(af, n, 6),
		AFAreaHeight:     u16At(af, n, 7),
	}

	layout, ok := parseLegacyAFInfoLayout(af)
	if !ok {
		return dst
	}
	dst.AFPointsInFocusBits = decodeBitWordsRange(af, n, layout.inFocusStart, layout.maskWordCount)
	if !eos {
		dst.PrimaryAFPoint = legacyAFInfoPrimary(af, n, layout.inFocusStart+layout.maskWordCount, afInfoCount)
	}
	areas := parseLegacyAFArea(af, layout)
	dst.AFArea = areas
	// AFInfo (0x0012) stores width/height/x/y directly in the AF area tuples.
	dst.AFPoints = areas
	return dst
}

// DecodeAFInfo2 decodes Canon AFInfo2/AFInfo3 words into AFInfo.
func DecodeAFInfo2(af []uint16, cfg AFInfo2DecodeConfig) AFInfo {
	n := len(af)
	dst := AFInfo{
		Source:           cfg.Source,
		AFAreaMode:       AFAreaMode(u16At(af, n, 1)),
		NumAFPoints:      u16At(af, n, 2),
		ValidAFPoints:    u16At(af, n, 3),
		CanonImageWidth:  u16At(af, n, 4),
		CanonImageHeight: u16At(af, n, 5),
		AFImageWidth:     u16At(af, n, 6),
		AFImageHeight:    u16At(af, n, 7),
	}

	num := int(dst.NumAFPoints)
	if num <= 0 {
		return dst
	}

	widthStart := 8
	heightStart := widthStart + num
	xStart := heightStart + num
	yStart := xStart + num
	bitsStart := yStart + num
	maskWordCount := bitWordCount(num)
	selectedStart := bitsStart + maskWordCount

	widthLen := rangeLen(n, widthStart, num)
	heightLen := rangeLen(n, heightStart, num)
	xLen := rangeLen(n, xStart, num)
	yLen := rangeLen(n, yStart, num)
	areaCount := min(yLen, min(xLen, min(heightLen, widthLen)))

	var pts []AFPoint
	if cfg.DecodeCoords {
		if areaCount > 0 {
			if cfg.DecodePoints {
				combined := make([]AFPoint, areaCount*2)
				dst.AFArea = combined[:areaCount]
				pts = combined[areaCount:]
			} else {
				dst.AFArea = make([]AFPoint, areaCount)
			}
			for i := 0; i < len(dst.AFArea); i++ {
				dst.AFArea[i] = NewAFPoint(
					int16(af[widthStart+i]),
					int16(af[heightStart+i]),
					int16(af[xStart+i]),
					int16(af[yStart+i]),
				)
			}
		}
	}

	wantSelected := cfg.EOS && cfg.DecodeSelected
	if cfg.DecodeInFocus || wantSelected {
		totalBits := 0
		if cfg.DecodeInFocus {
			totalBits += countBitWordsRange(af, n, bitsStart, maskWordCount)
		}
		if wantSelected {
			totalBits += countBitWordsRange(af, n, selectedStart, maskWordCount)
		}
		combinedBits := make([]int, 0, totalBits)

		if cfg.DecodeInFocus {
			startIdx := len(combinedBits)
			combinedBits = appendBitWordsRange(combinedBits, af, n, bitsStart, maskWordCount)
			dst.AFPointsInFocusBits = combinedBits[startIdx:]
		}
		if wantSelected {
			startIdx := len(combinedBits)
			combinedBits = appendBitWordsRange(combinedBits, af, n, selectedStart, maskWordCount)
			dst.AFPointsSelectedBits = combinedBits[startIdx:]
		}
	}

	if !(cfg.EOS && cfg.DecodeSelected) && !cfg.AFInfo3 {
		// Non-EOS AFInfo2 uses an unknown field of maskWordCount+1 at seq 13.
		dst.PrimaryAFPoint = u16At(af, n, selectedStart+maskWordCount+1)
	}

	if !cfg.DecodePoints || areaCount <= 0 {
		return dst
	}
	if pts == nil {
		pts = make([]AFPoint, areaCount)
	}
	xAdjust := int16(dst.CanonImageWidth / 2)
	yAdjust := int16(dst.CanonImageHeight / 2)
	for i := 0; i < areaCount; i++ {
		var w, h, x, y int16
		if cfg.DecodeCoords {
			area := dst.AFArea[i]
			w, h, x, y = area[0], area[1], area[2], area[3]
		} else {
			w = int16(af[widthStart+i])
			h = int16(af[heightStart+i])
			x = int16(af[xStart+i])
			y = int16(af[yStart+i])
		}
		x += xAdjust - (w / 2)
		y += yAdjust - (h / 2)
		pts[i] = NewAFPoint(w, h, x, y)
	}
	dst.AFPoints = pts
	return dst
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

func u16At(vals []uint16, n, idx int) uint16 {
	if idx < 0 || idx >= n {
		return 0
	}
	return vals[idx]
}

func rangeLen(n, start, count int) int {
	if count <= 0 || start < 0 || start >= n {
		return 0
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return 0
	}
	return end - start
}

// legacyAFInfoPrimary mirrors Canon.pm sequence handling for AFInfo:
// sequence 11 is either PrimaryAFPoint or an 8-word unknown block, and
// sequence 12 is always PrimaryAFPoint when enough payload remains.
func legacyAFInfoPrimary(vals []uint16, n, seq11Start, afInfoCount int) uint16 {
	if afInfoCount == 36 {
		return u16At(vals, n, seq11Start+8)
	}
	if seq11Start+1 < n {
		return vals[seq11Start+1]
	}
	return u16At(vals, n, seq11Start)
}

func decodeBitWordsRange(vals []uint16, n, start, count int) []int {
	capHint := countBitWordsRange(vals, n, start, count)
	if capHint == 0 {
		return nil
	}
	out := make([]int, 0, capHint)
	return appendBitWordsRange(out, vals, n, start, count)
}

func countBitWordsRange(vals []uint16, n, start, count int) int {
	if count <= 0 || start < 0 || start >= n {
		return 0
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return 0
	}
	total := 0
	for i := start; i < end; i++ {
		total += bits.OnesCount16(vals[i])
	}
	return total
}

func appendBitWordsRange(dst []int, vals []uint16, n, start, count int) []int {
	if count <= 0 || start < 0 || start >= n {
		return dst
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return dst
	}

	base := 0
	for i := start; i < end; i++ {
		word := vals[i]
		for word != 0 {
			bit := bits.TrailingZeros16(word)
			dst = append(dst, base+bit)
			word &= word - 1
		}
		base += 16
	}
	return dst
}
