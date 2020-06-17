package exiftool

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

// TODO: write tests for ParseExif

func TestParseExif(t *testing.T) {
	exifTests := []struct {
		filename string
		make     string
		model    string
	}{
		{"testImages/ARW.exif", "SONY", "SLT-A55V"},
		{"testImages/NEF.exif", "NIKON CORPORATION", "NIKON D7100"},
		{"testImages/CR2.exif", "Canon", "Canon EOS-1Ds Mark III"},
		{"testImages/Heic.exif", "Canon", "Canon EOS 6D"},
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
		})
	}
}
