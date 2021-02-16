package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta/xmp"
)

const (
	dir   = "../../../test/img/"
	dir2  = "test/samples/"
	name  = "CanonEOS7DII.xmp"
	name1 = "jpeg.xmp"
	name2 = "9.jpg"
)

func main() {
	xmp.DebugMode = true
	f, err := os.Open(dir + name2) //"retouch.xmp") //name)
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println(f.Close())
	}()

	xmp, err := xmp.Read(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xmp)
}
