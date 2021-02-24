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
	hm, err := heic.NewMetadata(f)
	if err != nil {
		fmt.Println(err)
		// Error retrieving Heic Metadata
	}
	var e *exif.Data
	hm.ExifDecodeFn = func(r io.Reader, header exif.Header) error {
		e, err = exif.ParseExif(f, header)
		fmt.Println(e, err, header)
		return nil
	}
	err = hm.DecodeExif(f)
	fmt.Println(e, err)

	fmt.Println(hm.FileType)
}
