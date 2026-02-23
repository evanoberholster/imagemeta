package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/evanoberholster/imagemeta/imagetype"
)

const headerSize = 64

type fixture struct {
	name string
	want imagetype.FileType
	buf  []byte
}

func main() {
	out := flag.String("out", defaultOutputPath(), "output path for generated test.dat")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		panic(err)
	}

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	cases := fixtures()
	for i, tc := range cases {
		if len(tc.buf) != headerSize {
			panic(fmt.Errorf("fixture %q has len=%d, want=%d", tc.name, len(tc.buf), headerSize))
		}

		got, err := imagetype.Buf(tc.buf)
		if err != nil {
			panic(fmt.Errorf("fixture %q Buf() error: %w", tc.name, err))
		}
		if got != tc.want {
			panic(fmt.Errorf("fixture %q Buf()=%s want=%s", tc.name, got, tc.want))
		}

		if _, err := f.Write(tc.buf); err != nil {
			panic(fmt.Errorf("fixture %q write error: %w", tc.name, err))
		}

		fmt.Printf("%02d %-18s -> %s\n", i, tc.name, got)
	}

	fmt.Printf("wrote %d records (%d bytes) to %s\n", len(cases), len(cases)*headerSize, *out)
}

func defaultOutputPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("imagetype", "test.dat")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "test.dat"))
}

func fixtures() []fixture {
	return []fixture{
		{name: ".CRW", want: imagetype.ImageCRW, buf: crwHeader()},
		{name: ".CR2/GPS", want: imagetype.ImageCR2, buf: cr2Header()},
		{name: ".CR2/7D", want: imagetype.ImageCR2, buf: cr2Header()},
		{name: ".CR3", want: imagetype.ImageCR3, buf: ftypHeader("crx ", "isom")},
		{name: ".JPG/GPS", want: imagetype.ImageJPEG, buf: prefixHeader([]byte{0xFF, 0xD8})},
		{name: ".JPG/NoExif", want: imagetype.ImageJPEG, buf: prefixHeader([]byte{0xFF, 0xD8})},
		{name: ".JPG/GoPro", want: imagetype.ImageJPEG, buf: prefixHeader([]byte{0xFF, 0xD8})},
		{name: ".JPEG", want: imagetype.ImageJPEG, buf: prefixHeader([]byte{0xFF, 0xD8})},
		{name: ".HEIC/iPhone", want: imagetype.ImageHEIC, buf: ftypHeader("heic")},
		{name: ".HEIC/Conv", want: imagetype.ImageHEIC, buf: ftypHeader("heim")},
		{name: ".HEIC/Alt", want: imagetype.ImageHEIC, buf: ftypHeader("heis")},
		{name: ".WEBP", want: imagetype.ImageWebP, buf: webpHeader()},
		{name: ".GPR/GoPro", want: imagetype.ImageGPR, buf: tiffSubtypeHeader(0x0039, 0, 0x0100)},
		{name: ".NEF/Nikon", want: imagetype.ImageNEF, buf: tiffSubtypeHeader(0x001B, 1, 0x0100)},
		{name: ".ARW/Sony", want: imagetype.ImageARW, buf: tiffSubtypeHeader(0x0012, 1, 0x0103)},
		{name: ".DNG/Adobe", want: imagetype.ImageDNG, buf: tiffSubtypeHeader(0x003F, 1, 0x0100)},
		{name: ".PNG", want: imagetype.ImagePNG, buf: prefixHeader([]byte{0x89, 0x50, 0x4E, 0x47})},
		{name: ".RW2", want: imagetype.ImagePanaRAW, buf: rw2Header()},
		{name: ".XMP", want: imagetype.ImageXMP, buf: prefixHeader([]byte("<x:xmpmeta"))},
		{name: ".PSD", want: imagetype.ImagePSD, buf: prefixHeader([]byte("8BPS"))},
		{name: ".JP2/JPEG2000", want: imagetype.ImageJP2K, buf: prefixHeader([]byte{0x00, 0x00, 0x00, 0x0C, 0x6A, 0x50, 0x20, 0x20, 0x0D, 0x0A, 0x87, 0x0A})},
		{name: ".BMP", want: imagetype.ImageBMP, buf: prefixHeader([]byte("BM"))},
	}
}

func prefixHeader(prefix []byte) []byte {
	buf := make([]byte, headerSize)
	copy(buf, prefix)
	return buf
}

func crwHeader() []byte {
	buf := prefixHeader([]byte{0x49, 0x49})
	copy(buf[6:], []byte("HEAPCCDR"))
	return buf
}

func cr2Header() []byte {
	buf := prefixHeader([]byte{0x49, 0x49, 0x2A, 0x00})
	copy(buf[8:], []byte{0x43, 0x52, 0x02, 0x00})
	return buf
}

func rw2Header() []byte {
	buf := prefixHeader([]byte{0x49, 0x49, 0x55, 0x00})
	copy(buf[8:], []byte{0x88, 0xE7, 0x74, 0xD8})
	return buf
}

func webpHeader() []byte {
	buf := prefixHeader([]byte("RIFF"))
	copy(buf[8:], []byte("WEBP"))
	return buf
}

func ftypHeader(major string, compatible ...string) []byte {
	if len(major) != 4 {
		panic(fmt.Errorf("major brand %q must be 4 bytes", major))
	}
	buf := prefixHeader([]byte{0x00, 0x00, 0x00, 0x20, 'f', 't', 'y', 'p'})
	copy(buf[8:12], []byte(major))
	copy(buf[12:16], []byte("0001"))

	offset := 16
	for _, brand := range compatible {
		if len(brand) != 4 {
			panic(fmt.Errorf("compatible brand %q must be 4 bytes", brand))
		}
		if offset+4 > len(buf) {
			break
		}
		copy(buf[offset:offset+4], []byte(brand))
		offset += 4
	}
	return buf
}

func tiffSubtypeHeader(entryCount uint16, firstValue uint32, secondTag uint16) []byte {
	buf := prefixHeader([]byte{0x49, 0x49, 0x2A, 0x00})

	binary.LittleEndian.PutUint16(buf[8:10], entryCount)
	binary.LittleEndian.PutUint16(buf[10:12], 0x00FE)
	binary.LittleEndian.PutUint16(buf[12:14], 4)
	binary.LittleEndian.PutUint32(buf[14:18], 1)
	binary.LittleEndian.PutUint32(buf[18:22], firstValue)
	binary.LittleEndian.PutUint16(buf[22:24], secondTag)

	return buf
}
