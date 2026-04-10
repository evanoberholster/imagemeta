package exif

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/bits"
	"strings"
	"sync"
	"time"
	"unicode/utf16"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	metacanon "github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

var (
	timezoneCache   = map[int32]*time.Location{}
	timezoneCacheMu sync.RWMutex
)

// parseStrUint parses the requested value from EXIF metadata.
func parseStrUint(buf []byte) (u uint) {
	for i := 0; i < len(buf); i++ {
		if buf[i] >= '0' && buf[i] <= '9' {
			u *= 10
			u += uint(buf[i] - '0')
		}
	}
	return
}

// trimNULBuffer trims input bytes into the expected EXIF representation.
func trimNULBuffer(buf []byte) []byte {
	for i := len(buf) - 1; i >= 0; i-- {
		if buf[i] == 0 || buf[i] == ' ' || buf[i] == '\n' {
			continue
		}
		return buf[:i+1]
	}
	return nil
}

// getLocation returns a cached fixed-zone location for an EXIF offset string.
func getLocation(offset int32, label []byte) *time.Location {
	timezoneCacheMu.RLock()
	if z, ok := timezoneCache[offset]; ok {
		timezoneCacheMu.RUnlock()
		return z
	}
	timezoneCacheMu.RUnlock()

	timezoneCacheMu.Lock()
	loc := time.FixedZone(string(label), int(offset))
	timezoneCache[offset] = loc
	timezoneCacheMu.Unlock()
	return loc
}

func exifASCIIText(buf []byte) string {
	trimmed := bytes.TrimSpace(trimNULBuffer(buf))
	if len(trimmed) == 0 {
		return ""
	}
	return string(trimmed)
}

func exifUTF16Text(buf []byte, bo binary.ByteOrder) string {
	if len(buf) < 2 {
		return ""
	}
	if len(buf)&1 != 0 {
		buf = buf[:len(buf)-1]
	}
	if len(buf) == 0 {
		return ""
	}

	switch {
	case len(buf) >= 2 && buf[0] == 0xfe && buf[1] == 0xff:
		bo = binary.BigEndian
		buf = buf[2:]
	case len(buf) >= 2 && buf[0] == 0xff && buf[1] == 0xfe:
		bo = binary.LittleEndian
		buf = buf[2:]
	}
	if len(buf) < 2 {
		return ""
	}
	if len(buf)&1 != 0 {
		buf = buf[:len(buf)-1]
	}

	u16 := make([]uint16, 0, len(buf)/2)
	for i := 0; i+1 < len(buf); i += 2 {
		v := bo.Uint16(buf[i : i+2])
		if v == 0 {
			break
		}
		u16 = append(u16, v)
	}
	if len(u16) == 0 {
		return ""
	}
	return strings.TrimSpace(string(utf16.Decode(u16)))
}

