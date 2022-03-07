package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
)

func main() {
	//flag.Parse()
	//if flag.NArg() != 1 {
	//	fmt.Fprintf(os.Stderr, "usage: main <file>\n")
	//	os.Exit(1)
	//}
	//f, err := os.Open(flag.Arg(0))
	f, err := os.Open("../../test/img/1.NEF")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	exif.InfoLogger = log.New(os.Stdout, "", log.Ltime)
	//var x xmp.XMP
	var e *exif.Data
	//exifDecodeFn := func(r io.Reader, m *meta.Metadata) error {
	//	e, err = e.ParseExifWithMetadata(f, m)
	//	return nil
	//}
	//xmpDecodeFn := func(r io.Reader, m *meta.Metadata) error {
	//	x, err = xmp.ParseXmp(r)
	//	return err
	//}

	exifFn := func(r io.Reader, header meta.ExifHeader) error {
		_, _ = f.Seek(0, 0)
		fmt.Println(header)
		e, err = exif.ParseExif(f, header)
		fmt.Println(e)
		return err
	}

	err = tiff.Scan(f, imagetype.ImageTiff, exifFn, nil)
	if err != nil {
		panic(err)
	}
	//m, err := imagemeta.NewMetadata(f, xmpDecodeFn, exifDecodeFn)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(m.Metadata)
	//fmt.Println(x)
	if e != nil {
		fmt.Println(e.ImageWidth())
		fmt.Println(e.ImageHeight())

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

		fmt.Println(e.ExposureValue())
		fmt.Println(e.ExposureBias())

		fmt.Println(e.GPSCoords())

		//c, _ := e.GPSCellID()
		//fmt.Println(c.ToToken())
		fmt.Println(e.DateTime(time.Local))
		fmt.Println(e.ModifyDate(time.Local))

		//fmt.Println(e.GPSDate(nil))
	}

	b, err := e.DebugJSON()
	fmt.Println(string(b), err)
}
