# Exif Tool

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]
[![Coverage Status](https://coveralls.io/repos/github/evanoberholster/exiftool/badge.svg?branch=master)](https://coveralls.io/github/evanoberholster/exiftool?branch=master)
[![Build][Build-Status-Image]][Build-Status-Url]

Provides decoding of basic exif and tiff encoded data.

Suggestions and pull requests are welcome.

Example usage:

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/evanoberholster/exiftool"
)

const imageFilename = "../../test/img/3.heic"

func main() {
	var err error

	f, err := os.Open(imageFilename)
	if err != nil {
	    panic(err)
	}

	eh, err := exiftool.SearchExifHeader(f)
	if err != nil {
	    panic(err)
	}
	f.Seek(0, 0)
	buf, _ := ioutil.ReadAll(f)
	cb := bytes.NewReader(buf)

	e, err := eh.ParseExif(cb)
	if err != nil {
	    fmt.Println(err)
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
	fmt.Println(e.XMLPacket())

	// Canon Makernotes
	fmt.Println(e.CanonCameraSettings())
	fmt.Println(e.CanonFileInfo())
	fmt.Println(e.CanonShotInfo())
	fmt.Println(e.CanonAFInfo())

	// Time
	fmt.Println(e.ModifyDate())
	fmt.Println(e.DateTime())
	fmt.Println(e.GPSTime())

	// GPS
	fmt.Println(e.GPSInfo())
	fmt.Println(e.GPSAltitude())

    // Metadata
	fmt.Println(e.Dimensions())
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

## Bencmarks

This was benchmarked without the retrival of values.
To run your own benchmarks see exifHeader_test.go

```go
BenchmarkParseExif100/.CR2/GPS-8         	   23035	     54230 ns/op	    9310 B/op	      56 allocs/op
BenchmarkParseExif100/.CR2/7D-8          	   24091	     48405 ns/op	    8957 B/op	      54 allocs/op
BenchmarkParseExif100/.CR3-8             	  176527	      6746 ns/op	     901 B/op	      14 allocs/op
BenchmarkParseExif100/.JPG/GPS-8         	   47270	     25754 ns/op	    5123 B/op	      32 allocs/op
BenchmarkParseExif100/.HEIC-8            	   50145	     25194 ns/op	    4882 B/op	      29 allocs/op
BenchmarkParseExif100/.GoPro/6-8         	   54031	     22543 ns/op	    3782 B/op	      28 allocs/op
BenchmarkParseExif100/.NEF/Nikon-8       	   22287	     54464 ns/op	   12417 B/op	      59 allocs/op
BenchmarkParseExif100/.ARW/Sony-8        	   28357	     42439 ns/op	    7671 B/op	      53 allocs/op
BenchmarkParseExif100/.DNG/Adobe-8       	   13603	     87055 ns/op	   18494 B/op	      87 allocs/op
```

## TODO

- Example Usage
- Create Thumbnail API
- Update ImageTypes API
- Stabalize API
- Documentation

## Based on and Inspired by

Significantly based on work by Dustin Oprea [https://github.com/dsoprea/go-exif](https://github.com/dsoprea/go-exif)

Inspired by Phil Harvey [http://exiftool.org](http://exiftool.org)

Some inspiration from RW Carlsen [https://github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif)

## LICENSE

Copyright (c) 2019, Evan Oberholster & Contributors

Copyright (c) 2019, Dustin Oprea

[License-Url]: https://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/badge/License-MIT-blue.svg?maxAge=2592000
[Godoc-Url]: https://godoc.org/github.com/evanoberholster/exiftool
[Godoc-Image]: https://godoc.org/github.com/evanoberholster/exiftool?status.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/evanoberholster/exiftool
[ReportCard-Image]: https://goreportcard.com/badge/github.com/evanoberholster/exiftool
[Build-Status-Url]: https://travis-ci.com/evanoberholster/exiftool?branch=master
[Build-Status-Image]: https://travis-ci.com/evanoberholster/exiftool.svg?branch=master