// rationalDuration converts a positive EXIF rational into a duration scaled by unit.
func rationalDuration(num uint32, den uint32, unit time.Duration) time.Duration {
	if num == 0 || den == 0 || unit <= 0 {
		return 0
	}
	unitU := uint64(unit)
	numU := uint64(num)
	denU := uint64(den)

	// Guard whole-part multiplication to avoid uint64 wrap.
	maxWhole := uint64(math.MaxInt64) / unitU
	whole := numU / denU
	if whole > maxWhole {
		return time.Duration(math.MaxInt64)
	}

	wholeUnits := whole * unitU

	// Compute (remainder * unit) / den as a 128-bit product to avoid overflow.
	remainder := numU % denU
	hi, lo := bits.Mul64(remainder, unitU)
	var frac uint64
	if hi >= denU {
		return time.Duration(math.MaxInt64)
	}
	frac, _ = bits.Div64(hi, lo, denU)
	if wholeUnits > uint64(math.MaxInt64)-frac {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(wholeUnits + frac)
}

// readTagBytes reads data from the underlying stream or parser buffers.
func (r *Reader) readTagBytes(t tag.Entry, max uint32) (buf []byte, truncated bool, err error) {
	if err = r.seekToTag(t); err != nil {
		return nil, false, err
	}

	size := t.Size()
	if size == 0 {
		return nil, false, nil
	}
	if max > 0 && size > max {
		size = max
		truncated = true
	}
	if size > uint32(len(r.state.buf)) {
		size = uint32(len(r.state.buf))
		truncated = true
	}

	buf, err = r.fastRead(int(size))
	if err != nil {
		return nil, false, err
	}

	remaining := int(t.Size() - size)
	if remaining > 0 {
		truncated = true
		if discardErr := r.discard(remaining); discardErr != nil {
			return nil, true, discardErr
		}
	}
	return buf, truncated, nil
}

// seekToTag moves reader state to the location required for parsing.
func (r *Reader) seekToTag(t tag.Entry) error {
	return r.discard(int(t.ValueOffset) - int(r.po))
}

// fastRead reads a bounded byte slice using the optimized buffered path.
func (r *Reader) fastRead(n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}
	if n < 0 || n > len(r.state.buf) {
		return nil, imagetype.ErrDataLength
	}
	if r.exifLength > 0 && int(r.po)+n > int(r.exifLength) {
		return nil, imagetype.ErrDataLength
	}
	buf, err := r.reader.Peek(n)
	if err != nil {
		return nil, err
	}
	readCount, err := r.reader.Discard(len(buf))
	r.po += uint32(readCount)
	return buf, err
}

// fastRead2 reads a bounded byte slice using the optimized buffered path.
func (r *Reader) fastRead2(buf []byte) (int, error) {
	l := len(buf)
	if l == 0 {
		return 0, nil
	}
	if l > len(r.state.buf) {
		return 0, imagetype.ErrDataLength
	}
	if r.exifLength > 0 && int(r.po)+l > int(r.exifLength) {
		return 0, imagetype.ErrDataLength
	}
	readCount, err := r.reader.Read(buf)
	r.po += uint32(readCount)
	if err != nil {
		return 0, err
	}
	buf = buf[:readCount]
	return readCount, nil
}

// discard advances the reader by discarding the requested number of bytes.
func (r *Reader) discard(n int) error {
	if n <= 0 {
		return nil
	}
	if r.exifLength > 0 && int(r.exifLength) < n+int(r.po) {
		n = int(r.exifLength) - int(r.po)
	}
	if n <= 0 {
		return nil
	}
	discarded, err := r.reader.Discard(n)
	r.po += uint32(discarded)
	return err
}

// readUint16 reads data from the underlying stream or parser buffers.
func (r *Reader) readUint16(directory tag.Directory) (uint16, error) {
	buf, err := r.fastRead(2)
	if err != nil || len(buf) < 2 {
		return 0, err
	}
	return directory.ByteOrder.Uint16(buf), nil
}

// readUint32 reads data from the underlying stream or parser buffers.
func (r *Reader) readUint32(directory tag.Directory) (uint32, error) {
	buf, err := r.fastRead(4)
	if err != nil || len(buf) < 4 {
		return 0, err
	}
	return directory.ByteOrder.Uint32(buf), nil
}

const canonAFWordsMax = 8192

func canonAFWordsBuffer(stack []uint16, unitCount uint32) ([]uint16, bool) {
	if unitCount == 0 {
		return stack[:0], false
	}
	wordCount := int(unitCount)
	truncated := false
	if unitCount > canonAFWordsMax {
		wordCount = canonAFWordsMax
		truncated = true
	}
	if wordCount <= len(stack) {
		return stack[:wordCount], truncated
	}
	return make([]uint16, wordCount), truncated
}

func canonAFInfoSource(id tag.ID) metacanon.AFInfoSource {
	switch metacanon.MakerNoteTag(id) {
	case metacanon.CanonAFInfo:
		return metacanon.AFInfoSourceAFInfo
	case metacanon.CanonAFInfo2:
		return metacanon.AFInfoSourceAFInfo2
	case metacanon.AFInfo3:
		return metacanon.AFInfoSourceAFInfo3
	default:
		return metacanon.AFInfoSourceUnknown
	}
}

