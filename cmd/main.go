package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/meta"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: main <file>\n")
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	//var x xmp.XMP
	var e *exif.Data
	exifDecodeFn := func(r io.Reader, header meta.ExifHeader) error {
		e, err = e.ParseExif(f, header)
		//fmt.Println("Item", e, err, header)
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
		fmt.Println(err, "here")
	}
	elapsed := time.Since(start)
	fmt.Println(m.Dimensions())
	fmt.Println(m)
	//fmt.Println(*e)
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

		c, _ := e.GPSCellID()
		fmt.Println(c.ToToken())
		fmt.Println(e.DateTime())
		//fmt.Println(e.ModifyDate())

		fmt.Println(e.GPSDate(nil))

		start = time.Now()
	}
	//elapsed = time.Since(start)
	//fmt.Println(elapsed)

}
