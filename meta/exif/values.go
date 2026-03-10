package exif

import (
	"math"
	"math/bits"
	"strings"
	"sync"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

const (
	hoursToSeconds   = 60 * minutesToSeconds
	minutesToSeconds = 60
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

// parseString parses the requested value from EXIF metadata.
func (r *Reader) parseString(t tag.Entry) string {
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		return string(trimNULBuffer(r.state.buf[:t.Size()]))
	}
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		buf, _, err := r.readTagBytes(t, uint32(len(r.state.buf)))
		if err != nil {
			return ""
		}
		return string(trimNULBuffer(buf))
	}
	return ""
}

// parseStringAllowUndefined parses the requested value from EXIF metadata.
func (r *Reader) parseStringAllowUndefined(t tag.Entry) string {
	if t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul) {
		return strings.TrimSpace(r.parseString(t))
	}
	if !t.IsType(tag.TypeUndefined) {
		return ""
	}
	buf := r.parseUndefinedBytes(t, 512)
	if len(buf) == 0 {
		return ""
	}
	// exiftool text dumps render non-printable bytes as '.'. Mirror that for
	// parity checks while still preserving raw bytes separately.
	trimmed := trimNULBuffer(buf)
	if len(trimmed) == 0 {
		return ""
	}
	if len(trimmed) <= 512 {
		var out [512]byte
		for i := 0; i < len(trimmed); i++ {
			b := trimmed[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
				continue
			}
			out[i] = '.'
		}
		return strings.TrimSpace(string(out[:len(trimmed)]))
	}
	out := make([]byte, len(trimmed))
	for i := 0; i < len(trimmed); i++ {
		b := trimmed[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
			continue
		}
		out[i] = '.'
	}
	return strings.TrimSpace(string(out))
}

// parseDisplayString keeps a printable representation that is close to exiftool
// text dumps by converting non-printable bytes to '.' and preserving trailing dots.
func (r *Reader) parseDisplayString(t tag.Entry, maxBytes uint32) string {
	var buf []byte
	switch {
	case t.IsEmbedded():
		if maxBytes == 0 {
			return ""
		}
		n := t.Size()
		if n > maxBytes {
			n = maxBytes
		}
		t.EmbeddedValue(r.state.buf[:4])
		buf = r.state.buf[:n]
	case t.IsType(tag.TypeUndefined), t.IsType(tag.TypeByte), t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		if maxBytes == 0 {
			return ""
		}
		var err error
		buf, _, err = r.readTagBytes(t, maxBytes)
		if err != nil || len(buf) == 0 {
			return ""
		}
	default:
		return ""
	}

	allPrintable := true
	for i := 0; i < len(buf); i++ {
		if buf[i] < 0x20 || buf[i] > 0x7e {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return strings.TrimRight(string(buf), " \t\r\n")
	}

	if len(buf) <= 512 {
		var out [512]byte
		for i := 0; i < len(buf); i++ {
			b := buf[i]
			if b >= 0x20 && b <= 0x7e {
				out[i] = b
				continue
			}
			out[i] = '.'
		}
		return strings.TrimRight(string(out[:len(buf)]), " \t\r\n")
	}

	out := make([]byte, len(buf))
	for i := 0; i < len(buf); i++ {
		b := buf[i]
		if b >= 0x20 && b <= 0x7e {
			out[i] = b
			continue
		}
		out[i] = '.'
	}
	return strings.TrimRight(string(out), " \t\r\n")
}

// parseUndefinedBytes parses the requested value from EXIF metadata.
func (r *Reader) parseUndefinedBytes(t tag.Entry, maxBytes uint32) []byte {
	if !t.IsType(tag.TypeUndefined) {
		return nil
	}
	return r.parseOpaqueBytes(t, maxBytes)
}

// parseOpaqueBytes parses the requested value from EXIF metadata.
func (r *Reader) parseOpaqueBytes(t tag.Entry, maxBytes uint32) []byte {
	if !(t.IsType(tag.TypeUndefined) || t.IsType(tag.TypeByte) || t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul)) {
		return nil
	}
	if maxBytes == 0 {
		return nil
	}
	if t.IsEmbedded() {
		n := t.Size()
		if n > maxBytes {
			n = maxBytes
		}
		t.EmbeddedValue(r.state.buf[:4])
		out := make([]byte, n)
		copy(out, r.state.buf[:n])
		return out
	}
	buf, _, err := r.readTagBytes(t, maxBytes)
	if err != nil || len(buf) == 0 {
		return nil
	}
	out := make([]byte, len(buf))
	copy(out, buf)
	return out
}

