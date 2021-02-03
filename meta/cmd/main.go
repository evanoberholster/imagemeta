package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
)

const testFilename = "../../../test/img/10.jpg"

func main() {
	f, err := os.Open(testFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println(testFilename)
	//buf, _ := ioutil.ReadAll(f)
	//cb := bytes.NewReader(buf)
	start := time.Now()
	m, err := meta.Scan(f, imagetype.ImageJPEG)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)

	x, err := xmp.Read(bytes.NewReader([]byte(m.XML())))
	fmt.Println(m.XML())
	fmt.Println(m.Size())
	fmt.Println(m.Header())
	fmt.Println(elapsed)
	fmt.Println(x, err)

}
