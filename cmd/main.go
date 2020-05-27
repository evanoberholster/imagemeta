package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/evanoberholster/exiftool"
)

const testFilename = "../../test/img/3.heic"

func main() {
	var err error

	f, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}

	eh, err := exiftool.SearchExifHeader(f)
	if err != nil {
		panic(err)
	}
	f.Seek(0, 0)

	//cb := bufra.NewBufReaderAt(f, 128*1024)
	buf, _ := ioutil.ReadAll(f)
	cb := bytes.NewReader(buf)
	start := time.Now()
	e, err := eh.ParseExif(cb)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(e)
	//start := time.Now()

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
