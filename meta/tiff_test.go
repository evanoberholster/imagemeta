package meta

import (
	"bufio"
	"encoding/binary"
	"os"
	"testing"
)

// Tests
func TestScan(t *testing.T) {
	exifHeaderTests := []struct {
		filename         string
		byteOrder        binary.ByteOrder
		firstIfdOffset   uint32
		tiffHeaderOffset uint32
	}{
		{"../testImages/ARW.exif", binary.LittleEndian, 0x0008, 0x00},
		{"../testImages/NEF.exif", binary.LittleEndian, 0x0008, 0x00},
		{"../testImages/CR2.exif", binary.LittleEndian, 0x0010, 0x00},
		{"../testImages/Heic.exif", binary.BigEndian, 0x0008, 0x1178},
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
			h, err := ScanTiff(br)
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
		})
	}
}
