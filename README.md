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
	"fmt"
	"io"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/xmp"
)

const testFilename = "image.jpg"

func main() {
	f, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	var x xmp.XMP
	var e exif.Exif
   
   // Exif Decode function
	exifDecodeFn := func(r io.Reader, header exif.Header) error {
		e, err = exif.ParseExif(f, header)
		return nil
	}
   // Xmp Decode function
	xmpDecodeFn := func(r io.Reader, header xmp.Header) error {
		fmt.Println(header)
		var err error
		x, err = xmp.ParseXmp(r)
		return err
	}
	start := time.Now()
	m, err := imagemeta.NewMetadata(f, xmpDecodeFn, exifDecodeFn)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	fmt.Println(m.Dimensions())
	fmt.Println(m)
	fmt.Println(elapsed)

	// XMP
	fmt.Println(x)

	// Exif
	fmt.Println(e)

    // Example Tags
	fmt.Println(e.CameraMake())
	fmt.Println(e.CameraModel())
	fmt.Println(e.Artist())
	fmt.Println(e.Copyright())
	fmt.Println(e.LensMake())
	fmt.Println(e.LensModel())
	fmt.Println(e.CameraSerial())
	fmt.Println(e.LensSerial())
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

}
```

## Contributing

Suggestions and pull requests are welcome.

## Benchmarks

This was benchmarked without the retrival of values.
To run your own benchmarks see benchmark_test.go

```go
BenchmarkScanExif100/NoExif.jpg-8              1249808          996 ns/op       4496 B/op          5 allocs/op
BenchmarkScanExif100/350D.CR2-8                  39280        31519 ns/op      10445 B/op         46 allocs/op
BenchmarkScanExif100/XT1.CR2-8                   37731        31582 ns/op      10444 B/op         46 allocs/op
BenchmarkScanExif100/60D.CR2-8                   27439        43459 ns/op      12593 B/op         52 allocs/op
BenchmarkScanExif100/6D.CR2-8                    26264        45286 ns/op      13185 B/op         57 allocs/op
BenchmarkScanExif100/7D.CR2-8                    26625        46062 ns/op      13216 B/op         57 allocs/op
BenchmarkScanExif100/5DMKIII.CR2-8               24404        48578 ns/op      13212 B/op         57 allocs/op
BenchmarkScanExif100/1.CR3-8                    138854         8470 ns/op       5157 B/op         17 allocs/op
BenchmarkScanExif100/1.jpg-8                     52980        22424 ns/op      31394 B/op         32 allocs/op
BenchmarkScanExif100/1.NEF-8                     24420        50230 ns/op      13598 B/op         61 allocs/op
BenchmarkScanExif100/3.NEF-8                     20294        58299 ns/op      17008 B/op         67 allocs/op
BenchmarkScanExif100/1.ARW-8                     30277        39593 ns/op      11928 B/op         56 allocs/op
BenchmarkScanExif100/4.RW2-8                     34719        34740 ns/op       8202 B/op         31 allocs/op
BenchmarkScanExif100/hero6.gpr-8                 31630        38285 ns/op      13606 B/op         39 allocs/op
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
- [ ] Add Avif, Heic and CR3 image metadata support (riff format images)
- [ ] Add Webp image metadata support
- [ ] Add Canon Exif Makernote support
- [ ] Add CRW image metadata support (ciff format images)
- [ ] Stabalize Imagemeta API
- [ ] Write tests for stable packages
- [ ] Create Thumbnail API
- [ ] Documentation

## Based on and Inspired by

Significantly based on work by Dustin Oprea [https://github.com/dsoprea/go-exif](https://github.com/dsoprea/go-exif)

Inspired by Phil Harvey [http://exiftool.org](http://exiftool.org)

Some inspiration from RW Carlsen [https://github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif)

## LICENSE

Copyright (c) 2020, Evan Oberholster & Contributors

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
