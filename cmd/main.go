package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
)

const testFilename = "../../test/img/1.heic"

const testFilename2 = "../testImages/Heic.exif"

func main() {
	f, err := os.Open(testFilename2)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println(testFilename)
	//var x xmp.XMP
	var e *exif.Data
	exifDecodeFn := func(r io.Reader, header exif.Header) error {
		e, err = exif.ParseExif(f, header)
		fmt.Println(e, err, header)
		return nil
	}
	//xmpDecodeFn := func(r io.Reader, header xmp.Header) error {
	//	fmt.Println(header)
	//	var err error
	//	x, err = xmp.ParseXmp(r)
	//	fmt.Println(x, err)
	//	return err
	//}
	start := time.Now()
	m, err := imagemeta.NewMetadata(f, nil, exifDecodeFn)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Println(m.Dimensions())
	fmt.Println(m)
	fmt.Println(elapsed)
	if e != nil {
		fmt.Println(e.Artist())
		fmt.Println(e.CameraMake())
		fmt.Println(e.CameraModel())
		fmt.Println(e.ISOSpeed())
		fmt.Println(e.FocalLength())
		fmt.Println(e.LensModel())
		fmt.Println(e.Aperture())
		fmt.Println(e.ShutterSpeed())

		fmt.Println(e.ExposureValue())
		fmt.Println(e.ExposureBias())

		fmt.Println(e.GPSCoords())
		fmt.Println(e.DateTime())
		//fmt.Println(e.ModifyDate())

		fmt.Println(e.GPSDate(nil))

		start = time.Now()
		//for t := range e.RangeTags() {
		//	//ifds.ExifIFD, 0, exififd.ISOSpeedRatings
		//	//if t.TagID == exififd.DateTimeDigitized {
		//	if a, ok := ifds.RootIfdTagIDMap[t.TagID]; ok {
		//		fmt.Println(t.TagID, a, t.UnitCount, t.ValueOffset)
		//	} else if a, ok := exififd.TagIDMap[t.TagID]; ok {
		//		fmt.Println(t.TagID, a, t.UnitCount, t.ValueOffset)
		//	} else if a, ok := gpsifd.TagIDMap[t.TagID]; ok {
		//		fmt.Println(t.TagID, a, t.UnitCount, t.ValueOffset)
		//	}
		//	//}
		//}
	}
	//elapsed = time.Since(start)
	//fmt.Println(elapsed)

}
