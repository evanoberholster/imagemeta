// +build !windows

package xmp

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMap = map[string]string{
	"jpeg.xmp": "{\"Aux\":{\"SerialNumber\":\"\",\"LensInfo\":\"\",\"Lens\":\"\",\"LensID\":0,\"LensSerialNumber\":\"\",\"ImageNumber\":0,\"ApproximateFocusDistance\":\"\",\"FlashCompensation\":\"\",\"Firmware\":\"\"},\"Exif\":{\"ExifVersion\":\"\",\"PixelXDimension\":0,\"PixelYDimension\":0,\"DateTimeOriginal\":\"2003-02-04T08:06:56Z\",\"CreateDate\":\"0001-01-01T00:00:00Z\",\"ExposureTime\":\"1/250\",\"ExposureProgram\":2,\"ExposureMode\":1,\"ExposureBias\":\"-7/10\",\"ISOSpeedRatings\":50,\"Flash\":{\"Fired\":false,\"Mode\":3,\"RedEyeMode\":false,\"Function\":false,\"Return\":0},\"MeteringMode\":\"Multi-segment\",\"Aperture\":\"\u0000\u0000\u0000\u00003.20\",\"FocalLength\":\"20.80mm\",\"GPSLatitude\":0,\"GPSLongitude\":0},\"Tiff\":{\"Make\":\"OLYMPUS CORPORATION\",\"Model\":\"C750UZ\",\"Software\":\"\",\"Copyright\":null,\"ImageDescription\":null,\"ImageWidth\":2288,\"ImageLength\":1712,\"Orientation\":1},\"Basic\":{\"CreateDate\":\"2007-08-16T11:57:04+01:00\",\"CreatorTool\":\"Adobe Photoshop CS3 Windows\",\"Label\":\"\",\"MetadataDate\":\"2007-08-16T11:57:04+01:00\",\"ModifyDate\":\"2007-08-16T11:57:04+01:00\",\"Rating\":0},\"DC\":{\"Contributor\":null,\"Coverage\":\"\",\"Creator\":[\"\",\"XMP SDK\"],\"Date\":\"0001-01-01T00:00:00Z\",\"Description\":null,\"Format\":\"image/jpeg\",\"Identifier\":\"\",\"Language\":null,\"Rights\":null,\"Source\":\"\",\"Subject\":null,\"Title\":[\"Un titre Francais\"]},\"CRS\":{\"RawFileName\":\"\"},\"MM\":{\"DocumentID\":\"544d6a6b-e74b-dc11-9e68-d4e6c4c1b201\",\"InstanceID\":\"554d6a6b-e74b-dc11-9e68-d4e6c4c1b201\",\"OriginalDocumentID\":\"00000000-0000-0000-0000-000000000000\"}}",
	"XMP.xmp":  "{\"Aux\":{\"SerialNumber\":\"1234567890abcd\",\"LensInfo\":\"16/1 35/1 0/0 0/0\",\"Lens\":\"EF16-35mm f/4L IS USM\",\"LensID\":507,\"LensSerialNumber\":\"987654321\",\"ImageNumber\":0,\"ApproximateFocusDistance\":\"\",\"FlashCompensation\":\"\",\"Firmware\":\"\"},\"Exif\":{\"ExifVersion\":\"\",\"PixelXDimension\":0,\"PixelYDimension\":0,\"DateTimeOriginal\":\"2015-10-01T06:42:43Z\",\"CreateDate\":\"0001-01-01T00:00:00Z\",\"ExposureTime\":\"\",\"ExposureMode\":0,\"ShutterSpeedValue\":\"\",\"ExposureProgram\":\"\",\"ISOSpeedRatings\":100,\"Flash\":{\"Fired\":false,\"Mode\":2,\"RedEyeMode\":false,\"Function\":false,\"Return\":0}},\"Tiff\":{\"Make\":\"Canon\",\"Model\":\"Canon EOS 6D\",\"Software\":\"\",\"Copyright\":null,\"ImageDescription\":null,\"ImageWidth\":5472,\"ImageLength\":3648,\"Orientation\":1},\"Basic\":{\"CreateDate\":\"2015-10-01T06:42:43Z\",\"CreatorTool\":\"\",\"Label\":\"\",\"MetadataDate\":\"2015-10-01T17:51:39+03:00\",\"ModifyDate\":\"2015-10-01T06:42:43Z\",\"Rating\":0},\"DC\":{\"Contributor\":null,\"Coverage\":\"\",\"Creator\":[\"\",\"Artist Name\"],\"Date\":\"0001-01-01T00:00:00Z\",\"Description\":null,\"Format\":\"image/x-canon-cr2\",\"Identifier\":\"\",\"Language\":null,\"Rights\":[\"\",\"Copyright Name\"],\"Source\":\"\",\"Subject\":null,\"Title\":null},\"CRS\":{\"RawFileName\":\"IMG_0001.CR2\"},\"MM\":{\"DocumentID\":\"00000000-0000-0000-0000-000000000000\",\"InstanceID\":\"00000000-0000-0000-0000-000000000000\",\"OriginalDocumentID\":\"00000000-0000-0000-0000-000000000000\"}}",
}

func TestGen(t *testing.T) {
	filename := "jpeg.xmp"
	f, err := os.Open("test" + string(os.PathSeparator) + filename)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	x, err := ParseXmp(f)
	if err != nil {
		t.Error(err)
	}

	j, err := json.Marshal(x)

	dat, err := os.Create("test" + string(os.PathSeparator) + filename + ".json")
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
}

func TestXmp(t *testing.T) {
	f, err := os.Open("test" + string(os.PathSeparator) + "jpeg.xmp")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	x, err := ParseXmp(f)
	if err != nil {
		t.Fatal(err)
	}

	x1 := XMP{}
	if err = json.Unmarshal([]byte(testMap["jpeg.xmp"]), &x1); err != nil {
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
