package canon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// LensWithTC returns lens name with teleconverter if applicable.
// Ported from ExifTool Canon.pm LensWithTC().
func LensWithTC(lens string, shortFocal float64) string {
	if strings.HasSuffix(lens, "x") {
		return lens
	}
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	m := re.FindStringSubmatch(lens)
	if len(m) < 2 {
		return lens
	}
	sf, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return lens
	}
	for _, tc := range []float64{1, 1.4, 2, 2.8} {
		if math.Abs(shortFocal-sf*tc) > 0.9 {
			continue
		}
		if tc > 1 {
			lens += " + " + strconv.FormatFloat(tc, 'f', -1, 64) + "x"
		}
		break
	}
	return lens
}

// CalcSensorDiag calculates sensor diagonal in mm from X/Y focal-plane resolution rationals.
// The input strings should be rational values in "num/den" format.
// Ported from ExifTool Canon.pm CalcSensorDiag().
func CalcSensorDiag(xResRational, yResRational string) (float64, bool) {
	xNum, xDen, ok := parseRationalParts(xResRational)
	if !ok {
		return 0, false
	}
	yNum, yDen, ok := parseRationalParts(yResRational)
	if !ok {
		return 0, false
	}

	if xNum%1000 != 0 || yNum%1000 != 0 {
		return 0, false
	}
	if xNum < 640000 || yNum < 480000 || xNum >= 10000000 || yNum >= 10000000 {
		return 0, false
	}
	if xDen < 61 || xDen >= 1500 || yDen < 61 || yDen >= 1000 {
		return 0, false
	}
	if xDen == yDen {
		return 0, false
	}
	return math.Sqrt(float64(xDen*xDen+yDen*yDen)) * 0.0254, true
}

func parseRationalParts(s string) (num, den int64, ok bool) {
	parts := strings.FieldsFunc(strings.TrimSpace(s), func(r rune) bool {
		return r == '/' || r == ' '
	})
	if len(parts) < 2 {
		return 0, 0, false
	}
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	d, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || d == 0 {
		return 0, 0, false
	}
	return n, d, true
}

// PrintLensID attempts to identify a specific lens if multiple lenses share the same LensType.
// The printConv map may contain keys "123", "123.1", "123.2", ... for alternate matches.
// Ported from ExifTool Canon.pm PrintLensID(), with simplified lens-model matching.
func PrintLensID(printConv map[string]string, lensType string, shortFocal, longFocal, maxAperture float64, lensModel string, userLens map[string]bool) string {
	lens := printConv[lensType]
	if lens != "" {
		if _, hasAlt := printConv[lensType+".1"]; !hasAlt {
			return LensWithTC(lens, shortFocal)
		}
		if idx := strings.Index(lens, " or "); idx >= 0 {
			lens = lens[:idx]
		}

		lenses := []string{lens}
		for i := 1; ; i++ {
			l := printConv[fmt.Sprintf("%s.%d", lensType, i)]
			if l == "" {
				break
			}
			lenses = append(lenses, l)
		}

		var user, maybe, likely, matches []string
		for _, l := range lenses {
			if userLens != nil && userLens[l] {
				user = append(user, l)
			}
		}

		for _, tc := range []float64{1, 1.4, 2, 2.8} {
			for _, l := range lenses {
				sf, lf, sa, la, ok := parseLensSpec(l)
				if !ok {
					continue
				}
				if ltc, ok := parseLensTCSuffix(l); ok {
					sf *= ltc
					lf *= ltc
					sa *= ltc
					la *= ltc
				}
				if math.Abs(shortFocal-sf*tc) > 0.9 {
					continue
				}
				tclens := l
				if tc > 1 {
					tclens += " + " + strconv.FormatFloat(tc, 'f', -1, 64) + "x"
				}
				maybe = append(maybe, tclens)
				if math.Abs(longFocal-lf*tc) > 0.9 {
					continue
				}
				likely = append(likely, tclens)
				if maxAperture > 0 {
					if maxAperture < sa*tc-0.18 || maxAperture > la*tc+0.18 {
						continue
					}
				}
				matches = append(matches, tclens)
			}
			if len(maybe) > 0 {
				break
			}
		}

		if len(user) > 0 {
			if len(user) == 1 {
				return LensWithTC(user[0], shortFocal)
			}
			var good []string
			for _, group := range [][]string{matches, likely, maybe} {
				for _, candidate := range group {
					if userLens[candidate] {
						good = append(good, candidate)
						continue
					}
					base := stripTCSuffix(candidate)
					if base != candidate && userLens[base] {
						good = append(good, candidate)
					}
				}
				if len(good) > 0 {
					return strings.Join(good, " or ")
				}
			}
			return LensWithTC(user[0], shortFocal)
		}

		if len(matches) > 1 && strings.Contains(lensModel, "| ") {
			if m := regexp.MustCompile(`\| [ACS]`).FindString(lensModel); m != "" {
				var best []string
				for _, v := range matches {
					if strings.Contains(v, m) {
						best = append(best, v)
					}
				}
				if len(best) > 0 {
					matches = best
				}
			}
		}
		if len(matches) == 0 {
			matches = likely
		}
		if len(matches) == 0 {
			matches = maybe
		}
		matches = matchLensModelApprox(matches, lensModel)
		if len(matches) > 0 {
			return strings.Join(matches, " or ")
		}
	} else if hasDigit(lensModel) {
		return lensModel
	}

	var suffix string
	if shortFocal > 0 {
		suffix += fmt.Sprintf(" %d", int(math.Round(shortFocal)))
		if longFocal > 0 && math.Abs(longFocal-shortFocal) > 0.01 {
			suffix += fmt.Sprintf("-%d", int(math.Round(longFocal)))
		}
		suffix += "mm"
	}
	if lensType == "-1" || lensType == "65535" {
		return "Unknown" + suffix
	}
	return fmt.Sprintf("Unknown (%s)%s", lensType, suffix)
}

