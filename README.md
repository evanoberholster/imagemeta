# Exif Tool

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Coverage Status][Coverage-Image]][Coverage-Url]
[![Build][Build-Status-Image]][Build-Status-Url]

This package provides for performance oriented decoding of exif and tiff encoded data.

See ([Godoc-Url])[Godoc-Url] for more doumentation.

Suggestions and pull requests are welcome.

Example usage:

```go
package main

import (
   "bytes"
   "fmt"
   "io/ioutil"
   "os"
   "time"

   "github.com/evanoberholster/exiftool"
)

const testFilename = "image.jpg"

func main() {
   var err error

   f, err := os.Open(testFilename)
   if err != nil {
      panic(err)
   }
   defer f.Close()

   start := time.Now()
   e, err := exiftool.ScanExif(f)
   if err != nil && err != exiftool.ErrNoExif {
      panic(err)
   }
   elapsed := time.Since(start)
   fmt.Println(elapsed)

   if err == exiftool.ErrNoExif {
      fmt.Println(e.XMLPacket())
      fmt.Println(e.Dimensions())
      return
   }

   // Strings
   fmt.Println(e.CameraMake())
   fmt.Println(e.CameraModel())
   fmt.Println(e.Artist())
   fmt.Println(e.Copyright())
   fmt.Println(e.LensMake())
   fmt.Println(e.LensModel())
   fmt.Println(e.CameraSerial())
   fmt.Println(e.LensSerial())
   //
   fmt.Println(e.Dimensions())
   fmt.Println(e.XMLPacket())
   //
   //// Makernotes
   fmt.Println(e.CanonCameraSettings())
   fmt.Println(e.CanonFileInfo())
   fmt.Println(e.CanonShotInfo())
   fmt.Println(e.CanonAFInfo())
   //
   //// Time
   fmt.Println(e.ModifyDate())
   fmt.Println(e.DateTime())
   fmt.Println(e.GPSTime())
   //
   //// GPS
   fmt.Println(e.GPSInfo())
   fmt.Println(e.GPSAltitude())
   //
   // Metadata
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

Images can be identified with: "github.com/evanoberholster/exiftool/imagetype" package.

Example:

```go
package main

import (
   "fmt"
   "os"

   "github.com/evanoberholster/exiftool/imagetype"
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

Benchmarks can be found with the exiftool/imagetype package

## TODO

- [x] Update ImageTypes API
- [-] Write Exif extraction for individual image types (jpg, heic)
- [ ] Write tests
- [ ] Include support for CRW image type (ciff format images)
- [ ] Create Thumbnail API
- [ ] Stabalize API
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
[Godoc-Url]: https://godoc.org/github.com/evanoberholster/exiftool
[Godoc-Image]: https://godoc.org/github.com/evanoberholster/exiftool?status.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/evanoberholster/exiftool
[ReportCard-Image]: https://goreportcard.com/badge/github.com/evanoberholster/exiftool
[Build-Status-Url]: https://travis-ci.com/evanoberholster/exiftool?branch=master
[Build-Status-Image]: https://travis-ci.com/evanoberholster/exiftool.svg?branch=master
