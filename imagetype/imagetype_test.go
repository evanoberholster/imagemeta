package imagetype

import (
	"bytes"
	"os"
	"testing"
)

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

func TestImageType(t *testing.T) {
	fileOffset := 32
	testDataFilename := "test.dat"

	var headerTests = []struct {
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

	// Open file
	f, err := os.Open(testDataFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, fileOffset)

	for i, header := range headerTests {
		t.Run(header.name, func(t *testing.T) {
			if n, err := f.ReadAt(buf, int64(i*fileOffset)); n != fileOffset || err != nil {
				t.Fatal(err)
			}
			// Search for Image Type
			imageType, err := Scan(bytes.NewReader(buf))
			if err != nil {
				t.Fatal(err)
			}

			if header.imageType != imageType.String() {
				t.Errorf("Incorrect Imagetype wanted %s got %s", header.imageType, imageType.String())
			}
		})
	}

}
