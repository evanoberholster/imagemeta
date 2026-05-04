package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/evanoberholster/imagemeta"
)

type result struct {
	Path      string `json:"path"`
	OK        bool   `json:"ok"`
	ImageType string `json:"imageType,omitempty"`
	Error     string `json:"error,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file-or-dir> [...]\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	paths, err := collectPaths(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	exitCode := 0
	for _, p := range paths {
		r := decodePath(p)
		if !r.OK {
			exitCode = 1
		}
		if err := enc.Encode(r); err != nil {
			fmt.Fprintf(os.Stderr, "encode: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}

func collectPaths(args []string) ([]string, error) {
	out := make([]string, 0, len(args))
	seen := make(map[string]struct{}, len(args))

	for _, in := range args {
		info, err := os.Stat(in)
		if err != nil {
			return nil, fmt.Errorf("stat %q: %w", in, err)
		}
		if !info.IsDir() {
			if includePath(in) {
				ap, absErr := filepath.Abs(in)
				if absErr != nil {
					return nil, fmt.Errorf("abs %q: %w", in, absErr)
				}
				if _, ok := seen[ap]; !ok {
					seen[ap] = struct{}{}
					out = append(out, ap)
				}
			}
			continue
		}

		err = filepath.WalkDir(in, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			if !includePath(path) {
				return nil
			}
			ap, absErr := filepath.Abs(path)
			if absErr != nil {
				return absErr
			}
			if _, ok := seen[ap]; ok {
				return nil
			}
			seen[ap] = struct{}{}
			out = append(out, ap)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk %q: %w", in, err)
		}
	}

	sort.Strings(out)
	return out, nil
}

func includePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".cr2", ".cr3", ".crw", ".tif", ".tiff", ".dng", ".nef", ".arw", ".heic", ".heif", ".avif", ".png", ".jxl":
		return true
	default:
		return false
	}
}

func decodePath(path string) result {
	f, err := os.Open(path)
	if err != nil {
		return result{Path: path, OK: false, Error: err.Error()}
	}

	e, err := imagemeta.Decode(f)
	if err != nil {
		_ = f.Close()
		return result{Path: path, OK: false, Error: err.Error()}
	}
	if err = f.Close(); err != nil && !errors.Is(err, io.EOF) {
		return result{Path: path, OK: false, Error: err.Error()}
	}
	return result{Path: path, OK: true, ImageType: e.ImageType.String()}
}