// parseUint16 parses the requested value from EXIF metadata.
func (r *Reader) parseUint16(t tag.Entry) uint16 {
	if !t.IsEmbedded() {
		return 0
	}
	switch t.Type {
	case tag.TypeShort:
		return t.EmbeddedShort()
	case tag.TypeLong:
		return uint16(t.EmbeddedLong())
	default:
		return 0
	}
}

// parseUint32 parses the requested value from EXIF metadata.
func (r *Reader) parseUint32(t tag.Entry) uint32 {
	if !t.IsEmbedded() {
		return 0
	}
	switch t.Type {
	case tag.TypeLong:
		return t.EmbeddedLong()
	case tag.TypeShort:
		return uint32(t.EmbeddedShort())
	default:
		return 0
	}
}

// parseRationalU parses the requested value from EXIF metadata.
func (r *Reader) parseRationalU(t tag.Entry) [2]uint32 {
	if !(t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		return [2]uint32{}
	}
	buf, _, err := r.readTagBytes(t, 8)
	if err != nil || len(buf) < 8 {
		return [2]uint32{}
	}
	return [2]uint32{t.ByteOrder.Uint32(buf[:4]), t.ByteOrder.Uint32(buf[4:8])}
}

// parseRationalValue parses the requested value from EXIF metadata.
func (r *Reader) parseRationalValue(t tag.Entry) tag.RationalU {
	parts := r.parseRationalU(t)
	return tag.RationalU{Numerator: parts[0], Denominator: parts[1]}
}

// parseUint16List parses the requested value from EXIF metadata.
func (r *Reader) parseUint16List(t tag.Entry, dst []uint16) int {
	if len(dst) == 0 {
		return 0
	}
	if !t.IsType(tag.TypeShort) {
		return 0
	}
	if t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	if t.IsEmbedded() {
		return t.EmbeddedShorts(dst[:n])
	}
	buf, _, err := r.readTagBytes(t, uint32(n*2))
	if err != nil {
		return 0
	}
	if got := len(buf) / 2; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := 2 * i
		dst[i] = t.ByteOrder.Uint16(buf[start : start+2])
	}
	return n
}

// parseUndefinedUint16List parses uint16 values from UNDEFINED or SHORT payloads.
func (r *Reader) parseUndefinedUint16List(t tag.Entry, dst []uint16) int {
	if len(dst) == 0 || t.UnitCount == 0 {
		return 0
	}
	if t.IsType(tag.TypeShort) {
		return r.parseUint16List(t, dst)
	}
	if !t.IsType(tag.TypeUndefined) {
		return 0
	}

	n := int(t.UnitCount / 2)
	if n > len(dst) {
		n = len(dst)
	}
	if n == 0 {
		return 0
	}

	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		if n > 2 {
			n = 2
		}
		for i := 0; i < n; i++ {
			start := i * 2
			dst[i] = t.ByteOrder.Uint16(r.state.buf[start : start+2])
		}
		return n
	}

	buf, _, err := r.readTagBytes(t, uint32(n*2))
	if err != nil {
		return 0
	}
	if got := len(buf) / 2; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 2
		dst[i] = t.ByteOrder.Uint16(buf[start : start+2])
	}
	return n
}

