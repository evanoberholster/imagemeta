package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/tiff"
	"github.com/evanoberholster/imagemeta/xmp"
)

func BenchmarkExif(b *testing.B) {
	f, err := os.Open("../../test/img/2.CR2")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	buf, _ := ioutil.ReadAll(f)
	cb := bytes.NewReader(buf)
	b.ResetTimer()
	b.ReportAllocs()

	var x xmp.XMP
	var e *exif.Data
	exifFn := func(r io.Reader, h meta.ExifHeader) error {
		cb.Seek(0, 0)
		e, err = exif.ParseExif(cb, h)
		return nil
	}
	//exifDecodeFn := func(r io.Reader, m *meta.Metadata) error {
	//	e, err = exif.ParseExif2(f, m)
	//	return nil
	//}
	//xmpDecodeFn := func(r io.Reader, m *meta.Metadata) error {
	//	x, err = xmp.ParseXmp(r)
	//	return err
	//}

	for i := 0; i < b.N; i++ {
		cb.Seek(0, 0)
		err := tiff.Scan(cb, imagetype.ImageTiff, exifFn, nil)
		//	m, err := imagemeta.NewMetadata(f, xmpDecodeFn, exifDecodeFn)
		if err != nil {
			fmt.Println(err)
		}
		//_ = m.Metadata
		_ = x
		if e != nil {

			_, _ = e.Artist()
			_, _ = e.Copyright()
			_ = e.CameraMake()
			_ = e.CameraModel()
			_, _ = e.CameraSerial()
			_, _ = e.Orientation()
			_, _ = e.LensMake()
			_, _ = e.LensModel()
			_, _ = e.LensSerial()
			_, _ = e.ISOSpeed()
			_, _ = e.FocalLength()
			_, _ = e.LensModel()
			_, _ = e.Aperture()
			_, _ = e.ShutterSpeed()
			_, _ = e.ExposureValue()
			_, _ = e.ExposureBias()

			_, _, _ = e.GPSCoords()
			c, _ := e.GPSCellID()
			_ = c.ToToken()
			_, _ = e.DateTime(time.Local)
			_, _ = e.ModifyDate(time.Local)

			//fmt.Println(e.GPSDate(nil))
		}
	}

}
