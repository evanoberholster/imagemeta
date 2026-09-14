//go:build !wasm

package imagehash

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Meta's PDQ reference data lives in the ThreatExchange repository. These
// vectors and tolerance mirror the reference implementation's test suite and
// are used to verify this package against Meta's original implementation.
//
// Fixtures are looked up in $IMAGEMETA_PDQ_DATA, then in testdata/pdq, and are
// otherwise downloaded to a cache directory. If they cannot be obtained the
// Meta comparison tests are skipped.
const (
	pdqBaseURL   = "https://raw.githubusercontent.com/facebook/ThreatExchange/main/pdq/data"
	pdqSentinel  = "reg-test-input/dih/bridge-1-original.jpg"
	pdqTolerance = 16
)

type pdqReferenceVector struct {
	path string
	hash string
}

type pdqQualityVector struct {
	path       string
	minQuality int
	maxQuality int
}

var pdqReferenceVectors = []pdqReferenceVector{
	{"misc-images/c.png", "e64cc9d91e623842f8d1f1d9a398e78c9f199a3bd87924f2b7e11e0bf061b064"},
	{"misc-images/small.jpg", "0007001f003f003f007f00ff00ff00ff01ff01ff01ff03ff03ff03ff03ff03ff"},
	{"misc-images/wee.jpg", "6227401f601ff4ccafcc9fad4b0d95d371a2eb7265a3285234d228ca94deeb2d"},
	{"reg-test-input/labelme-subset/q0003.jpg", "54a977c221d14c1c43ba5e6e21d4a13989a3553f1462611cbb85fda7be83b677"},
	{"reg-test-input/labelme-subset/q0004.jpg", "992d44af36d69e6ca6b812585928bac11def254ef5398c6d07466c9abcc65b92"},
	{"reg-test-input/labelme-subset/q0122.jpg", "cfb2009ddd21c6dab0046a7745b5984757a8a4535b3377aea2591d32b33ff940"},
	{"reg-test-input/labelme-subset/q0291.jpg", "a0fe94f1e5cc1cc8dd855948498dc9243f7ca27336f036d7f212b74bc103c9a7"},
	{"reg-test-input/labelme-subset/q0746.jpg", "1049d96239e24d4dca2c55512b8bdb77425f4dbcf575a0a95555aaab5554aaaa"},
	{"reg-test-input/labelme-subset/q1050.jpg", "489db672e9190276d452aeab41eba20f02375fe4092d88defdf491a5c55c5f70"},
	{"reg-test-input/labelme-subset/q2821.jpg", "b150231ffae4710ffcf4f18bb574b109a576f14bb8543189f8743289f174b109"},
}

var pdqDihedralVectors = []pdqReferenceVector{
	{"reg-test-input/dih/bridge-1-original.jpg", "d8f8f0cce0f4a84f0e370a22028f67f0b36e2ed596623e1d33e6b39c4e9c9b22"},
	{"reg-test-input/dih/bridge-2-rotate-90.jpg", "38a50efd71c83f429013d68d0ffffc52e34e0e15ada952a9d29684214aa9e5af"},
	{"reg-test-input/dih/bridge-3-rotate-180.jpg", "2dadda64b5a142e5d362209057da895ae63b8c7fc277b4b766b319361f893188"},
	{"reg-test-input/dih/bridge-4-rotate-270.jpg", "a5f0a457248995e8c9065c275aaa54d8b61ba4bdf8fcfc0387c32f8b0bfc4f05"},
	{"reg-test-input/dih/bridge-5-flipx.jpg", "d8f80f31e0f417b00e37f5dd028f980fb36ed12a9662c1e233e64c634e9c64dd"},
	{"reg-test-input/dih/bridge-6-flipy.jpg", "0dad259bb1a1bd18d362576556da32a1e63b7380c2374b4866b3c6c91b89ce77"},
	{"reg-test-input/dih/bridge-7-flip-plus-1.jpg", "f0a5e10271dcc0bd9c5309720fff018de34ef1e8ada9a956d2967ade1ea91a50"},
	{"reg-test-input/dih/bridge-8-flip-minus-1.jpg", "69f05aa8a4996a17c146a2da5aaaab07b61b5b60f8fc07fc83c3d0740bfcb0fa"},
}

