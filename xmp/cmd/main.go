package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/image-meta/xml"
)

const (
	dir  = "../../test/img/"
	dir2 = "test/samples/"
	name = "CanonEOS7D.xmp"
)

func main() {
	xml.DebugMode = true
	f, err := os.Open(dir2 + name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	xmp, err := xml.Read(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xmp)
}

func main2() {
	fmt.Println([]byte(":"))
}
