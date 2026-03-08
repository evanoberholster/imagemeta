package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta/meta/xmp"
)

const (
	filename = "../test/jpeg.xmp"
)

func main() {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println(f.Close())
	}()

	xmp, err := xmp.ParseXmpWithOptions(f, xmp.ParseOptions{Debug: true})
	if err != nil {
		fmt.Println(err)
	}
	// /xmp, err := xmp.ParseXmp(f)
	// /if err != nil {
	// fmt.Println(err)
	// /}
	fmt.Println(xmp)
}
