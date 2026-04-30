package exif

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	sonymk "github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

func TestParseSonyMakerNoteSamples(t *testing.T) {
	benchDir := resolveSonySampleDir(t)

	cases := []struct {
		file  string
		check func(t *testing.T, got *sonymk.Sony)
	}{
		{
			file: "SonyDSLR-A200.jpg",
			check: func(t *testing.T, got *sonymk.Sony) {
				if got.CreativeStyle != "Standard" {
					t.Fatalf("CreativeStyle = %q, want %q", got.CreativeStyle, "Standard")
				}
				if got.Quality != 2 {
					t.Fatalf("Quality = %d, want %d", got.Quality, 2)
				}
				if got.ImageStabilization != 1 {
					t.Fatalf("ImageStabilization = %d, want %d", got.ImageStabilization, 1)
				}
				if got.SonyModelID != 286 {
					t.Fatalf("SonyModelID = %d, want %d", got.SonyModelID, 286)
				}
				if got.LensType != 55 {
					t.Fatalf("LensType = %d, want %d", got.LensType, 55)
				}

				if got.CameraInfo2.AFPointSelected != 0 || got.CameraInfo2.FocusModeSetting != 3 || got.CameraInfo2.AFPoint != 4 {
					t.Fatalf("CameraInfo2 = %+v, want AFPointSelected=0 FocusModeSetting=3 AFPoint=4", got.CameraInfo2)
				}
				if got.CameraInfo2.AFStatusActiveSensor != 17 || got.CameraInfo2.AFStatusRight != 18 {
					t.Fatalf("CameraInfo2 = %+v, want AFStatusActiveSensor=17 AFStatusRight=18", got.CameraInfo2)
				}

				if got.FocusInfo.DriveMode2 != 1 || got.FocusInfo.ExposureProgram != 3 || got.FocusInfo.ISO != 48 || got.FocusInfo.FocusPosition != 94 {
					t.Fatalf("FocusInfo = %+v, want DriveMode2=1 ExposureProgram=3 ISO=48 FocusPosition=94", got.FocusInfo)
				}

				if got.CameraSettings.ExposureTime != 32 || got.CameraSettings.FNumber != 53 || got.CameraSettings.WhiteBalance != 2 {
					t.Fatalf("CameraSettings = %+v, want ExposureTime=32 FNumber=53 WhiteBalance=2", got.CameraSettings)
				}
				if got.CameraSettings.FlashMode != 4 || got.CameraSettings.FocusMode != 3 || got.CameraSettings.BatteryLevel != 77 || got.CameraSettings.Quality != 32 {
					t.Fatalf("CameraSettings = %+v, want FlashMode=4 FocusMode=3 BatteryLevel=77 Quality=32", got.CameraSettings)
				}

				if got.Tag9050.ShutterCount != 73 || got.Tag9050.LensType != 55 {
					t.Fatalf("Tag9050 = %+v, want ShutterCount=73 LensType=55", got.Tag9050)
				}
			},
		},
		{
			file: "SonySLT-A65.jpg",
			check: func(t *testing.T, got *sonymk.Sony) {
				if got.CreativeStyle != "Standard" {
					t.Fatalf("CreativeStyle = %q, want %q", got.CreativeStyle, "Standard")
				}
				if got.Quality != 2 {
					t.Fatalf("Quality = %d, want %d", got.Quality, 2)
				}
				if got.SonyModelID != 286 {
					t.Fatalf("SonyModelID = %d, want %d", got.SonyModelID, 286)
				}
				if got.LensType != 55 {
					t.Fatalf("LensType = %d, want %d", got.LensType, 55)
				}

				if got.CameraInfo3.FocalLength != 35 || got.CameraInfo3.AFStatusActiveSensor != -34 || got.CameraInfo3.AFPoint != 6 || got.CameraInfo3.FocusMode != 2 {
					t.Fatalf("CameraInfo3 = %+v, want FocalLength=35 AFStatusActiveSensor=-34 AFPoint=6 FocusMode=2", got.CameraInfo3)
				}

				if got.AFInfo.AFType != 1 || got.AFInfo.AFStatusActiveSensor != -34 || got.AFInfo.AFPoint != 6 || got.AFInfo.FocusMode != 2 {
					t.Fatalf("AFInfo = %+v, want AFType=1 AFStatusActiveSensor=-34 AFPoint=6 FocusMode=2", got.AFInfo)
				}
				if got.Tag9050.ShutterCount != 73 || got.Tag9050.LensType != 55 {
					t.Fatalf("Tag9050 = %+v, want ShutterCount=73 LensType=55", got.Tag9050)
				}
			},
		},
		{
			file: "SonyZV-E1.jpg",
			check: func(t *testing.T, got *sonymk.Sony) {
				if got.CreativeStyle != "Standard" {
					t.Fatalf("CreativeStyle = %q, want %q", got.CreativeStyle, "Standard")
				}
				if got.Quality != 6 {
					t.Fatalf("Quality = %d, want %d", got.Quality, 6)
				}
				if got.SonyModelID != 393 {
					t.Fatalf("SonyModelID = %d, want %d", got.SonyModelID, 393)
				}
				if got.LensType != 65535 {
					t.Fatalf("LensType = %d, want %d", got.LensType, 65535)
				}

				if got.Tag202A.FocalPlaneAFPointsUsed != 0 {
					t.Fatalf("Tag202A = %+v, want FocalPlaneAFPointsUsed=0", got.Tag202A)
				}
				if got.HiddenInfo.HiddenDataOffset != 7995380 || got.HiddenInfo.HiddenDataLength != 53248 {
					t.Fatalf("HiddenInfo = %+v, want HiddenDataOffset=7995380 HiddenDataLength=53248", got.HiddenInfo)
				}

				if got.Tag9416.SonyExposureTime2 != 6144 || got.Tag9416.SonyFNumber2 != 5389 || got.Tag9416.ReleaseMode2 != 1 {
					t.Fatalf("Tag9416 = %+v, want SonyExposureTime2=6144 SonyFNumber2=5389 ReleaseMode2=1", got.Tag9416)
				}
				if got.Tag9416.InternalSerialNumber != [6]uint8{0, 0, 0, 0, 128, 12} {
					t.Fatalf("Tag9416.InternalSerialNumber = %v, want %v", got.Tag9416.InternalSerialNumber, [6]uint8{0, 0, 0, 0, 128, 12})
				}

				if got.AFTracking != 2 {
					t.Fatalf("AFTracking = %d, want %d", got.AFTracking, 2)
				}
				if got.AFInfo.AFType != 1 || got.AFInfo.AFPoint != 6 || got.AFInfo.FocusMode != 3 {
					t.Fatalf("AFInfo = %+v, want AFType=1 AFPoint=6 FocusMode=3", got.AFInfo)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			samplePath := filepath.Join(benchDir, tc.file)
			if _, err := os.Stat(samplePath); err != nil {
				t.Skipf("sample not found: %s", samplePath)
			}

			f, err := os.Open(samplePath)
			if err != nil {
				t.Fatalf("open %s: %v", samplePath, err)
			}
			defer func() { _ = f.Close() }()

			parsed, err := Parse(f)
			if err != nil {
				t.Fatalf("parse %s: %v", samplePath, err)
			}
			if parsed.MakerNote.Sony == nil {
				t.Fatalf("Sony maker-note missing for %s", samplePath)
			}

			tc.check(t, parsed.MakerNote.Sony)
		})
	}
}

func resolveSonySampleDir(t *testing.T) string {
	t.Helper()

	candidates := []string{
		os.Getenv("IMAGEMETA_SONY_IMAGE_DIR"),
		os.Getenv("IMAGEMETA_BENCH_IMAGE_DIR"),
		filepath.Join("..", "..", "download_samples", "Sony", "Sony"),
		filepath.Join("download_samples", "Sony", "Sony"),
	}
	for _, dir := range candidates {
		if dir == "" {
			continue
		}
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	t.Fatalf("no Sony sample directory found")
	return ""
}

func TestDecodeARWSample(t *testing.T) {
	data := syntheticARW()

	it, err := imagetype.Buf(data[:64])
	if err != nil {
		t.Fatalf("Buf() error = %v", err)
	}
	if it != imagetype.ImageARW {
		t.Fatalf("Buf() = %s, want %s", it, imagetype.ImageARW)
	}

	r := NewReader(metalog.Logger)
	defer r.Close()
	header, err := ScanTiffHeader(bufio.NewReader(bytes.NewReader(data)), it)
	if err != nil {
		t.Fatalf("ScanTiffHeader() error = %v", err)
	}
	if err := r.DecodeTiff(bytes.NewReader(data), header); err != nil {
		t.Fatalf("DecodeTiff() error = %v", err)
	}
	if r.Exif.IFD0.Make != "Sony" {
		t.Fatalf("Make = %q, want Sony", r.Exif.IFD0.Make)
	}
	if r.Exif.IFD0.Model != "ILCE-7M5" {
		t.Fatalf("Model = %q, want ILCE-7M5", r.Exif.IFD0.Model)
	}
	if r.Exif.IFD0.ImageWidth != 7008 {
		t.Fatalf("ImageWidth = %d, want 7008", r.Exif.IFD0.ImageWidth)
	}
	if r.Exif.IFD0.ImageHeight != 4672 {
		t.Fatalf("ImageHeight = %d, want 4672", r.Exif.IFD0.ImageHeight)
	}
	if r.Exif.MakerNote.Sony == nil {
		t.Fatal("Sony maker-note missing")
	}
	if r.Exif.MakerNote.Sony.SonyModelID != 393 {
		t.Fatalf("SonyModelID = %d, want 393", r.Exif.MakerNote.Sony.SonyModelID)
	}
}

func syntheticARW() []byte {
	buf := make([]byte, 1024)
	copy(buf[:4], []byte{'I', 'I', 0x2a, 0x00})
	binary.LittleEndian.PutUint32(buf[4:8], 8)

	const ifd = 8
	const count = 0x13
	binary.LittleEndian.PutUint16(buf[ifd:ifd+2], count)
	pos := ifd + 2
	putEntry := func(id, typ uint16, units, value uint32) {
		binary.LittleEndian.PutUint16(buf[pos:pos+2], id)
		binary.LittleEndian.PutUint16(buf[pos+2:pos+4], typ)
		binary.LittleEndian.PutUint32(buf[pos+4:pos+8], units)
		binary.LittleEndian.PutUint32(buf[pos+8:pos+12], value)
		pos += 12
	}

	putEntry(0x00fe, 4, 1, 1)       // SubfileType
	putEntry(0x0103, 3, 1, 0x0006)  // Compression
	putEntry(0x010e, 2, 32, 0x00f2) // ImageDescription
	putEntry(0x010f, 2, 5, 0x0112)  // Make
	putEntry(0x0110, 2, 9, 0x0118)  // Model
	putEntry(0x0112, 3, 1, 8)       // Orientation
	putEntry(0x011a, 5, 1, 0x0122)  // XResolution
	putEntry(0x011b, 5, 1, 0x012a)  // YResolution
	putEntry(0x0128, 3, 1, 2)       // ResolutionUnit
	putEntry(0x0131, 2, 15, 0x0132) // Software
	putEntry(0x0132, 2, 20, 0x0142) // DateTime
	putEntry(0x014a, 4, 1, 0)       // SubIFDs
	putEntry(0x0201, 4, 1, 0)       // ThumbnailOffset
	putEntry(0x0202, 4, 1, 0)       // ThumbnailLength
	putEntry(0x0213, 3, 1, 2)       // YCbCrPositioning
	putEntry(0x8769, 4, 1, 0x0160)  // ExifIFDPointer
	putEntry(0x0100, 4, 1, 7008)    // ImageWidth
	putEntry(0x0101, 4, 1, 4672)    // ImageLength
	putEntry(0xc634, 1, 4, 0x100e8) // SonyRawFileType, collides with DNGAdobeData
	binary.LittleEndian.PutUint32(buf[pos:pos+4], 0)

	copy(buf[0x00f2:0x0112], []byte("                               \x00"))
	copy(buf[0x0112:0x0117], []byte("SONY\x00"))
	copy(buf[0x0118:0x0121], []byte("ILCE-7M5\x00"))
	binary.LittleEndian.PutUint32(buf[0x0122:0x0126], 350)
	binary.LittleEndian.PutUint32(buf[0x0126:0x012a], 1)
	binary.LittleEndian.PutUint32(buf[0x012a:0x012e], 350)
	binary.LittleEndian.PutUint32(buf[0x012e:0x0132], 1)
	copy(buf[0x0132:0x0141], []byte("ILCE-7M5 v1.00\x00"))
	copy(buf[0x0142:0x0156], []byte("2025:11:20 13:41:11\x00"))

	binary.LittleEndian.PutUint16(buf[0x0160:0x0162], 1)
	binary.LittleEndian.PutUint16(buf[0x0162:0x0164], 0x927c) // MakerNote
	binary.LittleEndian.PutUint16(buf[0x0164:0x0166], 7)      // UNDEFINED
	binary.LittleEndian.PutUint32(buf[0x0166:0x016a], 14)
	binary.LittleEndian.PutUint32(buf[0x016a:0x016e], 0x01c0)

	binary.LittleEndian.PutUint16(buf[0x01c0:0x01c2], 1)
	binary.LittleEndian.PutUint16(buf[0x01c2:0x01c4], 0xb001) // SonyModelID
	binary.LittleEndian.PutUint16(buf[0x01c4:0x01c6], 3)      // SHORT
	binary.LittleEndian.PutUint32(buf[0x01c6:0x01ca], 1)
	binary.LittleEndian.PutUint32(buf[0x01ca:0x01ce], 393)
	return buf
}
