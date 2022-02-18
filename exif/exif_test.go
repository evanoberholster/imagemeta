package exif

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/tag"
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
}{
	// TODO: Add test for RW2
	{"../testImages/JPEG.jpg", imagetype.ImageJPEG, "GoPro", "HERO4 Silver", 113, 2.8, 3, meta.NewShutterSpeed(1, 60), 0, 0, time.Unix(1476205190, 0)},
	{"../testImages/Hero8.GPR", imagetype.ImageTiff, "GoPro", "HERO8 Black", 317, 2.8, 3, meta.NewShutterSpeed(1, 240), 4000, 3000, time.Unix(1590641247, 0)},
	{"../testImages/ARW.exif", imagetype.ImageARW, "SONY", "SLT-A55V", 100, 13.0, 30.0, meta.NewShutterSpeed(1, 100), 4928, 3280, time.Unix(1508673260, 0)},
	{"../testImages/NEF.exif", imagetype.ImageNEF, "NIKON CORPORATION", "NIKON D7100", 100, 8.0, 50.0, meta.NewShutterSpeed(10, 300), 160, 120, time.Unix(1378201522, 0)},
	{"../testImages/CR2.exif", imagetype.ImageCR2, "Canon", "Canon EOS-1Ds Mark III", 100, 1.20, 50.0, meta.NewShutterSpeed(1, 40), 5616, 3744, time.Unix(1192715074, 0)},
	{"../testImages/Heic.exif", imagetype.ImageHEIF, "Canon", "Canon EOS 6D", 500, 5.0, 20.0, meta.NewShutterSpeed(1, 20), 3648, 5472, time.Unix(1575608513, 0)},
}

//func TestGenSamples(t *testing.T) {
//	for _, wantedExif := range exifTests {
//		f, err := os.Open(wantedExif.filename)
//		if err != nil {
//			panic(err)
//		}
//		defer func() {
//			err = f.Close()
//			if err != nil {
//				panic(err)
//			}
//		}()
//		buf, err := ioutil.ReadAll(f)
//		e, err := ScanExif(bytes.NewReader(buf))
//		if !assert.ErrorIs(t, err, nil) {
//			return
//		}
//
//		buf, err = json.Marshal(e)
//		if !assert.ErrorIs(t, err, nil) {
//			return
//		}
//		// Pretty JSON
//		buf = pretty.Pretty(buf)
//
//		dat, err := os.Create(wantedExif.filename + ".json")
//		if !assert.ErrorIs(t, err, nil) {
//			return
//		}
//		defer func() {
//			err = dat.Close()
//			if err != nil {
//				panic(err)
//			}
//		}()
//		if _, err := dat.Write(buf); err != nil {
//			err = f.Close()
//			panic(err)
//		}
//	}
//}
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
			e, err := ScanExif(cb)
			if !assert.ErrorIs(t, err, nil) {
				return
			}
			b1, err := e.MarshalJSON()
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

func TestParseExif(t *testing.T) {
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
			e, err := ScanExif(cb)
			if !assert.ErrorIs(t, err, nil) {
				return
			}

			var val interface{}

			// Camera Make
			assert.Equal(t, wantedExif.make, e.CameraMake(), "Camera Make")

			// Camera Model
			assert.Equal(t, wantedExif.model, e.CameraModel(), "Camera Model")

			// Dimensions
			w, h := e.Dimensions().Size()
			assert.Equal(t, wantedExif.width, w, "Image Width")
			assert.Equal(t, wantedExif.height, h, "Image Height")

			// ISO Speed
			val, _ = e.ISOSpeed()
			assert.Equal(t, wantedExif.ISOSpeed, val, "ISO Speed")

			// Aperture
			val, _ = e.Aperture()
			assert.Equal(t, wantedExif.aperture, val, "Aperture")

			// Shutter Speed
			val, _ = e.ShutterSpeed()
			assert.Equal(t, wantedExif.shutterSpeed, val, "Shutter Speed")

			// Focal Length
			val, _ = e.FocalLength()
			assert.Equal(t, wantedExif.focalLength, val, "Focal Length")

			// Created Date
			date, _ := e.DateTime()
			assert.Equal(t, wantedExif.createdDate.Unix(), date.Unix())

		})
	}
}

// jsonExif for testing purposes.

type jsonExif struct {
	Ifds   map[string]map[uint8]jsonIfds `json:"Ifds"`
	It     imagetype.ImageType           `json:"ImageType"`
	Make   string
	Model  string
	Width  uint16
	Height uint16
}

type jsonIfds struct {
	Tags []jsonTags `json:"Tags"`
}

func (je *jsonExif) addTag(ifd ifds.IFD, ifdIndex uint8, t tag.Tag, v interface{}) {
	if je.Ifds == nil {
		je.Ifds = make(map[string]map[uint8]jsonIfds)
	}
	ji, ok := je.Ifds[ifd.String()]
	if !ok {
		je.Ifds[ifd.String()] = make(map[uint8]jsonIfds)
		ji = je.Ifds[ifd.String()]
	}
	jm, ok := ji[ifdIndex]
	if !ok {
		ji[ifdIndex] = jsonIfds{make([]jsonTags, 0)}
		jm = ji[ifdIndex]
	}
	jm.insertSorted(jsonTags{Name: ifd.TagName(t.ID), Type: t.Type(), ID: t.ID, Count: t.UnitCount, Value: v})
	je.Ifds[ifd.String()][ifdIndex] = jm
}

func (ji *jsonIfds) insertSorted(e jsonTags) {
	i := sort.Search(len(ji.Tags), func(i int) bool { return ji.Tags[i].ID > e.ID })
	ji.Tags = append(ji.Tags, jsonTags{})
	copy(ji.Tags[i+1:], ji.Tags[i:])
	ji.Tags[i] = e
}

type jsonTags struct {
	ID    tag.ID
	Name  string
	Count uint16
	Type  tag.Type
	Value interface{}
}

func (jt jsonTags) MarshalJSON() ([]byte, error) {
	st := struct {
		ID    string
		Name  string
		Count uint16
		Type  string
		Value interface{} `json:"Val"`
	}{
		ID:    jt.ID.String(),
		Name:  jt.Name,
		Count: jt.Count,
		Type:  jt.Type.String(),
		Value: jt.Value,
	}
	return json.Marshal(st)
}

// MarshalJSON implements the JSONMarshaler interface that is used by encoding/json
// This is mostly used for testing and debuging.
func (e *Data) MarshalJSON() ([]byte, error) {
	je := jsonExif{It: e.imageType, Make: e.CameraMake(), Model: e.CameraModel(), Width: e.width, Height: e.height}
	for k, t := range e.tagMap {
		ifd, ifdIndex, _ := k.Val()
		value := e.GetTagValue(t)
		je.addTag(ifd, ifdIndex, t, value)
	}

	return json.Marshal(je)
}
