package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
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
		if filepath.Ext(name) == ".html" || filepath.Ext(name) == ".CRW" || filepath.Ext(name) == ".jp2" {
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
				fmt.Println(err)
			}

		}
		_ = exif
		fmt.Printf("\t%s\t\t%s\n", exif.Make, exif.Model)
		//fmt.Println(string(exif.ApplicationNotes))
		//fmt.Println(len(exif.ApplicationNotes))
	}
}

func main3() {
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("iPad2022_1.jpeg")

	f, err := os.Open(dir + "/" + "DJI.dng")
	//f, err := os.Open(dir + "/" + "iPhone11.heic")
	//f, err := os.Open(dir + "/" + "8.heic")
	//f, err := os.Open(dir + "/" + "2.CR2")
	//f, err := os.Open("test123.jpg")
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
