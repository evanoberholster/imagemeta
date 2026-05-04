# imagemeta

[![License][License-Image]][License-Url]
[![Go Reference][GoDoc-Image]][GoDoc-Url]
[![Go Report Card][ReportCard-Image]][ReportCard-Url]
[![Coverage][Coverage-Image]][Coverage-Url]
[![Build][Build-Status-Image]][Build-Status-Url]

`imagemeta` is a high-performance Go library for extracting EXIF and XMP metadata from images and camera RAW formats.

It is built for:

- ExifTool-aligned decoding behavior
- low-allocation parsing on hot paths
- broad support for modern and legacy camera formats
- structured MakerNote decoding (Canon, Nikon, Sony, Panasonic, Apple)

## Features

- Decode EXIF from JPEG, TIFF, CR2, CR3, DNG, NEF, ARW, RW2, HEIC/HEIF/AVIF, PNG, and CRW
- Parse XMP sidecars and embedded XMP blocks
- Read Canon/Nikon/Sony/Panasonic/Apple MakerNotes
- Extract embedded CR3 previews
- Fast image type detection
- Allocation-conscious perceptual image hashing (`imagehash`)

## Install

```bash
go get github.com/evanoberholster/imagemeta
```

## Supported Formats

- JPEG: `.jpg`, `.jpeg`
- TIFF + TIFF-based RAW: `.tif`, `.tiff`, `.cr2`, `.dng`, `.nef`, `.arw`, `.rw2`
- ISO-BMFF family: `.cr3`, `.heic`, `.heif`, `.avif`
- PNG (EXIF-in-TIFF payload)
- CIFF/CRW: `.crw`

## Quick Start

```go
package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta"
)

func main() {
	f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ex, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("type:", ex.ImageType)
	fmt.Println("make:", ex.Make)
	fmt.Println("model:", ex.Model)
	fmt.Println("datetime:", ex.DateTimeOriginal)
}
```

## Decode API

Top-level decode entry points:

- `Decode(io.ReadSeeker)`
- `DecodeJPEG(io.ReadSeeker)`
- `DecodeTiff(io.ReadSeeker)`
- `DecodeCR3(io.ReadSeeker)`
- `DecodeCRW(io.ReadSeeker)`
- `DecodeHeif(io.ReadSeeker)`
- `DecodePng(io.ReadSeeker)`
- `PreviewCR3(io.ReadSeeker)`

## XMP Parsing

Parse a sidecar `.xmp` file:

```go
x, err := xmp.ParseXmp(file)
```

Parse embedded XMP from an image:

```go
x, err := xmp.Parse(file)
```

Package: `github.com/evanoberholster/imagemeta/meta/xmp`

## Image Type Detection

Use `imagetype` for quick format detection before decode:

```go
t, err := imagetype.Scan(r)
```

Package: `github.com/evanoberholster/imagemeta/imagetype`

## Perceptual Hashing

`imagehash` provides 64-bit and 256-bit perceptual hashing with low-allocation paths.

Package: `github.com/evanoberholster/imagemeta/imagehash`

## Performance Notes

- Designed for low allocation and high throughput
- Uses pooled readers and fixed-size parsing structures on hot paths
- Suitable for large-batch metadata extraction pipelines

## Benchmarks

- General benchmarks: `bench_test.go`

| Benchmark | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| Canon/CR2 | 6646 | 2440 | 20 |
| Canon/CR3 | 8015 | 3723 | 22 |
| Canon/JPG | 2058 | 176 | 9 |
| Nikon/JPG | 1747 | 104 | 7 |
| Nikon/NEF | 10306 | 736 | 19 |
| Sony/ARW | 6584 | 1512 | 12 |

Note: benchmark results vary by CPU, Go version, and sample files.

## TODO

- Expand ExifTool parity coverage across MakerNote fields (Nikon, Sony, Panasonic, Apple, DJI)
- Resolve known reverse-offset parsing edge case in MakerNote decoding
- Improve malformed metadata resilience with more targeted fuzz and corpus tests
- Add more specific benchmarks to compare edge cases
- Document compatibility matrix by camera model/format 

## Contributing

Pull requests and issue reports are welcome.

Before opening a PR, run:

```bash
go test ./...
golangci-lint run
```

## Acknowledgements

`imagemeta` is heavily inspired and tested on:

- Phil Harvey's ExifTool: <https://exiftool.org>
- `rwcarlsen/goexif`

Additional references:

- `go4` (BMFF/HEIF parser concepts)
- `lclevy/canon_cr3`
- Lasse Heikkilä's HEIF thesis work

### LLMs were used to assist in writing code, the code has been reviewed and has passed stringent checks.

## License

MIT

[License-Url]: https://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/badge/License-MIT-blue.svg?maxAge=2592000
[GoDoc-Url]: https://pkg.go.dev/github.com/evanoberholster/imagemeta
[GoDoc-Image]: https://pkg.go.dev/badge/github.com/evanoberholster/imagemeta.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/evanoberholster/imagemeta
[ReportCard-Image]: https://goreportcard.com/badge/github.com/evanoberholster/imagemeta
[Coverage-Image]: https://coveralls.io/repos/github/evanoberholster/imagemeta/badge.svg?branch=master
[Coverage-Url]: https://coveralls.io/github/evanoberholster/imagemeta?branch=master
[Build-Status-Url]: https://github.com/evanoberholster/imagemeta/actions/workflows/golangci-lint.yml?query=branch%3Amaster
[Build-Status-Image]: https://github.com/evanoberholster/imagemeta/actions/workflows/golangci-lint.yml/badge.svg?branch=master
