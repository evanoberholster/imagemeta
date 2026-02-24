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

const benchmarkSampleLimit = 5

var benchmarkPriorityModelTokens = []struct {
	model string
	token string
}{
	{model: "EOS R", token: "Canon_EOS_R_"},
	{model: "EOS R6", token: "Canon_EOS_R6_"},
	{model: "EOS R7", token: "Canon_EOS_R7_"},
}

func init() {
	Logger = log.Level(zerolog.PanicLevel)
}

func BenchmarkCR3Samples(b *testing.B) {
	paths := benchmarkCR3SamplePaths(b, benchmarkSampleLimit)

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
				if err = readMetadataToEOF(r); err != nil {
					b.Fatalf("ReadMetadata(%q): %v", path, err)
				}
			}
		})
	}
}

func benchmarkCR3SamplePaths(tb testing.TB, limit int) []string {
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
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".cr3") {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	if len(paths) == 0 {
		tb.Skipf("no .CR3 files found in %q", dir)
	}

	slices.Sort(paths)

	selected := make([]string, 0, len(paths))
	selectedSet := make(map[string]struct{}, len(paths))
	for _, priority := range benchmarkPriorityModelTokens {
		path := firstMatchingSamplePath(paths, priority.token)
		if path == "" {
			tb.Logf("benchmark sample for %s not found (token %q)", priority.model, priority.token)
			continue
		}
		selected = append(selected, path)
		selectedSet[path] = struct{}{}
	}

	for _, path := range paths {
		if limit > 0 && len(selected) >= limit {
			break
		}
		if _, ok := selectedSet[path]; ok {
			continue
		}
		selected = append(selected, path)
	}

	if len(selected) == 0 {
		tb.Skipf("no benchmark samples selected from %q", dir)
	}
	return selected
}

func firstMatchingSamplePath(paths []string, token string) string {
	for _, path := range paths {
		if strings.Contains(filepath.Base(path), token) {
			return path
		}
	}
	return ""
}

func benchmarkSamplesDir(tb testing.TB) string {
	tb.Helper()

	candidates := []string{
		"../cmd/samples",
		"cmd/samples",
	}
	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return dir
		}
	}

	tb.Skip("cmd/samples directory not found")
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
