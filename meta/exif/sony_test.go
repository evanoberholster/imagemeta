package exif

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSonyMakerNoteSamples(t *testing.T) {
	benchDir := os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR")
	if benchDir == "" {
		benchDir = defaultBenchImageDir
	}

	cases := []struct {
		file                  string
		quality               uint32
		quality2              [2]uint16
		sharpness             int32
		creativeStyle         string
		dynamicRangeOptimizer uint32
		imageStabilization    uint32
		sonyModelID           uint16
		lensType              uint32
	}{
		{
			file:                  "1.ARW",
			quality:               0,
			sharpness:             0,
			creativeStyle:         "Standard",
			dynamicRangeOptimizer: 3,
			imageStabilization:    1,
			sonyModelID:           281,
			lensType:              191,
		},
		{
			file:                  "Sony.ARW",
			quality:               6,
			sharpness:             0,
			creativeStyle:         "Standard",
			dynamicRangeOptimizer: 0,
			imageStabilization:    0,
			sonyModelID:           340,
			lensType:              65535,
		},
		{
			file:                  "SonyA7V.ARW",
			quality:               6,
			quality2:              [2]uint16{1, 2},
			sharpness:             4,
			creativeStyle:         "Standard",
			dynamicRangeOptimizer: 3,
			imageStabilization:    1,
			sonyModelID:           407,
			lensType:              65535,
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
			if parsed.MakerNote.Sony == nil {
				t.Fatalf("Sony maker-note missing for %s", samplePath)
			}

			got := parsed.MakerNote.Sony
			if got.Quality != tc.quality {
				t.Fatalf("Quality = %d, want %d", got.Quality, tc.quality)
			}
			if got.Quality2 != tc.quality2 {
				t.Fatalf("Quality2 = %v, want %v", got.Quality2, tc.quality2)
			}
			if got.Sharpness != tc.sharpness {
				t.Fatalf("Sharpness = %d, want %d", got.Sharpness, tc.sharpness)
			}
			if got.CreativeStyle != tc.creativeStyle {
				t.Fatalf("CreativeStyle = %q, want %q", got.CreativeStyle, tc.creativeStyle)
			}
			if got.DynamicRangeOptimizer != tc.dynamicRangeOptimizer {
				t.Fatalf("DynamicRangeOptimizer = %d, want %d", got.DynamicRangeOptimizer, tc.dynamicRangeOptimizer)
			}
			if got.ImageStabilization != tc.imageStabilization {
				t.Fatalf("ImageStabilization = %d, want %d", got.ImageStabilization, tc.imageStabilization)
			}
			if got.SonyModelID != tc.sonyModelID {
				t.Fatalf("SonyModelID = %d, want %d", got.SonyModelID, tc.sonyModelID)
			}
			if got.LensType != tc.lensType {
				t.Fatalf("LensType = %d, want %d", got.LensType, tc.lensType)
			}
		})
	}
}
