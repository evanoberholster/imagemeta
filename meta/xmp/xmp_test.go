package xmp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
)

func TestParseXMPTextFixtures(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		assert func(t *testing.T, x XMP)
	}{
		{
			name: "acr sidecar",
			file: "acr_sidecar.xmp",
			assert: func(t *testing.T, x XMP) {
				if !x.IsParsed(CrsNS) {
					t.Fatal("CRS namespace not parsed")
				}
				if x.CRS.RawFileName != "IMG_9620.CR2" {
					t.Fatalf("CRS.RawFileName = %q", x.CRS.RawFileName)
				}
				if x.DC.Format.String() != "image/x-canon-cr2" {
					t.Fatalf("DC.Format = %q", x.DC.Format.String())
				}
				if x.Basic.ModifyDate.IsZero() {
					t.Fatal("Basic.ModifyDate is zero")
				}
				if x.Basic.ModifyDate.Format(time.RFC3339) != "2012-10-17T13:07:01+03:00" {
					t.Fatalf("Basic.ModifyDate = %s", x.Basic.ModifyDate.Format(time.RFC3339))
				}
				if x.Basic.MetadataDate.Format(time.RFC3339) != "2012-10-17T13:07:01+03:00" {
					t.Fatalf("Basic.MetadataDate = %s", x.Basic.MetadataDate.Format(time.RFC3339))
				}
			},
		},
		{
			name: "dng embedded sample",
			file: "dng_embedded.xmp",
			assert: func(t *testing.T, x XMP) {
				if x.Tiff.Make != "OLYMPUS CORPORATION" {
					t.Fatalf("Tiff.Make = %q", x.Tiff.Make)
				}
				if x.Exif.PixelXDimension != 2288 || x.Exif.PixelYDimension != 1712 {
					t.Fatalf("Exif dimensions = %dx%d", x.Exif.PixelXDimension, x.Exif.PixelYDimension)
				}
				if x.Basic.CreatorTool == "" {
					t.Fatal("Basic.CreatorTool empty")
				}
				if len(x.DC.Title) == 0 {
					t.Fatal("DC.Title empty")
				}
				if !x.IsParsed(XmpMMNS) {
					t.Fatal("MM namespace not parsed")
				}
				assertApproxFloat64(t, x.Exif.ExposureTime.Float64(), 0.004, 0.000001, "Exif.ExposureTime")
				assertApproxFloat64(t, float64(x.Exif.Aperture), 3.2, 0.0001, "Exif.FNumber")
				if x.MM.DocumentID.String() != "544d6a6b-e74b-dc11-9e68-d4e6c4c1b201" {
					t.Fatalf("MM.DocumentID = %s", x.MM.DocumentID.String())
				}
				if x.MM.InstanceID.String() != "554d6a6b-e74b-dc11-9e68-d4e6c4c1b201" {
					t.Fatalf("MM.InstanceID = %s", x.MM.InstanceID.String())
				}
				if x.MM.PreservedFileName != "P2040006.TIF" {
					t.Fatalf("MM.PreservedFileName = %q", x.MM.PreservedFileName)
				}
			},
		},
		{
			name: "lightroom sidecar",
			file: "lightroom_sidecar.xmp",
			assert: func(t *testing.T, x XMP) {
				if !x.IsParsed(CrsNS) {
					t.Fatal("CRS namespace not parsed")
				}
				if x.CRS.RawFileName != "_MG_1563.CR2" {
					t.Fatalf("CRS.RawFileName = %q", x.CRS.RawFileName)
				}
				if x.Basic.Rating != 4 {
					t.Fatalf("Basic.Rating = %d", x.Basic.Rating)
				}
				if x.Tiff.Model != "Canon EOS 6D" {
					t.Fatalf("Tiff.Model = %q", x.Tiff.Model)
				}
				if x.Exif.ExifVersion != "0230" {
					t.Fatalf("Exif.ExifVersion = %q", x.Exif.ExifVersion)
				}
				if x.Exif.Flash.Fired || x.Exif.Flash.Mode != 2 || x.Exif.Flash.Return != 0 || x.Exif.Flash.RedEyeMode || x.Exif.Flash.Function {
					t.Fatalf("Exif.Flash = %+v", x.Exif.Flash)
				}
				if x.Aux.Firmware != "1.1.6" {
					t.Fatalf("Aux.Firmware = %q", x.Aux.Firmware)
				}
				if x.Basic.Toolkit == "" {
					t.Fatal("Basic.Toolkit empty")
				}
				if !x.IsParsed(PhotoshopNS) {
					t.Fatal("Photoshop namespace not parsed")
				}
				if x.Photoshop.DateCreated.IsZero() || x.Photoshop.SidecarForExtension != "CR2" || x.Photoshop.EmbeddedXMPDigest == "" {
					t.Fatalf("Photoshop = %+v", x.Photoshop)
				}
				if x.CRS.WhiteBalance != "As Shot" || x.CRS.Sharpness != 40 || x.CRS.Saturation != 0 {
					t.Fatalf("CRS = %+v", x.CRS)
				}
				if !x.IsParsed(XmpMMNS) {
					t.Fatal("MM namespace not parsed")
				}
				if len(x.MM.History) == 0 || x.MM.HistoryAction != "saved" || x.MM.HistorySoftwareAgent == "" {
					t.Fatalf("MM history = %+v", x.MM)
				}
				if !x.IsParsed(LrNS) {
					t.Fatal("Lightroom namespace not parsed")
				}
				if len(x.Lightroom.HierarchicalSubject) == 0 {
					t.Fatal("Lightroom.HierarchicalSubject empty")
				}
				assertApproxFloat64(t, float64(x.Exif.ApertureValue), 16.0, 0.0001, "Exif.ApertureValue")
				assertApproxFloat64(t, float64(x.Exif.MaxApertureValue), math.Sqrt2, 0.0001, "Exif.MaxApertureValue")
				assertApproxFloat64(t, x.Exif.ShutterSpeedValue.Float64(), -0.378512, 0.000001, "Exif.ShutterSpeedValue")
				assertApproxFloat64(t, float64(x.Exif.GPS.Altitude.Signed()), 6.9, 0.0001, "Exif.GPS.Altitude")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("test", tt.file)
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			x, err := ParseXmp(f)
			if err != nil {
				t.Fatal(err)
			}
			tt.assert(t, x)
		})
	}
}

