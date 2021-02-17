// +build !windows

package xmp

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Gen(t *testing.T) {
	for _, v := range testXmp {
		t.Run(v.filename, func(t *testing.T) {
			f, err := os.Open("test" + string(os.PathSeparator) + v.filename)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = f.Close()
			}()

			x, err := ParseXmp(f)
			if err != nil {
				if err != io.EOF {
					t.Error(err)
				}
			}

			j, err := json.Marshal(x)

			dat, err := os.Create("test" + string(os.PathSeparator) + v.filename + ".json")
			if err != nil {
				panic(err)
			}
			defer func() {
				err = dat.Close()
				if err != nil {
					panic(err)
				}
			}()
			_, err = dat.Write(j)
			if err != nil {
				t.Fatal(err)
			}

		})
	}
}

var testXmp = []struct {
	filename string
}{
	{"jpeg.xmp"},
	{"1.xmp"},
}

func TestXmp(t *testing.T) {
	for _, v := range testXmp {
		t.Run(v.filename, func(t *testing.T) {
			f, err := os.Open("test" + string(os.PathSeparator) + v.filename)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = f.Close()
			}()

			x, err := ParseXmp(f)
			if err != nil {
				if err != io.EOF {
					t.Fatal(err)
				}
			}

			f2, err := os.Open("test" + string(os.PathSeparator) + v.filename + ".json")
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = f2.Close()
			}()

			x1 := XMP{}
			if err = json.NewDecoder(f2).Decode(&x1); err != nil {
				t.Fatal(err)
			}
			//
			//j, err := json.Marshal(x)
			//fmt.Println(string(j))

			basicTest(t, &x.Aux, &x1.Aux)
			basicTest(t, &x.Basic, &x1.Basic)
			basicTest(t, &x.CRS, &x1.CRS)
			basicTest(t, &x.DC, &x1.DC)
			basicTest(t, &x.Exif, &x1.Exif)
			basicTest(t, &x.MM, &x1.MM)
			basicTest(t, &x.Tiff, &x1.Tiff)
		})
	}

}

func basicTest(t *testing.T, a1 interface{}, a2 interface{}) {
	defer func() {
		if x := recover(); x != nil {
			t.Error("Testing paniced for", x)
		}
	}()
	s := reflect.ValueOf(a1).Elem()
	s1 := reflect.ValueOf(a2).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f1 := s1.Field(i)
		assert.Equalf(t, f1.Interface(), f.Interface(), "error message: %s/%s", s.Type().Name(), typeOfT.Field(i).Name)
	}
}