func parseLensSpec(lens string) (sf, lf, sa, la float64, ok bool) {
	re := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)(?:-(\d+(?:\.\d+)?))?mm.*?f/?(\d+(?:\.\d+)?)(?:-(\d+(?:\.\d+)?))?`)
	m := re.FindStringSubmatch(lens)
	if len(m) < 4 {
		return 0, 0, 0, 0, false
	}
	var err error
	sf, err = strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, 0, 0, 0, false
	}
	lf = sf
	if m[2] != "" {
		lf, err = strconv.ParseFloat(m[2], 64)
		if err != nil {
			return 0, 0, 0, 0, false
		}
	}
	sa, err = strconv.ParseFloat(m[3], 64)
	if err != nil {
		return 0, 0, 0, 0, false
	}
	la = sa
	if len(m) > 4 && m[4] != "" {
		la, err = strconv.ParseFloat(m[4], 64)
		if err != nil {
			return 0, 0, 0, 0, false
		}
	}
	return sf, lf, sa, la, true
}

func parseLensTCSuffix(s string) (float64, bool) {
	re := regexp.MustCompile(` \+ (\d+(?:\.\d+)?)x$`)
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func stripTCSuffix(s string) string {
	re := regexp.MustCompile(` \+ \d+(?:\.\d+)?x$`)
	return re.ReplaceAllString(s, "")
}

func matchLensModelApprox(matches []string, lensModel string) []string {
	if lensModel == "" || len(matches) < 2 {
		return matches
	}
	lm := strings.ToLower(strings.TrimSpace(lensModel))
	var filtered []string
	for _, m := range matches {
		if strings.Contains(strings.ToLower(m), lm) {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) > 0 {
		return filtered
	}
	return matches
}

func hasDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

// SwapWords swaps 16-bit words in each 32-bit integer.
// Ported from ExifTool Canon.pm SwapWords().
func SwapWords(vals []uint32) []uint32 {
	out := make([]uint32, len(vals))
	for i, v := range vals {
		out[i] = (v>>16 | v<<16) & 0xffffffff
	}
	return out
}

// Validate checks whether the first 16-bit value at offset matches one of valid values.
// Ported from ExifTool Canon.pm Validate().
func Validate(data []byte, offset int, vals ...uint16) bool {
	return ValidateWithByteOrder(data, offset, binary.LittleEndian, vals...)
}

// ValidateWithByteOrder checks first 16-bit value using explicit byte order.
func ValidateWithByteOrder(data []byte, offset int, order binary.ByteOrder, vals ...uint16) bool {
	if offset < 0 || offset+2 > len(data) {
		return false
	}
	got := order.Uint16(data[offset : offset+2])
	for _, v := range vals {
		if got == v {
			return true
		}
	}
	return false
}

// ValidateAFInfo validates Canon AFInfo binary payload.
// Ported from ExifTool Canon.pm ValidateAFInfo().
func ValidateAFInfo(data []byte, offset, size int) bool {
	if size < 24 || offset < 0 || offset+size > len(data) {
		return false
	}
	get16u := func(off int) uint16 {
		if off < 0 || off+2 > len(data) {
			return 0
		}
		return binary.LittleEndian.Uint16(data[off : off+2])
	}
	af := get16u(offset)
	switch af {
	case 1, 5, 7, 9, 15, 45, 53:
	default:
		return false
	}
	w1, h1 := float64(get16u(offset+4)), float64(get16u(offset+6))
	if h1 == 0 || w1 == 0 {
		return false
	}
	f1 := w1 / h1
	if math.Abs(f1-1.33) < 0.01 || math.Abs(f1-1.67) < 0.01 {
		return true
	}
	if math.Abs(f1-0.75) < 0.01 || math.Abs(f1-0.60) < 0.01 {
		return true
	}
	w2, h2 := float64(get16u(offset+8)), float64(get16u(offset+10))
	if h2 == 0 || w2 == 0 {
		return false
	}
	if w1 == h1 {
		return false
	}
	f2 := w2 / h2
	if math.Abs(1-f1/f2) < 0.01 {
		return true
	}
	if math.Abs(1-f1*f2) < 0.01 {
		return true
	}
	return false
}

// ReadODD reads Canon OriginalDecisionData block at file offset.
// Ported from ExifTool Canon.pm ReadODD() with file-agnostic ReaderAt API.
func ReadODD(r ReaderAt, offset int64) ([]byte, error) {
	if r == nil || offset <= 0 {
		return nil, nil
	}
	head := make([]byte, 8)
	if _, err := r.ReadAt(head, offset); err != nil {
		return nil, err
	}
	if !bytes.Equal(head[:4], []byte{0xff, 0xff, 0xff, 0xff}) {
		return nil, errors.New("invalid original decision data header")
	}

	order, version, ok := detectODDVersion(head[4:8])
	if !ok {
		return nil, fmt.Errorf("unsupported original decision data version")
	}
	buff := append([]byte{}, head...)

	switch version {
	case 1, 2:
		tail := make([]byte, 24)
		if _, err := r.ReadAt(tail, offset+8); err != nil {
			return nil, err
		}
		buff = append(buff, tail...)
		count := order.Uint32(tail[20:24])
		if count == 0 || count >= 20 {
			return nil, fmt.Errorf("invalid original decision data record count: %d", count)
		}
		records := make([]byte, count*32)
		if _, err := r.ReadAt(records, offset+8+24); err != nil {
			return nil, err
		}
		buff = append(buff, records...)
	case 3:
		pos := offset + 8
		for i := 0; i < 3; i++ {
			lenWord := make([]byte, 4)
			if _, err := r.ReadAt(lenWord, pos); err != nil {
				return nil, err
			}
			pos += 4
			buff = append(buff, lenWord...)
			l := order.Uint32(lenWord)
			if i == 2 && l >= 4 {
				l -= 4
			}
			if l > 0x10000 {
				return nil, fmt.Errorf("invalid original decision data segment length: %d", l)
			}
			seg := make([]byte, l)
			if _, err := r.ReadAt(seg, pos); err != nil {
				return nil, err
			}
			pos += int64(l)
			buff = append(buff, seg...)
		}
	default:
		return nil, fmt.Errorf("unsupported original decision data version %d", version)
	}
	return buff, nil
}

func detectODDVersion(versionBytes []byte) (binary.ByteOrder, uint32, bool) {
	le := binary.LittleEndian.Uint32(versionBytes)
	switch le {
	case 1, 2, 3:
		return binary.LittleEndian, le, true
	}
	be := binary.BigEndian.Uint32(versionBytes)
	switch be {
	case 1, 2, 3:
		return binary.BigEndian, be, true
	}
	return nil, 0, false
}

// ReaderAt abstracts *os.File for ReadODD.
type ReaderAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
}

var (
	cameraISOLookup = map[uint16]string{
		0:  "n/a",
		14: "Auto High",
		15: "Auto",
		16: "50",
		17: "100",
		18: "200",
		19: "400",
		20: "800",
	}
	cameraISOReverseLookup = map[string]uint16{
		"n/a":       0,
		"auto high": 14,
		"auto":      15,
		"50":        16,
		"100":       17,
		"200":       18,
		"400":       19,
		"800":       20,
	}
)

// CameraISO converts Canon CameraISO value to text.
// Ported from ExifTool Canon.pm CameraISO().
func CameraISO(val uint16) string {
	if val == 0x7fff {
		return ""
	}
	if val&0x4000 != 0 {
		return strconv.Itoa(int(val & 0x3fff))
	}
	if s, ok := cameraISOLookup[val]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", val)
}

// CameraISOInv converts CameraISO text back to Canon CameraISO code.
func CameraISOInv(val string) (uint16, bool) {
	v := strings.ToLower(strings.TrimSpace(val))
	if code, ok := cameraISOReverseLookup[v]; ok {
		return code, true
	}
	if n, err := strconv.Atoi(v); err == nil {
		return (uint16(n) & 0x3fff) | 0x4000, true
	}
	return 0, false
}

// PrintFocalRange prints a short/long focal range in mm.
// Ported from ExifTool Canon.pm PrintFocalRange().
func PrintFocalRange(short, long, scale float64) string {
	if scale == 0 {
		scale = 1
	}
	if math.Abs(short-long) < 0.0001 {
		return fmt.Sprintf("%.1f mm", short*scale)
	}
	return fmt.Sprintf("%.1f - %.1f mm", short*scale, long*scale)
}

// PrintAFPoints1D decodes 1D AF points packed value.
// Ported from ExifTool Canon.pm PrintAFPoints1D().
func PrintAFPoints1D(val []byte) string {
	if len(val) != 8 {
		return "Unknown"
	}
	focusPts := []int{
		0, 0,
		0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e, 0x10, 0, 0,
		0x21, 0x23, 0x25, 0x27, 0x29, 0x2b, 0x2d, 0x2f, 0x31, 0x33,
		0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4d, 0x50, 0x52, 0x54,
		0x61, 0x63, 0x65, 0x67, 0x69, 0x6b, 0x6d, 0x6f, 0x71, 0x73, 0, 0,
		0x84, 0x86, 0x88, 0x8a, 0x8c, 0x8e, 0x90, 0, 0, 0, 0, 0,
	}
	rows := []rune("  AAAAAAA  BBBBBBBBBBCCCCCCCCCCCDDDDDDDDDD  EEEEEEE     ")
	bits := bitStringLSB(val[1:])

	focus := int(val[0])
	var focusing string
	var points []string
	var lastRow rune
	col := 0

	for i, fp := range focusPts {
		row := rows[i]
		if row == lastRow {
			col++
		} else {
			col = 1
		}
		lastRow = row
		if focus == fp {
			focusing = fmt.Sprintf("%c%d", row, col)
		}
		if i < len(bits) && bits[i] {
			points = append(points, fmt.Sprintf("%c%d", row, col))
		}
	}
	if focusing == "" {
		if focus == 0xff {
			focusing = "Auto"
		} else {
			focusing = fmt.Sprintf("Unknown (0x%.2x)", focus)
		}
	}
	return focusing + " (" + strings.Join(points, ",") + ")"
}

func bitStringLSB(data []byte) []bool {
	out := make([]bool, 0, len(data)*8)
	for _, b := range data {
		for i := 0; i < 8; i++ {
			out = append(out, b&(1<<i) != 0)
		}
	}
	return out
}

// SerialTag defines one tag in a serial binary stream.
type SerialTag struct {
	Index   int
	Name    string
	Format  string
	Count   int
	Unknown bool
}

// SerialValue is one parsed serial tag value.
type SerialValue struct {
	Tag    SerialTag
	Value  interface{}
	Offset int
	Size   int
}

// ProcessSerialData parses a serial stream of binary tag data.
// Ported conceptually from ExifTool Canon.pm ProcessSerialData().
func ProcessSerialData(data []byte, tags []SerialTag, defaultFormat string) ([]SerialValue, error) {
	if defaultFormat == "" {
		defaultFormat = "int8u"
	}
	pos := 0
	values := make([]SerialValue, 0, len(tags))
	for _, tag := range tags {
		format := tag.Format
		if format == "" {
			format = defaultFormat
		}
		count := tag.Count
		if count <= 0 {
			count = 1
		}
		sizePer, err := formatSize(format)
		if err != nil {
			return nil, err
		}
		total := sizePer * count
		if pos+total > len(data) {
			break
		}
		val, err := readFormatValue(data[pos:pos+total], format, count)
		if err != nil {
			return nil, err
		}
		values = append(values, SerialValue{
			Tag:    tag,
			Value:  val,
			Offset: pos,
			Size:   total,
		})
		pos += total
	}
	return values, nil
}

func formatSize(format string) (int, error) {
	switch format {
	case "int8u", "int8s", "string":
		return 1, nil
	case "int16u", "int16s":
		return 2, nil
	case "int32u", "int32s":
		return 4, nil
	default:
		return 0, fmt.Errorf("unsupported serial format %q", format)
	}
}

func readFormatValue(data []byte, format string, count int) (interface{}, error) {
	switch format {
	case "string":
		return string(data), nil
	case "int8u":
		if count == 1 {
			return data[0], nil
		}
		out := make([]uint8, count)
		copy(out, data)
		return out, nil
	case "int8s":
		if count == 1 {
			return int8(data[0]), nil
		}
		out := make([]int8, count)
		for i := range out {
			out[i] = int8(data[i])
		}
		return out, nil
	case "int16u":
		if count == 1 {
			return binary.LittleEndian.Uint16(data[:2]), nil
		}
		out := make([]uint16, count)
		for i := range out {
			out[i] = binary.LittleEndian.Uint16(data[i*2 : i*2+2])
		}
		return out, nil
	case "int16s":
		if count == 1 {
			return int16(binary.LittleEndian.Uint16(data[:2])), nil
		}
		out := make([]int16, count)
		for i := range out {
			out[i] = int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
		}
		return out, nil
	case "int32u":
		if count == 1 {
			return binary.LittleEndian.Uint32(data[:4]), nil
		}
		out := make([]uint32, count)
		for i := range out {
			out[i] = binary.LittleEndian.Uint32(data[i*4 : i*4+4])
		}
		return out, nil
	case "int32s":
		if count == 1 {
			return int32(binary.LittleEndian.Uint32(data[:4])), nil
		}
		out := make([]int32, count)
		for i := range out {
			out[i] = int32(binary.LittleEndian.Uint32(data[i*4 : i*4+4]))
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported serial format %q", format)
	}
}

// ExifInfoRecord is one ExifInfo block.
type ExifInfoRecord struct {
	Tag     uint32
	Payload []byte
	Offset  int
}

// ProcessExifInfo parses CTMD-style ExifInfo records.
// Ported conceptually from ExifTool Canon.pm ProcessExifInfo().
func ProcessExifInfo(data []byte, knownTags map[uint32]struct{}) []ExifInfoRecord {
	var out []ExifInfoRecord
	for pos := 0; pos+8 <= len(data); {
		l := int(binary.LittleEndian.Uint32(data[pos : pos+4]))
		tag := binary.LittleEndian.Uint32(data[pos+4 : pos+8])
		if l < 8 || pos+l > len(data) {
			break
		}
		if knownTags != nil {
			if _, ok := knownTags[tag]; !ok {
				break
			}
		}
		payload := append([]byte{}, data[pos+8:pos+l]...)
		out = append(out, ExifInfoRecord{Tag: tag, Payload: payload, Offset: pos + 8})
		pos += l
	}
	return out
}

// CTMDRecord is one Canon Timed MetaData record.
type CTMDRecord struct {
	Type    uint16
	Header  [6]byte
	Payload []byte
	Offset  int
}

// ProcessCTMD parses Canon Timed MetaData records.
// Ported conceptually from ExifTool Canon.pm ProcessCTMD().
func ProcessCTMD(data []byte) ([]CTMDRecord, error) {
	var out []CTMDRecord
	pos := 0
	for pos+6 < len(data) {
		size := int(binary.LittleEndian.Uint32(data[pos : pos+4]))
		typ := binary.LittleEndian.Uint16(data[pos+4 : pos+6])
		if size < 12 {
			return nil, errors.New("short CTMD record")
		}
		if pos+size > len(data) {
			return nil, errors.New("truncated CTMD record")
		}
		var hdr [6]byte
		copy(hdr[:], data[pos+6:pos+12])
		payload := append([]byte{}, data[pos+12:pos+size]...)
		out = append(out, CTMDRecord{
			Type:    typ,
			Header:  hdr,
			Payload: payload,
			Offset:  pos,
		})
		pos += size
	}
	if pos != len(data) {
		return nil, errors.New("error parsing Canon CTMD data")
	}
	return out, nil
}

// FilterParam is one creative-filter parameter.
type FilterParam struct {
	ID     uint32
	Values []int32
}

// CreativeFilter is one creative-filter record.
type CreativeFilter struct {
	Number uint32
	Params []FilterParam
}

// ProcessFilters parses Canon creative filter structures.
// Ported conceptually from ExifTool Canon.pm ProcessFilters().
func ProcessFilters(data []byte) ([]CreativeFilter, error) {
	if len(data) < 8 {
		return nil, errors.New("creative filter data too short")
	}
	numFilters := int(binary.LittleEndian.Uint32(data[4:8]))
	pos := 8
	end := len(data)
	out := make([]CreativeFilter, 0, numFilters)

	for i := 0; i < numFilters; i++ {
		if pos+12 > end {
			return nil, fmt.Errorf("truncated data for filter %d", i)
		}
		fnum := binary.LittleEndian.Uint32(data[pos : pos+4])
		size := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		nparm := int(binary.LittleEndian.Uint32(data[pos+8 : pos+12]))
		next := pos + 4 + size
		if next > end {
			return nil, fmt.Errorf("invalid size (%d) for filter %d", size, i)
		}
		pos += 12

		filter := CreativeFilter{Number: fnum}
		for j := 0; j < nparm; j++ {
			if pos+12 > end {
				return nil, fmt.Errorf("truncated data for filter %d param %d", i, j)
			}
			tag := binary.LittleEndian.Uint32(data[pos : pos+4])
			count := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
			pos += 8
			if pos+4*count > end {
				return nil, fmt.Errorf("truncated value for filter %d param %d", i, j)
			}
			vals := make([]int32, count)
			for k := 0; k < count; k++ {
				vals[k] = int32(binary.LittleEndian.Uint32(data[pos+k*4 : pos+k*4+4]))
			}
			pos += 4 * count
			filter.Params = append(filter.Params, FilterParam{ID: tag, Values: vals})
		}
		pos = next
		out = append(out, filter)
	}
	return out, nil
}

// ProcessCMT3 extracts Canon static maker notes payload from CMT3 data.
// Ported conceptually from ExifTool Canon.pm ProcessCMT3().
func ProcessCMT3(data []byte) []byte {
	if len(data) <= 8 {
		return nil
	}
	clean := stripCanonTrailer(data)
	out := make([]byte, 0, len(clean))
	out = append(out, clean[8:]...)
	out = append(out, clean[:8]...)
	return out
}

func stripCanonTrailer(data []byte) []byte {
	if len(data) < 12 {
		return data
	}
	for z := 4; z <= 10; z++ {
		if len(data) < 4+z {
			continue
		}
		tail := data[len(data)-(4+z):]
		isII := bytes.Equal(tail[:4], []byte{'I', 'I', 0x2a, 0x00})
		isMM := bytes.Equal(tail[:4], []byte{'M', 'M', 0x00, 0x2a})
		if !isII && !isMM {
			continue
		}
		allZero := true
		for _, b := range tail[4:] {
			if b != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return data[:len(data)-(4+z)]
		}
	}
	return data
}

// WriteCanon appends the Canon TIFF-style maker-note footer when requested.
// Ported from ExifTool Canon.pm WriteCanon().
func WriteCanon(dirData []byte, addFooter bool, littleEndian bool) []byte {
	if !addFooter || len(dirData) == 0 {
		return dirData
	}
	footer := []byte{'I', 'I', 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00}
	if !littleEndian {
		footer = []byte{'M', 'M', 0x00, 0x2a, 0x00, 0x00, 0x00, 0x00}
	}
	out := make([]byte, 0, len(dirData)+len(footer))
	out = append(out, dirData...)
	out = append(out, footer...)
	return out
}
