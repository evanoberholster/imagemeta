# AGENTS.MD — imagemeta

## Project Overview

High-performance Go library for parsing EXIF metadata from JPEG, HEIC, AVIF, TIFF, DNG, CR2, CR3, CRW, NEF, ARW, RW2, and PNG. Targets parity with ExifTool as the reference implementation. Zero-allocation perceptual image hashing included.

## Core Principles

- **Idiomatic Go** — Use `gofmt`/`goimports` style. No JavaScript, Java, or Python patterns (no classes, no getters/setters, no builder patterns, no fluent interfaces).
- **Near-zero allocation** — Hot paths must avoid heap allocation. Use stack buffers (`[N]byte`), `sync.Pool`, and the `Acquire`/`Release` pattern. Always `b.ReportAllocs()` in benchmarks.
- **Functional options** — `ReaderOption func(*Reader)` pattern for optional config. Precompute lookup tables for option closures where possible.
- **Pooled resources** — Every pooled type has `Acquire*` / `Release*` functions. Always `defer Release*()` immediately after acquire.

## Project Structure

```
imagemeta/                  # Root module: github.com/evanoberholster/imagemeta
  imagemeta.go              # Public API: Decode(), DecodeCR3(), DecodeJPEG(), etc.
  imagetype/                # Image type detection (scan bytes, identify format)
  imagehash/                # Zero-alloc perceptual image hash (64/256 bit)
  preview/                  # Preview image extraction (CR3)
  meta/                     # Metadata model types (Exif, ExifHeader, etc.)
    exif/                   # Core EXIF IFD parser (main workhorse)
      tag/                  # Tag types, IFD types, value types
      makernote/            # MakerNote parsers
        canon/              # Canon makernote
        nikon/              # Nikon makernote
        sony/               # Sony makernote
        panasonic/          # Panasonic makernote
        apple/              # Apple makernote
    isobmff/                # ISO Base Media File Format (CR3, HEIC, AVIF)
    jpeg/                   # JPEG segment scanner
    png/                    # PNG chunk scanner
    xmp/                    # XMP XML parser (sax-based)
    logging/                # Zerolog-based logging mixin
    utils/                  # Byte order, bufio pool, limited reader
  cmd/                      # CLI tools
    compare_download_samples/  # Comparison tool for download_samples
  testImages/               # Small test images
  download_samples/         # Camera sample catalog (123 makes, NOT committed)
  golangci.yml              # Linter config
```

## Coding Conventions

### Imports
```
import (
    "stdlib"        // stdlib first, no blank line
    "io"
    "log/slog"
    "sync"

    "github.com/pkg/errors"              // third-party first group

    "github.com/evanoberholster/imagemeta/imagetype"  // internal second group
    "github.com/evanoberholster/imagemeta/meta"
)
```

### Error Handling
- **Sentinel errors**: Define in `var (...)` blocks using `errors.New` (stdlib), not `fmt.Errorf`:
  ```go
  var (
      ErrNoExif = meta.ErrNoExif
      ErrImageTypeNotFound = imagetype.ErrImageTypeNotFound
  )
  ```
- **Error wrapping**: Use `github.com/pkg/errors.Wrapf(err, "message")` with lowercase message, no trailing punctuation.
- **Error comparison**: Always use `errors.Is(err, ErrNoExif)`, never `==`.
- **Graceful degradation**: Skip invalid tags rather than aborting. Log non-fatal issues at warn level.
- **Conditional logging**: Always check level before creating log event: `if r.warnEnabled() { ... }`.

### Naming
- Exported types: PascalCase (`IFD0Tag`, `ExifIFDTags`, `Reader`)
- Unexported: camelCase (`state`, `eofReader`, `loggerMixin`)
- Constants: PascalCase exported (`ImageJPEG`), camelCase unexported
- Receiver names: single letter (`r *Reader`, `s *state`)
- Avoid stutter: `exif.Exif` not `exif.ExifReader`
- File names: short, single word where possible

### Pooling Pattern
```go
var parseReaderPool = sync.Pool{
    New: func() any { return &Reader{...} },
}

func AcquirePooledReader(l *slog.Logger) *Reader {
    r, ok := parseReaderPool.Get().(*Reader)
    if !ok || r == nil {
        r = &Reader{...}
    }
    r.resetDecodeState(true)
    return r
}

func ReleasePooledReader(r *Reader) {
    if r == nil { return }
    r.state.reset()
    statePool.Put(r.state)
    r.state = nil
    parseReaderPool.Put(r)
}
```

