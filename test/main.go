package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/evanoberholster/imagemeta/bmff"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/heic"
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
	hm := heic.NewMetadata(f)
	err = hm.GetMeta()
	if err != nil {
		fmt.Println(err)
		// Error retrieving Heic Metadata
	}
	hm.ExifDecodeFn = func(r io.Reader, header exif.Header) error {
		return nil
	}
	edata, err := hm.DecodeExif(f)
	fmt.Println(edata, err)
}
