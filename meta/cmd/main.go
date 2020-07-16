package main

import (
	"fmt"
	"os"
	"time"

	"github.com/evanoberholster/exiftool/imagetype"
	"github.com/evanoberholster/exiftool/meta"
)

const testFilename = "../../../../test/img/6.jpg"

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
	m, err := meta.Scan(f, imagetype.ImageUnknown)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)

	fmt.Println(m.XML())
	fmt.Println(m.Size())
	fmt.Println(m.Header())
	fmt.Println(elapsed)

}