// parseInt16List parses the requested value from EXIF metadata.
func (r *Reader) parseInt16List(t tag.Entry, dst []int16) int {
	if len(dst) == 0 {
		return 0
	}
	if !t.IsType(tag.TypeSignedShort) || t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	if t.IsEmbedded() {
		var shorts [2]uint16
		m := t.EmbeddedShorts(shorts[:])
		if m > n {
			m = n
		}
		for i := 0; i < m; i++ {
			dst[i] = int16(shorts[i])
		}
		return m
	}
	buf, _, err := r.readTagBytes(t, uint32(n*2))
	if err != nil {
		return 0
	}
	if got := len(buf) / 2; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 2
		dst[i] = int16(t.ByteOrder.Uint16(buf[start : start+2]))
	}
	return n
}

// parseUint32List parses the requested value from EXIF metadata.
func (r *Reader) parseUint32List(t tag.Entry, dst []uint32) int {
	if len(dst) == 0 {
		return 0
	}
	if !(t.IsType(tag.TypeLong) || t.IsType(tag.TypeShort) || t.IsType(tag.TypeIfd)) {
		return 0
	}
	if t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeLong, tag.TypeIfd:
			if n > 0 {
				dst[0] = t.EmbeddedLong()
				return 1
			}
		case tag.TypeShort:
			var shorts [2]uint16
			m := t.EmbeddedShorts(shorts[:])
			if m > n {
				m = n
			}
			for i := 0; i < m; i++ {
				dst[i] = uint32(shorts[i])
			}
			return m
		}
		return 0
	}
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd:
		buf, _, err := r.readTagBytes(t, uint32(n*4))
		if err != nil {
			return 0
		}
		if got := len(buf) / 4; got < n {
			n = got
		}
		for i := 0; i < n; i++ {
			start := i * 4
			dst[i] = t.ByteOrder.Uint32(buf[start : start+4])
		}
	case tag.TypeShort:
		buf, _, err := r.readTagBytes(t, uint32(n*2))
		if err != nil {
			return 0
		}
		if got := len(buf) / 2; got < n {
			n = got
		}
		for i := 0; i < n; i++ {
			start := i * 2
			dst[i] = uint32(t.ByteOrder.Uint16(buf[start : start+2]))
		}
	default:
		return 0
	}
	return n
}

// parseInt32List parses the requested value from EXIF metadata.
func (r *Reader) parseInt32List(t tag.Entry, dst []int32) int {
	if len(dst) == 0 {
		return 0
	}
	if !t.IsType(tag.TypeSignedLong) || t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	if t.IsEmbedded() {
		if n > 0 {
			dst[0] = int32(t.EmbeddedLong())
			return 1
		}
		return 0
	}
	buf, _, err := r.readTagBytes(t, uint32(n*4))
	if err != nil {
		return 0
	}
	if got := len(buf) / 4; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 4
		dst[i] = int32(t.ByteOrder.Uint32(buf[start : start+4]))
	}
	return n
}

// parseByteList parses the requested value from EXIF metadata.
func (r *Reader) parseByteList(t tag.Entry, dst []byte) int {
	if len(dst) == 0 {
		return 0
	}
	if !(t.IsType(tag.TypeByte) || t.IsType(tag.TypeUndefined) || t.IsType(tag.TypeASCII) || t.IsType(tag.TypeASCIINoNul)) {
		return 0
	}
	if t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		copy(dst[:n], r.state.buf[:n])
		return n
	}
	buf, _, err := r.readTagBytes(t, uint32(n))
	if err != nil {
		return 0
	}
	if len(buf) < n {
		n = len(buf)
	}
	copy(dst[:n], buf[:n])
	return n
}

