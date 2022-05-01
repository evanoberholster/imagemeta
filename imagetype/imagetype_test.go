package imagetype

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

// Tests
func TestScan(t *testing.T) {
	exifHeaderTests := []struct {
		filename  string
		imageType ImageType
	}{
		{"../testImages/ARW.exif", ImageTiff},
		{"../testImages/NEF.exif", ImageTiff},
		{"../testImages/CR2.exif", ImageCR2},
		{"../testImages/Heic.exif", ImageHEIF},
		{"../testImages/AVIF.avif", ImageAVIF},
		{"../testImages/AVIF2.avif", ImageAVIF},
		{"../testImages/CRW.CRW", ImageCRW},
		{"../testImages/XMP.xmp", ImageXMP},
		{"../testImages/GIF.gif", ImageGIF},
		{"../testImages/Unknown.exif", ImageUnknown},
		{"../testImages/ppm-ascii.ppm", ImagePPM},
		{"../testImages/ppm-raw.ppm", ImagePPM},
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
			if header.imageType == ImageUnknown {
				if err != ErrImageTypeNotFound {
					t.Fatal(err)
				}
			}

			if header.imageType != imageType {
				t.Errorf("Incorrect Imagetype wanted %s got %s", header.imageType.String(), imageType.String())
			}

			f.Seek(0, 0)
			imageType, err = ReadAt(f)
			if header.imageType == ImageUnknown {
				if err != ErrImageTypeNotFound {
					t.Fatal(err)
				}
			}

			if header.imageType != imageType {
				t.Errorf("Incorrect Imagetype wanted %s got %s", header.imageType.String(), imageType.String())
			}
		})
	}
}

func TestImageTypeIndices(t *testing.T) {
	cases := map[ImageType]struct {
		ext string
		mt  string
	}{
		ImageUnknown: {"", "application/octet-stream"},
		ImageJPEG:    {"jpg", "image/jpeg"},
		ImagePNG:     {"png", "image/png"},
		ImageGIF:     {"gif", "image/gif"},
		ImageBMP:     {"bmp", "image/bmp"},
		ImageWebP:    {"webp", "image/webp"},
		ImageHEIF:    {"heif", "image/heif"},
		ImageRAW:     {"RAW", "image/raw"},
		ImageTiff:    {"TIFF", "image/tiff"},
		ImageDNG:     {"DNG", "image/x-adobe-dng"},
		ImageNEF:     {"NEF", "image/x-nikon-nef"},
		ImagePanaRAW: {"RW2", "image/x-panasonic-raw"},
		ImageARW:     {"ARW", "image/x-sony-arw"},
		ImageCRW:     {"CRW", "image/x-canon-crw"},
		ImageGPR:     {"GPR", "image/x-gopro-gpr"},
		ImageCR3:     {"CR3", "image/x-canon-cr3"},
		ImageCR2:     {"CR2", "image/x-canon-cr2"},
		ImagePSD:     {"PSD", "image/vnd.adobe.photoshop"},
		ImageXMP:     {"XMP", "application/rdf+xml"},
		ImageAVIF:    {"avif", "image/avif"},
		ImagePPM:     {"ppm", "image/x-portable-pixmap"},
	}

	for it, exp := range cases {
		if it.Extension() != exp.ext {
			t.Errorf("%d.Extension() returned '%s', '%s' expected", it, it.Extension(), exp.ext)
		}

		if it.String() != exp.mt {
			t.Errorf("%d.String() returned '%s', '%s' expected", it, it.String(), exp.mt)
		}
	}
}

