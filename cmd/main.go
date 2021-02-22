package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/xmp"
)

const testFilename = "../../test/img/10.jpg"

func main() {
	f, err := os.Open(testFilename)
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
	var x xmp.XMP
	exifDecodeFn := func(r io.Reader, header exif.Header) error {
		exif, err := exif.ParseExif(f, header)
		fmt.Println(exif, err)
		return nil
	}
	xmpDecodeFn := func(r io.Reader, header xmp.Header) error {
		fmt.Println(header)
		var err error
		x, err = xmp.ParseXmp(r)
		fmt.Println(x, err)
		return err
	}
	start := time.Now()
	m, err := imagemeta.NewMetadata(f, xmpDecodeFn, exifDecodeFn)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Println(m.Dimensions())
	fmt.Println(m)
	fmt.Println(elapsed)

}
