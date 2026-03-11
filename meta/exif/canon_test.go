package exif

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	metacanon "github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func u16(v int16) uint16 {
	return uint16(v)
}

func parseAFInfo2ForTest(t *testing.T, words []uint16, model string, isAFInfo3 bool, opts ...AFInfoDecodeOptions) metacanon.AFInfo {
	t.Helper()
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	tagID := metacanon.CanonAFInfo2
	if isAFInfo3 {
		tagID = metacanon.AFInfo3
	}
	entry := tag.NewEntry(
		tag.ID(tagID),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	readerOpts := make([]ReaderOption, 0, 1)
	if len(opts) > 0 {
		readerOpts = append(readerOpts, WithAFInfoDecodeOptions(opts[0]))
	}
	r := NewReader(Logger, readerOpts...)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	r.Exif.IFD0.Model = model
	return r.parseCanonAFInfo2(entry)
}

func TestFillCanonAFInfoEOS(t *testing.T) {
	words := make([]uint16, 19)
	words[0] = 5 // NumAFPoints
	words[1] = 5 // ValidAFPoints
	words[6] = 10
	words[7] = 12

	words[8] = u16(-2)
	words[9] = u16(-1)
	words[10] = 0
	words[11] = 1
	words[12] = 2

	words[13] = 3
	words[14] = 2
	words[15] = 1
	words[16] = 0
	words[17] = u16(-1)

	words[18] = (1 << 0) | (1 << 4)

	var got metacanon.AFInfo
	fillCanonAFInfo(&got, words, "Canon EOS 6D", len(words))

	if got.PrimaryAFPoint != 0 {
		t.Fatalf("PrimaryAFPoint = %d, want 0 for EOS", got.PrimaryAFPoint)
	}
	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{0, 4}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [0 4]", got.AFPointsInFocusBits)
	}
	if got.AFPointsSelectedBits != nil {
		t.Fatalf("AFPointsSelectedBits = %v, want nil", got.AFPointsSelectedBits)
	}
	if len(got.AFPoints) != 5 {
		t.Fatalf("len(AFPoints) = %d, want 5", len(got.AFPoints))
	}
}

func TestFillCanonAFInfoNonEOSCount36PrimaryOffset(t *testing.T) {
	words := make([]uint16, 36)
	words[0] = 9 // NumAFPoints
	words[1] = 9 // ValidAFPoints
	words[6] = 8
	words[7] = 8

	// Sequence 10 mask starts at 8 + 2*9 = 26.
	words[26] = 1 << 3

	// Sequence 11/12 behavior for AFInfoCount==36:
	// seq11 PrimaryAFPoint is skipped, seq11 unknown[8], seq12 PrimaryAFPoint.
	words[27] = 2 // should be ignored
	words[35] = 6 // expected primary point

	var got metacanon.AFInfo
	fillCanonAFInfo(&got, words, "PowerShot G1", 36)

	if got.PrimaryAFPoint != 6 {
		t.Fatalf("PrimaryAFPoint = %d, want 6", got.PrimaryAFPoint)
	}
	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{3}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [3]", got.AFPointsInFocusBits)
	}
}

func TestFillCanonAFInfo2EOSSelectedBits(t *testing.T) {
	words := make([]uint16, 38)
	words[1] = 2 // AFAreaMode
	words[2] = 7 // NumAFPoints
	words[3] = 7 // ValidAFPoints
	words[4] = 100
	words[5] = 80

	// widths/heights/x/y blocks of NumAFPoints each.
	for i := 0; i < 7; i++ {
		words[8+i] = 4
		words[15+i] = 6
		words[22+i] = uint16(i)
		words[29+i] = uint16(i)
	}

	// in-focus and selected masks.
	words[36] = 1 << 1
	words[37] = 1 << 2

	got := parseAFInfo2ForTest(t, words, "Canon EOS R6", false)

	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{1}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [1]", got.AFPointsInFocusBits)
	}
	if !reflect.DeepEqual(got.AFPointsSelectedBits, []int{2}) {
		t.Fatalf("AFPointsSelectedBits = %v, want [2]", got.AFPointsSelectedBits)
	}
	if got.PrimaryAFPoint != 0 {
		t.Fatalf("PrimaryAFPoint = %d, want 0 for EOS AFInfo2", got.PrimaryAFPoint)
	}
}

