package exiftool

import (
	"os"
	"testing"
)

// TODO: write tests for ImageTypes

// Tests
func TestSearchImageType(t *testing.T) {
	exifHeaderTests := []struct {
		filename  string
		imageType ImageType
	}{
		{"testImages/ARW.exif", ImageTiff},
		{"testImages/NEF.exif", ImageTiff},
		{"testImages/CR2.exif", ImageCR2},
		{"testImages/Heic.exif", ImageHEIF},
	}
	for _, header := range exifHeaderTests {
		t.Run(header.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(header.filename)
			if err != nil {
				t.Fatal(err)
			}
			// Search for Image Type
			imageType, err := SearchImageType(f)
			if err != nil {
				t.Fatal(err)
			}

			if header.imageType != imageType {
				t.Errorf("Incorrect Byte Order wanted %s got %s", header.imageType.String(), imageType.String())
			}
		})
	}

}