// parseRationalUList parses the requested value from EXIF metadata.
func (r *Reader) parseRationalUList(t tag.Entry, dst []uint32) int {
	if len(dst) < 2 {
		return 0
	}
	if !(t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		return 0
	}
	if t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst)/2 {
		n = len(dst) / 2
	}
	if n == 0 {
		return 0
	}
	buf, _, err := r.readTagBytes(t, uint32(n*8))
	if err != nil {
		return 0
	}
	if got := len(buf) / 8; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 8
		dst[i*2] = t.ByteOrder.Uint32(buf[start : start+4])
		dst[i*2+1] = t.ByteOrder.Uint32(buf[start+4 : start+8])
	}
	return n
}

// parseRationalSList parses the requested value from EXIF metadata.
func (r *Reader) parseRationalSList(t tag.Entry, dst []int32) int {
	if len(dst) < 2 {
		return 0
	}
	if !(t.IsType(tag.TypeSignedRational) || t.IsType(tag.TypeRational)) {
		return 0
	}
	if t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst)/2 {
		n = len(dst) / 2
	}
	if n == 0 {
		return 0
	}
	buf, _, err := r.readTagBytes(t, uint32(n*8))
	if err != nil {
		return 0
	}
	if got := len(buf) / 8; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 8
		dst[i*2] = int32(t.ByteOrder.Uint32(buf[start : start+4]))
		dst[i*2+1] = int32(t.ByteOrder.Uint32(buf[start+4 : start+8]))
	}
	return n
}

// parseFloat64List parses the requested value from EXIF metadata.
func (r *Reader) parseFloat64List(t tag.Entry, dst []float64) int {
	if len(dst) == 0 {
		return 0
	}
	if !t.IsType(tag.TypeDouble) || t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	buf, _, err := r.readTagBytes(t, uint32(n*8))
	if err != nil {
		return 0
	}
	if got := len(buf) / 8; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 8
		dst[i] = math.Float64frombits(t.ByteOrder.Uint64(buf[start : start+8]))
	}
	return n
}

// parseFloat32List parses the requested value from EXIF metadata.
func (r *Reader) parseFloat32List(t tag.Entry, dst []float32) int {
	if len(dst) == 0 {
		return 0
	}
	if !t.IsType(tag.TypeFloat) || t.UnitCount == 0 {
		return 0
	}
	n := int(t.UnitCount)
	if n > len(dst) {
		n = len(dst)
	}
	buf, _, err := r.readTagBytes(t, uint32(n*4))
	if err != nil {
		return 0
	}
	if got := len(buf) / 4; got < n {
		n = got
	}
	for i := 0; i < n; i++ {
		start := i * 4
		dst[i] = math.Float32frombits(t.ByteOrder.Uint32(buf[start : start+4]))
	}
	return n
}

// parseDate parses the requested value from EXIF metadata.
func (r *Reader) parseDate(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 19 {
		return time.Time{}
	}
	// YYYY:MM:DD HH:MM:SS
	if buf[4] != ':' || buf[7] != ':' || buf[10] != ' ' || buf[13] != ':' || buf[16] != ':' {
		return time.Time{}
	}
	return time.Date(
		int(parseStrUint(buf[0:4])),
		time.Month(parseStrUint(buf[5:7])),
		int(parseStrUint(buf[8:10])),
		int(parseStrUint(buf[11:13])),
		int(parseStrUint(buf[14:16])),
		int(parseStrUint(buf[17:19])),
		0,
		time.UTC,
	)
}

// parseOffsetTime parses the requested value from EXIF metadata.
func (r *Reader) parseOffsetTime(t tag.Entry) *time.Location {
	if !t.IsType(tag.TypeASCII) {
		return time.UTC
	}
	buf, _, err := r.readTagBytes(t, 8)
	if err != nil || len(buf) < 6 {
		return time.UTC
	}
	if buf[3] != ':' {
		return time.UTC
	}
	offset := int(parseStrUint(buf[1:3]))*hoursToSeconds + int(parseStrUint(buf[4:6]))*minutesToSeconds
	switch buf[0] {
	case '-':
		return getLocation(int32(offset*-1), buf[:6])
	case '+':
		return getLocation(int32(offset), buf[:6])
	default:
		return time.UTC
	}
}

