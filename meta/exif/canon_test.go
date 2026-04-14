package exif

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func u16(v int16) uint16 {
	return uint16(v)
}

func parseAFInfo2ForTest(t *testing.T, words []uint16, model string, isAFInfo3 bool, opts ...AFInfoDecodeOptions) canon.AFInfo {
	t.Helper()
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	tagID := canon.CanonAFInfo2
	if isAFInfo3 {
		tagID = canon.AFInfo3
	}
	entry := tag.NewEntry(
		tag.ID(tagID),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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

func parseShotInfoForTest(t *testing.T, words []uint16, model string, focalUnits uint16) canon.ShotInfo {
	t.Helper()
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonShotInfo),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	r.Exif.IFD0.Model = model
	r.Exif.MakerNote.Canon = &canon.Canon{}
	r.Exif.MakerNote.Canon.CanonCameraSettings.FocalUnits = focalUnits
	return r.parseCanonShotInfo(entry)
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

	var got canon.AFInfo
	fillCanonAFInfo(&got, words, "Canon EOS 6D", len(words))

	if got.PrimaryAFPoint != 0 {
		t.Fatalf("PrimaryAFPoint = %d, want 0 for EOS", got.PrimaryAFPoint)
	}
	if got.Source != canon.AFInfoSourceAFInfo {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo", got.Source)
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

	var got canon.AFInfo
	fillCanonAFInfo(&got, words, "PowerShot G1", 36)

	if got.PrimaryAFPoint != 6 {
		t.Fatalf("PrimaryAFPoint = %d, want 6", got.PrimaryAFPoint)
	}
	if got.Source != canon.AFInfoSourceAFInfo {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo", got.Source)
	}
	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{3}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [3]", got.AFPointsInFocusBits)
	}
}

func TestFillCanonAFInfoNonEOSPrimarySequence12(t *testing.T) {
	words := make([]uint16, 21)
	words[0] = 5 // NumAFPoints
	words[1] = 5 // ValidAFPoints
	words[6] = 8
	words[7] = 8

	// Sequence 10 mask starts at 8 + 2*5 = 18.
	words[18] = 1 << 3

	// Canon.pm sequence 11 is an alternative: for non-EOS records that are not
	// AFInfoCount==36, seq11 is PrimaryAFPoint and seq12 follows immediately.
	words[19] = 2
	words[20] = 6

	var got canon.AFInfo
	fillCanonAFInfo(&got, words, "PowerShot G1", len(words))

	if got.PrimaryAFPoint != 6 {
		t.Fatalf("PrimaryAFPoint = %d, want 6 from sequence 12", got.PrimaryAFPoint)
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
	if got.Source != canon.AFInfoSourceAFInfo2 {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo2", got.Source)
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
	if got.Source != canon.AFInfoSourceAFInfo2 {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo2", got.Source)
	}
	if got.AFPointsSelectedBits != nil {
		t.Fatalf("AFPointsSelectedBits = %v, want nil for non-EOS AFInfo2", got.AFPointsSelectedBits)
	}

	gotAFInfo3 := parseAFInfo2ForTest(t, words, "PowerShot G1", true)
	if gotAFInfo3.PrimaryAFPoint != 0 {
		t.Fatalf("PrimaryAFPoint(AFInfo3) = %d, want 0", gotAFInfo3.PrimaryAFPoint)
	}
	if gotAFInfo3.Source != canon.AFInfoSourceAFInfo3 {
		t.Fatalf("Source(AFInfo3) = %v, want AFInfoSourceAFInfo3", gotAFInfo3.Source)
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

	if got.AFArea != nil {
		t.Fatalf("expected AFArea to be nil when AFInfoDecodeCoords is off")
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
	if got.Source != canon.AFInfoSourceAFInfo2 {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo2", got.Source)
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
	if got.Source != canon.AFInfoSourceAFInfo2 {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo2", got.Source)
	}
}

func TestParseCanonBatteryType(t *testing.T) {
	withHeader := func(payload []byte) []byte {
		raw := make([]byte, canonBatteryTypePayloadSize)
		copy(raw[:4], []byte{0xde, 0xad, 0xbe, 0xef})
		copy(raw[4:], payload)
		return raw
	}

	tests := []struct {
		name      string
		raw       []byte
		unitCount uint32
		want      string
	}{
		{
			name:      "invalid length",
			raw:       make([]byte, canonBatteryTypePayloadSize-1),
			unitCount: canonBatteryTypePayloadSize - 1,
			want:      "",
		},
		{
			name:      "empty payload after header",
			raw:       withHeader([]byte{0}),
			unitCount: canonBatteryTypePayloadSize,
			want:      "",
		},
		{
			name:      "nul terminated battery type",
			raw:       withHeader([]byte("LP-E6N\x00TRAILING")),
			unitCount: canonBatteryTypePayloadSize,
			want:      "LP-E6N",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			entry := tag.NewEntry(
				tag.ID(canon.BatteryType),
				tag.TypeUndefined,
				tc.unitCount,
				0,
				tag.MakerNoteIFD,
				0,
				utils.LittleEndian,
			)
			r := NewReader(Logger)
			defer r.Close()

			var br bytes.Reader
			br.Reset(tc.raw)
			r.Reset(&br)
			got := r.parseCanonBatteryType(entry)
			if got != tc.want {
				t.Fatalf("got = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseCanonLensModelTerminatedAtNUL(t *testing.T) {
	raw := []byte("EF70-200\x00TRAILING")
	entry := tag.NewEntry(
		tag.ID(canon.LensModel),
		tag.TypeASCII,
		uint32(len(raw)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	if ok := r.parseCanonTag(entry); !ok {
		t.Fatal("parseCanonTag returned false for LensModel")
	}

	if got := r.Exif.MakerNote.Canon.LensModel; got != "EF70-200" {
		t.Fatalf("LensModel = %q, want %q", got, "EF70-200")
	}
}

func TestParseCanonBatteryTypeTag(t *testing.T) {
	raw := make([]byte, canonBatteryTypePayloadSize)
	copy(raw[4:], []byte("LP-E6N\x00TRAILING"))
	entry := tag.NewEntry(
		tag.ID(canon.BatteryType),
		tag.TypeUndefined,
		uint32(len(raw)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	if ok := r.parseCanonTag(entry); !ok {
		t.Fatal("parseCanonTag returned false for BatteryType")
	}

	if got := r.Exif.MakerNote.Canon.BatteryType; got != "LP-E6N" {
		t.Fatalf("BatteryType = %q, want %q", got, "LP-E6N")
	}
}

func TestParseCanonImageUniqueIDFromByte(t *testing.T) {
	raw := []byte{
		0xe4, 0x2d, 0xcf, 0xf1, 0x86, 0x92, 0x4c, 0x33,
		0x97, 0x4d, 0xbf, 0x00, 0xe1, 0x7f, 0x61, 0x62,
	}
	entry := tag.NewEntry(
		tag.ID(canon.ImageUniqueID),
		tag.TypeByte,
		uint32(len(raw)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	if ok := r.parseCanonTag(entry); !ok {
		t.Fatal("parseCanonTag returned false for ImageUniqueID")
	}

	want := meta.UUIDFromString("e42dcff186924c33974dbf00e17f6162")
	if got := r.Exif.MakerNote.Canon.ImageUniqueID; got != want {
		t.Fatalf("ImageUniqueID = %v, want %v", got, want)
	}
}

func TestParseCanonImageUniqueIDAllZeroIsNilUUID(t *testing.T) {
	raw := make([]byte, 16)
	entry := tag.NewEntry(
		tag.ID(canon.ImageUniqueID),
		tag.TypeByte,
		uint32(len(raw)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)
	if ok := r.parseCanonTag(entry); !ok {
		t.Fatal("parseCanonTag returned false for ImageUniqueID")
	}
	if got := r.Exif.MakerNote.Canon.ImageUniqueID; got != meta.NilUUID {
		t.Fatalf("ImageUniqueID = %v, want NilUUID", got)
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
	if parsed.MakerNote.Canon == nil {
		t.Fatalf("MakerNote.Canon = nil for %s", samplePath)
	}
	if parsed.MakerNote.Canon.AFConfig == (canon.AFConfig{}) {
		t.Fatalf("AFConfig not parsed for %s", samplePath)
	}
}

func TestCanonShouldReplaceAFInfoPrefersPopulated(t *testing.T) {
	current := canon.AFInfo{
		Source:        canon.AFInfoSourceAFInfo2,
		NumAFPoints:   1053,
		ValidAFPoints: 651,
		AFArea:        make([]canon.AFPoint, 32),
	}
	candidate := canon.AFInfo{
		Source: canon.AFInfoSourceAFInfo2,
	}

	if canonShouldReplaceAFInfo(current, candidate) {
		t.Fatal("expected populated AFInfo to be retained over empty candidate")
	}
	if !canonShouldReplaceAFInfo(candidate, current) {
		t.Fatal("expected populated AFInfo candidate to replace empty current")
	}
}

func TestCanonShouldReplaceAFInfoPrefersHigherQuality(t *testing.T) {
	current := canon.AFInfo{
		Source:        canon.AFInfoSourceAFInfo2,
		NumAFPoints:   11,
		ValidAFPoints: 11,
		AFArea:        make([]canon.AFPoint, 11),
	}
	candidate := canon.AFInfo{
		Source:        canon.AFInfoSourceAFInfo2,
		NumAFPoints:   1053,
		ValidAFPoints: 1053,
		AFArea:        make([]canon.AFPoint, 1053),
	}

	if !canonShouldReplaceAFInfo(current, candidate) {
		t.Fatal("expected higher-quality AFInfo candidate to replace current")
	}
}

func TestCanonShouldReplaceAFInfoSourceTieBreak(t *testing.T) {
	current := canon.AFInfo{Source: canon.AFInfoSourceAFInfo3}
	candidate := canon.AFInfo{Source: canon.AFInfoSourceAFInfo2}
	if !canonShouldReplaceAFInfo(current, candidate) {
		t.Fatal("expected AFInfo2 to win source-priority tie over AFInfo3")
	}
}

func TestParseCanonAFInfo2LargePayload(t *testing.T) {
	const numAFPoints = 1053

	maskWords := canonBitWordCount(numAFPoints)
	totalWords := 8 + (4 * numAFPoints) + (2 * maskWords)
	words := make([]uint16, totalWords)
	words[1] = 22 // AFAreaMode
	words[2] = numAFPoints
	words[3] = numAFPoints
	words[4] = 6000
	words[5] = 4000
	words[6] = 6000
	words[7] = 4000

	widthStart := 8
	heightStart := widthStart + numAFPoints
	xStart := heightStart + numAFPoints
	yStart := xStart + numAFPoints
	inFocusStart := yStart + numAFPoints
	selectedStart := inFocusStart + maskWords

	for i := 0; i < numAFPoints; i++ {
		words[widthStart+i] = 130
		words[heightStart+i] = 131
		words[xStart+i] = uint16(i)
		words[yStart+i] = uint16(numAFPoints - 1 - i)
	}
	// Set representative bit flags spanning multiple words.
	words[inFocusStart+0] = 1 << 0
	words[inFocusStart+10] = 1 << 1
	words[inFocusStart+65] = 1 << 4

	// EOS-selected bitset.
	words[selectedStart+0] = 1 << 0
	words[selectedStart+65] = 1 << 4

	got := parseAFInfo2ForTest(t, words, "Canon EOS R3", false)
	if got.Source != canon.AFInfoSourceAFInfo2 {
		t.Fatalf("Source = %v, want AFInfoSourceAFInfo2", got.Source)
	}
	if int(got.NumAFPoints) != numAFPoints {
		t.Fatalf("NumAFPoints = %d, want %d", got.NumAFPoints, numAFPoints)
	}
	if len(got.AFArea) != numAFPoints {
		t.Fatalf("len(AFArea) = %d, want %d", len(got.AFArea), numAFPoints)
	}
	if len(got.AFPoints) != numAFPoints {
		t.Fatalf("len(AFPoints) = %d, want %d", len(got.AFPoints), numAFPoints)
	}
	if got.AFArea[0] != canon.NewAFPoint(130, 131, 0, int16(numAFPoints-1)) {
		t.Fatalf("AFArea[0] = %v, unexpected", got.AFArea[0])
	}
	last := got.AFArea[numAFPoints-1]
	if last != canon.NewAFPoint(130, 131, int16(numAFPoints-1), 0) {
		t.Fatalf("AFArea[last] = %v, unexpected", last)
	}
	if !reflect.DeepEqual(got.AFPointsInFocusBits, []int{0, 161, 1044}) {
		t.Fatalf("AFPointsInFocusBits = %v, want [0 161 1044]", got.AFPointsInFocusBits)
	}
	if !reflect.DeepEqual(got.AFPointsSelectedBits, []int{0, 1044}) {
		t.Fatalf("AFPointsSelectedBits = %v, want [0 1044]", got.AFPointsSelectedBits)
	}
}

func TestParseCanonAFInfoLegacySamples(t *testing.T) {
	tests := []struct {
		file string
		want canon.AFInfo
	}{
		{
			file: filepath.Join(defaultBenchImageDir, "350D.CR2"),
			want: canon.AFInfo{
				Source:              canon.AFInfoSourceAFInfo,
				NumAFPoints:         7,
				ValidAFPoints:       7,
				CanonImageWidth:     3456,
				CanonImageHeight:    2304,
				AFImageWidth:        3456,
				AFImageHeight:       2304,
				AFAreaWidth:         189,
				AFAreaHeight:        188,
				AFPointsInFocusBits: []int{3},
				AFArea: []canon.AFPoint{
					canon.NewAFPoint(189, 188, 0, -617),
					canon.NewAFPoint(189, 188, -1237, 0),
					canon.NewAFPoint(189, 188, -742, 0),
					canon.NewAFPoint(189, 188, 0, 0),
					canon.NewAFPoint(189, 188, 742, 0),
					canon.NewAFPoint(189, 188, 1237, 0),
					canon.NewAFPoint(189, 188, 0, 617),
				},
			},
		},
		{
			file: filepath.Join(defaultBenchImageDir, "XT1.CR2"),
			want: canon.AFInfo{
				Source:              canon.AFInfoSourceAFInfo,
				NumAFPoints:         7,
				ValidAFPoints:       7,
				CanonImageWidth:     3456,
				CanonImageHeight:    2304,
				AFImageWidth:        3456,
				AFImageHeight:       2304,
				AFAreaWidth:         189,
				AFAreaHeight:        188,
				AFPointsInFocusBits: []int{3},
				AFArea: []canon.AFPoint{
					canon.NewAFPoint(189, 188, 0, -617),
					canon.NewAFPoint(189, 188, -1237, 0),
					canon.NewAFPoint(189, 188, -742, 0),
					canon.NewAFPoint(189, 188, 0, 0),
					canon.NewAFPoint(189, 188, 742, 0),
					canon.NewAFPoint(189, 188, 1237, 0),
					canon.NewAFPoint(189, 188, 0, 617),
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(filepath.Base(tc.file), func(t *testing.T) {
			if _, err := os.Stat(tc.file); err != nil {
				t.Skipf("sample %s unavailable: %v", tc.file, err)
			}

			f, err := os.Open(tc.file)
			if err != nil {
				t.Fatalf("open %s: %v", tc.file, err)
			}
			defer func() {
				_ = f.Close()
			}()

			parsed, err := Parse(f)
			if err != nil {
				t.Fatalf("parse %s: %v", tc.file, err)
			}

			got := parsed.MakerNote.Canon.AFInfo
			if got.Source != tc.want.Source {
				t.Fatalf("Source = %v, want %v", got.Source, tc.want.Source)
			}
			if got.NumAFPoints != tc.want.NumAFPoints {
				t.Fatalf("NumAFPoints = %d, want %d", got.NumAFPoints, tc.want.NumAFPoints)
			}
			if got.ValidAFPoints != tc.want.ValidAFPoints {
				t.Fatalf("ValidAFPoints = %d, want %d", got.ValidAFPoints, tc.want.ValidAFPoints)
			}
			if got.CanonImageWidth != tc.want.CanonImageWidth || got.CanonImageHeight != tc.want.CanonImageHeight {
				t.Fatalf(
					"CanonImage = %dx%d, want %dx%d",
					got.CanonImageWidth,
					got.CanonImageHeight,
					tc.want.CanonImageWidth,
					tc.want.CanonImageHeight,
				)
			}
			if got.AFImageWidth != tc.want.AFImageWidth || got.AFImageHeight != tc.want.AFImageHeight {
				t.Fatalf(
					"AFImage = %dx%d, want %dx%d",
					got.AFImageWidth,
					got.AFImageHeight,
					tc.want.AFImageWidth,
					tc.want.AFImageHeight,
				)
			}
			if got.AFAreaWidth != tc.want.AFAreaWidth || got.AFAreaHeight != tc.want.AFAreaHeight {
				t.Fatalf(
					"AFAreaSize = %dx%d, want %dx%d",
					got.AFAreaWidth,
					got.AFAreaHeight,
					tc.want.AFAreaWidth,
					tc.want.AFAreaHeight,
				)
			}
			if !reflect.DeepEqual(got.AFPointsInFocusBits, tc.want.AFPointsInFocusBits) {
				t.Fatalf("AFPointsInFocusBits = %v, want %v", got.AFPointsInFocusBits, tc.want.AFPointsInFocusBits)
			}
			if !reflect.DeepEqual(got.AFArea, tc.want.AFArea) {
				t.Fatalf("AFArea = %v, want %v", got.AFArea, tc.want.AFArea)
			}
			if !reflect.DeepEqual(got.AFPoints, tc.want.AFArea) {
				t.Fatalf("AFPoints = %v, want %v", got.AFPoints, tc.want.AFArea)
			}
		})
	}
}

func TestParseCanonMaxAperture(t *testing.T) {
	tests := []struct {
		name string
		raw  uint16
		want float64
	}{
		{name: "zero", raw: 0, want: 0},
		{name: "negative code", raw: 0xffff, want: 0},
		{name: "ev one", raw: 0x20, want: math.Exp2(0.5)},
		{name: "ev two", raw: 0x40, want: 2},
		{name: "canon one third", raw: 0x0c, want: math.Exp2((1.0 / 3.0) * 0.5)},
		{name: "canon two thirds", raw: 0x14, want: math.Exp2((2.0 / 3.0) * 0.5)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := float64(parseCanonMaxAperture(tc.raw))
			if math.Abs(got-tc.want) > 1e-5 {
				t.Fatalf("parseCanonMaxAperture(%#x)=%.8f want %.8f", tc.raw, got, tc.want)
			}
		})
	}
}

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
	input := make([]uint16, 53)
	input[0] = uint16(len(input) * 2) // ExifTool Validate(): first word is byte length.
	copy(input[1:], unsigned[:52])
	raw := canonUint16WordsToBytes(input, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(input)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)

	if got.MacroMode != canon.MacroMode(signed[0]) {
		t.Fatalf("MacroMode = %d, want %d", got.MacroMode, canon.MacroMode(signed[0]))
	}
	if got.FocusMode != canon.FocusMode(signed[6]) {
		t.Fatalf("FocusMode = %d, want %d", got.FocusMode, signed[6])
	}
	if got.AFPoint != unsigned[18] {
		t.Fatalf("AFPoint = %#x, want %#x", got.AFPoint, unsigned[18])
	}
	if got.CanonExposureMode != canon.ExposureMode(signed[19]) {
		t.Fatalf("CanonExposureMode = %d, want %d", got.CanonExposureMode, signed[19])
	}
	if got.LensType != canon.CanonLensType(unsigned[21]) {
		t.Fatalf("LensType = %#x, want %#x", got.LensType, canon.CanonLensType(unsigned[21]))
	}
	if got.DisplayAperture != parseCanonDisplayAperture(unsigned[34]) {
		t.Fatalf("DisplayAperture = %v, want %v", got.DisplayAperture, parseCanonDisplayAperture(unsigned[34]))
	}
	if got.SRAWQuality != canon.SRAWQuality(signed[45]) {
		t.Fatalf("SRAWQuality = %d, want %d", got.SRAWQuality, canon.SRAWQuality(signed[45]))
	}
	if got.FocusBracketing != canon.FocusBracketing(signed[49]) {
		t.Fatalf("FocusBracketing = %d, want %d", got.FocusBracketing, canon.FocusBracketing(signed[49]))
	}
	if got.Clarity != signed[50] {
		t.Fatalf("Clarity = %d, want %d", got.Clarity, signed[50])
	}
	if got.HDRPQ != canon.HDRPQ(signed[51]) {
		t.Fatalf("HDRPQ = %d, want %d", got.HDRPQ, canon.HDRPQ(signed[51]))
	}
}

func TestParseCanonCameraSettingsShortInputProgressive(t *testing.T) {
	// CameraSettings stores a byte-size header first, followed by ExifTool
	// sequence entries 1..N.
	words := []uint16{8, 123, 456, 789}
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)
	if got.MacroMode != canon.MacroMode(123) {
		t.Fatalf("MacroMode = %d, want 123", got.MacroMode)
	}
	if got.SelfTimer != 456 {
		t.Fatalf("SelfTimer = %d, want 456", got.SelfTimer)
	}
	if got.Quality != canon.Quality(789) {
		t.Fatalf("Quality = %d, want 789", got.Quality)
	}
	if got.HDRPQ != 0 {
		t.Fatalf("HDRPQ = %d, want 0", got.HDRPQ)
	}
}

func TestParseCanonCameraSettingsApertureConversions(t *testing.T) {
	words := make([]uint16, 36)
	words[0] = uint16(len(words) * 2) // leading CameraSettings byte-size header
	words[1] = 2                      // [1] MacroMode
	words[26] = 0x40                  // [26] MaxAperture => f/2
	words[27] = 0x20                  // [27] MinAperture => f/1.41
	words[35] = 45                    // [35] DisplayAperture => 4.5

	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)
	if math.Abs(float64(got.MaxAperture)-2.0) > 1e-5 {
		t.Fatalf("MaxAperture = %.5f, want 2.0", got.MaxAperture)
	}
	if math.Abs(float64(got.MinAperture)-math.Exp2(0.5)) > 1e-5 {
		t.Fatalf("MinAperture = %.5f, want %.5f", got.MinAperture, math.Exp2(0.5))
	}
	if got.DisplayAperture != 4.5 {
		t.Fatalf("DisplayAperture = %.2f, want 4.50", got.DisplayAperture)
	}
}

func TestParseCanonCameraSettingsInvalidLengthReturnsZero(t *testing.T) {
	words := []uint16{123, 456, 789}
	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonCameraSettings),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()
	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonCameraSettings(entry)
	if got != (canon.CameraSettings{}) {
		t.Fatalf("parseCanonCameraSettings(invalid length) = %+v, want zero value", got)
	}
}

func TestParseCanonShotInfo(t *testing.T) {
	words := make([]uint16, 34)
	words[0] = uint16(len(words) * 2) // leading ShotInfo byte-size header
	words[1] = 0                      // [1] AutoISO => 100
	words[2] = 224                    // [2] BaseISO => 400
	words[3] = 96                     // [3] MeasuredEV
	words[4] = 0x40                   // [4] TargetAperture
	words[5] = 0x60                   // [5] TargetExposureTime
	words[6] = 32                     // [6] ExposureCompensation
	words[7] = uint16(canon.WhiteBalanceShade)
	words[8] = uint16(canon.SlowShutterNightScene)
	words[9] = 7
	words[10] = 3
	words[12] = 150
	words[13] = 64
	words[14] = 0x11
	words[15] = 5
	words[16] = 2
	words[17] = 1
	words[18] = 4
	words[19] = 250
	words[20] = 400
	words[21] = 0x40
	words[22] = 0x60
	words[23] = 64
	words[24] = 9
	words[26] = uint16(canon.CameraTypeEOSHighEnd)
	words[27] = uint16(canon.AutoRotateRotate90CW)
	words[28] = uint16(canon.NDFilterOn)
	words[29] = 15
	words[33] = 12

	got := parseShotInfoForTest(t, words, "Canon EOS 5D Mark IV", 10)

	if math.Abs(float64(got.AutoISOValue-100.0)) > 1e-5 {
		t.Fatalf("AutoISOValue = %.5f, want 100.0", got.AutoISOValue)
	}
	if math.Abs(float64(got.BaseISOValue-400.0)) > 1e-5 {
		t.Fatalf("BaseISOValue = %.5f, want 400.0", got.BaseISOValue)
	}
	if math.Abs(float64(got.ActualISO-400.0)) > 1e-5 {
		t.Fatalf("ActualISO = %.5f, want 400.0", got.ActualISO)
	}
	if got.WhiteBalance != canon.WhiteBalanceShade {
		t.Fatalf("WhiteBalance = %d, want %d", got.WhiteBalance, canon.WhiteBalanceShade)
	}
	if got.SlowShutter != canon.SlowShutterNightScene {
		t.Fatalf("SlowShutter = %d, want %d", got.SlowShutter, canon.SlowShutterNightScene)
	}
	if got.TargetApertureValue != canonShotAperture(int16(words[4])) {
		t.Fatalf("TargetApertureValue = %v, want %v", got.TargetApertureValue, canonShotAperture(int16(words[4])))
	}
	if got.TargetExposureTimeValue != canonShotExposureTime(int16(words[5]), false) {
		t.Fatalf("TargetExposureTimeValue = %v, want %v", got.TargetExposureTimeValue, canonShotExposureTime(int16(words[5]), false))
	}
	if got.CameraTemperatureC != 22 {
		t.Fatalf("CameraTemperatureC = %d, want 22", got.CameraTemperatureC)
	}
	if math.Abs(float64(got.FlashGuideNumberMeters-2.0)) > 1e-6 {
		t.Fatalf("FlashGuideNumberMeters = %.6f, want 2.0", got.FlashGuideNumberMeters)
	}
	if got.FNumberValue != canonShotAperture(int16(words[21])) {
		t.Fatalf("FNumberValue = %v, want %v", got.FNumberValue, canonShotAperture(int16(words[21])))
	}
	if got.ExposureTimeValue != canonShotExposureTime(int16(words[22]), false) {
		t.Fatalf("ExposureTimeValue = %v, want %v", got.ExposureTimeValue, canonShotExposureTime(int16(words[22]), false))
	}
	if got.CameraType != canon.CameraTypeEOSHighEnd {
		t.Fatalf("CameraType = %d, want %d", got.CameraType, canon.CameraTypeEOSHighEnd)
	}
	if got.AutoRotate != canon.AutoRotateRotate90CW {
		t.Fatalf("AutoRotate = %d, want %d", got.AutoRotate, canon.AutoRotateRotate90CW)
	}
	if got.NDFilter != canon.NDFilterOn {
		t.Fatalf("NDFilter = %d, want %d", got.NDFilter, canon.NDFilterOn)
	}
	if got.FocusDistance != canon.NewFocusDistance(250, 400) {
		t.Fatalf("FocusDistance = %v, want %v", got.FocusDistance, canon.NewFocusDistance(250, 400))
	}
	if math.Abs(float64(got.FocusDistance.UpperMeters()-2.5)) > 1e-6 {
		t.Fatalf("FocusDistance.UpperMeters = %.6f, want 2.5", got.FocusDistance.UpperMeters())
	}
	if math.Abs(float64(got.FocusDistance.LowerMeters()-4.0)) > 1e-6 {
		t.Fatalf("FocusDistance.LowerMeters = %.6f, want 4.0", got.FocusDistance.LowerMeters())
	}
}

func TestParseCanonShotInfoTruncated(t *testing.T) {
	words := []uint16{12, 0, 224, 96, 0x40, 0x60}
	got := parseShotInfoForTest(t, words, "Canon EOS R6", 10)

	if math.Abs(float64(got.ActualISO-400.0)) > 1e-5 {
		t.Fatalf("ActualISO = %.5f, want 400.0", got.ActualISO)
	}
	if got.TargetApertureValue != canonShotAperture(int16(words[4])) {
		t.Fatalf("TargetApertureValue = %v, want %v", got.TargetApertureValue, canonShotAperture(int16(words[4])))
	}
	if got.TargetExposureTimeValue != canonShotExposureTime(int16(words[5]), false) {
		t.Fatalf("TargetExposureTimeValue = %v, want %v", got.TargetExposureTimeValue, canonShotExposureTime(int16(words[5]), false))
	}
	if got.WhiteBalance != 0 {
		t.Fatalf("WhiteBalance = %d, want 0", got.WhiteBalance)
	}
	if got.FocusDistance != (canon.FocusDistance{}) {
		t.Fatalf("FocusDistance = %v, want zero", got.FocusDistance)
	}
	if got.FNumberValue != 0 {
		t.Fatalf("FNumberValue = %v, want 0", got.FNumberValue)
	}
}

func TestParseCanonShotInfoSkipsZeroFocusDistanceUpper(t *testing.T) {
	words := make([]uint16, 21)
	words[0] = uint16(len(words) * 2)
	words[19] = 0
	words[20] = 400

	got := parseShotInfoForTest(t, words, "Canon EOS R5", 10)
	if got.FocusDistance != (canon.FocusDistance{}) {
		t.Fatalf("FocusDistance = %v, want zero", got.FocusDistance)
	}
}

func TestParseCanonShotInfoLegacyExposureTime(t *testing.T) {
	words := make([]uint16, 23)
	words[0] = uint16(len(words) * 2)
	words[22] = 672 // [22] ExposureTime

	got := parseShotInfoForTest(t, words, "Canon EOS 20D", 10)
	want := canonShotExposureTime(int16(words[22]), true)
	if got.ExposureTimeValue != want {
		t.Fatalf("ExposureTimeValue = %v, want %v", got.ExposureTimeValue, want)
	}
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
			city := canon.TimeZoneCityNewYork
			daylightSavings := canon.DaylightSavingsOn

			words := []uint32{16, uint32(tz), uint32(city), uint32(daylightSavings)}
			raw := canonUint32WordsToBytes(words, tc.bo)

			entry := tag.NewEntry(
				tag.ID(canon.TimeInfo),
				tag.TypeLong,
				uint32(len(words)),
				0,
				tag.MakerNoteIFD,
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

func TestParseCanonFileInfoIndexMapping(t *testing.T) {
	words := make([]uint16, 62)
	words[0] = uint16(len(words) * 2)
	words[1] = 0x1234
	words[2] = 0x5678
	words[3] = uint16(canon.BracketModeWB)
	words[4] = u16(-2)
	words[5] = 3
	words[6] = uint16(canon.RawJpgQualityCRAW)
	words[7] = uint16(canon.RawJpgSizeMedium3)
	words[8] = uint16(canon.OnOffAutoAuto)
	words[9] = 2
	words[12] = u16(-1)
	words[13] = 4
	words[14] = uint16(canon.FilterEffectRed)
	words[15] = 0xffff
	words[16] = 75
	words[19] = uint16(canon.OnOffAutoOn)
	words[20] = 299
	words[21] = 267
	words[23] = uint16(canon.ShutterModeElectronic)
	words[25] = uint16(canon.OnOffAutoOff)
	words[32] = uint16(canon.OnOffAutoOff)
	words[61] = uint16(canon.CanonRFLensType(258))

	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonFileInfo),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonFileInfo(entry)
	if got.FileNumber != 0x56781234 {
		t.Fatalf("FileNumber = %#x, want %#x", got.FileNumber, uint32(0x56781234))
	}
	if got.BracketMode != canon.BracketModeWB {
		t.Fatalf("BracketMode = %d, want %d", got.BracketMode, canon.BracketModeWB)
	}
	if got.BracketValue != -2 {
		t.Fatalf("BracketValue = %d, want -2", got.BracketValue)
	}
	if got.RawJpgQuality != canon.RawJpgQualityCRAW {
		t.Fatalf("RawJpgQuality = %d, want %d", got.RawJpgQuality, canon.RawJpgQualityCRAW)
	}
	if got.RawJpgSize != canon.RawJpgSizeMedium3 {
		t.Fatalf("RawJpgSize = %d, want %d", got.RawJpgSize, canon.RawJpgSizeMedium3)
	}
	if got.LongExposureNoiseReduction2 != canon.OnOffAutoAuto {
		t.Fatalf("LongExposureNoiseReduction2 = %d, want %d", got.LongExposureNoiseReduction2, canon.OnOffAutoAuto)
	}
	if got.WBBracketMode != 2 {
		t.Fatalf("WBBracketMode = %d, want 2", got.WBBracketMode)
	}
	if got.WBBracketValueAB != -1 {
		t.Fatalf("WBBracketValueAB = %d, want -1", got.WBBracketValueAB)
	}
	if got.WBBracketValueGM != 4 {
		t.Fatalf("WBBracketValueGM = %d, want 4", got.WBBracketValueGM)
	}
	if got.FilterEffect != canon.FilterEffectRed {
		t.Fatalf("FilterEffect = %d, want %d", got.FilterEffect, canon.FilterEffectRed)
	}
	if got.ToningEffect != canon.ToningEffect(0xffff) {
		t.Fatalf("ToningEffect = %#x, want %#x", got.ToningEffect, uint16(0xffff))
	}
	if got.MacroMagnification != 75 {
		t.Fatalf("MacroMagnification = %d, want 75", got.MacroMagnification)
	}
	if got.LiveViewShooting != canon.OnOffAutoOn {
		t.Fatalf("LiveViewShooting = %d, want %d", got.LiveViewShooting, canon.OnOffAutoOn)
	}
	if got.FocusDistance != canon.NewFocusDistance(299, 267) {
		t.Fatalf("FocusDistance = %v, want %v", got.FocusDistance, canon.NewFocusDistance(299, 267))
	}
	if got.ShutterMode != canon.ShutterModeElectronic {
		t.Fatalf("ShutterMode = %d, want %d", got.ShutterMode, canon.ShutterModeElectronic)
	}
	if got.AntiFlicker != canon.OnOffAutoOff {
		t.Fatalf("AntiFlicker = %d, want %d", got.AntiFlicker, canon.OnOffAutoOff)
	}
	if got.RFLensType != canon.CanonRFLensType(258) {
		t.Fatalf("RFLensType = %d, want 258", got.RFLensType)
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
		tag.ID(canon.FaceDetect1),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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
	if got.FacePositions[0] != (canon.FacePosition{X: -10, Y: 20}) {
		t.Fatalf("FacePositions[0] = %+v, want {-10, 20}", got.FacePositions[0])
	}
	if got.FacePositions[1] != (canon.FacePosition{X: 30, Y: -40}) {
		t.Fatalf("FacePositions[1] = %+v, want {30, -40}", got.FacePositions[1])
	}
}

func TestParseCanonFaceDetect2IndexMapping(t *testing.T) {
	raw := []byte{0, 42, 3, 0, 0}
	entry := tag.NewEntry(
		tag.ID(canon.FaceDetect2),
		tag.TypeByte,
		uint32(len(raw)),
		0,
		tag.MakerNoteIFD,
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
		tag.ID(canon.FaceDetect3),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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
	words := make([]uint32, 12)
	words[0] = uint32(len(words) * 4)
	words[1] = 1  // PeripheralIlluminationCorr
	words[2] = 2  // AutoLightingOptimizer
	words[3] = 1  // HighlightTonePriority
	words[4] = 2  // LongExposureNoiseReduction
	words[5] = 3  // HighISONoiseReduction
	words[10] = 2 // DigitalLensOptimizer
	words[11] = 1 // DualPixelRaw

	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonLightingOpt),
		tag.TypeLong,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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
	words := []uint32{16, 1, 3, 1}
	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonLightingOpt),
		tag.TypeLong,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
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

func TestParseCanonMultiExpIndexMapping(t *testing.T) {
	words := []uint32{16, 0, 3, 5}
	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonMultiExp),
		tag.TypeLong,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonMultiExp(entry)
	if got.MultiExposure != 0 {
		t.Fatalf("MultiExposure = %d, want 0", got.MultiExposure)
	}
	if got.MultiExposureControl != 3 {
		t.Fatalf("MultiExposureControl = %d, want 3", got.MultiExposureControl)
	}
	if got.MultiExposureShots != 5 {
		t.Fatalf("MultiExposureShots = %d, want 5", got.MultiExposureShots)
	}
}

func TestParseCanonHDRInfoIndexMapping(t *testing.T) {
	words := []uint32{12, 0, 2}
	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonHDRInfo),
		tag.TypeLong,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonHDRInfo(entry)
	if got.HDR != 0 {
		t.Fatalf("HDR = %d, want 0", got.HDR)
	}
	if got.HDREffect != 2 {
		t.Fatalf("HDREffect = %d, want 2", got.HDREffect)
	}
}

func TestParseCanonProcessingInfoIndexMapping(t *testing.T) {
	words := []uint16{
		28,
		0,    // ToneCurve
		2,    // Sharpness
		0,    // SharpnessFrequency
		10,   // SensorRedLevel
		11,   // SensorBlueLevel
		12,   // WhiteBalanceRed
		13,   // WhiteBalanceBlue
		14,   // WhiteBalance
		3100, // ColorTemperature
		130,  // PictureStyle
		0,    // DigitalGain
		1,    // WBShiftAB
		2,    // WBShiftGM
		3,    // UnsharpMaskFineness
		4,    // UnsharpMaskThreshold
	}

	raw := canonUint16WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonProcessingInfo),
		tag.TypeShort,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonProcessingInfo(entry)
	if got.ToneCurve != 0 {
		t.Fatalf("ToneCurve = %d, want 0", got.ToneCurve)
	}
	if got.Sharpness != 2 {
		t.Fatalf("Sharpness = %d, want 2", got.Sharpness)
	}
	if got.SharpnessFrequency != 0 {
		t.Fatalf("SharpnessFrequency = %d, want 0", got.SharpnessFrequency)
	}
	if got.ColorTemperature != 3100 {
		t.Fatalf("ColorTemperature = %d, want 3100", got.ColorTemperature)
	}
	if got.PictureStyle != 130 {
		t.Fatalf("PictureStyle = %d, want 130", got.PictureStyle)
	}
	if got.DigitalGain != 0 {
		t.Fatalf("DigitalGain = %d, want 0", got.DigitalGain)
	}
	if got.UnsharpMaskThreshold != 4 {
		t.Fatalf("UnsharpMaskThreshold = %d, want 4", got.UnsharpMaskThreshold)
	}
}

func TestParseCanonAFMicroAdjIndexMapping(t *testing.T) {
	words := []uint32{
		16,
		0,
		7,
		3,
	}

	raw := canonUint32WordsToBytes(words, utils.LittleEndian)
	entry := tag.NewEntry(
		tag.ID(canon.CanonAFMicroAdj),
		tag.TypeLong,
		uint32(len(words)),
		0,
		tag.MakerNoteIFD,
		0,
		utils.LittleEndian,
	)

	r := NewReader(Logger)
	defer r.Close()

	var br bytes.Reader
	br.Reset(raw)
	r.Reset(&br)

	got := r.parseCanonAFMicroAdj(entry)
	if got.Mode != 0 {
		t.Fatalf("Mode = %d, want 0", got.Mode)
	}
	if got.ValueNumerator != 7 {
		t.Fatalf("ValueNumerator = %d, want 7", got.ValueNumerator)
	}
	if got.ValueDenominator != 3 {
		t.Fatalf("ValueDenominator = %d, want 3", got.ValueDenominator)
	}
}
