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

	var x xmp.XMP
	var e *exif.Data
	exifDecodeFn := func(r io.Reader, m *meta.Metadata) error {
		e, err = e.ParseExifWithMetadata(f, m)
		return nil
	}
	xmpDecodeFn := func(r io.Reader, m *meta.Metadata) error {
		x, err = xmp.ParseXmp(r)
		return err
	}

	m, err := imagemeta.NewMetadata(f, xmpDecodeFn, exifDecodeFn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m.Metadata)
	fmt.Println(x)
	if e != nil {
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
	  fmt.Println(e.ModifyDate())
	  fmt.Println(e.DateTime())
	  fmt.Println(e.GPSDate(time.UTC))
    fmt.Println(e.GPSDate(nil))
	
	  // GPS Tags
	  fmt.Println(e.GPSCoords())
	  fmt.Println(e.GPSAltitude())
	
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

		fmt.Println(e.GPSCoords())
		c, _ := e.GPSCellID()
		fmt.Println(c.ToToken())		
	}
}
```

## Contributing

Suggestions and pull requests are welcome.

## Benchmarks

This was benchmarked without the retrival of values.
To run your own benchmarks see bench_test.go

```go
BenchmarkImagemeta100/.CR2/60D         	   15552	     69759 ns/op	   11281 B/op	      30 allocs/op
BenchmarkImagemeta100/.CR2/GPS         	   14929	     83673 ns/op	   11332 B/op	      32 allocs/op
BenchmarkImagemeta100/.CR2/7D          	   14068	     83297 ns/op	   11333 B/op	      32 allocs/op
BenchmarkImagemeta100/.JPG/GPS         	   19404	     59828 ns/op	    5629 B/op	      24 allocs/op
BenchmarkImagemeta100/.JPF/GoPro6      	   25165	     47939 ns/op	    5607 B/op	      24 allocs/op
BenchmarkImagemeta100/.HEIC            	   31309	     35531 ns/op	   12608 B/op	      76 allocs/op
BenchmarkImagemeta100/.HEIC/CanonR5    	   14814	     81757 ns/op	   17655 B/op	      67 allocs/op
BenchmarkImagemeta100/.HEIC/CanonR6    	   13557	     86212 ns/op	   17367 B/op	      65 allocs/op
BenchmarkImagemeta100/.HEIC/iPhone11   	   15961	     68563 ns/op	   17101 B/op	      93 allocs/op
BenchmarkImagemeta100/.HEIC/iPhone12   	   18288	     64648 ns/op	   14561 B/op	      94 allocs/op
BenchmarkImagemeta100/.NEF/Nikon       	   13954	     84132 ns/op	   11443 B/op	      34 allocs/op
BenchmarkImagemeta100/.NEF/Nikon#01    	   14503	     96479 ns/op	   11442 B/op	      34 allocs/op
BenchmarkImagemeta100/.RW2/Panasonic   	   22882	     51855 ns/op	    5558 B/op	      22 allocs/op
BenchmarkImagemeta100/.ARW/Sony        	   16862	     71868 ns/op	    5856 B/op	      31 allocs/op
BenchmarkImagemeta100/.DNG/Adobe       	    8196	    131020 ns/op	   22138 B/op	      43 allocs/op
BenchmarkImagemeta100/.JPG/NoExif      	 1609324	       756 ns/op	     352 B/op	       3 allocs/op

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