// parseAperture parses the requested value from EXIF metadata.
func (r *Reader) parseAperture(t tag.Entry) meta.Aperture {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return meta.Aperture(float32(rat[0]) / float32(rat[1]))
}

// parseExposureTime parses the requested value from EXIF metadata.
func (r *Reader) parseExposureTime(t tag.Entry) meta.ExposureTime {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return meta.ExposureTime(float32(rat[0]) / float32(rat[1]))
}

// parseShutterSpeed parses the requested value from EXIF metadata.
func (r *Reader) parseShutterSpeed(t tag.Entry) meta.ShutterSpeed {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return meta.ShutterSpeed(float32(rat[0]) / float32(rat[1]))
}

// parseSignedRationalFloat32 parses a rational (signed or unsigned) as float32.
func (r *Reader) parseSignedRationalFloat32(t tag.Entry) float32 {
	switch {
	case t.IsType(tag.TypeSignedRational):
		var rat [2]int32
		if r.parseRationalSList(t, rat[:]) == 0 {
			return 0
		}
		if rat[1] == 0 {
			return 0
		}
		return float32(rat[0]) / float32(rat[1])
	case t.IsType(tag.TypeRational):
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return 0
		}
		return float32(rat[0]) / float32(rat[1])
	default:
		return 0
	}
}

// parseExposureBias parses the requested value from EXIF metadata.
func (r *Reader) parseExposureBias(t tag.Entry) meta.ExposureBias {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return meta.NewExposureBias(0, 0)
	}
	return meta.NewExposureBias(int16(rat[0]), int16(rat[1]))
}

// parseFocalLength parses the requested value from EXIF metadata.
func (r *Reader) parseFocalLength(t tag.Entry) meta.FocalLength {
	switch t.Type {
	case tag.TypeShort, tag.TypeLong:
		return meta.NewFocalLength(r.parseUint32(t), 1)
	case tag.TypeRational, tag.TypeSignedRational:
		rat := r.parseRationalU(t)
		if rat[1] == 0 {
			return 0
		}
		return meta.FocalLength(float32(rat[0]) / float32(rat[1]))
	default:
		return 0
	}
}

