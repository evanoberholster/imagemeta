package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/xmp"
)

const testFilename = "../../test/img/10.jpg"

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
	br := bufio.NewReader(f)

	var x xmp.XMP
	xmpDecodeFn := func(r io.Reader) error {
		var err error
		x, err = xmp.Read(r)
		return err
	}
	start := time.Now()
	m, err := imagemeta.ScanBuf2(br, imagetype.ImageJPEG, xmpDecodeFn)
	//m, err := meta.Scan(f, imagetype.ImageJPEG)
	if err != nil {
		panic(err)
	}

	//x, err = xmp.Read(bytes.NewReader([]byte(m.XML())))
	elapsed := time.Since(start)
	fmt.Println(m.XMP())
	fmt.Println(m.Size())
	fmt.Println(m.Header())
	fmt.Println(elapsed)
	fmt.Println(x, err)

}
