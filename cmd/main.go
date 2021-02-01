// Package main provides an example command
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
)

const testFilename = "../../test/img/2.CR2"

func main() {
	var err error

	f, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}

	buf, _ := ioutil.ReadAll(f)
	reader := bytes.NewReader(buf)

	start := time.Now()
	e, err := imagemeta.ScanExif(reader)
	if err != nil && err != imagemeta.ErrNoExif {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed)

	if err == imagemeta.ErrNoExif {
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

	fmt.Println(time.Since(start))
}
