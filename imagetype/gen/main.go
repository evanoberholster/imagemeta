package main

import (
	"fmt"
	"os"

	"github.com/evanoberholster/exiftool/imagetype"
)

var (
	dir        = "../../../test/img/"
	benchmarks = []struct {
		name      string
		fileName  string
		imageType string
	}{
		{".CRW", "0.CRW", "image/x-canon-crw"},
		{".CR2/GPS", "2.CR2", "image/x-canon-cr2"},
		{".CR2/7D", "7D2.CR2", "image/x-canon-cr2"},
		{".CR3", "1.CR3", "image/x-canon-cr3"},
		{".JPG/GPS", "17.jpg", "image/jpeg"},
		{".JPG/NoExif", "20.jpg", "image/jpeg"},
		{".JPG/GoPro", "hero6.jpg", "image/jpeg"},
		{".JPEG", "21.jpeg", "image/jpeg"},
		{".HEIC/iPhone", "1.heic", "image/heif"},
		{".HEIC/Conv", "3.heic", "image/heif"},
		{".HEIC/Alt", "4.heic", "image/heif"},
		{".WEBP", "4.webp", "image/webp"},
		{".GPR/GoPro", "hero6.gpr", "image/tiff"},
		{".NEF/Nikon", "2.NEF", "image/tiff"},
		{".ARW/Sony", "2.ARW", "image/tiff"},
		{".DNG/Adobe", "1.DNG", "image/tiff"},
		{".PNG", "0.png", "image/png"},
		{".RW2", "4.RW2", "image/x-panasonic-raw"},
		{".XMP", "test.xmp", "application/rdf+xml"},
	}
)

func main() {
	dat, err := os.Create("test.dat")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = dat.Close()
		if err != nil {
			panic(err)
		}
	}()
	buf := make([]byte, 32)
	// first 32 bytes of each file
	for _, filename := range benchmarks {
		f, err := os.Open(dir + filename.fileName)
		if err != nil {
			panic(err)
		}
		if n, err := f.ReadAt(buf, 0); n != 32 || err != nil {
			err = f.Close()
			panic(err)
		}
		fmt.Println(imagetype.Scan(f))
		if _, err := dat.Write(buf); err != nil {
			err = f.Close()
			panic(err)
		}
		if _, err := f.Close(); err != nil {
			panic(err)
		}
	}
}
