package exiftool

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/evanoberholster/exiftool/exif"
	"github.com/rs/zerolog"
)

// TODO: write tests for ParseExif

func TestParseExif(t *testing.T) {
	exifTests := []struct {
		filename    string
		make        string
		model       string
		ISOSpeed    int
		aperture    float32
		focalLength exif.FocalLength
	}{
		{"testImages/ARW.exif", "SONY", "SLT-A55V", 100, 13.0, 30.0},
		{"testImages/NEF.exif", "NIKON CORPORATION", "NIKON D7100", 100, 8.0, 50.0},
		{"testImages/CR2.exif", "Canon", "Canon EOS-1Ds Mark III", 100, 1.20, 50.0},
		{"testImages/Heic.exif", "Canon", "Canon EOS 6D", 500, 5.0, 20.0},
	}
	for _, wantedExif := range exifTests {
		t.Run(wantedExif.filename, func(t *testing.T) {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			// Open file
			f, err := os.Open(wantedExif.filename)
			if err != nil {
				t.Fatal(err)
			}
			// Search for Tiff header
			eh, err := SearchExifHeader(f)
			if err != nil {
				panic(err)
			}
			f.Seek(0, 0)

			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)
			e, err := eh.ParseExif(cb)
			if err != nil {
				fmt.Println(err)
			}
			if e.Make != wantedExif.make {
				t.Errorf("Incorrect Exif Make wanted %s got %s", wantedExif.make, e.Make)
			}
			if e.Model != wantedExif.model {
				t.Errorf("Incorrect Exif Model wanted %s got %s", wantedExif.model, e.Model)
			}
			isoSpeed, err := e.ISOSpeed()
			if err != nil || isoSpeed != wantedExif.ISOSpeed {
				t.Errorf("Incorrect ISO Speed wanted %d got %d", wantedExif.ISOSpeed, isoSpeed)
			}
			aperture, err := e.Aperture()
			if err != nil || aperture != wantedExif.aperture {
				t.Errorf("Incorrect Aperture wanted %0.2f got %0.2f", wantedExif.aperture, aperture)
			}
			focalLength, err := e.FocalLength()
			if err != nil || focalLength != wantedExif.focalLength {
				t.Errorf("Incorrect Focal Length wanted %s got %s", wantedExif.focalLength.String(), focalLength.String())
			}
		})
	}
}
