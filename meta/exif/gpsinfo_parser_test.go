package exif

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestParseGPSTagCoverage(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	parseEmbedded := func(id tag.ID, typ tag.Type, unitCount uint32, valueOffset uint32) {
		tg := tag.NewEntry(id, typ, unitCount, valueOffset, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}
	parseWithData := func(id tag.ID, typ tag.Type, unitCount uint32, data []byte) {
		r.setReader(bytes.NewReader(data))
		r.resetOffsets()
		tg := tag.NewEntry(id, typ, unitCount, 0, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}

	// Refs/signs and embedded primitives.
	parseWithData(tag.TagGPSVersionID, tag.TypeByte, 5, []byte{2, 3, 0, 0, 0})
	parseEmbedded(tag.TagGPSLatitudeRef, tag.TypeASCII, 2, packLE([4]byte{'S', 0, 0, 0}))
	parseEmbedded(tag.TagGPSLongitudeRef, tag.TypeASCII, 2, packLE([4]byte{'W', 0, 0, 0}))
	parseEmbedded(tag.TagGPSAltitudeRef, tag.TypeByte, 1, 1)
	parseEmbedded(tag.TagGPSSpeedRef, tag.TypeASCII, 2, packLE([4]byte{'M', 0, 0, 0}))
	parseEmbedded(tag.TagGPSTrackRef, tag.TypeASCII, 2, packLE([4]byte{'T', 0, 0, 0}))
	parseEmbedded(tag.TagGPSImgDirectionRef, tag.TypeASCII, 2, packLE([4]byte{'M', 0, 0, 0}))
	parseEmbedded(tag.TagGPSDestLatitudeRef, tag.TypeASCII, 2, packLE([4]byte{'S', 0, 0, 0}))
	parseEmbedded(tag.TagGPSDestLongitudeRef, tag.TypeASCII, 2, packLE([4]byte{'W', 0, 0, 0}))
	parseEmbedded(tag.TagGPSDestBearingRef, tag.TypeASCII, 2, packLE([4]byte{'M', 0, 0, 0}))
	parseEmbedded(tag.TagGPSDestDistanceRef, tag.TypeASCII, 2, packLE([4]byte{'N', 0, 0, 0}))
	parseEmbedded(tag.TagGPSDifferential, tag.TypeShort, 1, 1)

	// Coordinate and rational values.
	parseWithData(tag.TagGPSLatitude, tag.TypeRational, 3, packRationalList([2]uint32{33, 1}, [2]uint32{54, 1}, [2]uint32{0, 1}))
	parseWithData(tag.TagGPSLongitude, tag.TypeRational, 3, packRationalList([2]uint32{18, 1}, [2]uint32{24, 1}, [2]uint32{0, 1}))
	parseWithData(tag.TagGPSAltitude, tag.TypeRational, 1, packRationalList([2]uint32{1234, 10}))
	parseWithData(tag.TagGPSDOP, tag.TypeRational, 1, packRationalList([2]uint32{7, 2}))
	parseWithData(tag.TagGPSSpeed, tag.TypeRational, 1, packRationalList([2]uint32{60, 1}))
	parseWithData(tag.TagGPSTrack, tag.TypeRational, 1, packRationalList([2]uint32{123, 1}))
	parseWithData(tag.TagGPSImgDirection, tag.TypeRational, 1, packRationalList([2]uint32{271, 1}))
	parseWithData(tag.TagGPSDestLatitude, tag.TypeRational, 3, packRationalList([2]uint32{45, 1}, [2]uint32{30, 1}, [2]uint32{0, 1}))
	parseWithData(tag.TagGPSDestLongitude, tag.TypeRational, 3, packRationalList([2]uint32{9, 1}, [2]uint32{6, 1}, [2]uint32{0, 1}))
	parseWithData(tag.TagGPSDestBearing, tag.TypeRational, 1, packRationalList([2]uint32{180, 1}))
	parseWithData(tag.TagGPSDestDistance, tag.TypeRational, 1, packRationalList([2]uint32{2, 1}))
	parseWithData(tag.TagGPSHPositioningError, tag.TypeRational, 1, packRationalList([2]uint32{3, 2}))
	parseWithData(tag.TagGPSTimeStamp, tag.TypeRational, 3, packRationalList([2]uint32{12, 1}, [2]uint32{34, 1}, [2]uint32{113, 2}))

	// Text values, including EXIF-header encoded Undefined tags.
	parseWithData(tag.TagGPSSatellites, tag.TypeASCII, 5, []byte("7/12\x00"))
	parseWithData(tag.TagGPSStatus, tag.TypeASCII, 5, []byte("A\x00\x00\x00\x00"))
	parseWithData(tag.TagGPSMeasureMode, tag.TypeASCII, 5, []byte("3\x00\x00\x00\x00"))
	parseWithData(tag.TagGPSMapDatum, tag.TypeASCII, 7, []byte("WGS-84\x00"))
	parseWithData(tag.TagGPSProcessingMethod, tag.TypeUndefined, 12, []byte("ASCII\x00\x00\x00GPS\x00"))
	parseWithData(tag.TagGPSAreaInformation, tag.TypeUndefined, 18, []byte("ASCII\x00\x00\x00Cape Town\x00"))
	parseWithData(tag.TagGPSDateStamp, tag.TypeASCII, 11, []byte("2024:03:01\x00"))

	gps := r.Exif.GPS

	if got := gps.GPSVersion.String(); got != "2300" {
		t.Fatalf("GPSVersion = %q, want %q", got, "2300")
	}
	if got := gps.Latitude(); math.Abs(got-(-33.9)) > 0.00001 {
		t.Fatalf("Latitude = %v, want -33.9", got)
	}
	if got := gps.Longitude(); math.Abs(got-(-18.4)) > 0.00001 {
		t.Fatalf("Longitude = %v, want -18.4", got)
	}
	if got := gps.Altitude(); math.Abs(float64(got-(-123.4))) > 0.00001 {
		t.Fatalf("Altitude = %v, want -123.4", got)
	}
	if got := gps.DOP(); math.Abs(got-3.5) > 0.00001 {
		t.Fatalf("DOP = %v, want 3.5", got)
	}
	if got := gps.HPositioningError(); math.Abs(got-1.5) > 0.00001 {
		t.Fatalf("HPositioningError = %v, want 1.5", got)
	}
	if got := gps.SpeedWithRef(); got.Ref != "M" || math.Abs(got.Value.Float64()-96) > 1 {
		t.Fatalf("SpeedWithRef = %+v, want ref M and mph-converted value", got)
	}
	if got := gps.TrackWithRef(); got.Ref != "T" || got.Value.Numerator != 123 {
		t.Fatalf("TrackWithRef = %+v", got)
	}
	if got := gps.ImgDirectionWithRef(); got.Ref != "M" || got.Value.Numerator != 271 {
		t.Fatalf("ImgDirectionWithRef = %+v", got)
	}
	if got := gps.DestLatitudeSigned(); math.Abs(got-(-45.5)) > 0.00001 {
		t.Fatalf("DestLatitudeSigned = %v, want -45.5", got)
	}
	if got := gps.DestLongitudeSigned(); math.Abs(got-(-9.1)) > 0.00001 {
		t.Fatalf("DestLongitudeSigned = %v, want -9.1", got)
	}
	if got := gps.DestBearingWithRef(); got.Ref != "M" || got.Value.Numerator != 180 {
		t.Fatalf("DestBearingWithRef = %+v", got)
	}
	if got := gps.DestDistanceWithRef(); got.Ref != "N" || got.Value.Numerator != 3704 {
		t.Fatalf("DestDistanceWithRef = %+v", got)
	}
	if got := gps.Satellites(); got != "7/12" {
		t.Fatalf("Satellites = %q, want %q", got, "7/12")
	}
	if got := gps.Status(); got != "A" {
		t.Fatalf("Status = %q, want %q", got, "A")
	}
	if got := gps.MeasureMode(); got != "3" {
		t.Fatalf("MeasureMode = %q, want %q", got, "3")
	}
	if got := gps.MapDatum(); got != "WGS-84" {
		t.Fatalf("MapDatum = %q, want %q", got, "WGS-84")
	}
	if got := gps.ProcessingMethod(); got != "GPS" {
		t.Fatalf("ProcessingMethod = %q, want %q", got, "GPS")
	}
	if got := gps.AreaInformation(); got != "Cape Town" {
		t.Fatalf("AreaInformation = %q, want %q", got, "Cape Town")
	}
	if got := gps.Differential(); got != 1 {
		t.Fatalf("Differential = %d, want 1", got)
	}
	if got := gps.GPSTimestamp(); !got.Equal(time.Date(2024, time.March, 1, 12, 34, 56, 500000000, time.UTC)) {
		t.Fatalf("GPSTimestamp = %v, want 2024-03-01 12:34:56.5 +0000 UTC", got)
	}
}

func TestParseGPSTagUnknownTagReturnsFalse(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	tg := tag.NewEntry(0xfffe, tag.TypeLong, 1, 0, tag.GPSIFD, 0, utils.LittleEndian)
	if ok := r.parseGPSTag(tg); ok {
		t.Fatal("parseGPSTag(unknown) = true, want false")
	}
}

func TestParseGPSTagUnitConversions(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	parseEmbedded := func(id tag.ID, typ tag.Type, unitCount uint32, valueOffset uint32) {
		tg := tag.NewEntry(id, typ, unitCount, valueOffset, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}
	parseWithData := func(id tag.ID, typ tag.Type, unitCount uint32, data []byte) {
		r.setReader(bytes.NewReader(data))
		r.resetOffsets()
		tg := tag.NewEntry(id, typ, unitCount, 0, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}

	parseEmbedded(tag.TagGPSSpeedRef, tag.TypeASCII, 2, packLE([4]byte{'N', 0, 0, 0}))
	parseWithData(tag.TagGPSSpeed, tag.TypeRational, 1, packRationalList([2]uint32{10, 1}))
	if got := r.Exif.GPS.speed; math.Abs(got-18.52) > 0.00001 {
		t.Fatalf("speed(knots->km/h) = %v, want 18.52", got)
	}

	parseEmbedded(tag.TagGPSSpeedRef, tag.TypeASCII, 2, packLE([4]byte{'M', 0, 0, 0}))
	parseWithData(tag.TagGPSSpeed, tag.TypeRational, 1, packRationalList([2]uint32{10, 1}))
	if got := r.Exif.GPS.speed; math.Abs(got-16.0934) > 0.00001 {
		t.Fatalf("speed(mph->km/h) = %v, want 16.0934", got)
	}

	parseEmbedded(tag.TagGPSDestDistanceRef, tag.TypeASCII, 2, packLE([4]byte{'M', 0, 0, 0}))
	parseWithData(tag.TagGPSDestDistance, tag.TypeRational, 1, packRationalList([2]uint32{1, 1}))
	if got := r.Exif.GPS.destDistance; math.Abs(got-1609.34) > 0.00001 {
		t.Fatalf("destDistance(mi->m) = %v, want 1609.34", got)
	}

	parseEmbedded(tag.TagGPSDestDistanceRef, tag.TypeASCII, 2, packLE([4]byte{'K', 0, 0, 0}))
	parseWithData(tag.TagGPSDestDistance, tag.TypeRational, 1, packRationalList([2]uint32{1, 1}))
	if got := r.Exif.GPS.destDistance; math.Abs(got-1000) > 0.00001 {
		t.Fatalf("destDistance(km->m) = %v, want 1000", got)
	}

	parseEmbedded(tag.TagGPSDestDistanceRef, tag.TypeASCII, 2, packLE([4]byte{'N', 0, 0, 0}))
	parseWithData(tag.TagGPSDestDistance, tag.TypeRational, 1, packRationalList([2]uint32{1, 1}))
	if got := r.Exif.GPS.destDistance; math.Abs(got-1852) > 0.00001 {
		t.Fatalf("destDistance(nm->m) = %v, want 1852", got)
	}
}

func TestParseGPSTagOrderInsensitiveCoordinates(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	parseEmbedded := func(id tag.ID, typ tag.Type, unitCount uint32, valueOffset uint32) {
		tg := tag.NewEntry(id, typ, unitCount, valueOffset, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}
	parseWithData := func(id tag.ID, typ tag.Type, unitCount uint32, data []byte) {
		r.setReader(bytes.NewReader(data))
		r.resetOffsets()
		tg := tag.NewEntry(id, typ, unitCount, 0, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}

	parseWithData(tag.TagGPSLatitude, tag.TypeRational, 3, packRationalList([2]uint32{33, 1}, [2]uint32{54, 1}, [2]uint32{0, 1}))
	parseEmbedded(tag.TagGPSLatitudeRef, tag.TypeASCII, 2, packLE([4]byte{'S', 0, 0, 0}))
	if got := r.Exif.GPS.Latitude(); math.Abs(got-(-33.9)) > 0.00001 {
		t.Fatalf("Latitude after ref = %v, want -33.9", got)
	}

	parseWithData(tag.TagGPSLongitude, tag.TypeRational, 3, packRationalList([2]uint32{18, 1}, [2]uint32{24, 1}, [2]uint32{0, 1}))
	parseEmbedded(tag.TagGPSLongitudeRef, tag.TypeASCII, 2, packLE([4]byte{'W', 0, 0, 0}))
	if got := r.Exif.GPS.Longitude(); math.Abs(got-(-18.4)) > 0.00001 {
		t.Fatalf("Longitude after ref = %v, want -18.4", got)
	}

	parseWithData(tag.TagGPSDestLatitude, tag.TypeRational, 3, packRationalList([2]uint32{45, 1}, [2]uint32{30, 1}, [2]uint32{0, 1}))
	parseEmbedded(tag.TagGPSDestLatitudeRef, tag.TypeASCII, 2, packLE([4]byte{'S', 0, 0, 0}))
	if got := r.Exif.GPS.DestLatitude(); math.Abs(got-(-45.5)) > 0.00001 {
		t.Fatalf("DestLatitude after ref = %v, want -45.5", got)
	}

	parseWithData(tag.TagGPSDestLongitude, tag.TypeRational, 3, packRationalList([2]uint32{9, 1}, [2]uint32{6, 1}, [2]uint32{0, 1}))
	parseEmbedded(tag.TagGPSDestLongitudeRef, tag.TypeASCII, 2, packLE([4]byte{'W', 0, 0, 0}))
	if got := r.Exif.GPS.DestLongitude(); math.Abs(got-(-9.1)) > 0.00001 {
		t.Fatalf("DestLongitude after ref = %v, want -9.1", got)
	}

	parseWithData(tag.TagGPSAltitude, tag.TypeRational, 1, packRationalList([2]uint32{1234, 10}))
	parseEmbedded(tag.TagGPSAltitudeRef, tag.TypeByte, 1, 1)
	if got := r.Exif.GPS.Altitude(); math.Abs(float64(got-(-123.4))) > 0.00001 {
		t.Fatalf("Altitude after ref = %v, want -123.4", got)
	}
}

func TestParseGPSAltitudeRefExifToolVariants(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	parseEmbedded := func(id tag.ID, typ tag.Type, unitCount uint32, valueOffset uint32) {
		tg := tag.NewEntry(id, typ, unitCount, valueOffset, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}
	parseWithData := func(id tag.ID, typ tag.Type, unitCount uint32, data []byte) {
		r.setReader(bytes.NewReader(data))
		r.resetOffsets()
		tg := tag.NewEntry(id, typ, unitCount, 0, tag.GPSIFD, 0, utils.LittleEndian)
		if ok := r.parseGPSTag(tg); !ok {
			t.Fatalf("parseGPSTag(%v) = false, want true", id)
		}
	}

	parseWithData(tag.TagGPSAltitude, tag.TypeRational, 1, packRationalList([2]uint32{25, 1}))
	parseEmbedded(tag.TagGPSAltitudeRef, tag.TypeByte, 1, 2)
	if got := r.Exif.GPS.Altitude(); got != 25 {
		t.Fatalf("Altitude ref=2 = %v, want 25", got)
	}

	parseEmbedded(tag.TagGPSAltitudeRef, tag.TypeByte, 1, 3)
	if got := r.Exif.GPS.Altitude(); got != -25 {
		t.Fatalf("Altitude ref=3 = %v, want -25", got)
	}
}

func TestParseGPSDateStampNULSeparators(t *testing.T) {
	t.Parallel()

	r := NewReader(metalog.GetLogger())
	defer r.Close()

	data := []byte{'2', '0', '2', '4', 0, '0', '3', 0, '0', '1', 0}
	r.setReader(bytes.NewReader(data))
	r.resetOffsets()
	tg := tag.NewEntry(tag.TagGPSDateStamp, tag.TypeASCII, uint32(len(data)), 0, tag.GPSIFD, 0, utils.LittleEndian)
	if ok := r.parseGPSTag(tg); !ok {
		t.Fatal("parseGPSTag(GPSDateStamp) = false, want true")
	}

	want := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	if got := r.Exif.GPS.GPSTimestamp(); !got.Equal(want) {
		t.Fatalf("GPSTimestamp = %v, want %v", got, want)
	}
}

func packLE(b [4]byte) uint32 {
	return utils.LittleEndian.Uint32(b[:])
}

func packRationalList(vals ...[2]uint32) []byte {
	buf := make([]byte, len(vals)*8)
	for i := range vals {
		base := i * 8
		utils.LittleEndian.PutUint32(buf[base:base+4], vals[i][0])
		utils.LittleEndian.PutUint32(buf[base+4:base+8], vals[i][1])
	}
	return buf
}
