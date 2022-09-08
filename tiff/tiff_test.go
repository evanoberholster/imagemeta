package tiff

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

// Tests
func TestTiff(t *testing.T) {
	exifHeaderTests := []struct {
		filename         string
		byteOrder        meta.ByteOrder
		firstIfdOffset   uint32
		tiffHeaderOffset uint32
		imageType        imagetype.ImageType
	}{
		{"../testImages/ARW.exif", meta.LittleEndian, 0x0008, 0x00, imagetype.ImageTiff},
		{"../testImages/NEF.exif", meta.LittleEndian, 0x0008, 0x00, imagetype.ImageTiff},
		{"../testImages/CR2.exif", meta.LittleEndian, 0x0010, 0x00, imagetype.ImageCR2},
		{"../testImages/Heic.exif", meta.BigEndian, 0x0008, 0x1178, imagetype.ImageHEIF},
	}
	for _, header := range exifHeaderTests {
		t.Run(header.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(header.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			// Search for Tiff header
			br := bufio.NewReader(f)
			h, err := ScanTiffHeader(br, header.imageType)
			if err != nil {
				t.Fatal(err)
			}
			if h.ByteOrder != header.byteOrder {
				t.Errorf("Incorrect Byte Order wanted %s got %s", header.byteOrder, h.ByteOrder)
			}
			if h.FirstIfdOffset != header.firstIfdOffset {
				t.Errorf("Incorrect first Ifd Offset wanted 0x%04x got 0x%04x ", header.firstIfdOffset, h.FirstIfdOffset)
			}
			if h.TiffHeaderOffset != header.tiffHeaderOffset {
				t.Errorf("Incorrect tiff Header Offset wanted 0x%04x got 0x%04x ", header.tiffHeaderOffset, h.TiffHeaderOffset)
			}
			if h.ImageType != header.imageType {
				t.Errorf("Incorrect tiff Imagetype wanted: %s got: %s ", header.imageType, h.ImageType)
			}
			if !h.IsValid() {
				t.Errorf("Wanted valid tiff Header")
			}
		})
	}

	// Error No Tiff Header
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_, err := ScanTiffHeader(bytes.NewReader(buf), imagetype.ImageTiff)
	if err != meta.ErrNoExif {
		t.Errorf("Incorrect err wanted %s got %s ", meta.ErrNoExif, err)
	}
}
