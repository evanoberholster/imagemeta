# Exif Tool

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Coverage Status][Coverage-Image]][Coverage-Url]
[![Build][Build-Status-Image]][Build-Status-Url]

Image Metadata (Exif and XMP) extraction for JPEG, HEIC, WebP, AVIF, TIFF and Camera Raw in golang. Focus is on providing wide variety of features while being perfomance oriented.

## Documentation

See [Documentation](https://godoc.org/github.com/evanoberholster/imagemeta) for more information.

## Example Usage

Example usage:

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: main <file>\n")
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	m, err := imagemeta.Parse(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(m.Exif())
	fmt.Println(m.Xmp())
	fmt.Println(m.ImageType())
	fmt.Println(m.Dimensions())
	fmt.Println(jpeg.DecodeConfig(m.PreviewImage()))

	e, _ := m.Exif()
	if e != nil {
		// ImageWidth and ImageHeight
		fmt.Println(e.Dimensions().Size())

		fmt.Println(e.Artist())
		fmt.Println(e.Copyright())

		fmt.Println(e.CameraMake())
		fmt.Println(e.CameraModel())
		fmt.Println(e.CameraSerial())

		fmt.Println(e.Orientation())

		fmt.Println(e.LensMake())
		fmt.Println(e.LensModel())
		fmt.Println(e.LensSerial())

		fmt.Println(e.ISOSpeed())
		fmt.Println(e.FocalLength())
		fmt.Println(e.LensModel())
		fmt.Println(e.Aperture())
		fmt.Println(e.ShutterSpeed())

		fmt.Println(e.Dimensions().Size())

		fmt.Println(e.Artist())
		fmt.Println(e.Copyright())

		fmt.Println(e.ISOSpeed())
		fmt.Println(e.FocalLength())
		fmt.Println(e.LensModel())
		fmt.Println(e.Aperture())
		fmt.Println(e.ShutterSpeed())

		fmt.Println(e.Aperture())
		fmt.Println(e.ExposureBias())

		fmt.Println(e.Artist())
		fmt.Println(e.Copyright())

		fmt.Println(e.CameraMake())
		fmt.Println(e.CameraModel())
		fmt.Println(e.CameraSerial())

		fmt.Println(e.LensMake())
		fmt.Println(e.LensModel())
		fmt.Println(e.LensSerial())

		// Example Tags
		fmt.Println(e.Dimensions())

		// Makernote Tags
		fmt.Println(e.CanonCameraSettings())
		fmt.Println(e.CanonFileInfo())
		fmt.Println(e.CanonShotInfo())
		fmt.Println(e.CanonAFInfo())

		// Time Tags
		fmt.Println(e.DateTime(time.Local))
		fmt.Println(e.ModifyDate(time.Local))
		fmt.Println(e.GPSDate(time.UTC))

		// GPS Tags
		fmt.Println(e.GPSCoords())
		fmt.Println(e.GPSAltitude())
		fmt.Println(e.GPSCoords())
		c, _ := e.GPSCellID()
		fmt.Println(c.ToToken())

		// Other Tags
		fmt.Println(e.ExposureProgram())
		fmt.Println(e.MeteringMode())
		fmt.Println(e.ShutterSpeed())
		fmt.Println(e.Aperture())
		fmt.Println(e.FocalLength())
		fmt.Println(e.FocalLengthIn35mmFilm())
		fmt.Println(e.ISOSpeed())
		fmt.Println(e.Flash())
		fmt.Println(e.ExposureValue())
		fmt.Println(e.ExposureBias())
	}
}
```

## Imagehash
Comparison between PHash and PHashFast
```go
name      old time/op    new time/op    delta
PHash-12     400µs ± 8%     203µs ± 6%  -49.25%  (p=0.000 n=19+20)

name      old alloc/op   new alloc/op   delta
PHash-12     193kB ± 0%       6kB ± 0%  -96.81%  (p=0.000 n=19+19)

name      old allocs/op  new allocs/op  delta
PHash-12     4.68k ± 0%     0.13k ± 0%  -97.24%  (p=0.000 n=20+20)
```

## Contributing

Suggestions and pull requests are welcome.

## Benchmarks

This was benchmarked without the retrival of values.
To run your own benchmarks see bench_test.go

```go
BenchmarkImageMeta/.CR2/6D-12         			   80455	     20820 ns/op	   10218 B/op	      22 allocs/op
BenchmarkImageMeta/.CR2/7DMkII-12     	 		   75770	     17921 ns/op	   10219 B/op	      21 allocs/op
BenchmarkImageMeta/.CR2/CanonM6MkII-12         	   63640	     18128 ns/op	    9261 B/op	      21 allocs/op
BenchmarkImageMeta/.CR3/CanonR-12              	   72705	     16597 ns/op	    9236 B/op	      21 allocs/op
BenchmarkImageMeta/.CR3/CanonRP-12             	   59492	     17623 ns/op	    9216 B/op	      20 allocs/op
BenchmarkImageMeta/.CR3/CanonR3-12             	   57465	     20322 ns/op	   14632 B/op	      22 allocs/op
BenchmarkImageMeta/.CR3/CanonR5-12             	   55514	     19874 ns/op	   14632 B/op	      22 allocs/op
BenchmarkImageMeta/.CR3/CanonR6-12             	   76044	     16234 ns/op	    9234 B/op	      21 allocs/op
BenchmarkImageMeta/.HEIC/CanonR5-12            	   32150	     43434 ns/op	   10198 B/op	      21 allocs/op
BenchmarkImageMeta/.HEIC/CanonR6-12            	   34396	     37288 ns/op	   10196 B/op	      21 allocs/op
BenchmarkImageMeta/.JPG/GPS-12                 	  117636	      9807 ns/op	     280 B/op	       4 allocs/op
BenchmarkImageMeta/.JPG/GoPro6-12              	  154758	      6528 ns/op	     280 B/op	       4 allocs/op
BenchmarkImageMeta/.HEIC-12                    	    8442	    194884 ns/op	    4540 B/op	      15 allocs/op
BenchmarkImageMeta/.HEIC/iPhone11-12           	    6170	    185167 ns/op	    4569 B/op	      15 allocs/op
BenchmarkImageMeta/.HEIC/iPhone12-12           	   26725	     48433 ns/op	    1716 B/op	      11 allocs/op
BenchmarkImageMeta/.HEIC/iPhone13-12           	    6726	    152748 ns/op	    4561 B/op	      15 allocs/op
BenchmarkImageMeta/.NEF/Nikon-12               	   54534	     22582 ns/op	   10242 B/op	      23 allocs/op
BenchmarkImageMeta/.NEF/Nikon#01-12            	   53985	     23082 ns/op	   10244 B/op	      23 allocs/op
BenchmarkImageMeta/.RW2/Panasonic-12           	   52309	     20677 ns/op	    4555 B/op	      15 allocs/op
BenchmarkImageMeta/.ARW/Sony-12                	   99100	     11576 ns/op	    4769 B/op	      22 allocs/op
BenchmarkImageMeta/.WEBP/Webp-12               	 4978038	     285.9 ns/op	      24 B/op	       1 allocs/op
BenchmarkImageMeta/.DNG/Adobe-12               	   34141	     35345 ns/op	   20866 B/op	      30 allocs/op
BenchmarkImageMeta/.JPG/NoExif-12              	 1290793	     942.1 ns/op	     280 B/op	       4 allocs/op

```

## Imagetype Identification

Images can be identified with: "github.com/evanoberholster/imagemeta/imagetype" package.

Benchmarks can be found with the imagemeta/imagetype package

Example:

```go
package main

import (
   "fmt"
   "os"

   "github.com/evanoberholster/imagemeta/imagetype"
)

const imageFilename = "../../test/img/1.CR2"

func main() {
   var err error

   f, err := os.Open(imageFilename)
   if err != nil {
      panic(err)
   }
   defer f.Close()
   t, err := imagetype.Scan(f)
   if err != nil {
      panic(err)
   }
   fmt.Println(t)
}
```

## TODO

- [x] Stabilize ImageTypes API
- [x] Add Exif parsing for individual image types (jpg, heic, cr2, tiff, dng)
- [x] Add XMP parsing as "xmp" package
- [x] Add Avif, Heic and CR3 image metadata support (riff format images)
- [ ] Stabalize Imagemeta API
- [ ] Improve test coverage
- [ ] Create Thumbnail API
- [ ] Add Webp image metadata support
- [ ] Add Canon Exif Makernote support
- [ ] Add Nikon Exif Makernote support
- [ ] Add CRW image metadata support (ciff format images)
- [ ] Documentation

## Based on and Inspired by

Based on work by Dustin Oprea [https://github.com/dsoprea/go-exif](https://github.com/dsoprea/go-exif)

Inspired by Phil Harvey [http://exiftool.org](http://exiftool.org)

Some inspiration from RW Carlsen [https://github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif)

## Special Thanks to:
- The go4 Authors (https://github.com/go4org/go4) for their work on a BMFF parser and HEIF structure in golang.
- Laurent Clévy (@Lorenzo2472) (https://github.com/lclevy/canon_cr3) for Canon CR3 structure.
- Lasse Heikkilä (https://trepo.tuni.fi/bitstream/handle/123456789/24147/heikkila.pdf) for HEIF structure from his thesis.

## LICENSE

Copyright (c) 2020-2021, Evan Oberholster & Contributors

Copyright (c) 2019, Dustin Oprea

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
