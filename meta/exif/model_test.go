package exif

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func TestRationalUFloat64(t *testing.T) {
	t.Parallel()

	if got, want := (tag.RationalU{Numerator: 3, Denominator: 2}).Float64(), 1.5; got != want {
		t.Fatalf("tag.RationalU.Float64() = %v, want %v", got, want)
	}
	if got := (tag.RationalU{Numerator: 3, Denominator: 0}).Float64(); got != 0 {
		t.Fatalf("tag.RationalU.Float64() with zero denominator = %v, want 0", got)
	}
}

func TestLensInfoStringAndJSON(t *testing.T) {
	t.Parallel()

	info := &LensInfo{
		MinFocalLength:        tag.RationalU{Numerator: 24, Denominator: 1},
		MaxFocalLength:        tag.RationalU{Numerator: 70, Denominator: 1},
		MaxApertureAtMinFocal: tag.RationalU{Numerator: 28, Denominator: 10},
		MaxApertureAtMaxFocal: tag.RationalU{Numerator: 4, Denominator: 1},
	}

	if got, want := info.String(), "24 70 2.8 4"; got != want {
		t.Fatalf("LensInfo.String() = %q, want %q", got, want)
	}

	buf, err := json.Marshal(struct {
		LensInfo *LensInfo
	}{LensInfo: info})
	if err != nil {
		t.Fatalf("json.Marshal(LensInfo): %v", err)
	}
	if got, want := string(buf), `{"LensInfo":"24 70 2.8 4"}`; got != want {
		t.Fatalf("json.Marshal(LensInfo) = %s, want %s", got, want)
	}

	buf, err = json.Marshal(struct {
		LensInfo *LensInfo
	}{})
	if err != nil {
		t.Fatalf("json.Marshal(nil LensInfo): %v", err)
	}
	if got, want := string(buf), `{"LensInfo":null}`; got != want {
		t.Fatalf("json.Marshal(nil LensInfo) = %s, want %s", got, want)
	}
}

func TestGPSInfoSetDateAndTimeOrder(t *testing.T) {
	t.Parallel()

	base := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	delta := 12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond
	want := base.Add(delta)

	var a GPSInfo
	a.setDate(base)
	a.setTime(delta)
	if got := a.GPSTime(); !got.Equal(want) {
		t.Fatalf("date->time GPSTime = %v, want %v", got, want)
	}

	var b GPSInfo
	b.setTime(delta)
	if pending, ok := gpsPendingDelta(b.GPSTime()); !ok || pending != delta {
		t.Fatalf("gpsPendingDelta() = (%v,%v), want (%v,true)", pending, ok, delta)
	}
	b.setDate(base)
	if got := b.GPSTime(); !got.Equal(want) {
		t.Fatalf("time->date GPSTime = %v, want %v", got, want)
	}

	var c GPSInfo
	c.setTime(0)
	if got := c.GPSTime(); !got.IsZero() {
		t.Fatalf("zero-delta GPSTime = %v, want zero", got)
	}
}

func TestGPSInfoAccessorsAndBitset(t *testing.T) {
	t.Parallel()

	g := GPSInfo{
		latitude:         -33.9,
		longitude:        18.4,
		altitude:         -123.4,
		latitudeRef:      tag.GPSRefSouth,
		longitudeRef:     tag.GPSRefEast,
		altitudeRef:      tag.GPSRefBelowSeaLevel,
		speedRef:         tag.GPSRefKilometersPerHour,
		speed:            tag.RationalU{Numerator: 90, Denominator: 1}.Float64(),
		trackRef:         tag.GPSRefTrueDirection,
		track:            tag.RationalU{Numerator: 123, Denominator: 1}.Float64(),
		destLatitude:     45.5,
		destLongitude:    9.1,
		destLatitudeRef:  tag.GPSRefSouth,
		destLongitudeRef: tag.GPSRefEast,
		destDistanceRef:  tag.GPSRefKilometers,
		destDistance:     tag.RationalU{Numerator: 12, Denominator: 1}.Float64(),
		mapDatum:         "WGS-84",
		processingMethod: "GPS",
		areaInformation:  "Cape Town",
		differential:     1,
	}

	if got := g.Latitude(); got != -33.9 {
		t.Fatalf("Latitude() = %v, want -33.9", got)
	}
	if got := g.Longitude(); got != 18.4 {
		t.Fatalf("Longitude() = %v, want 18.4", got)
	}
	if got := g.Altitude(); got != -123.4 {
		t.Fatalf("Altitude() = %v, want -123.4", got)
	}
	if got := g.SpeedWithRef(); got.Ref != "K" || got.Value.Numerator != 90 || got.Value.Denominator != 1 {
		t.Fatalf("SpeedWithRef() = %+v", got)
	}
	if got := g.TrackWithRef(); got.Ref != "T" || got.Value.Numerator != 123 || got.Value.Denominator != 1 {
		t.Fatalf("TrackWithRef() = %+v", got)
	}
	if got := g.DestLatitude(); got != 45.5 {
		t.Fatalf("DestLatitude() = %v, want 45.5", got)
	}
	g.destLatitudeRef = tag.GPSRefNorth
	if got := g.DestLatitudeSigned(); got != 45.5 {
		t.Fatalf("DestLatitudeSigned() with north ref = %v, want 45.5", got)
	}
	if got := g.DestLongitude(); got != 9.1 {
		t.Fatalf("DestLongitude() = %v, want 9.1", got)
	}
	g.destLongitudeRef = tag.GPSRefEast
	if got := g.DestLongitudeSigned(); got != 9.1 {
		t.Fatalf("DestLongitudeSigned() with east ref = %v, want 9.1", got)
	}
	if got := g.DestDistanceWithRef(); got.Ref != "K" || got.Value.Numerator != 12 || got.Value.Denominator != 1 {
		t.Fatalf("DestDistanceWithRef() = %+v", got)
	}
	if got := g.MapDatum(); got != "WGS-84" {
		t.Fatalf("MapDatum() = %q, want %q", got, "WGS-84")
	}
	if got := g.ProcessingMethod(); got != "GPS" {
		t.Fatalf("ProcessingMethod() = %q, want %q", got, "GPS")
	}
	if got := g.AreaInformation(); got != "Cape Town" {
		t.Fatalf("AreaInformation() = %q, want %q", got, "Cape Town")
	}
	if got := g.Differential(); got != 1 {
		t.Fatalf("Differential() = %d, want 1", got)
	}
}

