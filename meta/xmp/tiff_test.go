package xmp

import (
	"strings"
	"testing"
)

func TestParseTiffAdditionalTags(t *testing.T) {
	const src = `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:tiff="http://ns.adobe.com/tiff/1.0/" tiff:Artist="Artist A" tiff:BitsPerSample="8 8 8" tiff:DateTime="2025-01-02T03:04:05Z" tiff:Compression="5" tiff:Copyright="(c) 2026" tiff:ImageDescription="desc" tiff:ImageWidth="4000" tiff:ImageHeight="3000" tiff:Make="Canon" tiff:Model="R6" tiff:NativeDigest="abcd" tiff:Orientation="1" tiff:PhotometricInterpretation="2" tiff:PlanarConfiguration="1" tiff:PrimaryChromaticities="0.64 0.33 0.30 0.60 0.15 0.06" tiff:ReferenceBlackWhite="0 255 128 255 128 255" tiff:ResolutionUnit="2" tiff:SamplesPerPixel="3" tiff:Software="Soft A" tiff:TransferFunction="0 1 2" tiff:WhitePoint="0.3127 0.3290" tiff:XResolution="300/1" tiff:YCbCrCoefficients="0.299 0.587 0.114" tiff:YCbCrPositioning="1" tiff:YCbCrSubSampling="2 2" tiff:YResolution="300/1"/></rdf:RDF></x:xmpmeta>`

	x, err := ParseXmp(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	got := x.Tiff
	if got.Artist != "Artist A" || got.BitsPerSample != "8 8 8" {
		t.Fatalf("artist/bits mismatch: %+v", got)
	}
	if got.DateTime.IsZero() || got.DateTime.Format("2006-01-02T15:04:05Z07:00") != "2025-01-02T03:04:05Z" {
		t.Fatalf("DateTime mismatch: %v", got.DateTime)
	}
	if got.Make != "Canon" || got.Model != "R6" || got.Software != "Soft A" {
		t.Fatalf("make/model/software mismatch: %+v", got)
	}
	if len(got.Copyright) != 1 || got.Copyright[0] != "(c) 2026" {
		t.Fatalf("copyright mismatch: %+v", got.Copyright)
	}
	if len(got.ImageDescription) != 1 || got.ImageDescription[0] != "desc" {
		t.Fatalf("image description mismatch: %+v", got.ImageDescription)
	}
	if got.ImageWidth != 4000 || got.ImageLength != 3000 {
		t.Fatalf("dimensions mismatch: %dx%d", got.ImageWidth, got.ImageLength)
	}
	if got.NativeDigest != "abcd" || got.Orientation != 1 || got.Compression != 5 {
		t.Fatalf("native/orientation/compression mismatch: %+v", got)
	}
	if got.PhotometricInterpretation != 2 || got.PlanarConfiguration != 1 {
		t.Fatalf("photo/planar mismatch: %+v", got)
	}
	if got.PrimaryChromaticities == "" || got.ReferenceBlackWhite == "" || got.TransferFunction == "" || got.WhitePoint == "" {
		t.Fatalf("color calibration fields mismatch: %+v", got)
	}
	if got.ResolutionUnit != 2 || got.SamplesPerPixel != 3 {
		t.Fatalf("resolution unit / spp mismatch: %+v", got)
	}
	assertApproxFloat64(t, got.XResolution, 300, 0.0001, "Tiff.XResolution")
	assertApproxFloat64(t, got.YResolution, 300, 0.0001, "Tiff.YResolution")
	if got.YCbCrCoefficients == "" || got.YCbCrPositioning != 1 || got.YCbCrSubSampling != "2 2" {
		t.Fatalf("ycbcr mismatch: %+v", got)
	}
}