func TestFillCanonAFInfo2NonEOSPrimaryOffset(t *testing.T) {
	words := make([]uint16, 40)
	words[1] = 4 // AFAreaMode
	words[2] = 7 // NumAFPoints
	words[3] = 7 // ValidAFPoints

	// in-focus mask at seq 12.
	words[36] = 1 << 5
	// seq 13 unknown has maskWordCount+1 values for non-EOS.
	words[39] = 6 // seq 14 PrimaryAFPoint

	got := parseAFInfo2ForTest(t, words, "PowerShot G1", false)

	if got.PrimaryAFPoint != 6 {
		t.Fatalf("PrimaryAFPoint = %d, want 6", got.PrimaryAFPoint)
	}
	if got.AFPointsSelectedBits != nil {
		t.Fatalf("AFPointsSelectedBits = %v, want nil for non-EOS AFInfo2", got.AFPointsSelectedBits)
	}

	gotAFInfo3 := parseAFInfo2ForTest(t, words, "PowerShot G1", true)
	if gotAFInfo3.PrimaryAFPoint != 0 {
		t.Fatalf("PrimaryAFPoint(AFInfo3) = %d, want 0", gotAFInfo3.PrimaryAFPoint)
	}
}

func TestParseCanonAFInfo2DecodeOptionsBitset(t *testing.T) {
	words := make([]uint16, 38)
	words[1] = 2 // AFAreaMode
	words[2] = 7 // NumAFPoints
	words[3] = 7 // ValidAFPoints
	words[4] = 100
	words[5] = 80
	for i := 0; i < 7; i++ {
		words[8+i] = 4
		words[15+i] = 6
		words[22+i] = uint16(i)
		words[29+i] = uint16(i)
	}
	words[36] = 1 << 1
	words[37] = 1 << 2

	got := parseAFInfo2ForTest(
		t,
		words,
		"Canon EOS R6",
		false,
		AFInfoDecodeInFocus|AFInfoDecodeSelected,
	)

	if got.AFAreaWidths != nil || got.AFAreaHeights != nil || got.AFAreaXPositions != nil || got.AFAreaYPositions != nil {
		t.Fatalf("expected coord slices to be nil when AFInfoDecodeCoords is off")
	}
	if got.AFPoints != nil {
		t.Fatalf("expected AFPoints to be nil when AFInfoDecodePoints is off")
	}
	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{1}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [1]", got.AFPointsInFocusBits)
	}
	if !reflect.DeepEqual(got.AFPointsSelectedBits, []int{2}) {
		t.Fatalf("AFPointsSelectedBits = %v, want [2]", got.AFPointsSelectedBits)
	}
}

func TestParseCanonAFInfo2DecodeOptionsNone(t *testing.T) {
	words := make([]uint16, 40)
	words[1] = 4 // AFAreaMode
	words[2] = 7 // NumAFPoints
	words[3] = 7 // ValidAFPoints
	words[36] = 1 << 5
	words[39] = 6 // seq 14 PrimaryAFPoint

	got := parseAFInfo2ForTest(t, words, "PowerShot G1", false, 0)

	if got.AFPointsInFocusBits != nil {
		t.Fatalf("expected in-focus bitsets to be nil when decode options are empty")
	}
	if got.AFPointsSelectedBits != nil {
		t.Fatalf("expected selected bitsets to be nil when decode options are empty")
	}
	if got.AFPoints != nil {
		t.Fatalf("expected AFPoints to be nil when decode options are empty")
	}
	if got.PrimaryAFPoint != 6 {
		t.Fatalf("PrimaryAFPoint = %d, want 6", got.PrimaryAFPoint)
	}
}

