package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main2() {
	dir := "../../test/img/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		//if f.Size() > 1*1024*1024 {
		name := f.Name()
		if filepath.Ext(name) != ".jpg" {
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
		//fmt.Println(f.Name())
		exif, err = imagemeta.DecodeJPEG(r)
		if err != nil {
			if err != imagemeta.ErrNoExif {
				panic(err)
			}

		}
		_ = exif
		fmt.Println(exif)
		//fmt.Println(string(exif.ApplicationNotes))
		//fmt.Println(len(exif.ApplicationNotes))

	}
}

func mainOld() {
	bmff.DebugLogger(bmff.STDLogger{})
	//f, err := os.Open("test.CR2")
	//f, err := os.Open("../../test/img/CanonR10_1.CR3")
	f, err := os.Open("../../test/img/1.heic")
	//f, err := os.Open("IMG_3001.jpeg")
	//f, err := os.Open("../testImages/Heic.exif")
	//f, err := os.Open("3.CR3")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	exif, err := imagemeta.DecodeHeif(f)
	fmt.Println(exif)
	f.Seek(0, 0)
	exif, err = imagemeta.DecodeTiff(f)
	fmt.Println(exif)

	//f.Seek(0, 0)
	//e, err := imagemeta.DecodeCR3(f)
	//fmt.Println(e)
}

func main() {
	exif2.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel)
	dir := "../../test/img/"
	//f, err := os.Open("../testImages/Heic.exif")
	f, err := os.Open(dir + "/" + "CanonR6_1.HIF")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	e, err := imagemeta.DecodeHeif(f)
	if err != nil {
		panic(err)
	}
	_ = e
	fmt.Println(e)
}
