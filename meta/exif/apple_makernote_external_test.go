package exif_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta"
)

func TestAppleMakerNoteFixture(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "download_samples", "Apple", "Apple", "Apple_iPadAir2.jpg")
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("fixture unavailable: %v", err)
	}
	defer func() { _ = f.Close() }()

	exifData, err := imagemeta.Decode(f)
	if err != nil {
		t.Fatalf("Decode(%s) = %v", path, err)
	}
	if exifData.MakerNote.Apple == nil {
		t.Fatal("Apple maker note was not parsed")
	}

	got := exifData.MakerNote.Apple
	if got.MakerNoteVersion != 2 {
		t.Fatalf("MakerNoteVersion = %d, want 2", got.MakerNoteVersion)
	}
	if !got.AEStable || !got.AFStable {
		t.Fatalf("AEStable=%v AFStable=%v, want both true", got.AEStable, got.AFStable)
	}
	if got.AETarget != 243 || got.AEAverage != 242 {
		t.Fatalf("AETarget/AEAverage = %d/%d, want 243/242", got.AETarget, got.AEAverage)
	}
	if got.RunTime.Flags != 1 || got.RunTime.Scale != 1_000_000_000 || got.RunTime.Epoch != 0 {
		t.Fatalf("RunTime = %+v, want flags=1 scale=1000000000 epoch=0", got.RunTime)
	}
	if got.RunTime.Value == 0 {
		t.Fatal("RunTime.Value should not be zero")
	}
}
