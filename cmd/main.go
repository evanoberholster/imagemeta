package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
)

func main() {
	f, err := os.Open("../testImages/Heic.exif")
	if err != nil {
		log.Fatal(err)
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
	//b, err := e.DebugJSON()
	//fmt.Println(string(b), err)
}
