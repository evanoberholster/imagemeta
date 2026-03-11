package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/exif"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/rs/zerolog"
)

func init() {
	imagemeta.SetLogger(zerolog.ConsoleWriter{Out: os.Stdout}, zerolog.DebugLevel)
	exif.Logger = exif.Logger.Level(zerolog.WarnLevel)
	isobmff.Logger = isobmff.Logger.Level(zerolog.DebugLevel)
}

func main() {
	//f, err := os.Open(dir + "/" + "DJI.dng")
	//f, err := os.Open(dir + "/" + "iPhone11.heic")
	f, err := os.Open("1.cr3")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	e, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}