### Zero-Alloc Hot Path Guidelines
1. Prefer `[N]byte` arrays on stack over `make([]byte, N)`.
2. Use fixed-size tag queues (`[128]tag.Entry`) instead of dynamic slices.
3. Embedded values (≤4 bytes) read directly from IFD entry without allocation.
4. Returned byte slices reference parser buffers — document they are only valid until next read.
5. Use `sync.Pool` for everything: readers, states, bufio readers, pixel arrays.
6. Avoid `fmt.Sprintf` in hot paths; use structured slog fields.
7. Avoid `interface{}` allocations; use concrete types.
8. Use insertion sort for small fixed-size queues (avoids allocating `sort.Interface`).

### Code Generation
- `//go:generate msgp` for MessagePack serialization (`*_gen.go`)
- `//go:generate stringer -type=TypeName` for enum string methods
- `//go:generate avo` for assembly generation
- Generated files: `*_gen.go`, `zz_generated.*.go`

## Testing

### Standard Tests
- Run: `go test ./...`
- Use `t.Parallel()` in all tests.
- Use `t.Fatal` / `t.Fatalf` for errors (not `t.Error`).
- Golden file comparisons in `testImages/` (`.exif`, `.json` sidecars).
- External tests: `*_external_test.go` for package-external blackbox tests.

### download_samples Tests
The `download_samples/` directory is a catalog of camera makes (not committed to git). Tests that validate against real camera samples go here:
- `cmd/compare_download_samples/` — ExifTool comparison tool and tests.
- New makernote parsers must be tested against real samples from the appropriate manufacturer directory.
- Unexported functions should provide display-value formatting that matches ExifTool output.

### Benchmarks
- Every benchmark must call `b.ReportAllocs()` and `b.SetBytes(int64(len(data)))`.
- Load file into `[]byte` with `os.ReadFile`, use `bytes.NewReader(data)`, call `r.Seek(0, 0)` between iterations.
- Table-driven with `b.Run(name, ...)` sub-benchmarks.
- Environment variable `IMAGEMETA_BENCH_IMAGE_DIR` overrides the benchmark image directory.
- Target: JPEG: ≤64 B/op, 2 allocs/op. CR2: ≤1600 B/op, ≤20 allocs/op. CR3: ≤3800 B/op, ≤22 allocs/op.

## Linting

Run `golangci-lint run` before committing. Config in `golangci.yml`. Key linters:
- `gofmt` / `goimports` — formatting
- `govet` — includes shadow variable detection
- `gosec` — security (excluded for `_test.go` and `_gen.go`)
- `errcheck` — checks type assertions and blank assigns
- `errorlint` — enforces `errors.Is` / `errors.As`
- `misspell` — US locale
- `nolintlint` — requires explanation and specific linter ID
- `unparam` / `unused` / `ineffassign` — dead code detection

## Key Dependencies

| Dependency | Purpose |
|---|---|
| `github.com/pkg/errors` | Error wrapping (Wrapf) |
| `log/slog` | Structured logging |
| `github.com/klauspost/cpuid/v2` | CPU feature detection (SIMD) |
| `github.com/tinylib/msgp` | MessagePack codegen |
| `github.com/mmcloughlin/avo` | Assembly codegen |
| `github.com/orisano/gosax` | SAX XML parser for XMP |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/nfnt/resize` | Image resize for hashing |

## Entry Points

Primary public API in `imagemeta.go`:
- `Decode(r io.ReadSeeker) (exif.Exif, error)` — auto-detect and decode
- `DecodeCR3(r io.ReadSeeker) (exif.Exif, error)`
- `DecodeTiff(r io.ReadSeeker) (exif.Exif, error)` — covers CR2, DNG, NEF, ARW, RW2
- `DecodeJPEG(r io.ReadSeeker) (exif.Exif, error)`
- `DecodeHeif(r io.ReadSeeker) (exif.Exif, error)`
- `DecodePng(r io.ReadSeeker) (exif.Exif, error)`
- `DecodeCRW(r io.ReadSeeker) (exif.Exif, error)`
- `PreviewCR3(r io.ReadSeeker) ([]byte, error)`

Internal decode flow:
1. `imagetype.ScanBuf(rr)` — detect image type from bytes
2. Dispatch to format-specific path (JPEG segment scan, TIFF header scan, ISOBMFF FTYP+metadata scan)
3. `exif.Reader.DecodeTiff` / `DecodeIfdAppend` — parse IFD tree
4. Tag values resolved through `drainQueuedTags` -> `parseTag` dispatch

## Go Version

Go 1.25. Module: `github.com/evanoberholster/imagemeta`
