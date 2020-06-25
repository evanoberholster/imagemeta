package imagetype

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
		{"../testImages/ARW.exif", ImageTiff},
		{"../testImages/NEF.exif", ImageTiff},
		{"../testImages/CR2.exif", ImageCR2},
		{"../testImages/Heic.exif", ImageHEIF},
		{"../testImages/CRW.CRW", ImageCRW},
		{"../testImages/XMP.xmp", ImageXMP},
	}
	for _, header := range exifHeaderTests {
		t.Run(header.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(header.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			// Search for Image Type
			imageType, err := Scan(f)
			if err != nil {
				t.Fatal(err)
			}

			if header.imageType != imageType {
				t.Errorf("Incorrect Imagetype wanted %s got %s", header.imageType.String(), imageType.String())
			}
		})
	}

}
