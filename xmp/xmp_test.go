package xmp

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMap = map[string]string{
	"test/samples/image1.jpeg.xmp": "{\"Aux\":{\"SerialNumber\":\"\",\"LensInfo\":\"\",\"Lens\":\"\",\"LensID\":0,\"LensSerialNumber\":\"\",\"ImageNumber\":0,\"ApproximateFocusDistance\":\"\",\"FlashCompensation\":\"\",\"Firmware\":\"\"},\"Exif\":{\"ExifVersion\":\"\",\"PixelXDimension\":0,\"PixelYDimension\":0,\"DateTimeOriginal\":\"2003-02-04T08:06:56Z\",\"CreateDate\":\"0001-01-01T00:00:00Z\",\"ExposureTime\":\"\",\"ExposureMode\":0,\"ShutterSpeedValue\":\"\",\"ExposureProgram\":\"\",\"ISOSpeedRatings\":50,\"Flash\":{\"Fired\":false,\"Mode\":3,\"RedEyeMode\":false,\"Function\":false,\"Return\":0}},\"Tiff\":{\"Make\":\"OLYMPUS CORPORATION\",\"Model\":\"C750UZ\",\"Software\":\"\",\"Copyright\":null,\"ImageDescription\":null,\"ImageWidth\":2288,\"ImageLength\":1712,\"Orientation\":1},\"Basic\":{\"CreateDate\":\"2007-08-16T11:57:04+01:00\",\"CreatorTool\":\"Adobe Photoshop CS3 Windows\",\"Label\":\"\",\"MetadataDate\":\"2007-08-16T11:57:04+01:00\",\"ModifyDate\":\"2007-08-16T11:57:04+01:00\",\"Rating\":0},\"DC\":{\"Contributor\":null,\"Coverage\":\"\",\"Creator\":[\"\",\"XMP SDK\"],\"Date\":\"0001-01-01T00:00:00Z\",\"Description\":null,\"Format\":\"image/jpeg\",\"Identifier\":\"\",\"Language\":null,\"Rights\":null,\"Source\":\"\",\"Subject\":null,\"Title\":[\"Un titre Francais\"]},\"CRS\":{\"RawFileName\":\"\"},\"MM\":{\"DocumentID\":\"544d6a6b-e74b-dc11-9e68-d4e6c4c1b201\",\"InstanceID\":\"554d6a6b-e74b-dc11-9e68-d4e6c4c1b201\",\"OriginalDocumentID\":\"00000000-0000-0000-0000-000000000000\"}}",
}

func TestImageJpeg(t *testing.T) {
	f, err := os.Open("test/samples/image1.jpeg.xmp")
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println(f.Close())
	}()

	x, err := Read(f)
	if err != nil {
		fmt.Println(err)
	}

	x1 := XMP{}
	if err = json.Unmarshal([]byte(testMap["test/samples/image1.jpeg.xmp"]), &x1); err != nil {
		t.Fatal(err)
	}

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
