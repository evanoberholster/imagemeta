package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/imagemeta/meta/xmp"
)

const (
	dir       = "../../../test/img/"
	dir2      = "../test/samples/" //samples/"
	filename1 = "../test/samples/bluesquare.avi.xmp"
	filename2 = "../test/1.xmp"
	filename3 = "../test/samples/RAW_SONY_SLTA55V.xmp"
	filename4 = "../test/samples/CanonEOS7D.xmp"
	filename5 = "../test/jpeg.xmp"
	name      = "CanonEOS7DII.xmp"
	name1     = "jpeg.xmp"
	name2     = "9.jpg"
	mholtTest = "../../cmd/img1.jpg"
)

func main() {
	xmp.DebugMode = true
	f, err := os.Open(filename2)
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println(f.Close())
	}()

	xmp, err := xmp.ParseXmp(f)
	if err != nil {
		fmt.Println(err)
	}
	// /xmp, err := xmp.ParseXmp(f)
	// /if err != nil {
	// fmt.Println(err)
	// /}
	fmt.Println(xmp)
}
