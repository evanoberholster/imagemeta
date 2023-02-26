# Imagemeta

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Coverage Status][Coverage-Image]][Coverage-Url]
[![Build][Build-Status-Image]][Build-Status-Url]

Image Metadata (Exif and XMP) extraction for JPEG, HEIC, AVIF, TIFF, and Camera Raw in golang. Imagetype identifcation. Zero allocation Perceptual Image Hash. Goal is features that are precise and performance oriented for working with images.

## Documentation

See [Documentation](https://godoc.org/github.com/evanoberholster/imagemeta) for more information.

## Example Usage

Example usage:

```go
    package main

    import (
    	"fmt"
    	"os"

    	"github.com/evanoberholster/imagemeta"
    )

	f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
```

## Imagehash
 Zero allocation PerceptualHash algorithm (64Bit and 256Bit) [github.com/evanoberholster/imagemeta/imagehash](github.com/evanoberholster/imagemeta/imagehash). Adapted from [https://github.com/corona10/goimagehash](https://github.com/corona10/goimagehash). Image will need to be resized to 64x64 prior to image hashing.

## Contributing

Issues, Suggestions and Pull Requests are welcome.

## Benchmarks

See BENCHMARK.md
To run your own benchmarks see bench_test.go

## Imagetype Identification

Images can be identified with: "github.com/evanoberholster/imagemeta/imagetype" package.

## TODO

- [x] Stabilize ImageTypes API
- [x] Add Exif parsing for individual image types (jpg, heic, cr2, tiff, dng)
- [x] Add CR3 and Heic image metadata support.
- [x] Add Avif image metadata support
- [ ] Add Canon Exif Makernote support
- [ ] Add Nikon Exif Makernote support 
- [ ] Add Camera Make and Model Lookup tables
- [ ] Add Preview Image extraction
- [ ] Refactor XMP parsing as "xmp" package
- [ ] Stabalize Imagemeta API
- [ ] Improve test coverage
- [ ] Add Webp image metadata support
- [ ] Add CRW image metadata support (ciff format images)
- [ ] Documentation

## Based on and Inspired by

Inspired by Phil Harvey [http://exiftool.org](http://exiftool.org), go-exif [https://github.com/dsoprea/go-exif](https://github.com/dsoprea/go-exif), and RW Carlsen [https://github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif)

## Special Thanks to:
- The go4 Authors (https://github.com/go4org/go4) for their work on a BMFF parser and HEIF structure in golang.
- Laurent Clévy (@Lorenzo2472) (https://github.com/lclevy/canon_cr3) for Canon CR3 structure.
- Lasse Heikkilä (https://trepo.tuni.fi/bitstream/handle/123456789/24147/heikkila.pdf) for HEIF structure from his thesis.
- Imagehash authors (https://github.com/corona10/goimagehash)

### Contributors
- Anders Brander [abrander](https://github.com/abrander)
- Dobrosław Żybort [matrixik](https://github.com/matrixik)

## LICENSE

Copyright (c) 2020-2023, Evan Oberholster & Contributors

[License-Url]: https://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/badge/License-MIT-blue.svg?maxAge=2592000
[Godoc-Url]: https://godoc.org/github.com/evanoberholster/imagemeta
[Godoc-Image]: https://godoc.org/github.com/evanoberholster/imagemeta?status.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/evanoberholster/imagemeta
[ReportCard-Image]: https://goreportcard.com/badge/github.com/evanoberholster/imagemeta
[Coverage-Image]: https://coveralls.io/repos/github/evanoberholster/imagemeta/badge.svg?branch=master
[Coverage-Url]: https://coveralls.io/github/evanoberholster/imagemeta?branch=master
[Build-Status-Url]: https://github.com/evanoberholster/imagemeta/actions?query=branch%3Amaster
[Build-Status-Image]: https://github.com/evanoberholster/imagemeta/workflows/Build/badge.svg?branch=master