var pdqQualityVectors = []pdqQualityVector{
	{"misc-images/small.jpg", 0, 10},
	{"misc-images/c.png", 50, 100},
	{"misc-images/wee.jpg", 50, 100},
	{"reg-test-input/labelme-subset/q0003.jpg", 0, 10},
	{"reg-test-input/labelme-subset/q0004.jpg", 0, 10},
	{"reg-test-input/labelme-subset/q0122.jpg", 50, 100},
	{"reg-test-input/labelme-subset/q0291.jpg", 50, 100},
	{"reg-test-input/labelme-subset/q0746.jpg", 50, 100},
	{"reg-test-input/labelme-subset/q1050.jpg", 50, 100},
	{"reg-test-input/labelme-subset/q2821.jpg", 50, 100},
	{"reg-test-input/dih/bridge-1-original.jpg", 50, 100},
	{"reg-test-input/dih/bridge-2-rotate-90.jpg", 50, 100},
	{"reg-test-input/dih/bridge-3-rotate-180.jpg", 50, 100},
	{"reg-test-input/dih/bridge-4-rotate-270.jpg", 50, 100},
	{"reg-test-input/dih/bridge-5-flipx.jpg", 50, 100},
	{"reg-test-input/dih/bridge-6-flipy.jpg", 50, 100},
	{"reg-test-input/dih/bridge-7-flip-plus-1.jpg", 50, 100},
	{"reg-test-input/dih/bridge-8-flip-minus-1.jpg", 50, 100},
}

var pdqFixtures = []string{
	"misc-images/c.png",
	"misc-images/small.jpg",
	"misc-images/wee.jpg",
	"reg-test-input/labelme-subset/q0003.jpg",
	"reg-test-input/labelme-subset/q0004.jpg",
	"reg-test-input/labelme-subset/q0122.jpg",
	"reg-test-input/labelme-subset/q0291.jpg",
	"reg-test-input/labelme-subset/q0746.jpg",
	"reg-test-input/labelme-subset/q1050.jpg",
	"reg-test-input/labelme-subset/q2821.jpg",
	"reg-test-input/dih/bridge-1-original.jpg",
	"reg-test-input/dih/bridge-2-rotate-90.jpg",
	"reg-test-input/dih/bridge-3-rotate-180.jpg",
	"reg-test-input/dih/bridge-4-rotate-270.jpg",
	"reg-test-input/dih/bridge-5-flipx.jpg",
	"reg-test-input/dih/bridge-6-flipy.jpg",
	"reg-test-input/dih/bridge-7-flip-plus-1.jpg",
	"reg-test-input/dih/bridge-8-flip-minus-1.jpg",
}

var (
	pdqDataOnce sync.Once
	pdqDataDir  string
	pdqDataErr  error
)

// pdqFixtureDir returns a directory containing Meta's ThreatExchange PDQ test
// data, downloading it to a cache directory if a local copy is not available.
// The test is skipped if the fixtures cannot be located or fetched.
func pdqFixtureDir(t *testing.T) string {
	t.Helper()
	pdqDataOnce.Do(func() {
		dirs := []string{os.Getenv("IMAGEMETA_PDQ_DATA"), filepath.Join("testdata", "pdq")}
		for _, dir := range dirs {
			if dir == "" {
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(pdqSentinel))); err == nil {
				pdqDataDir = dir
				return
			}
		}
		cache := os.Getenv("IMAGEMETA_PDQ_CACHE")
		if cache == "" {
			base, err := os.UserCacheDir()
			if err != nil {
				pdqDataErr = fmt.Errorf("locate cache directory: %w", err)
				return
			}
			cache = filepath.Join(base, "imagemeta", "pdq-data")
		}
		if err := pdqDownloadFixtures(cache); err != nil {
			pdqDataErr = err
			return
		}
		pdqDataDir = cache
	})
	if pdqDataErr != nil {
		t.Skipf("Meta PDQ fixtures unavailable: %v\nSet IMAGEMETA_PDQ_DATA to a local copy of ThreatExchange pdq/data", pdqDataErr)
	}
	return pdqDataDir
}

