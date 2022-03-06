package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/heic"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
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
	parseJPG(f)
	//parseCR2(f)
	//parseHeic(f)
	//parseCR3(f)
}

func parseJPG(f meta.Reader) {
	exifFn := func(r io.Reader, header meta.ExifHeader) error {

		fmt.Println(header)
		fmt.Println(header.TiffHeaderOffset, header.FirstIfdOffset, header.FirstIfdOffset-header.TiffHeaderOffset)
		return nil
	}

	_, err := jpeg.ScanJPEG(f, exifFn, nil)
	fmt.Println(err)
}

func parseCR2(f meta.Reader) {

	exifFn := func(r io.Reader, header meta.ExifHeader) error {

		fmt.Println(header)
		fmt.Println(header.TiffHeaderOffset, header.FirstIfdOffset, header.FirstIfdOffset-header.TiffHeaderOffset)
		return nil
	}

	err := tiff.Scan(f, imagetype.ImageCR2, exifFn, nil)
	fmt.Println(err)
}

func parseHeic(f meta.Reader) {
	var err error
	m := &meta.Metadata{}

	var e *exif.Data
	var x xmp.XMP

	m.ExifFn = func(r io.Reader, m *meta.Metadata) error {
		e, err = exif.ParseExif(f, m.ExifHeader)
		return nil
	}
	m.XmpFn = func(r io.Reader, m *meta.Metadata) error {
		x, err = xmp.ParseXmp(r)
		return nil
	}
	hm, err := heic.NewMetadata(f, m)
	if err != nil {

		fmt.Println(err)
		// Error retrieving Heic Metadata
	}
	_, err = hm.ReadExifHeader(f)
	if err != nil {
		fmt.Println(err)
	}
	if err != meta.ErrNoExif {
		if err = m.ExifFn(f, m); err != nil {
			panic(err)
		}
		printJSON(e)
		printExif(e)
	}

	if _, err = hm.ReadXmpHeader(f); err == nil {
		_, err = f.Seek(int64(hm.XmpHeader.Offset), 0)
		if err = m.XmpFn(f, m); err != nil {
			panic(err)
		}
		fmt.Println(x)
	}

	fmt.Println(m.XmpHeader)
	fmt.Println(m.ExifHeader)
	fmt.Println(m.It, m.Dim)
}

//
//func parseCR3(f meta.Reader) {
//	m, err := cr3.NewMetadata(f, meta.Metadata{})
//	fmt.Println(m, err)
//	var XMP xmp.XMP
//	m.XmpDecodeFn = func(r io.Reader, header meta.XmpHeader) error {
//		start := time.Now()
//		XMP, err = xmp.ParseXmp(r)
//		elapsed := time.Since(start)
//		fmt.Println(elapsed)
//		return err
//	}
//
//	var e *exif.Data
//	m.ExifDecodeFn = func(r io.Reader, header meta.ExifHeader) (err error) {
//		e, err = e.ParseExif(f, header)
//		return err
//	}
//
//	if err = m.DecodeExif(f); err != nil {
//		fmt.Println(err)
//	}
//	if err = m.DecodeXMP(f); err != nil {
//		fmt.Println(err)
//	}
//	_ = XMP
//
//	printExif(e)
//}

func printJSON(e *exif.Data) {
	buf, err := e.DebugJSON()
	fmt.Println(string(buf), err)
}

func printExif(e *exif.Data) {
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

		fmt.Println(e.DateTime(nil))
	}
}
