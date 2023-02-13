package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/isobmff"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	imagemeta.SetLogger(zerolog.ConsoleWriter{Out: os.Stdout}, zerolog.TraceLevel)
	exif2.Logger = exif2.Logger.Level(zerolog.DebugLevel)
	isobmff.Logger = isobmff.Logger.Level(zerolog.DebugLevel)
}

func main2() {
	exif2.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
	jpeg.Logger = exif2.Logger
	bmff.Logger = exif2.Logger
	dir := "../../test/img/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		//if f.Size() > 1*1024*1024 {
		name := f.Name()
		if !(filepath.Ext(name) == ".CR3" || filepath.Ext(name) == ".CR2") {
			continue
		}
		//fmt.Println(filepath.Ext(name))
		r, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			panic(err)
		}
		defer func() {
			err = r.Close()
			if err != nil {
				panic(err)
			}
		}()
		var exif exif2.Exif
		fmt.Println(name)
		//fmt.Println(f.Name())
		exif, err = imagemeta.Decode(r)
		if err != nil {
			if err != imagemeta.ErrNoExif {
				//fmt.Println(err)
				panic(err)
			}

		}
		_ = exif

		fmt.Println(exif.Model)
		//fmt.Println(string(exif.ApplicationNotes))
		//fmt.Println(len(exif.ApplicationNotes))
	}
}

func main() {
	exif2.Logger = exif2.Logger.Level(zerolog.TraceLevel)
	isobmff.Logger = isobmff.Logger.Level(zerolog.TraceLevel)
	dir := "../../test/img/"
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("img1.jpg")
	//f, err := os.Open(dir + "/" + "CanonR6_1.CR3")
	f, err := os.Open(dir + "/" + "iPhone12.heic")
	//f, err := os.Open(dir + "/" + "14.JPG")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	_ = e
	fmt.Println(e)
}

func main3() {
	imagemeta.SetLogger(zerolog.ConsoleWriter{Out: os.Stdout}, zerolog.TraceLevel)
	exif2.Logger = exif2.Logger.Level(zerolog.ErrorLevel)
	isobmff.Logger = isobmff.Logger.Level(zerolog.ErrorLevel)
	dir := "../../test/img/"
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("img1.heif")
	f, err := os.Open(dir + "/" + "CanonR5_1.CR3")
	//f, err := os.Open(dir + "/" + "14.JPG")
	//f, err := os.Open(dir + "/" + "CanonR5_1.HIF")
	//f, err := os.Open(dir + "/" + "iPhone12.heic")

	//f, err := os.Open(dir + "/" + "5.heic")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ir := exif2.NewIfdReader(f)
	defer ir.Close()

	r := isobmff.NewReader(f)
	r.ExifReader = ir.DecodeIfd
	defer r.Close()
	if err := r.ReadFTYP(); err != nil {
		panic(err)
	}
	if err := r.ReadMetadata(); err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(ir.Exif)
	//e, err := imagemeta.Decode(f)
	//if err != nil {
	//	panic(err)
	//}
	//_ = e
	//fmt.Println(e)
}