func TestGPSInfoVersionIDFormatting(t *testing.T) {
	t.Parallel()

	if got := (GPSInfo{}).GPSVersion.String(); got != "" {
		t.Fatalf("GPSVersion() zero = %q, want empty", got)
	}

	g := GPSInfo{GPSVersion: GPSVersion{2, 3, 0, 0}}
	if got := g.GPSVersion.String(); got != "2300" {
		t.Fatalf("GPSVersion() = %q, want %q", got, "2300")
	}
	if got := (GPSInfo{GPSVersion: GPSVersion{2, 0, 0, 0}}).GPSVersion.String(); got != "2000" {
		t.Fatalf("GPSVersion(2000) = %q, want %q", got, "2000")
	}
	if got := (GPSInfo{GPSVersion: GPSVersion{2, 1, 0, 0}}).GPSVersion.String(); got != "2100" {
		t.Fatalf("GPSVersion(2100) = %q, want %q", got, "2100")
	}
	if got := (GPSInfo{GPSVersion: GPSVersion{2, 2, 0, 0}}).GPSVersion.String(); got != "2200" {
		t.Fatalf("GPSVersion(2200) = %q, want %q", got, "2200")
	}
	if got := (GPSInfo{GPSVersion: GPSVersion{9, 9, 9, 9}}).GPSVersion.String(); got == "" {
		t.Fatal("GPSVersion(custom) should not be empty")
	}
}

func TestGPSPendingDeltaNonSentinel(t *testing.T) {
	t.Parallel()

	if d, ok := gpsPendingDelta(time.Time{}); ok || d != 0 {
		t.Fatalf("gpsPendingDelta(zero) = (%v,%v), want (0,false)", d, ok)
	}
	nonSentinel := time.Date(2024, time.March, 1, 12, 0, 0, 0, time.UTC)
	if d, ok := gpsPendingDelta(nonSentinel); ok || d != 0 {
		t.Fatalf("gpsPendingDelta(non-sentinel) = (%v,%v), want (0,false)", d, ok)
	}
	nonSentinelDay := time.Date(1, time.January, 2, 12, 0, 0, 0, time.UTC)
	if d, ok := gpsPendingDelta(nonSentinelDay); ok || d != 0 {
		t.Fatalf("gpsPendingDelta(non-sentinel day) = (%v,%v), want (0,false)", d, ok)
	}
}

