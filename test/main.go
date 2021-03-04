package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/cr3"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/heic"
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
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	stdLogger := bmff.STDLogger{}
	bmff.DebugLogger(stdLogger)
	//parseHeic(f)
	parseCR3(f)
}

func parseHeic(f meta.Reader) {
	hm, err := heic.NewMetadata(f)
	if err != nil {
		fmt.Println(err)
		// Error retrieving Heic Metadata
	}
	var e *exif.Data
	hm.ExifDecodeFn = func(r io.Reader, header exif.Header) error {
		e, err = exif.ParseExif(f, header)
		fmt.Println(e, err, header)
		return nil
	}
	err = hm.DecodeExif(f)
	fmt.Println(e, err)
	PrintExif(e)
}

func parseCR3(f meta.Reader) {
	m, err := cr3.NewMetadata(f)
	fmt.Println(m, err)
	var XMP xmp.XMP
	m.XmpDecodeFn = func(r io.Reader, header xmp.Header) error {
		start := time.Now()
		XMP, err = xmp.ParseXmp(r)
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		return err
	}

	e := &exif.Data{}
	m.ExifDecodeFn = func(r io.Reader, header exif.Header) (err error) {
		err = e.ParseExif(f, header)
		return err
	}

	m.DecodeExif(f)
	m.DecodeXMP(f)
	_ = XMP

	PrintExif(e)
}

func PrintExif(e *exif.Data) {
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

		fmt.Println(e.GPSDate(nil))
		fmt.Println(e.GPSCoords())

		fmt.Println(e.DateTime())
	}
}
