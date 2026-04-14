package exif

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/makernote/nikon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func parseNikonBlockForTest(t *testing.T, tagID tag.ID, raw []byte, model string) *nikon.Nikon {
	t.Helper()

	entry := tag.NewEntry(
		tagID,
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
	r.Exif.IFD0.Model = model
	if ok := r.parseNikonTag(entry); !ok {
		t.Fatalf("parseNikonTag returned false for 0x%04x", uint16(tagID))
	}
	return r.nikonMakerNote()
}

func TestParseNikonISOInfoUsesByteOffsets(t *testing.T) {
	raw := []byte{
		0x54, 0x01, 0x0c, 0x00,
		0x00, 0x00,
		0x54, 0x01, 0x0c, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}

	got := parseNikonBlockForTest(t, tag.ID(nikon.ISOInfo), raw, "NIKON D300S").ISOInfo
	if math.IsInf(got.ISO, 0) || math.IsInf(got.ISO2, 0) {
		t.Fatalf("ISOInfo produced non-finite ISO values: %+v", got)
	}
	if got.ISO != 400 {
		t.Fatalf("ISO = %v, want 400", got.ISO)
	}
	if got.ISO2 != 400 {
		t.Fatalf("ISO2 = %v, want 400", got.ISO2)
	}
	if got.ISOExpansion != 0 || got.ISOExpansion2 != 0 {
		t.Fatalf("unexpected ISO expansion values: %+v", got)
	}
}

func TestParseNikonFileInfoLittleEndianHeuristic(t *testing.T) {
	raw := []byte{
		'0', '1', '0', '0',
		0x00, 0x00,
		0x64, 0x00,
		0x7c, 0x10,
	}

	got := parseNikonBlockForTest(t, tag.ID(nikon.FileInfo), raw, "NIKON Z 9").FileInfo
	if got.FileInfoVersion != "0100" {
		t.Fatalf("FileInfoVersion = %q, want 0100", got.FileInfoVersion)
	}
	if got.MemoryCardNumber != 0 {
		t.Fatalf("MemoryCardNumber = %d, want 0", got.MemoryCardNumber)
	}
	if got.DirectoryNumber != 100 {
		t.Fatalf("DirectoryNumber = %d, want 100", got.DirectoryNumber)
	}
	if got.FileNumber != 4220 {
		t.Fatalf("FileNumber = %d, want 4220", got.FileNumber)
	}
}

func TestParseNikonAFInfo2V0400(t *testing.T) {
	raw := make([]byte, 0x4b)
	copy(raw[:4], []byte("0400"))
	raw[4] = 2
	raw[5] = 204
	raw[7] = 1
	utils.LittleEndian.PutUint16(raw[0x3e:0x40], 8256)
	utils.LittleEndian.PutUint16(raw[0x40:0x42], 5504)
	utils.LittleEndian.PutUint16(raw[0x42:0x44], 4128)
	utils.LittleEndian.PutUint16(raw[0x44:0x46], 3043)
	utils.LittleEndian.PutUint16(raw[0x46:0x48], 291)
	utils.LittleEndian.PutUint16(raw[0x48:0x4a], 323)
	raw[0x4a] = 1

	got := parseNikonBlockForTest(t, tag.ID(nikon.AFInfo2), raw, "NIKON Z 9").AFInfo2
	if got.AFInfo2Version != "0400" {
		t.Fatalf("AFInfo2Version = %q, want 0400", got.AFInfo2Version)
	}
	if got.AFDetectionMethod != 2 {
		t.Fatalf("AFDetectionMethod = %d, want 2", got.AFDetectionMethod)
	}
	if got.AFAreaMode != 204 {
		t.Fatalf("AFAreaMode = %d, want 204", got.AFAreaMode)
	}
	if got.AFCoordinatesAvailable != 1 {
		t.Fatalf("AFCoordinatesAvailable = %d, want 1", got.AFCoordinatesAvailable)
	}
	if got.AFImageWidth != 8256 || got.AFImageHeight != 5504 {
		t.Fatalf("AFImage = %dx%d, want 8256x5504", got.AFImageWidth, got.AFImageHeight)
	}
	if got.AFAreaXPosition != 4128 || got.AFAreaYPosition != 3043 {
		t.Fatalf("AFArea position = (%d,%d), want (4128,3043)", got.AFAreaXPosition, got.AFAreaYPosition)
	}
	if got.AFAreaWidth != 291 || got.AFAreaHeight != 323 {
		t.Fatalf("AFArea size = %dx%d, want 291x323", got.AFAreaWidth, got.AFAreaHeight)
	}
	if got.FocusResult != 1 {
		t.Fatalf("FocusResult = %d, want 1", got.FocusResult)
	}
}

func TestParseNikonMakerNoteSamples(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	cases := []struct {
		file            string
		version         string
		quality         string
		whiteBalance    string
		focusMode       string
		flashSetting    string
		serialNumber    string
		lens            string
		shutterMode     uint16
		afInfo2Version  string
		afDetection     uint8
		afAreaMode      uint8
		afCoordAvail    uint8
		afImageWidth    uint16
		afImageHeight   uint16
		afAreaX         uint16
		afAreaY         uint16
		directoryNumber uint16
		fileNumber      uint16
	}{
		{
			file:            "NikonZ9.NEF",
			version:         "0211",
			quality:         "RAW",
			whiteBalance:    "NATURAL AUTO",
			focusMode:       "AF-C",
			flashSetting:    "NORMAL",
			serialNumber:    "3002822",
			lens:            "100 400 4.5 5.6",
			shutterMode:     16,
			afInfo2Version:  "0400",
			afDetection:     2,
			afAreaMode:      204,
			afCoordAvail:    1,
			afImageWidth:    8256,
			afImageHeight:   5504,
			afAreaX:         4128,
			afAreaY:         3043,
			directoryNumber: 100,
			fileNumber:      4220,
		},
		{
			file:            "NikonZ50II.NEF",
			version:         "0211",
			quality:         "RAW",
			whiteBalance:    "SUNNY",
			focusMode:       "AF-C",
			flashSetting:    "NORMAL",
			serialNumber:    "3002943",
			lens:            "16 50 3.5 6.3",
			shutterMode:     81,
			afInfo2Version:  "0402",
			afDetection:     2,
			afAreaMode:      207,
			afCoordAvail:    1,
			afImageWidth:    5568,
			afImageHeight:   3712,
			afAreaX:         2256,
			afAreaY:         2741,
			directoryNumber: 100,
			fileNumber:      422,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			samplePath := filepath.Join(benchDir, tc.file)
			if _, err := os.Stat(samplePath); err != nil {
				t.Skipf("sample not found: %s", samplePath)
			}

			f, err := os.Open(samplePath)
			if err != nil {
				t.Fatalf("open %s: %v", samplePath, err)
			}
			defer func() { _ = f.Close() }()

			parsed, err := Parse(f)
			if err != nil {
				t.Fatalf("parse %s: %v", samplePath, err)
			}
			if parsed.MakerNote.Nikon == nil {
				t.Fatalf("Nikon maker-note missing for %s", samplePath)
			}

			got := parsed.MakerNote.Nikon
			if got.MakerNoteVersion != tc.version {
				t.Fatalf("MakerNoteVersion = %q, want %q", got.MakerNoteVersion, tc.version)
			}
			if got.Quality != tc.quality {
				t.Fatalf("Quality = %q, want %q", got.Quality, tc.quality)
			}
			if got.WhiteBalance != tc.whiteBalance {
				t.Fatalf("WhiteBalance = %q, want %q", got.WhiteBalance, tc.whiteBalance)
			}
			if got.FocusMode != tc.focusMode {
				t.Fatalf("FocusMode = %q, want %q", got.FocusMode, tc.focusMode)
			}
			if got.FlashSetting != tc.flashSetting {
				t.Fatalf("FlashSetting = %q, want %q", got.FlashSetting, tc.flashSetting)
			}
			if got.SerialNumber != tc.serialNumber {
				t.Fatalf("SerialNumber = %q, want %q", got.SerialNumber, tc.serialNumber)
			}
			if got.Lens != tc.lens {
				t.Fatalf("Lens = %q, want %q", got.Lens, tc.lens)
			}
			if got.ShutterMode != tc.shutterMode {
				t.Fatalf("ShutterMode = %d, want %d", got.ShutterMode, tc.shutterMode)
			}
			if got.AFInfo2.AFInfo2Version != tc.afInfo2Version {
				t.Fatalf("AFInfo2Version = %q, want %q", got.AFInfo2.AFInfo2Version, tc.afInfo2Version)
			}
			if got.AFInfo2.AFDetectionMethod != tc.afDetection {
				t.Fatalf("AFDetectionMethod = %d, want %d", got.AFInfo2.AFDetectionMethod, tc.afDetection)
			}
			if got.AFInfo2.AFAreaMode != tc.afAreaMode {
				t.Fatalf("AFAreaMode = %d, want %d", got.AFInfo2.AFAreaMode, tc.afAreaMode)
			}
			if got.AFInfo2.AFCoordinatesAvailable != tc.afCoordAvail {
				t.Fatalf("AFCoordinatesAvailable = %d, want %d", got.AFInfo2.AFCoordinatesAvailable, tc.afCoordAvail)
			}
			if got.AFInfo2.AFImageWidth != tc.afImageWidth || got.AFInfo2.AFImageHeight != tc.afImageHeight {
				t.Fatalf("AFImage = %dx%d, want %dx%d", got.AFInfo2.AFImageWidth, got.AFInfo2.AFImageHeight, tc.afImageWidth, tc.afImageHeight)
			}
			if got.AFInfo2.AFAreaXPosition != tc.afAreaX || got.AFInfo2.AFAreaYPosition != tc.afAreaY {
				t.Fatalf("AFArea position = (%d,%d), want (%d,%d)", got.AFInfo2.AFAreaXPosition, got.AFInfo2.AFAreaYPosition, tc.afAreaX, tc.afAreaY)
			}
			if got.FileInfo.DirectoryNumber != tc.directoryNumber {
				t.Fatalf("DirectoryNumber = %d, want %d", got.FileInfo.DirectoryNumber, tc.directoryNumber)
			}
			if got.FileInfo.FileNumber != tc.fileNumber {
				t.Fatalf("FileNumber = %d, want %d", got.FileInfo.FileNumber, tc.fileNumber)
			}
		})
	}
}