func canonShouldReplaceAFInfo(current, candidate metacanon.AFInfo) bool {
	curHas := canonAFInfoHasData(current)
	candHas := canonAFInfoHasData(candidate)
	switch {
	case candHas && !curHas:
		return true
	case !candHas && curHas:
		return false
	case !candHas && !curHas:
		return canonAFInfoSourcePriority(candidate.Source) > canonAFInfoSourcePriority(current.Source)
	}

	curScore := canonAFInfoQualityScore(current)
	candScore := canonAFInfoQualityScore(candidate)
	if candScore != curScore {
		return candScore > curScore
	}

	return canonAFInfoSourcePriority(candidate.Source) > canonAFInfoSourcePriority(current.Source)
}

func canonAFInfoHasData(v metacanon.AFInfo) bool {
	return v.NumAFPoints != 0 ||
		v.ValidAFPoints != 0 ||
		v.CanonImageWidth != 0 ||
		v.CanonImageHeight != 0 ||
		len(v.AFArea) != 0 ||
		len(v.AFPointsInFocusBits) != 0 ||
		len(v.AFPointsSelectedBits) != 0 ||
		v.PrimaryAFPoint != 0
}

func canonAFInfoQualityScore(v metacanon.AFInfo) int {
	score := int(v.NumAFPoints) + int(v.ValidAFPoints)
	score += len(v.AFArea)
	score += len(v.AFPoints)
	score += len(v.AFPointsInFocusBits)
	score += len(v.AFPointsSelectedBits)
	if v.CanonImageWidth != 0 && v.CanonImageHeight != 0 {
		score += 8
	}
	if v.AFImageWidth != 0 && v.AFImageHeight != 0 {
		score += 8
	}
	if v.AFAreaWidth != 0 || v.AFAreaHeight != 0 {
		score += 4
	}
	return score
}

func canonAFInfoSourcePriority(source metacanon.AFInfoSource) int {
	switch source {
	case metacanon.AFInfoSourceAFInfo2:
		return 3
	case metacanon.AFInfoSourceAFInfo3:
		return 2
	case metacanon.AFInfoSourceAFInfo:
		return 1
	default:
		return 0
	}
}

func canonU16At(vals []uint16, n, idx int) uint16 {
	if idx < 0 || idx >= n {
		return 0
	}
	return vals[idx]
}

func canonI16At(vals []uint16, n, idx int) int16 {
	return int16(canonU16At(vals, n, idx))
}

func canonBitWordCount(pointCount int) int {
	if pointCount <= 0 {
		return 0
	}
	return (pointCount + 15) / 16
}

