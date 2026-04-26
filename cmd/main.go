package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
	"github.com/tidwall/pretty"
)

func init() {
	logging.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel)
}

func main() {
	path := "1.cr3"
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

	colored := pretty.Pretty(buf)
	fmt.Println(string(colored))
	//fmt.Printf("%s\n", string(buf))
}
