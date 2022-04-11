package exif

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/pretty"
)

var exifTests = []struct {
	filename     string
	imageType    imagetype.ImageType
	make         string
	model        string
	ISOSpeed     uint32
	aperture     meta.Aperture
	focalLength  meta.FocalLength
	shutterSpeed meta.ShutterSpeed
	width        uint32
	height       uint32
	createdDate  time.Time
	orientation  meta.Orientation
	header       meta.ExifHeader
}{
	// TODO: Add test for RW2
	{"../testImages/JPEG.jpg", imagetype.ImageJPEG, "GoPro", "HERO4 Silver", 113, 2.8, 3, meta.NewShutterSpeed(1, 60), 0, 0, time.Unix(1476205190, 0), meta.OrientationHorizontal, meta.NewExifHeader(binary.LittleEndian, 13746, 12, 0, imagetype.ImageJPEG)},
	{"../testImages/Hero8.GPR", imagetype.ImageTiff, "GoPro", "HERO8 Black", 317, 2.8, 3, meta.NewShutterSpeed(1, 240), 4000, 3000, time.Unix(1590641247, 0), meta.OrientationMirrorHorizontal, meta.NewExifHeader(binary.BigEndian, 8, 8, 0, imagetype.ImageGPR)},
	{"../testImages/ARW.exif", imagetype.ImageARW, "SONY", "SLT-A55V", 100, 13.0, 30.0, meta.NewShutterSpeed(1, 100), 4928, 3280, time.Unix(1508673260, 0), meta.OrientationHorizontal, meta.NewExifHeader(binary.LittleEndian, 8, 0, 0, imagetype.ImageARW)},
	{"../testImages/NEF.exif", imagetype.ImageNEF, "NIKON CORPORATION", "NIKON D7100", 100, 8.0, 50.0, meta.NewShutterSpeed(10, 300), 160, 120, time.Unix(1378201522, 0), meta.OrientationHorizontal, meta.NewExifHeader(binary.LittleEndian, 8, 0, 0, imagetype.ImageNEF)},
	{"../testImages/CR2.exif", imagetype.ImageCR2, "Canon", "Canon EOS-1Ds Mark III", 100, 1.20, 50.0, meta.NewShutterSpeed(1, 40), 5616, 3744, time.Unix(1192715074, 0), meta.OrientationHorizontal, meta.NewExifHeader(binary.LittleEndian, 16, 0, 0, imagetype.ImageCR2)},
	{"../testImages/Heic.exif", imagetype.ImageHEIF, "Canon", "Canon EOS 6D", 500, 5.0, 20.0, meta.NewShutterSpeed(1, 20), 3648, 5472, time.Unix(1575608513, 0), meta.OrientationHorizontal, meta.NewExifHeader(binary.BigEndian, 8, 4472, 0, imagetype.ImageHEIF)},
}

func TestGenSamples(t *testing.T) {
	for _, wantedExif := range exifTests {
		f, err := os.Open(wantedExif.filename)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}()
		buf, err := ioutil.ReadAll(f)
		e, err := ParseExif(bytes.NewReader(buf), wantedExif.header)
		if err != nil {
			t.Error(wantedExif.filename)
			panic(err)
		}

		buf, err = e.DebugJSON()
		if err != nil {
			panic(err)
		}
		// Pretty JSON
		buf = pretty.Pretty(buf)

		dat, err := os.Create(wantedExif.filename + ".json")
		if !assert.ErrorIs(t, err, nil) {
			return
		}
		defer func() {
			err = dat.Close()
			if err != nil {
				panic(err)
			}
		}()
		if _, err := dat.Write(buf); err != nil {
			err = f.Close()
			panic(err)
		}
	}
}

//
func TestPreviouslyParsedExif(t *testing.T) {
	for _, wantedExif := range exifTests {
		t.Run(wantedExif.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(wantedExif.filename)
			if err != nil {
				t.Fatal(err)
			}
			buf, _ := ioutil.ReadAll(f)
			if err := f.Close(); err != nil {
				panic(err)
			}
			cb := bytes.NewReader(buf)
			e, err := ParseExif(cb, wantedExif.header)
			if !assert.ErrorIs(t, err, nil) {
				return
			}
			b1, err := e.DebugJSON()
			if !assert.ErrorIs(t, err, nil) {
				return
			}

			// Open file
			f, err = os.Open(wantedExif.filename + ".json")
			if err != nil {
				t.Fatal(err)
			}
			b2, _ := ioutil.ReadAll(f)
			if err := f.Close(); err != nil {
				panic(err)
			}

			b2 = pretty.Ugly(b2)
			if !bytes.Equal(b1, b2) {
				t.Errorf("Please review: Incorrect Exif Data wanted length %d got length %d", len(b2), len(b1))
			}
		})
	}
}

//
//func TestParseExif(t *testing.T) {
//	for _, wantedExif := range exifTests {
//		t.Run(wantedExif.filename, func(t *testing.T) {
//
//			// Open file
//			f, err := os.Open(wantedExif.filename)
//			if err != nil {
//				t.Fatal(err)
//			}
//			buf, _ := ioutil.ReadAll(f)
//			if err := f.Close(); err != nil {
//				panic(err)
//			}
//			cb := bytes.NewReader(buf)
//			e, err := ScanExif2(cb)
//			if !assert.ErrorIs(t, err, nil) {
//				return
//			}
//
//			var val interface{}
//
//			// Camera Make
//			assert.Equal(t, wantedExif.make, e.CameraMake(), "Camera Make")
//
//			// Camera Model
//			assert.Equal(t, wantedExif.model, e.CameraModel(), "Camera Model")
//
//			// Dimensions
//			w, h := e.Dimensions().Size()
//			assert.Equal(t, wantedExif.width, w, "Image Width")
//			assert.Equal(t, wantedExif.height, h, "Image Height")
//
//			// ISO Speed
//			val, _ = e.ISOSpeed()
//			assert.Equal(t, wantedExif.ISOSpeed, val, "ISO Speed")
//
//			// Aperture
//			val, _ = e.Aperture()
//			assert.Equal(t, wantedExif.aperture, val, "Aperture")
//
//			// Shutter Speed
//			val, _ = e.ShutterSpeed()
//			assert.Equal(t, wantedExif.shutterSpeed, val, "Shutter Speed")
//
//			// Focal Length
//			val, _ = e.FocalLength()
//			assert.Equal(t, wantedExif.focalLength, val, "Focal Length")
//
//			// Created Date
//			date, _ := e.DateTime(nil)
//			assert.Equal(t, wantedExif.createdDate.Unix(), date.Unix())
//
//			// Orientation
//			orientation, _ := e.Orientation()
//			assert.Equal(t, wantedExif.orientation, orientation)
//		})
//	}
//}
//
// jsonExif for testing purposes.