func TestParseCanonAFConfigTagPresent(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	candidates := []string{
		"1.CR2",
		"1.CR3",
		"2.CR3",
		"CanonR3.CR3",
		"CanonR8.CR3",
		"EOSR6III.CR3",
	}

	var samplePath string
	for i := range candidates {
		p := filepath.Join(benchDir, candidates[i])
		if _, err := os.Stat(p); err == nil {
			samplePath = p
			break
		}
	}
	if samplePath == "" {
		t.Skipf("no AFConfig sample found in %s", benchDir)
	}

	f, err := os.Open(samplePath)
	if err != nil {
		t.Fatalf("open %s: %v", samplePath, err)
	}
	defer func() {
		_ = f.Close()
	}()

	parsed, err := Parse(f)
	if err != nil {
		t.Fatalf("parse %s: %v", samplePath, err)
	}

	if !parsed.MakerNote.HasTagParsed(uint16(metacanon.CanonAFConfig)) {
		t.Fatalf("expected CanonAFConfig (0x%04x) to be marked parsed for %s", uint16(metacanon.CanonAFConfig), samplePath)
	}
	if parsed.MakerNote.Canon.AFConfig.AFConfigTool == 0 {
		t.Fatalf("expected non-zero AFConfigTool for %s", samplePath)
	}
}

var sinkCanonCameraSettings metacanon.CameraSettings

func benchmarkCanonCameraSettingsData() ([64]int16, [64]uint16) {
	var signed [64]int16
	var unsigned [64]uint16
	for i := range signed {
		v := int16(i*13 - 377)
		signed[i] = v
		unsigned[i] = uint16(v)
	}

	// Ensure unsigned fields keep high bits in decode.
	unsigned[18] = 0xf0f1
	signed[18] = int16(unsigned[18])
	unsigned[21] = 0x8123
	signed[21] = int16(unsigned[21])
	unsigned[22] = 0x8abc
	signed[22] = int16(unsigned[22])
	unsigned[23] = 0xfedc
	signed[23] = int16(unsigned[23])
	unsigned[24] = 0xabcd
	signed[24] = int16(unsigned[24])
	unsigned[28] = 0xa55a
	signed[28] = int16(unsigned[28])
	unsigned[34] = 0xff10
	signed[34] = int16(unsigned[34])
	unsigned[35] = 0x7f11
	signed[35] = int16(unsigned[35])
	unsigned[36] = 0xbeef
	signed[36] = int16(unsigned[36])

	return signed, unsigned
}

func canonUint16WordsToBytes(words []uint16, bo utils.ByteOrder) []byte {
	out := make([]byte, len(words)*2)
	for i := range words {
		bo.PutUint16(out[i*2:], words[i])
	}
	return out
}

func TestParseCanonCameraSettingsIndexMapping(t *testing.T) {
	signed, unsigned := benchmarkCanonCameraSettingsData()
	raw := canonUint16WordsToBytes(unsigned[:52], utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonCameraSettings),
		tag.TypeShort,
		52,
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)

	if got.MacroMode != metacanon.MacroMode(signed[0]) {
		t.Fatalf("MacroMode = %d, want %d", got.MacroMode, metacanon.MacroMode(signed[0]))
	}
	if got.FocusMode != metacanon.FocusMode(signed[6]) {
		t.Fatalf("FocusMode = %d, want %d", got.FocusMode, signed[6])
	}
	if got.AFPoint != unsigned[18] {
		t.Fatalf("AFPoint = %#x, want %#x", got.AFPoint, unsigned[18])
	}
	if got.CanonExposureMode != metacanon.ExposureMode(signed[19]) {
		t.Fatalf("CanonExposureMode = %d, want %d", got.CanonExposureMode, signed[19])
	}
	if got.LensType != unsigned[21] {
		t.Fatalf("LensType = %#x, want %#x", got.LensType, unsigned[21])
	}
	if got.DisplayAperture != unsigned[34] {
		t.Fatalf("DisplayAperture = %#x, want %#x", got.DisplayAperture, unsigned[34])
	}
	if got.SRAWQuality != metacanon.SRAWQuality(signed[45]) {
		t.Fatalf("SRAWQuality = %d, want %d", got.SRAWQuality, metacanon.SRAWQuality(signed[45]))
	}
	if got.FocusBracketing != metacanon.FocusBracketing(signed[49]) {
		t.Fatalf("FocusBracketing = %d, want %d", got.FocusBracketing, metacanon.FocusBracketing(signed[49]))
	}
	if got.Clarity != signed[50] {
		t.Fatalf("Clarity = %d, want %d", got.Clarity, signed[50])
	}
	if got.HDRPQ != metacanon.HDRPQ(signed[51]) {
		t.Fatalf("HDRPQ = %d, want %d", got.HDRPQ, metacanon.HDRPQ(signed[51]))
	}
}

