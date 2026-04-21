package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
)

func init() {
	logging.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel)
}

func main() {
	path := "2.CR3"
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