func TestImageType(t *testing.T) {

	str := "image/jpeg"
	ext := "jpg"
	it := ImageJPEG

	if it.IsUnknown() || !ImageUnknown.IsUnknown() {
		t.Errorf("Error Imagetype should not be Unknown")
	}

	itbuf, err := it.MarshalText()
	if err != nil {
		t.Errorf("Error Imagetype could not be marshalled")
	}

	if !bytes.Equal(itbuf, []byte(str)) {
		t.Errorf("Incorrect Imagetype Marshall wanted %s got %s", str, string(itbuf))
	}

	it2 := FromString(str)
	if it2 != it {
		t.Errorf("Incorrect Imagetype FromString wanted %s got %s", str, it2.String())
	}

	err = it.UnmarshalText(itbuf)
	if err != nil {
		t.Errorf("Error Imagetype could not be unmarshalled")
	}

	if it2 != it {
		t.Errorf("Incorrect Imagetype Unmarshal wanted %s got %s", str, it.String())
	}
	if ext != it.Extension() {
		t.Errorf("Incorrect Imagetype extension wanted %s got %s", ext, it.Extension())
	}

	// Unknown
	it = ImageType(255)
	if it.Extension() != "" {
		t.Errorf("Incorrect Imagetype extension wanted %s got %s", "", it.Extension())
	}
	if it.String() != ImageUnknown.String() {
		t.Errorf("Incorrect Imagetype extension wanted %s got %s", ImageUnknown.String(), it.String())
	}

	// FromString
	it = FromString(".jpg")
	if it != ImageJPEG {
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageJPEG, it)
	}

	it = FromString("hello")
	if it != ImageUnknown {
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageUnknown, it)
	}

	err = it.EncodeMsg(msgp.NewWriterSize(&msgp.Writer{}, 0))
	if err != nil {
		t.Errorf("Incorrect Error for EncodeMsg wanted %s got %s", err, err)

	}

}

func TestScanImageType(t *testing.T) {
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
		{".PSD", "0.psd", "image/vnd.adobe.photoshop"},
		{".JP2/JPEG2000", "0.jp2", "image/jpeg"},
		{".BMP", "0.bmp", "image/bmp"},
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

	// Image Unknown
	imageType, err := Scan(bytes.NewReader([]byte("abcdefghijklmnop1234567890abcdefghijklmnopqrs")))
	if imageType != ImageUnknown && err != ErrImageTypeNotFound {
		t.Errorf("Incorrect Error wanted %s got %s", ErrImageTypeNotFound.Error(), err.Error())
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageUnknown, imageType.String())
	}

	r := bytes.NewReader([]byte(""))
	//  Image Unknown - Empty ByteSlice
	imageType, err = Scan(r)
	if err != io.EOF {
		t.Fatal(err)
	}

	if imageType != ImageUnknown {
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageUnknown, imageType.String())
	}

	imageType, err = ReadAt(r)
	if imageType != ImageUnknown && err != ErrImageTypeNotFound {
		t.Errorf("Incorrect Error wanted %s got %s", ErrImageTypeNotFound.Error(), err.Error())
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageUnknown, imageType.String())
	}

	buf = make([]byte, 10)
	imageType, err = Buf(buf)
	if imageType != ImageUnknown && err != ErrDataLength {
		t.Errorf("Incorrect Error wanted %s got %s", ErrDataLength.Error(), err.Error())
		t.Errorf("Incorrect Imagetype wanted %s got %s", ImageUnknown, imageType.String())
	}
}

func TestMsgp(t *testing.T) {
	var err error
	it, it2 := ImageJPEG, ImageUnknown
	b, _ := it.MarshalMsg(nil)

	b, _ = it2.UnmarshalMsg(b)
	if it != it2 {
		t.Errorf("Incorrect Imagetype wanted %s got %s", it, it2)
	}

	if _, err = it2.UnmarshalMsg(b); err != msgp.ErrShortBytes {
		t.Error(err)
	}

	it3 := ImageUnknown
	var buf bytes.Buffer

	if err = msgp.Encode(&buf, &it3); err != nil {
		t.Error(err)
	}

	if err = msgp.Decode(&buf, &it3); err != nil {
		t.Error(err)
	}

	// Test EOF error
	if err = msgp.Decode(&buf, &it3); !errors.Is(err, io.EOF) {
		t.Error(err)
	}

}