func TestParseIntegratesJPEG(t *testing.T) {
	packet, err := os.ReadFile(filepath.Join("test", "acr_sidecar.xmp"))
	if err != nil {
		t.Fatal(err)
	}

	jpegBytes := makeJPEGWithXMP(packet)
	x, err := Parse(bytes.NewReader(jpegBytes))
	if err != nil {
		t.Fatal(err)
	}
	if !x.IsParsed(CrsNS) {
		t.Fatal("CRS namespace not parsed")
	}
	if x.CRS.RawFileName != "IMG_9620.CR2" {
		t.Fatalf("CRS.RawFileName = %q", x.CRS.RawFileName)
	}
}

func TestParseIntegratesISOBMFF(t *testing.T) {
	packet, err := os.ReadFile(filepath.Join("test", "acr_sidecar.xmp"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		majorBrand string
	}{
		{name: "cr3", majorBrand: "crx "},
		{name: "heic", majorBrand: "heic"},
		{name: "jxl", majorBrand: "jxl "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := makeBMFFWithXMP(tt.majorBrand, packet)

			x, err := Parse(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			if !x.IsParsed(CrsNS) {
				t.Fatal("CRS namespace not parsed")
			}
			if x.CRS.RawFileName != "IMG_9620.CR2" {
				t.Fatalf("CRS.RawFileName = %q", x.CRS.RawFileName)
			}
		})
	}
}

func TestParseDNGFallback(t *testing.T) {
	packet, err := os.ReadFile(filepath.Join("test", "dng_embedded.xmp"))
	if err != nil {
		t.Fatal(err)
	}

	data := makeDNGWithXMP(packet)
	x, err := Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if x.Tiff.Make != "OLYMPUS CORPORATION" {
		t.Fatalf("Tiff.Make = %q", x.Tiff.Make)
	}
}

func TestParseNoXMP(t *testing.T) {
	_, err := Parse(bytes.NewReader([]byte("not an xmp packet")))
	if !errors.Is(err, ErrNoXMP) {
		t.Fatalf("err = %v, want ErrNoXMP", err)
	}
}

func TestParseWithOptionsMatchesDefault(t *testing.T) {
	packet, err := os.ReadFile(filepath.Join("test", "acr_sidecar.xmp"))
	if err != nil {
		t.Fatal(err)
	}

	gotDefault, err := ParseXmp(bytes.NewReader(packet))
	if err != nil {
		t.Fatal(err)
	}

	gotDebug, err := ParseXmpWithOptions(bytes.NewReader(packet), ParseOptions{Debug: true})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(gotDefault, gotDebug) {
		t.Fatal("ParseXmpWithOptions(Debug=true) result differs from ParseXmp")
	}
}

func TestParseLeavesMissingNamespacesZero(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmp:CreateDate="2025-01-02T03:04:05Z"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if x.Basic.CreateDate.IsZero() {
		t.Fatal("Basic.CreateDate is zero")
	}
	if !reflect.DeepEqual(x.Exif, Exif{}) ||
		!reflect.DeepEqual(x.Aux, Aux{}) ||
		!reflect.DeepEqual(x.Tiff, Tiff{}) ||
		!reflect.DeepEqual(x.DC, DublinCore{}) ||
		!reflect.DeepEqual(x.CRS, CRS{}) ||
		!reflect.DeepEqual(x.MM, XMPMM{}) ||
		!reflect.DeepEqual(x.Photoshop, Photoshop{}) ||
		!reflect.DeepEqual(x.DynamicMedia, DynamicMedia{}) ||
		!reflect.DeepEqual(x.Lightroom, Lightroom{}) ||
		!reflect.DeepEqual(x.Regions, RegionInfo{}) ||
		x.IsParsed(CrsNS) ||
		x.IsParsed(XmpMMNS) ||
		x.IsParsed(PhotoshopNS) ||
		x.IsParsed(XmpDMNS) ||
		x.IsParsed(LrNS) ||
		x.IsParsed(MwgRSNS) {
		t.Fatalf("unexpected non-zero namespaces: %+v", x)
	}
}

func TestParseDecodesXMLEntities(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title><rdf:Alt><rdf:li xml:lang="x-default">A &amp; B &#x1F44D;</rdf:li></rdf:Alt></dc:title></rdf:Description></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(x.DC.Title) != 1 || x.DC.Title[0] != "A & B 👍" {
		t.Fatalf("decoded title = %v", x.DC.Title)
	}
}

