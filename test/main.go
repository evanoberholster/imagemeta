package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/bmff"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: main <file>\n")
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	bmff.Debug = true
	hm := imagemeta.NewHeifMetadata(bufio.NewReader(f))
	hm.GetMeta()

}
