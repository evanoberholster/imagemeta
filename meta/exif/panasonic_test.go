package exif_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta"
)

const defaultPanasonicBenchImageDir = "/home/evanoberholster/go/src/github.com/evanoberholster/test/img"

func TestDecodePanasonicMakerNoteSample(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultPanasonicBenchImageDir
	}

	samplePath := filepath.Join(benchDir, "4.RW2")
	if _, err := os.Stat(samplePath); err != nil {
		t.Skipf("sample not found: %s", samplePath)
	}

	f, err := os.Open(samplePath)
	if err != nil {
		t.Fatalf("open %s: %v", samplePath, err)
	}
	defer func() { _ = f.Close() }()

	parsed, err := imagemeta.Decode(f)
	if err != nil {
		t.Fatalf("decode %s: %v", samplePath, err)
	}
	if parsed.MakerNote.Panasonic == nil {
		t.Fatalf("Panasonic maker-note missing for %s", samplePath)
	}

	got := parsed.MakerNote.Panasonic
	if got.ImageQuality != 7 {
		t.Fatalf("ImageQuality = %d, want 7", got.ImageQuality)
	}
	if got.FirmwareVersion != "0 1 0 0" {
		t.Fatalf("FirmwareVersion = %q, want %q", got.FirmwareVersion, "0 1 0 0")
	}
	if got.WhiteBalance != 1 {
		t.Fatalf("WhiteBalance = %d, want 1", got.WhiteBalance)
	}
	if got.FocusMode != 1 {
		t.Fatalf("FocusMode = %d, want 1", got.FocusMode)
	}
	if got.AFAreaMode != [2]uint8{32, 0} {
		t.Fatalf("AFAreaMode = %v, want %v", got.AFAreaMode, [2]uint8{32, 0})
	}
	if got.ImageStabilization != 2 {
		t.Fatalf("ImageStabilization = %d, want 2", got.ImageStabilization)
	}
	if got.MacroMode != 1 {
		t.Fatalf("MacroMode = %d, want 1", got.MacroMode)
	}
	if got.ShootingMode != 6 {
		t.Fatalf("ShootingMode = %d, want 6", got.ShootingMode)
	}
	if got.Audio != 2 {
		t.Fatalf("Audio = %d, want 2", got.Audio)
	}
	if got.FlashBias != 0 {
		t.Fatalf("FlashBias = %v, want 0", got.FlashBias)
	}
	if got.PanasonicExifVersion != "0270" {
		t.Fatalf("PanasonicExifVersion = %q, want %q", got.PanasonicExifVersion, "0270")
	}
	if got.ColorEffect != 1 {
		t.Fatalf("ColorEffect = %d, want 1", got.ColorEffect)
	}
	if got.TimeSincePowerOn != 60.25 {
		t.Fatalf("TimeSincePowerOn = %v, want 60.25", got.TimeSincePowerOn)
	}
	if got.BurstMode != 0 {
		t.Fatalf("BurstMode = %d, want 0", got.BurstMode)
	}
	if got.SequenceNumber != 0 {
		t.Fatalf("SequenceNumber = %d, want 0", got.SequenceNumber)
	}
	if got.ContrastMode != 0 {
		t.Fatalf("ContrastMode = %d, want 0", got.ContrastMode)
	}
	if got.NoiseReduction != 0 {
		t.Fatalf("NoiseReduction = %d, want 0", got.NoiseReduction)
	}
	if got.SelfTimer != 1 {
		t.Fatalf("SelfTimer = %d, want 1", got.SelfTimer)
	}
	if got.Rotation != 1 {
		t.Fatalf("Rotation = %d, want 1", got.Rotation)
	}
	if got.TravelDay != 65535 {
		t.Fatalf("TravelDay = %d, want 65535", got.TravelDay)
	}
	if got.BatteryLevel != 1 {
		t.Fatalf("BatteryLevel = %d, want 1", got.BatteryLevel)
	}
	if got.TextStamp != 1 {
		t.Fatalf("TextStamp = %d, want 1", got.TextStamp)
	}
	if got.PanasonicImageWidth != 3648 {
		t.Fatalf("PanasonicImageWidth = %d, want 3648", got.PanasonicImageWidth)
	}
	if got.PanasonicImageHeight != 2736 {
		t.Fatalf("PanasonicImageHeight = %d, want 2736", got.PanasonicImageHeight)
	}
	if got.MakerNoteVersion != "0130" {
		t.Fatalf("MakerNoteVersion = %q, want %q", got.MakerNoteVersion, "0130")
	}
	if got.SceneMode != 0 {
		t.Fatalf("SceneMode = %d, want 0", got.SceneMode)
	}
	if got.WBRedLevel != 1875 {
		t.Fatalf("WBRedLevel = %d, want 1875", got.WBRedLevel)
	}
	if got.WBGreenLevel != 1054 {
		t.Fatalf("WBGreenLevel = %d, want 1054", got.WBGreenLevel)
	}
	if got.WBBlueLevel != 1681 {
		t.Fatalf("WBBlueLevel = %d, want 1681", got.WBBlueLevel)
	}
}