func TestParseHandlesXPacketAndDerivedFrom(t *testing.T) {
	const src = `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Adobe XMP Core 7.0-c000">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about=""
 xmlns:xmp="http://ns.adobe.com/xap/1.0/"
 xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
 xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
 xmlns:xmpDM="http://ns.adobe.com/xmp/1.0/DynamicMedia/"
 xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
 xmpDM:pick="1"
 xmpDM:good="true">
 <xmpMM:DerivedFrom stRef:documentID="0CAE8A6363AD4E49964DA145CBBD23D1" stRef:originalDocumentID="0CAE8A6363AD4E49964DA145CBBD23D1"/>
 <lr:weightedFlatSubject><rdf:Bag><rdf:li>Wedding</rdf:li></rdf:Bag></lr:weightedFlatSubject>
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>
<?xpacket end="w"?>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if x.Basic.Toolkit != "Adobe XMP Core 7.0-c000" {
		t.Fatalf("Basic.Toolkit = %q", x.Basic.Toolkit)
	}
	if !x.IsParsed(XmpDMNS) {
		t.Fatal("DynamicMedia namespace not parsed")
	}
	if x.DynamicMedia.Pick != 1 || !x.DynamicMedia.Good {
		t.Fatalf("DynamicMedia = %+v", x.DynamicMedia)
	}
	if !x.IsParsed(XmpMMNS) {
		t.Fatal("MM namespace not parsed")
	}
	if x.MM.DerivedFromDocumentID.String() == "" || x.MM.DerivedFromOriginalDocumentID.String() == "" {
		t.Fatalf("DerivedFrom = %+v", x.MM)
	}
	if !x.IsParsed(LrNS) {
		t.Fatal("Lightroom namespace not parsed")
	}
	if len(x.Lightroom.WeightedFlatSubject) != 1 || x.Lightroom.WeightedFlatSubject[0] != "Wedding" {
		t.Fatalf("Lightroom.WeightedFlatSubject = %v", x.Lightroom.WeightedFlatSubject)
	}
}

func TestParseExifAliasTags(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" exif:ExifImageWidth="4000" exif:ExifImageHeight="3000" exif:ISO="640" exif:ExposureCompensation="+1/3" exif:FlashFired="True" exif:FlashReturn="0" exif:FlashMode="2" exif:FlashFunction="False" exif:FlashRedEyeMode="False"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if x.Exif.PixelXDimension != 4000 || x.Exif.PixelYDimension != 3000 {
		t.Fatalf("Exif dimensions = %dx%d", x.Exif.PixelXDimension, x.Exif.PixelYDimension)
	}
	if x.Exif.ISOSpeedRatings != 640 {
		t.Fatalf("Exif.ISOSpeedRatings = %d", x.Exif.ISOSpeedRatings)
	}
	if x.Exif.ExposureBias.String() != "+1/3" {
		t.Fatalf("Exif.ExposureBias = %s", x.Exif.ExposureBias.String())
	}
	if !x.Exif.Flash.Fired || x.Exif.Flash.Mode != 2 || x.Exif.Flash.Function || x.Exif.Flash.RedEyeMode {
		t.Fatalf("Exif.Flash = %+v", x.Exif.Flash)
	}
}

func makeJPEGWithXMP(packet []byte) []byte {
	const xmpPrefix = "http://ns.adobe.com/xap/1.0/\x00"

	segmentLen := 2 + len(xmpPrefix) + len(packet)
	out := make([]byte, 0, 4+segmentLen+72)
	out = append(out, 0xFF, 0xD8) // SOI

	out = append(out, 0xFF, 0xE1)
	out = append(out, byte(segmentLen>>8), byte(segmentLen))
	out = append(out, xmpPrefix...)
	out = append(out, packet...)

	// The jpeg scanner stops at SOS and peeks 64 bytes per iteration.
	out = append(out, 0xFF, 0xDA, 0x00, 0x04, 0x00, 0x00)
	out = append(out, make([]byte, 64)...)
	return out
}

func makeBMFFWithXMP(majorBrand string, packet []byte) []byte {
	const xpacketUUID = "be7acfcb-97a9-42e8-9c71-999491e3afac"

	uuidPayload := append(meta.UUIDFromString(xpacketUUID).Bytes(), packet...)

	out := make([]byte, 0, 8+8+8+len(uuidPayload))
	out = append(out, makeFTYP(majorBrand)...)
	out = append(out, makeBox("uuid", uuidPayload)...)
	return out
}

func makeFTYP(major string) []byte {
	payload := make([]byte, 8)
	copy(payload[:4], []byte(major))
	copy(payload[4:8], []byte("0001"))
	return makeBox("ftyp", payload)
}

func makeBox(boxType string, payload []byte) []byte {
	out := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(out[:4], uint32(len(out)))
	copy(out[4:8], []byte(boxType))
	copy(out[8:], payload)
	return out
}

func makeDNGWithXMP(packet []byte) []byte {
	// Header layout mirrors imagetype/gen tiffSubtypeHeader with DNG markers.
	header := make([]byte, 64)
	header[0], header[1], header[2], header[3] = 0x49, 0x49, 0x2A, 0x00
	binary.LittleEndian.PutUint16(header[8:10], 0x003F)
	binary.LittleEndian.PutUint16(header[10:12], 0x00FE)
	binary.LittleEndian.PutUint16(header[12:14], 4)
	binary.LittleEndian.PutUint32(header[14:18], 1)
	binary.LittleEndian.PutUint32(header[18:22], 1)
	binary.LittleEndian.PutUint16(header[22:24], 0x0100)

	out := make([]byte, 0, len(header)+len(packet))
	out = append(out, header...)
	out = append(out, packet...)
	return out
}

func assertApproxFloat64(t *testing.T, got, want, epsilon float64, name string) {
	t.Helper()
	if math.Abs(got-want) > epsilon {
		t.Fatalf("%s = %.8f, want %.8f (eps %.8f)", name, got, want, epsilon)
	}
}
