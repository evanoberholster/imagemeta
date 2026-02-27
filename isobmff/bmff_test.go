package isobmff

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	Logger = log.Level(zerolog.PanicLevel)
}

func BenchmarkISOBMFFSamples(b *testing.B) {
	paths := benchmarkSamplePaths(b)

	for i, path := range paths {
		path := path
		label := benchmarkSampleLabel(i, path)
		b.Run(label, func(b *testing.B) {
			data, err := os.ReadFile(path)
			if err != nil {
				b.Fatalf("ReadFile(%q): %v", path, err)
			}

			reader := bytes.NewReader(data)
			r := NewReader(reader, nil, nil, nil)
			b.Cleanup(r.Close)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if _, err = reader.Seek(0, io.SeekStart); err != nil {
					b.Fatalf("Seek(%q): %v", path, err)
				}
				r.reset(reader)

				if err = r.ReadFTYP(); err != nil {
					b.Fatalf("ReadFTYP(%q): %v", path, err)
				}
				if err = readMetadataToEOFOrBufLength(r); err != nil {
					b.Fatalf("ReadMetadata(%q): %v", path, err)
				}
			}
		})
	}
}

func benchmarkSamplePaths(tb testing.TB) []string {
	tb.Helper()

	dir := benchmarkSamplesDir(tb)
	entries, err := os.ReadDir(dir)
	if err != nil {
		tb.Fatalf("ReadDir(%q): %v", dir, err)
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	if len(paths) == 0 {
		tb.Skipf("no sample files found in %q", dir)
	}

	slices.Sort(paths)
	return paths
}

func benchmarkSamplesDir(tb testing.TB) string {
	tb.Helper()

	candidates := []string{
		"samples",
		"isobmff/samples",
	}
	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return dir
		}
	}

	tb.Skip("isobmff/samples directory not found")
	return ""
}

func benchmarkSampleLabel(index int, path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.ReplaceAll(base, "[", "")
	base = strings.ReplaceAll(base, "]", "")
	base = strings.ReplaceAll(base, "__", "_")
	base = strings.ReplaceAll(base, "Canon_", "")

	if len(base) > 40 {
		base = base[:40]
	}
	return fmt.Sprintf("%02d_%s", index+1, base)
}

func readMetadataToEOF(r *Reader) error {
	for {
		err := r.ReadMetadata()
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
}

func TestBrandFromBufAdditionalBrands(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		want brand
	}{
		{name: "heif", buf: []byte("heif"), want: brandHeif},
		{name: "avis", buf: []byte("avis"), want: brandAvis},
		{name: "3gp6", buf: []byte("3gp6"), want: brand3GP6},
		{name: "3g2a", buf: []byte("3g2a"), want: brand3G2A},
		{name: "M4V", buf: []byte("M4V "), want: brandM4V},
		{name: "mp71", buf: []byte("mp71"), want: brandMp71},
		{name: "dash", buf: []byte("dash"), want: brandDash},
		{name: "jxl", buf: []byte("jxl "), want: brandJxl},
		{name: "qt", buf: []byte("qt  "), want: brandQt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := brandFromBuf(tt.buf)
			if got != tt.want {
				t.Fatalf("brandFromBuf(%q) = %v, want %v", tt.buf, got, tt.want)
			}
			if got.String() != string(tt.buf) {
				t.Fatalf("brand string = %q, want %q", got.String(), string(tt.buf))
			}
		})
	}
}

func TestReadFTYPSkipsJXLSignatureBox(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x0C, // size
		'J', 'X', 'L', ' ', // type
		0x0D, 0x0A, 0x87, 0x0A, // JXL signature payload
		0x00, 0x00, 0x00, 0x10, // size
		'f', 't', 'y', 'p', // type
		'a', 'v', 'i', 'f', // major brand
		'0', '0', '0', '1', // minor version
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if r.ftyp.MajorBrand != brandAvif {
		t.Fatalf("MajorBrand = %v, want %v", r.ftyp.MajorBrand, brandAvif)
	}
}

func TestReadFTYPErrorPrefix(t *testing.T) {
	r := NewReader(bytes.NewReader(nil), nil, nil, nil)
	t.Cleanup(r.Close)

	err := r.ReadFTYP()
	if err == nil {
		t.Fatal("expected ReadFTYP error")
	}
	if !strings.Contains(err.Error(), "ReadFTYP:") {
		t.Fatalf("ReadFTYP error = %q, want prefix %q", err.Error(), "ReadFTYP:")
	}
}

func TestReadMetadataFromSamples(t *testing.T) {
	paths := benchmarkSamplePaths(t)
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", path, err)
			}

			r := NewReader(bytes.NewReader(data), nil, nil, nil)
			t.Cleanup(r.Close)

			if err := r.ReadFTYP(); err != nil {
				t.Fatalf("ReadFTYP(%q): %v", path, err)
			}
			if err := readMetadataToEOFOrBufLength(r); err != nil {
				t.Fatalf("ReadMetadata(%q): %v", path, err)
			}
		})
	}
}

func readMetadataToEOFOrBufLength(r *Reader) error {
	for {
		err := r.ReadMetadata()
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) || errors.Is(err, ErrBufLength) {
			return nil
		}
		return err
	}
}