func TestParseCanonCameraSettingsShortInputProgressive(t *testing.T) {
	words := []uint16{123, 456, 789}
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)
	if got.MacroMode != metacanon.MacroMode(123) {
		t.Fatalf("MacroMode = %d, want 123", got.MacroMode)
	}
	if got.SelfTimer != 456 {
		t.Fatalf("SelfTimer = %d, want 456", got.SelfTimer)
	}
	if got.Quality != metacanon.Quality(789) {
		t.Fatalf("Quality = %d, want 789", got.Quality)
	}
	if got.HDRPQ != 0 {
		t.Fatalf("HDRPQ = %d, want 0", got.HDRPQ)
	}
}

func TestParseCanonCameraSettingsTooShortReturnsZero(t *testing.T) {
	words := []uint16{123}
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)
	if got != (metacanon.CameraSettings{}) {
		t.Fatalf("parseCanonCameraSettings(too short) = %+v, want zero value", got)
	}
}

func BenchmarkParseCanonCameraSettings(b *testing.B) {
	_, unsigned := benchmarkCanonCameraSettingsData()
	input := unsigned[:52]
	raw := canonUint16WordsToBytes(input, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(input)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader

	b.ReportAllocs()
	b.SetBytes(int64(len(raw)))

	var dst metacanon.CameraSettings
	for i := 0; i < b.N; i++ {
		br.Reset(raw)
		r.Reset(&br)
		dst = r.parseCanonCameraSettings(entry)
	}
	sinkCanonCameraSettings = dst
}

func TestParseCanonTimeInfo(t *testing.T) {
	cases := []struct {
		name string
		bo   utils.ByteOrder
	}{
		{name: "little-endian", bo: utils.LittleEndian},
		{name: "big-endian", bo: utils.BigEndian},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tz := int32(-300)
			city := metacanon.TimeZoneCityNewYork
			daylightSavings := metacanon.DaylightSavingsOn

			raw := make([]byte, 12)
			tc.bo.PutUint32(raw[0:4], uint32(tz))
			tc.bo.PutUint32(raw[4:8], uint32(city))
			tc.bo.PutUint32(raw[8:12], uint32(daylightSavings))

			entry := tag.NewEntry(
				tag.ID(metacanon.TimeInfo),
				tag.TypeLong,
				3,
				0,
				ifd.MakerNoteIFD,
				0,
				tc.bo,
			)

			r := NewReader(Logger)
			defer r.Close()
			var br bytes.Reader
			br.Reset(raw)
			r.Reset(&br)

			got := r.parseCanonTimeInfo(entry)

			if got.TimeZone != tz {
				t.Fatalf("TimeZone = %d, want %d", got.TimeZone, tz)
			}
			if got.TimeZoneCity != city {
				t.Fatalf("TimeZoneCity = %d, want %d", got.TimeZoneCity, city)
			}
			if got.DaylightSavings != daylightSavings {
				t.Fatalf("DaylightSavings = %d, want %d", got.DaylightSavings, daylightSavings)
			}
		})
	}
}

func TestParseCanonFaceDetect1IndexMapping(t *testing.T) {
	words := make([]uint16, 26)
	words[2] = 2
	words[3] = 640
	words[4] = 480
	words[8] = u16(-10)
	words[9] = u16(20)
	words[10] = u16(30)
	words[11] = u16(-40)

	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.FaceDetect1),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonFaceDetect1(entry)
	if got.FacesDetected != 2 {
		t.Fatalf("FacesDetected = %d, want 2", got.FacesDetected)
	}
	if got.FaceDetectFrameSize != [2]uint16{640, 480} {
		t.Fatalf("FaceDetectFrameSize = %v, want [640 480]", got.FaceDetectFrameSize)
	}
	if got.FacePositions[0] != (metacanon.FacePosition{X: -10, Y: 20}) {
		t.Fatalf("FacePositions[0] = %+v, want {-10, 20}", got.FacePositions[0])
	}
	if got.FacePositions[1] != (metacanon.FacePosition{X: 30, Y: -40}) {
		t.Fatalf("FacePositions[1] = %+v, want {30, -40}", got.FacePositions[1])
	}
}