// parseLensInfo parses the requested value from EXIF metadata.
func (r *Reader) parseLensInfo(t tag.Entry) LensInfo {
	if t.IsEmbedded() {
		return LensInfo{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 32 {
		return LensInfo{}
	}
	return LensInfo{
		t.ByteOrder.Uint32(buf[:4]),
		t.ByteOrder.Uint32(buf[4:8]),
		t.ByteOrder.Uint32(buf[8:12]),
		t.ByteOrder.Uint32(buf[12:16]),
		t.ByteOrder.Uint32(buf[16:20]),
		t.ByteOrder.Uint32(buf[20:24]),
		t.ByteOrder.Uint32(buf[24:28]),
		t.ByteOrder.Uint32(buf[28:32]),
	}
}

// parseSceneType parses Exif SceneType from BYTE/UNDEFINED/ASCII or integer encodings.
func (r *Reader) parseSceneType(t tag.Entry) uint16 {
	switch {
	case t.IsType(tag.TypeShort), t.IsType(tag.TypeLong):
		return uint16(r.parseUint32(t))
	case t.IsType(tag.TypeASCII), t.IsType(tag.TypeASCIINoNul):
		s := strings.TrimSpace(r.parseString(t))
		if s == "" {
			return 0
		}
		return uint16(parseStrUint([]byte(s)))
	case t.IsType(tag.TypeByte), t.IsType(tag.TypeUndefined):
		var b [4]byte
		if r.parseByteList(t, b[:]) == 0 {
			return 0
		}
		return uint16(b[0])
	default:
		return 0
	}
}

// parseGPSRef parses and normalizes GPS reference tags.
func (r *Reader) parseGPSRef(t tag.Entry) tag.GPSRef {
	first := byte(0)
	ok := false
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		first = r.state.buf[0]
		ok = true
	} else {
		buf, _, err := r.readTagBytes(t, 1)
		if err == nil && len(buf) > 0 {
			first = buf[0]
			ok = true
		}
	}
	if !ok {
		return tag.GPSRefUnknown
	}

	switch t.ID {
	case tag.TagGPSAltitudeRef:
		// ExifTool documents altitude ref values:
		// 0 above sea level, 1 below, 2 above sea-level reference, 3 below.
		if first == 1 || first == 3 {
			return tag.GPSRefBelowSeaLevel
		}
		if first == 0 || first == 2 {
			return tag.GPSRefAboveSeaLevel
		}
		return tag.GPSRefUnknown
	case tag.TagGPSLatitudeRef, tag.TagGPSDestLatitudeRef:
		if first == 'S' || first == 's' {
			return tag.GPSRefSouth
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefNorth
		}
		return tag.GPSRefUnknown
	case tag.TagGPSLongitudeRef, tag.TagGPSDestLongitudeRef:
		if first == 'W' || first == 'w' {
			return tag.GPSRefWest
		}
		if first == 'E' || first == 'e' {
			return tag.GPSRefEast
		}
		return tag.GPSRefUnknown
	case tag.TagGPSSpeedRef:
		if first == 'K' || first == 'k' {
			return tag.GPSRefKilometersPerHour
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMilesPerHour
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefKnots
		}
		return tag.GPSRefUnknown
	case tag.TagGPSTrackRef, tag.TagGPSImgDirectionRef, tag.TagGPSDestBearingRef:
		if first == 'T' || first == 't' {
			return tag.GPSRefTrueDirection
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMagneticDirection
		}
		return tag.GPSRefUnknown
	case tag.TagGPSDestDistanceRef:
		if first == 'K' || first == 'k' {
			return tag.GPSRefKilometers
		}
		if first == 'M' || first == 'm' {
			return tag.GPSRefMiles
		}
		if first == 'N' || first == 'n' {
			return tag.GPSRefNauticalMiles
		}
		return tag.GPSRefUnknown
	}
	return tag.GPSRefUnknown
}

// parseGPSCoord parses the requested value from EXIF metadata.
func (r *Reader) parseGPSCoord(t tag.Entry) float64 {
	if t.UnitCount != 3 {
		return 0
	}
	if !(t.IsType(tag.TypeRational) || t.IsType(tag.TypeSignedRational)) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	dNum := t.ByteOrder.Uint32(buf[:4])
	dDen := t.ByteOrder.Uint32(buf[4:8])
	mNum := t.ByteOrder.Uint32(buf[8:12])
	mDen := t.ByteOrder.Uint32(buf[12:16])
	sNum := t.ByteOrder.Uint32(buf[16:20])
	sDen := t.ByteOrder.Uint32(buf[20:24])
	if dDen == 0 || mDen == 0 || sDen == 0 {
		return 0
	}
	coord := (float64(dNum) / float64(dDen))
	coord += (float64(mNum) / float64(mDen) / 60.0)
	coord += (float64(sNum) / float64(sDen) / 3600.0)
	return coord
}

// parseGPSAltitude parses the requested value from EXIF metadata.
func (r *Reader) parseGPSAltitude(t tag.Entry) float32 {
	rat := r.parseRationalU(t)
	if rat[1] == 0 {
		return 0
	}
	return float32(rat[0]) / float32(rat[1])
}

// parseGPSTimeStamp parses the requested value from EXIF metadata.
func (r *Reader) parseGPSTimeStamp(t tag.Entry) time.Duration {
	if t.UnitCount != 3 || !t.IsType(tag.TypeRational) {
		return 0
	}
	buf, _, err := r.readTagBytes(t, 24)
	if err != nil || len(buf) < 24 {
		return 0
	}
	v := [6]uint32{
		t.ByteOrder.Uint32(buf[:4]),
		t.ByteOrder.Uint32(buf[4:8]),
		t.ByteOrder.Uint32(buf[8:12]),
		t.ByteOrder.Uint32(buf[12:16]),
		t.ByteOrder.Uint32(buf[16:20]),
		t.ByteOrder.Uint32(buf[20:24]),
	}
	return rationalDuration(v[0], v[1], time.Hour) +
		rationalDuration(v[2], v[3], time.Minute) +
		rationalDuration(v[4], v[5], time.Second)
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

// parseGPSDateStamp parses the requested value from EXIF metadata.
func (r *Reader) parseGPSDateStamp(t tag.Entry) time.Time {
	if !t.IsType(tag.TypeASCII) {
		return time.Time{}
	}
	buf, _, err := r.readTagBytes(t, 32)
	if err != nil || len(buf) < 10 {
		return time.Time{}
	}
	if len(buf) >= 19 && buf[10] == ' ' && buf[13] == ':' && buf[16] == ':' {
		return time.Date(
			int(parseStrUint(buf[0:4])),
			time.Month(parseStrUint(buf[5:7])),
			int(parseStrUint(buf[8:10])),
			int(parseStrUint(buf[11:13])),
			int(parseStrUint(buf[14:16])),
			int(parseStrUint(buf[17:19])),
			0,
			time.UTC,
		)
	}
	if buf[4] == ':' && buf[7] == ':' {
		return time.Date(
			int(parseStrUint(buf[0:4])),
			time.Month(parseStrUint(buf[5:7])),
			int(parseStrUint(buf[8:10])),
			0, 0, 0, 0, time.UTC,
		)
	}
	return time.Time{}
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
	if n < 0 {
		return nil, imagetype.ErrDataLength
	}
	if n == 0 {
		return nil, nil
	}
	if n > len(r.state.buf) {
		return nil, imagetype.ErrDataLength
	}
	if r.exifLength > 0 && int(r.po)+n > int(r.exifLength) {
		return nil, imagetype.ErrDataLength
	}
	if br, ok := r.reader.(BufferedReader); ok {
		buf, err := br.Peek(n)
		if err != nil {
			return nil, err
		}
		readCount, err := br.Discard(len(buf))
		r.po += uint32(readCount)
		return buf, err
	}
	readCount, err := r.reader.Read(r.state.buf[:n])
	r.po += uint32(readCount)
	if err != nil {
		return nil, err
	}
	return r.state.buf[:readCount], nil
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
	if br, ok := r.reader.(BufferedReader); ok {
		discarded, err := br.Discard(n)
		r.po += uint32(discarded)
		return err
	}
	for n > 0 {
		chunk := n
		if chunk > len(r.state.discardBuf) {
			chunk = len(r.state.discardBuf)
		}
		readCount, err := r.reader.Read(r.state.discardBuf[:chunk])
		r.po += uint32(readCount)
		n -= readCount
		if err != nil {
			return err
		}
	}
	return nil
}

// readUint16 reads data from the underlying stream or parser buffers.
func (r *Reader) readUint16(directory ifd.Directory) (uint16, error) {
	buf, err := r.fastRead(2)
	if err != nil || len(buf) < 2 {
		return 0, err
	}
	return directory.ByteOrder.Uint16(buf), nil
}

// readUint32 reads data from the underlying stream or parser buffers.
func (r *Reader) readUint32(directory ifd.Directory) (uint32, error) {
	buf, err := r.fastRead(4)
	if err != nil || len(buf) < 4 {
		return 0, err
	}
	return directory.ByteOrder.Uint32(buf), nil
}

// TODO(simd): accelerate trim/ASCII scan and offset-date parsing with simd/archsimd
// when GOEXPERIMENT=simd integration is prioritized.
