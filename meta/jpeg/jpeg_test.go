// Copyright (c) 2018-2023 Evan Oberholster. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package jpeg

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

var (
	dir            = "../assets/"
	benchmarksJPEG = []struct {
		fileName  string
		noExifErr bool
	}{
		{"a1.jpg", false},
		{"a2.jpg", true},
		{"JPEG.jpg", false},
		{"NoExif.jpg", true},
	}
)

func BenchmarkScanJPEG100(b *testing.B) {
	for _, bm := range benchmarksJPEG {
		f, err := os.Open(dir + bm.fileName)
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()
		buf, readErr := io.ReadAll(f)
		if readErr != nil {
			b.Fatal(readErr)
		}
		r := bytes.NewReader(buf)
		b.ReportAllocs()
		b.ResetTimer()

		b.Run(bm.fileName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err = r.Seek(0, 0); err != nil {
					b.Fatal(err)
				}
				if err := ScanJPEG(r, nil, nil); err != nil {
					if !bm.noExifErr {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

func TestScanJPEG(t *testing.T) {
	testJPEGs := []struct {
		filename string
		exif     bool
		header   meta.ExifHeader
		width    uint32
		height   uint32
	}{
		{"../assets/JPEG.jpg", true, meta.NewExifHeader(utils.LittleEndian, 13746, 12, 13872, imagetype.ImageJPEG), 1000, 563},
		{"../assets/NoExif.jpg", true, meta.NewExifHeader(utils.BigEndian, 8, 30, 140, imagetype.ImageJPEG), 50, 50},
		{"../assets/a2.jpg", false, meta.NewExifHeader(utils.LittleEndian, 13746, 12, 13872, imagetype.ImageJPEG), 1024, 1280},
		{"../assets/a1.jpg", true, meta.NewExifHeader(utils.BigEndian, 8, 30, 752, imagetype.ImageJPEG), 389, 259},
	}

	for _, jpg := range testJPEGs {
		t.Run(jpg.filename, func(t *testing.T) {
			// Open file
			f, err := os.Open(jpg.filename)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skip(err)
				}
				t.Fatal(err)
			}
			defer f.Close()

			testExifHeaderfn := func(r io.Reader, eh meta.ExifHeader) error {
				metaExifHeaderEqual(t, jpg.header, eh)
				return nil
			}
			testXmpHeaderFn := func(r io.Reader) error {
				return nil
			}

			err = ScanJPEG(f, testExifHeaderfn, testXmpHeaderFn)
			if jpg.exif && err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestScanJPEGWithReaderAtUsesIndependentExifSection(t *testing.T) {
	data := testJPEG(
		testSegment(markerAPP1, append(append([]byte(exifPrefix), testTIFFHeader()...), bytes.Repeat([]byte{0xa5}, 96)...)),
		testSegment(markerAPP1, append([]byte(xmpPrefix), []byte("<x:xmpmeta></x:xmpmeta>")...)),
	)

	var sawExif bool
	var sawXMP bool
	stream := onlyReader{r: bytes.NewReader(data)}
	err := ScanJPEGWithReaderAt(stream, bytes.NewReader(data), func(r io.Reader, h meta.ExifHeader) error {
		sawExif = true
		if h.TiffHeaderOffset != 12 {
			t.Fatalf("TiffHeaderOffset = %d, want 12", h.TiffHeaderOffset)
		}
		seeker, ok := r.(io.Seeker)
		if !ok {
			t.Fatalf("Exif reader = %T, want io.Seeker", r)
		}
		buf := make([]byte, 8)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return err
		}
		again := make([]byte, 8)
		if _, err := io.ReadFull(r, again); err != nil {
			return err
		}
		if !bytes.Equal(buf, again) {
			t.Fatalf("seek reread mismatch: %x != %x", buf, again)
		}
		return nil
	}, func(r io.Reader) error {
		sawXMP = true
		buf, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		if !bytes.HasPrefix(buf, []byte("<x:xmpmeta")) {
			t.Fatalf("XMP payload = %q", buf)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sawExif || !sawXMP {
		t.Fatalf("callbacks: exif=%t xmp=%t", sawExif, sawXMP)
	}
}

func TestScanJPEGStreamingExifCallbackDrainsRemainder(t *testing.T) {
	data := testJPEG(
		testSegment(markerAPP1, append(append([]byte(exifPrefix), testTIFFHeader()...), bytes.Repeat([]byte{0xb6}, 96)...)),
		testSegment(markerAPP1, append([]byte(xmpPrefix), []byte("<x:xmpmeta></x:xmpmeta>")...)),
	)

	var sawXMP bool
	err := ScanJPEG(onlyReader{r: bytes.NewReader(data)}, func(r io.Reader, _ meta.ExifHeader) error {
		buf := make([]byte, 8)
		_, err := io.ReadFull(r, buf)
		return err
	}, func(r io.Reader) error {
		sawXMP = true
		_, err := io.ReadAll(r)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sawXMP {
		t.Fatal("XMP callback was not called after partial Exif read")
	}
}

func TestScanJPEGDoesNotStopAtDQT(t *testing.T) {
	data := testJPEG(
		testSegment(markerDQT, bytes.Repeat([]byte{0}, 64)),
		testSegment(markerAPP1, append(append([]byte(exifPrefix), testTIFFHeader()...), bytes.Repeat([]byte{0xc7}, 96)...)),
	)

	var sawExif bool
	err := ScanJPEG(onlyReader{r: bytes.NewReader(data)}, func(r io.Reader, _ meta.ExifHeader) error {
		sawExif = true
		_, err := io.ReadAll(r)
		return err
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !sawExif {
		t.Fatal("Exif callback was not called after DQT")
	}
}

func TestScanJPEGShortEOIAfterImageData(t *testing.T) {
	data := []byte{0xff, byte(markerSOI)}
	data = append(data, testSegment(markerAPP1, append(append([]byte(exifPrefix), testTIFFHeader()...), bytes.Repeat([]byte{0xd8}, 32)...))...)
	data = append(data, testSegment(markerSOS, []byte{3, 1, 0, 2, 0x11, 3, 0x11, 0, 0x3f, 0})...)
	data = append(data, 0, 0x1f, 0xff, byte(markerEOI))

	var sawExif bool
	err := ScanJPEG(onlyReader{r: bytes.NewReader(data)}, func(r io.Reader, _ meta.ExifHeader) error {
		sawExif = true
		_, err := io.ReadAll(r)
		return err
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !sawExif {
		t.Fatal("Exif callback was not called")
	}
}

func TestScanJPEGExtendedXMP(t *testing.T) {
	const guid = "0123456789abcdef0123456789abcdef"
	part0 := []byte("<x:xmpmeta>")
	part1 := []byte("</x:xmpmeta>")
	size := len(part0) + len(part1)
	data := testJPEG(
		testExtendedXMPSegment(guid, uint32(size), uint32(len(part0)), part1),
		testExtendedXMPSegment(guid, uint32(size), 0, part0),
	)

	var got []byte
	err := ScanJPEG(onlyReader{r: bytes.NewReader(data)}, nil, func(r io.Reader) error {
		var err error
		got, err = io.ReadAll(r)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := append(append([]byte(nil), part0...), part1...); !bytes.Equal(got, want) {
		t.Fatalf("extended XMP = %q, want %q", got, want)
	}
}

func TestScanJPEG_nextMarkerNoInfiniteLoop(t *testing.T) {
	inputs := [][]byte{
		{0xFF, 0x00},
		{0xFF, 0xFF},
	}
	for _, input := range inputs {
		name := fmt.Sprintf("%02x_%02x", input[0], input[1])
		t.Run(name, func(t *testing.T) {
			done := make(chan error, 1)
			go func() {
				done <- ScanJPEG(bytes.NewReader(input), nil, nil)
			}()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			case <-time.After(time.Second):
				t.Fatal("ScanJPEG hung")
			}
		})
	}
}

func TestScanJPEG_metadataScanLimit(t *testing.T) {
	data := make([]byte, 0, 2+maxMetadataScanBytes+1024)
	data = append(data, 0xFF, byte(markerSOI))
	data = append(data, bytes.Repeat([]byte{0}, maxMetadataScanBytes+1024)...)

	start := time.Now()
	err := ScanJPEG(bytes.NewReader(data), nil, nil)
	if !errors.Is(err, ErrMetadataScanLimit) {
		t.Fatalf("err = %v, want ErrMetadataScanLimit", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("scan took %v, expected to hit limit quickly", elapsed)
	}
}

func TestScanJPEG_contextCancel(t *testing.T) {
	data := append([]byte{0xFF, byte(markerSOI)}, bytes.Repeat([]byte{0}, 512*1024)...)
	r := &slowReader{r: bytes.NewReader(data), delay: 5 * time.Millisecond}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := ScanJPEGContext(ctx, r, nil, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("cancel took %v, expected prompt return", elapsed)
	}
}

func TestScanMetadataSyntheticAPPFamilies(t *testing.T) {
	data := testJPEG(
		testSegment(markerAPP0, testJFIFPayload(1, 2, 1, 300, 200)),
		testSegment(markerAPP2, testMPFPayload(2, 100, 200)),
		testSegment(markerAPP13, testPhotoshopPayload(
			testPhotoshopResource(0x040a, []byte{1}),
			testPhotoshopResource(0x0404, append(
				testIPTCDataset(2, 25, []byte("codex")),
				testIPTCDataset(2, 80, []byte("Jane Doe"))...,
			)),
		)),
		testSegment(markerAPP14, testAdobePayload(100, 1, 2, 1)),
	)

	got, err := ScanMetadata(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if got.JFIF == nil || got.JFIF.MajorVersion != 1 || got.JFIF.MinorVersion != 2 ||
		got.JFIF.XResolution != 300 || got.JFIF.YResolution != 200 {
		t.Fatalf("JFIF = %+v", got.JFIF)
	}
	if got.MPF == nil || got.MPF.NumberOfImages != 2 || len(got.MPF.Images) != 2 {
		t.Fatalf("MPF = %+v", got.MPF)
	}
	if got.MPF.Images[0].MPImageStart != 128 || got.MPF.Images[1].MPImageStart != 228 {
		t.Fatalf("MPImageStart = %d/%d, want 128/228", got.MPF.Images[0].MPImageStart, got.MPF.Images[1].MPImageStart)
	}
	if got.Photoshop == nil || !got.Photoshop.CopyrightFlagSet || got.Photoshop.CopyrightFlag != 1 {
		t.Fatalf("Photoshop = %+v", got.Photoshop)
	}
	if got.IPTC == nil || len(got.IPTC.Keywords) != 1 || got.IPTC.Keywords[0] != "codex" || got.IPTC.ByLine != "Jane Doe" {
		t.Fatalf("IPTC = %+v", got.IPTC)
	}
	if got.Adobe == nil || got.Adobe.DCTEncodeVersion != 100 || got.Adobe.ColorTransform != 1 {
		t.Fatalf("Adobe = %+v", got.Adobe)
	}
}

func TestScanMetadataAssemblesICCChunks(t *testing.T) {
	profile := testICCProfile()
	split := 70
	data := testJPEG(
		testSegment(markerAPP2, testICCPayload(2, 2, profile[split:])),
		testSegment(markerAPP2, testICCPayload(1, 2, profile[:split])),
	)

	got, err := ScanMetadata(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if got.ICC == nil {
		t.Fatal("ICC profile was not parsed")
	}
	if got.ICC.ProfileFileSignature != "acsp" || got.ICC.ProfileVersion != 4 ||
		got.ICC.ProfileClass != "mntr" || got.ICC.ColorSpaceData != "RGB " {
		t.Fatalf("ICC = %+v", got.ICC)
	}
}

func TestScanMetadataCanonCIFFSample(t *testing.T) {
	f, err := os.Open("../../download_samples/Canon/Canon/CanonPowerShotA5.jpg")
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip(err)
		}
		t.Fatal(err)
	}
	defer f.Close()

	got, err := ScanMetadataWithReaderAt(f, f)
	if err != nil {
		t.Fatal(err)
	}
	if got.JFIF == nil || got.JFIF.XResolution != 180 || got.JFIF.YResolution != 180 {
		t.Fatalf("JFIF = %+v", got.JFIF)
	}
	if got.CIFF == nil {
		t.Fatal("CIFF was not parsed")
	}
	if got.CIFF.FileFormat != 65536 ||
		got.CIFF.ImageWidth != 512 ||
		got.CIFF.ImageHeight != 384 ||
		got.CIFF.FileNumber != 45 ||
		got.CIFF.Model != "Canon PowerShot A5" {
		t.Fatalf("CIFF = %+v", got.CIFF)
	}
}

func TestParseCIFFCRW(t *testing.T) {
	f, err := os.Open("../../testImages/CRW.CRW")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := ParseCIFFReader(f)
	if err != nil {
		t.Fatal(err)
	}
	if got.Make != "Canon" {
		t.Fatalf("Make = %q, want Canon", got.Make)
	}
	if got.Model != "Canon EOS D30" {
		t.Fatalf("Model = %q, want Canon EOS D30", got.Model)
	}
	if got.ImageWidth == 0 || got.ImageHeight == 0 {
		t.Fatalf("image dimensions = %dx%d, want non-zero", got.ImageWidth, got.ImageHeight)
	}
	if got.DateTimeOriginal.IsZero() {
		t.Fatal("DateTimeOriginal is zero")
	}
	if got.FocalLength == 0 {
		t.Fatal("FocalLength is zero")
	}
}

func TestParseCIFFReaderWithoutReaderAt(t *testing.T) {
	f, err := os.Open("../../testImages/CRW.CRW")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := ParseCIFFReader(readSeekerOnly{f})
	if err != nil {
		t.Fatal(err)
	}
	if got.Model != "Canon EOS D30" {
		t.Fatalf("Model = %q, want Canon EOS D30", got.Model)
	}
}

func TestScanMetadataCanonAPPFamilies(t *testing.T) {
	tests := []struct {
		name string
		path string
		want func(t *testing.T, m Metadata)
	}{
		{
			name: "MPF",
			path: "../../download_samples/Canon/Canon/CanonEOS_R8.jpg",
			want: func(t *testing.T, m Metadata) {
				if m.MPF == nil || m.MPF.NumberOfImages != 2 || len(m.MPF.Images) != 2 {
					t.Fatalf("MPF = %+v", m.MPF)
				}
				if got := m.MPF.Images[1].MPImageStart; got != 6063104 {
					t.Fatalf("MPImageStart = %d, want 6063104", got)
				}
			},
		},
		{
			name: "ICC/JFIF",
			path: "../../download_samples/Canon/Canon/CanonMP220.jpg",
			want: func(t *testing.T, m Metadata) {
				if m.JFIF == nil || m.JFIF.XResolution != 301 {
					t.Fatalf("JFIF = %+v", m.JFIF)
				}
				if m.ICC == nil || m.ICC.ProfileVersion != 528 || m.ICC.ProfileDescription != "sRGB IEC61966-2.1" {
					t.Fatalf("ICC = %+v", m.ICC)
				}
			},
		},
		{
			name: "Photoshop/IPTC/Adobe",
			path: "../../download_samples/Canon/Canon/CanonIXUS285HS.jpg",
			want: func(t *testing.T, m Metadata) {
				if m.Photoshop == nil || m.Photoshop.PhotoshopThumbnailLength != 12685 || m.Photoshop.IPTCDigest == "" {
					t.Fatalf("Photoshop = %+v", m.Photoshop)
				}
				if m.IPTC == nil || len(m.IPTC.Keywords) != 12 || m.IPTC.ByLine != "Dave Stevenson" {
					t.Fatalf("IPTC = %+v", m.IPTC)
				}
				if m.Adobe == nil || m.Adobe.ColorTransform != 1 {
					t.Fatalf("Adobe = %+v", m.Adobe)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skip(err)
				}
				t.Fatal(err)
			}
			defer f.Close()
			got, err := ScanMetadataWithReaderAt(f, f)
			if err != nil {
				t.Fatal(err)
			}
			tt.want(t, got)
		})
	}
}

type readSeekerOnly struct {
	io.ReadSeeker
}

type onlyReader struct {
	r *bytes.Reader
}

func (r onlyReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

type slowReader struct {
	r     io.Reader
	delay time.Duration
}

func (s *slowReader) Read(p []byte) (int, error) {
	time.Sleep(s.delay)
	return s.r.Read(p)
}

func testJPEG(segments ...[]byte) []byte {
	out := []byte{0xff, byte(markerSOI)}
	for _, segment := range segments {
		out = append(out, segment...)
	}
	out = append(out, testSegment(markerSOS, bytes.Repeat([]byte{0}, 64))...)
	return out
}

func testSegment(marker markerType, payload []byte) []byte {
	out := []byte{0xff, byte(marker), 0, 0}
	binary.BigEndian.PutUint16(out[2:4], uint16(len(payload)+2))
	return append(out, payload...)
}

func testTIFFHeader() []byte {
	return []byte{'I', 'I', '*', 0, 8, 0, 0, 0}
}

func testJFIFPayload(major, minor, unit uint8, xres, yres uint16) []byte {
	payload := make([]byte, 14)
	copy(payload, []byte(jfifPrefix))
	payload[5] = major
	payload[6] = minor
	payload[7] = unit
	binary.BigEndian.PutUint16(payload[8:10], xres)
	binary.BigEndian.PutUint16(payload[10:12], yres)
	return payload
}

func testAdobePayload(version, flags0, flags1 uint16, transform uint8) []byte {
	payload := make([]byte, 12)
	copy(payload, []byte(adobePrefix))
	binary.BigEndian.PutUint16(payload[5:7], version)
	binary.BigEndian.PutUint16(payload[7:9], flags0)
	binary.BigEndian.PutUint16(payload[9:11], flags1)
	payload[11] = transform
	return payload
}

func testMPFPayload(imageCount uint32, starts ...uint32) []byte {
	const (
		ifdOffset = 8
		tagCount  = 3
		entrySize = 12
	)
	entriesStart := ifdOffset + 2
	nextIFDStart := entriesStart + tagCount*entrySize
	mpEntryOffset := nextIFDStart + 4
	tiffLen := mpEntryOffset + len(starts)*16
	tiff := make([]byte, tiffLen)
	copy(tiff, []byte{'I', 'I', '*', 0})
	binary.LittleEndian.PutUint32(tiff[4:8], ifdOffset)
	binary.LittleEndian.PutUint16(tiff[ifdOffset:ifdOffset+2], tagCount)

	putEntry := func(index int, tagID, typ uint16, count uint32, raw []byte) {
		pos := entriesStart + index*entrySize
		binary.LittleEndian.PutUint16(tiff[pos:pos+2], tagID)
		binary.LittleEndian.PutUint16(tiff[pos+2:pos+4], typ)
		binary.LittleEndian.PutUint32(tiff[pos+4:pos+8], count)
		copy(tiff[pos+8:pos+12], raw)
	}

	putEntry(0, 0xb000, 7, 4, []byte("0100"))
	var raw [4]byte
	binary.LittleEndian.PutUint32(raw[:], imageCount)
	putEntry(1, 0xb001, 4, 1, raw[:])
	binary.LittleEndian.PutUint32(raw[:], uint32(mpEntryOffset))
	putEntry(2, 0xb002, 7, uint32(len(starts)*16), raw[:])

	for i, start := range starts {
		pos := mpEntryOffset + i*16
		binary.LittleEndian.PutUint32(tiff[pos:pos+4], 0x02000000)
		binary.LittleEndian.PutUint32(tiff[pos+4:pos+8], 10)
		binary.LittleEndian.PutUint32(tiff[pos+8:pos+12], start)
	}
	return append([]byte(mpfPrefix), tiff...)
}

func testPhotoshopPayload(resources ...[]byte) []byte {
	payload := []byte(photoshopPrefix)
	for _, resource := range resources {
		payload = append(payload, resource...)
	}
	return payload
}

func testPhotoshopResource(id uint16, data []byte) []byte {
	out := []byte{'8', 'B', 'I', 'M', 0, 0, 0, 0}
	binary.BigEndian.PutUint16(out[4:6], id)
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(data)))
	out = append(out, size[:]...)
	out = append(out, data...)
	if len(out)&1 != 0 {
		out = append(out, 0)
	}
	return out
}

func testIPTCDataset(record, dataset uint8, value []byte) []byte {
	out := []byte{0x1c, record, dataset, 0, 0}
	binary.BigEndian.PutUint16(out[3:5], uint16(len(value)))
	return append(out, value...)
}

func testICCPayload(seq, total uint8, chunk []byte) []byte {
	payload := append([]byte(iccPrefix+"\x00"), seq, total)
	return append(payload, chunk...)
}

func testICCProfile() []byte {
	profile := make([]byte, 132)
	copy(profile[4:8], []byte("lcms"))
	binary.BigEndian.PutUint32(profile[8:12], 0x00040000)
	copy(profile[12:16], []byte("mntr"))
	copy(profile[16:20], []byte("RGB "))
	copy(profile[20:24], []byte("XYZ "))
	copy(profile[36:40], []byte("acsp"))
	return profile
}

func testExtendedXMPSegment(guid string, size, offset uint32, chunk []byte) []byte {
	payload := make([]byte, 0, xmpExtHeaderLen+len(chunk))
	payload = append(payload, []byte(xmpPrefixExt)...)
	payload = append(payload, []byte(guid)...)
	var word [4]byte
	binary.BigEndian.PutUint32(word[:], size)
	payload = append(payload, word[:]...)
	binary.BigEndian.PutUint32(word[:], offset)
	payload = append(payload, word[:]...)
	payload = append(payload, chunk...)
	return testSegment(markerAPP1, payload)
}

func metaExifHeaderEqual(t *testing.T, h1 meta.ExifHeader, h2 meta.ExifHeader) {
	if h1.ByteOrder != h2.ByteOrder {
		t.Errorf("Incorrect Byte Order wanted %s got %s", h1.ByteOrder, h2.ByteOrder)
	}
	if h1.FirstIfdOffset != h2.FirstIfdOffset {
		t.Errorf("Incorrect first Ifd Offset wanted %d got %d ", h1.FirstIfdOffset, h2.FirstIfdOffset)
	}
	if h1.TiffHeaderOffset != h2.TiffHeaderOffset {
		t.Errorf("Incorrect tiff Header Offset wanted %d got %d ", h1.TiffHeaderOffset, h2.TiffHeaderOffset)
	}
	if h1.ExifLength != h2.ExifLength {
		t.Errorf("Incorrect Exif Length wanted %d got %d ", h1.ExifLength, h2.ExifLength)
	}
	if h1.ImageType != h2.ImageType {
		t.Errorf("Incorrect Exif Header Imagetype wanted %s got %s ", h1.ImageType, h2.ImageType)
	}
	if !h2.IsValid() {
		t.Errorf("Wanted valid tiff Header")
	}

}

//func TestScanMarkers(t *testing.T) {
//	data := []byte{0, markerFirstByte, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	r := bytes.NewReader(data)
//	m := jpegReader{br: bufio.NewReader(r)}
//
//	// Test discard
//	m.discard(0)
//	if m.discarded != 0 {
//		t.Errorf("Incorrect Metadata.discard wanted %d got %d", 0, m.discarded)
//	}
//	// Test Scan Markers
//	buf, _ := m.br.Peek(16)
//	err := m.scanMarkers(buf)
//	if err != nil {
//		t.Errorf("Incorrect Scan Markers error wanted %s got %s", err.Error(), ErrNoJPEGMarker)
//	}
//
//	data = []byte{markerFirstByte, markerSOI, markerFirstByte, markerEOI, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//	r = bytes.NewReader(data)
//	m = jpegReader{br: bufio.NewReader(r)}
//
//	// Test SOI
//	buf, _ = m.br.Peek(16)
//	err = m.scanMarkers(buf)
//	if m.discarded != 2 || m.pos != 1 || err != nil {
//		t.Errorf("Incorrect JPEG Start of Image error wanted discarded %d got %d", 2, m.discarded)
//	}
//
//	// Test EOI
//	buf, _ = m.br.Peek(16)
//	err = m.scanMarkers(buf)
//	if m.discarded != 4 || m.pos != 0 || err != nil {
//		t.Errorf("Incorrect JPEG End of Image error wanted discarded %d got %d", 4, m.discarded)
//	}
//
//	// Test Scan JPEG
//	err = ScanJPEG(bytes.NewReader(data), nil, nil)
//	if err != ErrNoJPEGMarker {
//		t.Errorf("Incorrect JPEG error at discarded %d wanted %s got %s", m.discarded, ErrNoJPEGMarker, err.Error())
//	}
//}
