package exiftool

import (
	"encoding/binary"
	"os"
	"testing"
)

// Tests
func TestSearchExifHeader(t *testing.T) {
	exifHeaderTests := []struct {
		filename         string
		byteOrder        binary.ByteOrder
		firstIfdOffset   uint32
		tiffHeaderOffset uint32
	}{
		{"testImages/ARW.exif", binary.LittleEndian, 0x0008, 0x00},
		{"testImages/NEF.exif", binary.LittleEndian, 0x0008, 0x00},
		{"testImages/CR2.exif", binary.LittleEndian, 0x0010, 0x00},
		{"testImages/Heic.exif", binary.BigEndian, 0x0008, 0x1178},
	}
	for _, header := range exifHeaderTests {
		t.Run(header.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(header.filename)
			if err != nil {
				t.Fatal(err)
			}
			// Search for Tiff header
			eh, err := SearchExifHeader(f)
			if err != nil {
				t.Fatal(err)
			}
			if eh.byteOrder != header.byteOrder {
				t.Errorf("Incorrect Byte Order wanted %s got %s", header.byteOrder, eh.byteOrder)
			}
			if eh.firstIfdOffset != header.firstIfdOffset {
				t.Errorf("Incorrect first Ifd Offset wanted 0x%04x got 0x%04x ", header.firstIfdOffset, eh.firstIfdOffset)
			}
			if eh.tiffHeaderOffset != header.tiffHeaderOffset {
				t.Errorf("Incorrect tiff Header Offset wanted 0x%04x got 0x%04x ", header.tiffHeaderOffset, eh.tiffHeaderOffset)
			}
		})
	}
}
