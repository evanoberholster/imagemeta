package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/metadata"
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
	br := bufio.NewReader(f)

	//xmpDecodeFn := func(r io.Reader) error {
	//	ioutil.ReadAll(r)
	//	return nil
	//}
	start := time.Now()
	m, err := metadata.ScanBuf2(br, imagetype.ImageJPEG, nil)
	//m, err := meta.Scan(f, imagetype.ImageJPEG)
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
