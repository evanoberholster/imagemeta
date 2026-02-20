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
		imageType FileType
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

			if _, err = f.Seek(0, 0); err != nil {
				t.Error(err)
			}
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
	cases := map[FileType]struct {
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
		ImageRAW:     {"raw", "image/raw"},
		ImageTiff:    {"tiff", "image/tiff"},
		ImageDNG:     {"dng", "image/x-adobe-dng"},
		ImageNEF:     {"nef", "image/x-nikon-nef"},
		ImagePanaRAW: {"rw2", "image/x-panasonic-raw"},
		ImageARW:     {"arw", "image/x-sony-arw"},
		ImageCRW:     {"crw", "image/x-canon-crw"},
		ImageGPR:     {"gpr", "image/x-gopro-gpr"},
		ImageCR3:     {"cr3", "image/x-canon-cr3"},
		ImageCR2:     {"cr2", "image/x-canon-cr2"},
		ImagePSD:     {"psd", "image/vnd.adobe.photoshop"},
		ImageXMP:     {"xmp", "application/rdf+xml"},
		ImageAVIF:    {"avif", "image/avif"},
		ImagePPM:     {"ppm", "image/x-portable-pixmap"},
		ImageHEIC:    {"heic", "image/heic"},
		ImageJXR:     {"jxr", "image/vnd.ms-photo"},
		ImageFITS:    {"fits", "image/fits"},
		ImageDCM:     {"dcm", "application/dicom"},
	}

	for it, exp := range cases {
		if it.FileTypeExtension() != FileTypeExtension(exp.ext) {
			t.Errorf("%d.FileTypeExtension() returned '%s', '%s' expected", it, it.FileTypeExtension(), exp.ext)
		}

		if it.MIMEType() != MIMEType(exp.mt) {
			t.Errorf("%d.MIMEType() returned '%s', '%s' expected", it, it.MIMEType(), exp.mt)
		}
	}
}

func TestImageTypeFamilyAndContainer(t *testing.T) {
	cases := []struct {
		imageType FileType
		mediaType MediaType
		baseType  BaseType
	}{
		{ImageUnknown, MediaTypeUnknown, BaseTypeUnknown},
		{ImageJPEG, MediaTypeRaster, BaseTypeJPEG},
		{ImageAPNG, MediaTypeRaster, BaseTypePNG},
		{ImageHEIC, MediaTypeRaster, BaseTypeISOBMFF},
		{ImageAVIF, MediaTypeRaster, BaseTypeISOBMFF},
		{ImageDNG, MediaTypeRaw, BaseTypeTIFF},
		{ImageCR3, MediaTypeRaw, BaseTypeISOBMFF},
		{ImageXMP, MediaTypeMetadata, BaseTypeXML},
		{ImageSVG, MediaTypeVector, BaseTypeSVG},
		{ImagePPM, MediaTypeRaster, BaseTypeNetpbm},
	}

	for _, tc := range cases {
		if got := tc.imageType.MediaType(); got != tc.mediaType {
			t.Errorf("%s.MediaType() = %s, expected %s", tc.imageType, got, tc.mediaType)
		}
		if got := tc.imageType.Family(); got != tc.mediaType {
			t.Errorf("%s.Family() = %s, expected %s", tc.imageType, got, tc.mediaType)
		}
		if got := tc.imageType.BaseType(); got != tc.baseType {
			t.Errorf("%s.BaseType() = %s, expected %s", tc.imageType, got, tc.baseType)
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
	it = FileType(255)
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

	for input, expected := range map[string]FileType{
		"jpg":                       ImageJPEG,
		"foo/bar/photo.jpeg":        ImageJPEG,
		"image/jpeg; charset=utf-8": ImageJPEG,
		"image/heic":                ImageHEIC,
		"image/heif; q=1.0":         ImageHEIF,
		"image/jp2":                 ImageJP2K,
		".dcm":                      ImageDCM,
	} {
		if got := FromString(input); got != expected {
			t.Errorf("FromString(%q) = %s, expected %s", input, got, expected)
		}
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
		{".JP2/JPEG2000", "0.jp2", "image/jp2"},
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

func TestBufDetectsJPEGXL(t *testing.T) {
	container := []byte{
		0x00, 0x00, 0x00, 0x0C, // box size
		0x4A, 0x58, 0x4C, 0x20, // "JXL "
		0x0D, 0x0A, 0x87, 0x0A, // signature
	}
	codestream := []byte{
		0xFF, 0x0A, // codestream magic
	}

	for _, testCase := range []struct {
		name string
		buf  []byte
	}{
		{name: "container", buf: container},
		{name: "codestream", buf: codestream},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			buf := make([]byte, searchHeaderLength)
			copy(buf, testCase.buf)

			imageType, err := Buf(buf)
			if err != nil {
				t.Fatalf("Buf() returned unexpected error: %v", err)
			}
			if imageType != ImageJXL {
				t.Fatalf("Buf() = %s, expected %s", imageType, ImageJXL)
			}
		})
	}
}

func TestBufDetectsAdditionalMagicNumbers(t *testing.T) {
	cases := []struct {
		name     string
		header   []byte
		expected FileType
	}{
		{name: "ICO", header: []byte{0x00, 0x00, 0x01, 0x00}, expected: ImageICO},
		{name: "CUR", header: []byte{0x00, 0x00, 0x02, 0x00}, expected: ImageCUR},
		{name: "DDS", header: []byte("DDS "), expected: ImageDDS},
		{name: "EXR", header: []byte{0x76, 0x2F, 0x31, 0x01}, expected: ImageEXR},
		{name: "DPX/BigEndian", header: []byte("SDPX"), expected: ImageDPX},
		{name: "DPX/LittleEndian", header: []byte("XPDS"), expected: ImageDPX},
		{name: "MNG", header: []byte{0x8A, 0x4D, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, expected: ImageMNG},
		{name: "JNG", header: []byte{0x8B, 0x4A, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, expected: ImageJNG},
		{name: "FITS", header: []byte("SIMPLE  ="), expected: ImageFITS},
		{name: "RAF", header: []byte("FUJIFILMCCD-RAW "), expected: ImageRAF},
		{name: "XCF", header: []byte("gimp xcf "), expected: ImageXCF},
		{name: "FLIF", header: []byte("FLIF"), expected: ImageFLIF},
		{name: "BPG", header: []byte{0x42, 0x50, 0x47, 0xFB}, expected: ImageBPG},
		{name: "HDR/Radiance", header: []byte("#?RADIANCE"), expected: ImageHDR},
		{name: "HDR/RGBE", header: []byte("#?RGBE"), expected: ImageHDR},
		{
			name: "DJVU",
			header: []byte{
				'A', 'T', '&', 'T', 'F', 'O', 'R', 'M',
				0x00, 0x00, 0x00, 0x00,
				'D', 'J', 'V', 'U',
			},
			expected: ImageDJVU,
		},
		{name: "PBM", header: []byte("P1\n"), expected: ImagePBM},
		{name: "PGM", header: []byte("P5 "), expected: ImagePGM},
		{name: "PPM", header: []byte("P6 "), expected: ImagePPM},
		{name: "PAM", header: []byte("P7\t"), expected: ImagePAM},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := make([]byte, searchHeaderLength)
			copy(buf, tc.header)

			got, err := Buf(buf)
			if err != nil {
				t.Fatalf("Buf() returned unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("Buf() = %s, expected %s", got, tc.expected)
			}
		})
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