func pdqDownloadFixtures(dir string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	for _, rel := range pdqFixtures {
		dest := filepath.Join(dir, filepath.FromSlash(rel))
		if _, err := os.Stat(dest); err == nil {
			continue
		}
		if err := pdqDownload(client, strings.TrimSuffix(pdqBaseURL, "/")+"/"+rel, dest); err != nil {
			return fmt.Errorf("download %s: %w", rel, err)
		}
	}
	return nil
}

func pdqDownload(client *http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	if err = os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".pdq-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, dest)
}

// loadPDQFixture decodes a fixture image, honoring its JPEG bounds.
func loadPDQFixture(t *testing.T, dir, rel string) image.Image {
	t.Helper()
	f, err := os.Open(filepath.Join(dir, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("open %s: %v", rel, err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode %s: %v", rel, err)
	}
	// Normalise the decoded bounds to the pixel origin so the PDQ algorithm
	// sees a consistent coordinate space.
	b := img.Bounds()
	if b.Min != (image.Point{}) {
		normalised := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		for y := 0; y < b.Dy(); y++ {
			for x := 0; x < b.Dx(); x++ {
				normalised.Set(x, y, img.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return normalised
	}
	return img
}

func TestPDQ256MetaReferenceVectors(t *testing.T) {
	dir := pdqFixtureDir(t)
	for _, tc := range pdqReferenceVectors {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			img := loadPDQFixture(t, dir, tc.path)
			got, err := NewPDQ256(img)
			if err != nil {
				t.Fatalf("NewPDQ256: %v", err)
			}
			want, err := pdqHashFromHex(tc.hash)
			if err != nil {
				t.Fatalf("bad reference hash: %v", err)
			}
			if d := got.Distance(want); d > pdqTolerance {
				t.Errorf("hamming %d > tolerance %d\n  got:  %s\n  want: %s", d, pdqTolerance, got, want)
			} else {
				t.Logf("hamming=%d  %s", d, got)
			}
		})
	}
}

func TestPDQ256MetaDihedralVectors(t *testing.T) {
	dir := pdqFixtureDir(t)
	for _, tc := range pdqDihedralVectors {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			img := loadPDQFixture(t, dir, tc.path)
			got, err := NewPDQ256(img)
			if err != nil {
				t.Fatalf("NewPDQ256: %v", err)
			}
			want, err := pdqHashFromHex(tc.hash)
			if err != nil {
				t.Fatalf("bad reference hash: %v", err)
			}
			if d := got.Distance(want); d > pdqTolerance {
				t.Errorf("hamming %d > tolerance %d\n  got:  %s\n  want: %s", d, pdqTolerance, got, want)
			}
		})
	}
}

func TestPDQ256MetaQuality(t *testing.T) {
	dir := pdqFixtureDir(t)
	for _, tc := range pdqQualityVectors {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			img := loadPDQFixture(t, dir, tc.path)
			_, quality, err := NewPDQ256WithQuality(img)
			if err != nil {
				t.Fatalf("NewPDQ256WithQuality: %v", err)
			}
			if quality < tc.minQuality || quality > tc.maxQuality {
				t.Errorf("quality = %d, want in [%d, %d]", quality, tc.minQuality, tc.maxQuality)
			}
		})
	}
}