func TestParseNikonLegacyISOInfoSamples(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	cases := []struct {
		file string
		iso  float64
		iso2 float64
	}{
		{file: "1.NEF", iso: 400, iso2: 400},
		{file: "2.NEF", iso: 3200, iso2: 3200},
		{file: "3.NEF", iso: 100, iso2: 100},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			samplePath := filepath.Join(benchDir, tc.file)
			if _, err := os.Stat(samplePath); err != nil {
				t.Skipf("sample not found: %s", samplePath)
			}

			f, err := os.Open(samplePath)
			if err != nil {
				t.Fatalf("open %s: %v", samplePath, err)
			}
			defer func() { _ = f.Close() }()

			parsed, err := Parse(f)
			if err != nil {
				t.Fatalf("parse %s: %v", samplePath, err)
			}
			if parsed.MakerNote.Nikon == nil {
				t.Fatalf("Nikon maker-note missing for %s", samplePath)
			}

			got := parsed.MakerNote.Nikon.ISOInfo
			if math.IsInf(got.ISO, 0) || math.IsInf(got.ISO2, 0) {
				t.Fatalf("ISOInfo contains non-finite ISO values: %+v", got)
			}
			if got.ISO != tc.iso {
				t.Fatalf("ISO = %v, want %v", got.ISO, tc.iso)
			}
			if got.ISO2 != tc.iso2 {
				t.Fatalf("ISO2 = %v, want %v", got.ISO2, tc.iso2)
			}
		})
	}
}