func TestParseCanonFaceDetect2IndexMapping(t *testing.T) {
	raw := []byte{0, 42, 3, 0, 0}
	entry := tag.NewEntry(
		tag.ID(metacanon.FaceDetect2),
		tag.TypeByte,
		uint32(len(raw)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonFaceDetect2(entry)
	if got.FaceWidth != 42 {
		t.Fatalf("FaceWidth = %d, want 42", got.FaceWidth)
	}
	if got.FacesDetected != 3 {
		t.Fatalf("FacesDetected = %d, want 3", got.FacesDetected)
	}
}

func TestParseCanonFaceDetect3IndexMapping(t *testing.T) {
	words := []uint16{0, 1, 1, 4}
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.FaceDetect3),
		tag.TypeShort,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonFaceDetect3(entry)
	if got.FacesDetected != 4 {
		t.Fatalf("FacesDetected = %d, want 4", got.FacesDetected)
	}
}

func canonUint32WordsToBytes(words []uint32, bo utils.ByteOrder) []byte {
	out := make([]byte, len(words)*4)
	for i := range words {
		bo.PutUint32(out[i*4:], words[i])
	}
	return out
}

func TestParseCanonLightingOptIndexMapping(t *testing.T) {
	words := make([]uint32, 11)
	words[0] = 1  // PeripheralIlluminationCorr
	words[1] = 2  // AutoLightingOptimizer
	words[2] = 1  // HighlightTonePriority
	words[3] = 2  // LongExposureNoiseReduction
	words[4] = 3  // HighISONoiseReduction
	words[9] = 2  // DigitalLensOptimizer
	words[10] = 1 // DualPixelRaw

	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonLightingOpt),
		tag.TypeLong,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonLightingOpt(entry)
	if got.PeripheralIlluminationCorr != 1 {
		t.Fatalf("PeripheralIlluminationCorr = %d, want 1", got.PeripheralIlluminationCorr)
	}
	if got.AutoLightingOptimizer != 2 {
		t.Fatalf("AutoLightingOptimizer = %d, want 2", got.AutoLightingOptimizer)
	}
	if got.HighlightTonePriority != 1 {
		t.Fatalf("HighlightTonePriority = %d, want 1", got.HighlightTonePriority)
	}
	if got.LongExposureNoiseReduction != 2 {
		t.Fatalf("LongExposureNoiseReduction = %d, want 2", got.LongExposureNoiseReduction)
	}
	if got.HighISONoiseReduction != 3 {
		t.Fatalf("HighISONoiseReduction = %d, want 3", got.HighISONoiseReduction)
	}
	if got.DigitalLensOptimizer != 2 {
		t.Fatalf("DigitalLensOptimizer = %d, want 2", got.DigitalLensOptimizer)
	}
	if got.DualPixelRaw != 1 {
		t.Fatalf("DualPixelRaw = %d, want 1", got.DualPixelRaw)
	}
}

func TestParseCanonLightingOptShortPayload(t *testing.T) {
	words := []uint32{1, 3, 1}
	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(metacanon.CanonLightingOpt),
		tag.TypeLong,
		uint32(len(words)),
		0,
		ifd.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonLightingOpt(entry)
	if got.PeripheralIlluminationCorr != 1 {
		t.Fatalf("PeripheralIlluminationCorr = %d, want 1", got.PeripheralIlluminationCorr)
	}
	if got.AutoLightingOptimizer != 3 {
		t.Fatalf("AutoLightingOptimizer = %d, want 3", got.AutoLightingOptimizer)
	}
	if got.HighlightTonePriority != 1 {
		t.Fatalf("HighlightTonePriority = %d, want 1", got.HighlightTonePriority)
	}
	if got.HighISONoiseReduction != 0 {
		t.Fatalf("HighISONoiseReduction = %d, want 0", got.HighISONoiseReduction)
	}
	if got.DigitalLensOptimizer != 0 {
		t.Fatalf("DigitalLensOptimizer = %d, want 0", got.DigitalLensOptimizer)
	}
}
