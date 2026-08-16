package exif

import (
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

// trimTrailingNULBytes removes trailing NUL bytes without trimming any other
// characters. This is used for fields like IFD0 Make where the raw text should
// be preserved apart from EXIF string terminators.
func trimTrailingNULBytes(buf []byte) []byte {
	for i := len(buf) - 1; i >= 0; i-- {
		if buf[i] == 0 {
			continue
		}
		return buf[:i+1]
	}
	return nil
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

func trimASCIIWhitespace(buf []byte) []byte {
	start := 0
	for start < len(buf) {
		switch buf[start] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			start++
		default:
			goto trimRight
		}
	}
	return nil

trimRight:
	end := len(buf)
	for end > start {
		switch buf[end-1] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			end--
		default:
			return buf[start:end]
		}
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
	trimmed := trimASCIIWhitespace(trimNULBuffer(buf))
	if len(trimmed) == 0 {
		return ""
	}
	return string(trimmed)
}

func trimAtNUL(buf []byte) []byte {
	for i, b := range buf {
		if b == 0 {
			return buf[:i]
		}
	}
	return buf
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
	unitU, ok := meta.SafecastInt64ToUint64(int64(unit))
	if !ok {
		return 0
	}
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
	value, ok := meta.SafecastUint64ToInt64(wholeUnits + frac)
	if !ok {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(value)
}

// readTagBytes reads data from the underlying stream or parser buffers.
func (r *Reader) readTagBytes(t tag.Entry, maxBytes uint32) (buf []byte, err error) {
	if err = r.seekToTag(t); err != nil {
		r.warnTagReadError(t, err, "failed seeking to tag payload")
		return nil, err
	}

	size := t.Size()
	if size == 0 {
		return nil, nil
	}
	if maxBytes > 0 && size > maxBytes {
		size = maxBytes
	}
	if size > uint32(len(r.state.buf)) {
		size = uint32(len(r.state.buf))
	}

	buf, err = r.fastRead(int(size))
	if err != nil {
		r.warnTagReadError(t, err, "failed reading tag payload")
		return nil, err
	}

	remaining := int(t.Size() - size)
	if remaining > 0 {
		if discardErr := r.discard(remaining); discardErr != nil {
			r.warnTagReadError(t, discardErr, "failed discarding unread tag payload")
			return nil, discardErr
		}
	}
	return buf, nil
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
	if delta, ok := meta.SafecastIntToUint32(readCount); ok {
		r.po += delta
	}
	return buf, err
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
	if delta, ok := meta.SafecastIntToUint32(discarded); ok {
		r.po += delta
	}
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

func (r *Reader) canonModelID() metacanon.CanonCameraModel {
	if r.Exif.MakerNote.Canon == nil {
		return metacanon.CanonModelUnknown
	}
	return metacanon.CanonCameraModel(r.Exif.MakerNote.Canon.ModelID)
}

func (r *Reader) canonShotInfoLegacyExposureTime() bool {
	modelID := r.canonModelID()
	if modelID != metacanon.CanonModelEOS20D && modelID != metacanon.CanonModelEOSDigitalRebelXT {
		return false
	}
	if r.Exif.MakerNote.Canon == nil {
		return false
	}
	return r.Exif.MakerNote.Canon.CanonCameraSettings.FocalUnits > 1
}
