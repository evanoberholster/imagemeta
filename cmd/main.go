package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/isobmff"
	"github.com/rs/zerolog"
)

func init() {
	imagemeta.SetLogger(zerolog.ConsoleWriter{Out: os.Stdout}, zerolog.WarnLevel)
	//exif2.Logger = exif2.Logger.Level(zerolog.DebugLevel)
	//isobmff.Logger = isobmff.Logger.Level(zerolog.DebugLevel)
}

var (
	dir = "../../test/img/"
)

func main() {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		//if f.Size() > 1*1024*1024 {
		name := f.Name()
		if !(filepath.Ext(name) == ".CR3" || filepath.Ext(name) == ".CR2" || filepath.Ext(name) == ".ARW" || filepath.Ext(name) == ".NEF" || filepath.Ext(name) == ".HIF" || filepath.Ext(name) == ".GPR" || filepath.Ext(name) == ".RW2") {
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
		fmt.Print(name + "\t\t")
		//fmt.Println(f.Name())
		exif, err = imagemeta.Decode(r)
		if err != nil {
			if err != imagemeta.ErrNoExif {
				//fmt.Println(err)
				panic(err)
			}

		}
		_ = exif
		fmt.Printf("\t%s\t\t%s\n", exif.Make(), exif.Model())
		//fmt.Println(string(exif.ApplicationNotes))
		//fmt.Println(len(exif.ApplicationNotes))
	}
}

func main3() {
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("test123.jpg")

	//f, err := os.Open(dir + "/" + "iPhone13.heic")
	f, err := os.Open(dir + "/" + "CanonR6_1.HIF")
	//f, err := os.Open(dir + "/" + "2.CR2")
	//f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}

func main2() {
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("img1.heif")
	//f, err := os.Open(dir + "/" + "CanonSL3_1.CR3")
	//f, err := os.Open(dir + "/" + "14.JPG")
	//f, err := os.Open(dir + "/" + "CanonR5_1.HIF")
	f, err := os.Open(dir + "/" + "iPhone11.heic")

	//f, err := os.Open(dir + "/" + "5.heic")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	start := time.Now()
	ir := exif2.NewIfdReader(exif2.Logger)
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
	if err := r.ReadMetadata(); err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
	fmt.Println(r)
	fmt.Println(ir.Exif)
	//e, err := imagemeta.Decode(f)
	//if err != nil {
	//	panic(err)
	//}
	//_ = e
	//fmt.Println(e)
}
