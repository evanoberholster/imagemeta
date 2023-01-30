// Copyright (c) 2018-2023 Evan Oberholster. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package jpeg

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var (
	dir            = "../assets/"
	benchmarksJPEG = []struct {
		fileName  string
		noExifErr bool
	}{
		{"a1.jpg", false},
		{"a2.jpg", true},
		{"JPEG.jpg", false},
		{"NoExif.jpg", true},
	}
)

func BenchmarkScanJPEG100(b *testing.B) {
	for _, bm := range benchmarksJPEG {
		f, err := os.Open(dir + bm.fileName)
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()
		buf, _ := io.ReadAll(f)
		r := bytes.NewReader(buf)
		b.ReportAllocs()
		b.ResetTimer()

		b.Run(bm.fileName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r.Seek(0, 0)
				if err := ScanJPEG(r, nil, nil); err != nil {
					if !bm.noExifErr {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

func TestScanJPEG(t *testing.T) {
	testJPEGs := []struct {
		filename string
		exif     bool
		header   meta.ExifHeader
		width    uint32
		height   uint32
	}{
		{"../assets/JPEG.jpg", true, meta.NewExifHeader(utils.LittleEndian, 13746, 12, 13872, imagetype.ImageJPEG), 1000, 563},
		{"../assets/NoExif.jpg", true, meta.NewExifHeader(utils.BigEndian, 8, 30, 140, imagetype.ImageJPEG), 50, 50},
		{"../assets/a2.jpg", false, meta.NewExifHeader(utils.LittleEndian, 13746, 12, 13872, imagetype.ImageJPEG), 1024, 1280},
		{"../assets/a1.jpg", true, meta.NewExifHeader(utils.BigEndian, 8, 30, 752, imagetype.ImageJPEG), 389, 259},
	}

	for _, jpg := range testJPEGs {
		t.Run(jpg.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(jpg.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			testExifHeaderfn := func(r io.Reader, eh meta.ExifHeader) error {
				metaExifHeaderEqual(t, jpg.header, eh)
				return nil
			}
			testXmpHeaderFn := func(r io.Reader) error {
				return nil
			}

			err = ScanJPEG(f, testExifHeaderfn, testXmpHeaderFn)
			if jpg.exif && err != nil {
				t.Fatal(err)
			}
		})
	}

}

func metaExifHeaderEqual(t *testing.T, h1 meta.ExifHeader, h2 meta.ExifHeader) {
	if h1.ByteOrder != h2.ByteOrder {
		t.Errorf("Incorrect Byte Order wanted %s got %s", h1.ByteOrder, h2.ByteOrder)
	}
	if h1.FirstIfdOffset != h2.FirstIfdOffset {
		t.Errorf("Incorrect first Ifd Offset wanted %d got %d ", h1.FirstIfdOffset, h2.FirstIfdOffset)
	}
	if h1.TiffHeaderOffset != h2.TiffHeaderOffset {
		t.Errorf("Incorrect tiff Header Offset wanted %d got %d ", h1.TiffHeaderOffset, h2.TiffHeaderOffset)
	}
	if h1.ExifLength != h2.ExifLength {
		t.Errorf("Incorrect Exif Length wanted %d got %d ", h1.ExifLength, h2.ExifLength)
	}
	if h1.ImageType != h2.ImageType {
		t.Errorf("Incorrect Exif Header Imagetype wanted %s got %s ", h1.ImageType, h2.ImageType)
	}
	if !h2.IsValid() {
		t.Errorf("Wanted valid tiff Header")
	}

}

//func TestScanMarkers(t *testing.T) {
//	data := []byte{0, markerFirstByte, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	r := bytes.NewReader(data)
//	m := jpegReader{br: bufio.NewReader(r)}
//
//	// Test discard
//	m.discard(0)
//	if m.discarded != 0 {
//		t.Errorf("Incorrect Metadata.discard wanted %d got %d", 0, m.discarded)
//	}
//	// Test Scan Markers
//	buf, _ := m.br.Peek(16)
//	err := m.scanMarkers(buf)
//	if err != nil {
//		t.Errorf("Incorrect Scan Markers error wanted %s got %s", err.Error(), ErrNoJPEGMarker)
//	}
//
//	data = []byte{markerFirstByte, markerSOI, markerFirstByte, markerEOI, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	r = bytes.NewReader(data)
//	m = jpegReader{br: bufio.NewReader(r)}
//
//	// Test SOI
//	buf, _ = m.br.Peek(16)
//	err = m.scanMarkers(buf)
//	if m.discarded != 2 || m.pos != 1 || err != nil {
//		t.Errorf("Incorrect JPEG Start of Image error wanted discarded %d got %d", 2, m.discarded)
//	}
//
//	// Test EOI
//	buf, _ = m.br.Peek(16)
//	err = m.scanMarkers(buf)
//	if m.discarded != 4 || m.pos != 0 || err != nil {
//		t.Errorf("Incorrect JPEG End of Image error wanted discarded %d got %d", 4, m.discarded)
//	}
//
//	// Test Scan JPEG
//	err = ScanJPEG(bytes.NewReader(data), nil, nil)
//	if err != ErrNoJPEGMarker {
//		t.Errorf("Incorrect JPEG error at discarded %d wanted %s got %s", m.discarded, ErrNoJPEGMarker, err.Error())
//	}
//}