func TestGPSInfoMarshalJSON(t *testing.T) {
	t.Parallel()

	g := GPSInfo{
		date:              time.Date(2024, time.March, 1, 12, 34, 56, 0, time.UTC),
		latitude:          1.23,
		longitude:         4.56,
		altitude:          7.8,
		destLatitude:      9.1,
		destLongitude:     2.3,
		dop:               1.5,
		speed:             42,
		track:             180,
		imgDirection:      270,
		destBearing:       90,
		destDistance:      1000,
		hPositioningError: 0.9,
		satellites:        "7/12",
		status:            "A",
		measureMode:       "3",
		mapDatum:          "WGS-84",
		processingMethod:  "GPS",
		areaInformation:   "Cape Town",
		differential:      1,
		GPSVersion:        GPSVersion2300,
		speedRef:          tag.GPSRefKilometersPerHour,
		trackRef:          tag.GPSRefTrueDirection,
		imgDirectionRef:   tag.GPSRefMagneticDirection,
		destBearingRef:    tag.GPSRefTrueDirection,
		destDistanceRef:   tag.GPSRefKilometers,
	}
	g.destLatitudeRef = tag.GPSRefSouth
	g.destLongitudeRef = tag.GPSRefWest

	buf, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("json.Marshal(GPSInfo): %v", err)
	}
	got := string(buf)
	if !strings.Contains(got, `"ProcessingMethod":"GPS"`) {
		t.Fatalf("ProcessingMethod missing in JSON: %s", got)
	}
	if !strings.Contains(got, `"AreaInformation":"Cape Town"`) {
		t.Fatalf("AreaInformation missing in JSON: %s", got)
	}
	if !strings.Contains(got, `"GPSVersion":"2300"`) {
		t.Fatalf("GPSVersion missing in JSON: %s", got)
	}
}

func TestIFD0AndExifIFDTimeTagsBitsetAndSelection(t *testing.T) {
	t.Parallel()

	base := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	orig := time.Date(2024, time.January, 3, 4, 5, 6, 0, time.UTC)
	create := time.Date(2024, time.January, 4, 5, 6, 7, 0, time.UTC)

	var ifd0 IFD0TimeTags
	ifd0.setDate(base)
	ifd0.setSubSec("125", 125)
	ifd0.markTagParsed(tag.TagDateTime)
	ifd0.markTagParsed(tag.TagSubSecTime)

	var exifTime ExifIFDTimeTags
	exifTime.setDate(tag.TagDateTimeOriginal, orig)
	exifTime.setDate(tag.TagDateTimeDigitized, create)
	exifTime.setSubSec(tag.TagSubSecTimeOriginal, "50", 50)
	exifTime.setSubSec(tag.TagSubSecTimeDigitized, "75", 75)
	exifTime.markTagParsed(tag.TagSubSecTimeOriginal)
	exifTime.markTagParsed(tag.TagSubSecTimeDigitized)

	if !ifd0.HasTagParsed(tag.TagDateTime) {
		t.Fatal("HasTagParsed(DateTime) should be true")
	}
	if !ifd0.HasSubSecTime() {
		t.Fatal("HasSubSecTime should be true")
	}
	if !exifTime.HasSubSecTimeOriginal() || !exifTime.HasSubSecTimeDigitized() {
		t.Fatal("ExifIFD HasSubSec* methods should both be true")
	}
	if ifd0.HasTagParsed(tag.TagOffsetTime) {
		t.Fatal("HasTagParsed(OffsetTime) should be false")
	}

	if got, want := ifd0.GetModifyDate(), base.Add(125*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetModifyDate() = %v, want %v", got, want)
	}
	if got, want := exifTime.GetDateTimeOriginal(), orig.Add(500*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetDateTimeOriginal() = %v, want %v", got, want)
	}
	if got, want := exifTime.GetCreateDate(), create.Add(750*time.Millisecond); !got.Equal(want) {
		t.Fatalf("GetCreateDate() = %v, want %v", got, want)
	}

	ex := Exif{
		IFD0:    IFD0Tag{IFD0TimeTags: ifd0},
		ExifIFD: ExifIFDTags{ExifIFDTimeTags: exifTime},
	}
	if got, want := ex.SelectedDate(), orig.Add(500*time.Millisecond); !got.Equal(want) {
		t.Fatalf("SelectedDate() = %v, want %v", got, want)
	}
	ex.ExifIFD.DateTimeOriginal = time.Time{}
	if got, want := ex.SelectedDate(), create.Add(750*time.Millisecond); !got.Equal(want) {
		t.Fatalf("SelectedDate() fallback create = %v, want %v", got, want)
	}
	ex.ExifIFD.CreateDate = time.Time{}
	if got, want := ex.SelectedDate(), base.Add(125*time.Millisecond); !got.Equal(want) {
		t.Fatalf("SelectedDate() fallback modify = %v, want %v", got, want)
	}
}

func TestApplyTimeParts(t *testing.T) {
	t.Parallel()

	if got := applyTimeParts(time.Time{}, "123", nil); !got.IsZero() {
		t.Fatalf("applyTimeParts(zero) = %v, want zero", got)
	}

	base := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	tz := time.FixedZone("+02:00", 2*60*60)
	got := applyTimeParts(base, "250", tz)
	want := time.Date(2024, time.January, 2, 3, 4, 5, 250000000, tz)
	if !got.Equal(want) {
		t.Fatalf("applyTimeParts() = %v, want %v", got, want)
	}
}

