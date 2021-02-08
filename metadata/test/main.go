package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/evanoberholster/imagemeta/metadata"
	"github.com/evanoberholster/imagemeta/metadata/bmff"
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
	defer f.Close()
	bmff.Debug = true
	hm := metadata.NewHeifMetadata(bufio.NewReader(f))
	hm.GetMeta()

}