func canonRangeLen(n, start, count int) int {
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

func canonDecodeUniformAFArea(vals []uint16, n, xStart, yStart, count int, w, h int16) []metacanon.AFPoint {
	pointCount := canonRangeLen(n, xStart, count)
	if yLen := canonRangeLen(n, yStart, count); yLen < pointCount {
		pointCount = yLen
	}
	if pointCount == 0 {
		return nil
	}

	areas := make([]metacanon.AFPoint, pointCount)
	for i := 0; i < pointCount; i++ {
		areas[i] = metacanon.NewAFPoint(w, h, int16(vals[xStart+i]), int16(vals[yStart+i]))
	}
	return areas
}

// canonLegacyAFInfoPrimary mirrors Canon.pm sequence handling for AFInfo:
// sequence 11 is either PrimaryAFPoint or an 8-word unknown block, and
// sequence 12 is always PrimaryAFPoint when enough payload remains.
func canonLegacyAFInfoPrimary(vals []uint16, n, seq11Start, afInfoCount int) uint16 {
	if afInfoCount == 36 {
		return canonU16At(vals, n, seq11Start+8)
	}
	if seq11Start+1 < n {
		return vals[seq11Start+1]
	}
	return canonU16At(vals, n, seq11Start)
}

func canonDecodeBitWordsRange(vals []uint16, n, start, count int) []int {
	capHint := canonCountBitWordsRange(vals, n, start, count)
	if capHint == 0 {
		return nil
	}
	out := make([]int, 0, capHint)
	return canonAppendBitWordsRange(out, vals, n, start, count)
}

func canonCountBitWordsRange(vals []uint16, n, start, count int) int {
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

func canonAppendBitWordsRange(dst []int, vals []uint16, n, start, count int) []int {
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

func (r *Reader) canonModelName() string {
	if model := r.Exif.IFD0.Model; model != "" {
		return model
	}
	if r.Exif.MakerNote.Canon != nil {
		return r.Exif.MakerNote.Canon.ImageType
	}
	return ""
}

func canonModelIsEOS(model string) bool {
	return strings.Contains(model, "EOS")
}

func (r *Reader) canonShotInfoLegacyExposureTime() bool {
	model := r.canonModelName()
	if !strings.Contains(model, "EOS 20D") && !strings.Contains(model, "EOS 350D") {
		return false
	}
	if r.Exif.MakerNote.Canon == nil {
		return false
	}
	return r.Exif.MakerNote.Canon.CanonCameraSettings.FocalUnits > 1
}

func canonShotISO(code int16) float32 {
	if code == 0 {
		return 100
	}
	return float32(100.0 * math.Exp2(float64(code-160)/32.0))
}

func canonShotActualISO(autoISO, baseISO float32) float32 {
	if autoISO <= 0 || baseISO <= 0 {
		return 0
	}
	return (autoISO * baseISO) / 100.0
}

func canonShotAperture(code int16) meta.Aperture {
	if code == 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(canonEV(code) * 0.5))
}

func canonShotExposureTime(code int16, legacy20D350D bool) meta.ExposureTime {
	if code == 0 {
		return 0
	}
	if legacy20D350D {
		return meta.ExposureTime(math.Exp2(float64(code-640) / 32.0))
	}
	return meta.ExposureTime(math.Exp2(-canonEV(code)))
}

func canonShotCameraTemperature(raw int16, model string) int16 {
	if raw == 0 || !canonModelIsEOS(model) {
		return 0
	}
	return raw - 128
}

func canonShotFlashGuideNumber(raw int16) float32 {
	if raw < 0 {
		return 0
	}
	return float32(raw) / 32.0
}

// parseCanonMaxAperture converts Canon CameraSettings MaxAperture/MinAperture
// codes to f-numbers using ExifTool's CanonEv conversion.
//
// ExifTool Canon.pm:
//
//	ValueConv => exp(CanonEv($val)*log(2)/2)
func parseCanonMaxAperture(raw uint16) meta.Aperture {
	code := int16(raw)
	if code <= 0 {
		return 0
	}
	ev := canonEV(code)
	return meta.Aperture(math.Exp2(ev * 0.5))
}

// canonEV decodes Canon's hex-based EV codes (modulo 0x20).
func canonEV(code int16) float64 {
	val := int(code)
	sign := 1.0
	if val < 0 {
		val = -val
		sign = -1
	}

	frac := val & 0x1f
	base := val - frac
	fracEV := float64(frac)

	// ExifTool CanonEv special-cases Canon 1/3 and 2/3 encodings.
	switch frac {
	case 0x0c:
		fracEV = 32.0 / 3.0
	case 0x14:
		fracEV = 64.0 / 3.0
	}
	return sign * (float64(base) + fracEV) / 32.0
}

// parseCanonDisplayAperture converts DisplayAperture (sequence 35) as ExifTool:
// RawConv => '$val ? $val : undef', ValueConv => '$val / 10'.
func parseCanonDisplayAperture(raw uint16) meta.Aperture {
	if raw == 0 {
		return 0
	}
	return meta.Aperture(float32(raw) / 10.0)
}

func canonTerminateAtNUL(s string) string {
	start := 0
	for start < len(s) {
		switch s[start] {
		case ' ', '\t', '\n', '\r':
			start++
		default:
			goto findEnd
		}
	}

findEnd:
	end := len(s)
	for i := start; i < end; i++ {
		if s[i] == 0 {
			end = i
			break
		}
	}
	for end > start {
		switch s[end-1] {
		case ' ', '\t', '\n', '\r':
			end--
		default:
			return s[start:end]
		}
	}
	return s[start:end]
}

func canonHexBytes(b []byte) string {
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