func TestExifIFDTimeTagsMarshalJSONOffsetStrings(t *testing.T) {
	t.Parallel()

	tt := ExifIFDTimeTags{
		DateTimeOriginal:   time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC),
		OffsetTimeOriginal: time.FixedZone("+01:00", 60*60),
	}

	buf, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("json.Marshal(ExifIFDTimeTags): %v", err)
	}
	got := string(buf)
	if !strings.Contains(got, `"OffsetTimeOriginal":"+01:00"`) {
		t.Fatalf("OffsetTimeOriginal JSON = %s", got)
	}
	if !strings.Contains(got, `"OffsetTimeDigitized":null`) {
		t.Fatalf("OffsetTimeDigitized nil should marshal as null: %s", got)
	}
}

func TestExifIFDExifToolDisplayHelpers(t *testing.T) {
	t.Parallel()

	exif := ExifIFDTags{
		DigitalZoomRatio:           float32(tag.RationalU{Numerator: 3, Denominator: 2}.Float64()),
		FocalPlaneXResolution:      5152000.0 / 243.0,
		FocalPlaneYResolution:      3864000.0 / 183.0,
		FocalPlaneResolutionUnit:   meta.ResolutionUnitInches,
		SubjectDistance:            0.53,
		ExposureIndex:              160,
		focalPlaneXResolutionState: rationalStateValue,
		focalPlaneYResolutionState: rationalStateValue,
		subjectDistanceState:       rationalStateValue,
		exposureIndexState:         rationalStateValue,
	}

	if got, want := exif.ExifToolDigitalZoomRatio(), "1.5"; got != want {
		t.Fatalf("ExifToolDigitalZoomRatio() = %q, want %q", got, want)
	}
	if got, want := exif.ExifToolFocalPlaneXResolution(), "21201.64609"; got != want {
		t.Fatalf("ExifToolFocalPlaneXResolution() = %q, want %q", got, want)
	}
	if got, want := exif.ExifToolFocalPlaneYResolution(), "21114.7541"; got != want {
		t.Fatalf("ExifToolFocalPlaneYResolution() = %q, want %q", got, want)
	}
	if got, want := exif.ExifToolExposureIndex(), "160"; got != want {
		t.Fatalf("ExifToolExposureIndex() = %q, want %q", got, want)
	}
	if got, want := exif.ExifToolFocalPlaneResolutionUnit(), "inches"; got != want {
		t.Fatalf("ExifToolFocalPlaneResolutionUnit() = %q, want %q", got, want)
	}
	if got, want := exif.ExifToolSubjectDistance(), "0.53 m"; got != want {
		t.Fatalf("ExifToolSubjectDistance() = %q, want %q", got, want)
	}

	exif.SubjectDistance = 0
	exif.subjectDistanceState = rationalStateInfinite
	if got, want := exif.ExifToolSubjectDistance(), "undef"; got != want {
		t.Fatalf("ExifToolSubjectDistance() undef = %q, want %q", got, want)
	}

	exif.SubjectDistance = 0
	exif.subjectDistanceState = rationalStateMissing
	if got := exif.ExifToolSubjectDistance(); got != "" {
		t.Fatalf("ExifToolSubjectDistance() zero = %q, want empty", got)
	}
}

func TestExifMarshalJSONIncludesSiblingFieldsAlongsideFlattenedTime(t *testing.T) {
	t.Parallel()

	ex := Exif{}
	ex.IFD0.Make = "Canon"
	ex.IFD0.Model = "EOS R6"
	ex.IFD0.ModifyDate = time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	ex.ExifIFD.LensModel = "RF24-70mm F2.8 L IS USM"
	ex.ExifIFD.DateTimeOriginal = time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	ex.IFD1 = &ImageIFD{ImageWidth: 160}

	buf, err := json.Marshal(ex)
	if err != nil {
		t.Fatalf("json.Marshal(Exif): %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Fatalf("json.Unmarshal(Exif JSON): %v", err)
	}

	ifd0, _ := got["IFD0"].(map[string]any)
	if ifd0["Make"] != "Canon" || ifd0["Model"] != "EOS R6" || ifd0["ModifyDate"] != "2024-01-02T03:04:05Z" {
		t.Fatalf("IFD0 JSON = %#v", ifd0)
	}

	exifIFD, _ := got["ExifIFD"].(map[string]any)
	if exifIFD["LensModel"] != "RF24-70mm F2.8 L IS USM" || exifIFD["DateTimeOriginal"] != "2024-01-02T03:04:05Z" {
		t.Fatalf("ExifIFD JSON = %#v", exifIFD)
	}

	ifd1, _ := got["IFD1"].(map[string]any)
	if ifd1["ImageWidth"] != float64(160) || ifd1["ModifyDate"] != "0001-01-01T00:00:00Z" {
		t.Fatalf("IFD1 JSON = %#v", ifd1)
	}
}
