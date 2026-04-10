package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/exif"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/rs/zerolog"
)

func init() {
	imagemeta.SetLogger(zerolog.ConsoleWriter{Out: os.Stdout}, zerolog.DebugLevel)
	exif.Logger = exif.Logger.Level(zerolog.ErrorLevel)
	isobmff.Logger = isobmff.Logger.Level(zerolog.ErrorLevel)
}

func main() {
	path := "1.NEF"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
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
	buf, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	//colored := pretty.Color(pretty.Pretty(buf), nil)
	fmt.Printf("%s\n", string(buf))
}
